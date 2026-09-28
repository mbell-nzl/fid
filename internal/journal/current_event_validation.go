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
	if !event.IsV2() {
		return nil
	}
	for _, effect := range event.Effects {
		switch effect.Op {
		case "insert", "update", "delete", "activity":
		default:
			return fmt.Errorf("project event %q: unsupported effect op %q", event.EventID, effect.Op)
		}
	}
	rule, current := currentEventRules[event.EventType]
	if !current {
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
		// Retained historical writers are not assigned new generic semantics by
		// Phase 22. They still cannot smuggle a current domain table under an
		// unrelated family. Their existing SQL/precondition validation remains.
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
	// The pinned pre-adoption archive has its own accepted grammar. Its
	// projection has no recovery runtime, so Phase 22 replacements do not apply.
	var recoverySchema int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='run_runtime'`).Scan(&recoverySchema); err != nil {
		return err
	}
	if recoverySchema == 0 {
		return nil
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	commandMatches := event.CommandName == currentEventCommands[event.EventType] || (event.EventType == "run.recovered" && event.CommandName == "run.force-recover")
	if event.Aggregate.Type != rule.aggregate || !commandMatches {
		return fmt.Errorf("current Event %s: aggregate/command does not match family", event.EventType)
	}
	if err := validateCurrentStages(event, rule); err != nil {
		return err
	}
	for _, effect := range event.Effects {
		if err := validateCurrentEffect(event, effect); err != nil {
			return fmt.Errorf("current Event %s %s/%s: %w", event.EventType, effect.Table, effect.Op, err)
		}
		if err := validateCurrentPrior(tx, event, effect); err != nil {
			return err
		}
	}
	if err := validateCurrentBindings(event); err != nil {
		return err
	}
	if err := validateCurrentCardinality(event); err != nil {
		return err
	}
	if err := validateCurrentStopBundles(tx, event); err != nil {
		return err
	}
	for _, effect := range event.Effects {
		for column, value := range effect.Values {
			if (strings.HasSuffix(column, "_json") || effect.Table == "activity_log" && column == "details") && currentEventValue(value) != nil {
				if err := validateCurrentJSON(event, effect, column, currentEventValue(value).(string)); err != nil {
					return fmt.Errorf("current Event %s %s/%s: invalid JSON column %q: %w", event.EventType, effect.Table, effect.Op, column, err)
				}
			}
		}
	}
	if err := validateCurrentComponentValues(event); err != nil {
		return err
	}
	if err := validateCurrentReferences(tx, event); err != nil {
		return err
	}
	if err := validateCurrentStopReferences(tx, event); err != nil {
		return err
	}
	return validateCurrentIdempotency(event)
}

func validateCurrentStages(event Event, rule currentEventRule) error {
	stage, count := 0, 0
	for _, effect := range event.Effects {
		key := currentEffectKey(effect.Table, effect.Op)
		// Automatic boundaries emit repeated complete Stop/Attention/link bundles.
		// Bundle structure and identities are separately validated before mutation.
		if key == "stops/insert" && stage < len(rule.stages) && rule.stages[stage].allowed["attention_stop_links/insert"] {
			for candidate := 0; candidate < stage; candidate++ {
				if rule.stages[candidate].allowed[key] {
					stage = candidate
					count = 0
					break
				}
			}
		}
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

func validateCurrentBindings(event Event) error {
	entityTable := map[string]string{"gate.execution-started": "execution_attempts", "gate.receipt-recorded": "execution_receipts", "stop.raised": "stops", "stop.resolved": "stop_resolutions", "attention.resolved": "attention_items", "concern.recorded": "judgment_concerns", "concern.disposed": "judgment_concerns", "deviation.recorded": "deviations"}[event.EventType]
	for _, effect := range event.Effects {
		row, err := mergeEffectValues(effect.PK, effect.Values)
		if err != nil {
			return err
		}
		if effect.Table == entityTable {
			id := row["id"]
			if effect.Table == "stop_resolutions" {
				id = row["stop_id"]
			}
			if id != event.Aggregate.ID {
				return fmt.Errorf("current Event: entity/aggregate identity mismatch")
			}
		}
		if effect.Table == "run_observations" {
			kind := "liveness"
			if event.EventType == "run.checkpointed" {
				kind = "checkpoint"
			}
			if row["run_id"] != event.Aggregate.ID || row["kind"] != kind || row["actor_id"] != event.ActorV2.ActorID {
				return fmt.Errorf("current observation binding mismatch")
			}
			var observation map[string]any
			if json.Unmarshal([]byte(row["observation_json"].(string)), &observation) != nil || observation["id"] != row["id"] || observation["run_id"] != row["run_id"] || observation["kind"] != row["kind"] {
				return fmt.Errorf("current observation payload binding mismatch")
			}
		}
		if strings.HasPrefix(event.EventType, "activation.") && effect.Table != "idempotency_results" {
			if id, found := row["activation_id"]; found && id != event.Aggregate.ID {
				return fmt.Errorf("current activation binding mismatch")
			}
		}
		if event.EventType != "run.recovered" && strings.HasPrefix(event.EventType, "run.") {
			if id, found := row["run_id"]; found && id != event.Aggregate.ID {
				return fmt.Errorf("current Run binding mismatch")
			}
			if effect.Table == "runs" && row["id"] != event.Aggregate.ID {
				return fmt.Errorf("current Run identity mismatch")
			}
		}
	}
	return nil
}
