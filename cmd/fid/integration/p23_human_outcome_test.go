//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/mbell-nzl/fid/internal/evidence"
)

// TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess proves
// provenanced Human Outcome recording end to end through the built binary:
// VAL-043 authority admission and its four negatives, VAL-045's deterministic
// non-proof fields, VAL-046's readiness-Stop raising and resolution, VAL-025's
// final nine-cell evidence-type closure, VAL-009's completing state, VAL-010's
// linked-evidence re-verification, and VAL-081's re-confirmed nullable binding.
func TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess(t *testing.T) {
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.observe-budget", "run.start"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, actions, p22ExecutionFixture{HumanObligations: true})
	// issue.transition-terminal-success (qa_pass) evaluates every required
	// Obligation regardless of its own guarded_actions, so it is deferred
	// until every Obligation this journey exercises is satisfied; the
	// non-terminal states below do not evaluate Obligations at all.
	for revision, state := range []string{"in_progress", "dev_complete", "qa_review"} {
		response := p23IssueTransition(t, a, state, revision+1, "human-outcome-"+state)
		assertP23HumanOutcomeCommandSuccess(t, "transition Issue to "+state, response)
	}

	// VAL-043 negative: agent actor type (actor-denied shares actor's exact
	// capabilities but is staged actor_type "agent").
	before := p21JournalEventCount(t, a.workspace)
	agentAttempt := p23RecordHumanOutcome(t, a, "agent-actor", "actor-denied", "release", "approved")
	assertP23AuthorityMissing(t, agentAttempt, before)

	// VAL-043 negative: standalone/unregistered actor. This is the same
	// TrustedActorUnavailableError code path an unfinalized Activation
	// produces, since neither actor ID resolves inside a finalized profile.
	before = p21JournalEventCount(t, a.workspace)
	unregistered := p23RecordHumanOutcome(t, a, "unregistered-actor", "human:unregistered", "release", "approved")
	assertP23AuthorityMissing(t, unregistered, before)

	// VAL-043 negative: missing capability, proven with a dedicated minimal
	// finalized Activation whose only human Actor omits human-outcome.record.
	assertP23AuthorityMissing(t, p23RecordHumanOutcomeMissingCapability(t), 0)

	// VAL-043 valid case + VAL-045: the recorded Outcome and its Evidence
	// never present deterministic proof.
	before = p21JournalEventCount(t, a.workspace)
	release := p23RecordHumanOutcome(t, a, "release-approved", a.actor, "release", "approved")
	assertP23HumanOutcomeSuccess(t, release)
	if after := p21JournalEventCount(t, a.workspace); after != before+1 {
		t.Fatalf("valid Human Outcome appended %d Event(s), want 1", after-before)
	}
	releaseOutcome := p21Path(t, p23Envelope(t, release.Stdout), "result", "outcome").(map[string]any)
	releaseOutcomeID := releaseOutcome["id"].(string)
	if releaseOutcome["actor_id"] != a.actor {
		t.Fatalf("recorded actor_id = %v, want %q", releaseOutcome["actor_id"], a.actor)
	}
	provenance := releaseOutcome["provenance"].(map[string]any)
	if provenance["activation_id"] != a.activation["id"] {
		t.Fatalf("recorded provenance activation_id = %v, want %v", provenance["activation_id"], a.activation["id"])
	}

	// VAL-025: link the real recorded approved Human Outcome to its matching
	// human-outcome Obligation and re-prove the two mismatches with the real
	// producer.
	gateMismatch := p23LinkHumanOutcome(t, a, "release-vs-gate", "OBL-p22-corpus", releaseOutcomeID)
	assertP23OutcomeNotPermitted(t, gateMismatch, "gate-evidence")
	matched := p23LinkHumanOutcome(t, a, "release-match", "OBL-human-release", releaseOutcomeID)
	assertP23HumanOutcomeSuccess(t, matched)
	linkedEnvelope := p23Envelope(t, matched.Stdout)
	linkedEvaluation := p21Path(t, linkedEnvelope, "result", "evaluation").(map[string]any)
	if linkedEvaluation["status"] != "satisfied" || linkedEvaluation["artifact_availability"] != "not-applicable" {
		t.Fatalf("linked Human Outcome evaluation = %#v", linkedEvaluation)
	}
	// VAL-045: the Evidence record itself never presents the Outcome as
	// deterministic or fid-executed proof.
	if got := p21PathString(t, linkedEnvelope, "result", "evidence", "producer"); got != "human" {
		t.Fatalf("linked Evidence producer = %q, want human", got)
	}
	if got := p21PathString(t, linkedEnvelope, "result", "evidence", "trust_source"); got != "activated-human" {
		t.Fatalf("linked Evidence trust_source = %q, want activated-human", got)
	}

	// VAL-046: recording the second class's approved Outcome makes it the
	// last unsatisfied Obligation, and evidence.link (the ordinary DEC-224
	// producer) raises human-outcome-required once the gate-evidence
	// Obligation is also satisfied.
	p22AcceptedCall(t, a, "human-outcome-budget", "run.observe-budget", map[string]any{"value": 0})
	gate, err := evidence.RunGate(context.Background(), p22BudgetGateOptions(t, a, "human-outcome-proof"))
	if err != nil || gate.Receipt == nil {
		t.Fatalf("run qualifying gate: result=%#v error=%v", gate, err)
	}
	beforeReadiness := p21JournalEventCount(t, a.workspace)
	linkedGate := p23RunInput(t, a, "evidence", "link", a.actor, "human-outcome-gate-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
	assertP23HumanOutcomeSuccess(t, linkedGate)
	if after := p21JournalEventCount(t, a.workspace); after != beforeReadiness+1 {
		t.Fatalf("evidence.link appended %d Event(s), want 1", after-beforeReadiness)
	}
	assertP23HumanOutcomeStopState(t, a, "security", "open")

	// VAL-009: envelope.complete stays blocked by human-outcome-required
	// (rank 9) ahead of the generic issue-not-terminal-success reason (rank
	// 19) while the last Obligation (OBL-human-security) is unsatisfied. The
	// mutating Claim and Run stay active here; the terminal Issue transition
	// below still needs them.
	beforeDenial := p21JournalEventCount(t, a.workspace)
	denied := p23RunInput(t, a, "envelope", "complete", a.actor, "human-outcome-complete-denied", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1, "repository": a.repository})
	assertP23HumanOutcomeRequiredDenial(t, denied, "OBL-human-security", "security")
	if after := p21JournalEventCount(t, a.workspace); after != beforeDenial {
		t.Fatalf("human-outcome-required denial appended Event(s): %d -> %d", beforeDenial, after)
	}

	// The security Human Outcome resolves the readiness Stop atomically with
	// its own recording Event (contracts.md section 10/11).
	beforeSecurity := p21JournalEventCount(t, a.workspace)
	security := p23RecordHumanOutcome(t, a, "security-approved", a.actor, "security", "approved")
	assertP23HumanOutcomeSuccess(t, security)
	if after := p21JournalEventCount(t, a.workspace); after != beforeSecurity+1 {
		t.Fatalf("security Human Outcome appended %d Event(s), want 1", after-beforeSecurity)
	}
	securityStops := p21Path(t, p23Envelope(t, security.Stdout), "result", "stops").([]any)
	if len(securityStops) != 1 {
		t.Fatalf("security Human Outcome resolved Stops = %#v, want exactly 1", securityStops)
	}
	assertP23HumanOutcomeStopState(t, a, "security", "resolved")
	securityOutcomeID := p21Path(t, p23Envelope(t, security.Stdout), "result", "outcome").(map[string]any)["id"].(string)
	securityMatched := p23LinkHumanOutcome(t, a, "security-match", "OBL-human-security", securityOutcomeID)
	assertP23HumanOutcomeSuccess(t, securityMatched)

	// VAL-010: re-verify with linked review-equivalent (Human Outcome) Evidence
	// present that artifact_availability stays the distinct not-applicable
	// value, never available/missing.
	securityEvaluation := p21Path(t, p23Envelope(t, securityMatched.Stdout), "result", "evaluation").(map[string]any)
	if securityEvaluation["artifact_availability"] != "not-applicable" {
		t.Fatalf("security evaluation artifact_availability = %v, want not-applicable", securityEvaluation["artifact_availability"])
	}

	assertP23HumanOutcomeCommandSuccess(t, "transition Issue to qa_pass", p23IssueTransition(t, a, "qa_pass", 4, "human-outcome-qa-pass"))
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "human-outcome-close")
	assertP23HumanOutcomeCommandSuccess(t, "complete Envelope after every Obligation satisfied", p23RunInput(t, a, "envelope", "complete", a.actor, "human-outcome-complete", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1, "repository": a.repository}))

	// VAL-025's remaining seven cells are supplied by Task 6 (receipt/review
	// positives and the four residual evidence-type-not-accepted mismatches),
	// reconciled here through the real producer's two Human Outcome cells.
	for _, mismatch := range []struct {
		name         string
		obligationID string
	}{
		{"security-vs-gate", "OBL-p22-corpus"},
	} {
		response := p23LinkHumanOutcome(t, a, mismatch.name, mismatch.obligationID, securityOutcomeID)
		assertP23OutcomeNotPermitted(t, response, "gate-evidence")
	}

	listed := a.p22Run(t, "human-outcome", "list", "--issue", a.issueID)
	if listed.ExitCode != 0 || listed.Stderr != "" {
		t.Fatalf("human-outcome list = %#v", listed)
	}
	var listEnvelope map[string]any
	if err := json.Unmarshal([]byte(listed.Stdout), &listEnvelope); err != nil {
		t.Fatal(err)
	}
	if items := p21Path(t, listEnvelope, "result", "items").([]any); len(items) != 2 {
		t.Fatalf("human-outcome list count = %d, want 2", len(items))
	}

	// VAL-081, re-confirmed: the Phase 22 nullable binding still fires from
	// task 9's per-type Condition dispatch after this task's own row.
	p23AssertHumanOutcomeConditionNullableBindingFires(t, a)
}

// TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation proves DEC-224's
// first named raising Event (VAL-046): activating an Envelope whose sole
// Obligation is already the human-outcome type raises human-outcome-required
// immediately, with no preceding evidence.link, review.record, or
// artifact.verify. The other three raising Events (evidence.linked,
// review.recorded, artifact.availability-observed) share this exact
// evaluation and Stop-raising path
// (workflow.HumanOutcomeReadinessStopEffects /
// humanOutcomeReadinessEffectsForEnvelope): evidence.linked is exercised
// above, review.recorded has its own dedicated fixture in
// TestP23HumanOutcomeRequiredRaisedByReviewRecord (its one producer-specific
// override, reviewObligationReadinessOverride, is not shared with the other
// three), and artifact.availability-observed is exercised in
// cmd/fid/artifact_test.go.
func TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation(t *testing.T) {
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p23-human-outcome-activation@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 23 Human Outcome Activation")
	p22Write(t, workspace, "base.txt", "base\n")
	p22Git(t, workspace, "add", "base.txt")
	p22Git(t, workspace, "commit", "-qm", "base")
	init := p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P23HA")
	projectID := p21PathString(t, init, "result", "project", "id")
	issue := p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Human-only Envelope", "--actor", "agent:developer", "--idempotency-key", "p23-human-activation-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")
	actor := "actor-human-outcome"
	capabilities := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.start"}
	activationID := p21Activation(t, binary, workspace, "prepare", "p23-human-activation-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p23-human-activation", "phase": 21})
	profiles := []any{
		map[string]any{"actor_id": actor, "actor_type": "human", "producer": "activated-harness", "provenance_ref": "activation:p23-human-activation", "session_ref": nil, "capabilities": capabilities},
	}
	profileDigest := p21Component(t, binary, workspace, activationID, "authority-profile", "p23-human-activation-authority", profiles)
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "p23-human-activation-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-p23-human-activation", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	digests := []string{gateDigest, profileDigest}
	sort.Strings(digests)
	p21Activation(t, binary, workspace, "commit", "p23-human-activation-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	p21Activation(t, binary, workspace, "finalize", "p23-human-activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})

	repository := p22CorpusRepository(t, workspace, projectID)
	source := map[string]any{"digest": "sha256:" + p22Zeros(64), "display_label": "Human-only Envelope", "external_id": nil, "path": "project-docs/p23-human-activation.md", "repository": repository, "section": nil, "source_type": "project-file"}
	obligation := map[string]any{"id": "OBL-human-only", "type": "human-outcome", "applicability": "required", "guarded_actions": []string{"envelope.complete"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "review_profile": nil, "outcome_class": "release", "accepted_evidence_types": []string{"human-outcome"}, "accepted_trust_sources": []string{"activated-human"}, "freshness_target": "none", "artifact_policy": "metadata-sufficient", "independence": "human", "required_lane": nil, "source": source}
	envelopeRequest := map[string]any{"outcome": "Human-only readiness", "non_goals": []string{}, "risk_result": map[string]any{"source": source, "result_digest": "sha256:" + p22Ones(64)}, "repository_scope": repository, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{}, "authority_grants": []any{map[string]any{"actions": capabilities, "actor_selector": map[string]any{"actor_id": actor, "required_capabilities": capabilities}, "scope": map[string]any{"issue_id": issueID}}}, "obligations": []any{obligation}, "budgets": []any{}}
	created := p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "create", "--actor", actor, "--input", p21Input(t, workspace, "p23-human-activation-envelope", envelopeRequest), "--idempotency-key", "p23-human-activation-envelope-create")
	envelopeID := p21PathString(t, created, "result", "envelope", "id")

	p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", actor, "--input", p21Input(t, workspace, "p23-human-activation-activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1}), "--idempotency-key", "p23-human-activation-activate")

	// envelope.activate's own result carries only envelope/evaluation
	// (public-interface.md section 4: "existing result; evaluation becomes
	// Evaluation[]"), so the raised Stop and Attention are observed from the
	// projection, exactly as task 9 proved for protected-decision-required.
	conn := p22FaultDB(t, workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var count int
	var detailsJSON string
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*), MAX(details_json) FROM stops WHERE code='human-outcome-required' AND subject_id=?`, envelopeID).Scan(&count, &detailsJSON))
	if count != 1 {
		t.Fatalf("human-outcome-required Stop count for Envelope %q = %d, want exactly 1", envelopeID, count)
	}
	var details struct {
		ObligationID string `json:"obligation_id"`
		OutcomeClass string `json:"outcome_class"`
	}
	p22FaultOK(t, json.Unmarshal([]byte(detailsJSON), &details))
	if details.ObligationID != "OBL-human-only" || details.OutcomeClass != "release" {
		t.Fatalf("human-outcome-required details = %#v, want obligation_id=OBL-human-only outcome_class=release", details)
	}
	var attentionCount int
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM attention_items WHERE code='human-outcome-required' AND subject_id=?`, envelopeID).Scan(&attentionCount))
	if attentionCount != 1 {
		t.Fatalf("human-outcome-required Attention count for Envelope %q = %d, want exactly 1", envelopeID, attentionCount)
	}
}

