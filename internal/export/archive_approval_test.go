package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/journal"
)

func approvalArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, h := handoffArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsFS(t), 33); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 33
	acceptMergeArchiveEvent(t, v, h)
	at := time.Date(2026, 11, 10, 9, 7, 0, 0, time.UTC)
	binding := approvalrecord.Binding{Algorithm: fingerprint.Algorithm, Fingerprint: strings.Repeat("f", 64), HandoffID: "HOF-100001", Base: strings.Repeat("b", 40), Head: strings.Repeat("a", 40)}
	raw, _ := json.Marshal(binding)
	return v, journal.Event{EventSchema: 2, EventID: "EVT-20261110-000005", EventType: "task.approved", CommandName: "task.approve", OccurredAt: at, ActorV2: journal.EventActor{ActorID: "mark", ActorType: "human", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{ID: "PDN-approve", Type: "decision", Sequence: 1}, Effects: []journal.Effect{{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-approve"}, Values: map[string]any{"kind": "approve", "actor": "human:mark", "subject_type": "task", "subject_id": "FID-a", "binding_id": binding.Fingerprint, "binding_json": string(raw), "reason": "Reviewed the full change", "decided_at": at.Format(time.RFC3339)}}}}
}

func TestApprovalAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name   string
			change func(*testing.T, *semanticValidator, *journal.Event)
			want   bool
		}{
			{name: "complete", want: true},
			{name: "agent", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.ActorV2.ActorType = "agent"
				e.Effects[0].Values["actor"] = "agent:mark"
			}},
			{name: "service", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.ActorV2.ActorType = "service"
				e.Effects[0].Values["actor"] = "service:mark"
			}},
			{name: "empty reason", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) { e.Effects[0].Values["reason"] = " " }},
			{name: "wrong state", change: func(t *testing.T, v *semanticValidator, _ *journal.Event) {
				if _, err := v.conn.Exec(`UPDATE issues SET status='in_progress' WHERE id='FID-a'`); err != nil {
					t.Fatal(err)
				}
			}},
			{name: "missing handoff", change: func(t *testing.T, v *semanticValidator, _ *journal.Event) {
				if _, err := v.conn.Exec(`DELETE FROM handoffs`); err != nil {
					t.Fatal(err)
				}
			}},
			{name: "wrong task", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Values["subject_id"] = "FID-other"
			}},
			{name: "wrong actor", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Values["actor"] = "human:someone"
			}},
			{name: "wrong kind", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) { e.Effects[0].Values["kind"] = "answer" }},
			{name: "wrong subject type", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Values["subject_type"] = "epic"
			}},
			{name: "wrong binding id", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Values["binding_id"] = "forged"
			}},
			{name: "legacy", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) { e.EventSchema = 1 }},
			{name: "other family", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) { e.EventType = "issue.updated" }},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := approvalArchiveFixture(t)
				if tc.change != nil {
					tc.change(t, v, &e)
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyLeanEmptyEvent(v, mode, e)
				if (err == nil) != tc.want {
					t.Fatalf("admit=%v want=%t", err, tc.want)
				}
				if !tc.want && handoffProjectionSnapshot(t, v) != before {
					t.Fatal("refusal changed state")
				}
			})
		}
	}
}

func TestApprovalForgedBindingsEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, field := range []string{"algorithm", "fingerprint", "handoff_id", "base", "head", "unknown"} {
			t.Run(mode+"/"+field, func(t *testing.T) {
				t.Parallel()
				v, e := approvalArchiveFixture(t)
				var binding map[string]any
				if err := json.Unmarshal([]byte(e.Effects[0].Values["binding_json"].(string)), &binding); err != nil {
					t.Fatal(err)
				}
				binding[field] = "forged"
				raw, _ := json.Marshal(binding)
				e.Effects[0].Values["binding_json"] = string(raw)
				before := handoffProjectionSnapshot(t, v)
				if err := applyLeanEmptyEvent(v, mode, e); err == nil || handoffProjectionSnapshot(t, v) != before {
					t.Fatalf("forged %s=%v", field, err)
				}
			})
		}
	}
}

func TestApprovalIsImmutableThroughLegacyAliases(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, op := range []string{"insert", "update", "delete"} {
			t.Run(mode+"/"+op, func(t *testing.T) {
				t.Parallel()
				v, e := approvalArchiveFixture(t)
				acceptMergeArchiveEvent(t, v, e)
				e.EventID = "EVT-20261110-000006"
				e.EventSchema = 1
				e.EventType = "issue.updated"
				e.Aggregate = journal.Aggregate{Type: "issue", ID: "FID-a", Sequence: 1}
				e.Effects[0].Table = "PERSON_DECISIONS"
				effect := &e.Effects[0]
				effect.Op = op
				effect.PK = map[string]any{"ID": "PDN-approve"}
				if op == "insert" {
					effect.PK["ID"] = "PDN-forged"
					delete(effect.Values, "kind")
					effect.Values["KiNd"] = "approve"
				} else {
					effect.Values = nil
					if op == "update" {
						effect.Values = map[string]any{"reason": "Changed later"}
					}
				}
				before := handoffProjectionSnapshot(t, v)
				if err := applyLeanEmptyEvent(v, mode, e); err == nil || handoffProjectionSnapshot(t, v) != before {
					t.Fatalf("mutable approval %s: %v", op, err)
				}
			})
		}
	}
}
