package gittest_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

// guardChildEnv marks the child process that TestGitDiscoveryStopsAtTheTempDirectory starts.
const guardChildEnv = "FID_GIT_GUARD_CHILD"

func TestMain(m *testing.M) {
	os.Exit(gittest.Main(m))
}

func TestUnbornHeadIsDiscoverableButNotObservable(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	git := func(args ...string) (string, error) {
		output, err := exec.Command("git", append([]string{"-C", workspace}, args...)...).Output()
		return strings.TrimSpace(string(output)), err
	}
	if _, err := git("-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-q", "--allow-empty", "-m", "base"); err != nil {
		t.Fatalf("commit: %v", err)
	}
	base, err := git("rev-parse", "HEAD")
	if err != nil {
		t.Fatalf("rev-parse HEAD: %v", err)
	}
	restore := gittest.UnbornHead(t, workspace)
	if err := gittest.CheckRepository(workspace); err != nil {
		t.Fatalf("unborn HEAD hid the repository: %v", err)
	}
	if head, err := git("rev-parse", "--verify", "-q", "HEAD"); err == nil {
		t.Fatalf("unborn HEAD resolves to %q, want no commit", head)
	}
	restore()
	if head, err := git("rev-parse", "HEAD"); err != nil || head != base {
		t.Fatalf("restored HEAD = %q, %v; want %q", head, err, base)
	}
}

