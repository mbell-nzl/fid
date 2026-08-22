package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveStopsAtContainingGitRoot(t *testing.T) {
	outside := t.TempDir()
	if err := os.Mkdir(filepath.Join(outside, ".fid"), 0o755); err != nil {
		t.Fatalf("Mkdir(.fid) error = %v", err)
	}
	repository := filepath.Join(outside, "repository")
	start := filepath.Join(repository, "nested", "directory")
	if err := os.MkdirAll(start, 0o755); err != nil {
		t.Fatalf("MkdirAll(start) error = %v", err)
	}
	if err := os.Mkdir(filepath.Join(repository, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}

	_, err := Resolve(start, "")
	if err == nil {
		t.Fatal("Resolve() error = nil, want discovery to stop at containing Git root")
	}
}

func TestResolveFindsNearestRootAndRequiresExplicitRootToContainFID(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".fid"), 0o755); err != nil {
		t.Fatalf("Mkdir(root .fid) error = %v", err)
	}
	nestedRoot := filepath.Join(root, "nested")
	if err := os.MkdirAll(filepath.Join(nestedRoot, ".fid"), 0o755); err != nil {
		t.Fatalf("MkdirAll(nested .fid) error = %v", err)
	}
	start := filepath.Join(nestedRoot, "child")
	if err := os.Mkdir(start, 0o755); err != nil {
		t.Fatalf("Mkdir(start) error = %v", err)
	}

	got, err := Resolve(start, "")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got != nestedRoot {
		t.Fatalf("Resolve() = %q, want nearest root %q", got, nestedRoot)
	}
	if _, err := Resolve(start, filepath.Join(root, "nested", "child")); err == nil {
		t.Fatal("Resolve(explicit child) error = nil, want direct-root rejection")
	}
}

func TestResolveRejectsSymlinkedFIDDirectory(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Mkdir(filepath.Join(outside, ".fid"), 0o755); err != nil {
		t.Fatalf("Mkdir(outside .fid) error = %v", err)
	}
	if err := os.Symlink(filepath.Join(outside, ".fid"), filepath.Join(root, ".fid")); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}

	if _, err := Resolve(root, ""); err == nil {
		t.Fatal("Resolve() error = nil, want symlink escape rejection")
	}
}
