package main

import (
	"bytes"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestRunRecoveryErrorsUseClosedRecoveryEnvelopes(t *testing.T) {
	tests := []struct {
		name     string
		cause    error
		runID    string
		digest   string
		wantCode string
	}{
		{"liveness unknown", &workflow.ExecutorLivenessUnknownError{RunID: "RUN-1", ExecutorID: "EXEC-1", ObservationID: "OBS-1"}, "RUN-1", "", "executor-liveness-unknown"},
		{"still live", &workflow.ExecutorStillLiveError{RunID: "RUN-1", ExecutorID: "EXEC-1", ObservationID: "OBS-1"}, "RUN-1", "", "executor-still-live"},
		{"state conflict", &workflow.RecoveryStateConflictError{Paths: []string{"base.txt"}, Reasons: []string{"changed"}}, "RUN-1", "sha256:digest", "recovery-state-conflict"},
		{"missing observation", &workflow.ObservationNotFoundError{ID: "OBS-1"}, "RUN-1", "", "observation-not-found"},
		{"linked successor", &workflow.RecoveryAlreadyLinkedError{RunID: "RUN-1", SuccessorRunID: "RUN-2"}, "RUN-1", "", "recovery-already-linked"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			cmd := newRootCommandAt(t.TempDir())
			cmd.SetErr(stderr)
			if err := cmd.PersistentFlags().Set("json", "true"); err != nil {
				t.Fatalf("enable JSON: %v", err)
			}
			if err := renderRunRecoveryError(cmd, "run.recover", testCase.cause, testCase.runID, testCase.digest); err == nil {
				t.Fatal("renderRunRecoveryError() error = nil, want CLI failure")
			}
			failure := decodeCommandJSONStructure(t, stderr.String()).(map[string]any)
			errorValue := failure["error"].(map[string]any)
			if failure["operation"] != "run.recover" || errorValue["code"] != testCase.wantCode {
				t.Fatalf("recovery error = %#v, want operation/code run.recover/%s", failure, testCase.wantCode)
			}
		})
	}
}
