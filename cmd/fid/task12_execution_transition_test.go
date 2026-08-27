package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestExecutionBoundNonterminalTransitionRequiresActiveEnvelopeClaimRunAndRevision(t *testing.T) {
	workspace, issueID, envelopeID, repository := executionTransitionFixture(t)
	before := journalEventCount(t, workspace)

	transition := executeCommandForTest(workspace,
		"--project", workspace, "--json", "issue", "transition", issueID,
		"--to", "in_progress", "--actor", "agent:actor-builder", "--expected-revision", "1", "--idempotency-key", "execution-transition-1",
	)
	if transition.err != nil || transition.stderr != "" {
		t.Fatalf("qualified nonterminal transition = %#v, want stdout-only success", transition)
	}
	if after := journalEventCount(t, workspace); after != before+1 {
		t.Fatalf("journal event count = %d, want %d", after, before+1)
	}

	for _, step := range []struct {
		toState  string
		revision string
	}{
		{toState: "dev_complete", revision: "2"},
		{toState: "qa_review", revision: "3"},
	} {
		result := executeCommandForTest(workspace,
			"--project", workspace, "--json", "issue", "transition", issueID,
			"--to", step.toState, "--actor", "agent:actor-builder", "--expected-revision", step.revision, "--idempotency-key", "execution-transition-"+step.toState,
		)
		if result.err != nil || result.stderr != "" {
			t.Fatalf("transition to %s = %#v, want stdout-only success", step.toState, result)
		}
	}

	terminalBefore := journalEventCount(t, workspace)
	terminal := executeCommandForTest(workspace,
		"--project", workspace, "--json", "issue", "transition", issueID,
		"--to", "qa_pass", "--actor", "agent:actor-builder", "--expected-revision", "4", "--idempotency-key", "execution-transition-terminal",
	)
	if terminal.err == nil || terminal.stdout != "" || !strings.Contains(terminal.stderr, `"obligation-evidence-absent"`) {
		t.Fatalf("terminal transition = %#v, want fail-closed missing Evidence rejection", terminal)
	}
	if after := journalEventCount(t, workspace); after != terminalBefore {
		t.Fatalf("terminal rejection changed journal event count to %d, want %d", after, terminalBefore)
	}
	if repository["worktree_id"] != "worktree-1" || envelopeID == "" {
		t.Fatalf("fixture repository/envelope = %#v/%q, want active execution binding", repository, envelopeID)
	}
}

func TestExecutionBoundTransitionReturnsEveryOrderedReasonWithoutMutation(t *testing.T) {
	workspace, issueID, _, _ := executionTransitionFixture(t)
	before := journalEventCount(t, workspace)

	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "issue", "transition", issueID,
		"--to", "missing-state", "--actor", "agent:actor-other", "--expected-revision", "2", "--idempotency-key", "execution-transition-all-reasons",
	)
	if result.err == nil || result.stdout != "" {
		t.Fatalf("multi-failure transition = %#v, want closed failure", result)
	}
	want := `"unsatisfied":["authority-unavailable","claim-missing","run-inactive","issue-revision-conflict","workflow-transition-invalid"]`
	if !strings.Contains(result.stderr, want) {
		t.Fatalf("multi-failure stderr = %q, want ordered reasons %s", result.stderr, want)
	}
	if after := journalEventCount(t, workspace); after != before {
		t.Fatalf("multi-failure rejection changed journal event count to %d, want %d", after, before)
	}
}

func executionTransitionFixture(t *testing.T) (string, string, string, map[string]any) {
	t.Helper()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	input := envelopeCreateInput(t, workspace, issueID, registration)
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "envelope.activate", "issue.transition", "run.start"}
	grant["actor_selector"] = map[string]any{
		"actor_id":              "actor-builder",
		"required_capabilities": []string{"claim.acquire", "envelope.activate", "issue.transition", "run.start"},
	}
	path := filepath.Join(workspace, "execution-transition-envelope.json")
	writeActivationInput(t, path, input)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", path, "--idempotency-key", "execution-transition-envelope-create")
	if created.err != nil {
		t.Fatalf("create execution envelope: %v; stderr=%q", created.err, created.stderr)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
	activatePath := filepath.Join(workspace, "execution-transition-envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder", "--input", activatePath, "--idempotency-key", "execution-transition-envelope-activate")
	if activated.err != nil {
		t.Fatalf("activate execution envelope: %v; stderr=%q", activated.err, activated.stderr)
	}
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "execution-transition-claim")
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	runPath := filepath.Join(workspace, "execution-transition-run.json")
	writeActivationInput(t, runPath, map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"claim_id": claim["id"], "claim_version": 1, "repository": repository,
	})
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", runPath, "--idempotency-key", "execution-transition-run-start")
	if started.err != nil {
		t.Fatalf("start execution Run: %v; stderr=%q", started.err, started.stderr)
	}
	return workspace, issueID, envelopeID, repository
}
