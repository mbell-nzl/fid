package project

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
)

// CheckKeptWorktree checks registration and repository identity locally.
// False with no error is confirmed missing/replaced; every uncertain result
// is an error. Callers must never release kept work on either result.
func CheckKeptWorktree(root, worktree string) (bool, error) {
	common, err := filepath.EvalSymlinks(filepath.Dir(config.LedgerDir(root)))
	if err != nil {
		return false, fmt.Errorf("resolve common directory: %w", err)
	}
	output, err := keptGit("--git-dir", common, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return false, err
	}
	registered := false
	for _, field := range strings.Split(string(output), "\x00") {
		if field == "worktree "+worktree {
			registered = true
		}
	}
	if !registered {
		return false, nil
	}
	if _, err := os.Stat(worktree); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("inspect kept worktree: %w", err)
	}
	if _, err := os.Lstat(filepath.Join(worktree, ".git")); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("inspect kept git entry: %w", err)
	}
	output, err = keptGit("-C", worktree, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return false, err
	}
	actual, err := filepath.EvalSymlinks(strings.TrimSuffix(string(output), "\n"))
	if err != nil {
		return false, fmt.Errorf("resolve kept common directory: %w", err)
	}
	return actual == common, nil
}

func keptGit(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	// The check is about the path's own registration and identity, never a
	// repository selected by a calling hook's environment.
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		switch key {
		case "GIT_DIR", "GIT_COMMON_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE":
			continue
		}
		cmd.Env = append(cmd.Env, entry)
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("check kept worktree with git: %w", err)
	}
	return output, nil
}
