//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"
)

// p23ObligationGuardActions is the authority the two DEC-249 guard probes need:
// a builder session, the gate runner, Review recording, and both guarded
// actions. Human Outcome recording is deliberately absent because neither probe
// uses a human-outcome Obligation.
func p23ObligationGuardActions() []string {
	return []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "issue.transition", "review.record", "run.close", "run.observe", "run.observe-budget", "run.start"}
}

// TestP23ReviewObligationGuardingBothActionsBlocksBothRealProcess is the
// blocking half of VAL-009 through the built binary (DEC-249). One Envelope
// carries a satisfied gate-evidence Obligation and an unsatisfied review
// Obligation that guards BOTH envelope.complete and
// issue.transition-terminal-success. Both actions are denied because of that
// review Obligation, with no Event and no projection change, while the read
// surface keeps reporting every Obligation.
func TestP23ReviewObligationGuardingBothActionsBlocksBothRealProcess(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, p23ObligationGuardActions(), p22ExecutionFixture{CoexistingReviewObligation: true, ReviewGuardsTerminalSuccess: true})
	assertP23GuardedActionsGuardedBy(t, a, "guard-both-authoring", map[string][]string{"OBL-p22-corpus": {"envelope.complete", "issue.transition-terminal-success"}, "OBL-p23-review": {"envelope.complete", "issue.transition-terminal-success"}})
	p23AdvanceToQAReview(t, a, "guard-both")
	p23SatisfyGateEvidence(t, a, "guard-both")

	// Reads report every Obligation: the gate is satisfied and the review is not.
	evaluated := p23EvidenceEvaluate(t, a)
	p23AssertEvaluationStatuses(t, evaluated, map[string]string{"OBL-p22-corpus": "satisfied", "OBL-p23-review": "unsatisfied"})

	// The terminal transition names the review Obligation through its reason
	// code only (the CLI denial carries codes, not subject IDs); the gate is
	// satisfied, so the single obligation-evidence-absent code can only come
	// from OBL-p23-review. envelope.complete names it by ID.
	assertP23GuardedActionsDenied(t, a, "guard-both", []string{"OBL-p23-review"}, []string{"OBL-p23-review"})
}

// TestP23UnsatisfiedObligationBlocksOnlyTheActionItGuardsRealProcess is the
// contrast half of VAL-009 through the built binary (DEC-249). The review
// Obligation guards only envelope.complete. While it is unsatisfied and the
// gate is satisfied, envelope.complete is denied naming it and
// issue.transition-terminal-success is permitted; the review Obligation then
// keeps blocking completion after the terminal transition.
//
// A permitted terminal transition commits a one-way state change, so this test
// owns a dedicated single-purpose fixture and performs qa_pass exactly once,
// as its last transition.
func TestP23UnsatisfiedObligationBlocksOnlyTheActionItGuardsRealProcess(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, p23ObligationGuardActions(), p22ExecutionFixture{CoexistingReviewObligation: true})
	assertP23GuardedActionsGuardedBy(t, a, "guard-contrast-authoring", map[string][]string{"OBL-p22-corpus": {"envelope.complete", "issue.transition-terminal-success"}, "OBL-p23-review": {"envelope.complete"}})
	p23AdvanceToQAReview(t, a, "guard-contrast")
	p23SatisfyGateEvidence(t, a, "guard-contrast")

	p23AssertEvaluationStatuses(t, p23EvidenceEvaluate(t, a), map[string]string{"OBL-p22-corpus": "satisfied", "OBL-p23-review": "unsatisfied"})

	// envelope.complete is denied by exactly the review Obligation it guards.
	beforeEvents, beforeProjection := p21JournalEventCount(t, a.workspace), p22FaultProjection(t, a.workspace)
	if got := p23CompleteDenialObligationIDs(t, a, "guard-contrast-complete-before"); !reflect.DeepEqual(got, []string{"OBL-p23-review"}) {
		t.Fatalf("envelope.complete blocked by Obligations %v, want exactly [OBL-p23-review]", got)
	}
	if afterEvents := p21JournalEventCount(t, a.workspace); afterEvents != beforeEvents || !reflect.DeepEqual(beforeProjection, p22FaultProjection(t, a.workspace)) {
		t.Fatalf("envelope.complete denial changed the journal (%d -> %d) or the projection", beforeEvents, afterEvents)
	}

	// The unsatisfied review Obligation does not guard the terminal transition,
	// so it is permitted and commits exactly one Event.
	beforeTerminal := p21JournalEventCount(t, a.workspace)
	assertP23HumanOutcomeCommandSuccess(t, "terminal qa_pass with an unsatisfied Obligation that does not guard it", p23IssueTransition(t, a, "qa_pass", 4, "guard-contrast-qa-pass"))
	if afterTerminal := p21JournalEventCount(t, a.workspace); afterTerminal != beforeTerminal+1 {
		t.Fatalf("terminal qa_pass appended %d Event(s), want 1", afterTerminal-beforeTerminal)
	}

	// The review Obligation still blocks the action it does guard, and reads
	// still report it unsatisfied.
	if got := p23CompleteDenialObligationIDs(t, a, "guard-contrast-complete-after"); !reflect.DeepEqual(got, []string{"OBL-p23-review"}) {
		t.Fatalf("envelope.complete after terminal qa_pass blocked by Obligations %v, want exactly [OBL-p23-review]", got)
	}
	p23AssertEvaluationStatuses(t, p23EvidenceEvaluate(t, a), map[string]string{"OBL-p22-corpus": "satisfied", "OBL-p23-review": "unsatisfied"})
}

