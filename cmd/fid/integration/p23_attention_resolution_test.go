//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/evidence"
)

// TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess
// completes the closed four-type judgment resolution Condition set (VAL-078,
// VAL-080) for the two types this task registers, review-recorded and
// finding-resolved, and proves attention.resolve keeps its exact existing
// input shape and retains every original raise fact (VAL-049). Neither
// Condition type has a Phase 23 producer, so both Stops are raised through
// the generic fid stop raise leaf with an accepted code, per the task's
// binding constraints; the qualifying fact for both is the real
// review.recorded Event a second independent Review produces when it resolves
// the Finding the first Review created.
func TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess(t *testing.T) {
	t.Parallel()
	actions := []string{
		"attention.resolve", "claim.acquire", "claim.release", "envelope.activate", "envelope.create",
		"gate.register", "review.record", "run.close", "run.observe", "run.start", "stop.raise", "stop.resolve",
	}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation")
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "attn-builder-before-review")

	firstClaim := p23AcquireClaim(t, a, "actor-denied", "review", "attn-review-first")
	firstRun := p23StartRun(t, a, "actor-denied", firstClaim, "attn-review-first")
	firstRunID := p21PathString(t, firstRun, "id")
	findingID := "FND-attn-condition"
	first := p23ReviewRecord(t, a, "actor-denied", "attn-review-first", p23ReviewInput(a, firstRunID, p23Number(t, firstRun["version"]), "changes-required", []any{p23Finding(a.issueID, findingID, "major", true)}, nil))
	p23AssertReviewSuccess(t, "first blocking Review", first)
	firstEventID := p21PathString(t, p23Envelope(t, first.Stdout), "event", "id")

	p23CloseAndRelease(t, a, "actor-denied", firstRunID, p21PathString(t, firstClaim, "id"), "attn-review-first")
	reworkClaim := p23AcquireClaim(t, a, a.actor, "mutation", "attn-builder-rework")
	reworkRun := p23StartRun(t, a, a.actor, reworkClaim, "attn-builder-rework")
	p23CloseAndRelease(t, a, a.actor, p21PathString(t, reworkRun, "id"), p21PathString(t, reworkClaim, "id"), "attn-builder-rework")

	secondClaim := p23AcquireClaim(t, a, "actor-denied", "review", "attn-review-second")
	secondRun := p23StartRun(t, a, "actor-denied", secondClaim, "attn-review-second")
	secondRunID := p21PathString(t, secondRun, "id")
	second := p23ReviewRecord(t, a, "actor-denied", "attn-review-second", p23ReviewInput(a, secondRunID, p23Number(t, secondRun["version"]), "pass", []any{}, []any{findingID}))
	p23AssertReviewSuccess(t, "resolving Review", second)
	secondEventID := p21PathString(t, p23Envelope(t, second.Stdout), "event", "id")

	runSubject := map[string]any{"type": "run", "id": secondRunID}

	// review-recorded: no discriminator field at all (DEC-236).
	reviewCondition := map[string]any{"type": "review-recorded", "subject": runSubject, "budget_name": nil, "capability_id": nil}
	reviewStop := p23RaiseAttentionCondition(t, a, "attn-review-recorded-stop", runSubject, reviewCondition)
	p23AssertReviewSuccess(t, "raise review-recorded Stop", reviewStop)
	reviewStopID := p21PathString(t, p23Envelope(t, reviewStop.Stdout), "result", "stop", "id")

	beforeReviewResolve := p21JournalEventCount(t, a.workspace)
	nonQualifyingReview := p23RunInput(t, a, "stop", "resolve", a.actor, "attn-review-recorded-wrong-fact", map[string]any{"stop_id": reviewStopID, "fact_event_id": firstEventID})
	assertP22ProcessFailure(t, "stop.resolve", "resolution-invalid", 4, nonQualifyingReview)
	if after := p21JournalEventCount(t, a.workspace); after != beforeReviewResolve {
		t.Fatalf("nonqualifying review-recorded resolve appended Event(s): %d -> %d", beforeReviewResolve, after)
	}
	reviewResolved := p23RunInput(t, a, "stop", "resolve", a.actor, "attn-review-recorded-fact", map[string]any{"stop_id": reviewStopID, "fact_event_id": secondEventID})
	p23AssertReviewSuccess(t, "resolve review-recorded Stop", reviewResolved)

	// finding-resolved: exactly one discriminator, finding_id.
	findingCondition := map[string]any{"type": "finding-resolved", "subject": runSubject, "budget_name": nil, "capability_id": nil, "finding_id": findingID}
	findingStop := p23RaiseAttentionCondition(t, a, "attn-finding-resolved-stop", runSubject, findingCondition)
	p23AssertReviewSuccess(t, "raise finding-resolved Stop", findingStop)
	findingStopID := p21PathString(t, p23Envelope(t, findingStop.Stdout), "result", "stop", "id")

	beforeFindingResolve := p21JournalEventCount(t, a.workspace)
	nonQualifyingFinding := p23RunInput(t, a, "stop", "resolve", a.actor, "attn-finding-resolved-wrong-fact", map[string]any{"stop_id": findingStopID, "fact_event_id": firstEventID})
	assertP22ProcessFailure(t, "stop.resolve", "resolution-invalid", 4, nonQualifyingFinding)
	if after := p21JournalEventCount(t, a.workspace); after != beforeFindingResolve {
		t.Fatalf("nonqualifying finding-resolved resolve appended Event(s): %d -> %d", beforeFindingResolve, after)
	}
	findingResolved := p23RunInput(t, a, "stop", "resolve", a.actor, "attn-finding-resolved-fact", map[string]any{"stop_id": findingStopID, "fact_event_id": secondEventID})
	p23AssertReviewSuccess(t, "resolve finding-resolved Stop", findingResolved)

	// VAL-049: attention.resolve keeps its exact existing input shape (the
	// sorted unique set of containing disposition Event IDs) and the original
	// raise facts, Stop links, and deduplication keys remain readable.
	reviewSubjectBefore, reviewCodeBefore, reviewRaisedAtBefore := p23StopIdentity(t, a, reviewStopID)
	reviewAttentionID, reviewVersion, reviewDispositionEventID := p23AttentionForStop(t, a, reviewStopID)

	beforeShapeReject := p21JournalEventCount(t, a.workspace)
	shapeRejected := p23RunInput(t, a, "attention", "resolve", a.actor, "attn-review-resolve-extra-field", map[string]any{"attention_id": reviewAttentionID, "expected_version": reviewVersion, "fact_event_ids": []string{reviewDispositionEventID}, "outcome_id": "OUT-not-a-field"})
	if shapeRejected.ExitCode == 0 || shapeRejected.Stdout != "" {
		t.Fatalf("attention.resolve with an unknown field = %#v, want a closed-schema rejection", shapeRejected)
	}
	if after := p21JournalEventCount(t, a.workspace); after != beforeShapeReject {
		t.Fatalf("rejected attention.resolve shape appended Event(s): %d -> %d", beforeShapeReject, after)
	}

	reviewAttentionResolved := p23RunInput(t, a, "attention", "resolve", a.actor, "attn-review-resolve-attention", map[string]any{"attention_id": reviewAttentionID, "expected_version": reviewVersion, "fact_event_ids": []string{reviewDispositionEventID}})
	p23AssertReviewSuccess(t, "resolve review-recorded Attention", reviewAttentionResolved)

	findingAttentionID, findingVersion, findingDispositionEventID := p23AttentionForStop(t, a, findingStopID)
	findingAttentionResolved := p23RunInput(t, a, "attention", "resolve", a.actor, "attn-finding-resolve-attention", map[string]any{"attention_id": findingAttentionID, "expected_version": findingVersion, "fact_event_ids": []string{findingDispositionEventID}})
	p23AssertReviewSuccess(t, "resolve finding-resolved Attention", findingAttentionResolved)

	reviewSubjectAfter, reviewCodeAfter, reviewRaisedAtAfter := p23StopIdentity(t, a, reviewStopID)
	if reviewSubjectAfter != reviewSubjectBefore || reviewCodeAfter != reviewCodeBefore || reviewRaisedAtAfter != reviewRaisedAtBefore {
		t.Fatalf("resolution altered the original raise facts: subject %q->%q code %q->%q raised_at %q->%q", reviewSubjectBefore, reviewSubjectAfter, reviewCodeBefore, reviewCodeAfter, reviewRaisedAtBefore, reviewRaisedAtAfter)
	}
	if factEventID := p23StopResolutionFact(t, a, reviewStopID); factEventID != secondEventID {
		t.Fatalf("review-recorded Stop resolution fact = %q, want the qualifying Event %q", factEventID, secondEventID)
	}
	if factEventID := p23StopResolutionFact(t, a, findingStopID); factEventID != secondEventID {
		t.Fatalf("finding-resolved Stop resolution fact = %q, want the qualifying Event %q", factEventID, secondEventID)
	}

	// VAL-047: the qualified capability reports the exercised resolution path.
	capabilities := a.p22Run(t, "capabilities")
	if !strings.Contains(capabilities.Stdout, `"capability_id":"attention.v1"`) || !strings.Contains(capabilities.Stdout, `"human_outcome_resolution":true`) {
		t.Fatalf("capabilities = %q, want attention.v1 human_outcome_resolution:true", capabilities.Stdout)
	}
}

// TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen
// proves VAL-048 (a wrong-class and a wrong-result Human Outcome both leave a
// human-outcome-recorded Stop open and resolve nothing) and VAL-050
// (recording a qualifying Outcome resolves exactly the matching Stop, never a
// Stop raised on an unrelated condition).
func TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.observe-budget", "run.start", "stop.raise", "stop.resolve"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, actions, p22ExecutionFixture{HumanObligations: true})
	for revision, state := range []string{"in_progress", "dev_complete", "qa_review"} {
		assertP23HumanOutcomeCommandSuccess(t, "transition Issue to "+state, p23IssueTransition(t, a, state, revision+1, "attn-human-outcome-"+state))
	}
	release := p23RecordHumanOutcome(t, a, "attn-release-approved", a.actor, "release", "approved")
	assertP23HumanOutcomeSuccess(t, release)
	releaseOutcomeID := p21PathString(t, p23Envelope(t, release.Stdout), "result", "outcome", "id")
	assertP23HumanOutcomeSuccess(t, p23LinkHumanOutcome(t, a, "attn-release-match", "OBL-human-release", releaseOutcomeID))

	p22AcceptedCall(t, a, "attn-human-outcome-budget", "run.observe-budget", map[string]any{"value": 0})
	gate, err := evidence.RunGate(context.Background(), p22BudgetGateOptions(t, a, "attn-human-outcome-proof"))
	if err != nil || gate.Receipt == nil {
		t.Fatalf("run qualifying gate: result=%#v error=%v", gate, err)
	}
	linkedGate := p23RunInput(t, a, "evidence", "link", a.actor, "attn-human-outcome-gate-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
	assertP23HumanOutcomeSuccess(t, linkedGate)
	assertP23HumanOutcomeStopState(t, a, "security", "open")

	// An unrelated Stop, raised on a different Condition type entirely
	// (this task's own review-recorded), proves VAL-050's negative: it must
	// stay open when the qualifying "security" Outcome resolves only its own
	// matching Stop.
	envelopeID := p23EnvelopeID(t, a)
	unrelatedCondition := map[string]any{"type": "review-recorded", "subject": map[string]any{"type": "envelope", "id": envelopeID}, "budget_name": nil, "capability_id": nil}
	unrelated := p23RaiseAttentionCondition(t, a, "attn-unrelated-stop", map[string]any{"type": "envelope", "id": envelopeID}, unrelatedCondition)
	assertP23HumanOutcomeSuccess(t, unrelated)
	unrelatedStopID := p21PathString(t, p23Envelope(t, unrelated.Stdout), "result", "stop", "id")

	// DEC-250 R1 (VAL-050): an operator-stop carrying the same
	// human-outcome-recorded Condition as the producer's readiness Stop. Only the
	// recorded Condition decides resolution, never the raising code, so the same
	// qualifying Outcome resolves it too.
	securitySubject := map[string]any{"type": "envelope", "id": envelopeID}
	matchingStopID := p23RaiseConditionStopID(t, a, "attn-matching-operator-stop", securitySubject, p23HumanOutcomeCondition(securitySubject, "security"))
	producerStopID := p23StopIDByDetail(t, a, "human-outcome-required", "outcome_class", "security")

	// VAL-048 negative one: an approved Outcome of an unrelated class never
	// matches the "security" Condition's outcome_class and resolves nothing.
	wrongClass := p23RecordHumanOutcome(t, a, "attn-wrong-class", a.actor, "unrelated-class", "approved")
	assertP23HumanOutcomeSuccess(t, wrongClass)
	if stops := p21Path(t, p23Envelope(t, wrongClass.Stdout), "result", "stops").([]any); len(stops) != 0 {
		t.Fatalf("wrong-class Outcome resolved %d Stop(s), want 0", len(stops))
	}
	assertP23HumanOutcomeStopState(t, a, "security", "open")

	// VAL-048 negative two: a non-approved result for the matching class
	// still leaves the Stop open and resolves nothing.
	wrongResult := p23RecordHumanOutcome(t, a, "attn-wrong-result", a.actor, "security", "rejected")
	assertP23HumanOutcomeSuccess(t, wrongResult)
	if stops := p21Path(t, p23Envelope(t, wrongResult.Stdout), "result", "stops").([]any); len(stops) != 0 {
		t.Fatalf("wrong-result Outcome resolved %d Stop(s), want 0", len(stops))
	}
	assertP23HumanOutcomeStopState(t, a, "security", "open")
	assertP23AttentionStopUnresolved(t, a, unrelatedStopID)

	// VAL-048 through the manual path (DEC-250 R1/R2): `fid stop resolve` with
	// the wrong-class Outcome Event and with the wrong-result Outcome Event each
	// reject exit 4 resolution-invalid, append no Event, and leave both the
	// producer's Stop and the matching operator-stop open.
	for _, target := range []struct{ name, stopID string }{{"producer", producerStopID}, {"operator", matchingStopID}} {
		assertP23StopResolveRejected(t, a, "attn-manual-wrong-class-"+target.name, target.stopID, p23EventID(t, wrongClass))
		assertP23StopResolveRejected(t, a, "attn-manual-wrong-result-"+target.name, target.stopID, p23EventID(t, wrongResult))
	}

	// The matching class and result resolves exactly the Stops whose recorded
	// Condition matches: the producer's Stop and the matching operator-stop, and
	// never the unrelated review-recorded Stop.
	security := p23RecordHumanOutcome(t, a, "attn-security-approved", a.actor, "security", "approved")
	assertP23HumanOutcomeSuccess(t, security)
	if stops := p21Path(t, p23Envelope(t, security.Stdout), "result", "stops").([]any); len(stops) != 2 {
		t.Fatalf("qualifying Outcome resolved %d Stop(s), want exactly the 2 matching-Condition Stops", len(stops))
	}
	assertP23HumanOutcomeStopState(t, a, "security", "resolved")
	assertP23StopOpen(t, a, matchingStopID, false)
	if got := p23StopResolutionFactEvent(t, a, matchingStopID); got != p23EventID(t, security) {
		t.Fatalf("matching operator-stop resolved by Event %q, want the recording Event %q", got, p23EventID(t, security))
	}
	assertP23AttentionStopUnresolved(t, a, unrelatedStopID)
}

