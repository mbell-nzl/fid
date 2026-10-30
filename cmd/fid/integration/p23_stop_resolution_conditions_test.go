//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/mbell-nzl/fid/internal/evidence"
)

// The tests in this file prove DEC-250 (review-7 BC-3) through the built
// binary. Every Condition type has exactly one qualifying-fact predicate that
// manual `fid stop resolve` and the atomic resolution inside the recording
// Event both use, and resolution is keyed on the Stop's recorded Condition,
// never on the code that raised it. Both a producer-raised Stop and an
// `operator-stop` carrying the Condition are therefore exercised, for each of
// `human-outcome-recorded` and `decision-satisfied` (VAL-048, VAL-050,
// VAL-078). An operator-stop is constructible through the existing
// `fid stop raise` leaf with an accepted code plus the Condition under test.

// TestP23StopResolveHumanOutcomeConditionKeysOnConditionNotRaisingCode covers
// the human-outcome-recorded Condition on a producer's readiness Stop and on
// operator-stop Stops of every subject type, through the manual path and
// through the atomic path.
func TestP23StopResolveHumanOutcomeConditionKeysOnConditionNotRaisingCode(t *testing.T) {
	t.Parallel()
	actions := []string{"attention.resolve", "claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.observe-budget", "run.start", "stop.raise", "stop.resolve"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, actions, p22ExecutionFixture{HumanObligations: true})
	for revision, state := range []string{"in_progress", "dev_complete", "qa_review"} {
		assertP23HumanOutcomeCommandSuccess(t, "transition Issue to "+state, p23IssueTransition(t, a, state, revision+1, "cond-"+state))
	}
	release := p23RecordHumanOutcome(t, a, "cond-release-approved", a.actor, "release", "approved")
	assertP23HumanOutcomeSuccess(t, release)
	releaseEventID := p23EventID(t, release)
	assertP23HumanOutcomeSuccess(t, p23LinkHumanOutcome(t, a, "cond-release-match", "OBL-human-release", p21PathString(t, p23Envelope(t, release.Stdout), "result", "outcome", "id")))

	// The security Outcome is recorded first and never linked as Evidence, so
	// the security Obligation stays unsatisfied and the readiness Stop is
	// still raised once the gate Evidence links. That leaves an open producer
	// Stop beside an already-durable qualifying Outcome.
	securityFirst := p23RecordHumanOutcome(t, a, "cond-security-first", a.actor, "security", "approved")
	assertP23HumanOutcomeSuccess(t, securityFirst)
	securityFirstEventID := p23EventID(t, securityFirst)
	if stops := p21Path(t, p23Envelope(t, securityFirst.Stdout), "result", "stops").([]any); len(stops) != 0 {
		t.Fatalf("Outcome recorded before any Stop resolved %d Stop(s), want 0", len(stops))
	}

	p22AcceptedCall(t, a, "cond-budget", "run.observe-budget", map[string]any{"value": 0})
	gate, err := evidence.RunGate(context.Background(), p22BudgetGateOptions(t, a, "cond-proof"))
	if err != nil || gate.Receipt == nil {
		t.Fatalf("run qualifying gate: result=%#v error=%v", gate, err)
	}
	linkedGate := p23RunInput(t, a, "evidence", "link", a.actor, "cond-gate-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}})
	assertP23HumanOutcomeSuccess(t, linkedGate)
	assertP23HumanOutcomeStopState(t, a, "security", "open")
	producerStopID := p23StopIDByDetail(t, a, "human-outcome-required", "outcome_class", "security")

	envelopeSubject := map[string]any{"type": "envelope", "id": p23EnvelopeID(t, a)}
	runSubject := map[string]any{"type": "run", "id": a.runID}
	activationSubject := map[string]any{"type": "activation", "id": p21PathString(t, a.activation, "id")}
	operatorEnvelope := p23RaiseConditionStopID(t, a, "cond-op-envelope", envelopeSubject, p23HumanOutcomeCondition(envelopeSubject, "security"))
	operatorRun := p23RaiseConditionStopID(t, a, "cond-op-run", runSubject, p23HumanOutcomeCondition(runSubject, "security"))
	operatorActivation := p23RaiseConditionStopID(t, a, "cond-op-activation", activationSubject, p23HumanOutcomeCondition(activationSubject, "security"))
	operatorRelease := p23RaiseConditionStopID(t, a, "cond-op-release", envelopeSubject, p23HumanOutcomeCondition(envelopeSubject, "release"))

	// VAL-078 non-qualifying facts, manual path: an approved Outcome of another
	// class, and an Event that is not an Outcome at all. Each rejects exit 4
	// resolution-invalid, appends no Event, and leaves the Stop open, on the
	// producer's Stop and on the operator-stop alike.
	for _, target := range []struct{ name, stopID string }{{"producer", producerStopID}, {"operator-envelope", operatorEnvelope}, {"operator-run", operatorRun}, {"operator-activation", operatorActivation}} {
		assertP23StopResolveRejected(t, a, "cond-wrong-class-"+target.name, target.stopID, releaseEventID)
		assertP23StopResolveRejected(t, a, "cond-wrong-type-"+target.name, target.stopID, p23EventID(t, linkedGate))
	}

	// Qualifying fact, manual path: the earlier-recorded approved security
	// Outcome is in scope for each subject type and resolves each Stop.
	for _, target := range []struct{ name, stopID string }{{"producer", producerStopID}, {"operator-envelope", operatorEnvelope}, {"operator-run", operatorRun}, {"operator-activation", operatorActivation}} {
		p23ResolveStopWithFact(t, a, "cond-manual-"+target.name, target.stopID, securityFirstEventID)
	}
	assertP23HumanOutcomeStopState(t, a, "security", "resolved")
	assertP23StopOpen(t, a, operatorRelease, true)

	// `stop resolve` leaves the Attention open. `attention.resolve` (its input
	// shape unchanged, VAL-049) revalidates each linked Stop's fact through the
	// same predicate and closes it.
	for _, stopID := range []string{producerStopID, operatorEnvelope} {
		attentionID, version, dispositionEventID := p23AttentionForStop(t, a, stopID)
		closed := p23RunInput(t, a, "attention", "resolve", a.actor, "cond-attention-"+stopID, map[string]any{"attention_id": attentionID, "expected_version": version, "fact_event_ids": []string{dispositionEventID}})
		p23AssertReviewSuccess(t, "resolve Attention of manually resolved Stop "+stopID, closed)
	}

	// Atomic path: a fresh set of matching operator-stops of every subject
	// type plus a different-Condition Stop. The next approved security Outcome
	// resolves exactly the matching ones in its own Event.
	atomicEnvelope := p23RaiseConditionStopID(t, a, "cond-atomic-envelope", envelopeSubject, p23HumanOutcomeCondition(envelopeSubject, "security"))
	atomicRun := p23RaiseConditionStopID(t, a, "cond-atomic-run", runSubject, p23HumanOutcomeCondition(runSubject, "security"))
	atomicActivation := p23RaiseConditionStopID(t, a, "cond-atomic-activation", activationSubject, p23HumanOutcomeCondition(activationSubject, "security"))
	unrelated := p23RaiseConditionStopID(t, a, "cond-atomic-unrelated", envelopeSubject, map[string]any{"type": "review-recorded", "subject": envelopeSubject, "budget_name": nil, "capability_id": nil})
	securitySecond := p23RecordHumanOutcome(t, a, "cond-security-second", a.actor, "security", "approved")
	assertP23HumanOutcomeSuccess(t, securitySecond)
	securitySecondEventID := p23EventID(t, securitySecond)
	resolved := p21Path(t, p23Envelope(t, securitySecond.Stdout), "result", "stops").([]any)
	if len(resolved) != 3 {
		t.Fatalf("qualifying Outcome resolved %d Stop(s), want the 3 matching operator-stops (envelope, run, activation)", len(resolved))
	}
	for _, stopID := range []string{atomicEnvelope, atomicRun, atomicActivation} {
		assertP23StopOpen(t, a, stopID, false)
		if got := p23StopResolutionFactEvent(t, a, stopID); got != securitySecondEventID {
			t.Fatalf("Stop %s resolved by Event %q, want the recording Event %q", stopID, got, securitySecondEventID)
		}
	}
	assertP23StopOpen(t, a, unrelated, true)
	assertP23StopOpen(t, a, operatorRelease, true)

	// Superseded and wrong-result facts, manual path: a rejected Outcome
	// supersedes the approved ones (section 7), so neither the rejected Event
	// (wrong result) nor the now-superseded approved Event qualifies.
	late := p23RaiseConditionStopID(t, a, "cond-late-operator", envelopeSubject, p23HumanOutcomeCondition(envelopeSubject, "security"))
	securityRejected := p23RecordHumanOutcome(t, a, "cond-security-rejected", a.actor, "security", "rejected")
	assertP23HumanOutcomeSuccess(t, securityRejected)
	assertP23StopResolveRejected(t, a, "cond-wrong-result", late, p23EventID(t, securityRejected))
	assertP23StopResolveRejected(t, a, "cond-superseded", late, securitySecondEventID)
	assertP23StopOpen(t, a, unrelated, true)
}

// TestP23StopResolveDecisionSatisfiedConditionKeysOnConditionNotRaisingCode
// covers the decision-satisfied Condition on a producer's protected-Decision
// Stop and on operator-stop Stops, through the manual and atomic paths, and
// proves `stop raise` rejects the same Condition on an activation subject.
func TestP23StopResolveDecisionSatisfiedConditionKeysOnConditionNotRaisingCode(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "decision.record-reference", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.observe-budget", "run.start", "stop.raise", "stop.resolve"}
	decisions := []any{
		map[string]any{"ordinal": 0, "decision_class": "release", "required_capabilities": []string{"decision.record-reference"}, "affected_actions": []string{"envelope.complete"}},
		map[string]any{"ordinal": 1, "decision_class": "security", "required_capabilities": []string{"decision.record-reference"}, "affected_actions": []string{"issue.transition-terminal-success"}},
	}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, p22ExecutionFixture{ProtectedDecisions: decisions})
	const human = "actor-human-outcome"
	envelopeSubject := map[string]any{"type": "envelope", "id": p23EnvelopeID(t, a)}
	runSubject := map[string]any{"type": "run", "id": a.runID}
	activationSubject := map[string]any{"type": "activation", "id": p21PathString(t, a.activation, "id")}

	// R3: an activation subject has no Envelope scope, so the Condition could
	// never qualify. `stop raise` rejects it with no Event.
	beforeRaise := p21JournalEventCount(t, a.workspace)
	rejectedRaise := p23RaiseAttentionCondition(t, a, "dec-activation-raise", activationSubject, p23DecisionCondition(activationSubject, "release"))
	assertP22ProcessFailure(t, "stop.raise", "invalid-argument", 2, rejectedRaise)
	if after := p21JournalEventCount(t, a.workspace); after != beforeRaise {
		t.Fatalf("rejected activation-subject decision-satisfied raise appended Event(s): %d -> %d", beforeRaise, after)
	}

	producerRelease := p23StopIDByDetail(t, a, "protected-decision-required", "decision_class", "release")
	producerSecurity := p23StopIDByDetail(t, a, "protected-decision-required", "decision_class", "security")
	operatorRelease := p23RaiseConditionStopID(t, a, "dec-op-release", envelopeSubject, p23DecisionCondition(envelopeSubject, "release"))
	operatorSecurity := p23RaiseConditionStopID(t, a, "dec-op-security", envelopeSubject, p23DecisionCondition(envelopeSubject, "security"))
	operatorSecurityRun := p23RaiseConditionStopID(t, a, "dec-op-security-run", runSubject, p23DecisionCondition(runSubject, "security"))
	unrelated := p23RaiseConditionStopID(t, a, "dec-unrelated", envelopeSubject, map[string]any{"type": "review-recorded", "subject": envelopeSubject, "budget_name": nil, "capability_id": nil})

	// VAL-078 non-qualifying facts, manual path, release class: a rejected and
	// a deferred reference, an approved reference with the wrong authority, an
	// approved Outcome of another class, a rejected Outcome, and the earlier
	// half of the pair on its own. None satisfies the protected Decision.
	source := p23DecisionSource(a, true)
	rejectedRef := p23RecordDecision(t, a, "dec-release-rejected", "release", "rejected", []string{"decision.record-reference"}, source)
	assertP23DecisionSuccess(t, rejectedRef, 0)
	deferredRef := p23RecordDecision(t, a, "dec-release-deferred", "release", "deferred", []string{"decision.record-reference"}, source)
	assertP23DecisionSuccess(t, deferredRef, 0)
	wrongAuthority := p23RecordDecision(t, a, "dec-release-wrong-authority", "release", "approved", []string{"stop.resolve"}, source)
	assertP23DecisionSuccess(t, wrongAuthority, 0)
	wrongClassOutcome := p23RecordHumanOutcome(t, a, "dec-unrelated-class", human, "unrelated-class", "approved")
	assertP23DecisionSuccess(t, wrongClassOutcome, 0)
	rejectedOutcome := p23RecordHumanOutcome(t, a, "dec-release-rejected-outcome", human, "release", "rejected")
	assertP23DecisionSuccess(t, rejectedOutcome, 0)
	referenceOnly := p23RecordDecision(t, a, "dec-release-approved", "release", "approved", []string{"decision.record-reference"}, source)
	assertP23DecisionSuccess(t, referenceOnly, 0)
	for _, fact := range []struct{ name, eventID string }{
		{"rejected-reference", p23EventID(t, rejectedRef)}, {"deferred-reference", p23EventID(t, deferredRef)},
		{"wrong-authority", p23EventID(t, wrongAuthority)}, {"wrong-class-outcome", p23EventID(t, wrongClassOutcome)},
		{"rejected-outcome", p23EventID(t, rejectedOutcome)}, {"reference-only", p23EventID(t, referenceOnly)},
	} {
		assertP23StopResolveRejected(t, a, "dec-release-producer-"+fact.name, producerRelease, fact.eventID)
		assertP23StopResolveRejected(t, a, "dec-release-operator-"+fact.name, operatorRelease, fact.eventID)
	}

	// Qualifying fact, atomic path: the approved Outcome completes the pair, so
	// its own Event resolves the producer's Stop and the operator-stop, and
	// nothing carrying another Condition.
	releaseOutcome := p23RecordHumanOutcome(t, a, "dec-release-approved-outcome", human, "release", "approved")
	assertP23DecisionSuccess(t, releaseOutcome, 2)
	releaseOutcomeEvent := p23EventID(t, releaseOutcome)
	for _, stopID := range []string{producerRelease, operatorRelease} {
		assertP23StopOpen(t, a, stopID, false)
		if got := p23StopResolutionFactEvent(t, a, stopID); got != releaseOutcomeEvent {
			t.Fatalf("Stop %s resolved by Event %q, want the pair-completing Event %q", stopID, got, releaseOutcomeEvent)
		}
	}
	assertP23StopOpen(t, a, producerSecurity, true)
	assertP23StopOpen(t, a, unrelated, true)

	// Manual path after the pair is complete: only the pair-completing Event
	// (the later of the two) qualifies. The earlier half, the reference, is
	// rejected. A run subject resolves through the Run's Envelope.
	lateOperator := p23RaiseConditionStopID(t, a, "dec-late-operator", envelopeSubject, p23DecisionCondition(envelopeSubject, "release"))
	lateRun := p23RaiseConditionStopID(t, a, "dec-late-run", runSubject, p23DecisionCondition(runSubject, "release"))
	for _, stopID := range []string{lateOperator, lateRun} {
		assertP23StopResolveRejected(t, a, "dec-earlier-half-"+stopID, stopID, p23EventID(t, referenceOnly))
		assertP23StopResolveRejected(t, a, "dec-rejected-ref-"+stopID, stopID, p23EventID(t, rejectedRef))
		p23ResolveStopWithFact(t, a, "dec-pair-completing-"+stopID, stopID, releaseOutcomeEvent)
	}

	// Security, opposite order: the Outcome is recorded first (unsatisfied, so
	// not qualifying), and the approved reference then completes the pair and
	// resolves the producer Stop, the envelope operator-stop, and the run
	// operator-stop (the Run's Envelope) atomically in its own Event.
	securityOutcome := p23RecordHumanOutcome(t, a, "dec-security-approved-outcome", human, "security", "approved")
	assertP23DecisionSuccess(t, securityOutcome, 0)
	for _, stopID := range []string{producerSecurity, operatorSecurity, operatorSecurityRun} {
		assertP23StopResolveRejected(t, a, "dec-security-half-"+stopID, stopID, p23EventID(t, securityOutcome))
	}
	securityReference := p23RecordDecision(t, a, "dec-security-approved", "security", "approved", []string{"decision.record-reference"}, source)
	assertP23DecisionSuccess(t, securityReference, 3)
	for _, stopID := range []string{producerSecurity, operatorSecurity, operatorSecurityRun} {
		assertP23StopOpen(t, a, stopID, false)
		if got := p23StopResolutionFactEvent(t, a, stopID); got != p23EventID(t, securityReference) {
			t.Fatalf("Stop %s resolved by Event %q, want the pair-completing reference Event %q", stopID, got, p23EventID(t, securityReference))
		}
	}
	assertP23StopOpen(t, a, unrelated, true)
}

