package main

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/review"
)

func TestDecodeClosedReviewInputRejectsMalformedInapplicability(t *testing.T) {
	t.Parallel()
	for _, payload := range []string{
		`{"inapplicable_findings":null}`,
		`{"inapplicable_findings":[],"inapplicable_findings":[]}`,
		`{"inapplicable_findings":[{"finding_id":"FND-1","reason":"one","reason":"two"}]}`,
		`{"inapplicable_findings":[{"finding_id":"FND-1","reason":"one","unknown":true}]}`,
	} {
		var request review.RecordRequest
		if err := decodeClosedReviewInput(strings.NewReader(payload), "-", &request); err == nil {
			t.Fatalf("decodeClosedReviewInput(%s) accepted malformed input", payload)
		}
	}
	var request review.RecordRequest
	if err := decodeClosedReviewInput(strings.NewReader(`{"inapplicable_findings":[]}`), "-", &request); err != nil {
		t.Fatalf("decodeClosedReviewInput(empty optional collection) error = %v", err)
	}
}

func TestFindingAddressDoesNotResolveOrUnblock(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID, _, source := reviewWorkspaceWithTwoActivatedReviewers(t)
	repository := runRepositoryForTest(t, workspace)

	reviewClaim := acquireReviewClaim(t, workspace, issueID, envelopeID, "actor-reviewer-one", "review", "finding-address-review")
	reviewRun := startReviewRun(t, workspace, issueID, envelopeID, reviewClaim, "actor-reviewer-one", "finding-address-review")
	recordPath := filepath.Join(workspace, "finding-address-review.json")
	recordInput := reviewRecordInput(issueID, reviewRun, repository, source, "changes-required", []any{
		reviewFinding(issueID, "FND-addressed", "major", true),
	}, []any{}, []any{})
	suppliedState := cloneJSONMap(t, recordInput)
	suppliedState["findings"].([]any)[0].(map[string]any)["state"] = "resolved"
	writeActivationInput(t, recordPath, suppliedState)
	beforeState := snapshotWorkspacePersistence(t, workspace)
	rejectedState := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-reviewer-one", "--input", recordPath, "--idempotency-key", "finding-supplied-state")
	if rejectedState.err == nil || !strings.Contains(rejectedState.stderr, `"code":"invalid-argument"`) || !snapshotWorkspacePersistence(t, workspace).equal(beforeState) {
		t.Fatalf("caller-supplied Finding state = %#v, want no-Event closed-input rejection", rejectedState)
	}
	writeActivationInput(t, recordPath, recordInput)
	recorded := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-reviewer-one", "--input", recordPath, "--idempotency-key", "finding-address-review")
	if recorded.err != nil || recorded.stderr != "" {
		t.Fatalf("record Review with Finding = %#v", recorded)
	}
	review := decodeCommandJSONStructure(t, recorded.stdout).(map[string]any)["result"].(map[string]any)["review"].(map[string]any)
	closeReviewRun(t, workspace, reviewRun, repository, "actor-reviewer-one", "finding-address-review")
	releaseReviewClaim(t, workspace, reviewClaim, "actor-reviewer-one", "finding-address-review")

	addressPath := filepath.Join(workspace, "finding-address.json")
	writeActivationInput(t, addressPath, map[string]any{"finding_id": "FND-addressed", "claim_id": "CLM-missing", "expected_claim_version": 1, "reason": "no Claim"})
	beforeNoClaim := snapshotWorkspacePersistence(t, workspace)
	noClaim := executeCommandForTest(workspace, "--project", workspace, "--json", "finding", "address", "--actor", "actor-builder", "--input", addressPath, "--idempotency-key", "finding-address-no-claim")
	if noClaim.err == nil || !strings.Contains(noClaim.stderr, `"code":"claim-not-found"`) || !snapshotWorkspacePersistence(t, workspace).equal(beforeNoClaim) {
		t.Fatalf("address Finding without Claim = %#v, want typed no-Event rejection", noClaim)
	}
	reviewModeClaim := acquireReviewClaim(t, workspace, issueID, envelopeID, "actor-reviewer-two", "review", "finding-address-review-mode")
	writeActivationInput(t, addressPath, map[string]any{"finding_id": "FND-addressed", "claim_id": reviewModeClaim["id"], "expected_claim_version": reviewModeClaim["version"], "reason": "wrong mode"})
	beforeMode := snapshotWorkspacePersistence(t, workspace)
	wrongMode := executeCommandForTest(workspace, "--project", workspace, "--json", "finding", "address", "--actor", "actor-reviewer-two", "--input", addressPath, "--idempotency-key", "finding-address-review-mode")
	if wrongMode.err == nil || !strings.Contains(wrongMode.stderr, `"code":"claim-mode-invalid"`) || !strings.Contains(wrongMode.stderr, `"required_mode":"mutation"`) || !snapshotWorkspacePersistence(t, workspace).equal(beforeMode) {
		t.Fatalf("address Finding under review Claim = %#v, want mutation-mode no-Event rejection", wrongMode)
	}
	releaseReviewClaim(t, workspace, reviewModeClaim, "actor-reviewer-two", "finding-address-review-mode")

	mutationClaim := acquireReviewClaim(t, workspace, issueID, envelopeID, "actor-builder", "mutation", "finding-address-rework")
	writeActivationInput(t, addressPath, map[string]any{
		"finding_id":             "FND-addressed",
		"claim_id":               mutationClaim["id"],
		"expected_claim_version": mutationClaim["version"],
		"reason":                 "implemented the requested correction",
	})
	addressed := executeCommandForTest(workspace, "--project", workspace, "--json", "finding", "address", "--actor", "actor-builder", "--input", addressPath, "--idempotency-key", "finding-address")
	if addressed.err != nil || addressed.stderr != "" {
		t.Fatalf("address Finding under active mutation Claim = %#v", addressed)
	}
	if !strings.Contains(addressed.stdout, `"action":"addressed"`) || !strings.Contains(addressed.stdout, `"state":"addressed"`) {
		t.Fatalf("address Finding result = %q, want addressed disposition and derived state", addressed.stdout)
	}

	got := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", review["id"].(string))
	if got.err != nil || !strings.Contains(got.stdout, `"blocking":true`) || !strings.Contains(got.stdout, `"state":"addressed"`) || strings.Contains(got.stdout, `"state":"resolved"`) {
		t.Fatalf("Review after address = %#v, want blocking Finding still unresolved in addressed state", got)
	}
	listed := executeCommandForTest(workspace, "--project", workspace, "--json", "finding", "list", "--issue", issueID, "--review", review["id"].(string), "--state", "addressed", "--severity", "major", "--limit", "1")
	if listed.err != nil || !strings.Contains(listed.stdout, `"id":"FND-addressed"`) || !strings.Contains(listed.stdout, `"state":"addressed"`) {
		t.Fatalf("filtered Finding list = %#v, want addressed Finding", listed)
	}
}

