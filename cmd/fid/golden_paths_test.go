package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// goldenTempPathLength leaves room for long TMPDIRs while keeping paths in
// refusal details below their 200-character limit. Each fixture child then has
// a fixed total path length before fid wraps it, regardless of random suffixes.
const goldenTempPathLength = 120

func newGoldenTempDir(t testing.TB, parent string) string {
	t.Helper()
	if parent == "" {
		parent = os.TempDir()
	}
	// Git reports resolved absolute paths, so budget that spelling too.
	parent, err := filepath.Abs(parent)
	if err != nil {
		t.Fatal(err)
	}
	parent, err = filepath.EvalSymlinks(parent)
	if err != nil {
		t.Fatal(err)
	}
	base, err := os.MkdirTemp(parent, "g-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(base); err != nil {
			t.Error(err)
		}
	})
	padding := goldenTempPathLength - len(base) - 1
	if padding < 1 {
		t.Fatalf("golden fixture base %q exceeds the %d-character path budget; use a shorter TMPDIR or fixture parent", base, goldenTempPathLength)
	}
	path := filepath.Join(base, strings.Repeat("p", padding))
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

func goldenTestParents(t *testing.T) []string {
	t.Helper()
	// Explicit roots keep this regression's short/long comparison meaningful
	// even when the suite itself runs under an unusually long TMPDIR.
	base, err := os.MkdirTemp("/tmp", "gw-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(base); err != nil {
			t.Error(err)
		}
	})
	var parents []string
	for _, name := range []string{"s", strings.Repeat("long", 16)} {
		parent := filepath.Join(base, name)
		if err := os.MkdirAll(parent, 0o700); err != nil {
			t.Fatal(err)
		}
		parents = append(parents, parent)
	}
	return parents
}

func TestGoldenTextIgnoresTempRootLength(t *testing.T) {
	t.Parallel()
	parents := goldenTestParents(t)
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"backup-config", []string{"status"}},
		{"backup-home", []string{"status"}},
		{"lapsed-worktree", []string{"show", "FID-71a0c2"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var outputs []string
			for _, parent := range parents {
				f := newStatusFixtureAt(t, 100*time.Hour, parent, "worktree")
				if tc.name == "backup-home" {
					f.configHome = ""
				}
				nextWork(f)
				f.hold("HLD-aaaaaa", "FID-71a0c2", "", 2*time.Hour)
				got := f.run(tc.args...)
				if got.err != nil || got.stderr != "" {
					t.Fatalf("%v = %#v", tc.args, got)
				}
				masked := f.placeholders(got.stdout)
				placeholder := "<CONFIG>"
				switch tc.name {
				case "backup-home":
					placeholder = "<HOME>"
				case "lapsed-worktree":
					placeholder = "<ROOT>"
				}
				if !strings.Contains(masked, placeholder) {
					t.Fatalf("missing %s in output:\n%s", placeholder, masked)
				}
				outputs = append(outputs, masked)
			}
			if outputs[0] != outputs[1] {
				t.Fatalf("masked output depends on temp root length:\n--- short\n%s\n--- long\n%s", outputs[0], outputs[1])
			}
		})
	}
}

func TestGoldenTempDirLengthAndCleanup(t *testing.T) {
	t.Parallel()
	for i, parent := range goldenTestParents(t) {
		var directory string
		t.Run([]string{"short", "long"}[i], func(t *testing.T) {
			directory = newGoldenTempDir(t, parent)
			if len(directory) != goldenTempPathLength {
				t.Fatalf("golden directory length = %d, want %d", len(directory), goldenTempPathLength)
			}
			if _, err := os.Stat(directory); err != nil {
				t.Fatal(err)
			}
		})
		if _, err := os.Stat(filepath.Dir(directory)); !os.IsNotExist(err) {
			t.Fatalf("golden directory's private parent survived cleanup: %v", err)
		}
	}
}
