//go:build integration

package integration

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestP22ResumeReadyCandidatesExecuteTheirActualOperations(t *testing.T) {
	t.Parallel()
	for _, operation := range []string{"run.inspect", "run.observe", "gate.run", "claim.renew", "claim.release"} {
		t.Run(operation, func(t *testing.T) {
			adapter := p22NewResumeActionAdapter(t, "activated-harness")
			p22AssertResumeActionState(t, adapter, operation, true)
			before := p21JournalEventCount(t, adapter.workspace)
			args := p22ResumeActionArgs(t, adapter, operation)
			stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, args...)
			if err != nil || stderr != "" || stdout == "" {
				t.Fatalf("advertised ready %s failed: err=%v stderr=%q stdout=%q", operation, err, stderr, stdout)
			}
			after := p21JournalEventCount(t, adapter.workspace)
			if operation == "run.inspect" && after != before {
				t.Fatalf("inspect appended Events: %d -> %d", before, after)
			}
			if operation != "run.inspect" && after <= before {
				t.Fatalf("mutation %s made no durable Event: %d -> %d", operation, before, after)
			}
		})
	}
}

func TestP22ResumeObserveRejectsUntrustedProducerDespiteCapabilities(t *testing.T) {
	t.Parallel()
	adapter := p22NewResumeActionAdapter(t, "kit")
	p22AssertResumeActionState(t, adapter, "run.observe", false)
	before := p21JournalEventCount(t, adapter.workspace)
	stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, p22ResumeActionArgs(t, adapter, "run.observe")...)
	if err == nil || stdout != "" || !strings.Contains(stderr, `"code":"authority-denied"`) {
		t.Fatalf("kit producer observation: err=%v stdout=%q stderr=%q", err, stdout, stderr)
	}
	if after := p21JournalEventCount(t, adapter.workspace); after != before {
		t.Fatalf("denied observation appended Events: %d -> %d", before, after)
	}
}

func p22AssertResumeActionState(t *testing.T, adapter *p22FIDRecoveryAdapter, operation string, ready bool) {
	t.Helper()
	run := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "get", adapter.runID)
	issueID := p21PathString(t, run, "result", "run", "issue_id")
	before := p21JournalEventCount(t, adapter.workspace)
	resume := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "resume", issueID)
	if after := p21JournalEventCount(t, adapter.workspace); after != before {
		t.Fatalf("resume appended Events: %d -> %d", before, after)
	}
	for _, value := range p21Path(t, resume, "result", "resume", "next_actions").([]any) {
		action := value.(map[string]any)
		if action["operation"] != operation {
			continue
		}
		reasons := action["reasons"].([]any)
		if (len(reasons) == 0) != ready {
			t.Fatalf("%s ready=%t: %#v", operation, ready, action)
		}
		if !ready && !p22HasReason(reasons, "authority-unavailable") {
			t.Fatalf("%s omitted authority reason: %#v", operation, action)
		}
		return
	}
	t.Fatalf("resume did not advertise %s: %#v", operation, resume)
}

func p22ResumeActionArgs(t *testing.T, adapter *p22FIDRecoveryAdapter, operation string) []string {
	t.Helper()
	args := append([]string{"--project", adapter.workspace, "--json"}, strings.Split(operation, ".")...)
	if operation == "run.inspect" {
		return append(args, adapter.runID)
	}
	var input map[string]any
	switch operation {
	case "run.observe":
		input = map[string]any{"run_id": adapter.runID, "expected_version": 1, "liveness": map[string]any{"executor_id": adapter.executorID, "state": "live", "source_ref": "process:resume-actions", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}}
	case "gate.run":
		input = map[string]any{"registration_id": adapter.registrationID, "registration_version": adapter.registrationVersion, "obligation_id": "OBL-p22-corpus", "run_id": adapter.runID, "expected_repository": map[string]any{"head_oid": adapter.repository["head_oid"], "dirty": adapter.repository["dirty"]}}
	case "claim.renew":
		input = map[string]any{"claim_id": adapter.claimID, "expected_version": 1, "ttl_seconds": 60}
	case "claim.release":
		input = map[string]any{"claim_id": adapter.claimID, "expected_version": 1, "reason": "completed"}
	default:
		t.Fatalf("unhandled action %s", operation)
	}
	return append(args, "--actor", adapter.actor, "--input", p21Input(t, adapter.workspace, "resume-action", input), "--idempotency-key", "resume-action")
}

