package main

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestDefaultConflictAfterTakeoverResumes(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		conflict bool
	}{
		{name: "clean"},
		{name: "conflict", conflict: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithLane(t, "assured")
			approveConflictReviews(t, f, task, "assured")
			if got := approveForTest(f, task, task+"\n", true); got.err != nil {
				t.Fatal(got)
			}
			moveConflictDefault(t, f, tc.conflict)
			builder := f.workspace
			f.workspace = gittest.AddWorktree(t, builder)
			taken, _ := lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:mark", "--reason", "Recover the task in another worktree"), "task.takeover")
			_, revision := taskState(t, f, task)
			before := readAllPolicyJournal(t, f.workspace)
			for _, claim := range []bool{false, true} {
				args, operation := []string{"--json", "next", "--explain"}, "next.get"
				if claim {
					args, operation = append(args, "--claim", "--actor", "human:mark"), "next.claim"
				}
				got, event := lifecycleResult(t, f.run(args...), operation)
				if got["kind"] != "resume" || got["task"].(map[string]any)["id"] != task {
					t.Errorf("claim=%t: takeover must resume, got %v", claim, got)
				}
				if got["rework"] != nil || got["rework_unavailable"] != nil {
					t.Errorf("claim=%t: takeover got conflict evidence: %v", claim, got)
				}
				if event != nil {
					t.Errorf("claim=%t: takeover resume appended an Event: %v", claim, event)
				}
				if claim && got["hold"].(map[string]any)["id"] != taken["hold_id"] {
					t.Errorf("takeover resume replaced the live hold: %v", got["hold"])
				}
			}
			if state, rev := taskState(t, f, task); state != "in_review" || rev != revision {
				t.Errorf("takeover resume changed task state/revision: %s/%v", state, rev)
			}
			f.workspace = builder
			old, event := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
			if old["kind"] != "none" || old["rework"] != nil || event != nil {
				t.Errorf("released builder reservation got work: %v; event=%v", old, event)
			}
			if readAllPolicyJournal(t, f.workspace) != before {
				t.Error("next after takeover changed journal bytes")
			}
		})
	}
}
