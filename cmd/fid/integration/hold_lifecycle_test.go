//go:build integration

package integration

import "testing"

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
