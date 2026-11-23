package main

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func selectionGoldens(t *testing.T, f *statusFixture, name string, args ...string) {
	t.Helper()
	for _, jsonMode := range []bool{false, true} {
		command, ext := append([]string{}, args...), ".txt"
		if jsonMode {
			command, ext = append([]string{"--json"}, command...), ".json"
		}
		got := f.run(command...)
		if got.err != nil {
			t.Fatal(got)
		}
		out := regexp.MustCompile(`RND-[0-9a-f]{6}`).ReplaceAllString(got.stdout, "RND-<new>")
		builder := filepath.Join(filepath.Dir(f.configHome), "worktree")
		if f.workspace != builder {
			out = strings.ReplaceAll(out, builder, "<BUILDER>")
		}
		out = regexp.MustCompile(`\bRCL-[0-9a-f]{6}\b`).ReplaceAllString(out, "RCL-<new>")
		handoffGolden(t, f, filepath.Join("..", "work-selection", name+ext), out)
	}
}

func TestWorkSelectionGoldens(t *testing.T) {
	t.Parallel()
	f := approvalDisplayFixture(t)
	builder := f.workspace
	f.workspace = gittest.AddWorktree(t, builder)
	selectionGoldens(t, f, "primary", "next", "--explain")
	selectionGoldens(t, f, "taken", "next", "--claim", "--actor", "agent:reviewer", "--model", "claude/opus-5.5")
	lifecycleResult(t, recordReview(t, f, reviewInput(), "primary", "human:courier"), "review.record")
	selectionGoldens(t, f, "separate", "next")
	in := reviewInput()
	in["stream"] = "separate"
	in["reviewer"].(map[string]any)["session_id"] = "separate-session"
	lifecycleResult(t, recordReview(t, f, in, "separate", "human:courier"), "review.record")
	f.workspace = builder
	selectionGoldens(t, f, "approval", "status")
	approveSuccess(t, approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true))
	selectionGoldens(t, f, "merge", "status")
	for i := range 4 {
		in := changesInput("primary")
		if i > 0 {
			in["findings"] = []any{}
		}
		lifecycleResult(t, recordReview(t, f, in, "changes-"+string(rune('0'+i)), "human:courier"), "review.record")
	}
	selectionGoldens(t, f, "rework-cap", "next", "--explain")
	selectionGoldens(t, f, "rework-cap-status", "status")
}
