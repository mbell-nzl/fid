//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

const p21CorpusDir = "testdata/p21_conformance"

var (
	p21CodeSpan = regexp.MustCompile("`([^`]+)`")
	p21Quoted   = regexp.MustCompile(`"([^"]+)"`)
	p21Kebab    = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
)

func TestPhase21ConformanceCorpusIsCompleteAndByteExact(t *testing.T) {
	t.Parallel()
	repoRoot := integrationRepositoryRoot()
	publicInterface := string(readP21File(t, filepath.Join(repoRoot, "project-docs/docs/phases/phase-21/public-interface.md")))
	expected := p21ExpectedCorpus(t, publicInterface)
	if os.Getenv("P21_WRITE_CONFORMANCE_FIXTURES") == "1" {
		p21WriteCorpus(t, expected)
	}

	for path, want := range expected {
		got, err := os.ReadFile(filepath.Join(p21CorpusDir, path))
		if err != nil {
			t.Fatalf("read committed conformance fixture %s: %v", path, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("conformance fixture %s differs from regenerated corpus\n%s", path, p21CorpusDiff(want, got))
		}
	}

	entries, err := p21CorpusEntries(p21CorpusDir)
	if err != nil {
		t.Fatalf("read committed corpus: %v", err)
	}
	if len(entries) != len(expected) {
		t.Fatalf("committed conformance fixture count = %d, want %d", len(entries), len(expected))
	}
	for _, entry := range entries {
		if _, ok := expected[entry]; !ok {
			t.Fatalf("unexpected committed conformance fixture %s", entry)
		}
	}

	manifest := p21DecodeFixture(t, filepath.Join(p21CorpusDir, "manifest.json"))
	observed := p21StringSlice(t, manifest["closed_values"], "manifest closed_values")
	closedMembers := p21ClosedMembers(publicInterface)
	missing := p21Missing(closedMembers, observed)
	if len(missing) != 0 {
		t.Fatalf("closed-value positive control missing %d member(s): %s", len(missing), strings.Join(missing, ", "))
	}
	for _, value := range observed {
		if !p21Kebab.MatchString(value) {
			t.Fatalf("closed public enum value %q is not lowercase kebab-case", value)
		}
		if strings.Contains(value, "_") {
			t.Fatalf("closed-value negative underscore control matched %q", value)
		}
	}
}

func TestPhase21RealProcessJourneyStartsOneActivatedEnvelope(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p21-corpus@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 21 corpus")
	p22Write(t, workspace, ".gitignore", ".fid/\n*.json\n*.zip\n")
	p22Git(t, workspace, "add", ".gitignore")
	p22Git(t, workspace, "commit", "-qm", "fixture")
	init := p21Run(t, binary, workspace, "--json", "init", "--prefix", "P21")
	projectID := p21PathString(t, init, "result", "project", "id")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "workflow", "define", "--name", "p21-journey", "--applies-to", "task",
		"--states", "unstarted,started,completed,cancelled",
		"--categories", "unstarted,started,completed,cancelled",
		"--transitions", "unstarted->started,started->completed", "--actor", "agent:developer", "--idempotency-key", "p21-workflow-define")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "workflow", "assign", "--name", "p21-journey", "--kind", "task", "--actor", "agent:developer", "--idempotency-key", "p21-workflow-assign")
	issue := p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Phase 21 conformance journey", "--actor", "agent:developer", "--idempotency-key", "p21-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")

	activationID := p21Activation(t, binary, workspace, "prepare", "p21-activation-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p21", "phase": 21})
	capabilities := []string{"artifact.verify", "claim.acquire", "claim.release", "claim.renew", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "issue.transition", "run.close", "run.start"}
	profileDigest := p21Component(t, binary, workspace, activationID, "authority-profile", "p21-authority", []any{
		map[string]any{"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:p21", "session_ref": nil, "capabilities": capabilities},
		map[string]any{"actor_id": "actor-harness", "actor_type": "agent", "producer": "activated-harness", "provenance_ref": "activation:p21", "session_ref": nil, "capabilities": []string{"run.observe"}},
	})
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "p21-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-unit", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	p21Activation(t, binary, workspace, "commit", "p21-activation-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	p21Activation(t, binary, workspace, "finalize", "p21-activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})

	repository := p22CorpusRepository(t, workspace, projectID)
	registerInput := p21Input(t, workspace, "gate-register", map[string]any{"definition_id": "gate-unit", "definition_version": 1, "scope": map[string]any{"issue_id": issueID}})
	registered := p21Run(t, binary, workspace, "--project", workspace, "--json", "gate", "register", "--actor", "actor-builder", "--input", registerInput, "--idempotency-key", "p21-gate-register")
	registrationID := p21PathString(t, registered, "result", "registration", "id")
	registrationVersion := p21Path(t, registered, "result", "registration", "version")

	if err := os.MkdirAll(filepath.Join(workspace, "project-docs/gates"), 0o700); err != nil {
		t.Fatal(err)
	}
	p22Write(t, workspace, "project-docs/gates/p21.json", "gate policy\n")
	source := map[string]any{"display_label": "Phase 21", "external_id": nil, "path": "project-docs/gates/p21.json", "section": nil, "source_type": "project-file"}
	envelopeInput := p21Input(t, workspace, "envelope-create", map[string]any{
		"outcome": "Exercise the one Issue walking skeleton", "non_goals": []string{}, "risk_result": map[string]any{"source": source, "result_digest": "sha256:" + strings.Repeat("1", 64)}, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{},
		"authority_grants": []any{
			map[string]any{"actions": capabilities[1:], "actor_selector": map[string]any{"actor_id": "actor-builder", "required_capabilities": capabilities[1:]}, "scope": map[string]any{"issue_id": issueID}},
			map[string]any{"actions": []string{"run.observe"}, "actor_selector": map[string]any{"actor_id": "actor-harness", "required_capabilities": []string{"run.observe"}}, "scope": map[string]any{"issue_id": issueID}},
		},
		"obligations": []any{map[string]any{"id": "OBL-1", "type": "gate-evidence", "applicability": "required", "guarded_actions": []string{"envelope.complete", "issue.transition-terminal-success"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": map[string]any{"id": registrationID, "version": registrationVersion}, "accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"}, "freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": source}},
		"budgets":     []any{map[string]any{"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 1, "warning_threshold": nil, "exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil}},
	})
	created := p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", envelopeInput, "--idempotency-key", "p21-envelope-create")
	envelopeID := p21PathString(t, created, "result", "envelope", "id")
	activateInput := p21Input(t, workspace, "envelope-activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder", "--input", activateInput, "--idempotency-key", "p21-envelope-activate")

	claimRequest := map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "scope": map[string]any{"paths": []string{"cmd/fid/integration"}}, "ttl_seconds": 60}
	claimInput := p21Input(t, workspace, "claim-acquire", claimRequest)
	otherClaimInput := p21Input(t, workspace, "claim-acquire-race", claimRequest)
	beforeClaimRace := p21JournalEventCount(t, workspace)
	claim := p21ConcurrentClaim(t, binary, workspace, claimInput, otherClaimInput)
	if afterClaimRace := p21JournalEventCount(t, workspace); afterClaimRace != beforeClaimRace+1 {
		t.Fatalf("two-process Claim race appended %d Events, want one", afterClaimRace-beforeClaimRace)
	}
	claimID := p21PathString(t, claim, "result", "claim", "id")
	runInput := p21Input(t, workspace, "run-start", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": claimID, "claim_version": 1, "executor": map[string]any{"id": "executor-builder", "source": "activated-harness"}, "obligation_id": nil, "expected_repository": p24ExpectedRepository(repository)})
	run := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", runInput, "--idempotency-key", "p21-run-start")
	runID := p21PathString(t, run, "result", "run", "id")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "transition", issueID, "--to", "started", "--actor", "agent:actor-builder", "--expected-revision", "1", "--idempotency-key", "p21-transition-nonterminal")

	gate := p21GateRun(t, binary, workspace, map[string]any{"registration_id": registrationID, "registration_version": registrationVersion, "obligation_id": "OBL-1", "run_id": runID, "expected_repository": map[string]any{"head_oid": repository["head_oid"], "dirty": repository["dirty"]}})
	receiptID := p21PathString(t, gate, "result", "receipt", "id")
	p21AssertBudgetBoundaryAndLaterGateDenial(t, binary, workspace, registrationID, registrationVersion, runID, repository)
	linkInput := p21Input(t, workspace, "evidence-link", map[string]any{"obligation_id": "OBL-1", "evidence_type": "receipt", "receipt_id": receiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}})
	p21Run(t, binary, workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", linkInput, "--idempotency-key", "p21-evidence-link")
	evaluateInput := p21Input(t, workspace, "evidence-evaluate", map[string]any{"obligation_id": "OBL-1"})
	evaluation := p21Run(t, binary, workspace, "--project", workspace, "--json", "evidence", "evaluate", "--input", evaluateInput)
	if status := p21PathString(t, evaluation, "result", "evaluation", "status"); status != "satisfied" {
		t.Fatalf("Gate Evidence evaluation status = %q, want satisfied", status)
	}
	p21AssertUndeclaredTransitionRejected(t, binary, workspace, issueID, "unstarted", 2, "p21-transition-undeclared-before-rebuild")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "transition", issueID, "--to", "completed", "--actor", "agent:actor-builder", "--expected-revision", "2", "--idempotency-key", "p21-transition-terminal")
	currentRun := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "get", runID)
	observeInput := p21Input(t, workspace, "run-observe", map[string]any{"run_id": runID, "expected_version": p21Path(t, currentRun, "result", "run", "version"), "liveness": map[string]any{"executor_id": "executor-builder", "state": "terminated", "source_ref": "p21:terminal", "observed_at": "2026-09-07T00:00:00Z"}})
	observed := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-harness", "--input", observeInput, "--idempotency-key", "p21-run-observe")
	closeInput := p21Input(t, workspace, "run-close", map[string]any{"run_id": runID, "expected_version": p21Path(t, observed, "result", "run", "version"), "observation_id": p21PathString(t, observed, "result", "observation", "id"), "expected_repository": p24ExpectedRepository(repository)})
	p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "close", "--actor", "actor-builder", "--input", closeInput, "--idempotency-key", "p21-run-close")
	releaseInput := p21Input(t, workspace, "claim-release", map[string]any{"claim_id": claimID, "expected_version": 1, "reason": "completed"})
	p21Run(t, binary, workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", releaseInput, "--idempotency-key", "p21-claim-release")
	completeInput := p21Input(t, workspace, "envelope-complete", map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "complete", "--actor", "actor-builder", "--input", completeInput, "--idempotency-key", "p21-envelope-complete")

	p21AssertJournalTypes(t, workspace, []string{
		"project.initialized", "workflow.defined", "workflow.assigned", "issue.created",
		"activation.prepared", "activation.component-staged", "activation.component-staged",
		"activation.resource-observed", "activation.resource-observed", "activation.resource-observed", "activation.resource-observed", "activation.resource-observed",
		"activation.committed", "activation.finalized", "gate.registered", "envelope.created",
		"envelope.activated", "claim.acquired", "run.started", "issue.transitioned",
		"gate.execution-started", "gate.receipt-recorded", "run.execution-blocked", "evidence.linked", "issue.transitioned",
		"run.observed", "run.closed", "claim.released", "envelope.completed",
	})
	p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "get", envelopeID)
	p21Run(t, binary, workspace, "--project", workspace, "--json", "receipt", "get", receiptID)
	p21Run(t, binary, workspace, "--project", workspace, "--json", "receipt", "list", "--run", runID, "--obligation", "OBL-1", "--outcome", "succeeded")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "resume", issueID)
	beforeArchive := filepath.Join(workspace, "p21-before-rebuild.zip")
	beforeExport := p21Run(t, binary, workspace, "--project", workspace, "--json", "export", "create", beforeArchive)
	p21Run(t, binary, workspace, "--project", workspace, "--json", "export", "validate", beforeArchive)

	if err := os.Remove(filepath.Join(workspace, ".fid", "fid.sqlite")); err != nil {
		t.Fatalf("delete derived projection: %v", err)
	}
	p21Run(t, binary, workspace, "--project", workspace, "--json", "rebuild")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "capabilities")
	completed := p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "get", envelopeID)
	if status := p21PathString(t, completed, "result", "envelope", "status"); status != "completed" {
		t.Fatalf("rebuilt Envelope status = %q, want completed", status)
	}
	p21Run(t, binary, workspace, "--project", workspace, "--json", "receipt", "get", receiptID)
	p21Run(t, binary, workspace, "--project", workspace, "--json", "receipt", "list", "--run", runID, "--obligation", "OBL-1", "--outcome", "succeeded")
	postRebuildEvaluation := p21Run(t, binary, workspace, "--project", workspace, "--json", "evidence", "evaluate", "--input", evaluateInput)
	if status := p21PathString(t, postRebuildEvaluation, "result", "evaluation", "status"); status != "satisfied" {
		t.Fatalf("rebuilt Gate Evidence evaluation status = %q, want satisfied", status)
	}
	p21Run(t, binary, workspace, "--project", workspace, "--json", "resume", issueID)
	afterArchive := filepath.Join(workspace, "p21-after-rebuild.zip")
	afterExport := p21Run(t, binary, workspace, "--project", workspace, "--json", "export", "create", afterArchive)
	p21Run(t, binary, workspace, "--project", workspace, "--json", "export", "validate", afterArchive)
	if beforeDigest, afterDigest := p21PathString(t, beforeExport, "result", "export", "archive_digest"), p21PathString(t, afterExport, "result", "export", "archive_digest"); beforeDigest != afterDigest {
		t.Fatalf("export digest after rebuild = %q, want %q", afterDigest, beforeDigest)
	}
	p21AssertUndeclaredTransitionRejected(t, binary, workspace, issueID, "started", 3, "p21-transition-undeclared-after-rebuild")
	p21AssertJournalTypes(t, workspace, []string{
		"project.initialized", "workflow.defined", "workflow.assigned", "issue.created",
		"activation.prepared", "activation.component-staged", "activation.component-staged",
		"activation.resource-observed", "activation.resource-observed", "activation.resource-observed", "activation.resource-observed", "activation.resource-observed",
		"activation.committed", "activation.finalized", "gate.registered", "envelope.created",
		"envelope.activated", "claim.acquired", "run.started", "issue.transitioned",
		"gate.execution-started", "gate.receipt-recorded", "run.execution-blocked", "evidence.linked", "issue.transitioned",
		"run.observed", "run.closed", "claim.released", "envelope.completed",
	})
}

