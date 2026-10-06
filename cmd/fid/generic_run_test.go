package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-start-claim")
	repository := runRepositoryForTest(t, workspace)
	path := filepath.Join(workspace, "run-start.json")
	writeActivationInput(t, path, runStartInput(issueID, envelopeID, claim["id"], repository))

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
		"version": "1", "obligation_id": nil,
	} {
		if got := stringifyClaimValue(run[field]); got != want {
			t.Fatalf("run.%s = %#v, want %#v", field, got, want)
		}
	}
	if executor, ok := run["executor"].(map[string]any); !ok || executor["id"] != "executor-builder" || executor["source"] != "activated-harness" {
		t.Fatalf("run.executor = %#v, want bound activated-harness identity", run["executor"])
	}
	if got := run["start_repository"].(map[string]any); got["worktree_id"] != repository["worktree_id"] || got["tree_oid"] != repository["tree_oid"] {
		t.Fatalf("start_repository = %#v, want exact input repository %#v", got, repository)
	}
	budgets := run["budgets"].([]any)
	if len(budgets) != 1 || stringifyClaimValue(budgets[0].(map[string]any)["consumed"]) != "0" || stringifyClaimValue(budgets[0].(map[string]any)["remaining"]) != "3" {
		t.Fatalf("run budgets = %#v, want one derived attempts budget with zero consumption and three remaining", budgets)
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
	t.Parallel()
	t.Run("qualifying", func(t *testing.T) {
		t.Parallel()
		workspace, issueID, envelopeID := runEnvelopeForTest(t)
		claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-binding-qualified-claim")
		repository := runRepositoryForTest(t, workspace)
		path := filepath.Join(workspace, "run-binding-qualified.json")
		writeActivationInput(t, path, runStartInput(issueID, envelopeID, claim["id"], repository))

		before := snapshotWorkspacePersistence(t, workspace)
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
		after := snapshotWorkspacePersistence(t, workspace)
		if after.equal(before) || after.eventCount() != before.eventCount()+1 {
			t.Fatalf("qualifying Run persistence changed=%t events=%d, want changed state and Event count %d", !after.equal(before), after.eventCount(), before.eventCount()+1)
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
			t.Parallel()
			workspace, issueID, envelopeID := runEnvelopeForTest(t)
			claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-binding-"+strings.ReplaceAll(testCase.name, " ", "-"))
			repository := runRepositoryForTest(t, workspace)
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

			before := snapshotWorkspacePersistence(t, workspace)
			result := executeCommandForTest(workspace,
				"--project", workspace, "--json", "run", "start", "--actor", testCase.actor,
				"--input", path, "--idempotency-key", "run-binding-"+strings.ReplaceAll(testCase.name, " ", "-"),
			)
			assertRunStartInvalidArgument(t, result)
			assertRunStateUnchanged(t, workspace, before)
		})
	}
}

func TestRunObserveRequiresTrustedBoundHarnessAndCurrentVersion(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-observe-claim")
	startPath := filepath.Join(workspace, "run-observe-start.json")
	writeActivationInput(t, startPath, runStartInput(issueID, envelopeID, claim["id"], runRepositoryForTest(t, workspace)))
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", startPath, "--idempotency-key", "run-observe-start")
	if started.err != nil {
		t.Fatalf("start Run: %v; stderr=%q", started.err, started.stderr)
	}
	run := runFromResponse(t, started.stdout)
	observePath := filepath.Join(workspace, "run-observe.json")
	request := map[string]any{"run_id": run["id"], "expected_version": 1, "liveness": map[string]any{"executor_id": "executor-builder", "state": "terminated", "source_ref": "harness://worker/1", "observed_at": time.Now().UTC().Format(time.RFC3339Nano)}}
	writeActivationInput(t, observePath, request)

	before := snapshotWorkspacePersistence(t, workspace)
	builder := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-builder", "--input", observePath, "--idempotency-key", "run-observe-builder")
	if builder.err == nil || !strings.Contains(builder.stderr, `"code":"authority-denied"`) {
		t.Fatalf("Builder observation = %#v, want typed authority rejection", builder)
	}
	assertRunStateUnchanged(t, workspace, before)

	harness := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-harness", "--input", observePath, "--idempotency-key", "run-observe-harness")
	if harness.err != nil || harness.stderr != "" {
		t.Fatalf("harness observation = %#v, want stdout-only success", harness)
	}
	payload := decodeCommandJSONStructure(t, harness.stdout).(map[string]any)
	result := payload["result"].(map[string]any)
	if result["run"].(map[string]any)["version"] != json.Number("2") || result["observation"].(map[string]any)["liveness"].(map[string]any)["state"] != "terminated" {
		t.Fatalf("observation result = %#v, want terminated observation and Run version 2", result)
	}

	request["expected_version"] = 1
	beforeStale := snapshotWorkspacePersistence(t, workspace)
	stale := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-harness", "--input", observePath, "--idempotency-key", "run-observe-stale")
	if stale.err == nil || !strings.Contains(stale.stderr, `"code":"run-version-conflict"`) {
		t.Fatalf("stale observation = %#v, want typed version conflict", stale)
	}
	assertRunStateUnchanged(t, workspace, beforeStale)

	request["expected_version"] = 2
	request["liveness"].(map[string]any)["executor_id"] = "executor-forged"
	writeActivationInput(t, observePath, request)
	beforeForged := snapshotWorkspacePersistence(t, workspace)
	forged := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-harness", "--input", observePath, "--idempotency-key", "run-observe-forged")
	if forged.err == nil || !strings.Contains(forged.stderr, `"code":"invalid-argument"`) {
		t.Fatalf("forged observation = %#v, want typed binding rejection", forged)
	}
	assertRunStateUnchanged(t, workspace, beforeForged)
}

