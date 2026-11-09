package journal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

const (
	lockDirectoryName      = "write.lock"
	lockMetadataFileName   = "owner.json"
	defaultLockWaitTimeout = 5 * time.Second
	defaultLockRetryDelay  = 10 * time.Millisecond
)

type LockOwner struct {
	Command string
}

type LockUnavailableError struct {
	LockPath string
}

func (err *LockUnavailableError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("write lock unavailable at %q", err.LockPath)
}

type WriteLock struct {
	ledgerDir string
	path      string
	released  bool
}

type lockMetadata struct {
	PID       int       `json:"pid"`
	Hostname  string    `json:"hostname"`
	StartedAt time.Time `json:"started_at"`
	Command   string    `json:"command"`
}

type lockOptions struct {
	waitTimeout time.Duration
	retryDelay  time.Duration
	now         func() time.Time
	sleep       func(time.Duration)
}

func WithLock(workspaceRoot string, owner LockOwner, fn func(*WriteLock) error) error {
	return withLockOptions(workspaceRoot, owner, defaultLockOptions(), fn)
}

// WithLedgerLock is WithLock for the ledger at ledgerDir, rather than the
// project at a root: a ledger being staged has no .git to find it by.
func WithLedgerLock(ledgerDir string, owner LockOwner, fn func(*WriteLock) error) error {
	return withLedgerLockOptions(ledgerDir, owner, defaultLockOptions(), fn)
}

func tryWithLedgerLock(ledgerDir string, owner LockOwner, fn func(*WriteLock) error) error {
	return withLedgerLockOptions(ledgerDir, owner, tryLockOptions(), fn)
}

func tryLockOptions() lockOptions {
	options := defaultLockOptions()
	options.waitTimeout = 0
	return options
}

// workspaceLedgerDir returns the ledger of the project at workspaceRoot, or
// "" for an empty root, which acquiring the lock refuses.
func workspaceLedgerDir(workspaceRoot string) string {
	if workspaceRoot == "" {
		return ""
	}
	return config.LedgerDir(workspaceRoot)
}

func withLockOptions(workspaceRoot string, owner LockOwner, options lockOptions, fn func(*WriteLock) error) error {
	return withLedgerLockOptions(workspaceLedgerDir(workspaceRoot), owner, options, fn)
}

func withLedgerLockOptions(ledgerDir string, owner LockOwner, options lockOptions, fn func(*WriteLock) error) error {
	if fn == nil {
		return fmt.Errorf("with lock: callback is required")
	}

	lock, err := acquireLedgerLock(ledgerDir, owner, options)
	if err != nil {
		return err
	}

	callbackErr := fn(lock)
	releaseErr := lock.release()
	if callbackErr != nil {
		return callbackErr
	}
	if releaseErr != nil {
		return releaseErr
	}

	return nil
}

func defaultLockOptions() lockOptions {
	return lockOptions{
		waitTimeout: defaultLockWaitTimeout,
		retryDelay:  defaultLockRetryDelay,
		now:         time.Now,
		sleep:       time.Sleep,
	}
}

func acquireLock(workspaceRoot string, owner LockOwner, options lockOptions) (*WriteLock, error) {
	return acquireLedgerLock(workspaceLedgerDir(workspaceRoot), owner, options)
}

func acquireLedgerLock(ledgerDir string, owner LockOwner, options lockOptions) (*WriteLock, error) {
	if ledgerDir == "" {
		return nil, fmt.Errorf("acquire lock: workspace root is required")
	}
	if options.now == nil {
		options.now = time.Now
	}
	if options.sleep == nil {
		options.sleep = time.Sleep
	}

	lockPath := filepath.Join(ledgerDir, lockDirectoryName)
	deadline := options.now().Add(options.waitTimeout)
	for {
		if err := os.Mkdir(lockPath, 0o755); err == nil {
			metadata, err := buildLockMetadata(owner, options.now)
			if err != nil {
				_ = os.Remove(lockPath)
				return nil, err
			}
			if err := writeLockMetadata(lockPath, metadata); err != nil {
				_ = os.RemoveAll(lockPath)
				return nil, err
			}

			return &WriteLock{
				ledgerDir: ledgerDir,
				path:      lockPath,
			}, nil
		} else if !os.IsExist(err) {
			return nil, fmt.Errorf("acquire lock %q: %w", lockPath, err)
		}

		if !options.now().Before(deadline) {
			return nil, &LockUnavailableError{LockPath: lockPath}
		}

		options.sleep(options.retryDelay)
	}
}

func buildLockMetadata(owner LockOwner, now func() time.Time) (lockMetadata, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return lockMetadata{}, fmt.Errorf("get hostname: %w", err)
	}

	metadata := lockMetadata{
		PID:       os.Getpid(),
		Hostname:  hostname,
		StartedAt: now().UTC(),
		Command:   owner.Command,
	}
	return metadata, nil
}

func writeLockMetadata(lockPath string, metadata lockMetadata) error {
	payload, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal lock metadata: %w", err)
	}

	metadataPath := filepath.Join(lockPath, lockMetadataFileName)
	if err := os.WriteFile(metadataPath, append(payload, '\n'), 0o644); err != nil {
		return fmt.Errorf("write lock metadata %q: %w", metadataPath, err)
	}

	return nil
}

func readLockMetadata(lockPath string) (lockMetadata, error) {
	metadataPath := filepath.Join(lockPath, lockMetadataFileName)
	payload, err := os.ReadFile(metadataPath)
	if err != nil {
		return lockMetadata{}, fmt.Errorf("read lock metadata %q: %w", metadataPath, err)
	}

	var metadata lockMetadata
	if err := json.Unmarshal(payload, &metadata); err != nil {
		return lockMetadata{}, fmt.Errorf("decode lock metadata %q: %w", metadataPath, err)
	}

	return metadata, nil
}

func (lock *WriteLock) release() error {
	if lock == nil {
		return fmt.Errorf("release lock: lock is nil")
	}
	if lock.released {
		return nil
	}
	if err := os.RemoveAll(lock.path); err != nil {
		return fmt.Errorf("release lock %q: %w", lock.path, err)
	}

	lock.released = true
	return nil
}
