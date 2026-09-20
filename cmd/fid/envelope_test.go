package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/workflow"
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
		{"duplicate obligation ID", func(value map[string]any) {
			value["obligations"] = append(value["obligations"].([]any), value["obligations"].([]any)[0])
		}},
		{"computed status", func(value map[string]any) {
			value["obligations"].([]any)[0].(map[string]any)["derived_status"] = "satisfied"
		}},
		{"other obligation type", func(value map[string]any) { value["obligations"].([]any)[0].(map[string]any)["type"] = "human-outcome" }},
		{"guard widening", func(value map[string]any) {
			value["obligations"].([]any)[0].(map[string]any)["guarded_actions"] = []string{"envelope.complete", "issue.transition-terminal-success", "run.close"}
		}},
		{"mutable budget", func(value map[string]any) { value["budgets"].([]any)[0].(map[string]any)["consumed"] = 1 }},
		{"budget output state empty", func(value map[string]any) { value["budgets"].([]any)[0].(map[string]any)["state"] = "" }},
		{"budget output state null", func(value map[string]any) { value["budgets"].([]any)[0].(map[string]any)["state"] = nil }},
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

func TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	input := envelopeCreateInput(t, workspace, issueID, registration)
	source := input["sources"].([]any)[0]
	input["obligations"] = []any{
		input["obligations"].([]any)[0],
		map[string]any{"id": "OBL-review", "type": "review", "applicability": "required", "guarded_actions": []string{"envelope.complete"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "review_profile": "standard", "accepted_evidence_types": []string{"review"}, "accepted_trust_sources": []string{"activated-actor"}, "freshness_target": "repository-tree", "artifact_policy": "metadata-sufficient", "independence": "distinct-actor", "required_lane": nil, "source": source},
		map[string]any{"id": "OBL-human", "type": "human-outcome", "applicability": "required", "guarded_actions": []string{"envelope.complete"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "outcome_class": "release", "accepted_evidence_types": []string{"human-outcome"}, "accepted_trust_sources": []string{"activated-human"}, "freshness_target": "none", "artifact_policy": "metadata-sufficient", "independence": "human", "required_lane": nil, "source": source},
	}
	input["protected_decisions"] = []any{map[string]any{"ordinal": 0, "decision_class": "release", "required_capabilities": []string{"envelope.complete"}, "affected_actions": []string{"envelope.complete"}}}
	inputPath := filepath.Join(workspace, "typed-envelope.json")
	writeActivationInput(t, inputPath, input)

	result := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "typed-envelope")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("typed envelope create = %#v, want success", result)
	}
	payload := decodeCommandJSONStructure(t, result.stdout).(map[string]any)
	envelope := payload["result"].(map[string]any)["envelope"].(map[string]any)
	obligations := envelope["obligations"].([]any)
	if len(obligations) != 3 || obligations[1].(map[string]any)["registration"] != nil || obligations[2].(map[string]any)["registration"] != nil {
		t.Fatalf("typed Envelope obligations = %#v, want three rows with null new-type registrations", obligations)
	}
	if decisions := envelope["protected_decisions"].([]any); len(decisions) != 1 {
		t.Fatalf("typed Envelope protected decisions = %#v, want one authoring row", decisions)
	}
}

func TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries(t *testing.T) {
	for _, test := range []struct {
		name            string
		obligationCount int
		decisionCount   int
	}{
		{name: "one obligation and zero decisions", obligationCount: 1, decisionCount: 0},
		{name: "sixteen obligations and sixteen decisions", obligationCount: 16, decisionCount: 16},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
			input := typedEnvelopeCreateInput(t, workspace, issueID, registration)
			input["obligations"] = commandObligationsForCount(t, input, issueID, test.obligationCount)
			input["protected_decisions"] = commandProtectedDecisionMaps(test.decisionCount)
			inputPath := filepath.Join(workspace, "envelope-boundary.json")
			writeActivationInput(t, inputPath, input)

			result := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "envelope-boundary-"+strings.ReplaceAll(test.name, " ", "-"))
			if result.err != nil || result.stderr != "" {
				t.Fatalf("boundary create = %#v, want success", result)
			}
			payload := decodeCommandJSONStructure(t, result.stdout).(map[string]any)
			envelope := payload["result"].(map[string]any)["envelope"].(map[string]any)
			if got := len(envelope["obligations"].([]any)); got != test.obligationCount {
				t.Fatalf("obligation count = %d, want %d", got, test.obligationCount)
			}
			if got := len(envelope["protected_decisions"].([]any)); got != test.decisionCount {
				t.Fatalf("protected decision count = %d, want %d", got, test.decisionCount)
			}
		})
	}
}

func TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(t *testing.T, input map[string]any, issueID string, registration gateRegistrationView)
	}{
		{name: "zero obligations", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["obligations"] = []any{}
		}},
		{name: "seventeen obligations", mutate: func(t *testing.T, input map[string]any, issueID string, registration gateRegistrationView) {
			input["obligations"] = commandObligationsForCount(t, input, issueID, 17)
		}},
		{name: "duplicate obligation ID", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			obligations := input["obligations"].([]any)
			obligations[1].(map[string]any)["id"] = obligations[0].(map[string]any)["id"]
		}},
		{name: "second gate-evidence obligation", mutate: func(t *testing.T, input map[string]any, issueID string, registration gateRegistrationView) {
			obligations := input["obligations"].([]any)
			input["obligations"] = append(obligations, commandObligationsForCount(t, input, issueID, 1)[0])
			input["obligations"].([]any)[3].(map[string]any)["id"] = "OBL-second-gate"
		}},
		{name: "duplicate review profile", mutate: func(_ *testing.T, input map[string]any, issueID string, _ gateRegistrationView) {
			obligations := input["obligations"].([]any)
			duplicate := commandReviewObligation(input["sources"].([]any)[0], issueID, "OBL-review-duplicate", "standard")
			input["obligations"] = append(obligations, duplicate)
		}},
		{name: "duplicate outcome class", mutate: func(_ *testing.T, input map[string]any, issueID string, _ gateRegistrationView) {
			obligations := input["obligations"].([]any)
			duplicate := commandHumanObligation(input["sources"].([]any)[0], issueID, "OBL-human-duplicate", "release")
			input["obligations"] = append(obligations, duplicate)
		}},
		{name: "non-primary subject", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["obligations"].([]any)[0].(map[string]any)["subject"] = map[string]any{"type": "issue", "id": "ISS-other"}
		}},
		{name: "seventeen protected decisions", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["protected_decisions"] = commandProtectedDecisionMaps(17)
		}},
		{name: "duplicate protected decision class", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["protected_decisions"] = []any{commandProtectedDecisionMap(0, "release", "envelope.complete"), commandProtectedDecisionMap(1, "release", "envelope.complete")}
		}},
		{name: "duplicate protected decision ordinal", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["protected_decisions"] = []any{commandProtectedDecisionMap(0, "release", "envelope.complete"), commandProtectedDecisionMap(0, "rollback", "envelope.complete")}
		}},
		{name: "unguarded protected decision action", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			for _, obligation := range input["obligations"].([]any) {
				obligation.(map[string]any)["guarded_actions"] = []string{"envelope.complete"}
			}
			input["protected_decisions"] = []any{commandProtectedDecisionMap(0, "release", "issue.transition-terminal-success")}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
			input := typedEnvelopeCreateInput(t, workspace, issueID, registration)
			test.mutate(t, input, issueID, registration)
			inputPath := filepath.Join(workspace, "envelope-rejected.json")
			writeActivationInput(t, inputPath, input)
			before := snapshotWorkspacePersistence(t, workspace)

			result := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "envelope-rejected-"+strings.ReplaceAll(test.name, " ", "-"))
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
				t.Fatalf("rejected command = %#v, want invalid-argument with no stdout", result)
			}
			if after := snapshotWorkspacePersistence(t, workspace); !before.equal(after) {
				t.Fatalf("rejected %s changed journal or projection", test.name)
			}
		})
	}
}

