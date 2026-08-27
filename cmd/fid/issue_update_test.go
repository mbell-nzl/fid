package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestIssueUpdateChangesFieldsAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Implement issue update", "--priority", "2", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	_, updateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "update", issueID,
		"--title", "Ship issue update",
		"--priority", "1",
		"--actor", "agent:developer",
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

	stdout := mustGetIssueJSON(t, workspace, issueID)

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

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Single field update", "--priority", "2", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	_, updateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "update", issueID,
		"--priority", "1",
		"--actor", "agent:developer",
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

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Reject invalid update", "--priority", "2", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	beforeStdout := mustGetIssueJSON(t, workspace, issueID)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "update", issueID, "--priority", "8", "--actor", "agent:developer")
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid update rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "error: update issue: priority must be between 0 and 5") {
		t.Fatalf("stderr = %q, want invalid priority error", stderr)
	}

	afterStdout := mustGetIssueJSON(t, workspace, issueID)
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

	_, createOutput := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Concurrency target", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createOutput)
	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	result := executeCommandForTest(workspace,
		"--json", "issue", "update", issueID,
		"--title", "Rejected stale update",
		"--expected-revision", "2",
		"--actor", "agent:developer",
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

	create := executeCommandForTest(workspace, "--json", "issue", "create", "--kind", "task", "--title", "Ordered changes", "--actor", "agent:developer")
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
		"--actor", "agent:developer",
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
				"--"+flag, "qa_lead", "--actor", "agent:developer",
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

func TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures(t *testing.T) {
	tests := []struct {
		operation string
		flag      string
	}{
		{operation: "create", flag: "owner-role"},
		{operation: "create", flag: "owner-actor"},
		{operation: "update", flag: "owner-role"},
		{operation: "update", flag: "owner-actor"},
	}

	for _, tt := range tests {
		t.Run(tt.operation+"/"+tt.flag, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)
			issueID := ""
			if tt.operation == "update" {
				created := executeCommandForTest(workspace, "--json", "issue", "create", "--kind", "task", "--title", "Owner flag target", "--actor", "agent:developer")
				if created.err != nil || created.stderr != "" {
					t.Fatalf("issue create = %#v, want stdout-only success", created)
				}
				var envelope struct {
					Result struct {
						Issue struct {
							ID string `json:"id"`
						} `json:"issue"`
					} `json:"result"`
				}
				if err := json.Unmarshal([]byte(created.stdout), &envelope); err != nil {
					t.Fatalf("Unmarshal(create stdout) error = %v", err)
				}
				issueID = envelope.Result.Issue.ID
			}

			journalPath := mustJournalPath(t, workspace)
			linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			args := []string{"--json", "issue", tt.operation}
			if tt.operation == "create" {
				args = append(args, "--kind", "task", "--title", "Removed owner input")
			} else {
				args = append(args, issueID, "--title", "Should not update")
			}
			args = append(args, "--"+tt.flag, "qa_lead", "--by", "developer")
			result := executeCommandForTest(workspace, args...)

			if result.err == nil {
				t.Fatalf("issue %s error = nil, want unknown %s rejection", tt.operation, tt.flag)
			}
			if result.stdout != "" {
				t.Fatalf("stdout = %q, want empty", result.stdout)
			}
			var failure struct {
				Version   string `json:"contract_version"`
				OK        bool   `json:"ok"`
				Operation string `json:"operation"`
				Error     struct {
					Code      string `json:"code"`
					Category  string `json:"category"`
					Retryable bool   `json:"retryable"`
					Message   string `json:"message"`
					Subject   any    `json:"subject"`
					Details   struct {
						Argument string   `json:"argument"`
						Value    any      `json:"value"`
						Allowed  []string `json:"allowed"`
					} `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil {
				t.Fatalf("Unmarshal(stderr) error = %v; stderr = %q", err, result.stderr)
			}
			if failure.Version != "fid-cli.v1" || failure.OK || failure.Operation != "issue."+tt.operation || failure.Error.Code != "invalid-argument" || failure.Error.Category != "invalid-input" || failure.Error.Retryable || failure.Error.Message != "invalid Issue mutation request" || failure.Error.Subject != nil {
				t.Fatalf("failure = %+v, want canonical invalid-argument envelope", failure)
			}
			if failure.Error.Details.Argument != "request" || failure.Error.Details.Value != nil || strings.Join(failure.Error.Details.Allowed, ",") != "accepted issue."+tt.operation+" input" {
				t.Fatalf("failure details = %+v, want request/null/accepted operation input", failure.Error.Details)
			}
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
				t.Fatalf("journal lines after rejected %s = %d, want %d", tt.operation, got, linesBefore)
			}
		})
	}
}

func TestIssueUpdateCanonicalChangesExcludeOwnerAndPreserveProjection(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	created := executeCommandForTest(workspace, "--json", "issue", "create", "--kind", "task", "--title", "Owner projection target", "--actor", "agent:developer")
	if created.err != nil || created.stderr != "" {
		t.Fatalf("issue create = %#v, want stdout-only success", created)
	}
	var createEnvelope struct {
		Result struct {
			Issue struct {
				ID string `json:"id"`
			} `json:"issue"`
		} `json:"result"`
		Event struct {
			ID         string `json:"id"`
			OccurredAt string `json:"occurred_at"`
		} `json:"event"`
	}
	if err := json.Unmarshal([]byte(created.stdout), &createEnvelope); err != nil {
		t.Fatalf("Unmarshal(create stdout) error = %v", err)
	}
	issueID := createEnvelope.Result.Issue.ID
	journalPath := mustJournalPath(t, workspace)
	createTail := lastCommandJournalEvent(t, journalPath)
	if createEnvelope.Event.ID != createTail.EventID || createEnvelope.Event.OccurredAt != createTail.OccurredAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("create event = %+v, want journal tail %q at %s", createEnvelope.Event, createTail.EventID, createTail.OccurredAt)
	}
	runFIDCommand(t, workspace, "label", "add", "api", "--by", "developer")

	beforeRole, beforeActor := commandIssueOwners(t, workspace, issueID)
	updated := executeCommandForTest(workspace,
		"--json", "issue", "update", issueID,
		"--title", "Owner projection updated",
		"--priority", "2",
		"--labels", "api",
		"--spec-md", "project-docs/specs/02_TRD.md",
		"--actor", "agent:developer",
	)
	if updated.err != nil || updated.stderr != "" {
		t.Fatalf("issue update = %#v, want stdout-only success", updated)
	}
	var updateEnvelope struct {
		Result struct {
			Changes []struct {
				Field string `json:"field"`
			} `json:"changes"`
		} `json:"result"`
		Event struct {
			ID         string `json:"id"`
			OccurredAt string `json:"occurred_at"`
		} `json:"event"`
	}
	if err := json.Unmarshal([]byte(updated.stdout), &updateEnvelope); err != nil {
		t.Fatalf("Unmarshal(update stdout) error = %v", err)
	}
	fields := make([]string, 0, len(updateEnvelope.Result.Changes))
	for _, change := range updateEnvelope.Result.Changes {
		if strings.Contains(change.Field, "owner") {
			t.Fatalf("update changes = %#v, want no owner member", updateEnvelope.Result.Changes)
		}
		fields = append(fields, change.Field)
	}
	if got, want := strings.Join(fields, ","), "labels,priority,spec_md_path,title"; got != want {
		t.Fatalf("update change fields = %q, want %q", got, want)
	}
	updateTail := lastCommandJournalEvent(t, journalPath)
	if updateEnvelope.Event.ID != updateTail.EventID || updateEnvelope.Event.OccurredAt != updateTail.OccurredAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("update event = %+v, want journal tail %q at %s", updateEnvelope.Event, updateTail.EventID, updateTail.OccurredAt)
	}
	afterRole, afterActor := commandIssueOwners(t, workspace, issueID)
	if !sameCommandOptionalString(afterRole, beforeRole) || !sameCommandOptionalString(afterActor, beforeActor) {
		t.Fatalf("projected owners after update = (%#v, %#v), want unchanged (%#v, %#v)", afterRole, afterActor, beforeRole, beforeActor)
	}
}

func commandIssueOwners(t *testing.T, workspace, issueID string) (*string, *string) {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projected database: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("close projected database: %v", err)
		}
	}()
	var role, actor sql.NullString
	if err := conn.QueryRow(`SELECT owner_role, owner_actor FROM issues WHERE id = ?`, issueID).Scan(&role, &actor); err != nil {
		t.Fatalf("query projected owners: %v", err)
	}
	return commandOptionalString(role), commandOptionalString(actor)
}

func commandOptionalString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func sameCommandOptionalString(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func TestIssueMutationHumanSuccessUsesStdoutOnly(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	created := executeCommandForTest(workspace, "issue", "create", "--kind", "task", "--title", "Human mutation channels", "--actor", "agent:developer")
	if created.err != nil || created.stderr != "" || !strings.Contains(created.stdout, "created issue ") {
		t.Fatalf("create = %#v, want stdout-only success", created)
	}
	issueID := parseCreatedIssueID(t, created.stdout)

	updated := executeCommandForTest(workspace, "issue", "update", issueID, "--priority", "2", "--actor", "agent:developer")
	if updated.err != nil || updated.stderr != "" || !strings.Contains(updated.stdout, "updated "+issueID) {
		t.Fatalf("update = %#v, want stdout-only success", updated)
	}

	if _, _, err, handled := runLegacyClaimFixture(workspace, []string{"issue", "claim", issueID, "--by", "developer"}); !handled || err != nil {
		t.Fatalf("fixture claim handled=%t error=%v", handled, err)
	}
	transitioned := executeCommandForTest(workspace, "issue", "transition", issueID, "--to", "in_progress", "--actor", "agent:developer")
	if transitioned.err != nil || transitioned.stderr != "" || !strings.Contains(transitioned.stdout, "transitioned "+issueID) {
		t.Fatalf("transition = %#v, want stdout-only success", transitioned)
	}
}

func TestIssueUpdateUsesStandaloneActorAttribution(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Role validation update", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	t.Run("missing actor", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "update", issueID, "--title", "No role")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --actor rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: --actor is required") {
			t.Fatalf("stderr = %q, want missing --actor error", stderr)
		}
	})

	t.Run("unregistered actor updates as attribution", func(t *testing.T) {
		_, stderr := runFIDCommand(t, workspace, "issue", "update", issueID, "--title", "Unregistered actor", "--actor", "agent:intern")
		if !strings.Contains(stderr, "updated "+issueID) {
			t.Fatalf("stderr = %q, want update confirmation", stderr)
		}
	})
}
