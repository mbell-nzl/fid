package evidence

import (
	"errors"
	"testing"
)

// TestHumanOutcomeAdmissionSatisfiesOnlyMatchingClassWithApprovedResult
// closes VAL-044's remaining cell (DEC-243): task 6's matrix
// (TestEvidenceAdmissionTask6PureProductionClassifierMatrix) proved every
// mismatched-Obligation-type cell, but could not record real Human Outcomes,
// so it never exercised a Human Outcome offered to its *matching*
// human-outcome Obligation across all four results. This proves that a
// matching class only ever admits with result "approved"; the other three
// results reject outcome-not-permitted exactly like a mismatched Obligation
// type does.
func TestHumanOutcomeAdmissionSatisfiesOnlyMatchingClassWithApprovedResult(t *testing.T) {
	t.Parallel()
	for _, result := range []string{"approved", "rejected", "changes-required", "acknowledged"} {
		t.Run(result, func(t *testing.T) {
			err := classifyEvidenceAdmission(evidenceAdmission{
				ObligationID:          "OBL-matching",
				ObligationType:        "human-outcome",
				AcceptedEvidenceTypes: []string{"human-outcome"},
				AcceptedTrustSources:  []string{"activated-human"},
				EvidenceType:          "human-outcome",
				TrustSource:           "activated-human",
				HumanOutcomeID:        "OUT-matching",
				HumanOutcomeResult:    result,
			})
			if result == "approved" {
				if err != nil {
					t.Fatalf("classifyEvidenceAdmission() error = %v, want admission for matching class with approved result", err)
				}
				return
			}
			var denied *OutcomeNotPermittedError
			if !errors.As(err, &denied) || denied.ObligationType != "human-outcome" || denied.OutcomeID != "OUT-matching" || denied.Result != result {
				t.Fatalf("classifyEvidenceAdmission() error = %#v, want exact outcome-not-permitted denial for result %q", err, result)
			}
		})
	}
}
