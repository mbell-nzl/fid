package journal

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestCurrentBudgetValuesAdmitReviewCyclesPair(t *testing.T) {
	t.Parallel()
	immutable := map[string]any{"name": "review-cycles", "scope": "envelope", "unit": "review-cycles", "maximum": json.Number("1"), "warning_threshold": nil, "exhaustion_action": "review-rework-exhausted", "consumed": nil, "remaining": nil}
	if err := validateCurrentValue("ImmutableBudget", immutable); err != nil {
		t.Fatalf("validateCurrentValue(ImmutableBudget) error = %v", err)
	}
	derived := map[string]any{"name": "review-cycles", "scope": "envelope", "unit": "review-cycles", "maximum": json.Number("1"), "warning_threshold": nil, "exhaustion_action": "review-rework-exhausted", "consumed": json.Number("1"), "remaining": json.Number("0"), "state": "exhausted"}
	if err := validateCurrentValue("Budget", derived); err != nil {
		t.Fatalf("validateCurrentValue(Budget) error = %v", err)
	}
	details := map[string]any{"envelope_id": "ENV-review", "budget_name": "review-cycles", "unit": "review-cycles", "maximum": json.Number("1"), "consumed": json.Number("1"), "remaining": json.Number("0"), "last_review_id": "REV-review"}
	if _, err := validateCurrentShape(details, currentStopDetailShapes["review-rework-exhausted"]); err != nil {
		t.Fatalf("validate review-rework-exhausted details: %v", err)
	}
}

func TestCurrentReviewCycleProjectionStartsKnownAndOwnsFIDSample(t *testing.T) {
	t.Parallel()
	conn := openReplayDB(t, filepath.Join(t.TempDir(), "review-cycle.sqlite"))
	defer closeProjectDB(t, conn)
	seedBudgetProjection(t, conn)
	policy := `{"name":"tokens","scope":"envelope","unit":"review-cycles","maximum":1,"warning_threshold":null,"exhaustion_action":"review-rework-exhausted","consumed":null,"remaining":null}`
	if _, err := conn.Exec(`UPDATE run_budget_bindings SET budget_json=? WHERE run_id='RUN-budget' AND name='tokens'`, policy); err != nil {
		t.Fatal(err)
	}
	tx, err := conn.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	budget, err := currentProjectedBudget(tx, "RUN-budget", "tokens")
	if err != nil {
		t.Fatal(err)
	}
	if budget["consumed"] != int64(0) || budget["remaining"] != int64(1) || budget["state"] != "available" {
		t.Fatalf("review-cycle initial projection = %#v, want known zero", budget)
	}
	effect := Effect{Op: "insert", Table: "budget_samples", PK: map[string]any{"id": "BSM-review"}, Values: map[string]any{
		"envelope_id": "ENV-projection", "envelope_revision": 1, "budget_name": "tokens", "run_id": "RUN-budget",
		"actor_id": "agent:reviewer", "source": "fid", "source_sample_id": "REV-review", "sample_kind": "delta",
		"previous_sample_id": nil, "value": 1, "measurement_ref_json": nil, "observed_at": "2026-09-20T00:00:00Z", "created_event_id": "EVT-review",
	}}
	if err := validateCurrentBudgetSampleChain(tx, effect, "RUN-budget"); err != nil {
		t.Fatalf("fid review-cycle sample ownership = %v", err)
	}
	event := Event{EventType: "review.recorded", Aggregate: Aggregate{ID: "REV-review"}, Effects: []Effect{
		{Op: "insert", Table: "review_records", PK: map[string]any{"id": "REV-review"}, Values: map[string]any{"run_id": "RUN-budget"}}, effect,
	}}
	if err := validateCurrentReferences(tx, event); err != nil {
		t.Fatalf("review sample Run reference = %v", err)
	}
}
