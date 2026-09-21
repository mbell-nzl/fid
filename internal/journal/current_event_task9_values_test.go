package journal

import "testing"

func TestTask9ProtectedDecisionStopDetailsAndConditionAreStrict(t *testing.T) {
	details := `{"envelope_id":"ENV-task9","decision_class":"release","required_capabilities":["decision.record-reference"],"affected_actions":["envelope.complete"]}`
	condition := `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-task9"},"budget_name":null,"capability_id":null,"decision_class":"release"}`
	effect := Effect{Table: "stops", Op: "insert", Values: map[string]any{"code": "protected-decision-required", "subject_type": "envelope", "subject_id": "ENV-task9"}}
	if err := validateCurrentJSON(Event{EventType: "stop.raised"}, effect, "details_json", details); err != nil {
		t.Fatalf("valid protected Decision Stop details rejected: %v", err)
	}
	if err := validateCurrentJSON(Event{EventType: "stop.raised"}, effect, "resolution_condition_json", condition); err != nil {
		t.Fatalf("valid protected Decision Condition rejected: %v", err)
	}

	invalid := []struct {
		name, column, raw string
	}{
		{"details-missing-class", "details_json", `{"envelope_id":"ENV-task9","required_capabilities":["decision.record-reference"],"affected_actions":["envelope.complete"]}`},
		{"details-extra", "details_json", `{"envelope_id":"ENV-task9","decision_class":"release","required_capabilities":["decision.record-reference"],"affected_actions":["envelope.complete"],"extra":true}`},
		{"details-capability-scalar", "details_json", `{"envelope_id":"ENV-task9","decision_class":"release","required_capabilities":"decision.record-reference","affected_actions":["envelope.complete"]}`},
		{"condition-missing-nullables", "resolution_condition_json", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-task9"},"decision_class":"release"}`},
		{"condition-extra", "resolution_condition_json", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-task9"},"budget_name":null,"capability_id":null,"decision_class":"release","extra":true}`},
		{"condition-missing-class", "resolution_condition_json", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-task9"},"budget_name":null,"capability_id":null}`},
		{"condition-outcome-class", "resolution_condition_json", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-task9"},"budget_name":null,"capability_id":null,"outcome_class":"release"}`},
		{"condition-two-new-fields", "resolution_condition_json", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-task9"},"budget_name":null,"capability_id":null,"decision_class":"release","outcome_class":"release"}`},
		{"condition-wrong-subject", "resolution_condition_json", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-other"},"budget_name":null,"capability_id":null,"decision_class":"release"}`},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateCurrentJSON(Event{EventType: "stop.raised"}, effect, tc.column, tc.raw); err == nil {
				t.Fatal("invalid protected Decision value accepted")
			}
		})
	}
}
