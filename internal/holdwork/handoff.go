package holdwork

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/policysnapshot"
	"github.com/mbell-nzl/fid/internal/project"
)

// Handoff records existing evidence and moves the task and its reservation in
// one Event. No check or network process is started.
func Handoff(options Options) (Result, error) {
	if options.IdempotencyKey == "" {
		return Result{}, fmt.Errorf("handoff requires an idempotency key")
	}
	return write(options, "task.handoff", func(ctx writeContext) (journal.Event, map[string]any, error) {
		hold, err := ctx.current()
		if err != nil {
			return journal.Event{}, nil, err
		}
		record, err := ctx.prepareHandoff(hold)
		if err != nil {
			return journal.Event{}, nil, err
		}
		raw, err := json.Marshal(record)
		if err != nil {
			return journal.Event{}, nil, err
		}
		var revision int
		if err := ctx.conn.QueryRow(`SELECT revision FROM issues WHERE id=?`, hold.TaskID).Scan(&revision); err != nil {
			return journal.Event{}, nil, err
		}
		end := holds.End(hold.ID, "review", ctx.now)
		end.Values["kept"] = true
		effects := []journal.Effect{
			{Op: "insert", Table: "handoffs", PK: map[string]any{"id": record.ID}, Values: map[string]any{"task_id": hold.TaskID, "hold_id": hold.ID, "policy_snapshot": record.PolicySnapshot, "record_json": string(raw)}},
			{Op: "update", Table: "issues", PK: map[string]any{"id": hold.TaskID}, Values: map[string]any{"status": "in_review", "revision": revision + 1, "updated_at": record.At}, Preconditions: map[string]any{"status": "in_progress", "revision": revision}},
			end,
		}
		return journal.Event{EventType: "task.handed-off", Aggregate: journal.Aggregate{ID: hold.TaskID, Type: "task"}, Effects: effects}, record.Fields(), nil
	})
}

func (ctx writeContext) prepareHandoff(hold holds.Hold) (handoffrecord.Record, error) {
	var record handoffrecord.Record
	loaded, err := policy.Load(ctx.options.WorkspaceRoot)
	if err != nil {
		return record, err
	}
	var status string
	var lane sql.NullString
	if err := ctx.conn.QueryRow(`SELECT status,lane FROM issues WHERE id=? AND kind='task'`, hold.TaskID).Scan(&status, &lane); err != nil {
		return record, err
	}
	if status != "in_progress" {
		return record, &Refusal{Code: "transition-not-permitted", TaskID: hold.TaskID, Message: hold.TaskID + " must be started before hand-off", Details: map[string]any{"from_state": status, "to_state": "in_review", "unsatisfied": []string{"task-in-progress"}}}
	}
	if !lane.Valid {
		return record, &Refusal{Code: "lane-not-set", TaskID: hold.TaskID, Message: hold.TaskID + " has no lane", Details: map[string]any{"task_id": hold.TaskID}}
	}
	if _, err := policy.Required(loaded.Effective, lane.String, hold.TaskID, nil); err != nil {
		return record, err
	}
	if !loaded.Effective.TestAdequacySet {
		return record, &checkrecord.Refusal{Code: "test-paths-not-configured", TaskID: hold.TaskID, Message: "no test paths set in .fid.toml", Details: map[string]any{}, Next: "add [review] test_adequacy, then merge .fid.toml and git fetch"}
	}
	observation, err := project.ObserveHead(ctx.options.WorkspaceRoot)
	if err != nil {
		return record, err
	}
	if observation.Dirty {
		return record, &checkrecord.Refusal{Code: "tree-dirty", TaskID: hold.TaskID, Message: hold.TaskID + " has uncommitted changes; commit them before hand-off", Details: map[string]any{"worktree": ctx.options.WorkspaceRoot}}
	}
	comparison, err := policy.Compare(ctx.options.WorkspaceRoot, hold.TaskID)
	if err != nil {
		return record, err
	}
	if comparison.Head != observation.HeadOID {
		return record, fmt.Errorf("repository changed during hand-off")
	}
	required, err := policy.Required(loaded.Effective, lane.String, hold.TaskID, comparison.Paths)
	if err != nil {
		return record, err
	}
	records, err := checkread.Read(ctx.conn, hold.TaskID)
	if err != nil {
		return record, err
	}
	evaluation := checkread.Evaluate(hold.TaskID, observation.TreeOID, required, loaded.Effective.Checks, records)
	if refusal := evaluation.Refusal(hold.TaskID); refusal != nil {
		refusal.Next = "fid check --hold " + hold.ID + " --actor " + ctx.options.Actor + " --idempotency-key <key>"
		return record, refusal
	}
	id, err := journal.GenerateID("HOF", func(candidate string) (bool, error) {
		var found int
		err := ctx.conn.QueryRow(`SELECT 1 FROM handoffs WHERE id=?`, candidate).Scan(&found)
		if err == sql.ErrNoRows {
			return false, nil
		}
		return err == nil, err
	})
	if err != nil {
		return record, err
	}
	if err := policysnapshot.Ensure(ctx.lock, ctx.conn, loaded, ctx.now); err != nil {
		return record, err
	}
	if err := ctx.conn.QueryRow(`SELECT sequence FROM policy_snapshots ORDER BY sequence DESC LIMIT 1`).Scan(&record.PolicySnapshot); err != nil {
		return record, err
	}
	record.View = handoffrecord.View{ID: id, PR: ctx.options.PR, Head: observation.HeadOID, Base: comparison.Base, RequiredChecks: required, ResultIDs: []string{}, At: holds.Timestamp(ctx.now)}
	record.TaskID, record.HoldID, record.Actor, record.Tree, record.Lane = hold.TaskID, hold.ID, ctx.options.Actor, observation.TreeOID, lane.String
	record.DefinitionDigests, record.ChangedPaths = map[string]string{}, comparison.Paths
	for _, name := range required {
		record.DefinitionDigests[name] = checkrecord.Digest(checkrecord.DefinitionFor(loaded.Effective.Checks[name]))
		record.ResultIDs = append(record.ResultIDs, evaluation.LatestResultIDs[name])
	}
	return record, nil
}
