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
// Discover runs `git rev-parse` twice, once for the top level and once for the
// git common directory, so each output is one path, whatever it holds. It
// refuses with ErrUnsafeLedger where fid could write outside <common>/fid or
// where git could track the ledger: where LedgerDir names a different ledger
// from git's, or where the common directory is inside the worktree but not
// within a .git directory.
func Discover(start string) (string, error) {
	directory, err := absoluteDirectory(start)
	if err != nil {
		return "", err
	}
	root, err := gitPath(directory, "--show-toplevel")
	if err != nil {
		return "", notARepository(directory, err)
	}
	common, err := gitPath(directory, "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", notARepository(directory, err)
	}
	if err := checkLedgerLocation(root, common); err != nil {
		return "", err
	}
	return root, nil
}

// gitPath returns the one path `git rev-parse` reports for arguments. Git
// ends its output with a newline, and a path may end in one too, so exactly
// one is stripped.
func gitPath(directory string, arguments ...string) (string, error) {
	output, err := gitOutput(directory, append([]string{"rev-parse"}, arguments...)...)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(output), "\n"), nil
}

func notARepository(directory string, err error) error {
	return fmt.Errorf("%w: %w", ErrNotARepository, unsafeRepositoryReason("not-a-repository", "%s is not in a git repository (%v)", directory, err))
}

// checkLedgerLocation refuses unless git and LedgerDir name the same ledger
// for the worktree root, and git can't track it. LedgerDir reads paths
// lexically while git resolves symlinks, so both sides are compared resolved:
// a .git that is itself a symlink is no disagreement.
func checkLedgerLocation(root, common string) error {
	resolvedGitLedger, agree := agreedLedger(root, common)
	if !agree {
		return unsafeLedger("git and the worktree name different ledgers: git names %s, but the worktree %s names %s; check GIT_DIR, GIT_COMMON_DIR, GIT_WORK_TREE and core.worktree", filepath.Join(common, config.LedgerDirName), root, config.LedgerDir(root))
	}
	resolvedRoot := root
	if resolved, err := filepath.EvalSymlinks(root); err == nil {
		resolvedRoot = resolved
	}
	resolvedCommon := filepath.Dir(resolvedGitLedger)
	if !within(resolvedRoot, resolvedCommon) {
		return nil
	}
	relative, _ := filepath.Rel(resolvedRoot, resolvedCommon) // within checked it
	for _, segment := range strings.Split(relative, string(filepath.Separator)) {
		if segment == ".git" {
			return nil
		}
	}
	return unsafeLedger("git could track the ledger %s: the git common directory %s is inside the worktree %s, outside any .git directory", resolvedGitLedger, resolvedCommon, root)
}

// agreedLedger returns the ledger that git, naming common for the worktree
// root, and LedgerDir both name, with symlinks resolved. It reports false
// where they name different ledgers.
func agreedLedger(root, common string) (string, bool) {
	gitLedger, gitErr := resolveParent(filepath.Join(common, config.LedgerDirName))
	ledger, err := resolveParent(config.LedgerDir(root))
	return gitLedger, gitErr == nil && err == nil && ledger == gitLedger
}

// within reports whether path is base or below it. Both are compared as
// given, so the caller resolves symlinks in both or in neither.
func within(base, path string) bool {
	relative, err := filepath.Rel(base, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
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
