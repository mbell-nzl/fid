// Package project resolves fid project roots through git.
package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
)

var ErrNotFound = errors.New("project not found")

// ErrNotARepository marks a discovery refusal: git found no worktree at or
// above the start directory. The error also carries a
// RepositoryUnavailableError with reason not-a-repository.
var ErrNotARepository = errors.New("not a git repository")

// ErrUnsafeLedger marks a discovery refusal: git and LedgerDir name different
// ledgers, or git could track the ledger. The error also carries a
// RepositoryUnavailableError with reason unsafe-path.
var ErrUnsafeLedger = errors.New("unsafe ledger location")

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
// git reports it, so symlinks are resolved. Discovery starts at start, as
// `git -C start` would.
//
// One `git rev-parse` reports the top level and the git common directory.
// Discover refuses with ErrUnsafeLedger where fid could write outside
// <common>/fid or where git could track the ledger: where LedgerDir names a
// different ledger from git's, or where the common directory is inside the
// worktree but not within a .git directory.
func Discover(start string) (string, error) {
	directory, err := absoluteDirectory(start)
	if err != nil {
		return "", err
	}
	output, err := gitText(directory, "rev-parse", "--show-toplevel", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrNotARepository, unsafeRepositoryReason("not-a-repository", "%s is not in a git repository (%v)", directory, err))
	}
	root, common, ok := strings.Cut(output, "\n")
	if !ok || root == "" || common == "" || strings.Contains(common, "\n") {
		return "", unsafeLedger("git rev-parse in %s reported %q, want a top level and a common directory", directory, output)
	}
	if err := checkLedgerLocation(root, common); err != nil {
		return "", err
	}
	return root, nil
}

// checkLedgerLocation refuses unless git and LedgerDir name the same ledger
// for the worktree root, and git can't track it. LedgerDir reads paths
// lexically while git resolves symlinks, so both sides are compared resolved:
// a .git that is itself a symlink is no disagreement.
func checkLedgerLocation(root, common string) error {
	gitLedger := filepath.Join(common, config.LedgerDirName)
	ledger := config.LedgerDir(root)
	resolvedGitLedger, gitErr := resolveParent(gitLedger)
	resolvedLedger, err := resolveParent(ledger)
	if gitErr != nil || err != nil || resolvedLedger != resolvedGitLedger {
		return unsafeLedger("git and the worktree name different ledgers: git names %s, but the worktree %s names %s; check GIT_DIR, GIT_COMMON_DIR, GIT_WORK_TREE and core.worktree", gitLedger, root, ledger)
	}
	resolvedRoot := root
	if resolved, err := filepath.EvalSymlinks(root); err == nil {
		resolvedRoot = resolved
	}
	resolvedCommon := filepath.Dir(resolvedGitLedger)
	relative, err := filepath.Rel(resolvedRoot, resolvedCommon)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return nil
	}
	for _, segment := range strings.Split(relative, string(filepath.Separator)) {
		if segment == ".git" {
			return nil
		}
	}
	return unsafeLedger("git could track the ledger %s: the git common directory %s is inside the worktree %s, outside any .git directory", resolvedGitLedger, resolvedCommon, root)
}

// resolveParent resolves symlinks in path's parent, which must exist; path
// itself need not.
func resolveParent(path string) (string, error) {
	parent, err := filepath.EvalSymlinks(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, filepath.Base(path)), nil
}

func unsafeLedger(format string, arguments ...any) error {
	return fmt.Errorf("%w: %w", ErrUnsafeLedger, unsafeRepositoryReason("unsafe-path", format, arguments...))
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
