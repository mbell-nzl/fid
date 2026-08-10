package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	internalpolicy "github.com/mbell-nzl/fid/internal/policy"
)

func TestPolicyShowDefaultsToJSONAndDispatchesFormats(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	defaultStdout, defaultStderr := runFIDCommand(t, workspace, "policy", "show")
	if defaultStderr != "" {
		t.Fatalf("default stderr = %q, want empty", defaultStderr)
	}

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "policy", "show", "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	if defaultStdout != jsonStdout {
		t.Fatalf("default stdout = %q, want json stdout %q", defaultStdout, jsonStdout)
	}

	want := decodeCommandJSONStructure(t, mustMarshalPolicyRowsJSON(t, internalpolicy.MatrixRows()))
	got := decodeCommandJSONStructure(t, jsonStdout)
	if !commandStructuresEqual(want, got) {
		t.Fatalf("json structure = %#v, want %#v", got, want)
	}

	if gotActions := commandPolicyActions(got.([]any)); !equalIssueIDs(gotActions, expectedPolicyActions()) {
		t.Fatalf("policy actions = %#v, want %#v", gotActions, expectedPolicyActions())
	}

	for _, formatName := range []string{"toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "policy", "show", "--format", formatName)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			var structured any
			switch formatName {
			case "toon":
				structured = decodeCommandTOONStructure(t, stdout)
			case "yaml":
				structured = decodeCommandYAMLStructure(t, stdout)
			}

			if !commandStructuresEqual(want, structured) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, structured, want)
			}
		})
	}

	markdownStdout, markdownStderr := runFIDCommand(t, workspace, "policy", "show", "--format", "markdown")
	if markdownStderr != "" {
		t.Fatalf("markdown stderr = %q, want empty", markdownStderr)
	}

	wantMarkdown := expectedPolicyMarkdown(internalpolicy.MatrixRows())
	if markdownStdout != wantMarkdown {
		t.Fatalf("markdown stdout = %q, want %q", markdownStdout, wantMarkdown)
	}
}

func TestPolicyShowIsReadOnlyAndDoesNotAcceptByFlag(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	journalPath := mustJournalPath(t, workspace)
	beforePayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}
	beforeLines := splitJournalLines(beforePayload)

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	activityBefore := countActivityLogRows(t, conn)

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show", "--format", "json")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if strings.TrimSpace(stdout) == "" {
		t.Fatal("stdout = empty, want policy matrix output")
	}

	afterPayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after show error = %v", journalPath, err)
	}
	afterLines := splitJournalLines(afterPayload)
	if len(afterLines) != len(beforeLines) {
		t.Fatalf("journal lines after policy show = %d, want unchanged %d", len(afterLines), len(beforeLines))
	}

	activityAfter := countActivityLogRows(t, conn)
	if activityAfter != activityBefore {
		t.Fatalf("activity_log rows after policy show = %d, want unchanged %d", activityAfter, activityBefore)
	}

	stdout, stderr, err = runFIDCommandExpectError(t, workspace, "policy", "show", "--by", "developer")
	if err == nil {
		t.Fatal("Execute(--by) error = nil, want unknown flag rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout with --by = %q, want empty", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr with --by = %q, want empty", stderr)
	}
	if !strings.Contains(err.Error(), "unknown flag: --by") {
		t.Fatalf("Execute(--by) error = %q, want unknown flag rejection", err.Error())
	}
}

