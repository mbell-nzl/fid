package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestMaintenanceVacuumCommandIsUnavailableAndPreservesState(t *testing.T) {
	workspace := clonePristineWorkspace(t)

	journalPath := mustJournalPath(t, workspace)
	journalBefore := mustReadCommandFile(t, journalPath)

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	activityBefore := countActivityLogRows(t, conn)

	for _, args := range [][]string{{"maintenance", "vacuum"}, {"maintenance", "vacuum", "--by", "tl"}} {
		result := executeCommandForTest(workspace, args...)
		if result.err == nil || result.stdout != "" {
			t.Fatalf("maintenance vacuum %v = %#v, want rejected without output", args, result)
		}
		want := "capability is unavailable"
		if containsTestArgument(args, "--by") {
			want = "unknown flag: --by"
		}
		if !strings.Contains(result.stderr, want) {
			t.Fatalf("maintenance vacuum %v = %#v, want %q", args, result, want)
		}
	}
	if journalAfter := mustReadCommandFile(t, journalPath); !bytes.Equal(journalAfter, journalBefore) {
		t.Fatal("maintenance vacuum changed journal bytes")
	}
	if activityAfter := countActivityLogRows(t, conn); activityAfter != activityBefore {
		t.Fatalf("activity_log rows after vacuum = %d, want %d", activityAfter, activityBefore)
	}
}

func TestMaintenanceVacuumCommandRejectsUnauthorizedRolesBeforeMutation(t *testing.T) {
	workspace := clonePristineWorkspace(t)

	journalPath := mustJournalPath(t, workspace)
	journalBefore, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}

	for _, args := range [][]string{{"maintenance", "vacuum"}, {"maintenance", "vacuum", "--by", "unknown"}} {
		result := executeCommandForTest(workspace, args...)
		if result.err == nil || result.stdout != "" {
			t.Fatalf("maintenance vacuum %v = %#v, want rejected without output", args, result)
		}
		want := "capability is unavailable"
		if containsTestArgument(args, "--by") {
			want = "unknown flag: --by"
		}
		if !strings.Contains(result.stderr, want) {
			t.Fatalf("maintenance vacuum %v = %#v, want %q", args, result, want)
		}
		if journalAfter := mustReadCommandFile(t, journalPath); !bytes.Equal(journalAfter, journalBefore) {
			t.Fatal("unavailable maintenance vacuum changed journal bytes")
		}
	}
}

func TestMaintenanceVacuumHasNoPublicPolicySurface(t *testing.T) {
	workspace := clonePristineWorkspace(t)

	result := executeCommandForTest(workspace, "policy", "show")
	if result.err == nil || result.stdout != "" || !strings.Contains(result.err.Error(), `unknown command "policy"`) {
		t.Fatalf("policy show = %#v, want absent public command", result)
	}
}
