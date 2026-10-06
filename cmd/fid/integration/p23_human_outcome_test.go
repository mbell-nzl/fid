//go:build integration

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestP24HumanOutcomeRejectsReferenceDigestMismatch(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "review.record", "run.close", "run.observe", "run.observe-budget", "run.start"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, actions,
		p22ExecutionFixture{HumanObligations: true, CoexistingReviewObligation: true})
	path := filepath.Join(a.workspace, "judgment-source.md")
	if err := os.WriteFile(path, []byte("actual source bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	input := map[string]any{
		"id": "OUT-reference-mismatch", "class": "release", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"result": "approved", "scope": map[string]any{"issue_id": a.issueID}, "reasoning": "reference check",
		"references": []any{map[string]any{"source_type": "project-file", "display_label": "judgment source",
			"path": "judgment-source.md", "section": nil, "external_id": nil, "digest": "sha256:" + p22Zeros(64)}},
	}
	before := p21JournalEventCount(t, a.workspace)
	response := p23RunInput(t, a, "human-outcome", "record", a.actor, "p24-reference-mismatch", input)
	assertP22ProcessFailure(t, "human-outcome.record", "source-digest-mismatch", 4, response)
	if after := p21JournalEventCount(t, a.workspace); after != before {
		t.Fatalf("digest mismatch appended %d Events", after-before)
	}
}

// TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess proves
// provenanced Human Outcome recording end to end through the built binary:
// VAL-043 authority admission and its four negatives, VAL-045's deterministic
// non-proof fields, VAL-046's readiness-Stop raising and resolution, VAL-025's
// final nine-cell evidence-type closure (both Human mismatches through the
// real producer), VAL-009's blocking stages (including the DEC-249 stage where
// the gate is satisfied and the review is not) and completing state with one
// Envelope carrying gate-evidence, review, and human-outcome Obligations at
// once, VAL-010's
// linked-evidence re-verification, and VAL-081's re-confirmed nullable binding.
func TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "review.record", "run.close", "run.observe", "run.observe-budget", "run.start"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, actions, p22ExecutionFixture{HumanObligations: true, CoexistingReviewObligation: true})
	a.task13ReplayCorpus = true
	missingOutcome := p23LinkHumanOutcome(t, a, "missing-outcome", "OBL-human-release", "OUT-missing")
	assertP22ProcessFailure(t, "evidence.link", "outcome-not-found", 3, missingOutcome)
	assertP23ThreeObligationTypesCoexist(t, a)
	resume := p21Run(t, a.binary, a.workspace, "--json", "resume", a.issueID, "--repository", a.workspace)
	evaluations := p21Path(t, resume, "result", "resume", "obligations").([]any)
	if len(evaluations) != 4 {
		t.Fatalf("resume Obligation evaluations = %#v, want all four", evaluations)
	}
	for index, want := range []string{"OBL-human-release", "OBL-human-security", "OBL-p22-corpus", "OBL-p23-review"} {
		if evaluations[index].(map[string]any)["obligation_id"] != want {
			t.Fatalf("resume Obligation order = %#v, want %s at %d", evaluations, want, index)
		}
	}
	// issue.transition-terminal-success (qa_pass) evaluates only the Obligations
	// whose guarded_actions contain it (DEC-249). In this fixture that is
	// OBL-p22-corpus alone; the review and both Human Obligations guard only
	// envelope.complete. A permitted qa_pass commits a one-way transition, so it
	// stays deferred until the journey's real terminal transition below and is
	// never used as a disposable probe once the gate is satisfied (stages 2 and 3
	// use the non-committing domain evaluator instead). The non-terminal states
	// below do not evaluate Obligations at all.
	for revision, state := range []string{"in_progress", "dev_complete", "qa_review"} {
		response := p23IssueTransition(t, a, state, revision+1, "human-outcome-"+state)
		assertP23HumanOutcomeCommandSuccess(t, "transition Issue to "+state, response)
	}

	// VAL-043 request 1 of 5: agent actor type (actor-denied shares actor's
	// exact capabilities but is staged actor_type "agent"). Each of requests 1-4
	// exits 5 with authority-missing and leaves the journal length unchanged.
	before := p21JournalEventCount(t, a.workspace)
	agentAttempt := p23RecordHumanOutcome(t, a, "agent-actor", "actor-denied", "release", "approved")
	assertP23AuthorityMissing(t, agentAttempt, before, p21JournalEventCount(t, a.workspace))

	// VAL-043 request 2 of 5: a Kit Activation that was prepared but never
	// finalized, whose human Actor DOES hold human-outcome.record. Two states
	// are driven: prepared with its components staged, and committed but still
	// not finalized. Neither may produce an authority-bearing Actor.
	for _, unfinalized := range p23UnfinalizedActivationAttempts(t) {
		assertP23AuthorityMissing(t, unfinalized.Response, unfinalized.Before, unfinalized.After)
	}

	// VAL-043 request 3 of 5: missing capability, proven with a dedicated
	// minimal finalized Activation whose only human Actor omits
	// human-outcome.record.
	missing := p23RecordHumanOutcomeMissingCapability(t)
	assertP23AuthorityMissing(t, missing.Response, missing.Before, missing.After)

	// VAL-043 request 4 of 5: standalone Issue-inventory attribution. The
	// attributed human:unregistered Actor resolves in no finalized profile.
	before = p21JournalEventCount(t, a.workspace)
	unregistered := p23RecordHumanOutcome(t, a, "unregistered-actor", "human:unregistered", "release", "approved")
	assertP23AuthorityMissing(t, unregistered, before, p21JournalEventCount(t, a.workspace))

	// VAL-043 request 5 of 5 (valid Kit-mode case) + VAL-045: the recorded
	// Outcome and its Evidence never present deterministic proof.
	before = p21JournalEventCount(t, a.workspace)
	release := p23RecordHumanOutcome(t, a, "release-approved", a.actor, "release", "approved")
	assertP23HumanOutcomeSuccess(t, release)
	p23AssertActivity(t, a.workspace, release, "human-outcome-recorded", a.issueID, 1)
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
	// human-outcome Obligation and re-prove both Human mismatches with the real
	// producer, each bracketed by an unchanged journal length and projection.
	assertP23OutcomeLinkNotPermitted(t, a, "release-vs-gate", "OBL-p22-corpus", releaseOutcomeID, "gate-evidence")
	assertP23OutcomeLinkNotPermitted(t, a, "release-vs-review", "OBL-p23-review", releaseOutcomeID, "review")
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

	// VAL-009 stage 1: the release Human Outcome is satisfied, but the
	// gate-evidence, review, and security human-outcome Obligations are not.
	// envelope.complete is denied by all three; the terminal transition is
	// denied by OBL-p22-corpus alone, the only unsatisfied Obligation that
	// guards it (DEC-249). Neither denial appends an Event or changes the
	// projection.
	assertP23GuardedActionsDenied(t, a, "human-outcome-stage-1", []string{"OBL-human-security", "OBL-p22-corpus", "OBL-p23-review"}, []string{"OBL-p22-corpus"})

	// VAL-009 stage 2 (review-7 MF-5, DEC-249): the gate-evidence Obligation is
	// satisfied by a real fid-executed receipt, and the review Obligation is not
	// (the security human-outcome Obligation is also still unsatisfied, which is
	// what keeps the readiness Stop unraised, see below). OBL-p23-review and
	// OBL-human-security guard envelope.complete alone, so it is denied naming
	// exactly those two with no Event and no projection change, while
	// issue.transition-terminal-success (guarded only by the now-satisfied
	// OBL-p22-corpus) is not blocked by any Obligation. The permitted qa_pass is
	// the journey's one real terminal transition further below; it is not used as
	// a disposable probe here, so the non-committing domain evaluator asserts the
	// terminal outcome. Satisfaction is never revoked, so this gate-satisfied,
	// review-unsatisfied state is only reachable before the review is linked.
	gate := p23RunQualifyingGateSession(t, a, "human-outcome-proof")
	beforeGateLink := p21JournalEventCount(t, a.workspace)
	linkedGate := p23RunInput(t, a, "evidence", "link", a.actor, "human-outcome-gate-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
	assertP23HumanOutcomeSuccess(t, linkedGate)
	if after := p21JournalEventCount(t, a.workspace); after != beforeGateLink+1 {
		t.Fatalf("gate evidence.link appended %d Event(s), want 1", after-beforeGateLink)
	}
	assertP23ObligationsUnsatisfied(t, a, "OBL-human-security", "OBL-p23-review")
	assertP23NoHumanOutcomeStop(t, a)
	assertP23EnvelopeCompleteDenied(t, a, "human-outcome-stage-2", []string{"OBL-human-security", "OBL-p23-review"})
	assertP23TerminalTransitionNotBlocked(t, a, "human-outcome-stage-2")

	// VAL-046: an independent Reviewer records a passing Review that is linked
	// to the review Obligation. That evidence.link (the ordinary DEC-224
	// producer) leaves the security human-outcome Obligation as the last
	// unsatisfied one, so it raises human-outcome-required atomically inside its
	// own single Event. VAL-010: the linked review Evidence's evaluation reports
	// the distinct not-applicable artifact availability, never available/missing.
	reviewEvaluation := p23SatisfyCoexistingReview(t, a, "human-outcome-review")
	if reviewEvaluation["status"] != "satisfied" || reviewEvaluation["artifact_availability"] != "not-applicable" {
		t.Fatalf("linked review evaluation = %#v, want satisfied/not-applicable", reviewEvaluation)
	}
	assertP23HumanOutcomeStopState(t, a, "security", "open")

	// VAL-009: envelope.complete stays blocked by human-outcome-required
	// (rank 9) ahead of the generic issue-not-terminal-success reason (rank
	// 19) while the last Obligation (OBL-human-security) is unsatisfied. The
	// builder's re-acquired mutating Claim and Run stay active here; the
	// terminal Issue transition below still needs them.
	beforeDenial := p21JournalEventCount(t, a.workspace)
	denied := p23RunInput(t, a, "envelope", "complete", a.actor, "human-outcome-complete-denied", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1})
	assertP23HumanOutcomeRequiredDenial(t, a, denied, "OBL-human-security", "security")
	if after := p21JournalEventCount(t, a.workspace); after != beforeDenial {
		t.Fatalf("human-outcome-required denial appended Event(s): %d -> %d", beforeDenial, after)
	}
	// Stage 3: only the security human-outcome Obligation is unsatisfied. It
	// guards envelope.complete alone, which the human-outcome-required denial
	// above already proves it blocks. Under DEC-249 it does not guard the
	// terminal transition, so qa_pass is not blocked by any Obligation here.
	// The one-way qa_pass is still reserved for the journey's real terminal step,
	// so the non-committing domain evaluator is used.
	assertP23ObligationsUnsatisfied(t, a, "OBL-human-security")
	assertP23TerminalTransitionNotBlocked(t, a, "human-outcome-stage-3")

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

	// VAL-010: with linked review Evidence (stage 2) and linked Human Outcome
	// Evidence both present, artifact_availability stays the distinct
	// not-applicable value for the two new Obligation types, never
	// available/missing, and gate-evidence keeps reporting available.
	securityEvaluation := p21Path(t, p23Envelope(t, securityMatched.Stdout), "result", "evaluation").(map[string]any)
	if securityEvaluation["artifact_availability"] != "not-applicable" {
		t.Fatalf("security evaluation artifact_availability = %v, want not-applicable", securityEvaluation["artifact_availability"])
	}
	allSatisfied := p23EvidenceEvaluate(t, a)
	for obligationID, availability := range map[string]string{"OBL-p22-corpus": "available", "OBL-p23-review": "not-applicable", "OBL-human-release": "not-applicable", "OBL-human-security": "not-applicable"} {
		evaluation := p23EvaluationByID(t, allSatisfied, obligationID)
		if evaluation["status"] != "satisfied" || evaluation["artifact_availability"] != availability {
			t.Fatalf("all-satisfied evaluation for %s = %#v, want satisfied/%s", obligationID, evaluation, availability)
		}
	}

	// VAL-009 completing state: gate-evidence, review, and both human-outcome
	// Obligations are all satisfied, so the terminal qa_pass transition
	// (issue.transition-terminal-success) and envelope.complete now succeed.
	assertP23HumanOutcomeCommandSuccess(t, "transition Issue to qa_pass", p23IssueTransition(t, a, "qa_pass", 4, "human-outcome-qa-pass"))
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "human-outcome-close")
	assertP23HumanOutcomeCommandSuccess(t, "complete Envelope after every Obligation satisfied", p23RunInput(t, a, "envelope", "complete", a.actor, "human-outcome-complete", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1}))

	// VAL-025's remaining seven cells are supplied by Task 6 (receipt/review
	// positives and the four residual evidence-type-not-accepted mismatches),
	// reconciled here through the real producer's two Human Outcome cells. The
	// security-class Outcome repeats both mismatches after completion.
	for _, mismatch := range []struct {
		name           string
		obligationID   string
		obligationType string
	}{
		{"security-vs-gate", "OBL-p22-corpus", "gate-evidence"},
		{"security-vs-review", "OBL-p23-review", "review"},
	} {
		assertP23OutcomeLinkNotPermitted(t, a, mismatch.name, mismatch.obligationID, securityOutcomeID, mismatch.obligationType)
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
// three), and artifact.availability-observed has its own dedicated fixture in
// TestP23HumanOutcomeRequiredRaisedByArtifactAvailability.
func TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p23-human-outcome-activation@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 23 Human Outcome Activation")
	p22Write(t, workspace, "base.txt", "base\n")
	if err := os.MkdirAll(filepath.Join(workspace, "project-docs"), 0o700); err != nil {
		t.Fatal(err)
	}
	p22Write(t, workspace, "project-docs/p23-human-activation.md", "Human readiness policy\n")
	p22Git(t, workspace, "add", "base.txt", "project-docs/p23-human-activation.md")
	p22Git(t, workspace, "commit", "-qm", "base")
	init := p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P23HA")
	_ = p21PathString(t, init, "result", "project", "id")
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

	source := map[string]any{"digest": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("Human readiness policy\n"))), "display_label": "Human-only Envelope", "external_id": nil, "path": "project-docs/p23-human-activation.md", "section": nil, "source_type": "project-file"}
	obligation := map[string]any{"id": "OBL-human-only", "type": "human-outcome", "applicability": "required", "guarded_actions": []string{"envelope.complete"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "outcome_class": "release", "accepted_evidence_types": []string{"human-outcome"}, "accepted_trust_sources": []string{"activated-human"}, "freshness_target": "none", "artifact_policy": "metadata-sufficient", "independence": "human", "required_lane": nil, "source": source}
	envelopeRequest := map[string]any{"outcome": "Human-only readiness", "non_goals": []string{}, "risk_result": map[string]any{"source": source, "result_digest": "sha256:" + p22Ones(64)}, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{}, "authority_grants": []any{map[string]any{"actions": capabilities, "actor_selector": map[string]any{"actor_id": actor, "required_capabilities": capabilities}, "scope": map[string]any{"issue_id": issueID}}}, "obligations": []any{obligation}, "budgets": []any{}}
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
	t.Parallel()
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
	return p23RunInput(t, a, "human-outcome", "record", actor, "p23-human-outcome-"+suffix, input)
}

