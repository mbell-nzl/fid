//go:build integration

package integration

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

// loopWorker drives only public commands, in a repository owned by the test.
type loopWorker struct {
	t                   *testing.T
	binary, root, actor string
}

func (w loopWorker) call(args ...string) map[string]any {
	w.t.Helper()
	return p21Run(w.t, w.binary, w.root, append([]string{"--json"}, args...)...)
}

func (w loopWorker) create(kind, title, epic string, flags ...string) string {
	w.t.Helper()
	args := []string{"issue", "create", "--kind", kind, "--title", title, "--actor", "human:planner", "--idempotency-key", title}
	if kind == "epic" {
		args = append(args, "--scheduling", "prioritised")
	} else {
		args = append(args, "--lane", "standard")
	}
	if epic != "" {
		args = append(args, "--epic", epic)
	}
	return p21PathString(w.t, w.call(append(args, flags...)...), "result", "issue", "id")
}

func (w loopWorker) next(scope string, claim bool) map[string]any {
	w.t.Helper()
	args := []string{"next"}
	if scope != "" {
		args = append(args, "--epic", scope)
	}
	if claim {
		args = append(args, "--claim", "--actor", w.actor, "--model", "codex/gpt-6.1-sol")
	}
	return w.call(args...)
}

func (w loopWorker) wantNext(scope string, claim bool, kind, task string) map[string]any {
	w.t.Helper()
	got := w.next(scope, claim)
	if p21PathString(w.t, got, "result", "kind") != kind {
		w.t.Fatalf("next = %v, want %s %s", got, kind, task)
	}
	if task == "" {
		if p21Path(w.t, got, "result", "task") != nil || p21Path(w.t, got, "result", "hold") != nil {
			w.t.Fatalf("none carries work: %v", got)
		}
	} else if p21PathString(w.t, got, "result", "task", "id") != task {
		w.t.Fatalf("next = %v, want task %s", got, task)
	}
	return got
}

func loopQueueReason(t *testing.T, got map[string]any, task, reason string) {
	t.Helper()
	for _, item := range p21Path(t, got, "result", "queue").([]any) {
		entry := item.(map[string]any)
		if entry["id"] == task {
			if entry["reason"] != reason {
				t.Fatalf("task %s reason = %v, want %s", task, entry["reason"], reason)
			}
			return
		}
	}
	t.Fatalf("task %s absent from queue: %v", task, got)
}

