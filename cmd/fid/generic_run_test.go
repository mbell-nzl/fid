package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets(t *testing.T) {
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-start-claim")
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	path := filepath.Join(workspace, "run-start.json")
	writeActivationInput(t, path, map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"claim_id": claim["id"], "claim_version": 1, "repository": repository,
	})

	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "run", "start", "--actor", "actor-builder",
		"--input", path, "--idempotency-key", "run-start-1",
	)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("run start = %#v, want stdout-only success", result)
	}
	run := runFromResponse(t, result.stdout)
	for field, want := range map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": "1", "claim_id": claim["id"],
		"claim_version": "1", "actor_id": "actor-builder", "status": "active", "terminal_repository": nil, "ended_at": nil,
	} {
		if got := stringifyClaimValue(run[field]); got != want {
			t.Fatalf("run.%s = %#v, want %#v", field, got, want)
		}
	}
	if got := run["start_repository"].(map[string]any); got["worktree_id"] != repository["worktree_id"] || got["tree_oid"] != repository["tree_oid"] {
		t.Fatalf("start_repository = %#v, want exact input repository %#v", got, repository)
	}
	budgets := run["budgets"].([]any)
	if len(budgets) != 1 || budgets[0].(map[string]any)["consumed"] != nil || budgets[0].(map[string]any)["remaining"] != nil {
		t.Fatalf("run budgets = %#v, want one immutable null-counter budget", budgets)
	}
	got := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "get", run["id"].(string))
	if got.err != nil || runFromResponse(t, got.stdout)["id"] != run["id"] {
		t.Fatalf("run get = %#v, want the started Run", got)
	}
	listed := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "list", "--issue", issueID, "--status", "active")
	if listed.err != nil {
		t.Fatalf("run list: %v; stderr=%q", listed.err, listed.stderr)
	}
	items := decodeCommandJSONStructure(t, listed.stdout).(map[string]any)["result"].(map[string]any)["items"].([]any)
	if len(items) != 1 || items[0].(map[string]any)["id"] != run["id"] {
		t.Fatalf("run list items = %#v, want the active Run", items)
	}
}

func TestRunStartBindingMatrixRejectsWithoutMutation(t *testing.T) {
	t.Run("qualifying", func(t *testing.T) {
		workspace, issueID, envelopeID := runEnvelopeForTest(t)
		claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-binding-qualified-claim")
		repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
		path := filepath.Join(workspace, "run-binding-qualified.json")
		writeActivationInput(t, path, runStartInput(issueID, envelopeID, claim["id"], repository))

		beforeRuns, beforeEvents := runStateSnapshot(t, workspace)
		started := executeCommandForTest(workspace,
			"--project", workspace, "--json", "run", "start", "--actor", "actor-builder",
			"--input", path, "--idempotency-key", "run-binding-qualified",
		)
		if started.err != nil || started.stderr != "" {
			t.Fatalf("qualifying run start = %#v, want stdout-only success", started)
		}
		if run := runFromResponse(t, started.stdout); run["status"] != "active" || run["budgets"] == nil {
			t.Fatalf("qualifying run = %#v, want active Run with bound budgets", run)
		}
		afterRuns, afterEvents := runStateSnapshot(t, workspace)
		if afterRuns == beforeRuns || afterEvents != beforeEvents+1 {
			t.Fatalf("qualifying Run state = (%q, %d), want changed runs and Event count %d", afterRuns, afterEvents, beforeEvents+1)
		}
	})

	tests := []struct {
		name         string
		actor        string
		releaseClaim bool
		mutate       func(map[string]any)
	}{
		{
			name: "inactive Claim", actor: "actor-builder", releaseClaim: true,
			mutate: func(map[string]any) {},
		},
		{
			name: "mismatched Issue", actor: "actor-builder",
			mutate: func(input map[string]any) { input["issue_id"] = "ISSUE-mismatched" },
		},
		{
			name: "mismatched Envelope", actor: "actor-builder",
			mutate: func(input map[string]any) { input["envelope_id"] = "ENVELOPE-mismatched" },
		},
		{
			name: "mismatched Envelope revision", actor: "actor-builder",
			mutate: func(input map[string]any) { input["envelope_revision"] = 2 },
		},
		{
			name: "mismatched Actor", actor: "actor-other",
			mutate: func(map[string]any) {},
		},
		{
			name: "mismatched repository", actor: "actor-builder",
			mutate: func(input map[string]any) { input["repository"].(map[string]any)["tree_oid"] = "tree-mismatched" },
		},
		{
			name: "immutable budgets", actor: "actor-builder",
			mutate: func(input map[string]any) { input["budgets"] = []any{map[string]any{"name": "override"}} },
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			workspace, issueID, envelopeID := runEnvelopeForTest(t)
			claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-binding-"+strings.ReplaceAll(testCase.name, " ", "-"))
			repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
			input := runStartInput(issueID, envelopeID, claim["id"], repository)
			if testCase.releaseClaim {
				releasePath := filepath.Join(workspace, "run-binding-release.json")
				writeActivationInput(t, releasePath, map[string]any{"claim_id": claim["id"], "expected_version": 1, "reason": "completed"})
				released := executeCommandForTest(workspace,
					"--project", workspace, "--json", "claim", "release", "--actor", "actor-builder",
					"--input", releasePath, "--idempotency-key", "run-binding-release",
				)
				if released.err != nil || released.stderr != "" {
					t.Fatalf("release Claim = %#v, want stdout-only success", released)
				}
			}
			testCase.mutate(input)
			path := filepath.Join(workspace, "run-binding-rejected.json")
			writeActivationInput(t, path, input)

			beforeRuns, beforeEvents := runStateSnapshot(t, workspace)
			result := executeCommandForTest(workspace,
				"--project", workspace, "--json", "run", "start", "--actor", testCase.actor,
				"--input", path, "--idempotency-key", "run-binding-"+strings.ReplaceAll(testCase.name, " ", "-"),
			)
			assertRunStartInvalidArgument(t, result)
			assertRunStateUnchanged(t, workspace, beforeRuns, beforeEvents)
		})
	}
}

