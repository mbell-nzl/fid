//go:build integration

package integration

import (
	"reflect"
	"strings"
	"testing"
)

// TestP23FindingDispositionRealProcessJourney proves the Task 5 lifecycle at
// the built-binary boundary. It keeps the exhaustive malformed-input matrix in
// the in-process contract tests, while exercising the externally observable
// authority, terminal-action, replay, and projection paths together.
func TestP23FindingDispositionRealProcessJourney(t *testing.T) {
	t.Parallel()
	actions := []string{
		"claim.acquire", "claim.release", "envelope.activate", "envelope.create",
		"concern.dispose", "concern.record", "finding.address", "gate.register", "review.record", "run.close", "run.observe", "run.start",
	}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation")
	a.task13ReplayCorpus = true
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "builder-before-review")

	firstClaim := p23AcquireClaim(t, a, "actor-denied", "review", "finding-creator")
	firstRun := p23StartRun(t, a, "actor-denied", firstClaim, "finding-creator")
	firstInput := p23ReviewInput(a, p21PathString(t, firstRun, "id"), p23Number(t, firstRun["version"]), "changes-required", []any{
		p23Finding(a.issueID, "FND-addressed", "major", true),
		p23Finding(a.issueID, "FND-inapplicable", "major", true),
		p23Finding(a.issueID, "FND-resolved", "major", true),
		p23Finding(a.issueID, "FND-same-inapplicable", "major", true),
		p23Finding(a.issueID, "FND-same-reviewer", "major", true),
	}, nil)

	t.Run("same_review_and_malformed_inapplicability_are_non_mutating", func(t *testing.T) {
		self := p23Clone(firstInput)
		self["findings"] = append(self["findings"].([]any), p23Finding(a.issueID, "FND-self", "major", true))
		self["resolves_finding_ids"] = []any{"FND-self"}
		self["inapplicable_findings"] = []any{map[string]any{"finding_id": "FND-self", "reason": "same review must fail"}}
		before := p21JournalEventCount(t, a.workspace)
		response := p23ReviewRecord(t, a, "actor-denied", "finding-self", self)
		assertP22ProcessFailure(t, "review.record", "invalid-argument", 2, response)
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("same-Review terminal request appended Event(s): %d -> %d", before, after)
		}

		malformed := p23Clone(firstInput)
		malformed["inapplicable_findings"] = nil
		response = p23ReviewRecord(t, a, "actor-denied", "finding-null", malformed)
		assertP22ProcessFailure(t, "review.record", "invalid-argument", 2, response)
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("null inapplicability collection appended Event(s): %d -> %d", before, after)
		}
	})

	first := p23ReviewRecord(t, a, "actor-denied", "finding-creator", firstInput)
	p23AssertReviewSuccess(t, "create Findings", first)
	p23AssertActivity(t, a.workspace, first, "review-recorded", a.issueID, 1)
	firstReview := p23ReviewFromResponse(t, first)
	for _, subject := range []struct{ kind, id string }{{"review", firstReview["id"].(string)}, {"finding", "FND-addressed"}} {
		input := map[string]any{"id": "CON-subject-" + subject.kind, "subject": map[string]any{"type": subject.kind, "id": subject.id}, "scope": map[string]any{"paths": []string{}}, "severity": "minor", "summary": "subject mapping", "reasoning": "Issue attribution", "references": []any{}}
		response := p23RunInput(t, a, "concern", "record", a.actor, "concern-subject-"+subject.kind, input)
		p23AssertReviewSuccess(t, "record "+subject.kind+" Concern", response)
		p23AssertActivity(t, a.workspace, response, "concern-recorded", a.issueID, 1)
	}
	if got := p23ReviewGet(t, a, firstReview["id"].(string)); !reflect.DeepEqual(got, firstReview) {
		t.Fatalf("first Review get result differs from its recorded typed result: got=%#v want=%#v", got, firstReview)
	}

	t.Run("omitted_and_empty_optional_collection_replay_the_task4_digest", func(t *testing.T) {
		before := p21JournalEventCount(t, a.workspace)
		empty := p23Clone(firstInput)
		empty["inapplicable_findings"] = []any{}
		replay := p23ReviewRecord(t, a, "actor-denied", "finding-creator", empty)
		p23AssertReviewSuccess(t, "empty optional collection replay", replay)
		if p21Path(t, p23Envelope(t, replay.Stdout), "result", "replayed") != true || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("empty optional collection replay = %#v, want Task 4 digest replay without an Event", replay)
		}
	})

	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, firstRun, "id"), p21PathString(t, firstClaim, "id"), "finding-creator")

	t.Run("address_requires_mutating_claim_and_stays_unresolved", func(t *testing.T) {
		before := p21JournalEventCount(t, a.workspace)
		noClaim := p23FindingAddress(t, a, a.actor, "finding-no-claim", map[string]any{"finding_id": "FND-addressed", "claim_id": "CLM-missing", "expected_claim_version": 1, "reason": "no active Claim"})
		if noClaim.ExitCode == 0 || noClaim.Stdout != "" || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("finding address without Claim = %#v, want typed no-Event rejection", noClaim)
		}

		reviewClaim := p23AcquireClaim(t, a, "actor-denied", "review", "finding-address-review")
		reviewMode := p23FindingAddress(t, a, "actor-denied", "finding-address-review", map[string]any{"finding_id": "FND-addressed", "claim_id": reviewClaim["id"], "expected_claim_version": reviewClaim["version"], "reason": "wrong claim mode"})
		assertP22ProcessFailure(t, "finding.address", "claim-mode-invalid", 4, reviewMode)
		if !strings.Contains(reviewMode.Stderr, `"required_mode":"mutation"`) || p21JournalEventCount(t, a.workspace) != before+1 {
			t.Fatalf("review Claim finding address = %#v, want mutation mode rejection only after Claim acquisition", reviewMode)
		}
		p23ReleaseClaim(t, a, "actor-denied", reviewClaim, "finding-address-review")

		mutationClaim := p23AcquireClaim(t, a, a.actor, "mutation", "finding-address-builder")
		addressed := p23FindingAddress(t, a, a.actor, "finding-address-builder", map[string]any{"finding_id": "FND-addressed", "claim_id": mutationClaim["id"], "expected_claim_version": mutationClaim["version"], "reason": "implemented the requested correction"})
		p23AssertReviewSuccess(t, "address Finding", addressed)
		p23AssertActivity(t, a.workspace, addressed, "finding-addressed", a.issueID, 1)
		addressedFinding := p23FindingByID(t, a, "FND-addressed")
		if addressedFinding["state"] != "addressed" || addressedFinding["blocking"] != true || !p23HasDisposition(t, addressedFinding, "addressed", a.actor, "", "implemented the requested correction") {
			t.Fatalf("address Finding projection = %#v, want unresolved blocking addressed disposition", addressedFinding)
		}
		p23ReleaseClaim(t, a, a.actor, mutationClaim, "finding-address-builder")
	})

	t.Run("non_independent_and_blocked_terminal_reviews_reject", func(t *testing.T) {
		builderClaim := p23AcquireClaim(t, a, a.actor, "review", "finding-non-independent")
		builderRun := p23StartRun(t, a, a.actor, builderClaim, "finding-non-independent")
		before := p21JournalEventCount(t, a.workspace)
		response := p23ReviewRecord(t, a, a.actor, "finding-non-independent", p23ReviewInput(a, p21PathString(t, builderRun, "id"), p23Number(t, builderRun["version"]), "pass", []any{}, []any{"FND-resolved"}))
		assertP22ProcessFailure(t, "review.record", "review-not-independent", 4, response)
		if p21JournalEventCount(t, a.workspace) != before {
			t.Fatal("non-independent terminal Review appended an Event")
		}
		p23CloseAndRelease(t, a, a.actor, p21PathString(t, builderRun, "id"), p21PathString(t, builderClaim, "id"), "finding-non-independent")

		blockedClaim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", "finding-blocked")
		blockedRun := p23StartRun(t, a, "actor-reviewer-two", blockedClaim, "finding-blocked")
		defer p23CloseAndRelease(t, a, "actor-reviewer-two", p21PathString(t, blockedRun, "id"), p21PathString(t, blockedClaim, "id"), "finding-blocked")
		before = p21JournalEventCount(t, a.workspace)
		response = p23ReviewRecord(t, a, "actor-reviewer-two", "finding-blocked", p23ReviewInput(a, p21PathString(t, blockedRun, "id"), p23Number(t, blockedRun["version"]), "blocked", []any{p23Finding(a.issueID, "FND-blocked", "minor", false)}, []any{"FND-resolved"}))
		assertP22ProcessFailure(t, "review.record", "review-verdict-invalid", 2, response)
		if p21JournalEventCount(t, a.workspace) != before {
			t.Fatal("blocked terminal Review appended an Event")
		}
	})

	mixedClaim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", "finding-mixed")
	mixedRun := p23StartRun(t, a, "actor-reviewer-two", mixedClaim, "finding-mixed")
	mixedInput := p23ReviewInput(a, p21PathString(t, mixedRun, "id"), p23Number(t, mixedRun["version"]), "pass", []any{}, []any{"FND-resolved"})
	mixedInput["inapplicable_findings"] = []any{map[string]any{"finding_id": "FND-inapplicable", "reason": "obsolete after the correction"}}
	beforeMixed := p21JournalEventCount(t, a.workspace)
	mixed := p23ReviewRecord(t, a, "actor-reviewer-two", "finding-mixed", mixedInput)
	p23AssertReviewSuccess(t, "mixed terminal Review", mixed)
	p23AssertActivity(t, a.workspace, mixed, "finding-resolved", a.issueID, 1)
	if after := p21JournalEventCount(t, a.workspace); after != beforeMixed+1 {
		t.Fatalf("mixed terminal Review changed Event count %d -> %d, want one atomic Event", beforeMixed, after)
	}
	mixedReview := p23ReviewFromResponse(t, mixed)
	if !reflect.DeepEqual(mixedReview["resolved_finding_ids"], []any{"FND-resolved"}) {
		t.Fatalf("mixed Review result = %#v, want resolved-only IDs", mixedReview)
	}

	t.Run("changed_inapplicability_reason_conflicts_without_an_event", func(t *testing.T) {
		before := p21JournalEventCount(t, a.workspace)
		changed := p23Clone(mixedInput)
		changed["inapplicable_findings"] = []any{map[string]any{"finding_id": "FND-inapplicable", "reason": "a different reason"}}
		response := p23ReviewRecord(t, a, "actor-reviewer-two", "finding-mixed", changed)
		assertP22ProcessFailure(t, "review.record", "idempotency-conflict", 4, response)
		if p21JournalEventCount(t, a.workspace) != before {
			t.Fatal("changed inapplicability reason appended an Event")
		}
	})

	resolved := p23FindingByID(t, a, "FND-resolved")
	inapplicable := p23FindingByID(t, a, "FND-inapplicable")
	if resolved["state"] != "resolved" || !p23HasDisposition(t, resolved, "resolved", "actor-reviewer-two", mixedReview["id"].(string), "resolved by qualifying Review") {
		t.Fatalf("resolved Finding = %#v, want qualifying Review disposition", resolved)
	}
	if inapplicable["state"] != "inapplicable" || !p23HasDisposition(t, inapplicable, "inapplicable", "actor-reviewer-two", mixedReview["id"].(string), "obsolete after the correction") {
		t.Fatalf("inapplicable Finding = %#v, want exact reason and qualifying Review", inapplicable)
	}
	if !reflect.DeepEqual(p23ReviewGet(t, a, mixedReview["id"].(string)), mixedReview) {
		t.Fatal("mixed Review get result differs from its recorded typed result")
	}
	beforeRebuildReview := p23ReviewGet(t, a, mixedReview["id"].(string))
	beforeRebuildFindings := p23FindingList(t, a)
	p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "rebuild")
	if after := p23ReviewGet(t, a, mixedReview["id"].(string)); !reflect.DeepEqual(after, beforeRebuildReview) {
		t.Fatalf("mixed Review changed across rebuild: before=%#v after=%#v", beforeRebuildReview, after)
	}
	if after := p23FindingList(t, a); !reflect.DeepEqual(after, beforeRebuildFindings) {
		t.Fatalf("Finding list changed across rebuild: before=%#v after=%#v", beforeRebuildFindings, after)
	}

	p23CloseAndRelease(t, a, "actor-reviewer-two", p21PathString(t, mixedRun, "id"), p21PathString(t, mixedClaim, "id"), "finding-mixed")
	sameClaim := p23AcquireClaim(t, a, "actor-denied", "review", "finding-same-reviewer")
	sameRun := p23StartRun(t, a, "actor-denied", sameClaim, "finding-same-reviewer")
	sameInput := p23ReviewInput(a, p21PathString(t, sameRun, "id"), p23Number(t, sameRun["version"]), "pass", []any{}, []any{"FND-same-reviewer"})
	sameInput["inapplicable_findings"] = []any{map[string]any{"finding_id": "FND-same-inapplicable", "reason": "no longer applies after rework"}}
	same := p23ReviewRecord(t, a, "actor-denied", "finding-same-reviewer", sameInput)
	p23AssertReviewSuccess(t, "same independent Reviewer terminal Review", same)
	sameReview := p23ReviewFromResponse(t, same)
	if finding := p23FindingByID(t, a, "FND-same-reviewer"); finding["state"] != "resolved" || !p23HasDisposition(t, finding, "resolved", "actor-denied", sameReview["id"].(string), "resolved by qualifying Review") {
		t.Fatalf("same Reviewer terminal Finding = %#v", finding)
	}
	if finding := p23FindingByID(t, a, "FND-same-inapplicable"); finding["state"] != "inapplicable" || !p23HasDisposition(t, finding, "inapplicable", "actor-denied", sameReview["id"].(string), "no longer applies after rework") {
		t.Fatalf("same Reviewer inapplicable Finding = %#v", finding)
	}
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, sameRun, "id"), p21PathString(t, sameClaim, "id"), "finding-same-reviewer")

	t.Run("second_terminal_and_unknown_findings_are_non_mutating", func(t *testing.T) {
		secondClaim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", "finding-second-terminal")
		secondRun := p23StartRun(t, a, "actor-reviewer-two", secondClaim, "finding-second-terminal")
		before := p21JournalEventCount(t, a.workspace)
		response := p23ReviewRecord(t, a, "actor-reviewer-two", "finding-second-terminal", p23ReviewInput(a, p21PathString(t, secondRun, "id"), p23Number(t, secondRun["version"]), "pass", []any{}, []any{"FND-resolved"}))
		assertP22ProcessFailure(t, "review.record", "finding-already-disposed", 4, response)
		if !strings.Contains(response.Stderr, `"finding_id":"FND-resolved"`) || !strings.Contains(response.Stderr, `"state":"resolved"`) || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("second terminal response = %#v, want exact no-Event terminal details", response)
		}
		secondInapplicable := p23ReviewInput(a, p21PathString(t, secondRun, "id"), p23Number(t, secondRun["version"]), "pass", []any{}, nil)
		secondInapplicable["inapplicable_findings"] = []any{map[string]any{"finding_id": "FND-inapplicable", "reason": "a second terminal disposition"}}
		response = p23ReviewRecord(t, a, "actor-reviewer-two", "finding-second-inapplicable", secondInapplicable)
		assertP22ProcessFailure(t, "review.record", "finding-already-disposed", 4, response)
		if !strings.Contains(response.Stderr, `"finding_id":"FND-inapplicable"`) || !strings.Contains(response.Stderr, `"state":"inapplicable"`) || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("second inapplicable response = %#v, want exact no-Event terminal details", response)
		}
		unknown := p23ReviewRecord(t, a, "actor-reviewer-two", "finding-unknown", p23ReviewInput(a, p21PathString(t, secondRun, "id"), p23Number(t, secondRun["version"]), "pass", []any{}, []any{"FND-missing"}))
		assertP22ProcessFailure(t, "review.record", "finding-not-found", 3, unknown)
		if !strings.Contains(unknown.Stderr, `"id":"FND-missing"`) || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("unknown Finding response = %#v, want exact no-Event not-found details", unknown)
		}
	})
}

