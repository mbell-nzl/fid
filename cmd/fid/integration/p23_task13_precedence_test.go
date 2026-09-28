//go:build integration

package integration

import (
	"reflect"
	"testing"
)

// TestP23FailurePrecedenceWithProtectedDecisionAndObligation proves the
// declared authority denial outranks an unsatisfied guarded Obligation.
func TestP23FailurePrecedenceWithProtectedDecisionAndObligation(t *testing.T) {
	t.Parallel()
	decision := map[string]any{"ordinal": 0, "decision_class": "security", "required_capabilities": []string{"decision.record-reference"}, "affected_actions": []string{"issue.transition-terminal-success"}}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, p23ObligationGuardActions(), p22ExecutionFixture{ProtectedDecisions: []any{decision}})
	p23AdvanceToQAReview(t, a, "task13-protected-precedence")
	p23AssertEvaluationStatuses(t, p23EvidenceEvaluate(t, a), map[string]string{"OBL-p22-corpus": "unsatisfied"})
	beforeEvents, beforeProjection := p21JournalEventCount(t, a.workspace), p22FaultProjection(t, a.workspace)
	denied := p23IssueTransition(t, a, "qa_pass", 4, "task13-protected-and-obligation")
	assertP22ProcessFailure(t, "issue.transition", "protected-decision-required", 5, denied)
	if p21JournalEventCount(t, a.workspace) != beforeEvents || !reflect.DeepEqual(p22FaultProjection(t, a.workspace), beforeProjection) {
		t.Fatal("protected Decision precedence denial changed the journal or projection")
	}
}

// TestP23FailurePrecedenceWithReadinessStopAndObligation opens the readiness
// Stop, then supersedes a previously satisfying Outcome. Both the Stop and an
// additional guarded Obligation are unsatisfied for envelope.complete.
func TestP23FailurePrecedenceWithReadinessStopAndObligation(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "review.record", "run.close", "run.observe", "run.observe-budget", "run.start"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, actions, p22ExecutionFixture{HumanObligations: true, CoexistingReviewObligation: true})
	release := p23RecordHumanOutcome(t, a, "precedence-release-approved", a.actor, "release", "approved")
	assertP23HumanOutcomeSuccess(t, release)
	releaseID := p21PathString(t, p23Envelope(t, release.Stdout), "result", "outcome", "id")
	assertP23HumanOutcomeSuccess(t, p23LinkHumanOutcome(t, a, "precedence-release-link", "OBL-human-release", releaseID))
	gate := p23RunQualifyingGateSession(t, a, "precedence-gate")
	assertP23HumanOutcomeSuccess(t, p23RunInput(t, a, "evidence", "link", a.actor, "precedence-gate-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}}))
	p23SatisfyCoexistingReview(t, a, "precedence-review")
	assertP23HumanOutcomeStopState(t, a, "security", "open")
	rejected := p23RecordHumanOutcome(t, a, "precedence-release-rejected", a.actor, "release", "rejected")
	assertP23HumanOutcomeSuccess(t, rejected)
	assertP23ObligationsUnsatisfied(t, a, "OBL-human-release", "OBL-human-security")
	beforeEvents, beforeProjection := p21JournalEventCount(t, a.workspace), p22FaultProjection(t, a.workspace)
	denied := p23RunInput(t, a, "envelope", "complete", a.actor, "precedence-stop-and-obligation", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1, "repository": a.repository})
	assertP22ProcessFailure(t, "envelope.complete", "human-outcome-required", 4, denied)
	if p21JournalEventCount(t, a.workspace) != beforeEvents || !reflect.DeepEqual(p22FaultProjection(t, a.workspace), beforeProjection) {
		t.Fatal("readiness Stop precedence denial changed the journal or projection")
	}
}
