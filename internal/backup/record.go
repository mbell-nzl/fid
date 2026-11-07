package backup

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

// appendRecord adds one line to the record with a single append, synced
// before it returns. Backups through one clone append while they hold it,
// so lines are in the order the attempts finished.
func appendRecord(path string, record Record) error {
	line, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encode backup record: %w", err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("open backup record: %w", err)
	}
	if _, err := file.Write(append(line, '\n')); err != nil {
		_ = file.Close()
		return fmt.Errorf("append backup record: %w", err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("sync backup record: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close backup record: %w", err)
	}
	return nil
}

const lockRetryDelay = 50 * time.Millisecond

// lock takes an exclusive flock on path, waiting up to wait for another
// backup to release it. The kernel releases it if the process dies, so a
// crashed backup never leaves it held.
func lock(path string, wait time.Duration, now func() time.Time) (func(), error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, failure(KindInternal, err, "open the backup lock")
	}
	deadline := now().Add(wait)
	for {
		err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB)
		if err == nil {
			return func() {
				_ = unix.Flock(int(file.Fd()), unix.LOCK_UN)
				_ = file.Close()
			}, nil
		}
		if !errors.Is(err, unix.EWOULDBLOCK) {
			_ = file.Close()
			return nil, failure(KindInternal, err, "take the backup lock")
		}
		if !now().Before(deadline) {
			_ = file.Close()
			return nil, &Error{Kind: KindLockUnavailable, Path: path,
				Message: fmt.Sprintf("another backup is using fid's clone of the backup repository; waited %s for %s", wait, path)}
		}
		time.Sleep(lockRetryDelay)
	}
}
