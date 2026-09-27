package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

// The V1 sync surface is gone. A markdown file must not create a pointer or Event.
func TestRequirementSyncCreatesPointersInOneJournalEvent(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	writeRequirementFixture(t, workspace, "project-docs/specs/sample_requirements.md", "# Requirements\n\n| PRD-GOV-006 | Pointer prose |\n")
	before := len(splitJournalLines(mustReadCommandFile(t, mustJournalPath(t, workspace))))
	_, _, err := runFIDCommandExpectError(t, workspace, "requirement", "sync")
	if err == nil {
		t.Fatal("removed requirement sync leaf was accepted")
	}
	if after := len(splitJournalLines(mustReadCommandFile(t, mustJournalPath(t, workspace)))); after != before {
		t.Fatalf("removed sync changed journal lines %d -> %d", before, after)
	}
}

func TestRequirementSyncRejectsMalformedInputBeforeAppend(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"--source-doc", "project-docs/specs/bad.md"},
		{"--source-doc", "project-docs/specs/bad.md", "--type", "OPS"},
		{"--source-doc", "project-docs/specs/bad.md", "--id", "TRD-ARCH-010"},
	} {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)
		before := len(splitJournalLines(mustReadCommandFile(t, mustJournalPath(t, workspace))))
		_, _, err := runFIDCommandExpectError(t, workspace, append([]string{"requirement", "sync"}, args...)...)
		if err == nil {
			t.Fatalf("removed sync accepted %v", args)
		}
		if after := len(splitJournalLines(mustReadCommandFile(t, mustJournalPath(t, workspace)))); after != before {
			t.Fatalf("removed sync changed journal lines %d -> %d", before, after)
		}
	}
}

func TestRequirementSyncAuthorizationAndPolicyShowParity(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	before := len(splitJournalLines(mustReadCommandFile(t, mustJournalPath(t, workspace))))
	for _, actor := range []string{"agent:developer", "agent:intern", "agent:human"} {
		_, _, err := runFIDCommandExpectError(t, workspace, "requirement", "sync", "--actor", actor)
		if err == nil {
			t.Fatalf("removed sync accepted actor %q", actor)
		}
	}
	if after := len(splitJournalLines(mustReadCommandFile(t, mustJournalPath(t, workspace)))); after != before {
		t.Fatalf("removed sync changed journal lines %d -> %d", before, after)
	}
}

func mustReadLastJournalEvent(t *testing.T, journalPath string) journal.Event {
	t.Helper()
	lines := splitJournalLines(mustReadCommandFile(t, journalPath))
	if len(lines) == 0 {
		t.Fatal("journal lines = 0, want at least one line")
	}
	var event journal.Event
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &event); err != nil {
		t.Fatalf("Unmarshal(last journal event) error = %v", err)
	}
	return event
}

func writeRequirementFixture(t *testing.T, workspace, relativePath, body string) {
	t.Helper()
	path := filepath.Join(workspace, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}
