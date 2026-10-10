package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

// writeP23Human projects only fields already present in the matching JSON
// result, including recorded reasoning. Referenced file content stays in
// the project; views show only its reference metadata.
func writeP23Human(cmd *cobra.Command, operation string, items ...map[string]any) error {
	_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(p23HumanLines(operation, items...), "\n"))
	return err
}

func p23HumanLines(operation string, items ...map[string]any) []string {
	kind := strings.Split(operation, ".")[0]
	heading := p23HumanHeading(kind)
	if operation == "requirement.link" {
		heading = "Issue links"
	}
	lines := []string{"Purpose: " + p23HumanPurpose(operation), "", heading}
	if operation == "requirement.get" {
		if len(items) > 0 {
			lines = append(lines, p23HumanItem(kind, items[0], strings.HasSuffix(operation, ".list"))...)
			items = items[1:]
		}
		lines = append(lines, "", "Issue links")
	}
	lines = append(lines, judgmentSectionDisclaimer(kind, len(items))...)
	if len(items) == 0 {
		lines = append(lines, "  None.")
	}
	for _, item := range items {
		lines = append(lines, p23HumanItem(kind, item, strings.HasSuffix(operation, ".list"))...)
	}
	return lines
}

func writeP23MutationHuman(cmd *cobra.Command, operation string, item map[string]any, budget *workflow.Budget, protected []workflow.ProtectedDecision, stops []workflow.Stop, attention []workflow.Attention) error {
	lines := p23HumanLines(operation, item)
	if budget != nil {
		lines = append(lines, "", "Review budget")
		lines = append(lines, budgetHumanLines([]workflow.Budget{*budget})...)
	}
	if len(protected) != 0 {
		lines = append(lines, "", "Protected Decisions")
		for _, decision := range protected {
			lines = append(lines, fmt.Sprintf("  Class %s: %s; affected actions: %s", decision.DecisionClass, decision.Status, strings.Join(decision.AffectedActions, ", ")))
			if decision.Status != "satisfied" {
				lines = append(lines, "    Resolution requires an approved Decision Reference with the required capabilities and an approved Human Outcome of the same class.")
			}
		}
	}
	if len(stops) != 0 {
		lines = append(lines, "", "Stops")
		for _, stop := range stops {
			lines = append(lines, stopHumanDescription(stop))
		}
	}
	if len(attention) != 0 {
		lines = append(lines, "", "Attention")
		for _, item := range attention {
			lines = append(lines, fmt.Sprintf("  Attention %s: %s; linked Stops: %s", item.ID, item.Code, strings.Join(item.StopIDs, ", ")))
		}
	}
	_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(lines, "\n"))
	return err
}

func p23HumanPurpose(operation string) string {
	switch operation {
	case "review.record":
		return "Record an independent Review of a Run and its bound Issue."
	case "review.get", "review.list":
		return "Inspect independent Reviews and their Findings."
	case "finding.address":
		return "Record a Builder assertion that a Finding was addressed."
	case "finding.list":
		return "Inspect Finding dispositions and current blocking state."
	case "concern.record", "concern.dispose", "concern.get", "concern.list":
		return "Inspect a raised concern and its disposition as judgment."
	case "deviation.record", "deviation.get", "deviation.list":
		return "Inspect a recorded deviation as judgment."
	case "decision.record-reference", "decision.list":
		return "Inspect Decision References and their authority."
	case "human-outcome.record", "human-outcome.get", "human-outcome.list":
		return "Inspect authorized Human judgments and their provenance."
	case "requirement.record-reference", "requirement.link", "requirement.get", "requirement.list":
		return "Inspect Requirement References and Issue links."
	default:
		return "Inspect recorded judgment state."
	}
}

func p23HumanHeading(kind string) string {
	switch kind {
	case "review":
		return "Reviews"
	case "finding":
		return "Findings"
	case "concern":
		return "Concerns"
	case "deviation":
		return "Deviations"
	case "decision":
		return "Decision References"
	case "human-outcome":
		return "Human Outcomes"
	case "requirement":
		return "Requirement References"
	default:
		return "Items"
	}
}

