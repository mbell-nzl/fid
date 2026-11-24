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
	if err := ValidateEventShape(event); err != nil {
		return err
	}
	if err := validateStateMoves(tx, event); err != nil {
		return err
	}
	if err := validateIssueIdentity(event); err != nil {
		return err
	}
	if err := validateIssueRows(event); err != nil {
		return err
	}
	if err := validateIssueHoldBindings(tx, event); err != nil {
		return err
	}
	// These tables have no legacy writer. The shape rules already refuse case
	// aliases, and Event-v1 Effects outside the Issue tables; these guards
	// stay as defence in depth before the Event-v1 compatibility return.
	for _, effect := range event.Effects {
		if err := validateEpicCloseOwnership(tx, event, effect); err != nil {
			return err
		}
		if strings.EqualFold(effect.Table, "questions") {
			for _, row := range []map[string]any{effect.PK, effect.Values, effect.Preconditions, effect.Prior} {
				for column, value := range row {
					if strings.EqualFold(column, "resolution") && value == "finished-before-answer" && (!event.IsV2() || event.EventType != "task.finished" || effect.Op != "update") {
						return fmt.Errorf("finished questions require their completion Event-v2 family")
					}
				}
			}
		}
		if strings.EqualFold(effect.Table, "completions") && (!event.IsV2() || event.EventType != "task.finished" || effect.Op != "insert") {
			return fmt.Errorf("completion records require their immutable Event-v2 family")
		}

		if strings.EqualFold(effect.Table, "holds") {
			for _, values := range []map[string]any{effect.PK, effect.Values, effect.Preconditions, effect.Prior} {
				for column := range values {
					if strings.EqualFold(column, "rework_at") && !validReworkColumn(event, effect, values[column]) {
						return fmt.Errorf("rework reservation requires a review round or default conflict")
					}
					if strings.EqualFold(column, "model") && (!event.IsV2() || effect.Op != "insert" || (event.EventType != "hold.taken" && event.EventType != "task.taken-over")) {
						return fmt.Errorf("build models require a validated Event-v2 hold insert")
					}
				}
			}
		}
		if strings.EqualFold(effect.Table, "person_decisions") {
			if err := validateApprovalDecisionEffect(tx, event, effect); err != nil {
				return err
			}
			if err := validateFindingDecisionEffect(tx, event, effect); err != nil {
				return err
			}
		}
		if strings.EqualFold(effect.Table, "finding_closures") && (!event.IsV2() || event.EventType != "review.round-recorded" && event.EventType != "finding.closed") {
			return fmt.Errorf("finding closures require their Event-v2 family")
		}
		if (strings.EqualFold(effect.Table, "review_rounds") || strings.EqualFold(effect.Table, "round_findings")) && (!event.IsV2() || event.EventType != "review.round-recorded") {
			return fmt.Errorf("review tables require their Event-v2 family")
		}
		if strings.EqualFold(effect.Table, "merge_observations") && (!event.IsV2() || event.EventType != "merge.observed") || strings.EqualFold(effect.Table, "ci_results") && (!event.IsV2() || event.EventType != "ci.result-recorded") {
			return fmt.Errorf("merge/CI tables require their validated Event-v2 family")
		}
	}
	if !event.IsV2() && (event.EventType == "merge.observed" || event.EventType == "ci.result-recorded") {
		return fmt.Errorf("merge/CI observations require Event-v2")
	}
	if !event.IsV2() && (event.EventType == "review.round-recorded" || event.EventType == "finding.closed") {
		return fmt.Errorf("review rounds require Event-v2")
	}
	if !event.IsV2() && event.EventType == "epic.closed" {
		return fmt.Errorf("epic close requires Event-v2")
	}
	if !event.IsV2() && (event.EventType == "task.approved" || event.EventType == "task.waived" || event.EventType == "task.finished") {
		return fmt.Errorf("approval requires Event-v2")
	}
	if err := validateReworkHold(tx, event); err != nil {
		return err
	}
	if !event.IsV2() {
		if event.EventType == "artifact.availability-observed" {
			return fmt.Errorf("artifact availability requires Event-v2")
		}
		return nil
	}
	if event.EventType == "task.default-conflict" {
		return validateDefaultConflict(tx, event)
	}
	if event.EventType == "work.accepted" || event.EventType == "rework.allowed" {
		return validateWorkDecision(event, tx)
	}
	// The shape rules have checked the family, its actor, aggregate and
	// Effect grammar. The family validators check the projection state.
	if event.EventType == "review.claimed" || event.EventType == "review.heartbeat" {
		return validateReviewClaim(tx, event)
	}
	if event.EventType == "merge.observed" || event.EventType == "ci.result-recorded" {
		return validateMergeCI(tx, event)
	}
	if event.EventType == "epic.closed" {
		return validateEpicClose(tx, event)
	}
	if event.EventType == "task.finished" {
		return validateCompletion(tx, event)
	}
	if event.EventType == "task.waived" {
		return validateWaiver(tx, event)
	}
	if event.EventType == "task.approved" {
		return validateApproval(tx, event)
	}
	if event.EventType == "finding.closed" {
		return validateFindingClosed(tx, event)
	}
	if event.EventType == "review.round-recorded" {
		return validateReviewRound(tx, event)
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
	if event.EventType == "artifact.availability-observed" {
		return validateArtifactAvailability(tx, event)
	}
	if err := validateHoldLifecycle(tx, event); err != nil {
		return err
	}
	return validatePolicySnapshot(tx, event)
}

func validateCurrentStages(event Event, rule currentEventRule) error {
	stage, count := 0, 0
	for _, effect := range event.Effects {
		key := currentEffectKey(effect.Table, effect.Op)
		for stage < len(rule.stages) && !rule.stages[stage].allowed[key] {
			if count < rule.stages[stage].minimum {
				return fmt.Errorf("current Event %s: Effect grammar: missing or reordered stage %d", event.EventType, stage+1)
			}
			stage++
			count = 0
		}
		if stage == len(rule.stages) {
			return fmt.Errorf("current Event %s: Effect grammar: unrelated or reordered Effect %s", event.EventType, key)
		}
		count++
	}
	for stage < len(rule.stages) {
		if count < rule.stages[stage].minimum {
			return fmt.Errorf("current Event %s: Effect grammar: missing stage %d", event.EventType, stage+1)
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
		if effect.Table == "person_decisions" {
			if _, present := row["binding_json"]; !present {
				row["binding_json"] = nil
			}
		}
		if effect.Table == "holds" {
			if _, present := row["rework_at"]; !present {
				row["rework_at"] = nil
			}
			if _, present := row["model"]; !present {
				row["model"] = ""
			}
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
