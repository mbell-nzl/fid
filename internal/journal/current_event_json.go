package journal

import (
	"fmt"
	"strings"
)

func validateCurrentJSON(event Event, effect Effect, column, raw string) error {
	value, err := decodeCurrentJSON(raw)
	if err != nil {
		return err
	}
	switch effect.Table + "." + column {
	case "activity_log.details":
		if _, ok := value.(map[string]any); !ok {
			return fmt.Errorf("activity details must be a JSON object")
		}
		return nil
	case "idempotency_results.result_json":
		return validateCurrentResultJSON(event, value)
	case "idempotency_results.request_field_digests_json":
		fields, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("idempotency field digests must be a JSON object")
		}
		for _, digest := range fields {
			if err := validateCurrentValue("Digest", digest); err != nil {
				return fmt.Errorf("invalid idempotency field digest: %w", err)
			}
		}
		return nil
	default:
		return fmt.Errorf("unregistered current JSON column %s.%s", effect.Table, column)
	}
}
func validateCurrentColumnDomain(column string, value any) error {
	if value == nil {
		return nil
	}
	if column == "request_digest" {
		// The existing Event idempotency grammar admits both hash spellings.
		return validateCurrentValue("Hash", value)
	}
	if strings.HasSuffix(column, "_at") {
		return validateCurrentValue("Timestamp", value)
	}
	if column == "id" || strings.HasSuffix(column, "_id") || column == "idempotency_key" {
		return validateCurrentValue("ID", value)
	}
	if column == "digest" || strings.HasSuffix(column, "_digest") || column == "semantic_key" {
		return validateCurrentValue("Digest", value)
	}
	if column == "declared_path" {
		return validateCurrentValue("Path", value)
	}
	return nil
}
func validateCurrentResultJSON(event Event, value any) error {
	if event.EventType != "artifact.availability-observed" {
		return fmt.Errorf("unregistered current result family")
	}
	_, err := validateCurrentShape(value, currentValueShapes["Artifact"])
	return err
}
