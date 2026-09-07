//go:build integration

package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func p22NewFIDActivationAdapter(t *testing.T, status string) p22RecoveryAdapter {
	t.Helper()
	binary, workspace := integrationFIDBinary(t), t.TempDir()
	p21Run(t, binary, workspace, "--json", "init", "--prefix", "ACT")
	id := p21Activation(t, binary, workspace, "prepare", "activation-corpus-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p22-activation", "phase": 21})
	if status == "committed" {
		profile := p21Component(t, binary, workspace, id, "authority-profile", "activation-corpus-profile", p22ActivationProfile())
		gates := p21Component(t, binary, workspace, id, "gate-definitions", "activation-corpus-gates", p22ActivationGates())
		digests := []string{profile, gates}
		sort.Strings(digests)
		p21Activation(t, binary, workspace, "commit", "activation-corpus-commit", map[string]any{"activation_id": id, "expected_status": "staged", "component_digests": digests})
	}
	response := p21Run(t, binary, workspace, "--json", "activation", "status", id)
	return &p22FIDRecoveryAdapter{binary: binary, workspace: workspace, actor: "bootstrap:corpus", activation: p21Path(t, response, "result", "activation").(map[string]any)}
}

func p22NewFakeActivationAdapter(status string) p22RecoveryAdapter {
	resources := make([]any, 0, 5)
	for _, input := range p21ActivationResources() {
		resource := p22CloneMap(input.(map[string]any))
		resource["version"], resource["state"] = 1, "pending"
		resource["observed_identity"], resource["evidence_ref"], resource["error_code"] = nil, nil, nil
		if status == "committed" {
			resource["version"], resource["state"], resource["observed_identity"], resource["evidence_ref"] = 2, "new", resource["new_identity"], "setup:observed"
		}
		resources = append(resources, resource)
	}
	digests := []any{}
	version := 1
	if status == "committed" {
		version = 9
		values := []string{p22ActivationContentDigest(p22ActivationProfile()), p22ActivationContentDigest(p22ActivationGates())}
		sort.Strings(values)
		for _, value := range values {
			digests = append(digests, value)
		}
	}
	activation := map[string]any{"id": "ACT-fake", "status": status, "version": version, "kit_identity": "p22-activation", "backend_identity": "fid", "phase": 21, "direction": "activate", "point_of_no_return": status == "committed", "prior_activation_id": nil, "component_digests": digests, "resources": resources, "recovery": nil, "created_at": "2026-09-07T00:00:00Z", "updated_at": "2026-09-07T00:00:00Z"}
	return &p22FakeRecoveryAdapter{context: &p22FIDRecoveryAdapter{actor: "bootstrap:corpus", activation: p22CloneMap(activation)}, model: &p22MemoryModel{activation: activation, replies: map[string]map[string]any{}, requests: map[string]string{}}}
}

