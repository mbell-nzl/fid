package project

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// Discovery asks git for the top level and the common directory in separate
// runs, so a newline in either path can't make git's report ambiguous. These
// tests put a newline followed by a slash in the paths, the hardest case:
// printed one per line by a single run, the two paths would also split at
// that newline into two other absolute paths.

// Repositories whose paths hold a newline work, including paths that end in
// one: discovery strips only the newline git ends each run's output with.
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
			return coreWorktreeLeadingBack(t, newlineRepository(t, "repository"), newlineDirectory(t, "worktree"))
		}},
		{"top level ending in a newline", func(t *testing.T) (string, string) {
			// gittest.InitRepository's own check trims that newline, so this
			// runs git init directly.
			worktree := newlineDirectory(t, "worktree\n")
			discoveryGit(t, worktree, "init", "-q", "--template=")
			return worktree, worktree
		}},
		{"top level ending in a carriage return", func(t *testing.T) (string, string) {
			worktree := newlineDirectory(t, "worktree\r")
			discoveryGit(t, worktree, "init", "-q", "--template=")
			return worktree, worktree
		}},
		{"common directory ending in a newline", func(t *testing.T) (string, string) {
			// Git trims that newline from a .git file, so .git is a symlink.
			worktree := newlineRepository(t, "worktree")
			common := filepath.Join(newlineDirectory(t, "git"), "repository.git\n")
			if err := os.Rename(filepath.Join(worktree, ".git"), common); err != nil {
				t.Fatalf("move the git directory: %v", err)
			}
			if err := os.Symlink(common, filepath.Join(worktree, ".git")); err != nil {
				t.Fatalf("Symlink(.git) error = %v", err)
			}
			return worktree, worktree
		}},
		{"core.worktree where a misreading of a single run's report agrees too", func(t *testing.T) (string, string) {
			layout := newMisreadableLayout(t, true)
			return layout.start, layout.root
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
	t.Run("core.worktree where only a misreading of a single run's report agrees", func(t *testing.T) {
		t.Parallel()
		layout := newMisreadableLayout(t, false)
		_, err := Discover(layout.start)
		assertUnsafeLedger(t, err, filepath.Join(layout.common, config.LedgerDirName), layout.root, config.LedgerDir(layout.root))
	})
}

func TestDiscoverRefusesWhereTheEnvironmentNamesAnotherLedgerInANewlinePath(t *testing.T) { //nolint:paralleltest // It sets GIT_DIR, GIT_COMMON_DIR and GIT_WORK_TREE.
	worktree := newlineRepository(t, "worktree")
	otherCommon := filepath.Join(newlineRepository(t, "other"), ".git")
	for _, variable := range []string{"GIT_DIR", "GIT_COMMON_DIR"} {
		t.Run(variable, func(t *testing.T) {
			t.Setenv(variable, otherCommon)
			_, err := Discover(worktree)
			assertUnsafeLedger(t, err, filepath.Join(otherCommon, config.LedgerDirName), worktree, config.LedgerDir(worktree))
		})
	}
	t.Run("GIT_DIR and GIT_WORK_TREE where only a misreading's top level holds the start", func(t *testing.T) {
		layout := newMisreadableEnvironment(t)
		t.Setenv("GIT_DIR", layout.common)
		t.Setenv("GIT_WORK_TREE", layout.root)
		_, err := Discover(layout.start)
		assertUnsafeLedger(t, err, filepath.Join(layout.common, config.LedgerDirName), layout.root, config.LedgerDir(layout.root))
	})
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

// misreadableLayout is a layout where the top level and common directory git
// means, root and common, joined by a newline as a single git run would print
// them, also split at an earlier newline into another top level and common
// directory, misreadRoot and misreadCommon, that LedgerDir agrees with.
type misreadableLayout struct {
	start, root, common, misreadRoot, misreadCommon string
}

// newMisreadableLayout builds the layout from #42's second review. Sibling
// directories a, b and s; s is a repository whose core.worktree names a\nb.
// Run from s, git means top level a\nb and common directory s/.git. Joined,
// those read a\nb\ns/.git, which also splits into top level a and common
// directory b\ns/.git, a bare repository that a/.git is a symlink to. Unless
// rootLeadsBack, a\nb is a repository of its own, so LedgerDir disagrees with
// what git means; with it, a\nb/.git is a symlink to s/.git, so LedgerDir
// agrees with both.
func newMisreadableLayout(t *testing.T, rootLeadsBack bool) misreadableLayout {
	t.Helper()
	base := gittest.ResolvedPath(t, t.TempDir())
	a, b, s := filepath.Join(base, "a"), filepath.Join(base, "b"), filepath.Join(base, "s")
	layout := misreadableLayout{start: s, root: a + "\n" + b, common: filepath.Join(s, ".git"), misreadRoot: a}
	layout.misreadCommon = b + "\n" + layout.common
	for _, directory := range []string{a, s, layout.root} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", directory, err)
		}
	}
	for _, repository := range []string{s, layout.root} {
		if repository == layout.root && rootLeadsBack {
			continue
		}
		if err := gittest.InitRepository(repository); err != nil {
			t.Fatalf("make %q a repository: %v", repository, err)
		}
	}
	if rootLeadsBack {
		if err := os.Symlink(layout.common, filepath.Join(layout.root, ".git")); err != nil {
			t.Fatalf("Symlink(.git) error = %v", err)
		}
	}
	discoveryGit(t, s, "config", "core.worktree", layout.root)
	bareRepository(t, layout.misreadCommon)
	if err := os.Symlink(layout.misreadCommon, filepath.Join(a, ".git")); err != nil {
		t.Fatalf("Symlink(.git) error = %v", err)
	}
	if root, common := discoveryGitText(t, s, "rev-parse", "--show-toplevel"), discoveryGitCommonDir(t, s); root != layout.root || common != layout.common {
		t.Fatalf("git from %q means %q and %q, want %q and %q", s, root, common, layout.root, layout.common)
	}
	layout.checkMisreading(t)
	return layout
}

