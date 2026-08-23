package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConcernRaiseCreatesConcernAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Raise CLI concern",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	_, completeStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	completionID := parseCompletionID(t, completeStderr)

	bodyPath := filepath.Join(workspace, "concern.md")
	if err := os.WriteFile(bodyPath, []byte("# Concern\n\nNeeds follow-up.\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", bodyPath, err)
	}

	_, stderr := runFIDCommand(
		t,
		workspace,
		"concern", "raise",
		"--on", issueID,
		"--by", "developer",
		"--in-completion", completionID,
		"--summary", "Potential race condition",
		"--body-md", bodyPath,
	)

	if !strings.Contains(stderr, "raised concern CON-") {
		t.Fatalf("stderr = %q, want concern confirmation", stderr)
	}
	if !strings.Contains(stderr, issueID) {
		t.Fatalf("stderr = %q, want issue id", stderr)
	}
	if !strings.Contains(stderr, `"Potential race condition"`) {
		t.Fatalf("stderr = %q, want summary", stderr)
	}

	concernsStdout, concernsStderr := runFIDCommand(t, workspace, "concern", "list", "--issue", issueID)
	if concernsStderr != "" {
		t.Fatalf("issue get concerns stderr = %q, want empty", concernsStderr)
	}

	var concerns []struct {
		ID                   string  `json:"id"`
		RaisedInCompletionID string  `json:"raised_in_completion_id"`
		RaisedByRole         string  `json:"raised_by_role"`
		Summary              string  `json:"summary"`
		BodyMDPath           *string `json:"body_md_path"`
		Status               string  `json:"status"`
	}
	if err := json.Unmarshal([]byte(concernsStdout), &concerns); err != nil {
		t.Fatalf("Unmarshal(concerns) error = %v", err)
	}
	if len(concerns) != 1 {
		t.Fatalf("concerns len = %d, want 1", len(concerns))
	}
	if !strings.HasPrefix(concerns[0].ID, "CON-") {
		t.Fatalf("concern id = %q, want CON- prefix", concerns[0].ID)
	}
	if concerns[0].RaisedInCompletionID != completionID {
		t.Fatalf("raised_in_completion_id = %q, want %q", concerns[0].RaisedInCompletionID, completionID)
	}
	if concerns[0].RaisedByRole != "developer" {
		t.Fatalf("raised_by_role = %q, want developer", concerns[0].RaisedByRole)
	}
	if concerns[0].Summary != "Potential race condition" {
		t.Fatalf("summary = %q, want summary", concerns[0].Summary)
	}
	if concerns[0].BodyMDPath == nil || *concerns[0].BodyMDPath != bodyPath {
		t.Fatalf("body_md_path = %#v, want %q", concerns[0].BodyMDPath, bodyPath)
	}
	if concerns[0].Status != "open" {
		t.Fatalf("status = %q, want open", concerns[0].Status)
	}

	openConcerns := 0
	for _, concern := range concerns {
		if concern.Status == "open" {
			openConcerns++
		}
	}
	if openConcerns != 1 {
		t.Fatalf("open concerns = %d, want 1", openConcerns)
	}
}

func TestConcernRaiseRequiresSummaryFlag(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"concern", "raise",
		"--on", "FID-task01",
		"--by", "developer",
		"--in-completion", "TC-000001",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want missing summary rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "error: --summary is required") {
		t.Fatalf("stderr = %q, want missing summary error", stderr)
	}
}

func TestConcernAddressUpdatesConcernAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Address CLI concern",
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

	_, raiseStderr := runFIDCommand(
		t,
		workspace,
		"concern", "raise",
		"--on", issueID,
		"--by", "developer",
		"--in-completion", firstCompletionID,
		"--summary", "Needs follow-up",
	)
	concernID := parseConcernID(t, raiseStderr)

	_, secondCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	addressCompletionID := parseCompletionID(t, secondCompleteStderr)

	_, addressStderr := runFIDCommand(
		t,
		workspace,
		"concern", "address", concernID,
		"--in-completion", addressCompletionID,
		"--by", "developer",
	)

	if !strings.Contains(addressStderr, "addressed concern "+concernID) {
		t.Fatalf("stderr = %q, want addressed concern confirmation", addressStderr)
	}
	if !strings.Contains(addressStderr, addressCompletionID) {
		t.Fatalf("stderr = %q, want addressing completion id", addressStderr)
	}

	concernsStdout, concernsStderr := runFIDCommand(t, workspace, "concern", "list", "--issue", issueID)
	if concernsStderr != "" {
		t.Fatalf("issue get concerns stderr = %q, want empty", concernsStderr)
	}

	var concerns []struct {
		ID                      string  `json:"id"`
		Status                  string  `json:"status"`
		UpdatedAt               string  `json:"updated_at"`
		RaisedInCompletionID    string  `json:"raised_in_completion_id"`
		AddressedInCompletionID *string `json:"addressed_in_completion_id"`
	}
	if err := json.Unmarshal([]byte(concernsStdout), &concerns); err != nil {
		t.Fatalf("Unmarshal(concerns) error = %v", err)
	}
	if len(concerns) != 1 {
		t.Fatalf("concerns len = %d, want 1", len(concerns))
	}
	if concerns[0].ID != concernID {
		t.Fatalf("concern id = %q, want %q", concerns[0].ID, concernID)
	}
	if concerns[0].Status != "addressed" {
		t.Fatalf("status = %q, want addressed", concerns[0].Status)
	}
	if concerns[0].AddressedInCompletionID == nil || *concerns[0].AddressedInCompletionID != addressCompletionID {
		t.Fatalf("addressed_in_completion_id = %#v, want %q", concerns[0].AddressedInCompletionID, addressCompletionID)
	}

	listStdout, listStderr := runFIDCommand(
		t,
		workspace,
		"concern", "list",
		"--issue", issueID,
	)
	if listStderr != "" {
		t.Fatalf("concern list stderr = %q, want empty", listStderr)
	}

	var listedConcerns []struct {
		ID                      string  `json:"id"`
		Status                  string  `json:"status"`
		AddressedInCompletionID *string `json:"addressed_in_completion_id"`
	}
	if err := json.Unmarshal([]byte(listStdout), &listedConcerns); err != nil {
		t.Fatalf("Unmarshal(concern list) error = %v", err)
	}
	if len(listedConcerns) != 1 {
		t.Fatalf("listed concerns len = %d, want 1", len(listedConcerns))
	}
	if listedConcerns[0].ID != concernID {
		t.Fatalf("listed concern id = %q, want %q", listedConcerns[0].ID, concernID)
	}
	if listedConcerns[0].Status != "addressed" {
		t.Fatalf("listed concern status = %q, want addressed", listedConcerns[0].Status)
	}
	if listedConcerns[0].AddressedInCompletionID == nil || *listedConcerns[0].AddressedInCompletionID != addressCompletionID {
		t.Fatalf("listed addressed_in_completion_id = %#v, want %q", listedConcerns[0].AddressedInCompletionID, addressCompletionID)
	}

	openConcerns := 0
	for _, concern := range listedConcerns {
		if concern.Status == "open" {
			openConcerns++
		}
	}
	if openConcerns != 0 {
		t.Fatalf("open concerns = %d, want 0", openConcerns)
	}
}

func TestConcernAddressRequiresCompletionFlag(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"concern", "address", "CON-test",
		"--by", "developer",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want missing completion rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "error: --in-completion is required") {
		t.Fatalf("stderr = %q, want missing completion error", stderr)
	}
}

