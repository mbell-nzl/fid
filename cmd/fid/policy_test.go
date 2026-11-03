package main

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"

	"github.com/mbell-nzl/fid/internal/config"
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

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
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
	grammarCheck(t, grammarRow{name: "absent policy with removed flag", operation: "fid", code: "invalid-argument", argument: "--by", value: "--by", allowed: []string{}, message: "unknown flag: --by"}, false, result, cli.ExitCode(result.err))
	if after := mustReadCommandFile(t, journalPath); !bytes.Equal(after, beforePayload) || countActivityLogRows(t, conn) != activityBefore {
		t.Fatal("absent policy with removed flag mutated state")
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
