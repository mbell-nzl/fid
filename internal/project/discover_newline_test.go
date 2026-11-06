package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// Git prints the top level and the common directory one per line, so a
// newline in either path leaves more than one place to split its report.
// These tests put a newline followed by a slash in the paths, the hardest
// case: each piece of a wrong split can look like an absolute path.

// Discovery reads git's report the one way that gives a top level holding
// the start directory, so repositories whose paths hold a newline work.
func TestDiscoverWorksWhereAPathHoldsANewline(t *testing.T) {
	t.Parallel()
	for _, layout := range []struct {
		name  string
		build func(t *testing.T) (start, root string)
	}{
		{"main worktree", func(t *testing.T) (string, string) {
			repository := newlineRepository(t, "repository\nname")
			return repository, repository
		}},
		{"subdirectory", func(t *testing.T) (string, string) {
			repository := newlineRepository(t, "repository")
			nested := filepath.Join(repository, "sub\n", "directory")
			if err := os.MkdirAll(nested, 0o755); err != nil {
				t.Fatalf("MkdirAll(nested) error = %v", err)
			}
			return nested, repository
		}},
		{"linked worktree", func(t *testing.T) (string, string) {
			repository := newlineRepository(t, "main")
			discoveryCommit(t, repository)
			linked := filepath.Join(newlineDirectory(t, "linked"), "worktree")
			discoveryGit(t, repository, "worktree", "add", "-q", "--detach", linked)
			return linked, linked
		}},
		{"separate git directory outside the worktree", func(t *testing.T) (string, string) {
			worktree := newlineDirectory(t, "worktree")
			separateGitDir(t, worktree, filepath.Join(newlineDirectory(t, "git"), "repository.git"))
			return worktree, worktree
		}},
		{"core.worktree naming a worktree whose .git leads back", func(t *testing.T) (string, string) {
			// No reading's top level holds the start directory here, so the
			// reading whose ledger LedgerDir agrees with is the one.
			return coreWorktreeLeadingBack(t, newlineRepository(t, "repository"), newlineDirectory(t, "worktree"))
		}},
	} {
		t.Run(layout.name, func(t *testing.T) {
			t.Parallel()
			start, root := layout.build(t)
			if got, err := Discover(start); err != nil || got != root {
				t.Fatalf("Discover(%q) = %q, %v; want top level %q", start, got, err, root)
			}
		})
	}
}

// Both refusals still fire where a path holds a newline.
func TestDiscoverRefusesWhereAPathHoldsANewline(t *testing.T) {
	t.Parallel()
	t.Run("git could track the ledger", func(t *testing.T) {
		t.Parallel()
		worktree := newlineDirectory(t, "worktree")
		common := filepath.Join(worktree, "meta\ndata")
		separateGitDir(t, worktree, common)
		_, err := Discover(worktree)
		assertUnsafeLedger(t, err, filepath.Join(common, config.LedgerDirName), worktree, common)
	})
	t.Run("core.worktree names another worktree", func(t *testing.T) {
		t.Parallel()
		worktree := newlineRepository(t, "worktree")
		other := newlineRepository(t, "other")
		discoveryGit(t, other, "config", "core.worktree", worktree)
		_, err := Discover(other)
		assertUnsafeLedger(t, err, worktree, filepath.Join(other, ".git"))
	})
}

func TestDiscoverRefusesWhereTheEnvironmentNamesAnotherLedgerInANewlinePath(t *testing.T) { //nolint:paralleltest // It sets GIT_DIR and GIT_COMMON_DIR.
	worktree := newlineRepository(t, "worktree")
	otherCommon := filepath.Join(newlineRepository(t, "other"), ".git")
	for _, variable := range []string{"GIT_DIR", "GIT_COMMON_DIR"} {
		t.Run(variable, func(t *testing.T) {
			t.Setenv(variable, otherCommon)
			_, err := Discover(worktree)
			assertUnsafeLedger(t, err, filepath.Join(otherCommon, config.LedgerDirName), worktree, config.LedgerDir(worktree))
		})
	}
}

// Where git's report splits into two readings that each name a ledger
// LedgerDir agrees with, and neither top level holds the start directory,
// fid can't tell which ledger git names, so it refuses.
func TestReadGitReportRefusesTwoReadingsThatBothAgree(t *testing.T) {
	t.Parallel()
	base := gittest.ResolvedPath(t, t.TempDir())
	a, b, c := filepath.Join(base, "a"), filepath.Join(base, "b"), filepath.Join(base, "c")
	// One reading: top level a, common directory b\nc. The other: top level
	// a\nb, common directory c. Each top level's .git leads to its common
	// directory.
	for _, reading := range [][2]string{{a, b + "\n" + c}, {a + "\n" + b, c}} {
		for _, directory := range reading {
			if err := os.MkdirAll(directory, 0o755); err != nil {
				t.Fatalf("MkdirAll(%q) error = %v", directory, err)
			}
		}
		if err := os.Symlink(reading[1], filepath.Join(reading[0], ".git")); err != nil {
			t.Fatalf("Symlink(.git) error = %v", err)
		}
	}
	start := filepath.Join(base, "start")
	inside := filepath.Join(a, "inside")
	for _, directory := range []string{start, inside} {
		if err := os.Mkdir(directory, 0o755); err != nil {
			t.Fatalf("Mkdir(%q) error = %v", directory, err)
		}
	}
	_, _, err := readGitReport(a+"\n"+b+"\n"+c, start)
	assertUnsafeLedger(t, err, a+"\n"+b+"\n"+c)
	// From inside a, only the first reading's top level holds the start
	// directory, so that reading wins.
	if root, common, err := readGitReport(a+"\n"+b+"\n"+c, inside); err != nil || root != a || common != b+"\n"+c {
		t.Fatalf("readGitReport from inside a = %q, %q, %v; want %q, %q", root, common, err, a, b+"\n"+c)
	}
}

// coreWorktreeLeadingBack points repository's core.worktree at worktree,
// whose .git file names repository's git directory, so git and LedgerDir
// agree. Run from repository, which is outside worktree, git reports
// worktree as the top level.
func coreWorktreeLeadingBack(t *testing.T, repository, worktree string) (start, root string) {
	t.Helper()
	discoveryGit(t, repository, "config", "core.worktree", worktree)
	if err := os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: "+filepath.Join(repository, ".git")+"\n"), 0o644); err != nil {
		t.Fatalf("write the .git file: %v", err)
	}
	return repository, worktree
}

// newlineDirectory returns a fresh directory, with symlinks resolved, whose
// path holds a newline followed by a slash.
func newlineDirectory(t *testing.T, name string) string {
	t.Helper()
	directory := filepath.Join(gittest.ResolvedPath(t, t.TempDir()), "line\n", name)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", directory, err)
	}
	return directory
}

func newlineRepository(t *testing.T, name string) string {
	t.Helper()
	repository := newlineDirectory(t, name)
	if err := gittest.InitRepository(repository); err != nil {
		t.Fatalf("make %q a repository: %v", repository, err)
	}
	return repository
}
