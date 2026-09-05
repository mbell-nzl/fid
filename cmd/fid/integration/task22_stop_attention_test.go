//go:build integration

package integration

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestStopAttentionExternalBoundaryKeepsRoutineStatusAndWarningsEffectFree(t *testing.T) {
	workspace := t.TempDir()
	binary := integrationFIDBinary(t)
	closureRun(t, workspace, binary, "--json", "init", "--prefix", "P22")

	journalPath := closureJournalPath(t, workspace)
	journalBefore := readStopAttentionJournal(t, journalPath)
	assertStopAttentionCount(t, workspace, 0)

	routine, _ := closureRun(t, workspace, binary, "--json", "maintenance", "check")
	if !strings.Contains(routine, `"status":"pass"`) {
		t.Fatalf("routine maintenance status = %q, want pass", routine)
	}
	assertStopAttentionEffectFree(t, workspace, journalPath, journalBefore, "routine maintenance status")

	setStartedIssueOldEnoughForWarning(t, workspace)
	warning, _ := closureRun(t, workspace, binary, "--json", "maintenance", "check")
	if !strings.Contains(warning, `"status":"warn"`) || !strings.Contains(warning, `"check":"started_without_claim"`) {
		t.Fatalf("maintenance warning = %q, want started_without_claim warning", warning)
	}
	assertStopAttentionEffectFree(t, workspace, journalPath, journalBefore, "maintenance warning")

	capabilities, _ := closureRun(t, workspace, binary, "--json", "capabilities")
	if !strings.Contains(capabilities, `"capability_id":"attention.v1"`) || !strings.Contains(capabilities, `"available":false`) {
		t.Fatalf("capabilities = %q, want unavailable attention.v1", capabilities)
	}
	help, _ := closureRun(t, workspace, binary, "--help")
	if strings.Contains(help, "\n  attention") {
		t.Fatalf("fid --help = %q, want no deferred attention leaf", help)
	}
	stdout, stderr, err := closureRunCommand(workspace, binary, "--json", "attention", "list")
	if err == nil || stdout != "" || stderr != "" {
		t.Fatalf("fid --json attention list = (%v, %q, %q), want unavailable no-leaf rejection", err, stdout, stderr)
	}
	assertStopAttentionEffectFree(t, workspace, journalPath, journalBefore, "unavailable Attention leaf")
}

func setStartedIssueOldEnoughForWarning(t *testing.T, workspace string) {
	t.Helper()
	database, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection for warning fixture: %v", err)
	}
	defer func() { _ = database.Close() }()
	if _, err := database.Exec(`INSERT INTO issues (id, kind, title, workflow_id, status, created_at, updated_at) VALUES ('P22-warning', 'task', 'warning fixture', 'task-default', 'in_progress', '2000-01-01T00:00:00Z', '2000-01-01T00:00:00Z')`); err != nil {
		t.Fatalf("seed started warning issue: %v", err)
	}
}

func assertStopAttentionEffectFree(t *testing.T, workspace, journalPath string, before []byte, operation string) {
	t.Helper()
	if after := readStopAttentionJournal(t, journalPath); !bytes.Equal(after, before) {
		t.Fatalf("%s changed journal bytes", operation)
	}
	assertStopAttentionCount(t, workspace, 0)
}

func readStopAttentionJournal(t *testing.T, path string) []byte {
	t.Helper()
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}
	return payload
}

func assertStopAttentionCount(t *testing.T, workspace string, want int) {
	t.Helper()
	database, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection for Attention count: %v", err)
	}
	defer func() { _ = database.Close() }()
	var got int
	if err := database.QueryRow(`SELECT COUNT(*) FROM attention_items WHERE status = 'open'`).Scan(&got); err != nil && err != sql.ErrNoRows {
		t.Fatalf("count open Attention: %v", err)
	}
	if got != want {
		t.Fatalf("open Attention count = %d, want %d", got, want)
	}
}
