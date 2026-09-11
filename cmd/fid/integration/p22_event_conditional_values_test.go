//go:build integration

package integration

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func p22RejectMissingConditionalFacts(t *testing.T, conn *sql.DB, event journal.Event) {
	t.Helper()
	mutations := map[string]func(*journal.Event){}
	if event.EventType == "run.recovered" {
		for _, field := range []struct{ name, table, column, result, key, value string }{
			{"successor-claim-worktree", "claims", "worktree_id", "claim", "worktree_id", "WT-forged"},
			{"successor-start-time", "runs", "started_at", "run", "started_at", "2026-01-01T00:00:00Z"},
		} {
			mutations[field.name] = func(e *journal.Event) {
				for i := range e.Effects {
					effect := &e.Effects[i]
					if effect.Table == field.table && effect.Op == "insert" {
						effect.Values[field.column] = field.value
					}
					if effect.Table == "idempotency_results" {
						var result map[string]any
						p22FaultOK(t, json.Unmarshal([]byte(effect.Values["result_json"].(string)), &result))
						result[field.result].(map[string]any)[field.key] = field.value
						raw, err := json.Marshal(result)
						p22FaultOK(t, err)
						effect.Values["result_json"] = string(raw)
					}
				}
			}
		}
		mutations["immutable-successor-claim-scope"] = func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table == "claims" && e.Effects[i].Op == "insert" {
					e.Effects[i].Values["scope_json"] = `{"paths":["forged.txt"]}`
				}
				if e.Effects[i].Table == "idempotency_results" {
					var result map[string]any
					p22FaultOK(t, json.Unmarshal([]byte(e.Effects[i].Values["result_json"].(string)), &result))
					result["claim"].(map[string]any)["scope"] = map[string]any{"paths": []string{"forged.txt"}}
					raw, err := json.Marshal(result)
					p22FaultOK(t, err)
					e.Effects[i].Values["result_json"] = string(raw)
				}
			}
		}
	}
	for i, effect := range event.Effects {
		if effect.Table == "budget_samples" && (event.EventType == "gate.execution-started" || event.EventType == "gate.receipt-recorded") {
			mutations["required-fid-charge"] = func(e *journal.Event) { e.Effects = append(e.Effects[:i], e.Effects[i+1:]...) }
		}
		if effect.Table == "stops" && (event.EventType == "run.budget-observed" || event.EventType == "gate.execution-started" || event.EventType == "gate.receipt-recorded") {
			mutations["required-budget-boundary"] = func(e *journal.Event) {
				kept := []journal.Effect{}
				for _, item := range e.Effects {
					if item.Table == "stops" || item.Table == "attention_items" || item.Table == "attention_stop_links" {
						continue
					}
					if item.Table == "idempotency_results" && e.EventType == "run.budget-observed" {
						var result map[string]any
						p22FaultOK(t, json.Unmarshal([]byte(item.Values["result_json"].(string)), &result))
						result["stops"], result["attention"] = []any{}, []any{}
						raw, err := json.Marshal(result)
						p22FaultOK(t, err)
						item.Values["result_json"] = string(raw)
					}
					kept = append(kept, item)
				}
				e.Effects = kept
			}
		}
		if effect.Table == "activation_components" && effect.Op == "update" {
			mutations["required-complete-component-transition"] = func(e *journal.Event) { e.Effects = append(e.Effects[:i], e.Effects[i+1:]...) }
		}
		if effect.Table == "activation_components" && effect.Op == "insert" {
			mutations["component-evidence-digest"] = func(e *journal.Event) {
				var values []map[string]any
				p22FaultOK(t, json.Unmarshal([]byte(e.Effects[i].Values["evidence_json"].(string)), &values))
				if _, found := values[0]["producer"]; found {
					values[0]["producer"] = "forged-producer"
				} else {
					values[0]["timeout_ms"] = float64(599999)
				}
				raw, err := json.Marshal(values)
				p22FaultOK(t, err)
				e.Effects[i].Values["evidence_json"] = string(raw)
			}
		}
		if effect.Table == "activation_gate_definitions" {
			mutations["definition-digest"] = func(e *journal.Event) { e.Effects[i].Values["definition_digest"] = "sha256:" + p22Zeros(64) }
		}
		if effect.Table == "kit_activations" && effect.Op == "update" && effect.PK["id"] != event.Aggregate.ID {
			mutations["required-predecessor-supersession"] = func(e *journal.Event) { e.Effects = append(e.Effects[:i], e.Effects[i+1:]...) }
		}
		if effect.Table == "stop_resolutions" && event.EventType == "activation.recovery-completed" {
			mutations["required-complete-Stop-disposition"] = func(e *journal.Event) { e.Effects = append(e.Effects[:i], e.Effects[i+1:]...) }
		}
		if effect.Table == "run_budget_bindings" && event.EventType == "run.recovered" {
			mutations["immutable-successor-budget-policy"] = func(e *journal.Event) { p22ChangeInheritedBudget(t, e, i) }
		}
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			raw, err := json.Marshal(event)
			p22FaultOK(t, err)
			var changed journal.Event
			p22FaultOK(t, json.Unmarshal(raw, &changed))
			mutate(&changed)
			if err := journal.DryRunProjectEvent(conn, changed); err == nil {
				t.Fatal("incomplete or altered conditional transaction accepted")
			}
		})
	}
}

// Keep the changed result consistent with the altered policy. Rejection must
// establish inheritance from the recorded predecessor, not merely compare the
// successor's two equally falsified copies.
func p22ChangeInheritedBudget(t *testing.T, event *journal.Event, index int) {
	t.Helper()
	var policy map[string]any
	p22FaultOK(t, json.Unmarshal([]byte(event.Effects[index].Values["budget_json"].(string)), &policy))
	policy["maximum"] = policy["maximum"].(float64) + 1
	raw, err := json.Marshal(policy)
	p22FaultOK(t, err)
	event.Effects[index].Values["budget_json"] = string(raw)
	for i, effect := range event.Effects {
		if effect.Table != "idempotency_results" {
			continue
		}
		var result map[string]any
		p22FaultOK(t, json.Unmarshal([]byte(effect.Values["result_json"].(string)), &result))
		for _, value := range result["run"].(map[string]any)["budgets"].([]any) {
			budget := value.(map[string]any)
			if budget["name"] != policy["name"] {
				continue
			}
			budget["maximum"] = policy["maximum"]
			if budget["consumed"] != nil {
				remaining := budget["maximum"].(float64) - budget["consumed"].(float64)
				if remaining < 0 {
					remaining = 0
				}
				budget["remaining"] = remaining
				state := "available"
				if remaining == 0 {
					state = "exhausted"
				} else if budget["warning_threshold"] != nil && remaining <= budget["warning_threshold"].(float64) {
					state = "warning"
				}
				budget["state"] = state
			}
		}
		raw, err := json.Marshal(result)
		p22FaultOK(t, err)
		event.Effects[i].Values["result_json"] = string(raw)
	}
}
