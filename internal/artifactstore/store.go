// Package artifactstore publishes immutable bytes by content digest. Both the
// old executor and lean checks share its publication and containment helpers.
package artifactstore

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type File struct {
	Digest    string
	SizeBytes int64
}

// Keep uses an os.Root for source traversal, including during revalidation:
// a symlink swap cannot cause even a temporary read outside the worktree.
func Keep(worktree, store, declared string, limit int64) (File, string, error) {
	return keep(worktree, store, declared, limit, nil)
}

func keep(worktree, store, declared string, limit int64, beforeRevalidate func()) (File, string, error) {
	if declared == "" || filepath.IsAbs(declared) || !filepath.IsLocal(declared) || strings.ContainsRune(declared, 0) {
		return File{}, "outside-worktree", nil
	}
	root, err := os.OpenRoot(worktree)
	if err != nil {
		return File{}, "outside-worktree", nil
	}
	defer func() { _ = root.Close() }()
	before, err := root.Stat(declared)
	if errors.Is(err, os.ErrNotExist) {
		return File{}, "missing", nil
	}
	if err != nil {
		return File{}, "outside-worktree", nil
	}
	if !before.Mode().IsRegular() {
		return File{}, "non-regular", nil
	}
	input, err := root.OpenFile(declared, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if errors.Is(err, os.ErrNotExist) {
		return File{}, "missing", nil
	}
	if err != nil {
		return File{}, "outside-worktree", nil
	}
	defer func() { _ = input.Close() }()
	opened, err := input.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
		return File{}, "non-regular", nil
	}
	if before.Size() > limit {
		return File{}, "oversized", nil
	}
	staging, err := os.CreateTemp(store, ".artifact-")
	if err != nil {
		return File{}, "", fmt.Errorf("keep file: create staging: %w", err)
	}
	defer func() { _ = os.Remove(staging.Name()) }()
	defer func() { _ = staging.Close() }()
	hash := sha256.New()
	size, err := io.Copy(io.MultiWriter(staging, hash), io.LimitReader(input, limit+1))
	if err != nil {
		return File{}, "", fmt.Errorf("keep file: copy: %w", err)
	}
	if size > limit {
		return File{}, "oversized", nil
	}
	if beforeRevalidate != nil {
		beforeRevalidate()
	}
	after, err := root.Stat(declared)
	if err != nil || !os.SameFile(before, after) || before.Size() != size || after.Size() != size || before.Mode() != after.Mode() || !before.ModTime().Equal(after.ModTime()) {
		return File{}, "changed-during-copy", nil
	}
	if err := staging.Sync(); err != nil {
		return File{}, "", fmt.Errorf("keep file: sync: %w", err)
	}
	if err := staging.Close(); err != nil {
		return File{}, "", fmt.Errorf("keep file: close: %w", err)
	}
	digest := fmt.Sprintf("sha256:%x", hash.Sum(nil))
	if err := Publish(staging.Name(), filepath.Join(store, strings.TrimPrefix(digest, "sha256:")), digest); err != nil {
		return File{}, "", err
	}
	return File{Digest: digest, SizeBytes: size}, "", nil
}