func TestRunStartBindsOptionalObligationAndRejectsReusedExecutor(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := observedEvidenceEnvelopeForTest(t, nil)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-executor-claim-1")
	repository := runRepositoryForTest(t, workspace)
	firstInput := runStartInput(issueID, envelopeID, claim["id"], repository)
	firstInput["obligation_id"] = testEvidenceObligationID
	firstPath := filepath.Join(workspace, "run-executor-first.json")
	writeActivationInput(t, firstPath, firstInput)
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", firstPath, "--idempotency-key", "run-executor-first")
	if started.err != nil {
		t.Fatalf("start bound Run: %v; stderr=%q", started.err, started.stderr)
	}
	first := runFromResponse(t, started.stdout)
	if first["obligation_id"] != testEvidenceObligationID {
		t.Fatalf("Run obligation = %#v, want exact bound obligation", first["obligation_id"])
	}

	closeRunForTest(t, workspace, first["id"].(string), repository, "run-executor-close")
	releasePath := filepath.Join(workspace, "run-executor-release.json")
	writeActivationInput(t, releasePath, map[string]any{"claim_id": claim["id"], "expected_version": 1, "reason": "failed"})
	if released := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", releasePath, "--idempotency-key", "run-executor-release"); released.err != nil {
		t.Fatalf("release first Claim: %#v", released)
	}
	second, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-executor-claim-2")
	secondInput := runStartInput(issueID, envelopeID, second["id"], repository)
	secondPath := filepath.Join(workspace, "run-executor-second.json")
	writeActivationInput(t, secondPath, secondInput)

	before := snapshotWorkspacePersistence(t, workspace)
	duplicate := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", secondPath, "--idempotency-key", "run-executor-duplicate")
	if duplicate.err == nil || !strings.Contains(duplicate.stderr, `"code":"duplicate-entity"`) || !strings.Contains(duplicate.stderr, `"subject":{"id":null,"type":"run"}`) || !strings.Contains(duplicate.stderr, `"values":["executor-builder"]`) {
		t.Fatalf("reused executor = %#v, want exact duplicate-entity Run envelope", duplicate)
	}
	assertRunStateUnchanged(t, workspace, before)

	secondInput["executor"] = map[string]any{"id": "executor-second", "source": "activated-harness"}
	secondInput["obligation_id"] = "OBL-missing"
	writeActivationInput(t, secondPath, secondInput)
	beforeWrongObligation := snapshotWorkspacePersistence(t, workspace)
	wrongObligation := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", secondPath, "--idempotency-key", "run-obligation-missing")
	if wrongObligation.err == nil || !strings.Contains(wrongObligation.stderr, `"code":"invalid-argument"`) {
		t.Fatalf("wrong obligation = %#v, want typed binding rejection", wrongObligation)
	}
	assertRunStateUnchanged(t, workspace, beforeWrongObligation)
}

func TestRetiredRunLeafAndUnsupportedRespawnLeavesAreAbsent(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-finish-claim")
	repository := runRepositoryForTest(t, workspace)
	startPath := filepath.Join(workspace, "run-finish-start.json")
	writeActivationInput(t, startPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", startPath, "--idempotency-key", "run-finish-start")
	if started.err != nil {
		t.Fatalf("start run: %v; stderr=%q", started.err, started.stderr)
	}
	before := journalEventCount(t, workspace)
	retired := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "finish")
	if retired.err == nil || retired.stdout != "" || !strings.Contains(retired.err.Error(), `unknown command "finish"`) {
		t.Fatalf("retired finish = %#v, want absent-leaf rejection", retired)
	}
	if after := journalEventCount(t, workspace); after != before {
		t.Fatalf("retired finish changed Event count to %d, want %d", after, before)
	}

	for _, absent := range []string{"run predecessor", "run continuation", "run respawn"} {
		for _, leaf := range liveExecutableLeafCommands(newRootCommandAt(workspace)) {
			if leaf == absent {
				t.Fatalf("unavailable Run leaf %q remains registered", absent)
			}
		}
	}
}

func runEnvelopeForTest(t *testing.T) (string, string, string) {
	t.Helper()
	return runEnvelopeForTestWithDefinition(t, nil)
}