func p23HumanItem(kind string, item map[string]any, list bool) []string {
	id := item["id"]
	if kind == "requirement" {
		if id == nil {
			return []string{fmt.Sprintf("  %v %v %v", item["requirement_id"], item["relation"], item["issue_id"])}
		}
		lines := []string{fmt.Sprintf("  Requirement %v", item["requirement_id"])}
		line := fmt.Sprintf("    Reference ID: %v; Status: %v", id, item["status"])
		if item["superseded_by_id"] != nil {
			line += fmt.Sprintf("; Superseded by: %v", item["superseded_by_id"])
		}
		return append(append(lines, line), p23HumanSource(item["source"])...)
	}
	lines := []string{fmt.Sprintf("  %s %v", p23HumanName(kind), id)}
	switch kind {
	case "review":
		lines = append(lines, fmt.Sprintf("    Issue: %v; Envelope: %v revision %v", item["issue_id"], item["envelope_id"], item["envelope_revision"]),
			fmt.Sprintf("    Run: %v; Reviewer: %v; Profile: %v", item["run_id"], item["actor_id"], item["profile"]),
			fmt.Sprintf("    Verdict: %v; Independence: %v", item["verdict"], item["independence"]))
		repo := p23HumanMap(item["repository"])
		head, _ := repo["head_oid"].(string)
		dirty, _ := repo["dirty"].(bool)
		lines = append(lines, "    Repository: "+repositoryHumanSummary(head, dirty))
		counts := map[string]int{}
		for _, finding := range p23HumanSlice(item["findings"]) {
			counts[fmt.Sprint(p23HumanMap(finding)["severity"])]++
		}
		lines = append(lines, fmt.Sprintf("    Findings by severity: blocker %d; major %d; minor %d; nit %d", counts["blocker"], counts["major"], counts["minor"], counts["nit"]))
		for _, finding := range p23HumanSlice(item["findings"]) {
			lines = append(lines, p23HumanItem("finding", p23HumanMap(finding), list)...)
		}
	case "finding":
		lines = append(lines, fmt.Sprintf("    Review: %v; Severity: %v; Blocking: %v; Disposition: %v", item["review_id"], item["severity"], item["blocking"], item["state"]))
		if item["state"] == "addressed" {
			lines = append(lines, "    Addressed is a Builder assertion, not a resolution.")
		}
	case "concern":
		lines = append(lines, "    Subject: "+p23HumanSubject(item["subject"]), fmt.Sprintf("    Severity: %v; Status: %v; Raised by: %v", item["severity"], item["status"], item["actor_id"]))
		if disposition := p23HumanMap(item["disposition"]); disposition != nil {
			lines = append(lines, fmt.Sprintf("    Disposition: %v; Disposed by: %v", disposition["action"], disposition["actor_id"]))
			if disposition["action"] == "carried-forward" {
				lines = append(lines, fmt.Sprintf("    Work now owned by Issue: %v", disposition["issue_id"]))
			}
		}
		lines = append(lines, p23HumanSources(item["references"])...)
	case "deviation":
		lines = append(lines, "    Subject: "+p23HumanSubject(item["subject"]), fmt.Sprintf("    Disposition: %v; Authority Actor: %v", item["disposition"], p23HumanMap(item["authority"])["actor_id"]))
		lines = append(lines, p23HumanSources(item["references"])...)
	case "decision":
		lines = append(lines, fmt.Sprintf("    Class: %v; Outcome: %v; Authority Actor: %v", item["class"], item["outcome"], p23HumanMap(item["authority"])["actor_id"]))
		lines = append(lines, "    Required capabilities: "+p23HumanValues(p23HumanMap(item["authority"])["required_capabilities"]))
		lines = append(lines, p23HumanSource(item["source"])...)
	case "human-outcome":
		provenance := p23HumanMap(item["provenance"])
		lines = append(lines, fmt.Sprintf("    Actor: %v; Class: %v; Result: %v", item["actor_id"], item["class"], item["result"]),
			fmt.Sprintf("    Activation: %v revision %v; Profile: %v", provenance["activation_id"], provenance["activation_version"], provenance["profile_ref"]))
		lines = append(lines, p23HumanSources(item["references"])...)
	}
	return append(lines, p23RecordedReasoning(kind, item, list)...)
}

