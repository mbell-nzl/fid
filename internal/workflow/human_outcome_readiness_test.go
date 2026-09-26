package workflow

import "testing"

func TestOnlyUnsatisfiedHumanOutcomeObligationRequiresExactlyOneRemainingHumanOutcome(t *testing.T) {
	t.Parallel()
	class := "release"
	humanOutcome := Obligation{ObligationInput: ObligationInput{ID: "OBL-human", Type: "human-outcome", OutcomeClass: &class, GuardedActions: []string{"envelope.complete"}}}
	gateEvidence := Obligation{ObligationInput: ObligationInput{ID: "OBL-gate", Type: "gate-evidence", GuardedActions: []string{"envelope.complete"}}}

	cases := []struct {
		name        string
		obligations []Obligation
		evaluations []ObligationEvaluation
		wantOK      bool
	}{
		{
			name:        "sole-unsatisfied-human-outcome",
			obligations: []Obligation{humanOutcome},
			evaluations: []ObligationEvaluation{{ObligationID: "OBL-human", Status: "unsatisfied"}},
			wantOK:      true,
		},
		{
			name:        "human-outcome-already-satisfied",
			obligations: []Obligation{humanOutcome},
			evaluations: []ObligationEvaluation{{ObligationID: "OBL-human", Status: "satisfied"}},
			wantOK:      false,
		},
		{
			name:        "another-obligation-also-unsatisfied",
			obligations: []Obligation{humanOutcome, gateEvidence},
			evaluations: []ObligationEvaluation{{ObligationID: "OBL-human", Status: "unsatisfied"}, {ObligationID: "OBL-gate", Status: "unsatisfied"}},
			wantOK:      false,
		},
		{
			name:        "sole-unsatisfied-is-not-human-outcome",
			obligations: []Obligation{humanOutcome, gateEvidence},
			evaluations: []ObligationEvaluation{{ObligationID: "OBL-human", Status: "satisfied"}, {ObligationID: "OBL-gate", Status: "unsatisfied"}},
			wantOK:      false,
		},
		{
			name:        "no-obligations",
			obligations: []Obligation{},
			evaluations: []ObligationEvaluation{},
			wantOK:      false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obligation, ok := onlyUnsatisfiedHumanOutcomeObligation(tc.obligations, tc.evaluations)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v (obligation=%#v)", ok, tc.wantOK, obligation)
			}
			if ok && obligation.ID != "OBL-human" {
				t.Fatalf("selected Obligation = %q, want OBL-human", obligation.ID)
			}
		})
	}
}
