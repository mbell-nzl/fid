package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestGateRegisterUnknownActorReturnsContractAuthorityFailure(t *testing.T) {
	t.Parallel()
	assertGateRegisterAuthorityFailure(t, "actor-not-activated", "authority-unknown-actor")
}

func TestGateRegisterDeniedActorReturnsContractAuthorityFailure(t *testing.T) {
	t.Parallel()
	assertGateRegisterAuthorityFailure(t, "actor-no-gate", "authority-denied")
}

func assertGateRegisterAuthorityFailure(t *testing.T, actor, code string) {
	t.Helper()
	workspace := finalizedGateWorkspace(t, nil)
	inputPath := filepath.Join(workspace, "gate-register-authority.json")
	writeActivationInput(t, inputPath, gateRegisterInput(t, workspace))
	before := journalEventCount(t, workspace)
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "register", "--actor", actor, "--input", inputPath, "--idempotency-key", "gate-register-authority-"+actor)
	if result.err == nil || cli.ExitCode(result.err) != 5 || result.stdout != "" {
		t.Fatalf("gate register result = %+v, want exit 5 with stderr only", result)
	}
	var failure struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Event     any    `json:"event"`
		Error     struct {
			Code     string         `json:"code"`
			Category string         `json:"category"`
			Subject  any            `json:"subject"`
			Details  map[string]any `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil {
		t.Fatalf("decode gate failure %q: %v", result.stderr, err)
	}
	wantDetails := map[string]any{"operation": "gate.register", "actor_id": actor,
		"required_actions": []any{"gate.register"}, "required_capabilities": []any{"gate.register"}}
	if failure.OK || failure.Operation != "gate.register" || failure.Event != nil || failure.Error.Code != code || failure.Error.Category != "authority" || failure.Error.Subject != nil || !reflect.DeepEqual(failure.Error.Details, wantDetails) {
		t.Fatalf("gate failure = %+v, want %s with details %#v and no Event", failure, code, wantDetails)
	}
	if after := journalEventCount(t, workspace); after != before {
		t.Fatalf("rejected gate register changed Event count: before=%d after=%d", before, after)
	}
}

func TestGateRegisterRequiresFinalizedActivatedDefinition(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	initializeObservedRepositoryForTest(t, workspace)
	activationID := prepareActivationForTest(t, workspace)
	profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "builders", []any{map[string]any{
		"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/builders", "session_ref": nil,
		"capabilities": []string{"gate.register"},
	}})
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", validGateDefinitions())
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	inputPath := filepath.Join(workspace, "gate-register.json")
	writeActivationInput(t, inputPath, gateRegisterInput(t, workspace))

	staged := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "register", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "gate-register-staged")
	if staged.err == nil || staged.stdout != "" || !strings.Contains(staged.stderr, `"code":"authority-unavailable"`) {
		t.Fatalf("staged registration stdout=%q stderr=%q err=%v, want finalized-activation rejection", staged.stdout, staged.stderr, staged.err)
	}

	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-commit-gate.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, workspace, activationID), "component_digests": digests})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-gate"); result.err != nil {
		t.Fatalf("commit activation: %v; stderr=%q", result.err, result.stderr)
	}
	finalizePath := filepath.Join(workspace, "activation-finalize-gate.json")
	writeActivationInput(t, finalizePath, map[string]any{"activation_id": activationID, "expected_status": "committed", "expected_version": activationVersionForTest(t, workspace, activationID)})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "finalize", "--input", finalizePath, "--idempotency-key", "activation-finalize-gate"); result.err != nil {
		t.Fatalf("finalize activation: %v; stderr=%q", result.err, result.stderr)
	}

	registered := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "register", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "gate-register-finalized")
	if registered.err != nil {
		t.Fatalf("finalized registration: %v; stderr=%q", registered.err, registered.stderr)
	}
	var envelope struct {
		Result struct {
			Registration struct {
				ID                 string `json:"id"`
				Version            int    `json:"version"`
				ActivationID       string `json:"activation_id"`
				DefinitionID       string `json:"definition_id"`
				DefinitionVersion  int    `json:"definition_version"`
				DefinitionDigest   string `json:"definition_digest"`
				BindingDigest      string `json:"binding_digest"`
				AuthorizationProof string `json:"authorization_provenance"`
			} `json:"registration"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(registered.stdout), &envelope); err != nil {
		t.Fatalf("decode registration: %v; stdout=%q", err, registered.stdout)
	}
	registration := envelope.Result.Registration
	if registration.ID == "" || registration.Version != 1 || registration.ActivationID != activationID || registration.DefinitionID != "gate-unit" || registration.DefinitionVersion != 1 || registration.DefinitionDigest == "" || registration.BindingDigest == "" || registration.AuthorizationProof != "activation:profile/builders" {
		t.Fatalf("registration = %+v, want resolved finalized definition", registration)
	}
}

