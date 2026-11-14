package holdwork

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/policysnapshot"
)

// Start validates the exact live session before moving its planned task and
// its epic. The journal Event holds both state changes atomically.
func Start(options Options) (Result, error) {
	if options.IdempotencyKey == "" {
		return Result{}, fmt.Errorf("start requires an idempotency key")
	}
	return write(options, "task.start", func(ctx writeContext) (journal.Event, map[string]any, error) {
		hold, err := ctx.current()
		if err != nil {
			return journal.Event{}, nil, err
		}
		loaded, err := policy.Load(options.WorkspaceRoot)
		if err != nil {
			return journal.Event{}, nil, err
		}
		var status string
		var lane, parent sql.NullString
		if err := ctx.conn.QueryRow(`SELECT status, lane, parent_id FROM issues WHERE id=? AND kind='task'`, hold.TaskID).Scan(&status, &lane, &parent); err != nil {
			return journal.Event{}, nil, fmt.Errorf("read held task: %w", err)
		}
		if !lane.Valid {
			return journal.Event{}, nil, &Refusal{Code: "lane-not-set", TaskID: hold.TaskID, Message: hold.TaskID + " has no lane", Details: map[string]any{"task_id": hold.TaskID}}
		}
		if status != "todo" && status != "in_progress" {
			return journal.Event{}, nil, &Refusal{Code: "transition-not-permitted", TaskID: hold.TaskID, Message: hold.TaskID + " cannot start from " + strings.ReplaceAll(status, "_", " "), Details: map[string]any{"from_state": status, "to_state": "in_progress", "unsatisfied": []string{"task-todo"}}}
		}
		if err := policysnapshot.Ensure(ctx.lock, ctx.conn, loaded, ctx.now); err != nil {
			return journal.Event{}, nil, err
		}
		fields := map[string]any{"task_id": hold.TaskID, "status": "in_progress", "hold_id": hold.ID, "lane": lane.String}
		if status == "in_progress" {
			return journal.Event{}, fields, nil
		}
		effects, err := ctx.startIssue(hold.TaskID)
		if err != nil {
			return journal.Event{}, nil, err
		}
		if parent.Valid {
			var state string
			if err := ctx.conn.QueryRow(`SELECT status FROM issues WHERE id=? AND kind='epic'`, parent.String).Scan(&state); err != nil {
				return journal.Event{}, nil, fmt.Errorf("read task epic: %w", err)
			}
			if state == "todo" {
				epic, err := ctx.startIssue(parent.String)
				if err != nil {
					return journal.Event{}, nil, err
				}
				effects = append(effects, epic...)
			}
		}
		return journal.Event{EventType: "task.started", Aggregate: journal.Aggregate{ID: hold.TaskID, Type: "task"}, Effects: effects}, fields, nil
	})
}

func (ctx writeContext) startIssue(id string) ([]journal.Effect, error) {
	var revision int
	if err := ctx.conn.QueryRow(`SELECT revision FROM issues WHERE id=?`, id).Scan(&revision); err != nil {
		return nil, fmt.Errorf("read starting issue: %w", err)
	}
	return []journal.Effect{{Op: "update", Table: "issues", PK: map[string]any{"id": id}, Values: map[string]any{"status": "in_progress", "revision": revision + 1, "updated_at": holdsTimestamp(ctx.now)}, Preconditions: map[string]any{"status": "todo", "revision": revision}}}, nil
}
