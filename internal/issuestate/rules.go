package issuestate

import "strings"

const (
	RuleAllChildTasksDoneOrCanceled = "all_child_tasks_done_or_cancelled"
	// Reserved rule names remain readable while kernel callers still exist.
	RuleRequiredGatesCurrent = "required_gates_current"
	RuleHasCompletion        = "has_completion"
)

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