func TestAddWorktreeSharesTheRepositorysCommonDirectory(t *testing.T) {
	t.Parallel()
	repository := gittest.NewWorkspace(t)
	worktree := gittest.AddWorktree(t, repository)
	commonDir := func(directory string) string {
		t.Helper()
		output, err := exec.Command("git", "-C", directory, "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
		if err != nil {
			t.Fatalf("git rev-parse --git-common-dir in %q: %v", directory, err)
		}
		return strings.TrimSpace(string(output))
	}
	if got, want := commonDir(worktree), commonDir(repository); got != want {
		t.Fatalf("linked worktree's common directory = %q, want the repository's %q", got, want)
	}
	if want := filepath.Join(gittest.ResolvedPath(t, repository), ".git"); commonDir(repository) != want {
		t.Fatalf("common directory = %q, want %q", commonDir(repository), want)
	}
	if gittest.ResolvedPath(t, worktree) == gittest.ResolvedPath(t, repository) {
		t.Fatalf("linked worktree %q is the repository itself", worktree)
	}
}

func TestGuardedGitWorkspaceIsItsOwnRepository(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	if err := gittest.CheckRepository(workspace); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(workspace, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatalf("Mkdir(nested) error = %v", err)
	}
	if err := gittest.CheckRepository(nested); err == nil {
		t.Fatal("check of a subdirectory = nil, want refusal because git reports its parent")
	}
}

// TestGitDiscoveryStopsAtTheTempDirectory runs this test binary again with its
// temp directory inside another repository, so the guard TestMain installs is
// what stops git. It covers a plain temp directory and a symlinked one that
// is also inherited after an empty ceiling entry.
func TestGitDiscoveryStopsAtTheTempDirectory(t *testing.T) {
	t.Parallel()
	if os.Getenv(guardChildEnv) != "" {
		directory := t.TempDir()
		if output, err := exec.Command("git", "-C", directory, "rev-parse", "--show-toplevel").Output(); err == nil {
			t.Fatalf("git found repository %q from temp directory %q, want ceiling refusal", strings.TrimSpace(string(output)), directory)
		}
		return
	}
	ancestor := gittest.NewWorkspace(t)
	temp := filepath.Join(ancestor, "tmp")
	if err := os.Mkdir(temp, 0o755); err != nil {
		t.Fatalf("Mkdir(temp) error = %v", err)
	}
	alias := filepath.Join(t.TempDir(), "tmp-alias")
	if err := os.Symlink(temp, alias); err != nil {
		t.Fatalf("Symlink(temp) error = %v", err)
	}
	control := exec.Command("git", "-C", temp, "rev-parse", "--show-toplevel")
	control.Env = environmentWithout(os.Environ(), gittest.CeilingDirectoriesEnv)
	wantAncestor, err := filepath.EvalSymlinks(ancestor)
	if err != nil {
		t.Fatalf("EvalSymlinks(ancestor) error = %v", err)
	}
	if output, err := control.Output(); err != nil || strings.TrimSpace(string(output)) != wantAncestor {
		t.Fatalf("control without a ceiling = %q, %v; want git to find %q", output, err, wantAncestor)
	}
	for _, child := range []struct{ name, tempDir, inherited string }{
		{name: "plain temp directory", tempDir: temp},
		{name: "symlinked temp directory after an inherited empty entry", tempDir: alias, inherited: string(os.PathListSeparator) + alias},
	} {
		command := exec.Command(os.Args[0], "-test.run=^TestGitDiscoveryStopsAtTheTempDirectory$", "-test.v")
		command.Env = append(environmentWithout(os.Environ(), gittest.CeilingDirectoriesEnv, "TMPDIR"), "TMPDIR="+child.tempDir, guardChildEnv+"=1")
		if child.inherited != "" {
			command.Env = append(command.Env, gittest.CeilingDirectoriesEnv+"="+child.inherited)
		}
		output, err := command.CombinedOutput()
		if err != nil || !strings.Contains(string(output), "--- PASS: TestGitDiscoveryStopsAtTheTempDirectory") {
			t.Fatalf("%s: child = %v, want git stopped at the temp directory\n%s", child.name, err, output)
		}
	}
}

func environmentWithout(environment []string, names ...string) []string {
	result := make([]string, 0, len(environment))
	for _, entry := range environment {
		name, _, _ := strings.Cut(entry, "=")
		if !slices.Contains(names, name) {
			result = append(result, entry)
		}
	}
	return result
}

// Every test process that makes repositories must guard git discovery first,
// so a package whose TestMain forgets fails at its first workspace.
func TestInitRepositoryRefusesUnlessGitDiscoveryIsGuarded(t *testing.T) { //nolint:paralleltest // changes GIT_CEILING_DIRECTORIES for the process.
	temp := os.TempDir()
	module := gittest.ModuleRoot()
	separator := string(os.PathListSeparator)
	for _, test := range []struct{ name, ceilings string }{
		{name: "unset"},
		{name: "temp directory only", ceilings: temp},
		{name: "module root only", ceilings: module},
		{name: "wrong order", ceilings: module + separator + temp},
		{name: "after another entry", ceilings: "/elsewhere" + separator + temp + separator + module},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv(gittest.CeilingDirectoriesEnv, test.ceilings)
			directory := t.TempDir()
			err := gittest.InitRepository(directory)
			if err == nil || !strings.Contains(err.Error(), "not guarded") {
				t.Fatalf("InitRepository() with %s = %q: error = %v, want a refusal naming the missing guard", gittest.CeilingDirectoriesEnv, test.ceilings, err)
			}
			if _, statErr := os.Lstat(filepath.Join(directory, ".git")); !os.IsNotExist(statErr) {
				t.Fatalf("refused InitRepository() left .git behind: %v", statErr)
			}
		})
	}
	t.Run("guarded with inherited entries after ours", func(t *testing.T) {
		t.Setenv(gittest.CeilingDirectoriesEnv, temp+separator+module+separator+"/elsewhere")
		if err := gittest.InitRepository(t.TempDir()); err != nil {
			t.Fatalf("InitRepository() = %v, want success", err)
		}
	})
}

// gittest imports testing and runs git, so it must stay out of fid itself.
func TestOnlyTestsImportGittest(t *testing.T) {
	t.Parallel()
	const importPath = "github.com/mbell-nzl/fid/internal/gittest"
	own, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("Abs(.) error = %v", err)
	}
	var importers []string
	for _, tree := range []string{"cmd", "internal"} {
		err := filepath.WalkDir(filepath.Join(gittest.ModuleRoot(), tree), func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() && (entry.Name() == "testdata" || path == own) {
				return filepath.SkipDir
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
			if err != nil {
				return err
			}
			for _, spec := range file.Imports {
				if spec.Path.Value == strconv.Quote(importPath) {
					importers = append(importers, path)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", tree, err)
		}
	}
	if len(importers) != 0 {
		t.Errorf("non-test files import %s:\n%s", importPath, strings.Join(importers, "\n"))
	}
}

func TestGitDiscoveryStopsAtTheModuleRoot(t *testing.T) {
	t.Parallel()
	packageDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	command := exec.Command("git", "rev-parse", "--show-toplevel")
	command.Dir = packageDirectory
	if output, err := command.Output(); err == nil {
		t.Fatalf("git found repository %q from package directory %q, want ceiling refusal", strings.TrimSpace(string(output)), packageDirectory)
	}
}
