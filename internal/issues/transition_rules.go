package issues

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func loadIssueTransitionState(conn *sql.DB, issueID string) (issueTransitionState, error) {
	row := conn.QueryRow(`
		SELECT id, kind, workflow_id, status, updated_at, revision
		FROM issues
		WHERE id = ?
	`, issueID)

	var (
		state     issueTransitionState
		updatedAt sql.NullString
	)
	if err := row.Scan(&state.ID, &state.Kind, &state.WorkflowID, &state.Status, &updatedAt, &state.Revision); err != nil {
		if err == sql.ErrNoRows {
			return issueTransitionState{}, &NotFoundError{IssueID: issueID}
		}
		return issueTransitionState{}, fmt.Errorf("query issue %q: %w", issueID, err)
	}
	if !updatedAt.Valid || strings.TrimSpace(updatedAt.String) == "" {
		return issueTransitionState{}, fmt.Errorf("query issue %q: updated_at is required", issueID)
	}

	state.UpdatedAt = updatedAt.String
	return state, nil
}

func evaluateTransitionRules(conn *sql.DB, current issueTransitionState, rules []workflow.RuleReference) ([]workflow.UnmetCondition, error) {
	unmet := make([]workflow.UnmetCondition, 0)
	for _, rule := range rules {
		result, err := evaluateTransitionRule(conn, current, rule)
		if err != nil {
			return nil, err
		}
		unmet = append(unmet, result.Unmet...)
	}

	return unmet, nil
}

func evaluateTransitionRule(conn *sql.DB, current issueTransitionState, rule workflow.RuleReference) (workflow.RuleResult, error) {
	switch rule.Name {
	case workflow.RuleAllChildTasksDoneOrCanceled:
		children, err := loadChildTaskStatuses(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.AllChildTasksDoneOrCancelled(children), nil
	case workflow.RuleLinkedToDecision:
		ids, err := loadApprovedDecisionReferenceIDs(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.LinkedToDecision(current.ID, ids), nil
	default:
		return workflow.RuleResult{Unmet: []workflow.UnmetCondition{{Rule: rule.Name}}}, nil
	}
}

func loadApprovedDecisionReferenceIDs(conn *sql.DB, issueID string) ([]string, error) {
	rows, err := conn.Query(`SELECT id FROM decision_refs WHERE scope_issue_id=? AND outcome='approved' ORDER BY id`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query approved Decision References for %q: %w", issueID, err)
	}
	defer func() { _ = rows.Close() }()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func loadChildTaskStatuses(conn *sql.DB, issueID string) ([]workflow.ChildTaskStatus, error) {
	rows, err := conn.Query(`
		SELECT id, status
		FROM issues
		WHERE parent_id = ?
		ORDER BY id
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query child tasks for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	children := []workflow.ChildTaskStatus{}
	for rows.Next() {
		var child workflow.ChildTaskStatus
		if err := rows.Scan(&child.IssueID, &child.State); err != nil {
			return nil, fmt.Errorf("scan child task for %q: %w", issueID, err)
		}
		children = append(children, child)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate child tasks for %q: %w", issueID, err)
	}

	return children, nil
}

func issueTransitionArgv(issueID, toState, actor, idempotencyKey string, expectedRevision *int) []string {
	argv := []string{"fid", "issue", "transition", issueID, "--to", toState}
	if expectedRevision != nil {
		argv = append(argv, "--expected-revision", fmt.Sprintf("%d", *expectedRevision))
	}
	return append(argv, "--actor", actor, "--idempotency-key", idempotencyKey)
}
