package journal

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
)

func currentValueDigest(value any) (string, error) {
	raw, err := MarshalCanonicalJSON(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("sha256:%x", sha256.Sum256(raw)), nil
}

func validateCurrentComponentValues(event Event) error {
	if event.EventType != "activation.component-staged" {
		return nil
	}
	component := currentInserted(event, "activation_components")
	value, err := decodeCurrentJSON(component["evidence_json"].(string))
	if err != nil {
		return err
	}
	digest, err := currentValueDigest(value)
	if err != nil || digest != component["content_digest"] {
		return fmt.Errorf("current component: evidence digest mismatch")
	}
	values := value.([]any)
	if len(values) == 0 {
		return fmt.Errorf("current component: empty evidence")
	}
	table, kind, key := "activation_actor_profiles", "ActorProfile", "actor_id"
	if component["component_type"] == "gate-definitions" {
		table, kind, key = "activation_gate_definitions", "GateDefinition", "definition_id"
	}
	index := 0
	for _, effect := range event.Effects {
		if effect.Table != table {
			continue
		}
		if index >= len(values) || effect.PK["activation_id"] != event.Aggregate.ID {
			return fmt.Errorf("current component: unexpected derived record")
		}
		recorded := values[index].(map[string]any)
		if effect.PK[key] != recorded[key] {
			return fmt.Errorf("current component: derived identity mismatch")
		}
		if kind == "ActorProfile" {
			profile := map[string]any{"actor_id": effect.PK["actor_id"]}
			for _, name := range []string{"actor_type", "producer", "provenance_ref", "session_ref"} {
				profile[name] = currentEventValue(effect.Values[name])
			}
			profile["capabilities"], err = decodeCurrentJSON(effect.Values["capabilities_json"].(string))
			if err != nil || !currentValuesEqual(profile, recorded) {
				return fmt.Errorf("current component: profile differs from evidence")
			}
		} else {
			definition, decodeErr := decodeCurrentJSON(effect.Values["definition_json"].(string))
			definitionDigest, hashErr := currentValueDigest(definition)
			if decodeErr != nil || hashErr != nil || !currentValuesEqual(definition, recorded) || !currentValuesEqual(effect.PK["definition_version"], recorded["definition_version"]) || effect.Values["definition_digest"] != definitionDigest {
				return fmt.Errorf("current component: definition differs from evidence")
			}
		}
		index++
	}
	if index != len(values) {
		return fmt.Errorf("current component: incomplete derived records")
	}
	return nil
}

func validateCurrentComponentTransitions(tx *sql.Tx, event Event) error {
	switch event.EventType {
	case "activation.committed", "activation.finalized", "activation.aborted", "activation.recovery-completed":
	default:
		return nil
	}
	var expected int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM activation_components WHERE activation_id=?`, event.Aggregate.ID).Scan(&expected); err != nil {
		return err
	}
	actual := 0
	seen := map[string]bool{}
	for _, effect := range event.Effects {
		if effect.Table == "activation_components" {
			key, err := MarshalCanonicalJSON(effect.PK)
			if err != nil || seen[string(key)] {
				return fmt.Errorf("current activation: duplicate component transition")
			}
			seen[string(key)] = true
			if effect.PK["activation_id"] != event.Aggregate.ID {
				return fmt.Errorf("current activation: unrelated component transition")
			}
			actual++
		}
	}
	if actual != expected {
		return fmt.Errorf("current activation: incomplete component transitions")
	}
	return nil
}
