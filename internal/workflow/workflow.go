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
