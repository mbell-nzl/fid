package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestConcernRaiseCreatesConcernAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	assertConcernMutationUnavailable(t, workspace, []string{"concern", "raise", "--on", "FID-task01", "--actor", "agent:developer", "--in-completion", "TC-000001", "--summary", "Potential race condition"})
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
		"--actor", "agent:developer",
		"--in-completion", "TC-000001",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want missing summary rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "capability is unavailable") {
		t.Fatalf("stderr = %q, want typed capability-unavailable error", stderr)
	}
}

func TestConcernAddressUpdatesConcernAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	assertConcernMutationUnavailable(t, workspace, []string{"concern", "address", "CON-test", "--actor", "human:reviewer", "--in-completion", "TC-000001"})
}

func TestConcernAddressRequiresCompletionFlag(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"concern", "address", "CON-test",
		"--actor", "agent:developer",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want missing completion rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "capability is unavailable") {
		t.Fatalf("stderr = %q, want typed capability-unavailable error", stderr)
	}
}

func TestConcernListRetainsIssueScopedV1ReadUntilRemoval(t *testing.T) {
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
			name:    "filter by issue",
			args:    []string{"concern", "list", "--issue", firstIssueID},
			wantIDs: []string{firstConcernID},
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

		})
	}
}

func TestConcernListGenericReadExcludesV1Rows(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, _ = seedConcernListIssueAndConcern(t, workspace, "First concern", true)
	_, _ = seedConcernListIssueAndConcern(t, workspace, "Second concern", false)

	stdout, stderr := runFIDCommand(t, workspace, "--json", "concern", "list")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	result := decodeCommandJSONStructure(t, stdout).(map[string]any)
	if got := result["operation"]; got != "concern.list" {
		t.Fatalf("operation = %v, want concern.list", got)
	}
	items := result["result"].(map[string]any)["items"]
	if !commandStructuresEqual([]any{}, items) {
		t.Fatalf("generic items = %#v, want no V1 rows", items)
	}
	human, stderr := runFIDCommand(t, workspace, "concern", "list")
	if stderr != "" || !strings.Contains(human, "Purpose: Inspect a raised concern") || !strings.Contains(human, "None.") {
		t.Fatalf("generic Human list = %q / %q", human, stderr)
	}
}

func assertConcernMutationUnavailable(t *testing.T, workspace string, args []string) {
	t.Helper()
	journalPath := mustJournalPath(t, workspace)
	before := mustReadCommandFile(t, journalPath)
	result := executeCommandForTest(workspace, args...)
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, "capability is unavailable") {
		t.Fatalf("concern mutation %v = %#v, want typed capability-unavailable rejection", args, result)
	}
	if after := mustReadCommandFile(t, journalPath); !bytes.Equal(after, before) {
		t.Fatalf("concern mutation %v changed journal bytes", args)
	}
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

	firstCompletionID := fmt.Sprintf("TC-%s-raised", strings.ToUpper(strings.ReplaceAll(summary, " ", "")))
	seedFrozenTaskCompletion(t, workspace, firstCompletionID, issueID)

	concernID := fmt.Sprintf("CON-%s", strings.ToUpper(strings.ReplaceAll(summary, " ", "")))
	addressCompletionID := ""
	if address {
		addressCompletionID = fmt.Sprintf("TC-%s-addressed", strings.ToUpper(strings.ReplaceAll(summary, " ", "")))
		seedFrozenTaskCompletion(t, workspace, addressCompletionID, issueID)
	}
	status := "open"
	var addressed any
	if address {
		status = "addressed"
		addressed = addressCompletionID
	}
	seedConcernProjection(t, workspace, concernID, issueID, firstCompletionID, summary, status, addressed)
	return issueID, concernID
}

func seedFrozenTaskCompletion(t *testing.T, workspace, completionID, issueID string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open frozen completion projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	if _, err := conn.Exec(`INSERT INTO task_completions (id, issue_id, attempt_number, developer_role, status, completed_at) VALUES (?, ?, ?, ?, ?, datetime('now'))`, completionID, issueID, 1, "developer", "dev_complete"); err != nil {
		t.Fatalf("seed frozen completion projection: %v", err)
	}
}

func seedConcernProjection(t *testing.T, workspace, concernID, issueID, raisedCompletionID, summary, status string, addressedCompletionID any) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open concern projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	if _, err := conn.Exec(`
		INSERT INTO concerns (id, issue_id, raised_in_completion_id, raised_by_role, summary, status, addressed_in_completion_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`, concernID, issueID, raisedCompletionID, "developer", summary, status, addressedCompletionID); err != nil {
		t.Fatalf("seed concern projection: %v", err)
	}
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
