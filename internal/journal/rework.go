package journal

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

// ReviewReworkEffects is shared by the writer and admission. A CHANGES round
// returns exactly its task's hand-off reservation; subsequent rounds on work
// already returned to the builder record evidence without moving it again.
func ReviewReworkEffects(conn reviewrecord.Queryer, r reviewrecord.Round) ([]Effect, error) {
	var installed int
	if err := conn.QueryRow(`SELECT count(*) FROM pragma_table_info('holds') WHERE name='rework_at'`).Scan(&installed); err != nil {
		return nil, err
	}
	// Pre-rework archive profiles retain their original recording semantics.
	if installed == 0 || r.Verdict != "CHANGES" {
		return nil, nil
	}
	var kind, state string
	var revision int
	if err := conn.QueryRow(`SELECT kind,status,revision FROM issues WHERE id=?`, r.TaskID).Scan(&kind, &state, &revision); err != nil {
		return nil, err
	}
	if kind != "task" || state != "in_review" && state != "in_progress" {
		return nil, fmt.Errorf("review changes require a task in_review or in_progress")
	}
	if state == "in_progress" {
		return nil, nil
	}
	var hold string
	if err := conn.QueryRow(`SELECT h.id FROM holds h JOIN handoffs f ON f.hold_id=h.id
 WHERE f.id=? AND f.task_id=? AND h.task_id=? AND h.kept=1 AND h.end_reason='review'
 AND h.ended_at IS NOT NULL AND h.rework_at IS NULL`, r.HandoffID, r.TaskID, r.TaskID).Scan(&hold); err != nil {
		return nil, fmt.Errorf("review changes require the task's own kept hand-off reservation: %w", err)
	}
	return []Effect{
		{Op: "update", Table: "issues", PK: map[string]any{"id": r.TaskID}, Values: map[string]any{"status": "in_progress", "revision": revision + 1, "updated_at": r.At}, Preconditions: map[string]any{"status": "in_review", "revision": revision}},
		{Op: "update", Table: "holds", PK: map[string]any{"id": hold}, Values: map[string]any{"rework_at": r.At}, Preconditions: map[string]any{"kept": true, "end_reason": "review", "rework_at": nil}},
	}, nil
}

// A kept rework reservation may be taken only as linked work in its original
// worktree. Takeover and cancellation have their own validated families.
func validateReworkHold(tx *sql.Tx, event Event) error {
	if event.EventType != "hold.taken" {
		return nil
	}
	var installed int
	if err := tx.QueryRow(`SELECT count(*) FROM pragma_table_info('holds') WHERE name='rework_at'`).Scan(&installed); err != nil || installed == 0 {
		return err
	}
	var insert *Effect
	for i := range event.Effects {
		e := &event.Effects[i]
		if strings.EqualFold(e.Table, "holds") && e.Op == "insert" {
			if insert != nil {
				return fmt.Errorf("taking work requires one new hold")
			}
			insert = e
		}
	}
	for _, e := range event.Effects {
		if !strings.EqualFold(e.Table, "holds") || e.Op != "update" {
			continue
		}
		// Resolve the affected row as SQLite does, including legacy key and
		// table spelling, before the Event-v1 compatibility path can return.
		row, found, err := readCurrentValues(tx, "holds", e.PK, []string{"id"})
		if err != nil {
			return err
		}
		if !found {
			continue
		}
		var task, worktree, state string
		var rework sql.NullString
		if err := tx.QueryRow(`SELECT h.task_id,h.worktree_id,h.rework_at,i.status FROM holds h JOIN issues i ON i.id=h.task_id WHERE h.id=? AND h.kept=1`, row["id"]).Scan(&task, &worktree, &rework, &state); err == sql.ErrNoRows {
			continue
		} else if err != nil {
			return err
		}
		if !rework.Valid {
			continue
		}
		if !event.IsV2() || insert == nil || state != "in_progress" || insert.Values["task_id"] != task || insert.Values["worktree_id"] != worktree || insert.Values["replaces_hold_id"] != row["id"] {
			return fmt.Errorf("rework requires a linked hold in the task's kept worktree")
		}
		// Rework has no legacy writer. Optional historical fields cannot
		// bypass the full new-hold checks, including the root/ID binding.
		if err := validateNewHold(event, *insert); err != nil {
			return err
		}
	}
	return nil
}

// Rework metadata belongs only to the returning round. New build holds may
// explicitly record null, as older inserts do implicitly through the default.
func validReworkColumn(event Event, effect Effect, value any) bool {
	if !event.IsV2() {
		return false
	}
	switch effect.Op {
	case "update":
		return event.EventType == "review.round-recorded" || event.EventType == "task.default-conflict"
	case "insert":
		return currentEventValue(value) == nil && (event.EventType == "hold.taken" || event.EventType == "task.taken-over")
	}
	return false
}
