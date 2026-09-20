package workflow

import (
	"errors"
	"testing"
)

func TestReviewClaimRejectsDirtyPathAdoption(t *testing.T) {
	t.Parallel()
	err := validateClaimModeAdoption(RecoveryClaim{ID: "CLM-review", Mode: "review"}, true)
	var invalid *ClaimModeInvalidError
	if !errors.As(err, &invalid) {
		t.Fatalf("review adoption error = %v, want ClaimModeInvalidError", err)
	}
	if invalid.ClaimID != "CLM-review" || invalid.Mode != "review" || invalid.RequiredMode != "mutation" {
		t.Fatalf("review adoption error = %#v, want review Claim and mutation requirement", invalid)
	}
	if err := validateClaimModeAdoption(RecoveryClaim{ID: "CLM-mutation", Mode: "mutation"}, true); err != nil {
		t.Fatalf("mutation adoption error = %v, want nil", err)
	}
	if err := validateClaimModeAdoption(RecoveryClaim{ID: "CLM-review", Mode: "review"}, false); err != nil {
		t.Fatalf("clean review transition error = %v, want nil", err)
	}
}