func p21AssertUndeclaredTransitionRejected(t *testing.T, binary, workspace, issueID, target string, revision int, key string) {
	t.Helper()
	before := p21JournalEventCount(t, workspace)
	stdout, stderr, err := closureRunCommand(workspace, binary, "--project", workspace, "--json", "issue", "transition", issueID, "--to", target, "--actor", "agent:actor-builder", "--expected-revision", fmt.Sprint(revision), "--idempotency-key", key)
	if got := processExitCode(t, err); got != 4 {
		t.Fatalf("undeclared transition %q exit code = %d, want 4; stdout=%q stderr=%q", target, got, stdout, stderr)
	}
	if stdout != "" || !strings.Contains(stderr, `"code":"transition-not-permitted"`) || !strings.Contains(stderr, "workflow-transition-invalid") {
		t.Fatalf("undeclared transition %q = stdout:%q stderr:%q, want transition-not-permitted/workflow-transition-invalid", target, stdout, stderr)
	}
	if after := p21JournalEventCount(t, workspace); after != before {
		t.Fatalf("undeclared transition %q appended Event count %d, want %d", target, after, before)
	}
}

func p21GateRun(t *testing.T, binary, workspace string, input map[string]any) map[string]any {
	t.Helper()
	path := p21Input(t, workspace, "gate-run", input)
	stdout, stderr, err := closureRunCommand(workspace, binary, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", path, "--idempotency-key", "p21-gate-run")
	if err != nil || stderr != "" {
		t.Fatalf("gate run = err:%v stdout:%q stderr:%q", err, stdout, stderr)
	}
	lines := strings.FieldsFunc(stdout, func(r rune) bool { return r == '\n' })
	if len(lines) == 0 {
		t.Fatal("gate run emitted no live frames")
	}
	var terminal map[string]any
	for _, line := range lines {
		var frame map[string]any
		if err := json.Unmarshal([]byte(line), &frame); err != nil {
			t.Fatalf("decode live frame: %v; line=%q", err, line)
		}
		if frame["stream"] == "terminal" {
			terminal, _ = frame["data"].(map[string]any)
		}
	}
	if terminal == nil || terminal["outcome"] != "succeeded" || terminal["replayed"] != false {
		t.Fatalf("gate terminal = %#v, want one non-replayed succeeded terminal frame", terminal)
	}
	response, ok := terminal["response"].(map[string]any)
	if !ok {
		t.Fatalf("gate terminal response = %#v, want success envelope", terminal["response"])
	}
	return response
}

func p21AssertBudgetBoundaryAndLaterGateDenial(t *testing.T, binary, workspace, registrationID string, registrationVersion any, runID string, repository map[string]any) {
	t.Helper()
	before := p21JournalBytes(t, workspace)
	if strings.Count(string(before), `"table":"execution_attempts"`) != 1 || strings.Count(string(before), `"table":"budget_samples"`) != 1 || strings.Count(string(before), `"table":"stops"`) != 1 || strings.Count(string(before), `"table":"attention_items"`) != 1 || strings.Count(string(before), `"table":"attention_stop_links"`) != 1 {
		t.Fatalf("boundary Event effects = %q, want one Attempt, budget sample, Stop, Attention, and link", before)
	}
	input := p21Input(t, workspace, "gate-run-after-exhaustion", map[string]any{"registration_id": registrationID, "registration_version": registrationVersion, "obligation_id": "OBL-1", "run_id": runID, "expected_repository": map[string]any{"head_oid": repository["head_oid"], "dirty": repository["dirty"]}})
	stdout, stderr, runErr := closureRunCommand(workspace, binary, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", input, "--idempotency-key", "p21-gate-run-after-exhaustion")
	if processExitCode(t, runErr) != 4 || stderr != "" || !strings.Contains(stdout, `"stream":"terminal"`) || !strings.Contains(stdout, `"outcome":"fid-error"`) || !strings.Contains(stdout, `"code":"budget-exhausted"`) || strings.Contains(stdout, `"attempt_id":"`) {
		t.Fatalf("later gate denial = exit:%d stdout:%q stderr:%q, want committed pre-start denial terminal", processExitCode(t, runErr), stdout, stderr)
	}
	after := p21JournalBytes(t, workspace)
	if !strings.Contains(string(after), `"event_type":"run.execution-blocked"`) || strings.Count(string(after), `"table":"execution_attempts"`) != 1 || strings.Count(string(after), `"table":"stops"`) != 2 {
		t.Fatalf("later denial journal = %q, want execution-blocked with no second Attempt and one additional Stop", after)
	}
}

func p21Activation(t *testing.T, binary, workspace, operation, key string, input map[string]any) string {
	t.Helper()
	if operation == "prepare" {
		input["expected_prior_activation_id"] = nil
		input["expected_prior_version"] = nil
		input["resources"] = p21ActivationResources()
	}
	if operation == "commit" || operation == "finalize" {
		activationID, _ := input["activation_id"].(string)
		activation := p21ActivationStatus(t, workspace, activationID)
		if operation == "commit" {
			p21ObserveNewActivationResources(t, binary, workspace, activation)
			activation = p21ActivationStatus(t, workspace, activationID)
		}
		input["expected_version"] = activation.Version
	}
	path := p21Input(t, workspace, "activation-"+operation, input)
	result := p21Run(t, binary, workspace, "--project", workspace, "--json", "activation", operation, "--input", path, "--idempotency-key", key)
	return p21PathString(t, result, "result", "activation", "id")
}

func p21Component(t *testing.T, binary, workspace, activationID, componentType, componentID string, content []any) string {
	t.Helper()
	activation := p21ActivationStatus(t, workspace, activationID)
	path := p21Input(t, workspace, componentID, map[string]any{"activation_id": activationID, "component_type": componentType, "component_id": componentID, "component_version": 1, "expected_version": activation.Version, "expected_prior_digest": nil, "content": content})
	result := p21Run(t, binary, workspace, "--project", workspace, "--json", "activation", "component", "--input", path, "--idempotency-key", "p21-"+componentID)
	return p21PathString(t, result, "result", "component", "digest")
}

func p21ActivationResources() []any {
	return []any{
		map[string]any{"component": "backend", "old_identity": nil, "new_identity": "fid"},
		map[string]any{"component": "kit", "old_identity": nil, "new_identity": "p21-kit"},
		map[string]any{"component": "overlay", "old_identity": nil, "new_identity": "p21-overlay"},
		map[string]any{"component": "project", "old_identity": nil, "new_identity": "p21-project"},
		map[string]any{"component": "provider", "old_identity": nil, "new_identity": "p21-provider"},
	}
}

func p21ActivationStatus(t *testing.T, workspace, activationID string) workflow.Activation {
	t.Helper()
	activation, err := workflow.ActivationStatus(workspace, filepath.Join(integrationRepositoryRoot(), "migrations"), activationID)
	if err != nil {
		t.Fatalf("load Activation %q: %v", activationID, err)
	}
	return activation
}

func p21ObserveNewActivationResources(t *testing.T, binary, workspace string, activation workflow.Activation) {
	t.Helper()
	for _, resource := range activation.Resources {
		input := p21Input(t, workspace, "activation-resource-"+resource.Component, map[string]any{"activation_id": activation.ID, "component": resource.Component, "expected_activation_version": activation.Version, "expected_version": resource.Version, "observed_identity": resource.NewIdentity, "state": "new", "evidence_ref": "integration:" + resource.Component, "error_code": nil})
		result := p21Run(t, binary, workspace, "--project", workspace, "--json", "activation", "resource", "--input", input, "--idempotency-key", "p21-activation-resource-"+activation.ID+"-"+resource.Component)
		activation = workflow.Activation{ID: activation.ID, Version: int(p21Path(t, result, "result", "activation", "version").(float64))}
	}
}

func p21Input(t *testing.T, workspace, name string, value map[string]any) string {
	t.Helper()
	contents, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal %s input: %v", name, err)
	}
	path := filepath.Join(workspace, name+".json")
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatalf("write %s input: %v", name, err)
	}
	return path
}

