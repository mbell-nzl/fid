package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/mergerecord"
)

func mergeArchiveFixture(t *testing.T) (*semanticValidator, journal.Event, journal.Event) {
	t.Helper()
	return mergeArchiveFixtureWithLane(t, "standard")
}

func mergeArchiveFixtureWithLane(t *testing.T, lane string) (*semanticValidator, journal.Event, journal.Event) {
	t.Helper()
	validator, handoff := handoffArchiveFixtureWithLane(t, lane)
	if err := migrateSemanticProjection(validator.conn, testMigrationsDir(t), 28); err != nil {
		t.Fatal(err)
	}
	var err error
	validator.tables, err = loadSemanticTables(validator.conn)
	if err != nil {
		t.Fatal(err)
	}
	validator.profile = 28
	acceptMergeArchiveEvent(t, validator, handoff)
	m := mergerecord.Merge{ID: "MRG-123456", TaskID: "FID-a", HandoffID: "HOF-100001", PR: "https://github.com/owner/repo/pull/63", Number: 63, Repository: "owner/repo", BaseBranch: "main", MergeCommit: strings.Repeat("d", 40), HeadSHA: strings.Repeat("a", 40), HeadTree: strings.Repeat("c", 40), MergedAt: "2026-11-10T09:05:00Z", ObservedAt: "2026-11-10T09:05:00Z"}
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	merge := archiveStartEvent(t)
	merge.EventID = "EVT-20261110-000004"
	merge.EventType = "merge.observed"
	merge.CommandName = "merge.observe"
	merge.Aggregate = journal.Aggregate{Type: "merge", ID: m.ID, Sequence: 1}
	merge.Idempotency = nil
	merge.Effects = []journal.Effect{{Op: "insert", Table: "merge_observations", PK: map[string]any{"id": m.ID}, Values: map[string]any{"task_id": m.TaskID, "handoff_id": m.HandoffID, "merge_commit": m.MergeCommit, "record_json": string(raw)}}}
	c := mergerecord.CIResult{ID: "CIR-123456", TaskID: m.TaskID, MergeID: m.ID, Repository: m.Repository, RunID: 10, Attempt: 1, WorkflowID: 20, Name: "CI", Path: ".github/workflows/ci.yml", Event: "push", HeadSHA: m.MergeCommit, Status: "completed", Conclusion: "failure", URL: "https://github.com/owner/repo/actions/runs/10", ObservedAt: m.ObservedAt}
	raw, err = json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	ci := archiveStartEvent(t)
	ci.EventID = "EVT-20261110-000005"
	ci.EventType = "ci.result-recorded"
	ci.CommandName = "ci.record"
	ci.Aggregate = journal.Aggregate{Type: "merge", ID: m.ID, Sequence: 2}
	ci.Idempotency = nil
	ci.Effects = []journal.Effect{{Op: "insert", Table: "ci_results", PK: map[string]any{"id": c.ID}, Values: map[string]any{"task_id": c.TaskID, "merge_id": c.MergeID, "repository": c.Repository, "run_id": c.RunID, "attempt": c.Attempt, "record_json": string(raw)}}}
	return validator, merge, ci
}

func acceptMergeArchiveEvent(t *testing.T, v *semanticValidator, event journal.Event) {
	t.Helper()
	raw, err := journal.MarshalCanonicalJSON(event)
	if err != nil {
		t.Fatal(err)
	}
	if err := v.accept(raw); err != nil {
		t.Fatal(err)
	}
}

func mutateObservationJSON(t *testing.T, e *journal.Event, name string, value any) {
	t.Helper()
	var fields map[string]any
	decoder := json.NewDecoder(strings.NewReader(e.Effects[0].Values["record_json"].(string)))
	decoder.UseNumber()
	if err := decoder.Decode(&fields); err != nil {
		t.Fatal(err)
	}
	fields[name] = value
	raw, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	e.Effects[0].Values["record_json"] = string(raw)
}

