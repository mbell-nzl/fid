package journal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestProjectionRefreshRechecksAfterWriterCommits(t *testing.T) {
	t.Parallel()
	for _, path := range []string{"catch-up", "rebuild", "missing database"} {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			testProjectionRefreshAfterWriter(t, path)
		})
	}
}

func testProjectionRefreshAfterWriter(t *testing.T, path string) {
	t.Helper()
	workspace, conn := newIdempotencyWorkspace(t)
	ledgerDir, migrations := config.LedgerDir(workspace), testProjectMigrationsDir(t)
	if err := WithLock(workspace, LockOwner{Command: "test writer"}, func(lock *WriteLock) error {
		_, err := CommitEvent(lock, conn, scaleIssueCreatedEvent(1, 1, 1))
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if path == "rebuild" {
		if _, err := conn.Exec(`UPDATE projection_meta SET projection_schema_version = 999`); err != nil {
			t.Fatal(err)
		}
	}
	closeProjectDB(t, conn)
	nextIssue, wantDecision := 2, fiddb.ProjectionRebuild
	switch path {
	case "catch-up":
		if err := WithLock(workspace, LockOwner{Command: "test append"}, func(lock *WriteLock) error {
			_, err := AppendEvent(lock, scaleIssueCreatedEvent(2, 2, 1))
			return err
		}); err != nil {
			t.Fatal(err)
		}
		nextIssue, wantDecision = 3, fiddb.ProjectionCatchUp
	case "missing database":
		if err := os.Remove(filepath.Join(ledgerDir, "fid.sqlite")); err != nil {
			t.Fatal(err)
		}
	}
	// Capture the reader's decision and checkpoint before a writer repairs
	// the projection and commits more events under the write lock.
	decision, meta, err := readProjectionFreshness(ledgerDir)
	if err != nil || decision != wantDecision {
		t.Fatalf("initial decision = %q, want %q; err=%v", decision, wantDecision, err)
	}
	if err := WithLock(workspace, LockOwner{Command: "test writer"}, func(lock *WriteLock) error {
		if _, err := EnsureProjectionFreshWithLock(lock, migrations); err != nil {
			return err
		}
		writer := openReplayDB(t, filepath.Join(ledgerDir, "fid.sqlite"))
		defer closeProjectDB(t, writer)
		for i := nextIssue; i <= 3; i++ {
			if _, err := CommitEvent(lock, writer, scaleIssueCreatedEvent(i, i, 1)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	result, err := refreshProjection(ledgerDir, migrations, decision, meta, nil)
	if err != nil {
		t.Fatalf("refresh after writer committed: %v", err)
	}
	if result.Decision != fiddb.ProjectionFresh || result.ReplayStats.EventsReplayed != 0 {
		t.Fatalf("refresh = %+v, want fresh with no events replayed", result)
	}
	assertProjectionIncludesWriter(t, workspace)
}

func assertProjectionIncludesWriter(t *testing.T, workspace string) {
	t.Helper()
	conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeProjectDB(t, conn)
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM issues`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("projected issues = %d, want 3", count)
	}
	current, exists, err := fiddb.ReadProjectionMeta(conn)
	if err != nil || !exists || current.LastEventID != scaleIssueCreatedEvent(3, 3, 1).EventID {
		t.Fatalf("final checkpoint = %+v exists=%t err=%v", current, exists, err)
	}
	inspection, err := inspectJournal(filepath.Join(config.LedgerDir(workspace), "journal"), current.JournalFile, current.JournalOffset)
	if err != nil || !inspection.checkpointSeen || inspection.tail.Offset != current.JournalOffset {
		t.Fatalf("final journal inspection = %+v err=%v", inspection, err)
	}
}

func TestFreshProjectionReadDoesNotWaitForWriter(t *testing.T) {
	t.Parallel()
	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")
	migrations := testProjectMigrationsDir(t)
	if _, err := RebuildWorkspace(workspace, migrations); err != nil {
		t.Fatal(err)
	}
	if err := WithLock(workspace, LockOwner{Command: "test writer"}, func(_ *WriteLock) error {
		result, err := EnsureProjectionFresh(workspace, migrations)
		if err == nil && result.Decision != fiddb.ProjectionFresh {
			t.Errorf("freshness decision = %q, want clean", result.Decision)
		}
		return err
	}); err != nil {
		t.Fatalf("fresh read while write lock held: %v", err)
	}
}
