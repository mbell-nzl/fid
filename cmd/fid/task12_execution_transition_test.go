package main

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestExecutionBoundNonterminalTransitionRequiresActiveEnvelopeClaimRunAndRevision(t *testing.T) {
	t.Parallel()
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

	assertExecutionTransitionRejectionWithoutMutation(t, workspace, issueID,
		"qa_pass", "agent:actor-builder", "4", "execution-transition-terminal",
		[]string{"obligation-evidence-absent"},
	)
	if repository["worktree_id"] != "worktree-1" || envelopeID == "" {
		t.Fatalf("fixture repository/envelope = %#v/%q, want active execution binding", repository, envelopeID)
	}
}

func TestExecutionBoundTransitionReturnsEveryOrderedReasonWithoutMutation(t *testing.T) {
	t.Parallel()
	workspace, issueID, _, _ := executionTransitionFixture(t)
	expectedRevision := 2
	evaluation, err := workflow.EvaluateExecutionTransition(workspace, issueMigrationsDir(), workflow.ExecutionTransitionRequest{
		IssueID: issueID, ActorID: "actor-other", ExpectedIssueRevision: &expectedRevision,
		ActualIssueRevision: 1, WorkflowEdgeValid: false,
	})
	if err != nil {
		t.Fatalf("evaluate multi-failure execution transition: %v", err)
	}
	assertExecutionTransitionReasonGolden(t, evaluation.Reasons, []workflow.EvaluationReason{
		{Code: "authority-unavailable", Detail: "issue.transition", SubjectID: "actor-other", SubjectType: "actor"},
		{Code: "claim-missing", Detail: "active mutation claim", SubjectID: issueID, SubjectType: "issue"},
		{Code: "run-inactive", Detail: "active run", SubjectID: issueID, SubjectType: "issue"},
		{Code: "issue-revision-conflict", Detail: "expected current Issue revision", SubjectID: issueID, SubjectType: "issue"},
		{Code: "workflow-transition-invalid", Detail: "workflow edge", SubjectID: issueID, SubjectType: "issue"},
	})
	assertExecutionTransitionRejectionWithoutMutation(t, workspace, issueID,
		"missing-state", "agent:actor-other", "2", "execution-transition-all-reasons",
		[]string{"authority-unavailable", "claim-missing", "run-inactive", "issue-revision-conflict", "workflow-transition-invalid"},
	)
}

func TestExecutionBoundTransitionRejectsMismatchedClaimWorktreeWithoutMutation(t *testing.T) {
	t.Parallel()
	workspace, issueID, _, _ := executionTransitionFixture(t)
	mutateExecutionTransitionFixture(t, workspace, `UPDATE claims SET worktree_id = 'other-worktree' WHERE issue_id = ?`, issueID)
	assertExecutionTransitionRejectionWithoutMutation(t, workspace, issueID,
		"in_progress", "agent:actor-builder", "1", "execution-transition-claim-worktree",
		[]string{"claim-missing", "run-inactive"},
	)
}

func TestExecutionBoundTransitionRetainsClaimHolderAndExpiryNoMutationGuards(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		updateClaim string
		want        []string
	}{
		{
			name:        "wrong holder",
			updateClaim: `UPDATE claims SET actor_id = 'actor-other' WHERE issue_id = ?`,
			want:        []string{"claim-missing", "run-inactive"},
		},
		{
			name:        "expired",
			updateClaim: `UPDATE claims SET expires_at = '2000-01-01T00:00:00Z' WHERE issue_id = ?`,
			want:        []string{"claim-expired", "run-inactive"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace, issueID, _, _ := executionTransitionFixture(t)
			mutateExecutionTransitionFixture(t, workspace, tt.updateClaim, issueID)
			assertExecutionTransitionRejectionWithoutMutation(t, workspace, issueID,
				"in_progress", "agent:actor-builder", "1", "execution-transition-claim-"+strings.ReplaceAll(tt.name, " ", "-"), tt.want,
			)
		})
	}
}