func p21ConcurrentClaim(t *testing.T, binary, workspace, firstInput, secondInput string) map[string]any {
	t.Helper()
	type process struct {
		stdout bytes.Buffer
		stderr bytes.Buffer
		cmd    *exec.Cmd
	}
	arguments := [][]string{
		{"--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder", "--input", firstInput, "--idempotency-key", "p21-claim-race-a"},
		{"--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder", "--input", secondInput, "--idempotency-key", "p21-claim-race-b"},
	}
	processes := make([]process, len(arguments))
	for index, args := range arguments {
		processes[index].cmd = exec.Command(binary, args...)
		processes[index].cmd.Dir = workspace
		processes[index].cmd.Stdout = &processes[index].stdout
		processes[index].cmd.Stderr = &processes[index].stderr
		if err := processes[index].cmd.Start(); err != nil {
			t.Fatalf("start Claim contender %d: %v", index, err)
		}
	}
	winners := []map[string]any{}
	for index := range processes {
		err := processes[index].cmd.Wait()
		if err == nil {
			var response map[string]any
			if decodeErr := json.Unmarshal(processes[index].stdout.Bytes(), &response); decodeErr != nil {
				t.Fatalf("decode winning Claim contender %d: %v; stdout=%q", index, decodeErr, processes[index].stdout.String())
			}
			winners = append(winners, response)
			continue
		}
		if code := processExitCode(t, err); code != 4 || processes[index].stdout.Len() != 0 || !strings.Contains(processes[index].stderr.String(), `"code":"claim-conflict"`) {
			t.Fatalf("losing Claim contender %d = exit %d stdout=%q stderr=%q, want typed conflict/no stdout", index, code, processes[index].stdout.String(), processes[index].stderr.String())
		}
	}
	if len(winners) != 1 {
		t.Fatalf("two real Claim processes produced %d winners, want one", len(winners))
	}
	return winners[0]
}

