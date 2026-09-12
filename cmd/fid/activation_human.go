package main

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func writeActivationStatusHuman(cmd *cobra.Command, activation workflow.Activation) error {
	lines := []string{"Activation " + activation.ID, "  Status: " + activation.Status,
		"  Kit: " + activation.KitIdentity, "  Backend: " + activation.BackendIdentity,
		"  This is an activation/install transaction, distinct from a Git commit."}
	terminal := activation.Status == "finalized" || activation.Status == "aborted"
	if terminal {
		lines = append(lines, "  Lifecycle is terminal; no activation recovery remains to select.")
	} else if activation.PointOfNoReturn {
		lines = append(lines, "  Recovery direction: complete-forward; irreversible commit",
			"  The activation commit has passed the point of no return. Recovery must finish installing the selected resources and finalize activation; returning to the prior state is unavailable.")
	} else {
		lines = append(lines, "  Recovery direction: restore-prior",
			"  The activation has not crossed its irreversible commit. Recovery returns resources to their state before this activation.")
	}
	if recovery := activation.Recovery; recovery != nil {
		lines = append(lines, "  Selected recovery: "+recovery.Direction, "  Reason: "+recovery.Reason)
		if recovery.CompletedEventID != nil {
			lines = append(lines, "  Recovery completion recorded: "+*recovery.CompletedEventID)
		}
	}
	lines = append(lines, "", "Resources (prior state → selected target)")
	for _, resource := range activation.Resources {
		prior := "absent (no resource was installed before this activation)"
		if resource.OldIdentity != nil {
			prior = *resource.OldIdentity
		}
		observed := "not observed"
		if resource.ObservedIdentity != nil {
			observed = *resource.ObservedIdentity
		}
		lines = append(lines, "  "+resource.Component, "    Prior: "+prior, "    Target: "+resource.NewIdentity,
			"    Observed: "+observed, "    State: "+resource.State)
	}
	if !terminal && !activation.PointOfNoReturn {
		lines = append(lines, "  Restoring an absent prior resource means leaving that resource absent, not restoring an unnamed version.")
	}
	_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(lines, "\n"))
	return err
}
