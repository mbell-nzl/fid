package holdwork

import (
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func (ctx writeContext) reviewHeartbeat() (journal.Event, map[string]any, error) {
	claims, err := reviewrecord.Claims(ctx.conn)
	if err != nil {
		return journal.Event{}, nil, err
	}
	var claim reviewrecord.Claim
	reason := "not-found"
	stamp := holds.Timestamp(ctx.now)
	for _, c := range claims {
		if c.ID != ctx.options.HoldID {
			continue
		}
		claim = c
		switch {
		case c.EndedAt != nil:
			reason = "ended"
		case c.ExpiresAt <= stamp:
			reason = "expired"
		case c.Actor != ctx.options.Actor:
			reason = "actor-mismatch"
		case c.WorktreeID != project.WorktreeID(ctx.options.WorkspaceRoot):
			reason = "worktree-mismatch"
		default:
			reason = ""
		}
	}
	if reason != "" {
		return journal.Event{}, nil, &Refusal{Code: "hold-not-current", TaskID: claim.TaskID, Message: "this review session is no longer current", Details: map[string]any{"hold_id": ctx.options.HoldID, "reason": reason}}
	}
	expires := holds.Timestamp(ctx.now.Add(time.Hour))
	effect := journal.Effect{Op: "update", Table: "review_claims", PK: map[string]any{"id": claim.ID}, Values: map[string]any{"last_heartbeat_at": stamp, "expires_at": expires}, Preconditions: map[string]any{"ended_at": nil, "expires_at": claim.ExpiresAt}}
	return journal.Event{EventType: "review.heartbeat", Aggregate: journal.Aggregate{ID: claim.ID, Type: "review-claim"}, Effects: []journal.Effect{effect}}, map[string]any{"task_id": claim.TaskID, "hold_id": claim.ID, "last_heartbeat_at": stamp, "expires_at": expires, "heartbeat_interval_seconds": holds.HeartbeatIntervalSeconds}, nil
}