func TestReviewRecordMixedTerminalDispositions(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID, _, source := reviewWorkspaceWithTwoActivatedReviewers(t)
	repository := runRepositoryForTest(t, workspace)

	firstClaim := acquireReviewClaim(t, workspace, issueID, envelopeID, "actor-reviewer-one", "review", "finding-terminal-first")
	firstRun := startReviewRun(t, workspace, issueID, envelopeID, firstClaim, "actor-reviewer-one", "finding-terminal-first")
	firstPath := filepath.Join(workspace, "finding-terminal-first.json")
	writeActivationInput(t, firstPath, reviewRecordInput(issueID, firstRun, repository, source, "changes-required", []any{
		reviewFinding(issueID, "FND-terminal-inapplicable", "major", true),
		reviewFinding(issueID, "FND-terminal-resolved", "major", true),
	}, []any{}, []any{}))
	first := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-reviewer-one", "--input", firstPath, "--idempotency-key", "finding-terminal-first")
	if first.err != nil || first.stderr != "" {
		t.Fatalf("record Review with terminal targets = %#v", first)
	}
	firstReview := decodeCommandJSONStructure(t, first.stdout).(map[string]any)["result"].(map[string]any)["review"].(map[string]any)
	closeReviewRun(t, workspace, firstRun, repository, "actor-reviewer-one", "finding-terminal-first")
	releaseReviewClaim(t, workspace, firstClaim, "actor-reviewer-one", "finding-terminal-first")

	secondClaim := acquireReviewClaim(t, workspace, issueID, envelopeID, "actor-reviewer-two", "review", "finding-terminal-second")
	secondRun := startReviewRun(t, workspace, issueID, envelopeID, secondClaim, "actor-reviewer-two", "finding-terminal-second")
	secondInput := reviewRecordInput(issueID, secondRun, repository, source, "pass", []any{}, []any{"FND-terminal-resolved"}, []any{})
	secondInput["inapplicable_findings"] = []any{map[string]any{"finding_id": "FND-terminal-inapplicable", "reason": "generated file is outside the accepted scope"}}
	secondPath := filepath.Join(workspace, "finding-terminal-second.json")
	writeActivationInput(t, secondPath, secondInput)
	second := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-reviewer-two", "--input", secondPath, "--idempotency-key", "finding-terminal-second")
	if second.err != nil || second.stderr != "" {
		t.Fatalf("record mixed terminal dispositions = %#v", second)
	}
	secondReview := decodeCommandJSONStructure(t, second.stdout).(map[string]any)["result"].(map[string]any)["review"].(map[string]any)
	if !reflect.DeepEqual(secondReview["resolved_finding_ids"], []any{"FND-terminal-resolved"}) {
		t.Fatalf("mixed Review resolved IDs = %#v, want resolved action only", secondReview["resolved_finding_ids"])
	}
	got := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", firstReview["id"].(string))
	if got.err != nil || !strings.Contains(got.stdout, `"state":"resolved"`) || !strings.Contains(got.stdout, `"state":"inapplicable"`) || !strings.Contains(got.stdout, `"reason":"generated file is outside the accepted scope"`) || !strings.Contains(got.stdout, `"qualifying_review_id":"`+secondReview["id"].(string)+`"`) {
		t.Fatalf("mixed terminal Finding readback = %#v", got)
	}
	beforeRebuild := got.stdout
	if rebuilt := executeCommandForTest(workspace, "--project", workspace, "--json", "rebuild"); rebuilt.err != nil || rebuilt.stderr != "" {
		t.Fatalf("rebuild mixed terminal dispositions = %#v", rebuilt)
	}
	got = executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", firstReview["id"].(string))
	if got.err != nil || got.stdout != beforeRebuild {
		t.Fatalf("mixed terminal readback after rebuild = %#v, want byte-identical projection", got)
	}
}
