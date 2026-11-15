package fingerprint_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestComputeUsesCommittedAttributesAndIgnoresDriverConfig(t *testing.T) {
	t.Parallel()
	root, _ := fixture(t)
	write(t, root, ".gitattributes", "file diff=custom\n")
	base := commit(t, root)
	write(t, root, "file", "after\n")
	head := commit(t, root)
	want := compute(t, root, base, head)
	for _, row := range []struct {
		name   string
		change func(*testing.T)
	}{
		{"driver config", func(t *testing.T) { git(t, root, "config", "diff.custom.binary", "true") }},
		{"worktree attributes", func(t *testing.T) { write(t, root, ".gitattributes", "file -diff\n") }},
		{"index attributes", func(t *testing.T) { git(t, root, "add", ".gitattributes") }},
		{"local attributes", func(t *testing.T) {
			info := filepath.Join(root, ".git", "info")
			if err := os.MkdirAll(info, 0o755); err != nil {
				t.Fatal(err)
			}
			write(t, info, "attributes", "file -diff\n")
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			row.change(t)
			index := filepath.Join(root, ".git", "index")
			before, err := os.ReadFile(index)
			if err != nil {
				t.Fatal(err)
			}
			if got := compute(t, root, base, head); got != want {
				t.Fatalf("%s changed recorded fingerprint: %+v != %+v", row.name, got, want)
			}
			after, err := os.ReadFile(index)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("fingerprinting changed caller index: %v", err)
			}
		})
	}
}

func TestComputeAppliesCommittedBinaryAttribute(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	write(t, root, "file", "after\n")
	head := commit(t, root)
	text := compute(t, root, base, head)
	write(t, root, ".gitattributes", "file -diff\n")
	binaryHead := commit(t, root)
	// Put the same attribute in both trees, so only file changes are hashed.
	git(t, root, "checkout", "-q", "--detach", base)
	write(t, root, ".gitattributes", "file -diff\n")
	binaryBase := commit(t, root)
	git(t, root, "checkout", "-q", "--detach", head)
	binary := compute(t, root, binaryBase, binaryHead)
	if binary.Fingerprint == text.Fingerprint {
		t.Fatal("committed binary attribute did not apply")
	}
}

func TestComputeObjectStorageLayouts(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	write(t, root, "file", "after\n")
	head := commit(t, root)
	want := compute(t, root, base, head)
	git(t, root, "gc", "--quiet")
	if got := compute(t, root, base, head); got != want {
		t.Fatalf("packing changed fingerprint: %+v != %+v", got, want)
	}
	linked := gittest.AddWorktree(t, root)
	if got := compute(t, linked, base, head); got != want {
		t.Fatalf("worktree changed fingerprint: %+v != %+v", got, want)
	}
	clone := filepath.Join(t.TempDir(), "shared")
	git(t, root, "clone", "--quiet", "--shared", root, clone)
	if err := gittest.CheckRepository(clone); err != nil {
		t.Fatal(err)
	}
	if got := compute(t, clone, base, head); got != want {
		t.Fatalf("alternates changed fingerprint: %+v != %+v", got, want)
	}
}

func TestComputeRefusesNULTextPatches(t *testing.T) {
	t.Parallel()
	root, _ := fixture(t)
	// A committed attribute can force NUL-containing data into a text hunk.
	// patch-id silently ignores bytes after NUL, even in --verbatim mode.
	write(t, root, ".gitattributes", "file diff\n")
	base := commit(t, root)
	for _, content := range []string{"after\x00A\n", "after\x00B\n"} {
		write(t, root, "file", content)
		refuse(t, root, base, commit(t, root), fingerprint.GitFailure)
	}
}
