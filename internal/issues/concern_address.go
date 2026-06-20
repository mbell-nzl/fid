package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type AddressConcernOptions struct {
	WorkspaceRoot           string
	MigrationsDir           string
	ConcernID               string
	ActorRole               string
	AddressedInCompletionID string
}

type AddressConcernResult struct {
	ID                     string
	IssueID                string
	AddressingCompletionID string
	Status                 string
	UpdatedAt              string
}

type concernAddressState struct {
	ID        string
	IssueID   string
	Status    string
	UpdatedAt string
}

func AddressConcern(options AddressConcernOptions) (AddressConcernResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return AddressConcernResult{}, fmt.Errorf("address concern: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return AddressConcernResult{}, fmt.Errorf("address concern: migrations dir is required")
	}
	concernID := strings.TrimSpace(options.ConcernID)
	if concernID == "" {
		return AddressConcernResult{}, fmt.Errorf("address concern: concern id is required")
	}
	completionID := strings.TrimSpace(options.AddressedInCompletionID)
	if completionID == "" {
		return AddressConcernResult{}, fmt.Errorf("address concern: completion id is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionConcernAddress); err != nil {
		return AddressConcernResult{}, err
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return AddressConcernResult{}, err
	}

	commandArgv := issueConcernAddressArgv(concernID, completionID, options.ActorRole)
	var result AddressConcernResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		return addressConcernLocked(lock, options, concernID, completionID, commandArgv, &result)
	})
	if err != nil {
		return AddressConcernResult{}, err
	}
	return result, nil
}

func addressConcernLocked(
	lock *journal.WriteLock,
	options AddressConcernOptions,
	concernID, completionID string,
	commandArgv []string,
	result *AddressConcernResult,
) error {
	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return err
	}
	defer closeIssueDB(conn)

	current, err := loadConcernAddressState(conn, concernID)
	if err != nil {
		return err
	}
	if current.Status != "open" {
		return fmt.Errorf("address concern: concern %q is status %q; only open concerns can be addressed", concernID, current.Status)
	}
	if err := requireAddressCompletionBelongsToIssue(conn, completionID, current.IssueID); err != nil {
		return err
	}

	updatedAt := nextIssueUpdateTime(current.UpdatedAt, time.Now().UTC()).Format(time.RFC3339)
	event := issueConcernAddressEvent(concernID, completionID, options.ActorRole, current.Status, current.UpdatedAt, updatedAt, commandArgv)
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

	*result = AddressConcernResult{
		ID:                     concernID,
		IssueID:                current.IssueID,
		AddressingCompletionID: completionID,
		Status:                 "addressed",
		UpdatedAt:              updatedAt,
	}
	return nil
}

func loadConcernAddressState(conn *sql.DB, concernID string) (concernAddressState, error) {
	var state concernAddressState
	err := conn.QueryRow(`
		SELECT id, issue_id, status, updated_at
		FROM concerns
		WHERE id = ?
		LIMIT 1
	`, concernID).Scan(&state.ID, &state.IssueID, &state.Status, &state.UpdatedAt)
	if err == nil {
		if strings.TrimSpace(state.UpdatedAt) == "" {
			return concernAddressState{}, fmt.Errorf("address concern: concern %q missing updated_at", concernID)
		}
		return state, nil
	}
	if err == sql.ErrNoRows {
		return concernAddressState{}, fmt.Errorf("address concern: concern %q not found", concernID)
	}
	return concernAddressState{}, fmt.Errorf("query concern %q: %w", concernID, err)
}

func requireAddressCompletionBelongsToIssue(conn *sql.DB, completionID, issueID string) error {
	var completionIssueID string
	err := conn.QueryRow(`SELECT issue_id FROM task_completions WHERE id = ? LIMIT 1`, completionID).Scan(&completionIssueID)
	if err == nil {
		if completionIssueID == issueID {
			return nil
		}
		return fmt.Errorf("address concern: completion %q does not belong to issue %q", completionID, issueID)
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("address concern: completion %q not found", completionID)
	}
	return fmt.Errorf("query completion %q: %w", completionID, err)
}

func issueConcernAddressEvent(
	concernID, completionID, actorRole, priorStatus, priorUpdatedAt, updatedAt string,
	commandArgv []string,
) journal.Event {
	return journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "issue.concern_addressed",
		Actor: journal.Actor{
			Role: actorRole,
		},
		Command: journal.Command{
			Argv:       commandArgv,
			CWD:        ".",
			FIDVersion: fidVersion,
		},
		Effects: []journal.Effect{{
			Op:    "update",
			Table: "concerns",
			PK:    map[string]any{"id": concernID},
			Values: map[string]any{
				"status":                     "addressed",
				"addressed_in_completion_id": completionID,
				"updated_at":                 updatedAt,
			},
			Prior: map[string]any{
				"status":     priorStatus,
				"updated_at": priorUpdatedAt,
			},
			Preconditions: map[string]any{
				"status":     priorStatus,
				"updated_at": priorUpdatedAt,
			},
		}},
	}
}

func issueConcernAddressArgv(concernID, completionID, actorRole string) []string {
	return []string{
		"fid", "concern", "address", concernID,
		"--in-completion", completionID,
		"--by", actorRole,
	}
}
