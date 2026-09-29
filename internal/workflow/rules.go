package workflow

import "strings"

const (
	RuleAllChildTasksDoneOrCanceled = "all_child_tasks_done_or_cancelled"
	// RuleRequiredGatesCurrent is reserved for the v0.2b trusted-gate path.
	// It has no evaluator in this phase.
	RuleRequiredGatesCurrent = "required_gates_current"
	// RuleHasCompletion is reserved for structured completion requirements.
	// It has no evaluator in this phase.
	RuleHasCompletion    = "has_completion"
	RuleLinkedToDecision = "linked_to_decision"
)

// IsImplementedRule reports whether a rule has an evaluator in this release.
func IsImplementedRule(name string) bool {
	return name == RuleAllChildTasksDoneOrCanceled ||
		name == RuleLinkedToDecision
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

type ChildTaskStatus struct {
	IssueID string
	State   string
}

func LinkedToDecision(issueID string, decisionReferenceIDs []string) RuleResult {
	ids := uniqueNonEmpty(decisionReferenceIDs)
	if len(ids) > 0 {
		return RuleResult{Satisfied: true}
	}
	return newRuleResult([]UnmetCondition{{Rule: RuleLinkedToDecision, EntityIDs: singleEntityID(issueID)}})
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
