//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestP23ReviewRecordRealProcessJourney keeps the VAL-015 real-process
// evidence at the built-binary boundary. The Run-only history direction is
// deliberately absent: DEC-239 confines that otherwise unreachable state to
// internal/review's isolated durable-history query fixture.
func TestP23ReviewRecordRealProcessJourney(t *testing.T) {
	t.Parallel()
	actions := []string{
		"claim.acquire", "claim.release", "envelope.activate", "envelope.create",
		"gate.register", "review.record", "run.close", "run.observe", "run.start",
	}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation")

	t.Run("mutating_claim_actor_is_rejected_without_an_event", func(t *testing.T) {
		before := p21JournalEventCount(t, a.workspace)
		response := p23ReviewRecord(t, a, a.actor, "builder-conflict", p23ReviewInput(a, a.runID, a.runVersion, "changes-required", []any{p23Finding(a.issueID, "FND-builder", "major", true)}, nil))
		assertP22ProcessFailure(t, "review.record", "review-not-independent", 4, response)
		for _, want := range []string{`"type":"claim"`, `"id":"` + a.claimID + `"`, `"conflicting_actor_id":"` + a.actor + `"`} {
			if !strings.Contains(response.Stderr, want) {
				t.Fatalf("review independence rejection = %q, want %s", response.Stderr, want)
			}
		}
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("mutating Claim rejection appended Event(s): %d -> %d", before, after)
		}
	})

	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "builder-before-review")

	firstClaim := p23AcquireClaim(t, a, "actor-denied", "review", "reviewer-first")
	firstRun := p23StartRun(t, a, "actor-denied", firstClaim, "reviewer-first")
	firstInput := p23ReviewInput(a, p21PathString(t, firstRun, "id"), p23Number(t, firstRun["version"]), "changes-required", []any{
		p23Finding(a.issueID, "FND-prior", "major", true),
		p23Finding(a.issueID, "FND-prior-two", "major", true),
	}, nil)

	t.Run("repository_mismatch_is_typed_and_non_mutating", func(t *testing.T) {
		mismatch := p23Clone(firstInput)
		repository := p23Clone(mismatch["repository"].(map[string]any))
		repository["head_oid"] = strings.Repeat("0", 40)
		mismatch["repository"] = repository
		before := p21JournalEventCount(t, a.workspace)
		response := p23ReviewRecord(t, a, "actor-denied", "repository-mismatch", mismatch)
		assertP22ProcessFailure(t, "review.record", "repository-mismatch", 4, response)
		if !strings.Contains(response.Stderr, `"expected"`) || !strings.Contains(response.Stderr, `"observed"`) {
			t.Fatalf("repository mismatch = %q, want expected and observed identities", response.Stderr)
		}
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("repository mismatch appended Event(s): %d -> %d", before, after)
		}
	})

	beforeFirst := p21JournalEventCount(t, a.workspace)
	first := p23ReviewRecord(t, a, "actor-denied", "reviewer-first", firstInput)
	p23AssertReviewSuccess(t, "first independent Review", first)
	if after := p21JournalEventCount(t, a.workspace); after != beforeFirst+1 {
		t.Fatalf("first Review Event count = %d, want %d", after, beforeFirst+1)
	}
	firstReview := p23ReviewFromResponse(t, first)

	t.Run("replay_duplicate_and_digest_conflict_preserve_history", func(t *testing.T) {
		before := p21JournalEventCount(t, a.workspace)
		replay := p23ReviewRecord(t, a, "actor-denied", "reviewer-first", firstInput)
		p23AssertReviewSuccess(t, "review replay", replay)
		if p21Path(t, p23Envelope(t, replay.Stdout), "result", "replayed") != true {
			t.Fatalf("review replay = %q, want replayed:true", replay.Stdout)
		}
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("review replay appended Event(s): %d -> %d", before, after)
		}

		conflictInput := p23Clone(firstInput)
		conflictInput["profile"] = "review-alternate"
		response := p23ReviewRecord(t, a, "actor-denied", "reviewer-first", conflictInput)
		assertP22ProcessFailure(t, "review.record", "idempotency-conflict", 4, response)
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("idempotency conflict appended Event(s): %d -> %d", before, after)
		}

		response = p23ReviewRecord(t, a, "actor-denied", "same-run", firstInput)
		assertP22ProcessFailure(t, "review.record", "duplicate-entity", 4, response)
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("second Review for Run appended Event(s): %d -> %d", before, after)
		}
	})

	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, firstRun, "id"), p21PathString(t, firstClaim, "id"), "reviewer-first")
	reworkClaim := p23AcquireClaim(t, a, a.actor, "mutation", "builder-rework")
	reworkRun := p23StartRun(t, a, a.actor, reworkClaim, "builder-rework")
	p23CloseAndRelease(t, a, a.actor, p21PathString(t, reworkRun, "id"), p21PathString(t, reworkClaim, "id"), "builder-rework")

	secondClaim := p23AcquireClaim(t, a, "actor-denied", "review", "reviewer-second")
	secondRun := p23StartRun(t, a, "actor-denied", secondClaim, "reviewer-second")
	beforeSecond := p21JournalEventCount(t, a.workspace)
	second := p23ReviewRecord(t, a, "actor-denied", "reviewer-second", p23ReviewInput(a, p21PathString(t, secondRun, "id"), p23Number(t, secondRun["version"]), "pass", []any{}, []any{"FND-prior", "FND-prior-two"}))
	p23AssertReviewSuccess(t, "same Reviewer post-rework Review", second)
	if after := p21JournalEventCount(t, a.workspace); after != beforeSecond+1 {
		t.Fatalf("second Review plus two resolutions changed Event count %d -> %d, want one atomic Event", beforeSecond, after)
	}
	secondReview := p23ReviewFromResponse(t, second)

	firstRead := p23ReviewGet(t, a, firstReview["id"].(string))
	secondRead := p23ReviewGet(t, a, secondReview["id"].(string))
	if firstRead["id"] != firstReview["id"] || firstRead["run_id"] != firstReview["run_id"] || firstRead["verdict"] != "changes-required" || !reflect.DeepEqual(firstRead["repository"], firstReview["repository"]) {
		t.Fatalf("first review readback = %#v, want immutable recorded binding %#v", firstRead, firstReview)
	}
	findings, ok := firstRead["findings"].([]any)
	if !ok || len(findings) != 2 {
		t.Fatalf("first review findings = %#v, want two projected Findings", firstRead["findings"])
	}
	for _, value := range findings {
		finding := value.(map[string]any)
		if finding["state"] != "resolved" || len(finding["dispositions"].([]any)) != 1 {
			t.Fatalf("resolved Finding = %#v, want one resolved disposition", finding)
		}
	}
	if !reflect.DeepEqual(secondRead, secondReview) || !reflect.DeepEqual(secondRead["resolved_finding_ids"], []any{"FND-prior", "FND-prior-two"}) {
		t.Fatalf("second review readback = %#v, want exact record result %#v with both resolved Finding IDs", secondRead, secondReview)
	}
}

