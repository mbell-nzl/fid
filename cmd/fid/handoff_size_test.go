package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// A hand-off records every changed path. One whose Event would reach the
// journal's line budget is refused before append, so it cannot stop every
// later write (RESET 1.9).
func TestHandoffEventSizeBudget(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "[]")
	stem := strings.Repeat("p", 200)
	for i := range 6000 {
		writeBackupTestFile(t, filepath.Join(f.workspace, "wide", fmt.Sprintf("%s%05d.txt", stem, i)), "x\n")
	}
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-q", "-m", "A wide change")
	requiredNames(t, checkRequired(t, f, "checks"))
	before := snapshotWorkspacePersistence(t, f.workspace)
	got := f.run(append([]string{"--json"}, handoffArgs("wide")...)...)
	if got.err == nil || got.stdout != "" || !strings.Contains(got.stderr, "exceeds the journal's 1 MiB line budget") {
		t.Fatalf("oversized hand-off: error=%v, stdout bytes=%d, stderr=%.300s", got.err, len(got.stdout), got.stderr)
	}
	if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
		t.Fatal("refused hand-off changed the journal or projection")
	}
	// A later write must still be able to scan the journal and allocate IDs.
	lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "After hand-off", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "after-handoff"), "issue.create")
}
