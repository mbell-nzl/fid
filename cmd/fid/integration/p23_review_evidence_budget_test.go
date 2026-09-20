//go:build integration

package integration

import (
	"strings"
	"testing"
)

func TestP23ReviewEvidenceRealProcessJourney(t *testing.T) {
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "evidence.link", "gate.register", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation", p22ExecutionFixture{ReviewObligation: true})
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "evidence-builder")

	changesClaim := p23AcquireClaim(t, a, "actor-denied", "review", "evidence-changes")
	changesRun := p23StartRun(t, a, "actor-denied", changesClaim, "evidence-changes")
	changes := p23ReviewRecord(t, a, "actor-denied", "evidence-changes", p23ReviewInput(a, p21PathString(t, changesRun, "id"), p23Number(t, changesRun["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-evidence", "major", true)}, nil))
	p23AssertReviewSuccess(t, "changes-required Review", changes)
	changesReview := p23ReviewFromResponse(t, changes)
	before := p21JournalEventCount(t, a.workspace)
	changesLink := p23EvidenceLink(t, a, "evidence-changes", changesReview["id"].(string))
	assertP22ProcessFailure(t, "evidence.link", "review-changes-required", 4, changesLink)
	if p21JournalEventCount(t, a.workspace) != before || !p23HasJSON(changesLink.Stderr, `"FND-evidence"`) {
		t.Fatalf("changes-required Evidence link = %#v, want typed no-Event denial naming blocker", changesLink)
	}
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, changesRun, "id"), p21PathString(t, changesClaim, "id"), "evidence-changes")

	blockedClaim := p23AcquireClaim(t, a, "actor-denied", "review", "evidence-blocked")
	blockedRun := p23StartRun(t, a, "actor-denied", blockedClaim, "evidence-blocked")
	blocked := p23ReviewRecord(t, a, "actor-denied", "evidence-blocked", p23ReviewInput(a, p21PathString(t, blockedRun, "id"), p23Number(t, blockedRun["version"]), "blocked", []any{p23Finding(a.issueID, "FND-blocked", "minor", false)}, nil))
	p23AssertReviewSuccess(t, "blocked Review", blocked)
	blockedReview := p23ReviewFromResponse(t, blocked)
	before = p21JournalEventCount(t, a.workspace)
	blockedLink := p23EvidenceLink(t, a, "evidence-blocked", blockedReview["id"].(string))
	assertP22ProcessFailure(t, "evidence.link", "review-changes-required", 4, blockedLink)
	if p21JournalEventCount(t, a.workspace) != before || !p23HasJSON(blockedLink.Stderr, `"blocking_finding_ids":[]`) {
		t.Fatalf("blocked Evidence link = %#v, want typed no-Event denial with no blocking Finding IDs", blockedLink)
	}
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, blockedRun, "id"), p21PathString(t, blockedClaim, "id"), "evidence-blocked")

	passClaim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", "evidence-pass")
	passRun := p23StartRun(t, a, "actor-reviewer-two", passClaim, "evidence-pass")
	pass := p23ReviewRecord(t, a, "actor-reviewer-two", "evidence-pass", p23ReviewInput(a, p21PathString(t, passRun, "id"), p23Number(t, passRun["version"]), "pass", []any{}, nil))
	p23AssertReviewSuccess(t, "pass Review", pass)
	passReview := p23ReviewFromResponse(t, pass)
	linked := p23EvidenceLink(t, a, "evidence-pass", passReview["id"].(string))
	p23AssertReviewSuccess(t, "pass Review Evidence link", linked)
	if p21Path(t, p23Envelope(t, linked.Stdout), "result", "evidence", "evidence_type") != "review" {
		t.Fatalf("pass Review Evidence = %q, want review evidence", linked.Stdout)
	}
	evaluation := p23EvidenceEvaluate(t, a)
	if !p23HasJSON(evaluation.Stdout, `"code":"finding-unresolved"`) || !p23HasJSON(evaluation.Stdout, `"subject_id":"FND-evidence"`) {
		t.Fatalf("linked pass evaluation = %q, want unresolved blocking Finding", evaluation.Stdout)
	}
	p23CloseAndRelease(t, a, "actor-reviewer-two", p21PathString(t, passRun, "id"), p21PathString(t, passClaim, "id"), "evidence-pass")

	lateClaim := p23AcquireClaim(t, a, "actor-reviewer-two", "mutation", "evidence-late-mutation")
	before = p21JournalEventCount(t, a.workspace)
	late := p23EvidenceLink(t, a, "evidence-late-mutation", passReview["id"].(string))
	assertP22ProcessFailure(t, "evidence.link", "review-not-independent", 4, late)
	if p21JournalEventCount(t, a.workspace) != before || !p23HasJSON(late.Stderr, `"id":"`+lateClaim["id"].(string)+`"`) {
		t.Fatalf("link-time Review independence = %#v, want no-Event Claim conflict", late)
	}
}

