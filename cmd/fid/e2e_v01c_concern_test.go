package main

import "testing"

func TestV01CConcernRebuildParity(t *testing.T) {
	t.Parallel()
	workspace := newGitWorkspace(t)
	runInitCommand(t, workspace, "init")
	assertRemovedV1CommandLeavesJournalUnchanged(t, workspace, "issue", "complete", "FID-legacy")
}

func TestV01CClosureSmoke(t *testing.T) {
	t.Parallel()
	workspace := newGitWorkspace(t)
	runInitCommand(t, workspace, "init")
	assertRemovedV1CommandLeavesJournalUnchanged(t, workspace, "dev", "start", "FID-legacy")
}
