package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// ValidateIssueCreateResult checks optional creation lanes at both retry and
// replay boundaries. Results written before lanes existed may omit the field.
func ValidateIssueCreateResult(raw string) error {
	var result struct {
		Issue struct {
			Kind string
			Lane any
		}
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return fmt.Errorf("decode stored Issue create result: %w", err)
	}
	if result.Issue.Lane != nil && (result.Issue.Kind != "task" || !validTaskLane(result.Issue.Lane)) {
		return fmt.Errorf("invalid Issue create result: lane must be routine, standard or assured and belong to a task")
	}
	return nil
}

func validTaskLane(value any) bool {
	lane, ok := value.(string)
	return ok && (lane == "routine" || lane == "standard" || lane == "assured")
}

// ValidateIssueUpdateResult checks lane values before a stored update result
// can be rendered or accepted into the projection by replay or archive import.
func ValidateIssueUpdateResult(raw string) error {
	var result struct {
		Changes []struct {
			Field string
			From  any
			To    any
		} `json:"changes"`
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return fmt.Errorf("decode stored Issue update result: %w", err)
	}
	for _, change := range result.Changes {
		if change.Field != "lane" {
			continue
		}
		if !validTaskLane(change.To) {
			return fmt.Errorf("invalid lane change result: to must be routine, standard or assured")
		}
		if change.From != nil && !validTaskLane(change.From) {
			return fmt.Errorf("invalid lane change result: from must be null, routine, standard or assured")
		}
	}
	return nil
}

// Lane changes use the existing issue Events, preserving their revision and
// idempotency contracts. Replay enforces the same authority as the writer.
func validateTaskLane(tx *sql.Tx, event Event) error {
	if event.EventType == "issue.created" || event.EventType == "issue.updated" {
		for _, effect := range event.Effects {
			if effect.Table != "idempotency_results" {
				continue
			}
			raw, ok := effect.Values["result_json"].(string)
			if !ok {
				return fmt.Errorf("invalid Issue mutation result")
			}
			validate := ValidateIssueUpdateResult
			if event.EventType == "issue.created" {
				validate = ValidateIssueCreateResult
			}
			if err := validate(raw); err != nil {
				return err
			}
		}
	}
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
