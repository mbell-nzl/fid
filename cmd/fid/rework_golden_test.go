package main

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestReworkGolden(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	lifecycleResult(t, recordReview(t, f, changesInput("primary"), "changes", "human:courier"), "review.record")
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"next", []string{"next"}},
		{"show", []string{"show", "FID-9f13aa"}},
		{"status", []string{"status"}},
		{"taken", []string{"next", "--claim", "--actor", "agent:fixer", "--model", "codex/gpt-6.1-sol"}},
	} {
		for _, jsonMode := range []bool{false, true} {
			args, ext := tc.args, ".txt"
			if jsonMode {
				args, ext = append([]string{"--json"}, args...), ".json"
			}
			// Read the taken hold before taking it for text, so both renderers
			// describe the same rework return (repeating next is a resume).
			if tc.name == "taken" && !jsonMode {
				continue
			}
			got := f.run(args...)
			if got.err != nil || got.stderr != "" {
				t.Fatal(got)
			}
			handoffGolden(t, f, filepath.Join("..", "rework", tc.name+ext), reworkGoldenOutput(f, got.stdout))
		}
	}
}

// shortCommit matches the abbreviated commits rework output names.
var shortCommit = regexp.MustCompile(`\b[0-9a-f]{8}\b`)

// reworkGoldenOutput names the fixture's paths before masking short commits:
// a temp directory whose random name is eight digits would otherwise lose
// part of its path to the commit mask and never match <ROOT>.
func reworkGoldenOutput(f *statusFixture, stdout string) string {
	return shortCommit.ReplaceAllString(f.placeholders(stdout), "<commit>")
}

func TestReworkGoldenOutputMasksRootsBeforeCommits(t *testing.T) {
	t.Parallel()
	workspace := filepath.Join(t.TempDir(), "12345678", "worktree")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	f := &statusFixture{t: t, workspace: workspace, configHome: filepath.Join(t.TempDir(), "config")}
	got := reworkGoldenOutput(f, `{"worktree":"`+workspace+`","head":"0123abcd"}`)
	if want := `{"worktree":"<ROOT>","head":"<commit>"}`; got != want {
		t.Fatalf("reworkGoldenOutput = %s, want %s", got, want)
	}
}
