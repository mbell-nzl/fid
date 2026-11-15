package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Each subject has its own review reservation and linked PR in one ledger.
func reconcileSubjects(t *testing.T) (*statusFixture, []string, []doneeval.GitHubFacts) {
	t.Helper()
	f, task, hold := doneBuildFixture(t, "routine", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	ids := []string{task}
	facts := []doneeval.GitHubFacts{doneFacts(t, f, task)}
	for i := 1; i < 5; i++ {
		other := *f
		other.workspace = gittest.AddWorktree(t, f.workspace)
		lane := "routine"
		if i == 2 {
			lane = "standard"
		} // Missing review is a refused gate.
		created, _ := lifecycleResult(t, other.run("--json", "issue", "create", "--kind", "task", "--title", fmt.Sprintf("Reconcile %d", i), "--lane", lane, "--actor", "human:mark", "--idempotency-key", fmt.Sprintf("create-%d", i)), "issue.create")
		id := created["issue"].(map[string]any)["id"].(string)
		took, _ := lifecycleResult(t, other.run("--json", "takeover", id, "--actor", "human:builder", "--model", "codex/gpt-6.1-sol", "--reason", "Build", "--idempotency-key", "take-"+id), "task.takeover")
		h := took["hold_id"].(string)
		lifecycleResult(t, other.run("--json", "start", "--hold", h, "--actor", "human:builder", "--idempotency-key", "start-"+id), "task.start")
		lifecycleResult(t, other.run("--json", "check", "--hold", h, "--actor", "human:builder", "--idempotency-key", "check-"+id), "check.execute")
		lifecycleResult(t, other.run("--json", "handoff", "--hold", h, "--actor", "human:builder", "--pr", fmt.Sprintf("https://github.com/mbell-nzl/fid/pull/%d", 63+i), "--idempotency-key", "handoff-"+id), "task.handoff")
		remote := doneFacts(t, &other, id)
		remote.PR.Number = int64(63 + i)
		remote.PR.MergeCommit = fmt.Sprintf("%040x", i+10)
		remote.CI = ptrV8CI(completedCI(*remote.PR))
		remote.CI.Runs[0].Conclusion = "success"
		facts, ids = append(facts, remote), append(ids, id)
	}
	facts[1].CI.Runs = []github.WorkflowRun{} // Waiting for CI.
	facts[3].PR.Merged, facts[3].PR.State = false, "OPEN"
	return f, ids, facts
}

func TestReconcileMixedOutcomesAndRetry(t *testing.T) {
	t.Parallel()
	f, ids, facts := reconcileSubjects(t)
	stub := newReconcileStub(t, f.workspace)
	for _, remote := range facts {
		stub.add(remote)
	}
	stub.requests[reconcilePRKey(facts[4])] = map[string]any{"error": "offline for this PR"}
	first := stub.run("--json", "done", "--actor", "agent:finisher")
	assertReconciled(t, first, ids, []string{"completed", "waiting", "refused", "waiting", "github-unavailable"})
	if !strings.Contains(first.stdout, `"code":"done-review-required"`) || !strings.Contains(first.stdout, `"gate":"primary-review"`) || !strings.Contains(first.stdout, `"code":"done-not-merged"`) {
		t.Fatalf("missing reasons: %s", first.stdout)
	}
	log, err := os.ReadFile(stub.log)
	if err != nil {
		t.Fatal(err)
	}
	observed := []string{}
	for _, line := range strings.Split(strings.TrimSpace(string(log)), "\n") {
		if strings.HasPrefix(line, "pr view ") {
			observed = append(observed, strings.Fields(line)[2])
		}
	}
	ordered := append([]string{}, ids...)
	sort.Strings(ordered)
	want := []string{}
	for _, id := range ordered {
		for i, subject := range ids {
			if id == subject {
				want = append(want, fmt.Sprint(63+i))
			}
		}
	}
	if !reflect.DeepEqual(observed, want) {
		t.Fatalf("observation order=%v want=%v", observed, want)
	}
	before := startEventCount(t, f.workspace)
	again := stub.run("--json", "done", "--actor", "agent:finisher")
	assertReconciled(t, again, ids[1:], []string{"waiting", "refused", "waiting", "github-unavailable"})
	if startEventCount(t, f.workspace) != before {
		t.Fatal("unchanged observations appended Events")
	}
	facts[1].CI = ptrV8CI(completedCI(*facts[1].PR))
	facts[1].CI.Runs[0].Conclusion = "success"
	stub.add(facts[1])
	last := stub.run("done", "--actor", "human:mark")
	if last.err != nil || last.stderr != "" || !strings.Contains(last.stdout, "Completed "+ids[1]) || !strings.Contains(last.stdout, "Refused "+ids[2]) || !strings.Contains(last.stdout, "Waiting "+ids[3]) || !strings.Contains(last.stdout, "GitHub unavailable "+ids[4]) {
		t.Fatalf("text=%#v", last)
	}
	// One CI Event and one completion, with no repeated merge/other CI receipt.
	if got := startEventCount(t, f.workspace) - before; got != 2 {
		t.Fatalf("new Events=%d want=2", got)
	}
}

func TestReconcileMovesDuringObservationForOneTaskOnly(t *testing.T) {
	t.Parallel()
	f, ids, facts := reconcileSubjects(t)
	calls := 0
	got := runDone(t, f, func(ctx context.Context, args, env []string) github.Result {
		i := 0
		if args[0] == "pr" {
			calls++
			for n, remote := range facts {
				if args[2] == fmt.Sprint(remote.PR.Number) {
					i = n
					break
				}
			}
			if i == 1 {
				lifecycleResult(t, executeCommandForTest(gittest.AddWorktree(t, f.workspace), "--json", "takeover", ids[i], "--actor", "human:mark", "--reason", "Changed during observation"), "task.takeover")
			}
		} else {
			for n, remote := range facts {
				if strings.Contains(args[5], remote.PR.MergeCommit) {
					i = n
					break
				}
			}
		}
		return doneRunner(t, facts[i], nil)(ctx, args, env)
	}, "--json", "done", "--actor", "agent:finisher")
	// All subjects have the same head/tree, so the shared commit-tree response is valid.
	assertReconciled(t, got, ids, []string{"completed", "refused", "refused", "waiting", "completed"})
	if calls != 5 || !strings.Contains(got.stdout, `"code":"done-observation-stale"`) {
		t.Fatalf("calls=%d result=%#v", calls, got)
	}
}

func TestReconcileEmptyAndCommandFailures(t *testing.T) {
	t.Parallel()
	f, _, _ := doneBuildFixture(t, "routine", false) // A started task is excluded.
	runner := func(context.Context, []string, []string) github.Result {
		t.Error("empty reconciliation contacted GitHub")
		return github.Result{}
	}
	got := runDone(t, f, runner, "--json", "done", "--actor", "agent:finisher")
	assertReconciled(t, got, nil, nil)
	before := startEventCount(t, f.workspace)
	for _, tc := range []struct {
		args []string
		code string
	}{
		{[]string{"--json", "done"}, "invalid-argument"},
		{[]string{"--json", "done", "--actor", "service:bot"}, "authority-denied"},
		{[]string{"--json", "done", "--actor", "agent:finisher", "--idempotency-key", "batch"}, "invalid-argument"},
		{[]string{"--json", "done", "one", "two", "--actor", "agent:finisher"}, "invalid-argument"},
	} {
		assertJSONFailureCode(t, runDone(t, f, runner, tc.args...).stderr, "task.done", tc.code)
	}
	if startEventCount(t, f.workspace) != before {
		t.Fatal("empty/invalid reconciliation wrote Events")
	}
}

func assertReconciled(t *testing.T, got commandExecutionResult, ids, states []string) {
	t.Helper()
	fields, _ := lifecycleResult(t, got, "task.done")
	outcomes, ok := fields["outcomes"].([]any)
	if !ok || len(outcomes) != len(ids) {
		t.Fatalf("outcomes=%s want %d", got.stdout, len(ids))
	}
	byID := map[string]string{}
	previous := ""
	for _, raw := range outcomes {
		row := raw.(map[string]any)
		id := row["task_id"].(string)
		if id <= previous {
			t.Fatalf("unstable order: %s", got.stdout)
		}
		previous = id
		byID[id] = row["status"].(string)
	}
	for i, id := range ids {
		if byID[id] != states[i] {
			t.Fatalf("%s=%s want=%s: %s", id, byID[id], states[i], got.stdout)
		}
	}
}

// The existing strict gh stub sees only exact expected requests; git rejects fetch.
type reconcileStub struct {
	t                       *testing.T
	root, bin, fixture, log string
	requests                map[string]any
}

func newReconcileStub(t *testing.T, root string) *reconcileStub {
	t.Helper()
	dir := t.TempDir()
	script, err := os.ReadFile("testdata/gh_stub.py")
	if err != nil {
		t.Fatal(err)
	}
	writeDoneExecutable(t, filepath.Join(dir, "gh"), script)
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	writeDoneExecutable(t, filepath.Join(dir, "git"), []byte("#!/bin/sh\nfor arg do\n if [ \"$arg\" = fetch ]; then exit 97; fi\ndone\nexec '"+git+"' \"$@\"\n"))
	return &reconcileStub{t: t, root: root, bin: dir, fixture: filepath.Join(dir, "responses.json"), log: filepath.Join(dir, "requests.log"), requests: map[string]any{}}
}
func reconcilePRKey(f doneeval.GitHubFacts) string {
	return fmt.Sprintf("pr view %d --repo mbell-nzl/fid --json number,url,state,mergedAt,mergeCommit,headRefOid,headRefName,baseRefName,baseRefOid,isCrossRepository,body", f.PR.Number)
}
func (s *reconcileStub) add(f doneeval.GitHubFacts) {
	pr := f.PR
	var mergedAt, merge any
	if pr.Merged {
		mergedAt = pr.MergedAt
		merge = map[string]any{"oid": pr.MergeCommit}
	}
	s.requests[reconcilePRKey(f)] = map[string]any{"body": map[string]any{"number": pr.Number, "url": pr.URL, "state": pr.State, "mergedAt": mergedAt, "mergeCommit": merge, "headRefOid": pr.HeadSHA, "headRefName": "task", "baseRefName": pr.BaseRef, "baseRefOid": pr.BaseSHA, "isCrossRepository": false, "body": pr.Body}}
	s.requests["api --hostname github.com --method GET repos/mbell-nzl/fid/git/commits/"+pr.HeadSHA+" --jq {sha, tree: .tree.sha}"] = map[string]any{"body": map[string]any{"sha": pr.HeadSHA, "tree": f.HeadTree}}
	s.requests["api --hostname github.com --method GET repos/mbell-nzl/fid/actions/runs?head_sha="+pr.MergeCommit+"&per_page=100&page=1 --jq {total_count, runs: [.workflow_runs[] | {id, run_attempt, workflow_id, name, path, event, head_sha, status, conclusion, html_url}]}"] = map[string]any{"body": map[string]any{"total_count": len(f.CI.Runs), "runs": f.CI.Runs}}
}
func (s *reconcileStub) run(args ...string) commandExecutionResult {
	s.t.Helper()
	raw, err := json.Marshal(map[string]any{"requests": s.requests})
	if err != nil {
		s.t.Fatal(err)
	}
	if err := os.WriteFile(s.fixture, raw, 0o600); err != nil {
		s.t.Fatal(err)
	}
	binary, err := os.Executable()
	if err != nil {
		s.t.Fatal(err)
	}
	cmd := exec.CommandContext(s.t.Context(), binary, args...)
	cmd.Dir = s.root
	cmd.Env = append(os.Environ(), runAsFIDVariable+"=1", "PATH="+s.bin+string(os.PathListSeparator)+os.Getenv("PATH"), "FID_TEST_GH_RESPONSES="+s.fixture, "FID_TEST_GH_LOG="+s.log)
	if testing.CoverMode() != "" {
		cmd.Env = append(cmd.Env, "GOCOVERDIR="+s.t.TempDir())
	}
	var out, errout bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errout
	err = cmd.Run()
	return commandExecutionResult{stdout: out.String(), stderr: errout.String(), err: err}
}

func TestReconcileCommandFailureIsNonzero(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "routine", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	facts := doneFacts(t, f, task)
	var got commandExecutionResult
	if err := journal.WithLock(f.workspace, journal.LockOwner{Command: "held for reconciliation test"}, func(*journal.WriteLock) error {
		got = runDone(t, f, doneRunner(t, facts, nil), "--json", "done", "--actor", "agent:finisher")
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	assertJSONFailureCode(t, got.stderr, "task.done", "lock-unavailable")
	if cli.ExitCode(got.err) == 0 || got.stdout != "" {
		t.Fatalf("lock failure=%#v", got)
	}
	noLedger := &statusFixture{t: t, workspace: gittest.NewWorkspace(t)}
	got = runDone(t, noLedger, nil, "--json", "done", "--actor", "agent:finisher")
	assertJSONFailureCode(t, got.stderr, "task.done", "project-not-found")
	if cli.ExitCode(got.err) == 0 || got.stdout != "" {
		t.Fatalf("no ledger=%#v", got)
	}
}