// newMisreadableEnvironment builds a layout for GIT_DIR=common and
// GIT_WORK_TREE=root, run from start. Sibling directories a, b and g; g is a
// bare repository, root is a\nb and start is a/start. Git means top level
// a\nb and common directory g. Joined, those read a\nb\ng, which also splits
// into top level a, which holds the start, and common directory b\ng, a bare
// repository that a/.git is a symlink to.
func newMisreadableEnvironment(t *testing.T) misreadableLayout {
	t.Helper()
	base := gittest.ResolvedPath(t, t.TempDir())
	a, b, g := filepath.Join(base, "a"), filepath.Join(base, "b"), filepath.Join(base, "g")
	layout := misreadableLayout{start: filepath.Join(a, "start"), root: a + "\n" + b, common: g, misreadRoot: a, misreadCommon: b + "\n" + g}
	for _, directory := range []string{layout.start, layout.root} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", directory, err)
		}
	}
	bareRepository(t, g)
	bareRepository(t, layout.misreadCommon)
	if err := os.Symlink(layout.misreadCommon, filepath.Join(a, ".git")); err != nil {
		t.Fatalf("Symlink(.git) error = %v", err)
	}
	layout.checkMisreading(t)
	return layout
}

// checkMisreading fails unless the layout is misreadable: what git means,
// joined by a newline, is also the misreading joined by one, and the
// misreading's top level has a .git leading to its common directory, so
// LedgerDir agrees with it.
func (layout misreadableLayout) checkMisreading(t *testing.T) {
	t.Helper()
	if joined := layout.root + "\n" + layout.common; joined != layout.misreadRoot+"\n"+layout.misreadCommon {
		t.Fatalf("%q doesn't also read as top level %q and common directory %q", joined, layout.misreadRoot, layout.misreadCommon)
	}
	if dotGit := gittest.ResolvedPath(t, filepath.Join(layout.misreadRoot, ".git")); dotGit != layout.misreadCommon {
		t.Fatalf("%s/.git leads to %q, want %q", layout.misreadRoot, dotGit, layout.misreadCommon)
	}
}

func bareRepository(t *testing.T, directory string) {
	t.Helper()
	if output, err := exec.Command("git", "init", "-q", "--bare", "--template=", directory).CombinedOutput(); err != nil {
		t.Fatalf("git init --bare %q: %v: %s", directory, err, output)
	}
}
