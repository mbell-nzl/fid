package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFinalSchemaNamedReadsAndExportMatchAfterRebuild(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID, evaluatePath := task21CompletedWorkspace(t)
	archivePath := filepath.Join(t.TempDir(), "task21-export.zip")
	before := task21SnapshotReads(t, workspace, issueID, envelopeID, evaluatePath, archivePath)

	if err := os.Remove(filepath.Join(workspace, ".fid", "fid.sqlite")); err != nil {
		t.Fatalf("delete projection: %v", err)
	}
	rebuilt := executeCommandForTest(workspace, "--project", workspace, "--json", "rebuild")
	if rebuilt.err != nil || rebuilt.stderr != "" {
		t.Fatalf("rebuild final-schema projection = %#v", rebuilt)
	}
	after := task21SnapshotReads(t, workspace, issueID, envelopeID, evaluatePath, archivePath)

	if !bytes.Equal(before.exportBytes, after.exportBytes) {
		t.Fatal("export bytes differ after final-schema rebuild")
	}
	if before.capabilities != after.capabilities || before.envelope != after.envelope || before.evaluation != after.evaluation || before.resume != after.resume || before.validation != after.validation {
		t.Fatalf("named reads differ after rebuild: before=%#v after=%#v", before, after)
	}
}

type task21ReadSnapshot struct {
	capabilities string
	envelope     string
	evaluation   string
	resume       string
	validation   string
	exportBytes  []byte
}

func task21SnapshotReads(t *testing.T, workspace, issueID, envelopeID, evaluatePath, archivePath string) task21ReadSnapshot {
	t.Helper()
	capabilities := task21Read(t, workspace, "--project", workspace, "--json", "capabilities")
	envelope := task21Read(t, workspace, "--project", workspace, "--json", "envelope", "get", envelopeID)
	evaluation := task21Read(t, workspace, "--project", workspace, "--json", "evidence", "evaluate", "--input", evaluatePath)
	resume := task21Read(t, workspace, "--project", workspace, "--json", "resume", issueID)
	_ = task21Read(t, workspace, "--project", workspace, "--json", "export", "create", archivePath)
	validation := task21Read(t, workspace, "--project", workspace, "--json", "export", "validate", archivePath)
	exportBytes, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("read export archive: %v", err)
	}
	return task21ReadSnapshot{capabilities, envelope, evaluation, resume, validation, exportBytes}
}

func task21Read(t *testing.T, workspace string, args ...string) string {
	t.Helper()
	result := executeCommandForTest(workspace, args...)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("read %v = %#v", args, result)
	}
	return result.stdout
}

func task21CompletedWorkspace(t *testing.T) (string, string, string, string) {
	t.Helper()
	workspace, issueID, envelopeID := observedGateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "task21-parity-gate")
	if gateRun.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("run qualifying Gate = %#v, terminal=%#v", gateRun, terminal)
	}
	repository := runRepositoryForTest(t, workspace)
	obligationID := testEvidenceObligationID
	linkPath := filepath.Join(workspace, "task21-link-evidence.json")
	writeActivationInput(t, linkPath, map[string]any{"obligation_id": obligationID, "evidence_type": "receipt", "receipt_id": terminal.Data.ReceiptID, "review_id": nil, "human_outcome_id": nil, "target_repository": repository, "required_artifact_ids": []string{}})
	task21Mutate(t, workspace, "evidence", "link", "--actor", "actor-builder", "--input", linkPath, "--idempotency-key", "task21-link-evidence")
	evaluatePath := filepath.Join(workspace, "task21-evaluate-evidence.json")
	writeActivationInput(t, evaluatePath, map[string]any{"obligation_id": obligationID, "target_repository": repository})

	for _, step := range []struct{ state, revision string }{{"in_progress", "1"}, {"dev_complete", "2"}, {"qa_review", "3"}, {"qa_pass", "4"}} {
		task21Mutate(t, workspace, "issue", "transition", issueID, "--to", step.state, "--actor", "agent:actor-builder", "--expected-revision", step.revision, "--idempotency-key", "task21-"+step.state)
	}
	runID, claimID := task18RunAndClaim(t, workspace, issueID)
	closeRunForTest(t, workspace, runID, repository, "task21-close-run")
	releasePath := filepath.Join(workspace, "task21-release-claim.json")
	writeActivationInput(t, releasePath, map[string]any{"claim_id": claimID, "expected_version": 1, "reason": "completed"})
	task21Mutate(t, workspace, "claim", "release", "--actor", "actor-builder", "--input", releasePath, "--idempotency-key", "task21-release-claim")
	completePath := filepath.Join(workspace, "task21-complete-envelope.json")
	writeActivationInput(t, completePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1, "repository": repository})
	task21Mutate(t, workspace, "envelope", "complete", "--actor", "actor-builder", "--input", completePath, "--idempotency-key", "task21-complete-envelope")
	return workspace, issueID, envelopeID, evaluatePath
}

func task21Mutate(t *testing.T, workspace string, args ...string) {
	t.Helper()
	result := executeCommandForTest(workspace, append([]string{"--project", workspace, "--json"}, args...)...)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("mutation %v = %#v", args, result)
	}
}
