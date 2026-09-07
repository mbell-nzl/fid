package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestActivationPrepareCreatesPreparedActivation(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	inputPath := filepath.Join(workspace, "activation-prepare.json")
	writeActivationInput(t, inputPath, map[string]any{
		"backend_identity":             "fid",
		"direction":                    "activate",
		"kit_identity":                 "kit",
		"phase":                        21,
		"expected_prior_activation_id": nil,
		"expected_prior_version":       nil,
		"resources":                    activationResources(),
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

func TestActivationReplaysAcrossLifecycleWithoutPreemptingLifecycleChecks(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	preparePath := filepath.Join(workspace, "activation-replay-prepare.json")
	writeActivationInput(t, preparePath, activationPrepareInput("kit"))
	prepare := assertActivationReplay(t, workspace, "prepare", preparePath, "activation-replay-prepare")
	activationID := activationIDFromReplay(t, prepare)

	newKey := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "prepare", "--input", preparePath, "--idempotency-key", "activation-replay-new-key")
	assertActivationTypedFailure(t, newKey, "activation-conflict")

	profilePath := filepath.Join(workspace, "activation-replay-profile.json")
	writeActivationInput(t, profilePath, map[string]any{
		"activation_id": activationID, "component_type": "authority-profile", "component_id": "replay-builders",
		"component_version": 1, "expected_prior_digest": nil, "expected_version": activationVersionForTest(t, workspace, activationID),
		"content": []any{map[string]any{
			"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/replay-builders", "session_ref": nil,
			"capabilities": []string{"claim.acquire", "gate.register"},
		}},
	})
	profile := assertActivationReplay(t, workspace, "component", profilePath, "activation-replay-profile")
	profileDigest := componentDigestFromReplay(t, profile)

	writeActivationInput(t, profilePath, map[string]any{
		"activation_id": activationID, "component_type": "authority-profile", "component_id": "replay-builders",
		"component_version": 1, "expected_prior_digest": nil, "expected_version": 2,
		"content": []any{map[string]any{
			"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/replay-builders-v2", "session_ref": nil,
			"capabilities": []string{"claim.acquire", "gate.register"},
		}},
	})
	assertActivationStableDigestConflict(t, workspace, "component", profilePath, "activation-replay-profile")

	gatePath := filepath.Join(workspace, "activation-replay-gates.json")
	writeActivationInput(t, gatePath, map[string]any{
		"activation_id": activationID, "component_type": "gate-definitions", "component_id": "replay-gates",
		"component_version": 1, "expected_prior_digest": nil, "expected_version": activationVersionForTest(t, workspace, activationID), "content": validGateDefinitions(),
	})
	gate := assertActivationReplay(t, workspace, "component", gatePath, "activation-replay-gates")
	gateDigest := componentDigestFromReplay(t, gate)
	observeActivationResourcesForTest(t, workspace, activationID, "new")

	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-replay-commit.json")
	writeActivationInput(t, commitPath, activationCommitInput(t, workspace, activationID, digests))
	assertActivationReplay(t, workspace, "commit", commitPath, "activation-replay-commit")

	finalizePath := filepath.Join(workspace, "activation-replay-finalize.json")
	writeActivationInput(t, finalizePath, activationFinalizeInput(t, workspace, activationID))
	assertActivationReplay(t, workspace, "finalize", finalizePath, "activation-replay-finalize")

	abortWorkspace := clonePristineWorkspace(t)
	abortPreparePath := filepath.Join(abortWorkspace, "activation-replay-abort-prepare.json")
	writeActivationInput(t, abortPreparePath, activationPrepareInput("kit"))
	abortPrepare := assertActivationReplay(t, abortWorkspace, "prepare", abortPreparePath, "activation-replay-abort-prepare")
	abortPath := filepath.Join(abortWorkspace, "activation-replay-abort.json")
	writeActivationInput(t, abortPath, map[string]any{"activation_id": activationIDFromReplay(t, abortPrepare), "expected_status": "prepared", "expected_version": 1, "reason": "replay"})
	assertActivationReplay(t, abortWorkspace, "abort", abortPath, "activation-replay-abort")
}

func assertActivationReplay(t *testing.T, workspace, operation, inputPath, key string) map[string]any {
	t.Helper()
	args := []string{"--project", workspace, "--json", "activation", operation, "--input", inputPath, "--idempotency-key", key}
	before := journalEventCount(t, workspace)
	first := executeCommandForTest(workspace, args...)
	if first.err != nil || first.stderr != "" {
		t.Fatalf("first activation %s = %#v, want stdout-only success", operation, first)
	}
	afterFirst := journalEventCount(t, workspace)
	if afterFirst != before+1 {
		t.Fatalf("first activation %s Event count = %d, want %d", operation, afterFirst, before+1)
	}
	second := executeCommandForTest(workspace, args...)
	t.Logf("pre-repair second activation %s stdout=%q stderr=%q", operation, second.stdout, second.stderr)
	if second.err != nil || second.stderr != "" {
		t.Fatalf("second activation %s = %#v, want stdout-only replay", operation, second)
	}
	if afterSecond := journalEventCount(t, workspace); afterSecond != afterFirst {
		t.Fatalf("replayed activation %s Event count = %d, want %d", operation, afterSecond, afterFirst)
	}

	firstEnvelope := decodeActivationEnvelope(t, first.stdout)
	secondEnvelope := decodeActivationEnvelope(t, second.stdout)
	firstResult, ok := firstEnvelope["result"].(map[string]any)
	if !ok {
		t.Fatalf("first activation %s result = %#v, want object", operation, firstEnvelope["result"])
	}
	firstResult["replayed"] = true
	if !reflect.DeepEqual(firstEnvelope, secondEnvelope) {
		t.Fatalf("activation %s replay = %#v, want original Event/entity/domain fields with only replayed:true from %#v", operation, secondEnvelope, firstEnvelope)
	}
	return secondEnvelope
}

func assertActivationStableDigestConflict(t *testing.T, workspace, operation, inputPath, key string) {
	t.Helper()
	args := []string{"--project", workspace, "--json", "activation", operation, "--input", inputPath, "--idempotency-key", key}
	before := journalEventCount(t, workspace)
	first := executeCommandForTest(workspace, args...)
	assertActivationTypedFailure(t, first, "invalid-argument")
	second := executeCommandForTest(workspace, args...)
	assertActivationTypedFailure(t, second, "invalid-argument")
	if first.stderr != second.stderr {
		t.Fatalf("same-key different-digest failures differ: first=%q second=%q", first.stderr, second.stderr)
	}
	if after := journalEventCount(t, workspace); after != before {
		t.Fatalf("same-key different-digest Event count = %d, want %d", after, before)
	}
}

func assertActivationTypedFailure(t *testing.T, result commandExecutionResult, code string) {
	t.Helper()
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"`+code+`"`) {
		t.Fatalf("activation failure = %#v, want stderr-only %s envelope", result, code)
	}
}

func decodeActivationEnvelope(t *testing.T, raw string) map[string]any {
	t.Helper()
	var envelope map[string]any
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		t.Fatalf("decode activation envelope: %v; stdout=%q", err, raw)
	}
	return envelope
}

func activationIDFromReplay(t *testing.T, envelope map[string]any) string {
	t.Helper()
	return activationReplayString(t, envelope, "activation", "id")
}

func componentDigestFromReplay(t *testing.T, envelope map[string]any) string {
	t.Helper()
	return activationReplayString(t, envelope, "component", "digest")
}

func activationReplayString(t *testing.T, envelope map[string]any, object, field string) string {
	t.Helper()
	result, ok := envelope["result"].(map[string]any)
	if !ok {
		t.Fatalf("activation replay result = %#v, want object", envelope["result"])
	}
	value, ok := result[object].(map[string]any)
	if !ok {
		t.Fatalf("activation replay %s = %#v, want object", object, result[object])
	}
	text, ok := value[field].(string)
	if !ok || text == "" {
		t.Fatalf("activation replay %s.%s = %#v, want nonempty string", object, field, value[field])
	}
	return text
}

func TestActivationComponentStagesValidatedProfilesAndRejectsCallerAugmentation(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	componentPath := filepath.Join(workspace, "activation-profile.json")
	writeActivationInput(t, componentPath, map[string]any{
		"activation_id":         activationID,
		"component_type":        "authority-profile",
		"component_id":          "builders",
		"component_version":     1,
		"expected_prior_digest": nil,
		"expected_version":      activationVersionForTest(t, workspace, activationID),
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

	for _, testCase := range []struct {
		name  string
		apply func(map[string]any)
	}{
		{"capabilities", func(profile map[string]any) { profile["caller_capabilities"] = []string{"envelope.create"} }},
		{"provenance", func(profile map[string]any) { profile["caller_provenance_ref"] = "caller:override" }},
	} {
		t.Run("caller "+testCase.name, func(t *testing.T) {
			profile := map[string]any{
				"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit",
				"provenance_ref": "activation:profile/builders", "session_ref": nil,
				"capabilities": []string{"claim.acquire"},
			}
			testCase.apply(profile)
			writeActivationInput(t, componentPath, map[string]any{
				"activation_id": activationID, "component_type": "authority-profile", "component_id": "builders",
				"component_version": 2, "expected_prior_digest": nil, "expected_version": activationVersionForTest(t, workspace, activationID), "content": []any{profile},
			})
			before := journalEventCount(t, workspace)
			result = executeCommandForTest(workspace,
				"--project", workspace, "--json", "activation", "component",
				"--input", componentPath, "--idempotency-key", "activation-component-reject-"+testCase.name,
			)
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
				t.Fatalf("caller %s result = %#v, want typed closed-input rejection", testCase.name, result)
			}
			if after := journalEventCount(t, workspace); after != before {
				t.Fatalf("caller %s Event count = %d, want no Event after %d", testCase.name, after, before)
			}
		})
	}
}

func TestActivationComponentResponseDigestsRemainSortedAfterSmallerSecondStage(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	profileContent := []any{map[string]any{"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/sorted", "session_ref": nil, "capabilities": []string{"claim.acquire"}}}
	firstDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "sorted-profile", profileContent)
	var smaller []any
	for timeout := 1; timeout < 10000; timeout++ {
		candidate := validGateDefinitions()
		candidate[0].(map[string]any)["timeout_ms"] = timeout
		if activationContentDigest(t, candidate) < firstDigest {
			smaller = candidate
			break
		}
	}
	if smaller == nil {
		t.Fatal("could not construct a lexicographically smaller second component digest")
	}
	path := filepath.Join(workspace, "sorted-gates.json")
	writeActivationInput(t, path, map[string]any{"activation_id": activationID, "component_type": "gate-definitions", "component_id": "sorted-gates", "component_version": 1, "expected_prior_digest": nil, "expected_version": activationVersionForTest(t, workspace, activationID), "content": smaller})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "component", "--input", path, "--idempotency-key", "sorted-gates")
	if result.err != nil {
		t.Fatalf("stage smaller second component: %#v", result)
	}
	var envelope struct {
		Result struct {
			Activation struct {
				ComponentDigests []string `json:"component_digests"`
			} `json:"activation"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil || len(envelope.Result.Activation.ComponentDigests) != 2 || !sort.StringsAreSorted(envelope.Result.Activation.ComponentDigests) {
		t.Fatalf("staged component digests = %#v, decode error = %v, want sorted exact pair", envelope.Result.Activation.ComponentDigests, err)
	}
}

func activationContentDigest(t *testing.T, content []any) string {
	t.Helper()
	raw, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal component content: %v", err)
	}
	canonical, err := journal.MarshalCanonicalJSON(json.RawMessage(raw))
	if err != nil {
		t.Fatalf("canonical component content: %v", err)
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func prepareActivationForTest(t *testing.T, workspace string) string {
	t.Helper()
	inputPath := filepath.Join(workspace, "activation-prepare.json")
	writeActivationInput(t, inputPath, activationPrepareInput("kit"))
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
	t.Parallel()
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
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	commitPath := filepath.Join(workspace, "activation-commit.json")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	writeActivationInput(t, commitPath, activationCommitInput(t, workspace, activationID, digests))

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
	writeActivationInput(t, abortPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, workspace, activationID), "reason": "retry"})
	result = executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "abort", "--input", abortPath, "--idempotency-key", "activation-abort-after-commit")
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"activation-direction-conflict"`) {
		t.Fatalf("abort after commit result error=%v stdout=%q stderr=%q, want forward-only direction failure", result.err, result.stdout, result.stderr)
	}
}

func TestCommittedActivationCanOnlyRecoverForward(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "forward-builders", []any{map[string]any{
		"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/forward-builders", "session_ref": nil,
		"capabilities": []string{"claim.acquire", "gate.register"},
	}})
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "forward-gates", validGateDefinitions())
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-forward-commit.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, workspace, activationID), "component_digests": digests})
	committed := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-forward-commit")
	if committed.err != nil {
		t.Fatalf("commit activation: %v; stderr=%q", committed.err, committed.stderr)
	}

	abortPath := filepath.Join(workspace, "activation-forward-abort.json")
	writeActivationInput(t, abortPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, workspace, activationID), "reason": "must complete forward"})
	aborted := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "abort", "--input", abortPath, "--idempotency-key", "activation-forward-abort")
	assertActivationTypedFailure(t, aborted, "activation-direction-conflict")
}

func TestActivationRecoveryPreservesDirectionAndIdentityAcrossReplay(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "recovery-builders", []any{map[string]any{"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/recovery-builders", "session_ref": nil, "capabilities": []string{"claim.acquire", "gate.register"}}})
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "recovery-gates", validGateDefinitions())
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-recovery-commit.json")
	writeActivationInput(t, commitPath, activationCommitInput(t, workspace, activationID, digests))
	if committed := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-recovery-commit"); committed.err != nil {
		t.Fatalf("commit recovery fixture: %v; stderr=%q", committed.err, committed.stderr)
	}

	activation, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), activationID)
	if err != nil {
		t.Fatalf("read committed activation: %v", err)
	}
	options := workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), BootstrapActor: "bootstrap-recovery", IdempotencyKey: "activation-recovery-select"}
	request := workflow.ActivationRecoveryRequest{ActivationID: activationID, ExpectedVersion: activation.Version, ExpectedStatus: "committed", Direction: "complete-forward", Reason: "response loss recovery"}
	selected, err := workflow.RecoverActivation(options, request)
	if err != nil || selected.Replayed || selected.Activation.Recovery == nil || selected.Activation.Recovery.Direction != "complete-forward" || selected.Activation.Recovery.SelectedEventID != selected.EventID {
		t.Fatalf("RecoverActivation() = (%#v, %v), want selected durable forward recovery", selected, err)
	}
	replayed, err := workflow.RecoverActivation(options, request)
	if err != nil || !replayed.Replayed || replayed.Activation.ID != activationID || replayed.Activation.Recovery == nil || replayed.Activation.Recovery.Direction != "complete-forward" || replayed.EventID != selected.EventID {
		t.Fatalf("RecoverActivation(replay) = (%#v, %v), want same identity and direction", replayed, err)
	}
	resourceVersions := make([]workflow.ActivationResourceVersion, 0, len(selected.Activation.Resources))
	for _, resource := range selected.Activation.Resources {
		resourceVersions = append(resourceVersions, workflow.ActivationResourceVersion{Component: resource.Component, Version: resource.Version})
	}
	completed, err := workflow.CompleteActivationRecovery(workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), BootstrapActor: "bootstrap-recovery", IdempotencyKey: "activation-recovery-complete"}, workflow.ActivationRecoveryCompleteRequest{ActivationID: activationID, ExpectedVersion: selected.Activation.Version, SelectedEventID: selected.EventID, ResourceVersions: resourceVersions})
	if err != nil || completed.Activation.Status != "finalized" || completed.Activation.ID != activationID || completed.Activation.Recovery == nil || completed.Activation.Recovery.CompletedEventID == nil || *completed.Activation.Recovery.CompletedEventID != completed.EventID {
		t.Fatalf("CompleteActivationRecovery() = (%#v, %v), want finalized selected activation", completed, err)
	}
}

func TestIncompleteUpgradeBlocksNormalWorkWhileInspectionAndRestoreRemainAvailable(t *testing.T) {
	t.Parallel()
	workspace, _, _ := buildFinalizedEnvelopeWorkspace(t, validGateDefinitions())
	if _, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), "missing"); err == nil {
		t.Fatal("missing activation unexpectedly inspected")
	}
	connActivation := currentFinalizedActivationForTest(t, workspace)
	upgradeID := prepareUpgradeActivationForTest(t, workspace, connActivation)
	if _, err := workflow.ResolveTrustedActor(workspace, issueMigrationsDir(), "actor-builder"); err == nil {
		t.Fatal("incomplete upgrade left predecessor actor trusted")
	}
	registrationPath := filepath.Join(workspace, "incomplete-upgrade-register.json")
	writeActivationInput(t, registrationPath, gateRegisterInput(t, workspace))
	blocked := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "register", "--actor", "actor-builder", "--input", registrationPath, "--idempotency-key", "incomplete-upgrade-register")
	if blocked.err == nil || !strings.Contains(blocked.stderr, `"code":"authority-unavailable"`) {
		t.Fatalf("normal gate registration during incomplete upgrade = %#v, want authority-unavailable", blocked)
	}
	if inspected, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), upgradeID); err != nil || inspected.Status != "prepared" {
		t.Fatalf("ActivationStatus(incomplete upgrade) = (%#v, %v), want readable prepared activation", inspected, err)
	}
	upgrade, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), upgradeID)
	if err != nil {
		t.Fatalf("read upgrade: %v", err)
	}
	selected, err := workflow.RecoverActivation(workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), BootstrapActor: "bootstrap-upgrade", IdempotencyKey: "upgrade-restore"}, workflow.ActivationRecoveryRequest{ActivationID: upgradeID, ExpectedVersion: upgrade.Version, ExpectedStatus: "prepared", Direction: "restore-prior", Reason: "restore prior trust"})
	if err != nil || selected.Activation.Recovery == nil || selected.Activation.Recovery.Direction != "restore-prior" {
		t.Fatalf("RecoverActivation(restore) = (%#v, %v)", selected, err)
	}
	resourceVersions := make([]workflow.ActivationResourceVersion, 0, len(selected.Activation.Resources))
	for _, resource := range selected.Activation.Resources {
		resourceVersions = append(resourceVersions, workflow.ActivationResourceVersion{Component: resource.Component, Version: resource.Version})
	}
	observeActivationResourcesForTest(t, workspace, upgradeID, "old")
	recoveredActivation, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), upgradeID)
	if err != nil {
		t.Fatalf("read restored upgrade: %v", err)
	}
	resourceVersions = resourceVersions[:0]
	for _, resource := range recoveredActivation.Resources {
		resourceVersions = append(resourceVersions, workflow.ActivationResourceVersion{Component: resource.Component, Version: resource.Version})
	}
	if _, err := workflow.CompleteActivationRecovery(workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), BootstrapActor: "bootstrap-upgrade", IdempotencyKey: "upgrade-restore-complete"}, workflow.ActivationRecoveryCompleteRequest{ActivationID: upgradeID, ExpectedVersion: recoveredActivation.Version, SelectedEventID: recoveredActivation.Recovery.SelectedEventID, ResourceVersions: resourceVersions}); err != nil {
		t.Fatalf("complete restore-prior: %v", err)
	}
	if actor, err := workflow.ResolveTrustedActor(workspace, issueMigrationsDir(), "actor-builder"); err != nil || actor.ProvenanceRef == "" {
		t.Fatalf("restored predecessor actor = (%#v, %v), want original active trust", actor, err)
	}
}

func currentFinalizedActivationForTest(t *testing.T, workspace string) workflow.Activation {
	t.Helper()
	activation, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), activationIDFromFinalizedJournal(t, workspace))
	if err != nil {
		t.Fatalf("read finalized activation: %v", err)
	}
	return activation
}

func activationIDFromFinalizedJournal(t *testing.T, workspace string) string {
	t.Helper()
	event := lastCommandJournalEvent(t, mustJournalPath(t, workspace))
	if event.EventType != "gate.registered" {
		t.Fatalf("last Event = %q, want gate registration after finalized activation", event.EventType)
	}
	// The registration Event carries the finalized activation ID in its Effect.
	for _, effect := range event.Effects {
		if effect.Table == "gate_registrations" {
			if id, ok := effect.Values["activation_id"].(string); ok {
				return id
			}
		}
	}
	t.Fatal("gate registration Event omitted activation ID")
	return ""
}

func prepareUpgradeActivationForTest(t *testing.T, workspace string, prior workflow.Activation) string {
	t.Helper()
	resources := make([]any, 0, len(prior.Resources))
	for _, resource := range prior.Resources {
		old := resource.NewIdentity
		resources = append(resources, map[string]any{"component": resource.Component, "old_identity": old, "new_identity": old + ":upgrade"})
	}
	inputPath := filepath.Join(workspace, "activation-upgrade.json")
	writeActivationInput(t, inputPath, map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "kit-upgrade", "phase": 21, "expected_prior_activation_id": prior.ID, "expected_prior_version": prior.Version, "resources": resources})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "prepare", "--input", inputPath, "--idempotency-key", "activation-upgrade-prepare")
	if result.err != nil {
		t.Fatalf("prepare upgrade: %v; stderr=%q", result.err, result.stderr)
	}
	return activationIDFromReplay(t, decodeActivationEnvelope(t, result.stdout))
}

func TestActivationLifecycleRejectsSkippedAndExtraTransitionsWithoutEvents(t *testing.T) {
	t.Parallel()
	for _, testCase := range []struct {
		name, status, operation, wantCode string
	}{
		{"prepared skips commit", "prepared", "commit", "invalid-argument"},
		{"prepared skips finalize", "prepared", "finalize", "invalid-argument"},
		{"staged skips finalize", "staged", "finalize", "invalid-argument"},
		{"committed stages extra component", "committed", "component", "invalid-argument"},
		{"committed repeats commit", "committed", "commit", "invalid-argument"},
		{"committed abort requires forward recovery", "committed", "abort", "activation-direction-conflict"},
		{"finalized stages extra component", "finalized", "component", "invalid-argument"},
		{"finalized repeats commit", "finalized", "commit", "invalid-argument"},
		{"finalized repeats finalize", "finalized", "finalize", "invalid-argument"},
		{"finalized abort requires forward recovery", "finalized", "abort", "activation-direction-conflict"},
		{"aborted stages extra component", "aborted", "component", "invalid-argument"},
		{"aborted skips commit", "aborted", "commit", "invalid-argument"},
		{"aborted skips finalize", "aborted", "finalize", "invalid-argument"},
		{"aborted repeats abort", "aborted", "abort", "invalid-argument"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := activationLifecycleFixtureAt(t, testCase.status)
			before := journalEventCount(t, fixture.workspace)
			result := rejectActivationLifecycleTransition(t, fixture, testCase.operation, testCase.name)
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"`+testCase.wantCode+`"`) {
				t.Fatalf("%s result = %#v, want %s rejection", testCase.name, result, testCase.wantCode)
			}
			if after := journalEventCount(t, fixture.workspace); after != before {
				t.Fatalf("%s Event count = %d, want no Event after %d", testCase.name, after, before)
			}
		})
	}
}