func TestExecutionBoundTransitionPredicateMatrix(t *testing.T) {
	t.Parallel()
	rows := []struct {
		name     string
		prepare  func(t *testing.T, workspace, issueID, envelopeID string)
		toState  string
		actor    string
		revision string
		want     []string
	}{
		{
			name:    "qualifying",
			toState: "in_progress", actor: "agent:actor-builder", revision: "1",
		},
		{
			name: "envelope inactive",
			prepare: func(t *testing.T, workspace, _, envelopeID string) {
				mutateExecutionTransitionFixture(t, workspace, `UPDATE work_envelopes SET status = 'draft' WHERE id = ?`, envelopeID)
			},
			toState: "in_progress", actor: "agent:actor-builder", revision: "1",
			want: []string{"authority-unavailable", "envelope-inactive"},
		},
		{
			name: "actor capability missing",
			prepare: func(t *testing.T, workspace, _, _ string) {
				mutateExecutionTransitionFixture(t, workspace, `UPDATE activation_actor_profiles SET capabilities_json = '[]' WHERE actor_id = 'actor-builder'`)
			},
			toState: "in_progress", actor: "agent:actor-builder", revision: "1",
			want: []string{"authority-unavailable"},
		},
		{
			name: "claim worktree mismatched",
			prepare: func(t *testing.T, workspace, issueID, _ string) {
				mutateExecutionTransitionFixture(t, workspace, `UPDATE claims SET worktree_id = 'other-worktree' WHERE issue_id = ?`, issueID)
			},
			toState: "in_progress", actor: "agent:actor-builder", revision: "1",
			want: []string{"claim-missing", "run-inactive"},
		},
		{
			name: "usable claim without active run",
			prepare: func(t *testing.T, workspace, issueID, _ string) {
				mutateExecutionTransitionFixture(t, workspace, `UPDATE runs SET status = 'failed' WHERE issue_id = ?`, issueID)
			},
			toState: "in_progress", actor: "agent:actor-builder", revision: "1",
			want: []string{"run-inactive"},
		},
		{
			name: "run repository mismatched",
			prepare: func(t *testing.T, workspace, issueID, _ string) {
				mutateExecutionTransitionFixture(t, workspace, `UPDATE runs SET start_repository_json = json_set(start_repository_json, '$.worktree_id', 'other-worktree') WHERE issue_id = ?`, issueID)
			},
			toState: "in_progress", actor: "agent:actor-builder", revision: "1",
			want: []string{"run-inactive"},
		},
		{
			name:    "Issue revision mismatched",
			toState: "in_progress", actor: "agent:actor-builder", revision: "2",
			want: []string{"issue-revision-conflict"},
		},
		{
			name:    "workflow edge invalid",
			toState: "missing-state", actor: "agent:actor-builder", revision: "1",
			want: []string{"workflow-transition-invalid"},
		},
	}

	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			workspace, issueID, envelopeID, _ := executionTransitionFixture(t)
			if row.prepare != nil {
				row.prepare(t, workspace, issueID, envelopeID)
			}
			if row.want == nil {
				before := executionTransitionIssueSnapshot(t, workspace, issueID)
				beforeEvents := journalEventCount(t, workspace)
				result := executeCommandForTest(workspace,
					"--project", workspace, "--json", "issue", "transition", issueID,
					"--to", row.toState, "--actor", row.actor, "--expected-revision", row.revision, "--idempotency-key", "execution-transition-matrix-qualifying",
				)
				if result.err != nil || result.stderr != "" {
					t.Fatalf("qualifying transition = %#v, want stdout-only success", result)
				}
				after := executionTransitionIssueSnapshot(t, workspace, issueID)
				if after.State != "in_progress" || after.Revision != before.Revision+1 {
					t.Fatalf("qualifying Issue snapshot = %#v, want state in_progress and revision %d", after, before.Revision+1)
				}
				if afterEvents := journalEventCount(t, workspace); afterEvents != beforeEvents+1 {
					t.Fatalf("qualifying journal event count = %d, want %d", afterEvents, beforeEvents+1)
				}
				return
			}
			assertExecutionTransitionRejectionWithoutMutation(t, workspace, issueID, row.toState, row.actor, row.revision, "execution-transition-matrix-"+strings.ReplaceAll(row.name, " ", "-"), row.want)
		})
	}
}

func assertExecutionTransitionReasonGolden(t *testing.T, got, want []workflow.EvaluationReason) {
	t.Helper()
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal execution transition reasons: %v", err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal execution transition golden: %v", err)
	}
	if string(gotJSON) != string(wantJSON) {
		t.Fatalf("execution transition reasons = %s, want byte-exact %s", gotJSON, wantJSON)
	}
}

type executionTransitionIssueState struct {
	State    string
	Revision int
}

