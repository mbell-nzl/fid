package journal

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/workrecord"
)

func validateWorkDecision(event Event, conn *sql.Tx) error {
	allowance := event.EventType == "rework.allowed"
	command, kind, keys := "task.accept", "accept", "kind actor subject_type subject_id reason binding_id decided_at"
	if allowance {
		command, kind, keys = "task.rework-allow", "rework-allowance", keys+" binding_json"
	}
	if event.CommandName != command || event.Repository != nil || event.Idempotency != nil || len(event.Effects) != 1 {
		return fmt.Errorf("work decision requires exactly one decision and its command")
	}
	e := event.Effects[0]
	if err := validateCurrentEffect(event, e); err != nil {
		return err
	}
	if !currentKeys(e.Values, keys) {
		return fmt.Errorf("work decision requires a complete immutable record")
	}
	d := e.Values
	if d["kind"] != kind || d["actor"] != "human:"+event.ActorV2.ActorID || d["subject_id"] != event.Aggregate.ID || d["binding_id"] != event.Aggregate.ID || d["decided_at"] != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("work decision differs from Event")
	}
	reason, ok := d["reason"].(string)
	if !ok || strings.TrimSpace(reason) == "" {
		return fmt.Errorf("work decision requires a reason")
	}
	var subjectKind, category string
	if err := conn.QueryRow(`SELECT i.kind,s.category FROM issues i JOIN workflow_states s ON s.workflow_id=i.workflow_id AND s.state=i.status WHERE i.id=?`, event.Aggregate.ID).Scan(&subjectKind, &category); err != nil {
		return err
	}
	if (subjectKind != "task" && (allowance || subjectKind != "subtask")) || category == "completed" || category == "cancelled" || d["subject_type"] != subjectKind {
		return fmt.Errorf("work decision requires an open task (or subtask for acceptance)")
	}
	if allowance {
		// The shared decision column is nullable for other kinds.
		raw, _ := d["binding_json"].(string)
		_, err := workrecord.AllowanceRounds(raw)
		return err
	}
	var n int
	if err := conn.QueryRow(`SELECT count(*) FROM person_decisions WHERE kind='accept' AND subject_id=?`, event.Aggregate.ID).Scan(&n); err != nil {
		return err
	}
	if n != 0 {
		return fmt.Errorf("work already accepted")
	}
	return nil
}

// The creation actor belongs to the Event, not its user-editable activity.
// Rebuilding the projection backfills attribution for every existing task.
func projectIssueCreator(tx *sql.Tx, event Event) error {
	if !event.IsV2() || event.EventType != "issue.created" {
		return nil
	}
	var installed int
	if err := tx.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='issue_creators'`).Scan(&installed); err != nil {
		return err
	}
	if installed == 0 {
		return nil
	}
	_, err := tx.Exec(`INSERT INTO issue_creators(issue_id,actor) VALUES(?,?)`, event.Aggregate.ID, event.ActorV2.ActorType+":"+event.ActorV2.ActorID)
	return err
}
