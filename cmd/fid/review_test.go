package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestReviewRecordRejectsNonIndependentActor is the Task 4 vertical tracer:
// recording under the mutating actor must be rejected before an Event is added.
func TestReviewRecordRejectsNonIndependentActor(t *testing.T) {
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "review-independent-builder")
	repository := runRepositoryForTest(t, workspace)
	runPath := filepath.Join(workspace, "review-independent-run.json")
	writeActivationInput(t, runPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace,
		"--project", workspace, "--json", "run", "start", "--actor", "actor-builder",
		"--input", runPath, "--idempotency-key", "review-independent-run",
	)
	if started.err != nil || started.stderr != "" {
		t.Fatalf("start mutating Run = %#v, want stdout-only success", started)
	}
	run := runFromResponse(t, started.stdout)
	inputPath := filepath.Join(workspace, "review-record.json")
	writeActivationInput(t, inputPath, map[string]any{"run_id": run["id"]})

	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "review", "record", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "review-not-independent",
	)
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"review-not-independent"`) {
		t.Fatalf("review record by mutating actor = %#v, want typed review-not-independent rejection", result)
	}
}
