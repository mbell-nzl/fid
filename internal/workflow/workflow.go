package workflow

type RuleReference struct {
	Name          string
	RequiredRoles []string
}

type Definition struct {
	ID          string
	States      map[string]string
	Transitions map[string]map[string][]RuleReference
}

var definitions = map[string]Definition{
	"task-default": {
		ID: "task-default",
		States: map[string]string{
			"todo":         "unstarted",
			"in_progress":  "started",
			"dev_complete": "started",
			"qa_review":    "started",
			"qa_block":     "started",
			"failed_gate":  "started",
			"qa_pass":      "completed",
			"done":         "completed",
			"cancelled":    "cancelled",
		},
		Transitions: map[string]map[string][]RuleReference{
			"todo": {
				"in_progress": nil,
			},
			"in_progress": {
				"dev_complete": {
					{Name: RuleAllACsVerified, RequiredRoles: []string{"developer"}},
					{Name: RuleHasCompletionReport},
				},
				"failed_gate": nil,
			},
			"failed_gate": {
				"in_progress": nil,
			},
			"dev_complete": {
				"qa_review": nil,
			},
			"qa_review": {
				"qa_pass": {
					{Name: RuleAllACsVerified, RequiredRoles: []string{"qa_lead"}},
				},
				"qa_block": nil,
			},
			"qa_block": {
				"in_progress": nil,
			},
			"qa_pass": {
				"done": nil,
			},
		},
	},
	"epic-default": {
		ID: "epic-default",
		States: map[string]string{
			"todo":        "unstarted",
			"in_progress": "started",
			"wrapped":     "started",
			"closed":      "completed",
			"cancelled":   "cancelled",
		},
		Transitions: map[string]map[string][]RuleReference{
			"todo": {
				"in_progress": nil,
			},
			"in_progress": {
				"wrapped": {
					{Name: RuleAllChildTasksDoneOrCanceled},
				},
			},
			"wrapped": {
				"closed": nil,
			},
		},
	},
	"subtask-default": {
		ID: "subtask-default",
		States: map[string]string{
			"todo":        "unstarted",
			"in_progress": "started",
			"done":        "completed",
			"cancelled":   "cancelled",
		},
		Transitions: map[string]map[string][]RuleReference{
			"todo": {
				"in_progress": nil,
			},
			"in_progress": {
				"done": nil,
			},
		},
	},
}

func Workflow(workflowID string) (Definition, bool) {
	definition, ok := definitions[workflowID]
	if !ok {
		return Definition{}, false
	}
	return cloneDefinition(definition), true
}

func CategoryForState(workflowID, state string) (string, bool) {
	definition, ok := definitions[workflowID]
	if !ok {
		return "", false
	}
	category, ok := definition.States[state]
	return category, ok
}

func IsAllowedEdge(workflowID, fromState, toState string) bool {
	_, ok := RulesForEdge(workflowID, fromState, toState)
	return ok
}

func RulesForEdge(workflowID, fromState, toState string) ([]RuleReference, bool) {
	definition, ok := definitions[workflowID]
	if !ok {
		return nil, false
	}
	if !knownState(definition, fromState) || !knownState(definition, toState) {
		return nil, false
	}
	if toState == "cancelled" && fromState != "cancelled" {
		return nil, true
	}
	destinations, ok := definition.Transitions[fromState]
	if !ok {
		return nil, false
	}
	rules, ok := destinations[toState]
	if !ok {
		return nil, false
	}
	return append([]RuleReference(nil), rules...), true
}

func knownState(definition Definition, state string) bool {
	_, ok := definition.States[state]
	return ok
}

func cloneDefinition(definition Definition) Definition {
	cloned := Definition{
		ID:          definition.ID,
		States:      make(map[string]string, len(definition.States)),
		Transitions: make(map[string]map[string][]RuleReference, len(definition.Transitions)),
	}
	for state, category := range definition.States {
		cloned.States[state] = category
	}
	for fromState, destinations := range definition.Transitions {
		clonedDestinations := make(map[string][]RuleReference, len(destinations))
		for toState, rules := range destinations {
			clonedDestinations[toState] = append([]RuleReference(nil), rules...)
		}
		cloned.Transitions[fromState] = clonedDestinations
	}
	return cloned
}
