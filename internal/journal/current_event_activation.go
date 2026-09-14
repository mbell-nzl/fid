package journal

import (
	"database/sql"
	"fmt"
	"strings"
)

func validateCurrentActivationReferences(tx *sql.Tx, event Event) error {
	if !strings.HasPrefix(event.EventType, "activation.") {
		return nil
	}
	if err := validateCurrentComponentTransitions(tx, event); err != nil {
		return err
	}
	var primary map[string]any
	for _, effect := range event.Effects {
		if effect.Table == "kit_activations" && effect.PK["id"] == event.Aggregate.ID {
			primary, _ = mergeEffectValues(effect.PK, effect.Values)
		}
	}
	if primary == nil {
		return fmt.Errorf("current activation: missing primary activation Effect")
	}
	if err := validateCurrentActivationAgreement(tx, event, primary); err != nil {
		return err
	}
	var prior map[string]any
	if event.EventType != "activation.prepared" {
		var err error
		prior, err = currentReference(tx, "kit_activations", map[string]any{"id": event.Aggregate.ID})
		if err != nil {
			return err
		}
	}
	for _, effect := range event.Effects {
		if effect.Table == "kit_activations" && effect.PK["id"] != event.Aggregate.ID {
			if prior == nil || effect.PK["id"] != prior["prior_activation_id"] || effect.Values["status"] != "superseded" {
				return fmt.Errorf("current activation: unrelated predecessor Effect")
			}
		}
		if effect.Table == "activation_recoveries" && effect.Op == "insert" && effect.Values["selected_event_id"] != event.EventID {
			return fmt.Errorf("current activation: selection Event mismatch")
		}
		if effect.Table == "stops" && (effect.Values["subject_type"] != "activation" || effect.Values["subject_id"] != event.Aggregate.ID || effect.Values["code"] != "activation-incomplete") {
			return fmt.Errorf("current activation: incomplete Stop binding mismatch")
		}
		if effect.Table == "stop_resolutions" {
			stop, err := currentReference(tx, "stops", map[string]any{"id": effect.PK["stop_id"]})
			if err != nil {
				return err
			}
			if stop["subject_type"] != "activation" || stop["subject_id"] != event.Aggregate.ID || stop["code"] != "activation-incomplete" || effect.Values["fact_event_id"] != event.EventID {
				return fmt.Errorf("current activation: Stop disposition binding mismatch")
			}
		}
	}
	if event.EventType == "activation.recovery-completed" {
		selection, err := currentReference(tx, "activation_recoveries", map[string]any{"activation_id": event.Aggregate.ID})
		if err != nil {
			return err
		}
		status := "aborted"
		if selection["direction"] == "complete-forward" {
			status = "finalized"
		}
		if primary["status"] != status {
			return fmt.Errorf("current activation: completion direction mismatch")
		}
		if err := validateCurrentActivationStopDispositions(tx, event); err != nil {
			return err
		}
	}
	if event.EventType == "activation.recovery-selected" {
		selection := currentInserted(event, "activation_recoveries")
		direction := "restore-prior"
		if prior["point_of_no_return"] == true || currentValuesEqual(prior["point_of_no_return"], 1) {
			direction = "complete-forward"
		}
		if selection["direction"] != direction {
			return fmt.Errorf("current activation: recovery selection direction does not match point of no return")
		}
	}
	return nil
}

func validateCurrentActivationAgreement(tx *sql.Tx, event Event, primary map[string]any) error {
	state := "new"
	switch event.EventType {
	case "activation.committed", "activation.finalized":
	case "activation.recovery-completed":
		if primary["status"] == "aborted" {
			state = "old"
		}
	default:
		return nil
	}
	rows, err := tx.Query(`SELECT state,old_identity,new_identity,observed_identity FROM activation_resources WHERE activation_id=?`, event.Aggregate.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	count := 0
	for rows.Next() {
		var actual, newID string
		var old, observed sql.NullString
		if err := rows.Scan(&actual, &old, &newID, &observed); err != nil {
			return err
		}
		if actual != state || (state == "old" && old != observed) || (state == "new" && (!observed.Valid || observed.String != newID)) {
			return fmt.Errorf("current activation: recorded resources do not agree with transition")
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if count != 5 {
		return fmt.Errorf("current activation: incomplete recorded resource set")
	}
	return nil
}

func validateCurrentActivationStopDispositions(tx *sql.Tx, event Event) error {
	rows, err := tx.Query(`SELECT s.id FROM stops s LEFT JOIN stop_resolutions r ON r.stop_id = s.id WHERE s.subject_type = 'activation' AND s.subject_id = ? AND s.code = 'activation-incomplete' AND r.stop_id IS NULL`, event.Aggregate.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	required := map[string]bool{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		required[id] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, effect := range event.Effects {
		if effect.Table != "stop_resolutions" {
			continue
		}
		id, _ := effect.PK["stop_id"].(string)
		if !required[id] {
			return fmt.Errorf("current activation: unexpected Stop disposition")
		}
		delete(required, id)
	}
	if len(required) != 0 {
		return fmt.Errorf("current activation: incomplete Stop disposition")
	}
	return nil
}
