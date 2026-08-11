package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV01CConcernRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Phase 9 concern rebuild parity",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	_, firstCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	firstCompletionID := parseCompletionID(t, firstCompleteStderr)

	_, raiseStderr := runFIDCommand(
		t,
		workspace,
		"concern", "raise",
		"--on", issueID,
		"--by", "developer",
		"--in-completion", firstCompletionID,
		"--summary", "First concern for parity",
	)
	firstConcernID := parseConcernID(t, raiseStderr)

	_, secondCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	secondCompletionID := parseCompletionID(t, secondCompleteStderr)

	_, addressStderr := runFIDCommand(
		t,
		workspace,
		"concern", "address", firstConcernID,
		"--in-completion", secondCompletionID,
		"--by", "developer",
	)
	if !strings.Contains(addressStderr, "addressed concern "+firstConcernID) {
		t.Fatalf("address stderr = %q, want concern confirmation", addressStderr)
	}

	_, thirdCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	thirdCompletionID := parseCompletionID(t, thirdCompleteStderr)

	_, raiseStderr = runFIDCommand(
		t,
		workspace,
		"concern", "raise",
		"--on", issueID,
		"--by", "developer",
		"--in-completion", thirdCompletionID,
		"--summary", "Second concern for parity",
	)
	secondConcernID := parseConcernID(t, raiseStderr)

	beforeConcernRows := snapshotConcernRows(t, workspace)
	beforeConcernList := mustGetConcernListJSON(t, workspace)
	beforeIssueJSON := mustGetIssueJSON(t, workspace, issueID)
	beforeIssueListJSON := mustGetIssueListJSON(t, workspace)

	if len(beforeConcernRows) != 2 {
		t.Fatalf("before concern rows len = %d, want 2", len(beforeConcernRows))
	}
	firstConcernRow := concernRowByID(t, beforeConcernRows, firstConcernID)
	if firstConcernRow.Status != "addressed" {
		t.Fatalf("before first concern row = %#v, want addressed %q", firstConcernRow, firstConcernID)
	}
	if firstConcernRow.AddressedInCompletionID == nil || *firstConcernRow.AddressedInCompletionID != secondCompletionID {
		t.Fatalf("before first addressed completion = %#v, want %q", firstConcernRow.AddressedInCompletionID, secondCompletionID)
	}
	secondConcernRow := concernRowByID(t, beforeConcernRows, secondConcernID)
	if secondConcernRow.Status != "open" {
		t.Fatalf("before second concern row = %#v, want open %q", secondConcernRow, secondConcernID)
	}

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStdout, rebuildStderr := runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")
	if !strings.Contains(rebuildStdout, "rebuild complete") || rebuildStderr != "" {
		t.Fatalf("rebuild stdout/stderr = %q/%q, want stdout completion and empty stderr", rebuildStdout, rebuildStderr)
	}

	if got := snapshotConcernRows(t, workspace); !reflect.DeepEqual(got, beforeConcernRows) {
		t.Fatalf("concern rows after rebuild = %#v, want %#v", got, beforeConcernRows)
	}
	if got := mustGetConcernListJSON(t, workspace); got != beforeConcernList {
		t.Fatalf("concern list after rebuild = %q, want %q", got, beforeConcernList)
	}
	if got := mustGetIssueJSON(t, workspace, issueID); got != beforeIssueJSON {
		t.Fatalf("issue JSON after rebuild = %q, want %q", got, beforeIssueJSON)
	}
	if got := mustGetIssueListJSON(t, workspace); got != beforeIssueListJSON {
		t.Fatalf("issue list JSON after rebuild = %q, want %q", got, beforeIssueListJSON)
	}
}