func p21Run(t *testing.T, binary, workspace string, args ...string) map[string]any {
	t.Helper()
	stdout, stderr, err := closureRunCommand(workspace, binary, args...)
	if err != nil || stderr != "" {
		t.Fatalf("fid %s = err:%v stdout:%q stderr:%q", strings.Join(args, " "), err, stdout, stderr)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil {
		t.Fatalf("decode fid %s response: %v; stdout=%q", strings.Join(args, " "), err, stdout)
	}
	return envelope
}

func p21PathString(t *testing.T, value map[string]any, path ...string) string {
	t.Helper()
	item := p21Path(t, value, path...)
	text, ok := item.(string)
	if !ok || text == "" {
		t.Fatalf("path %s = %#v, want non-empty string", strings.Join(path, "."), item)
	}
	return text
}

func p21Path(t *testing.T, value map[string]any, path ...string) any {
	t.Helper()
	var item any = value
	for _, key := range path {
		object, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("path %s traverses %#v, want object", strings.Join(path, "."), item)
		}
		item = object[key]
	}
	return item
}

func p21JournalEventCount(t *testing.T, workspace string) int {
	t.Helper()
	return len(p21JournalEventTypes(t, workspace))
}

func p21AssertJournalTypes(t *testing.T, workspace string, want []string) {
	t.Helper()
	got := p21JournalEventTypes(t, workspace)
	if len(got) != len(want) {
		t.Fatalf("journal Event count = %d, want %d; types=%v", len(got), len(want), got)
	}
	for index, eventType := range want {
		if got[index] != eventType {
			t.Fatalf("journal Event %d = %q, want %q; types=%v", index+1, got[index], eventType, got)
		}
	}
}

