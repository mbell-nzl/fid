//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/project"
)

const p22CorpusDir = "testdata/p22_conformance"

type p22ConformanceCase struct {
	Name         string   `json:"name"`
	Operation    string   `json:"operation"`
	Kind         string   `json:"kind"`
	ResultFields []string `json:"result_fields"`
	Event        bool     `json:"event"`
	FailureCode  string   `json:"failure_code"`
	ExitCode     int      `json:"exit_code"`
}

type p22ConformanceResponse struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type p22RecoveryAdapter interface {
	Name() string
	Execute(t *testing.T, fixture p22ConformanceCase) p22ConformanceResponse
}

func TestP22DualAdapterRecoveryConformance(t *testing.T) {
	fixtures := p22RecoveryFixtures(t)
	for _, newAdapter := range []func(*testing.T) p22RecoveryAdapter{
		func(_ *testing.T) p22RecoveryAdapter { return p22FakeRecoveryAdapter{} },
		func(t *testing.T) p22RecoveryAdapter { return p22NewFIDRecoveryAdapter(t) },
	} {
		adapter := newAdapter(t)
		t.Run(adapter.Name(), func(t *testing.T) {
			for _, fixture := range fixtures {
				t.Run(fixture.Name, func(t *testing.T) {
					assertP22ConformanceResponse(t, fixture, adapter.Execute(t, fixture))
				})
			}
		})
	}
}

func p22RecoveryFixtures(t *testing.T) []p22ConformanceCase {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(p22CorpusDir, "manifest.json"))
	if err != nil {
		t.Fatalf("read Phase 22 recovery conformance manifest: %v", err)
	}
	var manifest struct {
		Cases []string `json:"cases"`
	}
	if err := json.Unmarshal(contents, &manifest); err != nil {
		t.Fatalf("decode Phase 22 recovery conformance manifest: %v", err)
	}
	if len(manifest.Cases) != 10 {
		t.Fatalf("Phase 22 recovery conformance fixture count = %d, want 10", len(manifest.Cases))
	}
	fixtures := make([]p22ConformanceCase, 0, len(manifest.Cases))
	for _, path := range manifest.Cases {
		contents, err := os.ReadFile(filepath.Join(p22CorpusDir, path))
		if err != nil {
			t.Fatalf("read Phase 22 recovery fixture %s: %v", path, err)
		}
		var fixture p22ConformanceCase
		if err := json.Unmarshal(contents, &fixture); err != nil {
			t.Fatalf("decode Phase 22 recovery fixture %s: %v", path, err)
		}
		if fixture.Name == "" || fixture.Operation == "" || fixture.Kind == "" {
			t.Fatalf("Phase 22 recovery fixture %s is incomplete: %#v", path, fixture)
		}
		fixtures = append(fixtures, fixture)
	}
	return fixtures
}

func assertP22ConformanceResponse(t *testing.T, fixture p22ConformanceCase, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != fixture.ExitCode {
		t.Fatalf("%s %s exit = %d, want %d; stdout=%q stderr=%q", fixture.Operation, fixture.Kind, response.ExitCode, fixture.ExitCode, response.Stdout, response.Stderr)
	}
	if fixture.Kind == "accepted" {
		if response.Stderr != "" {
			t.Fatalf("%s accepted stderr = %q, want empty", fixture.Operation, response.Stderr)
		}
		var envelope map[string]any
		if err := json.Unmarshal([]byte(response.Stdout), &envelope); err != nil {
			t.Fatalf("decode %s accepted stdout: %v; stdout=%q", fixture.Operation, err, response.Stdout)
		}
		if envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != true || envelope["operation"] != fixture.Operation {
			t.Fatalf("%s accepted envelope = %#v", fixture.Operation, envelope)
		}
		result, ok := envelope["result"].(map[string]any)
		if !ok {
			t.Fatalf("%s accepted result = %#v, want object", fixture.Operation, envelope["result"])
		}
		for _, field := range fixture.ResultFields {
			if _, ok := result[field]; !ok {
				t.Fatalf("%s accepted result missing %q: %#v", fixture.Operation, field, result)
			}
		}
		if fixture.Event {
			if _, ok := envelope["event"].(map[string]any); !ok {
				t.Fatalf("%s accepted Event = %#v, want Event object", fixture.Operation, envelope["event"])
			}
		} else if envelope["event"] != nil {
			t.Fatalf("%s accepted Event = %#v, want null", fixture.Operation, envelope["event"])
		}
		return
	}
	if response.Stdout != "" {
		t.Fatalf("%s unknown-field stdout = %q, want empty", fixture.Operation, response.Stdout)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatalf("decode %s unknown-field stderr: %v; stderr=%q", fixture.Operation, err, response.Stderr)
	}
	errorValue, ok := envelope["error"].(map[string]any)
	if !ok || envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != false || envelope["operation"] != fixture.Operation || errorValue["code"] != fixture.FailureCode {
		t.Fatalf("%s unknown-field envelope = %#v", fixture.Operation, envelope)
	}
}

