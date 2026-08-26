package main

import (
	"encoding/json"
	"os"
	"path/filepath"
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
		"--actor", "agent:developer",
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

	firstConcernID := "CON-v01c-first"
	assertConcernMutationUnavailable(t, workspace, []string{"concern", "raise", "--on", issueID, "--actor", "agent:developer", "--in-completion", firstCompletionID, "--summary", "First concern for parity"})
	seedConcernProjection(t, workspace, firstConcernID, issueID, firstCompletionID, "First concern for parity", "open", nil)

	_, secondCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	secondCompletionID := parseCompletionID(t, secondCompleteStderr)

	assertConcernMutationUnavailable(t, workspace, []string{"concern", "address", firstConcernID, "--actor", "human:reviewer", "--in-completion", secondCompletionID})
	updateConcernProjection(t, workspace, firstConcernID, secondCompletionID)

	_, thirdCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	thirdCompletionID := parseCompletionID(t, thirdCompleteStderr)

	secondConcernID := "CON-v01c-second"
	assertConcernMutationUnavailable(t, workspace, []string{"concern", "raise", "--on", issueID, "--actor", "agent:developer", "--in-completion", thirdCompletionID, "--summary", "Second concern for parity"})
	seedConcernProjection(t, workspace, secondConcernID, issueID, thirdCompletionID, "Second concern for parity", "open", nil)

	beforeConcernRows := snapshotConcernRows(t, workspace)

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

	rebuildStdout, rebuildStderr := runFIDCommand(t, workspace, "rebuild")
	if !strings.Contains(rebuildStdout, "rebuild complete") || rebuildStderr != "" {
		t.Fatalf("rebuild stdout/stderr = %q/%q, want stdout completion and empty stderr", rebuildStdout, rebuildStderr)
	}

	if got := snapshotConcernRows(t, workspace); len(got) != 0 {
		t.Fatalf("concern rows after rebuild = %#v, want no non-journal historical rows", got)
	}
	if got := decodeConcernListState(t, mustGetConcernListJSON(t, workspace)); got.ConcernsRaised || got.OpenConcernsCount != 0 {
		t.Fatalf("concern list state after rebuild = %#v, want false/0 for non-journal historical rows", got)
	}
	assertIssueConcernState(t, mustGetIssueJSON(t, workspace, issueID), false, 0)
	assertIssueListConcernState(t, mustGetIssueListJSON(t, workspace), issueID, false, 0)
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
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	addStdout, _ := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Concern lifecycle is visible end to end",
		"--ordinal", "1",
		"--by", "developer",
	)
	acID := parseAddedAcceptanceCriterionID(t, addStdout)

	assertIssueACVerifyUnavailable(t, workspace, issueID, acID, "agent:qa_lead")
	seedHistoricalACVerification(t, workspace, acID, "qa_lead")

	acListStdout, acListStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "list", issueID,
		"--json",
	)
	if acListStderr != "" {
		t.Fatalf("ac list stderr = %q, want empty", acListStderr)
	}
	criteriaValue := decodeReadEnvelope(t, acListStdout, "issue.ac.list").Result["acceptance_criteria"]
	criteria, ok := criteriaValue.([]any)
	if !ok {
		t.Fatalf("acceptance_criteria = %#v, want array", criteriaValue)
	}
	if len(criteria) != 1 {
		t.Fatalf("ac list criteria len = %d, want 1", len(criteria))
	}

	if got := decodeConcernListState(t, mustGetConcernListJSON(t, workspace)); got.ConcernsRaised || got.OpenConcernsCount != 0 {
		t.Fatalf("before raise concern state = %#v, want false/0", got)
	}

	_, firstCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	firstCompletionID := parseCompletionID(t, firstCompleteStderr)

	concernID := "CON-v01c-smoke"
	assertConcernMutationUnavailable(t, workspace, []string{"concern", "raise", "--on", issueID, "--actor", "agent:developer", "--in-completion", firstCompletionID, "--summary", "Smoke concern"})
	seedConcernProjection(t, workspace, concernID, issueID, firstCompletionID, "Smoke concern", "open", nil)

	if got := decodeConcernListState(t, mustGetConcernListJSON(t, workspace)); !got.ConcernsRaised || got.OpenConcernsCount != 1 {
		t.Fatalf("after raise concern state = %#v, want true/1", got)
	}

	openConcernsStdout, openConcernsStderr := runFIDCommand(
		t,
		workspace,
		"concern", "list",
		"--issue", issueID,
		"--status", "open",
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

	assertConcernMutationUnavailable(t, workspace, []string{"concern", "address", concernID, "--actor", "human:reviewer", "--in-completion", secondCompletionID})
	updateConcernProjection(t, workspace, concernID, secondCompletionID)

	if got := decodeConcernListState(t, mustGetConcernListJSON(t, workspace)); got.ConcernsRaised || got.OpenConcernsCount != 0 {
		t.Fatalf("after address concern state = %#v, want false/0", got)
	}

	addressedConcernsStdout, addressedConcernsStderr := runFIDCommand(
		t,
		workspace,
		"concern", "list",
		"--issue", issueID,
		"--status", "addressed",
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

	issueConcernsStdout, issueConcernsStderr := runFIDCommand(t, workspace, "concern", "list", "--issue", issueID)
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

	stdout, stderr := runFIDCommand(t, workspace, "concern", "list")
	if stderr != "" {
		t.Fatalf("concern list stderr = %q, want empty", stderr)
	}
	return stdout
}

func mustGetIssueListJSON(t *testing.T, workspace string) string {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "--json", "issue", "list")
	if stderr != "" {
		t.Fatalf("issue list stderr = %q, want empty", stderr)
	}
	result := decodeCanonicalResult(t, stdout)
	payload, err := json.Marshal(result["issues"])
	if err != nil {
		t.Fatalf("Marshal(issue list result) error = %v", err)
	}
	return string(payload)
}

