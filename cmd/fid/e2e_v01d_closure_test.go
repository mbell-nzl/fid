package main

import "testing"

func TestV01DClosureLifecycleSmokeAndRebuildParity(t *testing.T) {
	t.Parallel()
	workspace := newGitWorkspace(t)
	runInitCommand(t, workspace, "init")
	assertRemovedV1CommandLeavesJournalUnchanged(t, workspace, "run", "record", "FID-legacy")
}