type p22FakeRecoveryAdapter struct{}

func (p22FakeRecoveryAdapter) Name() string { return "fake" }

func (p22FakeRecoveryAdapter) Execute(t *testing.T, fixture p22ConformanceCase) p22ConformanceResponse {
	t.Helper()
	if fixture.Kind == "accepted" {
		result := map[string]any{}
		for _, field := range fixture.ResultFields {
			result[field] = field + "-fake"
		}
		event := any(nil)
		if fixture.Event {
			event = map[string]any{"id": "EVT-fake", "occurred_at": "2026-09-06T00:00:00Z"}
		}
		return p22FixtureResponse(t, map[string]any{"contract_version": "fid-cli.v1", "ok": true, "operation": fixture.Operation, "result": result, "event": event}, true, fixture.ExitCode)
	}
	return p22FixtureResponse(t, map[string]any{"contract_version": "fid-cli.v1", "ok": false, "operation": fixture.Operation, "error": map[string]any{"code": fixture.FailureCode}}, false, fixture.ExitCode)
}

func p22FixtureResponse(t *testing.T, envelope map[string]any, stdout bool, exitCode int) p22ConformanceResponse {
	t.Helper()
	encoded, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("marshal fake recovery response: %v", err)
	}
	value := string(encoded) + "\n"
	if stdout {
		return p22ConformanceResponse{Stdout: value, ExitCode: exitCode}
	}
	return p22ConformanceResponse{Stderr: value, ExitCode: exitCode}
}

type p22FIDRecoveryAdapter struct {
	binary       string
	workspace    string
	projectID    string
	actor        string
	repository   map[string]any
	claimID      string
	claimVersion int
	runID        string
	runVersion   int
	observation  string
	digest       string
}

func (adapter *p22FIDRecoveryAdapter) Name() string { return "fid" }

