package journal

import (
	"database/sql"
	"fmt"
	"strings"
)

// Read only the projection inside the Event transaction, after its Effects but
// before commit. A mismatching cached response rolls the complete Event back.
func validateCurrentAppliedResults(tx *sql.Tx, event Event) error {
	if _, current := currentEventRules[event.EventType]; !current || !event.IsV2() {
		return nil
	}
	var schema int
	if err := tx.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='run_runtime'").Scan(&schema); err != nil {
		return err
	}
	if schema == 0 {
		return nil
	}
	if err := validateCurrentBudgetFacts(tx, event); err != nil {
		return err
	}
	var result map[string]any
	for _, effect := range event.Effects {
		if effect.Table == "idempotency_results" {
			v, err := decodeCurrentJSON(effect.Values["result_json"].(string))
			if err != nil {
				return err
			}
			result = v.(map[string]any)
		}
	}
	if result == nil {
		return fmt.Errorf("current Event missing result")
	}
	expected, err := currentExpectedResult(tx, event)
	if err != nil {
		return err
	}
	expected["idempotency_key"] = event.Idempotency.Key
	expected["replayed"] = false
	if !currentValuesEqual(result, expected) {
		for key, want := range expected {
			if !currentValuesEqual(result[key], want) {
				return fmt.Errorf("current Event %s: stored result %s disagrees with recorded facts", event.EventType, currentResultDifference(key, result[key], want))
			}
		}
		return fmt.Errorf("current Event: stored result fields mismatch")
	}
	return nil
}

