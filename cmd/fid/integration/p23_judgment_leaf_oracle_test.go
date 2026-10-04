//go:build integration

package integration

import (
	"reflect"
	"testing"
)

// These three shapes are the closed Deviation, DecisionRef, and HumanOutcome
// values in public-interface.md section 3. Requests supply every prose value.
func p23JudgmentModel(t *testing.T, pairs map[string]p23HumanLeafPair) map[string]map[string]any {
	t.Helper()
	deviationPair := pairs["deviation.record"]
	deviation := p23LeafRequest(t, deviationPair)
	deviationValue := map[string]any{
		"id": deviation["id"], "subject": deviation["subject"], "actor_id": p23LeafFlag(t, deviationPair, "--actor"),
		"expected": deviation["expected"], "observed": deviation["observed"], "authority": deviation["authority"],
		"disposition": deviation["disposition"], "reasoning": deviation["reasoning"], "references": p24ModelSources(t, deviationPair, deviation["references"]),
		"supersedes_deviation_id": deviation["supersedes_deviation_id"], "created_at": "<generated:timestamp>", "event_id": "<generated:EventRef>",
	}
	decisionPair := pairs["decision.record-reference"]
	decision := p23LeafRequest(t, decisionPair)
	decisionValue := map[string]any{
		"id": decision["id"], "class": decision["class"], "source": p24ModelSource(t, decisionPair, decision["source"]),
		"authority": decision["authority"], "outcome": decision["outcome"], "scope": decision["scope"],
		"created_at": "<generated:timestamp>", "event_id": "<generated:EventRef>",
	}
	outcomePair := pairs["human-outcome.record"]
	outcome := p23LeafRequest(t, outcomePair)
	outcomeValue := map[string]any{
		"id": outcome["id"], "actor_id": p23LeafFlag(t, outcomePair, "--actor"),
		"provenance": map[string]any{"activation_id": "<generated:ACT>", "activation_version": float64(10), "profile_ref": "activation:p22-corpus"},
		"class":      outcome["class"], "subject": outcome["subject"], "result": outcome["result"],
		"scope": outcome["scope"], "reasoning": outcome["reasoning"], "references": p24ModelSources(t, outcomePair, outcome["references"]),
		"created_at": "<generated:timestamp>", "event_id": "<generated:EventRef>",
	}
	return map[string]map[string]any{
		"deviation.record":          {"deviation": deviationValue, "idempotency_key": p23LeafFlag(t, deviationPair, "--idempotency-key"), "replayed": true},
		"deviation.list":            {"items": []any{deviationValue}},
		"decision.record-reference": {"decision": decisionValue, "protected_decisions": []any{}, "stops": []any{}, "attention": []any{}, "idempotency_key": p23LeafFlag(t, decisionPair, "--idempotency-key"), "replayed": true},
		"decision.list":             {"items": []any{decisionValue}},
		"human-outcome.record":      {"outcome": outcomeValue, "protected_decisions": []any{}, "stops": []any{}, "attention": []any{}, "idempotency_key": p23LeafFlag(t, outcomePair, "--idempotency-key"), "replayed": true},
		"human-outcome.list":        {"items": []any{outcomeValue}},
	}
}

func p23CheckJudgmentIdentity(t *testing.T, accepted map[string]p23LeafOutcome) {
	t.Helper()
	for _, row := range []struct{ mutation, list, key string }{
		{"deviation.record", "deviation.list", "deviation"},
		{"decision.record-reference", "decision.list", "decision"},
		{"human-outcome.record", "human-outcome.list", "outcome"},
	} {
		mutation, ok := accepted[row.mutation]
		if !ok {
			continue // A -run subtest selection may omit this family.
		}
		if _, ok := accepted[row.list]; !ok {
			continue
		}
		items := accepted[row.list].Semantics["items"].([]any)
		if len(items) != 1 || !reflect.DeepEqual(mutation.Semantics[row.key], items[0]) {
			t.Fatalf("%s/list identity drift: %#v %#v", row.mutation, mutation.Semantics[row.key], items)
		}
	}
}