// TestP23HumanOutcomeRequiredRaisedByReviewRecord proves DEC-224's third
// named raising Event (VAL-046) with a dedicated fixture, rather than resting
// on the shared-core code review alone: review.record itself raises
// human-outcome-required when resolving the envelope's last open blocking
// Finding retroactively satisfies an already-linked review Obligation via
// reviewObligationReadinessOverride - the one override path unique to this
// producer and not exercised by the evidence.linked, artifact.availability-
// observed, or envelope.activated fixtures elsewhere in this file and in
// artifact_test.go. Sequence: an independent Review opens a blocking Finding;
// a second independent pass Review is linked as this Obligation's Evidence
// while that Finding is still open (so the Obligation stays unsatisfied); a
// third Review resolves the Finding, which must raise the Stop atomically
// inside that same Event, leaving the sole remaining human-outcome Obligation
// as its target.
func TestP23HumanOutcomeRequiredRaisedByReviewRecord(t *testing.T) {
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "evidence.link", "gate.register", "human-outcome.record", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, actions, p22ExecutionFixture{ReviewObligation: true, HumanObligations: true})
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "review-raise-builder")

	// Pre-satisfy the security class so exactly the release class remains
	// unsatisfied once the review Obligation resolves.
	security := p23RecordHumanOutcome(t, a, "review-raise-security", a.actor, "security", "approved")
	assertP23HumanOutcomeSuccess(t, security)
	securityID := p21Path(t, p23Envelope(t, security.Stdout), "result", "outcome").(map[string]any)["id"].(string)
	assertP23HumanOutcomeSuccess(t, p23LinkHumanOutcome(t, a, "review-raise-security", "OBL-human-security", securityID))

	// Review 1: changes-required, opens the blocking Finding that stays
	// unresolved while a later pass Review's Evidence is linked.
	blockClaim := p23AcquireClaim(t, a, "actor-denied", "review", "review-raise-block")
	blockRun := p23StartRun(t, a, "actor-denied", blockClaim, "review-raise-block")
	blocking := p23ReviewRecord(t, a, "actor-denied", "review-raise-block", p23ReviewInput(a, p21PathString(t, blockRun, "id"), p23Number(t, blockRun["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-review-raise", "major", true)}, nil))
	p23AssertReviewSuccess(t, "blocking Review", blocking)
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, blockRun, "id"), p21PathString(t, blockClaim, "id"), "review-raise-block")

	// Review 2: an independent pass Review, linked as this Obligation's
	// review Evidence while the blocking Finding is still open, so the
	// Obligation stays unsatisfied (finding-unresolved).
	passClaim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", "review-raise-pass")
	passRun := p23StartRun(t, a, "actor-reviewer-two", passClaim, "review-raise-pass")
	pass := p23ReviewRecord(t, a, "actor-reviewer-two", "review-raise-pass", p23ReviewInput(a, p21PathString(t, passRun, "id"), p23Number(t, passRun["version"]), "pass", []any{}, nil))
	p23AssertReviewSuccess(t, "pass Review", pass)
	passReview := p23ReviewFromResponse(t, pass)
	linked := p23RunInput(t, a, "evidence", "link", "actor-reviewer-two", "review-raise-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "review", "receipt_id": nil, "review_id": passReview["id"], "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
	p23AssertReviewSuccess(t, "linked pass Review Evidence", linked)
	p23CloseAndRelease(t, a, "actor-reviewer-two", p21PathString(t, passRun, "id"), p21PathString(t, passClaim, "id"), "review-raise-pass")

	// Review 3: resolves the still-open blocking Finding. This review.record
	// Event itself must retroactively satisfy the linked review Obligation
	// (reviewObligationReadinessOverride), leaving OBL-human-release the sole
	// unsatisfied Obligation, so human-outcome-required raises atomically
	// inside this same Event - observed from the response, not a later query.
	resolveClaim := p23AcquireClaim(t, a, "actor-denied", "review", "review-raise-resolve")
	resolveRun := p23StartRun(t, a, "actor-denied", resolveClaim, "review-raise-resolve")
	beforeResolve := p21JournalEventCount(t, a.workspace)
	resolve := p23ReviewRecord(t, a, "actor-denied", "review-raise-resolve", p23ReviewInput(a, p21PathString(t, resolveRun, "id"), p23Number(t, resolveRun["version"]), "pass", []any{}, []any{"FND-review-raise"}))
	p23AssertReviewSuccess(t, "resolving Review", resolve)
	if after := p21JournalEventCount(t, a.workspace); after != beforeResolve+1 {
		t.Fatalf("resolving Review Event count = %d, want %d", after, beforeResolve+1)
	}
	resolveEnvelope := p23Envelope(t, resolve.Stdout)
	stops := p21Path(t, resolveEnvelope, "result", "stops").([]any)
	attention := p21Path(t, resolveEnvelope, "result", "attention").([]any)
	if len(stops) != 1 || len(attention) != 1 {
		t.Fatalf("resolving Review stops/attention = %#v %#v, want exactly one raised human-outcome-required Stop", stops, attention)
	}
	stop := stops[0].(map[string]any)
	if stop["code"] != "human-outcome-required" {
		t.Fatalf("resolving Review Stop code = %v, want human-outcome-required", stop["code"])
	}
	details := stop["details"].(map[string]any)
	if details["obligation_id"] != "OBL-human-release" || details["outcome_class"] != "release" {
		t.Fatalf("resolving Review Stop details = %#v, want OBL-human-release/release", details)
	}
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, resolveRun, "id"), p21PathString(t, resolveClaim, "id"), "review-raise-resolve")

	// The security class never became the sole unsatisfied Obligation (the
	// review Obligation was still open when it was linked), so it never
	// raised its own Stop at all - only the release class did, above. The
	// qualifying release Human Outcome resolves that one.
	assertP23HumanOutcomeStopState(t, a, "release", "open")
	release := p23RecordHumanOutcome(t, a, "review-raise-release", a.actor, "release", "approved")
	assertP23HumanOutcomeSuccess(t, release)
	assertP23HumanOutcomeStopState(t, a, "release", "resolved")
}

