package main

import (
	"strings"
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
			assertTakeoverConflictStatus(t, f, task, tc.conflict)
			builder := f.workspace
			f.workspace = gittest.AddWorktree(t, builder)
			taken, _ := lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:mark", "--reason", "Recover the task in another worktree"), "task.takeover")
			assertTakeoverConflictStatus(t, f, task, false)
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
			assertTakeoverConflictStatus(t, f, task, false)
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

func assertTakeoverConflictStatus(t *testing.T, f *statusFixture, task string, conflict bool) {
	t.Helper()
	status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	wantKind, wantAction := "merge", "https://github.com/mbell-nzl/fid/pull/63"
	if conflict {
		wantKind, wantAction = "default-conflict", "fid next --claim in the builder's worktree"
	}
	items := 0
	for _, value := range status["needs_you"].([]any) {
		item := value.(map[string]any)
		subject, ok := item["subject"].(map[string]any)
		if !ok || subject["id"] != task {
			continue
		}
		items++
		if item["kind"] != wantKind || item["next_action"] != wantAction {
			t.Errorf("status action for task=%v, want kind=%s action=%s", item, wantKind, wantAction)
		}
	}
	if items != 1 {
		t.Errorf("status has %d actions for task, want 1: %v", items, status)
	}
	text := f.run("status")
	if text.err != nil {
		t.Fatal(text)
	}
	if conflict {
		if !strings.Contains(text.stdout, wantAction) || !strings.Contains(text.stdout, "PR no longer merges") {
			t.Errorf("status text must name the builder's fid next --claim action:\n%s", text.stdout)
		}
	} else if strings.Contains(text.stdout, "PR no longer merges") || strings.Contains(text.stdout, "fid next --claim") || !strings.Contains(text.stdout, task+": merge PR #63") {
		t.Errorf("status text must keep the existing merge action without a builder update:\n%s", text.stdout)
	}
}
