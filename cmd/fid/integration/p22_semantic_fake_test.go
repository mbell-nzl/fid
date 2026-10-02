//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// The fake consumes requests, never expected responses. It owns only memory;
// the shared request builder does not access a fid workspace or production API.
type p22FakeRecoveryAdapter struct {
	context *p22FIDRecoveryAdapter
	model   *p22MemoryModel
}

type p22MemoryModel struct {
	activation                    map[string]any
	denials                       map[string]map[string]any
	run, claim                    map[string]any
	liveness, observation, digest string
	sequence                      int
	replies                       map[string]map[string]any
	requests                      map[string]string
	stop, attention               map[string]any
	consumed                      int
	sampleID, terminalEvent       string
	budgets                       []p22FakeBudget
}

type p22FakeBudget struct {
	name, unit, action string
	maximum, consumed  int
	remaining          *int
}

func p22FakeInt(value int) *int { return &value }

func p22NewFakeRecoveryAdapter() *p22FakeRecoveryAdapter {
	return &p22FakeRecoveryAdapter{
		context: &p22FIDRecoveryAdapter{actor: "actor-recovery", runID: "RUN-prior", runVersion: 1, claimID: "CLM-prior", claimVersion: 1, executorID: "executor-corpus-prior", measurement: map[string]any{"path": "measurement", "digest": "sha256:" + strings.Repeat("0", 64)}},
		model: &p22MemoryModel{
			run:     map[string]any{"id": "RUN-prior", "version": 1, "status": "active", "claim_id": "CLM-prior", "predecessor_run_id": nil},
			claim:   map[string]any{"id": "CLM-prior", "version": 1, "status": "active"},
			budgets: []p22FakeBudget{{name: "tokens", unit: "input-tokens", action: "budget-exhausted", maximum: 10}},
			digest:  "sha256:" + strings.Repeat("1", 64), replies: map[string]map[string]any{}, requests: map[string]string{},
		},
	}
}

func (*p22FakeRecoveryAdapter) Name() string { return "fake" }

func (adapter *p22FakeRecoveryAdapter) Execute(t *testing.T, fixture p22ConformanceCase) p22ConformanceResponse {
	t.Helper()
	args, input := adapter.context.p22Args(t, fixture)
	invalidFlag := false
	for _, arg := range args {
		invalidFlag = invalidFlag || arg == "--unknown-field"
	}
	response := adapter.model.execute(t, fixture.Operation, "p22-corpus-"+fixture.Name, adapter.context.actor, input, invalidFlag)
	if response.ExitCode == 0 {
		adapter.context.p22Advance(t, fixture.Operation, response.Stdout)
	}
	return response
}

