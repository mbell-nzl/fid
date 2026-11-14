package workflow

// definitions is the fixed reset workflow reference used by the existing
// transition engine and later migration-parity tests.
var definitions = map[string]Definition{
	"task-default": {
		ID:          "task-default",
		States:      map[string]string{"todo": "unstarted", "in_progress": "started", "in_review": "started", "done": "completed", "cancelled": "cancelled"},
		Transitions: map[string]map[string][]RuleReference{"todo": {"in_progress": nil}, "in_progress": {"in_review": nil}, "in_review": {"in_progress": nil, "done": nil}},
	},
	"epic-default": {
		ID:          "epic-default",
		States:      map[string]string{"todo": "unstarted", "in_progress": "started", "done": "completed", "cancelled": "cancelled"},
		Transitions: map[string]map[string][]RuleReference{"todo": {"in_progress": nil}, "in_progress": {"done": nil}},
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

// Workflow returns an isolated fixed reset fixture definition.
func Workflow(workflowID string) (Definition, bool) {
	definition, ok := definitions[workflowID]
	if !ok {
		return Definition{}, false
	}
	return cloneDefinition(definition), true
}

// CategoryForState returns a category from the fixed reset fixture.
func CategoryForState(workflowID, state string) (string, bool) {
	definition, ok := definitions[workflowID]
	if !ok {
		return "", false
	}
	return definition.CategoryForState(state)
}

// IsAllowedEdge reports whether an edge exists in the fixed reset fixture.
func IsAllowedEdge(workflowID, fromState, toState string) bool {
	definition, ok := definitions[workflowID]
	return ok && definition.IsAllowedEdge(fromState, toState)
}

// RulesForEdge returns the ordered fixed reset fixture rules for an edge.
func RulesForEdge(workflowID, fromState, toState string) ([]RuleReference, bool) {
	definition, ok := definitions[workflowID]
	if !ok {
		return nil, false
	}
	return definition.RulesForEdge(fromState, toState)
}

// FrozenDefinitions returns an isolated copy of the reset workflow reference.
func FrozenDefinitions() map[string]Definition {
	fixtures := make(map[string]Definition, len(definitions))
	for workflowID, definition := range definitions {
		fixtures[workflowID] = cloneDefinition(definition)
	}
	return fixtures
}