func p22NewFIDRecoveryAdapter(t *testing.T) *p22FIDRecoveryAdapter {
	t.Helper()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p22-corpus@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 22 corpus")
	p22Write(t, workspace, "base.txt", "base\n")
	p22Write(t, workspace, ".gitignore", "*.json\n")
	p22Git(t, workspace, "add", "base.txt", ".gitignore")
	p22Git(t, workspace, "commit", "-qm", "base")
	init := p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P22")
	projectID := p21PathString(t, init, "result", "project", "id")
	issue := p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Phase 22 conformance", "--actor", "agent:developer", "--idempotency-key", "p22-corpus-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")
	actor := "actor-recovery"
	capabilities := []string{"artifact.verify", "claim.acquire", "claim.expire", "envelope.activate", "envelope.create", "gate.register", "run.checkpoint", "run.observe", "run.recover", "run.start"}
	activationID := p21Activation(t, binary, workspace, "prepare", "p22-corpus-activation-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p22-corpus", "phase": 21})
	profileDigest := p21Component(t, binary, workspace, activationID, "authority-profile", "p22-corpus-authority", []any{
		map[string]any{"actor_id": "actor-denied", "actor_type": "agent", "producer": "activated-harness", "provenance_ref": "activation:p22-corpus", "session_ref": nil, "capabilities": capabilities},
		map[string]any{"actor_id": actor, "actor_type": "agent", "producer": "activated-harness", "provenance_ref": "activation:p22-corpus", "session_ref": nil, "capabilities": capabilities},
	})
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "p22-corpus-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-p22-corpus", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	digests := []string{gateDigest, profileDigest}
	sort.Strings(digests)
	p21Activation(t, binary, workspace, "commit", "p22-corpus-activation-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	p21Activation(t, binary, workspace, "finalize", "p22-corpus-activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})
	repository := p22CorpusRepository(t, workspace, projectID)
	registered := p21Run(t, binary, workspace, "--project", workspace, "--json", "gate", "register", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-gate-register", map[string]any{"definition_id": "gate-p22-corpus", "definition_version": 1, "worktree_id": repository["worktree_id"], "scope": map[string]any{"issue_id": issueID}, "repository": repository}), "--idempotency-key", "p22-corpus-gate-register")
	registrationID := p21PathString(t, registered, "result", "registration", "id")
	registrationVersion := p21Path(t, registered, "result", "registration", "version")
	source := map[string]any{"digest": "sha256:" + p22Zeros(64), "display_label": "Phase 22 corpus", "external_id": nil, "path": "project-docs/p22-corpus.md", "repository": repository, "section": nil, "source_type": "project-file"}
	envelopeInput := p21Input(t, workspace, "p22-corpus-envelope", map[string]any{"outcome": "Recover one interrupted Run", "non_goals": []string{}, "risk_result": map[string]any{"source": source, "result_digest": "sha256:" + p22Ones(64)}, "repository_scope": repository, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{}, "authority_grants": []any{map[string]any{"actions": capabilities[1:], "actor_selector": map[string]any{"actor_id": actor, "required_capabilities": capabilities[1:]}, "scope": map[string]any{"issue_id": issueID}}}, "obligations": []any{map[string]any{"id": "OBL-p22-corpus", "type": "gate-evidence", "applicability": "required", "guarded_actions": []string{"envelope.complete", "issue.transition-terminal-success"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": map[string]any{"id": registrationID, "version": registrationVersion}, "accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"}, "freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": source}}, "budgets": []any{map[string]any{"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 3, "warning_threshold": nil, "exhaustion_action": "unavailable", "consumed": nil, "remaining": nil}}})
	created := p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "create", "--actor", actor, "--input", envelopeInput, "--idempotency-key", "p22-corpus-envelope-create")
	envelopeID := p21PathString(t, created, "result", "envelope", "id")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-envelope-activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1}), "--idempotency-key", "p22-corpus-envelope-activate")
	claim := p21Run(t, binary, workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-claim", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "worktree_id": repository["worktree_id"], "scope": map[string]any{"paths": []string{"base.txt"}}, "ttl_seconds": 60}), "--idempotency-key", "p22-corpus-claim-acquire")
	claimID := p21PathString(t, claim, "result", "claim", "id")
	run := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "start", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-run", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": claimID, "claim_version": 1, "executor": map[string]any{"id": "executor-corpus-prior", "source": "activated-harness"}, "obligation_id": nil, "repository": repository}), "--idempotency-key", "p22-corpus-run-start")
	return &p22FIDRecoveryAdapter{binary: binary, workspace: workspace, projectID: projectID, actor: actor, repository: repository, claimID: claimID, claimVersion: 1, runID: p21PathString(t, run, "result", "run", "id"), runVersion: 1}
}

func p22CorpusRepository(t *testing.T, workspace, projectID string) map[string]any {
	t.Helper()
	observation, err := project.ObserveRepository(workspace)
	if err != nil {
		t.Fatalf("observe corpus repository: %v", err)
	}
	return map[string]any{"base_oid": observation.HeadOID, "capture_source": "git", "dirty": observation.Dirty, "dirty_digest": observation.DirtyDigest, "head_oid": observation.HeadOID, "project_id": projectID, "tree_oid": observation.TreeOID, "worktree_id": observation.WorktreeID}
}

func (adapter *p22FIDRecoveryAdapter) Execute(t *testing.T, fixture p22ConformanceCase) p22ConformanceResponse {
	t.Helper()
	args, input := adapter.p22Args(t, fixture)
	if input != nil {
		args = append(args, "--input", p21Input(t, adapter.workspace, fixture.Name, input), "--idempotency-key", "p22-corpus-"+fixture.Name)
	}
	stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, args...)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	} else {
		response.ExitCode = 0
	}
	if fixture.Kind == "accepted" {
		adapter.p22Advance(t, fixture.Operation, stdout)
	}
	return response
}

