package issues

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type IssueConcern struct {
	ID          string `json:"id"`
	SubjectType string `json:"subject_type"`
	SubjectID   string `json:"subject_id"`
	ActorID     string `json:"actor_id"`
	Summary     string `json:"summary"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type IssueLink struct {
	Relation    string  `json:"relation"`
	TargetType  string  `json:"target_type"`
	TargetID    string  `json:"target_id"`
	TargetTitle *string `json:"target_title"`
}

func issueParentTitle(conn *sql.DB, parentID *string) (*string, error) {
	if parentID == nil {
		return nil, nil
	}

	var title string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, *parentID).Scan(&title); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query parent issue %q: %w", *parentID, err)
	}
	return &title, nil
}

func issueLabels(conn *sql.DB, issueID string) ([]string, error) {
	rows, err := conn.Query(`SELECT label FROM issue_labels WHERE issue_id = ? ORDER BY label`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query issue labels for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	labels := []string{}
	for rows.Next() {
		var label string
		if err := rows.Scan(&label); err != nil {
			return nil, fmt.Errorf("scan issue label for %q: %w", issueID, err)
		}
		labels = append(labels, label)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate issue labels for %q: %w", issueID, err)
	}
	return labels, nil
}

func issueAcceptanceCriteria(conn *sql.DB, issueID string) ([]AcceptanceCriterion, error) {
	rows, err := conn.Query(`
		SELECT id, ordinal, text
		FROM acceptance_criteria
		WHERE issue_id = ?
		ORDER BY ordinal, id
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query acceptance criteria for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	criteria := []AcceptanceCriterion{}
	for rows.Next() {
		var criterion AcceptanceCriterion
		if err := rows.Scan(&criterion.ID, &criterion.Ordinal, &criterion.Text); err != nil {
			return nil, fmt.Errorf("scan acceptance criterion for %q: %w", issueID, err)
		}

		criteria = append(criteria, criterion)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate acceptance criteria for %q: %w", issueID, err)
	}
	return criteria, nil
}

func issueConcerns(conn *sql.DB, issueID string) ([]IssueConcern, error) {
	rows, err := conn.Query(`SELECT id, subject_type, subject_id, actor_id, summary, status, created_at FROM concerns WHERE subject_type='issue' AND subject_id=? ORDER BY created_at, id`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query concerns for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)
	concerns := []IssueConcern{}
	for rows.Next() {
		var concern IssueConcern
		if err := rows.Scan(&concern.ID, &concern.SubjectType, &concern.SubjectID, &concern.ActorID, &concern.Summary, &concern.Status, &concern.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan concern for %q: %w", issueID, err)
		}
		concerns = append(concerns, concern)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate concerns for %q: %w", issueID, err)
	}
	return concerns, nil
}

func issueLinks(conn *sql.DB, issueID string) ([]IssueLink, error) {
	rows, err := conn.Query(`
		SELECT l.relation, l.to_type, l.to_id, target_issue.title
		FROM links l
		LEFT JOIN issues target_issue
		  ON l.to_type = 'issue' AND target_issue.id = l.to_id
		WHERE l.from_type = 'issue' AND l.from_id = ?
		ORDER BY l.created_at, l.relation, l.to_type, l.to_id
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query issue links for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	links := []IssueLink{}
	for rows.Next() {
		var link IssueLink
		var targetTitle sql.NullString
		if err := rows.Scan(&link.Relation, &link.TargetType, &link.TargetID, &targetTitle); err != nil {
			return nil, fmt.Errorf("scan issue link for %q: %w", issueID, err)
		}
		link.TargetTitle = nullStringPointer(targetTitle)
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate issue links for %q: %w", issueID, err)
	}
	return links, nil
}

func closeIssueRows(rows *sql.Rows) {
	_ = rows.Close()
}

func normalizeDateOnlyPointer(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	if len(trimmed) >= len("2006-01-02") {
		dateOnly := trimmed[:len("2006-01-02")]
		if _, err := time.Parse("2006-01-02", dateOnly); err == nil {
			return &dateOnly
		}
	}

	return &trimmed
}