func TestActivationAbortPreCommitLeavesNoTrustedActor(t *testing.T) {
	t.Parallel()
	for _, status := range []string{"prepared", "staged"} {
		t.Run(status, func(t *testing.T) {
			fixture := activationLifecycleFixtureAt(t, status)
			inputPath := filepath.Join(fixture.workspace, "activation-abort-"+status+".json")
			writeActivationInput(t, inputPath, map[string]any{
				"activation_id": fixture.activationID, "expected_status": status, "expected_version": activationVersionForTest(t, fixture.workspace, fixture.activationID), "reason": "operator-abort",
			})
			before := journalEventCount(t, fixture.workspace)
			result := executeCommandForTest(fixture.workspace,
				"--project", fixture.workspace, "--json", "activation", "abort",
				"--input", inputPath, "--idempotency-key", "activation-abort-"+status,
			)
			if result.err != nil || result.stderr != "" || activationStatusFromResult(t, result.stdout) != "aborted" {
				t.Fatalf("abort from %s result = %#v, want aborted success", status, result)
			}
			if after := journalEventCount(t, fixture.workspace); after != before+1 {
				t.Fatalf("abort from %s Event count = %d, want one Event after %d", status, after, before)
			}
			if _, err := workflow.ResolveTrustedActor(fixture.workspace, issueMigrationsDir(), "actor-builder"); err == nil {
				t.Fatalf("abort from %s left actor-builder available as a trusted actor", status)
			}
		})
	}
}