func observationSnapshot(t *testing.T, v *semanticValidator) string {
	t.Helper()
	var snapshot strings.Builder
	for _, table := range []string{"merge_observations", "ci_results", "projection_meta", "issues", "handoffs", "idempotency_results"} {
		rows, err := v.conn.Query("SELECT * FROM " + table + " ORDER BY rowid")
		if err != nil {
			t.Fatal(err)
		}
		columns, err := rows.Columns()
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			if err := rows.Scan(pointers...); err != nil {
				t.Fatal(err)
			}
			fmt.Fprint(&snapshot, values)
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return snapshot.String()
}

func TestMergeAndCIArchiveAndReplayBindings(t *testing.T) {
	t.Parallel()
	for _, family := range []string{"merge", "ci"} {
		for _, mode := range []string{"live", "replay", "archive"} {
			for _, name := range []string{"valid", "task", "handoff", "PR", "number", "repository", "head", "tree", "time", "unknown field", "aggregate", "service", "command", "retry result", "missing effect", "duplicate effect", "legacy", "uppercase relabel", "not a task", "missing merge", "status", "conclusion", "attempt", "workflow", "name", "path", "event", "URL"} {
				if family == "merge" && (name == "missing merge" || name == "status" || name == "conclusion" || name == "attempt" || name == "workflow" || name == "name" || name == "path" || name == "event" || name == "URL") {
					continue
				}
				if family == "ci" && (name == "handoff" || name == "PR" || name == "number" || name == "tree") {
					continue
				}
				t.Run(family+"/"+mode+"/"+name, func(t *testing.T) {
					t.Parallel()
					v, merge, ci := mergeArchiveFixture(t)
					event := merge
					if family == "ci" {
						acceptMergeArchiveEvent(t, v, merge)
						event = ci
					}
					switch name {
					case "not a task":
						if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-epic','epic','Epic','epic-default','todo','tl','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`); err != nil {
							t.Fatal(err)
						}
						mutateObservationJSON(t, &event, "task_id", "FID-epic")
						event.Effects[0].Values["task_id"] = "FID-epic"
					case "missing merge":
						mutateObservationJSON(t, &event, "merge_id", "MRG-absent")
						event.Effects[0].Values["merge_id"] = "MRG-absent"
						event.Aggregate.ID = "MRG-absent"
					case "status":
						mutateObservationJSON(t, &event, "status", "in_progress")
					case "conclusion":
						mutateObservationJSON(t, &event, "conclusion", "green")
					case "attempt":
						mutateObservationJSON(t, &event, "attempt", 0)
						event.Effects[0].Values["attempt"] = 0
					case "workflow":
						mutateObservationJSON(t, &event, "workflow_id", 0)
					case "name":
						mutateObservationJSON(t, &event, "name", "")
					case "path":
						mutateObservationJSON(t, &event, "path", "")
					case "event":
						mutateObservationJSON(t, &event, "event", "")
					case "URL":
						mutateObservationJSON(t, &event, "url", "https://github.com/other/repo/actions/runs/10")
					case "task":
						mutateObservationJSON(t, &event, "task_id", "FID-other")
						event.Effects[0].Values["task_id"] = "FID-other"
					case "handoff":
						mutateObservationJSON(t, &event, "handoff_id", "HOF-other")
						event.Effects[0].Values["handoff_id"] = "HOF-other"
					case "PR":
						mutateObservationJSON(t, &event, "pr", "https://github.com/owner/repo/pull/64")
						mutateObservationJSON(t, &event, "number", 64)
					case "number":
						mutateObservationJSON(t, &event, "number", 64)
					case "repository":
						mutateObservationJSON(t, &event, "repository", "other/repo")
					case "head":
						if family == "merge" {
							mutateObservationJSON(t, &event, "merge_commit", "abc")
							event.Effects[0].Values["merge_commit"] = "abc"
						} else {
							mutateObservationJSON(t, &event, "head_sha", strings.Repeat("f", 40))
						}
					case "tree":
						mutateObservationJSON(t, &event, "head_tree", strings.Repeat("f", 40))
					case "time":
						mutateObservationJSON(t, &event, "observed_at", "2026-11-10T09:06:00Z")
					case "unknown field":
						mutateObservationJSON(t, &event, "surprise", true)
					case "aggregate":
						event.Aggregate.ID = "MRG-other"
					case "service":
						event.ActorV2.ActorType = "service"
					case "command":
						event.CommandName = "task.start"
					case "retry result":
						takeoverArchiveRetry(&event)
					case "missing effect":
						event.Effects = nil
					case "legacy":
						event.EventSchema = journal.EventSchemaV1
					case "uppercase relabel":
						event.EventType = "issue.updated"
						event.Aggregate.Type = "issue"
						event.Effects[0].Table = strings.ToUpper(event.Effects[0].Table)
					case "duplicate effect":
						event.Effects = append(event.Effects, event.Effects[0])
					}
					before := observationSnapshot(t, v)
					var err error
					switch mode {
					case "archive":
						var raw []byte
						raw, err = journal.MarshalCanonicalJSON(event)
						if err == nil {
							err = v.accept(raw)
						}
					case "live":
						err = journal.DryRunProjectEvent(v.conn, event)
					case "replay":
						// This isolated journal begins at sequence one, for its
						// aggregate and its date; the private projection already
						// contains the referenced hand-off/merge.
						event.Aggregate.Sequence = 1
						event.EventID = "EVT-20261110-000001"
						directory := t.TempDir()
						var raw []byte
						raw, err = journal.MarshalCanonicalJSON(event)
						if err == nil {
							err = os.WriteFile(filepath.Join(directory, "2026-11-10.jsonl"), append(raw, '\n'), 0o600)
						}
						if err == nil {
							_, err = journal.ReplayJournal(v.conn, directory, nil)
						}
					}
					if (err == nil) != (name == "valid") {
						t.Fatalf("accept=%v", err)
					}
					if name != "valid" && observationSnapshot(t, v) != before {
						t.Fatal("refusal changed projection")
					}
				})
			}
		}
	}
}

func TestCIIdentityIntegersSurviveWireReplay(t *testing.T) {
	t.Parallel()
	v, merge, ci := mergeArchiveFixture(t)
	acceptMergeArchiveEvent(t, v, merge)
	// GitHub identities are integers. Values beyond float64's exact range
	// must retain their complete decimal spelling in both JSON and Effects.
	const identity int64 = 9007199254740993
	mutateObservationJSON(t, &ci, "run_id", identity)
	mutateObservationJSON(t, &ci, "attempt", identity)
	mutateObservationJSON(t, &ci, "workflow_id", identity)
	mutateObservationJSON(t, &ci, "url", fmt.Sprintf("https://github.com/owner/repo/actions/runs/%d", identity))
	ci.Effects[0].Values["run_id"] = identity
	ci.Effects[0].Values["attempt"] = identity
	acceptMergeArchiveEvent(t, v, ci)
	var run, attempt int64
	if err := v.conn.QueryRow(`SELECT run_id,attempt FROM ci_results`).Scan(&run, &attempt); err != nil {
		t.Fatal(err)
	}
	if run != identity || attempt != identity {
		t.Fatalf("replayed identities=%d/%d", run, attempt)
	}
}
