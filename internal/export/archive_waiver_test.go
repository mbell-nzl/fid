package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

func waiverArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, h := handoffArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsDir(t), 38); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 38
	acceptMergeArchiveEvent(t, v, h)
	at := time.Date(2026, 11, 10, 9, 7, 0, 0, time.UTC)
	binding := waiverrecord.Binding{Lane: "routine", Algorithm: fingerprint.Algorithm, Fingerprint: strings.Repeat("f", 64), HandoffID: "HOF-100001", Base: strings.Repeat("b", 40), Head: strings.Repeat("a", 40)}
	raw, _ := json.Marshal(binding)
	return v, journal.Event{EventSchema: 2, EventID: "EVT-20261110-000005", EventType: "task.waived", CommandName: "task.waive", OccurredAt: at, ActorV2: journal.EventActor{ActorID: "mark", ActorType: "human", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{ID: "PDN-waive", Type: "decision", Sequence: 1}, Effects: []journal.Effect{{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-waive"}, Values: map[string]any{"kind": "waive", "actor": "human:mark", "subject_type": "task", "subject_id": "FID-a", "binding_id": binding.Fingerprint, "binding_json": string(raw), "reason": "Reviewed the full change", "decided_at": at.Format(time.RFC3339)}}}}
}

func TestWaiverAdmissionEveryPath(t *testing.T) {
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
			{name: "extra decision", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				extra := e.Effects[0]
				extra.PK = map[string]any{"id": "PDN-extra"}
				e.Effects = append(e.Effects, extra)
			}},
			{name: "wrong time", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Values["decided_at"] = "2026-11-10T09:08:00Z"
			}},
			{name: "legacy", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) { e.EventSchema = 1 }},
			{name: "other family", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) { e.EventType = "issue.updated" }},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := waiverArchiveFixture(t)
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

func TestWaiverForgedBindingsEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, field := range []string{"algorithm", "fingerprint", "handoff_id", "base", "head", "lane", "unknown"} {
			t.Run(mode+"/"+field, func(t *testing.T) {
				t.Parallel()
				v, e := waiverArchiveFixture(t)
				var binding map[string]any
				if err := json.Unmarshal([]byte(e.Effects[0].Values["binding_json"].(string)), &binding); err != nil {
					t.Fatal(err)
				}
				binding[field] = "forged"
				if field == "base" || field == "head" {
					binding[field] = strings.Repeat("e", 40)
				}
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

func TestWaiverIsImmutableThroughLegacyAliases(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, op := range []string{"insert", "update", "delete"} {
			t.Run(mode+"/"+op, func(t *testing.T) {
				t.Parallel()
				v, e := waiverArchiveFixture(t)
				acceptMergeArchiveEvent(t, v, e)
				e.EventID = "EVT-20261110-000006"
				e.EventSchema = 1
				e.EventType = "issue.updated"
				e.Aggregate = journal.Aggregate{Type: "issue", ID: "FID-a", Sequence: 1}
				e.Effects[0].Table = "PERSON_DECISIONS"
				effect := &e.Effects[0]
				effect.Op = op
				effect.PK = map[string]any{"ID": "PDN-waive"}
				if op == "insert" {
					effect.PK["ID"] = "PDN-forged"
					delete(effect.Values, "kind")
					effect.Values["KiNd"] = "waive"
				} else {
					effect.Values = nil
					if op == "update" {
						effect.Values = map[string]any{"reason": "Changed later"}
					}
				}
				before := handoffProjectionSnapshot(t, v)
				if err := applyLeanEmptyEvent(v, mode, e); err == nil || handoffProjectionSnapshot(t, v) != before {
					t.Fatalf("mutable waiver %s: %v", op, err)
				}
			})
		}
	}
}

func TestWaiverLaneAndHardFloorAdmission(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, lane := range []string{"standard", "assured", "routine"} {
			t.Run(mode+"/"+lane, func(t *testing.T) {
				t.Parallel()
				v, e := waiverArchiveFixture(t)
				var binding waiverrecord.Binding
				if err := json.Unmarshal([]byte(e.Effects[0].Values["binding_json"].(string)), &binding); err != nil {
					t.Fatal(err)
				}
				binding.Lane = lane
				e.Effects[0].Values["binding_json"] = encodedLean(binding)
				if lane == "routine" {
					var raw string
					if err := v.conn.QueryRow(`SELECT record_json FROM handoffs`).Scan(&raw); err != nil {
						t.Fatal(err)
					}
					var h handoffrecord.Record
					if err := json.Unmarshal([]byte(raw), &h); err != nil {
						t.Fatal(err)
					}
					h.Lane = "assured"
					h.ChangedPaths = []string{".fid.toml", ".github/workflows/ci.yml"}
					if _, err := v.conn.Exec(`UPDATE handoffs SET record_json=?`, encodedLean(h)); err != nil {
						t.Fatal(err)
					}
				}
				before := handoffProjectionSnapshot(t, v)
				if err := applyLeanEmptyEvent(v, mode, e); err == nil || handoffProjectionSnapshot(t, v) != before {
					t.Fatalf("forged lane accepted: %v", err)
				}
			})
		}
	}
}