func p21JournalEventTypes(t *testing.T, workspace string) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(workspace, ".fid", "journal"))
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}
	var types []string
	for _, entry := range entries {
		contents := readP21File(t, filepath.Join(workspace, ".fid", "journal", entry.Name()))
		for _, line := range bytes.FieldsFunc(contents, func(r rune) bool { return r == '\n' }) {
			var event struct {
				Type string `json:"event_type"`
			}
			if err := json.Unmarshal(line, &event); err != nil {
				t.Fatalf("decode journal Event: %v", err)
			}
			if event.Type == "" {
				t.Fatal("journal Event has empty type")
			}
			types = append(types, event.Type)
		}
	}
	return types
}

func p21JournalBytes(t *testing.T, workspace string) []byte {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(workspace, ".fid", "journal"))
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}
	var journal []byte
	for _, entry := range entries {
		journal = append(journal, readP21File(t, filepath.Join(workspace, ".fid", "journal", entry.Name()))...)
	}
	return journal
}

func p21ExpectedCorpus(t *testing.T, publicInterface string) map[string][]byte {
	t.Helper()
	operations := p21OperationRows(t, publicInterface)
	for _, operation := range []p21OperationFixture{
		{Operation: "issue.ac.verify", Command: "fid issue ac verify", Kind: "fail-closed"},
		{Operation: "concern.raise", Command: "fid concern raise", Kind: "fail-closed"},
		{Operation: "concern.address", Command: "fid concern address", Kind: "fail-closed"},
		{Operation: "maintenance.vacuum", Command: "fid maintenance vacuum", Kind: "fail-closed"},
	} {
		operations = append(operations, operation)
	}
	if len(operations) == 0 {
		t.Fatal("operation extraction produced no rows")
	}
	failures := p21FailureCodes(t, publicInterface)
	if len(failures) == 0 {
		t.Fatal("failure extraction produced no codes")
	}
	live := p21LiveVariants()

	expected := map[string][]byte{}
	expected["manifest.json"] = p21JSON(t, map[string]any{
		"closed_values": p21ClosedMembers(publicInterface),
		"failures":      failures,
		"live":          live,
		"operations":    operations,
	})
	for _, operation := range operations {
		expected["operations/"+operation.Operation+".json"] = p21JSON(t, operation)
	}
	for _, code := range failures {
		expected["failures/"+code+".json"] = p21JSON(t, map[string]string{"failure_code": code})
	}
	for _, variant := range live {
		encoded, err := json.Marshal(p21LiveFixture(variant))
		if err != nil {
			t.Fatalf("marshal live fixture %s: %v", variant, err)
		}
		expected["live/"+variant+".jsonl"] = append(encoded, '\n')
	}
	return expected
}