func p23RecordHumanOutcome(t *testing.T, a *p22FIDRecoveryAdapter, suffix, actor, class, result string) p22ConformanceResponse {
	t.Helper()
	input := map[string]any{
		"id": "OUT-" + suffix, "class": class, "subject": map[string]any{"type": "issue", "id": a.issueID},
		"result": result, "scope": map[string]any{"issue_id": a.issueID}, "reasoning": "real-process Human Outcome " + suffix, "references": []any{},
	}
	stdout, stderr, err := closureRunCommand(a.workspace, a.binary, "--project", a.workspace, "--json", "human-outcome", "record", "--actor", actor, "--input", p21Input(t, a.workspace, "p23-human-outcome-"+suffix, input), "--idempotency-key", "p23-human-outcome-"+suffix)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return response
}

func p23LinkHumanOutcome(t *testing.T, a *p22FIDRecoveryAdapter, suffix, obligationID, outcomeID string) p22ConformanceResponse {
	t.Helper()
	return p23RunInput(t, a, "evidence", "link", a.actor, "p23-human-outcome-link-"+suffix, map[string]any{"obligation_id": obligationID, "evidence_type": "human-outcome", "receipt_id": nil, "review_id": nil, "human_outcome_id": outcomeID, "target_repository": a.repository, "required_artifact_ids": []string{}})
}

