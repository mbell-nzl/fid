package doneeval

import "testing"

func TestMissingPolicyDoesNotWaiveAssuredGates(t *testing.T) {
	t.Parallel()
	for _, missing := range []string{"policy", "test paths"} {
		t.Run(missing, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", true)
			if missing == "policy" {
				s.Policy = nil
			} else {
				s.Policy.TestAdequacySet = false
			}
			s.Approval.State = "missing"
			s.Rounds = nil
			s.GitHub.PR.Body = ""
			got := Evaluate(s)
			if got.Status != Refused || !got.NoOtherModelReview {
				t.Fatalf("outcome=%+v", got)
			}
			for _, gate := range []Gate{PrimaryGate, SeparateGate, ApprovalGate, RollbackGate} {
				if got.Gates[gate].State != Missing {
					t.Fatalf("known Assured gate %s must remain missing: %+v", gate, got.Gates[gate])
				}
			}
			if got.Gates[ApprovalGate].Reason != ApprovalMissing || got.Gates[RollbackGate].Reason != RollbackMissing {
				t.Fatalf("independent evidence gates=%+v", got.Gates)
			}
		})
	}
}
