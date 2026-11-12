package workflow

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// Registry resolves workflow definitions from the projected workflow tables.
type Registry struct {
	workflows map[string]Definition
	edges     map[transitionKey][]RuleReference
}

// NextState describes a workflow edge's target state and category.
type NextState struct {
	State    string `json:"state"`
	Category string `json:"category"`
}

// NextStates lists every allowed edge target, without evaluating its rules.
func (registry *Registry) NextStates(workflowID, fromState string) []NextState {
	states := []NextState{}
	definition, ok := registry.workflows[workflowID]
	if !ok {
		return states
	}
	for state, category := range definition.States {
		if _, edge := registry.RulesForEdge(workflowID, fromState, state); edge {
			states = append(states, NextState{State: state, Category: category})
		}
	}
	sort.Slice(states, func(i, j int) bool { return states[i].State < states[j].State })
	return states
}

type transitionKey struct {
	workflowID string
	fromState  string
	toState    string
}

// Load builds a read-only workflow registry from the projected workflow tables.
func Load(conn *sql.DB) (*Registry, error) {
	if conn == nil {
		return nil, fmt.Errorf("load workflow registry: database connection is required")
	}
	return LoadFrom(conn)
}

// LoadFrom reads a registry using the caller's projection snapshot.
func LoadFrom(conn fiddb.Queryer) (*Registry, error) {
	rows, err := conn.Query(`SELECT id FROM workflows ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("load workflow registry: query workflows: %w", err)
	}
	defer func() { _ = rows.Close() }()

	workflowIDs := make([]string, 0)
	for rows.Next() {
		var workflowID string
		if err := rows.Scan(&workflowID); err != nil {
			return nil, fmt.Errorf("load workflow registry: scan workflow: %w", err)
		}
		workflowIDs = append(workflowIDs, workflowID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("load workflow registry: iterate workflows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("load workflow registry: close workflows: %w", err)
	}

	registry := &Registry{
		workflows: make(map[string]Definition),
		edges:     make(map[transitionKey][]RuleReference),
	}
	for _, workflowID := range workflowIDs {
		definition, err := loadDefinition(conn, registry, workflowID)
		if err != nil {
			return nil, err
		}
		registry.workflows[workflowID] = definition
	}
	return registry, nil
}

// Workflow returns an isolated definition from the loaded registry.
func (registry *Registry) Workflow(workflowID string) (Definition, bool) {
	definition, ok := registry.workflows[workflowID]
	if !ok {
		return Definition{}, false
	}
	return cloneDefinition(definition), true
}

// CategoryForState returns a state category from the loaded registry.
func (registry *Registry) CategoryForState(workflowID, state string) (string, bool) {
	definition, ok := registry.workflows[workflowID]
	if !ok {
		return "", false
	}
	return definition.CategoryForState(state)
}

// IsAllowedEdge reports whether an edge exists in the loaded registry.
func (registry *Registry) IsAllowedEdge(workflowID, fromState, toState string) bool {
	_, ok := registry.RulesForEdge(workflowID, fromState, toState)
	return ok
}

// RulesForEdge returns the ordered rules for an edge in the loaded registry.
func (registry *Registry) RulesForEdge(workflowID, fromState, toState string) ([]RuleReference, bool) {
	definition, ok := registry.workflows[workflowID]
	if !ok || !knownState(definition, fromState) || !knownState(definition, toState) {
		return nil, false
	}
	rules, ok := registry.edges[transitionKey{
		workflowID: workflowID,
		fromState:  fromState,
		toState:    toState,
	}]
	if !ok {
		return nil, false
	}
	return cloneRuleReferences(rules), true
}

func loadDefinition(conn fiddb.Queryer, registry *Registry, workflowID string) (Definition, error) {
	definition := Definition{
		ID:          workflowID,
		States:      make(map[string]string),
		Transitions: make(map[string]map[string][]RuleReference),
	}
	if err := loadStates(conn, &definition); err != nil {
		return Definition{}, err
	}
	if err := loadTransitions(conn, registry, &definition); err != nil {
		return Definition{}, err
	}
	return definition, loadTransitionRequirements(conn, registry, &definition)
}

func loadStates(conn fiddb.Queryer, definition *Definition) error {
	rows, err := conn.Query(`
		SELECT state, category
		FROM workflow_states
		WHERE workflow_id = ?
		ORDER BY ordinal, state
	`, definition.ID)
	if err != nil {
		return fmt.Errorf("load workflow %q states: %w", definition.ID, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var state, category string
		if err := rows.Scan(&state, &category); err != nil {
			return fmt.Errorf("load workflow %q states: %w", definition.ID, err)
		}
		definition.States[state] = category
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("load workflow %q states: %w", definition.ID, err)
	}
	return nil
}

func loadTransitions(conn fiddb.Queryer, registry *Registry, definition *Definition) error {
	rows, err := conn.Query(`
		SELECT from_state, to_state
		FROM workflow_transitions
		WHERE workflow_id = ?
		ORDER BY from_state, to_state
	`, definition.ID)
	if err != nil {
		return fmt.Errorf("load workflow %q transitions: %w", definition.ID, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var fromState, toState string
		if err := rows.Scan(&fromState, &toState); err != nil {
			return fmt.Errorf("load workflow %q transitions: %w", definition.ID, err)
		}
		registry.edges[transitionKey{
			workflowID: definition.ID,
			fromState:  fromState,
			toState:    toState,
		}] = nil
		if toState == "cancelled" && fromState != "cancelled" {
			continue
		}
		if definition.Transitions[fromState] == nil {
			definition.Transitions[fromState] = make(map[string][]RuleReference)
		}
		definition.Transitions[fromState][toState] = nil
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("load workflow %q transitions: %w", definition.ID, err)
	}
	return nil
}

func loadTransitionRequirements(conn fiddb.Queryer, registry *Registry, definition *Definition) error {
	rows, err := conn.Query(`
		SELECT from_state, to_state, rule, params_json
		FROM workflow_transition_requires
		WHERE workflow_id = ?
		ORDER BY from_state, to_state, rowid
	`, definition.ID)
	if err != nil {
		return fmt.Errorf("load workflow %q transition requirements: %w", definition.ID, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var fromState, toState, name string
		var paramsJSON sql.NullString
		if err := rows.Scan(&fromState, &toState, &name, &paramsJSON); err != nil {
			return fmt.Errorf("load workflow %q transition requirements: %w", definition.ID, err)
		}
		if paramsJSON.Valid && strings.TrimSpace(paramsJSON.String) != "" {
			return fmt.Errorf("load workflow %q transition %q -> %q rule %q: params_json is unsupported", definition.ID, fromState, toState, name)
		}
		rule := RuleReference{Name: name, Unimplemented: isReservedUnimplementedRule(name)}
		key := transitionKey{
			workflowID: definition.ID,
			fromState:  fromState,
			toState:    toState,
		}
		registry.edges[key] = append(registry.edges[key], rule)
		if toState != "cancelled" || fromState == "cancelled" {
			definition.Transitions[fromState][toState] = append(
				definition.Transitions[fromState][toState], rule,
			)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("load workflow %q transition requirements: %w", definition.ID, err)
	}
	return nil
}

func isReservedUnimplementedRule(name string) bool {
	return name == RuleRequiredGatesCurrent || name == RuleHasCompletion
}
