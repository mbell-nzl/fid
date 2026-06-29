package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestDevStartCommandPrintsCompositeConfirmation(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI dev start",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"dev", "start", issueID,
		"--model", "claude-sonnet-4-6",
		"--by", "developer",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"dev start " + issueID,
		"claimed by developer",
		"transitioned: todo → in_progress",
		"agent run: RUN-",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %q, want substring %q", stderr, want)
		}
	}

	issueStdout, issueStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if issueStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", issueStderr)
	}
	var got struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal([]byte(issueStdout), &got); err != nil {
		t.Fatalf("Unmarshal(issue get) error = %v", err)
	}
	if got.Status != "in_progress" {
		t.Fatalf("issue status = %q, want in_progress", got.Status)
	}
}

func TestDevGateCommandRecordsReportedResultAndDoesNotExecuteCommand(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI dev gate",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)
	runFIDCommand(
		t,
		workspace,
		"dev", "start", issueID,
		"--model", "claude-sonnet-4-6",
		"--by", "developer",
	)

	forbiddenPath := filepath.Join(workspace, "should-not-exist.txt")
	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"dev", "gate", issueID,
		"--cmd", "touch should-not-exist.txt",
		"--exit-code", "1",
		"--output-excerpt", "touch failed",
		"--by", "developer",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"gate recorded on " + issueID + " (attempt 1)",
		"cmd: touch should-not-exist.txt",
		"exit_code: 1",
		"gates_passed: 0, gates_failed: 1",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %q, want substring %q", stderr, want)
		}
	}
	if _, err := os.Stat(forbiddenPath); !os.IsNotExist(err) {
		t.Fatalf("Stat(%q) error = %v, want command not executed", forbiddenPath, err)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var count int
	if err := conn.QueryRow(`
		SELECT COUNT(*)
		FROM gate_runs
		WHERE issue_id = ? AND command = ? AND exit_code = 1
	`, issueID, "touch should-not-exist.txt").Scan(&count); err != nil {
		t.Fatalf("count gate_runs: %v", err)
	}
	if count != 1 {
		t.Fatalf("gate_runs count = %d, want 1", count)
	}
}

func TestDevCommandsPolicyShowMatchesRowsAndDeniedRoles(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show", "--format", "json")
	if stderr != "" {
		t.Fatalf("policy show stderr = %q, want empty", stderr)
	}

	rows := decodeCommandJSONStructure(t, stdout).([]any)
	for _, action := range []string{"dev.start", "dev.gate", "dev.finish", "dev.block"} {
		if !commandPolicyCellAllowed(t, rows, action, "developer") {
			t.Fatalf("policy show developer/%s = false, want true", action)
		}
		if !commandPolicyCellAllowed(t, rows, action, "human") {
			t.Fatalf("policy show human/%s = false, want true", action)
		}
		if commandPolicyCellAllowed(t, rows, action, "qa_reviewer") {
			t.Fatalf("policy show qa_reviewer/%s = true, want false", action)
		}
	}

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI dev policy",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBeforeDenied := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	deniedStdout, deniedStderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"dev", "start", issueID,
		"--model", "claude-sonnet-4-6",
		"--by", "qa_reviewer",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want permission denial")
	}
	if deniedStdout != "" {
		t.Fatalf("stdout = %q, want empty", deniedStdout)
	}
	if !strings.Contains(deniedStderr, `error: role "qa_reviewer" not permitted for action "dev.start"`) {
		t.Fatalf("stderr = %q, want permission denial", deniedStderr)
	}

	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforeDenied {
		t.Fatalf("journal lines after denied dev start = %d, want unchanged %d", got, linesBeforeDenied)
	}

	runFIDCommand(
		t,
		workspace,
		"dev", "start", issueID,
		"--model", "claude-sonnet-4-6",
		"--by", "developer",
	)
	linesBeforeGateDenied := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	deniedStdout, deniedStderr, err = runFIDCommandExpectError(
		t,
		workspace,
		"dev", "gate", issueID,
		"--cmd", "go test ./...",
		"--exit-code", "0",
		"--by", "qa_reviewer",
	)
	if err == nil {
		t.Fatal("Execute() gate error = nil, want permission denial")
	}
	if deniedStdout != "" {
		t.Fatalf("gate stdout = %q, want empty", deniedStdout)
	}
	if !strings.Contains(deniedStderr, `error: role "qa_reviewer" not permitted for action "dev.gate"`) {
		t.Fatalf("gate stderr = %q, want permission denial", deniedStderr)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforeGateDenied {
		t.Fatalf("journal lines after denied dev gate = %d, want unchanged %d", got, linesBeforeGateDenied)
	}
}
