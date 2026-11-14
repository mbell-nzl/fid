package issuestate

import (
	"database/sql"
	"fmt"
)

// ResolvedIssueWorkflow identifies the workflow and initial state for a new issue.
type ResolvedIssueWorkflow struct {
	ID           string
	InitialState string
}

// ResolveIssueWorkflow selects the fixed built-in workflow for an issue kind.
func ResolveIssueWorkflow(conn *sql.DB, kind string) (ResolvedIssueWorkflow, error) {
	workflowID := defaultIssueWorkflowID(kind)
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
