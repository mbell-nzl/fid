package backup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

// appendRecord adds one line to the record, synced before it returns. It
// holds an flock on the record while it repairs and appends, so attempts
// that record without holding fid's clone, such as one with no target, never
// interleave with another's repair. Backups through one clone append while
// they hold it, so lines are in the order the attempts finished.
func appendRecord(path string, record Record) error {
	line, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encode backup record: %w", err)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("open backup record: %w", err)
	}
	// Closing the file releases the flock.
	if err := unix.Flock(int(file.Fd()), unix.LOCK_EX); err != nil {
		_ = file.Close()
		return fmt.Errorf("lock backup record: %w", err)
	}
	if err := repairTail(file); err != nil {
		_ = file.Close()
		return fmt.Errorf("repair backup record: %w", err)
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

// repairTail makes the record end with a complete line again after a write
// was cut short, by a full disk or a killed process, so the next line starts
// on its own. A tail that is a whole record missing only its newline gets
// the newline. Any other tail is a fragment no reader can use, so it is cut
// back to the last complete line; the attempt that wrote it reported then
// that recording failed.
func repairTail(file *os.File) error {
	info, err := file.Stat()
	if err != nil || info.Size() == 0 {
		return err
	}
	last := make([]byte, 1)
	if _, err := file.ReadAt(last, info.Size()-1); err != nil {
		return err
	}
	if last[0] == '\n' {
		return nil
	}
	payload, err := io.ReadAll(io.NewSectionReader(file, 0, info.Size()))
	if err != nil {
		return err
	}
	complete := bytes.LastIndexByte(payload, '\n') + 1
	var tail Record
	if json.Unmarshal(payload[complete:], &tail) == nil && tail.Schema == RecordSchema {
		_, err := file.Write([]byte{'\n'})
		return err
	}
	return file.Truncate(int64(complete))
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
