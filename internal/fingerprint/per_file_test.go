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
		patch := referenceFilePatch(t, root, base, head, name)
		fields := strings.Fields(referenceGit(t, root, patch, "patch-id", "--verbatim"))
		if len(fields) != 2 {
			t.Fatalf("reference patch-id: %q", fields)
		}
		ids = append(ids, fields[0])
	}
	slices.Sort(ids)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d\n", len(ids))+strings.Join(ids, "\n")+"\n")))
}

// Pin the reference renderer independently of production's aggregate diff.
// Fixtures use the same renderer so injected sections obey that contract too.
func referenceFilePatch(t *testing.T, root, base, head, name string) string {
	t.Helper()
	return referenceGit(t, root, "", "diff", "--patch", "--no-color", "--no-ext-diff", "--no-textconv",
		"--no-relative", "--src-prefix=a/", "--dst-prefix=b/", "--line-prefix=", "--unified=3", "--inter-hunk-context=0",
		"--diff-algorithm=myers", "--indent-heuristic", "--find-renames=50%", "-l0", "--full-index", "--binary",
		"--ignore-submodules=none", "--submodule=short", "--output-indicator-new=+", "--output-indicator-old=-",
		"--output-indicator-context= ", "-O"+os.DevNull, base, head, "--", name)
}

func referenceGit(t *testing.T, root, input string, args ...string) string {
	t.Helper()
	argv := []string{"--no-pager", "--no-replace-objects", "-C", root, "-c", "core.quotePath=true",
		"-c", "core.attributesFile=" + os.DevNull, "-c", "diff.suppressBlankEmpty=false"}
	command := exec.Command("git", append(argv, args...)...)
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "GIT_") || strings.HasPrefix(entry, "GIT_CEILING_DIRECTORIES=") {
			command.Env = append(command.Env, entry)
		}
	}
	command.Env = append(command.Env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull,
		"GIT_ATTR_NOSYSTEM=1", "GIT_NO_LAZY_FETCH=1", "GIT_OPTIONAL_LOCKS=0", "GIT_TERMINAL_PROMPT=0", "LC_ALL=C")
	command.Stdin = strings.NewReader(input)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("reference git %v: %v: %s", args, err, output)
	}
	return string(output)
}

//nolint:paralleltest // Changes the rendering configuration inherited by git.
func TestPerFileReferenceIgnoresHostileConfig(t *testing.T) {
	root, _ := fixture(t)
	original := strings.Repeat("before\n", 20)
	write(t, root, "file", original)
	base := commit(t, root)
	write(t, root, "file", strings.Repeat("before\n", 10)+"after\n"+strings.Repeat("before\n", 9))
	head := commit(t, root)
	names := []string{"file"}
	want := perFileFingerprint(t, root, base, head, names)
	for _, scope := range []string{"global", "environment", "repository"} {
		t.Run(scope, func(t *testing.T) {
			switch scope {
			case "global":
				dir := t.TempDir()
				write(t, dir, "gitconfig", "[diff]\n\tcontext = 0\n")
				t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(dir, "gitconfig"))
			case "environment":
				t.Setenv("GIT_CONFIG_COUNT", "1")
				t.Setenv("GIT_CONFIG_KEY_0", "diff.context")
				t.Setenv("GIT_CONFIG_VALUE_0", "0")
				t.Setenv("GIT_DIFF_OPTS", "--unified=0")
			case "repository":
				git(t, root, "config", "diff.context", "0")
			}
			if got := perFileFingerprint(t, root, base, head, names); got != want {
				t.Fatalf("hostile %s config changed reference: %s != %s", scope, got, want)
			}
			if got := compute(t, root, base, head); got.Fingerprint != want {
				t.Fatalf("reference differs from Compute: %s != %s", want, got.Fingerprint)
			}
		})
	}
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
		sections = append(sections, referenceFilePatch(t, root, base, head, name))
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
