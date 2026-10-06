//go:build integration

package integration

import (
	"strings"
	"testing"
)

func TestP23ReviewEvidenceRealProcessJourney(t *testing.T) {
	t.Parallel()
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
	before = p21JournalEventCount(t, a.workspace)
	conflict := p23RunInput(t, a, "evidence", "link", "actor-recovery", "evidence-pass", map[string]any{
		"obligation_id": "OBL-p22-corpus", "evidence_type": "review", "receipt_id": nil,
		"review_id": passReview["id"], "human_outcome_id": nil,
		"target_repository": a.repository, "required_artifact_ids": []string{"ART-conflicting"},
	})
	assertP22ProcessFailure(t, "evidence.link", "idempotency-conflict", 4, conflict)
	if after := p21JournalEventCount(t, a.workspace); after != before {
		t.Fatalf("Evidence key conflict appended Event(s): %d -> %d", before, after)
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
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "human-outcome.record", "review.record", "run.close", "run.observe", "run.observe-budget", "run.start", "stop.resolve"}
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
	humanInput := p21Input(t, a.workspace, "cycles-human-replay", input)
	human, humanStderr, humanErr := closureRunCommand(a.workspace, a.binary, "--project", a.workspace, "review", "record", "--actor", "actor-denied", "--input", humanInput, "--idempotency-key", "cycles-review")
	if humanErr != nil || humanStderr != "" || !strings.Contains(human, "review-rework-exhausted") || !strings.Contains(human, "human-outcome-recorded") {
		t.Fatalf("exhausted Review Human result = %q / %q / %v", human, humanStderr, humanErr)
	}
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
	condition := stops[0].(map[string]any)["resolution_condition"].(map[string]any)
	if condition["type"] != "human-outcome-recorded" || condition["outcome_class"] != "review-rework-extension" || condition["budget_name"] != nil || condition["capability_id"] != nil || condition["subject"].(map[string]any)["id"] != p23EnvelopeID(t, a) {
		t.Fatalf("review rework Stop Condition = %#v", condition)
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
	claimInput := map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": "mutation", "ttl_seconds": 60}
	deniedClaim := p23RunInput(t, a, "claim", "acquire", a.actor, "cycles-denied-claim", claimInput)
	assertP22ProcessFailure(t, "claim.acquire", "review-rework-exhausted", 4, deniedClaim)
	if !p23HasJSON(deniedClaim.Stderr, `"event":null`) || !p23HasJSON(deniedClaim.Stderr, `"replayed":false`) || p21JournalEventCount(t, a.workspace) != beforeDenials {
		t.Fatalf("exhausted Claim denial = %q, want non-committing response", deniedClaim.Stderr)
	}
	p23AssertNoncommittingDenialRetry(t, a, deniedClaim, func() p22ConformanceResponse {
		return p23RunInput(t, a, "claim", "acquire", a.actor, "cycles-denied-claim", claimInput)
	})
	runInput := map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "claim_id": "CLM-would-be-rework", "claim_version": 1, "executor": map[string]any{"id": "executor-would-be-rework", "source": "activated-harness"}, "obligation_id": nil, "expected_repository": p24ExpectedRepository(a.repository)}
	deniedRun := p23RunInput(t, a, "run", "start", a.actor, "cycles-denied-run", runInput)
	assertP22ProcessFailure(t, "run.start", "review-rework-exhausted", 4, deniedRun)
	p23AssertNoncommittingDenialRetry(t, a, deniedRun, func() p22ConformanceResponse {
		return p23RunInput(t, a, "run", "start", a.actor, "cycles-denied-run", runInput)
	})
	afterClaims, afterRuns := p23EntityCounts(t, a)
	if p21JournalEventCount(t, a.workspace) != beforeDenials || afterClaims != beforeClaims || afterRuns != beforeRuns {
		t.Fatalf("exhausted denials changed Event or entity counts: events %d -> %d, claims %d -> %d, runs %d -> %d", beforeDenials, p21JournalEventCount(t, a.workspace), beforeClaims, afterClaims, beforeRuns, afterRuns)
	}
	for _, tc := range []struct{ suffix, class, result, actor string }{
		{"wrong-class", "unrelated", "approved", "actor-human-outcome"},
		{"rejected", "review-rework-extension", "rejected", "actor-human-outcome"},
		{"non-human", "review-rework-extension", "approved", "actor-denied"},
	} {
		prior := p21JournalEventCount(t, a.workspace)
		outcome := p23RecordHumanOutcome(t, a, "rework-"+tc.suffix, tc.actor, tc.class, tc.result)
		if tc.suffix == "non-human" {
			assertP22ProcessFailure(t, "human-outcome.record", "authority-missing", 5, outcome)
			if p21JournalEventCount(t, a.workspace) != prior {
				t.Fatal("non-human Outcome appended an Event")
			}
		} else {
			p23AssertReviewSuccess(t, "non-qualifying Outcome", outcome)
		}
		assertP22ProcessFailure(t, "claim.acquire", "review-rework-exhausted", 4, p23RunInput(t, a, "claim", "acquire", a.actor, "cycles-still-denied-"+tc.suffix, claimInput))
		assertP22ProcessFailure(t, "run.start", "review-rework-exhausted", 4, p23RunInput(t, a, "run", "start", a.actor, "cycles-run-still-denied-"+tc.suffix, runInput))
	}
	other := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "issue", "create", "--kind", "task", "--title", "unselected rework issue", "--actor", "agent:developer", "--idempotency-key", "cycles-other-issue")
	otherID := p21PathString(t, other, "result", "issue", "id")
	outOfScope := p23RunInput(t, a, "human-outcome", "record", "actor-human-outcome", "cycles-out-of-scope", map[string]any{"id": "OUT-cycles-out-of-scope", "class": "review-rework-extension", "subject": map[string]any{"type": "issue", "id": otherID}, "result": "approved", "scope": map[string]any{"issue_id": otherID}, "reasoning": "approval for a different Issue", "references": []any{}})
	p23AssertReviewSuccess(t, "out-of-scope Outcome", outOfScope)
	assertP22ProcessFailure(t, "claim.acquire", "review-rework-exhausted", 4, p23RunInput(t, a, "claim", "acquire", a.actor, "cycles-still-denied-out-of-scope", claimInput))
	assertP22ProcessFailure(t, "run.start", "review-rework-exhausted", 4, p23RunInput(t, a, "run", "start", a.actor, "cycles-run-still-denied-out-of-scope", runInput))
	approved := p23RecordHumanOutcome(t, a, "rework-approved", "actor-human-outcome", "review-rework-extension", "approved")
	p23AssertReviewSuccess(t, "qualifying rework Outcome", approved)
	resolved := p21Path(t, p23Envelope(t, approved.Stdout), "result", "stops").([]any)
	if len(resolved) != 1 || resolved[0].(map[string]any)["id"] != stops[0].(map[string]any)["id"] || resolved[0].(map[string]any)["resolution"] == nil {
		t.Fatalf("Outcome did not resolve rework Stop atomically: %#v", resolved)
	}
	reworkClaim := p23AcquireClaim(t, a, a.actor, "mutation", "cycles-after-approval")
	reworkRun := p23StartRun(t, a, a.actor, reworkClaim, "cycles-after-approval")
	p23CloseAndRelease(t, a, a.actor, p21PathString(t, reworkRun, "id"), p21PathString(t, reworkClaim, "id"), "cycles-after-approval")
	for i := 0; i < 2; i++ {
		suffix := "cycles-later-" + string(rune('a'+i))
		reviewClaim := p23AcquireClaim(t, a, "actor-denied", "review", suffix)
		reviewRun := p23StartRun(t, a, "actor-denied", reviewClaim, suffix)
		later := p23ReviewRecord(t, a, "actor-denied", suffix, p23ReviewInput(a, p21PathString(t, reviewRun, "id"), p23Number(t, reviewRun["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-"+suffix, "major", true)}, nil))
		p23AssertReviewSuccess(t, "later changes-required Review", later)
		value := p21Path(t, p23Envelope(t, later.Stdout), "result", "budget").(map[string]any)
		if p23Number(t, value["consumed"]) != i+2 || p23Number(t, value["remaining"]) != 0 || len(p21Path(t, p23Envelope(t, later.Stdout), "result", "stops").([]any)) != 0 {
			t.Fatalf("later Review budget/Stops = %#v", p23Envelope(t, later.Stdout))
		}
		p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, reviewRun, "id"), p21PathString(t, reviewClaim, "id"), suffix)
	}
}