func assertP23AuthorityMissing(t *testing.T, response p22ConformanceResponse, before int) {
	t.Helper()
	if response.ExitCode != 5 || response.Stdout != "" {
		t.Fatalf("authority-missing denial = %#v", response)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatalf("decode authority-missing denial: %v; stderr=%q", err, response.Stderr)
	}
	if got := p21PathString(t, envelope, "error", "code"); got != "authority-missing" {
		t.Fatalf("denial code = %q, want authority-missing", got)
	}
}

func assertP23HumanOutcomeSuccess(t *testing.T, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("Human Outcome mutation = %#v", response)
	}
}

func assertP23HumanOutcomeCommandSuccess(t *testing.T, operation string, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("%s = %#v", operation, response)
	}
}

func assertP23OutcomeNotPermitted(t *testing.T, response p22ConformanceResponse, obligationType string) {
	t.Helper()
	if response.ExitCode != 4 || response.Stdout != "" {
		t.Fatalf("outcome-not-permitted denial = %#v", response)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatalf("decode outcome-not-permitted denial: %v; stderr=%q", err, response.Stderr)
	}
	if got := p21PathString(t, envelope, "error", "code"); got != "outcome-not-permitted" {
		t.Fatalf("denial code = %q, want outcome-not-permitted", got)
	}
	if got := p21PathString(t, envelope, "error", "details", "obligation_type"); got != obligationType {
		t.Fatalf("denial obligation_type = %q, want %q", got, obligationType)
	}
}

