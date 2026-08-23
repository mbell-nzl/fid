package issues

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/workflow"
)

const fidVersion = "0.1.0"

type CreateOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	Kind           string
	Title          string
	Actor          string
	IdempotencyKey string
	// ActorRole is retained temporarily for V1 in-process test callers. CLI
	// transport only supplies Actor and never uses role authority.
	ActorRole  string
	ParentID   *string
	Priority   *int
	Sequence   *int
	SpecMDPath *string
	Labels     []string
	// OwnerRole and OwnerActor remain compatibility inputs for retained V1
	// in-process callers. Canonical Issue commands never populate them.
	OwnerRole  *string
	OwnerActor *string
	Scheduling string
	now        func() time.Time
}

type CreateResult struct {
	ID             string
	Kind           string
	Title          string
	Status         string
	Revision       int
	EventID        string
	OccurredAt     string
	IdempotencyKey string
	Replayed       bool
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

type HierarchyInvalidError struct {
	Kind       string
	ParentID   *string
	ParentKind *string
	Rule       string
}

// StateConflictError records an optimistic-concurrency mismatch observed while
// the Issue mutation lock is held.
type StateConflictError struct {
	IssueID          string
	ExpectedRevision int
	ActualRevision   int
}

func (err *StateConflictError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("issue %q revision is %d, not expected revision %d", err.IssueID, err.ActualRevision, err.ExpectedRevision)
}

func (err *HierarchyInvalidError) Error() string {
	if err == nil {
		return "<nil>"
	}
	switch err.Kind {
	case "epic":
		return "issues: epic parent_id must be NULL"
	case "task":
		return "issues: task parent_id must reference an epic or be NULL"
	case "subtask":
		return "issues: subtask parent_id must reference a task"
	default:
		return fmt.Sprintf("issue hierarchy is invalid: %s", err.Rule)
	}
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
	Revision           int                   `json:"-"`
	OwnerRole          *string               `json:"owner_role"`
	OwnerActor         *string               `json:"owner_actor"`
	Labels             []string              `json:"labels"`
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
	Gates              []IssueGate           `json:"-"`
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

func registryStatusCategory(registry *workflow.Registry, workflowID, status string) (string, error) {
	category, ok := registry.CategoryForState(workflowID, status)
	if !ok {
		return "", fmt.Errorf("workflow %q does not define state %q", workflowID, status)
	}
	return category, nil
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
