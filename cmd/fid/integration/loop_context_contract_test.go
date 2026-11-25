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
	nullable := w.create("task", "Task with no brief or criteria", "")
	c.observe("show-no-brief", w.call("show", nullable), 0)
	c.observe("review-list-empty", w.call("review", "list", task), 0)
	w.build("", task, "context")
	reviewer := loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:reviewer"}
	reviewer.wantNext("", true, "review", task)
	round := loopRound(reviewer, task, "CHANGES", "context-changes", []string{"code"})
	c.observe("review-list-changes", w.call("review", "list", task), 0)
	finding := p21Path(t, round, "result", "round", "findings").([]any)[0].(map[string]any)["id"].(string)
	w.call("finding", "answer", finding, "--text", "Fixed with a retry regression", "--actor", w.actor, "--idempotency-key", "response")
	w.build("", task, "fixed-context")
	reviewer.wantNext("", true, "review", task)
	loopRound(reviewer, task, "APPROVE", "context-approve", []string{"code"}, finding)
	c.observe("review-list-resolved", w.call("review", "list", task), 0)
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
