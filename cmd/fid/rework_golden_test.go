package main

import (
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
			output := regexp.MustCompile(`\b[0-9a-f]{8}\b`).ReplaceAllString(got.stdout, "<commit>")
			handoffGolden(t, f, filepath.Join("..", "rework", tc.name+ext), output)
		}
	}
}
