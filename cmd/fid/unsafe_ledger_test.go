package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

const unsafeLedgerMessage = "the ledger location is unsafe"

// TestEveryProjectCommandRefusesWhereGitNamesAnotherLedger runs every project
// command where git and LedgerDir name different ledgers. With core.worktree,
// git reports another repository's worktree as the top level but keeps its
// own common directory, while LedgerDir reads that worktree's .git. Each
// command refuses, and no file in either repository changes.
func TestEveryProjectCommandRefusesWhereGitNamesAnotherLedger(t *testing.T) {
	t.Parallel()
	worktree := clonePristineWorkspace(t)
	start := pointCoreWorktreeAt(t, worktree)
	resolved := gittest.ResolvedPath(t, worktree)
	assertEveryProjectCommandRefuses(t, start, repositoryRefusal{
		message: unsafeLedgerMessage, reason: "unsafe-path",
		names: []string{filepath.Join(gittest.ResolvedPath(t, start), ".git", config.LedgerDirName), resolved, config.LedgerDir(resolved)},
	}, worktree, start)
}

// TestEveryProjectCommandRefusesWhereGitCouldTrackTheLedger runs every
// project command in a worktree whose git directory, and so the ledger, is an
// ordinary directory inside it, as `git init --separate-git-dir` makes it.
// Each command refuses, and no file in the worktree changes.
func TestEveryProjectCommandRefusesWhereGitCouldTrackTheLedger(t *testing.T) {
	t.Parallel()
	worktree := clonePristineWorkspace(t)
	common := filepath.Join(worktree, "metadata")
	if err := os.Rename(filepath.Join(worktree, ".git"), common); err != nil {
		t.Fatalf("move the git directory into the worktree: %v", err)
	}
	if err := os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: metadata\n"), 0o644); err != nil {
		t.Fatalf("write the .git file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(config.LedgerDir(worktree), config.ConfigFileName)); err != nil {
		t.Fatalf("the ledger LedgerDir names: %v", err)
	}
	resolved := gittest.ResolvedPath(t, worktree)
	assertEveryProjectCommandRefuses(t, worktree, repositoryRefusal{
		message: unsafeLedgerMessage, reason: "unsafe-path",
		names: []string{filepath.Join(resolved, "metadata", config.LedgerDirName), resolved},
	}, worktree)
}

// fid refuses before it stages anything: init and import, each in a fresh
// layout, create no ledger or staging directory anywhere. With the git
// directory inside the worktree, `git add --all` then tracks git's own files
// but nothing of fid's.
func TestRefusedInitAndImportLeaveNoLedgerForGitToTrack(t *testing.T) {
	t.Parallel()
	archive := filepath.Join(t.TempDir(), "fid-export.zip")
	if created := executeCommandForTest(clonePristineWorkspace(t), "--json", "export", "create", archive); created.err != nil {
		t.Fatalf("export create = %#v", created)
	}
	for _, command := range []string{"init", "export import"} {
		t.Run(command+" where git could track the ledger", func(t *testing.T) {
			t.Parallel()
			worktree := gittest.ResolvedPath(t, t.TempDir())
			common := filepath.Join(worktree, "metadata")
			if output, err := exec.Command("git", "init", "-q", "--template=", "--separate-git-dir="+common, worktree).CombinedOutput(); err != nil {
				t.Fatalf("git init --separate-git-dir: %v: %s", err, output)
			}
			assertRefusedBeforeStaging(t, worktree, unsafeLedgerCommand(command, worktree, archive), common, worktree)
			unsafeLedgerGit(t, worktree, "add", "--all")
			tracked := strings.Fields(unsafeLedgerGit(t, worktree, "ls-files"))
			if !slices.Contains(tracked, "metadata/HEAD") {
				t.Fatalf("git ls-files = %v, want git's own metadata/HEAD tracked, as this layout allows", tracked)
			}
			for _, path := range tracked {
				if strings.HasPrefix(path, "metadata/"+config.LedgerDirName) {
					t.Errorf("git tracked fid's %s", path)
				}
			}
		})
		t.Run(command+" where git names another ledger", func(t *testing.T) {
			t.Parallel()
			worktree := gittest.NewWorkspace(t)
			start := pointCoreWorktreeAt(t, worktree)
			assertRefusedBeforeStaging(t, start, unsafeLedgerCommand(command, start, archive), filepath.Join(start, ".git"), filepath.Join(worktree, ".git"))
		})
	}
}

