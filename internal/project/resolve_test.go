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

func TestDiscoverReturnsTheWorktreeTopLevelFromAnyDirectory(t *testing.T) {
	t.Parallel()
	repository := newDiscoveryRepository(t)
	nested := filepath.Join(repository, "nested", "directory")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("MkdirAll(nested) error = %v", err)
	}
	for _, start := range []string{repository, nested, filepath.Join(nested, "..")} {
		got, err := Discover(start)
		if err != nil || got != repository {
			t.Fatalf("Discover(%q) = %q, %v; want top level %q", start, got, err, repository)
		}
	}
}

func TestDiscoverResolvesSymlinksAsGitReportsThem(t *testing.T) {
	t.Parallel()
	repository := newDiscoveryRepository(t)
	nested := filepath.Join(repository, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatalf("Mkdir(nested) error = %v", err)
	}
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(nested, alias); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}
	got, err := Discover(alias)
	if err != nil || got != repository {
		t.Fatalf("Discover(symlink) = %q, %v; want real top level %q", got, err, repository)
	}
}

func TestDiscoverReturnsALinkedWorktreesOwnTopLevel(t *testing.T) {
	t.Parallel()
	repository := newDiscoveryRepository(t)
	discoveryGit(t, repository, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-q", "--allow-empty", "-m", "base")
	linked := filepath.Join(t.TempDir(), "linked")
	discoveryGit(t, repository, "worktree", "add", "-q", "--detach", linked)
	want, err := filepath.EvalSymlinks(linked)
	if err != nil {
		t.Fatalf("EvalSymlinks(linked) error = %v", err)
	}
	if got, err := Discover(linked); err != nil || got != want {
		t.Fatalf("Discover(linked worktree) = %q, %v; want %q", got, err, want)
	}
}

func TestDiscoverRefusesOutsideAGitRepository(t *testing.T) {
	t.Parallel()
	start := t.TempDir()
	_, err := Discover(start)
	if !errors.Is(err, ErrNotARepository) {
		t.Fatalf("Discover(outside) error = %v, want ErrNotARepository", err)
	}
	var unavailable *RepositoryUnavailableError
	if !errors.As(err, &unavailable) || unavailable.Reason != "not-a-repository" || !strings.Contains(unavailable.Detail, start) {
		t.Fatalf("Discover(outside) unavailable = %#v, want reason not-a-repository naming %q", unavailable, start)
	}
}

func TestDiscoverReportsAMissingStartAsAPathError(t *testing.T) {
	t.Parallel()
	missing := filepath.Join(t.TempDir(), "missing")
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatalf("WriteFile(file) error = %v", err)
	}
	for _, start := range []string{missing, file} {
		_, err := Discover(start)
		if err == nil || errors.Is(err, ErrNotARepository) || errors.Is(err, ErrNotFound) {
			t.Fatalf("Discover(%q) error = %v, want a path error", start, err)
		}
	}
}

func TestResolveFindsTheLedgerAtTheTopLevelFromAnyDirectory(t *testing.T) {
	t.Parallel()
	repository := newDiscoveryRepository(t)
	if err := os.Mkdir(config.LedgerDir(repository), 0o755); err != nil {
		t.Fatalf("Mkdir(.fid) error = %v", err)
	}
	nested := filepath.Join(repository, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatalf("Mkdir(nested) error = %v", err)
	}
	for _, start := range []string{repository, nested} {
		got, err := Resolve(start)
		if err != nil || got != repository {
			t.Fatalf("Resolve(%q) = %q, %v; want %q", start, got, err, repository)
		}
	}
}

func TestResolveIgnoresALedgerBelowTheTopLevel(t *testing.T) {
	t.Parallel()
	repository := newDiscoveryRepository(t)
	nested := filepath.Join(repository, "nested")
	if err := os.MkdirAll(config.LedgerDir(nested), 0o755); err != nil {
		t.Fatalf("MkdirAll(nested ledger) error = %v", err)
	}
	for _, start := range []string{repository, nested} {
		if _, err := Resolve(start); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Resolve(%q) error = %v, want ErrNotFound", start, err)
		}
	}
}

func TestResolveRefusesOutsideAGitRepositoryEvenWithALedger(t *testing.T) {
	t.Parallel()
	start := t.TempDir()
	if err := os.MkdirAll(config.LedgerDir(start), 0o755); err != nil {
		t.Fatalf("MkdirAll(ledger) error = %v", err)
	}
	if _, err := Resolve(start); !errors.Is(err, ErrNotARepository) {
		t.Fatalf("Resolve(outside) error = %v, want ErrNotARepository", err)
	}
}

func TestResolveRejectsSymlinkedFIDDirectory(t *testing.T) {
	t.Parallel()
	repository := newDiscoveryRepository(t)
	outside := t.TempDir()
	if err := os.MkdirAll(config.LedgerDir(outside), 0o755); err != nil {
		t.Fatalf("MkdirAll(outside ledger) error = %v", err)
	}
	if err := os.Symlink(config.LedgerDir(outside), config.LedgerDir(repository)); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}

	if _, err := Resolve(repository); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Resolve() error = %v, want symlink escape rejection", err)
	}
}

// newDiscoveryRepository returns a fresh guarded git repository, as the real
// path git reports for its top level.
func newDiscoveryRepository(t *testing.T) string {
	t.Helper()
	return gittest.ResolvedPath(t, gittest.NewWorkspace(t))
}

func discoveryGit(t *testing.T, repository string, arguments ...string) {
	t.Helper()
	if output, err := exec.Command("git", append([]string{"-C", repository}, arguments...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
}
