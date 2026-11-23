package journal

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// validateCurrentIdempotency also binds the stored result to the Event's
// aggregate, for the families whose retry returns that aggregate.
func validateCurrentIdempotency(event Event) error {
	row, err := validateIdempotencyBinding(event)
	if err != nil {
		return err
	}
	if row["result_entity_id"] != event.Aggregate.ID {
		return fmt.Errorf("current Event: idempotency result binding mismatch")
	}
	return nil
}

// validateIdempotencyBinding requires exactly one retry result, bound to the
// Event's own request, identity and sequence.
func validateIdempotencyBinding(event Event) (map[string]any, error) {
	var result *Effect
	for index := range event.Effects {
		if event.Effects[index].Table != "idempotency_results" {
			continue
		}
		if result != nil {
			return nil, fmt.Errorf("current Event: duplicate idempotency Effect")
		}
		result = &event.Effects[index]
	}
	if result == nil {
		return nil, fmt.Errorf("current Event: missing idempotency Effect")
	}
	row, err := mergeEffectValues(result.PK, result.Values)
	if err != nil {
		return nil, err
	}
	if row["result_event_id"] != event.EventID || !currentValuesEqual(row["aggregate_sequence"], event.Aggregate.Sequence) {
		return nil, fmt.Errorf("current Event: idempotency result binding mismatch")
	}
	fieldDigests := result.Values["request_field_digests_json"]
	if result.Op == "update" {
		fieldDigests = result.Preconditions["request_field_digests_json"]
	}
	fieldsJSON, ok := fieldDigests.(string)
	if !ok {
		return nil, fmt.Errorf("current Event: missing idempotency field digests")
	}
	if err := validateCurrentJSON(event, *result, "request_field_digests_json", fieldsJSON); err != nil {
		return nil, fmt.Errorf("current Event: invalid idempotency field digests: %w", err)
	}
	{
		id := event.Idempotency
		digest := row["request_digest"]
		if result.Op == "update" {
			digest = result.Preconditions["request_digest"]
		}
		if id == nil || id.Key == "" || !strings.HasPrefix(id.Scope, event.CommandName+":") || id.Scope == event.CommandName+":" || !currentDigest(id.RequestDigest) || row["operation_scope"] != id.Scope || row["idempotency_key"] != id.Key || digest != id.RequestDigest {
			return nil, fmt.Errorf("current Event: idempotency request binding mismatch")
		}
		if event.Repository != nil && id.Scope != event.CommandName+":"+event.Repository.ProjectID {
			return nil, fmt.Errorf("current Event: idempotency project mismatch")
		}
	}
	var body map[string]any
	resultJSON, ok := row["result_json"].(string)
	if !ok {
		return nil, fmt.Errorf("current Event: invalid original idempotency result")
	}
	if err := unmarshalCurrentJSON(resultJSON, &body); err != nil {
		return nil, fmt.Errorf("current Event: invalid original idempotency result: %w", err)
	}
	if body["replayed"] != false {
		return nil, fmt.Errorf("current Event: invalid original idempotency result")
	}
	return row, nil
}

func currentDigest(value string) bool {
	value = strings.TrimPrefix(value, "sha256:")
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
