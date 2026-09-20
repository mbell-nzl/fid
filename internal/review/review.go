// Package review owns immutable Review Record admission.
package review

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type RecordRequest struct {
	RunID string `json:"run_id"`
}

type RecordOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	ActorID       string
}

type NotIndependentError struct {
	ReviewActorID      string
	ConflictingType    string
	ConflictingID      string
	ConflictingActorID string
}

func (err *NotIndependentError) Error() string { return "review actor is not independent" }

// Record currently establishes the durable independence admission boundary.
// Subsequent Task 4 slices append the immutable Review Event after this check.
func Record(options RecordOptions, request RecordRequest) error {
	if options.WorkspaceRoot == "" || options.MigrationsDir == "" || options.ActorID == "" || request.RunID == "" {
		return fmt.Errorf("review record input is incomplete")
	}
	return journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid review record"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return fmt.Errorf("open review projection: %w", err)
		}
		defer func() { _ = conn.Close() }()
		return checkIndependence(conn, options.ActorID, request.RunID)
	})
}

func checkIndependence(conn *sql.DB, actorID, runID string) error {
	var issueID, envelopeID, claimID string
	var revision int
	err := conn.QueryRow(`SELECT issue_id, envelope_id, envelope_revision, claim_id FROM runs WHERE id = ?`, runID).Scan(&issueID, &envelopeID, &revision, &claimID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("review Run %q was not found", runID)
	}
	if err != nil {
		return fmt.Errorf("load review Run: %w", err)
	}
	var conflictingID string
	err = conn.QueryRow(`SELECT id FROM claims WHERE issue_id = ? AND envelope_id = ? AND envelope_revision = ? AND mode = 'mutation' AND actor_id = ? ORDER BY acquired_at, id LIMIT 1`, issueID, envelopeID, revision, actorID).Scan(&conflictingID)
	if err == nil {
		return &NotIndependentError{ReviewActorID: actorID, ConflictingType: "claim", ConflictingID: conflictingID, ConflictingActorID: actorID}
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("load mutating Claim history: %w", err)
	}
	err = conn.QueryRow(`SELECT r.id FROM runs r JOIN claims c ON c.id = r.claim_id WHERE r.issue_id = ? AND r.envelope_id = ? AND r.envelope_revision = ? AND c.mode = 'mutation' AND r.actor_id = ? ORDER BY r.started_at, r.id LIMIT 1`, issueID, envelopeID, revision, actorID).Scan(&conflictingID)
	if err == nil {
		return &NotIndependentError{ReviewActorID: actorID, ConflictingType: "run", ConflictingID: conflictingID, ConflictingActorID: actorID}
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("load mutating Run history: %w", err)
	}
	return nil
}
