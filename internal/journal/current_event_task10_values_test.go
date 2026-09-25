package journal

import "testing"

func TestTask10HumanOutcomeStopDetailsAndConditionAreStrict(t *testing.T) {
	details := `{"envelope_id":"ENV-task10","obligation_id":"OBL-human-release","outcome_class":"release"}`
	condition := `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-task10"},"budget_name":null,"capability_id":null,"outcome_class":"release"}`
	effect := Effect{Table: "stops", Op: "insert", Values: map[string]any{"code": "human-outcome-required", "subject_type": "envelope", "subject_id": "ENV-task10"}}
	if err := validateCurrentJSON(Event{EventType: "stop.raised"}, effect, "details_json", details); err != nil {
		t.Fatalf("valid human-outcome-required Stop details rejected: %v", err)
	}
	if err := validateCurrentJSON(Event{EventType: "stop.raised"}, effect, "resolution_condition_json", condition); err != nil {
		t.Fatalf("valid human-outcome-recorded Condition rejected: %v", err)
	}

	invalid := []struct {
		name, column, raw string
	}{
		{"details-missing-outcome-class", "details_json", `{"envelope_id":"ENV-task10","obligation_id":"OBL-human-release"}`},
		{"details-extra", "details_json", `{"envelope_id":"ENV-task10","obligation_id":"OBL-human-release","outcome_class":"release","extra":true}`},
		{"details-obligation-scalar", "details_json", `{"envelope_id":"ENV-task10","obligation_id":42,"outcome_class":"release"}`},
		{"condition-missing-nullables", "resolution_condition_json", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-task10"},"outcome_class":"release"}`},
		{"condition-extra", "resolution_condition_json", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-task10"},"budget_name":null,"capability_id":null,"outcome_class":"release","extra":true}`},
		{"condition-missing-outcome-class", "resolution_condition_json", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-task10"},"budget_name":null,"capability_id":null}`},
		{"condition-decision-class", "resolution_condition_json", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-task10"},"budget_name":null,"capability_id":null,"decision_class":"release"}`},
		{"condition-two-new-fields", "resolution_condition_json", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-task10"},"budget_name":null,"capability_id":null,"outcome_class":"release","decision_class":"release"}`},
		{"condition-wrong-subject", "resolution_condition_json", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-other"},"budget_name":null,"capability_id":null,"outcome_class":"release"}`},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateCurrentJSON(Event{EventType: "stop.raised"}, effect, tc.column, tc.raw); err == nil {
				t.Fatal("invalid human-outcome-recorded value accepted")
			}
		})
	}
}

// TestTask10HumanOutcomeConditionDispatchPreservesNullableBinding proves the
// Phase 22 budget_name/capability_id nullable-binding relation still fires
// through the per-type Condition dispatch (task 9) for the row this task adds:
// a human-outcome-recorded Condition carrying a stray budget_name is rejected
// exactly like the retained Phase 22 types (VAL-081, re-confirmed by task 10).
func TestTask10HumanOutcomeConditionDispatchPreservesNullableBinding(t *testing.T) {
	raw := `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":"attempts","capability_id":null,"outcome_class":"release"}`
	value, err := decodeCurrentJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	err = validateCurrentValue("Condition", value)
	if err == nil || err.Error() != "condition nullable binding mismatch" {
		t.Fatalf("binding error = %v, want \"condition nullable binding mismatch\"", err)
	}
}

func TestTask10HumanOutcomeConditionShapeAndNonStringDiscriminant(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		ok   bool
	}{
		{"canonical", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"outcome_class":"release"}`, true},
		{"wrong-field", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"decision_class":"release"}`, false},
		{"null-own-field", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"outcome_class":null}`, false},
		{"two-new-fields", `{"type":"human-outcome-recorded","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"outcome_class":"release","decision_class":"release"}`, false},
		{"non-string-type", `{"type":42,"subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"outcome_class":"release"}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := decodeCurrentJSON(tc.raw)
			if err != nil {
				if tc.ok {
					t.Fatalf("decode: %v", err)
				}
				return
			}
			got := validateCurrentValue("Condition", value) == nil
			if got != tc.ok {
				t.Fatalf("accepted = %v, want %v", got, tc.ok)
			}
		})
	}
}
