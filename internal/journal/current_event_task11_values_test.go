package journal

import "testing"

// TestTask11ReviewRecordedAndFindingResolvedConditionShapes completes the
// four-type judgment Condition encoding matrix (VAL-080): review-recorded
// carries no discriminator at all, and finding-resolved carries exactly
// finding_id. Every malformed combination rejects with the DEC-235 shape
// evidence class - closed object has wrong shape for an extra field, missing
// required field for the wrong one, and invalid ID for a null own field -
// never the nullable-binding relation message.
func TestTask11ReviewRecordedAndFindingResolvedConditionShapes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		raw  string
		ok   bool
	}{
		{"review-recorded canonical bare form", `{"type":"review-recorded","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null}`, true},
		{"review-recorded rejects finding_id", `{"type":"review-recorded","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null,"finding_id":"FND-pinned"}`, false},
		{"review-recorded rejects decision_class", `{"type":"review-recorded","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null,"decision_class":"release"}`, false},
		{"review-recorded rejects outcome_class", `{"type":"review-recorded","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null,"outcome_class":"release"}`, false},
		{"review-recorded rejects present-null finding_id", `{"type":"review-recorded","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null,"finding_id":null}`, false},
		{"finding-resolved canonical five-field form", `{"type":"finding-resolved","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null,"finding_id":"FND-pinned"}`, true},
		{"finding-resolved requires finding_id", `{"type":"finding-resolved","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null}`, false},
		{"finding-resolved rejects null finding_id", `{"type":"finding-resolved","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null,"finding_id":null}`, false},
		{"finding-resolved rejects the wrong field", `{"type":"finding-resolved","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null,"decision_class":"release"}`, false},
		{"finding-resolved rejects two discriminators", `{"type":"finding-resolved","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null,"finding_id":"FND-pinned","decision_class":"release"}`, false},
		{"finding-resolved subject may be an envelope", `{"type":"finding-resolved","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"finding_id":"FND-pinned"}`, true},
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

// TestTask11FindingResolvedConditionDispatchPreservesNullableBinding proves
// the Phase 22 budget_name/capability_id nullable-binding relation still
// fires through the per-type Condition dispatch for the last row this task
// adds (VAL-081): a finding-resolved Condition carrying a stray budget_name
// is rejected with the exact observed relation message, not a shape error.
func TestTask11FindingResolvedConditionDispatchPreservesNullableBinding(t *testing.T) {
	t.Parallel()
	raw := `{"type":"finding-resolved","subject":{"type":"run","id":"RUN-pinned"},"budget_name":"attempts","capability_id":null,"finding_id":"FND-pinned"}`
	value, err := decodeCurrentJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	err = validateCurrentValue("Condition", value)
	if err == nil || err.Error() != "condition nullable binding mismatch" {
		t.Fatalf("binding error = %v, want \"condition nullable binding mismatch\"", err)
	}
}

// TestTask11ConditionShapeRegistryIsClosedAtTenRows confirms the closed
// four-type judgment Condition set is complete: the six retained Phase 22
// rows plus the four judgment rows task 9, 10, and 11 register between them.
func TestTask11ConditionShapeRegistryIsClosedAtTenRows(t *testing.T) {
	t.Parallel()
	if len(currentConditionShapes) != 10 {
		t.Fatalf("currentConditionShapes has %d rows, want 10", len(currentConditionShapes))
	}
}
