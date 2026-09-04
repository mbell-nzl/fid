package main

import (
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

const testEvidenceObligationID = "OBL-1"

func TestEnvelopeCreateGetActivateRoundTripAndRejectsWidening(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	inputPath := filepath.Join(workspace, "envelope.json")
	input := envelopeCreateInput(t, workspace, issueID, registration)
	writeActivationInput(t, inputPath, input)

	created := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "envelope-create-1",
	)
	if created.err != nil || created.stderr != "" {
		t.Fatalf("envelope create = %#v, want stdout-only success", created)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")

	got := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "get", envelopeID)
	if got.err != nil || got.stderr != "" || assertEnvelopeShape(t, got.stdout, "draft") != envelopeID {
		t.Fatalf("envelope get = %#v, want exact draft envelope", got)
	}

	activatePath := filepath.Join(workspace, "envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder",
		"--input", activatePath, "--idempotency-key", "envelope-activate-1",
	)
	if activated.err != nil || activated.stderr != "" || assertEnvelopeShape(t, activated.stdout, "active") != envelopeID {
		t.Fatalf("envelope activate = %#v, want active exact envelope", activated)
	}

	for _, testCase := range []struct {
		name  string
		apply func(map[string]any)
	}{
		{"related issues", func(value map[string]any) { value["related_issue_ids"] = []string{"ISS-2"} }},
		{"multiple obligations", func(value map[string]any) {
			value["obligations"] = append(value["obligations"].([]any), value["obligations"].([]any)[0])
		}},
		{"computed status", func(value map[string]any) {
			value["obligations"].([]any)[0].(map[string]any)["derived_status"] = "satisfied"
		}},
		{"other obligation type", func(value map[string]any) { value["obligations"].([]any)[0].(map[string]any)["type"] = "human-outcome" }},
		{"guard widening", func(value map[string]any) {
			value["obligations"].([]any)[0].(map[string]any)["guarded_actions"] = []string{"envelope.complete", "issue.transition-terminal-success", "run.finish"}
		}},
		{"mutable budget", func(value map[string]any) { value["budgets"].([]any)[0].(map[string]any)["consumed"] = 1 }},
		{"routine lane", func(value map[string]any) { value["routine"] = map[string]any{"classification": "external"} }},
		{"revision", func(value map[string]any) { value["revision"] = 2 }},
		{"cancellation", func(value map[string]any) { value["cancellation"] = map[string]any{} }},
		{"protected decision", func(value map[string]any) { value["protected_decision"] = "DEC-1" }},
		{"stop", func(value map[string]any) { value["stop"] = map[string]any{} }},
		{"attention", func(value map[string]any) { value["attention"] = map[string]any{} }},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			candidate := cloneEnvelopeInput(t, input)
			testCase.apply(candidate)
			writeActivationInput(t, inputPath, candidate)
			before := journalEventCount(t, workspace)
			result := executeCommandForTest(workspace,
				"--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder",
				"--input", inputPath, "--idempotency-key", "envelope-reject-"+strings.ReplaceAll(testCase.name, " ", "-"),
			)
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
				t.Fatalf("%s result = %#v, want typed closed-input rejection", testCase.name, result)
			}
			if after := journalEventCount(t, workspace); after != before {
				t.Fatalf("%s event count = %d, want no Event after %d", testCase.name, after, before)
			}
		})
	}
}

func finalizedEnvelopeWorkspace(t *testing.T) (string, string, gateRegistrationView) {
	return finalizedEnvelopeWorkspaceWithDefinition(t, nil)
}

func finalizedEnvelopeWorkspaceWithDefinition(t *testing.T, modifyDefinition func(map[string]any)) (string, string, gateRegistrationView) {
	t.Helper()
	return finalizedEnvelopeWorkspaceWithDefinitions(t, gateDefinitionsForTest(modifyDefinition))
}

