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
