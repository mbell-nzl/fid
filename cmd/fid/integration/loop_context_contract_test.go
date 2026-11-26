//go:build integration

package integration

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestLoopContextContract(t *testing.T) {
	t.Parallel()
	c := &loopContract{t: t, decl: readLoopDeclaration(t)}
	root, binary := handoffJourneyRepo(t, true)
	w := loopWorker{t, binary, root, "agent:builder"}
	task := w.create("task", "Build from the approved brief", "", "--spec-md", "docs/briefs/2.5.md")
	w.call("issue", "ac", "add", task, "--ordinal", "1", "--text", "The first loop runs", "--actor", "human:planner", "--idempotency-key", "criterion")
	c.observe("show-brief", w.call("show", task), 0)
	c.observe("review-list-empty", w.call("review", "list", task), 0)
	w.build("", task, "context")
	c.observe("show-handed-off", w.call("show", task), 0)
	c.observe("status-review-due", w.call("status"), 0)
	reviewer := loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:reviewer"}
	reviewer.wantNext("", true, "review", task)
	round := loopRound(reviewer, task, "CHANGES", "context-changes", []string{"code"})
	c.observe("review-list-changes", w.call("review", "list", task), 0)
	c.observe("status-changes", w.call("status"), 0)
	finding := p21Path(t, round, "result", "round", "findings").([]any)[0].(map[string]any)["id"].(string)
	w.call("finding", "answer", finding, "--text", "Fixed with a retry regression", "--actor", w.actor)
	taken := w.wantNext("", true, "rework", task)
	hold := p21PathString(t, taken, "result", "hold", "id")
	loopCommit(t, root, "Fix context finding")
	w.call("check", "--hold", hold, "--actor", w.actor, "--idempotency-key", "fixed-context-check")
	w.call("handoff", "--hold", hold, "--actor", w.actor, "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "fixed-context-handoff")
	reviewer.wantNext("", true, "review", task)
	loopRound(reviewer, task, "APPROVE", "context-approve", []string{"code"}, finding)
	c.observe("review-list-resolved", w.call("review", "list", task), 0)
	c.observe("status-approved", w.call("status"), 0)
	nullable := w.create("task", "Task with no brief or criteria", "")
	c.observe("show-no-brief", w.call("show", nullable), 0)
}

func TestLoopReadinessContract(t *testing.T) {
	t.Parallel()
	for _, scenario := range []string{"routine", "assured"} {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			c := &loopContract{t: t, decl: readLoopDeclaration(t)}
			j := newFirstLoop(t, scenario)
			w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
			if scenario == "routine" {
				firstLoopRoutine(t, j)
			} else {
				firstLoopReplaceStub(t, j, "agent", `if work["kind"] == "review":`, `if work["kind"] == "review":
    if work["review"]["stream"] == "separate" and not (state / "separate-interrupted").exists():
        (state / "separate-interrupted").touch()
        sys.exit(42)`)
				j.run(t, 1)
				c.observe("status-primary-only", w.call("status"), 0)
				for _, call := range j.calls(t) {
					if call["kind"] == "ready" {
						t.Fatal("PR became ready before its required separate review")
					}
				}
			}
			j.run(t, 0)
			c.observe("status-"+scenario, w.call("status"), 0)
		})
	}
}

func TestLoopBackupContract(t *testing.T) {
	t.Parallel()
	c := &loopContract{t: t, decl: readLoopDeclaration(t)}
	j := newBackupJourney(t)
	for _, outcome := range []string{"pushed", "unchanged"} {
		run := backupJourneyFID(t, j.env, j.workspace, "--json", "backup")
		run.expectExit(t, 0)
		c.observe("backup-"+outcome, p23Envelope(t, run.stdout.String()), 0)
	}
}
