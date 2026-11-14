package journal

import (
	"database/sql"
	"fmt"
)

// Lane changes use the existing issue Events, preserving their revision and
// idempotency contracts. Replay enforces the same authority as the writer.
func validateTaskLane(tx *sql.Tx, event Event) error {
	for _, effect := range event.Effects {
		value, present := effect.Values["lane"]
		if effect.Table != "issues" || !present {
			continue
		}
		lane, ok := value.(string)
		rank := map[string]int{"routine": 1, "standard": 2, "assured": 3}
		if !ok || rank[lane] == 0 {
			return fmt.Errorf("invalid task lane")
		}
		if event.ActorV2.ActorType != "agent" && event.ActorV2.ActorType != "human" {
			return fmt.Errorf("task lane requires a work actor")
		}
		if effect.Op == "insert" {
			if event.EventType != "issue.created" || effect.Values["kind"] != "task" {
				return fmt.Errorf("lane belongs to a task creation")
			}
			continue
		}
		if effect.Op != "update" || event.EventType != "issue.updated" {
			return fmt.Errorf("lane belongs to an issue update")
		}
		var kind, status string
		var prior sql.NullString
		if err := tx.QueryRow(`SELECT kind,status,lane FROM issues WHERE id=?`, effect.PK["id"]).Scan(&kind, &status, &prior); err != nil {
			return fmt.Errorf("read lane subject: %w", err)
		}
		if kind != "task" {
			return fmt.Errorf("lane belongs to tasks only")
		}
		if rank[lane] < rank[prior.String] && (event.ActorV2.ActorType != "human" || status != "todo") {
			return fmt.Errorf("lane lowering requires a person and a todo task")
		}
	}
	return nil
}