func p23FindingAddress(t *testing.T, a *p22FIDRecoveryAdapter, actor, key string, input map[string]any) p22ConformanceResponse {
	t.Helper()
	return p23RunInput(t, a, "finding", "address", actor, key, input)
}

func p23ReleaseClaim(t *testing.T, a *p22FIDRecoveryAdapter, actor string, claim map[string]any, key string) {
	t.Helper()
	release := p23RunInput(t, a, "claim", "release", actor, key+"-release", map[string]any{"claim_id": claim["id"], "expected_version": claim["version"], "reason": "completed"})
	p23AssertReviewSuccess(t, "release Claim", release)
}

func p23FindingByID(t *testing.T, a *p22FIDRecoveryAdapter, id string) map[string]any {
	t.Helper()
	for _, finding := range p23FindingList(t, a) {
		if finding["id"] == id {
			return finding
		}
	}
	t.Fatalf("finding list omitted %q", id)
	return nil
}

func p23FindingList(t *testing.T, a *p22FIDRecoveryAdapter) []map[string]any {
	t.Helper()
	response := a.p22Run(t, "finding", "list", "--issue", a.issueID, "--limit", "64")
	p23AssertReviewSuccess(t, "finding list", response)
	items := p21Path(t, p23Envelope(t, response.Stdout), "result", "items").([]any)
	findings := make([]map[string]any, len(items))
	for index, value := range items {
		findings[index] = value.(map[string]any)
	}
	return findings
}

func p23HasDisposition(t *testing.T, finding map[string]any, action, actor, reviewID, reason string) bool {
	t.Helper()
	for _, value := range finding["dispositions"].([]any) {
		disposition := value.(map[string]any)
		qualifyingReviewMatches := disposition["qualifying_review_id"] == reviewID
		if reviewID == "" {
			qualifyingReviewMatches = disposition["qualifying_review_id"] == nil
		}
		if disposition["action"] == action && disposition["actor_id"] == actor && qualifyingReviewMatches && disposition["reason"] == reason {
			return true
		}
	}
	return false
}
