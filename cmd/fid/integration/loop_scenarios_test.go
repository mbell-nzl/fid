//go:build integration

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func loopCommit(t *testing.T, root, name string) {
	t.Helper()
	p22Write(t, root, "change.txt", name+"\n")
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", name)
}

func loopRound(w loopWorker, task, verdict, key string, profiles []string, resolved ...string) map[string]any {
	w.t.Helper()
	in := map[string]any{"reviewer": map[string]any{"actor": w.actor, "model": "claude/opus-5.5", "session_id": key},
		"stream": "primary", "profiles": profiles, "verdict": verdict, "purpose": "Review the loop scenario", "verification": "Public commands exercised", "findings": []any{}}
	if verdict == "CHANGES" {
		in["findings"] = []any{map[string]any{"profile": "code", "severity": "major", "summary": "Missing retry case", "location": nil, "evidence": "Retry case reproduced", "fix": "Handle the retry"}}
	}
	if len(resolved) > 0 {
		dispositions := []any{}
		for _, id := range resolved {
			dispositions = append(dispositions, map[string]any{"finding_id": id, "disposition": "resolved", "reason": "Regression verifies the fix"})
		}
		in["dispositions"] = dispositions
	}
	input := p21Input(w.t, w.t.TempDir(), "round", in)
	return w.call("review", "record", task, "--input", input, "--actor", "human:courier", "--idempotency-key", key)
}

