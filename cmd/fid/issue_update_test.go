package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIssueUpdateChangesFieldsAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Implement issue update", "--priority", "2", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	_, updateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "update", issueID,
		"--title", "Ship issue update",
		"--priority", "1",
		"--by", "developer",
	)
	gotLines := strings.Split(strings.TrimSpace(updateStderr), "\n")
	wantLines := []string{
		"updated " + issueID + ":",
		"priority 2 -> 1",
		`title "Implement issue update" -> "Ship issue update"`,
	}
	if len(gotLines) != len(wantLines) {
		t.Fatalf("stderr lines = %#v, want %#v", gotLines, wantLines)
	}
	for index, wantLine := range wantLines {
		if gotLines[index] != wantLine {
			t.Fatalf("stderr lines = %#v, want %#v", gotLines, wantLines)
		}
	}

	stdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

	var got struct {
		Title    string `json:"title"`
		Priority *int   `json:"priority"`
	}
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Title != "Ship issue update" {
		t.Fatalf("title = %q, want Ship issue update", got.Title)
	}
	if got.Priority == nil || *got.Priority != 1 {
		t.Fatalf("priority = %#v, want 1", got.Priority)
	}
}

func TestIssueUpdateSingleFieldPrintsExactlyOneChangedFieldLine(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Single field update", "--priority", "2", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	_, updateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "update", issueID,
		"--priority", "1",
		"--by", "developer",
	)

	gotLines := strings.Split(strings.TrimSpace(updateStderr), "\n")
	wantLines := []string{
		"updated " + issueID + ":",
		"priority 2 -> 1",
	}
	if len(gotLines) != len(wantLines) {
		t.Fatalf("stderr lines = %#v, want %#v", gotLines, wantLines)
	}
	for index, wantLine := range wantLines {
		if gotLines[index] != wantLine {
			t.Fatalf("stderr lines = %#v, want %#v", gotLines, wantLines)
		}
	}
}

func TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Reject invalid update", "--priority", "2", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	beforeStdout, beforeStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if beforeStderr != "" {
		t.Fatalf("issue get before stderr = %q, want empty", beforeStderr)
	}

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "update", issueID, "--priority", "8", "--by", "developer")
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid update rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "error: update issue: priority must be between 0 and 5") {
		t.Fatalf("stderr = %q, want invalid priority error", stderr)
	}

	afterStdout, afterStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if afterStderr != "" {
		t.Fatalf("issue get after stderr = %q, want empty", afterStderr)
	}
	if afterStdout != beforeStdout {
		t.Fatalf("issue JSON after rejected update = %q, want %q", afterStdout, beforeStdout)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore {
		t.Fatalf("journal lines after rejected update = %d, want %d", linesAfter, linesBefore)
	}
}