type activationLifecycleFixture struct {
	workspace        string
	activationID     string
	componentDigests []string
}

func activationLifecycleFixtureAt(t *testing.T, status string) activationLifecycleFixture {
	t.Helper()
	workspace, values := cloneReusableWorkspaceFixture(t, "activation-lifecycle:"+status, func() (string, []string) {
		fixture := buildActivationLifecycleFixtureAt(t, status)
		return fixture.workspace, append([]string{fixture.activationID}, fixture.componentDigests...)
	})
	return activationLifecycleFixture{workspace: workspace, activationID: values[0], componentDigests: append([]string(nil), values[1:]...)}
}

func buildActivationLifecycleFixtureAt(t *testing.T, status string) activationLifecycleFixture {
	t.Helper()
	fixture := activationLifecycleFixture{workspace: clonePristineWorkspace(t)}
	fixture.activationID = prepareActivationForTest(t, fixture.workspace)
	fixture.componentDigests = []string{
		"sha256:" + strings.Repeat("0", 64),
		"sha256:" + strings.Repeat("1", 64),
	}
	if status == "prepared" {
		return fixture
	}
	if status == "aborted" {
		abortPath := filepath.Join(fixture.workspace, "activation-lifecycle-abort.json")
		writeActivationInput(t, abortPath, map[string]any{"activation_id": fixture.activationID, "expected_status": "prepared", "expected_version": activationVersionForTest(t, fixture.workspace, fixture.activationID), "reason": "fixture-abort"})
		if result := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "activation", "abort", "--input", abortPath, "--idempotency-key", "activation-lifecycle-abort"); result.err != nil {
			t.Fatalf("abort lifecycle fixture: %#v", result)
		}
		return fixture
	}
	profileDigest := stageActivationComponentForTest(t, fixture.workspace, fixture.activationID, "authority-profile", "lifecycle-builders", []any{map[string]any{
		"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/lifecycle-builders", "session_ref": nil,
		"capabilities": []string{"claim.acquire", "gate.register"},
	}})
	gateDigest := stageActivationComponentForTest(t, fixture.workspace, fixture.activationID, "gate-definitions", "lifecycle-gates", validGateDefinitions())
	fixture.componentDigests = []string{profileDigest, gateDigest}
	sort.Strings(fixture.componentDigests)
	if status == "staged" {
		return fixture
	}
	observeActivationResourcesForTest(t, fixture.workspace, fixture.activationID, "new")
	commitPath := filepath.Join(fixture.workspace, "activation-lifecycle-commit.json")
	writeActivationInput(t, commitPath, activationCommitInput(t, fixture.workspace, fixture.activationID, fixture.componentDigests))
	if result := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-lifecycle-commit"); result.err != nil {
		t.Fatalf("commit lifecycle fixture: %#v", result)
	}
	if status == "committed" {
		return fixture
	}
	if status == "finalized" {
		finalizePath := filepath.Join(fixture.workspace, "activation-lifecycle-finalize.json")
		writeActivationInput(t, finalizePath, activationFinalizeInput(t, fixture.workspace, fixture.activationID))
		if result := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "activation", "finalize", "--input", finalizePath, "--idempotency-key", "activation-lifecycle-finalize"); result.err != nil {
			t.Fatalf("finalize lifecycle fixture: %#v", result)
		}
		return fixture
	}
	t.Fatalf("unknown lifecycle fixture status %q", status)
	return activationLifecycleFixture{}
}

