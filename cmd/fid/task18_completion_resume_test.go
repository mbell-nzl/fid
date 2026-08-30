package main

import (
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence(t *testing.T) {
	workspace, issueID, envelopeID := runEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/false"}
	})
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "task18-unsuccessful")
	if gateRun.err == nil || terminal.Data.ReceiptID == "" || terminal.Data.Outcome != "failed" {
		t.Fatalf("failed gate run = %#v, terminal=%#v, want failed Receipt", gateRun, terminal)
	}

	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	linkPath := filepath.Join(workspace, "task18-evidence-link.json")
	writeActivationInput(t, linkPath, map[string]any{
		"obligation_id": evidenceObligationID(t, workspace, envelopeID),
		"receipt_id":    terminal.Data.ReceiptID, "target_repository": repository,
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

func TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease(t *testing.T) {
	workspace, issueID, envelopeID := runEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "task18-qualified")
	if gateRun.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("qualified gate run = %#v, terminal=%#v", gateRun, terminal)
	}
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	linkPath := filepath.Join(workspace, "task18-qualified-evidence.json")
	obligationID := evidenceObligationID(t, workspace, envelopeID)
	writeActivationInput(t, linkPath, map[string]any{
		"obligation_id": obligationID, "receipt_id": terminal.Data.ReceiptID,
		"target_repository": repository, "required_artifact_ids": []string{},
	})
	linked := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", linkPath, "--idempotency-key", "task18-qualified-evidence")
	if linked.err != nil || linked.stderr != "" {
		t.Fatalf("link qualified Evidence = %#v", linked)
	}
	evaluatePath := filepath.Join(workspace, "task18-qualified-evaluate.json")
	writeActivationInput(t, evaluatePath, map[string]any{"obligation_id": obligationID, "target_repository": repository})
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
	finishPath := filepath.Join(workspace, "task18-finish-run.json")
	writeActivationInput(t, finishPath, map[string]any{"run_id": runID, "expected_status": "active", "outcome": "succeeded", "repository": repository})
	finished := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "finish", "--actor", "actor-builder", "--input", finishPath, "--idempotency-key", "task18-finish-run")
	if finished.err != nil || finished.stderr != "" {
		t.Fatalf("finish Run = %#v", finished)
	}
	releasePath := filepath.Join(workspace, "task18-release-claim.json")
	writeActivationInput(t, releasePath, map[string]any{"claim_id": claimID, "expected_version": 1, "reason": "completed"})
	released := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", releasePath, "--idempotency-key", "task18-release-claim")
	if released.err != nil || released.stderr != "" {
		t.Fatalf("release Claim = %#v", released)
	}

	completePath := filepath.Join(workspace, "task18-envelope-complete.json")
	writeActivationInput(t, completePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1, "repository": repository})
	beforeEvents := journalEventCount(t, workspace)
	completed := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "complete", "--actor", "actor-builder", "--input", completePath, "--idempotency-key", "task18-envelope-complete")
	if completed.err != nil || completed.stderr != "" {
		t.Fatalf("complete terminal Envelope = %#v", completed)
	}
	result := decodeCanonicalResult(t, completed.stdout)
	if result["envelope"].(map[string]any)["status"] != "completed" || !sameJSONMap(result["evaluation"].(map[string]any), wantEvaluation) {
		t.Fatalf("completion result = %#v, want completed Envelope and shared evaluator result %#v", result, wantEvaluation)
	}
	if afterEvents := journalEventCount(t, workspace); afterEvents != beforeEvents+1 {
		t.Fatalf("completion Event count = %d, want %d", afterEvents, beforeEvents+1)
	}
}

func TestResumeReadsDurableFactsWithoutMutation(t *testing.T) {
	workspace, issueID, envelopeID := runEvidenceEnvelopeForTest(t, nil)
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

func TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation(t *testing.T) {
	workspace, issueID, envelopeID := runEvidenceEnvelopeForTest(t, nil)
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	beforeEnvelope := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "get", envelopeID).stdout
	beforeIssue := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "get", issueID).stdout
	beforeEvents := journalEventCount(t, workspace)
	beforeEvidence, beforeLegacy := task18TableCounts(t, workspace)
	path := filepath.Join(workspace, "task18-incomplete-envelope.json")
	writeActivationInput(t, path, map[string]any{"envelope_id": envelopeID, "expected_revision": 1, "repository": repository})

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
	afterEvidence, afterLegacy := task18TableCounts(t, workspace)
	if afterEvidence != beforeEvidence || afterLegacy != beforeLegacy || journalEventCount(t, workspace) != beforeEvents {
		t.Fatalf("failed completion mutated counts Evidence/legacy/Events = %d/%d/%d, want %d/%d/%d", afterEvidence, afterLegacy, journalEventCount(t, workspace), beforeEvidence, beforeLegacy, beforeEvents)
	}
}

func task18RunAndClaim(t *testing.T, workspace, issueID string) (string, string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
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
	conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open task 18 count projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var evidence, legacy int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM evidence_records`).Scan(&evidence); err != nil {
		t.Fatalf("count Evidence: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM task_completions`).Scan(&legacy); err != nil {
		t.Fatalf("count legacy completions: %v", err)
	}
	return evidence, legacy
}
