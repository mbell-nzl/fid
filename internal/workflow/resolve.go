package workflow

import (
	"database/sql"
	"errors"
	"fmt"
)

// ResolvedIssueWorkflow identifies the workflow and initial state for a new issue.
type ResolvedIssueWorkflow struct {
	ID           string
	InitialState string
}

// GoverningHumanOutcome is the latest durable judgment for one class and
// Issue scope. Event identity, rather than a caller-authored timestamp, fixes
// the order.
type GoverningHumanOutcome struct {
	ID      string
	Result  string
	EventID string
}

// ResolveGoverningHumanOutcome selects by immutable Event identity so caller
// timestamps cannot change which judgment governs.
func ResolveGoverningHumanOutcome(conn *sql.DB, class, scopeIssueID string) (GoverningHumanOutcome, error) {
	var outcome GoverningHumanOutcome
	err := conn.QueryRow(`
		SELECT id,result,created_event_id
		FROM human_outcomes
		WHERE class=? AND scope_issue_id=?
		ORDER BY created_event_id DESC
		LIMIT 1`, class, scopeIssueID).Scan(&outcome.ID, &outcome.Result, &outcome.EventID)
	if err != nil {
		return GoverningHumanOutcome{}, fmt.Errorf("resolve governing Human Outcome: %w", err)
	}
	return outcome, nil
}

// HumanOutcomeSuperseded is the single supersession predicate (contracts.md
// section 7): an Outcome is superseded exactly when the governing Outcome for
// its class and scope is a different Outcome with a different result. Both the
// Obligation evaluator and the Stop resolution predicates call it.
func HumanOutcomeSuperseded(governing GoverningHumanOutcome, outcomeID, result string) bool {
	return governing.ID != outcomeID && governing.Result != result
}

// ResolveIssueWorkflow resolves kind, project, then built-in workflow assignment precedence.
func ResolveIssueWorkflow(conn *sql.DB, kind string) (ResolvedIssueWorkflow, error) {
	workflowID, err := assignedWorkflowID(conn, kind)
	if err != nil {
		return ResolvedIssueWorkflow{}, err
	}
	states, err := loadWorkflowStates(conn, workflowID)
	if err != nil {
		return ResolvedIssueWorkflow{}, err
	}
	initialState := InitialState(states)
	if initialState == "" {
		return ResolvedIssueWorkflow{}, fmt.Errorf("resolve workflow %q: no states", workflowID)
	}
	return ResolvedIssueWorkflow{ID: workflowID, InitialState: initialState}, nil
}

func assignedWorkflowID(conn *sql.DB, kind string) (string, error) {
	if kind == "task" || kind == "epic" {
		return defaultIssueWorkflowID(kind), nil
	}
	for _, scope := range []struct {
		query string
		args  []any
	}{
		{query: `SELECT workflow_id FROM workflow_assignments WHERE scope = 'kind' AND kind = ?`, args: []any{kind}},
		{query: `SELECT workflow_id FROM workflow_assignments WHERE scope = 'project'`, args: nil},
	} {
		var workflowID string
		err := conn.QueryRow(scope.query, scope.args...).Scan(&workflowID)
		if err == nil {
			return workflowID, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("resolve workflow assignment: %w", err)
		}
	}
	return defaultIssueWorkflowID(kind), nil
}

func loadWorkflowStates(conn *sql.DB, workflowID string) ([]State, error) {
	rows, err := conn.Query(`
		SELECT state, category
		FROM workflow_states
		WHERE workflow_id = ?
		ORDER BY ordinal`, workflowID)
	if err != nil {
		return nil, fmt.Errorf("load workflow %q states: %w", workflowID, err)
	}
	defer func() { _ = rows.Close() }()
	states := make([]State, 0)
	for rows.Next() {
		var state State
		if err := rows.Scan(&state.Name, &state.Category); err != nil {
			return nil, fmt.Errorf("scan workflow %q state: %w", workflowID, err)
		}
		states = append(states, state)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read workflow %q states: %w", workflowID, err)
	}
	return states, nil
}

func defaultIssueWorkflowID(kind string) string {
	switch kind {
	case "epic":
		return "epic-default"
	case "subtask":
		return "subtask-default"
	default:
		return "task-default"
	}
}