func (w loopWorker) build(scope, task, key string) map[string]any {
	w.t.Helper()
	taken := w.next(scope, true)
	kind := p21PathString(w.t, taken, "result", "kind")
	if (kind != "new" && kind != "resume") || p21PathString(w.t, taken, "result", "task", "id") != task {
		w.t.Fatalf("build selection = %v", taken)
	}
	hold := p21PathString(w.t, taken, "result", "hold", "id")
	w.call("start", "--hold", hold, "--actor", w.actor, "--idempotency-key", key+"-start")
	p22Write(w.t, w.root, key+".txt", key+"\n")
	closureRun(w.t, w.root, "git", "add", ".")
	closureRun(w.t, w.root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", key)
	w.call("check", "--hold", hold, "--actor", w.actor, "--idempotency-key", key+"-check")
	return w.call("handoff", "--hold", hold, "--actor", w.actor, "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", key+"-handoff")
}

func TestLoopEpicScope(t *testing.T) {
	t.Parallel()
	c := &loopContract{t: t, decl: readLoopDeclaration(t)}
	root, binary := handoffJourneyRepo(t, true)
	w := loopWorker{t, binary, root, "agent:side"}
	main, side, empty := w.create("epic", "Main", ""), w.create("epic", "Side", ""), w.create("epic", "Empty", "")
	w.call("issue", "update", main, "--priority", "1", "--actor", "human:planner", "--idempotency-key", "main-priority")
	mainTask := w.create("task", "Main task", main)
	sideTask := w.create("task", "Side task", side)
	orphan := w.create("task", "No epic", "")
	w.wantNext("", false, "new", mainTask)
	got := w.wantNext(side, false, "new", sideTask)
	c.observe("next-epic", got, 0)
	loopQueueReason(t, got, mainTask, "outside-scope")
	loopQueueReason(t, got, orphan, "outside-scope")
	c.observe("next-epic-none", w.wantNext(empty, true, "none", ""), 0)
	other := loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:main"}
	c.observe("next-claim-epic", w.wantNext(side, true, "new", sideTask), 0)
	w.wantNext(main, true, "none", "")
	other.wantNext("", true, "new", mainTask)
	w.build(side, sideTask, "side")
	other.build("", mainTask, "main")
	w.wantNext(side, true, "none", "")
	// Scope includes reviews: an empty epic never takes either due review.
	reviewer := loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:reviewer"}
	reviewer.wantNext(empty, true, "none", "")
	review := reviewer.wantNext(side, true, "review", sideTask)
	c.observe("next-epic-review", review, 0)
	loopQueueReason(t, review, mainTask, "outside-scope")
	// Changing scope never mixes another task into an occupied worktree.
	reviewer.wantNext(main, true, "none", "")
	w.wantNext(main, true, "none", "")
	if out, stderr, err := closureRunCommand(root, binary, "next", "--epic", empty, "--explain"); err != nil || stderr != "" || !strings.Contains(out, "outside scope") {
		t.Fatalf("scope explanation: %v %q %q", err, out, stderr)
	}
}

func TestLoopScopeKeepsOrderingAndBlockers(t *testing.T) {
	t.Parallel()
	for _, scheduling := range []string{"prioritised", "sequential"} {
		t.Run(scheduling, func(t *testing.T) {
			t.Parallel()
			root, binary := handoffJourneyRepo(t, true)
			w := loopWorker{t, binary, root, "agent:builder"}
			e := w.call("issue", "create", "--kind", "epic", "--title", "Ordered epic", "--scheduling", scheduling, "--actor", "human:planner", "--idempotency-key", "epic")
			scope := p21PathString(t, e, "result", "issue", "id")
			order := "--priority"
			if scheduling == "sequential" {
				order = "--sequence"
			}
			first := w.create("task", "Ordered first", scope, order, "1")
			second := w.create("task", "Ordered second", scope, order, "2")
			blocker := w.create("task", "External blocker", "")
			w.call("link", "add", blocker, "blocks", first, "--from-type", "issue", "--to-type", "issue", "--actor", "human:planner", "--idempotency-key", "block")
			kind, served := "new", second
			if scheduling == "sequential" {
				kind, served = "none", ""
			}
			got := w.wantNext(scope, false, kind, served)
			loopQueueReason(t, got, blocker, "outside-scope")
			w.call("link", "remove", blocker, "blocks", first, "--from-type", "issue", "--to-type", "issue", "--actor", "human:planner", "--idempotency-key", "unblock")
			w.wantNext(scope, true, "new", first)
			other := loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:second"}
			if scheduling == "sequential" {
				other.wantNext(scope, true, "none", "")
			} else {
				other.wantNext(scope, true, "new", second)
			}
		})
	}
}

func TestLoopEpicScopeRefusals(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepo(t, true)
	w := loopWorker{t, binary, root, "agent:builder"}
	task := w.create("task", "Not an epic", "")
	for _, test := range []struct {
		scope, code string
		exit        int
	}{{"FID-missing", "issue-not-found", 3}, {task, "unsupported-kind", 2}, {"", "invalid-argument", 2}} {
		for _, claim := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/claim=%t", test.code, claim), func(t *testing.T) {
				args := []string{"--json", "next", "--epic", test.scope}
				if claim {
					args = append(args, "--claim", "--actor", w.actor, "--model", "codex/gpt-6.1-sol")
				}
				before := p21JournalEventCount(t, root)
				out, stderr, err := closureRunCommand(root, binary, args...)
				if out != "" || processExitCode(t, err) != test.exit {
					t.Fatalf("scope refusal: %v %s %s", err, out, stderr)
				}
				got := p23Envelope(t, stderr)
				operation := "next.get"
				if claim && test.scope != "" {
					operation = "next.claim"
				}
				if got["operation"] != operation || p21PathString(t, got, "error", "code") != test.code {
					t.Fatalf("scope refusal = %v", got)
				}
				if p21JournalEventCount(t, root) != before {
					t.Fatal("invalid scope wrote an event")
				}
			})
		}
	}
}

func TestLoopScopedParallelClaims(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepo(t, true)
	w := loopWorker{t, binary, root, "agent:first"}
	scope := w.create("epic", "Shared epic", "")
	first, second := w.create("task", "First", scope), w.create("task", "Second", scope)
	roots := []string{root, gittest.AddWorktree(t, root)}
	type reply struct {
		out, stderr string
		err         error
	}
	replies := make([]reply, 2)
	var wg sync.WaitGroup
	for i := range roots {
		wg.Go(func() {
			replies[i].out, replies[i].stderr, replies[i].err = closureRunCommand(roots[i], binary, "--json", "next", "--epic", scope, "--claim", "--actor", fmt.Sprintf("agent:worker%d", i), "--model", "codex/gpt-6.1-sol")
		})
	}
	wg.Wait()
	seen := map[string]bool{}
	for i, r := range replies {
		if r.err != nil || r.stderr != "" {
			t.Fatalf("worker %d: %v %s", i, r.err, r.stderr)
		}
		got := p23Envelope(t, r.out)
		id := p21PathString(t, got, "result", "task", "id")
		if seen[id] || (id != first && id != second) {
			t.Fatalf("duplicate or out-of-scope task %s", id)
		}
		seen[id] = true
		worker := loopWorker{t, binary, roots[i], fmt.Sprintf("agent:worker%d", i)}
		repeated := worker.wantNext(scope, true, "resume", id)
		if p21PathString(t, repeated, "result", "hold", "id") != p21PathString(t, got, "result", "hold", "id") {
			t.Fatal("repeat changed hold")
		}
	}
}
