package issues

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type UpdateOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	ID            string
	ActorRole     string
	Title         *string
	Priority      *int
	Labels        *[]string
	SpecMDPath    *string
	// OwnerRole and OwnerActor are retained only for source compatibility with
	// V1 in-process callers. Canonical Issue updates deliberately ignore them.
	OwnerRole        *string
	OwnerActor       *string
	ExpectedRevision *int
}

type UpdateResult struct {
	ID         string
	Revision   int
	Changes    []FieldChange
	EventID    string
	OccurredAt string
}

type FieldChange struct {
	Field string
	From  any
	To    any
}

type issueUpdateState struct {
	ID         string
	Title      string
	Priority   *int
	Labels     []string
	SpecMDPath *string
	UpdatedAt  string
	Revision   int
}

func Update(options UpdateOptions) (UpdateResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return UpdateResult{}, fmt.Errorf("update issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return UpdateResult{}, fmt.Errorf("update issue: migrations dir is required")
	}
	if strings.TrimSpace(options.ID) == "" {
		return UpdateResult{}, fmt.Errorf("update issue: issue id is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionIssueUpdate); err != nil {
		return UpdateResult{}, err
	}

	request, err := normalizeUpdateRequest(options)
	if err != nil {
		return UpdateResult{}, err
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return UpdateResult{}, err
	}

	commandArgv := issueUpdateArgv(options.ID, options.ActorRole, request, options.ExpectedRevision)
	var result UpdateResult
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

		current, err := loadIssueUpdateState(conn, options.ID)
		if err != nil {
			return err
		}
		if options.ExpectedRevision != nil && *options.ExpectedRevision != current.Revision {
			return &StateConflictError{IssueID: current.ID, ExpectedRevision: *options.ExpectedRevision, ActualRevision: current.Revision}
		}
		if request.labels != nil {
			if err := requireExistingLabels(conn, *request.labels); err != nil {
				return fmt.Errorf("%s", strings.Replace(err.Error(), "create issue:", "update issue:", 1))
			}
		}

		scalarValues, scalarPreconditions, scalarPrior, changes := buildIssueUpdateChanges(current, request)
		labelEffects, labelChanged, labelChange := buildLabelEffects(current.ID, current.Labels, request.labels)
		if labelChanged {
			changes = append(changes, labelChange)
		}
		sort.Slice(changes, func(left, right int) bool {
			return changes[left].Field < changes[right].Field
		})
		if len(changes) == 0 {
			return fmt.Errorf("update issue: at least one field must change")
		}

		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return activityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}
		occurredAt := nextIssueUpdateTime(current.UpdatedAt, time.Now().UTC())

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.updated",
			OccurredAt:  occurredAt,
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
		scalarValues["updated_at"] = timestamp
		scalarPreconditions["updated_at"] = current.UpdatedAt
		scalarPrior["updated_at"] = current.UpdatedAt

		effects := []journal.Effect{
			{
				Op:            "update",
				Table:         "issues",
				PK:            map[string]any{"id": current.ID},
				Values:        scalarValues,
				Prior:         scalarPrior,
				Preconditions: scalarPreconditions,
			},
		}
		revisionEffects, err := BuildIssueRevisionEffects([]IssueRevision{{
			IssueID:  current.ID,
			Observed: current.Revision,
		}})
		if err != nil {
			return err
		}
		effects = append(effects, revisionEffects...)
		effects = append(effects, labelEffects...)
		effects = append(effects, journal.Effect{
			Op:    "activity",
			Table: "activity_log",
			PK:    map[string]any{"id": activityID},
			Values: map[string]any{
				"journal_event_id": event.EventID,
				"ref_type":         "issue",
				"ref_id":           current.ID,
				"event":            "issue_updated",
				"details": map[string]any{
					"changes": activityDetailChanges(changes),
				},
				"actor":       options.ActorRole,
				"occurred_at": timestamp,
			},
		})
		event.Effects = effects

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

		result = UpdateResult{
			ID:         current.ID,
			Revision:   current.Revision + 1,
			Changes:    changes,
			EventID:    appended.EventID,
			OccurredAt: appended.OccurredAt.UTC().Format(time.RFC3339),
		}
		return nil
	})
	if err != nil {
		return UpdateResult{}, err
	}

	return result, nil
}

func loadIssueUpdateState(conn *sql.DB, issueID string) (issueUpdateState, error) {
	row := conn.QueryRow(`
		SELECT id, title, priority, spec_md_path, updated_at, revision
		FROM issues
		WHERE id = ?
	`, issueID)

	var (
		state      issueUpdateState
		priority   sql.NullInt64
		specMDPath sql.NullString
		updatedAt  sql.NullString
	)
	if err := row.Scan(
		&state.ID,
		&state.Title,
		&priority,
		&specMDPath,
		&updatedAt,
		&state.Revision,
	); err != nil {
		if err == sql.ErrNoRows {
			return issueUpdateState{}, &NotFoundError{IssueID: issueID}
		}
		return issueUpdateState{}, fmt.Errorf("query issue %q: %w", issueID, err)
	}
	state.Priority = nullIntPointer(priority)
	state.SpecMDPath = nullStringPointer(specMDPath)
	if !updatedAt.Valid || strings.TrimSpace(updatedAt.String) == "" {
		return issueUpdateState{}, fmt.Errorf("query issue %q: updated_at is required", issueID)
	}
	state.UpdatedAt = updatedAt.String

	labels, err := issueLabels(conn, issueID)
	if err != nil {
		return issueUpdateState{}, err
	}
	state.Labels = labels
	return state, nil
}
