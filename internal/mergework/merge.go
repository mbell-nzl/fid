package mergework

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/mergerecord"
)

// RecordMerge binds a confirmed PR and head tree to a specific hand-off. The
// natural key (hand-off, merge commit) returns the original observation on retry.
func RecordMerge(options Options, handoff string, pr github.PullRequest, headTree string) (mergerecord.Merge, error) {
	var result mergerecord.Merge
	if !pr.Merged || pr.State != "MERGED" || pr.MergedAt == nil {
		return result, fmt.Errorf("merge observation requires a confirmed merge")
	}
	err := withWrite(options, func(ctx writeContext) (journal.Event, error) {
		record := mergerecord.Merge{TaskID: options.TaskID, HandoffID: handoff, PR: pr.URL, Number: pr.Number, Repository: pr.BaseRepository.String(), BaseBranch: pr.BaseRef, MergeCommit: pr.MergeCommit, HeadSHA: pr.HeadSHA, HeadTree: headTree, MergedAt: pr.MergedAt.UTC().Format(time.RFC3339), ObservedAt: ctx.at}
		existing, err := readMerge(ctx.conn, `SELECT record_json FROM merge_observations WHERE handoff_id=? AND merge_commit=?`, handoff, pr.MergeCommit)
		if err == nil {
			record.ID, record.ObservedAt = existing.ID, existing.ObservedAt
			if !reflect.DeepEqual(record, existing) {
				return journal.Event{}, fmt.Errorf("merge observation retry differs from original")
			}
			result = existing
			return journal.Event{}, nil
		}
		if err != sql.ErrNoRows {
			return journal.Event{}, err
		}
		record.ID, err = ctx.id("MRG", "merge_observations")
		if err != nil {
			return journal.Event{}, err
		}
		if err := record.Validate(); err != nil {
			return journal.Event{}, err
		}
		raw, err := journal.MarshalCanonicalJSON(record)
		if err != nil {
			return journal.Event{}, err
		}
		result = record
		effect := journal.Effect{Op: "insert", Table: "merge_observations", PK: map[string]any{"id": record.ID}, Values: map[string]any{"task_id": record.TaskID, "handoff_id": record.HandoffID, "merge_commit": record.MergeCommit, "record_json": string(raw)}}
		return journal.Event{EventType: "merge.observed", CommandName: "merge.observe", Aggregate: journal.Aggregate{ID: record.ID, Type: "merge"}, Effects: []journal.Effect{effect}}, nil
	})
	if err != nil {
		return mergerecord.Merge{}, fmt.Errorf("record merge observation: %w", err)
	}
	return result, nil
}

// RecordMergeLocked records within the caller's already-refreshed ledger lock.
// The caller must finish stale-observation checks before invoking it.
func RecordMergeLocked(lock *journal.WriteLock, conn *sql.DB, options Options, handoff string, pr github.PullRequest, headTree string) (mergerecord.Merge, error) {
	options.lock, options.conn = lock, conn
	return RecordMerge(options, handoff, pr, headTree)
}
