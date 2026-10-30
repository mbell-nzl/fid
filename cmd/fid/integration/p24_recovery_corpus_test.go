//go:build integration

package integration

import (
	"testing"
)

var p24RecoveryFields = map[string]string{
	"run.observe": "expected_version", "run.checkpoint": "reconciliation_digest",
	"claim.expire": "expected_version", "run.recover": "ttl_seconds",
	"run.observe-budget": "value", "stop.raise": "details",
	"run.abandon": "expected_version", "stop.resolve": "stop_id",
	"attention.resolve": "expected_version", "run.close": "expected_version",
}

func TestP24RecoveryKeyedLeafConflictAndCanonicalRetryCorpus(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	fixtures := p22RecoveryFixtures(t)
	type capturedCase struct {
		workspace string
		entry     p24CapturedKeyedCommand
	}
	captured := map[string]capturedCase{}
	var adapter *p22FIDRecoveryAdapter
	finish := func() {
		if adapter == nil {
			return
		}
		for operation, entry := range p24EndKeyedCapture(adapter.workspace) {
			if _, seen := captured[operation]; !seen {
				captured[operation] = capturedCase{workspace: adapter.workspace, entry: entry}
			}
		}
	}
	for _, fixture := range fixtures {
		if adapter == nil || fixture.Reset {
			finish()
			workspace := t.TempDir()
			p24StartKeyedCapture(workspace)
			adapter = p22NewFIDRecoveryAdapterAt(t, binary, workspace, "agent", 60, []string{}, nil)
		}
		response := adapter.Execute(t, fixture)
		if response.ExitCode != fixture.ExitCode {
			t.Fatalf("%s setup = exit %d stderr=%q; want %d", fixture.Name, response.ExitCode, response.Stderr, fixture.ExitCode)
		}
	}
	finish()
	for operation, name := range p24RecoveryFields {
		fixture, ok := captured[operation]
		if !ok {
			t.Fatalf("recovery fixture did not capture keyed %s", operation)
		}
		p24ExerciseCapturedLeaf(t, fixture.workspace, binary, operation, name, fixture.entry)
	}
	for operation, second := range map[string]string{"attention.resolve": "attention_id", "stop.raise": "resolver_capability"} {
		fixture := captured[operation]
		p24AssertFamilyConflict(t, fixture.workspace, binary, operation, fixture.entry, []string{p24RecoveryFields[operation], second})
		p24AssertFamilyConflict(t, fixture.workspace, binary, operation, fixture.entry, []string{"--actor"})
	}
	workspaces := map[string]bool{}
	for _, fixture := range captured {
		workspaces[fixture.workspace] = true
	}
	for workspace := range workspaces {
		p24AssertRebuildMapParity(t, workspace, binary)
	}
}
