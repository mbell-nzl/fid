//go:build integration

package integration

import (
	"github.com/mbell-nzl/fid/internal/gittest"
	"testing"
)

var p24LifecycleFields = map[string]string{"hold.park": "--reason", "task.takeover": "--reason"}

func TestHoldLifecycleKeyedCorpus(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	closureRun(t, workspace, "git", "init")
	closureRunFID(t, binary, workspace, "init", "--prefix", "FID")
	p24StartKeyedCapture(workspace)
	p24LifecycleCommands(t, binary, workspace)
	captured := p24EndKeyedCapture(workspace)
	for operation, field := range p24LifecycleFields {
		entry, found := captured[operation]
		if !found {
			t.Fatalf("missing %s capture", operation)
		}
		p24ExerciseCapturedLeaf(t, workspace, binary, operation, field, entry)
		p24AssertFamilyConflict(t, workspace, binary, operation, entry, []string{"--actor"})
	}
	p24AssertRebuildMapParity(t, workspace, binary)
}

func p24SingleLifecycle(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	p24LifecycleCommands(t, a.binary, a.workspace)
}

func p24LifecycleCommands(t *testing.T, binary, workspace string) {
	t.Helper()
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, workspace, append([]string{"--project", workspace, "--json"}, args...)...)
	}
	created := call("issue", "create", "--kind", "task", "--title", "Lifecycle corpus", "--actor", "human:mark", "--idempotency-key", "lifecycle-create")
	task := p21PathString(t, created, "result", "issue", "id")
	taken := call("takeover", task, "--actor", "human:mark", "--reason", "Resume task", "--idempotency-key", "lifecycle-takeover")
	hold := p21PathString(t, taken, "result", "hold_id")
	call("heartbeat", "--hold", hold, "--actor", "human:mark")
	call("park", "--hold", hold, "--actor", "human:mark", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "lifecycle-park")
}

func p24LifecycleCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	root, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
	p21Run(t, binary, root, "init", "--prefix", "FID", "--json")
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, root, append([]string{"--json"}, args...)...)
	}
	if code == "hold-not-current" {
		return p24CodePair{root, "hold.heartbeat", []string{"heartbeat", "--hold", "HLD-missing", "--actor", "agent:builder"}}
	}
	kind := "task"
	if code == "unsupported-kind" {
		kind = "epic"
	}
	createArgs := []string{"issue", "create", "--kind", kind, "--title", "Held work", "--actor", "human:mark", "--idempotency-key", "create-first"}
	if kind == "epic" {
		createArgs = append(createArgs, "--scheduling", "prioritised")
	}
	created := call(createArgs...)
	task := p21PathString(t, created, "result", "issue", "id")
	takeover := func(id string) []string {
		return []string{"takeover", id, "--actor", "human:mark", "--reason", "Agent stuck"}
	}
	if code == "unsupported-kind" {
		return p24CodePair{root, "task.takeover", takeover(task)}
	}
	taken := call(takeover(task)...)
	if code == "already-held" {
		return p24CodePair{root, "task.takeover", takeover(task)}
	}
	call("issue", "transition", task, "--to", "in_progress", "--actor", "human:mark", "--idempotency-key", "start-first")
	call("park", "--hold", p21PathString(t, taken, "result", "hold_id"), "--actor", "human:mark", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "park-first")
	second := call("issue", "create", "--kind", "task", "--title", "Other work", "--actor", "human:mark", "--idempotency-key", "create-second")
	return p24CodePair{root, "task.takeover", takeover(p21PathString(t, second, "result", "issue", "id"))}
}
