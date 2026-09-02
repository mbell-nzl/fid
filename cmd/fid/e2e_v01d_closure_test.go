package main

import "testing"

func TestV01DClosureLifecycleSmokeAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	assertRemovedV1CommandLeavesJournalUnchanged(t, workspace, "run", "record", "FID-legacy")
}
