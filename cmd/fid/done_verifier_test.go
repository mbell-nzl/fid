package main

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// V5 / check:done-refusals. Each journey uses public commands for all ledger
// writes; GitHub alone is replaced by the strict offline gh executable.
func TestV5DoneRefusesWithoutReceiptsAndAcceptsMergedGreenWork(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		run  func(*testing.T)
	}{
		{"missing required check", v5MissingCheck},
		{"removed kept Artifact", v5MissingArtifact},
		{"no merge commit CI", v5NoCI},
		{"major resolved by later round", func(t *testing.T) { v5OpenMajor(t, false) }},
		{"major closed by person", func(t *testing.T) { v5OpenMajor(t, true) }},
		{"stale Assured approval", v5StaleApproval},
		{"unmerged PR", func(t *testing.T) { v5RemoteGate(t, "unmerged") }},
		{"failed merge CI", func(t *testing.T) { v5RemoteGate(t, "failed") }},
		{"pending merge CI", func(t *testing.T) { v5RemoteGate(t, "pending") }},
		{"another repository", func(t *testing.T) { v5RemoteGate(t, "repository") }},
		{"another branch", func(t *testing.T) { v5RemoteGate(t, "branch") }},
		{"transition bypasses", v5TransitionBypasses},
		{"agent epic close", v5AgentEpicClose},
		{"reconciliation uses the same gates", v5Reconcile},
		{"Routine complete", func(t *testing.T) { v5CompleteLane(t, "routine", false) }},
		{"Routine tests complete", func(t *testing.T) { v5CompleteLane(t, "routine", true) }},
		{"Standard complete", func(t *testing.T) { v5CompleteLane(t, "standard", true) }},
		{"Assured complete", func(t *testing.T) { v5CompleteLane(t, "assured", true) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.run(t)
		})
	}
}

func v5MissingCheck(t *testing.T) {
	t.Helper()
	j := newV5Journey(t, "routine", false)
	j.refuse(t, j.task, "task.handoff", "checks-unsatisfied", 4, func() commandExecutionResult {
		return j.f.run(j.handoffArgs("missing-check")...)
	})
	j.refuseDone(t, "done-task-not-in-review", "task-not-in-review")
	j.assertState(t, "in_progress")
	j.checkAndHandoff(t, "checked")
	j.finish(t)
}

func v5MissingArtifact(t *testing.T) {
	t.Helper()
	j := newV5Journey(t, "standard", false)
	j.checkAndHandoff(t, "handoff")
	j.review(t, "primary", "primary", reviewInput())
	path := j.keptArtifact(t)
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	j.removeArtifact(t)
	j.refuseDone(t, "done-kept-artifact-unavailable", "kept-artifact-unavailable")
	writeBackupTestFile(t, path, string(bytes))
	j.finish(t)
}

func v5NoCI(t *testing.T) {
	t.Helper()
	j := newV5Journey(t, "routine", false)
	j.checkAndHandoff(t, "handoff")
	j.facts.CI.Runs = []github.WorkflowRun{}
	j.stub.add(j.facts)
	got := j.refuseDone(t, "done-ci-waiting", "ci-waiting")
	if !strings.Contains(got.stderr, "waiting for CI") {
		t.Fatalf("missing wait guidance: %s", got.stderr)
	}
	j.facts.CI.Runs = completedCI(*j.facts.PR).Runs
	j.facts.CI.Runs[0].Conclusion = "success"
	j.stub.add(j.facts)
	j.finish(t)
}

func v5OpenMajor(t *testing.T, person bool) {
	t.Helper()
	j := newV5Journey(t, "standard", false)
	j.checkAndHandoff(t, "handoff")
	round := j.review(t, "primary", "changes", changesInput("primary"))
	finding := round["findings"].([]any)[0].(map[string]any)["id"].(string)
	j.assertState(t, "in_progress")
	j.rework(t, false) // Same change keeps the round current; its major stays open.
	j.refuseDone(t, "done-open-major", "open-major")
	if person {
		lifecycleResult(t, closeFinding(j.f, finding, "human:mark", "Reviewed and accepted the fix", "close"), "finding.close")
	} else {
		in := reviewInput()
		in["dispositions"] = []any{disposition(finding, "resolved")}
		j.review(t, "primary", "resolved", in)
	}
	j.finish(t)
}

