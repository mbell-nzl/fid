package journal

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func validateQuestions(tx *sql.Tx, event Event) error {
	inserts, updates := 0, 0
	for _, effect := range event.Effects {
		if effect.Table != "questions" {
			continue
		}
		if err := validateCurrentEffect(event, effect); err != nil {
			return err
		}
		if effect.Op == "insert" {
			inserts++
			if err := validateQuestionInsert(tx, event, effect); err != nil {
				return err
			}
		} else {
			updates++
			if err := validateQuestionResolution(tx, event, effect); err != nil {
				return err
			}
		}
	}
	switch event.EventType {
	case "hold.waiting":
		if inserts != 1 || updates != 0 {
			return fmt.Errorf("waiting park requires one new question")
		}
	case "task.answered":
		if inserts != 0 || updates != 1 {
			return fmt.Errorf("answer requires one question resolution")
		}
	default:
		if inserts != 0 || updates > 1 || (updates == 1 && event.EventType != "issue.transitioned") {
			return fmt.Errorf("unrelated question Effect")
		}
	}
	// A cancellation must resolve the open question, even when the task's
	// waiting reservation was replaced by a person's live takeover hold.
	if event.EventType == "issue.transitioned" {
		cancelled := false
		for _, effect := range event.Effects {
			if effect.Table == "person_decisions" && effect.Values["kind"] == "cancel" {
				cancelled = true
			}
		}
		if cancelled {
			var installed int
			if err := tx.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='questions'`).Scan(&installed); err != nil {
				return err
			}
			if installed == 0 {
				return nil
			}
			var open int
			if err := tx.QueryRow(`SELECT count(*) FROM questions WHERE task_id=? AND resolved_at IS NULL`, event.Aggregate.ID).Scan(&open); err != nil {
				return err
			}
			if updates != open {
				return fmt.Errorf("cancellation must resolve its open question")
			}
		}
	}
	return nil
}

func validateQuestionInsert(tx *sql.Tx, event Event, effect Effect) error {
	row := effect.Values
	if event.EventType != "hold.waiting" || effect.Op != "insert" || len(row) != 8 || len(effect.PK) != 1 {
		return fmt.Errorf("question requires a complete waiting record")
	}
	var task, actor string
	if err := tx.QueryRow(`SELECT task_id,actor FROM holds WHERE id=?`, event.Aggregate.ID).Scan(&task, &actor); err != nil {
		return err
	}
	if row["task_id"] != task || row["hold_id"] != event.Aggregate.ID || row["asker"] != actor || row["asker"] != event.ActorV2.ActorType+":"+event.ActorV2.ActorID || row["asked_at"] != event.OccurredAt.UTC().Format(time.RFC3339) || row["resolved_at"] != nil || row["resolution"] != nil || row["decision_id"] != nil {
		return fmt.Errorf("question differs from waiting hold or Event")
	}
	if text, ok := row["text"].(string); !ok || strings.TrimSpace(text) == "" {
		return fmt.Errorf("question requires text")
	}
	return nil
}

func validateQuestionResolution(tx *sql.Tx, event Event, effect Effect) error {
	if effect.Op != "update" || !currentKeys(effect.Values, "resolved_at resolution decision_id") || !currentKeys(effect.Preconditions, "resolved_at") || effect.Preconditions["resolved_at"] != nil {
		return fmt.Errorf("question resolution changes immutable fields")
	}
	var task string
	var resolved *string
	if err := tx.QueryRow(`SELECT task_id,resolved_at FROM questions WHERE id=?`, effect.PK["id"]).Scan(&task, &resolved); err != nil {
		return err
	}
	if resolved != nil || task != event.Aggregate.ID || event.ActorV2.ActorType != "human" || effect.Values["resolved_at"] != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("question resolution differs from Event or already resolved")
	}
	kind, resolution, binding := "answer", "answered", effect.PK["id"]
	if event.EventType == "issue.transitioned" {
		kind, resolution, binding = "cancel", "cancelled", any(task)
	} else if event.EventType != "task.answered" {
		return fmt.Errorf("unrelated question resolution")
	}
	found := false
	for _, decision := range event.Effects {
		if decision.Table != "person_decisions" || decision.Values["kind"] != kind {
			continue
		}
		if decision.Values["subject_id"] != task || decision.Values["binding_id"] != binding {
			return fmt.Errorf("decision binds another question or task")
		}
		if kind == "answer" && effect.Values["decision_id"] != decision.PK["id"] {
			return fmt.Errorf("question answer differs from decision")
		}
		found = true
	}
	if !found || effect.Values["resolution"] != resolution || (kind == "cancel" && effect.Values["decision_id"] != nil) {
		return fmt.Errorf("question requires its person's decision")
	}
	return nil
}
