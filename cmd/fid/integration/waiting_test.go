//go:build integration

package integration

import (
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestWaitingWorktreeJourney(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	root := gittest.NewWorkspace(t)
	call := func(worktree string, args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, worktree, append([]string{"--json"}, args...)...)
	}
	call(root, "init", "--prefix", "FID")
	first := call(root, "issue", "create", "--kind", "task", "--title", "Waiting work", "--priority", "0", "--actor", "human:mark", "--idempotency-key", "first")
	task := p21PathString(t, first, "result", "issue", "id")
	second := call(root, "issue", "create", "--kind", "task", "--title", "Other work", "--actor", "human:mark", "--idempotency-key", "second")
	otherTask := p21PathString(t, second, "result", "issue", "id")
	other := gittest.AddWorktree(t, root)
	held := call(root, "next", "--claim", "--actor", "agent:builder")
	hold := p21PathString(t, held, "result", "hold", "id")
	parked := call(root, "park", "--hold", hold, "--actor", "agent:builder", "--waiting", "Include timeouts?", "--idempotency-key", "wait")
	question := p21PathString(t, parked, "result", "question_id")
	blocked := call(root, "next", "--claim", "--actor", "agent:builder")
	if p21PathString(t, blocked, "result", "kind") != "none" || p21PathString(t, blocked, "result", "waiting", "question_id") != question {
		t.Fatalf("waiting next: %v", blocked)
	}
	if available := call(other, "next", "--claim", "--actor", "agent:other"); p21PathString(t, available, "result", "task", "id") != otherTask {
		t.Fatalf("other worktree: %v", available)
	}
	answer := call(other, "answer", task, "Yes", "--actor", "human:mark", "--idempotency-key", "answer")
	resumed := call(root, "next", "--claim", "--actor", "agent:builder")
	if p21PathString(t, resumed, "result", "kind") != "resume" || p21PathString(t, resumed, "result", "hold", "replaces") != hold || p21PathString(t, resumed, "result", "answer", "decision_id") != p21PathString(t, answer, "result", "decision_id") {
		t.Fatalf("resume: %v", resumed)
	}
	repeated := call(root, "next", "--claim", "--actor", "agent:builder")
	if !reflect.DeepEqual(p21Path(t, resumed, "result", "answer"), p21Path(t, repeated, "result", "answer")) {
		t.Fatalf("repeated resume lost answer: %v", repeated)
	}
	call(root, "rebuild")
	replay := call(other, "answer", task, "Yes", "--actor", "human:mark", "--idempotency-key", "answer")
	if !reflect.DeepEqual(answer, replay) {
		t.Fatalf("rebuild changed answer retry: %v, %v", answer, replay)
	}
}