// TestTheEnvironmentNamingAnotherLedgerIsRefused covers the main way git and
// LedgerDir disagree: GIT_DIR or GIT_COMMON_DIR naming another repository.
// Naming the worktree's own git directory, as git does for hooks, works.
func TestTheEnvironmentNamingAnotherLedgerIsRefused(t *testing.T) { //nolint:paralleltest // It sets GIT_DIR and GIT_COMMON_DIR.
	worktree := clonePristineWorkspace(t)
	other := gittest.ResolvedPath(t, gittest.NewWorkspace(t))
	archive := filepath.Join(t.TempDir(), "fid-export.zip")
	if created := executeCommandForTest(worktree, "--json", "export", "create", archive); created.err != nil {
		t.Fatalf("export create = %#v", created)
	}
	for _, variable := range []string{"GIT_DIR", "GIT_COMMON_DIR"} {
		t.Run(variable, func(t *testing.T) {
			fresh := map[string]string{"init": gittest.NewWorkspace(t), "export import": gittest.NewWorkspace(t)}
			before := snapshotWorkspace(t, worktree)
			t.Setenv(variable, filepath.Join(other, ".git"))
			result := grammarExecute(worktree, "--json", "issue", "list")
			if cli.ExitCode(result.err) != 4 || result.stdout != "" {
				t.Fatalf("issue list with %s naming another repository = %#v, want refusal with exit 4", variable, result)
			}
			resolved := gittest.ResolvedPath(t, worktree)
			assertRepositoryRefusal(t, "issue.list", repositoryRefusal{message: unsafeLedgerMessage, reason: "unsafe-path",
				names: []string{filepath.Join(other, ".git", config.LedgerDirName), config.LedgerDir(resolved)}}, result.stderr)
			for command, start := range fresh {
				assertRefusedBeforeStaging(t, start, unsafeLedgerCommand(command, start, archive), filepath.Join(other, ".git"), filepath.Join(start, ".git"))
			}
			if !reflect.DeepEqual(before, snapshotWorkspace(t, worktree)) {
				t.Errorf("a refused command changed %s", worktree)
			}
		})
	}
	t.Run("GIT_DIR naming the worktree's own git directory", func(t *testing.T) {
		t.Setenv("GIT_DIR", filepath.Join(worktree, ".git"))
		if result := grammarExecute(worktree, "--json", "issue", "list"); result.err != nil {
			t.Fatalf("issue list with GIT_DIR naming its own git directory = %#v", result)
		}
	})
}

// Layouts whose common directory is within a .git directory, or outside the
// worktree, still work. A worktree made inside the main worktree shares its
// ledger, and a submodule keeps its own ledger in the superproject's
// .git/modules, where git never tracks it.
func TestWorktreesInsideTheMainWorktreeAndSubmodulesStillWork(t *testing.T) {
	t.Parallel()
	t.Run("worktree inside the main worktree", func(t *testing.T) {
		t.Parallel()
		worktree := clonePristineWorkspace(t)
		unsafeLedgerGit(t, worktree, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false",
			"commit", "-q", "--no-verify", "--allow-empty", "-m", "base")
		inner := filepath.Join(worktree, "inner")
		unsafeLedgerGit(t, worktree, "worktree", "add", "-q", "--detach", inner)
		created := executeCommandForTest(inner, "--json", "issue", "create", "--kind", "task", "--title", "From the inner worktree", "--actor", "agent:developer")
		if created.err != nil {
			t.Fatalf("issue create from a worktree inside the main worktree = %#v", created)
		}
		if listed := executeCommandForTest(worktree, "--json", "issue", "list"); listed.err != nil || !strings.Contains(listed.stdout, "From the inner worktree") {
			t.Fatalf("issue list from the main worktree = %#v, want the issue created from the inner worktree", listed)
		}
	})
	t.Run("submodule", func(t *testing.T) {
		t.Parallel()
		source := gittest.NewWorkspace(t)
		unsafeLedgerGit(t, source, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false",
			"commit", "-q", "--no-verify", "--allow-empty", "-m", "base")
		superproject := gittest.NewWorkspace(t)
		unsafeLedgerGit(t, superproject, "-c", "protocol.file.allow=always", "submodule", "add", "-q", source, "module")
		initialized := grammarExecute(filepath.Join(superproject, "module"), "--json", "init", "--prefix", "MOD")
		if initialized.err != nil {
			t.Fatalf("init in a submodule = %#v", initialized)
		}
		ledger := filepath.Join(gittest.ResolvedPath(t, superproject), ".git", "modules", "module", config.LedgerDirName)
		if journal, _ := decodeCanonicalResult(t, initialized.stdout)["journal_path"].(string); !strings.HasPrefix(journal, ledger+string(filepath.Separator)) {
			t.Fatalf("init in a submodule journal_path = %q, want it in %s", journal, ledger)
		}
		unsafeLedgerGit(t, superproject, "add", "--all")
		if tracked := strings.Fields(unsafeLedgerGit(t, superproject, "ls-files")); !slices.Equal(tracked, []string{".gitmodules", "module"}) {
			t.Fatalf("superproject git ls-files = %v, want only .gitmodules and module", tracked)
		}
	})
}

