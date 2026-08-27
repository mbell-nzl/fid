package main

import (
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
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

func TestExecutionBoundTransitionRejectsMismatchedClaimWorktreeWithoutMutation(t *testing.T) {
	workspace, issueID, _, _ := executionTransitionFixture(t)
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open execution database: %v", err)
	}
	if _, err := conn.Exec(`UPDATE work_claims SET worktree_id = 'other-worktree' WHERE issue_id = ?`, issueID); err != nil {
		_ = conn.Close()
		t.Fatalf("mismatch claim worktree: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close execution database: %v", err)
	}

	before := journalEventCount(t, workspace)
	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "issue", "transition", issueID,
		"--to", "in_progress", "--actor", "agent:actor-builder", "--expected-revision", "1", "--idempotency-key", "execution-transition-claim-worktree",
	)
	if result.err == nil || result.stdout != "" {
		t.Fatalf("mismatched Claim worktree transition = %#v, want closed failure", result)
	}
	want := `"unsatisfied":["claim-missing","run-inactive"]`
	if !strings.Contains(result.stderr, want) {
		t.Fatalf("mismatched Claim worktree stderr = %q, want ordered reasons %s", result.stderr, want)
	}
	if after := journalEventCount(t, workspace); after != before {
		t.Fatalf("mismatched Claim worktree changed journal event count to %d, want %d", after, before)
	}
}

func TestExecutionBoundTransitionRetainsClaimHolderAndExpiryNoMutationGuards(t *testing.T) {
	tests := []struct {
		name        string
		updateClaim string
		wantCode    string
	}{
		{
			name:        "wrong holder",
			updateClaim: `UPDATE work_claims SET actor_id = 'actor-other' WHERE issue_id = ?`,
			wantCode:    "claim-missing",
		},
		{
			name:        "expired",
			updateClaim: `UPDATE work_claims SET expires_at = '2000-01-01T00:00:00Z' WHERE issue_id = ?`,
			wantCode:    "claim-expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace, issueID, _, _ := executionTransitionFixture(t)
			conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
			if err != nil {
				t.Fatalf("open execution database: %v", err)
			}
			if _, err := conn.Exec(tt.updateClaim, issueID); err != nil {
				_ = conn.Close()
				t.Fatalf("prepare %s Claim: %v", tt.name, err)
			}
			if err := conn.Close(); err != nil {
				t.Fatalf("close execution database: %v", err)
			}

			before := journalEventCount(t, workspace)
			result := executeCommandForTest(workspace,
				"--project", workspace, "--json", "issue", "transition", issueID,
				"--to", "in_progress", "--actor", "agent:actor-builder", "--expected-revision", "1", "--idempotency-key", "execution-transition-claim-"+strings.ReplaceAll(tt.name, " ", "-"),
			)
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"`+tt.wantCode+`"`) {
				t.Fatalf("%s Claim transition = %#v, want %s rejection", tt.name, result, tt.wantCode)
			}
			if after := journalEventCount(t, workspace); after != before {
				t.Fatalf("journal events after %s Claim rejection = %d, want %d", tt.name, after, before)
			}
		})
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
