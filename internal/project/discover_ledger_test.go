package project

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// Discovery accepts every layout where git and LedgerDir name the same ledger
// and git can't track it: the common directory is within a .git directory,
// or outside the worktree. Discover returns the top level, as before.
func TestDiscoverAcceptsLayoutsWhereGitCannotTrackTheLedger(t *testing.T) {
	t.Parallel()
	for _, layout := range []struct {
		name  string
		build func(t *testing.T) (start, root string)
	}{
		{"main worktree", func(t *testing.T) (string, string) {
			repository := newDiscoveryRepository(t)
			return repository, repository
		}},
		{"linked worktree", func(t *testing.T) (string, string) {
			linked := gittest.ResolvedPath(t, gittest.AddWorktree(t, newDiscoveryRepository(t)))
			return linked, linked
		}},
		{"worktree inside the main worktree", func(t *testing.T) (string, string) {
			repository := newDiscoveryRepository(t)
			discoveryCommit(t, repository)
			inner := filepath.Join(repository, "inner")
			discoveryGit(t, repository, "worktree", "add", "-q", "--detach", inner)
			return inner, inner
		}},
		{"submodule", func(t *testing.T) (string, string) {
			source := newDiscoveryRepository(t)
			discoveryCommit(t, source)
			superproject := newDiscoveryRepository(t)
			discoveryGit(t, superproject, "-c", "protocol.file.allow=always", "submodule", "add", "-q", source, "module")
			module := filepath.Join(superproject, "module")
			if common := discoveryGitCommonDir(t, module); common != filepath.Join(superproject, ".git", "modules", "module") {
				t.Fatalf("submodule common directory = %q, want it in the superproject's .git", common)
			}
			return module, module
		}},
		{"separate git directory outside the worktree", func(t *testing.T) (string, string) {
			worktree := gittest.ResolvedPath(t, t.TempDir())
			separateGitDir(t, worktree, filepath.Join(t.TempDir(), "repository.git"))
			return worktree, worktree
		}},
		{"separate git directory within a .git directory inside the worktree", func(t *testing.T) (string, string) {
			worktree := gittest.ResolvedPath(t, t.TempDir())
			separateGitDir(t, worktree, filepath.Join(worktree, "vendor", ".git", "metadata"))
			return worktree, worktree
		}},
		{"core.worktree naming a worktree whose .git leads back", func(t *testing.T) (string, string) {
			return coreWorktreeLeadingBack(t, newDiscoveryRepository(t), gittest.ResolvedPath(t, t.TempDir()))
		}},
		{".git is a symlink", func(t *testing.T) (string, string) {
			worktree := gittest.ResolvedPath(t, t.TempDir())
			gitDir := filepath.Join(t.TempDir(), "repository.git")
			separateGitDir(t, worktree, gitDir)
			symlinkDotGit(t, worktree, gitDir)
			return worktree, worktree
		}},
	} {
		t.Run(layout.name, func(t *testing.T) {
			t.Parallel()
			start, root := layout.build(t)
			nested := filepath.Join(start, "nested")
			if err := os.Mkdir(nested, 0o755); err != nil {
				t.Fatalf("Mkdir(nested) error = %v", err)
			}
			for _, from := range []string{start, nested} {
				if got, err := Discover(from); err != nil || got != root {
					t.Fatalf("Discover(%q) = %q, %v; want top level %q", from, got, err, root)
				}
			}
		})
	}
}

// With core.worktree, git reports another directory's worktree as the top
// level, while the common directory stays where git found it. LedgerDir
// reads the top level's own .git, so git and LedgerDir name different
// ledgers, and discovery refuses.
func TestDiscoverRefusesWhereGitAndLedgerDirNameDifferentLedgers(t *testing.T) {
	t.Parallel()
	worktree := newDiscoveryRepository(t)
	other := newDiscoveryRepository(t)
	discoveryGit(t, other, "config", "core.worktree", worktree)
	_, err := Discover(other)
	assertUnsafeLedger(t, err, filepath.Join(other, ".git", config.LedgerDirName), worktree, config.LedgerDir(worktree))
}

