package main

import (
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
)

func TestDoneRequiresAReviewTask(t *testing.T) {
	t.Parallel()
	f, task, _ := doneBuildFixture(t, "standard", false)
	before := startEventCount(t, f.workspace)
	got := f.run("--json", "done", task, "--actor", "agent:finisher")
	if cli.ExitCode(got.err) != 4 || !strings.Contains(got.stderr, `"code":"done-task-not-in-review"`) || got.stdout != "" {
		t.Fatalf("done=%#v", got)
	}
	if startEventCount(t, f.workspace) != before {
		t.Fatal("refusal changed the journal")
	}
}

func TestIssueTransitionRequiresDoneForTasks(t *testing.T) {
	t.Parallel()
	for _, actor := range []string{"agent:finisher", "human:mark", "service:robot"} {
		t.Run(actor, func(t *testing.T) {
			t.Parallel()
			f, task, hold := doneBuildFixture(t, "standard", false)
			lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
			before := startEventCount(t, f.workspace)
			got := f.run("--json", "issue", "transition", task, "--to", "done", "--actor", actor, "--idempotency-key", "bypass")
			if cli.ExitCode(got.err) != 4 || !strings.Contains(got.stderr, "fid done") || got.stdout != "" {
				t.Fatalf("transition=%#v", got)
			}
			if startEventCount(t, f.workspace) != before {
				t.Fatal("transition bypassed completion")
			}
		})
	}
}

func TestDoneRefusesOtherKindsServicesAndMissingIDs(t *testing.T) {
	t.Parallel()
	f, task, _ := doneBuildFixture(t, "standard", false)
	epic, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "epic", "--scheduling", "prioritised", "--title", "Epic", "--actor", "human:mark", "--idempotency-key", "epic"), "issue.create")
	id := epic["issue"].(map[string]any)["id"].(string)
	for _, tc := range []struct {
		task, actor, code string
		exit              int
	}{{task, "service:bot", "authority-denied", 5}, {id, "human:mark", "invalid-argument", 2}, {"FID-absent", "agent:finisher", "issue-not-found", 3}} {
		t.Run(tc.code, func(t *testing.T) {
			t.Parallel()
			got := f.run("--json", "done", tc.task, "--actor", tc.actor)
			assertJSONFailureCode(t, got.stderr, "task.done", tc.code)
			if cli.ExitCode(got.err) != tc.exit {
				t.Fatalf("done=%#v", got)
			}
		})
	}
}
