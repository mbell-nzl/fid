package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func runDone(t *testing.T, f *statusFixture, run github.Runner, args ...string) commandExecutionResult {
	t.Helper()
	root := newRootCommandAt(f.workspace)
	old, _, err := root.Find([]string{"done"})
	if err != nil {
		t.Fatal(err)
	}
	root.RemoveCommand(old)
	root.AddCommand(newDoneCommand(run))
	out, errout := &bytes.Buffer{}, &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errout)
	root.SetArgs(args)
	err = root.Execute()
	return commandExecutionResult{stdout: out.String(), stderr: errout.String(), err: err}
}

func doneFacts(t *testing.T, f *statusFixture, task string) doneeval.GitHubFacts {
	t.Helper()
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	h, err := reviewrecord.LatestHandoff(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	// Policy is read from local origin/trunk; done resolves the GitHub identity
	// from its configured URL, independently of that local tracking ref.
	runGitForTest(t, f.workspace, "remote", "set-url", "origin", "https://github.com/mbell-nzl/fid.git")
	pr := github.PullRequest{Number: 63, URL: h.PR, State: "MERGED", Merged: true, MergeCommit: strings.Repeat("a", 40), HeadSHA: h.Head, BaseSHA: h.Base, BaseRef: "trunk", HeadRef: "task", Body: "## Rollback\nRevert the change.", BaseRepository: github.Repository{Owner: "mbell-nzl", Name: "fid"}}
	now := f.clock().UTC()
	pr.MergedAt = &now
	ci := completedCI(pr)
	ci.Runs[0].Conclusion = "success"
	return doneeval.GitHubFacts{PR: &pr, HeadTree: h.Tree, CI: &ci}
}

func doneRunner(t *testing.T, facts doneeval.GitHubFacts, onPR func()) github.Runner {
	t.Helper()
	return func(_ context.Context, args, _ []string) github.Result {
		var response any
		switch {
		case len(args) > 2 && args[0] == "pr" && args[1] == "view":
			if onPR != nil {
				onPR()
			}
			pr := facts.PR
			var mergedAt any
			var merge any
			if pr.Merged {
				mergedAt = pr.MergedAt
				merge = map[string]any{"oid": pr.MergeCommit}
			}
			response = map[string]any{"number": pr.Number, "url": pr.URL, "state": pr.State, "mergedAt": mergedAt, "mergeCommit": merge, "headRefOid": pr.HeadSHA, "headRefName": "task", "baseRefName": pr.BaseRef, "baseRefOid": pr.BaseSHA, "isCrossRepository": false, "body": pr.Body}
		case len(args) > 5 && args[0] == "api" && strings.Contains(args[5], "git/commits/"):
			response = map[string]any{"sha": facts.PR.HeadSHA, "tree": facts.HeadTree}
		case len(args) > 5 && args[0] == "api" && strings.Contains(args[5], "actions/runs?"):
			response = map[string]any{"total_count": len(facts.CI.Runs), "runs": facts.CI.Runs}
		default:
			t.Errorf("unknown gh request: %v", args)
			return github.Result{ExitCode: 1}
		}
		raw, err := json.Marshal(response)
		if err != nil {
			t.Fatal(err)
		}
		return github.Result{Stdout: raw}
	}
}

func TestDoneCompletesReviewedTasksAndReplays(t *testing.T) {
	t.Parallel()
	for _, lane := range []string{"standard", "assured"} {
		t.Run(lane, func(t *testing.T) {
			t.Parallel()
			f, task, hold := doneBuildFixture(t, lane, true)
			lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
			lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
			if lane == "assured" {
				in := reviewInput()
				in["stream"] = "separate"
				in["reviewer"].(map[string]any)["session_id"] = "separate-session"
				lifecycleResult(t, recordReviewTask(t, f, in, "separate", "human:courier", task), "review.record")
				approveSuccess(t, approveForTest(f, task, task+"\n", true))
			}
			facts := doneFacts(t, f, task)
			calls := 0
			runner := doneRunner(t, facts, func() { calls++ })
			args := []string{"--json", "done", task, "--actor", "agent:finisher", "--idempotency-key", "finish"}
			first := runDone(t, f, runner, args...)
			fields, _ := lifecycleResult(t, first, "task.done")
			if fields["status"] != "done" || calls != 1 {
				t.Fatalf("result=%v calls=%d", fields, calls)
			}
			before := startEventCount(t, f.workspace)
			again := runDone(t, f, func(context.Context, []string, []string) github.Result {
				t.Error("retry contacted GitHub")
				return github.Result{}
			}, args...)
			if again.err != nil || again.stdout != first.stdout || startEventCount(t, f.workspace) != before {
				t.Fatalf("retry=%#v", again)
			}
			if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
				t.Fatal(err)
			}
			conn, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = conn.Close() }()
			var state string
			var kept int
			if err := conn.QueryRow(`SELECT status FROM issues WHERE id=?`, task).Scan(&state); err != nil {
				t.Fatal(err)
			}
			if err := conn.QueryRow(`SELECT count(*) FROM holds WHERE task_id=? AND (ended_at IS NULL OR kept=1)`, task).Scan(&kept); err != nil {
				t.Fatal(err)
			}
			if state != "done" || kept != 0 {
				t.Fatalf("state=%s kept=%d", state, kept)
			}
		})
	}
}

func TestDoneWaitingReobservesAndRecordsMerge(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"not-merged", "no-ci", "pending", "failed"} {
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			f, task, hold := doneBuildFixture(t, "routine", false)
			lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
			facts := doneFacts(t, f, task)
			code := "done-ci-waiting"
			switch state {
			case "not-merged":
				facts.PR.Merged = false
				facts.PR.State = "OPEN"
				code = "done-not-merged"
			case "no-ci":
				facts.CI.Runs = []github.WorkflowRun{}
			case "pending":
				facts.CI.Runs[0].Status = "in_progress"
				facts.CI.Runs[0].Conclusion = ""
			case "failed":
				facts.CI.Runs[0].Conclusion = "failure"
				code = "done-ci-not-green"
			}
			calls := 0
			runner := doneRunner(t, facts, func() { calls++ })
			args := []string{"--json", "done", task, "--actor", "agent:finisher", "--idempotency-key", "retry"}
			for range 2 {
				got := runDone(t, f, runner, args...)
				if cli.ExitCode(got.err) != 4 || !strings.Contains(got.stderr, `"code":"`+code+`"`) {
					t.Fatalf("waiting=%#v", got)
				}
			}
			if calls != 2 {
				t.Fatalf("observations=%d", calls)
			}
			conn, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = conn.Close() }()
			var n int
			if err := conn.QueryRow(`SELECT count(*) FROM merge_observations WHERE task_id=?`, task).Scan(&n); err != nil {
				t.Fatal(err)
			}
			want := 1
			if state == "not-merged" {
				want = 0
			}
			if n != want {
				t.Fatalf("merge count=%d", n)
			}
			if err := conn.QueryRow(`SELECT count(*) FROM completions`).Scan(&n); err != nil {
				t.Fatal(err)
			}
			if n != 0 {
				t.Fatal("waiting completed task")
			}
		})
	}
}
