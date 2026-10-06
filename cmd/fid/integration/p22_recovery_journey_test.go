//go:build integration

package integration

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/project"
)

func TestP22RecoveryConformanceCorpusHasSharedCases(t *testing.T) {
	t.Parallel()
	fixture, err := os.ReadFile(filepath.Join("testdata", "p22_conformance", "manifest.json"))
	if err != nil {
		t.Fatalf("read Phase 22 recovery conformance manifest: %v", err)
	}
	var manifest struct {
		Cases []json.RawMessage `json:"cases"`
	}
	if err := json.Unmarshal(fixture, &manifest); err != nil {
		t.Fatalf("decode Phase 22 recovery conformance manifest: %v", err)
	}
	if len(manifest.Cases) == 0 {
		t.Fatal("Phase 22 recovery conformance manifest has no shared cases")
	}
}

func TestP22EarlyRecoveryJourney(t *testing.T) {
	t.Parallel()
	for _, dirty := range []bool{false, true} {
		name := "clean"
		if dirty {
			name = "dirty"
		}
		t.Run(name, func(t *testing.T) { p22RecoveryFreshnessJourney(t, dirty) })
	}
}

func p22RecoveryFreshnessJourney(t *testing.T, dirty bool) {
	t.Helper()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p22@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 22")
	p22Write(t, workspace, "base.txt", "base\n")
	if err := os.MkdirAll(filepath.Join(workspace, "project-docs"), 0o700); err != nil {
		t.Fatal(err)
	}
	p22Write(t, workspace, "project-docs/p22.md", "Phase 22 recovery policy\n")
	p22Write(t, workspace, ".gitignore", "*.json\n")
	p22Git(t, workspace, "add", "base.txt", ".gitignore", "project-docs/p22.md")
	p22Git(t, workspace, "commit", "-qm", "base")

	init := p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P22")
	projectID := p21PathString(t, init, "result", "project", "id")
	issue := p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Phase 22 recovery journey", "--actor", "agent:developer", "--idempotency-key", "p22-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")

	capabilities := []string{"artifact.verify", "claim.acquire", "claim.expire", "envelope.activate", "envelope.create", "evidence.link", "gate.register", "gate.run", "run.checkpoint", "run.observe", "run.recover", "run.start"}
	activationID := p21Activation(t, binary, workspace, "prepare", "p22-activation-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p22", "phase": 21})
	profileDigest := p21Component(t, binary, workspace, activationID, "authority-profile", "p22-authority", []any{map[string]any{"actor_id": "actor-recovery", "actor_type": "agent", "producer": "activated-harness", "provenance_ref": "activation:p22", "session_ref": nil, "capabilities": capabilities}})
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "p22-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-p22", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	digests := []string{gateDigest, profileDigest}
	sort.Strings(digests)
	p21Activation(t, binary, workspace, "commit", "p22-activation-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	p21Activation(t, binary, workspace, "finalize", "p22-activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})

	adoptedPaths := []string{}
	if dirty {
		p22Write(t, workspace, "base.txt", "base*\n")
		adoptedPaths = []string{"base.txt"}
	}
	repository := p22Repository(t, workspace, projectID)
	registered := p21Run(t, binary, workspace, "--project", workspace, "--json", "gate", "register", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-gate-register", map[string]any{"definition_id": "gate-p22", "definition_version": 1, "scope": map[string]any{"issue_id": issueID}}), "--idempotency-key", "p22-gate-register")
	registrationID := p21PathString(t, registered, "result", "registration", "id")
	registrationVersion := p21Path(t, registered, "result", "registration", "version")
	source := map[string]any{"digest": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("Phase 22 recovery policy\n"))), "display_label": "Phase 22", "external_id": nil, "path": "project-docs/p22.md", "section": nil, "source_type": "project-file"}
	envelopeInput := p21Input(t, workspace, "p22-envelope", map[string]any{
		"outcome": "Recover one interrupted Run", "non_goals": []string{}, "risk_result": map[string]any{"source": source, "result_digest": "sha256:" + p22Ones(64)}, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{},
		"authority_grants": []any{map[string]any{"actions": capabilities[1:], "actor_selector": map[string]any{"actor_id": "actor-recovery", "required_capabilities": capabilities[1:]}, "scope": map[string]any{"issue_id": issueID}}},
		"obligations":      []any{map[string]any{"id": "OBL-p22", "type": "gate-evidence", "applicability": "required", "guarded_actions": []string{"envelope.complete", "issue.transition-terminal-success"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": map[string]any{"id": registrationID, "version": registrationVersion}, "accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"}, "freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": source}},
		"budgets":          []any{map[string]any{"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 3, "warning_threshold": nil, "exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil}},
	})
	created := p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-recovery", "--input", envelopeInput, "--idempotency-key", "p22-envelope-create")
	envelopeID := p21PathString(t, created, "result", "envelope", "id")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-envelope-activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1}), "--idempotency-key", "p22-envelope-activate")

	claim := p21Run(t, binary, workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-claim", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "scope": map[string]any{"paths": []string{"base.txt"}}, "ttl_seconds": 60}), "--idempotency-key", "p22-claim-acquire")
	claimID := p21PathString(t, claim, "result", "claim", "id")
	run := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-run", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": claimID, "claim_version": 1, "executor": map[string]any{"id": "executor-prior", "source": "activated-harness"}, "obligation_id": nil, "repository": repository}), "--idempotency-key", "p22-run-start")
	runID := p21PathString(t, run, "result", "run", "id")
	gateInput := p21Input(t, workspace, "p22-gate-run", map[string]any{"registration_id": registrationID, "registration_version": registrationVersion, "obligation_id": "OBL-p22", "run_id": runID, "repository": repository})
	gate := p22GateRun(t, binary, workspace, "actor-recovery", gateInput, "p22-gate-run")
	receiptID := p21PathString(t, gate, "result", "receipt", "id")
	linked := p21Run(t, binary, workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-evidence-link", map[string]any{"obligation_id": "OBL-p22", "evidence_type": "receipt", "receipt_id": receiptID, "review_id": nil, "human_outcome_id": nil, "target_repository": repository, "required_artifact_ids": []string{}}), "--idempotency-key", "p22-evidence-link")
	evidenceID := p21PathString(t, linked, "result", "evidence", "id")

	terminatedAt := time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)
	observed := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-observe", map[string]any{"run_id": runID, "expected_version": 2, "liveness": map[string]any{"executor_id": "executor-prior", "state": "terminated", "source_ref": "process:p22", "observed_at": terminatedAt}}), "--idempotency-key", "p22-observe")
	observationID := p21PathString(t, observed, "result", "observation", "id")

	inspected := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "inspect", runID)
	digest := p21PathString(t, inspected, "result", "reconciliation", "digest")
	checkpointed := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "checkpoint", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-checkpoint", map[string]any{"run_id": runID, "expected_version": 3, "reconciliation_digest": digest}), "--idempotency-key", "p22-checkpoint")
	runVersion := int(p21Path(t, checkpointed, "result", "run", "version").(float64))
	inspected = p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "inspect", runID)
	digest = p21PathString(t, inspected, "result", "reconciliation", "digest")
	expired := p21Run(t, binary, workspace, "--project", workspace, "--json", "claim", "expire", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-expire", map[string]any{"claim_id": claimID, "expected_version": 1}), "--idempotency-key", "p22-expire")
	claimVersion := int(p21Path(t, expired, "result", "claim", "version").(float64))

	recovered := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "recover", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-recover", map[string]any{"prior_run_id": runID, "expected_run_version": runVersion, "expected_claim_version": claimVersion, "observation_id": observationID, "reconciliation_digest": digest, "adopted_paths": adoptedPaths, "authority": map[string]any{"mode": "ordinary", "actor_id": "actor-recovery", "reason": "trusted termination", "reference": nil, "grant_ref": nil}, "executor": map[string]any{"id": "executor-successor", "source": "activated-harness"}, "ttl_seconds": 60}), "--idempotency-key", "p22-recover")
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
		if value["consumed"] != float64(1) || value["remaining"] != float64(2) {
			t.Fatalf("recovered budget = %#v, want one preserved attempt charge and remaining two", value)
		}
	}

	p22Git(t, workspace, "branch", "-m", "recovery-renamed")
	freshResume := p21Run(t, binary, workspace, "--project", workspace, "--json", "resume", issueID, "--repository", workspace)
	fresh := p21Path(t, freshResume, "result", "resume", "obligations").([]any)[0].(map[string]any)
	if fresh["freshness"] != "fresh" || fresh["validity"] != "valid" {
		t.Fatalf("unchanged recovery/branch-rename resume = %#v, want retained fresh valid proof", fresh)
	}
	freshRepository := p21Path(t, freshResume, "result", "resume", "repository").(map[string]any)
	recoveredRepository := p21Path(t, recovered, "result", "run", "start_repository").(map[string]any)
	if freshRepository["dirty"] != dirty || freshRepository["dirty_digest"] != recoveredRepository["dirty_digest"] {
		t.Fatalf("unchanged recovery repository = %#v, want the successor's reconciled target %#v", freshRepository, recoveredRepository)
	}
	successorGateInput := p21Input(t, workspace, "p22-successor-gate-run", map[string]any{"registration_id": registrationID, "registration_version": registrationVersion, "obligation_id": "OBL-p22", "run_id": successorID, "repository": recoveredRepository})
	successorGate := p22GateRun(t, binary, workspace, "actor-recovery", successorGateInput, "p22-successor-gate-run")
	successorReceiptID := p21PathString(t, successorGate, "result", "receipt", "id")
	successorLinked := p21Run(t, binary, workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-recovery", "--input", p21Input(t, workspace, "p22-successor-evidence-link", map[string]any{"obligation_id": "OBL-p22", "evidence_type": "receipt", "receipt_id": successorReceiptID, "review_id": nil, "human_outcome_id": nil, "target_repository": recoveredRepository, "required_artifact_ids": []string{}}), "--idempotency-key", "p22-successor-evidence-link")
	if successorEvidenceID := p21PathString(t, successorLinked, "result", "evidence", "id"); successorEvidenceID == evidenceID {
		t.Fatalf("reconciled successor Evidence identity = %q, want new Evidence for successor Receipt", successorEvidenceID)
	}
	beforeRebuild := p21Run(t, binary, workspace, "--project", workspace, "--json", "resume", issueID, "--repository", workspace)
	p21Run(t, binary, workspace, "--project", workspace, "--json", "rebuild")
	afterRebuild := p21Run(t, binary, workspace, "--project", workspace, "--json", "resume", issueID, "--repository", workspace)
	if !reflect.DeepEqual(p21Path(t, beforeRebuild, "result", "resume"), p21Path(t, afterRebuild, "result", "resume")) {
		t.Fatalf("recovered resume changed across rebuild: before=%#v after=%#v", beforeRebuild, afterRebuild)
	}
	changed := "bas!\n"
	if dirty {
		changed = "base!\n"
	}
	p22Write(t, workspace, "base.txt", changed)
	staleResume := p21Run(t, binary, workspace, "--project", workspace, "--json", "resume", issueID, "--repository", workspace)
	stale := p21Path(t, staleResume, "result", "resume", "obligations").([]any)[0].(map[string]any)
	if stale["freshness"] != "stale" || stale["validity"] != "valid" || !p22HasReason(stale["reasons"].([]any), "evidence-stale") {
		t.Fatalf("same-HEAD dirty-byte resume = %#v, want valid stale proof with evidence-stale", stale)
	}
	repositoryState := p21Path(t, staleResume, "result", "resume", "repository").(map[string]any)
	if repositoryState["dirty"] != true || repositoryState["head_oid"] != repository["head_oid"] {
		t.Fatalf("resume repository = %#v, want current same-HEAD dirty observation", repositoryState)
	}
	evidence := p21Run(t, binary, workspace, "--project", workspace, "--json", "evidence", "get", evidenceID)
	if got := p21PathString(t, evidence, "result", "evidence", "id"); got != evidenceID {
		t.Fatalf("same-HEAD dirty staleness changed Evidence identity to %q, want %q", got, evidenceID)
	}
}

func p22HasReason(reasons []any, code string) bool {
	for _, item := range reasons {
		if item.(map[string]any)["code"] == code {
			return true
		}
	}
	return false
}

func p22GateRun(t *testing.T, binary, workspace, actor, inputPath, key string) map[string]any {
	t.Helper()
	stdout, stderr, err := closureRunCommand(workspace, binary, "--project", workspace, "--json", "gate", "run", "--actor", actor, "--input", inputPath, "--idempotency-key", key)
	if err != nil || stderr != "" {
		t.Fatalf("gate run = err:%v stdout:%q stderr:%q", err, stdout, stderr)
	}
	for _, line := range strings.FieldsFunc(stdout, func(r rune) bool { return r == '\n' }) {
		var frame map[string]any
		if err := json.Unmarshal([]byte(line), &frame); err != nil {
			t.Fatalf("decode gate frame: %v", err)
		}
		if frame["stream"] == "terminal" {
			data, ok := frame["data"].(map[string]any)
			if !ok || data["outcome"] != "succeeded" {
				t.Fatalf("gate terminal = %#v, want succeeded", data)
			}
			return data["response"].(map[string]any)
		}
	}
	t.Fatal("gate run emitted no terminal response")
	return nil
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
