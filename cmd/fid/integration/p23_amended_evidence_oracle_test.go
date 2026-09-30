//go:build integration

package integration

import "testing"

// public-interface.md section 2 and contracts.md section 13 fix the linked
// receipt's producer, trust source, and satisfied Evaluation.
func p23EvidenceLinkModel(t *testing.T, pair p23HumanLeafPair) map[string]map[string]any {
	t.Helper()
	input := p23LeafRequest(t, pair)
	value := map[string]any{
		"id": "<generated:EVD>", "obligation_id": input["obligation_id"], "evidence_type": input["evidence_type"],
		"receipt_id": input["receipt_id"], "review_id": input["review_id"], "human_outcome_id": input["human_outcome_id"],
		"producer": "fid", "trust_source": "fid-executed", "target_repository": input["target_repository"],
		"validity": "valid", "required_artifact_ids": input["required_artifact_ids"], "created_at": "<generated:timestamp>",
	}
	evaluation := map[string]any{
		"obligation_id": input["obligation_id"], "status": "satisfied", "validity": "valid",
		"freshness": "fresh", "artifact_availability": "available", "reasons": []any{},
	}
	return map[string]map[string]any{"evidence.link": {
		"evidence": value, "evaluation": evaluation,
		"idempotency_key": p23LeafFlag(t, pair, "--idempotency-key"), "replayed": true,
	}}
}