func TestEnvelopeCreateCommandRejectsEveryForbiddenTypedFieldCombinationWithoutPersistenceChange(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{name: "gate missing registration", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["registration"] = nil
		}},
		{name: "gate has review profile", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["review_profile"] = "profile"
		}},
		{name: "gate has outcome class", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["outcome_class"] = "outcome"
		}},
		{name: "gate evidence type mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["accepted_evidence_types"] = []string{"review"}
		}},
		{name: "gate trust source mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["accepted_trust_sources"] = []string{"activated-actor"}
		}},
		{name: "gate freshness mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["freshness_target"] = "none"
		}},
		{name: "gate artifact policy mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["artifact_policy"] = "metadata-sufficient"
		}},
		{name: "gate independence mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["independence"] = "distinct-actor"
		}},
		{name: "review has registration", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["registration"] = map[string]any{"id": "REG-test", "version": 1}
		}},
		{name: "review missing profile", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["review_profile"] = nil
		}},
		{name: "review has outcome class", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["outcome_class"] = "outcome"
		}},
		{name: "review evidence type mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["accepted_evidence_types"] = []string{"receipt"}
		}},
		{name: "review trust source mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["accepted_trust_sources"] = []string{"fid-executed"}
		}},
		{name: "review freshness mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["freshness_target"] = "none"
		}},
		{name: "review artifact policy mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["artifact_policy"] = "required-available"
		}},
		{name: "review independence mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["independence"] = "human"
		}},
		{name: "human has registration", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["registration"] = map[string]any{"id": "REG-test", "version": 1}
		}},
		{name: "human has review profile", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["review_profile"] = "profile"
		}},
		{name: "human missing outcome class", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["outcome_class"] = nil
		}},
		{name: "human evidence type mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["accepted_evidence_types"] = []string{"review"}
		}},
		{name: "human trust source mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["accepted_trust_sources"] = []string{"activated-actor"}
		}},
		{name: "human freshness mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["freshness_target"] = "repository-tree"
		}},
		{name: "human artifact policy mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["artifact_policy"] = "required-available"
		}},
		{name: "human independence mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["independence"] = "distinct-actor"
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
			input := typedEnvelopeCreateInput(t, workspace, issueID, registration)
			test.mutate(input)
			inputPath := filepath.Join(workspace, "envelope-forbidden-field.json")
			writeActivationInput(t, inputPath, input)
			before := snapshotWorkspacePersistence(t, workspace)

			result := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "envelope-forbidden-field-"+strings.ReplaceAll(test.name, " ", "-"))
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
				t.Fatalf("forbidden typed field command = %#v, want invalid-argument with no stdout", result)
			}
			if after := snapshotWorkspacePersistence(t, workspace); !before.equal(after) {
				t.Fatalf("forbidden typed field %s changed journal or projection", test.name)
			}
		})
	}
}

func commandObligationByType(t *testing.T, input map[string]any, wantType string) map[string]any {
	t.Helper()
	for _, value := range input["obligations"].([]any) {
		obligation := value.(map[string]any)
		if obligation["type"] == wantType {
			return obligation
		}
	}
	t.Fatalf("typed input has no %q Obligation", wantType)
	return nil
}

func TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState(t *testing.T) {
	t.Parallel()
	for _, state := range []string{`""`, `null`} {
		t.Run(state, func(t *testing.T) {
			var request workflow.EnvelopeCreateRequest
			err := decodeClosedEnvelopeInput(strings.NewReader(`{"budgets":[{"state":`+state+`}]}`), "-", &request)
			if err == nil || !strings.Contains(err.Error(), "budget state is output-only") {
				t.Fatalf("decodeClosedEnvelopeInput() error = %v, want output-only state rejection", err)
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
	return finalizedEnvelopeWorkspaceWithRepository(t, definitions, false)
}

func finalizedObservedEnvelopeWorkspaceWithDefinitions(t *testing.T, definitions []any) (string, string, gateRegistrationView) {
	return finalizedEnvelopeWorkspaceWithRepository(t, definitions, true)
}

func finalizedEnvelopeWorkspaceWithRepository(t *testing.T, definitions []any, observeRepository bool) (string, string, gateRegistrationView) {
	t.Helper()
	if observeRepository {
		return buildFinalizedEnvelopeWorkspaceWithRepository(t, definitions, true)
	}
	key := reusableWorkspaceFixtureKey(t, "finalized-envelope", map[string]any{"definitions": definitions, "observed_repository": observeRepository})
	workspace, values := cloneReusableWorkspaceFixture(t, key, func() (string, []string) {
		workspace, issueID, registration := buildFinalizedEnvelopeWorkspaceWithRepository(t, definitions, observeRepository)
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
	return buildFinalizedEnvelopeWorkspaceWithRepository(t, definitions, false)
}

func buildFinalizedEnvelopeWorkspaceWithRepository(t *testing.T, definitions []any, observeRepository bool) (string, string, gateRegistrationView) {
	t.Helper()
	// Reuse only facts produced before any path-bound repository observation.
	// Registration and every subsequent Envelope/Claim/Run use this clone's root.
	key := reusableWorkspaceFixtureKey(t, "unregistered-envelope", definitions)
	workspace, values := cloneReusableWorkspaceFixture(t, key, func() (string, []string) {
		root, issueID := buildUnregisteredEnvelopeWorkspace(t, definitions)
		return root, []string{issueID}
	})
	issueID := values[0]
	if observeRepository {
		initializeObservedRepositoryForTest(t, workspace)
	}
	registrationInput := gateRegisterInput(t, workspace)
	registrationInput["scope"] = map[string]any{"issue_id": issueID}
	registrationPath := filepath.Join(workspace, "gate-register-envelope.json")
	writeActivationInput(t, registrationPath, registrationInput)
	return workspace, issueID, registerGateForTest(t, workspace, registrationPath, "gate-register-envelope")
}

func buildUnregisteredEnvelopeWorkspace(t *testing.T, definitions []any) (string, string) {
	t.Helper()
	prepare := func() (string, []string) {
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
		}, map[string]any{
			"actor_id": "actor-harness", "actor_type": "service", "producer": "activated-harness", "provenance_ref": "activation:profile/envelope-harness", "session_ref": nil,
			"capabilities": []string{"run.observe"},
		}})
		return workspace, []string{issue.Result.Issue.ID, activationID, profileDigest}
	}
	workspace, values := cloneReusableWorkspaceFixture(t, reusableWorkspaceFixtureKey(t, "prepared-envelope", false), prepare)
	issueID, activationID, profileDigest := values[0], values[1], values[2]
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", definitions)
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-commit-envelope.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, workspace, activationID), "component_digests": digests})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-envelope"); result.err != nil {
		t.Fatalf("commit activation: %v; stderr=%q", result.err, result.stderr)
	}
	finalizePath := filepath.Join(workspace, "activation-finalize-envelope.json")
	writeActivationInput(t, finalizePath, map[string]any{"activation_id": activationID, "expected_status": "committed", "expected_version": activationVersionForTest(t, workspace, activationID)})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "finalize", "--input", finalizePath, "--idempotency-key", "activation-finalize-envelope"); result.err != nil {
		t.Fatalf("finalize activation: %v; stderr=%q", result.err, result.stderr)
	}
	return workspace, issueID
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
		"artifact.verify", "claim.acquire", "claim.expire", "claim.release", "claim.renew", "envelope.activate", "envelope.complete",
		"envelope.create", "evidence.link", "gate.register", "gate.run", "issue.transition", "run.close", "run.observe", "run.start",
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
			"exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil,
		}},
	}
}

func typedEnvelopeCreateInput(t *testing.T, workspace, issueID string, registration gateRegistrationView) map[string]any {
	t.Helper()
	input := envelopeCreateInput(t, workspace, issueID, registration)
	source := input["sources"].([]any)[0]
	input["obligations"] = []any{
		input["obligations"].([]any)[0],
		commandReviewObligation(source, issueID, "OBL-review", "standard"),
		commandHumanObligation(source, issueID, "OBL-human", "release"),
	}
	input["protected_decisions"] = []any{}
	return input
}

func commandObligationsForCount(t *testing.T, input map[string]any, issueID string, count int) []any {
	t.Helper()
	if count < 1 || count > 17 {
		t.Fatalf("command obligation count = %d, want 1..17", count)
	}
	source := input["sources"].([]any)[0]
	gate := input["obligations"].([]any)[0].(map[string]any)
	obligations := []any{gate}
	for index := 1; index < count; index++ {
		if index%2 == 1 {
			obligations = append(obligations, commandReviewObligation(source, issueID, "OBL-review-"+string(rune('a'+index)), "profile-"+string(rune('a'+index))))
			continue
		}
		obligations = append(obligations, commandHumanObligation(source, issueID, "OBL-human-"+string(rune('a'+index)), "outcome-"+string(rune('a'+index))))
	}
	return obligations
}