func TestGateRegisterRejectsUnsupportedPolicyWithoutPersistence(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name  string
		apply func(map[string]any)
		field string
	}{
		{"argv", func(definition map[string]any) { definition["argv"] = []string{} }, "argv"},
		{"shell", func(definition map[string]any) { definition["shell"] = true }, "shell"},
		{"sandbox", func(definition map[string]any) { definition["sandbox"] = "isolated" }, "sandbox"},
		{"network", func(definition map[string]any) { definition["network"] = "blocked" }, "network"},
		{"write", func(definition map[string]any) { definition["write"] = "denied" }, "write"},
		{"environment", func(definition map[string]any) {
			definition["environment"] = []any{map[string]any{"name": "HOME", "source": "ambient", "value": nil}}
		}, "environment"},
		{"timeout", func(definition map[string]any) { definition["timeout_ms"] = 0 }, "timeout_ms"},
		{"output", func(definition map[string]any) { definition["output_bytes"] = 0 }, "output_bytes"},
		{"excerpt", func(definition map[string]any) { definition["excerpt_bytes"] = 1025 }, "excerpt_bytes"},
		{"success codes", func(definition map[string]any) { definition["success_exit_codes"] = []int{} }, "success_exit_codes"},
		{"credential names", func(definition map[string]any) { definition["credential_names"] = []string{"B", "A"} }, "credential_names"},
		{"redact literals", func(definition map[string]any) { definition["redact_literals"] = []string{"secret", "secret"} }, "redact_literals"},
		{"artifact bytes", func(definition map[string]any) { definition["artifact_bytes"] = 0 }, "artifact_bytes"},
		{"artifact path", func(definition map[string]any) { definition["artifact_paths"] = []string{"../outside"} }, "artifact_paths"},
		{"cwd", func(definition map[string]any) { definition["cwd"] = "../outside" }, "cwd"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			workspace := finalizedGateWorkspace(t, nil)
			// Stage valid current facts first. This test exercises registration's
			// defensive rejection of a corrupted stored policy, not staging.
			definition := gateDefinitionsForTest(nil)[0].(map[string]any)
			testCase.apply(definition)
			raw, err := json.Marshal(definition)
			if err != nil {
				t.Fatal(err)
			}
			conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
			if err != nil {
				t.Fatal(err)
			}
			_, updateErr := conn.Exec(`UPDATE activation_gate_definitions SET definition_json = ? WHERE definition_id = 'gate-unit'`, string(raw))
			closeErr := conn.Close()
			if updateErr != nil || closeErr != nil {
				t.Fatalf("corrupt policy fixture: update=%v close=%v", updateErr, closeErr)
			}
			inputPath := filepath.Join(workspace, "gate-register.json")
			request := gateRegisterInput(t, workspace)
			writeActivationInput(t, inputPath, request)
			result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "register", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "gate-register-policy-"+strings.ReplaceAll(testCase.name, " ", "-"))
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"executor-policy-unsupported"`) || !strings.Contains(result.stderr, `"fields":["`+testCase.field+`"]`) {
				t.Fatalf("policy %s stdout=%q stderr=%q err=%v, want typed %q rejection", testCase.name, result.stdout, result.stderr, result.err, testCase.field)
			}
			if count := registrationCount(t, workspace); count != 0 {
				t.Fatalf("policy %s persisted %d registrations, want 0", testCase.name, count)
			}
		})
	}
}

func TestGateRegisterReusesBindingAndVersionsChangedBinding(t *testing.T) {
	t.Parallel()
	workspace := finalizedGateWorkspace(t, func(map[string]any) {})
	inputPath := filepath.Join(workspace, "gate-register.json")
	writeActivationInput(t, inputPath, gateRegisterInput(t, workspace))

	first := registerGateForTest(t, workspace, inputPath, "gate-register-first")
	second := registerGateForTest(t, workspace, inputPath, "gate-register-same-binding")
	if first.ID == "" || first.Version != 1 || second.ID != first.ID || second.Version != first.Version || registrationCount(t, workspace) != 1 {
		t.Fatalf("same binding first=%+v second=%+v count=%d, want one immutable registration", first, second, registrationCount(t, workspace))
	}

	changed := gateRegisterInput(t, workspace)
	changed["scope"] = map[string]any{"issue_id": "ISS-2"}
	writeActivationInput(t, inputPath, changed)
	third := registerGateForTest(t, workspace, inputPath, "gate-register-changed-binding")
	if third.ID != first.ID || third.Version != 2 || third.BindingDigest == first.BindingDigest || registrationCount(t, workspace) != 2 {
		t.Fatalf("changed binding first=%+v third=%+v count=%d, want immutable version 2", first, third, registrationCount(t, workspace))
	}

	resolved := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "resolve", first.ID, "--version", "1")
	if resolved.err != nil || !strings.Contains(resolved.stdout, `"version":1`) || !strings.Contains(resolved.stdout, `"definition_id":"gate-unit"`) {
		t.Fatalf("resolve registration stdout=%q stderr=%q err=%v, want immutable first version", resolved.stdout, resolved.stderr, resolved.err)
	}
}

func TestGateRegisterObservesRepositoryWithoutCallerIdentity(t *testing.T) {
	t.Parallel()
	workspace := finalizedGateWorkspace(t, nil)
	input := gateRegisterInput(t, workspace)
	delete(input, "worktree_id")
	delete(input, "repository")
	path := filepath.Join(workspace, "gate-observed.json")
	writeActivationInput(t, path, input)
	registered := registerGateForTest(t, workspace, path, "gate-observed")
	observed := testEnvelopeRepository(t, workspace)
	if registered.WorktreeID != observed["worktree_id"] {
		t.Fatalf("registration worktree = %q, want observed %q", registered.WorktreeID, observed["worktree_id"])
	}
	independent := p24IndependentBindingRepository(t, workspace)
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"repository": independent, "scope": input["scope"], "worktree_id": independent["worktree_id"]})
	if err != nil {
		t.Fatal(err)
	}
	if registered.BindingDigest != fmt.Sprintf("sha256:%x", sha256.Sum256(payload)) {
		t.Fatalf("binding digest = %q, want independent observation and scope", registered.BindingDigest)
	}
	var response struct {
		Result struct {
			Registration struct {
				ProjectID string `json:"project_id"`
			} `json:"registration"`
		} `json:"result"`
	}
	get := executeCommandForTest(workspace, "--json", "gate", "resolve", registered.ID, "--version", "1")
	if get.err != nil || json.Unmarshal([]byte(get.stdout), &response) != nil || response.Result.Registration.ProjectID != independent["project_id"] {
		t.Fatalf("Registration project: %#v", get)
	}
	p24AssertBindingEventRepository(t, workspace, independent)
}

func TestGateRegisterRejectsRuntimePolicyOverride(t *testing.T) {
	t.Parallel()
	workspace := finalizedGateWorkspace(t, func(map[string]any) {})
	inputPath := filepath.Join(workspace, "gate-register.json")
	request := gateRegisterInput(t, workspace)
	request["argv"] = []string{"sh", "-c", "override"}
	writeActivationInput(t, inputPath, request)

	result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "register", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "gate-register-runtime-override")
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) || registrationCount(t, workspace) != 0 {
		t.Fatalf("runtime policy override stdout=%q stderr=%q count=%d err=%v, want closed-input rejection without persistence", result.stdout, result.stderr, registrationCount(t, workspace), result.err)
	}
}

type gateRegistrationView struct {
	ID            string `json:"id"`
	Version       int    `json:"version"`
	BindingDigest string `json:"binding_digest"`
	WorktreeID    string `json:"worktree_id"`
}

func registerGateForTest(t *testing.T, workspace, inputPath, key string) gateRegistrationView {
	t.Helper()
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "register", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", key)
	if result.err != nil {
		t.Fatalf("register gate: %v; stderr=%q", result.err, result.stderr)
	}
	var envelope struct {
		Result struct {
			Registration gateRegistrationView `json:"registration"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("decode registration: %v; stdout=%q", err, result.stdout)
	}
	return envelope.Result.Registration
}

