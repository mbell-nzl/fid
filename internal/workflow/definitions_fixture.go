package workflow

// definitions is the frozen v0.1 workflow reference used by the existing
// transition engine and later migration-parity tests.
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

// Workflow returns an isolated frozen v0.1 fixture definition.
func Workflow(workflowID string) (Definition, bool) {
	definition, ok := definitions[workflowID]
	if !ok {
		return Definition{}, false
	}
	return cloneDefinition(definition), true
}

// CategoryForState returns a category from the frozen v0.1 fixture.
func CategoryForState(workflowID, state string) (string, bool) {
	definition, ok := definitions[workflowID]
	if !ok {
		return "", false
	}
	return definition.CategoryForState(state)
}

// IsAllowedEdge reports whether an edge exists in the frozen v0.1 fixture.
func IsAllowedEdge(workflowID, fromState, toState string) bool {
	definition, ok := definitions[workflowID]
	return ok && definition.IsAllowedEdge(fromState, toState)
}

// RulesForEdge returns the ordered frozen v0.1 fixture rules for an edge.
func RulesForEdge(workflowID, fromState, toState string) ([]RuleReference, bool) {
	definition, ok := definitions[workflowID]
	if !ok {
		return nil, false
	}
	return definition.RulesForEdge(fromState, toState)
}

// FrozenDefinitions returns an isolated copy of the v0.1 workflow reference.
func FrozenDefinitions() map[string]Definition {
	fixtures := make(map[string]Definition, len(definitions))
	for workflowID, definition := range definitions {
		fixtures[workflowID] = cloneDefinition(definition)
	}
	return fixtures
}
