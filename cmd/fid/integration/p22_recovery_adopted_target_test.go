//go:build integration

package integration

import (
	"testing"
	"time"
)

func TestP22RecoveryAdoptsChangedTargetForSuccessorEvidence(t *testing.T) {
	t.Parallel()
	adapter := p22NewResumeActionAdapter(t, "activated-harness")
	priorRunID := adapter.runID
	priorRepository := adapter.repository

	priorGate := p22GateRun(t, adapter.binary, adapter.workspace, adapter.actor, p21Input(t, adapter.workspace, "p22-adopted-gate-a", map[string]any{
		"registration_id": adapter.registrationID, "registration_version": adapter.registrationVersion, "obligation_id": "OBL-p22-corpus", "run_id": priorRunID, "expected_repository": map[string]any{"head_oid": priorRepository["head_oid"], "dirty": priorRepository["dirty"]},
	}), "p22-adopted-gate-a")
	priorReceiptID := p21PathString(t, priorGate, "result", "receipt", "id")
	priorEvidence := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "evidence", "link", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "p22-adopted-evidence-a", map[string]any{
		"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": priorReceiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{},
	}), "--idempotency-key", "p22-adopted-evidence-a")
	priorEvidenceID := p21PathString(t, priorEvidence, "result", "evidence", "id")
	priorRun := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "get", priorRunID)
	adapter.runVersion = p22Number(p21Path(t, priorRun, "result", "run", "version"))

	p22Write(t, adapter.workspace, "base.txt", "target B adopted after binding\n")
	adapter.p22Observe(t, "terminated")
	adapter.p22Expire(t)
	recovered := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "recover", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "p22-adopted-recover", map[string]any{
		"prior_run_id": priorRunID, "expected_run_version": adapter.runVersion, "expected_claim_version": adapter.claimVersion, "observation_id": adapter.observation, "reconciliation_digest": adapter.digest,
		"adopted_paths": []string{"base.txt"}, "authority": map[string]any{"mode": "ordinary", "actor_id": adapter.actor, "reason": "adopt observed target B", "reference": nil, "grant_ref": nil},
		"executor": map[string]any{"id": "executor-adopted-successor", "source": "activated-harness"}, "ttl_seconds": 60,
	}), "--idempotency-key", "p22-adopted-recover")
	successorRunID := p21PathString(t, recovered, "result", "run", "id")
	successorRepository := p21Path(t, recovered, "result", "run", "start_repository").(map[string]any)
	if successorRunID == priorRunID || successorRepository["dirty_digest"] == priorRepository["dirty_digest"] {
		t.Fatalf("recovered target = run %q repository %#v, want linked successor at changed target B", successorRunID, successorRepository)
	}

	successorGate := p22GateRun(t, adapter.binary, adapter.workspace, adapter.actor, p21Input(t, adapter.workspace, "p22-adopted-gate-b", map[string]any{
		"registration_id": adapter.registrationID, "registration_version": adapter.registrationVersion, "obligation_id": "OBL-p22-corpus", "run_id": successorRunID, "expected_repository": map[string]any{"head_oid": successorRepository["head_oid"], "dirty": successorRepository["dirty"]},
	}), "p22-adopted-gate-b")
	successorReceiptID := p21PathString(t, successorGate, "result", "receipt", "id")
	successorEvidence := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "evidence", "link", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "p22-adopted-evidence-b", map[string]any{
		"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": successorReceiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{},
	}), "--idempotency-key", "p22-adopted-evidence-b")
	successorEvidenceID := p21PathString(t, successorEvidence, "result", "evidence", "id")
	if successorReceiptID == priorReceiptID || successorEvidenceID == priorEvidenceID {
		t.Fatalf("successor facts = Receipt %q Evidence %q, want new immutable facts after target adoption", successorReceiptID, successorEvidenceID)
	}
	successorReceipt := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "receipt", "get", successorReceiptID)
	if got := p21PathString(t, successorReceipt, "result", "receipt", "registration_id"); got != adapter.registrationID {
		t.Fatalf("successor Registration = %q, want retained %q", got, adapter.registrationID)
	}
	if got := p21PathString(t, successorReceipt, "result", "receipt", "obligation_id"); got != "OBL-p22-corpus" {
		t.Fatalf("successor Obligation = %q, want retained OBL-p22-corpus", got)
	}
	prior := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "evidence", "get", priorEvidenceID)
	if got := p21PathString(t, prior, "result", "evidence", "id"); got != priorEvidenceID {
		t.Fatalf("recovery changed prior Evidence identity to %q, want %q", got, priorEvidenceID)
	}
}

func TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof(t *testing.T) {
	t.Parallel()
	adapter := p22NewResumeActionAdapter(t, "activated-harness")
	priorRunID := adapter.runID
	priorRepository := adapter.repository
	priorGate := p22GateRun(t, adapter.binary, adapter.workspace, adapter.actor, p21Input(t, adapter.workspace, "p22-adopted-terminal-gate-a", map[string]any{
		"registration_id": adapter.registrationID, "registration_version": adapter.registrationVersion, "obligation_id": "OBL-p22-corpus", "run_id": priorRunID, "expected_repository": map[string]any{"head_oid": priorRepository["head_oid"], "dirty": priorRepository["dirty"]},
	}), "p22-adopted-terminal-gate-a")
	priorReceiptID := p21PathString(t, priorGate, "result", "receipt", "id")
	priorRun := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "get", priorRunID)
	adapter.runVersion = p22Number(p21Path(t, priorRun, "result", "run", "version"))

	p22Write(t, adapter.workspace, "base.txt", "target B permits terminal transition\n")
	adapter.p22Observe(t, "terminated")
	adapter.p22Expire(t)
	recovered := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "recover", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "p22-adopted-terminal-recover", map[string]any{
		"prior_run_id": priorRunID, "expected_run_version": adapter.runVersion, "expected_claim_version": adapter.claimVersion, "observation_id": adapter.observation, "reconciliation_digest": adapter.digest,
		"adopted_paths": []string{"base.txt"}, "authority": map[string]any{"mode": "ordinary", "actor_id": adapter.actor, "reason": "adopt observed terminal target B", "reference": nil, "grant_ref": nil},
		"executor": map[string]any{"id": "executor-adopted-terminal", "source": "activated-harness"}, "ttl_seconds": 60,
	}), "--idempotency-key", "p22-adopted-terminal-recover")
	successorRunID := p21PathString(t, recovered, "result", "run", "id")
	successorRepository := p21Path(t, recovered, "result", "run", "start_repository").(map[string]any)
	if successorRunID == priorRunID || successorRepository["dirty_digest"] == priorRepository["dirty_digest"] {
		t.Fatalf("recovered terminal target = run %q repository %#v, want successor at target B", successorRunID, successorRepository)
	}

	successorGate := p22GateRun(t, adapter.binary, adapter.workspace, adapter.actor, p21Input(t, adapter.workspace, "p22-adopted-terminal-gate-b", map[string]any{
		"registration_id": adapter.registrationID, "registration_version": adapter.registrationVersion, "obligation_id": "OBL-p22-corpus", "run_id": successorRunID, "expected_repository": map[string]any{"head_oid": successorRepository["head_oid"], "dirty": successorRepository["dirty"]},
	}), "p22-adopted-terminal-gate-b")
	successorReceiptID := p21PathString(t, successorGate, "result", "receipt", "id")
	p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "evidence", "link", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "p22-adopted-terminal-evidence-b", map[string]any{
		"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": successorReceiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{},
	}), "--idempotency-key", "p22-adopted-terminal-evidence-b")
	if successorReceiptID == priorReceiptID {
		t.Fatalf("successor Receipt = prior Receipt %q, want new retained proof at B", priorReceiptID)
	}
	priorReceipt := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "receipt", "get", priorReceiptID)
	if got := p21PathString(t, priorReceipt, "result", "receipt", "id"); got != priorReceiptID {
		t.Fatalf("recovery changed prior Receipt identity to %q, want %q", got, priorReceiptID)
	}

	run := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "get", successorRunID)
	issueID := p21PathString(t, run, "result", "run", "issue_id")
	for _, step := range []struct{ state, revision string }{
		{state: "building", revision: "1"},
		{state: "dev_complete", revision: "2"},
		{state: "qa_review", revision: "3"},
		{state: "qa_pass", revision: "4"},
	} {
		p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "issue", "transition", issueID, "--to", step.state, "--actor", "agent:"+adapter.actor, "--expected-revision", step.revision, "--idempotency-key", "p22-adopted-terminal-"+step.state)
	}
	observed := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "observe", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "p22-adopted-terminal-observe-b", map[string]any{
		"run_id": successorRunID, "expected_version": p21Path(t, run, "result", "run", "version"),
		"liveness": map[string]any{"executor_id": "executor-adopted-terminal", "state": "terminated", "source_ref": "p22:adopted-target-b", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)},
	}), "--idempotency-key", "p22-adopted-terminal-observe-b")
	p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "close", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "p22-adopted-terminal-close-b", map[string]any{
		"run_id": successorRunID, "expected_version": p21Path(t, observed, "result", "run", "version"), "observation_id": p21PathString(t, observed, "result", "observation", "id"), "expected_repository": p24ExpectedRepository(successorRepository),
	}), "--idempotency-key", "p22-adopted-terminal-close-b")
	p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "envelope", "complete", "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "p22-adopted-terminal-complete-b", map[string]any{
		"envelope_id": p21PathString(t, run, "result", "run", "envelope_id"), "expected_revision": 1,
	}), "--idempotency-key", "p22-adopted-terminal-complete-b")
}
