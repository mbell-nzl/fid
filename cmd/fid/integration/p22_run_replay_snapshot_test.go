//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestP22RunReplayRetainsHistoricalBudgetSnapshot(t *testing.T) {
	for _, operation := range []string{"run.observe", "run.checkpoint", "run.recover", "run.abandon"} {
		t.Run(operation, func(t *testing.T) {
			adapter := p22NewResumeActionAdapter(t, "activated-harness")
			if operation != "run.observe" {
				adapter.p22Observe(t, "terminated")
			}
			if operation == "run.recover" {
				adapter.p22Expire(t)
			}
			fixture := p22ConformanceCase{Name: "historical-snapshot", Operation: operation}
			original := adapter.Execute(t, fixture)
			before := p22SnapshotResponse(t, original)
			if operation == "run.abandon" {
				p22StartAfterAbandon(t, adapter, before)
			}
			gate := p22GateRun(t, adapter.binary, adapter.workspace, adapter.actor, p21Input(t, adapter.workspace, "later-budget-charge", map[string]any{"run_id": adapter.runID, "registration_id": adapter.registrationID, "registration_version": adapter.registrationVersion, "obligation_id": "OBL-p22-corpus", "repository": adapter.repository}), "later-budget-charge")
			if p21PathString(t, gate, "result", "receipt", "outcome") != "succeeded" {
				t.Fatalf("later Gate failed: %#v", gate)
			}
			count := p21JournalEventCount(t, adapter.workspace)
			after := p22SnapshotResponse(t, adapter.Execute(t, fixture))
			if p21JournalEventCount(t, adapter.workspace) != count {
				t.Fatal("replay appended an Event")
			}
			beforeResult := before["result"].(map[string]any)
			afterResult := after["result"].(map[string]any)
			if afterResult["replayed"] != true {
				t.Fatalf("retry did not replay: %#v", afterResult)
			}
			beforeResult["replayed"] = true
			if !reflect.DeepEqual(beforeResult, afterResult) || !reflect.DeepEqual(before["event"], after["event"]) {
				t.Fatalf("historical replay changed after later Budget charge:\noriginal=%#v\nreplay=%#v", before, after)
			}
		})
	}
}

func p22SnapshotResponse(t *testing.T, response p22ConformanceResponse) map[string]any {
	t.Helper()
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("snapshot command: exit=%d stderr=%q stdout=%q", response.ExitCode, response.Stderr, response.Stdout)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(response.Stdout), &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func p22StartAfterAbandon(t *testing.T, adapter *p22FIDRecoveryAdapter, response map[string]any) {
	t.Helper()
	run := p21Path(t, response, "result", "run").(map[string]any)
	claim := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "claim", "acquire", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "after-abandon-claim", map[string]any{"issue_id": run["issue_id"], "envelope_id": run["envelope_id"], "envelope_revision": run["envelope_revision"], "worktree_id": adapter.repository["worktree_id"], "scope": map[string]any{"paths": []string{"base.txt"}}, "ttl_seconds": 60}), "--idempotency-key", "after-abandon-claim")
	started := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "start", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "after-abandon-run", map[string]any{"issue_id": run["issue_id"], "envelope_id": run["envelope_id"], "envelope_revision": run["envelope_revision"], "claim_id": p21PathString(t, claim, "result", "claim", "id"), "claim_version": 1, "executor": map[string]any{"id": "executor-after-abandon", "source": "activated-harness"}, "obligation_id": nil, "repository": adapter.repository}), "--idempotency-key", "after-abandon-run")
	adapter.runID = p21PathString(t, started, "result", "run", "id")
}
