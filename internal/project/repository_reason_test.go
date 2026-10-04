package project

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

func TestObserveRecordedRepositoryNotARepositoryReason(t *testing.T) {
	t.Parallel()
	_, err := ObserveRecordedRepository(t.TempDir(), "P-test")
	var unavailable *RepositoryUnavailableError
	if !errors.Is(err, ErrUnsafeRepository) || !errors.As(err, &unavailable) {
		t.Fatalf("error = %v, want typed unsafe repository", err)
	}
	if unavailable.Reason != "not-a-repository" || unavailable.Detail == "" {
		t.Fatalf("unavailable = %#v", unavailable)
	}
}

func TestRepositoryUnavailableCauseFixtures(t *testing.T) {
	t.Parallel()
	git := func(root string, args ...string) {
		t.Helper()
		command := exec.Command("git", append([]string{"-C", root}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	root := t.TempDir()
	git(root, "init", "-q")
	assertRepositoryReason(t, root, "not-a-repository") // unborn HEAD
	git(root, "config", "user.name", "Source observer")
	git(root, "config", "user.email", "observer@example.test")
	path := filepath.Join(root, "tracked.txt")
	if err := os.WriteFile(path, []byte("tracked\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	git(root, "add", "tracked.txt")
	git(root, "commit", "-qm", "base")
	git(root, "update-index", "--assume-unchanged", "tracked.txt")
	assertRepositoryReason(t, root, "tracked-flags")
	git(root, "update-index", "--no-assume-unchanged", "tracked.txt")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}
	assertRepositoryReason(t, root, "unsafe-path")
}

func TestRepositoryCaptureChangeReason(t *testing.T) {
	t.Parallel()
	reads := 0
	_, err := captureRepository(func() (Observation, error) {
		reads++
		return Observation{HeadOID: string(rune('0' + reads))}, nil
	})
	var unavailable *RepositoryUnavailableError
	if !errors.As(err, &unavailable) || unavailable.Reason != "changed-during-capture" || reads != 2 {
		t.Fatalf("capture = %v, reads %d", err, reads)
	}
}

func TestRepositoryUnmergedAndSizeClassification(t *testing.T) {
	t.Parallel()
	_, err := parseUnmergedPaths([]string{"100644 deadbeef 1\t../invalid"})
	var unavailable *RepositoryUnavailableError
	if !errors.As(err, &unavailable) || unavailable.Reason != "unmerged-paths" || unavailable.Detail != "invalid unmerged path" {
		t.Fatalf("unmerged parser = %v", err)
	}
	if err := checkRepositoryOutputSize(maxRepositoryOutputBytes); err != nil {
		t.Fatalf("limit boundary: %v", err)
	}
	err = checkRepositoryOutputSize(maxRepositoryOutputBytes + 1)
	if !errors.As(err, &unavailable) || unavailable.Reason != "unsafe-path" {
		t.Fatalf("oversize output = %v", err)
	}
	commandFailure := unsafeRepository("list unmerged entries: %v", errors.New("ls-files failed"))
	if !errors.As(commandFailure, &unavailable) || unavailable.Reason != "changed-during-capture" {
		t.Fatalf("ls-files failure = %v", commandFailure)
	}
}

func assertRepositoryReason(t *testing.T, root, reason string) {
	t.Helper()
	_, err := ObserveRecordedRepository(root, "P-test")
	var unavailable *RepositoryUnavailableError
	if !errors.As(err, &unavailable) || !errors.Is(err, ErrUnsafeRepository) || unavailable.Reason != reason {
		t.Fatalf("repository error = %v, want %s", err, reason)
	}
}