func TestIssueUpdateRejectsStaleExpectedRevision(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createOutput := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Concurrency target", "--by", "developer")
	issueID := parseCreatedIssueID(t, createOutput)
	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	result := executeCommandForTest(workspace,
		"--json", "issue", "update", issueID,
		"--title", "Rejected stale update",
		"--expected-revision", "2",
		"--by", "developer",
	)
	if result.err == nil {
		t.Fatal("Execute() error = nil, want stale revision rejection")
	}
	if result.stdout != "" {
		t.Fatalf("stdout = %q, want empty", result.stdout)
	}
	var failure struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Error     struct {
			Code    string `json:"code"`
			Details struct {
				Expected int `json:"expected_revision"`
				Actual   int `json:"actual_revision"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil {
		t.Fatalf("Unmarshal(stderr) error = %v; stderr = %q", err, result.stderr)
	}
	if failure.OK || failure.Operation != "issue.update" || failure.Error.Code != "state-conflict" {
		t.Fatalf("failure = %+v, want issue.update state-conflict", failure)
	}
	if failure.Error.Details.Expected != 2 || failure.Error.Details.Actual != 1 {
		t.Fatalf("state-conflict details = %+v, want expected 2 actual 1", failure.Error.Details)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
		t.Fatalf("journal lines after stale update = %d, want unchanged %d", got, linesBefore)
	}
}

func TestIssueUpdateJSONOrdersChangesByField(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	create := executeCommandForTest(workspace, "--json", "issue", "create", "--kind", "task", "--title", "Ordered changes", "--by", "developer")
	if create.err != nil {
		t.Fatalf("issue create error = %v, stderr = %q", create.err, create.stderr)
	}
	var created struct {
		Result struct {
			Issue struct {
				ID string `json:"id"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(create.stdout), &created); err != nil {
		t.Fatalf("Unmarshal(create stdout) error = %v", err)
	}

	updated := executeCommandForTest(workspace,
		"--json", "issue", "update", created.Result.Issue.ID,
		"--title", "Ordered changes updated",
		"--priority", "2",
		"--spec-md", "project-docs/specs/02_TRD.md",
		"--by", "developer",
	)
	if updated.err != nil {
		t.Fatalf("issue update error = %v, stderr = %q", updated.err, updated.stderr)
	}
	var envelope struct {
		Result struct {
			Changes []struct {
				Field string `json:"field"`
			} `json:"changes"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(updated.stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(update stdout) error = %v", err)
	}
	got := make([]string, 0, len(envelope.Result.Changes))
	for _, change := range envelope.Result.Changes {
		got = append(got, change.Field)
	}
	want := []string{"priority", "spec_md_path", "title"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("change fields = %#v, want %#v", got, want)
	}
}

func TestIssueCreateRejectsRemovedOwnerFlagsWithJSONFailure(t *testing.T) {
	for _, flag := range []string{"owner-role", "owner-actor"} {
		t.Run(flag, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)

			result := executeCommandForTest(workspace,
				"--json", "issue", "create",
				"--kind", "task", "--title", "Removed owner input",
				"--"+flag, "qa_lead", "--by", "developer",
			)
			if result.err == nil {
				t.Fatalf("issue create error = nil, want unknown %s rejection", flag)
			}
			if result.stdout != "" {
				t.Fatalf("stdout = %q, want empty", result.stdout)
			}
			var failure struct {
				Operation string `json:"operation"`
				Error     struct {
					Code string `json:"code"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil {
				t.Fatalf("Unmarshal(stderr) error = %v; stderr = %q", err, result.stderr)
			}
			if failure.Operation != "issue.create" || failure.Error.Code != "invalid-argument" {
				t.Fatalf("failure = %+v, want issue.create invalid-argument", failure)
			}
		})
	}
}

func TestIssueMutationHumanSuccessUsesStdoutOnly(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	created := executeCommandForTest(workspace, "issue", "create", "--kind", "task", "--title", "Human mutation channels", "--by", "developer")
	if created.err != nil || created.stderr != "" || !strings.Contains(created.stdout, "created issue ") {
		t.Fatalf("create = %#v, want stdout-only success", created)
	}
	issueID := parseCreatedIssueID(t, created.stdout)

	updated := executeCommandForTest(workspace, "issue", "update", issueID, "--priority", "2", "--by", "developer")
	if updated.err != nil || updated.stderr != "" || !strings.Contains(updated.stdout, "updated "+issueID) {
		t.Fatalf("update = %#v, want stdout-only success", updated)
	}

	claimed := executeCommandForTest(workspace, "issue", "claim", issueID, "--by", "developer")
	if claimed.err != nil {
		t.Fatalf("issue claim error = %v, stderr = %q", claimed.err, claimed.stderr)
	}
	transitioned := executeCommandForTest(workspace, "issue", "transition", issueID, "--to", "in_progress", "--by", "developer")
	if transitioned.err != nil || transitioned.stderr != "" || !strings.Contains(transitioned.stdout, "transitioned "+issueID) {
		t.Fatalf("transition = %#v, want stdout-only success", transitioned)
	}
}

func TestIssueUpdateRequiresKnownByRole(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Role validation update", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	t.Run("missing by", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "update", issueID, "--title", "No role")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --by rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: --by is required") {
			t.Fatalf("stderr = %q, want missing --by error", stderr)
		}
	})

	t.Run("unknown role", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "update", issueID, "--title", "Unknown role", "--by", "intern")
		if err == nil {
			t.Fatal("Execute() error = nil, want unknown role rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
			t.Fatalf("stderr = %q, want unknown role error", stderr)
		}
	})
}
