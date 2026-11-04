package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// TestStopPredicatesAcrossRunStartAndGateRunCommands exercises the two remaining
// guarded command operations. Stop admission and resolution validation have
// direct producer coverage; these rows seed admitted projection facts to isolate
// the shared predicate at each operation boundary.
func TestStopPredicatesAcrossRunStartAndGateRunCommands(t *testing.T) {
	t.Parallel()
	for _, surface := range []string{"run-start", "gate-run"} {
		for _, state := range []string{"applicable", "unrelated-subject", "unrelated-action", "resolved", "multiple"} {
			t.Run(surface+"/"+state, func(t *testing.T) {
				t.Parallel()
				workspace, issueID, envelopeID, marker := stopOperationEnvelopeForTest(t)
				repository := runRepositoryForTest(t, workspace)
				claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "stop-operation-"+surface+"-"+state+"-claim")

				var runID string
				if surface == "gate-run" {
					runID = startStopMatrixRun(t, workspace, issueID, envelopeID, claim["id"], repository, "stop-operation-"+state)
				}
				wantIDs := seedStopOperationMatrix(t, workspace, surface, envelopeID, runID, state)
				before := snapshotWorkspacePersistence(t, workspace)
				beforeRuns, beforeAttempts := stopOperationCounts(t, workspace)

				var result commandExecutionResult
				if surface == "run-start" {
					inputPath := filepath.Join(workspace, "stop-operation-run-start.json")
					writeActivationInput(t, inputPath, runStartInput(issueID, envelopeID, claim["id"], repository))
					result = executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "stop-operation-run-start-"+state)
				} else {
					inputPath := filepath.Join(workspace, "stop-operation-gate-run.json")
					binding := envelopeBindingForWorkspace(t, workspace)
					writeActivationInput(t, inputPath, map[string]any{"registration_id": binding.registrationID, "registration_version": binding.registrationVersion, "obligation_id": binding.obligationID, "run_id": runID, "expected_repository": map[string]any{"head_oid": repository["head_oid"], "dirty": repository["dirty"]}})
					result = executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "stop-operation-gate-run-"+state)
				}

				assertStopOperationResult(t, workspace, surface, result, wantIDs, before, beforeRuns, beforeAttempts, marker)
			})
		}
	}
}

func startStopMatrixRun(t *testing.T, workspace, issueID, envelopeID string, claimID any, repository map[string]any, key string) string {
	t.Helper()
	inputPath := filepath.Join(workspace, "stop-operation-active-run.json")
	writeActivationInput(t, inputPath, runStartInput(issueID, envelopeID, claimID, repository))
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", key+"-start")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("prepare active Run = %#v, want success", result)
	}
	return runFromResponse(t, result.stdout)["id"].(string)
}

func seedStopOperationMatrix(t *testing.T, workspace, surface, envelopeID, runID, state string) []string {
	t.Helper()
	subjectType, subjectID, action := "envelope", envelopeID, "run.start"
	if surface == "gate-run" {
		subjectType, subjectID, action = "run", runID, "gate.run"
	}
	if state == "unrelated-subject" {
		subjectID = "unrelated-" + subjectType
	}
	if state == "unrelated-action" {
		if action == "run.start" {
			action = "gate.run"
		} else {
			action = "run.start"
		}
	}
	condition, err := json.Marshal(map[string]any{"type": "work-ended", "subject": map[string]any{"type": subjectType, "id": subjectID}})
	if err != nil {
		t.Fatalf("marshal Stop condition: %v", err)
	}
	ids := []string{"STOP-operation-a"}
	if state == "multiple" {
		ids = append(ids, "STOP-operation-b")
	}
	for _, id := range ids {
		task18MutateEvidenceFixture(t, workspace, `INSERT INTO stops (id,subject_type,subject_id,code,blocked_actions_json,details_json,resolution_condition_json,resolver_capability,actor_id,raised_at,created_event_id) VALUES (?,?,?,?,?,'{"reason":"operation guard matrix","reference":"policy:matrix"}',?,'stop.resolve','actor-builder','2026-01-01T00:00:00Z','EVT-operation-raised')`, id, subjectType, subjectID, "operator-stop", "[\""+action+"\"]", string(condition))
		if state == "resolved" {
			task18MutateEvidenceFixture(t, workspace, `INSERT INTO stop_resolutions (stop_id,actor_id,fact_event_id,resolved_at,created_event_id) VALUES (?,'actor-builder','EVT-operation-fact','2026-01-01T00:00:01Z','EVT-operation-resolved')`, id)
		}
	}
	if state == "applicable" || state == "multiple" {
		return ids
	}
	return nil
}