func TestActivationTransitionsEmitPublicEventTypes(t *testing.T) {
	t.Parallel()
	for _, testCase := range []struct {
		status, wantEventType string
	}{
		{status: "aborted", wantEventType: "activation.aborted"},
		{status: "committed", wantEventType: "activation.committed"},
		{status: "finalized", wantEventType: "activation.finalized"},
	} {
		t.Run(testCase.status, func(t *testing.T) {
			fixture := activationLifecycleFixtureAt(t, testCase.status)
			event := lastCommandJournalEvent(t, mustJournalPath(t, fixture.workspace))
			if event.EventType != testCase.wantEventType {
				t.Fatalf("activation %s Event type = %q, want %q", testCase.status, event.EventType, testCase.wantEventType)
			}
		})
	}
}

func rejectActivationLifecycleTransition(t *testing.T, fixture activationLifecycleFixture, operation, name string) commandExecutionResult {
	t.Helper()
	inputPath := filepath.Join(fixture.workspace, "activation-reject-"+strings.ReplaceAll(name, " ", "-")+".json")
	switch operation {
	case "component":
		writeActivationInput(t, inputPath, map[string]any{
			"activation_id": fixture.activationID, "component_type": "authority-profile", "component_id": "late-profile",
			"component_version": 2, "expected_prior_digest": nil, "expected_version": activationVersionForTest(t, fixture.workspace, fixture.activationID),
			"content": []any{map[string]any{
				"actor_id": "actor-late", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/late", "session_ref": nil,
				"capabilities": []string{"claim.acquire"},
			}},
		})
	case "commit":
		writeActivationInput(t, inputPath, activationCommitInput(t, fixture.workspace, fixture.activationID, fixture.componentDigests))
	case "abort":
		writeActivationInput(t, inputPath, map[string]any{"activation_id": fixture.activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, fixture.workspace, fixture.activationID), "reason": "reject-transition"})
	case "finalize":
		writeActivationInput(t, inputPath, activationFinalizeInput(t, fixture.workspace, fixture.activationID))
	default:
		t.Fatalf("unknown activation operation %q", operation)
	}
	return executeCommandForTest(fixture.workspace,
		"--project", fixture.workspace, "--json", "activation", operation,
		"--input", inputPath, "--idempotency-key", "activation-reject-"+strings.ReplaceAll(name, " ", "-"),
	)
}

