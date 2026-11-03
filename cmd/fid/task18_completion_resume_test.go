package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := observedGateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/false"}
	})
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "task18-unsuccessful")
	if gateRun.err == nil || terminal.Data.ReceiptID == "" || terminal.Data.Outcome != "failed" {
		t.Fatalf("failed gate run = %#v, terminal=%#v, want failed Receipt", gateRun, terminal)
	}

	linkPath := filepath.Join(workspace, "task18-evidence-link.json")
	writeActivationInput(t, linkPath, map[string]any{
		"obligation_id": testEvidenceObligationID, "evidence_type": "receipt",
		"receipt_id": terminal.Data.ReceiptID, "review_id": nil, "human_outcome_id": nil,
		"required_artifact_ids": []string{},
	})
	linked := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", linkPath, "--idempotency-key", "task18-unsuccessful-evidence")
	if linked.err != nil || linked.stderr != "" {
		t.Fatalf("link failed Receipt Evidence = %#v, want immutable linked Evidence", linked)
	}

	for _, step := range []struct {
		state    string
		revision string
	}{
		{state: "in_progress", revision: "1"},
		{state: "dev_complete", revision: "2"},
		{state: "qa_review", revision: "3"},
	} {
		result := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "transition", issueID, "--to", step.state, "--actor", "agent:actor-builder", "--expected-revision", step.revision, "--idempotency-key", "task18-before-"+step.state)
		if result.err != nil || result.stderr != "" {
			t.Fatalf("transition to %s = %#v, want success", step.state, result)
		}
	}

	assertExecutionTransitionRejectionWithoutMutation(t, workspace, issueID, "qa_pass", "agent:actor-builder", "4", "task18-terminal-unsuccessful", []string{"receipt-unsuccessful"})
}

func TestTerminalSuccessTransitionRejectsEveryOtherNonqualifyingEvidenceState(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		withArtifact bool
		mutate       func(t *testing.T, workspace, evidenceID, artifactID string)
		wantReason   string
	}{
		{name: "missing", wantReason: "obligation-evidence-absent"},
		{
			name: "invalid",
			mutate: func(t *testing.T, workspace, evidenceID, artifactID string) {
				t.Helper()
				task18MutateEvidenceFixture(t, workspace, `UPDATE evidence_records SET validity = 'invalid' WHERE id = ?`, evidenceID)
			},
			wantReason: "evidence-invalid",
		},
		{
			name: "stale",
			mutate: func(t *testing.T, workspace, evidenceID, artifactID string) {
				t.Helper()
				staleRepository := map[string]any{}
				if err := json.Unmarshal([]byte(task18EvidenceRepository(t, workspace, evidenceID)), &staleRepository); err != nil {
					t.Fatalf("decode Evidence repository: %v", err)
				}
				staleRepository["tree_oid"] = strings.Repeat("f", 40)
				encoded, err := json.Marshal(staleRepository)
				if err != nil {
					t.Fatalf("encode stale Evidence repository: %v", err)
				}
				task18MutateEvidenceFixture(t, workspace, `UPDATE evidence_records SET target_repository_json = ? WHERE id = ?`, string(encoded), evidenceID)
			},
			wantReason: "evidence-stale",
		},
		{
			name:         "Artifact missing",
			withArtifact: true,
			mutate: func(t *testing.T, workspace, evidenceID, artifactID string) {
				t.Helper()
				if artifactID == "" {
					t.Fatal("Artifact-missing fixture has no required Artifact")
				}
				task18MutateEvidenceFixture(t, workspace, `UPDATE referenced_artifacts SET availability = 'missing' WHERE id = ?`, artifactID)
			},
			wantReason: "artifact-unavailable",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			workspace, issueID, _, evidenceID, artifactID := task18TerminalGuardFixture(t, testCase.withArtifact, testCase.name != "missing")
			if testCase.mutate != nil {
				testCase.mutate(t, workspace, evidenceID, artifactID)
			}
			assertExecutionTransitionRejectionWithoutMutation(t, workspace, issueID, "qa_pass", "agent:actor-builder", "4", "task18-terminal-"+strings.ReplaceAll(testCase.name, " ", "-"), []string{testCase.wantReason})
		})
	}
}

func TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := observedGateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "task18-qualified")
	if gateRun.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("qualified gate run = %#v, terminal=%#v", gateRun, terminal)
	}
	linkPath := filepath.Join(workspace, "task18-qualified-evidence.json")
	obligationID := testEvidenceObligationID
	writeActivationInput(t, linkPath, map[string]any{
		"obligation_id": obligationID, "evidence_type": "receipt", "receipt_id": terminal.Data.ReceiptID, "review_id": nil, "human_outcome_id": nil,
		"required_artifact_ids": []string{},
	})
	linked := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", linkPath, "--idempotency-key", "task18-qualified-evidence")
	if linked.err != nil || linked.stderr != "" {
		t.Fatalf("link qualified Evidence = %#v", linked)
	}
	evaluatePath := filepath.Join(workspace, "task18-qualified-evaluate.json")
	writeActivationInput(t, evaluatePath, map[string]any{"obligation_id": obligationID})
	evaluated := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "evaluate", "--input", evaluatePath)
	if evaluated.err != nil || evaluated.stderr != "" {
		t.Fatalf("evaluate qualified Evidence = %#v", evaluated)
	}
	wantEvaluation := decodeCanonicalResult(t, evaluated.stdout)["evaluation"].(map[string]any)

	for _, step := range []struct {
		state    string
		revision string
	}{
		{state: "in_progress", revision: "1"},
		{state: "dev_complete", revision: "2"},
		{state: "qa_review", revision: "3"},
		{state: "qa_pass", revision: "4"},
	} {
		result := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "transition", issueID, "--to", step.state, "--actor", "agent:actor-builder", "--expected-revision", step.revision, "--idempotency-key", "task18-qualified-"+step.state)
		if result.err != nil || result.stderr != "" {
			t.Fatalf("transition to %s = %#v", step.state, result)
		}
	}

	runID, claimID := task18RunAndClaim(t, workspace, issueID)
	closeRunForTest(t, workspace, runID, runRepositoryForTest(t, workspace), "task18-close-run")
	releasePath := filepath.Join(workspace, "task18-release-claim.json")
	writeActivationInput(t, releasePath, map[string]any{"claim_id": claimID, "expected_version": 1, "reason": "completed"})
	released := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", releasePath, "--idempotency-key", "task18-release-claim")
	if released.err != nil || released.stderr != "" {
		t.Fatalf("release Claim = %#v", released)
	}

	completePath := filepath.Join(workspace, "task18-envelope-complete.json")
	writeActivationInput(t, completePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	beforeEvents := journalEventCount(t, workspace)
	completed := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "complete", "--actor", "actor-builder", "--input", completePath, "--idempotency-key", "task18-envelope-complete")
	if completed.err != nil || completed.stderr != "" {
		t.Fatalf("complete terminal Envelope = %#v", completed)
	}
	result := decodeCanonicalResult(t, completed.stdout)
	evaluations := result["evaluation"].([]any)
	if result["envelope"].(map[string]any)["status"] != "completed" || len(evaluations) != 1 || !sameJSONMap(evaluations[0].(map[string]any), wantEvaluation) {
		t.Fatalf("completion result = %#v, want completed Envelope and one shared evaluator result %#v", result, wantEvaluation)
	}
	if afterEvents := journalEventCount(t, workspace); afterEvents != beforeEvents+1 {
		t.Fatalf("completion Event count = %d, want %d", afterEvents, beforeEvents+1)
	}
}

func TestResumeReadsDurableFactsWithoutMutation(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := observedEvidenceEnvelopeForTest(t, nil)
	beforeEvents := journalEventCount(t, workspace)
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "resume", issueID)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("fresh resume command = %#v, want bounded read", result)
	}
	resume := decodeCanonicalResult(t, result.stdout)["resume"].(map[string]any)
	if resume["envelope"].(map[string]any)["id"] != envelopeID || resume["issue"].(map[string]any)["id"] != issueID {
		t.Fatalf("resume projection = %#v, want Envelope %q and Issue %q", resume, envelopeID, issueID)
	}
	if afterEvents := journalEventCount(t, workspace); afterEvents != beforeEvents {
		t.Fatalf("resume Event count = %d, want unchanged %d", afterEvents, beforeEvents)
	}
}

func TestResumeReportsRealFindingsAndJudgmentSections(t *testing.T) {
	t.Parallel()
	workspace, issueID, _ := observedEvidenceEnvelopeForTest(t, nil)
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "resume", issueID)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("resume command = %#v", result)
	}
	resume := decodeCanonicalResult(t, result.stdout)["resume"].(map[string]any)
	for _, section := range []string{"findings", "reviews", "concerns", "deviations", "decisions", "human_outcomes", "protected_decisions"} {
		if _, ok := resume[section].([]any); !ok {
			t.Errorf("resume %s = %#v, want array", section, resume[section])
		}
	}
}

func TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery(t *testing.T) {
	t.Parallel()
	t.Run("active", func(t *testing.T) {
		t.Parallel()
		workspace, issueID, envelopeID := observedGateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
			definition["argv"] = []string{"/bin/true"}
		})
		terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "task18-resume-active")
		if gateRun.err != nil || terminal.Data.ReceiptID == "" {
			t.Fatalf("active resume gate run = %#v, terminal=%#v", gateRun, terminal)
		}
		beforeEvents := journalEventCount(t, workspace)
		result := executeCommandForTest(workspace, "--project", workspace, "--json", "resume", issueID)
		if result.err != nil || result.stderr != "" {
			t.Fatalf("active resume command = %#v, want bounded read", result)
		}
		resume := decodeCanonicalResult(t, result.stdout)["resume"].(map[string]any)
		assertTask18ResumeProjection(t, resume, envelopeID, issueID, "active", "todo", 1, true, "active", "incomplete", []string{"obligation-evidence-absent"}, []string{"obligation-evidence-absent", "issue-not-terminal-success"}, []string{"run.inspect", "run.observe", "gate.run", "claim.renew", "claim.release", "concern.record", "deviation.record", "decision.record-reference", "human-outcome.record"})
		if afterEvents := journalEventCount(t, workspace); afterEvents != beforeEvents {
			t.Fatalf("active resume Event count = %d, want unchanged %d", afterEvents, beforeEvents)
		}
	})

	t.Run("completed", func(t *testing.T) {
		t.Parallel()
		workspace, issueID, envelopeID := observedGateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
			definition["argv"] = []string{"/bin/true"}
		})
		terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "task18-resume-completed")
		if gateRun.err != nil || terminal.Data.ReceiptID == "" {
			t.Fatalf("completed resume gate run = %#v, terminal=%#v", gateRun, terminal)
		}
		repository := runRepositoryForTest(t, workspace)
		obligationID := testEvidenceObligationID
		linkPath := filepath.Join(workspace, "task18-resume-completed-evidence.json")
		writeActivationInput(t, linkPath, map[string]any{"obligation_id": obligationID, "evidence_type": "receipt", "receipt_id": terminal.Data.ReceiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}})
		linked := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", linkPath, "--idempotency-key", "task18-resume-completed-evidence")
		if linked.err != nil || linked.stderr != "" {
			t.Fatalf("completed resume Evidence link = %#v", linked)
		}
		task18AdvanceToTerminalGuard(t, workspace, issueID)
		if err := os.WriteFile(filepath.Join(workspace, "current-target.txt"), []byte("one byte changed\n"), 0o600); err != nil {
			t.Fatalf("change current observed target: %v", err)
		}
		beforeStaleTransition := journalEventCount(t, workspace)
		staleTransition := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "transition", issueID, "--to", "qa_pass", "--actor", "agent:actor-builder", "--expected-revision", "4", "--idempotency-key", "task18-resume-terminal-stale")
		if staleTransition.err == nil || staleTransition.stdout != "" || !strings.Contains(staleTransition.stderr, `"evidence-stale"`) {
			t.Fatalf("stale current-target terminal transition = %#v, want evidence-stale rejection", staleTransition)
		}
		if after := journalEventCount(t, workspace); after != beforeStaleTransition {
			t.Fatalf("stale terminal transition event count = %d, want unchanged %d", after, beforeStaleTransition)
		}
		if err := os.Remove(filepath.Join(workspace, "current-target.txt")); err != nil {
			t.Fatalf("restore current observed target: %v", err)
		}
		transitioned := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "transition", issueID, "--to", "qa_pass", "--actor", "agent:actor-builder", "--expected-revision", "4", "--idempotency-key", "task18-resume-terminal")
		if transitioned.err != nil || transitioned.stderr != "" {
			t.Fatalf("terminal transition for completed resume = %#v", transitioned)
		}
		runID, claimID := task18RunAndClaim(t, workspace, issueID)
		closeRunForTest(t, workspace, runID, repository, "task18-resume-close")
		releasePath := filepath.Join(workspace, "task18-resume-release.json")
		writeActivationInput(t, releasePath, map[string]any{"claim_id": claimID, "expected_version": 1, "reason": "completed"})
		if released := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", releasePath, "--idempotency-key", "task18-resume-release"); released.err != nil || released.stderr != "" {
			t.Fatalf("release completed resume Claim = %#v", released)
		}
		completePath := filepath.Join(workspace, "task18-resume-complete.json")
		writeActivationInput(t, completePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
		if err := os.WriteFile(filepath.Join(workspace, "current-target.txt"), []byte("one byte changed\n"), 0o600); err != nil {
			t.Fatalf("change complete current target: %v", err)
		}
		beforeStaleComplete := journalEventCount(t, workspace)
		staleComplete := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "complete", "--actor", "actor-builder", "--input", completePath, "--idempotency-key", "task18-resume-complete-stale")
		if staleComplete.err == nil || staleComplete.stdout != "" || !strings.Contains(staleComplete.stderr, `"evidence-stale"`) {
			t.Fatalf("stale current-target Envelope completion = %#v, want evidence-stale rejection", staleComplete)
		}
		if after := journalEventCount(t, workspace); after != beforeStaleComplete {
			t.Fatalf("stale Envelope completion event count = %d, want unchanged %d", after, beforeStaleComplete)
		}
		if err := os.Remove(filepath.Join(workspace, "current-target.txt")); err != nil {
			t.Fatalf("restore complete current target: %v", err)
		}
		if completed := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "complete", "--actor", "actor-builder", "--input", completePath, "--idempotency-key", "task18-resume-complete"); completed.err != nil || completed.stderr != "" {
			t.Fatalf("complete resume Envelope = %#v", completed)
		}

		beforeEvents := journalEventCount(t, workspace)
		result := executeCommandForTest(workspace, "--project", workspace, "--json", "resume", issueID)
		if result.err != nil || result.stderr != "" {
			t.Fatalf("completed resume command = %#v, want bounded read", result)
		}
		resume := decodeCanonicalResult(t, result.stdout)["resume"].(map[string]any)
		assertTask18ResumeProjection(t, resume, envelopeID, issueID, "completed", "qa_pass", 5, false, "succeeded", "complete", []string{}, []string{}, []string{})
		if afterEvents := journalEventCount(t, workspace); afterEvents != beforeEvents {
			t.Fatalf("completed resume Event count = %d, want unchanged %d", afterEvents, beforeEvents)
		}
	})
}

func TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := observedEvidenceEnvelopeForTest(t, nil)
	beforeEnvelope := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "get", envelopeID).stdout
	beforeIssue := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "get", issueID).stdout
	beforeEvents := journalEventCount(t, workspace)
	beforeEvidence, beforeConcerns := task18TableCounts(t, workspace)
	path := filepath.Join(workspace, "task18-incomplete-envelope.json")
	writeActivationInput(t, path, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})

	result := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "complete", "--actor", "actor-builder", "--input", path, "--idempotency-key", "task18-incomplete-envelope")
	if result.err == nil || result.stdout != "" {
		t.Fatalf("incomplete Envelope completion = %#v, want typed rejection", result)
	}
	if result.stderr == "" {
		t.Fatalf("incomplete Envelope completion has no failure envelope: %v", result.err)
	}
	failure := decodeCommandJSONStructure(t, result.stderr).(map[string]any)
	if failure["operation"] != "envelope.complete" || failure["error"].(map[string]any)["code"] != "obligation-unsatisfied" {
		t.Fatalf("completion failure = %#v, want obligation-unsatisfied", failure)
	}
	reasons := failure["error"].(map[string]any)["details"].(map[string]any)["reasons"].([]any)
	gotReasons := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		gotReasons = append(gotReasons, reason.(map[string]any)["code"].(string))
	}
	if len(gotReasons) != 2 || gotReasons[0] != "obligation-evidence-absent" || gotReasons[1] != "issue-not-terminal-success" {
		t.Fatalf("completion reasons = %#v, want every current reason", gotReasons)
	}
	if after := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "get", envelopeID).stdout; after != beforeEnvelope {
		t.Fatalf("Envelope changed after failed completion")
	}
	if after := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "get", issueID).stdout; after != beforeIssue {
		t.Fatalf("Issue changed after failed completion")
	}
	afterEvidence, afterConcerns := task18TableCounts(t, workspace)
	if afterEvidence != beforeEvidence || afterConcerns != beforeConcerns || journalEventCount(t, workspace) != beforeEvents {
		t.Fatalf("failed completion mutated counts Evidence/Concerns/Events = %d/%d/%d, want %d/%d/%d", afterEvidence, afterConcerns, journalEventCount(t, workspace), beforeEvidence, beforeConcerns, beforeEvents)
	}
}

