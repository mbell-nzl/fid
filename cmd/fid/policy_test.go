package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	internalpolicy "github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

func TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	if flag := newPolicyShowCommand().Flags().Lookup("format"); flag != nil {
		t.Fatalf("policy show --format flag = %#v, want absent", flag)
	}

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if !strings.HasSuffix(stdout, "\n") {
		t.Fatalf("stdout = %q, want one trailing newline", stdout)
	}

	rows := decodeCommandJSONStructure(t, stdout).([]any)
	if gotActions := commandPolicyActions(rows); !equalIssueIDs(gotActions, expectedPolicyActions()) {
		t.Fatalf("policy actions = %#v, want %#v", gotActions, expectedPolicyActions())
	}
	if _, enveloped := decodeCommandJSONStructure(t, stdout).(map[string]any); enveloped {
		t.Fatalf("policy show JSON = %s, want unenveloped legacy rows", stdout)
	}
}

func TestRemainingLegacyJSONReadsDoNotDeclareFormatFlags(t *testing.T) {
	t.Parallel()

	commands := map[string]*cobra.Command{
		"activity list":   newActivityListCommand(),
		"claim list":      newClaimListCommand(),
		"concern list":    newConcernListCommand(),
		"journal inspect": newJournalInspectCommand(),
		"sql":             newSQLCommand(),
		"stats cost":      newStatsCostCommand(),
		"stats tokens":    newStatsTokensCommand(),
	}
	for name, command := range commands {
		if flag := command.Flags().Lookup("format"); flag != nil {
			t.Fatalf("%s --format flag = %#v, want absent", name, flag)
		}
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

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show")
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

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show")
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

func TestPolicyShowIncludesWorkflowAuthoringRows(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	rows := decodeCommandJSONStructure(t, stdout).([]any)
	for _, action := range []string{"workflow.define", "workflow.assign"} {
		if !commandPolicyCellAllowed(t, rows, action, "tl") || !commandPolicyCellAllowed(t, rows, action, "human") {
			t.Fatalf("policy show missing allowed %s row", action)
		}
	}
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
