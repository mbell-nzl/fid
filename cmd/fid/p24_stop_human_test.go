package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func TestP24StopHumanBudgetLabel(t *testing.T) {
	t.Parallel()
	stop := workflow.Stop{ID: "STP-budget", Subject: workflow.Subject{Type: "run", ID: "RUN-one"}, Code: "budget-exhausted",
		BlockedActions: []string{"gate.run"}, Details: json.RawMessage(`{"envelope_id":"ENV-one","budget_name":"tokens","unit":"tokens","maximum":10,"consumed":10,"remaining":0}`),
		ResolutionCondition: workflow.Condition{Type: "budget-observed", Subject: workflow.Subject{Type: "run", ID: "RUN-one"}, BudgetName: stringPointer("tokens")}, ResolverCapability: "stop.resolve"}
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)
	if err := writeStopHuman(cmd, "stop.list", map[string]any{"items": []workflow.Stop{stop}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Budget: tokens") || strings.Contains(out.String(), "{") {
		t.Fatalf("expected labelled Budget, no JSON: %s", out.String())
	}
}