func v5StaleApproval(t *testing.T) {
	t.Helper()
	j := newV5Journey(t, "assured", false)
	j.checkAndHandoff(t, "handoff")
	j.review(t, "primary", "primary", reviewInput())
	j.review(t, "separate", "separate", reviewInput())
	j.approve(t, "first-approval")
	round := j.review(t, "primary", "changes", changesInput("primary"))
	finding := round["findings"].([]any)[0].(map[string]any)["id"].(string)
	j.rework(t, true)
	in := reviewInput()
	in["dispositions"] = []any{disposition(finding, "resolved")}
	j.review(t, "primary", "fixed-primary", in)
	j.review(t, "separate", "fixed-separate", reviewInput())
	j.refuseDone(t, "done-approval-stale", "approval-stale")
	j.approve(t, "second-approval")
	j.finish(t)
}

func v5RemoteGate(t *testing.T, state string) {
	t.Helper()
	j := newV5Journey(t, "standard", false)
	j.checkAndHandoff(t, "handoff")
	j.review(t, "primary", "primary", reviewInput())
	code, reason := "done-ci-waiting", "ci-waiting"
	switch state {
	case "unmerged":
		j.facts.PR.Merged, j.facts.PR.State = false, "OPEN"
		code, reason = "done-not-merged", "not-merged"
	case "failed":
		j.facts.CI.Runs[0].Conclusion = "failure"
		code, reason = "done-ci-not-green", "ci-not-green"
	case "pending":
		j.facts.CI.Runs[0].Status, j.facts.CI.Runs[0].Conclusion = "in_progress", ""
	case "repository":
		j.facts.PR.URL = "https://github.com/another/project/pull/63"
		code, reason = "done-handoff-required", "repository-mismatch"
	case "branch":
		j.facts.PR.BaseRef = "another-branch"
		code, reason = "done-handoff-required", "branch-mismatch"
	}
	j.stub.add(j.facts)
	j.refuseDone(t, code, reason)
	// A later successful attempt on this exact merge commit completes the task.
	j.facts.PR.Merged, j.facts.PR.State = true, "MERGED"
	j.facts.PR.URL, j.facts.PR.BaseRef = "https://github.com/mbell-nzl/fid/pull/63", "trunk"
	j.facts.CI.Runs[0].Attempt++
	j.facts.CI.Runs[0].Status, j.facts.CI.Runs[0].Conclusion = "completed", "success"
	j.stub.add(j.facts)
	j.finish(t)
}

func v5TransitionBypasses(t *testing.T) {
	t.Helper()
	for _, kind := range []string{"task", "epic"} {
		for _, actor := range []string{"agent:builder", "human:mark", "service:ci"} {
			t.Run(kind+"/"+actor, func(t *testing.T) {
				t.Parallel()
				j := newV5Journey(t, "routine", false)
				j.checkAndHandoff(t, "handoff")
				id := j.task
				if kind == "epic" {
					created, _ := lifecycleResult(t, j.f.run("--json", "issue", "create", "--kind", "epic", "--title", "V5 epic", "--scheduling", "prioritised", "--actor", "human:mark", "--idempotency-key", "epic"), "issue.create")
					id = created["issue"].(map[string]any)["id"].(string)
				}
				got := j.refuse(t, id, "issue.transition", "transition-not-permitted", 4, func() commandExecutionResult {
					return j.f.run("--json", "issue", "transition", id, "--to", "done", "--actor", actor, "--idempotency-key", "bypass")
				})
				guidance := id + " completes with fid done"
				if kind == "epic" {
					guidance = "fid done " + id
				}
				if !strings.Contains(got.stderr, "use-done") || !strings.Contains(got.stderr, guidance) {
					t.Fatalf("bypass refusal must direct every actor to done: %s", got.stderr)
				}
			})
		}
	}
}