func requiresAdditionalJudgmentActors(actions []string) bool {
	for _, action := range actions {
		if action == "review.record" || action == "concern.dispose" {
			return true
		}
	}
	return false
}

func containsReviewRecordAction(actions []string) bool {
	return containsAction(actions, "review.record")
}

func containsAction(actions []string, want string) bool {
	for _, action := range actions {
		if action == want {
			return true
		}
	}
	return false
}

func p23ReviewInput(a *p22FIDRecoveryAdapter, runID string, version int, verdict string, findings, resolves []any) map[string]any {
	return map[string]any{"run_id": runID, "expected_run_version": version, "profile": "review-default", "repository": a.repository, "coverage": map[string]any{"items": []any{map[string]any{"source": a.measurement, "assessed": true}}, "exclusions": []any{}}, "verdict": verdict, "findings": findings, "resolves_finding_ids": resolves, "evidence_refs": []any{}, "receipt_ids": []any{}}
}

func p23Finding(issueID, id, severity string, blocking bool) map[string]any {
	return map[string]any{"id": id, "subject": map[string]any{"type": "issue", "id": issueID}, "scope": map[string]any{"paths": []any{}}, "severity": severity, "summary": "real-process Review finding", "blocking": blocking}
}

func p23ReviewRecord(t *testing.T, a *p22FIDRecoveryAdapter, actor, key string, input map[string]any) p22ConformanceResponse {
	t.Helper()
	return p23RunInput(t, a, "review", "record", actor, key, input)
}

