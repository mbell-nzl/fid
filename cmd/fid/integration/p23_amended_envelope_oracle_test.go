//go:build integration

package integration

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

// Phase 21 Envelope grammar plus Phase 23's protected_decisions and derived
// Obligation status. All caller fields are taken from the fixture request.
func p23EnvelopeCreateModel(t *testing.T, pair p23HumanLeafPair) map[string]map[string]any {
	t.Helper()
	input := p23LeafRequest(t, pair)
	repository := p24ModelBindingRepository(t, pair.CWD)
	input["repository_scope"] = repository
	resolve := func(source map[string]any) {
		source["repository"] = repository
		if source["source_type"] == "project-file" {
			data, err := os.ReadFile(filepath.Join(pair.CWD, source["path"].(string)))
			if err != nil {
				t.Fatal(err)
			}
			source["digest"] = fmt.Sprintf("sha256:%x", sha256.Sum256(data))
		}
	}
	resolve(input["risk_result"].(map[string]any)["source"].(map[string]any))
	for _, source := range input["sources"].([]any) {
		resolve(source.(map[string]any))
	}
	for _, obligation := range input["obligations"].([]any) {
		resolve(obligation.(map[string]any)["source"].(map[string]any))
	}
	obligations := make([]any, 0, len(input["obligations"].([]any)))
	for _, item := range input["obligations"].([]any) {
		value := make(map[string]any)
		for key, field := range item.(map[string]any) {
			value[key] = field
		}
		value["derived_status"] = "unsatisfied"
		value["reasons"] = []any{}
		obligations = append(obligations, value)
	}
	budgets := make([]any, 0, len(input["budgets"].([]any)))
	for _, item := range input["budgets"].([]any) {
		value := make(map[string]any)
		for key, field := range item.(map[string]any) {
			value[key] = field
		}
		if value["unit"] == "attempts" {
			value["consumed"] = float64(0)
			value["remaining"] = value["maximum"]
			value["state"] = "available"
		} else {
			value["state"] = "unknown"
		}
		budgets = append(budgets, value)
	}
	envelope := map[string]any{
		"id": "<generated:ENV>", "current_revision": float64(1), "status": "draft",
		"outcome": input["outcome"], "non_goals": input["non_goals"], "risk_result": input["risk_result"],
		"repository_scope": input["repository_scope"], "sources": input["sources"],
		"primary_issue_id": input["primary_issue_id"], "related_issue_ids": input["related_issue_ids"],
		"authority_grants": input["authority_grants"], "obligations": obligations, "budgets": budgets,
		"protected_decisions": []any{}, "created_actor_id": p23LeafFlag(t, pair, "--actor"),
		"approved_actor_id": nil, "created_at": "<generated:timestamp>", "activated_at": nil, "completed_at": nil,
	}
	return map[string]map[string]any{"envelope.create": {"envelope": envelope, "idempotency_key": p23LeafFlag(t, pair, "--idempotency-key"), "replayed": true}}
}

func p24ModelBindingRepository(t *testing.T, workspace string) map[string]any {
	t.Helper()
	projectConfig, err := config.ReadProjectConfig(filepath.Join(workspace, ".fid", "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	repository := p22CorpusRepository(t, workspace, projectConfig.ProjectID)
	if value, ok := repository["dirty_digest"].(*string); ok {
		if value == nil {
			repository["dirty_digest"] = nil
		} else {
			repository["dirty_digest"] = *value
		}
	}
	return repository
}

func p23EnvelopeActivateModel(t *testing.T, pair p23HumanLeafPair) map[string]map[string]any {
	t.Helper()
	// The creation request is a fixture input in this workspace. Reading it is
	// independent of both the creation and activation responses.
	creation := p23HumanLeafPair{Operation: "envelope.create", CWD: pair.CWD, JSONArgs: []string{
		"--input", filepath.Join(pair.CWD, "p22-corpus-envelope.json"),
		"--actor", p23LeafFlag(t, pair, "--actor"), "--idempotency-key", "p22-corpus-envelope-create",
	}}
	envelope := p23EnvelopeCreateModel(t, creation)["envelope.create"]["envelope"].(map[string]any)
	envelope["status"] = "active"
	envelope["activated_at"] = "<generated:timestamp>"
	obligationID := envelope["obligations"].([]any)[0].(map[string]any)["id"]
	evaluation := map[string]any{
		"obligation_id": obligationID, "status": "unsatisfied", "validity": "invalid",
		"freshness": "stale", "artifact_availability": "missing",
		"reasons": []any{map[string]any{
			"code": "obligation-evidence-absent", "subject_type": "obligation", "subject_id": obligationID,
			"detail": "qualifying Evidence is absent",
		}},
	}
	return map[string]map[string]any{"envelope.activate": {
		"envelope": envelope, "evaluation": []any{evaluation},
		"idempotency_key": p23LeafFlag(t, pair, "--idempotency-key"), "replayed": true,
	}}
}

func p23CheckAmendedIdentity(t *testing.T, accepted map[string]p23LeafOutcome) {
	t.Helper()
	for _, operation := range []string{"envelope.create", "envelope.activate", "claim.acquire", "evidence.link", "evidence.evaluate"} {
		if _, ok := accepted[operation]; !ok {
			return // A -run subtest selection may omit this family.
		}
	}
	created := accepted["envelope.create"]
	initial := created.Semantics["envelope"].(map[string]any)
	active := accepted["envelope.activate"].Semantics["envelope"].(map[string]any)
	claim := accepted["claim.acquire"].Semantics["claim"].(map[string]any)
	linkEvaluation := accepted["evidence.link"].Semantics["evaluation"]
	evaluations := accepted["evidence.evaluate"].Semantics["evaluations"].([]any)
	if initial["id"] != active["id"] || initial["id"] != claim["envelope_id"] ||
		initial["created_at"] != active["created_at"] || len(evaluations) != 1 ||
		!reflect.DeepEqual(linkEvaluation, evaluations[0]) {
		t.Fatalf("amended Envelope/Claim/Evidence cross-case identities disagree")
	}
}