func p23RaiseAttentionCondition(t *testing.T, a *p22FIDRecoveryAdapter, key string, subject, condition map[string]any) p22ConformanceResponse {
	t.Helper()
	input := map[string]any{
		"subject": subject, "code": "operator-stop", "blocked_actions": []string{"run.recover"},
		"details":              map[string]any{"reason": "attention resolution condition fixture", "reference": "task-11"},
		"resolution_condition": condition, "resolver_capability": "stop.resolve",
	}
	return p23RunInput(t, a, "stop", "raise", a.actor, key, input)
}

func p23StopIdentity(t *testing.T, a *p22FIDRecoveryAdapter, stopID string) (subject, code, raisedAt string) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	p22FaultOK(t, conn.QueryRow(`SELECT subject_type||':'||subject_id, code, raised_at FROM stops WHERE id=?`, stopID).Scan(&subject, &code, &raisedAt))
	return subject, code, raisedAt
}

func p23StopResolutionFact(t *testing.T, a *p22FIDRecoveryAdapter, stopID string) string {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var factEventID string
	p22FaultOK(t, conn.QueryRow(`SELECT fact_event_id FROM stop_resolutions WHERE stop_id=?`, stopID).Scan(&factEventID))
	return factEventID
}

func p23AttentionForStop(t *testing.T, a *p22FIDRecoveryAdapter, stopID string) (attentionID string, version int, dispositionEventID string) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	p22FaultOK(t, conn.QueryRow(`SELECT a.id, a.version, r.created_event_id FROM attention_items a JOIN attention_stop_links l ON l.attention_id=a.id JOIN stop_resolutions r ON r.stop_id=l.stop_id WHERE l.stop_id=?`, stopID).Scan(&attentionID, &version, &dispositionEventID))
	return attentionID, version, dispositionEventID
}

func assertP23AttentionStopUnresolved(t *testing.T, a *p22FIDRecoveryAdapter, stopID string) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var resolved int
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM stop_resolutions WHERE stop_id=?`, stopID).Scan(&resolved))
	if resolved != 0 {
		t.Fatalf("unrelated Stop %s was resolved, want it to remain open", stopID)
	}
}