func p23RunInput(t *testing.T, a *p22FIDRecoveryAdapter, group, command, actor, key string, input map[string]any) p22ConformanceResponse {
	t.Helper()
	args := []string{"--project", a.workspace, "--json", group, command, "--actor", actor, "--input", p21Input(t, a.workspace, "p23-"+key, input), "--idempotency-key", key}
	stdout, stderr, err := closureRunCommand(a.workspace, a.binary, args...)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return response
}

func p23AcquireClaim(t *testing.T, a *p22FIDRecoveryAdapter, actor, mode, key string) map[string]any {
	t.Helper()
	response := p23RunInput(t, a, "claim", "acquire", actor, key, map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "worktree_id": a.repository["worktree_id"], "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": mode, "ttl_seconds": 3600})
	p23AssertReviewSuccess(t, "acquire "+mode+" Claim", response)
	return p21Path(t, p23Envelope(t, response.Stdout), "result", "claim").(map[string]any)
}

func p23StartRun(t *testing.T, a *p22FIDRecoveryAdapter, actor string, claim map[string]any, key string) map[string]any {
	t.Helper()
	response := p23RunInput(t, a, "run", "start", actor, key, map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "claim_id": claim["id"], "claim_version": claim["version"], "executor": map[string]any{"id": "executor-" + key, "source": "activated-harness"}, "obligation_id": nil, "repository": a.repository})
	p23AssertReviewSuccess(t, "start Run", response)
	return p21Path(t, p23Envelope(t, response.Stdout), "result", "run").(map[string]any)
}

func p23CloseAndRelease(t *testing.T, a *p22FIDRecoveryAdapter, actor, runID, claimID, key string) {
	t.Helper()
	run := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "run", "get", runID)
	version := p23Number(t, p21Path(t, run, "result", "run", "version"))
	executorID := p21PathString(t, run, "result", "run", "executor", "id")
	observed := p23RunInput(t, a, "run", "observe", actor, key+"-observe", map[string]any{"run_id": runID, "expected_version": version, "liveness": map[string]any{"executor_id": executorID, "state": "terminated", "source_ref": "process:p23-review", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
	p23AssertReviewSuccess(t, "observe Run", observed)
	observedEnvelope := p23Envelope(t, observed.Stdout)
	close := p23RunInput(t, a, "run", "close", actor, key+"-close", map[string]any{"run_id": runID, "expected_version": p21Path(t, observedEnvelope, "result", "run", "version"), "observation_id": p21PathString(t, observedEnvelope, "result", "observation", "id"), "repository": a.repository})
	p23AssertReviewSuccess(t, "close Run", close)
	release := p23RunInput(t, a, "claim", "release", actor, key+"-release", map[string]any{"claim_id": claimID, "expected_version": 1, "reason": "completed"})
	p23AssertReviewSuccess(t, "release Claim", release)
}

func p23ReviewGet(t *testing.T, a *p22FIDRecoveryAdapter, reviewID string) map[string]any {
	t.Helper()
	response := a.p22Run(t, "review", "get", reviewID)
	p23AssertReviewSuccess(t, "review get", response)
	return p21Path(t, p23Envelope(t, response.Stdout), "result", "review").(map[string]any)
}

func p23ReviewFromResponse(t *testing.T, response p22ConformanceResponse) map[string]any {
	t.Helper()
	return p21Path(t, p23Envelope(t, response.Stdout), "result", "review").(map[string]any)
}

func p23Envelope(t *testing.T, stdout string) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("decode Review response: %v; stdout=%q", err, stdout)
	}
	return result
}

func p23AssertReviewSuccess(t *testing.T, operation string, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("%s = %#v, want stdout-only success", operation, response)
	}
}

func p23EnvelopeID(t *testing.T, a *p22FIDRecoveryAdapter) string {
	t.Helper()
	run := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "run", "get", a.runID)
	return p21PathString(t, run, "result", "run", "envelope_id")
}

func p23Number(t *testing.T, value any) int {
	t.Helper()
	number, ok := value.(float64)
	if !ok {
		t.Fatalf("number = %#v, want float64 JSON number", value)
	}
	return int(number)
}

func p23Clone(value map[string]any) map[string]any {
	payload, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	var clone map[string]any
	if err := json.Unmarshal(payload, &clone); err != nil {
		panic(err)
	}
	return clone
}
