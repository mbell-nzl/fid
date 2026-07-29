package workflow

import (
	"database/sql"

	"github.com/mbell-nzl/fid/internal/journal"
)

// BuildDefinitionEffects validates a definition and returns its ordered journal effects.
func BuildDefinitionEffects(conn *sql.DB, request DefinitionRequest) ([]journal.Effect, error) {
	if err := ValidateDefinitionRequest(conn, request); err != nil {
		return nil, err
	}
	effects := buildWorkflowEffect(request)
	effects = append(effects, buildStateEffects(request)...)
	transitions := materializedTransitions(request)
	effects = append(effects, buildTransitionEffects(request.Name, transitions)...)
	effects = append(effects, buildRequirementEffects(request)...)
	return effects, nil
}

func buildWorkflowEffect(request DefinitionRequest) []journal.Effect {
	appliesTo := anyAppliesToValue(request.AppliesTo)
	return []journal.Effect{{
		Op:    "insert",
		Table: "workflows",
		PK:    map[string]any{"id": request.Name},
		Values: map[string]any{
			"name": request.Name, "applies_to": appliesTo, "active": true,
		},
	}}
}

func buildStateEffects(request DefinitionRequest) []journal.Effect {
	effects := make([]journal.Effect, 0, len(request.States))
	for index, state := range request.States {
		effects = append(effects, journal.Effect{
			Op: "insert", Table: "workflow_states",
			PK:     map[string]any{"workflow_id": request.Name, "state": state},
			Values: map[string]any{"category": request.Categories[index], "ordinal": index + 1},
		})
	}
	return effects
}

func materializedTransitions(request DefinitionRequest) []Transition {
	transitions := append([]Transition(nil), request.Transitions...)
	cancelledState := cancelledState(request)
	if cancelledState == "" {
		return transitions
	}
	known := make(map[string]bool, len(transitions))
	for _, transition := range transitions {
		known[transitionKeyFor(transition.From, transition.To)] = true
	}
	for index, state := range request.States {
		if request.Categories[index] == "cancelled" {
			continue
		}
		key := transitionKeyFor(state, cancelledState)
		if !known[key] {
			transitions = append(transitions, Transition{From: state, To: cancelledState})
			known[key] = true
		}
	}
	return transitions
}

func cancelledState(request DefinitionRequest) string {
	for index, category := range request.Categories {
		if category == "cancelled" {
			return request.States[index]
		}
	}
	return ""
}

func buildTransitionEffects(workflowID string, transitions []Transition) []journal.Effect {
	effects := make([]journal.Effect, 0, len(transitions))
	for _, transition := range transitions {
		effects = append(effects, journal.Effect{
			Op: "insert", Table: "workflow_transitions",
			PK: map[string]any{"workflow_id": workflowID, "from_state": transition.From, "to_state": transition.To},
		})
	}
	return effects
}

func buildRequirementEffects(request DefinitionRequest) []journal.Effect {
	effects := make([]journal.Effect, 0, len(request.Requirements))
	for _, requirement := range request.Requirements {
		paramsJSON := any(nil)
		if requirement.ParamsJSON != "" {
			paramsJSON = requirement.ParamsJSON
		}
		effects = append(effects, journal.Effect{
			Op: "insert", Table: "workflow_transition_requires",
			PK:     map[string]any{"workflow_id": request.Name, "from_state": requirement.From, "to_state": requirement.To, "rule": requirement.Rule},
			Values: map[string]any{"params_json": paramsJSON},
		})
	}
	return effects
}

func anyAppliesToValue(appliesTo string) any {
	if appliesTo == "any" {
		return nil
	}
	return appliesTo
}
