package issues

import "testing"

func TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites(t *testing.T) {
	assertJudgmentUnavailable(t, "human.outcomes.v1", func() error {
		_, err := RaiseConcern(RaiseConcernOptions{IssueID: "FID-1", RaisedInCompletionID: "TC-1", Summary: "deferred"})
		return err
	})
}
