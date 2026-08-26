package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV01CAllACsVerifiedGateAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Phase 7 AC gate parity",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)
	_, claimStderr := runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")
	if !strings.Contains(claimStderr, "claimed "+issueID+" by developer (expires ") {
		t.Fatalf("claim stderr = %q, want confirmation", claimStderr)
	}

	_, transitionStderr := runFIDCommand(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "in_progress",
		"--actor", "agent:developer",
	)
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": todo → in_progress") {
		t.Fatalf("todo -> in_progress stderr = %q, want confirmation", transitionStderr)
	}

	insertTask4CompletionRow(t, workspace, issueID)

	_, transitionStderr = runFIDCommand(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "dev_complete",
		"--actor", "agent:developer",
	)
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": in_progress → dev_complete") {
		t.Fatalf("in_progress -> dev_complete stderr = %q, want confirmation", transitionStderr)
	}

	_, transitionStderr = runFIDCommand(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "qa_review",
		"--actor", "agent:qa_lead",
	)
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": dev_complete → qa_review") {
		t.Fatalf("dev_complete -> qa_review stderr = %q, want confirmation", transitionStderr)
	}

	addStdout, _ := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "First QA approval criterion",
		"--ordinal", "1",
		"--by", "developer",
	)
	firstACID := parseAddedAcceptanceCriterionID(t, addStdout)

	addStdout, _ = runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Second QA approval criterion",
		"--ordinal", "2",
		"--by", "developer",
	)
	secondACID := parseAddedAcceptanceCriterionID(t, addStdout)

	journalPath := mustJournalPath(t, workspace)
	linesBeforeReject := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	statusBeforeReject := mustGetIssueJSON(t, workspace, issueID)

	stdout, stderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "qa_pass",
		"--actor", "agent:qa_lead",
	)
	if err == nil {
		t.Fatal("qa_review -> qa_pass error = nil, want gated rejection")
	}
	if stdout != "" {
		t.Fatalf("reject stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"all_acs_verified",
		"missing roles qa_lead",
		firstACID,
		secondACID,
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("reject stderr = %q, want substring %q", stderr, want)
		}
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforeReject {
		t.Fatalf("journal lines after reject = %d, want %d", got, linesBeforeReject)
	}
	if got := mustGetIssueJSON(t, workspace, issueID); got != statusBeforeReject {
		t.Fatalf("issue JSON after reject = %q, want %q", got, statusBeforeReject)
	}

	assertIssueACVerifyUnavailable(t, workspace, issueID, firstACID, "agent:qa_lead")
	seedHistoricalACVerification(t, workspace, firstACID, "qa_lead")
	assertIssueACVerifyUnavailable(t, workspace, issueID, secondACID, "agent:qa_lead")
	seedHistoricalACVerification(t, workspace, secondACID, "qa_lead")

	_, transitionStderr = runFIDCommand(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "qa_pass",
		"--actor", "agent:qa_lead",
	)
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": qa_review → qa_pass") {
		t.Fatalf("qa_review -> qa_pass stderr = %q, want confirmation", transitionStderr)
	}

	beforeAcceptanceRows := snapshotTask4ACRows(t, workspace)
	beforeAcceptanceJSON := mustGetAcceptanceSectionJSON(t, workspace, issueID)
	beforeACListJSON := mustGetTask4ACListJSON(t, workspace, issueID)
	assertTask4ACRowsEqual(t, snapshotTask4ACRows(t, workspace), beforeAcceptanceRows)
	if got := mustGetAcceptanceSectionJSON(t, workspace, issueID); got != beforeAcceptanceJSON {
		t.Fatalf("issue get acceptance after fail-closed verification = %q, want %q", got, beforeAcceptanceJSON)
	}
	if got := mustGetTask4ACListJSON(t, workspace, issueID); got != beforeACListJSON {
		t.Fatalf("ac list after fail-closed verification = %q, want %q", got, beforeACListJSON)
	}
}

type task4AcceptanceRow struct {
	ID      string
	IssueID string
	Ordinal int
	Text    string
}

func insertTask4CompletionRow(t *testing.T, workspace, issueID string) {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	if _, err := conn.Exec(`
		INSERT INTO task_completions (id, issue_id, attempt_number, developer_role, status, completed_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "TC-task4-001", issueID, 1, "developer", "dev_complete", "2026-06-14T12:00:00Z"); err != nil {
		t.Fatalf("insert task_completions: %v", err)
	}
}

func snapshotTask4ACRows(t *testing.T, workspace string) []task4AcceptanceRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT id, issue_id, ordinal, text
		FROM acceptance_criteria
		ORDER BY ordinal, id
	`)
	if err != nil {
		t.Fatalf("query acceptance_criteria: %v", err)
	}
	defer closeTask4Rows(t, rows)

	var snapshot []task4AcceptanceRow
	for rows.Next() {
		var row task4AcceptanceRow
		if err := rows.Scan(&row.ID, &row.IssueID, &row.Ordinal, &row.Text); err != nil {
			t.Fatalf("scan acceptance_criteria: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate acceptance_criteria: %v", err)
	}
	return snapshot
}

func mustGetAcceptanceSectionJSON(t *testing.T, workspace, issueID string) string {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "--json", "issue", "ac", "list", issueID)
	if stderr != "" {
		t.Fatalf("issue ac list stderr = %q, want empty", stderr)
	}
	return mustCanonicalAcceptanceCriteriaJSON(t, stdout)
}

func mustGetTask4ACListJSON(t *testing.T, workspace, issueID string) string {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "--json", "issue", "ac", "list", issueID)
	if stderr != "" {
		t.Fatalf("issue ac list stderr = %q, want empty", stderr)
	}
	return canonicalJSON(t, stdout)
}

func assertTask4ACRowsEqual(t *testing.T, got, want []task4AcceptanceRow) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("acceptance_criteria rows = %#v, want %#v", got, want)
	}
}

func canonicalJSON(t *testing.T, raw string) string {
	t.Helper()

	var value any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	normalized, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	return string(normalized)
}

func closeTask4Rows(t *testing.T, rows *sql.Rows) {
	t.Helper()

	if err := rows.Close(); err != nil {
		t.Fatalf("Rows.Close() error = %v", err)
	}
}
