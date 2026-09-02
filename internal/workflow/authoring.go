package workflow

import (
	"database/sql"
	"fmt"
	"regexp"
)

var (
	workflowSlugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
	stateNamePattern    = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)
)

// DefinitionRequest is a parsed workflow definition before journal effects exist.
type DefinitionRequest struct {
	Name         string
	AppliesTo    string
	States       []string
	Categories   []string
	Transitions  []Transition
	Requirements []Requirement
}

// Transition identifies one directed workflow edge.
type Transition struct {
	From string
	To   string
}

// Requirement attaches an implemented rule to one workflow edge.
type Requirement struct {
	From       string
	To         string
	Rule       string
	ParamsJSON string
}

// State preserves a workflow state's declared order for initial-state resolution.
type State struct {
	Name     string
	Category string
}

// ValidateDefinitionRequest checks a parsed definition without mutating the projection.
func ValidateDefinitionRequest(conn *sql.DB, request DefinitionRequest) error {
	if !workflowSlugPattern.MatchString(request.Name) {
		return fmt.Errorf("invalid workflow slug %q", request.Name)
	}
	if err := rejectExistingWorkflow(conn, request.Name); err != nil {
		return err
	}
	states, cancelled, err := validateStates(request)
	if err != nil {
		return err
	}
	if err := validateCancelledStates(request.AppliesTo, cancelled); err != nil {
		return err
	}
	_, err = validateTransitions(request.Transitions, states)
	if err != nil {
		return err
	}
	if len(request.Requirements) != 0 {
		return fmt.Errorf("workflow transition requirements are unavailable")
	}
	return nil
}

func rejectExistingWorkflow(conn *sql.DB, name string) error {
	if conn == nil {
		return fmt.Errorf("check workflow slug %q: database connection is required", name)
	}
	var existing string
	err := conn.QueryRow(`SELECT id FROM workflows WHERE id = ? LIMIT 1`, name).Scan(&existing)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("check workflow slug %q: %w", name, err)
	}
	return &DuplicateError{ID: name}
}

func validateStates(request DefinitionRequest) (map[string]bool, []string, error) {
	if len(request.States) == 0 {
		return nil, nil, fmt.Errorf("workflow states must not be empty")
	}
	if len(request.Categories) != len(request.States) {
		return nil, nil, fmt.Errorf("workflow states/categories count mismatch: %d states, %d categories", len(request.States), len(request.Categories))
	}
	states := make(map[string]bool, len(request.States))
	cancelled := make([]string, 0, 1)
	for index, name := range request.States {
		if !stateNamePattern.MatchString(name) {
			return nil, nil, fmt.Errorf("invalid workflow state %q", name)
		}
		if states[name] {
			return nil, nil, fmt.Errorf("duplicate workflow state %q", name)
		}
		states[name] = true
		category := request.Categories[index]
		if !isWorkflowCategory(category) {
			return nil, nil, fmt.Errorf("unknown workflow category %q", category)
		}
		if category == "cancelled" {
			cancelled = append(cancelled, name)
		}
	}
	return states, cancelled, nil
}

func validateCancelledStates(appliesTo string, cancelled []string) error {
	if len(cancelled) > 1 {
		return fmt.Errorf("multiple cancelled workflow states %q and %q", cancelled[0], cancelled[1])
	}
	if len(cancelled) == 0 && isIssueApplicable(appliesTo) {
		return fmt.Errorf("workflow applies_to %q must declare cancelled state: any issue may be cancelled from any non-cancelled state (PRD-WF-013)", appliesTo)
	}
	return nil
}

func validateTransitions(transitions []Transition, states map[string]bool) (map[string]bool, error) {
	known := make(map[string]bool, len(transitions))
	for _, transition := range transitions {
		if !states[transition.From] || !states[transition.To] {
			return nil, fmt.Errorf("workflow transition %q -> %q references unknown state", transition.From, transition.To)
		}
		key := transitionKeyFor(transition.From, transition.To)
		if known[key] {
			return nil, fmt.Errorf("duplicate workflow transition %q -> %q", transition.From, transition.To)
		}
		known[key] = true
	}
	return known, nil
}

func isWorkflowCategory(category string) bool {
	return category == "unstarted" || category == "started" || category == "completed" || category == "cancelled"
}

func isIssueApplicable(appliesTo string) bool {
	return appliesTo == "epic" || appliesTo == "task" || appliesTo == "subtask" || appliesTo == "any"
}

func transitionKeyFor(from, to string) string {
	return from + "\x00" + to
}

// InitialState returns the first unstarted state, or the first declared state.
func InitialState(states []State) string {
	for _, state := range states {
		if state.Category == "unstarted" {
			return state.Name
		}
	}
	if len(states) == 0 {
		return ""
	}
	return states[0].Name
}