func activationStatusFromResult(t *testing.T, output string) string {
	t.Helper()
	var envelope struct {
		Result struct {
			Activation struct {
				Status string `json:"status"`
			} `json:"activation"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(output), &envelope); err != nil || envelope.Result.Activation.Status == "" {
		t.Fatalf("decode activation mutation result: %v; stdout=%q", err, output)
	}
	return envelope.Result.Activation.Status
}

func TestActivationFinalizeAndStatusExposeOnlyFinalLifecycleState(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "builders", []any{map[string]any{
		"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/builders", "session_ref": nil, "capabilities": []string{"claim.acquire", "gate.register"},
	}})
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", validGateDefinitions())
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-commit.json")
	writeActivationInput(t, commitPath, activationCommitInput(t, workspace, activationID, digests))
	commit := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-finalize")
	if commit.err != nil {
		t.Fatalf("commit activation: %v; stderr=%q", commit.err, commit.stderr)
	}

	finalizePath := filepath.Join(workspace, "activation-finalize.json")
	writeActivationInput(t, finalizePath, activationFinalizeInput(t, workspace, activationID))
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
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	activationID := prepareActivationForTest(t, workspace)
	profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "builders", []any{map[string]any{
		"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/builders", "session_ref": nil, "capabilities": []string{"claim.acquire", "gate.register"},
	}})
	if _, err := workflow.ResolveTrustedActor(workspace, issueMigrationsDir(), "actor-builder"); err == nil {
		t.Fatal("staged actor resolved as trusted")
	}
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", validGateDefinitions())
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-commit-trust.json")
	writeActivationInput(t, commitPath, activationCommitInput(t, workspace, activationID, digests))
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-trust"); result.err != nil {
		t.Fatalf("commit: %v", result.err)
	}
	if _, err := workflow.ResolveTrustedActor(workspace, issueMigrationsDir(), "actor-builder"); err == nil {
		t.Fatal("committed actor resolved before finalize")
	}
	finalizePath := filepath.Join(workspace, "activation-finalize-trust.json")
	writeActivationInput(t, finalizePath, activationFinalizeInput(t, workspace, activationID))
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

func activationPrepareInput(kitIdentity string) map[string]any {
	return map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": kitIdentity, "phase": 21, "expected_prior_activation_id": nil, "expected_prior_version": nil, "resources": activationResources()}
}

func activationCommitInput(t *testing.T, workspace, activationID string, digests []string) map[string]any {
	t.Helper()
	return map[string]any{"activation_id": activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, workspace, activationID), "component_digests": digests}
}

func activationFinalizeInput(t *testing.T, workspace, activationID string) map[string]any {
	t.Helper()
	return map[string]any{"activation_id": activationID, "expected_status": "committed", "expected_version": activationVersionForTest(t, workspace, activationID)}
}

func activationResources() []any {
	return []any{
		map[string]any{"component": "kit", "old_identity": nil, "new_identity": "kit:new"},
		map[string]any{"component": "overlay", "old_identity": nil, "new_identity": "overlay:new"},
		map[string]any{"component": "provider", "old_identity": nil, "new_identity": "provider:new"},
		map[string]any{"component": "project", "old_identity": nil, "new_identity": "project:new"},
		map[string]any{"component": "backend", "old_identity": nil, "new_identity": "backend:new"},
	}
}

func stageActivationComponentForTest(t *testing.T, workspace, activationID, componentType, componentID string, content []any) string {
	t.Helper()
	inputPath := filepath.Join(workspace, componentID+".json")
	writeActivationInput(t, inputPath, map[string]any{
		"activation_id": activationID, "component_type": componentType, "component_id": componentID,
		"component_version": 1, "expected_prior_digest": nil, "expected_version": activationVersionForTest(t, workspace, activationID), "content": content,
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

func activationVersionForTest(t *testing.T, workspace, activationID string) int {
	t.Helper()
	activation, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), activationID)
	if err != nil {
		t.Fatalf("read activation version: %v", err)
	}
	return activation.Version
}

func observeActivationResourcesForTest(t *testing.T, workspace, activationID, state string) {
	t.Helper()
	activation, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), activationID)
	if err != nil {
		t.Fatalf("read prepared activation: %v", err)
	}
	for _, resource := range activation.Resources {
		observed := resource.NewIdentity
		if state == "old" && resource.OldIdentity != nil {
			observed = *resource.OldIdentity
		}
		identity := &observed
		if state == "old" && resource.OldIdentity == nil {
			identity = nil
		}
		_, err := workflow.ObserveActivationResource(workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), BootstrapActor: "test", IdempotencyKey: "observe-" + activationID + "-" + resource.Component}, workflow.ActivationResourceObservationRequest{ActivationID: activationID, Component: resource.Component, ExpectedActivationVersion: activation.Version, ExpectedVersion: resource.Version, ObservedIdentity: identity, State: state, EvidenceRef: "test:" + resource.Component})
		if err != nil {
			t.Fatalf("observe %s resource: %v", resource.Component, err)
		}
		activation, err = workflow.ActivationStatus(workspace, issueMigrationsDir(), activationID)
		if err != nil {
			t.Fatalf("reload activation after %s: %v", resource.Component, err)
		}
	}
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