func p23HumanOutcomeCondition(subject map[string]any, class string) map[string]any {
	return map[string]any{"type": "human-outcome-recorded", "subject": subject, "budget_name": nil, "capability_id": nil, "outcome_class": class}
}

func p23DecisionCondition(subject map[string]any, class string) map[string]any {
	return map[string]any{"type": "decision-satisfied", "subject": subject, "budget_name": nil, "capability_id": nil, "decision_class": class}
}

func p23EventID(t *testing.T, response p22ConformanceResponse) string {
	t.Helper()
	return p21PathString(t, p23Envelope(t, response.Stdout), "event", "id")
}

func p23RaiseConditionStopID(t *testing.T, a *p22FIDRecoveryAdapter, key string, subject, condition map[string]any) string {
	t.Helper()
	response := p23RaiseAttentionCondition(t, a, key, subject, condition)
	p23AssertReviewSuccess(t, "raise operator-stop "+key, response)
	return p21PathString(t, p23Envelope(t, response.Stdout), "result", "stop", "id")
}

func p23StopIDByDetail(t *testing.T, a *p22FIDRecoveryAdapter, code, field, value string) string {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var id string
	p22FaultOK(t, conn.QueryRow(`SELECT id FROM stops WHERE code=? AND json_extract(details_json, '$.'||?)=? ORDER BY id LIMIT 1`, code, field, value).Scan(&id))
	return id
}