func finalizedEnvelopeWorkspaceWithDefinitions(t *testing.T, definitions []any) (string, string, gateRegistrationView) {
	t.Helper()
	key := reusableWorkspaceFixtureKey(t, "finalized-envelope", definitions)
	workspace, values := cloneReusableWorkspaceFixture(t, key, func() (string, []string) {
		workspace, issueID, registration := buildFinalizedEnvelopeWorkspace(t, definitions)
		registrationJSON, err := json.Marshal(registration)
		if err != nil {
			t.Fatalf("encode cached gate registration: %v", err)
		}
		return workspace, []string{issueID, string(registrationJSON)}
	})
	var registration gateRegistrationView
	if err := json.Unmarshal([]byte(values[1]), &registration); err != nil {
		t.Fatalf("decode cached gate registration: %v", err)
	}
	return workspace, values[0], registration
}

func buildFinalizedEnvelopeWorkspace(t *testing.T, definitions []any) (string, string, gateRegistrationView) {
	t.Helper()
	workspace, values := cloneReusableWorkspaceFixture(t, "prepared-envelope", func() (string, []string) {
		workspace := clonePristineWorkspace(t)
		created := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Envelope issue", "--actor", "agent:developer")
		if created.err != nil {
			t.Fatalf("create issue: %v; stderr=%q", created.err, created.stderr)
		}
		var issue struct {
			Result struct {
				Issue struct {
					ID string `json:"id"`
				} `json:"issue"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(created.stdout), &issue); err != nil || issue.Result.Issue.ID == "" {
			t.Fatalf("decode created issue: %v; stdout=%q", err, created.stdout)
		}
		activationID := prepareActivationForTest(t, workspace)
		profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "envelope-builders", []any{map[string]any{
			"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/envelope-builders", "session_ref": nil,
			"capabilities": authorityCapabilities(),
		}})
		return workspace, []string{issue.Result.Issue.ID, activationID, profileDigest}
	})
	issueID, activationID, profileDigest := values[0], values[1], values[2]
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", definitions)
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-commit-envelope.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-envelope"); result.err != nil {
		t.Fatalf("commit activation: %v; stderr=%q", result.err, result.stderr)
	}
	finalizePath := filepath.Join(workspace, "activation-finalize-envelope.json")
	writeActivationInput(t, finalizePath, map[string]any{"activation_id": activationID, "expected_status": "committed"})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "finalize", "--input", finalizePath, "--idempotency-key", "activation-finalize-envelope"); result.err != nil {
		t.Fatalf("finalize activation: %v; stderr=%q", result.err, result.stderr)
	}
	registrationInput := gateRegisterInput(t, workspace)
	registrationInput["scope"] = map[string]any{"issue_id": issueID}
	registrationPath := filepath.Join(workspace, "gate-register-envelope.json")
	writeActivationInput(t, registrationPath, registrationInput)
	return workspace, issueID, registerGateForTest(t, workspace, registrationPath, "gate-register-envelope")
}

func gateDefinitionsForTest(modifyDefinition func(map[string]any)) []any {
	definitions := validGateDefinitions()
	if modifyDefinition != nil {
		modifyDefinition(definitions[0].(map[string]any))
	}
	return definitions
}

func authorityCapabilities() []string {
	return []string{
		"artifact.verify", "claim.acquire", "claim.release", "claim.renew", "envelope.activate", "envelope.complete",
		"envelope.create", "evidence.link", "gate.register", "gate.run", "issue.transition", "run.finish", "run.start",
	}
}

func envelopeCreateInput(t *testing.T, workspace, issueID string, registration gateRegistrationView) map[string]any {
	t.Helper()
	repository := testEnvelopeRepository(t, workspace)
	source := map[string]any{
		"digest": "sha256:" + strings.Repeat("0", 64), "display_label": "Gate policy", "external_id": nil,
		"path": "project-docs/gates/unit.json", "repository": repository, "section": nil, "source_type": "project-file",
	}
	return map[string]any{
		"outcome": "Prove the envelope value", "non_goals": []string{},
		"risk_result":      map[string]any{"source": source, "result_digest": "sha256:" + strings.Repeat("1", 64)},
		"repository_scope": repository, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{},
		"authority_grants": []any{map[string]any{
			"actions":        []string{"envelope.activate"},
			"actor_selector": map[string]any{"actor_id": "actor-builder", "required_capabilities": []string{"envelope.activate"}},
			"scope":          map[string]any{"issue_id": issueID},
		}},
		"obligations": []any{map[string]any{
			"id": testEvidenceObligationID, "type": "gate-evidence", "applicability": "required",
			"guarded_actions":         []string{"envelope.complete", "issue.transition-terminal-success"},
			"subject":                 map[string]any{"type": "issue", "id": issueID},
			"registration":            map[string]any{"id": registration.ID, "version": registration.Version},
			"accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"},
			"freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": source,
		}},
		"budgets": []any{map[string]any{
			"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 3, "warning_threshold": nil,
			"exhaustion_action": "unavailable", "consumed": nil, "remaining": nil,
		}},
	}
}

func testEnvelopeRepository(t *testing.T, workspace string) map[string]any {
	t.Helper()
	project, err := config.ReadProjectConfig(filepath.Join(workspace, ".fid", "config.toml"))
	if err != nil {
		t.Fatalf("read project config: %v", err)
	}
	return map[string]any{
		"base_oid": nil, "capture_source": "git", "dirty": false, "dirty_digest": nil,
		"head_oid": strings.Repeat("a", 40), "project_id": project.ProjectID, "tree_oid": strings.Repeat("b", 40), "worktree_id": "worktree-1",
	}
}

func assertEnvelopeShape(t *testing.T, output, wantStatus string) string {
	t.Helper()
	value := decodeCommandJSONStructure(t, output).(map[string]any)
	envelope := value["result"].(map[string]any)["envelope"].(map[string]any)
	if value["operation"] != "envelope.create" && value["operation"] != "envelope.get" && value["operation"] != "envelope.activate" {
		t.Fatalf("operation = %#v, want envelope operation", value["operation"])
	}
	if envelope["current_revision"] != json.Number("1") || envelope["status"] != wantStatus || len(envelope["related_issue_ids"].([]any)) != 0 {
		t.Fatalf("envelope = %#v, want revision 1 %s and no related issues", envelope, wantStatus)
	}
	obligation := envelope["obligations"].([]any)[0].(map[string]any)
	if obligation["type"] != "gate-evidence" || obligation["derived_status"] != "unsatisfied" ||
		strings.Join(anyStrings(t, obligation["guarded_actions"]), ",") != "envelope.complete,issue.transition-terminal-success" {
		t.Fatalf("obligation = %#v, want one gate-evidence dual-guard value", obligation)
	}
	budget := envelope["budgets"].([]any)[0].(map[string]any)
	if budget["consumed"] != nil || budget["remaining"] != nil {
		t.Fatalf("budget = %#v, want immutable null counters", budget)
	}
	return envelope["id"].(string)
}

func cloneEnvelopeInput(t *testing.T, input map[string]any) map[string]any {
	t.Helper()
	payload, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal envelope input: %v", err)
	}
	var cloned map[string]any
	if err := json.Unmarshal(payload, &cloned); err != nil {
		t.Fatalf("unmarshal envelope input: %v", err)
	}
	return cloned
}

func anyStrings(t *testing.T, value any) []string {
	t.Helper()
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("value = %T, want []any", value)
	}
	result := make([]string, len(items))
	for index, item := range items {
		result[index], ok = item.(string)
		if !ok {
			t.Fatalf("item %d = %T, want string", index, item)
		}
	}
	return result
}

func journalEventCount(t *testing.T, workspace string) int {
	t.Helper()
	return len(splitJournalLines(mustReadCommandFile(t, mustJournalPath(t, workspace))))
}
