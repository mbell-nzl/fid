package journal

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// All import and replay paths reach this validator before applying Effects.
// Legacy hold.taken records remain valid without the new optional columns.
func validateHoldLifecycle(tx *sql.Tx, event Event) error {
	var installed int
	if err := tx.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='person_decisions' AND type='table'`).Scan(&installed); err != nil {
		return err
	}
	if installed == 0 {
		return nil
	}
	commands := map[string]string{"hold.heartbeat": "hold.heartbeat", "hold.parked": "hold.park", "task.taken-over": "task.takeover"}
	if command := commands[event.EventType]; command != "" && command != event.CommandName {
		return fmt.Errorf("hold lifecycle command differs from family")
	}
	decisions, inserts, updates := 0, 0, 0
	for _, effect := range event.Effects {
		if effect.Table == "idempotency_results" && commands[event.EventType] != "" {
			if err := validateCurrentEffect(event, effect); err != nil {
				return err
			}
		}
		if effect.Table == "person_decisions" {
			decisions++
			if err := validatePersonDecision(tx, event, effect); err != nil {
				return err
			}
		}
		if effect.Table != "holds" {
			continue
		}
		switch effect.Op {
		case "insert":
			inserts++
			if _, modern := effect.Values["last_heartbeat_at"]; modern || event.EventType == "task.taken-over" {
				if err := validateNewHold(event, effect); err != nil {
					return err
				}
			}
			if event.EventType == "hold.taken" && effect.PK["id"] != event.Aggregate.ID {
				return fmt.Errorf("hold identity differs from aggregate")
			}
			if event.EventType == "task.taken-over" && effect.Values["task_id"] != event.Aggregate.ID {
				return fmt.Errorf("takeover hold differs from task")
			}
		case "update":
			updates++
			if err := validateHoldUpdate(tx, event, effect); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported hold Effect")
		}
	}
	if commands[event.EventType] != "" && event.Idempotency != nil {
		if err := validateCurrentIdempotency(event); err != nil {
			return err
		}
	}
	switch event.EventType {
	case "hold.heartbeat":
		if updates != 1 || len(event.Effects) != 1 || inserts != 0 || event.Idempotency != nil {
			return fmt.Errorf("heartbeat requires exactly one hold update")
		}
	case "hold.parked":
		if updates != 1 || inserts != 0 || event.Idempotency == nil {
			return fmt.Errorf("park requires one ended hold and retry key")
		}
	case "task.taken-over":
		if decisions != 1 || inserts != 1 {
			return fmt.Errorf("takeover requires one decision and one new hold")
		}
		return validateTakeover(tx, event)
	}
	return nil
}

func validatePersonDecision(tx *sql.Tx, event Event, effect Effect) error {
	if err := validateCurrentEffect(event, effect); err != nil {
		return err
	}
	if effect.Op != "insert" || len(effect.PK) != 1 || len(effect.Values) != 7 {
		return fmt.Errorf("person's decision requires a complete immutable record")
	}
	kind := ""
	switch event.EventType {
	case "task.taken-over":
		kind = "takeover"
	case "issue.transitioned":
		kind = "cancel"
	default:
		return fmt.Errorf("unrelated person's decision Effect")
	}
	row := effect.Values
	if event.ActorV2.ActorType != "human" || row["actor"] != "human:"+event.ActorV2.ActorID || row["kind"] != kind || row["subject_id"] != event.Aggregate.ID || row["binding_id"] != event.Aggregate.ID || row["decided_at"] != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("person's decision actor, subject or time differs from Event")
	}
	if reason, ok := row["reason"].(string); !ok || strings.TrimSpace(reason) == "" {
		return fmt.Errorf("person's decision requires a reason")
	}
	var subjectKind string
	if err := tx.QueryRow(`SELECT kind FROM issues WHERE id=?`, event.Aggregate.ID).Scan(&subjectKind); err != nil {
		return err
	}
	if row["subject_type"] != subjectKind || (subjectKind != "task" && subjectKind != "epic") {
		return fmt.Errorf("person's decision has wrong subject kind")
	}
	if kind == "cancel" {
		cancelled := false
		for _, change := range event.Effects {
			if change.Table == "issues" {
				var category string
				err := tx.QueryRow(`SELECT s.category FROM issues i JOIN workflow_states s ON s.workflow_id=i.workflow_id AND s.state=? WHERE i.id=?`, change.Values["status"], event.Aggregate.ID).Scan(&category)
				if err == nil && category == "cancelled" {
					cancelled = true
				}
			}
		}
		if !cancelled {
			return fmt.Errorf("cancel decision requires cancellation")
		}
	}
	return nil
}

func validateHoldUpdate(tx *sql.Tx, event Event, effect Effect) error {
	var actor, task, expiry string
	var ended sql.NullString
	var kept int
	err := tx.QueryRow(`SELECT actor,task_id,expires_at,ended_at,kept FROM holds WHERE id=?`, effect.PK["id"]).Scan(&actor, &task, &expiry, &ended, &kept)
	if err != nil {
		return fmt.Errorf("read changed hold: %w", err)
	}
	if event.EventType == "hold.heartbeat" || event.EventType == "hold.parked" {
		if effect.PK["id"] != event.Aggregate.ID || ended.Valid || actor != event.ActorV2.ActorType+":"+event.ActorV2.ActorID {
			return fmt.Errorf("hold session is not current")
		}
		expires, err := time.Parse(time.RFC3339, expiry)
		if err != nil {
			return err
		}
		if !event.OccurredAt.Before(expires) {
			return fmt.Errorf("hold session expired")
		}
		if event.EventType == "hold.heartbeat" {
			if len(effect.Values) != 2 || effect.Values["last_heartbeat_at"] != event.OccurredAt.UTC().Format(time.RFC3339) || effect.Values["expires_at"] != event.OccurredAt.Add(time.Hour).UTC().Format(time.RFC3339) {
				return fmt.Errorf("invalid heartbeat renewal")
			}
			return nil
		}
		var category string
		if err := tx.QueryRow(`SELECT s.category FROM issues i JOIN workflow_states s ON s.workflow_id=i.workflow_id AND s.state=i.status WHERE i.id=?`, task).Scan(&category); err != nil {
			return err
		}
		if len(effect.Values) != 4 || effect.Values["end_reason"] != "interrupted" || !currentValuesEqual(effect.Values["kept"], category == "started") {
			return fmt.Errorf("invalid interrupted park")
		}
		if reason, ok := effect.Values["park_reason"].(string); !ok || strings.TrimSpace(reason) == "" {
			return fmt.Errorf("park requires a reason")
		}
	} else if ended.Valid {
		if event.EventType == "issue.transitioned" && task != event.Aggregate.ID {
			return fmt.Errorf("transition releases another task reservation")
		}
		if kept != 1 || len(effect.Values) != 1 || !currentValuesEqual(effect.Values["kept"], false) {
			return fmt.Errorf("only a kept reservation can be released")
		}
		return nil
	} else {
		if len(effect.Values) != 2 {
			return fmt.Errorf("hold end changes immutable fields")
		}
		reason, ok := effect.Values["end_reason"].(string)
		if !ok {
			return fmt.Errorf("hold end requires a reason")
		}
		allowed := map[string]string{"hold.taken": "replaced", "task.taken-over": "takeover"}
		if event.EventType == "task.taken-over" && task != event.Aggregate.ID {
			allowed[event.EventType] = "replaced"
		}
		if event.EventType == "issue.transitioned" {
			if task != event.Aggregate.ID {
				return fmt.Errorf("transition ends another task's hold")
			}
			if reason != "task-completed" && reason != "task-cancelled" {
				return fmt.Errorf("invalid task hold end")
			}
		} else if reason != allowed[event.EventType] {
			return fmt.Errorf("invalid hold end reason")
		}
	}
	if effect.Values["ended_at"] != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("hold end differs from Event time")
	}
	return nil
}

func validateNewHold(event Event, effect Effect) error {
	if err := validateCurrentEffect(event, effect); err != nil {
		return err
	}
	row := effect.Values
	now := event.OccurredAt.UTC().Format(time.RFC3339)
	if row["actor"] != event.ActorV2.ActorType+":"+event.ActorV2.ActorID || row["since"] != now || row["last_heartbeat_at"] != now || row["expires_at"] != event.OccurredAt.Add(time.Hour).UTC().Format(time.RFC3339) || row["ended_at"] != nil || row["end_reason"] != nil || row["park_reason"] != nil || !currentValuesEqual(row["kept"], false) {
		return fmt.Errorf("new hold differs from Event actor or initial times")
	}
	root, ok := row["worktree_root"].(string)
	if !ok || root == "" || row["worktree_id"] != fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(root))) {
		return fmt.Errorf("new hold worktree identity differs from path")
	}
	return nil
}

func validateTakeover(tx *sql.Tx, event Event) error {
	var kind, category string
	if err := tx.QueryRow(`SELECT i.kind,s.category FROM issues i JOIN workflow_states s ON s.workflow_id=i.workflow_id AND s.state=i.status WHERE i.id=?`, event.Aggregate.ID).Scan(&kind, &category); err != nil {
		return err
	}
	if kind != "task" || category == "completed" || category == "cancelled" {
		return fmt.Errorf("takeover requires an open task")
	}
	var previous string
	err := tx.QueryRow(`SELECT id FROM holds WHERE task_id=? AND (ended_at IS NULL OR kept=1)`, event.Aggregate.ID).Scan(&previous)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, effect := range event.Effects {
		if effect.Table == "holds" && effect.Op == "insert" {
			replaces := currentEventValue(effect.Values["replaces_hold_id"])
			if previous == "" && replaces != nil || previous != "" && replaces != previous {
				return fmt.Errorf("takeover must link to the task's previous hold")
			}
		}
	}
	return nil
}