func decodeConcernListState(t *testing.T, listJSON string) issueConcernState {
	t.Helper()
	rows := decodeCommandJSONStructure(t, listJSON).([]any)
	open := 0
	for _, row := range rows {
		object, ok := row.(map[string]any)
		if ok && object["status"] == "open" {
			open++
		}
	}
	return issueConcernState{ConcernsRaised: open > 0, OpenConcernsCount: open}
}

func assertIssueConcernState(t *testing.T, issueJSON string, wantRaised bool, wantOpen int) {
	t.Helper()

	var issue struct {
		ConcernsRaised    bool `json:"concerns_raised"`
		OpenConcernsCount int  `json:"open_concerns_count"`
	}
	if err := json.Unmarshal([]byte(issueJSON), &issue); err != nil {
		t.Fatalf("Unmarshal(issue concern state) error = %v", err)
	}
	if issue.ConcernsRaised != wantRaised || issue.OpenConcernsCount != wantOpen {
		t.Fatalf("issue concern state = %#v, want %t/%d", issue, wantRaised, wantOpen)
	}
}

func assertIssueListConcernState(t *testing.T, issueListJSON, issueID string, wantRaised bool, wantOpen int) {
	t.Helper()

	var issues []struct {
		ID                string `json:"id"`
		ConcernsRaised    bool   `json:"concerns_raised"`
		OpenConcernsCount int    `json:"open_concerns_count"`
	}
	if err := json.Unmarshal([]byte(issueListJSON), &issues); err != nil {
		t.Fatalf("Unmarshal(issue list concern state) error = %v", err)
	}
	for _, issue := range issues {
		if issue.ID == issueID {
			if issue.ConcernsRaised != wantRaised || issue.OpenConcernsCount != wantOpen {
				t.Fatalf("issue list concern state = %#v, want %t/%d", issue, wantRaised, wantOpen)
			}
			return
		}
	}
	t.Fatalf("issue %q not found in issue list", issueID)
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
