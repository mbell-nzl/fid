package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV01EDevCompositeSmokeAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	journalPath := mustJournalPath(t, workspace)
	expectedJournalLines := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	finishReportPath := writeV01EMarkdown(t, workspace, "reports/dev-finish.md", "# finish\n")
	blockReasonPath := writeV01EMarkdown(t, workspace, "reports/dev-block.md", "# block\n")

	_, epicStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "epic",
		"--title", "Developer composite smoke epic",
		"--scheduling", "prioritised",
		"--by", "tl",
	)
	assertNoNestedLockFailure(t, epicStderr)
	epicID := parseCreatedIssueID(t, epicStderr)
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "epic create")

	_, finishCreateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Composite finish path",
		"--epic", epicID,
		"--priority", "1",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, finishCreateStderr)
	finishIssueID := parseCreatedIssueID(t, finishCreateStderr)
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "finish task create")
	assertV01EIssueState(t, workspace, finishIssueID, "todo", "unstarted")

	_, blockCreateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Composite block path",
		"--epic", epicID,
		"--priority", "2",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, blockCreateStderr)
	blockIssueID := parseCreatedIssueID(t, blockCreateStderr)
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "block task create")
	assertV01EIssueState(t, workspace, blockIssueID, "todo", "unstarted")

	_, addACStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", finishIssueID,
		"--text", "Developer gate passes before finish",
		"--ordinal", "1",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, addACStderr)
	finishACID := parseAddedAcceptanceCriterionID(t, addACStderr)
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "finish task ac add")

	deniedStdout, deniedStderr, deniedErr := runFIDCommandExpectError(
		t,
		workspace,
		"dev", "start", finishIssueID,
		"--model", "claude-sonnet-4-6",
		"--by", "qa_reviewer",
	)
	if deniedErr == nil {
		t.Fatal("qa_reviewer dev start error = nil, want permission denial")
	}
	if deniedStdout != "" {
		t.Fatalf("qa_reviewer dev start stdout = %q, want empty", deniedStdout)
	}
	assertNoNestedLockFailure(t, deniedStderr)
	if !strings.Contains(deniedStderr, `error: role "qa_reviewer" not permitted for action "dev.start"`) {
		t.Fatalf("qa_reviewer dev start stderr = %q, want permission denial", deniedStderr)
	}
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "qa_reviewer dev start denial")

	_, startFinishStderr := runFIDCommand(
		t,
		workspace,
		"dev", "start", finishIssueID,
		"--model", "claude-sonnet-4-6",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, startFinishStderr)
	for _, want := range []string{
		"dev start " + finishIssueID,
		"claimed by developer",
		"transitioned: todo → in_progress",
		"agent run: RUN-",
	} {
		if !strings.Contains(startFinishStderr, want) {
			t.Fatalf("dev start stderr = %q, want substring %q", startFinishStderr, want)
		}
	}
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "dev start finish path")
	assertV01EIssueState(t, workspace, finishIssueID, "in_progress", "started")
	finishClaims := mustReadIssueClaimsSection(t, workspace, finishIssueID)
	if len(finishClaims) != 1 || finishClaims[0].ClaimedBy != "developer" {
		t.Fatalf("finish path claims after dev start = %#v, want one developer claim", finishClaims)
	}
	finishRuns := mustReadV01ERunsForIssue(t, workspace, finishIssueID)
	if len(finishRuns) != 1 {
		t.Fatalf("finish path runs after dev start = %#v, want 1 run", finishRuns)
	}
	finishRunID := finishRuns[0].ID
	if finishRuns[0].EndedAt != "" {
		t.Fatalf("finish path run ended_at = %q, want open run", finishRuns[0].EndedAt)
	}
	if finishRuns[0].Role != "developer" || finishRuns[0].Model != "claude-sonnet-4-6" {
		t.Fatalf("finish path run row = %#v, want developer/claude-sonnet-4-6", finishRuns[0])
	}

	_, verifyACStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", finishIssueID, finishACID,
		"--status", "pass",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, verifyACStderr)
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "finish task ac verify")

	_, gateStderr := runFIDCommand(
		t,
		workspace,
		"dev", "gate", finishIssueID,
		"--cmd", "go test ./...",
		"--exit-code", "0",
		"--output-excerpt", "ok",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, gateStderr)
	for _, want := range []string{
		"gate recorded on " + finishIssueID + " (attempt 1)",
		"cmd: go test ./...",
		"exit_code: 0",
		"gates_passed: 1, gates_failed: 0",
	} {
		if !strings.Contains(gateStderr, want) {
			t.Fatalf("dev gate stderr = %q, want substring %q", gateStderr, want)
		}
	}
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "dev gate finish path")
	finishGateRuns := mustReadV01EGateRunsForIssue(t, workspace, finishIssueID)
	if len(finishGateRuns) != 1 {
		t.Fatalf("finish path gate_runs = %#v, want 1 row", finishGateRuns)
	}
	if finishGateRuns[0].AttemptNumber != 1 || finishGateRuns[0].AgentRunID != finishRunID {
		t.Fatalf("finish path gate row = %#v, want attempt 1 on run %q", finishGateRuns[0], finishRunID)
	}
	if finishGateRuns[0].Command != "go test ./..." || finishGateRuns[0].ExitCode != 0 || finishGateRuns[0].OutputExcerpt != "ok" {
		t.Fatalf("finish path gate row = %#v, want recorded go test success", finishGateRuns[0])
	}

	_, finishStderr := runFIDCommand(
		t,
		workspace,
		"dev", "finish", finishIssueID,
		"--completion-md", finishReportPath,
		"--tokens", "14823",
		"--commits", "abc123",
		"--files", "cmd/fid/e2e_v01e_dev_test.go",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, finishStderr)
	for _, want := range []string{
		"dev finish " + finishIssueID,
		"agent run: " + finishRunID + " closed (14823 tokens)",
		"completion: TC-",
		"transitioned: in_progress → dev_complete",
		"claim released",
	} {
		if !strings.Contains(finishStderr, want) {
			t.Fatalf("dev finish stderr = %q, want substring %q", finishStderr, want)
		}
	}
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "dev finish")
	assertV01EIssueState(t, workspace, finishIssueID, "dev_complete", "started")
	if claims := mustReadIssueClaimsSection(t, workspace, finishIssueID); len(claims) != 0 {
		t.Fatalf("finish path claims after dev finish = %#v, want empty", claims)
	}
	finishRuns = mustReadV01ERunsForIssue(t, workspace, finishIssueID)
	if len(finishRuns) != 1 || finishRuns[0].EndedAt == "" {
		t.Fatalf("finish path runs after dev finish = %#v, want one closed run", finishRuns)
	}
	finishCompletions := mustReadV01ECompletionsForIssue(t, workspace, finishIssueID)
	if len(finishCompletions) != 1 {
		t.Fatalf("finish path completions = %#v, want 1 row", finishCompletions)
	}
	if finishCompletions[0].AttemptNumber != 1 || finishCompletions[0].Status != "dev_complete" {
		t.Fatalf("finish path completion row = %#v, want attempt 1/dev_complete", finishCompletions[0])
	}
	if finishCompletions[0].NotesMDPath != finishReportPath || finishCompletions[0].ActualTokens != 14823 {
		t.Fatalf("finish path completion row = %#v, want notes %q and 14823 tokens", finishCompletions[0], finishReportPath)
	}
	if finishCompletions[0].GatesPassed != 1 || finishCompletions[0].GatesFailed != 0 {
		t.Fatalf("finish path completion gate counts = (%d, %d), want (1, 0)", finishCompletions[0].GatesPassed, finishCompletions[0].GatesFailed)
	}
	var finishGateSnapshot []v01EGateResultSnapshot
	if err := json.Unmarshal([]byte(finishCompletions[0].GateResults), &finishGateSnapshot); err != nil {
		t.Fatalf("Unmarshal(finish gate_results) error = %v", err)
	}
	if len(finishGateSnapshot) != 1 || finishGateSnapshot[0].Cmd != "go test ./..." || finishGateSnapshot[0].ExitCode != 0 {
		t.Fatalf("finish path gate_results = %#v, want one go test success snapshot", finishGateSnapshot)
	}

	_, startBlockStderr := runFIDCommand(
		t,
		workspace,
		"dev", "start", blockIssueID,
		"--model", "claude-sonnet-4-6",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, startBlockStderr)
	for _, want := range []string{
		"dev start " + blockIssueID,
		"claimed by developer",
		"transitioned: todo → in_progress",
		"agent run: RUN-",
	} {
		if !strings.Contains(startBlockStderr, want) {
			t.Fatalf("block path dev start stderr = %q, want substring %q", startBlockStderr, want)
		}
	}
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "dev start block path")
	assertV01EIssueState(t, workspace, blockIssueID, "in_progress", "started")
	blockClaims := mustReadIssueClaimsSection(t, workspace, blockIssueID)
	if len(blockClaims) != 1 || blockClaims[0].ClaimedBy != "developer" {
		t.Fatalf("block path claims after dev start = %#v, want one developer claim", blockClaims)
	}
	blockRuns := mustReadV01ERunsForIssue(t, workspace, blockIssueID)
	if len(blockRuns) != 1 || blockRuns[0].EndedAt != "" {
		t.Fatalf("block path runs after dev start = %#v, want one open run", blockRuns)
	}
	blockRunID := blockRuns[0].ID

	_, blockStderr := runFIDCommand(
		t,
		workspace,
		"dev", "block", blockIssueID,
		"--reason-md", blockReasonPath,
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, blockStderr)
	for _, want := range []string{
		"dev block " + blockIssueID,
		"transitioned: in_progress → failed_gate",
		"completion: TC-",
		"claim released",
	} {
		if !strings.Contains(blockStderr, want) {
			t.Fatalf("dev block stderr = %q, want substring %q", blockStderr, want)
		}
	}
	expectedJournalLines++
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "dev block")
	assertV01EIssueState(t, workspace, blockIssueID, "failed_gate", "started")
	if claims := mustReadIssueClaimsSection(t, workspace, blockIssueID); len(claims) != 0 {
		t.Fatalf("block path claims after dev block = %#v, want empty", claims)
	}
	blockRuns = mustReadV01ERunsForIssue(t, workspace, blockIssueID)
	if len(blockRuns) != 1 || blockRuns[0].ID != blockRunID || blockRuns[0].EndedAt != "" {
		t.Fatalf("block path runs after dev block = %#v, want same open run", blockRuns)
	}
	blockCompletions := mustReadV01ECompletionsForIssue(t, workspace, blockIssueID)
	if len(blockCompletions) != 1 {
		t.Fatalf("block path completions = %#v, want 1 row", blockCompletions)
	}
	if blockCompletions[0].AttemptNumber != 1 || blockCompletions[0].Status != "failed_gate" {
		t.Fatalf("block path completion row = %#v, want attempt 1/failed_gate", blockCompletions[0])
	}
	if blockCompletions[0].NotesMDPath != blockReasonPath {
		t.Fatalf("block path completion notes = %q, want %q", blockCompletions[0].NotesMDPath, blockReasonPath)
	}
	if blockCompletions[0].GatesPassed != 0 || blockCompletions[0].GatesFailed != 0 {
		t.Fatalf("block path completion gate counts = (%d, %d), want (0, 0)", blockCompletions[0].GatesPassed, blockCompletions[0].GatesFailed)
	}
	var blockGateSnapshot []v01EGateResultSnapshot
	if err := json.Unmarshal([]byte(blockCompletions[0].GateResults), &blockGateSnapshot); err != nil {
		t.Fatalf("Unmarshal(block gate_results) error = %v", err)
	}
	if len(blockGateSnapshot) != 0 {
		t.Fatalf("block path gate_results = %#v, want empty snapshot", blockGateSnapshot)
	}

	beforeRebuild := mustSnapshotV01EState(t, workspace, epicID, finishIssueID, blockIssueID)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStdout, rebuildStderr := runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")
	assertNoNestedLockFailure(t, rebuildStderr)
	if !strings.Contains(rebuildStdout, "rebuild complete") || rebuildStderr != "" {
		t.Fatalf("rebuild stdout/stderr = %q/%q, want stdout completion and empty stderr", rebuildStdout, rebuildStderr)
	}
	assertV01EJournalLines(t, journalPath, expectedJournalLines, "rebuild")

	if got := mustSnapshotV01EState(t, workspace, epicID, finishIssueID, blockIssueID); !reflect.DeepEqual(got, beforeRebuild) {
		t.Fatalf("projected state after rebuild = %#v, want %#v", got, beforeRebuild)
	}
}

