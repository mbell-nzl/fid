//go:build integration

package integration

import (
	"reflect"
	"testing"
)

// public-interface.md section 3 fixes Review, Finding, and disposition
// members; sections 4 and contracts.md section 5 fix the addressed transition.
func p23ReviewFindingModel(t *testing.T, pairs map[string]p23HumanLeafPair) map[string]map[string]any {
	t.Helper()
	reviewPair := pairs["review.record"]
	reviewInput := p23LeafRequest(t, reviewPair)
	findingInput := reviewInput["findings"].([]any)[0].(map[string]any)
	addressPair := pairs["finding.address"]
	addressInput := p23LeafRequest(t, addressPair)
	open := map[string]any{
		"id": findingInput["id"], "subject": findingInput["subject"], "scope": findingInput["scope"],
		"severity": findingInput["severity"], "summary": findingInput["summary"], "blocking": findingInput["blocking"],
		"review_id": "<generated:REV>", "created_at": "<generated:timestamp>",
		"dispositions": []any{}, "state": "open",
	}
	addressed := make(map[string]any, len(open))
	for key, value := range open {
		addressed[key] = value
	}
	addressed["state"] = "addressed"
	addressed["dispositions"] = []any{map[string]any{
		"id": "<generated:FDP>", "action": "addressed", "actor_id": p23LeafFlag(t, addressPair, "--actor"),
		"qualifying_review_id": nil, "reason": addressInput["reason"],
		"disposed_at": "<generated:timestamp>", "event_id": "<generated:EventRef>",
	}}
	review := map[string]any{
		"id": "<generated:REV>", "run_id": reviewInput["run_id"],
		"issue_id":    findingInput["subject"].(map[string]any)["id"],
		"envelope_id": "<generated:ENV>", "envelope_revision": float64(1),
		"actor_id": p23LeafFlag(t, reviewPair, "--actor"), "profile": reviewInput["profile"],
		"repository": reviewInput["repository"], "coverage": reviewInput["coverage"],
		"verdict": reviewInput["verdict"], "evidence_refs": reviewInput["evidence_refs"],
		"receipt_ids": reviewInput["receipt_ids"], "findings": []any{open},
		"resolved_finding_ids": []any{}, "independence": "distinct-actor",
		"recorded_at": "<generated:timestamp>", "event_id": "<generated:EventRef>",
	}
	current := make(map[string]any, len(review))
	for key, value := range review {
		current[key] = value
	}
	current["findings"] = []any{addressed}
	return map[string]map[string]any{
		"review.record":   {"review": review, "budget": nil, "stops": []any{}, "attention": []any{}, "idempotency_key": p23LeafFlag(t, reviewPair, "--idempotency-key"), "replayed": true},
		"review.get":      {"review": current},
		"review.list":     {"items": []any{current}},
		"finding.address": {"finding": addressed, "idempotency_key": p23LeafFlag(t, addressPair, "--idempotency-key"), "replayed": true},
		"finding.list":    {"items": []any{addressed}},
	}
}

func p23CheckReviewFindingIdentity(t *testing.T, accepted map[string]p23LeafOutcome) {
	t.Helper()
	for _, operation := range []string{"review.record", "review.get", "review.list", "finding.address", "finding.list"} {
		if _, ok := accepted[operation]; !ok {
			return
		}
	}
	initial := accepted["review.record"]
	recorded := initial.Semantics["review"].(map[string]any)
	current := accepted["review.get"].Semantics["review"].(map[string]any)
	listedReviews := accepted["review.list"].Semantics["items"].([]any)
	addressed := accepted["finding.address"].Semantics["finding"].(map[string]any)
	listedFindings := accepted["finding.list"].Semantics["items"].([]any)
	currentFindings := current["findings"].([]any)
	if len(listedReviews) != 1 || len(listedFindings) != 1 || len(currentFindings) != 1 ||
		!reflect.DeepEqual(current, listedReviews[0]) || !reflect.DeepEqual(addressed, listedFindings[0]) ||
		!reflect.DeepEqual(addressed, currentFindings[0]) || recorded["id"] != current["id"] ||
		recorded["event_id"] != current["event_id"] || recorded["recorded_at"] != current["recorded_at"] ||
		addressed["review_id"] != recorded["id"] {
		t.Fatalf("Review/Finding mutation and read identities disagree")
	}
}
