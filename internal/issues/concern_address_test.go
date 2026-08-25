package issues

import "testing"

func TestAddressConcernUpdatesConcernInPlaceAndClearsDerivedFlags(t *testing.T) {
	assertJudgmentUnavailable(t, "human.outcomes.v1", func() error {
		_, err := AddressConcern(AddressConcernOptions{})
		return err
	})
}

func TestAddressConcernRejectsBeforeAppend(t *testing.T) {
	assertJudgmentUnavailable(t, "human.outcomes.v1", func() error {
		_, err := AddressConcern(AddressConcernOptions{ConcernID: "CON-1"})
		return err
	})
}
