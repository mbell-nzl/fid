package journal

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const projectionRepositoryJSON = `{"base_oid":null,"capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":null,"project_id":"PRJ-projection","tree_oid":null,"worktree_id":"WT-projection"}`

func projectionFixtureJSON(value any) string {
	raw, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(raw)
}

func projectionFixtureDecode(raw string) any {
	var value any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		panic(err)
	}
	return value
}

func projectionReconciliation() map[string]any {
	v := map[string]any{"repository": projectionFixtureDecode(projectionRepositoryJSON), "checkpoint_id": nil, "paths": []any{}, "unmerged_paths": []any{}, "operation": "none", "ancestry": "unchanged"}
	v["digest"] = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(projectionFixtureJSON(v))))
	return v
}

func projectionCheckpoint() map[string]any {
	return map[string]any{"id": "OBS-checkpoint", "run_id": "RUN-checkpoint", "kind": "checkpoint", "actor_id": "agent:checkpoint", "observed_at": "2026-09-06T00:00:00Z", "reconciliation": projectionReconciliation()}
}

// These explicit expected bodies accompany the hand-authored fixture facts.
// They do not call the current Event validator or its result reconstruction.
func projectionResult(event Event) map[string]any {
	result := map[string]any{"idempotency_key": "projection-fixture", "replayed": false}
	switch event.EventType {
	case "run.checkpointed":
		result["checkpoint"], result["run"] = projectionCheckpoint(), projectionRun("checkpoint", event)
	case "run.closed", "run.abandoned":
		result["run"] = projectionRun("closure", event)
		if event.EventType == "run.abandoned" {
			result["claim"] = projectionClaim("closure", event)
		}
	case "run.recovered":
		result["run"], result["claim"] = projectionRun("successor", event), projectionClaim("successor", event)
		result["recovery"] = map[string]any{"id": "REC-projection", "prior_run_id": "RUN-prior", "successor_run_id": "RUN-successor", "prior_claim_id": "CLM-prior", "successor_claim_id": "CLM-successor", "observation_id": "OBS-liveness", "reconciliation": projectionReconciliation(), "adopted_paths": []any{}, "authority": map[string]any{"mode": "ordinary", "actor_id": "agent:recovery", "reason": "test", "reference": nil, "grant_ref": nil}, "event_id": event.EventID}
	case "run.budget-observed":
		result["sample"] = map[string]any{"id": "BUD-projection", "envelope_id": "ENV-projection", "envelope_revision": 1, "budget_name": "tokens", "run_id": "RUN-budget", "actor_id": "agent:budget", "source": "activated-harness", "source_sample_id": "sample-1", "sample_kind": "baseline", "previous_sample_id": nil, "value": 0, "measurement_ref": nil, "observed_at": "2026-09-06T00:01:00Z", "event_id": event.EventID}
		budget := projectionBudget()
		budget["consumed"], budget["remaining"], budget["state"] = 0, 10, "available"
		result["budget"], result["stops"], result["attention"] = budget, []any{}, []any{}
	case "run.execution-blocked":
		budget := projectionBudget()
		budget["state"] = "unknown"
		result["denial"] = map[string]any{"budget": budget, "stop_ids": []string{"STP-blocked"}, "attention_ids": []string{"ATT-blocked"}, "event_id": event.EventID, "event_occurred_at": "2026-09-06T00:01:00Z"}
	default:
		result["activation"] = projectionActivation(event)
		if event.EventType == "activation.resource-observed" {
			result["resource"] = projectionResources(event)[1]
		}
	}
	return result
}

func projectionRun(suffix string, event Event) map[string]any {
	actor, issue, started := "agent:"+suffix, "ISS-"+suffix, "2026-09-06T00:00:00Z"
	if suffix == "successor" {
		actor, issue, started = "agent:recovery", "ISS-projection", "2026-09-06T00:01:00Z"
	}
	run := map[string]any{"id": "RUN-" + suffix, "issue_id": issue, "envelope_id": strings.Replace(issue, "ISS-", "ENV-", 1), "envelope_revision": 1, "claim_id": "CLM-" + suffix, "claim_version": 1, "actor_id": actor, "start_repository": projectionFixtureDecode(projectionRepositoryJSON), "terminal_repository": nil, "status": "active", "budgets": []any{}, "started_at": started, "ended_at": nil, "version": 2, "executor": map[string]any{"id": "executor-" + suffix, "source": "activated-harness"}, "obligation_id": nil, "predecessor_run_id": nil, "successor_run_id": nil}
	if suffix == "successor" {
		run["version"], run["predecessor_run_id"] = 1, "RUN-prior"
	}
	for _, effect := range event.Effects {
		if effect.Table == "runs" && effect.PK["id"] == run["id"] && effect.Op == "update" {
			run["status"], run["ended_at"] = effect.Values["status"], effect.Values["ended_at"]
			run["terminal_repository"] = projectionFixtureDecode(effect.Values["terminal_repository_json"].(string))
		}
	}
	return run
}

