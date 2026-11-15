package fingerprint_test

import (
	"fmt"
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

func TestInspectKeepsEveryPathBeyondFortyFiles(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	var paths []string
	for i := range 45 {
		name := fmt.Sprintf("path-%02d.txt", i)
		paths = append(paths, name)
		write(t, root, name, "new file\n")
	}
	head := commit(t, root)
	_, stat, err := fingerprint.Inspect(root, base, head)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		if !strings.Contains(stat, path) {
			t.Errorf("path %q missing from full diffstat:\n%s", path, stat)
		}
	}
	if strings.Contains(stat, "...") || !strings.Contains(stat, "45 files changed") {
		t.Fatalf("incomplete diffstat:\n%s", stat)
	}
}