func TestConcernListFiltersByIssueAndStatus(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	firstIssueID, firstConcernID := seedConcernListIssueAndConcern(t, workspace, "First concern", true)
	secondIssueID, secondConcernID := seedConcernListIssueAndConcern(t, workspace, "Second concern", false)

	tests := []struct {
		name    string
		args    []string
		wantIDs []string
	}{
		{
			name:    "all concerns",
			args:    []string{"concern", "list"},
			wantIDs: []string{secondConcernID, firstConcernID},
		},
		{
			name:    "filter by issue",
			args:    []string{"concern", "list", "--issue", firstIssueID},
			wantIDs: []string{firstConcernID},
		},
		{
			name:    "filter by status",
			args:    []string{"concern", "list", "--status", "open"},
			wantIDs: []string{secondConcernID},
		},
		{
			name:    "combined filters",
			args:    []string{"concern", "list", "--issue", secondIssueID, "--status", "open"},
			wantIDs: []string{secondConcernID},
		},
		{
			name:    "empty result",
			args:    []string{"concern", "list", "--issue", secondIssueID, "--status", "addressed"},
			wantIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, tt.args...)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			rows := decodeCommandJSONStructure(t, stdout).([]any)
			if gotIDs := commandConcernListIDs(rows); !equalIssueIDs(gotIDs, tt.wantIDs) {
				t.Fatalf("stdout ids = %#v, want %#v", gotIDs, tt.wantIDs)
			}

			if tt.name == "all concerns" {
				first := rows[0].(map[string]any)
				if got := first["issue_id"]; got != secondIssueID {
					t.Fatalf("first concern issue_id = %v, want %q", got, secondIssueID)
				}
				if got := first["status"]; got != "open" {
					t.Fatalf("first concern status = %v, want open", got)
				}
			}
		})
	}
}

func TestConcernListUsesDirectLegacyJSON(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, _ = seedConcernListIssueAndConcern(t, workspace, "First concern", true)
	secondIssueID, _ := seedConcernListIssueAndConcern(t, workspace, "Second concern", false)

	stdout, stderr := runFIDCommand(t, workspace, "concern", "list")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	want := decodeCommandJSONStructure(t, stdout)
	if got := want.([]any); len(got) != 2 {
		t.Fatalf("concern list length = %d, want 2", len(got))
	}

	emptyStdout, emptyStderr := runFIDCommand(t, workspace, "concern", "list", "--issue", secondIssueID, "--status", "addressed")
	if emptyStderr != "" {
		t.Fatalf("empty stderr = %q, want empty", emptyStderr)
	}
	if got := decodeCommandJSONStructure(t, emptyStdout); !commandStructuresEqual([]any{}, got) {
		t.Fatalf("empty JSON structure = %#v, want empty list", got)
	}
}

func parseCompletionID(t *testing.T, stderr string) string {
	t.Helper()

	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "completion: ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "completion: "))
		}
	}
	t.Fatalf("stderr = %q, want completion line", stderr)
	return ""
}

func parseConcernID(t *testing.T, stderr string) string {
	t.Helper()

	fields := strings.Fields(stderr)
	for index := range fields {
		if fields[index] == "concern" && index+1 < len(fields) {
			return strings.TrimSpace(fields[index+1])
		}
	}
	t.Fatalf("stderr = %q, want concern line", stderr)
	return ""
}

func seedConcernListIssueAndConcern(t *testing.T, workspace, summary string, address bool) (string, string) {
	t.Helper()

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", summary,
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

	_, raiseStderr := runFIDCommand(
		t,
		workspace,
		"concern", "raise",
		"--on", issueID,
		"--by", "developer",
		"--in-completion", firstCompletionID,
		"--summary", summary,
	)
	concernID := parseConcernID(t, raiseStderr)

	if !address {
		return issueID, concernID
	}

	_, secondCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		issueID,
		"--by", "developer",
	)
	addressCompletionID := parseCompletionID(t, secondCompleteStderr)

	_, _ = runFIDCommand(
		t,
		workspace,
		"concern", "address", concernID,
		"--in-completion", addressCompletionID,
		"--by", "developer",
	)

	return issueID, concernID
}

func commandConcernListIDs(rows []any) []string {
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		object, ok := row.(map[string]any)
		if !ok {
			continue
		}
		id, ok := object["id"].(string)
		if ok {
			ids = append(ids, id)
		}
	}
	return ids
}
