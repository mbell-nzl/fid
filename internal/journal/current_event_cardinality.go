package journal

import (
	"database/sql"
	"fmt"
)

func validateCurrentCardinality(event Event) error {
	counts := map[string]int{}
	for _, effect := range event.Effects {
		key := currentEffectKey(effect.Table, effect.Op)
		counts[key]++
		maximum := 1
		switch effect.Table {
		case "budget_samples", "run_budget_bindings", "referenced_artifacts", "receipt_artifacts", "activation_actor_profiles", "activation_gate_definitions":
			maximum = len(event.Effects)
		case "stops", "attention_items", "attention_stop_links", "stop_resolutions":
			if event.EventType != "stop.raised" && event.EventType != "stop.resolved" && event.EventType != "attention.resolved" {
				maximum = len(event.Effects)
			}
		case "activation_components":
			if effect.Op == "update" {
				maximum = len(event.Effects)
			}
		case "activation_resources":
			if effect.Op == "insert" {
				maximum = 5
			}
		case "kit_activations":
			if event.EventType == "activation.finalized" || event.EventType == "activation.recovery-completed" {
				maximum = 2
			}
		}
		if counts[key] > maximum {
			return fmt.Errorf("current Event %s: too many %s Effects", event.EventType, key)
		}
	}
	if event.EventType == "run.recovered" && counts["runs/update"] != counts["run_runtime/update"] {
		return fmt.Errorf("current recovery: predecessor terminal and runtime Effects must occur together")
	}
	if event.EventType == "activation.prepared" {
		names := map[any]bool{}
		for _, e := range event.Effects {
			if e.Table == "activation_resources" {
				names[e.PK["component"]] = true
			}
		}
		if counts["activation_resources/insert"] != 5 || !names["kit"] || !names["backend"] || !names["overlay"] || !names["project"] || !names["provider"] {
			return fmt.Errorf("current activation: exactly five distinct resources required")
		}
	}
	if event.EventType == "activation.recovery-selected" && counts["stops/insert"] != 1 {
		return fmt.Errorf("current activation selection: incomplete Stop required")
	}
	return nil
}

// Each automatic budget boundary contributes one ordered Stop/Attention/link
// bundle. Attention reuse omits only the insert, never the link.
func validateCurrentStopBundles(tx *sql.Tx, event Event) error {
	for index := 0; index < len(event.Effects); index++ {
		stop := event.Effects[index]
		if stop.Table != "stops" {
			if event.Effects[index].Op == "insert" && (stop.Table == "attention_items" || stop.Table == "attention_stop_links") {
				return fmt.Errorf("current Event: orphan Attention Effect")
			}
			continue
		}
		var attention map[string]any
		if index+1 < len(event.Effects) && event.Effects[index+1].Table == "attention_items" {
			index++
			attention, _ = mergeEffectValues(event.Effects[index].PK, event.Effects[index].Values)
			if attention["status"] != "open" || !currentValuesEqual(attention["version"], 1) || attention["resolution_event_id"] != nil {
				return fmt.Errorf("current Event: newly raised Attention must be open")
			}
		}
		index++
		if index >= len(event.Effects) || event.Effects[index].Table != "attention_stop_links" || event.Effects[index].PK["stop_id"] != stop.PK["id"] {
			return fmt.Errorf("current Event: Stop requires its ordered Attention link")
		}
		link := event.Effects[index]
		if attention == nil {
			row, found, err := readCurrentEventValues(tx, "attention_items", map[string]any{"id": link.PK["attention_id"]}, []string{"id", "status", "subject_type", "subject_id", "code"})
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("current Event: missing Attention reference")
			}
			attention = row
		}
		if attention["id"] != link.PK["attention_id"] || attention["status"] != "open" || attention["subject_type"] != stop.Values["subject_type"] || attention["subject_id"] != stop.Values["subject_id"] || attention["code"] != stop.Values["code"] {
			return fmt.Errorf("current Event: Stop/Attention binding mismatch")
		}
	}
	return nil
}
