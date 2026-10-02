//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func (model *p22MemoryModel) denyGate(t *testing.T, key string) p22ConformanceResponse {
	t.Helper()
	if model.sampleID != "" {
		t.Fatal("semantic fake fixture must supply an unknown budget for gate denial; observed execution is a separate fixture")
	}
	code, budgetName := model.genericBudgetBlock()
	if code == "" {
		t.Fatal("semantic fake gate denial requires an unknown or exhausted generic budget")
	}
	if model.denials == nil {
		model.denials = map[string]map[string]any{}
	}
	response, replayed := model.denials[key]
	if !replayed {
		model.sequence++
		stopID, attentionID := fmt.Sprintf("STP-%d", model.sequence), fmt.Sprintf("ATT-%d", model.sequence)
		event := map[string]any{"id": fmt.Sprintf("EVT-%d", model.sequence), "occurred_at": "2026-09-07T00:00:00Z"}
		subject := map[string]any{"type": "run", "id": model.run["id"]}
		model.stop = map[string]any{"id": stopID, "subject": subject, "code": code}
		model.attention = map[string]any{"id": attentionID, "subject": subject, "code": code, "status": "open", "version": 1, "stop_ids": []any{stopID}}
		response = map[string]any{"contract_version": "fid-cli.v1", "operation": "gate.run", "ok": false, "error": map[string]any{"code": code, "category": "precondition", "retryable": false, "subject": subject, "details": map[string]any{"budget": map[string]any{"name": budgetName, "consumed": nil, "remaining": nil}, "stop_ids": []any{stopID}, "attention_ids": []any{attentionID}, "event": event, "replayed": false}}}
		model.denials[key] = p22CloneMap(response)
	}
	response = p22CloneMap(response)
	response["error"].(map[string]any)["details"].(map[string]any)["replayed"] = replayed
	frame := map[string]any{"frame_schema": "fid.live.v1", "sequence": 1, "stream": "terminal", "data": map[string]any{"attempt_id": nil, "receipt_id": nil, "outcome": "fid-error", "response": response, "replayed": replayed}}
	return p22FixtureResponse(t, frame, true, 4)
}

func (model *p22MemoryModel) genericBudgetBlock() (string, string) {
	for _, budget := range model.budgets {
		if budget.action == "review-rework-exhausted" {
			continue
		}
		if budget.remaining == nil {
			return "budget-data-unavailable", budget.name
		}
		if *budget.remaining == 0 {
			return "budget-exhausted", budget.name
		}
	}
	return "", ""
}

func assertP22LiveDenial(t *testing.T, response p22ConformanceResponse, replayed bool) map[string]any {
	t.Helper()
	if response.ExitCode != 4 || response.Stderr != "" {
		t.Fatalf("committed denial exit=%d stdout=%q stderr=%q", response.ExitCode, response.Stdout, response.Stderr)
	}
	lines := strings.Split(strings.TrimSpace(response.Stdout), "\n")
	if len(lines) != 1 {
		t.Fatalf("denial emitted %d frames, want sole terminal", len(lines))
	}
	var frame map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &frame); err != nil {
		t.Fatal(err)
	}
	data, ok := frame["data"].(map[string]any)
	if !ok || frame["frame_schema"] != "fid.live.v1" || frame["stream"] != "terminal" || data["attempt_id"] != nil || data["receipt_id"] != nil || data["outcome"] != "fid-error" || data["replayed"] != replayed {
		t.Fatalf("pre-start denial frame = %#v; replayed=%v", frame, replayed)
	}
	envelope, ok := data["response"].(map[string]any)
	if !ok || envelope["ok"] != false {
		t.Fatalf("denial implies success: %#v", data)
	}
	failure := envelope["error"].(map[string]any)
	details := failure["details"].(map[string]any)
	if failure["code"] != "budget-data-unavailable" || details["replayed"] != replayed || p22ResultPath(details, "budget.consumed") != nil || p22ResultPath(details, "budget.remaining") != nil || p22ResultPath(details, "stop_ids.#") != float64(1) || p22ResultPath(details, "attention_ids.#") != float64(1) || p22ResultPath(details, "event.id") == nil {
		t.Fatalf("denial durable facts = %#v", failure)
	}
	return envelope
}

func TestP22FakeGenericAdmissionExcludesReviewCycles(t *testing.T) {
	t.Parallel()
	zero := 0
	tests := []struct {
		name    string
		budgets []p22FakeBudget
		code    string
		budget  string
	}{
		{"unknown review only", []p22FakeBudget{{name: "review-cycles", action: "review-rework-exhausted"}}, "", ""},
		{"exhausted review only", []p22FakeBudget{{name: "review-cycles", action: "review-rework-exhausted", remaining: &zero}}, "", ""},
		{"unknown other after review", []p22FakeBudget{{name: "review-cycles", action: "review-rework-exhausted", remaining: &zero}, {name: "tokens", action: "budget-exhausted"}}, "budget-data-unavailable", "tokens"},
		{"exhausted other after review", []p22FakeBudget{{name: "review-cycles", action: "review-rework-exhausted"}, {name: "attempts", action: "budget-exhausted", remaining: &zero}}, "budget-exhausted", "attempts"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			model := p22NewFakeRecoveryAdapter().model
			model.budgets = tt.budgets
			code, budget := model.genericBudgetBlock()
			if code != tt.code || budget != tt.budget {
				t.Fatalf("generic gate block = (%q, %q), want (%q, %q)", code, budget, tt.code, tt.budget)
			}
		})
	}
	t.Run("shared gate denial keeps the other budget", func(t *testing.T) {
		t.Parallel()
		model := p22NewFakeRecoveryAdapter().model
		model.budgets = []p22FakeBudget{{name: "review-cycles", action: "review-rework-exhausted", remaining: &zero}, {name: "tokens", action: "budget-exhausted"}}
		response := model.denyGate(t, "mixed-budget-denial")
		denial := assertP22LiveDenial(t, response, false)
		if name := p22ResultPath(denial, "error.details.budget.name"); name != "tokens" {
			t.Fatalf("gate denied for budget %v, want tokens", name)
		}
	})
}
