package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func writeStopHuman(cmd *cobra.Command, operation string, result map[string]any) error {
	lines := []string{"Purpose: " + stopHumanPurpose(operation), "", operation}
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
	lines := []string{"Stop " + stop.ID,
		"  Subject: " + stop.Subject.Type + ":" + stop.Subject.ID,
		"  Code: " + stop.Code, "  Blocked actions: " + strings.Join(stop.BlockedActions, ", ")}
	var details map[string]any
	_ = json.Unmarshal(stop.Details, &details)
	for _, member := range stopDetailMembers(stop.Code) {
		lines = append(lines, "  "+member.label+": "+stopDetailValue(details[member.key]))
	}
	lines = append(lines, stopConditionHumanLines(stop.ResolutionCondition)...)
	lines = append(lines, "  Resolver capability: "+stop.ResolverCapability)
	if stop.Resolution != nil {
		lines = append(lines, "  Status: resolved", "  Fact Event: "+stop.Resolution.FactEventID, "  Resolution Event: "+stop.Resolution.EventID)
	} else {
		lines = append(lines, "  Status: open", "  Next: "+stopConditionNext(stop.ResolutionCondition)+"; after the fact: fid stop resolve, then fid attention resolve where no lifecycle Event disposed it.")
	}
	return strings.Join(lines, "\n")
}

type stopDetailMember struct{ key, label string }

func stopDetailMembers(code string) []stopDetailMember {
	budget := []stopDetailMember{{"envelope_id", "Envelope"}, {"budget_name", "Budget"}, {"unit", "Unit"}, {"maximum", "Maximum"}, {"consumed", "Consumed"}, {"remaining", "Remaining"}}
	switch code {
	case "executor-liveness-unknown":
		return []stopDetailMember{{"run_id", "Run"}, {"executor_id", "Executor"}, {"observation_id", "Observation"}}
	case "recovery-state-conflict":
		return []stopDetailMember{{"run_id", "Run"}, {"observation_digest", "Observation digest"}, {"paths", "Paths"}, {"reasons", "Reasons"}}
	case "budget-exhausted", "budget-data-unavailable":
		return budget
	case "review-rework-exhausted":
		return append(budget, stopDetailMember{"last_review_id", "Last Review"})
	case "protected-decision-required":
		return []stopDetailMember{{"envelope_id", "Envelope"}, {"decision_class", "Decision class"}, {"required_capabilities", "Required capabilities"}, {"affected_actions", "Affected actions"}}
	case "human-outcome-required":
		return []stopDetailMember{{"envelope_id", "Envelope"}, {"obligation_id", "Obligation"}, {"outcome_class", "Outcome class"}}
	case "activation-incomplete":
		return []stopDetailMember{{"activation_id", "Activation"}, {"point_of_no_return", "Point of no return"}, {"components", "Components"}}
	case "capability-unavailable":
		return []stopDetailMember{{"capability_id", "Capability"}, {"available", "Available"}, {"owner_phase", "Owner phase"}}
	case "operator-stop":
		return []stopDetailMember{{"reason", "Reason"}, {"reference", "Reference"}}
	default:
		return nil
	}
}

func stopDetailValue(value any) string {
	if value == nil {
		return "unknown"
	}
	if values, ok := value.([]any); ok {
		if len(values) == 0 {
			return "none"
		}
		parts := make([]string, 0, len(values))
		for _, item := range values {
			parts = append(parts, stopDetailValue(item))
		}
		return strings.Join(parts, ", ")
	}
	return cli.EscapeControlText(fmt.Sprint(value))
}

func stopConditionHumanLines(condition workflow.Condition) []string {
	lines := []string{"  Resolution condition: " + condition.Type + " for " + condition.Subject.Type + ":" + condition.Subject.ID}
	members := []struct {
		label string
		value *string
	}{
		{"Condition budget", condition.BudgetName}, {"Condition capability", condition.CapabilityID},
		{"Condition Decision class", condition.DecisionClass}, {"Condition Outcome class", condition.OutcomeClass}, {"Condition Finding", condition.FindingID},
	}
	for _, member := range members {
		if member.value != nil {
			lines = append(lines, "  "+member.label+": "+cli.EscapeControlText(*member.value))
		}
	}
	return lines
}

func stopConditionNext(condition workflow.Condition) string {
	switch condition.Type {
	case "executor-terminated":
		return "fid run observe"
	case "run-recovered":
		return "fid run recover"
	case "activation-recovered":
		return "fid activation recover"
	case "capability-available":
		return "wait for the capability " + cli.EscapeControlText(*condition.CapabilityID)
	case "budget-observed":
		return "fid run observe-budget"
	case "work-ended":
		return "fid run close or fid run abandon"
	case "review-recorded", "finding-resolved":
		return "fid review record"
	case "human-outcome-recorded":
		return "fid human-outcome record"
	case "decision-satisfied":
		return "fid decision record-reference then fid human-outcome record"
	default:
		return "inspect the resolution condition"
	}
}

func attentionHumanDescription(cmd *cobra.Command, item workflow.Attention, known map[string]workflow.Stop) (string, error) {
	lines := []string{"Attention " + item.ID, "  Status: " + item.Status, "  Code: " + item.Code, "  Linked Stops: " + strings.Join(item.StopIDs, ", ")}
	if item.ResolutionEventID != nil {
		lines = append(lines, "  Resolution Event: "+*item.ResolutionEventID)
	}
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

func stopHumanPurpose(operation string) string {
	if strings.HasPrefix(operation, "attention.") {
		return "Inspect Attention and the dispositions of its linked Stops."
	}
	return "Inspect Stops, blocked actions and qualifying resolution conditions."
}
