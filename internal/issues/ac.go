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
	"github.com/mbell-nzl/fid/internal/validation"
)

const AcceptanceCriteriaLimit = 1000

type AddAcceptanceCriterionOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	Text          string
	Ordinal       int
	ActorRole     string
}

type AddAcceptanceCriterionResult struct {
	ID            string
	IssueID       string
	Text          string
	Ordinal       int
	IssueRevision int
	EventID       string
	OccurredAt    string
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
	if strings.TrimSpace(options.Text) == "" {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: text is required")
	}
	text, err := validation.NormalizeText("acceptance criterion text", options.Text, 4096)
	if err != nil {
		return AddAcceptanceCriterionResult{}, err
	}
	if options.Ordinal <= 0 {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: ordinal must be a positive integer")
	}
	if options.Ordinal > 2147483647 {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: ordinal must be between 1 and 2147483647")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionIssueACAdd); err != nil {
		return AddAcceptanceCriterionResult{}, err
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return AddAcceptanceCriterionResult{}, err
	}

	commandArgv := issueACAddArgv(issueID, text, options.Ordinal, options.ActorRole)
	var result AddAcceptanceCriterionResult
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		issueRevision, err := acceptanceCriteriaIssueRevision(conn, issueID)
		if err != nil {
			return err
		}
		var existing string
		err = conn.QueryRow(`SELECT id FROM acceptance_criteria WHERE issue_id = ? AND ordinal = ?`, issueID, options.Ordinal).Scan(&existing)
		if err == nil {
			return &DuplicateAcceptanceCriterionError{IssueID: issueID, Ordinal: options.Ordinal}
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("query acceptance criterion ordinal: %w", err)
		}
		var criterionCount int
		if err := conn.QueryRow(`SELECT COUNT(*) FROM acceptance_criteria WHERE issue_id = ?`, issueID).Scan(&criterionCount); err != nil {
			return fmt.Errorf("count acceptance criteria for issue %q: %w", issueID, err)
		}
		if criterionCount >= AcceptanceCriteriaLimit {
			return &ResultLimitExceededError{Collection: "acceptance_criteria", Limit: AcceptanceCriteriaLimit, Observed: criterionCount + 1}
		}

		criterionID, err := journal.GenerateID("AC", func(candidate string) (bool, error) {
			return acceptanceCriterionIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}
		activityID, err := nextIssueActivityID(conn)
		if err != nil {
			return err
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
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
		})
		if err != nil {
			return err
		}

		timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
		event.Effects = []journal.Effect{
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
		}
		revisionEffects, err := BuildIssueRevisionEffects([]IssueRevision{{IssueID: issueID, Observed: issueRevision}})
		if err != nil {
			return err
		}
		event.Effects = append(event.Effects, revisionEffects...)
		event.Effects = append(event.Effects, issueActivityEffect(
			activityID,
			event.EventID,
			options.ActorRole,
			timestamp,
			issueID,
			"ac_added",
			map[string]any{
				"ac_id":   criterionID,
				"ordinal": options.Ordinal,
			},
		))

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
			ID:            criterionID,
			IssueID:       issueID,
			Text:          text,
			Ordinal:       options.Ordinal,
			IssueRevision: issueRevision + 1,
			EventID:       appended.EventID,
			OccurredAt:    appended.OccurredAt.UTC().Format(time.RFC3339Nano),
		}
		return nil
	})
	if err != nil {
		return AddAcceptanceCriterionResult{}, err
	}

	return result, nil
}

func acceptanceCriteriaIssueRevision(conn *sql.DB, issueID string) (int, error) {
	var kind string
	var revision int
	err := conn.QueryRow(`SELECT kind, revision FROM issues WHERE id = ?`, issueID).Scan(&kind, &revision)
	if err == sql.ErrNoRows {
		return 0, &NotFoundError{IssueID: issueID}
	}
	if err != nil {
		return 0, fmt.Errorf("query issue %q: %w", issueID, err)
	}
	if kind == "task" || kind == "subtask" {
		return revision, nil
	}
	return 0, fmt.Errorf(
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