func p23HumanValues(value any) string {
	values := p23HumanSlice(value)
	text := make([]string, 0, len(values))
	for _, item := range values {
		text = append(text, fmt.Sprint(item))
	}
	return strings.Join(text, ", ")
}

func p23HumanName(kind string) string {
	switch kind {
	case "review":
		return "Review"
	case "finding":
		return "Finding"
	case "concern":
		return "Concern"
	case "deviation":
		return "Deviation"
	case "decision":
		return "Decision Reference"
	case "human-outcome":
		return "Human Outcome"
	case "requirement":
		return "Requirement"
	default:
		return "Item"
	}
}

func p23HumanMap(value any) map[string]any {
	if value == nil {
		return nil
	}
	if mapped, ok := value.(map[string]any); ok {
		return mapped
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var mapped map[string]any
	if json.Unmarshal(data, &mapped) != nil {
		return nil
	}
	return mapped
}

func p23HumanSlice(value any) []any {
	if values, ok := value.([]any); ok {
		return values
	}
	if value == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var values []any
	if json.Unmarshal(data, &values) != nil {
		return nil
	}
	return values
}

func p23HumanSubject(value any) string {
	subject := p23HumanMap(value)
	return fmt.Sprintf("%v:%v", subject["type"], subject["id"])
}

func p23HumanSources(value any) []string {
	lines := []string{}
	for _, source := range p23HumanSlice(value) {
		lines = append(lines, p23HumanSource(source)...)
	}
	return lines
}

func p23HumanSource(value any) []string {
	source := p23HumanMap(value)
	if source == nil {
		return nil
	}
	return []string{fmt.Sprintf("    Reference: %v; path %v; digest %v", source["display_label"], source["path"], source["digest"])}
}

func recordedReasoningHumanLines(label string, value any, list bool) []string {
	text, ok := value.(string)
	if !ok {
		return nil
	}
	if list {
		text = strings.ReplaceAll(text, "\t", " ")
	}
	parts := strings.Split(text, "\n")
	for i, part := range parts {
		parts[i] = cli.EscapeControlText(part)
	}
	text = strings.Join(parts, "\n")
	if list {
		text = strings.Join(strings.Fields(text), " ")
		scalars := []rune(text)
		if len(scalars) > 80 {
			text = string(scalars[:80]) + "…"
		}
	}
	return strings.Split("    "+label+": "+strings.ReplaceAll(text, "\n", "\n      "), "\n")
}

// judgmentSectionDisclaimer is shared by leaf views and resume sections.
func judgmentSectionDisclaimer(kind string, count int) []string {
	if count == 0 {
		return nil
	}
	switch kind {
	case "concern", "deviation":
		return []string{"  These records express judgment, not deterministic proof."}
	case "human-outcome":
		return []string{"  These Outcomes express authorized Human judgment, not deterministic proof."}
	default:
		return nil
	}
}

func p23RecordedReasoning(kind string, item map[string]any, list bool) []string {
	lines := []string{}
	switch kind {
	case "concern":
		lines = append(lines, recordedReasoningHumanLines("Summary", item["summary"], list)...)
		lines = append(lines, recordedReasoningHumanLines("Reasoning", item["reasoning"], list)...)
		lines = append(lines, recordedReasoningHumanLines("Disposition reason", p23HumanMap(item["disposition"])["reason"], list)...)
	case "finding":
		lines = append(lines, recordedReasoningHumanLines("Summary", item["summary"], list)...)
		for _, disposition := range p23HumanSlice(item["dispositions"]) {
			lines = append(lines, recordedReasoningHumanLines("Disposition reason", p23HumanMap(disposition)["reason"], list)...)
		}
	case "deviation":
		lines = append(lines, recordedReasoningHumanLines("Expected", item["expected"], list)...)
		lines = append(lines, recordedReasoningHumanLines("Observed", item["observed"], list)...)
		lines = append(lines, recordedReasoningHumanLines("Reasoning", item["reasoning"], list)...)
	case "human-outcome":
		lines = append(lines, recordedReasoningHumanLines("Reasoning", item["reasoning"], list)...)
	}
	return lines
}
