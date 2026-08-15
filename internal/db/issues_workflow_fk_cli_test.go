package db_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestMigration0008PreservesIssueGetAndListRenderedBytes(t *testing.T) {
	workspace := t.TempDir()
	if err := os.MkdirAll(filepath.Join(workspace, ".fid", "journal"), 0o755); err != nil {
		t.Fatalf("MkdirAll workspace journal: %v", err)
	}

	v7Migrations := migration0007Directory(t)
	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, v7Migrations); err != nil {
		t.Fatalf("MigrateUp(version 7) error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeMigrationReadDB(t, conn)
	if _, err := conn.Exec(`
		INSERT INTO issues (id, kind, title, workflow_id, status, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, "FID-cli-parity", "task", "CLI parity fixture", "task-default", "in_progress", 1, "2026-07-19T03:00:00Z", "2026-07-19T03:01:00Z"); err != nil {
		t.Fatalf("insert issue fixture: %v", err)
	}
	if err := fiddb.MigrateUp(dbPath, sourceMigrationsDirectory(t)); err != nil {
		t.Fatalf("MigrateUp(current schema before canonical read): %v", err)
	}
	writeMigrationReadCheckpoint(t, conn, workspace)

	beforeGet := renderedIssueGet(t, workspace, v7Migrations, "FID-cli-parity")
	beforeList := renderedIssueList(t, workspace, v7Migrations)

	if err := fiddb.MigrateUp(dbPath, sourceMigrationsDirectory(t)); err != nil {
		t.Fatalf("MigrateUp(version 8) error = %v", err)
	}

	afterGet := renderedIssueGet(t, workspace, sourceMigrationsDirectory(t), "FID-cli-parity")
	afterList := renderedIssueList(t, workspace, sourceMigrationsDirectory(t))
	if afterGet != beforeGet {
		t.Fatalf("fid issue get JSON after migration = %q, want byte-identical %q", afterGet, beforeGet)
	}
	if afterList != beforeList {
		t.Fatalf("fid issue list JSON after migration = %q, want byte-identical %q", afterList, beforeList)
	}
}

func renderedIssueGet(t *testing.T, workspace, migrationsDir, issueID string) string {
	t.Helper()

	issue, err := issues.Get(issues.GetOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ID: issueID})
	if err != nil {
		t.Fatalf("fid issue get read: %v", err)
	}
	return renderMigrationRead(t, issue)
}

func renderedIssueList(t *testing.T, workspace, migrationsDir string) string {
	t.Helper()

	items, err := issues.List(issues.ListOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir})
	if err != nil {
		t.Fatalf("fid issue list read: %v", err)
	}
	return renderMigrationRead(t, items)
}

func renderMigrationRead(t *testing.T, value any) string {
	t.Helper()

	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		t.Fatalf("render JSON: %v", err)
	}
	return output.String()
}

func writeMigrationReadCheckpoint(t *testing.T, conn *sql.DB, workspace string) {
	t.Helper()

	event := journal.Event{
		EventID:     "EVT-20260719-000001",
		EventSchema: journal.EventSchemaV1,
		EventType:   "test.checkpoint",
		OccurredAt:  time.Date(2026, time.July, 19, 3, 2, 0, 0, time.UTC),
		Actor:       journal.Actor{Role: "developer"},
		Command:     journal.Command{Argv: []string{"fid", "test"}, CWD: ".", FIDVersion: "0.2.0"},
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal checkpoint event: %v", err)
	}
	hash, err := journal.HashEvent(event)
	if err != nil {
		t.Fatalf("hash checkpoint event: %v", err)
	}
	journalPath := filepath.Join(workspace, ".fid", "journal", "2026-07-19.jsonl")
	contents := append(encoded, '\n')
	if err := os.WriteFile(journalPath, contents, 0o644); err != nil {
		t.Fatalf("write checkpoint journal: %v", err)
	}

	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("begin checkpoint meta transaction: %v", err)
	}
	if err := fiddb.WriteProjectionMeta(tx, fiddb.ProjectionMeta{
		JournalFile:             "2026-07-19.jsonl",
		JournalOffset:           int64(len(contents)),
		LastEventID:             event.EventID,
		LastEventHash:           hash,
		EventSchema:             journal.EventSchemaV1,
		ProjectionSchemaVersion: fiddb.ProjectionSchemaVersion,
		ReplayStatus:            fiddb.ReplayStatusClean,
		UpdatedAt:               event.OccurredAt,
	}); err != nil {
		_ = tx.Rollback()
		t.Fatalf("write checkpoint meta: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit checkpoint meta: %v", err)
	}
}

func migration0007Directory(t *testing.T) string {
	t.Helper()

	directory := t.TempDir()
	for version := 1; version <= 7; version++ {
		pattern := filepath.Join(sourceMigrationsDirectory(t), fmt.Sprintf("%04d_*.sql", version))
		matches, err := filepath.Glob(pattern)
		if err != nil || len(matches) != 2 {
			t.Fatalf("glob migrations for version %d = %#v, %v; want up/down pair", version, matches, err)
		}
		for _, source := range matches {
			contents, err := os.ReadFile(source)
			if err != nil {
				t.Fatalf("read migration %q: %v", source, err)
			}
			destination := filepath.Join(directory, filepath.Base(source))
			if err := os.WriteFile(destination, contents, 0o644); err != nil {
				t.Fatalf("write migration %q: %v", destination, err)
			}
		}
	}
	return directory
}

func sourceMigrationsDirectory(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func closeMigrationReadDB(t *testing.T, conn *sql.DB) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
