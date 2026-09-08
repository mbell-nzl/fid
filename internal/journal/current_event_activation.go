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
	var primary map[string]any
	for _, effect := range event.Effects {
		if effect.Table == "kit_activations" && effect.PK["id"] == event.Aggregate.ID {
			primary, _ = mergeEffectValues(effect.PK, effect.Values)
		}
	}
	if primary == nil {
		return fmt.Errorf("current activation: missing primary activation Effect")
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
	}
	return nil
}