func TestPolicyShowMatchesWriteEnforcementForPermittedAndDeniedCells(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show", "--format", "json")
	if stderr != "" {
		t.Fatalf("policy show stderr = %q, want empty", stderr)
	}

	rows := decodeCommandJSONStructure(t, stdout).([]any)
	if got := commandPolicyCellAllowed(t, rows, "issue.create.epic", "developer"); got {
		t.Fatal("policy show developer/issue.create.epic = true, want false")
	}
	if got := commandPolicyCellAllowed(t, rows, "issue.create.epic", "tl"); !got {
		t.Fatal("policy show tl/issue.create.epic = false, want true")
	}
	for _, action := range []string{"workflow.define", "workflow.assign"} {
		if got := commandPolicyCellAllowed(t, rows, action, "tl"); !got {
			t.Fatalf("policy show tl/%s = false, want true", action)
		}
		if got := commandPolicyCellAllowed(t, rows, action, "human"); !got {
			t.Fatalf("policy show human/%s = false, want true", action)
		}
		if got := commandPolicyCellAllowed(t, rows, action, "developer"); got {
			t.Fatalf("policy show developer/%s = true, want false", action)
		}
	}

	deniedStdout, deniedStderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"issue", "create",
		"--kind", "epic",
		"--title", "Denied parity epic",
		"--scheduling", "sequential",
		"--by", "developer",
	)
	if err == nil {
		t.Fatal("developer epic create error = nil, want permission denial")
	}
	if deniedStdout != "" {
		t.Fatalf("developer epic create stdout = %q, want empty", deniedStdout)
	}
	if !strings.Contains(deniedStderr, `error: role "developer" not permitted for action "issue.create.epic"`) {
		t.Fatalf("developer epic create stderr = %q, want permission denial", deniedStderr)
	}

	allowedStdout, allowedStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "epic",
		"--title", "Permitted parity epic",
		"--scheduling", "sequential",
		"--by", "tl",
	)
	if allowedStdout != "" {
		t.Fatalf("tl epic create stdout = %q, want empty", allowedStdout)
	}
	if !strings.Contains(allowedStderr, "created issue ") || !strings.Contains(allowedStderr, "(epic)") {
		t.Fatalf("tl epic create stderr = %q, want create confirmation", allowedStderr)
	}
}

func TestPolicyShowIncludesWorkflowAuthoringRowsInEveryFormat(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	for _, formatName := range []string{"json", "toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "policy", "show", "--format", formatName)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}
			var rows []any
			switch formatName {
			case "json":
				rows = decodeCommandJSONStructure(t, stdout).([]any)
			case "toon":
				rows = decodeCommandTOONStructure(t, stdout).([]any)
			case "yaml":
				rows = decodeCommandYAMLStructure(t, stdout).([]any)
			}
			for _, action := range []string{"workflow.define", "workflow.assign"} {
				if !commandPolicyCellAllowed(t, rows, action, "tl") || !commandPolicyCellAllowed(t, rows, action, "human") {
					t.Fatalf("policy show %s missing allowed %s row", formatName, action)
				}
			}
		})
	}

	markdown, stderr := runFIDCommand(t, workspace, "policy", "show", "--format", "markdown")
	if stderr != "" {
		t.Fatalf("markdown stderr = %q, want empty", stderr)
	}
	for _, action := range []string{"workflow.define", "workflow.assign"} {
		if !strings.Contains(markdown, "| "+action+" | false | false | false | false | true | false | true |") {
			t.Fatalf("markdown = %q, want %s row", markdown, action)
		}
	}
}

func mustMarshalPolicyRowsJSON(t *testing.T, rows []internalpolicy.MatrixRow) string {
	t.Helper()

	payload, err := json.Marshal(rows)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	return string(payload)
}

func expectedPolicyActions() []string {
	rows := internalpolicy.MatrixRows()
	actions := make([]string, 0, len(rows))
	for _, row := range rows {
		actions = append(actions, row.Action)
	}
	return actions
}

func commandPolicyActions(rows []any) []string {
	actions := make([]string, 0, len(rows))
	for _, row := range rows {
		fields, ok := row.(map[string]any)
		if !ok {
			continue
		}
		action, _ := fields["action"].(string)
		actions = append(actions, action)
	}
	return actions
}

func expectedPolicyMarkdown(rows []internalpolicy.MatrixRow) string {
	var builder strings.Builder

	builder.WriteString("# Policy\n\n")
	builder.WriteString("| action | developer | qa_reviewer | qa_lead | tl_lite | tl | workflow_runner | human |\n")
	builder.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- |\n")

	for _, row := range rows {
		_, _ = fmt.Fprintf(&builder,
			"| %s | %t | %t | %t | %t | %t | %t | %t |\n",
			row.Action,
			row.Developer,
			row.QAReviewer,
			row.QALead,
			row.TLLite,
			row.TL,
			row.WorkflowRunner,
			row.Human,
		)
	}

	return builder.String()
}

func countActivityLogRows(t *testing.T, conn *sql.DB) int {
	t.Helper()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&count); err != nil {
		t.Fatalf("count activity_log rows: %v", err)
	}
	return count
}

func commandPolicyCellAllowed(t *testing.T, rows []any, action, role string) bool {
	t.Helper()

	for _, row := range rows {
		fields, ok := row.(map[string]any)
		if !ok {
			continue
		}
		if fields["action"] != action {
			continue
		}
		value, ok := fields[role].(bool)
		if !ok {
			t.Fatalf("policy row %q missing boolean role %q: %#v", action, role, fields)
		}
		return value
	}

	t.Fatalf("policy row %q not found", action)
	return false
}