func TestP23ReviewCyclesRealProcessJourney(t *testing.T) {
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.close", "run.observe", "run.observe-budget", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation", p22ExecutionFixture{ReviewObligation: true, ReviewCyclesMaximum: 1})
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "cycles-builder")
	blockedClaim := p23AcquireClaim(t, a, "actor-denied", "review", "cycles-blocked")
	blockedRun := p23StartRun(t, a, "actor-denied", blockedClaim, "cycles-blocked")
	blocked := p23ReviewRecord(t, a, "actor-denied", "cycles-blocked", p23ReviewInput(a, p21PathString(t, blockedRun, "id"), p23Number(t, blockedRun["version"]), "blocked", []any{p23Finding(a.issueID, "FND-cycle-blocked", "minor", false)}, nil))
	p23AssertReviewSuccess(t, "unmetered blocked Review", blocked)
	if p21Path(t, p23Envelope(t, blocked.Stdout), "result", "budget") != nil {
		t.Fatalf("blocked Review budget = %q, want no review-cycles charge", blocked.Stdout)
	}
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, blockedRun, "id"), p21PathString(t, blockedClaim, "id"), "cycles-blocked")

	claim := p23AcquireClaim(t, a, "actor-denied", "review", "cycles-review")
	run := p23StartRun(t, a, "actor-denied", claim, "cycles-review")
	input := p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-cycle", "major", true)}, nil)
	before := p21JournalEventCount(t, a.workspace)
	result := p23ReviewRecord(t, a, "actor-denied", "cycles-review", input)
	p23AssertReviewSuccess(t, "exhausting changes-required Review", result)
	if after := p21JournalEventCount(t, a.workspace); after != before+1 {
		t.Fatalf("exhausting Review Event count = %d, want %d", after, before+1)
	}
	resultEnvelope := p23Envelope(t, result.Stdout)
	budget := p21Path(t, resultEnvelope, "result", "budget").(map[string]any)
	if budget["name"] != "review-cycles" || p23Number(t, budget["consumed"]) != 1 || p23Number(t, budget["remaining"]) != 0 || budget["state"] != "exhausted" {
		t.Fatalf("exhausting Review budget = %#v", budget)
	}
	stops := p21Path(t, resultEnvelope, "result", "stops").([]any)
	attention := p21Path(t, resultEnvelope, "result", "attention").([]any)
	if len(stops) != 1 || len(attention) != 1 || stops[0].(map[string]any)["code"] != "review-rework-exhausted" {
		t.Fatalf("exhausting Review Stop/Attention = %#v %#v", stops, attention)
	}
	rebuilt := a.p22Run(t, "rebuild")
	p23AssertReviewSuccess(t, "rebuild after exhausting Review", rebuilt)
	replay := p23ReviewRecord(t, a, "actor-denied", "cycles-review", input)
	p23AssertReviewSuccess(t, "exhausting Review replay", replay)
	if p21Path(t, p23Envelope(t, replay.Stdout), "result", "replayed") != true || p21JournalEventCount(t, a.workspace) != before+1 {
		t.Fatalf("Review cycle replay = %#v, want no second charge", replay)
	}

	beforeObserveEvents := p21JournalEventCount(t, a.workspace)
	beforeObserveSamples := p23ReviewCycleSampleCount(t, a)
	observe := p23RunInput(t, a, "run", "observe-budget", "actor-denied", "cycles-observe", map[string]any{"run_id": p21PathString(t, run, "id"), "expected_version": run["version"], "budget_name": "review-cycles", "source_sample_id": "SMP-review-cycles", "sample_kind": "baseline", "previous_sample_id": nil, "value": 0, "measurement_ref": a.measurement, "observed_at": "2026-09-20T00:00:00Z"})
	assertP22ProcessFailure(t, "run.observe-budget", "invalid-argument", 2, observe)
	if observe.Stdout != "" || p21JournalEventCount(t, a.workspace) != beforeObserveEvents || p23ReviewCycleSampleCount(t, a) != beforeObserveSamples {
		t.Fatalf("observe fid-owned review-cycles = %#v, want typed no-sample denial", observe)
	}
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, run, "id"), p21PathString(t, claim, "id"), "cycles-review")

	beforeDenials := p21JournalEventCount(t, a.workspace)
	beforeClaims, beforeRuns := p23EntityCounts(t, a)
	deniedClaim := p23RunInput(t, a, "claim", "acquire", a.actor, "cycles-denied-claim", map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "worktree_id": a.repository["worktree_id"], "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": "mutation", "ttl_seconds": 60})
	assertP22ProcessFailure(t, "claim.acquire", "review-rework-exhausted", 4, deniedClaim)
	if !p23HasJSON(deniedClaim.Stderr, `"event":null`) || !p23HasJSON(deniedClaim.Stderr, `"replayed":false`) || p21JournalEventCount(t, a.workspace) != beforeDenials {
		t.Fatalf("exhausted Claim denial = %q, want non-committing response", deniedClaim.Stderr)
	}
	deniedRun := p23RunInput(t, a, "run", "start", a.actor, "cycles-denied-run", map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "claim_id": "CLM-would-be-rework", "claim_version": 1, "executor": map[string]any{"id": "executor-would-be-rework", "source": "activated-harness"}, "obligation_id": nil, "repository": a.repository})
	assertP22ProcessFailure(t, "run.start", "review-rework-exhausted", 4, deniedRun)
	afterClaims, afterRuns := p23EntityCounts(t, a)
	if p21JournalEventCount(t, a.workspace) != beforeDenials || afterClaims != beforeClaims || afterRuns != beforeRuns {
		t.Fatalf("exhausted denials changed Event or entity counts: events %d -> %d, claims %d -> %d, runs %d -> %d", beforeDenials, p21JournalEventCount(t, a.workspace), beforeClaims, afterClaims, beforeRuns, afterRuns)
	}
}