func task18RunAndClaim(t *testing.T, workspace, issueID string) (string, string) {
	t.Helper()
	binding := envelopeBindingForWorkspace(t, workspace)
	if binding.runID != "" && binding.claimID != "" {
		return binding.runID, binding.claimID
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open task 18 projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var runID, claimID string
	if err := conn.QueryRow(`SELECT id, claim_id FROM runs WHERE issue_id = ?`, issueID).Scan(&runID, &claimID); err != nil {
		t.Fatalf("load task 18 Run/Claim: %v", err)
	}
	return runID, claimID
}

func task18TableCounts(t *testing.T, workspace string) (int, int) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open task 18 count projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var evidence, concerns int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM evidence_records`).Scan(&evidence); err != nil {
		t.Fatalf("count Evidence: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM concerns`).Scan(&concerns); err != nil {
		t.Fatalf("count generic Concerns: %v", err)
	}
	return evidence, concerns
}

func task18QualifiedEvidenceFixture(t *testing.T, withArtifact, linkEvidence bool) (string, string, string, string, string) {
	t.Helper()
	workspace, issueID, envelopeID := observedGateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		if withArtifact {
			definition["artifact_paths"] = []string{"task18-evidence.txt"}
		}
	})
	if withArtifact {
		// This declared Gate output is material, not a project-source change.
		// Keep its generation outside the repository freshness target.
		if err := os.WriteFile(filepath.Join(workspace, ".git", "info", "exclude"), []byte("/task18-evidence.txt\n"), 0o600); err != nil {
			t.Fatalf("exclude generated Gate Artifact: %v", err)
		}
		if err := os.WriteFile(filepath.Join(workspace, "task18-evidence.txt"), []byte("terminal evidence\n"), 0o600); err != nil {
			t.Fatalf("write terminal Evidence Artifact: %v", err)
		}
	}
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "task18-qualified-"+strconv.FormatBool(withArtifact))
	if gateRun.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("qualified terminal gate run = %#v, terminal=%#v", gateRun, terminal)
	}
	obligationID := testEvidenceObligationID
	artifactID := ""
	if withArtifact {
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
		if err != nil {
			t.Fatalf("open terminal Artifact projection: %v", err)
		}
		if err := conn.QueryRow(`SELECT artifact_id FROM receipt_artifacts WHERE receipt_id = ?`, terminal.Data.ReceiptID).Scan(&artifactID); err != nil {
			_ = conn.Close()
			t.Fatalf("load terminal Evidence Artifact: %v", err)
		}
		if err := conn.Close(); err != nil {
			t.Fatalf("close terminal Artifact projection: %v", err)
		}
	}
	if !linkEvidence {
		return workspace, issueID, envelopeID, "", artifactID
	}
	linkPath := filepath.Join(workspace, "task18-terminal-evidence.json")
	artifactIDs := []string{}
	if artifactID != "" {
		artifactIDs = append(artifactIDs, artifactID)
	}
	writeActivationInput(t, linkPath, map[string]any{"obligation_id": obligationID, "evidence_type": "receipt", "receipt_id": terminal.Data.ReceiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": artifactIDs})
	linked := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", linkPath, "--idempotency-key", "task18-terminal-evidence-link")
	if linked.err != nil || linked.stderr != "" {
		t.Fatalf("link terminal Evidence = %#v", linked)
	}
	evidenceID := decodeCanonicalResult(t, linked.stdout)["evidence"].(map[string]any)["id"].(string)
	return workspace, issueID, envelopeID, evidenceID, artifactID
}

func task18TerminalGuardFixture(t *testing.T, withArtifact, linkEvidence bool) (string, string, string, string, string) {
	t.Helper()
	workspace, issueID, envelopeID, evidenceID, artifactID := task18QualifiedEvidenceFixture(t, withArtifact, linkEvidence)
	task18AdvanceToTerminalGuard(t, workspace, issueID)
	return workspace, issueID, envelopeID, evidenceID, artifactID
}

func task18AdvanceToTerminalGuard(t *testing.T, workspace, issueID string) {
	t.Helper()
	for _, step := range []struct {
		state    string
		revision string
	}{
		{state: "in_progress", revision: "1"},
		{state: "dev_complete", revision: "2"},
		{state: "qa_review", revision: "3"},
	} {
		result := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "transition", issueID, "--to", step.state, "--actor", "agent:actor-builder", "--expected-revision", step.revision, "--idempotency-key", "task18-before-"+step.state)
		if result.err != nil || result.stderr != "" {
			t.Fatalf("transition to %s = %#v, want success", step.state, result)
		}
	}
}