func runEnvelopeForTestWithDefinition(t *testing.T, modifyDefinition func(map[string]any)) (string, string, string) {
	t.Helper()
	workspace, issueID, envelopeID, registration := buildRunEnvelopeForTest(t, gateDefinitionsForTest(modifyDefinition))
	rememberEnvelopeBindingForTest(t, workspace, envelopeBindingForTest{
		registrationID: registration.ID, registrationVersion: registration.Version, obligationID: testEvidenceObligationID,
	})
	return workspace, issueID, envelopeID
}

func buildRunEnvelopeForTest(t *testing.T, definitions []any) (string, string, string, gateRegistrationView) {
	t.Helper()
	workspace, issueID, registration := finalizedEnvelopeWorkspaceWithDefinitions(t, definitions)
	input := envelopeCreateInput(t, workspace, issueID, registration)
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "claim.expire", "claim.release", "envelope.activate", "gate.run", "run.close", "run.observe", "run.start"}
	grant["actor_selector"] = map[string]any{
		"actor_id":              "actor-builder",
		"required_capabilities": []string{"claim.acquire", "claim.expire", "claim.release", "envelope.activate", "gate.run", "run.close", "run.observe", "run.start"},
	}
	input["authority_grants"] = append(input["authority_grants"].([]any), map[string]any{
		"actor_selector": map[string]any{"actor_id": "actor-harness", "required_capabilities": []string{"run.observe"}},
		"actions":        []string{"run.observe"}, "scope": map[string]any{"issue_id": issueID},
	})
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
	return workspace, issueID, envelopeID, registration
}

func runRepositoryForTest(t *testing.T, workspace string) map[string]any {
	t.Helper()
	if repository, ok := retainedRepositoryInputs.Load(workspace); ok {
		return cloneJSONMap(t, repository.(map[string]any))
	}
	return cloneJSONMap(t, testEnvelopeRepository(t, workspace))
}

func runFromResponse(t *testing.T, output string) map[string]any {
	t.Helper()
	return decodeCommandJSONStructure(t, output).(map[string]any)["result"].(map[string]any)["run"].(map[string]any)
}

func runStartInput(issueID, envelopeID, claimID any, repository map[string]any) map[string]any {
	return map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"claim_id": claimID, "claim_version": 1, "executor": map[string]any{"id": "executor-builder", "source": "activated-harness"}, "obligation_id": nil, "repository": repository,
	}
}

func closeRunForTest(t *testing.T, workspace, runID string, repository map[string]any, key string) {
	t.Helper()
	current := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "get", runID)
	if current.err != nil {
		t.Fatalf("get Run before close: %#v", current)
	}
	currentRun := runFromResponse(t, current.stdout)
	var err error
	version, err := strconv.Atoi(fmt.Sprint(currentRun["version"]))
	if err != nil {
		t.Fatalf("decode Run version: %v", err)
	}
	executorID := currentRun["executor"].(map[string]any)["id"]
	observePath := filepath.Join(workspace, key+"-observe.json")
	writeActivationInput(t, observePath, map[string]any{"run_id": runID, "expected_version": version, "liveness": map[string]any{"executor_id": executorID, "state": "terminated", "source_ref": "process:test", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
	observed := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-harness", "--input", observePath, "--idempotency-key", key+"-observe")
	if observed.err != nil {
		t.Fatalf("observe terminated Run: %#v", observed)
	}
	response := decodeCommandJSONStructure(t, observed.stdout).(map[string]any)["result"].(map[string]any)
	observation := response["observation"].(map[string]any)
	run := response["run"].(map[string]any)
	closePath := filepath.Join(workspace, key+"-close.json")
	writeActivationInput(t, closePath, map[string]any{"run_id": runID, "expected_version": run["version"], "observation_id": observation["id"], "repository": repository})
	closed := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "close", "--actor", "actor-builder", "--input", closePath, "--idempotency-key", key+"-close")
	if closed.err != nil {
		t.Fatalf("close Run: %#v", closed)
	}
}

func assertRunStateUnchanged(t *testing.T, workspace string, want workspacePersistenceSnapshot) {
	t.Helper()
	got := snapshotWorkspacePersistence(t, workspace)
	if !got.equal(want) {
		t.Fatalf("Run rejection changed persisted projection or journal; events=%d, want %d", got.eventCount(), want.eventCount())
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
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "run-table-claim")
	repository := runRepositoryForTest(t, workspace)
	path := filepath.Join(workspace, "run-table-start.json")
	writeActivationInput(t, path, runStartInput(issueID, envelopeID, claim["id"], repository))
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "run-table-start")
	if result.err != nil {
		t.Fatalf("start Run: %v; stderr=%q", result.err, result.stderr)
	}
	if strings.Contains(string(mustReadCommandFile(t, mustJournalPath(t, workspace))), "agent_runs") {
		t.Fatal("generic Run Event consulted historical agent_runs telemetry")
	}
}
