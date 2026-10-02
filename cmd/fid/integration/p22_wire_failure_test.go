//go:build integration

package integration

import (
	"encoding/json"
	"testing"
)

func TestP22AmendedFailureWireCorpus(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"executor-liveness-unknown", "executor-still-live", "recovery-state-conflict", "run-version-conflict", "attention-version-conflict", "recovery-already-linked", "stop-not-found", "attention-not-found", "observation-not-found", "stop-unresolved", "resolution-invalid", "budget-sample-conflict", "duplicate-entity", "recover-null-observation", "abandon-null-observation"} {
		t.Run(name, func(t *testing.T) {
			a := p22NewFIDRecoveryAdapterWithLease(t, "agent", 3600)
			fixture := p22ConformanceCase{Name: "wire-" + name, Operation: "run.recover", Kind: "rejected"}
			code := name
			switch name {
			case "executor-liveness-unknown", "executor-still-live":
				state := "unknown"
				if name == "executor-still-live" {
					state = "live"
				}
				a.p22Observe(t, state)
				a.p22Expire(t)
			case "recovery-state-conflict":
				a.p22Observe(t, "terminated")
				fixture.Operation = "run.checkpoint"
				fixture.Input = map[string]any{"reconciliation_digest": "sha256:" + p22Zeros(64)}
			case "run-version-conflict":
				fixture.Operation = "run.observe"
				fixture.Input = map[string]any{"expected_version": 999}
			case "attention-version-conflict", "resolution-invalid":
				p22AcceptedCall(t, a, "wire-stop", "stop.raise", nil)
				if name == "attention-version-conflict" {
					fixture.Operation = "attention.resolve"
					fixture.Input = map[string]any{"expected_version": 999, "fact_event_ids": []string{"EVT-missing"}}
				} else {
					fixture.Operation = "stop.resolve"
					fixture.Input = map[string]any{"fact_event_id": "EVT-missing"}
				}
			case "recovery-already-linked":
				a.p22Observe(t, "terminated")
				a.p22Expire(t)
				prior := a.runID
				p22AcceptedCall(t, a, "wire-first-recover", "run.recover", nil)
				read := p21Run(t, a.binary, a.workspace, "--json", "run", "get", prior)
				fixture.Input = map[string]any{"prior_run_id": prior, "expected_run_version": p21Path(t, read, "result", "run", "version")}
			case "stop-not-found":
				fixture.Operation = "stop.resolve"
				fixture.Input = map[string]any{"stop_id": "STP-missing", "fact_event_id": "EVT-missing"}
			case "attention-not-found":
				fixture.Operation = "attention.resolve"
				fixture.Input = map[string]any{"attention_id": "ATT-missing", "expected_version": 1, "fact_event_ids": []string{"EVT-missing"}}
			case "observation-not-found", "duplicate-entity":
				a.p22Observe(t, "terminated")
				a.p22Expire(t)
				if name == "observation-not-found" {
					fixture.Input = map[string]any{"observation_id": "OBS-missing"}
				} else {
					fixture.Input = map[string]any{"executor": map[string]any{"id": a.executorID, "source": "activated-harness"}}
				}
			case "stop-unresolved":
				a.p22Observe(t, "terminated")
				a.p22Expire(t)
				p22AcceptedCall(t, a, "wire-blocked-recover", "stop.raise", map[string]any{"blocked_actions": []string{"run.recover"}})
			case "budget-sample-conflict":
				p22AcceptedCall(t, a, "wire-baseline", "run.observe-budget", nil)
				fixture.Operation = "run.observe-budget"
				fixture.Input = map[string]any{"source_sample_id": "wire-baseline", "value": 4}
			case "recover-null-observation", "abandon-null-observation":
				a.p22Observe(t, "terminated")
				a.p22Expire(t)
				code = "invalid-argument"
				fixture.Input = map[string]any{"observation_id": nil}
				if name == "abandon-null-observation" {
					fixture.Operation = "run.abandon"
				}
			}
			before := p21JournalEventCount(t, a.workspace)
			response := a.Execute(t, fixture)
			p22AssertWireFailure(t, fixture.Operation, code, response)
			if p21JournalEventCount(t, a.workspace) != before {
				t.Fatal("wire rejection appended an Event")
			}
		})
	}
	t.Run("competing-prepare", func(t *testing.T) {
		a := p22NewFIDActivationAdapter(t, "prepared").(*p22FIDRecoveryAdapter)
		input := p21Input(t, a.workspace, "wire-competing", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "competing", "phase": 21, "expected_prior_activation_id": nil, "expected_prior_version": nil, "resources": p21ActivationResources()})
		before := p21JournalEventCount(t, a.workspace)
		response := a.p22Run(t, "activation", "prepare", "--input", input, "--idempotency-key", "wire-competing")
		p22AssertWireFailure(t, "activation.prepare", "activation-conflict", response)
		if p21JournalEventCount(t, a.workspace) != before {
			t.Fatal("competing prepare appended an Event")
		}
	})
}

func p22AssertWireFailure(t *testing.T, operation, code string, response p22ConformanceResponse) {
	t.Helper()
	want := map[string]struct {
		category, fields string
		exit             int
	}{
		"executor-liveness-unknown":  {"precondition", "run_id executor_id observation_id", 4},
		"executor-still-live":        {"precondition", "run_id executor_id observation_id", 4},
		"recovery-state-conflict":    {"conflict", "run_id observation_digest paths reasons", 4},
		"run-version-conflict":       {"conflict", "run_id expected_version actual_version", 4},
		"attention-version-conflict": {"conflict", "attention_id expected_version actual_version", 4},
		"recovery-already-linked":    {"conflict", "run_id successor_run_id", 4},
		"stop-not-found":             {"not-found", "id", 3}, "attention-not-found": {"not-found", "id", 3}, "observation-not-found": {"not-found", "id", 3},
		"stop-unresolved":        {"precondition", "stop_ids reasons", 4},
		"resolution-invalid":     {"precondition", "subject fact_event_ids reasons", 4},
		"budget-sample-conflict": {"conflict", "budget_name source_sample_id expected_previous_sample_id actual_previous_sample_id", 4},
		"duplicate-entity":       {"conflict", "entity_type identity", 4},
		"activation-conflict":    {"conflict", "active_id nonterminal_id requested_direction", 4},
		"invalid-argument":       {"invalid-input", "argument value allowed", 2},
	}[code]
	assertP22ProcessFailure(t, operation, code, want.exit, response)
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatal(err)
	}
	p22AssertKeys(t, envelope, "contract_version error ok operation")
	failure := envelope["error"].(map[string]any)
	p22AssertKeys(t, failure, "category code details message retryable subject")
	if failure["category"] != want.category || failure["retryable"] != false || failure["message"] == "" {
		t.Fatalf("closed failure metadata: %#v", failure)
	}
	p22AssertKeys(t, failure["details"].(map[string]any), want.fields)
	if code == "duplicate-entity" {
		p22AssertKeys(t, p22ResultPath(failure, "details.identity").(map[string]any), "fields values")
		if p22ResultPath(failure, "details.entity_type") != "run" || p22ResultPath(failure, "details.identity.fields.0") != "executor_id" || p22ResultPath(failure, "subject.type") != "run" || p22ResultPath(failure, "subject.id") == nil {
			t.Fatalf("duplicate executor identity/subject: %#v", failure)
		}
	}
}