func commandReviewObligation(source any, issueID, id, profile string) map[string]any {
	return map[string]any{
		"id": id, "type": "review", "applicability": "required", "guarded_actions": []string{"envelope.complete"},
		"subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "review_profile": profile,
		"accepted_evidence_types": []string{"review"}, "accepted_trust_sources": []string{"activated-actor"},
		"freshness_target": "repository-tree", "artifact_policy": "metadata-sufficient", "independence": "distinct-actor", "required_lane": nil, "source": source,
	}
}

func commandHumanObligation(source any, issueID, id, outcomeClass string) map[string]any {
	return map[string]any{
		"id": id, "type": "human-outcome", "applicability": "required", "guarded_actions": []string{"envelope.complete"},
		"subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "outcome_class": outcomeClass,
		"accepted_evidence_types": []string{"human-outcome"}, "accepted_trust_sources": []string{"activated-human"},
		"freshness_target": "none", "artifact_policy": "metadata-sufficient", "independence": "human", "required_lane": nil, "source": source,
	}
}

func commandProtectedDecisionMaps(count int) []any {
	decisions := make([]any, 0, count)
	for ordinal := 0; ordinal < count; ordinal++ {
		decisions = append(decisions, commandProtectedDecisionMap(ordinal, "decision-"+string(rune('a'+ordinal)), "envelope.complete"))
	}
	return decisions
}

func commandProtectedDecisionMap(ordinal int, decisionClass, action string) map[string]any {
	return map[string]any{
		"ordinal": ordinal, "decision_class": decisionClass, "required_capabilities": []string{"envelope.complete"}, "affected_actions": []string{action},
	}
}

func testEnvelopeRepository(t *testing.T, workspace string) map[string]any {
	t.Helper()
	projectConfig, err := config.ReadProjectConfig(filepath.Join(workspace, ".fid", "config.toml"))
	if err != nil {
		t.Fatalf("read project config: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, ".git")); err == nil {
		observation, err := project.ObserveRepository(workspace)
		if err != nil {
			t.Fatalf("observe test repository: %v", err)
		}
		return map[string]any{"base_oid": observation.HeadOID, "capture_source": "git", "dirty": observation.Dirty, "dirty_digest": observation.DirtyDigest, "head_oid": observation.HeadOID, "project_id": projectConfig.ProjectID, "tree_oid": observation.TreeOID, "worktree_id": observation.WorktreeID}
	}
	return map[string]any{
		"base_oid": nil, "capture_source": "git", "dirty": false, "dirty_digest": nil,
		"head_oid": strings.Repeat("a", 40), "project_id": projectConfig.ProjectID, "tree_oid": strings.Repeat("b", 40), "worktree_id": "worktree-1",
	}
}

func initializeObservedRepositoryForTest(t *testing.T, workspace string) {
	t.Helper()
	for _, arguments := range [][]string{{"init", "-q"}, {"config", "user.email", "fid@example.test"}, {"config", "user.name", "fid test"}} {
		if output, err := exec.Command("git", append([]string{"-C", workspace}, arguments...)...).CombinedOutput(); err != nil {
			t.Fatalf("initialize observed repository %v: %v (%s)", arguments, err, output)
		}
	}
	if err := os.WriteFile(filepath.Join(workspace, ".gitignore"), []byte(".fid/\n*.json\n"), 0o600); err != nil {
		t.Fatalf("write observed repository ignore: %v", err)
	}
	for _, arguments := range [][]string{{"add", ".gitignore"}, {"commit", "-qm", "fixture"}} {
		if output, err := exec.Command("git", append([]string{"-C", workspace}, arguments...)...).CombinedOutput(); err != nil {
			t.Fatalf("commit observed repository %v: %v (%s)", arguments, err, output)
		}
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
	if budget["consumed"] != json.Number("0") || budget["remaining"] != json.Number("3") || budget["state"] != "available" {
		t.Fatalf("budget = %#v, want derived counters and available state", budget)
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
