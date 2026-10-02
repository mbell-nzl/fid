package journal

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestRebuildReplayRollsBackFailedEventAndReleasesConnection(t *testing.T) {
	t.Parallel()
	workspace := newReplayWorkspace(t)
	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)
	conn.SetMaxOpenConns(1)
	first := testV1Event()
	first.Effects = []Effect{{Op: "insert", Table: "labels", PK: map[string]any{"name": "first"}}}
	second := first
	second.EventID = "EVT-failed-rebuild"
	second.Effects = []Effect{
		{Op: "insert", Table: "labels", PK: map[string]any{"name": "rolled-back"}},
		{Op: "insert", Table: "issues", PK: map[string]any{"id": "ISS-invalid"}, Values: map[string]any{
			"kind": "task", "title": "invalid FK", "workflow_id": "missing", "status": "todo",
		}},
	}
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-06-05.jsonl"), first, second)
	_, err := replayRebuild(conn, filepath.Join(workspace, ".fid", "journal"))
	var corrupt *CorruptJournalError
	if !errors.As(err, &corrupt) || corrupt.Line != 2 {
		t.Fatalf("replay error = %v, want second-Event constraint failure", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var names string
	if err := conn.QueryRowContext(ctx, "SELECT group_concat(name) FROM labels WHERE name IN ('first','rolled-back')").Scan(&names); err != nil {
		t.Fatalf("query after failed rebuild: %v", err)
	}
	if names != "first" {
		t.Fatalf("labels = %q, want only the committed first Event", names)
	}
	meta, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil || !found || meta.LastEventID != first.EventID {
		t.Fatalf("checkpoint after failure = %+v, found %v, error %v", meta, found, err)
	}
}

func TestRebuildFailurePreservesInstalledProjectionAndJournal(t *testing.T) {
	t.Parallel()
	workspace := newReplayWorkspace(t)
	path := filepath.Join(workspace, ".fid", "fid.sqlite")
	conn := openReplayDB(t, path)
	closeProjectDB(t, conn)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	writeReplayFixture(t, workspace, "conflicting-precondition.jsonl", "2026-06-05.jsonl")
	journalPath := filepath.Join(workspace, ".fid", "journal", "2026-06-05.jsonl")
	journalBefore, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = RebuildWorkspace(workspace, testProjectMigrationsDir(t))
	var conflict *PreconditionConflictError
	if !errors.As(err, &conflict) || conflict.Line != 2 || conflict.EventID != "EVT-20260605-000002" {
		t.Fatalf("rebuild error = %v, want located precondition conflict", err)
	}
	for name, want := range map[string][]byte{path: before, journalPath: journalBefore} {
		got, err := os.ReadFile(name)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("failed rebuild changed %s: %v", name, err)
		}
	}
	temporary, err := filepath.Glob(filepath.Join(workspace, ".fid", ".fid-rebuild-*"))
	if err != nil || len(temporary) != 0 {
		t.Fatalf("temporary projections after failure = %v, error %v", temporary, err)
	}
}
