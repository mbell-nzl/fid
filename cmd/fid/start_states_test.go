package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
)

func TestTaskLanePlanning(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Plan the task lane", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "lane-create"), "issue.create")
	issue := created["issue"].(map[string]any)
	id := issue["id"].(string)
	if issue["lane"] != "standard" || issue["status"] != "todo" {
		t.Fatalf("created = %v", issue)
	}
	for _, step := range []struct{ lane, actor, key, code string }{
		{"routine", "agent:builder", "agent-lower", "authority-denied"},
		{"routine", "human:mark", "person-lower", ""},
		{"assured", "agent:builder", "agent-raise", ""},
	} {
		result := f.run("--json", "issue", "update", id, "--lane", step.lane, "--actor", step.actor, "--idempotency-key", step.key)
		if step.code != "" {
			assertJSONFailureCode(t, result.stderr, "issue.update", step.code)
			if cli.ExitCode(result.err) != 5 {
				t.Fatal(result)
			}
			continue
		}
		lifecycleResult(t, result, "issue.update")
	}
	shown, _ := lifecycleResult(t, f.run("--json", "show", id), "issue.show")
	if shown["issue"].(map[string]any)["lane"] != "assured" {
		t.Fatalf("shown = %v", shown)
	}
	human := f.run("show", id)
	if human.err != nil || !strings.Contains(human.stdout, "Lane        Assured") {
		t.Fatalf("show text = %#v", human)
	}
}

func TestStartUsesCurrentHoldAndIsHarmlessToRepeat(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	now := time.Now().Add(time.Minute)
	f.clock = func() time.Time { return now }
	f.seed([]readyCommandRow{statusEpic("FID-e3a1c0", "Planned epic", "todo", nil), statusTask("FID-9f13aa", "Planned task", "todo", stringPointer("FID-e3a1c0"), nil, nil)}, nil)
	lifecycleResult(t, f.run("--json", "issue", "update", "FID-9f13aa", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "lane"), "issue.update")
	// The planned task is held through the same scheduler used by the loop.
	taken, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder"), "next.claim")
	id := taken["task"].(map[string]any)["id"].(string)
	if id != "FID-9f13aa" {
		t.Fatalf("taken task = %s", id)
	}
	hold := taken["hold"].(map[string]any)["id"].(string)
	args := []string{"--json", "start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start-1"}
	before := startEventCount(t, f.workspace)
	first := f.run(args...)
	started, event := lifecycleResult(t, first, "task.start")
	if started["status"] != "in_progress" || started["lane"] != "standard" || event == nil {
		t.Fatalf("start = %v, %v", started, event)
	}
	if again := f.run(args...); again.err != nil || again.stdout != first.stdout {
		t.Fatalf("same-key start = %#v", again)
	}
	args[len(args)-1] = "start-2"
	_, event = lifecycleResult(t, f.run(args...), "task.start")
	if event != nil || startEventCount(t, f.workspace) != before+1 {
		t.Fatal("repeat start appended an Event")
	}
	epic, _ := lifecycleResult(t, f.run("--json", "show", "FID-e3a1c0"), "issue.show")
	if epic["issue"].(map[string]any)["status"] != "in_progress" {
		t.Fatalf("epic = %v", epic)
	}
	refused := f.run("--json", "issue", "update", id, "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "started-lower")
	assertJSONFailureCode(t, refused.stderr, "issue.update", "transition-not-permitted")
}

func TestFixedTaskTransitionsRefuseManualMovement(t *testing.T) {
	t.Parallel()
	for _, step := range []struct{ id, state, to string }{
		{"FID-9f13aa", "todo", "in_progress"},
		{"FID-71a0c2", "in_progress", "in_review"},
		{"FID-71a0c2", "in_progress", "todo"},
		{"FID-e3a1c0", "todo", "in_progress"},
	} {
		t.Run(step.id+step.to, func(t *testing.T) {
			f := newLifecycleFixture(t)
			nextWork(f)
			before := startEventCount(t, f.workspace)
			result := f.run("--json", "issue", "transition", step.id, "--to", step.to, "--actor", "human:mark", "--idempotency-key", "manual")
			if cli.ExitCode(result.err) != 4 {
				t.Fatalf("manual transition = %#v", result)
			}
			assertJSONFailureCode(t, result.stderr, "issue.transition", "transition-not-permitted")
			if before != startEventCount(t, f.workspace) {
				t.Fatal("refused transition wrote an Event")
			}
		})
	}
}

func startEventCount(t *testing.T, root string) int {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(config.LedgerDir(root), "journal", "*.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		n += bytes.Count(b, []byte("\n"))
	}
	return n
}
