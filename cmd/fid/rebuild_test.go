package main

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestRebuildCommandReplaysJournalAndReportsCounts(t *testing.T) {
	workspace := newRebuildWorkspace(t)
	writeCommandFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"rebuild", "--by", "developer"})

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	defer mustChdir(t, cwd)
	mustChdir(t, workspace)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if got := stderr.String(); !strings.Contains(got, "replayed 2 events from 1 journal files") {
		t.Fatalf("stderr = %q, want replay count line", got)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var title string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, "FID-a1b2c3").Scan(&title); err != nil {
		t.Fatalf("query rebuilt issue: %v", err)
	}
	if title != "Implement replay" {
		t.Fatalf("rebuilt issue title = %q, want %q", title, "Implement replay")
	}
}

func TestRebuildCommandReportsCorruptJournal(t *testing.T) {
	workspace := newRebuildWorkspace(t)
	writeCommandFixture(t, workspace, "corrupt-tail.jsonl", "2026-06-05.jsonl")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"rebuild", "--by", "developer"})

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	defer mustChdir(t, cwd)
	mustChdir(t, workspace)

	err = cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want corrupt journal error")
	}
	if got := stderr.String(); !strings.Contains(got, "malformed journal entry at .fid/journal/2026-06-05.jsonl:2") {
		t.Fatalf("stderr = %q, want corrupt journal line", got)
	}
}

func newRebuildWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if err := os.MkdirAll(filepath.Join(workspace, ".fid", "journal"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	return workspace
}

func commandFixturePath(t *testing.T, name string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "journal", name))
}

func writeCommandFixture(t *testing.T, workspaceRoot, fixtureName, journalFileName string) {
	t.Helper()

	payload, err := os.ReadFile(commandFixturePath(t, fixtureName))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", fixtureName, err)
	}

	journalPath := filepath.Join(workspaceRoot, ".fid", "journal", journalFileName)
	if err := os.WriteFile(journalPath, payload, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", journalPath, err)
	}
}

func mustChdir(t *testing.T, path string) {
	t.Helper()

	if err := os.Chdir(path); err != nil {
		t.Fatalf("Chdir(%q) error = %v", path, err)
	}
}

func closeCommandDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
