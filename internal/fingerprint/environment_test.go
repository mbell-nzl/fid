package fingerprint_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/fingerprint"
)

//nolint:paralleltest // Changes the git environment inherited by child processes.
func TestComputeIgnoresHostileConfig(t *testing.T) {
	root, _ := fixture(t)
	write(t, root, "file", strings.Repeat("before\n\n", 20))
	write(t, root, "renamed", "rename contents\n")
	submoduleBase(t, root)
	base := commit(t, root)
	write(t, root, "file", strings.Repeat("before\n\n", 10)+"after\n\n"+strings.Repeat("before\n\n", 9))
	git(t, root, "mv", "renamed", "new-name")
	write(t, filepath.Join(root, "module"), "file", "submodule edit\n")
	commit(t, filepath.Join(root, "module"))
	write(t, root, "z file", "last\n")
	write(t, root, "é\tfile", "quoted\n")
	write(t, root, "binary", "\x00binary\n")
	head := commit(t, root)
	want := compute(t, root, base, head)
	order := filepath.Join(t.TempDir(), "order")
	write(t, filepath.Dir(order), "order", "z*\nfile\n*")
	attrs := filepath.Join(t.TempDir(), "attributes")
	write(t, filepath.Dir(attrs), "attributes", "* diff=hostile\n")
	global := filepath.Join(t.TempDir(), "gitconfig")
	write(t, filepath.Dir(global), "gitconfig", "")
	rows := [][2]string{
		{"diff.context", "20"}, {"diff.interHunkContext", "20"}, {"diff.algorithm", "histogram"},
		{"diff.indentHeuristic", "false"}, {"diff.noprefix", "true"}, {"diff.mnemonicPrefix", "true"},
		{"diff.renames", "copies"}, {"diff.renameLimit", "1"}, {"diff.external", "false"},
		{"color.ui", "always"}, {"core.pager", "false"}, {"diff.orderFile", order},
		{"core.quotePath", "false"}, {"core.attributesFile", attrs}, {"diff.hostile.textconv", "false"},
		{"diff.submodule", "log"}, {"diff.ignoreSubmodules", "all"}, {"diff.suppressBlankEmpty", "true"},
		{"diff.outputIndicatorNew", ">"}, {"diff.outputIndicatorOld", "<"}, {"diff.outputIndicatorContext", "!"},
		{"patchid.stable", "false"}, {"patchid.verbatim", "false"},
	}
	for _, scope := range []string{"repository", "global"} {
		t.Run(scope, func(t *testing.T) {
			for _, row := range rows {
				args := []string{"config"}
				if scope == "global" {
					args = append(args, "--file", global)
				}
				git(t, root, append(args, row[0], row[1])...)
			}
			t.Setenv("GIT_CONFIG_GLOBAL", global)
			t.Setenv("GIT_EXTERNAL_DIFF", "false")
			t.Setenv("GIT_DIFF_OPTS", "--unified=30")
			t.Setenv("GIT_PAGER", "false")
			if got := compute(t, root, base, head); got != want {
				t.Fatalf("hostile %s config changed result: %+v != %+v", scope, got, want)
			}
		})
		if scope == "repository" {
			for _, row := range rows {
				git(t, root, "config", "--unset", row[0])
			}
		}
	}
}

//nolint:paralleltest // Installs a git shim on PATH for real child-process failures.
func TestComputeChecksBothProcessesAndCompleteOutput(t *testing.T) {
	root, base := fixture(t)
	write(t, root, "file", "after\n")
	head := commit(t, root)
	installGitShim(t)
	for _, mode := range []string{"producer-failed", "consumer-failed", "no-id", "two-ids", "malformed-id", "wrong-source"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("FID_TEST_GIT_MODE", mode)
			refuse(t, root, base, head, fingerprint.GitFailure)
		})
	}
}

func installGitShim(t *testing.T) {
	t.Helper()
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	script := `#!/bin/sh
for arg do
  case "$arg:$FID_TEST_GIT_MODE" in
    diff:producer-failed) "$FID_TEST_REAL_GIT" "$@"; exit 17 ;;
    patch-id:consumer-failed) cat >/dev/null; printf '%040d %040d\n' 1 0; exit 18 ;;
    patch-id:no-id) cat >/dev/null; exit 0 ;;
    patch-id:two-ids) cat >/dev/null; printf '%040d %040d\n%040d %040d\n' 1 0 2 0; exit 0 ;;
    patch-id:malformed-id) cat >/dev/null; printf 'not a patch id\n'; exit 0 ;;
    patch-id:wrong-source) cat >/dev/null; printf '%040d %040d\n' 1 2; exit 0 ;;
  esac
done
exec "$FID_TEST_REAL_GIT" "$@"
`
	write(t, dir, "git", script)
	if err := os.Chmod(filepath.Join(dir, "git"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FID_TEST_REAL_GIT", realGit)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

//nolint:paralleltest // Sets a marker environment variable for a local promisor transport.
func TestComputeMissingPromisorObjectsNeverFetches(t *testing.T) {
	remote, _ := fixture(t)
	write(t, remote, "file", "only in remote\n")
	missing := commit(t, remote)
	root, base := fixture(t)
	marker := filepath.Join(t.TempDir(), "fetched")
	script := filepath.Join(t.TempDir(), "upload-pack")
	write(t, filepath.Dir(script), "upload-pack", "#!/bin/sh\ntouch \"$FID_TEST_FETCH_MARKER\"\nexit 1\n")
	if err := os.Chmod(script, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FID_TEST_FETCH_MARKER", marker)
	git(t, root, "remote", "add", "origin", remote)
	git(t, root, "config", "remote.origin.promisor", "true")
	git(t, root, "config", "remote.origin.partialclonefilter", "blob:none")
	git(t, root, "config", "remote.origin.uploadpack", script)
	for _, row := range [][2]string{{missing, base}, {base, missing}} {
		refuse(t, root, row[0], row[1], fingerprint.MissingObject)
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("fingerprinting invoked promisor transport: %v", err)
		}
	}
	// Positive control: this missing object really would invoke the configured
	// transport without the no-lazy-fetch setting, even though cat-file fails.
	command := exec.Command("git", "-C", root, "cat-file", "-e", missing)
	command.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=0")
	if output, err := command.CombinedOutput(); err == nil {
		t.Fatalf("unexpected fetch success: %s", output)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("promisor control did not invoke transport: %v", err)
	}
}

func TestComputeRejectsUnavailableBlob(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	write(t, root, "file", "after\n")
	head := commit(t, root)
	blob := git(t, root, "rev-parse", head+":file")
	if err := os.Remove(filepath.Join(root, ".git", "objects", blob[:2], blob[2:])); err != nil {
		t.Fatal(err)
	}
	refuse(t, root, base, head, fingerprint.GitFailure)
}

func TestComputeUnrelatedCommitsAreEmpty(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	git(t, root, "commit", "-q", "--allow-empty", "-m", "no change")
	refuse(t, root, base, git(t, root, "rev-parse", "HEAD"), fingerprint.EmptyChange)
}

func TestComputeRequiresFullCommitIDs(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	for _, invalid := range []string{base[:12], base + "^{commit}", strings.Repeat("z", 40), "--help", "", strings.Repeat("a", 64)} {
		t.Run(invalid, func(t *testing.T) { refuse(t, root, invalid, base, fingerprint.GitFailure) })
	}
}