func projectionClaim(suffix string, event Event) map[string]any {
	issue, actor, acquired, expires := "ISS-closure", "agent:closure", "2026-09-06T00:00:00Z", "2026-09-07T00:00:00Z"
	if suffix == "successor" {
		issue, actor, acquired, expires = "ISS-projection", "agent:recovery", "2026-09-06T00:01:00Z", "2026-09-07T00:01:00Z"
	}
	claim := map[string]any{"id": "CLM-" + suffix, "issue_id": issue, "envelope_id": strings.Replace(issue, "ISS-", "ENV-", 1), "envelope_revision": 1, "actor_id": actor, "worktree_id": "WT-" + suffix, "scope": map[string]any{"paths": []any{}}, "mode": "mutation", "version": 1, "status": "active", "acquired_at": acquired, "expires_at": expires, "released_at": nil, "release_reason": nil}
	if suffix == "successor" {
		claim["worktree_id"] = "WT-projection"
	}
	if event.EventType == "run.abandoned" {
		claim["version"], claim["status"], claim["released_at"], claim["release_reason"] = 2, "released", "2026-09-06T00:00:00Z", "abandoned"
	}
	return claim
}

func projectionBudget() map[string]any {
	return map[string]any{"name": "tokens", "scope": "envelope", "unit": "input-tokens", "maximum": 10, "warning_threshold": nil, "exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil}
}

func seedProjectionBudget(t *testing.T, conn *sql.DB, runID string) {
	t.Helper()
	if _, err := conn.Exec(`INSERT INTO run_budget_bindings (run_id,name,budget_json) VALUES (?, 'tokens', ?)`, runID, projectionFixtureJSON(projectionBudget())); err != nil {
		t.Fatal(err)
	}
}

func projectionResources(event Event) []any {
	resources := []any{}
	for _, component := range []string{"backend", "kit", "overlay", "project", "provider"} {
		v := map[string]any{"component": component, "version": 1, "old_identity": nil, "new_identity": component + ":new", "observed_identity": nil, "state": "pending", "evidence_ref": nil, "error_code": nil}
		if event.EventType == "activation.recovery-completed" {
			v["state"], v["evidence_ref"] = "old", "test:prior-agreement"
		}
		if event.EventType == "activation.resource-observed" && component == "kit" {
			v["version"], v["observed_identity"], v["state"], v["evidence_ref"] = 2, "kit:new", "new", "test:kit"
		}
		resources = append(resources, v)
	}
	return resources
}

func projectionActivation(event Event) map[string]any {
	v := map[string]any{"id": "ACT-projection", "prior_activation_id": nil, "kit_identity": "kit", "backend_identity": "fid", "phase": 21, "status": "prepared", "point_of_no_return": false, "direction": "activate", "created_at": "2026-09-06T00:00:00Z", "updated_at": "2026-09-06T00:01:00Z", "version": 2, "component_digests": []any{}, "resources": projectionResources(event), "recovery": nil}
	if event.EventType != "activation.resource-observed" {
		r := map[string]any{"direction": "restore-prior", "actor_attribution": "fid-bootstrap", "reason": "test", "selected_event_id": event.EventID, "completed_event_id": nil}
		if event.EventType == "activation.recovery-completed" {
			v["status"], r["selected_event_id"], r["completed_event_id"] = "aborted", "EVT-selected", event.EventID
		}
		v["recovery"] = r
	}
	return v
}

// Complete expected snapshot for the activation used by sequence/idempotency tests.
func testActivationResultJSON(key string, replayed bool) string {
	activation := map[string]any{"id": "ACT-1", "prior_activation_id": nil,
		"kit_identity": "kit-1", "backend_identity": "fid-1", "phase": 1,
		"status": "prepared", "point_of_no_return": false, "direction": "activate",
		"created_at": "2026-08-22T00:00:00Z", "updated_at": "2026-08-22T00:00:00Z",
		"version": 1, "component_digests": []any{}, "resources": projectionResources(Event{}), "recovery": nil}
	return projectionFixtureJSON(map[string]any{"activation": activation, "idempotency_key": key, "replayed": replayed})
}
