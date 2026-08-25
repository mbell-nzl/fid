package issues

import "testing"

func TestAddressConcernProjectsConcernAddressedActivityInSameEvent(t *testing.T) {
	assertJudgmentUnavailable(t, "human.outcomes.v1", func() error {
		_, err := AddressConcern(AddressConcernOptions{ConcernID: "CON-1", AddressedInCompletionID: "TC-1"})
		return err
	})
}

func TestAddressConcernRejectsWithoutActivityAppend(t *testing.T) {
	assertJudgmentUnavailable(t, "human.outcomes.v1", func() error {
		_, err := AddressConcern(AddressConcernOptions{ConcernID: "CON-1"})
		return err
	})
}
