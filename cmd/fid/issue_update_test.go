package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIssueUpdateChangesFieldsAndPrintsConfirmation(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Implement issue update", "--priority", "2", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	_, updateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "update", issueID,
		"--title", "Ship issue update",
		"--priority", "1",
		"--by", "developer",
	)
	gotLines := strings.Split(strings.TrimSpace(updateStderr), "\n")
	wantLines := []string{
		"updated " + issueID + ":",
		`title "Implement issue update" -> "Ship issue update"`,
		"priority 2 -> 1",
	}
	if len(gotLines) != len(wantLines) {
		t.Fatalf("stderr lines = %#v, want %#v", gotLines, wantLines)
	}
	for index, wantLine := range wantLines {
		if gotLines[index] != wantLine {
			t.Fatalf("stderr lines = %#v, want %#v", gotLines, wantLines)
		}
	}

	stdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

	var got struct {
		Title    string `json:"title"`
		Priority *int   `json:"priority"`
	}
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Title != "Ship issue update" {
		t.Fatalf("title = %q, want Ship issue update", got.Title)
	}
	if got.Priority == nil || *got.Priority != 1 {
		t.Fatalf("priority = %#v, want 1", got.Priority)
	}
}

func TestIssueUpdateSingleFieldPrintsExactlyOneChangedFieldLine(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Single field update", "--priority", "2", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	_, updateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "update", issueID,
		"--priority", "1",
		"--by", "developer",
	)

	gotLines := strings.Split(strings.TrimSpace(updateStderr), "\n")
	wantLines := []string{
		"updated " + issueID + ":",
		"priority 2 -> 1",
	}
	if len(gotLines) != len(wantLines) {
		t.Fatalf("stderr lines = %#v, want %#v", gotLines, wantLines)
	}
	for index, wantLine := range wantLines {
		if gotLines[index] != wantLine {
			t.Fatalf("stderr lines = %#v, want %#v", gotLines, wantLines)
		}
	}
}

func TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Reject invalid update", "--priority", "2", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	beforeStdout, beforeStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if beforeStderr != "" {
		t.Fatalf("issue get before stderr = %q, want empty", beforeStderr)
	}

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "update", issueID, "--priority", "8", "--by", "developer")
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid update rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "error: update issue: priority must be between 0 and 5") {
		t.Fatalf("stderr = %q, want invalid priority error", stderr)
	}

	afterStdout, afterStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if afterStderr != "" {
		t.Fatalf("issue get after stderr = %q, want empty", afterStderr)
	}
	if afterStdout != beforeStdout {
		t.Fatalf("issue JSON after rejected update = %q, want %q", afterStdout, beforeStdout)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore {
		t.Fatalf("journal lines after rejected update = %d, want %d", linesAfter, linesBefore)
	}
}

func TestIssueUpdateRequiresKnownByRole(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Role validation update", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	t.Run("missing by", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "update", issueID, "--title", "No role")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --by rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: --by is required") {
			t.Fatalf("stderr = %q, want missing --by error", stderr)
		}
	})

	t.Run("unknown role", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "update", issueID, "--title", "Unknown role", "--by", "intern")
		if err == nil {
			t.Fatal("Execute() error = nil, want unknown role rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
			t.Fatalf("stderr = %q, want unknown role error", stderr)
		}
	})
}