func TestRunFinishIsOneWayAndRecoveryRespawnAreAbsent(t *testing.T) {
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-finish-claim")
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	startPath := filepath.Join(workspace, "run-finish-start.json")
	writeActivationInput(t, startPath, map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"claim_id": claim["id"], "claim_version": 1, "repository": repository,
	})
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", startPath, "--idempotency-key", "run-finish-start")
	if started.err != nil {
		t.Fatalf("start run: %v; stderr=%q", started.err, started.stderr)
	}
	run := runFromResponse(t, started.stdout)

	finishPath := filepath.Join(workspace, "run-finish.json")
	writeActivationInput(t, finishPath, map[string]any{"run_id": run["id"], "expected_status": "active", "outcome": "succeeded", "repository": repository})
	finished := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "finish", "--actor", "actor-builder", "--input", finishPath, "--idempotency-key", "run-finish-1")
	if finished.err != nil || finished.stderr != "" {
		t.Fatalf("finish run = %#v, want stdout-only success", finished)
	}
	finishedRun := runFromResponse(t, finished.stdout)
	if finishedRun["status"] != "succeeded" || finishedRun["ended_at"] == nil || finishedRun["terminal_repository"] == nil {
		t.Fatalf("finished run = %#v, want one terminal succeeded Run", finishedRun)
	}

	before := journalEventCount(t, workspace)
	repeated := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "finish", "--actor", "actor-builder", "--input", finishPath, "--idempotency-key", "run-finish-repeat")
	if repeated.err == nil || repeated.stdout != "" {
		t.Fatalf("repeat finish = %#v, want closed failure", repeated)
	}
	if after := journalEventCount(t, workspace); after != before {
		t.Fatalf("repeat finish changed Event count to %d, want %d", after, before)
	}

	for _, absent := range []string{"run recover", "run predecessor", "run continuation", "run respawn"} {
		for _, leaf := range liveExecutableLeafCommands(newRootCommandAt(workspace)) {
			if leaf == absent {
				t.Fatalf("unavailable Run leaf %q remains registered", absent)
			}
		}
	}
}

