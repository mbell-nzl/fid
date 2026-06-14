package issues

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type AddAcceptanceCriterionOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	Text          string
	Ordinal       int
	ActorRole     string
}

type AddAcceptanceCriterionResult struct {
	ID      string
	IssueID string
	Text    string
	Ordinal int
}

func AddAcceptanceCriterion(options AddAcceptanceCriterionOptions) (AddAcceptanceCriterionResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: issue id is required")
	}
	text := strings.TrimSpace(options.Text)
	if text == "" {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: text is required")
	}
	if options.Ordinal <= 0 {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: ordinal must be a positive integer")
	}
	if err := policy.RequireKnownRole(options.WorkspaceRoot, options.ActorRole); err != nil {
		return AddAcceptanceCriterionResult{}, err
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return AddAcceptanceCriterionResult{}, err
	}

	commandArgv := issueACAddArgv(issueID, text, options.Ordinal, options.ActorRole)
	var result AddAcceptanceCriterionResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		if err := requireAcceptanceCriteriaIssue(conn, issueID); err != nil {
			return err
		}

		criterionID, err := journal.GenerateID("AC", func(candidate string) (bool, error) {
			return acceptanceCriterionIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		event := journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.ac_added",
			Actor: journal.Actor{
				Role: options.ActorRole,
			},
			Command: journal.Command{
				Argv:       commandArgv,
				CWD:        ".",
				FIDVersion: fidVersion,
			},
			Effects: []journal.Effect{
				{
					Op:    "insert",
					Table: "acceptance_criteria",
					PK:    map[string]any{"id": criterionID},
					Values: map[string]any{
						"issue_id": issueID,
						"ordinal":  options.Ordinal,
						"text":     text,
					},
				},
			},
		}

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
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

		result = AddAcceptanceCriterionResult{
			ID:      criterionID,
			IssueID: issueID,
			Text:    text,
			Ordinal: options.Ordinal,
		}
		return nil
	})
	if err != nil {
		return AddAcceptanceCriterionResult{}, err
	}

	return result, nil
}

func requireAcceptanceCriteriaIssue(conn *sql.DB, issueID string) error {
	kind, found, err := issueKindForID(conn, issueID)
	if err != nil {
		return err
	}
	if !found {
		return &NotFoundError{IssueID: issueID}
	}
	if kind == "task" || kind == "subtask" {
		return nil
	}
	return fmt.Errorf(
		"add acceptance criterion: issue %q is kind %q; acceptance criteria require task or subtask",
		issueID,
		kind,
	)
}

func acceptanceCriterionIDExists(conn *sql.DB, criterionID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM acceptance_criteria WHERE id = ? LIMIT 1`, criterionID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query acceptance criterion id %q: %w", criterionID, err)
}

func issueACAddArgv(issueID, text string, ordinal int, actorRole string) []string {
	return []string{
		"fid", "issue", "ac", "add", issueID,
		"--text", text,
		"--ordinal", fmt.Sprintf("%d", ordinal),
		"--by", actorRole,
	}
}