// GIT_DIR or GIT_COMMON_DIR in the environment make git name another
// repository's common directory, the main way git and LedgerDir disagree.
// Naming the worktree's own git directory, as git does for hooks, is fine.
func TestDiscoverRefusesWhereTheEnvironmentNamesAnotherLedger(t *testing.T) { //nolint:paralleltest // It sets GIT_DIR and GIT_COMMON_DIR.
	worktree := newDiscoveryRepository(t)
	linked := gittest.ResolvedPath(t, gittest.AddWorktree(t, worktree))
	gitDirs := map[string]string{}
	for _, start := range []string{worktree, linked} {
		gitDirs[start] = discoveryGitText(t, start, "rev-parse", "--absolute-git-dir")
	}
	other := newDiscoveryRepository(t)
	otherCommon := filepath.Join(other, ".git")
	for _, variable := range []string{"GIT_DIR", "GIT_COMMON_DIR"} {
		t.Run(variable, func(t *testing.T) {
			t.Setenv(variable, otherCommon)
			_, err := Discover(worktree)
			assertUnsafeLedger(t, err, filepath.Join(otherCommon, config.LedgerDirName), worktree, config.LedgerDir(worktree))
		})
	}
	t.Run("GIT_DIR naming the worktree's own git directory", func(t *testing.T) {
		for _, start := range []string{worktree, linked} {
			gitDir := gitDirs[start]
			t.Setenv("GIT_DIR", gitDir)
			if got, err := Discover(start); err != nil || got != start {
				t.Fatalf("Discover(%q) with GIT_DIR=%q = %q, %v; want %q", start, gitDir, got, err, start)
			}
		}
	})
}

// Where git's common directory sits inside the worktree but outside any .git
// directory, `git add --all` would track the ledger, so discovery refuses.
// Ignore rules wouldn't help: `git add -f` still adds an ignored file.
func TestDiscoverRefusesWhereGitCouldTrackTheLedger(t *testing.T) {
	t.Parallel()
	for _, layout := range []struct {
		name  string
		build func(t *testing.T, worktree string) (common string)
	}{
		{"separate git directory inside the worktree", func(t *testing.T, worktree string) string {
			common := filepath.Join(worktree, "metadata")
			separateGitDir(t, worktree, common)
			return common
		}},
		{".git is a symlink into the worktree", func(t *testing.T, worktree string) string {
			common := filepath.Join(worktree, "metadata")
			separateGitDir(t, worktree, common)
			symlinkDotGit(t, worktree, "metadata")
			return common
		}},
	} {
		t.Run(layout.name, func(t *testing.T) {
			t.Parallel()
			worktree := gittest.ResolvedPath(t, t.TempDir())
			common := layout.build(t, worktree)
			_, err := Discover(worktree)
			assertUnsafeLedger(t, err, filepath.Join(common, config.LedgerDirName), worktree, common)
		})
	}
}

// assertUnsafeLedger checks that err is the unsafe-ledger refusal: reason
// unsafe-path, with a detail naming every location given.
func assertUnsafeLedger(t *testing.T, err error, locations ...string) {
	t.Helper()
	var unavailable *RepositoryUnavailableError
	if !errors.Is(err, ErrUnsafeLedger) || !errors.As(err, &unavailable) || unavailable.Reason != "unsafe-path" {
		t.Fatalf("Discover error = %v, want ErrUnsafeLedger with reason unsafe-path", err)
	}
	if errors.Is(err, ErrNotARepository) || errors.Is(err, ErrNotFound) {
		t.Fatalf("Discover error = %v, want neither not-a-repository nor not-found", err)
	}
	for _, location := range locations {
		if !strings.Contains(unavailable.Detail, location) {
			t.Errorf("refusal detail %q does not name %q", unavailable.Detail, location)
		}
	}
}

// separateGitDir makes worktree a git repository whose git directory is
// gitDir, as `git init --separate-git-dir` does: worktree/.git is a file
// naming gitDir.
func separateGitDir(t *testing.T, worktree, gitDir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(gitDir), 0o755); err != nil {
		t.Fatalf("create the git directory's parent: %v", err)
	}
	if output, err := exec.Command("git", "init", "-q", "--template=", "--separate-git-dir="+gitDir, worktree).CombinedOutput(); err != nil {
		t.Fatalf("git init --separate-git-dir: %v: %s", err, output)
	}
}

// symlinkDotGit replaces worktree's .git file with a symlink to target.
func symlinkDotGit(t *testing.T, worktree, target string) {
	t.Helper()
	dotGit := filepath.Join(worktree, ".git")
	if err := os.Remove(dotGit); err != nil {
		t.Fatalf("remove .git file: %v", err)
	}
	if err := os.Symlink(target, dotGit); err != nil {
		t.Fatalf("Symlink(.git) error = %v", err)
	}
}

func discoveryCommit(t *testing.T, repository string) {
	t.Helper()
	discoveryGit(t, repository, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false",
		"commit", "-q", "--no-verify", "--allow-empty", "-m", "base")
}

func discoveryGitCommonDir(t *testing.T, start string) string {
	t.Helper()
	return discoveryGitText(t, start, "rev-parse", "--path-format=absolute", "--git-common-dir")
}

func discoveryGitText(t *testing.T, start string, arguments ...string) string {
	t.Helper()
	output, err := exec.Command("git", append([]string{"-C", start}, arguments...)...).Output()
	if err != nil {
		t.Fatalf("git %v in %q: %v", arguments, start, err)
	}
	return strings.TrimSuffix(string(output), "\n")
}
