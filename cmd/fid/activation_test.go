package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestActivationPrepareCreatesPreparedActivation(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	inputPath := filepath.Join(workspace, "activation-prepare.json")
	writeActivationInput(t, inputPath, map[string]any{
		"backend_identity": "fid",
		"direction":        "activate",
		"kit_identity":     "kit",
		"phase":            21,
	})

	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "activation", "prepare",
		"--input", inputPath, "--idempotency-key", "activation-prepare-1",
	)
	if result.err != nil {
		t.Fatalf("activation prepare error = %v, stderr = %q", result.err, result.stderr)
	}
	if result.stderr != "" {
		t.Fatalf("activation prepare stderr = %q, want empty", result.stderr)
	}

	var envelope struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Result    struct {
			Activation struct {
				ID      string `json:"id"`
				Status  string `json:"status"`
				Phase   int    `json:"phase"`
				PointOf bool   `json:"point_of_no_return"`
			} `json:"activation"`
			IdempotencyKey string `json:"idempotency_key"`
			Replayed       bool   `json:"replayed"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("decode activation prepare response: %v\nstdout=%q", err, result.stdout)
	}
	if !envelope.OK || envelope.Operation != "activation.prepare" || envelope.Result.Activation.ID == "" ||
		envelope.Result.Activation.Status != "prepared" || envelope.Result.Activation.Phase != 21 ||
		envelope.Result.Activation.PointOf || envelope.Result.IdempotencyKey != "activation-prepare-1" || envelope.Result.Replayed {
		t.Fatalf("activation prepare response = %+v, want prepared non-replayed activation", envelope)
	}
}

func TestActivationComponentStagesValidatedProfilesAndRejectsCallerAugmentation(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	componentPath := filepath.Join(workspace, "activation-profile.json")
	writeActivationInput(t, componentPath, map[string]any{
		"activation_id":         activationID,
		"component_type":        "authority-profile",
		"component_id":          "builders",
		"component_version":     1,
		"expected_prior_digest": nil,
		"content": []any{map[string]any{
			"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit",
			"provenance_ref": "activation:profile/builders", "session_ref": nil,
			"capabilities": []string{"claim.acquire", "gate.register"},
		}},
	})

	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "activation", "component",
		"--input", componentPath, "--idempotency-key", "activation-component-1",
	)
	if result.err != nil {
		t.Fatalf("activation component error = %v, stderr = %q", result.err, result.stderr)
	}
	var envelope struct {
		Result struct {
			Activation struct {
				Status           string   `json:"status"`
				ComponentDigests []string `json:"component_digests"`
			} `json:"activation"`
			Component struct {
				Type   string `json:"type"`
				ID     string `json:"id"`
				Digest string `json:"digest"`
			} `json:"component"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("decode activation component response: %v\nstdout=%q", err, result.stdout)
	}
	if envelope.Result.Activation.Status != "staged" || len(envelope.Result.Activation.ComponentDigests) != 1 ||
		envelope.Result.Component.Type != "authority-profile" || envelope.Result.Component.ID != "builders" || envelope.Result.Component.Digest == "" {
		t.Fatalf("activation component response = %+v, want staged profile component", envelope)
	}

	writeActivationInput(t, componentPath, map[string]any{
		"activation_id": activationID, "component_type": "authority-profile", "component_id": "builders",
		"component_version": 2, "expected_prior_digest": nil,
		"content": []any{map[string]any{
			"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit",
			"provenance_ref": "activation:profile/builders", "session_ref": nil,
			"capabilities": []string{"claim.acquire"}, "caller_capabilities": []string{"envelope.create"},
		}},
	})
	result = executeCommandForTest(workspace,
		"--project", workspace, "--json", "activation", "component",
		"--input", componentPath, "--idempotency-key", "activation-component-2",
	)
	if result.err == nil {
		t.Fatal("caller capability augmentation error = nil, want closed-input rejection")
	}
	if result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
		t.Fatalf("caller augmentation result stdout=%q stderr=%q, want typed closed-input rejection", result.stdout, result.stderr)
	}
}

