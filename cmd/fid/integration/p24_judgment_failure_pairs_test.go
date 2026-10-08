//go:build integration

package integration

import (
	"strings"
	"testing"
)

func p24FailureJudgmentAdapter(t *testing.T, cfg ...p22ExecutionFixture) *p22FIDRecoveryAdapter {
	t.Helper()
	actions := []string{"claim.acquire", "claim.release", "concern.dispose", "concern.record", "deviation.record", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "finding.address", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "review.record", "run.close", "run.observe", "run.observe-budget", "run.start"}
	return p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, actions, cfg...)
}

func p24JudgmentCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	a := p24FailureJudgmentAdapter(t)
	operation, input := "", map[string]any{}
	switch code {
	case "source-digest-mismatch":
		operation = "requirement.record-reference"
		input = map[string]any{"id": "REQ-pair", "source": p23DecisionSource(a, false), "status": "active"}
	case "finding-not-found":
		operation = "finding.address"
		input = map[string]any{"finding_id": "FND-missing", "claim_id": a.claimID, "expected_claim_version": a.claimVersion, "reason": "pair"}
	case "deviation-not-found":
		operation = "deviation.record"
		input = map[string]any{"id": "DEV-pair", "subject": map[string]any{"type": "issue", "id": a.issueID}, "expected": "expected", "observed": "observed", "authority": map[string]any{"actor_id": a.actor, "reference": nil}, "disposition": "pending-decision", "reasoning": "pair", "references": []any{}, "supersedes_deviation_id": "DEV-missing"}
	case "concern-not-found", "concern-already-disposed", "result-limit-exceeded":
		operation = "concern.dispose"
		input = map[string]any{"concern_id": "CON-pair", "action": "rejected", "issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "pair"}
		if code != "concern-not-found" {
			count := 1
			if code == "result-limit-exceeded" {
				count = 2
			}
			for i := 0; i < count; i++ {
				id := "CON-pair" + strings.Repeat("x", i)
				result := p23ConcernMutation(t, a, "record", id, map[string]any{"id": id, "subject": map[string]any{"type": "issue", "id": a.issueID}, "scope": map[string]any{"paths": []string{}}, "severity": "minor", "summary": "pair", "reasoning": "pair", "references": []any{}})
				p23AssertReviewSuccess(t, "pair Concern", result)
			}
			if code == "result-limit-exceeded" {
				return p24CodePair{a.workspace, "concern.list", []string{"concern", "list", "--subject", "issue:" + a.issueID, "--limit", "1"}}
			}
			p23AssertReviewSuccess(t, "dispose pair Concern", p23ConcernMutation(t, a, "dispose", "first-disposal", input))
		}
	}
	return p24InputPair(t, a.workspace, operation, a.actor, input)
}

func p24ObligationCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	cfg := p22ExecutionFixture{}
	if code == "human-outcome-required" {
		cfg.HumanObligations = true
	}
	a := p24FailureJudgmentAdapter(t, cfg)
	if code == "human-outcome-required" {
		// Only security remains unsatisfied after the release and gate links.
		release := p23RecordHumanOutcome(t, a, "pair-release", a.actor, "release", "approved")
		assertP23HumanOutcomeSuccess(t, release)
		releaseID := p21PathString(t, p23Envelope(t, release.Stdout), "result", "outcome", "id")
		assertP23HumanOutcomeSuccess(t, p23LinkHumanOutcome(t, a, "pair-release", "OBL-human-release", releaseID))
		p22AcceptedCall(t, a, "pair-baseline", "run.observe-budget", map[string]any{"value": 0})
		args, input := a.p22BuildArgs(t, p22ConformanceCase{Name: "pair-gate", Operation: "gate.run"})
		args = append(args, "--input", p21Input(t, a.workspace, "pair-qualifying-gate", input), "--idempotency-key", "pair-gate")
		out, stderr, err := closureRunCommand(a.workspace, a.binary, args...)
		if err != nil || stderr != "" {
			t.Fatalf("qualifying gate: %v %q %q", err, out, stderr)
		}
		receiptID := p24GateTerminal(t, out)["receipt_id"].(string)
		linked := p23RunInput(t, a, "evidence", "link", a.actor, "pair-gate-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": receiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}})
		assertP23HumanOutcomeSuccess(t, linked)
	}
	return p24InputPair(t, a.workspace, "envelope.complete", a.actor, map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1})
}

