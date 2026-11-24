package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func closureArchiveFixture(t *testing.T) (*semanticValidator, journal.Event, reviewrecord.Round) {
	t.Helper()
	v, initial := reviewArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsFS(t), 31); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 31
	acceptMergeArchiveEvent(t, v, initial)
	var r reviewrecord.Round
	if err := json.Unmarshal([]byte(initial.Effects[0].Values["record_json"].(string)), &r); err != nil {
		t.Fatal(err)
	}
	r.ID, r.Findings, r.Verdict = "RND-654321", []reviewrecord.Finding{}, "APPROVE"
	r.Dispositions = []reviewrecord.Disposition{{FindingID: "FND-1", Disposition: "resolved", Reason: "Regression catches the fault"}, {FindingID: "FND-2", Disposition: "not-applicable", Reason: "Outside this change"}, {FindingID: "FND-3", Disposition: "open", Reason: "Optional follow-up"}}
	e := initial
	e.EventID, e.Aggregate.ID = "EVT-20261110-000005", r.ID
	setReviewArchiveRound(t, &e, r)
	e.Idempotency.Key = "later"
	e.Effects[len(e.Effects)-1].PK["idempotency_key"] = "later"
	return v, e, r
}

func TestRoundDispositionAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name   string
			change func(*testing.T, *semanticValidator, *reviewrecord.Round)
			want   bool
		}{
			{name: "complete", want: true},
			{name: "open major", change: func(_ *testing.T, _ *semanticValidator, r *reviewrecord.Round) {
				r.Dispositions[0].Disposition = "open"
				r.Verdict = "CHANGES"
			}, want: true},
			{name: "omitted major", change: func(_ *testing.T, _ *semanticValidator, r *reviewrecord.Round) { r.Dispositions = r.Dispositions[1:] }},
			{name: "changes after closing all", change: func(_ *testing.T, _ *semanticValidator, r *reviewrecord.Round) { r.Verdict = "CHANGES" }},
			{name: "same round", change: func(_ *testing.T, _ *semanticValidator, r *reviewrecord.Round) {
				r.Findings = []reviewrecord.Finding{{ID: "FND-own", FindingInput: reviewrecord.FindingInput{Profile: "code", Severity: "major", Summary: "New fault", Evidence: "Observed", Fix: "Fix"}}}
				r.Dispositions[0].FindingID = "FND-own"
				r.Verdict = "CHANGES"
			}},
			{name: "future finding", change: func(_ *testing.T, _ *semanticValidator, r *reviewrecord.Round) {
				r.Dispositions[0].FindingID = "FND-future"
			}},
			{name: "another task", change: func(t *testing.T, v *semanticValidator, r *reviewrecord.Round) {
				if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-other','task','Other','task-default','todo','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z'); INSERT INTO review_rounds VALUES('RND-other','FID-other','HOF-100001','{}'); INSERT INTO round_findings VALUES('FND-other','RND-other','{}')`); err != nil {
					t.Fatal(err)
				}
				r.Dispositions[0].FindingID = "FND-other"
			}},
			{name: "builder", change: func(_ *testing.T, _ *semanticValidator, r *reviewrecord.Round) { r.Reviewer.Actor = "agent:builder" }},
			{name: "empty reason", change: func(_ *testing.T, _ *semanticValidator, r *reviewrecord.Round) { r.Dispositions[0].Reason = " " }},
			{name: "duplicate", change: func(_ *testing.T, _ *semanticValidator, r *reviewrecord.Round) {
				r.Dispositions = append(r.Dispositions, r.Dispositions[0])
			}},
			{name: "answered is not closure", change: func(_ *testing.T, _ *semanticValidator, r *reviewrecord.Round) {
				r.Dispositions[0].Disposition = "answered"
			}},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e, r := closureArchiveFixture(t)
				if tc.change != nil {
					tc.change(t, v, &r)
				}
				setReviewArchiveRound(t, &e, r)
				e.Idempotency.Key = "later"
				e.Effects[len(e.Effects)-1].PK["idempotency_key"] = "later"
				before := handoffProjectionSnapshot(t, v)
				err := applyLeanEmptyEvent(v, mode, e)
				if (err == nil) != tc.want {
					t.Fatalf("admit=%v want=%t", err, tc.want)
				}
				if !tc.want && handoffProjectionSnapshot(t, v) != before {
					t.Fatal("refusal changed state")
				}
				if tc.want {
					open, err := reviewrecord.OpenMajors(v.conn, r.TaskID)
					want := 0
					if tc.name == "open major" {
						want = 1
					}
					if err != nil || len(open) != want {
						t.Fatalf("open majors=%v/%v", open, err)
					}
				}
			})
		}
	}
}

