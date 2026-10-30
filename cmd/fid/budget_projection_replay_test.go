package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestBudgetDenialIncludesDerivedOutputState(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	path, _ := gateRunInputForTest(t, workspace, issueID, envelopeID, "state-denial")
	args := []string{"--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", path, "--idempotency-key"}
	for index := 0; index < 3; index++ {
		result := executeCommandForTest(workspace, append(args, fmt.Sprintf("state-charge-%d", index))...)
		if result.err != nil {
			t.Fatalf("charge %d: %#v", index, result)
		}
	}
	denied := executeCommandForTest(workspace, append(args, "state-denial")...)
	if denied.err == nil || denied.stderr != "" {
		t.Fatalf("denial: %#v", denied)
	}
	var terminal gateTerminalFrame
	if err := json.Unmarshal([]byte(strings.TrimSpace(denied.stdout)), &terminal); err != nil {
		t.Fatal(err)
	}
	details := terminal.Data.Response["error"].(map[string]any)["details"].(map[string]any)
	budget := details["budget"].(map[string]any)
	if budget["state"] != "exhausted" || budget["remaining"] != float64(0) || budget["consumed"] != float64(3) {
		t.Fatalf("denial budget = %#v, want exhausted/0 remaining/3 consumed", budget)
	}
	before := snapshotWorkspacePersistence(t, workspace)
	replayed := executeCommandForTest(workspace, append(args, "state-denial")...)
	if replayed.err == nil {
		t.Fatal("denial replay unexpectedly succeeded")
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(replayed.stdout)), &terminal); err != nil {
		t.Fatal(err)
	}
	replayDetails := terminal.Data.Response["error"].(map[string]any)["details"].(map[string]any)
	if !reflect.DeepEqual(replayDetails["budget"], budget) {
		t.Fatalf("denial replay changed budget: %#v", replayDetails)
	}
	assertRunStateUnchanged(t, workspace, before)
}

func TestBudgetProjectionAndHistoricalReplayAcrossCLI(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspaceWithDefinitions(t, gateDefinitionsForTest(func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	}))
	input := envelopeCreateInput(t, workspace, issueID, registration)
	input["budgets"].([]any)[0].(map[string]any)["warning_threshold"] = 1
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	actions := []string{"claim.acquire", "envelope.activate", "gate.run", "run.start"}
	grant["actions"] = actions
	grant["actor_selector"] = map[string]any{"actor_id": "actor-builder", "required_capabilities": actions}
	createPath := filepath.Join(workspace, "projection-create.json")
	writeActivationInput(t, createPath, input)
	createArgs := []string{"envelope", "create", "--actor", "actor-builder", "--input", createPath, "--idempotency-key", "projection-create"}
	created := budgetCLIResult(t, workspace, createArgs...)["envelope"].(map[string]any)
	envelopeID := created["id"].(string)
	assertProjectedBudget(t, created, "available", "0", "3")
	activatePath := filepath.Join(workspace, "projection-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := budgetCLIResult(t, workspace, "envelope", "activate", "--actor", "actor-builder", "--input", activatePath, "--idempotency-key", "projection-activate")
	assertProjectedBudget(t, activated["envelope"].(map[string]any), "available", "0", "3")
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "projection-claim")
	repository := runRepositoryForTest(t, workspace)
	startPath := filepath.Join(workspace, "projection-start.json")
	writeActivationInput(t, startPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	startArgs := []string{"run", "start", "--actor", "actor-builder", "--input", startPath, "--idempotency-key", "projection-start"}
	started := budgetCLIResult(t, workspace, startArgs...)["run"].(map[string]any)
	assertProjectedBudget(t, started, "available", "0", "3")
	runID := started["id"].(string)
	gatePath := filepath.Join(workspace, "projection-gate.json")
	writeActivationInput(t, gatePath, map[string]any{"registration_id": registration.ID, "registration_version": registration.Version, "obligation_id": testEvidenceObligationID, "run_id": runID, "expected_repository": map[string]any{"head_oid": repository["head_oid"], "dirty": repository["dirty"]}})
	for index, state := range []string{"available", "warning", "exhausted"} {
		result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", gatePath, "--idempotency-key", fmt.Sprintf("projection-gate-%d", index))
		if result.err != nil {
			t.Fatalf("charge %d: %#v", index, result)
		}
		for _, entity := range []string{"run", "envelope"} {
			id := runID
			if entity == "envelope" {
				id = envelopeID
			}
			value := budgetCLIResult(t, workspace, entity, "get", id)[entity].(map[string]any)
			assertProjectedBudget(t, value, state, fmt.Sprint(index+1), fmt.Sprint(2-index))
		}
	}
	before := snapshotWorkspacePersistence(t, workspace)
	for _, replay := range []struct {
		args     []string
		entity   string
		original map[string]any
	}{{createArgs, "envelope", created}, {startArgs, "run", started}} {
		got := budgetCLIResult(t, workspace, replay.args...)[replay.entity].(map[string]any)
		if !reflect.DeepEqual(got, replay.original) {
			t.Fatalf("%s replay changed historical result: got %#v want %#v", replay.entity, got, replay.original)
		}
	}
	assertRunStateUnchanged(t, workspace, before)
	assertBudgetPolicyEffectsRemainInputOnly(t, workspace)
}

