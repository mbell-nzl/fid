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

// guardGitDiscovery stops git, for this process and every child it starts,
// from climbing out of the temp directory or out of a package directory into
// a real repository. Git still inspects the directory it starts in, so tests
// never run fid with a real repository's root as the working directory.
func guardGitDiscovery(moduleRoot string) error {
	existing := filepath.SplitList(os.Getenv(gitCeilingDirectoriesEnv))
	var ceilings []string
	for _, ceiling := range []string{os.TempDir(), moduleRoot} {
		if !slices.Contains(existing, ceiling) {
			ceilings = append(ceilings, ceiling)
		}
	}
	// Ours go first: git does not resolve symlinks in entries after an empty one.
	ceilings = append(ceilings, existing...)
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

func TestGitDiscoveryStopsAtTheTempDirectory(t *testing.T) {
	t.Parallel()
	plain := t.TempDir()
	if err := checkGuardedGitRepository(plain); err == nil {
		t.Fatal("check of a plain temp directory = nil, want no repository")
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