func runEnvelopeForTest(t *testing.T) (string, string, string) {
	t.Helper()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	input := envelopeCreateInput(t, workspace, issueID, registration)
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "claim.release", "envelope.activate", "run.finish", "run.start"}
	grant["actor_selector"] = map[string]any{
		"actor_id":              "actor-builder",
		"required_capabilities": []string{"claim.acquire", "claim.release", "envelope.activate", "run.finish", "run.start"},
	}
	path := filepath.Join(workspace, "run-envelope.json")
	writeActivationInput(t, path, input)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", path, "--idempotency-key", "run-envelope-create")
	if created.err != nil {
		t.Fatalf("create run envelope: %v; stderr=%q", created.err, created.stderr)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
	activatePath := filepath.Join(workspace, "run-envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder", "--input", activatePath, "--idempotency-key", "run-envelope-activate")
	if activated.err != nil {
		t.Fatalf("activate run envelope: %v; stderr=%q", activated.err, activated.stderr)
	}
	return workspace, issueID, envelopeID
}

func runRepositoryForTest(t *testing.T, workspace, issueID, envelopeID string) map[string]any {
	t.Helper()
	output := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "get", envelopeID)
	if output.err != nil {
		t.Fatalf("get envelope: %v; stderr=%q", output.err, output.stderr)
	}
	return decodeCommandJSONStructure(t, output.stdout).(map[string]any)["result"].(map[string]any)["envelope"].(map[string]any)["repository_scope"].(map[string]any)
}

func runFromResponse(t *testing.T, output string) map[string]any {
	t.Helper()
	return decodeCommandJSONStructure(t, output).(map[string]any)["result"].(map[string]any)["run"].(map[string]any)
}

func runStartInput(issueID, envelopeID, claimID any, repository map[string]any) map[string]any {
	return map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"claim_id": claimID, "claim_version": 1, "repository": repository,
	}
}

func runStateSnapshot(t *testing.T, workspace string) (string, int) {
	t.Helper()
	listed := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "list")
	if listed.err != nil || listed.stderr != "" {
		t.Fatalf("run list snapshot = %#v, want stdout-only success", listed)
	}
	return listed.stdout, journalEventCount(t, workspace)
}

func assertRunStateUnchanged(t *testing.T, workspace, wantRuns string, wantEvents int) {
	t.Helper()
	gotRuns, gotEvents := runStateSnapshot(t, workspace)
	if gotRuns != wantRuns || gotEvents != wantEvents {
		t.Fatalf("Run state after rejection = (%q, %d), want (%q, %d)", gotRuns, gotEvents, wantRuns, wantEvents)
	}
}

func assertRunStartInvalidArgument(t *testing.T, result commandExecutionResult) {
	t.Helper()
	if result.err == nil || result.stdout != "" {
		t.Fatalf("run start = %#v, want stderr-only rejection", result)
	}
	failure := decodeCommandJSONStructure(t, result.stderr).(map[string]any)
	assertEnvelopeKeySet(t, []byte(result.stderr), []string{"contract_version", "error", "ok", "operation"})
	assertNestedEnvelopeKeySet(t, []byte(result.stderr), "error", []string{"category", "code", "details", "message", "retryable", "subject"})
	assertNestedEnvelopeKeySet(t, []byte(result.stderr), "error.details", []string{"allowed", "argument", "value"})
	errorValue := failure["error"].(map[string]any)
	if failure["contract_version"] != "fid-cli.v1" || failure["ok"] != false || failure["operation"] != "run.start" ||
		errorValue["category"] != "invalid-input" || errorValue["code"] != "invalid-argument" || errorValue["message"] != "invalid Run request" ||
		errorValue["retryable"] != false || errorValue["subject"] != nil {
		t.Fatalf("run start failure = %#v, want canonical invalid-argument envelope", failure)
	}
	details := errorValue["details"].(map[string]any)
	allowed := details["allowed"].([]any)
	if details["argument"] != "input" || details["value"] != nil || len(allowed) != 1 || allowed[0] != "accepted run.start input" {
		t.Fatalf("run start failure details = %#v, want accepted closed input", details)
	}
}

func TestGenericRunUsesNoHistoricalAgentRunsTable(t *testing.T) {
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-table-claim")
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	path := filepath.Join(workspace, "run-table-start.json")
	writeActivationInput(t, path, map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": claim["id"], "claim_version": 1, "repository": repository})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "run-table-start")
	if result.err != nil {
		t.Fatalf("start Run: %v; stderr=%q", result.err, result.stderr)
	}
	if strings.Contains(string(mustReadCommandFile(t, mustJournalPath(t, workspace))), "agent_runs") {
		t.Fatal("generic Run Event consulted historical agent_runs telemetry")
	}
}