func p24OutcomeCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	a := p24FailureJudgmentAdapter(t, p22ExecutionFixture{HumanObligations: true})
	if code == "authority-missing" {
		return p24InputPair(t, a.workspace, "human-outcome.record", "actor-denied", map[string]any{"id": "OUT-pair", "class": "release", "subject": map[string]any{"type": "issue", "id": a.issueID}, "result": "approved", "scope": map[string]any{"issue_id": a.issueID}, "reasoning": "pair", "references": []any{}})
	}
	id := "OUT-missing"
	if code == "outcome-not-permitted" {
		result := p23RecordHumanOutcome(t, a, "pair-rejected", a.actor, "release", "rejected")
		assertP23HumanOutcomeSuccess(t, result)
		id = p21PathString(t, p23Envelope(t, result.Stdout), "result", "outcome", "id")
	}
	return p24InputPair(t, a.workspace, "evidence.link", a.actor, map[string]any{"obligation_id": "OBL-human-release", "evidence_type": "human-outcome", "receipt_id": nil, "review_id": nil, "human_outcome_id": id, "required_artifact_ids": []string{}})
}

func p24ReviewClaimPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, nil, "review")
	p22Write(t, a.workspace, "base.txt", "review adoption\n")
	a.p22Observe(t, "terminated")
	a.p22Expire(t)
	return p24InputPair(t, a.workspace, "run.recover", a.actor, map[string]any{"prior_run_id": a.runID, "expected_run_version": a.runVersion, "expected_claim_version": a.claimVersion, "observation_id": a.observation, "reconciliation_digest": a.digest, "adopted_paths": []string{"base.txt"}, "authority": map[string]any{"mode": "ordinary", "actor_id": a.actor, "reason": "pair", "reference": nil, "grant_ref": nil}, "executor": map[string]any{"id": "pair-successor", "source": "activated-harness"}, "ttl_seconds": 60})
}

func p24ReviewCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	cfg := p22ExecutionFixture{ReviewObligation: true}
	if code == "evidence-type-not-accepted" {
		cfg = p22ExecutionFixture{CoexistingReviewObligation: true}
	}
	a := p24FailureJudgmentAdapter(t, cfg)
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "pair-builder")
	actor := "actor-denied"
	if code == "review-not-independent" {
		actor = a.actor
	}
	claim := p23AcquireClaim(t, a, actor, "review", "pair-review")
	run := p23StartRun(t, a, actor, claim, "pair-review")
	input := p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "pass", nil, nil)
	if code == "review-not-independent" {
		return p24InputPair(t, a.workspace, "review.record", actor, input)
	}
	if code == "review-verdict-invalid" {
		// Blocked reviews cannot resolve prior Findings; domain validation checks
		// the verdict before looking up the referenced Finding.
		input["verdict"], input["resolves_finding_ids"] = "blocked", []any{"FND-unresolved"}
		return p24InputPair(t, a.workspace, "review.record", actor, input)
	}
	if code == "finding-already-disposed" {
		// Record a Finding, resolve it in a subsequent independent Review, then
		// address the terminal Finding through an otherwise valid mutation Claim.
		input["verdict"], input["findings"] = "changes-required", []any{p23Finding(a.issueID, "FND-pair", "major", true)}
	}
	if code == "review-changes-required" {
		input["verdict"] = "blocked"
		input["findings"] = []any{p23Finding(a.issueID, "FND-blocked-pair", "minor", false)}
	}
	result := p23ReviewRecord(t, a, actor, "pair-first-review", input)
	p23AssertReviewSuccess(t, "pair Review", result)
	reviewID := p23ReviewFromResponse(t, result)["id"].(string)
	if code != "finding-already-disposed" {
		obligation := "OBL-p22-corpus"

		return p24InputPair(t, a.workspace, "evidence.link", a.actor, map[string]any{"obligation_id": obligation, "evidence_type": "review", "receipt_id": nil, "review_id": reviewID, "human_outcome_id": nil, "required_artifact_ids": []string{}})
	}
	p23CloseAndRelease(t, a, actor, p21PathString(t, run, "id"), p21PathString(t, claim, "id"), "pair-first")
	nextClaim := p23AcquireClaim(t, a, actor, "review", "pair-resolve")
	nextRun := p23StartRun(t, a, actor, nextClaim, "pair-resolve")
	resolved := p23ReviewRecord(t, a, actor, "pair-resolve", p23ReviewInput(a, p21PathString(t, nextRun, "id"), p23Number(t, nextRun["version"]), "pass", nil, []any{"FND-pair"}))
	p23AssertReviewSuccess(t, "resolve pair Finding", resolved)
	p23CloseAndRelease(t, a, actor, p21PathString(t, nextRun, "id"), p21PathString(t, nextClaim, "id"), "pair-resolve")
	mutation := p23AcquireClaim(t, a, a.actor, "mutation", "pair-address")
	return p24InputPair(t, a.workspace, "finding.address", a.actor, map[string]any{"finding_id": "FND-pair", "claim_id": mutation["id"], "expected_claim_version": mutation["version"], "reason": "pair"})
}