func p22NewResumeActionAdapter(t *testing.T, producer string) *p22FIDRecoveryAdapter {
	t.Helper()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p22-corpus@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 22 corpus")
	p22Write(t, workspace, "base.txt", "base\n")
	if err := os.MkdirAll(filepath.Join(workspace, "project-docs"), 0o700); err != nil {
		t.Fatal(err)
	}
	p22Write(t, workspace, "project-docs/p22-corpus.md", "Phase 22 recovery policy\n")
	p22Write(t, workspace, ".gitignore", "*.json\n")
	p22Git(t, workspace, "add", "base.txt", ".gitignore", "project-docs/p22-corpus.md")
	p22Git(t, workspace, "commit", "-qm", "base")
	init := p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P22")
	projectID := p21PathString(t, init, "result", "project", "id")
	issue := p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Phase 22 conformance", "--actor", "agent:developer", "--idempotency-key", "p22-corpus-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")
	actor := "actor-recovery"
	capabilities := []string{"artifact.verify", "attention.resolve", "claim.acquire", "claim.expire", "envelope.activate", "envelope.complete", "envelope.create", "gate.register", "gate.run", "issue.transition", "run.abandon", "run.checkpoint", "run.close", "run.observe", "run.observe-budget", "run.recover", "run.start", "stop.raise", "stop.resolve"}
	capabilities = append(capabilities, "claim.renew", "claim.release", "evidence.link")
	sort.Strings(capabilities)
	activationID := p21Activation(t, binary, workspace, "prepare", "p22-corpus-activation-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p22-corpus", "phase": 21})
	profileDigest := p21Component(t, binary, workspace, activationID, "authority-profile", "p22-corpus-authority", []any{
		map[string]any{"actor_id": "actor-denied", "actor_type": "agent", "producer": "activated-harness", "provenance_ref": "activation:p22-corpus", "session_ref": nil, "capabilities": capabilities},
		map[string]any{"actor_id": actor, "actor_type": "agent", "producer": producer, "provenance_ref": "activation:p22-corpus", "session_ref": nil, "capabilities": capabilities},
	})
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "p22-corpus-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-p22-corpus", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	digests := []string{gateDigest, profileDigest}
	sort.Strings(digests)
	p21Activation(t, binary, workspace, "commit", "p22-corpus-activation-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	p21Activation(t, binary, workspace, "finalize", "p22-corpus-activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})
	repository := p22CorpusRepository(t, workspace, projectID)
	registered := p21Run(t, binary, workspace, "--project", workspace, "--json", "gate", "register", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-gate-register", map[string]any{"definition_id": "gate-p22-corpus", "definition_version": 1, "scope": map[string]any{"issue_id": issueID}}), "--idempotency-key", "p22-corpus-gate-register")
	registrationID := p21PathString(t, registered, "result", "registration", "id")
	registrationVersion := p21Path(t, registered, "result", "registration", "version")
	source := map[string]any{"digest": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("Phase 22 recovery policy\n"))), "display_label": "Phase 22 corpus", "external_id": nil, "path": "project-docs/p22-corpus.md", "section": nil, "source_type": "project-file"}
	envelopeInput := p21Input(t, workspace, "p22-corpus-envelope", map[string]any{"outcome": "Recover one interrupted Run", "non_goals": []string{}, "risk_result": map[string]any{"source": source, "result_digest": "sha256:" + p22Ones(64)}, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{}, "authority_grants": []any{map[string]any{"actions": capabilities[1:], "actor_selector": map[string]any{"actor_id": actor, "required_capabilities": capabilities[1:]}, "scope": map[string]any{"issue_id": issueID}}}, "obligations": []any{map[string]any{"id": "OBL-p22-corpus", "type": "gate-evidence", "applicability": "required", "guarded_actions": []string{"envelope.complete", "issue.transition-terminal-success"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": map[string]any{"id": registrationID, "version": registrationVersion}, "accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"}, "freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": source}}, "budgets": []any{map[string]any{"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 3, "warning_threshold": nil, "exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil}}})
	created := p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "create", "--actor", actor, "--input", envelopeInput, "--idempotency-key", "p22-corpus-envelope-create")
	envelopeID := p21PathString(t, created, "result", "envelope", "id")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-envelope-activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1}), "--idempotency-key", "p22-corpus-envelope-activate")
	claim := p21Run(t, binary, workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-claim", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "scope": map[string]any{"paths": []string{"base.txt"}}, "ttl_seconds": 60}), "--idempotency-key", "p22-corpus-claim-acquire")
	claimID := p21PathString(t, claim, "result", "claim", "id")
	run := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "start", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-run", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": claimID, "claim_version": 1, "executor": map[string]any{"id": "executor-corpus-prior", "source": "activated-harness"}, "obligation_id": nil, "expected_repository": p24ExpectedRepository(repository)}), "--idempotency-key", "p22-corpus-run-start")
	source["repository"] = repository
	return &p22FIDRecoveryAdapter{activation: map[string]any{"id": activationID}, binary: binary, workspace: workspace, projectID: projectID, actor: actor, repository: repository, claimID: claimID, claimVersion: 1, runID: p21PathString(t, run, "result", "run", "id"), runVersion: 1, executorID: "executor-corpus-prior", measurement: source, registrationID: registrationID, registrationVersion: registrationVersion}
}