func assertP23StopOpen(t *testing.T, a *p22FIDRecoveryAdapter, stopID string, wantOpen bool) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var resolutions int
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM stop_resolutions WHERE stop_id=?`, stopID).Scan(&resolutions))
	if open := resolutions == 0; open != wantOpen {
		t.Fatalf("Stop %s open = %v, want %v", stopID, open, wantOpen)
	}
}

func p23StopResolutionFactEvent(t *testing.T, a *p22FIDRecoveryAdapter, stopID string) string {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var factEventID string
	p22FaultOK(t, conn.QueryRow(`SELECT fact_event_id FROM stop_resolutions WHERE stop_id=?`, stopID).Scan(&factEventID))
	return factEventID
}

// assertP23StopResolveRejected asserts the DEC-250 rejection contract: exit 4
// resolution-invalid, no Event appended, and the Stop still open.
func assertP23StopResolveRejected(t *testing.T, a *p22FIDRecoveryAdapter, key, stopID, factEventID string) {
	t.Helper()
	before := p21JournalEventCount(t, a.workspace)
	response := p23RunInput(t, a, "stop", "resolve", a.actor, key, map[string]any{"stop_id": stopID, "fact_event_id": factEventID})
	assertP22ProcessFailure(t, "stop.resolve", "resolution-invalid", 4, response)
	if after := p21JournalEventCount(t, a.workspace); after != before {
		t.Fatalf("%s: nonqualifying resolve appended Event(s): %d -> %d", key, before, after)
	}
	assertP23StopOpen(t, a, stopID, true)
}

func p23ResolveStopWithFact(t *testing.T, a *p22FIDRecoveryAdapter, key, stopID, factEventID string) {
	t.Helper()
	before := p21JournalEventCount(t, a.workspace)
	response := p23RunInput(t, a, "stop", "resolve", a.actor, key, map[string]any{"stop_id": stopID, "fact_event_id": factEventID})
	p23AssertReviewSuccess(t, "resolve Stop "+key, response)
	if after := p21JournalEventCount(t, a.workspace); after != before+1 {
		t.Fatalf("%s: qualifying resolve appended %d Event(s), want 1", key, after-before)
	}
	assertP23StopOpen(t, a, stopID, false)
	if got := p23StopResolutionFactEvent(t, a, stopID); got != factEventID {
		t.Fatalf("%s: Stop resolution fact = %q, want %q", key, got, factEventID)
	}
}
