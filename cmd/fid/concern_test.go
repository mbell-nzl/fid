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
		"--format", "json",
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
			args:    []string{"concern", "list", "--format", "json"},
			wantIDs: []string{secondConcernID, firstConcernID},
		},
		{
			name:    "filter by issue",
			args:    []string{"concern", "list", "--issue", firstIssueID, "--format", "json"},
			wantIDs: []string{firstConcernID},
		},
		{
			name:    "filter by status",
			args:    []string{"concern", "list", "--status", "open", "--format", "json"},
			wantIDs: []string{secondConcernID},
		},
		{
			name:    "combined filters",
			args:    []string{"concern", "list", "--issue", secondIssueID, "--status", "open", "--format", "json"},
			wantIDs: []string{secondConcernID},
		},
		{
			name:    "empty result",
			args:    []string{"concern", "list", "--issue", secondIssueID, "--status", "addressed", "--format", "json"},
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

func TestConcernListDispatchesFormats(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, _ = seedConcernListIssueAndConcern(t, workspace, "First concern", true)
	secondIssueID, _ := seedConcernListIssueAndConcern(t, workspace, "Second concern", false)

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "concern", "list", "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}
	want := decodeCommandJSONStructure(t, jsonStdout)

	for _, formatName := range []string{"toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "concern", "list", "--format", formatName)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			var got any
			switch formatName {
			case "toon":
				got = decodeCommandTOONStructure(t, stdout)
			case "yaml":
				got = decodeCommandYAMLStructure(t, stdout)
			}

			if !commandStructuresEqual(want, got) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, got, want)
			}
		})
	}

	t.Run("markdown", func(t *testing.T) {
		stdout, stderr := runFIDCommand(t, workspace, "concern", "list", "--format", "markdown")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}

		rows := want.([]any)
		expected := concernListMarkdownExpectation(t, rows)
		if stdout != expected {
			t.Fatalf("markdown stdout = %q, want %q", stdout, expected)
		}
	})

	for _, formatName := range []string{"json", "toon", "yaml", "markdown"} {
		t.Run("empty_"+formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(
				t,
				workspace,
				"concern", "list",
				"--issue", secondIssueID,
				"--status", "addressed",
				"--format", formatName,
			)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			switch formatName {
			case "json":
				if got := decodeCommandJSONStructure(t, stdout); !commandStructuresEqual([]any{}, got) {
					t.Fatalf("empty json structure = %#v, want empty list", got)
				}
			case "toon":
				if got := decodeCommandTOONStructure(t, stdout); !commandStructuresEqual([]any{}, got) {
					t.Fatalf("empty toon structure = %#v, want empty list", got)
				}
			case "yaml":
				if got := decodeCommandYAMLStructure(t, stdout); !commandStructuresEqual([]any{}, got) {
					t.Fatalf("empty yaml structure = %#v, want empty list", got)
				}
			case "markdown":
				expected := concernListMarkdownExpectation(t, nil)
				if stdout != expected {
					t.Fatalf("empty markdown stdout = %q, want %q", stdout, expected)
				}
			}
		})
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

func concernListMarkdownExpectation(t *testing.T, rows []any) string {
	t.Helper()

	var builder strings.Builder
	builder.WriteString("# Concerns\n\n")
	if len(rows) == 0 {
		builder.WriteString("(no concerns)\n")
		return builder.String()
	}

	for _, row := range rows {
		item, ok := row.(map[string]any)
		if !ok {
			t.Fatalf("row = %#v, want object", row)
		}
		builder.WriteString("- ")
		builder.WriteString(item["id"].(string))
		builder.WriteString(" | issue ")
		builder.WriteString(item["issue_id"].(string))
		builder.WriteString(" | ")
		builder.WriteString(item["status"].(string))
		builder.WriteString(" | by ")
		builder.WriteString(item["raised_by_role"].(string))
		builder.WriteString(" | completion ")
		builder.WriteString(item["raised_in_completion_id"].(string))
		builder.WriteString("\n")
		builder.WriteString("  ")
		builder.WriteString(item["summary"].(string))
		builder.WriteString("\n")
		if bodyMDPath, ok := item["body_md_path"].(string); ok && bodyMDPath != "" {
			builder.WriteString("  Body: ")
			builder.WriteString(bodyMDPath)
			builder.WriteString("\n")
		}
	}
	return builder.String()
}