func (adapter *p22FIDRecoveryAdapter) p22Args(t *testing.T, fixture p22ConformanceCase) ([]string, map[string]any) {
	t.Helper()
	args := []string{"--project", adapter.workspace, "--json"}
	unknown := fixture.Kind == "unknown-field"
	var input map[string]any
	switch fixture.Operation {
	case "run.observe":
		args = append(args, "run", "observe", "--actor", adapter.actor)
		input = map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "liveness": map[string]any{"executor_id": "executor-corpus-prior", "state": "terminated", "source_ref": "process:p22-corpus", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}}
	case "run.inspect":
		args = append(args, "run", "inspect", adapter.runID)
		if unknown {
			args = append(args, "--unknown-field")
		}
		return args, nil
	case "run.checkpoint":
		args = append(args, "run", "checkpoint", "--actor", adapter.actor)
		input = map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "reconciliation_digest": adapter.digest}
	case "claim.expire":
		args = append(args, "claim", "expire", "--actor", adapter.actor)
		input = map[string]any{"claim_id": adapter.claimID, "expected_version": adapter.claimVersion}
	case "run.recover":
		args = append(args, "run", "recover", "--actor", adapter.actor)
		input = map[string]any{"prior_run_id": adapter.runID, "expected_run_version": adapter.runVersion, "expected_claim_version": adapter.claimVersion, "observation_id": adapter.observation, "reconciliation_digest": adapter.digest, "adopted_paths": []string{}, "authority": map[string]any{"mode": "ordinary", "actor_id": adapter.actor, "reason": "trusted termination", "reference": nil, "grant_ref": nil}, "executor": map[string]any{"id": "executor-corpus-successor", "source": "activated-harness"}, "ttl_seconds": 60}
	default:
		t.Fatalf("unsupported Phase 22 conformance operation %q", fixture.Operation)
	}
	if unknown {
		input["unknown_field"] = true
	}
	return args, input
}

func (adapter *p22FIDRecoveryAdapter) p22Advance(t *testing.T, operation, stdout string) {
	t.Helper()
	var envelope map[string]any
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil {
		t.Fatalf("decode %s accepted response: %v; stdout=%q", operation, err, stdout)
	}
	result := p21Path(t, envelope, "result").(map[string]any)
	switch operation {
	case "run.observe":
		adapter.observation = p21PathString(t, envelope, "result", "observation", "id")
		adapter.runVersion = int(p21Path(t, envelope, "result", "run", "version").(float64))
	case "run.inspect":
		adapter.digest = p21PathString(t, envelope, "result", "reconciliation", "digest")
	case "run.checkpoint":
		adapter.runVersion = int(p21Path(t, envelope, "result", "run", "version").(float64))
	case "claim.expire":
		adapter.claimVersion = int(p21Path(t, envelope, "result", "claim", "version").(float64))
	case "run.recover":
		if p21PathString(t, envelope, "result", "recovery", "prior_run_id") != adapter.runID {
			t.Fatalf("recovery response = %#v, want prior Run %q", result, adapter.runID)
		}
	}
}

func TestP22RealBinaryRecoveryFailureCorpus(t *testing.T) {
	for _, fixture := range []struct {
		name, operation, code string
		exit                  int
		run                   func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse
	}{
		{"not-found", "run.inspect", "run-not-found", 3, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			return adapter.p22Run(t, "run", "inspect", "RUN-missing")
		}},
		{"authority", "run.observe", "authority-denied", 5, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			return adapter.p22RunInput(t, "run", "observe", "actor-denied", map[string]any{"run_id": adapter.runID, "expected_version": 1, "liveness": map[string]any{"executor_id": "executor-corpus-prior", "state": "terminated", "source_ref": "process:p22-authority", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
		}},
		{"liveness", "run.recover", "executor-still-live", 4, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			adapter.p22Observe(t, "live")
			adapter.p22Expire(t)
			return adapter.p22Recover(t)
		}},
		{"revision", "run.observe", "run-version-conflict", 4, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			return adapter.p22RunInput(t, "run", "observe", adapter.actor, map[string]any{"run_id": adapter.runID, "expected_version": 99, "liveness": map[string]any{"executor_id": "executor-corpus-prior", "state": "terminated", "source_ref": "process:p22-revision", "observed_at": time.Now().UTC().Format(time.RFC3339Nano)}})
		}},
		{"conflict", "run.checkpoint", "recovery-state-conflict", 4, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			adapter.p22Observe(t, "terminated")
			return adapter.p22RunInput(t, "run", "checkpoint", adapter.actor, map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "reconciliation_digest": "sha256:" + p22Zeros(64)})
		}},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			response := fixture.run(t, p22NewFIDRecoveryAdapter(t))
			assertP22ProcessFailure(t, fixture.operation, fixture.code, fixture.exit, response)
		})
	}
}