func prepareActivationForTest(t *testing.T, workspace string) string {
	t.Helper()
	inputPath := filepath.Join(workspace, "activation-prepare.json")
	writeActivationInput(t, inputPath, map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "kit", "phase": 21})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "prepare", "--input", inputPath, "--idempotency-key", "activation-prepare-helper")
	if result.err != nil {
		t.Fatalf("prepare activation: %v; stderr=%q", result.err, result.stderr)
	}
	var envelope struct {
		Result struct {
			Activation struct {
				ID string `json:"id"`
			} `json:"activation"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil || envelope.Result.Activation.ID == "" {
		t.Fatalf("decode prepared activation: %v; stdout=%q", err, result.stdout)
	}
	return envelope.Result.Activation.ID
}

func TestActivationCommitIsPointOfNoReturn(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "builders", []any{map[string]any{
		"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/builders", "session_ref": nil,
		"capabilities": []string{"claim.acquire", "gate.register"},
	}})
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", []any{map[string]any{
		"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"go", "test", "./..."}, "credential_names": []string{}, "cwd": ".",
		"definition_id": "gate-unit", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024,
		"redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed",
	}})
	commitPath := filepath.Join(workspace, "activation-commit.json")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})

	result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-1")
	if result.err != nil {
		t.Fatalf("activation commit error = %v, stderr = %q", result.err, result.stderr)
	}
	var committed struct {
		Result struct {
			Activation struct {
				Status  string `json:"status"`
				PointOf bool   `json:"point_of_no_return"`
			} `json:"activation"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &committed); err != nil || committed.Result.Activation.Status != "committed" || !committed.Result.Activation.PointOf {
		t.Fatalf("commit response error=%v response=%q, want committed point of no return", err, result.stdout)
	}

	abortPath := filepath.Join(workspace, "activation-abort.json")
	writeActivationInput(t, abortPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "reason": "retry"})
	result = executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "abort", "--input", abortPath, "--idempotency-key", "activation-abort-after-commit")
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"activation-recovery-unavailable"`) {
		t.Fatalf("abort after commit result error=%v stdout=%q stderr=%q, want recovery-unavailable failure", result.err, result.stdout, result.stderr)
	}
}

func TestActivationFinalizeAndStatusExposeOnlyFinalLifecycleState(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "builders", []any{map[string]any{
		"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/builders", "session_ref": nil, "capabilities": []string{"claim.acquire", "gate.register"},
	}})
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", validGateDefinitions())
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-commit.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	commit := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-finalize")
	if commit.err != nil {
		t.Fatalf("commit activation: %v; stderr=%q", commit.err, commit.stderr)
	}

	finalizePath := filepath.Join(workspace, "activation-finalize.json")
	writeActivationInput(t, finalizePath, map[string]any{"activation_id": activationID, "expected_status": "committed"})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "finalize", "--input", finalizePath, "--idempotency-key", "activation-finalize-1")
	if result.err != nil {
		t.Fatalf("finalize activation: %v; stderr=%q", result.err, result.stderr)
	}
	result = executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "status", activationID)
	if result.err != nil {
		t.Fatalf("activation status: %v; stderr=%q", result.err, result.stderr)
	}
	var envelope struct {
		Result struct {
			Activation struct {
				Status           string   `json:"status"`
				PointOf          bool     `json:"point_of_no_return"`
				ComponentDigests []string `json:"component_digests"`
			} `json:"activation"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil || envelope.Result.Activation.Status != "finalized" || !envelope.Result.Activation.PointOf || len(envelope.Result.Activation.ComponentDigests) != 2 {
		t.Fatalf("final status error=%v stdout=%q, want finalized activation", err, result.stdout)
	}
}

func TestFinalizedActivationIsTheOnlyTrustedActorSource(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "builders", []any{map[string]any{
		"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/builders", "session_ref": nil, "capabilities": []string{"claim.acquire", "gate.register"},
	}})
	if _, err := workflow.ResolveTrustedActor(workspace, issueMigrationsDir(), "actor-builder"); err == nil {
		t.Fatal("staged actor resolved as trusted")
	}
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", validGateDefinitions())
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-commit-trust.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-trust"); result.err != nil {
		t.Fatalf("commit: %v", result.err)
	}
	if _, err := workflow.ResolveTrustedActor(workspace, issueMigrationsDir(), "actor-builder"); err == nil {
		t.Fatal("committed actor resolved before finalize")
	}
	finalizePath := filepath.Join(workspace, "activation-finalize-trust.json")
	writeActivationInput(t, finalizePath, map[string]any{"activation_id": activationID, "expected_status": "committed"})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "finalize", "--input", finalizePath, "--idempotency-key", "activation-finalize-trust"); result.err != nil {
		t.Fatalf("finalize: %v", result.err)
	}
	actor, err := workflow.ResolveTrustedActor(workspace, issueMigrationsDir(), "actor-builder")
	if err != nil || actor.ActorType != "agent" || actor.Producer != "kit" || actor.ProvenanceRef != "activation:profile/builders" || !reflect.DeepEqual(actor.Capabilities, []string{"claim.acquire", "gate.register"}) {
		t.Fatalf("trusted actor=%+v err=%v, want finalized profile values", actor, err)
	}
}

func validGateDefinitions() []any {
	return []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"go", "test", "./..."}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-unit", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}}
}

func stageActivationComponentForTest(t *testing.T, workspace, activationID, componentType, componentID string, content []any) string {
	t.Helper()
	inputPath := filepath.Join(workspace, componentID+".json")
	writeActivationInput(t, inputPath, map[string]any{
		"activation_id": activationID, "component_type": componentType, "component_id": componentID,
		"component_version": 1, "expected_prior_digest": nil, "content": content,
	})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "component", "--input", inputPath, "--idempotency-key", "component-"+componentID)
	if result.err != nil {
		t.Fatalf("stage %s component: %v; stderr=%q", componentType, result.err, result.stderr)
	}
	var envelope struct {
		Result struct {
			Component struct {
				Digest string `json:"digest"`
			} `json:"component"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil || envelope.Result.Component.Digest == "" {
		t.Fatalf("decode %s component: %v; stdout=%q", componentType, err, result.stdout)
	}
	return envelope.Result.Component.Digest
}

func writeActivationInput(t *testing.T, path string, value map[string]any) {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal activation input: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("write activation input: %v", err)
	}
}
