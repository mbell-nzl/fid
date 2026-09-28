package cli

import "testing"

func TestPhase23FailureCodesAreRegistered(t *testing.T) {
	t.Parallel()
	want := []Code{
		CodeReviewNotIndependent, CodeReviewVerdictInvalid, CodeReviewChangesRequired,
		CodeReviewReworkExhausted, CodeFindingAlreadyDisposed, CodeFindingNotFound,
		CodeReviewNotFound, CodeConcernNotFound, CodeDeviationNotFound,
		CodeOutcomeNotFound, CodeRequirementNotFound, CodeConcernAlreadyDisposed,
		CodeClaimModeInvalid, CodeProtectedDecisionRequired, CodeOutcomeNotPermitted,
		CodeEvidenceTypeNotAccepted, CodeSourceDigestMismatch, CodeRepositoryMismatch,
		CodeHumanOutcomeRequired,
	}
	registered := map[Code]bool{}
	for _, code := range RegisteredCodes() {
		registered[code] = true
	}
	for _, code := range want {
		if !registered[code] {
			t.Errorf("Phase 23 failure code %q is not registered", code)
		}
	}
	if registered[Code("finding-unresolved")] {
		t.Fatal("evaluation reason finding-unresolved is registered as a failure code")
	}
}