func p23LinkHumanOutcome(t *testing.T, a *p22FIDRecoveryAdapter, suffix, obligationID, outcomeID string) p22ConformanceResponse {
	t.Helper()
	return p23RunInput(t, a, "evidence", "link", a.actor, "p23-human-outcome-link-"+suffix, map[string]any{"obligation_id": obligationID, "evidence_type": "human-outcome", "receipt_id": nil, "review_id": nil, "human_outcome_id": outcomeID, "target_repository": a.repository, "required_artifact_ids": []string{}})
}

// assertP23AuthorityMissing asserts a VAL-043 negative: exit 5 with
// authority-missing and a journal Event count (before -> after) that did not
// move, so no Human Outcome Event was appended.
func assertP23AuthorityMissing(t *testing.T, response p22ConformanceResponse, before, after int) {
	t.Helper()
	if after != before {
		t.Fatalf("authority-missing denial appended Event(s): %d -> %d; response=%#v", before, after, response)
	}
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

// assertP23OutcomeLinkNotPermitted links a recorded Human Outcome to an
// Obligation of the wrong type through the real producer and asserts the
// outcome-not-permitted denial with its exact obligation_type detail, bracketed
// by an unchanged journal length and an unchanged whole-database projection.
func assertP23OutcomeLinkNotPermitted(t *testing.T, a *p22FIDRecoveryAdapter, suffix, obligationID, outcomeID, obligationType string) {
	t.Helper()
	beforeEvents, beforeProjection := p21JournalEventCount(t, a.workspace), p22FaultProjection(t, a.workspace)
	response := p23LinkHumanOutcome(t, a, suffix, obligationID, outcomeID)
	assertP22ProcessFailure(t, "evidence.link", "outcome-not-permitted", 4, response)
	assertP23OutcomeNotPermitted(t, response, obligationType)
	if afterEvents := p21JournalEventCount(t, a.workspace); afterEvents != beforeEvents {
		t.Fatalf("%s outcome-not-permitted denial appended Event(s): %d -> %d", suffix, beforeEvents, afterEvents)
	}
	if !reflect.DeepEqual(beforeProjection, p22FaultProjection(t, a.workspace)) {
		t.Fatalf("%s outcome-not-permitted denial changed the projection", suffix)
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

func assertP23HumanOutcomeRequiredDenial(t *testing.T, a *p22FIDRecoveryAdapter, response p22ConformanceResponse, obligationID, outcomeClass string) {
	t.Helper()
	assertP22ProcessFailure(t, "envelope.complete", "human-outcome-required", 4, response)
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
	p23AssertStopBlockedDetails(t, a, response, "human-outcome-required", "outcome_class", outcomeClass)
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

// p23IsolatedHumanOutcomeAttempt is one Human Outcome request executed in a
// dedicated workspace, with the journal Event count observed immediately
// before and after the request.
type p23IsolatedHumanOutcomeAttempt struct {
	Name     string
	Response p22ConformanceResponse
	Before   int
	After    int
}

// p23IsolatedActivationStage says how far the dedicated Activation progresses.
type p23IsolatedActivationStage string

const (
	// p23ActivationPrepared stops after the prepare and its two components.
	p23ActivationPrepared p23IsolatedActivationStage = "prepared"
	// p23ActivationCommitted stops after commit, before finalize.
	p23ActivationCommitted p23IsolatedActivationStage = "committed"
	// p23ActivationFinalized runs the whole prepare, commit, finalize sequence.
	p23ActivationFinalized p23IsolatedActivationStage = "finalized"
)

// p23UnfinalizedActivationAttempts drives VAL-043's unfinalized-Activation
// request twice: a prepared Activation whose components are staged, and a
// committed Activation that is still not finalized. In both the sole human Actor
// declares human-outcome.record and a real Issue exists, so only the missing
// finalization can deny the request.
func p23UnfinalizedActivationAttempts(t *testing.T) []p23IsolatedHumanOutcomeAttempt {
	t.Helper()
	return []p23IsolatedHumanOutcomeAttempt{
		p23IsolatedHumanOutcomeRequest(t, "unfinalized-prepared", p23ActivationPrepared, []string{"evidence.link", "human-outcome.record"}),
		p23IsolatedHumanOutcomeRequest(t, "unfinalized-committed", p23ActivationCommitted, []string{"evidence.link", "human-outcome.record"}),
	}
}

func p23RecordHumanOutcomeMissingCapability(t *testing.T) p23IsolatedHumanOutcomeAttempt {
	t.Helper()
	return p23IsolatedHumanOutcomeRequest(t, "missing-capability", p23ActivationFinalized, []string{"evidence.link"})
}

func p23IsolatedHumanOutcomeRequest(t *testing.T, name string, stage p23IsolatedActivationStage, capabilities []string) p23IsolatedHumanOutcomeAttempt {
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
	issue := p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Human Outcome "+name, "--actor", "agent:developer", "--idempotency-key", "p23-human-outcome-"+name+"-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")
	actorID := "actor-" + name
	prefix := "p23-human-outcome-" + name
	activationID := p21Activation(t, binary, workspace, "prepare", prefix+"-activation-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": prefix, "phase": 21})
	profiles := []any{
		map[string]any{"actor_id": actorID, "actor_type": "human", "producer": "activated-harness", "provenance_ref": "activation:" + prefix, "session_ref": nil, "capabilities": capabilities},
	}
	profileDigest := p21Component(t, binary, workspace, activationID, "authority-profile", prefix+"-authority", profiles)
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", prefix+"-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-" + prefix, "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	if stage != p23ActivationPrepared {
		digests := []string{gateDigest, profileDigest}
		sort.Strings(digests)
		p21Activation(t, binary, workspace, "commit", prefix+"-activation-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	}
	if stage == p23ActivationFinalized {
		p21Activation(t, binary, workspace, "finalize", prefix+"-activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})
	}
	input := map[string]any{
		"id": "OUT-" + name, "class": "release", "subject": map[string]any{"type": "issue", "id": issueID},
		"result": "approved", "scope": map[string]any{"issue_id": issueID}, "reasoning": "real-process " + name + " negative", "references": []any{},
	}
	inputPath := p21Input(t, workspace, prefix+"-record", input)
	before := p21JournalEventCount(t, workspace)
	stdout, stderr, err := closureRunCommand(workspace, binary, "--project", workspace, "--json", "human-outcome", "record", "--actor", actorID, "--input", inputPath, "--idempotency-key", prefix+"-record")
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return p23IsolatedHumanOutcomeAttempt{Name: name, Response: response, Before: before, After: p21JournalEventCount(t, workspace)}
}

// assertP23ThreeObligationTypesCoexist reads the durable Envelope projection and
// asserts one Envelope carries a gate-evidence, a review, and human-outcome
// Obligations at the same time, each under its own ID.
func assertP23ThreeObligationTypesCoexist(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`SELECT id, obligation_type FROM envelope_obligations WHERE envelope_id=? ORDER BY id`, p23EnvelopeID(t, a))
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	got := map[string]string{}
	for rows.Next() {
		var id, obligationType string
		p22FaultOK(t, rows.Scan(&id, &obligationType))
		got[id] = obligationType
	}
	p22FaultOK(t, rows.Err())
	want := map[string]string{"OBL-human-release": "human-outcome", "OBL-human-security": "human-outcome", "OBL-p22-corpus": "gate-evidence", "OBL-p23-review": "review"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Envelope Obligations = %#v, want %#v", got, want)
	}
}

// assertP23GuardedActionsDenied proves both VAL-009 guarded actions are denied
// while an Obligation that guards them is unsatisfied (DEC-249): envelope.complete
// rejects naming exactly the completeUnsatisfied Obligations, and the terminal
// qa_pass transition (issue.transition-terminal-success) rejects for exactly the
// terminalUnsatisfied Obligations, which are the unsatisfied ones that guard it.
// Neither denial appends an Event or changes the projection.
func assertP23GuardedActionsDenied(t *testing.T, a *p22FIDRecoveryAdapter, key string, completeUnsatisfied, terminalUnsatisfied []string) {
	t.Helper()
	beforeEvents, beforeProjection := p21JournalEventCount(t, a.workspace), p22FaultProjection(t, a.workspace)
	complete := p23RunInput(t, a, "envelope", "complete", a.actor, key+"-complete", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1})
	if complete.ExitCode != 4 || complete.Stdout != "" {
		t.Fatalf("%s envelope.complete denial = %#v, want exit 4 with no stdout", key, complete)
	}
	var denial map[string]any
	if err := json.Unmarshal([]byte(complete.Stderr), &denial); err != nil {
		t.Fatalf("decode %s envelope.complete denial: %v; stderr=%q", key, err, complete.Stderr)
	}
	if got := p21PathString(t, denial, "error", "code"); got != "obligation-unsatisfied" {
		t.Fatalf("%s envelope.complete denial code = %q, want obligation-unsatisfied", key, got)
	}
	var blocking []string
	for _, reason := range p21Path(t, denial, "error", "details", "reasons").([]any) {
		if item := reason.(map[string]any); item["subject_type"] == "obligation" {
			blocking = append(blocking, item["subject_id"].(string))
		}
	}
	sort.Strings(blocking)
	if !reflect.DeepEqual(blocking, completeUnsatisfied) {
		t.Fatalf("%s envelope.complete blocked by Obligations %v, want exactly %v", key, blocking, completeUnsatisfied)
	}
	assertP23TerminalTransitionDenied(t, a, key, terminalUnsatisfied)
	if afterEvents := p21JournalEventCount(t, a.workspace); afterEvents != beforeEvents {
		t.Fatalf("%s guarded-action denials appended Event(s): %d -> %d", key, beforeEvents, afterEvents)
	}
	if !reflect.DeepEqual(beforeProjection, p22FaultProjection(t, a.workspace)) {
		t.Fatalf("%s guarded-action denials changed the projection", key)
	}
}

// assertP23ObligationsUnsatisfied asserts the evidence.evaluate read reports
// every Obligation and that exactly the named ones are unsatisfied.
func assertP23ObligationsUnsatisfied(t *testing.T, a *p22FIDRecoveryAdapter, unsatisfiedIDs ...string) {
	t.Helper()
	response := p23EvidenceEvaluate(t, a)
	unsatisfied := map[string]bool{}
	for _, id := range unsatisfiedIDs {
		unsatisfied[id] = true
	}
	seen := 0
	for _, item := range p21Path(t, p23Envelope(t, response.Stdout), "result", "evaluations").([]any) {
		evaluation := item.(map[string]any)
		seen++
		want := "satisfied"
		if unsatisfied[evaluation["obligation_id"].(string)] {
			want = "unsatisfied"
		}
		if evaluation["status"] != want {
			t.Fatalf("evaluation for %v = %v, want %s (only %v unsatisfied)", evaluation["obligation_id"], evaluation["status"], want, unsatisfiedIDs)
		}
	}
	if seen != 4 {
		t.Fatalf("evidence.evaluate reported %d Obligations, want every one of the 4", seen)
	}
}

// assertP23EnvelopeCompleteDenied attempts envelope.complete and asserts it is
// denied because of exactly the named unsatisfied Obligations, each of which
// guards envelope.complete (DEC-249), bracketed by an unchanged journal length
// and an unchanged whole-database projection.
func assertP23EnvelopeCompleteDenied(t *testing.T, a *p22FIDRecoveryAdapter, key string, want []string) {
	t.Helper()
	beforeEvents, beforeProjection := p21JournalEventCount(t, a.workspace), p22FaultProjection(t, a.workspace)
	if got := p23CompleteDenialObligationIDs(t, a, key+"-complete"); !reflect.DeepEqual(got, want) {
		t.Fatalf("%s envelope.complete blocked by Obligations %v, want exactly %v", key, got, want)
	}
	if afterEvents := p21JournalEventCount(t, a.workspace); afterEvents != beforeEvents {
		t.Fatalf("%s envelope.complete denial appended Event(s): %d -> %d", key, beforeEvents, afterEvents)
	}
	if !reflect.DeepEqual(beforeProjection, p22FaultProjection(t, a.workspace)) {
		t.Fatalf("%s envelope.complete denial changed the projection", key)
	}
}

// assertP23NoHumanOutcomeStop asserts no human-outcome-required Stop exists:
// the readiness Stop is raised only when a Human Obligation is the sole
// unsatisfied Obligation, never while other Obligations are also unsatisfied.
func assertP23NoHumanOutcomeStop(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var count int
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM stops WHERE code='human-outcome-required'`).Scan(&count))
	if count != 0 {
		t.Fatalf("human-outcome-required Stop count = %d, want 0 while other Obligations remain unsatisfied", count)
	}
}

// assertP23TerminalTransitionNotBlocked evaluates the terminal qa_pass
// transition with the domain evaluator, which reads durable state and commits
// nothing, and asserts no reason at all blocks it. Journal and projection are
// unchanged. It never performs the one-way transition.
func assertP23TerminalTransitionNotBlocked(t *testing.T, a *p22FIDRecoveryAdapter, key string) {
	t.Helper()
	beforeEvents, beforeProjection := p21JournalEventCount(t, a.workspace), p22FaultProjection(t, a.workspace)
	revision := 4
	evaluation, err := workflow.EvaluateExecutionTransition(a.workspace, filepath.Join(integrationRepositoryRoot(), "migrations"), workflow.ExecutionTransitionRequest{IssueID: a.issueID, ActorID: a.actor, ExpectedIssueRevision: &revision, ActualIssueRevision: revision, WorkflowEdgeValid: true, TerminalSuccess: true})
	p22FaultOK(t, err)
	if !evaluation.Bound || len(evaluation.Reasons) != 0 {
		t.Fatalf("%s terminal transition evaluation = %#v, want bound with no blocking reason", key, evaluation)
	}
	if afterEvents := p21JournalEventCount(t, a.workspace); afterEvents != beforeEvents || !reflect.DeepEqual(beforeProjection, p22FaultProjection(t, a.workspace)) {
		t.Fatalf("%s terminal transition evaluation changed the journal (%d -> %d) or the projection", key, beforeEvents, afterEvents)
	}
}

// assertP23TerminalTransitionDenied attempts the terminal qa_pass transition
// (issue.transition-terminal-success) and asserts it is denied because of
// exactly the named unsatisfied Obligations, with no Event and no projection
// change. A human-outcome Obligation surfaces as human-outcome-absent and every
// other Obligation as obligation-evidence-absent.
func assertP23TerminalTransitionDenied(t *testing.T, a *p22FIDRecoveryAdapter, key string, unsatisfied []string) {
	t.Helper()
	beforeEvents, beforeProjection := p21JournalEventCount(t, a.workspace), p22FaultProjection(t, a.workspace)
	terminal := p23IssueTransition(t, a, "qa_pass", 4, key+"-qa-pass")
	if terminal.ExitCode == 0 || terminal.Stdout != "" {
		t.Fatalf("%s terminal qa_pass transition = %#v, want a denial", key, terminal)
	}
	var denial map[string]any
	if err := json.Unmarshal([]byte(terminal.Stderr), &denial); err != nil {
		t.Fatalf("decode %s qa_pass denial: %v; stderr=%q", key, err, terminal.Stderr)
	}
	if got := p21PathString(t, denial, "error", "code"); got != "transition-not-permitted" {
		t.Fatalf("%s qa_pass denial code = %q, want transition-not-permitted", key, got)
	}
	var want, got []string
	for _, id := range unsatisfied {
		if len(id) >= 10 && id[:10] == "OBL-human-" {
			want = append(want, "human-outcome-absent")
		} else {
			want = append(want, "obligation-evidence-absent")
		}
	}
	for _, reason := range p21Path(t, denial, "error", "details", "unsatisfied").([]any) {
		if code := reason.(string); code == "human-outcome-absent" || code == "obligation-evidence-absent" {
			got = append(got, code)
		} else {
			t.Fatalf("%s qa_pass denial carries non-Obligation reason %q in %v", key, code, got)
		}
	}
	sort.Strings(want)
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s qa_pass denial reasons = %v, want exactly %v", key, got, want)
	}
	if afterEvents := p21JournalEventCount(t, a.workspace); afterEvents != beforeEvents {
		t.Fatalf("%s qa_pass denial appended Event(s): %d -> %d", key, beforeEvents, afterEvents)
	}
	if !reflect.DeepEqual(beforeProjection, p22FaultProjection(t, a.workspace)) {
		t.Fatalf("%s qa_pass denial changed the projection", key)
	}
}

// p23SatisfyCoexistingReview satisfies OBL-p23-review. A single Claim may be
// active per Issue, so the builder's mutating session is closed first; an
// independent Reviewer then records a passing Review that is linked as
// Evidence, and its review session is closed. A fresh mutating Claim and Run
// for the builder Actor are then re-acquired and adopted by the adapter so the
// terminal Issue transition stays reachable. It returns the linked review
// Evidence's evaluation.
func p23SatisfyCoexistingReview(t *testing.T, a *p22FIDRecoveryAdapter, key string) map[string]any {
	t.Helper()
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, key+"-builder")
	claim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", key)
	run := p23StartRun(t, a, "actor-reviewer-two", claim, key)
	pass := p23ReviewRecord(t, a, "actor-reviewer-two", key, p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "pass", []any{}, nil))
	p23AssertReviewSuccess(t, "coexisting pass Review", pass)
	review := p23ReviewFromResponse(t, pass)
	beforeLink := p21JournalEventCount(t, a.workspace)
	linked := p23RunInput(t, a, "evidence", "link", a.actor, key+"-link", map[string]any{"obligation_id": "OBL-p23-review", "evidence_type": "review", "receipt_id": nil, "review_id": review["id"], "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
	p23AssertReviewSuccess(t, "coexisting review Evidence link", linked)
	if after := p21JournalEventCount(t, a.workspace); after != beforeLink+1 {
		t.Fatalf("review evidence.link appended %d Event(s), want 1", after-beforeLink)
	}
	p23CloseAndRelease(t, a, "actor-reviewer-two", p21PathString(t, run, "id"), p21PathString(t, claim, "id"), key)
	p23ReacquireBuilderSession(t, a, key+"-builder-again")
	return p21Path(t, p23Envelope(t, linked.Stdout), "result", "evaluation").(map[string]any)
}

// p23RunQualifyingGateSession runs the registered gate through fid's own gate
// runner inside a dedicated mutating session held by actor-denied. The
// builder's session is closed first (one Claim may be active per Issue) and
// re-acquired afterwards. The session is held by actor-denied because, once a review-capable
// fixture adds further authority grants, internal/evidence
// validateGateRunReferences (execution.go:381) reads only the first
// envelope_authority_grants row and requires it to name the running Actor; the
// grants sort by actor_id, so that first row is actor-denied's, whose profile
// carries the same gate.run and run.observe-budget capabilities. The receipt is
// real and fid-executed; only its running Actor differs from the builder.
func p23RunQualifyingGateSession(t *testing.T, a *p22FIDRecoveryAdapter, key string) evidence.GateRunResult {
	t.Helper()
	const runner = "actor-denied"
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, key+"-builder")
	claim := p23AcquireClaim(t, a, runner, "mutation", key+"-claim")
	run := p23StartRun(t, a, runner, claim, key+"-run")
	runID := p21PathString(t, run, "id")
	sample := p23RunInput(t, a, "run", "observe-budget", runner, key+"-budget", map[string]any{"run_id": runID, "expected_version": run["version"], "budget_name": "tokens", "source_sample_id": key + "-sample", "sample_kind": "baseline", "previous_sample_id": nil, "value": 0, "measurement_ref": a.measurement, "observed_at": time.Now().UTC().Format(time.RFC3339Nano)})
	p23AssertReviewSuccess(t, "observe budget baseline for gate Run", sample)
	request, err := json.Marshal(map[string]any{"run_id": runID, "registration_id": a.registrationID, "registration_version": a.registrationVersion, "obligation_id": "OBL-p22-corpus", "repository": a.repository})
	p22FaultOK(t, err)
	options := evidence.GateRunOptions{WorkspaceRoot: a.workspace, MigrationsDir: filepath.Join(integrationRepositoryRoot(), "migrations"), ActorID: runner, IdempotencyKey: key}
	p22FaultOK(t, json.Unmarshal(request, &options.Request))
	gate, err := evidence.RunGate(context.Background(), options)
	if err != nil || gate.Receipt == nil {
		t.Fatalf("run qualifying gate: result=%#v error=%v", gate, err)
	}
	p23CloseAndRelease(t, a, runner, runID, p21PathString(t, claim, "id"), key)
	p23ReacquireBuilderSession(t, a, key+"-builder-again")
	return gate
}

// p23ReacquireBuilderSession gives the builder Actor a fresh mutating Claim and
// Run and adopts them into the adapter, for the terminal Issue transition.
func p23ReacquireBuilderSession(t *testing.T, a *p22FIDRecoveryAdapter, key string) {
	t.Helper()
	claim := p23AcquireClaim(t, a, a.actor, "mutation", key)
	run := p23StartRun(t, a, a.actor, claim, key)
	a.claimID, a.claimVersion = p21PathString(t, claim, "id"), p23Number(t, claim["version"])
	a.runID, a.runVersion = p21PathString(t, run, "id"), p23Number(t, run["version"])
	a.executorID = p21PathString(t, run, "executor", "id")
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

// TestP23HumanOutcomeRequiredRaisedByArtifactAvailability proves that the
// artifact.availability-observed producer raises once when its Artifact becomes
// available and the release Human Obligation is the only one still unsatisfied.
func TestP23HumanOutcomeRequiredRaisedByArtifactAvailability(t *testing.T) {
	t.Parallel()
	const artifactName = "artifact.json" // ignored by the fixture Git worktree
	actions := []string{"artifact.verify", "claim.acquire", "claim.release", "envelope.activate", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "review.record", "run.close", "run.observe", "run.observe-budget", "run.start"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{artifactName}, actions, p22ExecutionFixture{HumanObligations: true})
	content := []byte("readiness artifact\n")
	p22Write(t, a.workspace, artifactName, string(content))
	gate := p23RunQualifyingGateSession(t, a, "artifact-readiness-gate")
	if len(gate.Receipt.Artifacts) != 1 {
		t.Fatalf("gate Receipt Artifacts = %#v, want one", gate.Receipt.Artifacts)
	}
	artifact := gate.Receipt.Artifacts[0]
	path := filepath.Join(a.workspace, ".fid", "artifacts", artifact.Digest[len("sha256:"):])
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	verify := func(key, availability string) string {
		t.Helper()
		response := p23RunInput(t, a, "artifact", "verify", a.actor, key, map[string]any{"artifact_id": artifact.ID})
		assertP23HumanOutcomeSuccess(t, response)
		envelope := p23Envelope(t, response.Stdout)
		if got := p21PathString(t, envelope, "result", "artifact", "availability"); got != availability {
			t.Fatalf("artifact.verify availability = %q, want %q", got, availability)
		}
		if changed := p21Path(t, envelope, "result", "changed"); changed != true {
			t.Fatalf("artifact.verify changed = %v, want true", changed)
		}
		eventID := p21PathString(t, envelope, "event", "id")
		event, err := journal.FindEvent(filepath.Join(a.workspace, ".fid", "journal"), eventID)
		p22FaultOK(t, err)
		if event.Event.EventType != "artifact.availability-observed" {
			t.Fatalf("artifact.verify Event type = %q", event.Event.EventType)
		}
		return eventID
	}
	verify("artifact-readiness-missing", "missing")
	linked := p23RunInput(t, a, "evidence", "link", a.actor, "artifact-readiness-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{artifact.ID}})
	assertP23HumanOutcomeSuccess(t, linked)
	security := p23RecordHumanOutcome(t, a, "artifact-readiness-security", a.actor, "security", "approved")
	assertP23HumanOutcomeSuccess(t, security)
	securityID := p21PathString(t, p23Envelope(t, security.Stdout), "result", "outcome", "id")
	assertP23HumanOutcomeSuccess(t, p23LinkHumanOutcome(t, a, "artifact-readiness-security", "OBL-human-security", securityID))

	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	assertCounts := func(want int) {
		t.Helper()
		var stops, attention int
		envelopeID := p23EnvelopeID(t, a)
		p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM stops WHERE code='human-outcome-required' AND subject_id=?`, envelopeID).Scan(&stops))
		p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM attention_items WHERE code='human-outcome-required' AND subject_id=?`, envelopeID).Scan(&attention))
		if stops != want || attention != want {
			t.Fatalf("readiness Stop/Attention counts = %d/%d, want %d/%d", stops, attention, want, want)
		}
	}
	assertEvaluations := func(gateStatus string) {
		t.Helper()
		evaluations := p21Path(t, p23Envelope(t, p23EvidenceEvaluate(t, a).Stdout), "result", "evaluations").([]any)
		want := map[string]string{"OBL-p22-corpus": gateStatus, "OBL-human-security": "satisfied", "OBL-human-release": "unsatisfied"}
		if len(evaluations) != len(want) {
			t.Fatalf("artifact readiness evaluations = %d, want %d", len(evaluations), len(want))
		}
		for _, item := range evaluations {
			evaluation := item.(map[string]any)
			id := evaluation["obligation_id"].(string)
			if evaluation["status"] != want[id] {
				t.Fatalf("artifact readiness evaluation for %s = %v, want %s", id, evaluation["status"], want[id])
			}
		}
	}
	assertCounts(0)
	assertEvaluations("unsatisfied")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	firstAvailable := verify("artifact-readiness-available", "available")
	assertCounts(1)
	assertEvaluations("satisfied")
	var detailsJSON string
	p22FaultOK(t, conn.QueryRow(`SELECT details_json FROM stops WHERE code='human-outcome-required' AND subject_id=?`, p23EnvelopeID(t, a)).Scan(&detailsJSON))
	var details map[string]any
	p22FaultOK(t, json.Unmarshal([]byte(detailsJSON), &details))
	if details["obligation_id"] != "OBL-human-release" || details["outcome_class"] != "release" {
		t.Fatalf("artifact readiness Stop details = %#v", details)
	}
	assertP23HumanOutcomeStopState(t, a, "release", "open")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	verify("artifact-readiness-missing-again", "missing")
	assertCounts(1)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	secondAvailable := verify("artifact-readiness-available-again", "available")
	if secondAvailable == firstAvailable {
		t.Fatal("repeat qualifying availability flip reused the first Event")
	}
	assertCounts(1)
	assertP23HumanOutcomeStopState(t, a, "release", "open")
	release := p23RecordHumanOutcome(t, a, "artifact-readiness-release", a.actor, "release", "approved")
	assertP23HumanOutcomeSuccess(t, release)
	assertP23HumanOutcomeStopState(t, a, "release", "resolved")
}
