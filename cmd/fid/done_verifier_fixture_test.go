package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/github"
)

// V5 writes all ledger evidence through public commands. The existing strict gh
// stub supplies only GitHub facts and rejects git fetch in its child processes.
type v5Journey struct {
	f          *statusFixture
	task, hold string
	prNumber   int64
	stub       *reconcileStub
	facts      doneeval.GitHubFacts
}

func newV5Journey(t *testing.T, lane string, testsChanged bool) *v5Journey {
	t.Helper()
	f := newJournalStatusFixture(t)
	f.clock = time.Now // Match real-terminal approval children, even across midnight.
	appendPolicy(t, f.workspace, "\n[review]\ntest_adequacy=['**/*_test.go']\n[checks.evidence]\nrun=['true']\ntimeout='1m'\nartifacts=['notes.md']\n[lanes.routine]\npaths=['**/*.md','**/*_test.go']\nchecks=['evidence']\n[lanes.standard]\nchecks=['evidence']\n[lanes.assured]\nchecks=['evidence']\n")
	policyOrigin(t, f.workspace)
	runGitForTest(t, f.workspace, "remote", "set-url", "origin", "https://github.com/mbell-nzl/fid.git")
	j := &v5Journey{f: f, prNumber: 63, stub: newReconcileStub(t, f.workspace)}
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "V5 completion subject", "--lane", lane, "--actor", "human:mark", "--idempotency-key", "create"), "issue.create")
	j.task = created["issue"].(map[string]any)["id"].(string)
	taken, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if taken["kind"] != "new" || taken["task"].(map[string]any)["id"] != j.task {
		t.Fatalf("new work=%v", taken)
	}
	j.hold = taken["hold"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "start", "--hold", j.hold, "--actor", "agent:builder", "--idempotency-key", "start"), "task.start")
	runGitForTest(t, f.workspace, "checkout", "-qb", "task")
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "V5 task change\n")
	if testsChanged {
		writeBackupTestFile(t, filepath.Join(f.workspace, "change_test.go"), "package change\n")
	}
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-qm", "Task change")
	return j
}

func (j *v5Journey) checkAndHandoff(t *testing.T, key string) {
	t.Helper()
	lifecycleResult(t, j.f.run("--json", "check", "--hold", j.hold, "--actor", "agent:builder", "--idempotency-key", key+"-check"), "check.execute")
	h, _ := lifecycleResult(t, j.f.run(j.handoffArgs(key)...), "task.handoff")
	// The head/base come from the command's receipt, never an internal reader.
	now := time.Now().UTC()
	pr := github.PullRequest{Number: j.prNumber, URL: h["pr"].(string), State: "MERGED", Merged: true, MergedAt: &now, MergeCommit: fmt.Sprintf("%040x", j.prNumber), HeadSHA: h["head"].(string), BaseSHA: h["base"].(string), BaseRef: "trunk", BaseRepository: github.Repository{Owner: "mbell-nzl", Name: "fid"}, Body: "## Rollback\nRevert this change.\n"}
	ci := completedCI(pr)
	ci.Runs[0].Conclusion = "success"
	tree := strings.TrimSpace(identityGit(t, j.f.workspace, "rev-parse", pr.HeadSHA+"^{tree}"))
	j.facts = doneeval.GitHubFacts{PR: &pr, HeadTree: tree, CI: &ci}
	j.stub.add(j.facts)
	j.assertState(t, "in_review")
}

func (j *v5Journey) handoffArgs(key string) []string {
	return []string{"--json", "handoff", "--hold", j.hold, "--actor", "agent:builder", "--pr", fmt.Sprintf("https://github.com/mbell-nzl/fid/pull/%d", j.prNumber), "--idempotency-key", key}
}

func (j *v5Journey) review(t *testing.T, stream, key string, input map[string]any) map[string]any {
	t.Helper()
	input["stream"] = stream
	input["reviewer"].(map[string]any)["session_id"] = key
	fields, _ := lifecycleResult(t, recordReviewTask(t, j.f, input, key, "human:courier", j.task), "review.record")
	return fields["round"].(map[string]any)
}

func (j *v5Journey) approve(t *testing.T, key string) {
	t.Helper()
	approveSuccess(t, approvePublicForTest(t, j.f.workspace, j.task, j.task+"\n", true, true, true, "--idempotency-key", key))
}

