package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/doneeval"
)

func doneRehandoff(t *testing.T, f *statusFixture, task string, edit bool) {
	t.Helper()
	next, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "human:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if next["kind"] != "rework" || next["task"].(map[string]any)["id"] != task {
		t.Fatalf("rework=%v", next)
	}
	hold := next["hold"].(map[string]any)["id"].(string)
	if edit {
		writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "Edited after review\n")
		runGitForTest(t, f.workspace, "add", ".")
		runGitForTest(t, f.workspace, "commit", "-qm", "Fix the change")
	}
	lifecycleResult(t, f.run("--json", "check", "--hold", hold, "--actor", "human:builder", "--idempotency-key", "recheck"), "check.execute")
	args := doneHandoffArgs(hold)
	args[len(args)-1] = "rehandoff"
	lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), "task.handoff")
}

func TestDoneGatherMajorClosureAuthorities(t *testing.T) {
	t.Parallel()
	for _, authority := range []string{"later round", "person"} {
		t.Run(authority, func(t *testing.T) {
			t.Parallel()
			f, task, hold := doneBuildFixture(t, "standard", false)
			lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
			fields, _ := lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
			finding := fields["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)["id"].(string)
			target, remote := doneRemoteFixture(t, f, task)
			open := doneeval.Evaluate(gatherDone(t, f, task, target, remote))
			if open.Gates[doneeval.MajorsGate].Reason != doneeval.OpenMajor {
				t.Fatalf("open=%+v", open)
			}
			if authority == "person" {
				lifecycleResult(t, f.run("--json", "finding", "close", finding, "--actor", "human:mark", "--reason", "Resolved after inspection", "--idempotency-key", "close"), "finding.close")
			}
			doneRehandoff(t, f, task, false)
			if authority == "later round" {
				in := reviewInput()
				in["dispositions"] = []any{disposition(finding, "resolved")}
				lifecycleResult(t, recordReviewTask(t, f, in, "resolved", "human:courier", task), "review.record")
			}
			target, remote = doneRemoteFixture(t, f, task)
			s := gatherDone(t, f, task, target, remote)
			got := doneeval.Evaluate(s)
			if got.Status != doneeval.Complete || len(s.OpenMajors) != 0 || s.Rounds[0].Findings[0].Closure == nil {
				t.Fatalf("closed=%+v snapshot=%+v", got, s)
			}
			if authority == "later round" && s.Rounds[0].Findings[0].Closure.RoundID == "" {
				t.Fatal("lost round closure authority")
			}
			if authority == "person" && s.Rounds[0].Findings[0].Closure.DecisionID == "" {
				t.Fatal("lost person closure authority")
			}
		})
	}
}

func TestDoneGatherEditStalesReviewsApprovalAndOtherModel(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "assured", true)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "first-primary", "human:courier", task), "review.record")
	in := reviewInput()
	in["stream"] = "separate"
	in["reviewer"].(map[string]any)["session_id"] = "first-separate"
	lifecycleResult(t, recordReviewTask(t, f, in, "first-separate", "human:courier", task), "review.record")
	approveSuccess(t, approveForTest(f, task, task+"\n", true, "--idempotency-key", "before-edit"))
	fields, _ := lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
	finding := fields["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)["id"].(string)
	doneRehandoff(t, f, task, true)
	target, remote := doneRemoteFixture(t, f, task)
	stale := doneeval.Evaluate(gatherDone(t, f, task, target, remote))
	if stale.Status != doneeval.Refused || stale.Gates[doneeval.PrimaryGate].Reason != doneeval.ReviewStale || stale.Gates[doneeval.SeparateGate].Reason != doneeval.ReviewStale || stale.Gates[doneeval.ApprovalGate].Reason != doneeval.ApprovalStale || stale.Gates[doneeval.MajorsGate].Reason != doneeval.OpenMajor || !stale.NoOtherModelReview {
		t.Fatalf("stale=%+v", stale)
	}
	lifecycleResult(t, f.run("--json", "finding", "close", finding, "--actor", "human:mark", "--reason", "Fix verified", "--idempotency-key", "close"), "finding.close")
	for _, stream := range []string{"primary", "separate"} {
		in := reviewInput()
		in["stream"] = stream
		in["reviewer"].(map[string]any)["session_id"] = "new-" + stream
		in["reviewer"].(map[string]any)["model"] = "codex/gpt-6.1-sol"
		lifecycleResult(t, recordReviewTask(t, f, in, "new-"+stream, "human:courier", task), "review.record")
	}
	stillStale := doneeval.Evaluate(gatherDone(t, f, task, target, remote))
	if stillStale.Status != doneeval.Refused || stillStale.Gates[doneeval.ApprovalGate].Reason != doneeval.ApprovalStale || !stillStale.NoOtherModelReview {
		t.Fatalf("new reviews=%+v", stillStale)
	}
	approveSuccess(t, approveForTest(f, task, task+"\n", true, "--idempotency-key", "after-edit"))
	got := doneeval.Evaluate(gatherDone(t, f, task, target, remote))
	if got.Status != doneeval.Complete || !got.NoOtherModelReview {
		t.Fatalf("reapproved=%+v", got)
	}
}

func TestDoneGatherUsesFrozenPolicyAndRefusesRaisedLane(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "standard", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
	target, remote := doneRemoteFixture(t, f, task)
	original := gatherDone(t, f, task, target, remote)
	path := filepath.Join(f.workspace, ".fid.toml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	changed := strings.Replace(string(raw), "[lanes.standard]\nchecks=['evidence']", "[lanes.standard]\nchecks=['retroactive']", 1)
	writeBackupTestFile(t, path, changed+"\n[checks.retroactive]\nrun=['false']\ntimeout='1m'\n")
	policyCommit(t, f.workspace)
	identityGit(t, f.workspace, "push", "-q", "origin", "HEAD:trunk")
	still := gatherDone(t, f, task, target, remote)
	if got := doneeval.Evaluate(still); got.Status != doneeval.Complete || !reflect.DeepEqual(still.Policy, original.Policy) {
		t.Fatalf("later policy changed gates: %+v", got)
	}
	lifecycleResult(t, f.run("--json", "issue", "update", task, "--lane", "assured", "--actor", "human:mark", "--idempotency-key", "raise"), "issue.update")
	got := doneeval.Evaluate(gatherDone(t, f, task, target, remote))
	if got.Status != doneeval.Refused || got.Gates[doneeval.LaneGate].Reason != doneeval.LaneChanged {
		t.Fatalf("raised lane=%+v", got)
	}
}

func TestDoneGatherMissingTaskAndHandoffAreEvidenceFailures(t *testing.T) {
	t.Parallel()
	f, task, _ := doneBuildFixture(t, "routine", false)
	for _, tc := range []struct {
		task   string
		gate   doneeval.Gate
		reason doneeval.Reason
	}{{"FID-absent", doneeval.TaskGate, doneeval.TaskUnavailable}, {task, doneeval.HandoffGate, doneeval.HandoffMissing}} {
		t.Run(tc.task, func(t *testing.T) {
			t.Parallel()
			got := doneeval.Evaluate(gatherDone(t, f, tc.task, doneeval.Target{}, doneeval.GitHubFacts{}))
			if got.Status != doneeval.Refused || got.Gates[tc.gate].Reason != tc.reason {
				t.Fatalf("missing=%+v", got)
			}
		})
	}
}