func p22ActivationContentDigest(content any) string {
	encoded, _ := json.Marshal(content)
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func p22ActivationProfile() []any {
	return []any{map[string]any{"actor_id": "actor-activation", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:corpus", "session_ref": nil, "capabilities": []string{"claim.acquire"}}}
}

func p22ActivationGates() []any {
	return []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-activation", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}}
}

func (adapter *p22FIDRecoveryAdapter) p22ActivationArgs(t *testing.T, fixture p22ConformanceCase) ([]string, map[string]any) {
	t.Helper()
	a := adapter.activation
	leaf := strings.TrimPrefix(fixture.Operation, "activation.")
	args := []string{"--project", adapter.workspace, "--json", "activation", leaf}
	if leaf == "status" {
		args = append(args, a["id"].(string))
		if fixture.Kind == "unknown-field" {
			args = append(args, "--unknown-field")
		}
		return args, nil
	}
	args = append(args, "--actor", adapter.actor)
	input := map[string]any{"activation_id": a["id"], "expected_version": a["version"]}
	switch leaf {
	case "resource":
		component := "backend"
		if value, ok := fixture.Input["component"].(string); ok {
			component = value
		}
		resource := p22ActivationResource(a, component)
		if resource == nil {
			t.Fatalf("request builder has no resource %s", component)
		}
		state := "new"
		if value, ok := fixture.Input["state"].(string); ok {
			state = value
		}
		identity := resource["new_identity"]
		if state == "old" {
			identity = resource["old_identity"]
		}
		input = map[string]any{"activation_id": a["id"], "component": component, "expected_activation_version": a["version"], "expected_version": resource["version"], "observed_identity": identity, "state": state, "evidence_ref": "corpus:" + component, "error_code": nil}
	case "recover":
		input["expected_status"], input["direction"], input["reason"] = a["status"], p22ActivationDirection(a), "recover activation corpus"
	case "recover-complete":
		input["selected_event_id"] = "EVT-not-selected"
		if recovery, ok := a["recovery"].(map[string]any); ok {
			input["selected_event_id"] = recovery["selected_event_id"]
		}
		versions := []any{}
		for _, item := range a["resources"].([]any) {
			r := item.(map[string]any)
			versions = append(versions, map[string]any{"component": r["component"], "version": r["version"]})
		}
		input["resource_versions"] = versions
	default:
		t.Fatalf("unsupported activation corpus leaf %q", leaf)
	}
	if fixture.Kind == "unknown-field" {
		input["unknown_field"] = true
	}
	return args, input
}

func p22ActivationResource(activation map[string]any, component string) map[string]any {
	for _, item := range activation["resources"].([]any) {
		r := item.(map[string]any)
		if r["component"] == component {
			return r
		}
	}
	return nil
}

func p22ActivationDirection(activation map[string]any) string {
	if activation["point_of_no_return"] == true {
		return "complete-forward"
	}
	return "restore-prior"
}

func (model *p22MemoryModel) executeActivation(t *testing.T, operation, key, actor string, input map[string]any, invalidFlag bool) p22ConformanceResponse {
	t.Helper()
	failure := func(code string, details map[string]any) p22ConformanceResponse {
		return model.activationFailure(t, operation, code, details)
	}
	if invalidFlag {
		return failure("invalid-argument", nil)
	}
	allowed := map[string]string{"activation.resource": "activation_id component expected_activation_version expected_version observed_identity state evidence_ref error_code", "activation.recover": "activation_id expected_version expected_status direction reason", "activation.recover-complete": "activation_id expected_version selected_event_id resource_versions", "activation.status": ""}
	for field := range input {
		if !strings.Contains(" "+allowed[operation]+" ", " "+field+" ") {
			return failure("invalid-argument", nil)
		}
	}
	encoded, _ := json.Marshal(input)
	request := actor + string(encoded)
	if cached, ok := model.replies[key]; ok {
		if model.requests[key] != request {
			return failure("idempotency-conflict", nil)
		}
		replay := p22CloneMap(cached)
		replay["result"].(map[string]any)["replayed"] = true
		return p22FixtureResponse(t, replay, true, 0)
	}
	if operation == "activation.status" {
		return p22FixtureResponse(t, map[string]any{"contract_version": "fid-cli.v1", "ok": true, "operation": operation, "event": nil, "result": map[string]any{"activation": p22CloneMap(model.activation)}}, true, 0)
	}
	a := p22CloneMap(model.activation)
	versionField := "expected_version"
	if operation == "activation.resource" {
		versionField = "expected_activation_version"
	}
	if p22Number(input[versionField]) != p22Number(a["version"]) {
		return failure("activation-version-conflict", map[string]any{"activation_id": a["id"], "expected_version": input[versionField], "actual_version": a["version"]})
	}
	eventID := fmt.Sprintf("EVT-activation-%06d", model.sequence+1)
	var resource map[string]any
	switch operation {
	case "activation.resource":
		resource = p22ActivationResource(a, input["component"].(string))
		if resource == nil {
			return failure("invalid-argument", nil)
		}
		if p22Number(input["expected_version"]) != p22Number(resource["version"]) {
			return failure("activation-version-conflict", map[string]any{"activation_id": a["id"], "expected_version": input["expected_version"], "actual_version": resource["version"]})
		}
		if input["state"] != "old" && input["state"] != "new" {
			return failure("invalid-argument", nil)
		}
		identity := resource["new_identity"]
		if input["state"] == "old" {
			identity = resource["old_identity"]
		}
		if input["observed_identity"] != identity {
			return failure("invalid-argument", nil)
		}
		resource["version"] = p22Number(resource["version"]) + 1
		for _, field := range []string{"observed_identity", "state", "evidence_ref", "error_code"} {
			resource[field] = input[field]
		}
	case "activation.recover":
		if input["direction"] != p22ActivationDirection(a) {
			return failure("activation-direction-conflict", map[string]any{"activation_id": a["id"], "point_of_no_return": a["point_of_no_return"], "required_direction": p22ActivationDirection(a)})
		}
		a["recovery"] = map[string]any{"direction": input["direction"], "actor_attribution": actor, "reason": input["reason"], "selected_event_id": eventID, "completed_event_id": nil}
	case "activation.recover-complete":
		recovery, ok := a["recovery"].(map[string]any)
		if !ok || recovery["selected_event_id"] != input["selected_event_id"] {
			return failure("activation-incomplete", map[string]any{"activation_id": a["id"], "components": []any{}, "reasons": []any{}})
		}
		versions, ok := input["resource_versions"].([]any)
		if !ok || len(versions) != 5 {
			return failure("invalid-argument", nil)
		}
		seen := map[string]bool{}
		for _, item := range versions {
			v := item.(map[string]any)
			component, _ := v["component"].(string)
			r := p22ActivationResource(a, component)
			if r == nil || seen[component] {
				return failure("invalid-argument", nil)
			}
			seen[component] = true
			if p22Number(v["version"]) != p22Number(r["version"]) {
				return failure("activation-version-conflict", map[string]any{"activation_id": a["id"], "expected_version": v["version"], "actual_version": r["version"]})
			}
		}
		state := "old"
		if recovery["direction"] == "complete-forward" {
			state = "new"
		}
		missing := []any{}
		for _, item := range a["resources"].([]any) {
			r := item.(map[string]any)
			if r["state"] != state {
				missing = append(missing, r["component"])
			}
		}
		if len(missing) > 0 {
			reasons := []any{}
			for _, component := range missing {
				reasons = append(reasons, map[string]any{"code": "activation-incomplete", "subject_type": "activation", "subject_id": a["id"], "detail": component})
			}
			return failure("activation-incomplete", map[string]any{"activation_id": a["id"], "components": missing, "reasons": reasons})
		}
		recovery["completed_event_id"] = eventID
		a["status"] = "aborted"
		if state == "new" {
			a["status"] = "finalized"
		}
	default:
		return failure("invalid-argument", nil)
	}
	model.sequence++
	a["version"] = p22Number(a["version"]) + 1
	a["updated_at"] = fmt.Sprintf("2026-09-07T00:00:%02dZ", model.sequence)
	model.activation = a
	result := map[string]any{"activation": p22CloneMap(a), "idempotency_key": key, "replayed": false}
	if resource != nil {
		result["resource"] = p22CloneMap(resource)
	}
	response := map[string]any{"contract_version": "fid-cli.v1", "ok": true, "operation": operation, "event": map[string]any{"id": eventID, "occurred_at": a["updated_at"]}, "result": result}
	model.requests[key], model.replies[key] = request, p22CloneMap(response)
	return p22FixtureResponse(t, response, true, 0)
}

func (model *p22MemoryModel) activationFailure(t *testing.T, operation, code string, details map[string]any) p22ConformanceResponse {
	category, exit := "precondition", 4
	if code == "activation-version-conflict" || code == "idempotency-conflict" {
		category = "conflict"
	}
	var subject any = map[string]any{"type": "activation", "id": model.activation["id"]}
	if code == "invalid-argument" {
		category, exit, subject = "invalid-input", 2, nil
		details = map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted activation input"}}
	}
	return p22FixtureResponse(t, map[string]any{"contract_version": "fid-cli.v1", "ok": false, "operation": operation, "error": map[string]any{"code": code, "category": category, "retryable": false, "message": "activation request rejected", "subject": subject, "details": details}}, false, exit)
}
