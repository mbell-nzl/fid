//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/project"
)

func TestP22EarlyRecoveryJourney(t *testing.T) {
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p22@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 22")
	p22Write(t, workspace, "base.txt", "base\n")
	p22Write(t, workspace, ".gitignore", "*.json\n")
	p22Git(t, workspace, "add", "base.txt", ".gitignore")
	p22Git(t, workspace, "commit", "-qm", "base")

	init := p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P22")
	projectID := p21PathString(t, init, "result", "project", "id")
	issue := p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Phase 22 recovery journey", "--actor", "agent:developer", "--idempotency-key", "p22-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")

	capabilities := []string{"artifact.verify", "claim.acquire", "claim.expire", "envelope.activate", "envelope.create", "gate.register", "run.checkpoint", "run.observe", "run.recover", "run.start"}
	activationID := p21Activation(t, binary, workspace, "prepare", "p22-activation-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p22", "phase": 21})
	profileDigest := p21Component(t, binary, workspace, activationID, "authority-profile", "p22-authority", []any{map[string]any{"actor_id": "actor-recovery", "actor_type": "agent", "producer": "activated-harness", "provenance_ref": "activation:p22", "session_ref": nil, "capabilities": capabilities}})
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "p22-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-p22", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	digests := []string{gateDigest, profileDigest}
	sort.Strings(digests)
	p21Activation(t, binary, workspace, "commit", "p22-activation-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	p21Activation(t, binary, workspace, "finalize", "p22-activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})

	repository := p22Repository(t, workspace, projectID)
	registered := p21Run(t, binary, workspace, "--project", workspace, "--json", "gate", "register", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-gate-register", map[string]any{"definition_id": "gate-p22", "definition_version": 1, "worktree_id": repository["worktree_id"], "scope": map[string]any{"issue_id": issueID}, "repository": repository}), "--idempotency-key", "p22-gate-register")
	registrationID := p21PathString(t, registered, "result", "registration", "id")
	registrationVersion := p21Path(t, registered, "result", "registration", "version")
	source := map[string]any{"digest": "sha256:" + p22Zeros(64), "display_label": "Phase 22", "external_id": nil, "path": "project-docs/p22.md", "repository": repository, "section": nil, "source_type": "project-file"}
	envelopeInput := p21Input(t, workspace, "p22-envelope", map[string]any{
		"outcome": "Recover one interrupted Run", "non_goals": []string{}, "risk_result": map[string]any{"source": source, "result_digest": "sha256:" + p22Ones(64)}, "repository_scope": repository, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{},
		"authority_grants": []any{map[string]any{"actions": capabilities[1:], "actor_selector": map[string]any{"actor_id": "actor-recovery", "required_capabilities": capabilities[1:]}, "scope": map[string]any{"issue_id": issueID}}},
		"obligations":      []any{map[string]any{"id": "OBL-p22", "type": "gate-evidence", "applicability": "required", "guarded_actions": []string{"envelope.complete", "issue.transition-terminal-success"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": map[string]any{"id": registrationID, "version": registrationVersion}, "accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"}, "freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": source}},
		"budgets":          []any{map[string]any{"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 3, "warning_threshold": nil, "exhaustion_action": "unavailable", "consumed": nil, "remaining": nil}},
	})
	created := p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-recovery", "--input", envelopeInput, "--idempotency-key", "p22-envelope-create")
	envelopeID := p21PathString(t, created, "result", "envelope", "id")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-envelope-activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1}), "--idempotency-key", "p22-envelope-activate")

	claim := p21Run(t, binary, workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-claim", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "worktree_id": repository["worktree_id"], "scope": map[string]any{"paths": []string{"base.txt"}}, "ttl_seconds": 60}), "--idempotency-key", "p22-claim-acquire")
	claimID := p21PathString(t, claim, "result", "claim", "id")
	run := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-run", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": claimID, "claim_version": 1, "executor": map[string]any{"id": "executor-prior", "source": "activated-harness"}, "obligation_id": nil, "repository": repository}), "--idempotency-key", "p22-run-start")
	runID := p21PathString(t, run, "result", "run", "id")

	terminatedAt := time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)
	observed := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-observe", map[string]any{"run_id": runID, "expected_version": 1, "liveness": map[string]any{"executor_id": "executor-prior", "state": "terminated", "source_ref": "process:p22", "observed_at": terminatedAt}}), "--idempotency-key", "p22-observe")
	observationID := p21PathString(t, observed, "result", "observation", "id")

	inspected := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "inspect", runID)
	digest := p21PathString(t, inspected, "result", "reconciliation", "digest")
	checkpointed := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "checkpoint", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-checkpoint", map[string]any{"run_id": runID, "expected_version": 2, "reconciliation_digest": digest}), "--idempotency-key", "p22-checkpoint")
	runVersion := int(p21Path(t, checkpointed, "result", "run", "version").(float64))
	expired := p21Run(t, binary, workspace, "--project", workspace, "--json", "claim", "expire", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-expire", map[string]any{"claim_id": claimID, "expected_version": 1}), "--idempotency-key", "p22-expire")
	claimVersion := int(p21Path(t, expired, "result", "claim", "version").(float64))

	recovered := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "recover", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-recover", map[string]any{"prior_run_id": runID, "expected_run_version": runVersion, "expected_claim_version": claimVersion, "observation_id": observationID, "reconciliation_digest": digest, "adopted_paths": []string{}, "authority": map[string]any{"mode": "ordinary", "actor_id": "actor-recovery", "reason": "trusted termination", "reference": nil, "grant_ref": nil}, "executor": map[string]any{"id": "executor-successor", "source": "activated-harness"}, "ttl_seconds": 60}), "--idempotency-key", "p22-recover")
	successorID := p21PathString(t, recovered, "result", "run", "id")
	if successorID == runID || p21PathString(t, recovered, "result", "recovery", "prior_run_id") != runID {
		t.Fatalf("recovered Run = %#v, want linked successor for %q", recovered, runID)
	}

	successor := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "inspect", successorID)
	if got := p21PathString(t, successor, "result", "run", "predecessor_run_id"); got != runID {
		t.Fatalf("successor predecessor_run_id = %q, want %q", got, runID)
	}
	for _, budget := range p21Path(t, successor, "result", "run", "budgets").([]any) {
		value := budget.(map[string]any)
		if value["consumed"] != nil || value["remaining"] != nil {
			t.Fatalf("early recovery budget = %#v, want unavailable null counters", value)
		}
	}
}

func p22Repository(t *testing.T, workspace, projectID string) map[string]any {
	t.Helper()
	observation, err := project.ObserveRepository(workspace)
	if err != nil {
		t.Fatalf("observe journey repository: %v", err)
	}
	return map[string]any{"base_oid": observation.HeadOID, "capture_source": "git", "dirty": observation.Dirty, "dirty_digest": observation.DirtyDigest, "head_oid": observation.HeadOID, "project_id": projectID, "tree_oid": observation.TreeOID, "worktree_id": observation.WorktreeID}
}

func p22Git(t *testing.T, workspace string, arguments ...string) {
	t.Helper()
	closureRun(t, workspace, "git", arguments...)
}

func p22Write(t *testing.T, workspace, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(workspace, name), []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func p22Zeros(count int) string { return strings.Repeat("0", count) }
func p22Ones(count int) string  { return strings.Repeat("1", count) }