func assertP23HumanOutcomeRequiredDenial(t *testing.T, response p22ConformanceResponse, obligationID, outcomeClass string) {
	t.Helper()
	if response.ExitCode != 4 || response.Stdout != "" {
		t.Fatalf("human-outcome-required denial = %#v", response)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatalf("decode human-outcome-required denial: %v; stderr=%q", err, response.Stderr)
	}
	if got := p21PathString(t, envelope, "error", "code"); got != "human-outcome-required" {
		t.Fatalf("denial code = %q, want human-outcome-required", got)
	}
	if got := p21PathString(t, envelope, "error", "details", "obligation_id"); got != obligationID {
		t.Fatalf("denial obligation_id = %q, want %q", got, obligationID)
	}
	if got := p21PathString(t, envelope, "error", "details", "outcome_class"); got != outcomeClass {
		t.Fatalf("denial outcome_class = %q, want %q", got, outcomeClass)
	}
}

func assertP23HumanOutcomeStopState(t *testing.T, a *p22FIDRecoveryAdapter, outcomeClass, want string) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var status string
	err := conn.QueryRow(`SELECT CASE WHEN r.stop_id IS NULL THEN 'open' ELSE 'resolved' END FROM stops s LEFT JOIN stop_resolutions r ON r.stop_id=s.id WHERE s.code='human-outcome-required' AND json_extract(s.details_json,'$.outcome_class')=?`, outcomeClass).Scan(&status)
	p22FaultOK(t, err)
	if status != want {
		t.Fatalf("human-outcome-required Stop for class %q = %q, want %q", outcomeClass, status, want)
	}
}

