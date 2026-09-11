package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
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
	sort.Strings(result)
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
		if effect.Table == "run_observations" && effect.Values["kind"] == "liveness" {
			runtime, err := currentReference(tx, "run_runtime", map[string]any{"run_id": effect.Values["run_id"]})
			if err != nil {
				return err
			}
			value, err := decodeCurrentJSON(effect.Values["observation_json"].(string))
			if err != nil {
				return err
			}
			observation, ok := value.(map[string]any)
			if !ok {
				return fmt.Errorf("current observation: invalid payload")
			}
			liveness, ok := observation["liveness"].(map[string]any)
			if !ok {
				return fmt.Errorf("current observation: invalid liveness payload")
			}
			if liveness["executor_id"] != runtime["executor_id"] {
				return fmt.Errorf("current observation: executor binding mismatch")
			}
		}
		if effect.Table == "budget_samples" {
			policyRunID := effect.Values["run_id"]
			var run map[string]any
			var err error
			if event.EventType == "run.recovered" {
				run = currentInserted(event, "runs")
				if effect.Values["run_id"] != run["id"] {
					return fmt.Errorf("current intervention: successor Run binding mismatch")
				}
				// Its exact policy inheritance is validated by recovery references;
				// the successor binding does not exist until Effects are applied.
				policyRunID = event.Aggregate.ID
			} else {
				run, err = currentReference(tx, "runs", map[string]any{"id": effect.Values["run_id"]})
				if err != nil {
					return err
				}
			}
			if err := validateCurrentBudgetSampleChain(tx, effect, policyRunID); err != nil {
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
	if err := validateCurrentRecoverySnapshot(event, recovery, successor, claim); err != nil {
		return err
	}
	if err := validateCurrentRecoveryEffects(tx, event, prior, successor); err != nil {
		return err
	}
	return nil
}

// A recovery records the state it actually fenced.  Terminal predecessors and
// already disposed Claims are legitimate omissions; active rows must carry the
// matching captured mutation in the same Event.
func validateCurrentRecoveryEffects(tx *sql.Tx, event Event, prior, successor map[string]any) error {
	updates := map[string]int{}
	for _, effect := range event.Effects {
		updates[currentEffectKey(effect.Table, effect.Op)]++
	}
	activeRun := prior["status"] == "active"
	if (updates["runs/update"] == 1) != activeRun || (updates["run_runtime/update"] == 1) != activeRun {
		return fmt.Errorf("current recovery: predecessor terminalization does not match stored Run state")
	}
	claim, err := currentReference(tx, "claims", map[string]any{"id": prior["claim_id"]})
	if err != nil {
		return err
	}
	var previousScope, nextScope any
	successorClaim := currentInserted(event, "claims")
	if json.Unmarshal([]byte(claim["scope_json"].(string)), &previousScope) != nil || json.Unmarshal([]byte(successorClaim["scope_json"].(string)), &nextScope) != nil || !currentValuesEqual(previousScope, nextScope) || successorClaim["mode"] != claim["mode"] {
		return fmt.Errorf("current recovery: successor Claim changes immutable scope")
	}
	if (updates["claims/update"] == 1) != (claim["status"] == "active") {
		return fmt.Errorf("current recovery: Claim disposal does not match stored Claim state")
	}
	var expected int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM run_budget_bindings WHERE run_id = ?`, prior["id"]).Scan(&expected); err != nil {
		return err
	}
	if updates["run_budget_bindings/insert"] != expected {
		return fmt.Errorf("current recovery: successor budget bindings are incomplete")
	}
	for _, effect := range event.Effects {
		if effect.Table != "run_budget_bindings" || effect.Op != "insert" {
			continue
		}
		var recorded string
		if err := tx.QueryRow(`SELECT budget_json FROM run_budget_bindings WHERE run_id = ? AND name = ?`, prior["id"], effect.PK["name"]).Scan(&recorded); err != nil {
			return err
		}
		var previous, inherited any
		if json.Unmarshal([]byte(recorded), &previous) != nil || json.Unmarshal([]byte(effect.Values["budget_json"].(string)), &inherited) != nil || !currentValuesEqual(previous, inherited) || effect.PK["run_id"] != successor["id"] {
			return fmt.Errorf("current recovery: successor budget binding does not match predecessor")
		}
	}
	return nil
}

func validateCurrentRecoverySnapshot(event Event, recovery, successor, claim map[string]any) error {
	value, err := decodeCurrentJSON(recovery["reconciliation_json"].(string))
	if err != nil {
		return err
	}
	repository := value.(map[string]any)["repository"].(map[string]any)
	start, err := decodeCurrentJSON(successor["start_repository_json"].(string))
	if err != nil {
		return err
	}
	if !currentValuesEqual(start, repository) || claim["worktree_id"] != repository["worktree_id"] {
		return fmt.Errorf("current recovery: repository/Claim worktree mismatch")
	}
	if successor["actor_id"] != event.ActorV2.ActorID || !currentValuesEqual(successor["claim_version"], claim["version"]) || successor["started_at"] != claim["acquired_at"] {
		return fmt.Errorf("current recovery: successor Actor/version/time mismatch")
	}
	for _, effect := range event.Effects {
		if effect.Op != "update" {
			continue
		}
		if effect.Table == "runs" {
			terminal, err := decodeCurrentJSON(effect.Values["terminal_repository_json"].(string))
			if err != nil {
				return err
			}
			if !currentValuesEqual(terminal, repository) || effect.Values["ended_at"] != successor["started_at"] {
				return fmt.Errorf("current recovery: predecessor terminal snapshot mismatch")
			}
		}
		if effect.Table == "claims" && effect.Values["released_at"] != successor["started_at"] {
			return fmt.Errorf("current recovery: Claim disposal time mismatch")
		}
	}
	return nil
}
