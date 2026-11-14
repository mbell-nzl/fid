package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
)

type currentTable struct {
	primary []string
	columns map[string]string
}

func currentTableColumns(primary, columns string) currentTable {
	result := currentTable{primary: strings.Fields(primary), columns: map[string]string{}}
	for _, column := range strings.Fields(columns) {
		name, kind, _ := strings.Cut(column, ":")
		result.columns[name] = kind
	}
	return result
}

func validateCurrentEvent(tx *sql.Tx, event Event) error {
	if err := ValidateEventSupport(event); err != nil {
		return err
	}
	if !event.IsV2() {
		if event.EventType == "artifact.availability-observed" {
			return fmt.Errorf("artifact availability requires Event-v2")
		}
		return nil
	}
	for _, effect := range event.Effects {
		switch effect.Op {
		case "insert", "update", "delete", "activity":
		default:
			return fmt.Errorf("project event %q: unsupported effect op %q", event.EventID, effect.Op)
		}
	}
	if event.EventType == "task.handed-off" {
		return validateHandoff(tx, event)
	}
	if strings.HasPrefix(event.EventType, "check.") {
		return validateLeanCheck(tx, event)
	}
	if err := validateTaskStart(tx, event); err != nil {
		return err
	}
	if err := validateTaskLane(tx, event); err != nil {
		return err
	}
	retained, found := retainedCurrentRules[event.EventType]
	if !found {
		return fmt.Errorf("current Event: unsupported Event family %q", event.EventType)
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	if event.Aggregate.Type != retained.aggregate {
		return fmt.Errorf("current Event: retained aggregate mismatch")
	}
	if event.EventType == "artifact.availability-observed" {
		return validateArtifactAvailability(tx, event)
	}
	if leanEventFamilies[event.EventType] {
		if err := validateCurrentStages(event, retained); err != nil {
			return err
		}
		if err := validateHoldLifecycle(tx, event); err != nil {
			return err
		}
		return validatePolicySnapshot(tx, event)
	}
	if err := validateHoldLifecycle(tx, event); err != nil {
		return err
	}
	// Retained Issue writers keep their historical SQL/precondition semantics.
	// They cannot write a lean domain table owned by another family.
	for _, effect := range event.Effects {
		if _, relevant := currentEventTables[effect.Table]; !relevant || effect.Table == "idempotency_results" {
			continue
		}
		allowed := false
		for _, stage := range retained.stages {
			allowed = allowed || stage.allowed[currentEffectKey(effect.Table, effect.Op)]
		}
		if !allowed {
			return fmt.Errorf("current Event %s: unrelated domain Effect", event.EventType)
		}
	}
	return nil
}

func validateCurrentStages(event Event, rule currentEventRule) error {
	stage, count := 0, 0
	for _, effect := range event.Effects {
		key := currentEffectKey(effect.Table, effect.Op)
		for stage < len(rule.stages) && !rule.stages[stage].allowed[key] {
			if count < rule.stages[stage].minimum {
				return fmt.Errorf("current Event %s: missing or reordered Effect stage %d", event.EventType, stage+1)
			}
			stage++
			count = 0
		}
		if stage == len(rule.stages) {
			return fmt.Errorf("current Event %s: unrelated or reordered Effect %s", event.EventType, key)
		}
		count++
	}
	for stage < len(rule.stages) {
		if count < rule.stages[stage].minimum {
			return fmt.Errorf("current Event %s: missing Effect stage %d", event.EventType, stage+1)
		}
		stage++
		count = 0
	}
	return nil
}

func validateCurrentEffect(event Event, effect Effect) error {
	table, found := currentEventTables[effect.Table]
	if !found {
		return fmt.Errorf("unknown table")
	}
	if len(effect.PK) != len(table.primary) {
		return fmt.Errorf("incomplete primary key")
	}
	for _, key := range table.primary {
		if _, found := effect.PK[key]; !found {
			return fmt.Errorf("wrong primary key")
		}
	}
	for _, values := range []map[string]any{effect.PK, effect.Values, effect.Preconditions, effect.Prior} {
		for column, value := range values {
			value = currentEventValue(value)
			kind, found := table.columns[column]
			if !found || !currentScalar(kind, value) {
				return fmt.Errorf("invalid column or scalar %q", column)
			}
			if strings.HasSuffix(column, "_json") && value != nil {
				if _, err := decodeCurrentJSON(value.(string)); err != nil {
					return fmt.Errorf("invalid JSON column %q: %w", column, err)
				}
			}
		}
	}
	if effect.Op == "insert" {
		row, err := mergeEffectValues(effect.PK, effect.Values)
		if err != nil {
			return err
		}
		if len(row) != len(table.columns) || len(effect.Prior) != 0 || len(effect.Preconditions) != 0 {
			return fmt.Errorf("insert must contain complete row and no prior/precondition")
		}
		if version, exists := row["version"]; exists && !currentValuesEqual(version, 1) {
			return fmt.Errorf("current Event: inserted version must be one")
		}
	}
	if effect.Op == "update" {
		if err := validateCurrentUpdate(event, effect); err != nil {
			return err
		}
		if len(effect.Values) == 0 {
			return fmt.Errorf("empty update")
		}
		if next, found := effect.Values["version"]; found {
			previous, ok := effect.Preconditions["version"]
			if !ok || !currentNextVersion(previous, next) {
				return fmt.Errorf("version update requires exact previous version plus one")
			}
		}
	}
	if created, found := effect.Values["created_event_id"]; found && created != event.EventID {
		return fmt.Errorf("created Event identity mismatch")
	}
	for _, values := range []map[string]any{effect.PK, effect.Values, effect.Preconditions, effect.Prior} {
		for column, value := range values {
			if err := validateCurrentColumnDomain(column, currentEventValue(value)); err != nil {
				return err
			}
		}
	}
	return nil
}

func currentScalar(kind string, value any) bool {
	if value == nil {
		return strings.HasSuffix(kind, "?")
	}
	if strings.HasPrefix(kind, "s") {
		_, ok := value.(string)
		return ok
	}
	if kind == "b" {
		if _, ok := value.(bool); ok {
			return true
		}
		encoded, err := json.Marshal(value)
		return err == nil && (string(encoded) == "0" || string(encoded) == "1")
	}
	switch n := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
		return true
	case float64:
		return !math.IsNaN(n) && !math.IsInf(n, 0) && math.Trunc(n) == n && n >= math.MinInt64 && n < math.MaxInt64
	case json.Number:
		_, err := n.Int64()
		return err == nil
	default:
		return false
	}
}

func currentEventValue(value any) any {
	for value != nil {
		ref := reflect.ValueOf(value)
		if ref.Kind() != reflect.Pointer {
			return value
		}
		if ref.IsNil() {
			return nil
		}
		value = ref.Elem().Interface()
	}
	return nil
}

func currentNextVersion(previous, next any) bool {
	p, err := json.Marshal(previous)
	if err != nil {
		return false
	}
	n, err := json.Marshal(next)
	if err != nil {
		return false
	}
	var a, b int64
	return json.Unmarshal(p, &a) == nil && json.Unmarshal(n, &b) == nil && a > 0 && a < math.MaxInt64 && b == a+1
}
