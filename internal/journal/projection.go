package journal

import (
	"fmt"
	"os"
	"path/filepath"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

type EnsureProjectionResult struct {
	Decision    fiddb.ProjectionFreshness
	ReplayStats ReplayStats
}

func EnsureProjectionFresh(workspaceRoot, migrationsDir string) (EnsureProjectionResult, error) {
	return ensureProjectionFresh(workspaceLedgerDir(workspaceRoot), migrationsDir, nil)
}

func EnsureProjectionFreshWithLock(lock *WriteLock, migrationsDir string) (EnsureProjectionResult, error) {
	if lock == nil {
		return EnsureProjectionResult{}, fmt.Errorf("ensure projection fresh with lock: lock is nil")
	}
	return ensureProjectionFresh(lock.ledgerDir, migrationsDir, lock)
}

func ensureProjectionFresh(ledgerDir, migrationsDir string, heldLock *WriteLock) (EnsureProjectionResult, error) {
	decision, meta, err := readProjectionFreshness(ledgerDir)
	if err != nil {
		return EnsureProjectionResult{}, err
	}
	return refreshProjection(ledgerDir, migrationsDir, decision, meta, heldLock)
}

func readProjectionFreshness(ledgerDir string) (fiddb.ProjectionFreshness, fiddb.ProjectionMeta, error) {
	dbPath := filepath.Join(ledgerDir, "fid.sqlite")
	journalDir := filepath.Join(ledgerDir, "journal")

	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			return fiddb.ProjectionRebuild, fiddb.ProjectionMeta{}, nil
		}
		return "", fiddb.ProjectionMeta{}, fmt.Errorf("stat sqlite database %q: %w", dbPath, err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		return "", fiddb.ProjectionMeta{}, err
	}
	defer closeReplayDB(conn)

	meta, metaExists, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		return "", fiddb.ProjectionMeta{}, err
	}

	inspection, err := inspectJournal(journalDir, meta.JournalFile, meta.JournalOffset)
	if err != nil {
		return "", fiddb.ProjectionMeta{}, err
	}

	decision := fiddb.DecideProjectionFreshness(fiddb.ProjectionFreshnessInput{
		MetaExists:                    metaExists,
		ReplayStatus:                  meta.ReplayStatus,
		StoredProjectionSchemaVersion: meta.ProjectionSchemaVersion,
		CurrentProjectionSchema:       fiddb.ProjectionSchemaVersion,
		CheckpointExists:              inspection.checkpointSeen,
		CheckpointMatches: metaExists &&
			inspection.checkpointSeen &&
			meta.LastEventID == inspection.checkpoint.EventID &&
			meta.LastEventHash == inspection.checkpoint.EventHash,
		JournalHasNewerEvents: inspection.tail.HasEvents &&
			(!metaExists ||
				meta.JournalFile != inspection.tail.File ||
				meta.JournalOffset != inspection.tail.Offset),
	})
	return decision, meta, nil
}

func refreshProjection(ledgerDir, migrationsDir string, decision fiddb.ProjectionFreshness, meta fiddb.ProjectionMeta, heldLock *WriteLock) (EnsureProjectionResult, error) {
	if decision == fiddb.ProjectionFresh {
		return EnsureProjectionResult{Decision: decision}, nil
	}
	if heldLock == nil {
		var result EnsureProjectionResult
		err := WithLedgerLock(ledgerDir, LockOwner{Command: "fid catch-up"}, func(lock *WriteLock) error {
			// A writer or another reader may have refreshed the projection while
			// we waited. Read its state again before replaying or replacing it.
			var err error
			result, err = ensureProjectionFresh(ledgerDir, migrationsDir, lock)
			return err
		})
		return result, err
	}

	switch decision {
	case fiddb.ProjectionCatchUp:
		stats, err := catchUpWorkspaceLocked(heldLock, migrationsDir, meta)
		if err != nil {
			return EnsureProjectionResult{}, err
		}
		return EnsureProjectionResult{Decision: decision, ReplayStats: stats}, nil
	case fiddb.ProjectionRebuild:
		stats, err := rebuildWorkspaceLocked(heldLock, migrationsDir)
		if err != nil {
			return EnsureProjectionResult{}, err
		}
		return EnsureProjectionResult{Decision: decision, ReplayStats: stats}, nil
	default:
		return EnsureProjectionResult{}, fmt.Errorf("ensure projection fresh: unsupported decision %q", decision)
	}
}

func RebuildWorkspace(workspaceRoot, migrationsDir string) (ReplayStats, error) {
	return RebuildLedger(workspaceLedgerDir(workspaceRoot), migrationsDir)
}

// RebuildLedger is RebuildWorkspace for the ledger at ledgerDir rather than
// the project at a root, for a ledger being staged, which has no .git to find
// it by.
func RebuildLedger(ledgerDir, migrationsDir string) (ReplayStats, error) {
	return rebuildWorkspaceWithOwner(ledgerDir, migrationsDir, LockOwner{
		Command: "fid rebuild",
	})
}

func rebuildWorkspaceWithOwner(ledgerDir, migrationsDir string, owner LockOwner) (ReplayStats, error) {
	var stats ReplayStats
	err := tryWithLedgerLock(ledgerDir, owner, func(lock *WriteLock) error {
		replayed, err := rebuildWorkspaceLocked(lock, migrationsDir)
		stats = replayed
		return err
	})
	if err != nil {
		return ReplayStats{}, err
	}

	return stats, nil
}

func rebuildWorkspaceLocked(lock *WriteLock, migrationsDir string) (ReplayStats, error) {
	dbPath := filepath.Join(lock.ledgerDir, "fid.sqlite")
	temporary, err := os.CreateTemp(filepath.Dir(dbPath), ".fid-rebuild-*.sqlite")
	if err != nil {
		return ReplayStats{}, fmt.Errorf("create rebuild projection: %w", err)
	}
	temporaryPath := temporary.Name()
	if err := temporary.Close(); err != nil {
		return ReplayStats{}, fmt.Errorf("close rebuild projection: %w", err)
	}
	if err := os.Remove(temporaryPath); err != nil {
		return ReplayStats{}, fmt.Errorf("prepare rebuild projection: %w", err)
	}
	defer func() { _ = os.Remove(temporaryPath) }()
	if err := fiddb.MigrateUp(temporaryPath, migrationsDir); err != nil {
		return ReplayStats{}, err
	}

	conn, err := fiddb.OpenForRebuild(temporaryPath)
	if err != nil {
		return ReplayStats{}, err
	}

	replayed, err := replayRebuild(conn, filepath.Join(lock.ledgerDir, "journal"))
	closeReplayDB(conn)
	if err != nil {
		return ReplayStats{}, err
	}
	if err := os.Rename(temporaryPath, dbPath); err != nil {
		return ReplayStats{}, fmt.Errorf("install rebuilt projection: %w", err)
	}

	return replayed, nil
}

func catchUpWorkspaceLocked(lock *WriteLock, migrationsDir string, meta fiddb.ProjectionMeta) (ReplayStats, error) {
	dbPath := filepath.Join(lock.ledgerDir, "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, migrationsDir); err != nil {
		return ReplayStats{}, err
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		return ReplayStats{}, err
	}
	defer closeReplayDB(conn)

	replayed, err := ReplayJournal(conn, filepath.Join(lock.ledgerDir, "journal"), &ReplayOptions{
		StartFile:   meta.JournalFile,
		StartOffset: meta.JournalOffset,
	})
	if err != nil {
		return ReplayStats{}, err
	}

	return replayed, nil
}
