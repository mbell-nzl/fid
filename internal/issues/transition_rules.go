package issues

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func loadIssueTransitionState(conn *sql.DB, issueID string) (issueTransitionState, error) {
	row := conn.QueryRow(`
		SELECT id, kind, workflow_id, status, updated_at
		FROM issues
		WHERE id = ?
	`, issueID)

	var (
		state     issueTransitionState
		updatedAt sql.NullString
	)
	if err := row.Scan(&state.ID, &state.Kind, &state.WorkflowID, &state.Status, &updatedAt); err != nil {
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
	case workflow.RuleAllACsVerified:
		criteria, err := loadAcceptanceCriterionVerifications(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.AllACsVerified(rule.RequiredRoles, criteria), nil
	case workflow.RuleHasCompletionReport:
		status, err := loadCompletionReportStatus(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.HasCompletionReport(status), nil
	case workflow.RuleConcernsAddressed:
		status, err := loadConcernStatus(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.ConcernsAddressed(status), nil
	case workflow.RuleAllChildTasksDoneOrCanceled:
		children, err := loadChildTaskStatuses(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.AllChildTasksDoneOrCancelled(children), nil
	default:
		return workflow.RuleResult{}, fmt.Errorf("transition issue %q: unsupported workflow rule %q", current.ID, rule.Name)
	}
}

func loadAcceptanceCriterionVerifications(conn *sql.DB, issueID string) ([]workflow.AcceptanceCriterionVerification, error) {
	criteria, err := issueAcceptanceCriteria(conn, issueID, nil)
	if err != nil {
		return nil, err
	}

	result := make([]workflow.AcceptanceCriterionVerification, 0, len(criteria))
	for _, criterion := range criteria {
		passingRoles := make([]string, 0, len(criterion.LatestVerifications))
		for role, raw := range criterion.LatestVerifications {
			verification, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			status, _ := verification["status"].(string)
			if strings.TrimSpace(status) == "pass" {
				passingRoles = append(passingRoles, role)
			}
		}
		sort.Strings(passingRoles)
		result = append(result, workflow.AcceptanceCriterionVerification{
			AcceptanceCriterionID: criterion.ID,
			PassingRoles:          passingRoles,
		})
	}

	return result, nil
}

func loadCompletionReportStatus(conn *sql.DB, issueID string) (workflow.CompletionReportStatus, error) {
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM task_completions WHERE issue_id = ?`, issueID).Scan(&count); err != nil {
		return workflow.CompletionReportStatus{}, fmt.Errorf("query completion reports for %q: %w", issueID, err)
	}

	return workflow.CompletionReportStatus{
		IssueID: issueID,
		Exists:  count > 0,
	}, nil
}

func loadConcernStatus(conn *sql.DB, issueID string) (workflow.ConcernStatus, error) {
	rows, err := conn.Query(`
		SELECT id
		FROM concerns
		WHERE issue_id = ? AND status = ?
		ORDER BY id
	`, issueID, "open")
	if err != nil {
		return workflow.ConcernStatus{}, fmt.Errorf("query open concerns for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	openConcernIDs := []string{}
	for rows.Next() {
		var concernID string
		if err := rows.Scan(&concernID); err != nil {
			return workflow.ConcernStatus{}, fmt.Errorf("scan open concern for %q: %w", issueID, err)
		}
		openConcernIDs = append(openConcernIDs, concernID)
	}
	if err := rows.Err(); err != nil {
		return workflow.ConcernStatus{}, fmt.Errorf("iterate open concerns for %q: %w", issueID, err)
	}

	return workflow.ConcernStatus{OpenConcernIDs: openConcernIDs}, nil
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

func issueTransitionArgv(issueID, toState, actorRole string) []string {
	return []string{"fid", "issue", "transition", issueID, "--to", toState, "--by", actorRole}
}
