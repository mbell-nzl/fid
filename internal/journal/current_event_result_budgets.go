package journal

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"time"
)

func currentProjectedBudget(tx *sql.Tx, runID, name any, excluding ...string) (map[string]any, error) {
	binding, err := currentPublicRow(tx, "run_budget_bindings", map[string]any{"run_id": runID, "name": name})
	if err != nil {
		return nil, err
	}
	budget := binding["budget"].(map[string]any)
	run, err := currentReference(tx, "runs", map[string]any{"id": runID})
	if err != nil {
		return nil, err
	}
	known := currentEnum(budget["unit"], "attempts", "failures", "milliseconds", "tool-failures", "interventions", "review-cycles")
	consumed := int64(0)
	excluded := ""
	if len(excluding) != 0 {
		excluded = excluding[0]
	}
	rows, err := tx.Query("SELECT sample_kind,value FROM budget_samples WHERE envelope_id=? AND envelope_revision=? AND budget_name=? AND created_event_id<>? ORDER BY observed_at,id", run["envelope_id"], run["envelope_revision"], name, excluded)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var kind string
		var value sql.NullInt64
		if err := rows.Scan(&kind, &value); err != nil {
			return nil, err
		}
		if !value.Valid {
			known = false
			continue
		}
		if kind == "baseline" {
			if known && value.Int64 < consumed {
				return nil, fmt.Errorf("budget baseline decreases")
			}
			known, consumed = true, value.Int64
			continue
		}
		if !known {
			continue
		}
		if value.Int64 < 0 || value.Int64 > math.MaxInt64-consumed {
			return nil, fmt.Errorf("budget consumption overflow")
		}
		consumed += value.Int64
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	budget["consumed"], budget["remaining"], budget["state"] = nil, nil, "unknown"
	if known {
		maximum, _ := currentIntegerValue(budget["maximum"])
		remaining := maximum - consumed
		if remaining < 0 {
			remaining = 0
		}
		state := "available"
		if remaining == 0 {
			state = "exhausted"
		} else if budget["warning_threshold"] != nil {
			warning, _ := currentIntegerValue(budget["warning_threshold"])
			if remaining <= warning {
				state = "warning"
			}
		}
		budget["consumed"], budget["remaining"], budget["state"] = consumed, remaining, state
	}
	return budget, nil
}
func currentEventIDs(event Event, table, key string) []string {
	seen := map[string]bool{}
	ids := []string{}
	for _, effect := range event.Effects {
		if effect.Table == table {
			v, _ := effect.PK[key].(string)
			if !seen[v] {
				seen[v] = true
				ids = append(ids, v)
			}
		}
	}
	sort.Strings(ids)
	return ids
}
func currentEventObjects(tx *sql.Tx, event Event, table, kind, key string) ([]any, error) {
	result := []any{}
	for _, id := range currentEventIDs(event, table, key) {
		v, err := currentProjectedValue(tx, kind, table, map[string]any{key: id})
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}
func currentBudgetEventResult(tx *sql.Tx, event Event) (map[string]any, error) {
	result := map[string]any{}
	if event.EventType == "run.execution-blocked" {
		stop := currentInserted(event, "stops")
		v, err := decodeCurrentJSON(stop["details_json"].(string))
		if err != nil {
			return nil, err
		}
		budget, err := currentProjectedBudget(tx, event.Aggregate.ID, v.(map[string]any)["budget_name"])
		if err != nil {
			return nil, err
		}
		result["denial"] = map[string]any{"budget": budget, "stop_ids": currentEventIDs(event, "stops", "id"), "attention_ids": currentEventIDs(event, "attention_stop_links", "attention_id"), "event_id": event.EventID, "event_occurred_at": event.OccurredAt.UTC().Format(time.RFC3339Nano)}
		return result, nil
	}
	sample := currentInserted(event, "budget_samples")
	value, err := currentProjectedValue(tx, "BudgetSample", "budget_samples", map[string]any{"id": sample["id"]})
	if err != nil {
		return nil, err
	}
	result["sample"] = value
	budget, err := currentProjectedBudget(tx, event.Aggregate.ID, sample["budget_name"])
	if err != nil {
		return nil, err
	}
	result["budget"] = budget
	stops, err := currentEventObjects(tx, event, "stops", "Stop", "id")
	if err != nil {
		return nil, err
	}
	result["stops"] = stops
	attention := []any{}
	for _, id := range currentEventIDs(event, "attention_stop_links", "attention_id") {
		v, err := currentProjectedValue(tx, "Attention", "attention_items", map[string]any{"id": id})
		if err != nil {
			return nil, err
		}
		attention = append(attention, v)
	}
	result["attention"] = attention
	return result, nil
}

func currentReviewBudgetResult(tx *sql.Tx, event Event) (map[string]any, error) {
	result := map[string]any{"budget": nil, "stops": []any{}, "attention": []any{}}
	sample := currentInserted(event, "budget_samples")
	if sample == nil {
		return result, nil
	}
	budget, err := currentProjectedBudget(tx, sample["run_id"], sample["budget_name"])
	if err != nil {
		return nil, err
	}
	result["budget"] = budget
	result["stops"], err = currentEventObjects(tx, event, "stops", "Stop", "id")
	if err != nil {
		return nil, err
	}
	attention := []any{}
	for _, id := range currentEventIDs(event, "attention_stop_links", "attention_id") {
		value, err := currentProjectedValue(tx, "Attention", "attention_items", map[string]any{"id": id})
		if err != nil {
			return nil, err
		}
		attention = append(attention, value)
	}
	result["attention"] = attention
	return result, nil
}
