package journal

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestReplayJournalAppliesGoldenV1JournalInFileOrder(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)

	stats, err := ReplayJournal(conn, filepath.Join(workspace, ".fid", "journal"), nil)
	if err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}

	if stats.EventsReplayed != 2 {
		t.Fatalf("replayed events = %d, want %d", stats.EventsReplayed, 2)
	}
	if stats.FilesReplayed != 1 {
		t.Fatalf("replayed files = %d, want %d", stats.FilesReplayed, 1)
	}
	if stats.UpgradedEvents != 0 {
		t.Fatalf("upgraded events = %d, want %d", stats.UpgradedEvents, 0)
	}

	var roleCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM roles WHERE name = ?`, "developer").Scan(&roleCount); err != nil {
		t.Fatalf("query roles: %v", err)
	}
	if roleCount != 1 {
		t.Fatalf("projected role count = %d, want %d", roleCount, 1)
	}

	var title string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, "FID-a1b2c3").Scan(&title); err != nil {
		t.Fatalf("query issues: %v", err)
	}
	if title != "Implement replay" {
		t.Fatalf("projected issue title = %q, want %q", title, "Implement replay")
	}
}

func TestRebuildWorkspaceRecreatesIdenticalStateAfterDeletingDatabase(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	conn := openReplayDB(t, dbPath)

	stats, err := ReplayJournal(conn, filepath.Join(workspace, ".fid", "journal"), nil)
	if err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if stats.EventsReplayed != 2 {
		t.Fatalf("initial replayed events = %d, want %d", stats.EventsReplayed, 2)
	}

	before := snapshotReplayState(t, conn)
	closeProjectDB(t, conn)

	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer")
	if err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}
	if rebuildStats.EventsReplayed != 2 {
		t.Fatalf("rebuilt events = %d, want %d", rebuildStats.EventsReplayed, 2)
	}

	rebuiltConn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() rebuilt error = %v", err)
	}
	defer closeProjectDB(t, rebuiltConn)

	after := snapshotReplayState(t, rebuiltConn)
	if before != after {
		t.Fatalf("rebuilt state = %q, want %q", after, before)
	}
}

func TestReplayJournalReportsMalformedTrailingLine(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "corrupt-tail.jsonl", "2026-06-05.jsonl")

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)

	_, err := ReplayJournal(conn, filepath.Join(workspace, ".fid", "journal"), nil)
	if err == nil {
		t.Fatal("ReplayJournal() error = nil, want corruption error")
	}

	var corruptErr *CorruptJournalError
	if !errors.As(err, &corruptErr) {
		t.Fatalf("ReplayJournal() error = %T, want *CorruptJournalError", err)
	}
	if corruptErr.File != filepath.ToSlash(filepath.Join(".fid", "journal", "2026-06-05.jsonl")) {
		t.Fatalf("corrupt file = %q, want %q", corruptErr.File, filepath.ToSlash(filepath.Join(".fid", "journal", "2026-06-05.jsonl")))
	}
	if corruptErr.Line != 2 {
		t.Fatalf("corrupt line = %d, want %d", corruptErr.Line, 2)
	}
}

func TestReplayJournalStopsOnPreconditionConflict(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "conflicting-precondition.jsonl", "2026-06-05.jsonl")

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)

	_, err := ReplayJournal(conn, filepath.Join(workspace, ".fid", "journal"), nil)
	if err == nil {
		t.Fatal("ReplayJournal() error = nil, want precondition conflict")
	}

	var conflictErr *PreconditionConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("ReplayJournal() error = %T, want *PreconditionConflictError", err)
	}
	if conflictErr.EventID != "EVT-20260605-000002" {
		t.Fatalf("conflict event_id = %q, want %q", conflictErr.EventID, "EVT-20260605-000002")
	}
	if conflictErr.Table != "issues" {
		t.Fatalf("conflict table = %q, want %q", conflictErr.Table, "issues")
	}
	if conflictErr.Line != 2 {
		t.Fatalf("conflict line = %d, want %d", conflictErr.Line, 2)
	}
}

func TestEnsureProjectionFreshUsesProjectionMetaMatrix(t *testing.T) {
	t.Parallel()

	t.Run("clean match stays clean", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer"); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionFresh {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionFresh)
		}
		if result.ReplayStats.EventsReplayed != 0 {
			t.Fatalf("fresh replayed events = %d, want %d", result.ReplayStats.EventsReplayed, 0)
		}
	})

	t.Run("missing meta rebuilds", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer"); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		conn := openReplayDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
		if _, err := conn.Exec(`DELETE FROM projection_meta`); err != nil {
			t.Fatalf("delete projection_meta: %v", err)
		}
		closeProjectDB(t, conn)

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionRebuild {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
		}
		if result.ReplayStats.EventsReplayed != 2 {
			t.Fatalf("rebuilt events = %d, want %d", result.ReplayStats.EventsReplayed, 2)
		}
	})

	t.Run("schema mismatch rebuilds", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer"); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		conn := openReplayDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
		if _, err := conn.Exec(`UPDATE projection_meta SET projection_schema_version = ? WHERE key = 'default'`, 999); err != nil {
			t.Fatalf("update projection_meta schema version: %v", err)
		}
		closeProjectDB(t, conn)

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionRebuild {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
		}
	})

	t.Run("position mismatch rebuilds", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer"); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		conn := openReplayDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
		if _, err := conn.Exec(`UPDATE projection_meta SET journal_offset = ? WHERE key = 'default'`, 99999); err != nil {
			t.Fatalf("update projection_meta journal_offset: %v", err)
		}
		closeProjectDB(t, conn)

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionRebuild {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
		}
	})

	t.Run("hash mismatch rebuilds", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer"); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		conn := openReplayDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
		if _, err := conn.Exec(`UPDATE projection_meta SET last_event_hash = ? WHERE key = 'default'`, "sha256:bad"); err != nil {
			t.Fatalf("update projection_meta last_event_hash: %v", err)
		}
		closeProjectDB(t, conn)

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionRebuild {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
		}
	})

	t.Run("journal advance catches up", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer"); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		appendReplayFixtureLine(t, workspace, "replay-catchup-event.json", "2026-06-05.jsonl")

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionCatchUp {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionCatchUp)
		}
		if result.ReplayStats.EventsReplayed != 1 {
			t.Fatalf("caught-up events = %d, want %d", result.ReplayStats.EventsReplayed, 1)
		}
	})
}

func newReplayWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	mustMkdirAll(t, filepath.Join(workspace, ".fid", "journal"))
	return workspace
}

func openReplayDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()

	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	return conn
}

func replayFixturePath(t *testing.T, name string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "journal", name))
}

func writeReplayFixture(t *testing.T, workspaceRoot, fixtureName, journalFileName string) {
	t.Helper()

	payload, err := os.ReadFile(replayFixturePath(t, fixtureName))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", fixtureName, err)
	}

	journalPath := filepath.Join(workspaceRoot, ".fid", "journal", journalFileName)
	if err := os.WriteFile(journalPath, payload, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", journalPath, err)
	}
}

func appendReplayFixtureLine(t *testing.T, workspaceRoot, fixtureName, journalFileName string) {
	t.Helper()

	payload, err := os.ReadFile(replayFixturePath(t, fixtureName))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", fixtureName, err)
	}

	journalPath := filepath.Join(workspaceRoot, ".fid", "journal", journalFileName)
	file, err := os.OpenFile(journalPath, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Fatalf("OpenFile(%q) error = %v", journalPath, err)
	}
	defer closeProjectDB(t, file)

	if _, err := file.Write(append(payload, '\n')); err != nil {
		t.Fatalf("append replay fixture line: %v", err)
	}
}

func snapshotReplayState(t *testing.T, conn *sql.DB) string {
	t.Helper()

	var roleName string
	if err := conn.QueryRow(`SELECT name FROM roles ORDER BY name LIMIT 1`).Scan(&roleName); err != nil {
		t.Fatalf("query snapshot role: %v", err)
	}

	var title string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, "FID-a1b2c3").Scan(&title); err != nil {
		t.Fatalf("query snapshot issue: %v", err)
	}

	meta, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta() error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta() found = false, want true")
	}

	return strings.Join([]string{
		roleName,
		title,
		meta.JournalFile,
		meta.LastEventID,
		meta.LastEventHash,
		meta.ReplayStatus,
	}, "|")
}