func TestLoopContractScenarios(t *testing.T) {
	t.Parallel()
	c := &loopContract{t: t, decl: readLoopDeclaration(t)}
	root, binary := handoffJourneyRepo(t, true)
	w := loopWorker{t, binary, root, "agent:builder"}
	task := w.create("task", "Loop task", "")
	call := func(name string, args ...string) map[string]any {
		t.Helper()
		got := w.call(args...)
		c.observe(name, got, 0)
		return got
	}
	call("next-new", "next")
	taken := w.wantNext("", true, "new", task)
	c.observe("next-claim-new", taken, 0)
	hold := p21PathString(t, taken, "result", "hold", "id")
	call("start", "start", "--hold", hold, "--actor", w.actor, "--idempotency-key", "start")
	call("heartbeat", "heartbeat", "--hold", hold, "--actor", w.actor)
	call("park-interrupted", "park", "--hold", hold, "--actor", w.actor, "--interrupted", "--reason", "Usage limit", "--idempotency-key", "interrupt")
	other := loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:other"}
	other.wantNext("", true, "none", "")
	call("next-interrupted", "next")
	taken = w.wantNext("", true, "resume", task)
	c.observe("next-claim-resume", taken, 0)
	newHold := p21PathString(t, taken, "result", "hold", "id")
	if newHold == hold {
		t.Fatal("interruption reused ended hold")
	}
	hold = newHold
	call("park-waiting", "park", "--hold", hold, "--actor", w.actor, "--waiting", "Include retries?", "--idempotency-key", "waiting")
	c.observe("next-waiting", w.wantNext("", true, "none", ""), 0)
	call("answer", "answer", task, "Include them.", "--actor", "human:mark", "--idempotency-key", "answer")
	taken = w.wantNext("", true, "resume", task)
	c.observe("next-answer", taken, 0)
	if p21PathString(t, taken, "result", "answer", "text") != "Include them." {
		t.Fatalf("lost answer: %v", taken)
	}
	hold = p21PathString(t, taken, "result", "hold", "id")
	loopCommit(t, root, "First implementation")
	call("check", "check", "--hold", hold, "--actor", w.actor, "--idempotency-key", "check")
	h := call("handoff", "handoff", "--hold", hold, "--actor", w.actor, "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "handoff")
	c.observe("next-review-wait", w.wantNext("", true, "none", ""), 0)
	// Two eligible reviewers race; only one receives this review session.
	reviewer := loopClaimReview(t, c, w, other, task)
	changes := loopRound(reviewer, task, "CHANGES", "changes", []string{"code"})
	c.observe("review-changes", changes, 0)
	finding := p21Path(t, changes, "result", "round", "findings").([]any)[0].(map[string]any)["id"].(string)
	other.wantNext("", true, "none", "")
	taken = w.wantNext("", true, "rework", task)
	c.observe("next-rework", taken, 0)
	hold = p21PathString(t, taken, "result", "hold", "id")
	loopCommit(t, root, "Fix requested change")
	call("check-rework", "check", "--hold", hold, "--actor", w.actor, "--idempotency-key", "check-fix")
	h2 := call("handoff-rework", "handoff", "--hold", hold, "--actor", w.actor, "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "handoff-fix")
	if p21PathString(t, h, "result", "handoff_id") == p21PathString(t, h2, "result", "handoff_id") {
		t.Fatal("rework reused hand-off")
	}
	reviewer.wantNext("", true, "review", task)
	c.observe("review-approve", loopRound(reviewer, task, "APPROVE", "approve", []string{"code"}, finding), 0)
	loopDoneContract(t, c, w, task, hold, h2)
	loopCapContract(t, c)
	c.sets()
}

func loopClaimReview(t *testing.T, c *loopContract, builder, other loopWorker, task string) loopWorker {
	t.Helper()
	roots := []string{other.root, gittest.AddWorktree(t, builder.root)}
	type reply struct {
		out, stderr string
		err         error
	}
	replies := make([]reply, 2)
	var wg sync.WaitGroup
	for i := range roots {
		wg.Go(func() {
			replies[i].out, replies[i].stderr, replies[i].err = closureRunCommand(roots[i], builder.binary, "--json", "next", "--claim", "--actor", fmt.Sprintf("agent:reviewer%d", i), "--model", "claude/opus-5.5")
		})
	}
	wg.Wait()
	winner := -1
	for i, r := range replies {
		if r.err != nil || r.stderr != "" {
			t.Fatalf("reviewer %d: %v %s", i, r.err, r.stderr)
		}
		got := p23Envelope(t, r.out)
		switch p21PathString(t, got, "result", "kind") {
		case "review":
			if winner >= 0 || p21PathString(t, got, "result", "task", "id") != task {
				t.Fatal("review served twice or wrong task")
			}
			winner = i
			c.observe("next-review", got, 0)
		case "none":
			c.observe("next-review-reserved", got, 0)
		default:
			t.Fatalf("unexpected review selection: %v", got)
		}
	}
	if winner < 0 {
		t.Fatal("neither reviewer got work")
	}
	w := loopWorker{t, builder.binary, roots[winner], fmt.Sprintf("agent:reviewer%d", winner)}
	c.observe("heartbeat-review", w.call("heartbeat", "--hold", p21PathString(t, w.next("", true), "result", "hold", "id"), "--actor", w.actor), 0)
	return w
}

func loopDoneContract(t *testing.T, c *loopContract, w loopWorker, task, hold string, h map[string]any) {
	t.Helper()
	j := newDoneObservation(t, w.root, w.binary, task, hold, h)
	green := j.ci
	j.ci = []any{}
	for _, phase := range []string{"waiting", "completed", "empty"} {
		if phase == "completed" {
			j.ci = green
		}
		out, stderr, err := j.run(t, "--json", "done", "--actor", "agent:finisher")
		if err != nil || stderr != "" {
			t.Fatalf("done %s: %v %s", phase, err, stderr)
		}
		c.observe("done-"+phase, p23Envelope(t, out), 0)
	}
	// A keyed explicit finish in a separate fixture pins status: done too.
	j2 := newDoneJourney(t, "standard")
	out, stderr, err := j2.run(t, "--json", "done", "--actor", "agent:finisher")
	if err != nil || stderr != "" {
		t.Fatalf("refused batch: %v %s", err, stderr)
	}
	c.observe("done-refused", p23Envelope(t, out), 0)
	j2.review(t, "primary", "APPROVE", "primary")
	j2.requests[donePRRequest] = map[string]any{"error": "offline"}
	out, stderr, err = j2.run(t, "--json", "done", "--actor", "agent:finisher")
	if err != nil || stderr != "" {
		t.Fatalf("offline batch: %v %s", err, stderr)
	}
	c.observe("done-github-unavailable", p23Envelope(t, out), 0)
	delete(j2.requests, donePRRequest)
	out, stderr, err = j2.run(t, "--json", "done", j2.task, "--actor", "agent:finisher", "--idempotency-key", "finish")
	if err != nil || stderr != "" {
		t.Fatalf("explicit done: %v %s", err, stderr)
	}
	c.observe("done-task", p23Envelope(t, out), 0)
}

func TestLoopContractReviewProfiles(t *testing.T) {
	t.Parallel()
	c := &loopContract{t: t, decl: readLoopDeclaration(t)}
	root, binary := handoffJourneyRepo(t, true)
	path := filepath.Join(root, ".fid.toml")
	raw := readP21File(t, path)
	policy := strings.Replace(string(raw), "test_adequacy=[]", "test_adequacy=['**/*_test.go']", 1) + "\n[lanes.assured]\nchecks=['test']\n"
	if err := os.WriteFile(path, []byte(policy), 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Review profiles policy")
	closureRun(t, root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	w := loopWorker{t, binary, root, "agent:builder"}
	task := w.create("task", "Assured review", "")
	w.call("issue", "update", task, "--lane", "assured", "--actor", "human:planner", "--idempotency-key", "assured")
	p22Write(t, root, "change_test.go", "package example\n")
	w.build("", task, "profiles")
	reviewer := loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:reviewer"}
	for _, stream := range []string{"primary", "separate"} {
		got := reviewer.wantNext("", false, "review", task)
		c.observe("next-review-"+stream+"-profiles", got, 0)
		reviewer.wantNext("", true, "review", task)
		in := map[string]any{"reviewer": map[string]any{"actor": reviewer.actor, "model": "claude/opus-5.5", "session_id": stream},
			"stream": stream, "profiles": []string{"code", "test-adequacy"}, "verdict": "APPROVE", "purpose": "Review all profiles", "verification": "Code and regression inspected", "findings": []any{}}
		input := p21Input(t, t.TempDir(), "round", in)
		c.observe("review-"+stream+"-profiles", reviewer.call("review", "record", task, "--input", input, "--actor", "human:courier", "--idempotency-key", stream), 0)
	}
	reviewer.wantNext("", true, "none", "")
}

func loopCapContract(t *testing.T, c *loopContract) {
	t.Helper()
	root, binary := handoffJourneyRepo(t, true)
	path := filepath.Join(root, ".fid.toml")
	raw := readP21File(t, path)
	raw = []byte(strings.Replace(string(raw), "[lanes.standard]\n", "[lanes.standard]\nrework_cap=1\n", 1))
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Cap at one")
	closureRun(t, root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	w := loopWorker{t, binary, root, "agent:builder"}
	task := w.create("task", "Capped task", "")
	w.build("", task, "capped")
	reviewer := loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:reviewer"}
	for round := 1; round <= 2; round++ {
		reviewer.wantNext("", true, "review", task)
		loopRound(reviewer, task, "CHANGES", fmt.Sprintf("changes-%d", round), []string{"code"})
		if round == 1 {
			taken := w.wantNext("", true, "rework", task)
			c.observe("next-at-cap", taken, 0)
			hold := p21PathString(t, taken, "result", "hold", "id")
			loopCommit(t, root, "Rework at cap")
			w.call("check", "--hold", hold, "--actor", w.actor, "--idempotency-key", "cap-check")
			w.call("handoff", "--hold", hold, "--actor", w.actor, "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "cap-handoff")
		}
	}
	got := w.wantNext("", true, "none", "")
	c.observe("next-past-cap", got, 0)
	if p21PathString(t, got, "result", "waiting", "kind") != "rework-cap" {
		t.Fatalf("cap waiting = %v", got)
	}
	loopQueueReason(t, got, task, "rework-cap-reached")
	unrelated := w.create("task", "Other work", "")
	w.wantNext("", true, "none", "")
	reviewer.wantNext("", true, "new", unrelated)
}
