package export

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/config"
	"golang.org/x/sys/unix"
)

func safeImportTarget(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", &ImportTargetConflictError{TargetRoot: root, Reason: "unsafe-path"}
	}
	info, err := os.Lstat(abs)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", &ImportTargetConflictError{TargetRoot: abs, Reason: "unsafe-path"}
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil || resolved != abs {
		return "", &ImportTargetConflictError{TargetRoot: abs, Reason: "unsafe-path"}
	}
	return abs, nil
}

func withImportLock(target string, fn func() error) error {
	// Lock the target directory inode itself: no pathname can be unlinked by a
	// contender, no source-tree file is introduced, and process exit releases it.
	file, err := os.Open(target)
	if err != nil {
		return fmt.Errorf("import: open target lock: %w", err)
	}
	defer func() { _ = file.Close() }()
	if err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		if errors.Is(err, unix.EWOULDBLOCK) {
			return fmt.Errorf("import: target lock unavailable")
		}
		return fmt.Errorf("import: acquire target lock: %w", err)
	}
	defer func() { _ = unix.Flock(int(file.Fd()), unix.LOCK_UN) }()
	return fn()
}

func snapshotImportInput(target, archivePath, digest string) (string, error) {
	input, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("import: open archive: %w", err)
	}
	defer func() { _ = input.Close() }()
	output, err := os.CreateTemp(filepath.Dir(config.LedgerDir(target)), config.LedgerStagingPrefix)
	if err != nil {
		return "", fmt.Errorf("import: create archive snapshot: %w", err)
	}
	path := output.Name()
	defer func() {
		if output != nil {
			_ = output.Close()
		}
	}()
	if _, err := io.Copy(output, input); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("import: snapshot archive: %w", err)
	}
	if err := output.Close(); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("import: close archive snapshot: %w", err)
	}
	output = nil
	got, err := fileDigest(path)
	if err != nil || got != digest {
		_ = os.Remove(path)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("import: archive changed while snapshotting")
	}
	return path, nil
}

func ensureEmptyImportTarget(target string) error {
	info, err := os.Lstat(config.LedgerDir(target))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return &ImportTargetConflictError{TargetRoot: target, Reason: "unsafe-path"}
	}
	return &ImportTargetConflictError{TargetRoot: target, Reason: "populated"}
}

// publishImport renames the staged ledger into place as target's ledger.
func publishImport(stage, target string) error {
	if err := unix.Renameat2(unix.AT_FDCWD, stage, unix.AT_FDCWD, config.LedgerDir(target), unix.RENAME_NOREPLACE); err != nil {
		if errors.Is(err, unix.EEXIST) {
			return &ImportTargetConflictError{TargetRoot: target, Reason: "populated"}
		}
		return fmt.Errorf("import: atomically publish staged project: %w", err)
	}
	return nil
}
func syncTree(root string) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("import: staged symlink is unsafe")
		}
		if entry.IsDir() || entry.Type().IsRegular() {
			return syncPath(path)
		}
		return fmt.Errorf("import: staged special file is unsafe")
	})
}
func syncPath(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("import: open for sync %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()
	if err := file.Sync(); err != nil {
		return fmt.Errorf("import: sync %q: %w", path, err)
	}
	return nil
}
func callImportFault(fault func(ImportStage) error, stage ImportStage) error {
	if fault == nil {
		return nil
	}
	return fault(stage)
}
