package workflow

import "strings"

const (
	RuleLinkedToDecision = "linked_to_decision"
)

// IsImplementedRule reports whether a rule has an evaluator in this release.
func IsImplementedRule(name string) bool {
	return name == RuleAllChildTasksDoneOrCanceled ||
		name == RuleLinkedToDecision
}

func LinkedToDecision(issueID string, decisionReferenceIDs []string) RuleResult {
	ids := uniqueNonEmpty(decisionReferenceIDs)
	if len(ids) > 0 {
		return RuleResult{Satisfied: true}
	}
	return newRuleResult([]UnmetCondition{{Rule: RuleLinkedToDecision, EntityIDs: singleEntityID(issueID)}})
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

func singleEntityID(value string) []string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return []string{trimmed}
}
