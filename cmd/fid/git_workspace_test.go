package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const gitCeilingDirectoriesEnv = "GIT_CEILING_DIRECTORIES"

// gitGuardChildEnv marks the child process that TestGitDiscoveryStopsAtTheTempDirectory starts.
const gitGuardChildEnv = "FID_GIT_GUARD_CHILD"

// guardGitDiscovery stops git, for this process and every child it starts,
// from climbing out of the temp directory or out of a package directory into
// a real repository. Git still inspects the directory it starts in, so tests
// never run fid with a real repository's root as the working directory.
func guardGitDiscovery(moduleRoot string) error {
	// Ours always go first, even if inherited later in the list: git does not
	// resolve symlinks in entries after an empty one.
	required := []string{os.TempDir(), moduleRoot}
	ceilings := append([]string(nil), required...)
	for _, inherited := range filepath.SplitList(os.Getenv(gitCeilingDirectoriesEnv)) {
		if !slices.Contains(required, inherited) {
			ceilings = append(ceilings, inherited)
		}
	}
	if err := os.Setenv(gitCeilingDirectoriesEnv, strings.Join(ceilings, string(os.PathListSeparator))); err != nil {
		return fmt.Errorf("set %s: %w", gitCeilingDirectoriesEnv, err)
	}
	return nil
}

// initGuardedGitRepository makes directory a git repository and checks that
// git reports directory itself as the repository.
func initGuardedGitRepository(directory string) error {
	if output, err := exec.Command("git", "init", "-q", "--template=", directory).CombinedOutput(); err != nil {
		return fmt.Errorf("git init %q: %w: %s", directory, err, output)
	}
	return checkGuardedGitRepository(directory)
}

// checkGuardedGitRepository fails unless git reports directory itself as the
// repository's top level.
func checkGuardedGitRepository(directory string) error {
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

// newGitWorkspace returns an empty temp directory that is a guarded git
// repository, ready for fid init.
func newGitWorkspace(t *testing.T) string {
	t.Helper()
	workspace := t.TempDir()
	if err := initGuardedGitRepository(workspace); err != nil {
		t.Fatalf("create git workspace: %v", err)
	}
	return workspace
}

// newNamedGitWorkspace returns a guarded git repository called name inside a
// fresh temp directory, for tests where the directory's name matters.
func newNamedGitWorkspace(t *testing.T, name string) string {
	t.Helper()
	workspace := filepath.Join(t.TempDir(), name)
	if err := os.Mkdir(workspace, 0o755); err != nil {
		t.Fatalf("create git workspace %q: %v", name, err)
	}
	if err := initGuardedGitRepository(workspace); err != nil {
		t.Fatalf("create git workspace %q: %v", name, err)
	}
	return workspace
}

// unbornGitHead points HEAD at a branch with no commits, so git still finds
// the repository but fid cannot observe it, and returns a function that puts
// HEAD back. Tests use it for a repository fid can discover but not inspect.
func unbornGitHead(t *testing.T, repository string) (restore func()) {
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

func TestUnbornGitHeadIsDiscoverableButNotObservable(t *testing.T) {
	t.Parallel()
	workspace := newGitWorkspace(t)
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
	restore := unbornGitHead(t, workspace)
	if err := checkGuardedGitRepository(workspace); err != nil {
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

func TestGuardedGitWorkspaceIsItsOwnRepository(t *testing.T) {
	t.Parallel()
	workspace := newGitWorkspace(t)
	if err := checkGuardedGitRepository(workspace); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(workspace, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatalf("Mkdir(nested) error = %v", err)
	}
	if err := checkGuardedGitRepository(nested); err == nil {
		t.Fatal("check of a subdirectory = nil, want refusal because git reports its parent")
	}
}

// TestGitDiscoveryStopsAtTheTempDirectory runs this test binary again with its
// temp directory inside another repository, so the guard TestMain installs is
// what stops git. It covers a plain temp directory and a symlinked one that
// is also inherited after an empty ceiling entry.
func TestGitDiscoveryStopsAtTheTempDirectory(t *testing.T) {
	t.Parallel()
	if os.Getenv(gitGuardChildEnv) != "" {
		directory := t.TempDir()
		if output, err := exec.Command("git", "-C", directory, "rev-parse", "--show-toplevel").Output(); err == nil {
			t.Fatalf("git found repository %q from temp directory %q, want ceiling refusal", strings.TrimSpace(string(output)), directory)
		}
		return
	}
	ancestor := newGitWorkspace(t)
	temp := filepath.Join(ancestor, "tmp")
	if err := os.Mkdir(temp, 0o755); err != nil {
		t.Fatalf("Mkdir(temp) error = %v", err)
	}
	alias := filepath.Join(t.TempDir(), "tmp-alias")
	if err := os.Symlink(temp, alias); err != nil {
		t.Fatalf("Symlink(temp) error = %v", err)
	}
	control := exec.Command("git", "-C", temp, "rev-parse", "--show-toplevel")
	control.Env = environmentWithout(os.Environ(), gitCeilingDirectoriesEnv)
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
		command.Env = append(environmentWithout(os.Environ(), gitCeilingDirectoriesEnv, "TMPDIR"), "TMPDIR="+child.tempDir, gitGuardChildEnv+"=1")
		if child.inherited != "" {
			command.Env = append(command.Env, gitCeilingDirectoriesEnv+"="+child.inherited)
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
