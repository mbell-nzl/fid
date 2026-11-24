package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/journal"
)

func completionArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, m, ci := mergeArchiveFixtureWithLane(t, "routine")
	if err := migrateSemanticProjection(v.conn, testMigrationsFS(t), 34); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 34
	acceptMergeArchiveEvent(t, v, m)
	mutateObservationJSON(t, &ci, "conclusion", "success")
	acceptMergeArchiveEvent(t, v, ci)
	at := time.Date(2026, 11, 10, 9, 8, 0, 0, time.UTC)
	r := completionrecord.Record{ID: "CMP-123456", TaskID: "FID-a", HandoffID: "HOF-100001", MergeID: "MRG-123456", CIResultIDs: []string{"CIR-123456"}, ObservedCIResultIDs: []string{"CIR-123456"}, RoundIDs: []string{}, Actor: "agent:finisher", At: at.Format(time.RFC3339)}
	raw, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	e := journal.Event{EventSchema: 2, EventID: "EVT-20261110-000006", EventType: "task.finished", CommandName: "task.done", OccurredAt: at, ActorV2: journal.EventActor{ActorType: "agent", ActorID: "finisher", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{Type: "task", ID: r.TaskID, Sequence: 2}, Effects: []journal.Effect{
		{Op: "insert", Table: "completions", PK: map[string]any{"id": r.ID}, Values: map[string]any{"task_id": r.TaskID, "handoff_id": r.HandoffID, "merge_id": r.MergeID, "record_json": string(raw)}},
		{Op: "update", Table: "issues", PK: map[string]any{"id": r.TaskID}, Values: map[string]any{"status": "done", "revision": 3, "updated_at": r.At}, Preconditions: map[string]any{"status": "in_review", "revision": 2}},
		{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-1"}, Values: map[string]any{"kept": false}, Preconditions: map[string]any{"kept": true}},
	}}
	return v, e
}

func TestCompletionAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, name := range []string{"valid", "other task", "other handoff", "other merge", "missing CI", "wrong observed CI", "duplicate CI", "missing rounds", "extra rounds", "wrong approval", "wrong fingerprint", "model disclosure", "unknown field", "actor", "command", "legacy", "relabel family", "relabel table", "no move", "wrong state", "wrong revision", "other release", "release fields", "omitted release", "duplicate release", "missing question", "cancelled question", "answered question"} {
			t.Run(mode+"/"+name, func(t *testing.T) {
				t.Parallel()
				v, e := completionArchiveFixture(t)
				switch name {
				case "other task":
					mutateObservationJSON(t, &e, "task_id", "FID-other")
				case "other handoff":
					mutateObservationJSON(t, &e, "handoff_id", "HOF-other")
					e.Effects[0].Values["handoff_id"] = "HOF-other"
				case "other merge":
					mutateObservationJSON(t, &e, "merge_id", "MRG-other")
					e.Effects[0].Values["merge_id"] = "MRG-other"
				case "missing CI":
					mutateObservationJSON(t, &e, "ci_result_ids", []string{"CIR-missing"})
				case "wrong observed CI":
					mutateObservationJSON(t, &e, "observed_ci_result_ids", []string{"CIR-missing"})
				case "duplicate CI":
					mutateObservationJSON(t, &e, "ci_result_ids", []string{"CIR-123456", "CIR-123456"})
				case "missing rounds":
					mutateObservationJSON(t, &e, "round_ids", nil)
				case "extra rounds":
					mutateObservationJSON(t, &e, "round_ids", []string{"RND-missing"})
				case "wrong approval":
					mutateObservationJSON(t, &e, "approval_id", "PDN-missing")
				case "wrong fingerprint":
					mutateObservationJSON(t, &e, "fingerprint", map[string]any{"algorithm": "patch-id-verbatim-per-file-v1", "fingerprint": strings.Repeat("a", 64), "base": strings.Repeat("d", 40), "head": strings.Repeat("a", 40)})
				case "model disclosure":
					mutateObservationJSON(t, &e, "no_other_model_review", true)
				case "unknown field":
					mutateObservationJSON(t, &e, "invented", true)
				case "actor":
					e.ActorV2.ActorType = "service"
				case "command":
					e.CommandName = "issue.transition"
				case "legacy":
					e.EventSchema = 1
				case "relabel family":
					e.EventType = "issue.updated"
				case "relabel table":
					e.Effects[0].Table = "Completions"
				case "no move":
					e.Effects = append(e.Effects[:1], e.Effects[2:]...)
				case "wrong state":
					e.Effects[1].Values["status"] = "cancelled"
				case "wrong revision":
					e.Effects[1].Values["revision"] = 9
				case "other release":
					e.Effects[2].PK["id"] = "HLD-other"
				case "release fields":
					e.Effects[2].Values["end_reason"] = "task-completed"
				case "omitted release":
					e.Effects = e.Effects[:2]
				case "duplicate release":
					e.Effects = append(e.Effects, e.Effects[2])
				case "missing question", "cancelled question", "answered question":
					if _, err := v.conn.Exec(`INSERT INTO questions(id,task_id,hold_id,text,asker,asked_at) VALUES('QST-open','FID-a','HLD-1','Question?','agent:builder','2026-11-10T09:00:00Z')`); err != nil {
						t.Fatal(err)
					}
					if name != "missing question" {
						resolution := "cancelled"
						if name == "answered question" {
							resolution = "answered"
						}
						e.Effects = append(e.Effects, journal.Effect{Op: "update", Table: "questions", PK: map[string]any{"id": "QST-open"}, Values: map[string]any{"resolved_at": e.OccurredAt.Format(time.RFC3339), "resolution": resolution, "decision_id": nil}, Preconditions: map[string]any{"resolved_at": nil}})
					}
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyLeanEmptyEvent(v, mode, e)
				if (err == nil) != (name == "valid") {
					t.Fatalf("admit=%v", err)
				}
				if name != "valid" && handoffProjectionSnapshot(t, v) != before {
					t.Fatal("refused completion changed projection")
				}
			})
		}
	}
}