func p23EvidenceLink(t *testing.T, a *p22FIDRecoveryAdapter, key, reviewID string) p22ConformanceResponse {
	t.Helper()
	return p23RunInput(t, a, "evidence", "link", "actor-recovery", key, map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "review", "receipt_id": nil, "review_id": reviewID, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
}

func p23EvidenceEvaluate(t *testing.T, a *p22FIDRecoveryAdapter) p22ConformanceResponse {
	t.Helper()
	response := a.p22Run(t, "evidence", "evaluate", "--input", p21Input(t, a.workspace, "p23-evaluate", map[string]any{"obligation_id": nil, "target_repository": a.repository}))
	p23AssertReviewSuccess(t, "evidence evaluate", response)
	return response
}

func p23EntityCounts(t *testing.T, a *p22FIDRecoveryAdapter) (claims, runs int) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	p22FaultOK(t, conn.QueryRow(`SELECT count(*) FROM claims WHERE issue_id=?`, a.issueID).Scan(&claims))
	p22FaultOK(t, conn.QueryRow(`SELECT count(*) FROM runs WHERE issue_id=?`, a.issueID).Scan(&runs))
	return claims, runs
}

func p23ReviewCycleSampleCount(t *testing.T, a *p22FIDRecoveryAdapter) int {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var count int
	p22FaultOK(t, conn.QueryRow(`SELECT count(*) FROM budget_samples WHERE envelope_id=? AND budget_name='review-cycles'`, p23EnvelopeID(t, a)).Scan(&count))
	return count
}

func p23HasJSON(value, want string) bool { return strings.Contains(value, want) }
