package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func discoveredBlockerFixture(t *testing.T) (*statusFixture, string, string, string) {
	t.Helper()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	target := epicDoneChild(t, f, epic, "target")
	lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "Original", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "original"), "epic.validation.add")
	approveTestEpic(t, f, epic)
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Discovered blocker", "--lane", "standard", "--epic", epic, "--actor", "agent:builder", "--idempotency-key", "blocker"), "issue.create")
	task := created["issue"].(map[string]any)["id"].(string)
	for _, relation := range []string{"blocks", "discovered_from"} {
		lifecycleResult(t, f.run("--json", "link", "add", task, relation, target, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:builder", "--idempotency-key", relation), "link.add")
	}
	taken := selectionNext(t, f, "agent:builder")
	if taken["kind"] != "new" || taken["task"].(map[string]any)["id"] != task {
		t.Fatalf("discovered blocker not served under current approval: %v", taken)
	}
	return f, epic, task, taken["hold"].(map[string]any)["id"].(string)
}

func staleBlockerApproval(t *testing.T, f *statusFixture, epic, action string) {
	t.Helper()
	args := []string{"--json", "epic", "validation", action, epic}
	if action != "add" {
		args = append(args, "V1")
	}
	if action != "remove" {
		args = append(args, "--claim", "Changed", "--verifier", "review:code")
	}
	lifecycleResult(t, f.run(append(args, "--actor", "agent:builder", "--idempotency-key", "stale")...), "epic.validation."+action)
}

func TestEpicStaleApprovalPreservesStartedDiscoveredBlocker(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, kind, edit string }{
		{"live resume", "resume", "add"},
		{"interrupted resume", "resume", "remove"},
		{"rework", "rework", "update"},
		{"review", "review", "remove"},
		{"judge", "judge", "add"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, epic, task, hold := discoveredBlockerFixture(t)
			lifecycleResult(t, f.run("--json", "start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start"), "task.start")
			actor := "agent:builder"
			if tc.name == "interrupted resume" {
				lifecycleResult(t, f.run("--json", "park", "--hold", hold, "--actor", actor, "--interrupted", "--reason", "Session stopped", "--idempotency-key", "park"), "hold.park")
			}
			if tc.kind != "resume" {
				runGitForTest(t, f.workspace, "checkout", "-b", "blocker")
				writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "Fix the blocker\n")
				runGitForTest(t, f.workspace, "add", ".")
				runGitForTest(t, f.workspace, "commit", "-m", "Fix discovered blocker")
				lifecycleResult(t, f.run("--json", "check", "--hold", hold, "--actor", actor, "--idempotency-key", "check"), "check.execute")
				reviewHandoffTask(t, f, "handoff", task)
			}
			if tc.kind == "rework" {
				lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
			}
			if tc.kind == "review" || tc.kind == "judge" {
				if tc.kind == "judge" {
					lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "approve-review", "human:courier", task), "review.record")
				}
				f.workspace = gittest.AddWorktree(t, f.workspace)
				actor = "agent:independent"
			}
			staleBlockerApproval(t, f, epic, tc.edit)
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if got := f.run("rebuild"); got.err != nil {
						t.Fatalf("rebuild: %#v", got)
					}
				}
				for _, claim := range []bool{false, true} {
					args := []string{"--json", "next", "--epic", epic, "--actor", actor}
					op := "next.get"
					if claim {
						args = append(args, "--claim", "--model", "openai/gpt-6.1-sol")
						op = "next.claim"
					}
					got, _ := lifecycleResult(t, f.run(args...), op)
					kind := tc.kind
					if rebuild && kind == "rework" {
						kind = "resume"
					}
					if got["kind"] != kind || got["task"].(map[string]any)["id"] != task {
						t.Fatalf("started blocker stranded after %s (rebuild=%v): %v", tc.edit, rebuild, got)
					}
				}
			}
		})
	}
}

func TestEpicStaleApprovalStillBlocksClaimedUnstartedBlocker(t *testing.T) {
	t.Parallel()
	f, epic, task, _ := discoveredBlockerFixture(t)
	staleBlockerApproval(t, f, epic, "add")
	got := selectionNext(t, f, "agent:builder")
	if got["kind"] != "none" || got["task"] != nil || got["waiting"] != nil {
		t.Fatalf("unstarted blocker served under stale approval: %v", got)
	}
	for _, raw := range got["queue"].([]any) {
		entry := raw.(map[string]any)
		if entry["id"] == task {
			if entry["reason"] != "not-accepted" || entry["working"] == nil {
				t.Fatalf("unstarted hold lost its acceptance gate or reservation: %v", entry)
			}
			return
		}
	}
	t.Fatal("unstarted blocker missing from queue")
}
