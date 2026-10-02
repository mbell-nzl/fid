package main

import "testing"

func TestExecutionTransitionFixtureDoesNotInheritExpiredClaim(t *testing.T) {
	t.Parallel()
	// This test is serial so no parallel case can clone the template while its
	// age is simulated. No sleep or dependence on suite duration is necessary.
	executionTransitionFixture(t)
	fixture := reusableWorkspaceFixtureFor("execution-transition")
	fixture.Lock()
	mutateExecutionTransitionFixture(t, fixture.root,
		`UPDATE claims SET expires_at = '2000-01-01T00:00:00Z'`)
	fixture.Unlock()

	workspace, issueID, _, _ := executionTransitionFixture(t)
	assertExecutionTransitionRejectionWithoutMutation(t, workspace, issueID,
		"missing-state", "agent:actor-builder", "1", "fixture-age-transition",
		[]string{"workflow-transition-invalid"},
	)
}
