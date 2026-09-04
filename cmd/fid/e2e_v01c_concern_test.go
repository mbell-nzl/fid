package main

import "testing"

func TestV01CConcernRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	assertRemovedV1CommandLeavesJournalUnchanged(t, workspace, "issue", "complete", "FID-legacy")
}

func TestV01CClosureSmoke(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	assertRemovedV1CommandLeavesJournalUnchanged(t, workspace, "dev", "start", "FID-legacy")
}