func (model *p22MemoryModel) execute(t *testing.T, operation, key, actor string, input map[string]any, invalidFlag bool) p22ConformanceResponse {
	t.Helper()
	if strings.HasPrefix(operation, "activation.") {
		return model.executeActivation(t, operation, key, actor, input, invalidFlag)
	}
	failure := func(code string, exit int) p22ConformanceResponse {
		return p22FixtureResponse(t, map[string]any{"contract_version": "fid-cli.v1", "ok": false, "operation": operation, "error": map[string]any{"code": code}}, false, exit)
	}
	if invalidFlag {
		return failure("invalid-argument", 2)
	}
	if operation == "gate.run" {
		return model.denyGate(t, key)
	}
	allowed := map[string]string{
		"run.observe": "run_id expected_version liveness", "run.inspect": "", "run.checkpoint": "run_id expected_version reconciliation_digest",
		"claim.expire": "claim_id expected_version", "run.recover": "prior_run_id expected_run_version expected_claim_version observation_id reconciliation_digest adopted_paths authority executor ttl_seconds",
		"stop.raise": "subject code blocked_actions details resolution_condition resolver_capability", "stop.list": "", "attention.list": "",
		"stop.resolve": "stop_id fact_event_id", "attention.resolve": "attention_id expected_version fact_event_ids",
		"run.observe-budget": "run_id expected_version budget_name source_sample_id sample_kind previous_sample_id value measurement_ref observed_at",
		"run.close":          "run_id expected_version observation_id repository", "run.abandon": "run_id expected_version expected_claim_version observation_id authority repository",
	}
	keys, known := allowed[operation]
	if !known {
		return failure("invalid-argument", 2)
	}
	for field := range input {
		if !strings.Contains(" "+keys+" ", " "+field+" ") {
			return failure("invalid-argument", 2)
		}
	}
	requestBytes, _ := json.Marshal(input)
	request := actor + string(requestBytes)
	if cached, ok := model.replies[key]; ok {
		if model.requests[key] != request {
			return failure("idempotency-conflict", 4)
		}
		copy := p22CloneMap(cached)
		copy["result"].(map[string]any)["replayed"] = true
		return p22FixtureResponse(t, copy, true, 0)
	}
	read := operation == "run.inspect" || operation == "stop.list" || operation == "attention.list"
	if !read && actor != "actor-recovery" {
		return failure("authority-denied", 5)
	}
	if expected, present := input["expected_version"]; present {
		current := model.run["version"]
		if operation == "claim.expire" {
			current = model.claim["version"]
		}
		if operation == "attention.resolve" {
			current = model.attention["version"]
		}
		if p22Number(expected) != p22Number(current) {
			if operation == "attention.resolve" {
				return failure("attention-version-conflict", 4)
			}
			return failure("run-version-conflict", 4)
		}
	}
	var result map[string]any
	switch operation {
	case "run.observe":
		liveness, ok := input["liveness"].(map[string]any)
		if !ok {
			return failure("invalid-argument", 2)
		}
		model.liveness, _ = liveness["state"].(string)
		model.sequence++
		model.observation = fmt.Sprintf("OBS-%d", model.sequence)
		model.run["version"] = p22Number(model.run["version"]) + 1
		result = map[string]any{"observation": map[string]any{"id": model.observation, "run_id": model.run["id"], "kind": "liveness", "liveness": liveness}, "run": p22CloneMap(model.run)}
	case "run.inspect":
		result = map[string]any{"run": p22CloneMap(model.run), "reconciliation": map[string]any{"digest": model.digest, "paths": []any{}}, "liveness": nil}
	case "run.checkpoint":
		if input["reconciliation_digest"] != model.digest {
			return failure("recovery-state-conflict", 4)
		}
		model.run["version"] = p22Number(model.run["version"]) + 1
		result = map[string]any{"checkpoint": map[string]any{"id": "OBS-checkpoint", "run_id": model.run["id"], "kind": "checkpoint"}, "run": p22CloneMap(model.run)}
	case "claim.expire":
		model.claim["status"] = "expired"
		model.claim["version"] = p22Number(model.claim["version"]) + 1
		result = map[string]any{"claim": p22CloneMap(model.claim)}
	case "run.recover":
		if p22Number(input["expected_run_version"]) != p22Number(model.run["version"]) {
			return failure("run-version-conflict", 4)
		}
		if model.liveness == "live" {
			return failure("executor-still-live", 4)
		}
		if model.liveness != "terminated" {
			return failure("executor-liveness-unknown", 4)
		}
		if input["reconciliation_digest"] != model.digest {
			return failure("recovery-state-conflict", 4)
		}
		prior := model.run["id"]
		model.sequence++
		model.run = map[string]any{"id": fmt.Sprintf("RUN-successor-%d", model.sequence), "version": 1, "status": "active", "predecessor_run_id": prior, "claim_id": fmt.Sprintf("CLM-successor-%d", model.sequence)}
		model.claim = map[string]any{"id": model.run["claim_id"], "version": 1, "status": "active"}
		result = map[string]any{"recovery": map[string]any{"id": fmt.Sprintf("REC-%d", model.sequence), "prior_run_id": prior, "successor_run_id": model.run["id"]}, "run": p22CloneMap(model.run), "claim": p22CloneMap(model.claim)}
		model.liveness, model.observation = "unknown", ""
	case "run.observe-budget":
		if input["budget_name"] != "tokens" || input["measurement_ref"] == nil || p22Number(input["value"]) < 0 {
			return failure("invalid-argument", 2)
		}
		previous := input["previous_sample_id"]
		if (model.sampleID == "" && previous != nil) || (model.sampleID != "" && previous != model.sampleID) {
			return failure("budget-sample-conflict", 4)
		}
		value := p22Number(input["value"])
		if input["sample_kind"] == "baseline" {
			model.consumed = value
		} else {
			model.consumed += value
		}
		model.sequence++
		model.sampleID = fmt.Sprintf("BSM-%d", model.sequence)
		model.run["version"] = p22Number(model.run["version"]) + 1
		remaining := max(0, 10-model.consumed)
		model.budgets[0].consumed, model.budgets[0].remaining = model.consumed, p22FakeInt(remaining)
		state := "available"
		if remaining == 0 {
			state = "exhausted"
		}
		result = map[string]any{"sample": map[string]any{"id": model.sampleID, "run_id": model.run["id"], "source": "activated-harness", "source_sample_id": input["source_sample_id"], "sample_kind": input["sample_kind"], "previous_sample_id": previous, "value": value}, "budget": map[string]any{"name": "tokens", "consumed": model.consumed, "remaining": remaining, "state": state}, "stops": []any{}, "attention": []any{}}
	case "stop.raise":
		model.sequence++
		model.stop = map[string]any{"id": fmt.Sprintf("STP-%d", model.sequence), "subject": input["subject"], "code": input["code"], "details": input["details"], "blocked_actions": input["blocked_actions"], "resolution_condition": input["resolution_condition"], "resolver_capability": input["resolver_capability"], "actor_id": actor}
		if model.attention != nil && model.attention["status"] == "open" && reflect.DeepEqual(model.attention["subject"], input["subject"]) && model.attention["code"] == input["code"] {
			model.attention["stop_ids"] = append(model.attention["stop_ids"].([]any), model.stop["id"])
		} else {
			model.attention = map[string]any{"id": fmt.Sprintf("ATT-%d", model.sequence), "subject": input["subject"], "code": input["code"], "status": "open", "version": 1, "stop_ids": []any{model.stop["id"]}}
		}
		result = map[string]any{"stop": p22CloneMap(model.stop), "attention": p22CloneMap(model.attention)}
	case "stop.list", "attention.list":
		items := []any{}
		item := model.stop
		if operation == "attention.list" {
			item = model.attention
		}
		if item != nil {
			items = append(items, p22CloneMap(item))
		}
		result = map[string]any{"items": items}
	case "run.close", "run.abandon":
		if model.liveness != "terminated" || input["observation_id"] != model.observation {
			return failure("executor-liveness-unknown", 4)
		}
		if model.run["status"] != "active" {
			return failure("run-inactive", 4)
		}
		status := "failed"
		if operation == "run.abandon" {
			status = "abandoned"
			model.claim["status"] = "released"
			model.claim["version"] = p22Number(model.claim["version"]) + 1
		}
		model.run["status"], model.run["version"] = status, p22Number(model.run["version"])+1
		result = map[string]any{"run": p22CloneMap(model.run)}
		if operation == "run.abandon" {
			result["claim"] = p22CloneMap(model.claim)
		}
	case "stop.resolve":
		if model.stop == nil || input["stop_id"] != model.stop["id"] || input["fact_event_id"] != model.terminalEvent {
			return failure("resolution-invalid", 4)
		}
		model.stop["resolution"] = map[string]any{"actor_id": actor, "fact_event_id": model.terminalEvent, "event_id": fmt.Sprintf("EVT-%d", model.sequence+1)}
		result = map[string]any{"stop": p22CloneMap(model.stop)}
	case "attention.resolve":
		resolution, ok := model.stop["resolution"].(map[string]any)
		if !ok || input["attention_id"] != model.attention["id"] || !reflect.DeepEqual(input["fact_event_ids"], []string{fmt.Sprint(resolution["event_id"])}) {
			return failure("resolution-invalid", 4)
		}
		model.attention["status"], model.attention["version"] = "resolved", p22Number(model.attention["version"])+1
		result = map[string]any{"attention": p22CloneMap(model.attention)}
	}
	var event any
	if !read {
		model.sequence++
		event = map[string]any{"id": fmt.Sprintf("EVT-%d", model.sequence), "occurred_at": "2026-09-07T00:00:00Z"}
		result["idempotency_key"], result["replayed"] = key, false
		if operation == "run.close" || operation == "run.abandon" {
			model.terminalEvent = event.(map[string]any)["id"].(string)
		}
	}
	envelope := map[string]any{"contract_version": "fid-cli.v1", "ok": true, "operation": operation, "result": result, "event": event}
	if event != nil {
		model.replies[key], model.requests[key] = p22CloneMap(envelope), request
	}
	return p22FixtureResponse(t, envelope, true, 0)
}

