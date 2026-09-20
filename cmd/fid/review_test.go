package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestReviewRecordRejectsNonIndependentActor is the Task 4 vertical tracer:
// recording under the mutating actor must be rejected before an Event is added.
func TestReviewRecordRejectsNonIndependentActor(t *testing.T) {
	workspace := t.TempDir()
	inputPath := filepath.Join(workspace, "review-record.json")
	writeActivationInput(t, inputPath, map[string]any{})

	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "review", "record", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "review-not-independent",
	)
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `\"code\":\"review-not-independent\"`) {
		t.Fatalf("review record by mutating actor = %#v, want typed review-not-independent rejection", result)
	}
}
