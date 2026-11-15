package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func TestDoneRejectsSupersededObservations(t *testing.T) {
	t.Parallel()
	for _, change := range []string{"takeover", "hand-off", "lane", "cancel"} {
		t.Run(change, func(t *testing.T) {
			t.Parallel()
			f, task, hold := doneBuildFixture(t, "routine", false)
			lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
			facts := doneFacts(t, f, task)
			callback := func() {
				switch change {
				case "takeover":
					lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:mark", "--reason", "New work session"), "task.takeover")
				case "hand-off":
					lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
					doneRehandoff(t, f, task, false)
				case "lane":
					lifecycleResult(t, f.run("--json", "issue", "update", task, "--lane", "assured", "--actor", "human:mark", "--idempotency-key", "raise"), "issue.update")
				case "cancel":
					lifecycleResult(t, f.run("--json", "issue", "transition", task, "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed"), "issue.transition")
				}
			}
			got := runDone(t, f, doneRunner(t, facts, callback), "--json", "done", task, "--actor", "agent:finisher")
			assertJSONFailureCode(t, got.stderr, "task.done", "done-observation-stale")
			conn, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = conn.Close() }()
			for _, table := range []string{"completions", "merge_observations", "ci_results"} {
				var n int
				if err := conn.QueryRow("SELECT count(*) FROM " + table).Scan(&n); err != nil {
					t.Fatal(err)
				}
				if n != 0 {
					t.Fatalf("stale observation wrote %s", table)
				}
			}
		})
	}
}

func TestDoneAfterTakeoverNeedsAnotherHandoff(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "routine", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	facts := doneFacts(t, f, task)
	lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:mark", "--reason", "Continue working"), "task.takeover")
	got := runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher")
	assertJSONFailureCode(t, got.stderr, "task.done", "done-handoff-required")
}

func TestDoneGitHubUnavailableKeepsKnowledgeBoundary(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "routine", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	facts := doneFacts(t, f, task)
	unavailable := func(context.Context, []string, []string) github.Result {
		return github.Result{ExitCode: 1, Stderr: []byte("network unavailable")}
	}
	args := []string{"done", task, "--actor", "agent:finisher", "--idempotency-key", "retry"}
	got := runDone(t, f, unavailable, args...)
	if !strings.Contains(got.stderr, "GitHub confirmation unavailable") || strings.Contains(got.stderr, "merged,") {
		t.Fatalf("unknown merge=%#v", got)
	}
	facts.CI.Runs = []github.WorkflowRun{}
	got = runDone(t, f, doneRunner(t, facts, nil), args...)
	if !strings.Contains(got.stderr, "waiting for CI") {
		t.Fatalf("pending=%#v", got)
	}
	got = runDone(t, f, unavailable, args...)
	if !strings.Contains(got.stderr, "merged, waiting for confirmation") || cli.ExitCode(got.err) != 4 {
		t.Fatalf("known merge=%#v", got)
	}
	facts.CI.Runs = completedCI(*facts.PR).Runs
	facts.CI.Runs[0].Conclusion = "success"
	lifecycleResult(t, runDone(t, f, doneRunner(t, facts, nil), append([]string{"--json"}, args...)...), "task.done")
}

func TestDoneResolvesQuestionWithoutAnAnswer(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "routine", false)
	lifecycleResult(t, f.run("--json", "park", "--hold", hold, "--actor", "human:builder", "--waiting", "Include retries?", "--idempotency-key", "question"), "hold.park")
	took, _ := lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:builder", "--reason", "Finish the current change"), "task.takeover")
	args := doneHandoffArgs(took["hold_id"].(string))
	lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), "task.handoff")
	facts := doneFacts(t, f, task)
	lifecycleResult(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
	shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	q := shown["questions"].([]any)[0].(map[string]any)
	if q["resolution"] != "finished-before-answer" || q["answer"] != nil || q["resolved_at"] == nil {
		t.Fatalf("question=%v", q)
	}
}

func TestDoneRefusesMissingReviewApprovalAndRemoteGates(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"review", "separate", "approval", "rollback", "tree", "artifact", "major", "lane"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			lane := "standard"
			if name == "separate" || name == "approval" || name == "rollback" {
				lane = "assured"
			}
			f, task, hold := doneBuildFixture(t, lane, false)
			lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
			code := "done-review-required"
			if name != "review" {
				input := reviewInput()
				if name == "major" {
					input = changesInput("primary")
				}
				fields, _ := lifecycleResult(t, recordReviewTask(t, f, input, "primary", "human:courier", task), "review.record")
				_ = fields
			}
			if lane == "assured" && name != "separate" {
				in := reviewInput()
				in["stream"] = "separate"
				in["reviewer"].(map[string]any)["session_id"] = "separate"
				lifecycleResult(t, recordReviewTask(t, f, in, "separate", "human:courier", task), "review.record")
			}
			if name == "rollback" {
				approveSuccess(t, approveForTest(f, task, task+"\n", true))
			}
			facts := doneFacts(t, f, task)
			switch name {
			case "separate":
				code = "done-separate-review-required"
			case "approval":
				code = "done-approval-missing"
			case "rollback":
				facts.PR.Body = "No rollback heading"
				code = "done-rollback-missing"
			case "tree":
				facts.HeadTree = strings.Repeat("c", 40)
				code = "done-handoff-required"
			case "artifact":
				if err := os.Remove(keptTaskCheckPath(t, f, task)); err != nil {
					t.Fatal(err)
				}
				code = "done-kept-artifact-unavailable"
			case "major":
				doneRehandoff(t, f, task, false)
				facts = doneFacts(t, f, task)
				code = "done-open-major"
			case "lane":
				lifecycleResult(t, f.run("--json", "issue", "update", task, "--lane", "assured", "--actor", "human:mark", "--idempotency-key", "raise"), "issue.update")
				code = "done-handoff-required"
			}
			got := runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher")
			assertJSONFailureCode(t, got.stderr, "task.done", code)
		})
	}
}

func TestDoneMissingLocalObjectsHasARepair(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "standard", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
	facts := doneFacts(t, f, task)
	head := facts.PR.HeadSHA
	if err := os.Remove(filepath.Join(f.workspace, ".git", "objects", head[:2], head[2:])); err != nil {
		t.Fatal(err)
	}
	got := runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher")
	assertJSONFailureCode(t, got.stderr, "task.done", "done-evidence-required")
}

func TestDoneGitHubFailureDoesNotConfirmAnEarlierHandoff(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "routine", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	facts := doneFacts(t, f, task)
	facts.CI.Runs = []github.WorkflowRun{}
	first := runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher")
	assertJSONFailureCode(t, first.stderr, "task.done", "done-ci-waiting")
	lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
	next, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "human:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	hold = next["hold"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "check", "--hold", hold, "--actor", "human:builder", "--idempotency-key", "recheck"), "check.execute")
	lifecycleResult(t, f.run("--json", "handoff", "--hold", hold, "--actor", "human:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/64", "--idempotency-key", "new-pr"), "task.handoff")
	unavailable := func(context.Context, []string, []string) github.Result { return github.Result{ExitCode: 1} }
	got := runDone(t, f, unavailable, "--json", "done", task, "--actor", "agent:finisher")
	assertJSONFailureCode(t, got.stderr, "task.done", "github-unavailable")
	if strings.Contains(got.stderr, `"merged":true`) || strings.Contains(got.stderr, "merged, waiting") {
		t.Fatalf("old merge confirmed new hand-off: %s", got.stderr)
	}
}