func p22CloneMap(value map[string]any) map[string]any {
	encoded, _ := json.Marshal(value)
	var clone map[string]any
	_ = json.Unmarshal(encoded, &clone)
	return clone
}

func p22Number(value any) int {
	switch number := value.(type) {
	case int:
		return number
	case float64:
		return int(number)
	}
	return -1
}

func p22SameResult(left, right map[string]any) bool {
	return reflect.DeepEqual(p22CloneMap(left), p22CloneMap(right))
}

func (model *p22MemoryModel) chargeReviewCycle(envelopeID, reviewID string) {
	for i := range model.budgets {
		budget := &model.budgets[i]
		if budget.action != "review-rework-exhausted" {
			continue
		}
		previous := budget.maximum - budget.consumed
		budget.consumed++
		budget.remaining = p22FakeInt(max(0, budget.maximum-budget.consumed))
		if previous <= 0 || *budget.remaining != 0 {
			return
		}
		model.sequence++
		subject := map[string]any{"type": "envelope", "id": envelopeID}
		model.stop = map[string]any{
			"id": fmt.Sprintf("STP-%d", model.sequence), "subject": subject, "code": "review-rework-exhausted",
			"details":         map[string]any{"envelope_id": envelopeID, "budget_name": budget.name, "unit": budget.unit, "maximum": budget.maximum, "consumed": budget.consumed, "remaining": *budget.remaining, "last_review_id": reviewID},
			"blocked_actions": []string{"claim.acquire", "run.start"}, "resolver_capability": "stop.resolve",
			"resolution_condition": map[string]any{"type": "human-outcome-recorded", "subject": subject, "budget_name": nil, "capability_id": nil, "outcome_class": "review-rework-extension"},
		}
		model.attention = map[string]any{"id": fmt.Sprintf("ATT-%d", model.sequence), "subject": subject, "code": "review-rework-exhausted", "status": "open", "version": 1, "stop_ids": []any{model.stop["id"]}}
		return
	}
}

