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
	clonePristineWorkspaceInto(t, workspace)

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
	clonePristineWorkspaceInto(t, workspace)

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
	clonePristineWorkspaceInto(t, workspace)

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

func TestDevFinishCommandClosesRunAndReleasesClaim(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	issueID := seedCommandDevFinishReadyIssue(t, workspace)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"dev", "finish", issueID,
		"--completion-md", "project-docs/docs/phases/phase-15/tasks/task-3-done.md",
		"--tokens", "14823",
		"--commits", "abc123, def456",
		"--files", "internal/dev/finish.go",
		"--files", "cmd/fid/dev.go",
		"--human-intervention",
		"--intervention-note", "paired on final review",
		"--by", "developer",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"dev finish " + issueID,
		"agent run: RUN-",
		"completion: TC-",
		"transitioned: in_progress → dev_complete",
		"claim released",
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
	if got.Status != "dev_complete" {
		t.Fatalf("issue status = %q, want dev_complete", got.Status)
	}
}

func TestDevBlockCommandRecordsFailedGateAndLeavesRunUntouched(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	issueID := seedCommandStartedIssue(t, workspace, "CLI dev block")
	runFIDCommand(
		t,
		workspace,
		"dev", "gate", issueID,
		"--cmd", "golangci-lint run ./...",
		"--exit-code", "1",
		"--output-excerpt", "internal/dev/block.go: line too long",
		"--by", "developer",
	)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"dev", "block", issueID,
		"--reason-md", "project-docs/docs/phases/phase-15/tasks/task-3-done.md",
		"--by", "developer",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"dev block " + issueID,
		"transitioned: in_progress → failed_gate",
		"completion: TC-",
		"claim released",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %q, want substring %q", stderr, want)
		}
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var status string
	if err := conn.QueryRow(`SELECT status FROM issues WHERE id = ?`, issueID).Scan(&status); err != nil {
		t.Fatalf("query issues row: %v", err)
	}
	if status != "failed_gate" {
		t.Fatalf("issues.status = %q, want failed_gate", status)
	}

	var openRunCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM agent_runs WHERE issue_id = ? AND ended_at IS NULL`, issueID).Scan(&openRunCount); err != nil {
		t.Fatalf("count open agent_runs: %v", err)
	}
	if openRunCount != 1 {
		t.Fatalf("open agent_runs = %d, want 1", openRunCount)
	}
}

func TestDevFinishAndBlockDeniedRolesRejectBeforeAppend(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	finishIssueID := seedCommandDevFinishReadyIssue(t, workspace)
	blockIssueID := seedCommandStartedIssue(t, workspace, "CLI denied dev block")
	journalPath := mustJournalPath(t, workspace)

	linesBeforeFinishDenied := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	deniedStdout, deniedStderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"dev", "finish", finishIssueID,
		"--completion-md", "project-docs/docs/phases/phase-15/tasks/task-3-done.md",
		"--by", "qa_reviewer",
	)
	if err == nil {
		t.Fatal("Execute() finish error = nil, want permission denial")
	}
	if deniedStdout != "" {
		t.Fatalf("finish stdout = %q, want empty", deniedStdout)
	}
	if !strings.Contains(deniedStderr, `error: role "qa_reviewer" not permitted for action "dev.finish"`) {
		t.Fatalf("finish stderr = %q, want permission denial", deniedStderr)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforeFinishDenied {
		t.Fatalf("journal lines after denied dev finish = %d, want unchanged %d", got, linesBeforeFinishDenied)
	}

	linesBeforeBlockDenied := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	deniedStdout, deniedStderr, err = runFIDCommandExpectError(
		t,
		workspace,
		"dev", "block", blockIssueID,
		"--reason-md", "project-docs/docs/phases/phase-15/tasks/task-3-done.md",
		"--by", "qa_reviewer",
	)
	if err == nil {
		t.Fatal("Execute() block error = nil, want permission denial")
	}
	if deniedStdout != "" {
		t.Fatalf("block stdout = %q, want empty", deniedStdout)
	}
	if !strings.Contains(deniedStderr, `error: role "qa_reviewer" not permitted for action "dev.block"`) {
		t.Fatalf("block stderr = %q, want permission denial", deniedStderr)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforeBlockDenied {
		t.Fatalf("journal lines after denied dev block = %d, want unchanged %d", got, linesBeforeBlockDenied)
	}
}

func seedCommandDevFinishReadyIssue(t *testing.T, workspace string) string {
	t.Helper()

	issueID := seedCommandStartedIssue(t, workspace, "CLI dev finish")
	_, addStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "developer verification required",
		"--ordinal", "1",
		"--by", "developer",
	)
	acID := parseAddedAcceptanceCriterionID(t, addStderr)
	runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, acID,
		"--status", "pass",
		"--by", "developer",
	)
	runFIDCommand(
		t,
		workspace,
		"dev", "gate", issueID,
		"--cmd", "go test ./...",
		"--exit-code", "0",
		"--output-excerpt", "ok",
		"--by", "developer",
	)
	return issueID
}

func seedCommandStartedIssue(t *testing.T, workspace, title string) string {
	t.Helper()

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", title,
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
	return issueID
}