func TestV01CClosureSmoke(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Phase 9 v0.1c closure smoke",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	_, addStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Concern lifecycle is visible end to end",
		"--ordinal", "1",
		"--by", "developer",
	)
	acID := parseAddedAcceptanceCriterionID(t, addStderr)

	_, verifyStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, acID,
		"--by", "qa_lead",
		"--status", "pass",
	)
	if !strings.Contains(verifyStderr, "verified "+acID+" on "+issueID+" by qa_lead: pass") {
		t.Fatalf("verify stderr = %q, want confirmation", verifyStderr)
	}

	acListStdout, acListStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "list", issueID,
		"--format", "json",
	)
	if acListStderr != "" {
		t.Fatalf("ac list stderr = %q, want empty", acListStderr)
	}
	criteria := decodeCommandJSONStructure(t, acListStdout).([]any)
	if len(criteria) != 1 {
		t.Fatalf("ac list criteria len = %d, want 1", len(criteria))
	}

	if got := decodeIssueConcernState(t, mustGetIssueJSON(t, workspace, issueID)); got.ConcernsRaised || got.OpenConcernsCount != 0 {
		t.Fatalf("before raise concern state = %#v, want false/0", got)
	}
	if got := decodeIssueListConcernState(t, mustGetIssueListJSON(t, workspace), issueID); got.ConcernsRaised || got.OpenConcernsCount != 0 {
		t.Fatalf("issue list before raise concern state = %#v, want false/0", got)
	}

	_, firstCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	firstCompletionID := parseCompletionID(t, firstCompleteStderr)

	_, raiseStderr := runFIDCommand(
		t,
		workspace,
		"concern", "raise",
		"--on", issueID,
		"--by", "developer",
		"--in-completion", firstCompletionID,
		"--summary", "Smoke concern",
	)
	concernID := parseConcernID(t, raiseStderr)

	if got := decodeIssueConcernState(t, mustGetIssueJSON(t, workspace, issueID)); !got.ConcernsRaised || got.OpenConcernsCount != 1 {
		t.Fatalf("after raise concern state = %#v, want true/1", got)
	}
	if got := decodeIssueListConcernState(t, mustGetIssueListJSON(t, workspace), issueID); !got.ConcernsRaised || got.OpenConcernsCount != 1 {
		t.Fatalf("issue list after raise concern state = %#v, want true/1", got)
	}

	openConcernsStdout, openConcernsStderr := runFIDCommand(
		t,
		workspace,
		"concern", "list",
		"--issue", issueID,
		"--status", "open",
		"--format", "json",
	)
	if openConcernsStderr != "" {
		t.Fatalf("open concern list stderr = %q, want empty", openConcernsStderr)
	}
	openConcerns := decodeCommandJSONStructure(t, openConcernsStdout).([]any)
	if len(openConcerns) != 1 {
		t.Fatalf("open concern list len = %d, want 1", len(openConcerns))
	}
	if got := openConcerns[0].(map[string]any)["id"]; got != concernID {
		t.Fatalf("open concern id = %v, want %q", got, concernID)
	}

	_, secondCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	secondCompletionID := parseCompletionID(t, secondCompleteStderr)

	_, addressStderr := runFIDCommand(
		t,
		workspace,
		"concern", "address", concernID,
		"--in-completion", secondCompletionID,
		"--by", "developer",
	)
	if !strings.Contains(addressStderr, "addressed concern "+concernID) {
		t.Fatalf("address stderr = %q, want concern confirmation", addressStderr)
	}

	if got := decodeIssueConcernState(t, mustGetIssueJSON(t, workspace, issueID)); got.ConcernsRaised || got.OpenConcernsCount != 0 {
		t.Fatalf("after address concern state = %#v, want false/0", got)
	}
	if got := decodeIssueListConcernState(t, mustGetIssueListJSON(t, workspace), issueID); got.ConcernsRaised || got.OpenConcernsCount != 0 {
		t.Fatalf("issue list after address concern state = %#v, want false/0", got)
	}

	addressedConcernsStdout, addressedConcernsStderr := runFIDCommand(
		t,
		workspace,
		"concern", "list",
		"--issue", issueID,
		"--status", "addressed",
		"--format", "json",
	)
	if addressedConcernsStderr != "" {
		t.Fatalf("addressed concern list stderr = %q, want empty", addressedConcernsStderr)
	}
	addressedConcerns := decodeCommandJSONStructure(t, addressedConcernsStdout).([]any)
	if len(addressedConcerns) != 1 {
		t.Fatalf("addressed concern list len = %d, want 1", len(addressedConcerns))
	}
	if got := addressedConcerns[0].(map[string]any)["id"]; got != concernID {
		t.Fatalf("addressed concern id = %v, want %q", got, concernID)
	}
	if got := addressedConcerns[0].(map[string]any)["status"]; got != "addressed" {
		t.Fatalf("addressed concern status = %v, want addressed", got)
	}
	if got := addressedConcerns[0].(map[string]any)["addressed_in_completion_id"]; got != secondCompletionID {
		t.Fatalf("addressed concern addressed_in_completion_id = %v, want %q", got, secondCompletionID)
	}

	issueConcernsStdout, issueConcernsStderr := runFIDCommand(
		t,
		workspace,
		"issue", "get", issueID,
		"--section", "concerns",
		"--format", "json",
	)
	if issueConcernsStderr != "" {
		t.Fatalf("issue concerns stderr = %q, want empty", issueConcernsStderr)
	}
	issueConcerns := decodeCommandJSONStructure(t, issueConcernsStdout).([]any)
	if len(issueConcerns) != 1 {
		t.Fatalf("issue concerns len = %d, want 1", len(issueConcerns))
	}
	if got := issueConcerns[0].(map[string]any)["addressed_in_completion_id"]; got != secondCompletionID {
		t.Fatalf("issue concern addressed_in_completion_id = %v, want %q", got, secondCompletionID)
	}
}

