package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/mbell-nzl/fid/internal/display"
)

// goldenTempPathWidth leaves room for long TMPDIRs while keeping paths in
// refusal details below their 200-character limit. Each fixture child then has
// a fixed rendered path width before fid wraps it, regardless of random suffixes
// or non-ASCII parents. Use the renderer's width measure, rather than byte length.
const goldenTempPathWidth = 120

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
	// Wrapping splits at whitespace before the fixture can mask paths.
	if strings.ContainsFunc(parent, unicode.IsSpace) {
		t.Fatalf("text goldens need a TMPDIR without whitespace; got %q", parent)
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
	padding := goldenTempPathWidth - display.VisibleWidth(base) - 1
	if padding < 1 {
		t.Fatalf("golden fixture base %q exceeds the %d-column path budget; use a shorter TMPDIR or fixture parent", base, goldenTempPathWidth)
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
	for _, name := range []string{"s", strings.Repeat("long", 16), strings.Repeat("🧪", 18)} {
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
			for i, output := range outputs[1:] {
				if outputs[0] != output {
					t.Fatalf("masked output depends on temp root %q:\n--- short\n%s\n--- other\n%s", parents[i+1], outputs[0], output)
				}
			}
		})
	}
}

func TestGoldenTempDirWidthAndCleanup(t *testing.T) {
	t.Parallel()
	for i, parent := range goldenTestParents(t) {
		var directory string
		t.Run([]string{"short", "long", "non-ASCII"}[i], func(t *testing.T) {
			directory = newGoldenTempDir(t, parent)
			if width := display.VisibleWidth(directory); width != goldenTempPathWidth {
				t.Fatalf("golden directory width = %d, want %d", width, goldenTempPathWidth)
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

func TestGoldenTempDirRejectsWhitespace(t *testing.T) {
	t.Parallel()
	parent := goldenTestParents(t)[0]
	for _, tc := range []struct{ name, directory string }{
		{"space", "a b"},
		{"tab", "a\tb"},
		{"non-ASCII-space", "a\u2003b"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(parent, tc.directory)
			if err := os.Mkdir(path, 0o700); err != nil {
				t.Fatal(err)
			}
			recorder := &goldenTempFatalTB{TB: t}
			finished := make(chan struct{})
			go func() {
				defer close(finished)
				newGoldenTempDir(recorder, path)
			}()
			<-finished
			want := fmt.Sprintf("text goldens need a TMPDIR without whitespace; got %q", path)
			if recorder.message != want {
				t.Fatalf("allocator fatal message = %q, want %q", recorder.message, want)
			}
			entries, err := os.ReadDir(path)
			if err != nil || len(entries) != 0 {
				t.Fatalf("allocator created temporary entries before refusing: %v, %v", entries, err)
			}
		})
	}
}

// Record Fatalf while ending only the allocator's goroutine, so the enclosing
// test can assert the failure without marking itself failed.
type goldenTempFatalTB struct {
	testing.TB
	message string
}

func (t *goldenTempFatalTB) Fatalf(format string, args ...any) {
	t.message = fmt.Sprintf(format, args...)
	runtime.Goexit()
}