// A repository whose path holds a newline works: discovery asks git for each
// path in its own run. Where git could track the ledger, init still refuses,
// and git tracks nothing of fid's.
func TestARepositoryWhosePathHoldsANewlineWorks(t *testing.T) {
	t.Parallel()
	base := gittest.ResolvedPath(t, t.TempDir())
	worktree := filepath.Join(base, "line\n", "repository\nname")
	nested := filepath.Join(worktree, "sub\n", "directory")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("MkdirAll(nested) error = %v", err)
	}
	if err := gittest.InitRepository(worktree); err != nil {
		t.Fatalf("make the worktree a repository: %v", err)
	}
	initialized := grammarExecute(nested, "--json", "init", "--prefix", "FID")
	if initialized.err != nil {
		t.Fatalf("init where the path holds a newline = %#v", initialized)
	}
	if root := decodeCanonicalResult(t, initialized.stdout)["project"].(map[string]any)["root"]; root != worktree {
		t.Fatalf("init project.root = %v, want %q", root, worktree)
	}
	if listed := executeCommandForTest(nested, "--json", "issue", "list"); listed.err != nil || !strings.Contains(listed.stdout, `"operation":"issue.list"`) {
		t.Fatalf("issue list where the path holds a newline = %#v", listed)
	}

	tracked := filepath.Join(base, "line\n", "tracked")
	common := filepath.Join(tracked, "meta\ndata")
	if err := os.MkdirAll(tracked, 0o755); err != nil {
		t.Fatalf("MkdirAll(tracked) error = %v", err)
	}
	if output, err := exec.Command("git", "init", "-q", "--template=", "--separate-git-dir="+common, tracked).CombinedOutput(); err != nil {
		t.Fatalf("git init --separate-git-dir: %v: %s", err, output)
	}
	assertRefusedBeforeStaging(t, tracked, unsafeLedgerCommand("init", tracked, ""), common, tracked)
	unsafeLedgerGit(t, tracked, "add", "--all")
	trackedPaths := strings.Split(unsafeLedgerGit(t, tracked, "ls-files", "-z"), "\x00")
	if !slices.Contains(trackedPaths, "meta\ndata/HEAD") {
		t.Fatalf("git ls-files = %q, want git's own meta\\ndata/HEAD tracked, as this layout allows", trackedPaths)
	}
	for _, path := range trackedPaths {
		if strings.HasPrefix(path, "meta\ndata/"+config.LedgerDirName) {
			t.Errorf("git tracked fid's %q", path)
		}
	}
}

