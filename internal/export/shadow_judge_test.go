package export

import (
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func mergeHoldArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, h := handoffArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsFS(t), 44); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 44
	acceptMergeArchiveEvent(t, v, h)
	at := time.Date(2026, 11, 10, 9, 7, 0, 0, time.UTC)
	return v, journal.Event{EventSchema: 2, EventID: "EVT-20261110-000005", EventType: "task.held", CommandName: "task.hold", OccurredAt: at, ActorV2: journal.EventActor{ActorType: "human", ActorID: "mark", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{ID: "PDN-hold", Type: "decision", Sequence: 1}, Effects: []journal.Effect{{Op: "insert", Table: "merge_holds", PK: map[string]any{"id": "PDN-hold"}, Values: map[string]any{"task_id": "FID-a", "handoff_id": "HOF-100001", "actor": "human:mark", "reason": "Wait for evidence", "at": at.Format(time.RFC3339)}}}}
}

func TestMergeHoldAdmissionEveryPath(t *testing.T) {
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
			{name: "blank reason", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) { e.Effects[0].Values["reason"] = " " }},
			{name: "wrong handoff", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Values["handoff_id"] = "HOF-other"
			}},
			{name: "wrong task", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Values["task_id"] = "FID-other"
			}},
			{name: "wrong actor", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Values["actor"] = "human:other"
			}},
			{name: "wrong time", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Values["at"] = "2026-11-10T09:06:59Z"
			}},
			{name: "not review", change: func(t *testing.T, v *semanticValidator, _ *journal.Event) {
				if _, err := v.conn.Exec(`UPDATE issues SET status='in_progress' WHERE id='FID-a'`); err != nil {
					t.Fatal(err)
				}
			}},
			{name: "extra row", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects = append(e.Effects, e.Effects[0])
			}},
			{name: "legacy", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) { e.EventSchema = 1 }},
			{name: "forged first call", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) { e.Effects[0].Table = "merge_calls" }},
			{name: "alias first call", change: func(_ *testing.T, _ *semanticValidator, e *journal.Event) {
				e.Effects[0].Table = strings.ToUpper("merge_calls")
			}},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := mergeHoldArchiveFixture(t)
				if tc.change != nil {
					tc.change(t, v, &e)
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyLeanEmptyEvent(v, mode, e)
				if (err == nil) != tc.want {
					t.Fatalf("admit=%v want=%t", err, tc.want)
				}
				if !tc.want && before != handoffProjectionSnapshot(t, v) {
					t.Fatal("refusal changed state")
				}
			})
		}
	}
}
