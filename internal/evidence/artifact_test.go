package evidence

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestArtifactCandidatePathRejectsEveryUnsafeCandidate(t *testing.T) {
	workspace := t.TempDir()
	cwd := filepath.Join(workspace, "cwd")
	outside := filepath.Join(workspace, "outside")
	for _, path := range []string{cwd, outside} {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatalf("create fixture directory %q: %v", path, err)
		}
	}
	if err := os.Mkdir(filepath.Join(cwd, "directory"), 0o755); err != nil {
		t.Fatalf("create directory candidate: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cwd, "target"), []byte("safe"), 0o600); err != nil {
		t.Fatalf("write symlink target: %v", err)
	}
	if err := os.Symlink("target", filepath.Join(cwd, "symlink")); err != nil {
		t.Fatalf("create symlink candidate: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outside, "target"), []byte("outside"), 0o600); err != nil {
		t.Fatalf("write outside target: %v", err)
	}
	if err := os.Symlink("../outside", filepath.Join(cwd, "escape")); err != nil {
		t.Fatalf("create escaping parent symlink: %v", err)
	}
	if err := syscall.Mkfifo(filepath.Join(cwd, "pipe"), 0o600); err != nil {
		t.Fatalf("create non-regular candidate: %v", err)
	}

	for path, wantReason := range map[string]string{
		"directory":     "directory",
		"escape/target": "outside-cwd",
		"missing":       "missing",
		"pipe":          "non-regular",
		"symlink":       "symlink",
	} {
		_, reason := artifactCandidatePath(workspace, "cwd", path)
		if reason != wantReason {
			t.Fatalf("candidate %q rejection = %q, want %q", path, reason, wantReason)
		}
	}
}

func TestPublishArtifactCandidateRejectsChangedDuringCopy(t *testing.T) {
	workspace := t.TempDir()
	candidate := filepath.Join(workspace, "candidate.txt")
	if err := os.WriteFile(candidate, []byte("before"), 0o600); err != nil {
		t.Fatalf("write candidate: %v", err)
	}
	store := filepath.Join(workspace, "store")
	if err := os.Mkdir(store, 0o700); err != nil {
		t.Fatalf("create Artifact store: %v", err)
	}
	_, reason, err := publishArtifactCandidate(candidate, "candidate.txt", store, 1024, nil, func() {
		if err := os.WriteFile(candidate, []byte("after mutation"), 0o600); err != nil {
			t.Fatalf("mutate candidate during copy: %v", err)
		}
	})
	if err != nil {
		t.Fatalf("publish changed candidate: %v", err)
	}
	if reason != "changed-during-copy" {
		t.Fatalf("changed candidate reason = %q, want changed-during-copy", reason)
	}
}

func TestPublishArtifactCandidateRejectsOversizedContent(t *testing.T) {
	workspace := t.TempDir()
	candidate := filepath.Join(workspace, "candidate.txt")
	if err := os.WriteFile(candidate, []byte("oversized"), 0o600); err != nil {
		t.Fatalf("write candidate: %v", err)
	}
	store := filepath.Join(workspace, "store")
	if err := os.Mkdir(store, 0o700); err != nil {
		t.Fatalf("create Artifact store: %v", err)
	}
	_, reason, err := publishArtifactCandidate(candidate, "candidate.txt", store, 1, nil, nil)
	if err != nil {
		t.Fatalf("publish oversized candidate: %v", err)
	}
	if reason != "oversized" {
		t.Fatalf("oversized candidate reason = %q, want oversized", reason)
	}
}
