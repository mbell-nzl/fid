package mergework

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/mergerecord"
)

// RecordCI records completed attempts, including failures and superseded
// attempts. Pending attempts produce no Event. A malformed batch changes nothing.
// Natural keys replace a command retry cache: a pending observation stays retryable.
func RecordCI(options Options, mergeID string, repo github.Repository, facts github.CI) ([]mergerecord.CIResult, error) {
	results := []mergerecord.CIResult{}
	err := withWrite(options, func(ctx writeContext) (journal.Event, error) {
		merge, err := readMerge(ctx.conn, `SELECT record_json FROM merge_observations WHERE id=?`, mergeID)
		if err != nil {
			return journal.Event{}, err
		}
		if merge.TaskID != options.TaskID || merge.Repository != repo.String() || merge.MergeCommit != facts.Commit {
			return journal.Event{}, fmt.Errorf("CI facts differ from merge observation")
		}
		effects := []journal.Effect{}
		seen := map[[2]int64]mergerecord.CIResult{}
		for _, run := range facts.Runs {
			record := mergerecord.CIResult{ID: "pending", TaskID: options.TaskID, MergeID: mergeID, Repository: repo.String(), RunID: run.ID, Attempt: run.Attempt, WorkflowID: run.WorkflowID, Name: run.Name, Path: run.Path, Event: run.Event, HeadSHA: run.HeadSHA, Status: run.Status, Conclusion: run.Conclusion, URL: run.URL, ObservedAt: ctx.at}
			if record.HeadSHA != merge.MergeCommit {
				return journal.Event{}, fmt.Errorf("CI head differs from merge commit")
			}
			if run.Status != "completed" {
				if run.Conclusion != "" || !slices.Contains([]string{"queued", "in_progress", "waiting", "pending", "requested"}, run.Status) {
					return journal.Event{}, fmt.Errorf("invalid pending CI status")
				}
				record.Status, record.Conclusion = "completed", "success"
				if err := record.Validate(); err != nil {
					return journal.Event{}, err
				}
				continue
			}
			if err := record.Validate(); err != nil {
				return journal.Event{}, err
			}
			key := [2]int64{run.ID, run.Attempt}
			if prior, exists := seen[key]; exists {
				record.ID, record.ObservedAt = prior.ID, prior.ObservedAt
				if !reflect.DeepEqual(prior, record) {
					return journal.Event{}, fmt.Errorf("inconsistent duplicate CI attempt")
				}
				continue
			}
			var raw string
			err := ctx.conn.QueryRow(`SELECT record_json FROM ci_results WHERE repository=? AND run_id=? AND attempt=? AND merge_id=?`, repo.String(), run.ID, run.Attempt, mergeID).Scan(&raw)
			switch err {
			case nil:
				var existing mergerecord.CIResult
				if err := json.Unmarshal([]byte(raw), &existing); err != nil {
					return journal.Event{}, err
				}
				if err := existing.Validate(); err != nil {
					return journal.Event{}, err
				}
				record.ID, record.ObservedAt = existing.ID, existing.ObservedAt
				if !reflect.DeepEqual(existing, record) {
					return journal.Event{}, fmt.Errorf("CI attempt retry differs from original")
				}
			case sql.ErrNoRows:
				record.ID, err = ctx.id("CIR", "ci_results")
				if err != nil {
					return journal.Event{}, err
				}
				encoded, err := journal.MarshalCanonicalJSON(record)
				if err != nil {
					return journal.Event{}, err
				}
				effects = append(effects, journal.Effect{Op: "insert", Table: "ci_results", PK: map[string]any{"id": record.ID}, Values: map[string]any{"task_id": record.TaskID, "merge_id": record.MergeID, "repository": record.Repository, "run_id": record.RunID, "attempt": record.Attempt, "record_json": string(encoded)}})
			default:
				return journal.Event{}, err
			}
			seen[key] = record
			results = append(results, record)
		}
		if len(effects) == 0 {
			return journal.Event{}, nil
		}
		return journal.Event{EventType: "ci.result-recorded", CommandName: "ci.record", Aggregate: journal.Aggregate{ID: mergeID, Type: "merge"}, Effects: effects}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("record CI results: %w", err)
	}
	return results, nil
}
