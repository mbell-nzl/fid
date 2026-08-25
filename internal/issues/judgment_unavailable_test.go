package issues

import (
	"errors"
	"testing"
)

func assertJudgmentUnavailable(t *testing.T, capabilityID string, call func() error) {
	t.Helper()
	err := call()
	var unavailable *CapabilityUnavailableError
	if !errors.As(err, &unavailable) {
		t.Fatalf("error = %v, want CapabilityUnavailableError", err)
	}
	if unavailable.CapabilityID != capabilityID || unavailable.OwnerPhase != 23 {
		t.Fatalf("unavailable = %#v, want capability %q owned by phase 23", unavailable, capabilityID)
	}
}

func TestJudgmentMutationsAreUnavailableBeforeAnyWriterWork(t *testing.T) {
	tests := []struct {
		name string
		call func() error
		want string
	}{
		{
			name: "acceptance criterion verification",
			call: func() error {
				_, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{})
				return err
			},
			want: "reviews.coverage.v1",
		},
		{
			name: "concern raise",
			call: func() error {
				_, err := RaiseConcern(RaiseConcernOptions{})
				return err
			},
			want: "human.outcomes.v1",
		},
		{
			name: "concern address",
			call: func() error {
				_, err := AddressConcern(AddressConcernOptions{})
				return err
			},
			want: "human.outcomes.v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			var unavailable *CapabilityUnavailableError
			if !errors.As(err, &unavailable) {
				t.Fatalf("error = %v, want CapabilityUnavailableError", err)
			}
			if unavailable.CapabilityID != tt.want || unavailable.OwnerPhase != 23 {
				t.Fatalf("unavailable = %#v, want capability %q owned by phase 23", unavailable, tt.want)
			}
		})
	}
}
