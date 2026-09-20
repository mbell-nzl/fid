package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestReviewRecordRejectsNonIndependentActor is the Task 4 vertical tracer:
// recording under the mutating actor must be rejected before an Event is added.
func TestReviewRecordRejectsNonIndependentActor(t *testing.T) {
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "review-independent-builder")
	repository := runRepositoryForTest(t, workspace)
	runPath := filepath.Join(workspace, "review-independent-run.json")
	writeActivationInput(t, runPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace,
		"--project", workspace, "--json", "run", "start", "--actor", "actor-builder",
		"--input", runPath, "--idempotency-key", "review-independent-run",
	)
	if started.err != nil || started.stderr != "" {
		t.Fatalf("start mutating Run = %#v, want stdout-only success", started)
	}
	run := runFromResponse(t, started.stdout)
	inputPath := filepath.Join(workspace, "review-record.json")
	writeActivationInput(t, inputPath, map[string]any{"run_id": run["id"]})

	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "review", "record", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "review-not-independent",
	)
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"review-not-independent"`) {
		t.Fatalf("review record by mutating actor = %#v, want typed review-not-independent rejection", result)
	}
}

func TestReviewRecordCreatesImmutableReview(t *testing.T) {
	workspace, issueID, registration := finalizedObservedEnvelopeWorkspaceWithDefinitions(t, gateDefinitionsForTest(nil))
	input := envelopeCreateInput(t, workspace, issueID, registration)
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "envelope.activate", "review.record", "run.start"}
	grant["actor_selector"] = map[string]any{"actor_id": "actor-builder", "required_capabilities": []string{"claim.acquire", "envelope.activate", "review.record", "run.start"}}
	path := filepath.Join(workspace, "review-envelope.json")
	writeActivationInput(t, path, input)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", path, "--idempotency-key", "review-envelope-create")
	if created.err != nil {
		t.Fatalf("create envelope = %#v", created)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
	activatePath := filepath.Join(workspace, "review-envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder", "--input", activatePath, "--idempotency-key", "review-envelope-activate")
	if activated.err != nil {
		t.Fatalf("activate envelope = %#v", activated)
	}
	claim, _ := acquireGenericClaimForTestWithMode(t, workspace, issueID, envelopeID, "review-record-claim", 60, "review")
	repository := runRepositoryForTest(t, workspace)
	runPath := filepath.Join(workspace, "review-record-run.json")
	writeActivationInput(t, runPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", runPath, "--idempotency-key", "review-record-run")
	if started.err != nil {
		t.Fatalf("start review Run = %#v", started)
	}
	run := runFromResponse(t, started.stdout)
	source := input["sources"].([]any)[0]
	recordPath := filepath.Join(workspace, "review-record-input.json")
	findings := []any{
		map[string]any{"id": "FND-review-one", "subject": map[string]any{"type": "issue", "id": issueID}, "scope": map[string]any{"paths": []any{}}, "severity": "minor", "summary": "first non-blocking finding", "blocking": false},
		map[string]any{"id": "FND-review-two", "subject": map[string]any{"type": "issue", "id": issueID}, "scope": map[string]any{"paths": []any{}}, "severity": "nit", "summary": "second non-blocking finding", "blocking": false},
	}
	writeActivationInput(t, recordPath, map[string]any{"run_id": run["id"], "expected_run_version": run["version"], "profile": "review-default", "repository": repository, "coverage": map[string]any{"items": []any{map[string]any{"source": source, "assessed": true}}, "exclusions": []any{}}, "verdict": "pass", "findings": findings, "resolves_finding_ids": []any{}, "evidence_refs": []any{}, "receipt_ids": []any{}})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-builder", "--input", recordPath, "--idempotency-key", "review-record")
	if result.err != nil || result.stderr != "" || !strings.Contains(result.stdout, `"verdict":"pass"`) {
		t.Fatalf("record Review = %#v, want immutable pass Review", result)
	}
	recorded := decodeCommandJSONStructure(t, result.stdout).(map[string]any)["result"].(map[string]any)["review"].(map[string]any)
	got := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", recorded["id"].(string))
	if got.err != nil || got.stderr != "" || !strings.Contains(got.stdout, `"verdict":"pass"`) {
		t.Fatalf("get recorded Review = %#v, want immutable pass Review", got)
	}
	if !strings.Contains(got.stdout, `"FND-review-one"`) || !strings.Contains(got.stdout, `"FND-review-two"`) {
		t.Fatalf("get recorded Review = %#v, want both immutable Findings", got)
	}
}
