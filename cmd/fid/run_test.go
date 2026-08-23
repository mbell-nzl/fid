package main

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestRunRecordCreatesAgentRunAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI run record",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"run", "record",
		"--issue", issueID,
		"--role", "developer",
		"--model", "claude-sonnet-4-6",
		"--tokens-in", "1200",
		"--tokens-out", "3400",
		"--by", "developer",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}

	runID := parseRecordedRunID(t, stderr)
	if !strings.Contains(stderr, "recorded agent run "+runID+" on "+issueID) {
		t.Fatalf("stderr = %q, want run confirmation", stderr)
	}
	if !strings.Contains(stderr, "role: developer, model: claude-sonnet-4-6") {
		t.Fatalf("stderr = %q, want role/model summary", stderr)
	}
	if !strings.Contains(stderr, "tokens: 1200 in, 3400 out") {
		t.Fatalf("stderr = %q, want token summary", stderr)
	}
	if !strings.Contains(stderr, "cost: $") {
		t.Fatalf("stderr = %q, want cost summary", stderr)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM agent_runs WHERE id = ?`, runID).Scan(&count); err != nil {
		t.Fatalf("count agent_runs row: %v", err)
	}
	if count != 1 {
		t.Fatalf("agent_runs rows for %s = %d, want 1", runID, count)
	}
}

func TestRunRecordUnknownModelWarnsAndStillRecords(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI unknown model run",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"run", "record",
		"--issue", issueID,
		"--role", "developer",
		"--model", "mystery-model",
		"--tokens-in", "500",
		"--tokens-out", "700",
		"--by", "developer",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `warning: model "mystery-model"`) {
		t.Fatalf("stderr = %q, want unknown-model warning", stderr)
	}
	if !strings.Contains(stderr, "recorded agent run RUN-") {
		t.Fatalf("stderr = %q, want run confirmation", stderr)
	}
	if !strings.Contains(stderr, "cost: unknown") {
		t.Fatalf("stderr = %q, want unknown cost summary", stderr)
	}

	runID := parseRecordedRunID(t, stderr)

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var priceIn sql.NullFloat64
	var priceOut sql.NullFloat64
	var cost sql.NullFloat64
	if err := conn.QueryRow(`
		SELECT price_in_per_mtok, price_out_per_mtok, cost_usd
		FROM agent_runs
		WHERE id = ?
	`, runID).Scan(&priceIn, &priceOut, &cost); err != nil {
		t.Fatalf("query agent_runs row: %v", err)
	}
	if priceIn.Valid || priceOut.Valid || cost.Valid {
		t.Fatalf("unknown-model pricing row = (%#v, %#v, %#v), want NULLs", priceIn, priceOut, cost)
	}
}

func TestRunRecordRejectsUnknownActorRoleBeforeAppend(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI unknown actor role",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"run", "record",
		"--issue", issueID,
		"--role", "developer",
		"--model", "claude-sonnet-4-6",
		"--by", "intern",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want unknown role rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
		t.Fatalf("stderr = %q, want unknown role error", stderr)
	}
	if !strings.Contains(stderr, "known roles: developer, human, qa_lead, qa_reviewer, tl, tl_lite, workflow_runner") {
		t.Fatalf("stderr = %q, want known roles list", stderr)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore {
		t.Fatalf("journal lines after unknown actor role = %d, want unchanged %d", linesAfter, linesBefore)
	}
}

func TestRunRecordPolicyShowMatchesAllowedRoles(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show")
	if stderr != "" {
		t.Fatalf("policy show stderr = %q, want empty", stderr)
	}

	rows := decodeCommandJSONStructure(t, stdout).([]any)
	for _, role := range []string{"developer", "qa_reviewer", "qa_lead", "tl_lite", "tl", "workflow_runner", "human"} {
		if !commandPolicyCellAllowed(t, rows, "run.record", role) {
			t.Fatalf("policy show %s/run.record = false, want true", role)
		}
	}

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI run policy parity",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	for _, role := range []string{"developer", "qa_reviewer", "qa_lead", "tl_lite", "tl", "workflow_runner", "human"} {
		t.Run(role, func(t *testing.T) {
			roleStdout, roleStderr := runFIDCommand(
				t,
				workspace,
				"run", "record",
				"--issue", issueID,
				"--role", "developer",
				"--model", "claude-sonnet-4-6",
				"--tokens-in", "10",
				"--tokens-out", "20",
				"--by", role,
			)
			if roleStdout != "" {
				t.Fatalf("stdout = %q, want empty", roleStdout)
			}
			if !strings.Contains(roleStderr, "recorded agent run RUN-") {
				t.Fatalf("stderr = %q, want run confirmation", roleStderr)
			}
		})
	}
}

func parseRecordedRunID(t *testing.T, stderr string) string {
	t.Helper()

	fields := strings.Fields(stderr)
	for index := range fields {
		if fields[index] == "run" && index > 1 && fields[index-1] == "agent" && index+1 < len(fields) {
			return fields[index+1]
		}
	}
	t.Fatalf("stderr = %q, want recorded agent run line with id", stderr)
	return ""
}
