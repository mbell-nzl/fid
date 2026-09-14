//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestP22ForceRecoveryBudgetConditionalEvents(t *testing.T) {
	for _, maximum := range []int{1, 2} {
		a := p22ForceBudgetFixture(t, maximum)
		p22AcceptedCall(t, a, "force-event", "run.recover", p22ForceBudgetInput(t, a, "force-event"))
		prefix, control := p22CurrentEventPrefix(t, a.workspace, "run.recovered")
		check := p22ForceEventChecker(t, a.binary, a.workspace, prefix)
		cases := []struct {
			name   string
			mutate func(*journal.Event)
		}{
			{"valid", nil},
			{"missing-charge-self-consistent", func(e *journal.Event) {
				p22RemoveForceEffects(e, "budget_samples", "stops", "attention_items", "attention_stop_links")
				p22ForceCachedBudget(e, 0, maximum, "available")
			}},
			{"wrong-value-self-consistent", func(e *journal.Event) {
				p22SetEventValue("budget_samples", "value", 0)(e)
				p22RemoveForceEffects(e, "stops", "attention_items", "attention_stop_links")
				p22ForceCachedBudget(e, 0, maximum, "available")
			}},
			{"wrong-source", p22SetEventValue("budget_samples", "source", "external")},
			{"wrong-source-id", p22SetEventValue("budget_samples", "source_sample_id", "REC-aaaaaa")},
			{"wrong-run", p22SetEventValue("budget_samples", "run_id", control.Aggregate.ID)},
			{"wrong-actor", p22SetEventValue("budget_samples", "actor_id", "actor-denied")},
			{"wrong-time", p22SetEventValue("budget_samples", "observed_at", "2026-01-01T00:00:00Z")},
			{"wrong-kind", p22SetEventValue("budget_samples", "sample_kind", "baseline")},
			{"wrong-predecessor", p22SetEventValue("budget_samples", "previous_sample_id", "BSM-aaaaaa")},
			{"measurement", p22SetEventValue("budget_samples", "measurement_ref_json", "{}")},
			{"unrelated-budget", p22SetEventValue("budget_samples", "budget_name", "attempts")},
			{"duplicate-charge", func(e *journal.Event) {
				for i, effect := range e.Effects {
					if effect.Table == "budget_samples" {
						raw, _ := json.Marshal(effect)
						var duplicate journal.Effect
						p22FaultOK(t, json.Unmarshal(raw, &duplicate))
						duplicate.PK["id"] = "BSM-aaaaaa"
						duplicate.Values["id"] = "BSM-aaaaaa"
						e.Effects = append(e.Effects[:i+1], append([]journal.Effect{duplicate}, e.Effects[i+1:]...)...)
						return
					}
				}
			}},
		}
		if maximum == 1 {
			cases = append(cases, struct {
				name   string
				mutate func(*journal.Event)
			}{"missing-boundary", func(e *journal.Event) { p22RemoveForceEffects(e, "stops", "attention_items", "attention_stop_links") }}, struct {
				name   string
				mutate func(*journal.Event)
			}{"wrong-stop-subject", p22SetEventValue("stops", "subject_id", control.Aggregate.ID)})
		}
		if maximum == 1 {
			cases = append(cases, struct {
				name   string
				mutate func(*journal.Event)
			}{"extra-boundary", func(e *journal.Event) {
				for i, effect := range e.Effects {
					if effect.Table == "attention_stop_links" {
						var stop, link journal.Effect
						for _, candidate := range e.Effects {
							if candidate.Table == "stops" {
								raw, err := json.Marshal(candidate)
								p22FaultOK(t, err)
								p22FaultOK(t, json.Unmarshal(raw, &stop))
								break
							}
						}
						raw, err := json.Marshal(effect)
						p22FaultOK(t, err)
						p22FaultOK(t, json.Unmarshal(raw, &link))
						stop.PK["id"] = "STP-aaaaaa"
						stop.Values["id"] = "STP-aaaaaa"
						link.PK["stop_id"] = "STP-aaaaaa"
						e.Effects = append(e.Effects[:i+1], append([]journal.Effect{stop, link}, e.Effects[i+1:]...)...)
						return
					}
				}
			}})
		}
		raw, err := json.Marshal(control)
		p22FaultOK(t, err)
		for _, test := range cases {
			t.Run(fmt.Sprint(maximum)+"/"+test.name, func(t *testing.T) {
				var event journal.Event
				p22FaultOK(t, json.Unmarshal(raw, &event))
				if test.mutate != nil {
					test.mutate(&event)
					changed, err := json.Marshal(event)
					p22FaultOK(t, err)
					if bytes.Equal(raw, changed) {
						t.Fatal("mutation had no effect")
					}
				}
				for _, mode := range []string{"live", "replay", "rebuild"} {
					t.Run(mode, func(t *testing.T) {
						check(t, event, mode, test.mutate == nil)
					})
				}
			})
		}
	}
}

func p22RemoveForceEffects(e *journal.Event, tables ...string) {
	kept := []journal.Effect{}
	for _, effect := range e.Effects {
		remove := false
		for _, table := range tables {
			if effect.Table == table {
				remove = true
			}
		}
		if !remove {
			kept = append(kept, effect)
		}
	}
	e.Effects = kept
}

func p22ForceCachedBudget(e *journal.Event, consumed, remaining int, state string) {
	p22MutateEventJSON("idempotency_results", "result_json", func(v map[string]any) {
		for _, raw := range v["run"].(map[string]any)["budgets"].([]any) {
			b := raw.(map[string]any)
			if b["name"] == "interventions" {
				b["consumed"] = consumed
				b["remaining"] = remaining
				b["state"] = state
			}
		}
	})(e)
}
