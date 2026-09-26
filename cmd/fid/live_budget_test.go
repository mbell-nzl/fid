package main

import (
	"strings"
	"testing"
)

func TestBudgetStopControlHumanNamesDurableIntervention(t *testing.T) {
	t.Parallel()
	got := renderLiveControl(map[string]any{"kind": "budget-stop", "stop_ids": []string{"STP-crossing"}, "attention_ids": []string{"ATT-crossing"}, "budget": map[string]any{"name": "attempts", "state": "exhausted", "remaining": 0}})
	for _, want := range []string{"budget-stop", "STP-crossing", "ATT-crossing"} {
		if !strings.Contains(got, want) {
			t.Fatalf("Human budget control = %q, missing %s", got, want)
		}
	}
}