func assertStopOperationResult(t *testing.T, workspace, surface string, result commandExecutionResult, wantIDs []string, before workspacePersistenceSnapshot, beforeRuns, beforeAttempts int, marker string) {
	t.Helper()
	afterRuns, afterAttempts := stopOperationCounts(t, workspace)
	if len(wantIDs) == 0 {
		if result.err != nil || result.stderr != "" {
			t.Fatalf("unrelated or resolved Stop blocked %s: %#v", surface, result)
		}
		if surface == "run-start" {
			if afterRuns != beforeRuns+1 || afterAttempts != beforeAttempts {
				t.Fatalf("permitted run.start counts = runs/attempts %d/%d, want %d/%d", afterRuns, afterAttempts, beforeRuns+1, beforeAttempts)
			}
			return
		}
		if afterRuns != beforeRuns || afterAttempts != beforeAttempts+1 {
			t.Fatalf("permitted gate.run counts = runs/attempts %d/%d, want %d/%d", afterRuns, afterAttempts, beforeRuns, beforeAttempts+1)
		}
		if _, err := os.Stat(marker); err != nil {
			t.Fatalf("permitted gate.run did not spawn marker process: %v", err)
		}
		return
	}

	if result.err == nil || result.stdout != "" {
		t.Fatalf("applicable Stop allowed %s: %#v", surface, result)
	}
	assertRunStateUnchanged(t, workspace, before)
	if afterRuns != beforeRuns || afterAttempts != beforeAttempts {
		t.Fatalf("Stop-blocked %s counts = runs/attempts %d/%d, want %d/%d", surface, afterRuns, afterAttempts, beforeRuns, beforeAttempts)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("Stop-blocked gate spawn marker stat = %v, want no spawned process", err)
	}
	failure := decodeCommandJSONStructure(t, result.stderr).(map[string]any)["error"].(map[string]any)
	details := failure["details"].(map[string]any)
	ids := make([]string, 0, len(details["stop_ids"].([]any)))
	for _, item := range details["stop_ids"].([]any) {
		ids = append(ids, item.(string))
	}
	if failure["code"] != "stop-unresolved" || !reflect.DeepEqual(ids, wantIDs) {
		t.Fatalf("Stop-blocked %s failure = %#v, want IDs %v", surface, failure, wantIDs)
	}
	reasons := details["reasons"].([]any)
	if len(reasons) != len(wantIDs) {
		t.Fatalf("Stop-blocked %s reasons = %#v, want %d", surface, reasons, len(wantIDs))
	}
	for index, item := range reasons {
		reason := item.(map[string]any)
		if reason["code"] != "operator-stop" || reason["detail"] != map[string]string{"run-start": "run.start", "gate-run": "gate.run"}[surface] || reason["subject_id"] != wantIDs[index] || reason["subject_type"] != "stop" {
			t.Fatalf("Stop-blocked %s reason = %#v, want exact Stop reason", surface, reason)
		}
	}
}

func stopOperationCounts(t *testing.T, workspace string) (int, int) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open Stop operation projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var runs, attempts int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM runs`).Scan(&runs); err != nil {
		t.Fatalf("count Runs: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts`).Scan(&attempts); err != nil {
		t.Fatalf("count Attempts: %v", err)
	}
	return runs, attempts
}

func stopOperationEnvelopeForTest(t *testing.T) (string, string, string, string) {
	t.Helper()
	// The definition stays identical across cases; executor cwd is each clone's
	// root. An absolute temporary path would defeat the definition-keyed cache.
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/usr/bin/touch", "gate-spawned"}
	})
	return workspace, issueID, envelopeID, filepath.Join(workspace, "gate-spawned")
}

func TestStopOperationFixtureExecutesInEachClonedWorkspace(t *testing.T) {
	t.Parallel()
	first, issueID, envelopeID, _ := stopOperationEnvelopeForTest(t)
	second, secondIssue, secondEnvelope, _ := stopOperationEnvelopeForTest(t)
	for index, row := range []struct{ root, issue, envelope string }{{first, issueID, envelopeID}, {second, secondIssue, secondEnvelope}} {
		marker := filepath.Join(row.root, "gate-spawned")
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("clone starts with another execution's marker: %v", err)
		}
		repository := runRepositoryForTest(t, row.root)
		claim, _ := acquireGenericClaimForTest(t, row.root, row.issue, row.envelope, "fixture-claim")
		runID := startStopMatrixRun(t, row.root, row.issue, row.envelope, claim["id"], repository, "fixture")
		input := filepath.Join(row.root, "fixture-gate.json")
		binding := envelopeBindingForWorkspace(t, row.root)
		writeActivationInput(t, input, map[string]any{"registration_id": binding.registrationID, "registration_version": binding.registrationVersion, "obligation_id": binding.obligationID, "run_id": runID, "expected_repository": map[string]any{"head_oid": repository["head_oid"], "dirty": repository["dirty"]}})
		result := executeCommandForTest(row.root, "--project", row.root, "--json", "gate", "run", "--actor", "actor-builder", "--input", input, "--idempotency-key", "fixture-gate")
		if result.err != nil {
			t.Fatalf("execute clone %d: %v: %s", index, result.err, result.stderr)
		}
		if _, err := os.Stat(marker); err != nil {
			t.Fatalf("execution did not create marker in its own clone: %v", err)
		}
	}
	third, _, _, _ := stopOperationEnvelopeForTest(t)
	if _, err := os.Stat(filepath.Join(third, "gate-spawned")); !os.IsNotExist(err) {
		t.Fatalf("execution contaminated reusable template: %v", err)
	}
}
