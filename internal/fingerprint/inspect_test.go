package fingerprint_test

import (
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/fingerprint"
)

func TestInspectShowsOnlyFullPinnedDiffstat(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	write(t, root, "file", "after\n")
	head := commit(t, root)
	fp, stat, err := fingerprint.Inspect(root, base, head)
	if err != nil || fp.Algorithm != fingerprint.Algorithm || !strings.Contains(stat, "file changed") || strings.Contains(stat, "diff --git") {
		t.Fatalf("inspect=%v, %q, %v", fp, stat, err)
	}
	computed, err := fingerprint.Compute(root, base, head)
	if err != nil || computed != fp {
		t.Fatalf("inspect binding=%v compute=%v/%v", fp, computed, err)
	}
}

func TestInspectKeepsEveryPathDespiteUserStatSettings(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	name := strings.Repeat("long-name-", 25) + ".txt"
	write(t, root, name, "new file\n")
	write(t, root, "file", "after\n")
	head := commit(t, root)
	git(t, root, "config", "diff.statNameWidth", "4")
	git(t, root, "config", "diff.statGraphWidth", "2")
	git(t, root, "config", "diff.statCount", "1")
	_, stat, err := fingerprint.Inspect(root, base, head)
	if err != nil || !strings.Contains(stat, name) || !strings.Contains(stat, " file ") || !strings.Contains(stat, "2 files changed") || strings.Contains(stat, "...") {
		t.Fatalf("full diffstat=%q/%v", stat, err)
	}
}