func finalizedGateWorkspace(t *testing.T, modifyDefinition func(map[string]any)) string {
	t.Helper()
	return finalizedGateWorkspaceWithDefinitions(t, gateDefinitionsForTest(modifyDefinition))
}

func finalizedGateWorkspaceWithDefinitions(t *testing.T, definitions []any) string {
	t.Helper()
	key := reusableWorkspaceFixtureKey(t, "finalized-gate", definitions)
	workspace, _ := cloneReusableWorkspaceFixture(t, key, func() (string, []string) {
		return buildFinalizedGateWorkspace(t, definitions), nil
	})
	initializeObservedRepositoryForTest(t, workspace)
	return workspace
}

func buildFinalizedGateWorkspace(t *testing.T, definitions []any) string {
	t.Helper()
	workspace, values := cloneReusableWorkspaceFixture(t, "prepared-gate", func() (string, []string) {
		workspace := clonePristineWorkspace(t)
		activationID := prepareActivationForTest(t, workspace)
		profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "builders", []any{map[string]any{
			"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/builders", "session_ref": nil,
			"capabilities": []string{"gate.register"},
		}, map[string]any{
			"actor_id": "actor-no-gate", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/builders", "session_ref": nil,
			"capabilities": []string{"artifact.verify"},
		}})
		return workspace, []string{activationID, profileDigest}
	})
	activationID, profileDigest := values[0], values[1]
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", definitions)
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-commit-gate.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, workspace, activationID), "component_digests": digests})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-gate"); result.err != nil {
		t.Fatalf("commit activation: %v; stderr=%q", result.err, result.stderr)
	}
	finalizePath := filepath.Join(workspace, "activation-finalize-gate.json")
	writeActivationInput(t, finalizePath, map[string]any{"activation_id": activationID, "expected_status": "committed", "expected_version": activationVersionForTest(t, workspace, activationID)})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "finalize", "--input", finalizePath, "--idempotency-key", "activation-finalize-gate"); result.err != nil {
		t.Fatalf("finalize activation: %v; stderr=%q", result.err, result.stderr)
	}
	return workspace
}

func registrationCount(t *testing.T, workspace string) int {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open registration database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM gate_registrations`).Scan(&count); err != nil && err != sql.ErrNoRows {
		t.Fatalf("count registrations: %v", err)
	}
	return count
}

func gateRegisterInput(t *testing.T, workspace string) map[string]any {
	t.Helper()
	return map[string]any{
		"definition_id": "gate-unit", "definition_version": 1,
		"scope": map[string]any{"issue_id": "ISS-1"},
	}
}
