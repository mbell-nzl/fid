//go:build integration

package integration

import (
	"encoding/json"
	"testing"
)

func TestP23ReviewClaimRejectsDirtyPathAdoption(t *testing.T) {
	t.Run("recover", func(t *testing.T) {
		a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, nil, "review")
		p22Write(t, a.workspace, "base.txt", "review recovery adoption\n")
		a.p22Observe(t, "terminated")
		a.p22Expire(t)
		before := p21JournalEventCount(t, a.workspace)
		response := a.p22RunInput(t, "run", "recover", a.actor, map[string]any{"prior_run_id": a.runID, "expected_run_version": a.runVersion, "expected_claim_version": a.claimVersion, "observation_id": a.observation, "reconciliation_digest": a.digest, "adopted_paths": []string{"base.txt"}, "authority": map[string]any{"mode": "ordinary", "actor_id": a.actor, "reason": "review recovery", "reference": nil, "grant_ref": nil}, "executor": map[string]any{"id": "executor-review-successor", "source": "activated-harness"}, "ttl_seconds": 60})
		assertP22ProcessFailure(t, "run.recover", "claim-mode-invalid", 4, response)
		if p21JournalEventCount(t, a.workspace) != before {
			t.Fatal("review recovery rejection appended an Event")
		}
	})

	t.Run("abandon", func(t *testing.T) {
		a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, nil, "review")
		p22Write(t, a.workspace, "base.txt", "review abandon adoption\n")
		a.p22Observe(t, "terminated")
		a.repository = p23CurrentRepository(t, a)
		before := p21JournalEventCount(t, a.workspace)
		response := a.p22RunInput(t, "run", "abandon", a.actor, map[string]any{"run_id": a.runID, "expected_version": a.runVersion, "expected_claim_version": a.claimVersion, "observation_id": a.observation, "authority": map[string]any{"mode": "ordinary", "actor_id": a.actor, "reason": "review abandon", "reference": nil, "grant_ref": nil}, "repository": a.repository})
		assertP22ProcessFailure(t, "run.abandon", "claim-mode-invalid", 4, response)
		if p21JournalEventCount(t, a.workspace) != before {
			t.Fatal("review abandon rejection appended an Event")
		}
	})
}

func TestP23ReviewRunRetainsPhase22Lifecycle(t *testing.T) {
	t.Run("close", func(t *testing.T) {
		a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, nil, "review")
		a.p22Observe(t, "terminated")
		response := a.p22RunInput(t, "run", "close", a.actor, map[string]any{"run_id": a.runID, "expected_version": a.runVersion, "observation_id": a.observation, "repository": a.repository})
		p23AssertRunSuccess(t, "run.close", response)
	})

	t.Run("abandon", func(t *testing.T) {
		a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, nil, "review")
		a.p22Observe(t, "terminated")
		response := a.p22RunInput(t, "run", "abandon", a.actor, map[string]any{"run_id": a.runID, "expected_version": a.runVersion, "expected_claim_version": a.claimVersion, "observation_id": a.observation, "authority": map[string]any{"mode": "ordinary", "actor_id": a.actor, "reason": "review abandoned", "reference": nil, "grant_ref": nil}, "repository": a.repository})
		p23AssertRunSuccess(t, "run.abandon", response)
	})

	t.Run("recover", func(t *testing.T) {
		a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, nil, "review")
		a.p22Observe(t, "terminated")
		a.p22Expire(t)
		response := a.p22RunInput(t, "run", "recover", a.actor, map[string]any{"prior_run_id": a.runID, "expected_run_version": a.runVersion, "expected_claim_version": a.claimVersion, "observation_id": a.observation, "reconciliation_digest": a.digest, "adopted_paths": []string{}, "authority": map[string]any{"mode": "ordinary", "actor_id": a.actor, "reason": "review recovered", "reference": nil, "grant_ref": nil}, "executor": map[string]any{"id": "executor-review-successor", "source": "activated-harness"}, "ttl_seconds": 60})
		p23AssertRunSuccess(t, "run.recover", response)
		var envelope map[string]any
		if err := json.Unmarshal([]byte(response.Stdout), &envelope); err != nil {
			t.Fatalf("decode review recovery: %v", err)
		}
		if mode := p21Path(t, envelope, "result", "claim", "mode"); mode != "review" {
			t.Fatalf("recovered Claim mode = %#v, want review", mode)
		}
	})
}

func p23CurrentRepository(t *testing.T, a *p22FIDRecoveryAdapter) map[string]any {
	t.Helper()
	response := a.p22Run(t, "run", "inspect", a.runID)
	if response.ExitCode != 0 {
		t.Fatalf("inspect review Run = %#v", response)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stdout), &envelope); err != nil {
		t.Fatalf("decode review Run inspection: %v", err)
	}
	a.p22Advance(t, "run.inspect", response.Stdout)
	return p21Path(t, envelope, "result", "reconciliation", "repository").(map[string]any)
}

func p23AssertRunSuccess(t *testing.T, operation string, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("%s = %#v, want success", operation, response)
	}
}
