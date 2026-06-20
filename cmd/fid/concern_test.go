package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConcernRaiseCreatesConcernAndPrintsConfirmation(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Raise CLI concern",
		"--by", "developer",
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

	concernsStdout, concernsStderr := runFIDCommand(
		t,
		workspace,
		"issue", "get", issueID,
		"--section", "concerns",
		"--format", "json",
	)
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

	issueStdout, issueStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if issueStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", issueStderr)
	}

	var issue struct {
		ConcernsRaised    bool `json:"concerns_raised"`
		OpenConcernsCount int  `json:"open_concerns_count"`
	}
	if err := json.Unmarshal([]byte(issueStdout), &issue); err != nil {
		t.Fatalf("Unmarshal(issue) error = %v", err)
	}
	if !issue.ConcernsRaised {
		t.Fatal("concerns_raised = false, want true")
	}
	if issue.OpenConcernsCount != 1 {
		t.Fatalf("open_concerns_count = %d, want 1", issue.OpenConcernsCount)
	}
}

func TestConcernRaiseRequiresSummaryFlag(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Address CLI concern",
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

	concernsStdout, concernsStderr := runFIDCommand(
		t,
		workspace,
		"issue", "get", issueID,
		"--section", "concerns",
		"--format", "json",
	)
	if concernsStderr != "" {
		t.Fatalf("issue get concerns stderr = %q, want empty", concernsStderr)
	}

	var concerns []struct {
		ID                   string `json:"id"`
		Status               string `json:"status"`
		UpdatedAt            string `json:"updated_at"`
		RaisedInCompletionID string `json:"raised_in_completion_id"`
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

	issueStdout, issueStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if issueStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", issueStderr)
	}

	var issue struct {
		ConcernsRaised    bool `json:"concerns_raised"`
		OpenConcernsCount int  `json:"open_concerns_count"`
	}
	if err := json.Unmarshal([]byte(issueStdout), &issue); err != nil {
		t.Fatalf("Unmarshal(issue) error = %v", err)
	}
	if issue.ConcernsRaised {
		t.Fatal("concerns_raised = true, want false")
	}
	if issue.OpenConcernsCount != 0 {
		t.Fatalf("open_concerns_count = %d, want 0", issue.OpenConcernsCount)
	}
}

func TestConcernAddressRequiresCompletionFlag(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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