func p23RecordHumanOutcomeMissingCapability(t *testing.T) p22ConformanceResponse {
	t.Helper()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p23-human-outcome@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 23 Human Outcome")
	p22Write(t, workspace, "base.txt", "base\n")
	p22Git(t, workspace, "add", "base.txt")
	p22Git(t, workspace, "commit", "-qm", "base")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P23HO")
	activationID := p21Activation(t, binary, workspace, "prepare", "p23-human-outcome-activation-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p23-human-outcome", "phase": 21})
	profiles := []any{
		map[string]any{"actor_id": "actor-no-capability", "actor_type": "human", "producer": "activated-harness", "provenance_ref": "activation:p23-human-outcome", "session_ref": nil, "capabilities": []string{"evidence.link"}},
	}
	profileDigest := p21Component(t, binary, workspace, activationID, "authority-profile", "p23-human-outcome-authority", profiles)
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "p23-human-outcome-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-p23-human-outcome", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	digests := []string{gateDigest, profileDigest}
	if digests[0] > digests[1] {
		digests[0], digests[1] = digests[1], digests[0]
	}
	p21Activation(t, binary, workspace, "commit", "p23-human-outcome-activation-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	p21Activation(t, binary, workspace, "finalize", "p23-human-outcome-activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})
	input := map[string]any{
		"id": "OUT-missing-capability", "class": "release", "subject": map[string]any{"type": "issue", "id": "ISS-p23-human-outcome"},
		"result": "approved", "scope": map[string]any{"issue_id": "ISS-p23-human-outcome"}, "reasoning": "real-process missing-capability negative", "references": []any{},
	}
	stdout, stderr, err := closureRunCommand(workspace, binary, "--project", workspace, "--json", "human-outcome", "record", "--actor", "actor-no-capability", "--input", p21Input(t, workspace, "p23-human-outcome-missing-capability", input), "--idempotency-key", "p23-human-outcome-missing-capability")
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return response
}

// p23AssertHumanOutcomeConditionNullableBindingFires complements
// internal/journal's direct assertion
// (TestTask10HumanOutcomeConditionDispatchPreservesNullableBinding, which
// proves a stray budget_name is rejected) by observing that the real,
// built-binary-recorded Stop carries the canonical null-bound Condition form
// - i.e. the accepted producer-side shape the journal-layer probe guards.
func p23AssertHumanOutcomeConditionNullableBindingFires(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	class := "security"
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var count int
	err := conn.QueryRow(`SELECT COUNT(*) FROM stops WHERE code='human-outcome-required' AND json_extract(resolution_condition_json,'$.type')='human-outcome-recorded' AND json_extract(resolution_condition_json,'$.budget_name') IS NULL AND json_extract(resolution_condition_json,'$.capability_id') IS NULL AND json_extract(resolution_condition_json,'$.outcome_class')=?`, class).Scan(&count)
	p22FaultOK(t, err)
	if count == 0 {
		t.Fatal("no recorded human-outcome-required Stop carries the canonical null-bound Condition")
	}
}