// A single git run would print the top level and the common directory one
// per line. In these layouts that report also splits at an earlier newline
// into another top level and common directory, which LedgerDir agrees with.
// fid asks git for each path in its own run, so it refuses where the paths
// git means disagree, works where they agree, and never writes into the
// misreading's ledger.
func TestAMisreadingOfGitsReportCannotMoveTheLedger(t *testing.T) {
	t.Parallel()
	archive := filepath.Join(t.TempDir(), "fid-export.zip")
	if created := executeCommandForTest(clonePristineWorkspace(t), "--json", "export", "create", archive); created.err != nil {
		t.Fatalf("export create = %#v", created)
	}
	for _, command := range []string{"init", "export import"} {
		t.Run(command+" where only the misreading agrees", func(t *testing.T) {
			t.Parallel()
			layout := newMisreadableLayout(t, false)
			assertRefusedLeavingNoLedger(t, layout, command, archive,
				filepath.Join(layout.common, config.LedgerDirName), layout.root, config.LedgerDir(layout.root))
		})
	}
	t.Run("init where the misreading agrees too", func(t *testing.T) {
		t.Parallel()
		layout := newMisreadableLayout(t, true)
		initialized := grammarExecute(layout.start, "--json", "init", "--prefix", "FID")
		if initialized.err != nil {
			t.Fatalf("init = %#v", initialized)
		}
		result := decodeCanonicalResult(t, initialized.stdout)
		// journal_path names the ledger through the root's .git, a symlink.
		ledger := filepath.Join(layout.common, config.LedgerDirName)
		if journal, _ := result["journal_path"].(string); !strings.HasPrefix(gittest.ResolvedPath(t, journal), ledger+string(filepath.Separator)) {
			t.Fatalf("init journal_path = %q, want it in git's common directory's ledger %s", journal, ledger)
		}
		if root := result["project"].(map[string]any)["root"]; root != layout.root {
			t.Fatalf("init project.root = %q, want %q", root, layout.root)
		}
		if listed := executeCommandForTest(layout.start, "--json", "issue", "list"); listed.err != nil {
			t.Fatalf("issue list after init = %#v", listed)
		}
		if _, err := os.Lstat(filepath.Join(layout.misreadCommon, config.LedgerDirName)); !os.IsNotExist(err) {
			t.Fatalf("Lstat(the misreading's ledger) error = %v, want it not to exist", err)
		}
	})
}

// The same with GIT_DIR and GIT_WORK_TREE, where only the misreading's top
// level holds the start directory.
func TestAMisreadingOfGitsReportUnderGitDirAndGitWorkTreeCannotMoveTheLedger(t *testing.T) { //nolint:paralleltest // It sets GIT_DIR and GIT_WORK_TREE.
	archive := filepath.Join(t.TempDir(), "fid-export.zip")
	if created := executeCommandForTest(clonePristineWorkspace(t), "--json", "export", "create", archive); created.err != nil {
		t.Fatalf("export create = %#v", created)
	}
	for _, command := range []string{"init", "export import"} {
		t.Run(command, func(t *testing.T) {
			layout := newMisreadableEnvironment(t)
			t.Setenv("GIT_DIR", layout.common)
			t.Setenv("GIT_WORK_TREE", layout.root)
			assertRefusedLeavingNoLedger(t, layout, command, archive,
				filepath.Join(layout.common, config.LedgerDirName), layout.root, config.LedgerDir(layout.root))
		})
	}
}

