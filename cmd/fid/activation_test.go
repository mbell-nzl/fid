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
	t.Parallel()
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

func TestActivationReplaysAcrossLifecycleWithoutPreemptingLifecycleChecks(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	preparePath := filepath.Join(workspace, "activation-replay-prepare.json")
	writeActivationInput(t, preparePath, map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "kit", "phase": 21})
	prepare := assertActivationReplay(t, workspace, "prepare", preparePath, "activation-replay-prepare")
	activationID := activationIDFromReplay(t, prepare)

	newKey := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "prepare", "--input", preparePath, "--idempotency-key", "activation-replay-new-key")
	assertActivationTypedFailure(t, newKey, "activation-conflict")

	profilePath := filepath.Join(workspace, "activation-replay-profile.json")
	writeActivationInput(t, profilePath, map[string]any{
		"activation_id": activationID, "component_type": "authority-profile", "component_id": "replay-builders",
		"component_version": 1, "expected_prior_digest": nil,
		"content": []any{map[string]any{
			"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/replay-builders", "session_ref": nil,
			"capabilities": []string{"claim.acquire", "gate.register"},
		}},
	})
	profile := assertActivationReplay(t, workspace, "component", profilePath, "activation-replay-profile")
	profileDigest := componentDigestFromReplay(t, profile)

	writeActivationInput(t, profilePath, map[string]any{
		"activation_id": activationID, "component_type": "authority-profile", "component_id": "replay-builders",
		"component_version": 1, "expected_prior_digest": nil,
		"content": []any{map[string]any{
			"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/replay-builders-v2", "session_ref": nil,
			"capabilities": []string{"claim.acquire", "gate.register"},
		}},
	})
	assertActivationStableDigestConflict(t, workspace, "component", profilePath, "activation-replay-profile")

	gatePath := filepath.Join(workspace, "activation-replay-gates.json")
	writeActivationInput(t, gatePath, map[string]any{
		"activation_id": activationID, "component_type": "gate-definitions", "component_id": "replay-gates",
		"component_version": 1, "expected_prior_digest": nil, "content": validGateDefinitions(),
	})
	gate := assertActivationReplay(t, workspace, "component", gatePath, "activation-replay-gates")
	gateDigest := componentDigestFromReplay(t, gate)

	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-replay-commit.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	assertActivationReplay(t, workspace, "commit", commitPath, "activation-replay-commit")

	finalizePath := filepath.Join(workspace, "activation-replay-finalize.json")
	writeActivationInput(t, finalizePath, map[string]any{"activation_id": activationID, "expected_status": "committed"})
	assertActivationReplay(t, workspace, "finalize", finalizePath, "activation-replay-finalize")

	abortWorkspace := clonePristineWorkspace(t)
	abortPreparePath := filepath.Join(abortWorkspace, "activation-replay-abort-prepare.json")
	writeActivationInput(t, abortPreparePath, map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "kit", "phase": 21})
	abortPrepare := assertActivationReplay(t, abortWorkspace, "prepare", abortPreparePath, "activation-replay-abort-prepare")
	abortPath := filepath.Join(abortWorkspace, "activation-replay-abort.json")
	writeActivationInput(t, abortPath, map[string]any{"activation_id": activationIDFromReplay(t, abortPrepare), "expected_status": "prepared", "reason": "replay"})
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
				"component_version": 2, "expected_prior_digest": nil, "content": []any{profile},
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
		{"committed abort remains unavailable", "committed", "abort", "activation-recovery-unavailable"},
		{"finalized stages extra component", "finalized", "component", "invalid-argument"},
		{"finalized repeats commit", "finalized", "commit", "invalid-argument"},
		{"finalized repeats finalize", "finalized", "finalize", "invalid-argument"},
		{"finalized abort remains unavailable", "finalized", "abort", "activation-recovery-unavailable"},
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
				"activation_id": fixture.activationID, "expected_status": status, "reason": "operator-abort",
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
		writeActivationInput(t, abortPath, map[string]any{"activation_id": fixture.activationID, "expected_status": "prepared", "reason": "fixture-abort"})
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
	commitPath := filepath.Join(fixture.workspace, "activation-lifecycle-commit.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": fixture.activationID, "expected_status": "staged", "component_digests": fixture.componentDigests})
	if result := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-lifecycle-commit"); result.err != nil {
		t.Fatalf("commit lifecycle fixture: %#v", result)
	}
	if status == "committed" {
		return fixture
	}
	if status == "finalized" {
		finalizePath := filepath.Join(fixture.workspace, "activation-lifecycle-finalize.json")
		writeActivationInput(t, finalizePath, map[string]any{"activation_id": fixture.activationID, "expected_status": "committed"})
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
			"component_version": 2, "expected_prior_digest": nil,
			"content": []any{map[string]any{
				"actor_id": "actor-late", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/late", "session_ref": nil,
				"capabilities": []string{"claim.acquire"},
			}},
		})
	case "commit":
		writeActivationInput(t, inputPath, map[string]any{"activation_id": fixture.activationID, "expected_status": "staged", "component_digests": fixture.componentDigests})
	case "abort":
		writeActivationInput(t, inputPath, map[string]any{"activation_id": fixture.activationID, "expected_status": "staged", "reason": "reject-transition"})
	case "finalize":
		writeActivationInput(t, inputPath, map[string]any{"activation_id": fixture.activationID, "expected_status": "committed"})
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
