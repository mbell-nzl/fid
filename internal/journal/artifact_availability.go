package journal

import (
	"database/sql"
	"fmt"
)

func validateArtifactAvailability(tx *sql.Tx, event Event) error {
	if event.CommandName != "artifact.verify" || event.Repository != nil || event.ActorV2.ActorType != "agent" && event.ActorV2.ActorType != "human" {
		return fmt.Errorf("invalid Artifact availability command or actor")
	}
	if err := validateCurrentStages(event, retainedCurrentRules[event.EventType]); err != nil {
		return err
	}
	updates, results := 0, 0
	for _, effect := range event.Effects {
		if err := validateCurrentEffect(event, effect); err != nil {
			return err
		}
		if err := validateCurrentPrior(tx, event, effect); err != nil {
			return err
		}
		switch effect.Table {
		case "referenced_artifacts":
			updates++
			if effect.PK["id"] != event.Aggregate.ID || !currentEnum(effect.Values["availability"], "available", "missing") || !currentEnum(effect.Preconditions["availability"], "available", "missing") || effect.Values["availability"] == effect.Preconditions["availability"] {
				return fmt.Errorf("invalid Artifact availability transition")
			}
		case "idempotency_results":
			results++
		}
	}
	if updates > 1 || results != 1 {
		return fmt.Errorf("artifact availability requires at most one update and one result")
	}
	return validateCurrentIdempotency(event)
}

func validateArtifactResult(tx *sql.Tx, event Event) error {
	var result map[string]any
	updates := 0
	for _, effect := range event.Effects {
		if effect.Table == "referenced_artifacts" {
			updates++
		}
		if effect.Table == "idempotency_results" {
			value, err := decodeCurrentJSON(effect.Values["result_json"].(string))
			if err != nil {
				return err
			}
			var ok bool
			result, ok = value.(map[string]any)
			if !ok {
				return fmt.Errorf("artifact result requires an object")
			}
		}
	}
	if _, err := validateCurrentShape(result, "artifact:Artifact changed:Bool idempotency_key:ID replayed:False EventID:Text EventOccurredAt:Text"); err != nil {
		return err
	}
	artifact, found, err := readCurrentEventValues(tx, "referenced_artifacts", map[string]any{"id": event.Aggregate.ID}, []string{"id", "digest", "media_type", "size_bytes", "retention_class", "availability", "location_policy", "created_at", "expires_at"})
	if err != nil {
		return err
	}
	if !found || !currentValuesEqual(result["artifact"], artifact) || result["changed"] != (updates == 1) || result["idempotency_key"] != event.Idempotency.Key || result["EventID"] != "" || result["EventOccurredAt"] != "" {
		return fmt.Errorf("artifact result differs from recorded facts")
	}
	return nil
}