func (j *v5Journey) rework(t *testing.T, edit bool) {
	t.Helper()
	taken, _ := lifecycleResult(t, j.f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if taken["kind"] != "rework" || taken["task"].(map[string]any)["id"] != j.task {
		t.Fatalf("rework=%v", taken)
	}
	j.hold = taken["hold"].(map[string]any)["id"].(string)
	if edit {
		writeBackupTestFile(t, filepath.Join(j.f.workspace, "notes.md"), "Edited V5 change\n")
		runGitForTest(t, j.f.workspace, "add", ".")
		runGitForTest(t, j.f.workspace, "commit", "-qm", "Fix review finding")
	}
	j.checkAndHandoff(t, "rework")
}

func (j *v5Journey) finish(t *testing.T) map[string]any {
	t.Helper()
	fields, _ := lifecycleResult(t, j.stub.run("--json", "done", j.task, "--actor", "agent:finisher"), "task.done")
	if fields["status"] != "done" || fields["completion"].(map[string]any)["task_id"] != j.task {
		t.Fatalf("completion=%v", fields)
	}
	j.assertState(t, "done")
	shown, _ := lifecycleResult(t, j.f.run("--json", "show", j.task), "issue.show")
	if shown["kept"] != nil || shown["working"] != nil {
		t.Fatalf("completion kept its worktree: %v", shown)
	}
	return fields["completion"].(map[string]any)
}

func (j *v5Journey) assertState(t *testing.T, want string) {
	t.Helper()
	if state, _ := taskState(t, j.f, j.task); state != want {
		t.Fatalf("task state=%s want=%s", state, want)
	}
}

// Merge and CI observations may be recorded on refusal. Task state, revision,
// checks, hand-off, questions and worktree ownership must remain unchanged.
func (j *v5Journey) state(t *testing.T, id string) map[string]any {
	t.Helper()
	fields, _ := lifecycleResult(t, j.f.run("--json", "show", id), "issue.show")
	state := map[string]any{}
	for _, key := range []string{"issue", "working", "kept", "handoff", "questions", "checks"} {
		state[key] = fields[key]
	}
	return state
}

func (j *v5Journey) refuse(t *testing.T, id, op, code string, exit int, run func() commandExecutionResult) commandExecutionResult {
	t.Helper()
	before := j.state(t, id)
	got := run()
	assertJSONFailureCode(t, got.stderr, op, code)
	if v5ExitCode(got.err) != exit || got.stdout != "" {
		t.Fatalf("refusal=%#v want exit=%d", got, exit)
	}
	if after := j.state(t, id); !reflect.DeepEqual(before, after) {
		t.Fatalf("refusal changed task or reservation:\nbefore=%v\nafter=%v", before, after)
	}
	return got
}

func v5ExitCode(err error) int {
	var process interface{ ExitCode() int }
	if errors.As(err, &process) {
		return process.ExitCode()
	}
	return cli.ExitCode(err)
}

func (j *v5Journey) refuseDone(t *testing.T, code, reason string) commandExecutionResult {
	t.Helper()
	got := j.refuse(t, j.task, "task.done", code, 4, func() commandExecutionResult {
		return j.stub.run("--json", "done", j.task, "--actor", "agent:finisher", "--idempotency-key", "finish")
	})
	if !strings.Contains(got.stderr, `"reason":"`+reason+`"`) {
		t.Fatalf("wrong refusal reason: %s", got.stderr)
	}
	return got
}

func (j *v5Journey) keptArtifact(t *testing.T) string {
	t.Helper()
	shown, _ := lifecycleResult(t, j.f.run("--json", "show", j.task), "issue.show")
	checks := shown["checks"].([]any)
	result := checks[0].(map[string]any)["result"].(map[string]any)
	digest := result["artifacts"].([]any)[0].(map[string]any)["digest"].(string)
	return filepath.Join(config.LedgerDir(j.f.workspace), "artifacts", strings.TrimPrefix(digest, "sha256:"))
}

func (j *v5Journey) removeArtifact(t *testing.T) {
	t.Helper()
	if err := os.Remove(j.keptArtifact(t)); err != nil {
		t.Fatal(err)
	}
}