type v01EIssueView struct {
	Status         string `json:"status"`
	StatusCategory string `json:"status_category"`
}

type v01EStateSnapshot struct {
	Issues      []v01EIssueRow
	Claims      []v01EClaimRow
	Completions []v01ECompletionRow
	Runs        []v01ERunRow
	GateRuns    []v01EGateRunRow
}

type v01EIssueRow struct {
	ID         string
	Title      string
	Kind       string
	ParentID   string
	WorkflowID string
	Status     string
	OwnerRole  string
}

type v01EClaimRow struct {
	RefType    string
	RefID      string
	ClaimedBy  string
	ClaimedAt  string
	TTLSeconds int
	Note       string
}

type v01ECompletionRow struct {
	ID                    string
	IssueID               string
	AttemptNumber         int
	Status                string
	NotesMDPath           string
	ActualTokens          int
	CommitHashes          string
	FilesTouched          string
	GatesPassed           int
	GatesFailed           int
	GateResults           string
	HumanIntervention     bool
	HumanInterventionNote string
	CompletedAt           string
}

type v01ERunRow struct {
	ID           string
	IssueID      string
	CompletionID string
	Role         string
	Model        string
	StartedAt    string
	EndedAt      string
}

type v01EGateRunRow struct {
	ID            string
	IssueID       string
	AttemptNumber int
	AgentRunID    string
	RecordedBy    string
	Command       string
	ExitCode      int
	OutputExcerpt string
	RanAt         string
}