func (adapter *p22FIDRecoveryAdapter) p22Run(t *testing.T, command ...string) p22ConformanceResponse {
	t.Helper()
	args := append([]string{"--project", adapter.workspace, "--json"}, command...)
	stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, args...)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return response
}

func (adapter *p22FIDRecoveryAdapter) p22RunInput(t *testing.T, group, command, actor string, input map[string]any) p22ConformanceResponse {
	t.Helper()
	args := []string{"--project", adapter.workspace, "--json", group, command, "--actor", actor, "--input", p21Input(t, adapter.workspace, "p22-failure-"+command, input), "--idempotency-key", "p22-failure-" + command + "-" + actor}
	stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, args...)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return response
}

func (adapter *p22FIDRecoveryAdapter) p22Observe(t *testing.T, state string) {
	t.Helper()
	response := adapter.p22RunInput(t, "run", "observe", adapter.actor, map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "liveness": map[string]any{"executor_id": "executor-corpus-prior", "state": state, "source_ref": "process:p22-failure", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
	if response.ExitCode != 0 {
		t.Fatalf("observe %s = exit %d stdout=%q stderr=%q", state, response.ExitCode, response.Stdout, response.Stderr)
	}
	adapter.p22Advance(t, "run.observe", response.Stdout)
	response = adapter.p22Run(t, "run", "inspect", adapter.runID)
	if response.ExitCode != 0 {
		t.Fatalf("inspect after observe = exit %d stdout=%q stderr=%q", response.ExitCode, response.Stdout, response.Stderr)
	}
	adapter.p22Advance(t, "run.inspect", response.Stdout)
}

func (adapter *p22FIDRecoveryAdapter) p22Expire(t *testing.T) {
	t.Helper()
	response := adapter.p22RunInput(t, "claim", "expire", adapter.actor, map[string]any{"claim_id": adapter.claimID, "expected_version": adapter.claimVersion})
	if response.ExitCode != 0 {
		t.Fatalf("expire Claim = exit %d stdout=%q stderr=%q", response.ExitCode, response.Stdout, response.Stderr)
	}
	adapter.p22Advance(t, "claim.expire", response.Stdout)
}

func (adapter *p22FIDRecoveryAdapter) p22Recover(t *testing.T) p22ConformanceResponse {
	t.Helper()
	return adapter.p22RunInput(t, "run", "recover", adapter.actor, map[string]any{"prior_run_id": adapter.runID, "expected_run_version": adapter.runVersion, "expected_claim_version": adapter.claimVersion, "observation_id": adapter.observation, "reconciliation_digest": adapter.digest, "adopted_paths": []string{}, "authority": map[string]any{"mode": "ordinary", "actor_id": adapter.actor, "reason": "failure corpus", "reference": nil, "grant_ref": nil}, "executor": map[string]any{"id": "executor-corpus-failure-successor", "source": "activated-harness"}, "ttl_seconds": 60})
}

func assertP22ProcessFailure(t *testing.T, operation, code string, exit int, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != exit || response.Stdout != "" {
		t.Fatalf("%s failure = exit %d stdout=%q stderr=%q, want exit %d and no stdout", operation, response.ExitCode, response.Stdout, response.Stderr, exit)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatalf("decode %s failure: %v; stderr=%q", operation, err, response.Stderr)
	}
	errorValue, ok := envelope["error"].(map[string]any)
	if !ok || envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != false || envelope["operation"] != operation || errorValue["code"] != code {
		t.Fatalf("%s failure envelope = %#v", operation, envelope)
	}
	if fmt.Sprint(errorValue["code"]) != code {
		t.Fatalf("%s failure code = %#v, want %q", operation, errorValue["code"], code)
	}
}
