package workflow

type RuleReference struct {
	Name          string
	RequiredRoles []string
	Unimplemented bool
}

type Definition struct {
	ID          string
	States      map[string]string
	Transitions map[string]map[string][]RuleReference
}

func (definition Definition) CategoryForState(state string) (string, bool) {
	category, ok := definition.States[state]
	return category, ok
}

func (definition Definition) IsAllowedEdge(fromState, toState string) bool {
	_, ok := definition.RulesForEdge(fromState, toState)
	return ok
}

func (definition Definition) RulesForEdge(fromState, toState string) ([]RuleReference, bool) {
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
	return cloneRuleReferences(rules), true
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
			clonedDestinations[toState] = cloneRuleReferences(rules)
		}
		cloned.Transitions[fromState] = clonedDestinations
	}
	return cloned
}

func cloneRuleReferences(rules []RuleReference) []RuleReference {
	if rules == nil {
		return nil
	}
	cloned := make([]RuleReference, len(rules))
	for index, rule := range rules {
		cloned[index] = RuleReference{
			Name:          rule.Name,
			RequiredRoles: append([]string(nil), rule.RequiredRoles...),
			Unimplemented: rule.Unimplemented,
		}
	}
	return cloned
}
