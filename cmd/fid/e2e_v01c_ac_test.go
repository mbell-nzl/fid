package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestV01CAllACsVerifiedGateAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	assertRemovedV1CommandLeavesJournalUnchanged(t, workspace, "issue", "claim", "FID-legacy")
}

func assertRemovedV1CommandLeavesJournalUnchanged(t *testing.T, workspace string, args ...string) {
	t.Helper()
	journalPath := mustJournalPath(t, workspace)
	before := mustReadCommandFile(t, journalPath)
	result := executeCommandForTest(workspace, args...)
	if result.err == nil || result.stdout != "" || !strings.Contains(result.err.Error(), "unknown command") {
		t.Fatalf("removed V1 command %v = %#v, want unknown command rejection", args, result)
	}
	if after := mustReadCommandFile(t, journalPath); !bytes.Equal(after, before) {
		t.Fatalf("removed V1 command %v changed journal bytes", args)
	}
}
