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
	root, common, err := readGitReport(output, directory)
	if err != nil {
		return "", err
	}
	if err := checkLedgerLocation(root, common); err != nil {
		return "", err
	}
	return root, nil
}

// readGitReport splits git's report of the top level and the common
// directory, which git prints one per line. A newline in either path leaves
// more than one way to read it. The right reading then is the one whose top
// level holds the start directory, as it does wherever git found the
// repository from there, and only one reading can. Failing that, where
// GIT_WORK_TREE or core.worktree names a worktree elsewhere, it is the
// reading whose ledger LedgerDir agrees with. If no reading agrees, git and
// LedgerDir disagree however the report is read; if several do, fid can't
// tell which ledger git names. Both refuse.
func readGitReport(output, directory string) (root, common string, err error) {
	var readings [][2]string
	for index := 1; index < len(output)-1; index++ {
		if output[index] == '\n' {
			readings = append(readings, [2]string{output[:index], output[index+1:]})
		}
	}
	if len(readings) == 0 {
		return "", "", unsafeLedger("git rev-parse in %s reported %s, want a top level and a common directory", directory, output)
	}
	if len(readings) == 1 {
		return readings[0][0], readings[0][1], nil
	}
	if start, err := filepath.EvalSymlinks(directory); err == nil {
		holding := filterReadings(readings, func(reading [2]string) bool { return within(reading[0], start) })
		if len(holding) == 1 {
			return holding[0][0], holding[0][1], nil
		}
	}
	agreeing := filterReadings(readings, func(reading [2]string) bool {
		_, agree := agreedLedger(reading[0], reading[1])
		return agree
	})
	switch len(agreeing) {
	case 0:
		return "", "", unsafeLedger("git and the worktree name different ledgers, however git's report of the top level and common directory, %s, is read; check GIT_DIR, GIT_COMMON_DIR, GIT_WORK_TREE and core.worktree", output)
	case 1:
		return agreeing[0][0], agreeing[0][1], nil
	}
	return "", "", unsafeLedger("git's report of the top level and common directory, %s, reads more than one way, each naming a ledger the worktree agrees with, so fid can't tell which ledger git names", output)
}

func filterReadings(readings [][2]string, keep func([2]string) bool) [][2]string {
	var kept [][2]string
	for _, reading := range readings {
		if keep(reading) {
			kept = append(kept, reading)
		}
	}
	return kept
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
