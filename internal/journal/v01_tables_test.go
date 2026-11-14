package journal

import (
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

	if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}

	conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if _, err := conn.Exec(`UPDATE projection_meta SET projection_schema_version = ? WHERE key = 'default'`, 1); err != nil {
		t.Fatalf("downgrade projection_meta schema version: %v", err)
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

	rebuiltConn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() rebuilt error = %v", err)
	}
	defer closeProjectDB(t, rebuiltConn)

	meta, found, err := fiddb.ReadProjectionMeta(rebuiltConn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta() error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta() found = false, want true")
	}
	if meta.ProjectionSchemaVersion != fiddb.ProjectionSchemaVersion {
		t.Fatalf("projection_schema_version = %d, want %d", meta.ProjectionSchemaVersion, fiddb.ProjectionSchemaVersion)
	}
}