func personClosureEvent(t *testing.T) journal.Event {
	t.Helper()
	at := time.Date(2026, 11, 10, 9, 7, 0, 0, time.UTC)
	e := journal.Event{EventSchema: 2, EventID: "EVT-20261110-000005", EventType: "finding.closed", CommandName: "finding.close", OccurredAt: at, ActorV2: journal.EventActor{ActorID: "mark", ActorType: "human", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{ID: "FND-1", Type: "finding", Sequence: 1}, Effects: []journal.Effect{
		{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-close"}, Values: map[string]any{"kind": "close-finding", "actor": "human:mark", "subject_type": "task", "subject_id": "FID-a", "binding_id": "FND-1", "reason": "Accepted risk", "decided_at": at.Format(time.RFC3339)}},
		reviewwork.ClosureEffect("FND-1", "", "PDN-close", "closed", "Accepted risk"),
	}}
	return e
}

func TestPersonClosureAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name   string
			change func(*journal.Event)
			want   bool
		}{
			{name: "complete", want: true},
			{name: "agent", change: func(e *journal.Event) { e.ActorV2.ActorType = "agent"; e.Effects[0].Values["actor"] = "agent:mark" }},
			{name: "service", change: func(e *journal.Event) { e.ActorV2.ActorType = "service"; e.Effects[0].Values["actor"] = "service:mark" }},
			{name: "reason", change: func(e *journal.Event) { e.Effects[0].Values["reason"] = " "; e.Effects[1].Values["reason"] = " " }},
			{name: "missing finding", change: func(e *journal.Event) {
				e.Aggregate.ID = "FND-future"
				e.Effects[0].Values["binding_id"] = "FND-future"
				e.Effects[1].PK["finding_id"] = "FND-future"
			}},
			{name: "wrong subject", change: func(e *journal.Event) { e.Effects[0].Values["subject_id"] = "FID-other" }},
			{name: "wrong binding", change: func(e *journal.Event) { e.Effects[0].Values["binding_id"] = "FND-2" }},
			{name: "wrong kind", change: func(e *journal.Event) { e.Effects[0].Values["kind"] = "answer" }},
			{name: "wrong closure", change: func(e *journal.Event) { e.Effects[1].Values["decision_id"] = "PDN-other" }},
			{name: "legacy", change: func(e *journal.Event) { e.EventSchema = 1 }},
			{name: "other family", change: func(e *journal.Event) { e.EventType = "issue.updated" }},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, _, _ := closureArchiveFixture(t)
				e := personClosureEvent(t)
				if tc.change != nil {
					tc.change(&e)
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyLeanEmptyEvent(v, mode, e)
				if (err == nil) != tc.want {
					t.Fatalf("admit=%v want=%t", err, tc.want)
				}
				if !tc.want && handoffProjectionSnapshot(t, v) != before {
					t.Fatal("refusal changed state")
				}
				if tc.want {
					before = handoffProjectionSnapshot(t, v)
					e.EventID = "EVT-20261110-000006"
					e.Effects[0].PK["id"] = "PDN-again"
					e.Effects[1].Values["decision_id"] = "PDN-again"
					if err := applyLeanEmptyEvent(v, mode, e); err == nil || handoffProjectionSnapshot(t, v) != before {
						t.Fatalf("double closure=%v", err)
					}
					// A later raw pass must not conceal the original immutable CHANGES.
					rounds, err := reviewrecord.Read(v.conn, "FID-a")
					if err != nil || len(rounds) != 1 || !strings.Contains(rounds[0].Verdict, "CHANGES") {
						t.Fatal(rounds, err)
					}
				}
			})
		}
	}
}

func TestPersonClosureRefusesLegacyDecisionOnlyWriter(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name, column string
			inKey        bool
		}{
			{"values", "kind", false},
			{"values alias", "KiNd", false},
			{"key", "kind", true},
			{"key alias", "KiNd", true},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, _, _ := closureArchiveFixture(t)
				e := personClosureEvent(t)
				e.EventSchema = 1
				e.EventType = "issue.updated"
				e.Aggregate.Type = "issue"
				e.Aggregate.ID = "FID-a"
				e.Effects = e.Effects[:1]
				e.Effects[0].Table = "PERSON_DECISIONS"
				delete(e.Effects[0].Values, "kind")
				row := e.Effects[0].Values
				if tc.inKey {
					row = e.Effects[0].PK
				}
				row[tc.column] = "close-finding"
				before := handoffProjectionSnapshot(t, v)
				if err := applyLeanEmptyEvent(v, mode, e); err == nil || handoffProjectionSnapshot(t, v) != before {
					t.Fatalf("legacy person's closure=%v", err)
				}
			})
		}
	}
}

func TestClosureDecisionIsImmutableThroughLegacyAliases(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			v, _, _ := closureArchiveFixture(t)
			acceptMergeArchiveEvent(t, v, personClosureEvent(t))
			e := journal.Event{OccurredAt: personClosureEvent(t).OccurredAt, EventSchema: 1, EventID: "EVT-20261110-000006", EventType: "issue.updated", Aggregate: journal.Aggregate{ID: "FID-a", Type: "issue", Sequence: 1}, Effects: []journal.Effect{{Op: "update", Table: "PERSON_DECISIONS", PK: map[string]any{"ID": "PDN-close"}, Values: map[string]any{"reason": "Changed after closure"}}}}
			before := handoffProjectionSnapshot(t, v)
			if err := applyLeanEmptyEvent(v, mode, e); err == nil || handoffProjectionSnapshot(t, v) != before {
				t.Fatalf("mutated immutable decision: %v", err)
			}
		})
	}
}
