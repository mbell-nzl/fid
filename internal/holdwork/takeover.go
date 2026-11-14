package holdwork

import (
	"database/sql"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/project"
)

func Takeover(options Options) (Result, error) {
	return write(options, "task.takeover", func(ctx writeContext) (journal.Event, map[string]any, error) {
		var kind, status string
		if err := ctx.conn.QueryRow(`SELECT kind,status FROM issues WHERE id=?`, options.TaskID).Scan(&kind, &status); err != nil {
			if err == sql.ErrNoRows {
				return journal.Event{}, nil, &issues.NotFoundError{IssueID: options.TaskID}
			}
			return journal.Event{}, nil, err
		}
		if kind != "task" {
			return journal.Event{}, nil, &Refusal{Code: "unsupported-kind", TaskID: options.TaskID, Message: "takeover requires a task", Details: map[string]any{"kind": kind}}
		}
		category, err := taskCategory(ctx.conn, options.TaskID)
		if err != nil {
			return journal.Event{}, nil, err
		}
		if category == "completed" || category == "cancelled" {
			return journal.Event{}, nil, &Refusal{Code: "transition-not-permitted", TaskID: options.TaskID, Message: "this task is closed", Details: map[string]any{"from_state": status, "to_state": status, "unsatisfied": []string{"task-open"}}}
		}
		open, err := holds.Open(ctx.conn)
		if err != nil {
			return journal.Event{}, nil, err
		}
		effects := []journal.Effect{}
		var previous *string
		for _, hold := range open {
			if hold.TaskID == options.TaskID {
				if hold.Live(ctx.now) && hold.WorktreeID == project.WorktreeID(options.WorkspaceRoot) && hold.Actor == options.Actor {
					return journal.Event{}, nil, &Refusal{Code: "already-held", TaskID: options.TaskID, Message: options.TaskID + " is already yours in this worktree", Details: map[string]any{"hold_id": hold.ID}}
				}
				previous = &hold.ID
				effects = append(effects, hold.Release("takeover", ctx.now))
				continue
			}
			if hold.WorktreeID != project.WorktreeID(options.WorkspaceRoot) {
				continue
			}
			otherCategory, err := taskCategory(ctx.conn, hold.TaskID)
			if err != nil {
				return journal.Event{}, nil, err
			}
			if hold.Live(ctx.now) || otherCategory == "started" || hold.Kept {
				return journal.Event{}, nil, &Refusal{Code: "worktree-busy", TaskID: options.TaskID, Message: "this worktree already has " + hold.TaskID, Details: map[string]any{"task_id": hold.TaskID, "worktree": options.WorkspaceRoot}}
			}
			effects = append(effects, hold.Release(holds.EndReplaced, ctx.now))
		}
		decisionID, decision, err := persondecision.Effect(ctx.conn, "takeover", options.Actor, "task", options.TaskID, options.Reason, ctx.now)
		if err != nil {
			return journal.Event{}, nil, err
		}
		id, err := holds.NewID(ctx.conn)
		if err != nil {
			return journal.Event{}, nil, err
		}
		hold := holds.New(id, options.TaskID, project.WorktreeID(options.WorkspaceRoot), options.WorkspaceRoot, options.Actor, ctx.now, previous)
		effects = append(effects, decision, hold.Insert())
		return journal.Event{EventType: "task.taken-over", Aggregate: journal.Aggregate{ID: options.TaskID, Type: "task"}, Effects: effects}, map[string]any{"task_id": options.TaskID, "hold_id": id, "previous_hold_id": previous, "decision_id": decisionID}, nil
	})
}
