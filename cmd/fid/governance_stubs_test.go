package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestGovernanceStubCommandsReturnVersionGateBeforeValidation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	journalPath := mustJournalPath(t, workspace)
	journalBefore := mustReadCommandFile(t, journalPath)
	countsBefore := governanceStubTableCounts(t, workspace)

	for _, testCase := range []struct {
		name       string
		args       []string
		wantStderr string
	}{
		{
			name:       "decision create accepts future flag before validation",
			args:       []string{"decision", "create", "--future-flag", "value"},
			wantStderr: "error: commands for decision ship in v0.2\n",
		},
		{
			name:       "decision get without future required argument",
			args:       []string{"decision", "get"},
			wantStderr: "error: commands for decision ship in v0.2\n",
		},
		{
			name:       "decision list accepts future actor flag",
			args:       []string{"decision", "list", "--by", "developer"},
			wantStderr: "error: commands for decision ship in v0.2\n",
		},
		{
			name:       "escalation create accepts future flag",
			args:       []string{"escalation", "create", "--future-flag", "value"},
			wantStderr: "error: commands for escalation ship in v0.2\n",
		},
		{
			name:       "escalation resolve without future required argument",
			args:       []string{"escalation", "resolve"},
			wantStderr: "error: commands for escalation ship in v0.2\n",
		},
		{
			name:       "escalation get accepts future actor flag",
			args:       []string{"escalation", "get", "--by", "developer"},
			wantStderr: "error: commands for escalation ship in v0.2\n",
		},
		{
			name:       "escalation list without future flags",
			args:       []string{"escalation", "list"},
			wantStderr: "error: commands for escalation ship in v0.2\n",
		},
		{
			name:       "triage record accepts future flag",
			args:       []string{"triage", "record", "--future-flag", "value"},
			wantStderr: "error: commands for triage ship in v0.2\n",
		},
		{
			name:       "triage start without future required argument",
			args:       []string{"triage", "start"},
			wantStderr: "error: commands for triage ship in v0.2\n",
		},
		{
			name:       "triage verdict accepts future actor flag",
			args:       []string{"triage", "verdict", "--by", "developer"},
			wantStderr: "error: commands for triage ship in v0.2\n",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, testCase.args...)
			if err == nil {
				t.Fatal("Execute() error = nil, want version gate")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if stderr != testCase.wantStderr {
				t.Fatalf("stderr = %q, want %q", stderr, testCase.wantStderr)
			}
			if journalAfter := mustReadCommandFile(t, journalPath); !bytes.Equal(journalAfter, journalBefore) {
				t.Fatal("version gate changed journal bytes")
			}
			if countsAfter := governanceStubTableCounts(t, workspace); !governanceStubCountsEqual(countsAfter, countsBefore) {
				t.Fatalf("governance table counts after gate = %#v, want %#v", countsAfter, countsBefore)
			}
		})
	}
}

func TestGovernanceStubVersionGateWinsBeforeWorkspaceValidation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "escalation", "create", "--future-flag", "value")
	if err == nil {
		t.Fatal("Execute() error = nil, want version gate")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if stderr != "error: commands for escalation ship in v0.2\n" {
		t.Fatalf("stderr = %q, want version gate", stderr)
	}
}

func TestGovernanceStubFamiliesAppearInRootHelp(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()

	stdout, stderr := runFIDCommand(t, workspace, "--help")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	for _, family := range []string{"decision", "escalation", "triage"} {
		if !strings.Contains(stdout, family) {
			t.Fatalf("root help = %q, want %q command family", stdout, family)
		}
	}
	if !strings.Contains(stdout, "Governance commands ship in v0.2.") {
		t.Fatalf("root help = %q, want version-gated governance description", stdout)
	}
}

func TestRequirementSyncRemainsShippedAlongsideGovernanceStubGuards(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	writeRequirementFixture(t, workspace, "project-docs/specs/control.md", `# Control

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | Requirement sync remains shipped. | Must | v0.1 |
`)

	journalPath := mustJournalPath(t, workspace)
	journalBefore := mustReadCommandFile(t, journalPath)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"requirement", "sync",
		"--source-doc", filepath.ToSlash(filepath.Join("project-docs", "specs", "control.md")),
		"--by", "tl",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "synced requirements from project-docs/specs/control.md") {
		t.Fatalf("stderr = %q, want sync confirmation", stderr)
	}
	if journalAfter := mustReadCommandFile(t, journalPath); len(splitJournalLines(journalAfter)) != len(splitJournalLines(journalBefore))+1 {
		t.Fatalf("journal lines after requirement sync = %d, want %d", len(splitJournalLines(journalAfter)), len(splitJournalLines(journalBefore))+1)
	}
}

func governanceStubTableCounts(t *testing.T, workspace string) map[string]int {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)

	counts := make(map[string]int, 3)
	for _, table := range []string{"decisions", "escalations", "triage_events"} {
		var count int
		if err := conn.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		counts[table] = count
	}
	return counts
}

func governanceStubCountsEqual(got, want map[string]int) bool {
	for table, wantCount := range want {
		if got[table] != wantCount {
			return false
		}
	}
	return true
}
