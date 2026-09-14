package workflow

import "strings"

const (
	RuleAllACsVerified              = "all_acs_verified"
	RuleHasCompletionReport         = "has_completion_report"
	RuleConcernsAddressed           = "concerns_addressed"
	RuleAllChildTasksDoneOrCanceled = "all_child_tasks_done_or_cancelled"
	// RuleRequiredGatesCurrent is reserved for the v0.2b trusted-gate path.
	// It has no evaluator in this phase.
	RuleRequiredGatesCurrent = "required_gates_current"
	// RuleHasCompletion is reserved for structured completion requirements.
	// It has no evaluator in this phase.
	RuleHasCompletion = "has_completion"
	// RuleLinkedToDecision remains reserved for project Decision scope (Phase 23).
	// It has no evaluator in this phase.
	RuleLinkedToDecision = "linked_to_decision"
)

// IsImplementedRule reports whether a rule has an evaluator in this release.
func IsImplementedRule(name string) bool {
	return name == RuleAllACsVerified ||
		name == RuleHasCompletionReport ||
		name == RuleConcernsAddressed ||
		name == RuleAllChildTasksDoneOrCanceled
}

type RuleResult struct {
	Satisfied bool
	Unmet     []UnmetCondition
}

type UnmetCondition struct {
	Rule         string
	EntityIDs    []string
	MissingRoles []string
}

type AcceptanceCriterionVerification struct {
	AcceptanceCriterionID string
	PassingRoles          []string
}

type CompletionReportStatus struct {
	IssueID string
	Exists  bool
}

type ConcernStatus struct {
	OpenConcernIDs []string
}

type ChildTaskStatus struct {
	IssueID string
	State   string
}

func AllACsVerified(requiredRoles []string, criteria []AcceptanceCriterionVerification) RuleResult {
	required := uniqueNonEmpty(requiredRoles)
	unmet := make([]UnmetCondition, 0)
	for _, criterion := range criteria {
		passingRoles := stringSet(criterion.PassingRoles)
		missingRoles := make([]string, 0)
		for _, role := range required {
			if !passingRoles[role] {
				missingRoles = append(missingRoles, role)
			}
		}
		if len(missingRoles) == 0 {
			continue
		}
		unmet = append(unmet, UnmetCondition{
			Rule:         RuleAllACsVerified,
			EntityIDs:    singleEntityID(criterion.AcceptanceCriterionID),
			MissingRoles: missingRoles,
		})
	}
	return newRuleResult(unmet)
}

func HasCompletionReport(status CompletionReportStatus) RuleResult {
	if status.Exists {
		return RuleResult{Satisfied: true}
	}
	return newRuleResult([]UnmetCondition{{
		Rule:      RuleHasCompletionReport,
		EntityIDs: singleEntityID(status.IssueID),
	}})
}

func ConcernsAddressed(status ConcernStatus) RuleResult {
	openConcernIDs := uniqueNonEmpty(status.OpenConcernIDs)
	if len(openConcernIDs) == 0 {
		return RuleResult{Satisfied: true}
	}
	return newRuleResult([]UnmetCondition{{
		Rule:      RuleConcernsAddressed,
		EntityIDs: openConcernIDs,
	}})
}

func AllChildTasksDoneOrCancelled(children []ChildTaskStatus) RuleResult {
	blockingChildIDs := make([]string, 0)
	for _, child := range children {
		state := strings.TrimSpace(child.State)
		if state == "done" || state == "cancelled" {
			continue
		}
		if issueID := strings.TrimSpace(child.IssueID); issueID != "" {
			blockingChildIDs = append(blockingChildIDs, issueID)
		}
	}
	if len(blockingChildIDs) == 0 {
		return RuleResult{Satisfied: true}
	}
	return newRuleResult([]UnmetCondition{{
		Rule:      RuleAllChildTasksDoneOrCanceled,
		EntityIDs: blockingChildIDs,
	}})
}

func newRuleResult(unmet []UnmetCondition) RuleResult {
	if len(unmet) == 0 {
		return RuleResult{Satisfied: true}
	}
	cloned := make([]UnmetCondition, 0, len(unmet))
	for _, condition := range unmet {
		cloned = append(cloned, UnmetCondition{
			Rule:         condition.Rule,
			EntityIDs:    append([]string(nil), condition.EntityIDs...),
			MissingRoles: append([]string(nil), condition.MissingRoles...),
		})
	}
	return RuleResult{
		Satisfied: false,
		Unmet:     cloned,
	}
}

func uniqueNonEmpty(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		result = append(result, trimmed)
	}
	return result
}

func stringSet(values []string) map[string]bool {
	set := make(map[string]bool, len(values))
	for _, value := range uniqueNonEmpty(values) {
		set[value] = true
	}
	return set
}

func singleEntityID(value string) []string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return []string{trimmed}
}