type p21OperationFixture struct {
	Operation string `json:"operation"`
	Command   string `json:"command"`
	Kind      string `json:"kind"`
}

func p21OperationRows(t *testing.T, source string) []p21OperationFixture {
	t.Helper()
	section := p21Section(t, source, "### 4.1 Retained standalone and inventory operations", "### 4.2 Operational-kernel operations")
	section += p21Section(t, source, "### 4.2 Operational-kernel operations", "## 5. Fail-Closed And Absent Leaves")
	seen := map[string]bool{}
	var operations []p21OperationFixture
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 3 {
			t.Fatalf("malformed public operation row: %q", line)
		}
		values := p21CodeSpan.FindAllStringSubmatch(cells[1], -1)
		if len(values) != 2 {
			t.Fatalf("operation row must contain operation and command: %q", line)
		}
		operation, command := values[0][1], values[1][1]
		if seen[operation] {
			t.Fatalf("duplicate public operation %q", operation)
		}
		seen[operation] = true
		operations = append(operations, p21OperationFixture{Operation: operation, Command: command, Kind: "accepted"})
	}
	return operations
}

func p21FailureCodes(t *testing.T, source string) []string {
	t.Helper()
	section := p21Section(t, source, "## 7. Complete Failure And Exit Contract", "## 8. Exact `fid.live.v1` Frames")
	seen := map[string]bool{}
	var codes []string
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		cell := strings.Split(line, "|")[1]
		for _, match := range p21CodeSpan.FindAllStringSubmatch(cell, -1) {
			code := match[1]
			if seen[code] {
				continue
			}
			seen[code] = true
			codes = append(codes, code)
		}
	}
	return codes
}

