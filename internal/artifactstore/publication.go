package artifactstore

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Publish(stagingPath, destination, digest string) error {
	info, err := os.Lstat(destination)
	if err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("retain Artifact: existing digest path is not regular")
		}
		return Verify(destination, digest)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("retain Artifact: inspect digest destination: %w", err)
	}
	if err := os.Rename(stagingPath, destination); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("retain Artifact: publish by digest: %w", err)
		}
		return Verify(destination, digest)
	}
	return nil
}

func Verify(path, digest string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("verify Artifact digest: inspect: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("verify Artifact digest: stored content is not regular")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("verify Artifact digest: open: %w", err)
	}
	defer func() { _ = file.Close() }()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("verify Artifact digest: read: %w", err)
	}
	if "sha256:"+fmt.Sprintf("%x", hash.Sum(nil)) != digest {
		return fmt.Errorf("verify Artifact digest: stored content does not match digest")
	}
	return nil
}

func Contained(root, candidate string) bool {
	relative, err := filepath.Rel(root, candidate)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

// Root resolves and contains the configured Artifact store inside the ledger.
func Root(ledgerRoot, root string) (string, error) {
	fidRoot, err := filepath.EvalSymlinks(ledgerRoot)
	if err != nil {
		return "", fmt.Errorf("retain Artifacts: resolve fid root: %w", err)
	}
	store := filepath.Join(fidRoot, root)
	if err := os.MkdirAll(store, 0o700); err != nil {
		return "", fmt.Errorf("retain Artifacts: create store: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(store)
	if err != nil || !Contained(fidRoot, resolved) {
		return "", fmt.Errorf("retain Artifacts: configured store escapes fid root")
	}
	return resolved, nil
}
