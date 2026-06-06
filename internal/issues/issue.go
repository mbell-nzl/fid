package issues

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

const fidVersion = "0.1.0"

type CreateOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Kind          string
	Title         string
	ActorRole     string
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
	return fmt.Sprintf("kind %q is not yet supported in Phase 1; use --kind task", err.Kind)
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
}

type AcceptanceCriterion struct {
	ID                  string         `json:"id"`
	Ordinal             int            `json:"ordinal"`
	Text                string         `json:"text"`
	LatestVerifications map[string]any `json:"latest_verifications"`
}

func Create(options CreateOptions) (CreateResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return CreateResult{}, fmt.Errorf("create issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return CreateResult{}, fmt.Errorf("create issue: migrations dir is required")
	}
	if strings.TrimSpace(options.Title) == "" {
		return CreateResult{}, fmt.Errorf("create issue: title is required")
	}
	if err := policy.RequireKnownRole(options.WorkspaceRoot, options.ActorRole); err != nil {
		return CreateResult{}, err
	}

	kind := strings.TrimSpace(options.Kind)
	if kind != "task" {
		return CreateResult{}, &UnsupportedKindError{Kind: kind}
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return CreateResult{}, err
	}

	projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return CreateResult{}, err
	}

	var result CreateResult
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   fmt.Sprintf("fid issue create --kind %s --title %s --by %s", kind, options.Title, options.ActorRole),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		issueID, err := journal.GenerateID(projectConfig.Prefix, func(candidate string) (bool, error) {
			return issueIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		timestamp := time.Now().UTC().Format(time.RFC3339)
		event := journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.created",
			Actor: journal.Actor{
				Role: options.ActorRole,
			},
			Command: journal.Command{
				Argv:       []string{"fid", "issue", "create", "--kind", kind, "--title", options.Title, "--by", options.ActorRole},
				CWD:        ".",
				FIDVersion: fidVersion,
			},
			Effects: []journal.Effect{
				{
					Op:    "insert",
					Table: "issues",
					PK: map[string]any{
						"id": issueID,
					},
					Values: map[string]any{
						"kind":        kind,
						"title":       options.Title,
						"workflow_id": "task-default",
						"status":      "todo",
						"owner_role":  "developer",
						"created_at":  timestamp,
						"updated_at":  timestamp,
					},
				},
			},
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}

		eventHash, err := journal.HashEvent(appended)
		if err != nil {
			return fmt.Errorf("hash issue event %q: %w", appended.EventID, err)
		}

		journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
		journalPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
		info, err := os.Stat(journalPath)
		if err != nil {
			return fmt.Errorf("stat issue journal %q: %w", journalPath, err)
		}

		if err := journal.ProjectEvent(conn, journal.ProjectionInput{
			Event:     appended,
			EventHash: eventHash,
			File:      journalFile,
			Line:      countJournalLines(journalPath),
			Offset:    info.Size(),
		}); err != nil {
			return err
		}

		result = CreateResult{
			ID:    issueID,
			Kind:  kind,
			Title: options.Title,
		}
		return nil
	})
	if err != nil {
		return CreateResult{}, err
	}

	return result, nil
}

func Get(options GetOptions) (Issue, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return Issue{}, fmt.Errorf("get issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return Issue{}, fmt.Errorf("get issue: migrations dir is required")
	}
	if strings.TrimSpace(options.ID) == "" {
		return Issue{}, fmt.Errorf("get issue: issue id is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return Issue{}, err
	}

	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return Issue{}, err
	}
	defer closeIssueDB(conn)

	row := conn.QueryRow(`
		SELECT id, kind, title, status, workflow_id, parent_id, priority, sequence,
		       spec_md_path, est_tokens, owner_role, owner_actor, in_inbox, created_at, updated_at
		FROM issues
		WHERE id = ?
	`, options.ID)

	var (
		record                           Issue
		parentID, specMDPath, ownerRole  sql.NullString
		ownerActor, createdAt, updatedAt sql.NullString
		priority, sequence, estTokens    sql.NullInt64
	)
	if err := row.Scan(
		&record.ID,
		&record.Kind,
		&record.Title,
		&record.Status,
		&record.WorkflowID,
		&parentID,
		&priority,
		&sequence,
		&specMDPath,
		&estTokens,
		&ownerRole,
		&ownerActor,
		&record.InInbox,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Issue{}, &NotFoundError{IssueID: options.ID}
		}
		return Issue{}, fmt.Errorf("query issue %q: %w", options.ID, err)
	}

	record.StatusCategory = statusCategory(record.Status)
	record.Labels = []string{}
	record.AcceptanceCriteria = []AcceptanceCriterion{}
	record.ConcernsRaised = false
	record.OpenConcernsCount = 0
	record.ParentID = nullStringPointer(parentID)
	record.SpecMDPath = nullStringPointer(specMDPath)
	record.OwnerRole = nullStringPointer(ownerRole)
	record.OwnerActor = nullStringPointer(ownerActor)
	record.CreatedAt = nullStringPointer(createdAt)
	record.UpdatedAt = nullStringPointer(updatedAt)
	record.Priority = nullIntPointer(priority)
	record.Sequence = nullIntPointer(sequence)
	record.EstTokens = nullIntPointer(estTokens)

	return record, nil
}

func issueIDExists(conn *sql.DB, issueID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM issues WHERE id = ? LIMIT 1`, issueID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("query issue id %q: %w", issueID, err)
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
