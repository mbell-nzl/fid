package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func writeStopHuman(cmd *cobra.Command, operation string, result map[string]any) error {
	lines := []string{operation}
	stops := map[string]workflow.Stop{}
	if stop, ok := result["stop"].(workflow.Stop); ok {
		stops[stop.ID] = stop
		lines = append(lines, stopHumanDescription(stop))
	}
	attention := []workflow.Attention{}
	if item, ok := result["attention"].(workflow.Attention); ok {
		attention = append(attention, item)
	}
	switch items := result["items"].(type) {
	case []workflow.Stop:
		if len(items) == 0 {
			lines = append(lines, "No Stops match the filters.")
		}
		for _, stop := range items {
			lines = append(lines, stopHumanDescription(stop))
		}
	case []workflow.Attention:
		if len(items) == 0 {
			lines = append(lines, "No Attention items match the filters.")
		}
		attention = append(attention, items...)
	}
	for _, item := range attention {
		description, err := attentionHumanDescription(cmd, item, stops)
		if err != nil {
			return err
		}
		lines = append(lines, description)
	}
	_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(lines, "\n"))
	return err
}

func stopHumanDescription(stop workflow.Stop) string {
	reason := stop.Code
	var details struct {
		Reason string `json:"reason"`
	}
	if json.Unmarshal(stop.Details, &details) == nil && details.Reason != "" {
		reason += ": " + details.Reason
	} else if len(stop.Details) > 0 {
		reason += ": " + string(stop.Details)
	}
	condition := stop.ResolutionCondition
	qualifier := fmt.Sprintf("%s for %s:%s", condition.Type, condition.Subject.Type, condition.Subject.ID)
	if condition.BudgetName != nil {
		qualifier += " budget " + *condition.BudgetName
	}
	if condition.CapabilityID != nil {
		qualifier += " capability " + *condition.CapabilityID
	}
	if stop.Resolution != nil {
		return fmt.Sprintf("Stop %s resolved: %s; qualifying resolution: %s; fact %s, disposition Event %s.", stop.ID, reason, qualifier, stop.Resolution.FactEventID, stop.Resolution.EventID)
	}
	return fmt.Sprintf("Stop %s blocks %s because %s; qualifying resolution: %s (resolver capability %s).", stop.ID, strings.Join(stop.BlockedActions, ", "), reason, qualifier, stop.ResolverCapability)
}

func attentionHumanDescription(cmd *cobra.Command, item workflow.Attention, known map[string]workflow.Stop) (string, error) {
	if item.ResolutionEventID != nil {
		return fmt.Sprintf("Attention %s for %s resolved by Attention disposition Event %s; linked Stops: %s.", item.ID, item.Code, *item.ResolutionEventID, strings.Join(item.StopIDs, ", ")), nil
	}
	lines := []string{fmt.Sprintf("Attention %s open for %s; qualifying resolution requires dispositions for every linked Stop:", item.ID, item.Code)}
	for _, id := range item.StopIDs {
		stop, ok := known[id]
		if !ok {
			workspace, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return "", err
			}
			stop, err = workflow.GetStop(workspace, issueMigrationsDir(), id)
			if err != nil {
				return "", fmt.Errorf("read Attention %s linked Stop: %w", item.ID, err)
			}
		}
		lines = append(lines, "  "+stopHumanDescription(stop))
	}
	return strings.Join(lines, "\n"), nil
}
