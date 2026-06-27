package journal

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWithLockWritesOwnerMetadata(t *testing.T) {
	workspace := t.TempDir()
	mustMkdirAll(t, filepath.Join(workspace, ".fid"))

	owner := LockOwner{
		Command:   "fid issue create --by developer",
		ActorRole: "developer",
	}
	err := withLockOptions(workspace, owner, defaultLockOptions(), func(lock *WriteLock) error {
		metadata, readErr := readLockMetadata(lock.path)
		if readErr != nil {
			return readErr
		}

		if metadata.PID <= 0 {
			t.Fatalf("lock pid = %d, want > 0", metadata.PID)
		}
		if strings.TrimSpace(metadata.Hostname) == "" {
			t.Fatal("lock hostname is empty")
		}
		if metadata.StartedAt.IsZero() {
			t.Fatal("lock started_at is zero")
		}
		if metadata.Command != owner.Command {
			t.Fatalf("lock command = %q, want %q", metadata.Command, owner.Command)
		}
		if metadata.ActorRole == nil || *metadata.ActorRole != owner.ActorRole {
			t.Fatalf("lock actor role = %v, want %q", metadata.ActorRole, owner.ActorRole)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("withLockOptions() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(workspace, ".fid", lockDirectoryName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("lock directory exists after release: stat error = %v, want %v", err, os.ErrNotExist)
	}
}

func TestWithLockTimesOutWithHolderNamedError(t *testing.T) {
	workspace := t.TempDir()
	mustMkdirAll(t, filepath.Join(workspace, ".fid"))

	heldLock, err := acquireLock(workspace, LockOwner{
		Command:   "fid issue create --by developer",
		ActorRole: "developer",
	}, defaultLockOptions())
	if err != nil {
		t.Fatalf("acquireLock() error = %v", err)
	}
	t.Cleanup(func() {
		if releaseErr := heldLock.release(); releaseErr != nil {
			t.Fatalf("heldLock.release() error = %v", releaseErr)
		}
	})

	err = withLockOptions(workspace, LockOwner{
		Command:   "fid issue create --by qa_lead",
		ActorRole: "qa_lead",
	}, lockOptions{
		waitTimeout: 20 * time.Millisecond,
		retryDelay:  time.Millisecond,
		now:         time.Now,
		sleep:       time.Sleep,
	}, func(lock *WriteLock) error {
		t.Fatalf("withLockOptions() unexpectedly acquired %v", lock)
		return nil
	})
	if err == nil {
		t.Fatal("withLockOptions() error = nil, want timeout error")
	}

	lockPath := filepath.Join(workspace, ".fid", lockDirectoryName)
	if !strings.Contains(err.Error(), lockPath) {
		t.Fatalf("timeout error = %q, want lock path %q", err, lockPath)
	}
	if !strings.Contains(err.Error(), "developer") {
		t.Fatalf("timeout error = %q, want holder role", err)
	}
}

func TestWithLockDefaultWaitsForShortHeldLock(t *testing.T) {
	workspace := t.TempDir()
	mustMkdirAll(t, filepath.Join(workspace, ".fid"))

	heldLock, err := acquireLock(workspace, LockOwner{
		Command:   "fid issue claim FID-123 --by developer",
		ActorRole: "developer",
	}, defaultLockOptions())
	if err != nil {
		t.Fatalf("acquireLock() error = %v", err)
	}

	releaseErrs := make(chan error, 1)
	go func() {
		time.Sleep(300 * time.Millisecond)
		releaseErrs <- heldLock.release()
	}()

	acquired := false
	err = withLockOptions(workspace, LockOwner{
		Command:   "fid issue claim FID-123 --by qa_reviewer",
		ActorRole: "qa_reviewer",
	}, defaultLockOptions(), func(lock *WriteLock) error {
		acquired = true
		return nil
	})
	if err != nil {
		t.Fatalf("withLockOptions() error = %v", err)
	}
	if !acquired {
		t.Fatal("withLockOptions() callback was not called")
	}
	if err := <-releaseErrs; err != nil {
		t.Fatalf("heldLock.release() error = %v", err)
	}
}
