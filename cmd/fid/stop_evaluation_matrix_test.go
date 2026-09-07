package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestStopPredicatesAcrossIssueAndEnvelopeCommands(t *testing.T) {
	for _, surface := range []string{"issue", "envelope"} {
		for _, state := range []string{"applicable", "unrelated-subject", "unrelated-action", "resolved", "multiple"} {
			t.Run(surface+"/"+state, func(t *testing.T) {
				workspace, issueID, envelopeID, _, _ := task18TerminalGuardFixture(t, false, true)
				if surface == "envelope" {
					result := stopMatrixTransition(workspace, issueID)
					if result.err != nil {
						t.Fatalf("qualifying Issue transition: %#v", result)
					}
				}
				wantIDs := seedStopEvaluationMatrix(t, workspace, envelopeID, state)
				beforeIssue := executionTransitionIssueSnapshot(t, workspace, issueID)
				beforeEvents := journalEventCount(t, workspace)
				if surface == "issue" {
					assertStopMatrixEvaluator(t, workspace, issueID, wantIDs)
				}
				var result commandExecutionResult
				if surface == "issue" {
					result = stopMatrixTransition(workspace, issueID)
				} else {
					path := filepath.Join(workspace, "stop-matrix-complete.json")
					writeActivationInput(t, path, map[string]any{"envelope_id": envelopeID, "expected_revision": 1, "repository": runRepositoryForTest(t, workspace)})
					result = executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "complete", "--actor", "actor-builder", "--input", path, "--idempotency-key", "stop-matrix-complete")
				}
				assertStopMatrixCommand(t, workspace, surface, issueID, envelopeID, result, wantIDs, beforeEvents)
				if len(wantIDs) > 0 && executionTransitionIssueSnapshot(t, workspace, issueID) != beforeIssue {
					t.Fatal("Stop rejection changed the Issue")
				}
			})
		}
	}
}

func stopMatrixTransition(workspace, issueID string) commandExecutionResult {
	return executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "transition", issueID, "--to", "qa_pass", "--actor", "agent:actor-builder", "--expected-revision", "4", "--idempotency-key", "stop-matrix-transition")
}

// These are admitted projection-state fixtures for guard selection, not evidence
// of Stop admission/resolution. The producer's validation is covered separately.
func seedStopEvaluationMatrix(t *testing.T, workspace, envelopeID, state string) []string {
	t.Helper()
	subjectID := envelopeID
	if state == "unrelated-subject" {
		subjectID = "ENV-unrelated"
	}
	actions := `["envelope.complete","issue.transition-terminal-success"]`
	if state == "unrelated-action" {
		actions = `["gate.run"]`
	}
	condition, err := json.Marshal(map[string]any{"type": "work-ended", "subject": map[string]any{"type": "envelope", "id": subjectID}})
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{"STOP-matrix-a"}
	if state == "multiple" {
		ids = append(ids, "STOP-matrix-b")
	}
	for _, id := range ids {
		task18MutateEvidenceFixture(t, workspace, `INSERT INTO stops (id,subject_type,subject_id,code,blocked_actions_json,details_json,resolution_condition_json,resolver_capability,actor_id,raised_at,created_event_id) VALUES (?,'envelope',?,'operator-stop',?,'{"reason":"guard matrix","reference":"policy:matrix"}',?,'stop.resolve','actor-builder','2026-01-01T00:00:00Z','EVT-matrix-raised')`, id, subjectID, actions, string(condition))
		if state == "resolved" {
			task18MutateEvidenceFixture(t, workspace, `INSERT INTO stop_resolutions (stop_id,actor_id,fact_event_id,resolved_at,created_event_id) VALUES (?,'actor-builder','EVT-matrix-fact','2026-01-01T00:00:01Z','EVT-matrix-resolved')`, id)
		}
	}
	if state == "applicable" || state == "multiple" {
		return ids
	}
	return []string{}
}

func assertStopMatrixEvaluator(t *testing.T, workspace, issueID string, want []string) {
	t.Helper()
	revision := 4
	result, err := workflow.EvaluateExecutionTransition(workspace, issueMigrationsDir(), workflow.ExecutionTransitionRequest{IssueID: issueID, ActorID: "actor-builder", ExpectedIssueRevision: &revision, ActualIssueRevision: revision, WorkflowEdgeValid: true, TerminalSuccess: true})
	if err != nil {
		t.Fatal(err)
	}
	got := []string{}
	for _, reason := range result.Reasons {
		if reason.SubjectType == "stop" {
			got = append(got, reason.SubjectID)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("shared evaluator Stop IDs = %v, want %v; reasons=%#v", got, want, result.Reasons)
	}
}

func assertStopMatrixCommand(t *testing.T, workspace, surface, issueID, envelopeID string, result commandExecutionResult, want []string, beforeEvents int) {
	t.Helper()
	if len(want) == 0 {
		if result.err != nil || result.stderr != "" {
			t.Fatalf("unrelated/resolved Stop blocked %s: %#v", surface, result)
		}
		if surface == "issue" {
			if executionTransitionIssueSnapshot(t, workspace, issueID).State != "qa_pass" {
				t.Fatal("permitted transition did not reach qa_pass")
			}
		} else if decodeCanonicalResult(t, result.stdout)["envelope"].(map[string]any)["status"] != "completed" {
			t.Fatal("permitted completion did not complete Envelope")
		}
		return
	}
	if result.err == nil || result.stdout != "" {
		t.Fatalf("applicable Stop allowed %s: %#v", surface, result)
	}
	if journalEventCount(t, workspace) != beforeEvents {
		t.Fatal("Stop rejection appended an Event")
	}
	failure := decodeCommandJSONStructure(t, result.stderr).(map[string]any)["error"].(map[string]any)
	details := failure["details"].(map[string]any)
	if surface == "issue" {
		if failure["code"] != "transition-not-permitted" {
			t.Fatalf("Issue Stop error = %#v", failure)
		}
		codes := details["unsatisfied"].([]any)
		if len(codes) != len(want) {
			t.Fatalf("Issue guard reasons = %#v, want %d Stop reasons", codes, len(want))
		}
		for _, code := range codes {
			if code != "operator-stop" {
				t.Fatalf("Issue guard reason %v", code)
			}
		}
		return
	}
	if failure["code"] != "obligation-unsatisfied" {
		t.Fatalf("Envelope Stop error = %#v", failure)
	}
	got := []string{}
	for _, value := range details["reasons"].([]any) {
		reason := value.(map[string]any)
		if reason["subject_type"] == "stop" {
			got = append(got, reason["subject_id"].(string))
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Envelope Stop IDs = %v, want %v", got, want)
	}
	read := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "get", envelopeID)
	if read.err != nil || decodeCanonicalResult(t, read.stdout)["envelope"].(map[string]any)["status"] != "active" {
		t.Fatalf("Stop rejection changed Envelope: %#v", read)
	}
}
