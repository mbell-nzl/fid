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
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want missing summary rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("error = %v, stderr = %q, want removed-leaf rejection", err, stderr)
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
		"concern", "address",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want missing completion rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("error = %v, stderr = %q, want removed-leaf rejection", err, stderr)
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
			args:    []string{"--json", "concern", "list", "--subject", "issue:" + firstIssueID},
			wantIDs: []string{firstConcernID},
		},
		{
			name:    "combined filters",
			args:    []string{"--json", "concern", "list", "--subject", "issue:" + secondIssueID, "--status", "open"},
			wantIDs: []string{secondConcernID},
		},
		{
			name:    "empty result",
			args:    []string{"--json", "concern", "list", "--subject", "issue:" + secondIssueID, "--status", "disposed"},
			wantIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, tt.args...)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			rows := decodeCommandJSONStructure(t, stdout).(map[string]any)["result"].(map[string]any)["items"].([]any)
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
	if got := commandConcernListIDs(items.([]any)); len(got) != 2 {
		t.Fatalf("generic items = %#v, want both generic Concerns", items)
	}
	human, stderr := runFIDCommand(t, workspace, "concern", "list")
	if stderr != "" || !strings.Contains(human, "CON-FIRSTCONCERN") || !strings.Contains(human, "CON-SECONDCONCERN") {
		t.Fatalf("generic Human list = %q / %q", human, stderr)
	}
}

func assertConcernMutationUnavailable(t *testing.T, workspace string, args []string) {
	t.Helper()
	journalPath := mustJournalPath(t, workspace)
	before := mustReadCommandFile(t, journalPath)
	result := executeCommandForTest(workspace, args[:2]...)
	if result.err == nil || result.stdout != "" || !strings.Contains(result.err.Error(), "unknown command") {
		t.Fatalf("removed concern mutation %v = %#v, want unknown-command rejection", args, result)
	}
	if after := mustReadCommandFile(t, journalPath); !bytes.Equal(after, before) {
		t.Fatalf("concern mutation %v changed journal bytes", args)
	}
}

func seedConcernListIssueAndConcern(t *testing.T, workspace, summary string, dispose bool) (string, string) {
	t.Helper()
	_, createStderr := runFIDCommand(t, workspace, "issue", "create",
		"--kind", "task", "--title", summary, "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)
	concernID := fmt.Sprintf("CON-%s", strings.ToUpper(strings.ReplaceAll(summary, " ", "")))
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open concern projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	status := "open"
	if dispose {
		status = "disposed"
	}
	if status == "open" {
		_, err = conn.Exec(`INSERT INTO concerns
			(id,subject_type,subject_id,actor_id,scope_json,severity,summary,reasoning,references_json,status,created_at,created_event_id)
			VALUES (?,'issue',?,'developer','{"paths":[]}','major',?,'Needs review','[]','open','2026-09-20T00:00:00Z',?)`,
			concernID, issueID, summary, "EVT-"+concernID)
	} else {
		_, err = conn.Exec(`INSERT INTO concerns
			(id,subject_type,subject_id,actor_id,scope_json,severity,summary,reasoning,references_json,status,
			 disposition_action,disposition_reason,disposed_actor_id,disposed_at,disposed_event_id,created_at,created_event_id)
			VALUES (?,'issue',?,'developer','{"paths":[]}','major',?,'Needs review','[]','disposed',
			 'rejected','Addressed','qa_lead','2026-09-21T00:00:00Z',?,'2026-09-20T00:00:00Z',?)`,
			concernID, issueID, summary, "EVT-dispose-"+concernID, "EVT-"+concernID)
	}
	if err != nil {
		t.Fatalf("seed generic Concern: %v", err)
	}
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