type v01EGateResultSnapshot struct {
	Cmd           string  `json:"cmd"`
	ExitCode      int     `json:"exit_code"`
	RanAt         string  `json:"ran_at"`
	OutputExcerpt *string `json:"output_excerpt"`
}

func writeV01EMarkdown(t *testing.T, workspace, relativePath, body string) string {
	t.Helper()

	path := filepath.Join(workspace, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
	return filepath.ToSlash(relativePath)
}

func assertV01EJournalLines(t *testing.T, journalPath string, want int, step string) {
	t.Helper()

	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != want {
		t.Fatalf("journal lines after %s = %d, want %d", step, got, want)
	}
}

func assertV01EIssueState(t *testing.T, workspace, issueID, wantStatus, wantCategory string) {
	t.Helper()

	var got v01EIssueView
	if err := json.Unmarshal([]byte(mustGetIssueJSON(t, workspace, issueID)), &got); err != nil {
		t.Fatalf("Unmarshal(issue %s) error = %v", issueID, err)
	}
	if got.Status != wantStatus {
		t.Fatalf("issue %s status = %q, want %q", issueID, got.Status, wantStatus)
	}
	if wantCategory != "" && got.StatusCategory != wantCategory {
		t.Fatalf("issue %s status_category = %q, want %q", issueID, got.StatusCategory, wantCategory)
	}
}

func mustReadV01ERunsForIssue(t *testing.T, workspace, issueID string) []v01ERunRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT id, issue_id, COALESCE(completion_id, ''), role, model, COALESCE(started_at, ''), COALESCE(ended_at, '')
		FROM agent_runs
		WHERE issue_id = ?
		ORDER BY id
	`, issueID)
	if err != nil {
		t.Fatalf("query agent_runs for %s: %v", issueID, err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []v01ERunRow
	for rows.Next() {
		var row v01ERunRow
		if err := rows.Scan(&row.ID, &row.IssueID, &row.CompletionID, &row.Role, &row.Model, &row.StartedAt, &row.EndedAt); err != nil {
			t.Fatalf("scan agent_runs for %s: %v", issueID, err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate agent_runs for %s: %v", issueID, err)
	}
	return snapshot
}

func mustReadV01EGateRunsForIssue(t *testing.T, workspace, issueID string) []v01EGateRunRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT id, issue_id, attempt_number, COALESCE(agent_run_id, ''), recorded_by_role, command, exit_code, COALESCE(output_excerpt, ''), ran_at
		FROM gate_runs
		WHERE issue_id = ?
		ORDER BY id
	`, issueID)
	if err != nil {
		t.Fatalf("query gate_runs for %s: %v", issueID, err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []v01EGateRunRow
	for rows.Next() {
		var row v01EGateRunRow
		if err := rows.Scan(&row.ID, &row.IssueID, &row.AttemptNumber, &row.AgentRunID, &row.RecordedBy, &row.Command, &row.ExitCode, &row.OutputExcerpt, &row.RanAt); err != nil {
			t.Fatalf("scan gate_runs for %s: %v", issueID, err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate gate_runs for %s: %v", issueID, err)
	}
	return snapshot
}

func mustReadV01ECompletionsForIssue(t *testing.T, workspace, issueID string) []v01ECompletionRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT id, issue_id, attempt_number, COALESCE(status, ''), COALESCE(notes_md_path, ''), COALESCE(actual_tokens, -1),
			COALESCE(commit_hashes, ''), COALESCE(files_touched, ''), gates_passed, gates_failed, COALESCE(gate_results, ''),
			human_intervention, COALESCE(human_intervention_note, ''), COALESCE(completed_at, '')
		FROM task_completions
		WHERE issue_id = ?
		ORDER BY id
	`, issueID)
	if err != nil {
		t.Fatalf("query task_completions for %s: %v", issueID, err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []v01ECompletionRow
	for rows.Next() {
		var row v01ECompletionRow
		if err := rows.Scan(
			&row.ID,
			&row.IssueID,
			&row.AttemptNumber,
			&row.Status,
			&row.NotesMDPath,
			&row.ActualTokens,
			&row.CommitHashes,
			&row.FilesTouched,
			&row.GatesPassed,
			&row.GatesFailed,
			&row.GateResults,
			&row.HumanIntervention,
			&row.HumanInterventionNote,
			&row.CompletedAt,
		); err != nil {
			t.Fatalf("scan task_completions for %s: %v", issueID, err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate task_completions for %s: %v", issueID, err)
	}
	return snapshot
}

func mustSnapshotV01EState(t *testing.T, workspace, epicID, finishIssueID, blockIssueID string) v01EStateSnapshot {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	return v01EStateSnapshot{
		Issues:      mustSnapshotV01EIssues(t, conn, epicID, finishIssueID, blockIssueID),
		Claims:      mustSnapshotV01EClaims(t, conn, epicID, finishIssueID, blockIssueID),
		Completions: mustSnapshotV01ECompletions(t, conn, finishIssueID, blockIssueID),
		Runs:        mustSnapshotV01ERuns(t, conn, finishIssueID, blockIssueID),
		GateRuns:    mustSnapshotV01EGateRuns(t, conn, finishIssueID, blockIssueID),
	}
}

func mustSnapshotV01EIssues(t *testing.T, conn *sql.DB, epicID, finishIssueID, blockIssueID string) []v01EIssueRow {
	t.Helper()

	rows, err := conn.Query(`
		SELECT id, title, kind, COALESCE(parent_id, ''), workflow_id, status, COALESCE(owner_role, '')
		FROM issues
		WHERE id IN (?, ?, ?)
		ORDER BY id
	`, epicID, finishIssueID, blockIssueID)
	if err != nil {
		t.Fatalf("query issues snapshot: %v", err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []v01EIssueRow
	for rows.Next() {
		var row v01EIssueRow
		if err := rows.Scan(&row.ID, &row.Title, &row.Kind, &row.ParentID, &row.WorkflowID, &row.Status, &row.OwnerRole); err != nil {
			t.Fatalf("scan issues snapshot: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate issues snapshot: %v", err)
	}
	return snapshot
}

func mustSnapshotV01EClaims(t *testing.T, conn *sql.DB, epicID, finishIssueID, blockIssueID string) []v01EClaimRow {
	t.Helper()

	rows, err := conn.Query(`
		SELECT ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, COALESCE(note, '')
		FROM claims
		WHERE ref_type = 'issue' AND ref_id IN (?, ?, ?)
		ORDER BY ref_type, ref_id, claimed_by, claimed_at
	`, epicID, finishIssueID, blockIssueID)
	if err != nil {
		t.Fatalf("query claims snapshot: %v", err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []v01EClaimRow
	for rows.Next() {
		var row v01EClaimRow
		if err := rows.Scan(&row.RefType, &row.RefID, &row.ClaimedBy, &row.ClaimedAt, &row.TTLSeconds, &row.Note); err != nil {
			t.Fatalf("scan claims snapshot: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate claims snapshot: %v", err)
	}
	return snapshot
}

func mustSnapshotV01ECompletions(t *testing.T, conn *sql.DB, finishIssueID, blockIssueID string) []v01ECompletionRow {
	t.Helper()

	rows, err := conn.Query(`
		SELECT id, issue_id, attempt_number, COALESCE(status, ''), COALESCE(notes_md_path, ''), COALESCE(actual_tokens, -1),
			COALESCE(commit_hashes, ''), COALESCE(files_touched, ''), gates_passed, gates_failed, COALESCE(gate_results, ''),
			human_intervention, COALESCE(human_intervention_note, ''), COALESCE(completed_at, '')
		FROM task_completions
		WHERE issue_id IN (?, ?)
		ORDER BY id
	`, finishIssueID, blockIssueID)
	if err != nil {
		t.Fatalf("query task_completions snapshot: %v", err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []v01ECompletionRow
	for rows.Next() {
		var row v01ECompletionRow
		if err := rows.Scan(
			&row.ID,
			&row.IssueID,
			&row.AttemptNumber,
			&row.Status,
			&row.NotesMDPath,
			&row.ActualTokens,
			&row.CommitHashes,
			&row.FilesTouched,
			&row.GatesPassed,
			&row.GatesFailed,
			&row.GateResults,
			&row.HumanIntervention,
			&row.HumanInterventionNote,
			&row.CompletedAt,
		); err != nil {
			t.Fatalf("scan task_completions snapshot: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate task_completions snapshot: %v", err)
	}
	return snapshot
}

func mustSnapshotV01ERuns(t *testing.T, conn *sql.DB, finishIssueID, blockIssueID string) []v01ERunRow {
	t.Helper()

	rows, err := conn.Query(`
		SELECT id, issue_id, COALESCE(completion_id, ''), role, model, COALESCE(started_at, ''), COALESCE(ended_at, '')
		FROM agent_runs
		WHERE issue_id IN (?, ?)
		ORDER BY id
	`, finishIssueID, blockIssueID)
	if err != nil {
		t.Fatalf("query agent_runs snapshot: %v", err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []v01ERunRow
	for rows.Next() {
		var row v01ERunRow
		if err := rows.Scan(&row.ID, &row.IssueID, &row.CompletionID, &row.Role, &row.Model, &row.StartedAt, &row.EndedAt); err != nil {
			t.Fatalf("scan agent_runs snapshot: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate agent_runs snapshot: %v", err)
	}
	return snapshot
}

func mustSnapshotV01EGateRuns(t *testing.T, conn *sql.DB, finishIssueID, blockIssueID string) []v01EGateRunRow {
	t.Helper()

	rows, err := conn.Query(`
		SELECT id, issue_id, attempt_number, COALESCE(agent_run_id, ''), recorded_by_role, command, exit_code, COALESCE(output_excerpt, ''), ran_at
		FROM gate_runs
		WHERE issue_id IN (?, ?)
		ORDER BY id
	`, finishIssueID, blockIssueID)
	if err != nil {
		t.Fatalf("query gate_runs snapshot: %v", err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []v01EGateRunRow
	for rows.Next() {
		var row v01EGateRunRow
		if err := rows.Scan(&row.ID, &row.IssueID, &row.AttemptNumber, &row.AgentRunID, &row.RecordedBy, &row.Command, &row.ExitCode, &row.OutputExcerpt, &row.RanAt); err != nil {
			t.Fatalf("scan gate_runs snapshot: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate gate_runs snapshot: %v", err)
	}
	return snapshot
}