func v5AgentEpicClose(t *testing.T) {
	t.Helper()
	j := newV5Journey(t, "routine", false)
	created, _ := lifecycleResult(t, j.f.run("--json", "issue", "create", "--kind", "epic", "--title", "V5 eligible epic", "--scheduling", "prioritised", "--actor", "human:mark", "--idempotency-key", "epic"), "issue.create")
	id := created["issue"].(map[string]any)["id"].(string)
	j.refuse(t, id, "task.done", "authority-denied", 5, func() commandExecutionResult {
		return j.stub.run("--json", "done", id, "--actor", "agent:builder", "--reason", "All tasks finished")
	})
	approveTestEpic(t, j.f, id)
	closed, _ := lifecycleResult(t, j.stub.run("--json", "done", id, "--actor", "human:mark", "--reason", "All tasks finished"), "task.done")
	if closed["status"] != "done" || closed["decision"].(map[string]any)["kind"] != "close-epic" {
		t.Fatalf("person's epic close=%v", closed)
	}
}

func v5CompleteLane(t *testing.T, lane string, testsChanged bool) {
	t.Helper()
	j := newV5Journey(t, lane, testsChanged)
	j.checkAndHandoff(t, "handoff")
	if lane != "routine" || testsChanged {
		j.refuseDone(t, "done-review-required", "review-missing")
		in := reviewInput()
		if lane == "routine" {
			in["profiles"] = []string{"test-adequacy"}
		}
		j.review(t, "primary", "primary", in)
	}
	if lane == "assured" {
		j.refuseDone(t, "done-separate-review-required", "review-missing")
		j.review(t, "separate", "separate", reviewInput())
		j.refuseDone(t, "done-approval-missing", "approval-missing")
		j.approve(t, "approval")
	}
	j.finish(t)
}

func v5Reconcile(t *testing.T) {
	t.Helper()
	blocked := newV5Journey(t, "standard", false)
	blocked.checkAndHandoff(t, "handoff")
	// Another worktree has a complete task in the same ledger. It requires no
	// review because it is Routine without changed test paths.
	other := *blocked.f
	other.workspace = gittest.AddWorktree(t, blocked.f.workspace)
	ready := &v5Journey{f: &other, prNumber: 64, stub: blocked.stub}
	created, _ := lifecycleResult(t, other.run("--json", "issue", "create", "--kind", "task", "--title", "V5 ready task", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "ready"), "issue.create")
	ready.task = created["issue"].(map[string]any)["id"].(string)
	taken, _ := lifecycleResult(t, other.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	ready.hold = taken["hold"].(map[string]any)["id"].(string)
	lifecycleResult(t, other.run("--json", "start", "--hold", ready.hold, "--actor", "agent:builder", "--idempotency-key", "ready-start"), "task.start")
	ready.checkAndHandoff(t, "ready-handoff")
	blocked.refuseDone(t, "done-review-required", "review-missing")
	before := blocked.state(t, blocked.task)
	outcomes := assertReconciled(t, blocked.stub.run("--json", "done", "--actor", "agent:finisher"), []string{blocked.task, ready.task}, []string{"refused", "completed"})
	if outcomes[blocked.task]["refusal"].(map[string]any)["code"] != "done-review-required" || !reflect.DeepEqual(before, blocked.state(t, blocked.task)) {
		t.Fatalf("reconciliation bypassed review: %v", outcomes)
	}
	ready.assertState(t, "done")
	blocked.assertState(t, "in_review")
	blocked.review(t, "primary", "primary", reviewInput())
	blocked.finish(t)
	assertReconciled(t, blocked.stub.run("--json", "done", "--actor", "agent:finisher"), nil, nil)
}