func TestP22FakeReviewReworkStopCondition(t *testing.T) {
	t.Parallel()
	model := p22NewFakeRecoveryAdapter().model
	model.budgets = []p22FakeBudget{{name: "review-cycles", unit: "review-cycles", action: "review-rework-exhausted", maximum: 1, remaining: p22FakeInt(1)}}
	model.chargeReviewCycle("ENV-review", "REV-first")
	want := map[string]any{
		"subject":             map[string]any{"type": "envelope", "id": "ENV-review"},
		"code":                "review-rework-exhausted",
		"blocked_actions":     []any{"claim.acquire", "run.start"},
		"resolver_capability": "stop.resolve",
		"details":             map[string]any{"envelope_id": "ENV-review", "budget_name": "review-cycles", "unit": "review-cycles", "maximum": 1, "consumed": 1, "remaining": 0, "last_review_id": "REV-first"},
		"resolution_condition": map[string]any{
			"type": "human-outcome-recorded", "subject": map[string]any{"type": "envelope", "id": "ENV-review"},
			"budget_name": nil, "capability_id": nil, "outcome_class": "review-rework-extension",
		},
	}
	for field, expected := range want {
		gotJSON, _ := json.Marshal(model.stop[field])
		wantJSON, _ := json.Marshal(expected)
		if string(gotJSON) != string(wantJSON) {
			t.Fatalf("fake rework Stop %s = %#v, want %#v", field, model.stop[field], expected)
		}
	}
	firstID := model.stop["id"]
	model.chargeReviewCycle("ENV-review", "REV-second")
	if model.stop["id"] != firstID || model.budgets[0].consumed != 2 || *model.budgets[0].remaining != 0 {
		t.Fatalf("fake rework charge retripped Stop or changed remaining: stop=%#v budget=%#v", model.stop, model.budgets[0])
	}
}