// p23AdvanceToQAReview walks the Issue through the three nonterminal states.
// Nonterminal transitions evaluate no Obligation.
func p23AdvanceToQAReview(t *testing.T, a *p22FIDRecoveryAdapter, key string) {
	t.Helper()
	for revision, state := range []string{"in_progress", "dev_complete", "qa_review"} {
		assertP23HumanOutcomeCommandSuccess(t, "transition Issue to "+state, p23IssueTransition(t, a, state, revision+1, key+"-"+state))
	}
}

// p23SatisfyGateEvidence runs the registered gate through fid's own runner and
// links its real Receipt as the qualifying gate-evidence for OBL-p22-corpus.
func p23SatisfyGateEvidence(t *testing.T, a *p22FIDRecoveryAdapter, key string) {
	t.Helper()
	gate := p23RunQualifyingGateSession(t, a, key+"-proof")
	linked := p23RunInput(t, a, "evidence", "link", a.actor, key+"-gate-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
	assertP23HumanOutcomeSuccess(t, linked)
}

// assertP23GuardedActionsGuardedBy reads the durable Envelope projection and
// asserts each Obligation's guarded_actions, so a probe cannot silently run
// against a differently guarded fixture.
func assertP23GuardedActionsGuardedBy(t *testing.T, a *p22FIDRecoveryAdapter, key string, want map[string][]string) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`SELECT id, json_extract(obligation_json,'$.guarded_actions') FROM envelope_obligations WHERE envelope_id=? ORDER BY id`, p23EnvelopeID(t, a))
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	got := map[string][]string{}
	for rows.Next() {
		var id, actionsJSON string
		p22FaultOK(t, rows.Scan(&id, &actionsJSON))
		var actions []string
		p22FaultOK(t, json.Unmarshal([]byte(actionsJSON), &actions))
		got[id] = actions
	}
	p22FaultOK(t, rows.Err())
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: Envelope Obligation guarded_actions = %#v, want %#v", key, got, want)
	}
}

// p23AssertEvaluationStatuses asserts the evidence.evaluate read reports every
// Obligation of the Envelope, each with the wanted status.
func p23AssertEvaluationStatuses(t *testing.T, response p22ConformanceResponse, want map[string]string) {
	t.Helper()
	got := map[string]string{}
	for _, item := range p21Path(t, p23Envelope(t, response.Stdout), "result", "evaluations").([]any) {
		evaluation := item.(map[string]any)
		got[evaluation["obligation_id"].(string)] = evaluation["status"].(string)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("evidence.evaluate reported statuses %v, want every Obligation as %v", got, want)
	}
}

// p23CompleteDenialObligationIDs attempts envelope.complete, asserts the
// obligation-unsatisfied denial, and returns the sorted Obligation IDs named by
// its reasons.
func p23CompleteDenialObligationIDs(t *testing.T, a *p22FIDRecoveryAdapter, key string) []string {
	t.Helper()
	complete := p23RunInput(t, a, "envelope", "complete", a.actor, key, map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1})
	if complete.ExitCode != 4 || complete.Stdout != "" {
		t.Fatalf("%s envelope.complete = %#v, want exit 4 with no stdout", key, complete)
	}
	var denial map[string]any
	if err := json.Unmarshal([]byte(complete.Stderr), &denial); err != nil {
		t.Fatalf("decode %s denial: %v; stderr=%q", key, err, complete.Stderr)
	}
	if got := p21PathString(t, denial, "error", "code"); got != "obligation-unsatisfied" {
		t.Fatalf("%s denial code = %q, want obligation-unsatisfied", key, got)
	}
	var blocking []string
	for _, reason := range p21Path(t, denial, "error", "details", "reasons").([]any) {
		if item := reason.(map[string]any); item["subject_type"] == "obligation" {
			blocking = append(blocking, item["subject_id"].(string))
		}
	}
	sort.Strings(blocking)
	return blocking
}
