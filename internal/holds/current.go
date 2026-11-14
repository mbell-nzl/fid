package holds

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

const HeartbeatIntervalSeconds = 300

// Live distinguishes an ended parked hold from one still held by a session.
func (hold Hold) Live(now time.Time) bool { return hold.EndedAt == nil && !hold.Lapsed(now) }

// LastHeartbeat preserves the initial sign of life for older records.
func (hold Hold) LastHeartbeat() time.Time {
	if hold.LastHeartbeatAt.IsZero() {
		return hold.Since
	}
	return hold.LastHeartbeatAt
}

// Release ends a session or releases its parked reservation, preserving history.
func (hold Hold) Release(reason string, at time.Time) journal.Effect {
	if hold.EndedAt == nil {
		return End(hold.ID, reason, at)
	}
	return journal.Effect{Op: "update", Table: "holds", PK: map[string]any{"id": hold.ID}, Values: map[string]any{"kept": false}, Preconditions: map[string]any{"kept": true}}
}

// CurrentByID includes ended holds: stale commands must distinguish ended
// from missing without ever substituting another session's hold.
func CurrentByID(conn *sql.DB, id string) (Hold, error) {
	var hold Hold
	var ended, reason, replaces sql.NullString
	var since, expires, heartbeat string
	err := conn.QueryRow(`SELECT id, task_id, worktree_id, worktree_root, actor, since, expires_at,
 COALESCE(last_heartbeat_at, since), ended_at, park_reason, kept, replaces_hold_id FROM holds WHERE id = ?`, id).
		Scan(&hold.ID, &hold.TaskID, &hold.WorktreeID, &hold.Worktree, &hold.Actor, &since, &expires, &heartbeat, &ended, &reason, &hold.Kept, &replaces)
	if err != nil {
		return Hold{}, err
	}
	for _, stamp := range []struct {
		raw    string
		target *time.Time
	}{{since, &hold.Since}, {expires, &hold.ExpiresAt}, {heartbeat, &hold.LastHeartbeatAt}} {
		if *stamp.target, err = time.Parse(time.RFC3339, stamp.raw); err != nil {
			return Hold{}, fmt.Errorf("read hold time: %w", err)
		}
	}
	if ended.Valid {
		hold.EndedAt = &ended.String
	}
	if reason.Valid {
		hold.ParkReason = &reason.String
	}
	if replaces.Valid {
		hold.Replaces = &replaces.String
	}
	return hold, nil
}
