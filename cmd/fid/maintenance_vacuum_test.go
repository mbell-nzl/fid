package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestMaintenanceVacuumCommandAuthorizesCompactsAndPreservesState(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	journalPath := mustJournalPath(t, workspace)
	journalBefore := mustReadCommandFile(t, journalPath)

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	activityBefore := countActivityLogRows(t, conn)

	for _, role := range []string{"tl", "human"} {
		t.Run(role, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "maintenance", "vacuum", "--by", role)
			if err != nil {
				t.Fatalf("maintenance vacuum error = %v", err)
			}
			if stdout != "" {
				t.Fatalf("maintenance vacuum stdout = %q, want empty", stdout)
			}
			if !regexp.MustCompile(`(?m)^vacuum complete — \.fid/fid\.sqlite: \d+ B → \d+ B$`).MatchString(stderr) {
				t.Fatalf("maintenance vacuum stderr = %q, want size confirmation", stderr)
			}
		})
	}
	if journalAfter := mustReadCommandFile(t, journalPath); !bytes.Equal(journalAfter, journalBefore) {
		t.Fatal("maintenance vacuum changed journal bytes")
	}
	if activityAfter := countActivityLogRows(t, conn); activityAfter != activityBefore {
		t.Fatalf("activity_log rows after vacuum = %d, want %d", activityAfter, activityBefore)
	}
}

func TestMaintenanceVacuumCommandRejectsUnauthorizedRolesBeforeMutation(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	journalPath := mustJournalPath(t, workspace)
	journalBefore, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}

	for _, testCase := range []struct {
		role       string
		wantStderr string
	}{
		{role: "", wantStderr: "error: --by is required"},
		{role: "unknown", wantStderr: `error: role "unknown" not in project profile`},
		{role: "developer", wantStderr: `error: role "developer" not permitted for action "maintenance.vacuum"`},
		{role: "qa_reviewer", wantStderr: `error: role "qa_reviewer" not permitted for action "maintenance.vacuum"`},
		{role: "qa_lead", wantStderr: `error: role "qa_lead" not permitted for action "maintenance.vacuum"`},
		{role: "tl_lite", wantStderr: `error: role "tl_lite" not permitted for action "maintenance.vacuum"`},
		{role: "workflow_runner", wantStderr: `error: role "workflow_runner" not permitted for action "maintenance.vacuum"`},
	} {
		role := testCase.role
		t.Run(role, func(t *testing.T) {
			args := []string{"maintenance", "vacuum"}
			if role != "" {
				args = append(args, "--by", role)
			}
			stdout, stderr, commandErr := runFIDCommandExpectError(t, workspace, args...)
			if commandErr == nil {
				t.Fatal("maintenance vacuum error = nil, want role rejection")
			}
			if stdout != "" {
				t.Fatalf("maintenance vacuum stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, testCase.wantStderr) {
				t.Fatalf("maintenance vacuum stderr = %q, want %q", stderr, testCase.wantStderr)
			}
			if journalAfter := mustReadCommandFile(t, journalPath); !bytes.Equal(journalAfter, journalBefore) {
				t.Fatal("rejected maintenance vacuum changed journal bytes")
			}
		})
	}
}

func TestMaintenanceVacuumPolicyShowMatchesEnforcement(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show")
	if stderr != "" {
		t.Fatalf("policy show stderr = %q, want empty", stderr)
	}
	rows := decodeCommandJSONStructure(t, stdout).([]any)
	for _, testCase := range []struct {
		role string
		want bool
	}{
		{role: "tl", want: true},
		{role: "human", want: true},
		{role: "workflow_runner", want: false},
	} {
		if got := commandPolicyCellAllowed(t, rows, "maintenance.vacuum", testCase.role); got != testCase.want {
			t.Fatalf("policy show maintenance.vacuum/%s = %t, want %t", testCase.role, got, testCase.want)
		}
	}
}
