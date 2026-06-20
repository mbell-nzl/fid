package issues

import (
	"database/sql"
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

type RaiseConcernOptions struct {
	WorkspaceRoot        string
	MigrationsDir        string
	IssueID              string
	ActorRole            string
	RaisedInCompletionID string
	Summary              string
	BodyMDPath           *string
}

type RaiseConcernResult struct {
	ID                   string
	IssueID              string
	RaisedInCompletionID string
	RaisedByRole         string
	Summary              string
	BodyMDPath           *string
	Status               string
	CreatedAt            string
	UpdatedAt            string
}

func RaiseConcern(options RaiseConcernOptions) (RaiseConcernResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return RaiseConcernResult{}, fmt.Errorf("raise concern: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return RaiseConcernResult{}, fmt.Errorf("raise concern: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return RaiseConcernResult{}, fmt.Errorf("raise concern: issue id is required")
	}
	completionID := strings.TrimSpace(options.RaisedInCompletionID)
	if completionID == "" {
		return RaiseConcernResult{}, fmt.Errorf("raise concern: completion id is required")
	}
	summary := strings.TrimSpace(options.Summary)
	if summary == "" {
		return RaiseConcernResult{}, fmt.Errorf("raise concern: summary is required")
	}
	if err := policy.RequireKnownRole(options.WorkspaceRoot, options.ActorRole); err != nil {
		return RaiseConcernResult{}, err
	}

	bodyMDPath := trimmedOptionalString(options.BodyMDPath)
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return RaiseConcernResult{}, err
	}

	commandArgv := issueConcernRaiseArgv(issueID, options.ActorRole, completionID, summary, bodyMDPath)
	var result RaiseConcernResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		return raiseConcernLocked(lock, options, issueID, completionID, summary, bodyMDPath, commandArgv, &result)
	})
	if err != nil {
		return RaiseConcernResult{}, err
	}

	return result, nil
}

func raiseConcernLocked(
	lock *journal.WriteLock,
	options RaiseConcernOptions,
	issueID, completionID, summary string,
	bodyMDPath *string,
	commandArgv []string,
	result *RaiseConcernResult,
) error {
	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return err
	}
	defer closeIssueDB(conn)

	if _, found, err := issueKindForID(conn, issueID); err != nil {
		return err
	} else if !found {
		return &NotFoundError{IssueID: issueID}
	}
	if err := requireCompletionBelongsToIssue(conn, completionID, issueID); err != nil {
		return err
	}

	concernID, err := journal.GenerateID("CON", func(candidate string) (bool, error) {
		return concernIDExists(conn, candidate)
	})
	if err != nil {
		return err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	event := issueConcernRaiseEvent(concernID, issueID, completionID, options.ActorRole, summary, bodyMDPath, timestamp, commandArgv)
	if err := journal.DryRunProjectEvent(conn, event); err != nil {
		return err
	}

	appended, err := journal.AppendEvent(lock, event)
	if err != nil {
		return err
	}
	if err := projectIssueEvent(conn, options.WorkspaceRoot, appended); err != nil {
		return err
	}

	*result = RaiseConcernResult{
		ID:                   concernID,
		IssueID:              issueID,
		RaisedInCompletionID: completionID,
		RaisedByRole:         options.ActorRole,
		Summary:              summary,
		BodyMDPath:           bodyMDPath,
		Status:               "open",
		CreatedAt:            timestamp,
		UpdatedAt:            timestamp,
	}
	return nil
}

func requireCompletionBelongsToIssue(conn *sql.DB, completionID, issueID string) error {
	var completionIssueID string
	err := conn.QueryRow(`SELECT issue_id FROM task_completions WHERE id = ? LIMIT 1`, completionID).Scan(&completionIssueID)
	if err == nil {
		if completionIssueID == issueID {
			return nil
		}
		return fmt.Errorf(
			"raise concern: completion %q does not belong to issue %q",
			completionID,
			issueID,
		)
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("raise concern: completion %q not found", completionID)
	}
	return fmt.Errorf("query completion %q: %w", completionID, err)
}

func concernIDExists(conn *sql.DB, concernID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM concerns WHERE id = ? LIMIT 1`, concernID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query concern id %q: %w", concernID, err)
}

func issueConcernRaiseEvent(
	concernID, issueID, completionID, actorRole, summary string,
	bodyMDPath *string,
	timestamp string,
	commandArgv []string,
) journal.Event {
	values := map[string]any{
		"issue_id":                issueID,
		"raised_in_completion_id": completionID,
		"raised_by_role":          actorRole,
		"summary":                 summary,
		"status":                  "open",
		"created_at":              timestamp,
		"updated_at":              timestamp,
	}
	if bodyMDPath != nil {
		values["body_md_path"] = *bodyMDPath
	}

	return journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "issue.concern_raised",
		Actor: journal.Actor{
			Role: actorRole,
		},
		Command: journal.Command{
			Argv:       commandArgv,
			CWD:        ".",
			FIDVersion: fidVersion,
		},
		Effects: []journal.Effect{
			{
				Op:     "insert",
				Table:  "concerns",
				PK:     map[string]any{"id": concernID},
				Values: values,
			},
		},
	}
}

func projectIssueEvent(conn *sql.DB, workspaceRoot string, appended journal.Event) error {
	eventHash, err := journal.HashEvent(appended)
	if err != nil {
		return fmt.Errorf("hash issue event %q: %w", appended.EventID, err)
	}

	journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return fmt.Errorf("stat issue journal %q: %w", journalPath, err)
	}

	return journal.ProjectEvent(conn, journal.ProjectionInput{
		Event:     appended,
		EventHash: eventHash,
		File:      journalFile,
		Line:      countJournalLines(journalPath),
		Offset:    info.Size(),
	})
}

func issueConcernRaiseArgv(issueID, actorRole, completionID, summary string, bodyMDPath *string) []string {
	argv := []string{
		"fid", "concern", "raise",
		"--on", issueID,
		"--by", actorRole,
		"--in-completion", completionID,
		"--summary", summary,
	}
	if bodyMDPath != nil {
		argv = append(argv, "--body-md", *bodyMDPath)
	}
	return argv
}
