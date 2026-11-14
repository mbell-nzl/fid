package export

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

func validateSemanticMap(tableName string, table semanticTable, values map[string]any, required bool) error {
	if required && len(values) == 0 {
		return fmt.Errorf("must be non-empty")
	}
	if required && !validSemanticPrimaryKey(tableName, values, table.primary) {
		return fmt.Errorf("must contain exactly the table primary-key columns")
	}
	for column, value := range values {
		kind, found := table.columns[column]
		if !found {
			return fmt.Errorf("unknown column %q", column)
		}
		if tableName == "activity_log" && column == "details" {
			continue
		}
		if _, boolean := value.(bool); boolean && strings.Contains(kind, "INT") && !semanticBooleanColumn(tableName, column) {
			return fmt.Errorf("column %q requires an integer, not a boolean", column)
		}
		if err := validateSemanticValue(kind, value); err != nil {
			return fmt.Errorf("column %q: %w", column, err)
		}
		if strings.HasSuffix(column, "_json") && value != nil {
			encoded, ok := value.(string)
			if !ok {
				return fmt.Errorf("column %q: requires JSON text", column)
			}
			if err := validateSemanticJSONText(encoded); err != nil {
				return fmt.Errorf("column %q: %w", column, err)
			}
		}
	}
	return nil
}

func semanticBooleanColumn(table, column string) bool {
	switch table + "." + column {
	case "holds.kept":
		return true
	default:
		return false
	}
}

func validSemanticPrimaryKey(tableName string, values map[string]any, primary map[string]struct{}) bool {
	if len(primary) != 0 {
		return sameSemanticKeys(values, primary)
	}
	return false
}

func sameSemanticKeys(values map[string]any, expected map[string]struct{}) bool {
	if len(values) != len(expected) {
		return false
	}
	for key := range values {
		if _, found := expected[key]; !found {
			return false
		}
	}
	return true
}

func validateSemanticJSONText(value string) error {
	if err := rejectDuplicateJSONKeys([]byte(value)); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func validateSemanticValue(kind string, value any) error {
	if value == nil {
		return nil
	}
	if strings.Contains(kind, "BOOL") {
		if _, ok := value.(bool); ok {
			return nil
		}
		if n, ok := value.(float64); ok && (n == 0 || n == 1) {
			return nil
		}
		return fmt.Errorf("requires boolean JSON value or zero/one")
	}
	if strings.Contains(kind, "INT") {
		if _, ok := value.(bool); ok {
			return nil
		}
		number, ok := value.(float64)
		if !ok || math.Trunc(number) != number {
			return fmt.Errorf("requires integer JSON value")
		}
		return nil
	}
	if strings.Contains(kind, "TEXT") || kind == "" {
		if _, ok := value.(string); !ok {
			return fmt.Errorf("requires string JSON value")
		}
	}
	return nil
}
