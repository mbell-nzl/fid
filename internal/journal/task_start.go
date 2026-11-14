package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func validateTaskStart(tx *sql.Tx, event Event) error {
	if event.EventType != "task.started" {
		return nil
	}
	if event.CommandName != "task.start" || event.Repository != nil || event.ActorV2.ActorType != "agent" && event.ActorV2.ActorType != "human" {
		return fmt.Errorf("task start requires a work actor and its command, without a repository base")
	}
	if err := validateCurrentIdempotency(event); err != nil {
		return err
	}
	var stored struct {
		Fields     map[string]any `json:"fields"`
		OccurredAt string         `json:"occurred_at"`
	}
	last := event.Effects[len(event.Effects)-1]
	raw, ok := last.Values["result_json"].(string)
	if !ok || json.Unmarshal([]byte(raw), &stored) != nil {
		return fmt.Errorf("invalid task start result")
	}
	var parent, lane sql.NullString
	var kind, status string
	if err := tx.QueryRow(`SELECT kind,status,lane,parent_id FROM issues WHERE id=?`, event.Aggregate.ID).Scan(&kind, &status, &lane, &parent); err != nil {
		return fmt.Errorf("read starting task: %w", err)
	}
	if kind != "task" || status != "todo" || !lane.Valid || len(stored.Fields) != 4 || stored.Fields["task_id"] != event.Aggregate.ID || stored.Fields["status"] != "in_progress" || stored.Fields["lane"] != lane.String || stored.OccurredAt != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("task start result differs from planned task")
	}
	var task, actor, since, expires string
	var ended sql.NullString
	if err := tx.QueryRow(`SELECT task_id,actor,since,expires_at,ended_at FROM holds WHERE id=?`, stored.Fields["hold_id"]).Scan(&task, &actor, &since, &expires, &ended); err != nil {
		return fmt.Errorf("read starting hold: %w", err)
	}
	stamp := event.OccurredAt.UTC().Format(time.RFC3339)
	if task != event.Aggregate.ID || actor != event.ActorV2.ActorType+":"+event.ActorV2.ActorID || ended.Valid || stamp < since || stamp >= expires {
		return fmt.Errorf("task start requires its actor's current hold")
	}
	ids := []string{event.Aggregate.ID}
	if parent.Valid {
		var parentState string
		if err := tx.QueryRow(`SELECT status FROM issues WHERE id=? AND kind='epic'`, parent.String).Scan(&parentState); err != nil {
			return err
		}
		if parentState == "todo" {
			ids = append(ids, parent.String)
		}
	}
	if len(event.Effects) != len(ids)+1 {
		return fmt.Errorf("task start must move exactly its task and unstarted epic")
	}
	for i, id := range ids {
		effect := event.Effects[i]
		var revision int
		if err := tx.QueryRow(`SELECT revision FROM issues WHERE id=?`, id).Scan(&revision); err != nil {
			return err
		}
		if effect.Table != "issues" || effect.Op != "update" || !currentKeys(effect.PK, "id") || effect.PK["id"] != id || !currentKeys(effect.Values, "status revision updated_at") || !currentKeys(effect.Preconditions, "status revision") || len(effect.Prior) != 0 || effect.Values["status"] != "in_progress" || effect.Values["updated_at"] != stamp || !currentValuesEqual(effect.Values["revision"], revision+1) || effect.Preconditions["status"] != "todo" || !currentValuesEqual(effect.Preconditions["revision"], revision) {
			return fmt.Errorf("invalid task or epic start Effect")
		}
	}
	return nil
}