func executionTransitionIssueSnapshot(t *testing.T, workspace, issueID string) executionTransitionIssueState {
	t.Helper()
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "get", issueID)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("Issue snapshot = %#v, want stdout-only success", result)
	}
	envelope := decodeCommandJSONStructure(t, result.stdout).(map[string]any)
	issue := envelope["result"].(map[string]any)["issue"].(map[string]any)
	revision, err := issue["revision"].(json.Number).Int64()
	if err != nil {
		t.Fatalf("Issue snapshot revision = %#v, want integer: %v", issue["revision"], err)
	}
	state, ok := issue["status"].(string)
	if !ok {
		t.Fatalf("Issue snapshot state = %#v, want string", issue["status"])
	}
	return executionTransitionIssueState{State: state, Revision: int(revision)}
}

func assertExecutionTransitionRejectionWithoutMutation(t *testing.T, workspace, issueID, toState, actor, revision, key string, wantReasons []string) {
	t.Helper()
	beforeIssue := executionTransitionIssueSnapshot(t, workspace, issueID)
	beforeEvents := journalEventCount(t, workspace)
	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "issue", "transition", issueID,
		"--to", toState, "--actor", actor, "--expected-revision", revision, "--idempotency-key", key,
	)
	if result.err == nil || result.stdout != "" {
		t.Fatalf("execution-bound rejection = %#v, want stderr-only failure", result)
	}
	envelope := decodeCommandJSONStructure(t, result.stderr).(map[string]any)
	assertExactKeys(t, envelope, []string{"contract_version", "error", "ok", "operation"})
	errorObject := envelope["error"].(map[string]any)
	assertExactKeys(t, errorObject, []string{"category", "code", "details", "message", "retryable", "subject"})
	if envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != false || envelope["operation"] != "issue.transition" ||
		errorObject["category"] != "precondition" || errorObject["code"] != "transition-not-permitted" || errorObject["message"] != "Issue transition is not permitted" || errorObject["retryable"] != false {
		t.Fatalf("execution-bound failure envelope = %#v, want canonical transition rejection", envelope)
	}
	wantDetails := map[string]any{"from_state": beforeIssue.State, "to_state": toState, "unsatisfied": stringValues(wantReasons)}
	if gotDetails, ok := errorObject["details"].(map[string]any); !ok || !sameJSONMap(gotDetails, wantDetails) {
		t.Fatalf("execution-bound failure details = %#v, want %#v", errorObject["details"], wantDetails)
	}
	if gotSubject, ok := errorObject["subject"].(map[string]any); !ok || !sameJSONMap(gotSubject, map[string]any{"type": "issue", "id": issueID}) {
		t.Fatalf("execution-bound failure subject = %#v, want Issue %q", errorObject["subject"], issueID)
	}
	afterIssue := executionTransitionIssueSnapshot(t, workspace, issueID)
	if afterIssue != beforeIssue {
		t.Fatalf("rejected execution-bound transition Issue snapshot = %#v, want unchanged %#v", afterIssue, beforeIssue)
	}
	if afterEvents := journalEventCount(t, workspace); afterEvents != beforeEvents {
		t.Fatalf("rejected execution-bound transition journal event count = %d, want %d", afterEvents, beforeEvents)
	}
}

func stringValues(values []string) []any {
	result := make([]any, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}

func mutateExecutionTransitionFixture(t *testing.T, workspace, statement string, args ...any) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open execution database: %v", err)
	}
	if _, err := conn.Exec(statement, args...); err != nil {
		_ = conn.Close()
		t.Fatalf("mutate execution fixture: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close execution database: %v", err)
	}
}

func executionTransitionFixture(t *testing.T) (string, string, string, map[string]any) {
	t.Helper()
	workspace, values := cloneReusableWorkspaceFixture(t, "execution-transition", func() (string, []string) {
		workspace, issueID, envelopeID, repository := buildExecutionTransitionFixture(t)
		repositoryJSON, err := json.Marshal(repository)
		if err != nil {
			t.Fatalf("encode execution-transition repository: %v", err)
		}
		return workspace, []string{issueID, envelopeID, string(repositoryJSON)}
	})
	repository := make(map[string]any)
	if err := json.Unmarshal([]byte(values[2]), &repository); err != nil {
		t.Fatalf("decode execution-transition repository: %v", err)
	}
	return workspace, values[0], values[1], repository
}

func buildExecutionTransitionFixture(t *testing.T) (string, string, string, map[string]any) {
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
	repository := runRepositoryForTest(t, workspace)
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