func p21LiveVariants() []string {
	return []string{
		"stdout", "stderr", "control-output-truncated", "control-redaction-applied",
		"control-artifact-rejected", "terminal-failed", "terminal-fid-error",
		"terminal-signaled", "terminal-spawn-failed", "terminal-succeeded",
		"terminal-timed-out", "replay-succeeded", "replay-gate-unsuccessful",
	}
}

func p21LiveFixture(variant string) map[string]any {
	stream := "terminal"
	data := map[string]any{"outcome": "succeeded", "replayed": false}
	switch {
	case variant == "stdout" || variant == "stderr":
		stream = variant
		data = map[string]any{"encoding": "utf-8", "text": "conformance"}
	case strings.HasPrefix(variant, "control-"):
		stream = "control"
		data = map[string]any{"kind": strings.TrimPrefix(variant, "control-")}
	case strings.HasPrefix(variant, "terminal-"):
		data["outcome"] = strings.TrimPrefix(variant, "terminal-")
	case variant == "replay-succeeded":
		data["replayed"] = true
	case variant == "replay-gate-unsuccessful":
		data["outcome"] = "failed"
		data["replayed"] = true
	}
	return map[string]any{"frame_schema": "fid.live.v1", "sequence": 1, "stream": stream, "data": data}
}

