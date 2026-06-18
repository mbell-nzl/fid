package issues

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var supportedIssueSections = []string{"acceptance", "concerns", "completions", "links", "claims"}

type UnsupportedSectionError struct {
	Section string
}

func (err *UnsupportedSectionError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("unsupported issue section %q; supported sections: %s", err.Section, strings.Join(supportedIssueSections, ", "))
}

type IssueConcern struct {
	ID                   string  `json:"id"`
	RaisedInCompletionID string  `json:"raised_in_completion_id"`
	RaisedByRole         string  `json:"raised_by_role"`
	Summary              string  `json:"summary"`
	BodyMDPath           *string `json:"body_md_path"`
	Status               string  `json:"status"`
	CreatedAt            *string `json:"created_at"`
	UpdatedAt            *string `json:"updated_at"`
}

type IssueCompletion struct {
	ID                string   `json:"id"`
	AttemptNumber     int      `json:"attempt_number"`
	DeveloperRole     *string  `json:"developer_role"`
	Model             *string  `json:"model"`
	Status            *string  `json:"status"`
	ActualTokens      *int     `json:"actual_tokens"`
	GateResults       []string `json:"gate_results"`
	HumanIntervention bool     `json:"human_intervention"`
	CompletedAt       *string  `json:"completed_at"`
}

type IssueLink struct {
	Relation    string  `json:"relation"`
	TargetType  string  `json:"target_type"`
	TargetID    string  `json:"target_id"`
	TargetTitle *string `json:"target_title"`
}

type IssueClaim struct {
	RefType    string  `json:"ref_type"`
	RefID      string  `json:"ref_id"`
	ClaimedBy  string  `json:"claimed_by"`
	ClaimedAt  string  `json:"claimed_at"`
	TTLSeconds int     `json:"ttl_seconds"`
	ExpiresAt  *string `json:"expires_at"`
	Note       *string `json:"note"`
}

func SupportedSections() []string {
	return append([]string(nil), supportedIssueSections...)
}

func SectionValue(issue Issue, section string) (any, error) {
	switch strings.TrimSpace(strings.ToLower(section)) {
	case "acceptance":
		return issue.AcceptanceCriteria, nil
	case "concerns":
		return issue.Concerns, nil
	case "completions":
		return issue.Completions, nil
	case "links":
		return issue.Links, nil
	case "claims":
		return issue.Claims, nil
	default:
		return nil, &UnsupportedSectionError{Section: strings.TrimSpace(section)}
	}
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

func issueAcceptanceCriteria(conn *sql.DB, issueID string, ownerRole *string) ([]AcceptanceCriterion, error) {
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

		criterion.LatestVerifications, err = latestVerificationsByRole(conn, criterion.ID)
		if err != nil {
			return nil, err
		}
		if len(criterion.LatestVerifications) == 0 && ownerRole != nil {
			criterion.LatestVerifications[*ownerRole] = nil
		}
		criteria = append(criteria, criterion)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate acceptance criteria for %q: %w", issueID, err)
	}
	return criteria, nil
}

