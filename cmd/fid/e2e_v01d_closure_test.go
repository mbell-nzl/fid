package main

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestV01DClosureLifecycleSmokeAndRebuildParity(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init")
	assertRemovedV1CommandLeavesJournalUnchanged(t, workspace, "run", "record", "FID-legacy")
}
