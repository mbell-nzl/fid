package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func TestRunAndGateStopErrorsUseStopUnresolvedEnvelope(t *testing.T) {
	cause := &workflow.StopBlockedError{Action: "gate.run", Stops: []workflow.Stop{{ID: "STOP-guard", Code: "operator-stop"}}}
	for name, render := range map[string]func(*cobra.Command, string, error) error{
		"run":  renderRunError,
		"gate": renderGateError,
	} {
		t.Run(name, func(t *testing.T) {
			var stderr bytes.Buffer
			root := &cobra.Command{}
			root.PersistentFlags().Bool("json", true, "")
			cmd := &cobra.Command{}
			root.AddCommand(cmd)
			cmd.SetErr(&stderr)
			if err := render(cmd, "gate.run", cause); err == nil {
				t.Fatal("render Stop error = nil, want nonzero failure")
			}
			if !strings.Contains(stderr.String(), `"code":"stop-unresolved"`) || !strings.Contains(stderr.String(), `"STOP-guard"`) {
				t.Fatalf("typed Stop failure = %q, want stop-unresolved STOP-guard", stderr.String())
			}
		})
	}
}

func TestReviewReworkStopErrorsRenderTypedBudgetDenial(t *testing.T) {
	consumed, remaining := int64(1), int64(0)
	cause := &workflow.StopBlockedError{
		Action:       "claim.acquire",
		Stops:        []workflow.Stop{{ID: "STP-review", Subject: workflow.Subject{Type: "envelope", ID: "ENV-review"}, Code: "review-rework-exhausted"}},
		AttentionIDs: []string{"ATT-review"},
		Budget:       &workflow.Budget{Name: "reviews", Scope: "envelope", Unit: "review-cycles", Maximum: 1, ExhaustionAction: "review-rework-exhausted", Consumed: &consumed, Remaining: &remaining, State: "exhausted"},
	}
	for name, render := range map[string]func(*cobra.Command, string, error) error{
		"claim": renderClaimError,
		"run":   renderRunError,
	} {
		t.Run(name, func(t *testing.T) {
			var stderr bytes.Buffer
			root := &cobra.Command{}
			root.PersistentFlags().Bool("json", true, "")
			cmd := &cobra.Command{}
			root.AddCommand(cmd)
			cmd.SetErr(&stderr)
			if err := render(cmd, cause.Action, cause); err == nil {
				t.Fatal("render review exhaustion = nil, want nonzero failure")
			}
			for _, want := range []string{`"code":"review-rework-exhausted"`, `"id":"ENV-review"`, `"unit":"review-cycles"`, `"remaining":0`, `"STP-review"`, `"ATT-review"`, `"event":null`, `"replayed":false`} {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("typed review exhaustion = %q, want %s", stderr.String(), want)
				}
			}
		})
	}
}