func latestVerificationsByRole(conn *sql.DB, criterionID string) (map[string]any, error) {
	rows, err := conn.Query(`
		SELECT verifier_role, status, verified_at
		FROM ac_verifications av
		WHERE av.ac_id = ?
		  AND av.rowid = (
		    SELECT inner_av.rowid
		    FROM ac_verifications inner_av
		    WHERE inner_av.ac_id = av.ac_id AND inner_av.verifier_role = av.verifier_role
		    ORDER BY inner_av.rowid DESC
		    LIMIT 1
		  )
		ORDER BY verifier_role
	`, criterionID)
	if err != nil {
		return nil, fmt.Errorf("query acceptance criterion verifications for %q: %w", criterionID, err)
	}
	defer closeIssueRows(rows)

	latest := map[string]any{}
	for rows.Next() {
		var role string
		var status string
		var verifiedAt string
		if err := rows.Scan(&role, &status, &verifiedAt); err != nil {
			return nil, fmt.Errorf("scan acceptance criterion verification for %q: %w", criterionID, err)
		}
		latest[role] = map[string]any{
			"status":      status,
			"verified_at": verifiedAt,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate acceptance criterion verifications for %q: %w", criterionID, err)
	}
	return latest, nil
}

func issueConcerns(conn *sql.DB, issueID string) ([]IssueConcern, error) {
	rows, err := conn.Query(`
		SELECT id, raised_in_completion_id, raised_by_role, summary, body_md_path, status, created_at, updated_at
		FROM concerns
		WHERE issue_id = ?
		ORDER BY created_at DESC, id DESC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query concerns for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	concerns := []IssueConcern{}
	for rows.Next() {
		var concern IssueConcern
		var bodyMDPath, createdAt, updatedAt sql.NullString
		if err := rows.Scan(
			&concern.ID,
			&concern.RaisedInCompletionID,
			&concern.RaisedByRole,
			&concern.Summary,
			&bodyMDPath,
			&concern.Status,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan concern for %q: %w", issueID, err)
		}
		concern.BodyMDPath = nullStringPointer(bodyMDPath)
		concern.CreatedAt = nullStringPointer(createdAt)
		concern.UpdatedAt = nullStringPointer(updatedAt)
		concerns = append(concerns, concern)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate concerns for %q: %w", issueID, err)
	}
	return concerns, nil
}

func issueCompletions(conn *sql.DB, issueID string) ([]IssueCompletion, error) {
	rows, err := conn.Query(`
		SELECT id, attempt_number, developer_role, model, status, actual_tokens, gate_results, human_intervention, completed_at
		FROM task_completions
		WHERE issue_id = ?
		ORDER BY attempt_number DESC, id DESC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query completions for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	completions := []IssueCompletion{}
	for rows.Next() {
		var completion IssueCompletion
		var developerRole, model, status, gateResultsJSON, completedAt sql.NullString
		var actualTokens sql.NullInt64
		if err := rows.Scan(
			&completion.ID,
			&completion.AttemptNumber,
			&developerRole,
			&model,
			&status,
			&actualTokens,
			&gateResultsJSON,
			&completion.HumanIntervention,
			&completedAt,
		); err != nil {
			return nil, fmt.Errorf("scan completion for %q: %w", issueID, err)
		}
		completion.DeveloperRole = nullStringPointer(developerRole)
		completion.Model = nullStringPointer(model)
		completion.Status = nullStringPointer(status)
		completion.ActualTokens = nullIntPointer(actualTokens)
		completion.CompletedAt = nullStringPointer(completedAt)
		completion.GateResults = []string{}
		if gateResultsJSON.Valid && strings.TrimSpace(gateResultsJSON.String) != "" {
			if err := json.Unmarshal([]byte(gateResultsJSON.String), &completion.GateResults); err != nil {
				return nil, fmt.Errorf("decode completion gate results for %q: %w", completion.ID, err)
			}
		}
		completions = append(completions, completion)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate completions for %q: %w", issueID, err)
	}
	return completions, nil
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

func issueClaims(conn *sql.DB, issueID string) ([]IssueClaim, error) {
	rows, err := conn.Query(`
		SELECT ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
		ORDER BY claimed_at DESC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query claims for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	claims := []IssueClaim{}
	for rows.Next() {
		var claim IssueClaim
		var note sql.NullString
		if err := rows.Scan(&claim.RefType, &claim.RefID, &claim.ClaimedBy, &claim.ClaimedAt, &claim.TTLSeconds, &note); err != nil {
			return nil, fmt.Errorf("scan claim for %q: %w", issueID, err)
		}
		claim.Note = nullStringPointer(note)
		if claimedAt, err := time.Parse(time.RFC3339, claim.ClaimedAt); err == nil {
			expiresAt := claimedAt.Add(time.Duration(claim.TTLSeconds) * time.Second).Format(time.RFC3339)
			claim.ExpiresAt = &expiresAt
		}
		claims = append(claims, claim)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate claims for %q: %w", issueID, err)
	}
	return claims, nil
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