// assertRefusedLeavingNoLedger runs command from the layout's start. It must
// refuse with unsafe-path, naming every location in names, and leave no
// ledger or staging directory anywhere in the layout, the misreading's
// common directory included.
func assertRefusedLeavingNoLedger(t *testing.T, layout misreadableLayout, command, archive string, names ...string) {
	t.Helper()
	args := unsafeLedgerCommand(command, layout.start, archive)
	result := grammarExecute(layout.start, args...)
	if cli.ExitCode(result.err) != 4 || result.stdout != "" {
		t.Errorf("fid %s = %#v, want the unsafe-path refusal", strings.Join(args, " "), result)
	}
	operation := map[string]string{"init": "project.init", "export import": "export.import"}[command]
	assertRepositoryRefusal(t, operation, repositoryRefusal{message: unsafeLedgerMessage, reason: "unsafe-path", names: names}, result.stderr)
	err := filepath.WalkDir(layout.base, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Name() == config.LedgerDirName || strings.HasPrefix(entry.Name(), config.LedgerStagingPrefix) {
			t.Errorf("fid %s left %q", command, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %q: %v", layout.base, err)
	}
}

// unsafeLedgerCommand returns the arguments for init, or for an import of
// archive into start.
func unsafeLedgerCommand(command, start, archive string) []string {
	if command == "init" {
		return []string{"--json", "init", "--prefix", "FID"}
	}
	return []string{"--json", "export", "import", archive, "--repository", start, "--idempotency-key", "unsafe-ledger"}
}

// assertRefusedBeforeStaging runs fid from start. It must refuse with
// unsafe-path and leave no ledger or staging directory in any of commonDirs.
func assertRefusedBeforeStaging(t *testing.T, start string, args []string, commonDirs ...string) {
	t.Helper()
	result := grammarExecute(start, args...)
	if cli.ExitCode(result.err) != 4 || result.stdout != "" || !strings.Contains(result.stderr, `"reason":"unsafe-path"`) {
		t.Errorf("fid %s = %#v, want the unsafe-path refusal", strings.Join(args, " "), result)
	}
	for _, common := range commonDirs {
		entries, err := os.ReadDir(common)
		if err != nil {
			t.Fatalf("ReadDir(%q) error = %v", common, err)
		}
		for _, entry := range entries {
			if entry.Name() == config.LedgerDirName || strings.HasPrefix(entry.Name(), config.LedgerStagingPrefix) {
				t.Errorf("fid %s left %s", args[1], filepath.Join(common, entry.Name()))
			}
		}
	}
}

// pointCoreWorktreeAt makes a new repository whose core.worktree is worktree,
// and returns it. Run from there, git reports worktree as the top level but
// the new repository's .git as the common directory.
func pointCoreWorktreeAt(t *testing.T, worktree string) string {
	t.Helper()
	repository := gittest.NewWorkspace(t)
	unsafeLedgerGit(t, repository, "config", "core.worktree", worktree)
	return repository
}

func unsafeLedgerGit(t *testing.T, directory string, args ...string) string {
	t.Helper()
	output, err := exec.Command("git", append([]string{"-C", directory}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v in %q: %v: %s", args, directory, err, output)
	}
	return string(output)
}

// misreadableLayout is a layout where the top level and common directory git
// means, root and common, joined by a newline as a single git run would print
// them, also split at an earlier newline into another top level and common
// directory, misreadRoot and misreadCommon, that LedgerDir agrees with. Every
// path is under base.
type misreadableLayout struct {
	base, start, root, common, misreadRoot, misreadCommon string
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
	layout := misreadableLayout{base: base, start: s, root: a + "\n" + b, common: filepath.Join(s, ".git"), misreadRoot: a}
	layout.misreadCommon = b + "\n" + layout.common
	for _, directory := range []string{a, s, layout.root} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", directory, err)
		}
	}
	if err := gittest.InitRepository(s); err != nil {
		t.Fatalf("make %q a repository: %v", s, err)
	}
	if rootLeadsBack {
		if err := os.Symlink(layout.common, filepath.Join(layout.root, ".git")); err != nil {
			t.Fatalf("Symlink(.git) error = %v", err)
		}
	} else if err := gittest.InitRepository(layout.root); err != nil {
		t.Fatalf("make %q a repository: %v", layout.root, err)
	}
	unsafeLedgerGit(t, s, "config", "core.worktree", layout.root)
	unsafeLedgerGit(t, base, "init", "-q", "--bare", "--template=", layout.misreadCommon)
	if err := os.Symlink(layout.misreadCommon, filepath.Join(a, ".git")); err != nil {
		t.Fatalf("Symlink(.git) error = %v", err)
	}
	if root := strings.TrimSuffix(unsafeLedgerGit(t, s, "rev-parse", "--show-toplevel"), "\n"); root != layout.root {
		t.Fatalf("git from %q means top level %q, want %q", s, root, layout.root)
	}
	if common := strings.TrimSuffix(unsafeLedgerGit(t, s, "rev-parse", "--path-format=absolute", "--git-common-dir"), "\n"); common != layout.common {
		t.Fatalf("git from %q means common directory %q, want %q", s, common, layout.common)
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
	layout := misreadableLayout{base: base, start: filepath.Join(a, "start"), root: a + "\n" + b, common: g, misreadRoot: a, misreadCommon: b + "\n" + g}
	for _, directory := range []string{layout.start, layout.root} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", directory, err)
		}
	}
	for _, repository := range []string{g, layout.misreadCommon} {
		unsafeLedgerGit(t, base, "init", "-q", "--bare", "--template=", repository)
	}
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