func TestP23ReviewReworkManualResolutionWithPriorOutcome(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "human-outcome.record", "review.record", "run.close", "run.observe", "run.start", "stop.resolve"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation", p22ExecutionFixture{ReviewObligation: true, ReviewCyclesMaximum: 1})
	prior := p23RecordHumanOutcome(t, a, "rework-prior", "actor-human-outcome", "review-rework-extension", "approved")
	p23AssertReviewSuccess(t, "prior qualifying Outcome", prior)
	factID := p21PathString(t, p23Envelope(t, prior.Stdout), "event", "id")
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "manual-builder")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "manual-review")
	run := p23StartRun(t, a, "actor-denied", claim, "manual-review")
	recorded := p23ReviewRecord(t, a, "actor-denied", "manual-review", p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-manual", "major", true)}, nil))
	p23AssertReviewSuccess(t, "manual exhaustion Review", recorded)
	stop := p21Path(t, p23Envelope(t, recorded.Stdout), "result", "stops").([]any)[0].(map[string]any)
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, run, "id"), p21PathString(t, claim, "id"), "manual-review")
	resolved := p23RunInput(t, a, "stop", "resolve", a.actor, "manual-rework-resolve", map[string]any{"stop_id": stop["id"], "fact_event_id": factID})
	p23AssertReviewSuccess(t, "manual rework Stop resolution", resolved)
	claim = p23AcquireClaim(t, a, a.actor, "mutation", "manual-after-approval")
	_ = p23StartRun(t, a, a.actor, claim, "manual-after-approval")
}

// TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild is the VAL-027
// linked-recovery leg (review-7 MF-4): one changes-required Review charges one
// review-cycles sample keyed on the Review ID, and neither run.recover into a
// linked successor nor fid rebuild adds, removes, or changes that charge.
func TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.expire", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.observe", "run.recover", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "review", p22ExecutionFixture{ReviewObligation: true, ReviewCyclesMaximum: 3})
	priorRunID := a.runID
	run := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "run", "get", priorRunID)
	a.runVersion = p23Number(t, p21Path(t, run, "result", "run", "version"))

	recorded := p23ReviewRecord(t, a, a.actor, "cycles-recovery-review", p23ReviewInput(a, priorRunID, a.runVersion, "changes-required", []any{p23Finding(a.issueID, "FND-cycle-recovery", "major", true)}, nil))
	p23AssertReviewSuccess(t, "changes-required Review in the review Run", recorded)
	reviewID := p23ReviewFromResponse(t, recorded)["id"].(string)
	charged := p23ReviewCycleSamples(t, a)
	if len(charged) != 1 || charged[0].SourceSampleID != reviewID || charged[0].Source != "fid" || charged[0].Value != 1 || charged[0].RunID != priorRunID {
		t.Fatalf("review-cycles samples after recording = %#v, want exactly one fid sample keyed on Review %q with consumed 1 on Run %q", charged, reviewID, priorRunID)
	}

	run = p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "run", "get", priorRunID)
	a.runVersion = p23Number(t, p21Path(t, run, "result", "run", "version"))
	a.p22Observe(t, "terminated")
	a.p22Expire(t)
	recovered := a.p22Recover(t)
	p23AssertRunSuccess(t, "run.recover of the review Run", recovered)
	recoveredEnvelope := p23Envelope(t, recovered.Stdout)
	successorID := p21PathString(t, recoveredEnvelope, "result", "run", "id")
	if p21PathString(t, recoveredEnvelope, "result", "recovery", "prior_run_id") != priorRunID || successorID == priorRunID {
		t.Fatalf("recovery = %q, want a linked successor of Run %q", recovered.Stdout, priorRunID)
	}
	p23AssertReviewCyclesUnchanged(t, "after linked recovery", charged, p23ReviewCycleSamples(t, a))

	rebuilt := a.p22Run(t, "rebuild")
	p23AssertReviewSuccess(t, "rebuild after linked recovery", rebuilt)
	p23AssertReviewCyclesUnchanged(t, "after rebuild", charged, p23ReviewCycleSamples(t, a))
}

type p23ReviewCycleSample struct {
	ID, RunID, Source, SourceSampleID string
	Value                             int
}

func p23ReviewCycleSamples(t *testing.T, a *p22FIDRecoveryAdapter) []p23ReviewCycleSample {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`SELECT id, run_id, source, source_sample_id, value FROM budget_samples WHERE envelope_id=? AND budget_name='review-cycles' ORDER BY id`, p23EnvelopeID(t, a))
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	var samples []p23ReviewCycleSample
	for rows.Next() {
		var sample p23ReviewCycleSample
		p22FaultOK(t, rows.Scan(&sample.ID, &sample.RunID, &sample.Source, &sample.SourceSampleID, &sample.Value))
		samples = append(samples, sample)
	}
	p22FaultOK(t, rows.Err())
	return samples
}

// p23AssertReviewCyclesUnchanged compares the sample count and the consumed sum,
// then the exact rows, so a second charge or a re-keyed charge both fail.
func p23AssertReviewCyclesUnchanged(t *testing.T, stage string, want, got []p23ReviewCycleSample) {
	t.Helper()
	consumed := func(samples []p23ReviewCycleSample) (total int) {
		for _, sample := range samples {
			total += sample.Value
		}
		return total
	}
	if len(got) != len(want) || consumed(got) != consumed(want) {
		t.Fatalf("review-cycles %s = count %d consumed %d, want count %d consumed %d; rows %#v", stage, len(got), consumed(got), len(want), consumed(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("review-cycles %s row %d = %#v, want %#v", stage, i, got[i], want[i])
		}
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