func p21ClosedMembers(source string) []string {
	sections := []string{
		p21Section(nil, source, "## 3. Closed Result Types", "## 4. Exact Accepted Leaves"),
		p21Section(nil, source, "## 7. Complete Failure And Exit Contract", "## 8. Exact `fid.live.v1` Frames"),
		p21Section(nil, source, "## 8. Exact `fid.live.v1` Frames", ""),
	}
	seen := map[string]bool{}
	var members []string
	for _, section := range sections {
		for _, match := range p21Quoted.FindAllStringSubmatch(section, -1) {
			value := match[1]
			if !p21Kebab.MatchString(value) || seen[value] {
				continue
			}
			seen[value] = true
			members = append(members, value)
		}
	}
	sort.Strings(members)
	return members
}

func p21Section(t *testing.T, source, start, end string) string {
	if t != nil {
		t.Helper()
	}
	startAt := strings.Index(source, start)
	if startAt < 0 {
		if t != nil {
			t.Fatalf("missing public-interface section %q", start)
		}
		return ""
	}
	section := source[startAt:]
	if end == "" {
		return section
	}
	endAt := strings.Index(section, end)
	if endAt < 0 {
		if t != nil {
			t.Fatalf("missing public-interface section boundary %q", end)
		}
		return ""
	}
	return section[:endAt]
}

func p21JSON(t *testing.T, value any) []byte {
	t.Helper()
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	return append(encoded, '\n')
}

func p21CorpusEntries(root string) ([]string, error) {
	var entries []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		entries = append(entries, relative)
		return nil
	})
	sort.Strings(entries)
	return entries, err
}

func p21WriteCorpus(t *testing.T, expected map[string][]byte) {
	t.Helper()
	for path, contents := range expected {
		fixturePath := filepath.Join(p21CorpusDir, path)
		if err := os.MkdirAll(filepath.Dir(fixturePath), 0o755); err != nil {
			t.Fatalf("create conformance fixture directory for %s: %v", path, err)
		}
		if err := os.WriteFile(fixturePath, contents, 0o644); err != nil {
			t.Fatalf("write conformance fixture %s: %v", path, err)
		}
	}
}

func p21DecodeFixture(t *testing.T, path string) map[string]any {
	t.Helper()
	contents := readP21File(t, path)
	var fixture map[string]any
	if err := json.Unmarshal(contents, &fixture); err != nil {
		t.Fatalf("decode fixture %s: %v", path, err)
	}
	return fixture
}

func p21StringSlice(t *testing.T, value any, label string) []string {
	t.Helper()
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("%s = %#v, want array", label, value)
	}
	stringsOnly := make([]string, 0, len(items))
	for _, item := range items {
		stringValue, ok := item.(string)
		if !ok {
			t.Fatalf("%s contains %#v, want string", label, item)
		}
		stringsOnly = append(stringsOnly, stringValue)
	}
	return stringsOnly
}

func p21Missing(expected, observed []string) []string {
	present := map[string]bool{}
	for _, value := range observed {
		present[value] = true
	}
	var missing []string
	for _, value := range expected {
		if !present[value] {
			missing = append(missing, value)
		}
	}
	return missing
}

func p21CorpusDiff(want, got []byte) string {
	return fmt.Sprintf("want:\n%s\ngot:\n%s", want, got)
}

func readP21File(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return contents
}
