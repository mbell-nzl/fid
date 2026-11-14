package issuestate

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
