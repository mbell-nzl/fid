package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Resolve the immutable bindings, not contemporary authorization policy. Replay
// must reproduce accepted facts even after their live authority has changed.
func currentReference(tx *sql.Tx, table string, pk map[string]any) (map[string]any, error) {
	definition := currentEventTables[table]
	row, found, err := readCurrentEventValues(tx, table, pk, sortedMapKeysString(definition.columns))
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("current Event: missing reference to %s %v", table, pk)
	}
	return row, nil
}
func sortedMapKeysString(values map[string]string) []string {
	result := make([]string, 0, len(values))
	for name := range values {
		result = append(result, name)
	}
	return result
}
func currentInserted(event Event, table string) map[string]any {
	for _, effect := range event.Effects {
		if effect.Table == table && effect.Op == "insert" {
			row, _ := mergeEffectValues(effect.PK, effect.Values)
			return row
		}
	}
	return nil
}
func currentMatching(left, right map[string]any, names ...string) bool {
	for _, name := range names {
		if !currentValuesEqual(left[name], right[name]) {
			return false
		}
	}
	return true
}
func validateCurrentReferences(tx *sql.Tx, event Event) error {
	if event.EventType == "run.recovered" {
		if err := validateCurrentRecoveryReferences(tx, event); err != nil {
			return err
		}
	}
	if err := validateCurrentActivationReferences(tx, event); err != nil {
		return err
	}
	runID := event.Aggregate.ID
	if event.EventType == "gate.execution-started" {
		runID, _ = currentInserted(event, "execution_attempts")["run_id"].(string)
	}
	if event.EventType == "gate.receipt-recorded" {
		receipt := currentInserted(event, "execution_receipts")
		attempt, err := currentReference(tx, "execution_attempts", map[string]any{"id": receipt["attempt_id"]})
		if err != nil {
			return err
		}
		runID, _ = attempt["run_id"].(string)
		if !currentValuesEqual(receipt["started_at"], attempt["started_at"]) {
			return fmt.Errorf("current Receipt: Attempt start mismatch")
		}
	}
	for _, effect := range event.Effects {
		if effect.Table == "budget_samples" {
			run, err := currentReference(tx, "runs", map[string]any{"id": effect.Values["run_id"]})
			if err != nil {
				return err
			}
			if event.EventType != "run.recovered" && effect.Values["run_id"] != runID {
				return fmt.Errorf("current Budget: Run binding mismatch")
			}
			if !currentMatching(effect.Values, run, "envelope_id", "envelope_revision") {
				return fmt.Errorf("current Budget: Envelope binding mismatch")
			}
			if previous := currentEventValue(effect.Values["previous_sample_id"]); previous != nil {
				prior, err := currentReference(tx, "budget_samples", map[string]any{"id": previous})
				if err != nil {
					return err
				}
				if !currentMatching(effect.Values, prior, "envelope_id", "envelope_revision", "budget_name") {
					return fmt.Errorf("current Budget: prior chain binding mismatch")
				}
			}
		}
		if (event.EventType == "gate.execution-started" || event.EventType == "gate.receipt-recorded") && effect.Table == "run_runtime" && effect.PK["run_id"] != runID {
			return fmt.Errorf("current execution: runtime binding mismatch")
		}
		if effect.Table == "receipt_artifacts" && effect.PK["receipt_id"] != event.Aggregate.ID {
			return fmt.Errorf("current Receipt: Artifact link mismatch")
		}
		if event.EventType == "run.abandoned" && effect.Table == "claims" {
			run, err := currentReference(tx, "runs", map[string]any{"id": runID})
			if err != nil {
				return err
			}
			if effect.PK["id"] != run["claim_id"] {
				return fmt.Errorf("current abandonment: Claim binding mismatch")
			}
		}
		if event.EventType == "claim.expired" && effect.Table == "claims" && effect.PK["id"] != event.Aggregate.ID {
			return fmt.Errorf("current expiry: Claim identity mismatch")
		}
	}
	return nil
}

func validateCurrentRecoveryReferences(tx *sql.Tx, event Event) error {
	recovery := currentInserted(event, "run_recoveries")
	successor := currentInserted(event, "runs")
	claim := currentInserted(event, "claims")
	runtime := currentInserted(event, "run_runtime")
	if recovery["prior_run_id"] != event.Aggregate.ID || recovery["successor_run_id"] != successor["id"] || recovery["successor_claim_id"] != claim["id"] || successor["claim_id"] != claim["id"] || runtime["run_id"] != successor["id"] || successor["id"] == event.Aggregate.ID {
		return fmt.Errorf("current recovery: predecessor/successor binding mismatch")
	}
	prior, err := currentReference(tx, "runs", map[string]any{"id": event.Aggregate.ID})
	if err != nil {
		return err
	}
	if recovery["prior_claim_id"] != prior["claim_id"] || !currentMatching(successor, prior, "issue_id", "envelope_id", "envelope_revision") || !currentMatching(claim, successor, "issue_id", "envelope_id", "envelope_revision", "actor_id") {
		return fmt.Errorf("current recovery: Issue/Envelope/Claim binding mismatch")
	}
	observation, err := currentReference(tx, "run_observations", map[string]any{"id": recovery["observation_id"]})
	if err != nil {
		return err
	}
	if observation["run_id"] != event.Aggregate.ID || observation["kind"] != "liveness" {
		return fmt.Errorf("current recovery: observation binding mismatch")
	}
	for _, effect := range event.Effects {
		if effect.Op == "update" {
			id := effect.PK["id"]
			if effect.Table == "run_runtime" {
				id = effect.PK["run_id"]
			}
			if (effect.Table == "runs" || effect.Table == "run_runtime") && id != event.Aggregate.ID {
				return fmt.Errorf("current recovery: unrelated predecessor update")
			}
			if effect.Table == "claims" && id != prior["claim_id"] {
				return fmt.Errorf("current recovery: unrelated Claim disposal")
			}
		}
		if effect.Table == "run_budget_bindings" && effect.PK["run_id"] != successor["id"] {
			return fmt.Errorf("current recovery: successor budget binding mismatch")
		}
	}
	var authority map[string]any
	if json.Unmarshal([]byte(recovery["authority_json"].(string)), &authority) != nil || authority["actor_id"] != event.ActorV2.ActorID || authority["reason"] != recovery["reason"] {
		return fmt.Errorf("current recovery: authority attribution mismatch")
	}
	mode := "ordinary"
	if event.CommandName == "run.force-recover" {
		mode = "force"
	}
	if authority["mode"] != mode {
		return fmt.Errorf("current recovery: authority mode mismatch")
	}
	return nil
}
