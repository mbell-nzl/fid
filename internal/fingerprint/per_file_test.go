package fingerprint_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestComputeMixedBinaryChangesBindEveryPath(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"mode", "binary"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			root := gittest.NewWorkspace(t)
			for _, name := range []string{"a.bin", "z-one", "z-two"} {
				write(t, root, name, "\x00before\n")
			}
			base := commit(t, root)
			var previous string
			for _, name := range []string{"z-one", "z-two"} {
				git(t, root, "reset", "--hard", base)
				write(t, root, "a.bin", "\x00after\n")
				if kind == "mode" {
					if err := os.Chmod(filepath.Join(root, name), 0o755); err != nil {
						t.Fatal(err)
					}
				} else {
					write(t, root, name, "\x00after\n")
				}
				head := commit(t, root)
				got := compute(t, root, base, head)
				if got.Fingerprint == previous {
					t.Fatalf("moving %s change from z-one to z-two kept fingerprint %s", kind, got.Fingerprint)
				}
				if repeat := compute(t, root, base, head); repeat != got {
					t.Fatalf("unstable mixed change: %+v != %+v", repeat, got)
				}
				previous = got.Fingerprint
			}
		})
	}
}

func TestComputeSectionBoundariesCannotBeForged(t *testing.T) {
	t.Parallel()
	root := gittest.NewWorkspace(t)
	// Newlines in paths must be quoted; header-like content must keep its
	// added/deleted/context prefix, even when the final newline is absent.
	names := []string{"a\ndiff --git forged", "z.bin"}
	write(t, root, names[0], "diff --git a/deleted b/deleted\ndiff --git a/context b/context\nold\n")
	write(t, root, names[1], "\x00before\n")
	base := commit(t, root)
	write(t, root, names[0], "diff --git a/context b/context\ndiff --git a/added b/added")
	write(t, root, names[1], "\x00after\n")
	head := commit(t, root)
	want := perFileFingerprint(t, root, base, head, names)
	got, err := fingerprint.Compute(root, base, head)
	if err != nil || got.Fingerprint != want {
		t.Fatalf("header-like content changed section boundaries: %+v, %v; want %s", got, err, want)
	}
}

// Construct the reference from independent path-limited diffs, so a bug in
// splitting the aggregate cannot also affect the expected result.
func perFileFingerprint(t *testing.T, root, base, head string, names []string) string {
	t.Helper()
	var ids []string
	for _, name := range names {
		patch := git(t, root, "-c", "core.quotePath=true", "diff", "--binary", "--full-index", base, head, "--", name) + "\n"
		command := exec.Command("git", "patch-id", "--verbatim")
		command.Stdin = strings.NewReader(patch)
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("reference patch-id: %v: %s", err, output)
		}
		fields := strings.Fields(string(output))
		if len(fields) != 2 {
			t.Fatalf("reference patch-id: %q", output)
		}
		ids = append(ids, fields[0])
	}
	slices.Sort(ids)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d\n", len(ids))+strings.Join(ids, "\n")+"\n")))
}

//nolint:paralleltest // Installs a git shim and supplies reordered producer output.
func TestComputePerFileOrderAndCompleteSections(t *testing.T) {
	root, base := fixture(t)
	write(t, root, "file", "after\n")
	write(t, root, "z", "\x00after\n")
	head := commit(t, root)
	names := []string{"file", "z"}
	want := perFileFingerprint(t, root, base, head, names)
	var sections []string
	for _, name := range names {
		sections = append(sections, git(t, root, "diff", "--binary", "--full-index", base, head, "--", name)+"\n")
	}
	installGitShim(t)
	dir := t.TempDir()
	t.Setenv("FID_TEST_DIFF", filepath.Join(dir, "diff"))
	t.Setenv("FID_TEST_GIT_MODE", "fixture-diff")
	for _, row := range []struct{ name, patch string }{
		{"forward", sections[0] + sections[1]},
		{"reverse", sections[1] + sections[0]},
	} {
		t.Run(row.name, func(t *testing.T) {
			write(t, dir, "diff", row.patch)
			if got := compute(t, root, base, head); got.Fingerprint != want {
				t.Fatalf("%s sections: %s, want %s", row.name, got.Fingerprint, want)
			}
		})
	}
	for _, patch := range []string{"unexpected output\n", "unexpected prefix\n" + sections[0]} {
		write(t, dir, "diff", patch)
		refuse(t, root, base, head, fingerprint.GitFailure)
	}
	// A valid first section must never hide a failed or incomplete later one.
	t.Setenv("FID_TEST_LATER_SECTION", "1")
	for _, mode := range []string{"consumer-failed", "no-id", "two-ids", "malformed-id", "wrong-source", "oversized-id"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("FID_TEST_GIT_MODE", mode)
			refuse(t, root, base, head, fingerprint.GitFailure)
		})
	}
}
