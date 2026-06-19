package issues

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

const fidVersion = "0.1.0"

type CreateOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Kind          string
	Title         string
	ActorRole     string
	ParentID      *string
	Priority      *int
	Sequence      *int
	SpecMDPath    *string
	EstTokens     *int
	DueDate       *string
	Labels        []string
	InInbox       bool
	OwnerRole     *string
	OwnerActor    *string
	Scheduling    string
}

type CreateResult struct {
	ID    string
	Kind  string
	Title string
}

type GetOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	ID            string
}

type UnsupportedKindError struct {
	Kind string
}

func (err *UnsupportedKindError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("kind %q must be epic, task, or subtask", err.Kind)
}

type NotFoundError struct {
	IssueID string
}

func (err *NotFoundError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("issue %q not found", err.IssueID)
}

type Issue struct {
	ID                 string                `json:"id"`
	Kind               string                `json:"kind"`
	Title              string                `json:"title"`
	Status             string                `json:"status"`
	StatusCategory     string                `json:"status_category"`
	WorkflowID         string                `json:"workflow_id"`
	ParentID           *string               `json:"parent_id"`
	Priority           *int                  `json:"priority"`
	Sequence           *int                  `json:"sequence"`
	SpecMDPath         *string               `json:"spec_md_path"`
	EstTokens          *int                  `json:"est_tokens"`
	DueDate            *string               `json:"due_date"`
	OwnerRole          *string               `json:"owner_role"`
	OwnerActor         *string               `json:"owner_actor"`
	Labels             []string              `json:"labels"`
	InInbox            bool                  `json:"in_inbox"`
	ConcernsRaised     bool                  `json:"concerns_raised"`
	OpenConcernsCount  int                   `json:"open_concerns_count"`
	TriageStatus       *string               `json:"triage_status"`
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptance_criteria"`
	CreatedAt          *string               `json:"created_at"`
	UpdatedAt          *string               `json:"updated_at"`
	ParentTitle        *string               `json:"-"`
	Concerns           []IssueConcern        `json:"-"`
	Completions        []IssueCompletion     `json:"completions"`
	Links              []IssueLink           `json:"-"`
	Claims             []IssueClaim          `json:"-"`
}

type AcceptanceCriterion struct {
	ID                  string         `json:"id"`
	Ordinal             int            `json:"ordinal"`
	Text                string         `json:"text"`
	LatestVerifications map[string]any `json:"latest_verifications"`
}

func countJournalLines(path string) int {
	payload, err := os.ReadFile(path)
	if err != nil {
		return 1
	}
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func issueKindForID(conn *sql.DB, issueID string) (string, bool, error) {
	var kind string
	err := conn.QueryRow(`SELECT kind FROM issues WHERE id = ? LIMIT 1`, issueID).Scan(&kind)
	if err == nil {
		return kind, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	return "", false, fmt.Errorf("query parent issue %q: %w", issueID, err)
}

func normalizedLabels(input []string) ([]string, error) {
	if len(input) == 0 {
		return nil, nil
	}

	seen := make(map[string]bool, len(input))
	labels := make([]string, 0, len(input))
	for _, raw := range input {
		label := strings.TrimSpace(raw)
		if label == "" {
			return nil, fmt.Errorf("create issue: label must not be blank")
		}
		if seen[label] {
			continue
		}
		seen[label] = true
		labels = append(labels, label)
	}
	return labels, nil
}

func trimmedOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func statusCategory(status string) string {
	switch status {
	case "todo":
		return "unstarted"
	case "in_progress", "dev_complete", "qa_review", "qa_block", "failed_gate", "wrapped":
		return "started"
	case "qa_pass", "done", "closed":
		return "completed"
	case "cancelled":
		return "cancelled"
	default:
		return "unstarted"
	}
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	cloned := value.String
	return &cloned
}

func nullIntPointer(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}
	cloned := int(value.Int64)
	return &cloned
}

func closeIssueDB(conn *sql.DB) {
	_ = conn.Close()
}
