// Package gittest makes the guarded git repositories that fid's tests use as
// workspaces. Only tests import it.
//
// A test process calls Guard, or Main from its TestMain, before it makes a
// repository, so git can't climb from a temp directory or a package
// directory into a real repository.
package gittest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

// CeilingDirectoriesEnv names the variable that stops git climbing above a
// directory while it looks for a repository.
const CeilingDirectoriesEnv = "GIT_CEILING_DIRECTORIES"

// ModuleRoot returns the root of fid's module, found from this file's path.
func ModuleRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// Main guards git discovery and then runs the tests. A package's TestMain
// passes its result to os.Exit.
func Main(m interface{ Run() int }) int {
	if err := Guard(); err != nil {
		fmt.Fprintf(os.Stderr, "guard git discovery: %v\n", err)
		return 1
	}
	return m.Run()
}

// Guard stops git, for this process and every child it starts, from climbing
// out of the temp directory or out of a package directory into a real
// repository. Git still inspects the directory it starts in, so tests never
// run fid with a real repository's root as the working directory.
func Guard() error {
	// Ours always go first, even if inherited later in the list: git does not
	// resolve symlinks in entries after an empty one.
	required := requiredCeilings()
	ceilings := append([]string(nil), required...)
	for _, inherited := range filepath.SplitList(os.Getenv(CeilingDirectoriesEnv)) {
		if !slices.Contains(required, inherited) {
			ceilings = append(ceilings, inherited)
		}
	}
	if err := os.Setenv(CeilingDirectoriesEnv, strings.Join(ceilings, string(os.PathListSeparator))); err != nil {
		return fmt.Errorf("set %s: %w", CeilingDirectoriesEnv, err)
	}
	return nil
}

func requiredCeilings() []string {
	return []string{os.TempDir(), ModuleRoot()}
}

// checkGuard fails unless Guard has run in this process: the ceilings it
// sets come first, in order.
func checkGuard() error {
	required := requiredCeilings()
	value := os.Getenv(CeilingDirectoriesEnv)
	if entries := filepath.SplitList(value); len(entries) < len(required) || !slices.Equal(entries[:len(required)], required) {
		return fmt.Errorf("git discovery is not guarded: %s is %q, want it to start with %q; call gittest.Guard or gittest.Main from TestMain", CeilingDirectoriesEnv, value, strings.Join(required, string(os.PathListSeparator)))
	}
	return nil
}

// InitRepository makes directory a git repository and checks that git
// reports directory itself as the repository. It refuses unless this
// process has guarded git discovery.
func InitRepository(directory string) error {
	if err := checkGuard(); err != nil {
		return err
	}
	if output, err := exec.Command("git", "init", "-q", "--template=", directory).CombinedOutput(); err != nil {
		return fmt.Errorf("git init %q: %w: %s", directory, err, output)
	}
	return CheckRepository(directory)
}

// CheckRepository fails unless git reports directory itself as the
// repository's top level.
func CheckRepository(directory string) error {
	want, err := filepath.EvalSymlinks(directory)
	if err != nil {
		return fmt.Errorf("resolve test repository %q: %w", directory, err)
	}
	output, err := exec.Command("git", "-C", directory, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("git rev-parse --show-toplevel in %q: %w", directory, err)
	}
	if got := strings.TrimSpace(string(output)); got != want {
		return fmt.Errorf("git reports repository %q for test directory %q", got, want)
	}
	return nil
}

// NewWorkspace returns an empty temp directory that is a guarded git
// repository, ready for fid init.
func NewWorkspace(t testing.TB) string {
	t.Helper()
	workspace := t.TempDir()
	if err := InitRepository(workspace); err != nil {
		t.Fatalf("create git workspace: %v", err)
	}
	return workspace
}

// NewNamedWorkspace returns a guarded git repository called name inside a
// fresh temp directory, for tests where the directory's name matters.
func NewNamedWorkspace(t testing.TB, name string) string {
	t.Helper()
	workspace := filepath.Join(t.TempDir(), name)
	if err := os.Mkdir(workspace, 0o755); err != nil {
		t.Fatalf("create git workspace %q: %v", name, err)
	}
	if err := InitRepository(workspace); err != nil {
		t.Fatalf("create git workspace %q: %v", name, err)
	}
	return workspace
}

// ResolvedPath returns path with symlinks resolved: the form git, and so
// fid, reports a repository's top level in, even when TMPDIR is a symlink.
func ResolvedPath(t testing.TB, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("resolve %q: %v", path, err)
	}
	return resolved
}

// UnbornHead points HEAD at a branch with no commits, so git still finds the
// repository but fid cannot observe it, and returns a function that puts
// HEAD back. Tests use it for a repository fid can discover but not inspect.
func UnbornHead(t testing.TB, repository string) (restore func()) {
	t.Helper()
	git := func(args ...string) string {
		t.Helper()
		output, err := exec.Command("git", append([]string{"-C", repository}, args...)...).Output()
		if err != nil {
			t.Fatalf("git %v in %q: %v", args, repository, err)
		}
		return strings.TrimSpace(string(output))
	}
	head := git("symbolic-ref", "HEAD")
	git("symbolic-ref", "HEAD", "refs/heads/fid-test-unborn")
	return func() {
		t.Helper()
		git("symbolic-ref", "HEAD", head)
	}
}
