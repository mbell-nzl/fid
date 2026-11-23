package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/defaultconflict"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func defaultConflictArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, round := reviewArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsDir(t), 38); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 38
	var r reviewrecord.Round
	if err := json.Unmarshal([]byte(round.Effects[0].Values["record_json"].(string)), &r); err != nil {
		t.Fatal(err)
	}
	r.Verdict, r.Findings = "APPROVE", []reviewrecord.Finding{}
	setReviewArchiveRound(t, &round, r)
	acceptMergeArchiveEvent(t, v, round)
	at := time.Date(2026, 11, 10, 9, 7, 0, 0, time.UTC)
	record := defaultconflict.Record{Evidence: defaultconflict.Evidence{Cause: "default-conflict", DefaultRef: "origin/main", DefaultCommit: strings.Repeat("c", 40), HandoffID: r.HandoffID, Head: r.Binding.Head, ConflictingPaths: []string{"notes.md"}}, Binding: fingerprint.Result{Algorithm: r.Binding.Algorithm, Fingerprint: r.Binding.Fingerprint, Base: r.Binding.Base, Head: r.Binding.Head}}
	effects, err := journal.DefaultConflictEffects(v.conn, "FID-a", "HLD-1", "DFC-123456", record, at)
	if err != nil {
		t.Fatal(err)
	}
	return v, journal.Event{EventID: "EVT-20261110-000005", EventSchema: 2, EventType: "task.default-conflict", CommandName: "next.claim", OccurredAt: at, ActorV2: journal.EventActor{ActorID: "builder", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{ID: "DFC-123456", Type: "default-conflict", Sequence: 1}, Effects: effects}
}

func TestDefaultConflictAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name     string
			change   func(*journal.Event)
			evidence func(*defaultconflict.Record)
			sql      string
			valid    bool
		}{
			{name: "complete without git objects", valid: true},
			{name: "raw path ordering", valid: true, evidence: func(r *defaultconflict.Record) {
				r.ConflictingPaths = []string{"x💛", "x�"}
				r.PathBytes = [][]byte{[]byte("x💛"), {'x', 0xff}}
			}},
			{name: "mismatched path bytes", evidence: func(r *defaultconflict.Record) { r.PathBytes = [][]byte{[]byte("other")} }},
			{name: "missing state", change: func(e *journal.Event) { e.Effects = append(e.Effects[:1], e.Effects[2]) }},
			{name: "missing reservation", change: func(e *journal.Event) { e.Effects = e.Effects[:2] }},
			{name: "wrong revision", change: func(e *journal.Event) { e.Effects[1].Values["revision"] = 99 }},
			{name: "wrong reservation", change: func(e *journal.Event) { e.Effects[2].PK["id"] = "HLD-other" }},
			{name: "before handoff", change: func(e *journal.Event) {
				e.OccurredAt = e.OccurredAt.Add(-time.Hour)
				stamp := e.OccurredAt.UTC().Format(time.RFC3339)
				e.Effects[1].Values["updated_at"] = stamp
				e.Effects[2].Values["rework_at"] = stamp
			}},
			{name: "wrong time", change: func(e *journal.Event) { e.Effects[2].Values["rework_at"] = "2026-11-10T09:08:00Z" }},
			{name: "wrong command", change: func(e *journal.Event) { e.CommandName = "review.record" }},
			{name: "service actor", change: func(e *journal.Event) { e.ActorV2.ActorType = "service" }},
			{name: "wrong aggregate", change: func(e *journal.Event) { e.Aggregate.ID = "DFC-other" }},
			{name: "another family", change: func(e *journal.Event) {
				e.EventType = "issue.updated"
				e.Aggregate.Type = "issue"
				e.Aggregate.ID = "FID-a"
			}},
			{name: "legacy family", change: func(e *journal.Event) { e.EventSchema = 1 }},
			{name: "wrong head", evidence: func(r *defaultconflict.Record) { r.Head = strings.Repeat("d", 40) }},
			{name: "wrong handoff", evidence: func(r *defaultconflict.Record) { r.HandoffID = "HOF-old" }},
			{name: "missing commit", evidence: func(r *defaultconflict.Record) { r.DefaultCommit = "" }},
			{name: "wrong cause", evidence: func(r *defaultconflict.Record) { r.Cause = "minor-fix" }},
			{name: "not origin default", evidence: func(r *defaultconflict.Record) { r.DefaultRef = "main" }},
			{name: "no conflicts", evidence: func(r *defaultconflict.Record) { r.ConflictingPaths = []string{} }},
			{name: "repeated path", evidence: func(r *defaultconflict.Record) { r.ConflictingPaths = []string{"a", "a"} }},
			{name: "stale reviews", evidence: func(r *defaultconflict.Record) { r.Binding.Fingerprint = strings.Repeat("e", 64) }},
			{name: "binding wrong base", evidence: func(r *defaultconflict.Record) { r.Binding.Base = strings.Repeat("e", 40) }},
			{name: "released reservation", sql: `UPDATE holds SET kept=0`},
			{name: "review due", sql: `DELETE FROM review_rounds`},
			{name: "recorded merge", sql: `INSERT INTO merge_observations VALUES('MRG-123456','FID-a','HOF-100001','aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa','{}')`},
			{name: "wrong state", sql: `UPDATE issues SET status='done'`},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := defaultConflictArchiveFixture(t)
				valid := e
				raw, _ := json.Marshal(e)
				if err := json.Unmarshal(raw, &e); err != nil {
					t.Fatal(err)
				}
				if tc.sql != "" {
					if _, err := v.conn.Exec(tc.sql); err != nil {
						t.Fatal(err)
					}
				}
				if tc.change != nil {
					tc.change(&e)
				}
				if tc.evidence != nil {
					var r defaultconflict.Record
					if err := json.Unmarshal([]byte(e.Effects[0].Values["record_json"].(string)), &r); err != nil {
						t.Fatal(err)
					}
					tc.evidence(&r)
					e.Effects[0].Values["record_json"] = encodedLean(r)
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyReworkEvent(v, mode, e)
				if (err == nil) != tc.valid {
					t.Fatalf("%s admission=%v, want valid=%t", mode, err, tc.valid)
				}
				if !tc.valid {
					if handoffProjectionSnapshot(t, v) != before {
						t.Fatal("refused default conflict changed projection")
					}
					if tc.sql == "" {
						if err := applyReworkEvent(v, mode, valid); err != nil {
							t.Fatalf("valid write after refusal: %v", err)
						}
					}
				}
			})
		}
	}
}