func TestBudgetUnknownProjectionPreservesNullCounters(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspaceWithDefinitions(t, gateDefinitionsForTest(nil))
	input := envelopeCreateInput(t, workspace, issueID, registration)
	input["budgets"].([]any)[0].(map[string]any)["unit"] = "input-tokens"
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	actions := []string{"claim.acquire", "envelope.activate", "run.start"}
	grant["actions"] = actions
	grant["actor_selector"] = map[string]any{"actor_id": "actor-builder", "required_capabilities": actions}
	path := filepath.Join(workspace, "unknown-budget.json")
	writeActivationInput(t, path, input)
	created := budgetCLIResult(t, workspace, "envelope", "create", "--actor", "actor-builder", "--input", path, "--idempotency-key", "unknown-create")["envelope"].(map[string]any)
	envelopeID := created["id"].(string)
	writeActivationInput(t, path, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := budgetCLIResult(t, workspace, "envelope", "activate", "--actor", "actor-builder", "--input", path, "--idempotency-key", "unknown-activate")["envelope"].(map[string]any)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "unknown-claim")
	writeActivationInput(t, path, runStartInput(issueID, envelopeID, claim["id"], runRepositoryForTest(t, workspace)))
	started := budgetCLIResult(t, workspace, "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "unknown-start")["run"].(map[string]any)
	read := budgetCLIResult(t, workspace, "run", "get", started["id"].(string))["run"].(map[string]any)
	for _, value := range []map[string]any{created, activated, started, read} {
		budget := value["budgets"].([]any)[0].(map[string]any)
		if budget["state"] != "unknown" || budget["consumed"] != nil || budget["remaining"] != nil {
			t.Fatalf("external Budget without observation = %#v, want unknown with null counters", budget)
		}
	}
	assertBudgetPolicyEffectsRemainInputOnly(t, workspace)
}

func budgetCLIResult(t *testing.T, workspace string, args ...string) map[string]any {
	t.Helper()
	result := executeCommandForTest(workspace, append([]string{"--project", workspace, "--json"}, args...)...)
	if result.err != nil {
		t.Fatalf("%v: %#v", args, result)
	}
	return decodeCommandJSONStructure(t, result.stdout).(map[string]any)["result"].(map[string]any)
}

func assertProjectedBudget(t *testing.T, value map[string]any, state, consumed, remaining string) {
	t.Helper()
	budgets := value["budgets"].([]any)
	if len(budgets) != 1 {
		t.Fatalf("budgets = %#v, want one", budgets)
	}
	budget := budgets[0].(map[string]any)
	if budget["state"] != state || budget["consumed"] != json.Number(consumed) || budget["remaining"] != json.Number(remaining) {
		t.Fatalf("budget = %#v, want state=%s consumed=%s remaining=%s", budget, state, consumed, remaining)
	}
}

func assertBudgetPolicyEffectsRemainInputOnly(t *testing.T, workspace string) {
	t.Helper()
	seen := map[string]int{}
	for _, line := range splitJournalLines(mustReadCommandFile(t, mustJournalPath(t, workspace))) {
		var event journal.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatal(err)
		}
		for _, effect := range event.Effects {
			if effect.Table != "envelope_budgets" && effect.Table != "run_budget_bindings" {
				continue
			}
			var budget map[string]any
			if err := json.Unmarshal([]byte(effect.Values["budget_json"].(string)), &budget); err != nil {
				t.Fatal(err)
			}
			if _, present := budget["state"]; present || budget["consumed"] != nil || budget["remaining"] != nil {
				t.Fatalf("%s persisted derived projection: %#v", effect.Table, budget)
			}
			seen[effect.Table]++
		}
	}
	if seen["envelope_budgets"] != 1 || seen["run_budget_bindings"] != 1 {
		t.Fatalf("policy Effects = %#v, want one immutable policy per owner", seen)
	}
}
