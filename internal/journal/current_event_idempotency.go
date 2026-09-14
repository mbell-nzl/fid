package journal

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

func validateCurrentIdempotency(event Event) error {
	var result *Effect
	for index := range event.Effects {
		if event.Effects[index].Table != "idempotency_results" {
			continue
		}
		if result != nil {
			return fmt.Errorf("current Event: duplicate idempotency Effect")
		}
		result = &event.Effects[index]
	}
	if result == nil {
		return fmt.Errorf("current Event: missing idempotency Effect")
	}
	row, err := mergeEffectValues(result.PK, result.Values)
	if err != nil {
		return err
	}
	if row["result_event_id"] != event.EventID || !currentValuesEqual(row["aggregate_sequence"], event.Aggregate.Sequence) || row["result_entity_id"] != currentResultEntity(event) {
		return fmt.Errorf("current Event: idempotency result binding mismatch")
	}
	if event.EventType == "gate.receipt-recorded" {
		if result.Op != "update" {
			return fmt.Errorf("current Receipt: invalid idempotency transition")
		}
	}
	{
		id := event.Idempotency
		digest := row["request_digest"]
		if result.Op == "update" {
			digest = result.Preconditions["request_digest"]
		}
		if id == nil || id.Key == "" || !strings.HasPrefix(id.Scope, event.CommandName+":") || id.Scope == event.CommandName+":" || !currentDigest(id.RequestDigest) || row["operation_scope"] != id.Scope || row["idempotency_key"] != id.Key || digest != id.RequestDigest {
			return fmt.Errorf("current Event: idempotency request binding mismatch")
		}
		if event.Repository != nil && id.Scope != event.CommandName+":"+event.Repository.ProjectID {
			return fmt.Errorf("current Event: idempotency project mismatch")
		}
	}
	var body map[string]any
	if json.Unmarshal([]byte(row["result_json"].(string)), &body) != nil || body["replayed"] != false {
		return fmt.Errorf("current Event: invalid original idempotency result")
	}
	return nil
}

func currentDigest(value string) bool {
	value = strings.TrimPrefix(value, "sha256:")
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func currentResultEntity(event Event) any {
	table := map[string]string{"run.observed": "run_observations", "run.checkpointed": "run_observations", "run.recovered": "run_recoveries", "run.budget-observed": "budget_samples", "run.execution-blocked": "stops"}[event.EventType]
	if table != "" {
		for _, effect := range event.Effects {
			if effect.Table == table && effect.Op == "insert" {
				return effect.PK["id"]
			}
		}
	}
	return event.Aggregate.ID
}