func task18MutateEvidenceFixture(t *testing.T, workspace, statement string, args ...any) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open terminal Evidence projection: %v", err)
	}
	if _, err := conn.Exec(statement, args...); err != nil {
		_ = conn.Close()
		t.Fatalf("mutate terminal Evidence fixture: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close terminal Evidence projection: %v", err)
	}
}

func task18EvidenceRepository(t *testing.T, workspace, evidenceID string) string {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open Evidence repository projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var repository string
	if err := conn.QueryRow(`SELECT target_repository_json FROM evidence_records WHERE id = ?`, evidenceID).Scan(&repository); err != nil {
		t.Fatalf("load Evidence repository: %v", err)
	}
	return repository
}

func assertTask18ResumeProjection(t *testing.T, resume map[string]any, envelopeID, issueID, envelopeStatus, issueStatus string, issueRevision int, wantClaim bool, wantRunStatus, completionStatus string, wantObligationReasons, wantCompletionReasons, wantActions []string) {
	t.Helper()
	envelope := resume["envelope"].(map[string]any)
	issue := resume["issue"].(map[string]any)
	if envelope["id"] != envelopeID || envelope["status"] != envelopeStatus || issue["id"] != issueID || issue["status"] != issueStatus || issue["revision"] != json.Number(strconv.Itoa(issueRevision)) {
		t.Fatalf("resume Envelope/Issue = %#v/%#v, want Envelope %q/%s and Issue %q/%s revision %d", envelope, issue, envelopeID, envelopeStatus, issueID, issueStatus, issueRevision)
	}
	claim := resume["claim"]
	if wantClaim {
		claimMap, ok := claim.(map[string]any)
		if !ok || claimMap["issue_id"] != issueID || claimMap["envelope_id"] != envelopeID || claimMap["status"] != "active" {
			t.Fatalf("resume Claim = %#v, want active Claim for Issue %q and Envelope %q", claim, issueID, envelopeID)
		}
	} else if claim != nil {
		t.Fatalf("resume Claim = %#v, want null", claim)
	}
	run := resume["run"]
	if wantRunStatus == "" {
		if run != nil {
			t.Fatalf("resume Run = %#v, want null", run)
		}
	} else if runMap, ok := run.(map[string]any); !ok || runMap["issue_id"] != issueID || runMap["envelope_id"] != envelopeID || runMap["status"] != wantRunStatus {
		t.Fatalf("resume Run = %#v, want %s Run for Issue %q and Envelope %q", run, wantRunStatus, issueID, envelopeID)
	}
	obligations := resume["obligations"].([]any)
	if len(obligations) != 1 {
		t.Fatalf("resume Obligations = %#v, want exactly one", obligations)
	}
	obligation := obligations[0].(map[string]any)
	if obligation["status"] != map[bool]string{true: "satisfied", false: "unsatisfied"}[len(wantObligationReasons) == 0] || !reflect.DeepEqual(task18ReasonCodes(obligation["reasons"].([]any)), wantObligationReasons) {
		t.Fatalf("resume Obligation = %#v, want reasons %#v", obligation, wantObligationReasons)
	}
	completion := resume["completion"].(map[string]any)
	if completion["status"] != completionStatus || !reflect.DeepEqual(task18ReasonCodes(completion["reasons"].([]any)), wantCompletionReasons) {
		t.Fatalf("resume completion = %#v, want %s/%#v", completion, completionStatus, wantCompletionReasons)
	}
	actionItems := resume["next_actions"].([]any)
	gotActions := make([]string, 0, len(actionItems))
	for _, item := range actionItems {
		gotActions = append(gotActions, item.(map[string]any)["operation"].(string))
	}
	if !reflect.DeepEqual(gotActions, wantActions) {
		t.Fatalf("resume next actions = %#v, want %#v", gotActions, wantActions)
	}
	for _, forbidden := range []string{"recover", "respawn", "run-recover", "run-respawn"} {
		for _, action := range gotActions {
			if action == forbidden {
				t.Fatalf("resume next actions = %#v, must not expose %q", gotActions, forbidden)
			}
		}
	}
}

func task18ReasonCodes(reasons []any) []string {
	result := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		result = append(result, reason.(map[string]any)["code"].(string))
	}
	return result
}
