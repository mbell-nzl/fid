package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/spf13/cobra"
)

func TestPublicPolicyCommandIsAbsent(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	result := executeCommandForTest(workspace, "policy", "show")
	if result.err == nil || result.stdout != "" || !strings.Contains(result.err.Error(), `unknown command "policy"`) {
		t.Fatalf("policy show = %#v, want absent public command", result)
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

	result := executeCommandForTest(workspace, "policy", "show")
	if result.err == nil || result.stdout != "" || !strings.Contains(result.err.Error(), `unknown command "policy"`) {
		t.Fatalf("policy show = %#v, want absent public command", result)
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

	result = executeCommandForTest(workspace, "policy", "show", "--by", "developer")
	if result.err == nil || result.stdout != "" || !strings.Contains(result.err.Error(), `unknown command "policy"`) {
		t.Fatalf("policy show with V1 flag = %#v, want absent public command", result)
	}
}

func TestPolicyShowDoesNotGrantStandaloneIssueAuthority(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	attributedStdout, attributedStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "epic",
		"--title", "Denied parity epic",
		"--scheduling", "sequential",
		"--actor", "agent:developer",
	)
	if attributedStdout != "" {
		t.Fatalf("developer epic create stdout = %q, want empty", attributedStdout)
	}
	if !strings.Contains(attributedStderr, "created issue ") || !strings.Contains(attributedStderr, "(epic)") {
		t.Fatalf("developer epic create stderr = %q, want standalone attribution confirmation", attributedStderr)
	}

	allowedStdout, allowedStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "epic",
		"--title", "Permitted parity epic",
		"--scheduling", "sequential",
		"--actor", "agent:tl",
	)
	if allowedStdout != "" {
		t.Fatalf("tl epic create stdout = %q, want empty", allowedStdout)
	}
	if !strings.Contains(allowedStderr, "created issue ") || !strings.Contains(allowedStderr, "(epic)") {
		t.Fatalf("tl epic create stderr = %q, want create confirmation", allowedStderr)
	}
}

func TestWorkflowAuthoringDoesNotNeedPublicPolicyInspection(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	result := executeCommandForTest(workspace, "policy", "show")
	if result.err == nil || result.stdout != "" || !strings.Contains(result.err.Error(), `unknown command "policy"`) {
		t.Fatalf("policy show = %#v, want absent public command", result)
	}
}

func countActivityLogRows(t *testing.T, conn *sql.DB) int {
	t.Helper()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&count); err != nil {
		t.Fatalf("count activity_log rows: %v", err)
	}
	return count
}

// Retained V1 dev/completion fixtures continue to use this internal matrix helper
// until their replacement-bound callers are removed in task 20.
func commandPolicyCellAllowed(t *testing.T, rows []any, action, role string) bool {
	t.Helper()

	for _, row := range rows {
		fields, ok := row.(map[string]any)
		if !ok || fields["action"] != action {
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
