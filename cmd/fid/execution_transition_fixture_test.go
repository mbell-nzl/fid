package main

import "testing"

func TestExecutionTransitionFixtureDoesNotInheritExpiredClaim(t *testing.T) {
	t.Parallel()
	// Expire a prior root's Claim explicitly, without waiting for suite time.
	priorWorkspace, priorIssueID, _, _ := executionTransitionFixture(t)
	mutateExecutionTransitionFixture(t, priorWorkspace,
		`UPDATE claims SET expires_at = '2000-01-01T00:00:00Z'`)
	assertExecutionTransitionRejectionWithoutMutation(t, priorWorkspace, priorIssueID,
		"in_progress", "agent:actor-builder", "1", "fixture-expired-prior-transition",
		[]string{"claim-expired", "run-inactive"},
	)

	workspace, issueID, _, _ := executionTransitionFixture(t)
	if workspace == priorWorkspace {
		t.Fatal("fresh execution fixture reused the expired Claim's root")
	}
	assertExecutionTransitionRejectionWithoutMutation(t, workspace, issueID,
		"missing-state", "agent:actor-builder", "1", "fixture-age-transition",
		[]string{"workflow-transition-invalid"},
	)
}
