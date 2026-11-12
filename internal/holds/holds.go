// Package holds keeps the hold record: a task handed out by fid next --claim
// to one actor in one worktree, with its own ID (ADR-0014). A hold is open
// until something ends it. Whether an open hold has lapsed is read from the
// clock against its expiry, never stored.
package holds

import (
	"database/sql"
	"fmt"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Lifetime is how long a hold lasts before it lapses (D8). Heartbeats renew
// it from task 1.5a.
const Lifetime = time.Hour

// How a hold ended: a new hold replaced it, or its task was completed or
// cancelled in the same Event.
const (
	EndReplaced      = "replaced"
	EndTaskCompleted = "task-completed"
	EndTaskCancelled = "task-cancelled"
)

// Hold is one hold. Worktree is the worktree's root path; WorktreeID is
// sha256 of that path, so a moved worktree is another worktree. Replaces
// names the lapsed hold a resume replaced.
type Hold struct {
	ID         string
	TaskID     string
	WorktreeID string
	Worktree   string
	Actor      string
	Since      time.Time
	ExpiresAt  time.Time
	Replaces   *string
}

// Lapsed reports whether the hold had expired at now.
func (hold Hold) Lapsed(now time.Time) bool {
	return !now.Before(hold.ExpiresAt)
}

// New is a hold on taskID for actor in the worktree at root, taken at since,
// to the second, and lasting Lifetime.
func New(id, taskID, worktreeID, root, actor string, since time.Time, replaces *string) Hold {
	since = since.UTC().Truncate(time.Second)
	return Hold{ID: id, TaskID: taskID, WorktreeID: worktreeID, Worktree: root, Actor: actor,
		Since: since, ExpiresAt: since.Add(Lifetime), Replaces: replaces}
}

// Open reads every hold that hasn't ended, by ID.
func Open(conn fiddb.Queryer) ([]Hold, error) {
	rows, err := conn.Query(`SELECT id, task_id, worktree_id, worktree_root, actor, since, expires_at, replaces_hold_id
		FROM holds WHERE ended_at IS NULL ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("read open holds: %w", err)
	}
	defer func() { _ = rows.Close() }()
	open := []Hold{}
	for rows.Next() {
		var (
			hold             Hold
			since, expiresAt string
			replaces         sql.NullString
		)
		if err := rows.Scan(&hold.ID, &hold.TaskID, &hold.WorktreeID, &hold.Worktree, &hold.Actor, &since, &expiresAt, &replaces); err != nil {
			return nil, fmt.Errorf("scan open hold: %w", err)
		}
		if hold.Since, err = time.Parse(time.RFC3339, since); err != nil {
			return nil, fmt.Errorf("read hold %s: since: %w", hold.ID, err)
		}
		if hold.ExpiresAt, err = time.Parse(time.RFC3339, expiresAt); err != nil {
			return nil, fmt.Errorf("read hold %s: expires_at: %w", hold.ID, err)
		}
		if replaces.Valid {
			hold.Replaces = &replaces.String
		}
		open = append(open, hold)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate open holds: %w", err)
	}
	return open, nil
}

// NewID is an unused hold ID, such as HLD-4c1e9a.
func NewID(conn *sql.DB) (string, error) {
	return journal.GenerateID("HLD", func(candidate string) (bool, error) {
		var found int
		err := conn.QueryRow(`SELECT 1 FROM holds WHERE id = ?`, candidate).Scan(&found)
		if err == sql.ErrNoRows {
			return false, nil
		}
		return err == nil, err
	})
}

// Insert is the Effect that records the hold.
func (hold Hold) Insert() journal.Effect {
	var replaces any
	if hold.Replaces != nil {
		replaces = *hold.Replaces
	}
	return journal.Effect{Op: "insert", Table: "holds", PK: map[string]any{"id": hold.ID}, Values: map[string]any{
		"task_id": hold.TaskID, "worktree_id": hold.WorktreeID, "worktree_root": hold.Worktree, "actor": hold.Actor,
		"since": Timestamp(hold.Since), "expires_at": Timestamp(hold.ExpiresAt), "replaces_hold_id": replaces,
		"ended_at": nil, "end_reason": nil,
	}}
}

// End is the Effect that ends the open hold id at at, for reason.
func End(id, reason string, at time.Time) journal.Effect {
	return journal.Effect{Op: "update", Table: "holds", PK: map[string]any{"id": id},
		Values:        map[string]any{"ended_at": Timestamp(at), "end_reason": reason},
		Preconditions: map[string]any{"ended_at": nil}}
}

// EndForTask is the Effect that ends taskID's open hold, if it has one,
// because the task moved into category at at. It is empty unless the
// category is completed or cancelled.
func EndForTask(conn *sql.DB, taskID, category string, at time.Time) ([]journal.Effect, error) {
	reason := map[string]string{"completed": EndTaskCompleted, "cancelled": EndTaskCancelled}[category]
	if reason == "" {
		return nil, nil
	}
	var id string
	err := conn.QueryRow(`SELECT id FROM holds WHERE task_id = ? AND ended_at IS NULL`, taskID).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read the open hold on %s: %w", taskID, err)
	}
	return []journal.Effect{End(id, reason, at)}, nil
}

// Timestamp is how a hold's times are stored and shown in JSON: UTC, to
// the second.
func Timestamp(at time.Time) string {
	return at.UTC().Truncate(time.Second).Format(time.RFC3339)
}
