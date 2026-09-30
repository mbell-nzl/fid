//go:build integration

package integration

import (
	"testing"
	"time"
)

// Phase 21 Claim shape, amended by Phase 23 public-interface.md section 2.
func p23AmendedModel(t *testing.T, operation string, pair p23HumanLeafPair) map[string]map[string]any {
	t.Helper()
	if operation == "envelope.create" {
		return p23EnvelopeCreateModel(t, pair)
	}
	if operation == "envelope.activate" {
		return p23EnvelopeActivateModel(t, pair)
	}
	if operation == "evidence.link" {
		return p23EvidenceLinkModel(t, pair)
	}
	if operation != "claim.acquire" {
		if operation == "evidence.evaluate" {
			// contracts.md section 13: the no-ID form returns every evaluation
			// in Obligation ID order. This fixture has one satisfied Obligation.
			return map[string]map[string]any{operation: {"evaluations": []any{map[string]any{
				"obligation_id": "OBL-p22-corpus", "status": "satisfied", "validity": "valid",
				"freshness": "fresh", "artifact_availability": "available", "reasons": []any{},
			}}}}
		}
		return nil
	}
	input := p23LeafRequest(t, pair)
	claim := map[string]any{
		"id": "<generated:CLM>", "issue_id": input["issue_id"], "envelope_id": input["envelope_id"],
		"envelope_revision": input["envelope_revision"], "actor_id": p23LeafFlag(t, pair, "--actor"),
		"worktree_id": input["worktree_id"], "scope": input["scope"], "mode": input["mode"],
		"version": float64(1), "status": "active", "acquired_at": "<generated:timestamp>",
		"expires_at": "<generated:timestamp>", "released_at": nil, "release_reason": nil,
	}
	return map[string]map[string]any{operation: {"claim": claim, "idempotency_key": p23LeafFlag(t, pair, "--idempotency-key"), "replayed": true}}
}

func p23CheckClaimLease(t *testing.T, outcome p23LeafOutcome) {
	t.Helper()
	claim := outcome.Semantics["claim"].(map[string]any)
	start, err := time.Parse(time.RFC3339Nano, claim["acquired_at"].(string))
	if err != nil {
		t.Fatal(err)
	}
	end, err := time.Parse(time.RFC3339Nano, claim["expires_at"].(string))
	if err != nil {
		t.Fatal(err)
	}
	if end.Sub(start) != time.Hour {
		t.Fatalf("claim lease = %v, want one hour", end.Sub(start))
	}
}
