// Package project resolves fid project roots through git.
package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/config"
)

var ErrNotFound = errors.New("project not found")

// ErrNotARepository marks a discovery refusal: git found no worktree at or
// above the start directory. The error also carries a
// RepositoryUnavailableError with reason not-a-repository.
var ErrNotARepository = errors.New("not a git repository")

// Resolve returns the fid project root for start: the top level of its git
// worktree, which must directly contain the ledger.
func Resolve(start string) (string, error) {
	root, err := Discover(start)
	if err != nil {
		return "", err
	}
	if !isProjectRoot(root) {
		return "", fmt.Errorf("%w: %s", ErrNotFound, root)
	}
	return root, nil
}

// Discover returns the top level of the git worktree containing start, as
// one `git rev-parse --show-toplevel` reports it, so symlinks are resolved.
// Discovery starts at start, as `git -C start` would.
func Discover(start string) (string, error) {
	directory, err := absoluteDirectory(start)
	if err != nil {
		return "", err
	}
	root, err := gitText(directory, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrNotARepository, unsafeRepositoryReason("not-a-repository", "%s is not in a git repository (%v)", directory, err))
	}
	return root, nil
}

func absoluteDirectory(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve project path %q: %w", path, err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("inspect project path %q: %w", abs, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("inspect project path %q: not a directory", abs)
	}
	return filepath.Clean(abs), nil
}

func isProjectRoot(root string) bool {
	info, err := os.Lstat(config.LedgerDir(root))
	return err == nil && info.IsDir() && info.Mode()&fs.ModeSymlink == 0
}
