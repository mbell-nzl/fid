package artifactstore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKeepFilesSafelyAndWithoutChangingBinaryContent(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ name, reason string }{
		{"kept", ""}, {"missing", "missing"}, {"../outside", "outside-worktree"},
		{"link", "outside-worktree"}, {"large", "oversized"}, {"directory", "non-regular"},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			root, store := t.TempDir(), t.TempDir()
			data := []byte{0, 255, 254, '\n'}
			if err := os.WriteFile(filepath.Join(root, "kept"), data, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(filepath.Join(t.TempDir(), "outside"), filepath.Join(root, "link")); err != nil {
				t.Fatal(err)
			}
			file, err := os.Create(filepath.Join(root, "large"))
			if err != nil {
				t.Fatal(err)
			}
			if err := file.Truncate((64 << 20) + 1); err != nil {
				t.Fatal(err)
			}
			if err := file.Close(); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(filepath.Join(root, "directory"), 0o700); err != nil {
				t.Fatal(err)
			}
			got, reason, err := Keep(root, store, row.name, 64<<20)
			if err != nil || reason != row.reason {
				t.Fatalf("keep=%#v/%s/%v", got, reason, err)
			}
			if reason == "" {
				saved, err := os.ReadFile(filepath.Join(store, got.Digest[len("sha256:"):]))
				if err != nil || string(saved) != string(data) {
					t.Fatalf("saved=%v/%v", saved, err)
				}
			}
		})
	}
}

func TestKeepRefusesAFileChangedDuringCopy(t *testing.T) {
	t.Parallel()
	root, store := t.TempDir(), t.TempDir()
	path := filepath.Join(root, "artifact")
	if err := os.WriteFile(path, []byte("before"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, reason, err := keep(root, store, "artifact", 64<<20, func() {
		if err := os.WriteFile(path, []byte("changed"), 0o600); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil || reason != "changed-during-copy" {
		t.Fatalf("changed file=%s/%v", reason, err)
	}
	entries, err := os.ReadDir(store)
	if err != nil || len(entries) != 0 {
		t.Fatalf("refused bytes were published: %v/%v", entries, err)
	}
}
