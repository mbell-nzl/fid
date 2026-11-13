package holdwork

import (
	"database/sql"
	"fmt"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/questions"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func (ctx writeContext) current() (holds.Hold, error) {
	hold, err := holds.CurrentByID(ctx.conn, ctx.options.HoldID)
	reason := ""
	switch {
	case err == sql.ErrNoRows:
		reason = "not-found"
	case err != nil:
		return holds.Hold{}, err
	case hold.EndedAt != nil:
		reason = "ended"
	case hold.Lapsed(ctx.now):
		reason = "expired"
	case hold.Actor != ctx.options.Actor:
		reason = "actor-mismatch"
	case hold.WorktreeID != project.WorktreeID(ctx.options.WorkspaceRoot):
		reason = "worktree-mismatch"
	}
	if reason != "" {
		return holds.Hold{}, &Refusal{Code: "hold-not-current", TaskID: hold.TaskID, Message: "this hold is no longer current", Details: map[string]any{"hold_id": ctx.options.HoldID, "reason": reason}}
	}
	return hold, nil
}

func Heartbeat(options Options) (Result, error) {
	return write(options, "hold.heartbeat", func(ctx writeContext) (journal.Event, map[string]any, error) {
		hold, err := ctx.current()
		if err != nil {
			return journal.Event{}, nil, err
		}
		stamp, expires := holds.Timestamp(ctx.now), holds.Timestamp(ctx.now.Add(holds.Lifetime))
		effect := journal.Effect{Op: "update", Table: "holds", PK: map[string]any{"id": hold.ID}, Values: map[string]any{"last_heartbeat_at": stamp, "expires_at": expires}, Preconditions: map[string]any{"ended_at": nil, "expires_at": holds.Timestamp(hold.ExpiresAt)}}
		return journal.Event{EventType: "hold.heartbeat", Aggregate: journal.Aggregate{ID: hold.ID, Type: "hold"}, Effects: []journal.Effect{effect}}, map[string]any{"task_id": hold.TaskID, "hold_id": hold.ID, "last_heartbeat_at": stamp, "expires_at": expires, "heartbeat_interval_seconds": holds.HeartbeatIntervalSeconds}, nil
	})
}

func Park(options Options) (Result, error) {
	return write(options, "hold.park", func(ctx writeContext) (journal.Event, map[string]any, error) {
		hold, err := ctx.current()
		if err != nil {
			return journal.Event{}, nil, err
		}
		if options.Question != "" {
			return ctx.parkWaiting(hold)
		}
		category, err := taskCategory(ctx.conn, hold.TaskID)
		if err != nil {
			return journal.Event{}, nil, err
		}
		kept := category == "started"
		history, err := questions.Read(ctx.conn, hold.TaskID)
		if err != nil {
			return journal.Event{}, nil, err
		}
		for _, q := range history {
			kept = kept || q.ResolvedAt == nil
		}
		effect := holds.End(hold.ID, "interrupted", ctx.now)
		effect.Values["kept"] = kept
		effect.Values["park_reason"] = options.Reason
		return journal.Event{EventType: "hold.parked", Aggregate: journal.Aggregate{ID: hold.ID, Type: "hold"}, Effects: []journal.Effect{effect}}, map[string]any{"task_id": hold.TaskID, "hold_id": hold.ID, "mode": "interrupted", "reason": options.Reason, "kept_for_worktree": kept}, nil
	})
}

func taskCategory(conn *sql.DB, taskID string) (string, error) {
	var workflowID, status string
	if err := conn.QueryRow(`SELECT workflow_id,status FROM issues WHERE id=?`, taskID).Scan(&workflowID, &status); err != nil {
		return "", fmt.Errorf("read held task: %w", err)
	}
	registry, err := workflow.Load(conn)
	if err != nil {
		return "", err
	}
	category, known := registry.CategoryForState(workflowID, status)
	if !known {
		return "", fmt.Errorf("unknown held task state %s", status)
	}
	return category, nil
}
