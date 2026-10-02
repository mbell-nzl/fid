package workflow

import (
	"strings"
	"testing"
)

func TestBudgetPolicyValueStripsDerivedOutput(t *testing.T) {
	t.Parallel()
	consumed, remaining := int64(2), int64(3)
	value := budgetPolicyValue(Budget{Name: "attempts", Scope: "envelope", Unit: "attempts", Maximum: 5, ExhaustionAction: "budget-exhausted", Consumed: &consumed, Remaining: &remaining, State: "warning"})
	if value.Consumed != nil || value.Remaining != nil || value.State != "" {
		t.Fatalf("budgetPolicyValue() = %#v, want no derived output", value)
	}
	if raw := string(mustRunJSON(value)); strings.Contains(raw, "state") {
		t.Fatalf("policy JSON = %s, want no derived state", raw)
	}
}
