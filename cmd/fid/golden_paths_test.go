package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGoldenTextIgnoresTempRootLength(t *testing.T) {
	t.Parallel()
	base, err := os.MkdirTemp("", "gw-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(base); err != nil {
			t.Error(err)
		}
	})
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"backup-config", []string{"status"}},
		{"lapsed-worktree", []string{"show", "FID-71a0c2"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var outputs []string
			for _, name := range []string{"s", strings.Repeat("long", 16)} {
				parent := filepath.Join(base, name)
				if err := os.MkdirAll(parent, 0o700); err != nil {
					t.Fatal(err)
				}
				f := newStatusFixtureAt(t, 100*time.Hour, parent)
				nextWork(f)
				f.hold("HLD-aaaaaa", "FID-71a0c2", "", 2*time.Hour)
				got := f.run(tc.args...)
				if got.err != nil || got.stderr != "" {
					t.Fatalf("%v = %#v", tc.args, got)
				}
				masked := f.placeholders(got.stdout)
				placeholder := "<CONFIG>"
				if tc.name == "lapsed-worktree" {
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
