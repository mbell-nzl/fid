// Package project resolves fid project roots without crossing repository boundaries.
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

// Resolve returns the selected fid project root. An explicit path must directly
// contain .fid; discovery otherwise walks from start up to its Git boundary.
func Resolve(start, explicit string) (string, error) {
	if strings.TrimSpace(explicit) != "" {
		root, err := absoluteDirectory(explicit)
		if err != nil {
			return "", err
		}
		if isProjectRoot(root) {
			return root, nil
		}
		return "", fmt.Errorf("%w: %s", ErrNotFound, root)
	}

	directory, err := absoluteDirectory(start)
	if err != nil {
		return "", err
	}
	for {
		if isProjectRoot(directory) {
			return directory, nil
		}
		if containsGitBoundary(directory) {
			break
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			break
		}
		directory = parent
	}
	return "", fmt.Errorf("%w: %s", ErrNotFound, start)
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
	info, err := os.Lstat(filepath.Join(root, config.WorkspaceDirName))
	return err == nil && info.IsDir() && info.Mode()&fs.ModeSymlink == 0
}

func containsGitBoundary(root string) bool {
	_, err := os.Lstat(filepath.Join(root, ".git"))
	return err == nil
}