type concernParityRow struct {
	ID                      string
	IssueID                 string
	Status                  string
	RaisedInCompletionID    string
	AddressedInCompletionID *string
}

type issueConcernState struct {
	ConcernsRaised    bool `json:"concerns_raised"`
	OpenConcernsCount int  `json:"open_concerns_count"`
}

func snapshotConcernRows(t *testing.T, workspace string) []concernParityRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT id, issue_id, status, raised_in_completion_id, addressed_in_completion_id
		FROM concerns
		ORDER BY id
	`)
	if err != nil {
		t.Fatalf("query concerns: %v", err)
	}
	defer closeTask4Rows(t, rows)

	var snapshot []concernParityRow
	for rows.Next() {
		var row concernParityRow
		var addressedInCompletionID *string
		if err := rows.Scan(
			&row.ID,
			&row.IssueID,
			&row.Status,
			&row.RaisedInCompletionID,
			&addressedInCompletionID,
		); err != nil {
			t.Fatalf("scan concerns: %v", err)
		}
		row.AddressedInCompletionID = addressedInCompletionID
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate concerns: %v", err)
	}
	return snapshot
}

func mustGetConcernListJSON(t *testing.T, workspace string) string {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "concern", "list", "--format", "json")
	if stderr != "" {
		t.Fatalf("concern list stderr = %q, want empty", stderr)
	}
	return stdout
}

func mustGetIssueListJSON(t *testing.T, workspace string) string {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "issue", "list", "--format", "json")
	if stderr != "" {
		t.Fatalf("issue list stderr = %q, want empty", stderr)
	}
	return stdout
}

func decodeIssueConcernState(t *testing.T, issueJSON string) issueConcernState {
	t.Helper()

	var result issueConcernState
	if err := json.Unmarshal([]byte(issueJSON), &result); err != nil {
		t.Fatalf("Unmarshal(issue) error = %v", err)
	}
	return result
}

func decodeIssueListConcernState(t *testing.T, listJSON, issueID string) issueConcernState {
	t.Helper()

	rows := decodeCommandJSONStructure(t, listJSON).([]any)
	for _, row := range rows {
		object, ok := row.(map[string]any)
		if !ok {
			continue
		}
		if object["id"] != issueID {
			continue
		}
		openConcernsCount, err := object["open_concerns_count"].(json.Number).Int64()
		if err != nil {
			t.Fatalf("open_concerns_count for %q error = %v", issueID, err)
		}
		return issueConcernState{
			ConcernsRaised:    object["concerns_raised"].(bool),
			OpenConcernsCount: int(openConcernsCount),
		}
	}
	t.Fatalf("issue %q not found in issue list", issueID)
	return issueConcernState{}
}

func concernRowByID(t *testing.T, rows []concernParityRow, concernID string) concernParityRow {
	t.Helper()

	for _, row := range rows {
		if row.ID == concernID {
			return row
		}
	}

	t.Fatalf("concern %q not found in rows %#v", concernID, rows)
	return concernParityRow{}
}