func currentResultDifference(prefix string, actual, expected any) string {
	if want, ok := expected.(map[string]any); ok {
		got, ok := actual.(map[string]any)
		if !ok {
			return prefix
		}
		for _, key := range sortedMapKeys(want) {
			if !currentValuesEqual(got[key], want[key]) {
				return currentResultDifference(prefix+"."+key, got[key], want[key])
			}
		}
	}
	return prefix
}
func currentExpectedResult(tx *sql.Tx, event Event) (map[string]any, error) {
	result := map[string]any{}
	var failure error
	put := func(key, kind, table string, pk map[string]any) {
		if failure != nil {
			return
		}
		var v map[string]any
		v, failure = currentProjectedValue(tx, kind, table, pk)
		result[key] = v
	}
	id := func(value any) map[string]any { return map[string]any{"id": value} }
	switch event.EventType {
	case "review.recorded":
		put("review", "Review", "review_records", id(event.Aggregate.ID))
		if failure == nil {
			var budgetResult map[string]any
			budgetResult, failure = currentReviewBudgetResult(tx, event)
			if failure == nil {
				result["budget"], result["stops"], result["attention"] = budgetResult["budget"], budgetResult["stops"], budgetResult["attention"]
			}
		}
	case "finding.addressed":
		put("finding", "Finding", "review_findings", id(event.Aggregate.ID))
	case "run.started", "run.closed", "run.observed", "run.checkpointed", "run.abandoned":
		put("run", "Run", "runs", id(event.Aggregate.ID))
		if event.EventType == "run.abandoned" && failure == nil {
			put("claim", "Claim", "claims", id(result["run"].(map[string]any)["claim_id"]))
		}
		if event.EventType == "run.observed" || event.EventType == "run.checkpointed" {
			kind, key := "Observation", "observation"
			if event.EventType == "run.checkpointed" {
				kind, key = "Checkpoint", "checkpoint"
			}
			put(key, kind, "run_observations", id(currentInserted(event, "run_observations")["id"]))
		}
	case "run.recovered":
		recovery := currentInserted(event, "run_recoveries")
		put("recovery", "Recovery", "run_recoveries", id(recovery["id"]))
		put("run", "Run", "runs", id(recovery["successor_run_id"]))
		put("claim", "Claim", "claims", id(recovery["successor_claim_id"]))
	case "claim.expired":
		put("claim", "Claim", "claims", id(event.Aggregate.ID))
	case "stop.raised", "stop.resolved":
		put("stop", "Stop", "stops", id(event.Aggregate.ID))
		if event.EventType == "stop.raised" {
			for _, e := range event.Effects {
				if e.Table == "attention_stop_links" {
					put("attention", "Attention", "attention_items", id(e.PK["attention_id"]))
				}
			}
		}
	case "attention.resolved":
		put("attention", "Attention", "attention_items", id(event.Aggregate.ID))
	case "gate.execution-started":
		put("attempt", "Attempt", "execution_attempts", id(event.Aggregate.ID))
		result["receipt"] = nil
	case "gate.receipt-recorded":
		put("receipt", "Receipt", "execution_receipts", id(event.Aggregate.ID))
	case "run.budget-observed", "run.execution-blocked":
		return currentBudgetEventResult(tx, event)
	default:
		if !strings.HasPrefix(event.EventType, "activation.") {
			return nil, fmt.Errorf("unknown current result")
		}
		put("activation", "Activation", "kit_activations", id(event.Aggregate.ID))
		if event.EventType == "activation.component-staged" {
			for _, e := range event.Effects {
				if e.Table == "activation_components" && e.Op == "insert" {
					row, err := currentReference(tx, e.Table, e.PK)
					if err != nil {
						return nil, err
					}
					result["component"] = map[string]any{"type": row["component_type"], "id": row["component_id"], "version": row["component_version"], "digest": row["content_digest"]}
				}
			}
		}
		if event.EventType == "activation.resource-observed" {
			for _, e := range event.Effects {
				if e.Table == "activation_resources" {
					put("resource", "ActivationResource", e.Table, e.PK)
				}
			}
		}
	}
	return result, failure
}
func currentPublicRow(tx *sql.Tx, table string, pk map[string]any) (map[string]any, error) {
	row, err := currentReference(tx, table, pk)
	if err != nil {
		return nil, err
	}
	result := map[string]any{}
	for column, value := range row {
		name := column
		if strings.HasSuffix(column, "_json") {
			name = strings.TrimSuffix(column, "_json")
			if value != nil {
				value, err = decodeCurrentJSON(value.(string))
				if err != nil {
					return nil, err
				}
			}
		}
		if column == "created_event_id" {
			name = "event_id"
		}
		if column == "terminal_status" {
			name = "outcome"
		}
		if currentEventTables[table].columns[column] == "b" {
			value = currentValuesEqual(value, 1) || value == true
		}
		result[name] = value
	}
	return result, nil
}
func currentProjectedValue(tx *sql.Tx, kind, table string, pk map[string]any) (map[string]any, error) {
	row, err := currentPublicRow(tx, table, pk)
	if err != nil {
		return nil, err
	}
	switch kind {
	case "Observation", "Checkpoint":
		return row["observation"].(map[string]any), nil
	case "Run":
		err = currentFillRun(tx, row)
	case "Activation":
		err = currentFillActivation(tx, row)
	case "Stop", "Attention":
		err = currentFillStopAttention(tx, kind, row)
	case "Receipt":
		err = currentFillReceipt(tx, row)
	case "Review":
		err = currentFillReview(tx, row)
	case "Finding":
		err = currentFillFinding(tx, row)
	}
	if err != nil {
		return nil, err
	}
	result := map[string]any{}
	for _, field := range strings.Fields(currentValueShapes[kind]) {
		name, _, _ := strings.Cut(field, ":")
		v, ok := row[name]
		if !ok {
			return nil, fmt.Errorf("missing projected %s.%s", kind, name)
		}
		result[name] = v
	}
	return result, nil
}
func currentQueryValues(tx *sql.Tx, query string, args ...any) ([]any, error) {
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	result := []any{}
	for rows.Next() {
		var value any
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, rows.Err()
}
func currentCollection(tx *sql.Tx, kind, table, pkColumn, query string, args ...any) ([]any, error) {
	ids, err := currentQueryValues(tx, query, args...)
	if err != nil {
		return nil, err
	}
	result := []any{}
	for _, id := range ids {
		value, err := currentProjectedValue(tx, kind, table, map[string]any{pkColumn: id})
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}
