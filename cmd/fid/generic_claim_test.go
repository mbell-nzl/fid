package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	envelopeInput := envelopeCreateInput(t, workspace, issueID, registration)
	grant := envelopeInput["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "envelope.activate"}
	grant["actor_selector"] = map[string]any{
		"actor_id": "actor-builder", "required_capabilities": []string{"claim.acquire", "envelope.activate"},
	}
	envelopePath := filepath.Join(workspace, "claim-envelope.json")
	writeActivationInput(t, envelopePath, envelopeInput)
	created := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder",
		"--input", envelopePath, "--idempotency-key", "claim-envelope-create",
	)
	if created.err != nil {
		t.Fatalf("create envelope: %v; stderr=%q", created.err, created.stderr)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
	activatePath := filepath.Join(workspace, "claim-envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder",
		"--input", activatePath, "--idempotency-key", "claim-envelope-activate",
	)
	if activated.err != nil {
		t.Fatalf("activate envelope: %v; stderr=%q", activated.err, activated.stderr)
	}

	claimPath := filepath.Join(workspace, "claim-acquire.json")
	writeActivationInput(t, claimPath, map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 60,
	})
	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder",
		"--input", claimPath, "--idempotency-key", "claim-acquire-1",
	)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("claim acquire = %#v, want stdout-only success", result)
	}

	payload := decodeCommandJSONStructure(t, result.stdout).(map[string]any)
	claim := payload["result"].(map[string]any)["claim"].(map[string]any)
	for field, want := range map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": "1", "actor_id": "actor-builder",
		"worktree_id": "worktree-1", "mode": "mutation", "version": "1", "status": "active",
	} {
		if got := stringifyClaimValue(claim[field]); got != want {
			t.Fatalf("claim.%s = %#v, want %#v", field, got, want)
		}
	}
}

func TestClaimAcquireAcceptsReviewMode(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := claimEnvelopeForTest(t)
	path := filepath.Join(workspace, "claim-acquire-review.json")
	worktreeID := testEnvelopeRepository(t, workspace)["worktree_id"]
	writeActivationInput(t, path, map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"worktree_id": worktreeID, "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "mode": "review", "ttl_seconds": 60,
	})
	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder",
		"--input", path, "--idempotency-key", "claim-acquire-review",
	)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("review Claim acquire = %#v, want stdout-only success", result)
	}
	if claim := claimFromResponse(t, result.stdout); claim["mode"] != "review" {
		t.Fatalf("Claim mode = %#v, want review", claim["mode"])
	}
}

func TestClaimAcquireKeepsOneActiveClaimAcrossModes(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := claimEnvelopeForTest(t)
	review, _ := acquireGenericClaimForTestWithMode(t, workspace, issueID, envelopeID, "claim-mode-review", 60, "review")

	mutationPath := filepath.Join(workspace, "claim-mode-mutation.json")
	worktreeID := testEnvelopeRepository(t, workspace)["worktree_id"]
	writeActivationInput(t, mutationPath, map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"worktree_id": worktreeID, "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "mode": "mutation", "ttl_seconds": 60,
	})
	beforeConflict := snapshotWorkspacePersistence(t, workspace)
	competing := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder", "--input", mutationPath, "--idempotency-key", "claim-mode-mutation-conflict")
	if competing.err == nil || competing.stdout != "" || !strings.Contains(competing.stderr, `"code":"claim-conflict"`) {
		t.Fatalf("mutation acquire under active review Claim = %#v, want conflict", competing)
	}
	assertClaimStateUnchanged(t, workspace, beforeConflict)

	releasePath := filepath.Join(workspace, "claim-mode-release-review.json")
	writeActivationInput(t, releasePath, map[string]any{"claim_id": review["id"], "expected_version": 1, "reason": "completed"})
	released := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", releasePath, "--idempotency-key", "claim-mode-release-review")
	if released.err != nil || released.stderr != "" {
		t.Fatalf("release review Claim = %#v, want stdout-only success", released)
	}

	acquired := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder", "--input", mutationPath, "--idempotency-key", "claim-mode-mutation")
	if acquired.err != nil || acquired.stderr != "" || claimFromResponse(t, acquired.stdout)["mode"] != "mutation" {
		t.Fatalf("mutation acquire after review release = %#v, want active mutation Claim", acquired)
	}
	mutation := claimFromResponse(t, acquired.stdout)
	writeActivationInput(t, releasePath, map[string]any{"claim_id": mutation["id"], "expected_version": 1, "reason": "completed"})
	released = executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", releasePath, "--idempotency-key", "claim-mode-release-mutation")
	if released.err != nil || released.stderr != "" {
		t.Fatalf("release mutation Claim = %#v, want stdout-only success", released)
	}
	reviewAgainPath := filepath.Join(workspace, "claim-mode-review-again.json")
	writeActivationInput(t, reviewAgainPath, map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"worktree_id": worktreeID, "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "mode": "review", "ttl_seconds": 60,
	})
	reviewAgain := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder", "--input", reviewAgainPath, "--idempotency-key", "claim-mode-review-after-mutation")
	if reviewAgain.err != nil || reviewAgain.stderr != "" || claimFromResponse(t, reviewAgain.stdout)["mode"] != "review" {
		t.Fatalf("review acquire after mutation release = %#v, want active review Claim", reviewAgain)
	}
}

func TestClaimAcquireBindingMatrixRejectsWithoutMutation(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := claimEnvelopeForTest(t)
	qualifying := map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 60,
	}
	path := filepath.Join(workspace, "claim-acquire-binding-matrix.json")
	writeActivationInput(t, path, qualifying)
	acquired := executeCommandForTest(workspace,
		"--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder",
		"--input", path, "--idempotency-key", "claim-acquire-binding-qualified",
	)
	if acquired.err != nil || acquired.stderr != "" {
		t.Fatalf("qualifying acquire = %#v, want stdout-only success", acquired)
	}
	_ = claimFromResponse(t, acquired.stdout)

	tests := []struct {
		name     string
		input    map[string]any
		actor    string
		wantCode string
	}{
		{
			name: "missing Issue",
			input: map[string]any{
				"issue_id": "", "envelope_id": envelopeID, "envelope_revision": 1,
				"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 60,
			},
			actor: "actor-builder", wantCode: "invalid-argument",
		},
		{
			name: "mismatched Issue",
			input: map[string]any{
				"issue_id": "ISSUE-mismatched", "envelope_id": envelopeID, "envelope_revision": 1,
				"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 60,
			},
			actor: "actor-builder", wantCode: "invalid-argument",
		},
		{
			name: "missing Envelope",
			input: map[string]any{
				"issue_id": issueID, "envelope_id": "", "envelope_revision": 1,
				"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 60,
			},
			actor: "actor-builder", wantCode: "invalid-argument",
		},
		{
			name: "mismatched Envelope revision",
			input: map[string]any{
				"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 2,
				"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 60,
			},
			actor: "actor-builder", wantCode: "invalid-argument",
		},
		{
			name: "Actor authority",
			input: map[string]any{
				"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
				"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 60,
			},
			actor: "actor-other", wantCode: "authority-denied",
		},
		{
			name: "mismatched worktree",
			input: map[string]any{
				"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
				"worktree_id": "worktree-other", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 60,
			},
			actor: "actor-builder", wantCode: "invalid-argument",
		},
		{
			name: "missing scope",
			input: map[string]any{
				"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
				"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{}}, "ttl_seconds": 60,
			},
			actor: "actor-builder", wantCode: "invalid-argument",
		},
		{
			name: "invalid TTL",
			input: map[string]any{
				"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
				"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 0,
			},
			actor: "actor-builder", wantCode: "invalid-argument",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			before := snapshotWorkspacePersistence(t, workspace)
			writeActivationInput(t, path, testCase.input)
			result := executeCommandForTest(workspace,
				"--project", workspace, "--json", "claim", "acquire", "--actor", testCase.actor,
				"--input", path, "--idempotency-key", "claim-acquire-binding-"+strings.ReplaceAll(testCase.name, " ", "-"),
			)
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"`+testCase.wantCode+`"`) {
				t.Fatalf("%s acquire = %#v, want typed %s rejection", testCase.name, result, testCase.wantCode)
			}
			assertClaimStateUnchanged(t, workspace, before)
		})
	}
}

func stringifyClaimValue(value any) any {
	if number, ok := value.(interface{ String() string }); ok {
		return number.String()
	}
	return value
}

func assertClaimStateUnchanged(t *testing.T, workspace string, want workspacePersistenceSnapshot) {
	t.Helper()
	got := snapshotWorkspacePersistence(t, workspace)
	if !got.equal(want) {
		t.Fatalf("Claim rejection changed persisted projection or journal; events=%d, want %d", got.eventCount(), want.eventCount())
	}
}

func TestClaimConflictsExpiryAndWrongHolderFailWithoutMutation(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := claimEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "claim-acquire-conflict")
	before := journalEventCount(t, workspace)

	conflictPath := filepath.Join(workspace, "claim-conflict.json")
	writeActivationInput(t, conflictPath, map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"internal/claims/generic.go"}}, "ttl_seconds": 60,
	})
	conflict := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder", "--input", conflictPath, "--idempotency-key", "claim-acquire-conflict-2")
	if conflict.err == nil || !strings.Contains(conflict.stderr, `"code":"claim-conflict"`) {
		t.Fatalf("conflict = %#v, want typed claim conflict", conflict)
	}

	wrongHolderPath := filepath.Join(workspace, "claim-wrong-holder.json")
	writeActivationInput(t, wrongHolderPath, map[string]any{"claim_id": claim["id"], "expected_version": 1, "ttl_seconds": 60})
	wrongHolder := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "renew", "--actor", "actor-other", "--input", wrongHolderPath, "--idempotency-key", "claim-renew-other")
	if wrongHolder.err == nil || !strings.Contains(wrongHolder.stderr, `"code":"authority-denied"`) {
		t.Fatalf("wrong holder = %#v, want closed authority rejection", wrongHolder)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open claim database: %v", err)
	}
	if _, err := conn.Exec(`UPDATE claims SET expires_at = '2000-01-01T00:00:00Z' WHERE id = ?`, claim["id"]); err != nil {
		_ = conn.Close()
		t.Fatalf("expire claim: %v", err)
	}
	_ = conn.Close()

	expired := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "renew", "--actor", "actor-builder", "--input", wrongHolderPath, "--idempotency-key", "claim-renew-expired")
	if expired.err == nil || !strings.Contains(expired.stderr, `"code":"claim-expired"`) {
		t.Fatalf("expired = %#v, want typed expiry rejection", expired)
	}
	if after := journalEventCount(t, workspace); after != before {
		t.Fatalf("journal events after closed claim paths = %d, want %d", after, before)
	}
}

func TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "expired-unresolved-claim")
	repository := runRepositoryForTest(t, workspace)
	startPath := filepath.Join(workspace, "expired-unresolved-run.json")
	writeActivationInput(t, startPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace,
		"--project", workspace, "--json", "run", "start", "--actor", "actor-builder",
		"--input", startPath, "--idempotency-key", "expired-unresolved-run",
	)
	if started.err != nil || started.stderr != "" {
		t.Fatalf("start unresolved Run = %#v, want stdout-only success", started)
	}

	expirePath := filepath.Join(workspace, "expired-unresolved-expire.json")
	writeActivationInput(t, expirePath, map[string]any{"claim_id": claim["id"], "expected_version": 1})
	expired := executeCommandForTest(workspace,
		"--project", workspace, "--json", "claim", "expire", "--actor", "actor-builder",
		"--input", expirePath, "--idempotency-key", "expired-unresolved-expire",
	)
	if expired.err != nil || expired.stderr != "" {
		t.Fatalf("expire Claim = %#v, want explicit expiry success", expired)
	}

	before := snapshotWorkspacePersistence(t, workspace)
	reacquirePath := filepath.Join(workspace, "expired-unresolved-reacquire.json")
	writeActivationInput(t, reacquirePath, map[string]any{
		"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1,
		"worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": 60,
	})
	reacquired := executeCommandForTest(workspace,
		"--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder",
		"--input", reacquirePath, "--idempotency-key", "expired-unresolved-reacquire",
	)
	if reacquired.err == nil || !strings.Contains(reacquired.stderr, `"code":"claim-conflict"`) {
		t.Fatalf("reacquire with unresolved executor = %#v, want typed claim conflict", reacquired)
	}
	assertClaimStateUnchanged(t, workspace, before)
}

func TestClaimRenewalAfterRunPreservesLineageAndFencesClaimVersions(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	envelopeInput := envelopeCreateInput(t, workspace, issueID, registration)
	grant := envelopeInput["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "claim.release", "claim.renew", "envelope.activate", "run.start"}
	grant["actor_selector"] = map[string]any{
		"actor_id": "actor-builder",
		"required_capabilities": []string{
			"claim.acquire", "claim.release", "claim.renew", "envelope.activate", "run.start",
		},
	}
	envelopePath := filepath.Join(workspace, "claim-renewal-run-envelope.json")
	writeActivationInput(t, envelopePath, envelopeInput)
	created := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder",
		"--input", envelopePath, "--idempotency-key", "claim-renewal-run-envelope-create",
	)
	if created.err != nil || created.stderr != "" {
		t.Fatalf("create envelope = %#v, want stdout-only success", created)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
	activatePath := filepath.Join(workspace, "claim-renewal-run-envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder",
		"--input", activatePath, "--idempotency-key", "claim-renewal-run-envelope-activate",
	)
	if activated.err != nil || activated.stderr != "" {
		t.Fatalf("activate envelope = %#v, want stdout-only success", activated)
	}

	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "claim-renewal-run-acquire")
	startPath := filepath.Join(workspace, "claim-renewal-run-start.json")
	repository := runRepositoryForTest(t, workspace)
	writeActivationInput(t, startPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace,
		"--project", workspace, "--json", "run", "start", "--actor", "actor-builder",
		"--input", startPath, "--idempotency-key", "claim-renewal-run-start",
	)
	if started.err != nil || started.stderr != "" {
		t.Fatalf("start Run = %#v, want stdout-only success", started)
	}
	run := runFromResponse(t, started.stdout)
	if run["id"] == "" || run["issue_id"] != issueID || run["envelope_id"] != envelopeID ||
		stringifyClaimValue(run["envelope_revision"]) != "1" || run["claim_id"] != claim["id"] ||
		stringifyClaimValue(run["claim_version"]) != "1" || run["actor_id"] != "actor-builder" ||
		!reflect.DeepEqual(run["executor"], map[string]any{"id": "executor-builder", "source": "activated-harness"}) ||
		!reflect.DeepEqual(run["start_repository"], repository) {
		t.Fatalf("started Run = %#v, want accepted Claim binding and immutable lineage", run)
	}
	lineage := map[string]any{}
	for _, field := range []string{
		"id", "issue_id", "envelope_id", "envelope_revision", "claim_id", "claim_version",
		"actor_id", "executor", "obligation_id", "start_repository", "predecessor_run_id", "successor_run_id",
	} {
		lineage[field] = run[field]
	}

	renewPath := filepath.Join(workspace, "claim-renewal-run-renew.json")
	writeActivationInput(t, renewPath, map[string]any{"claim_id": claim["id"], "expected_version": 1, "ttl_seconds": 120})
	renewed := executeCommandForTest(workspace,
		"--project", workspace, "--json", "claim", "renew", "--actor", "actor-builder",
		"--input", renewPath, "--idempotency-key", "claim-renewal-run-renew",
	)
	if renewed.err != nil || renewed.stderr != "" || stringifyClaimValue(claimFromResponse(t, renewed.stdout)["version"]) != "2" {
		t.Fatalf("renew Claim = %#v, want accepted version 2", renewed)
	}

	assertRenewedClaimRunLineage(t, workspace, run["id"].(string), lineage)
	staleReleasePath := filepath.Join(workspace, "claim-renewal-run-release-stale.json")
	writeActivationInput(t, staleReleasePath, map[string]any{"claim_id": claim["id"], "expected_version": 1, "reason": "completed"})
	beforeStale := snapshotWorkspacePersistence(t, workspace)
	stale := executeCommandForTest(workspace,
		"--project", workspace, "--json", "claim", "release", "--actor", "actor-builder",
		"--input", staleReleasePath, "--idempotency-key", "claim-renewal-run-release-stale",
	)
	if stale.err == nil || stale.stdout != "" || !strings.Contains(stale.stderr, `"code":"claim-stale"`) {
		t.Fatalf("stale Claim release = %#v, want typed claim-stale rejection", stale)
	}
	assertClaimStateUnchanged(t, workspace, beforeStale)
	assertRenewedClaimRunLineage(t, workspace, run["id"].(string), lineage)

	currentReleasePath := filepath.Join(workspace, "claim-renewal-run-release-current.json")
	writeActivationInput(t, currentReleasePath, map[string]any{"claim_id": claim["id"], "expected_version": 2, "reason": "completed"})
	current := executeCommandForTest(workspace,
		"--project", workspace, "--json", "claim", "release", "--actor", "actor-builder",
		"--input", currentReleasePath, "--idempotency-key", "claim-renewal-run-release-current",
	)
	if current.err != nil || current.stderr != "" {
		t.Fatalf("current Claim release = %#v, want stdout-only success", current)
	}
	currentClaim := claimFromResponse(t, current.stdout)
	if stringifyClaimValue(currentClaim["version"]) != "3" || currentClaim["status"] != "released" || currentClaim["release_reason"] != "completed" {
		t.Fatalf("current Claim release = %#v, want released Claim version 3", currentClaim)
	}
	assertRenewedClaimRunLineage(t, workspace, run["id"].(string), lineage)
}

func assertRenewedClaimRunLineage(t *testing.T, workspace, runID string, want map[string]any) {
	t.Helper()
	got := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "get", runID)
	if got.err != nil || got.stderr != "" {
		t.Fatalf("get renewed-Claim Run = %#v, want stdout-only success", got)
	}
	run := runFromResponse(t, got.stdout)
	for field, wantValue := range want {
		if !reflect.DeepEqual(run[field], wantValue) {
			t.Fatalf("Run.%s after Claim renewal = %#v, want immutable %#v", field, run[field], wantValue)
		}
	}
}

func TestClaimRenewReleaseVersionAndReplayAreExact(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := claimEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "claim-acquire-lifecycle")
	renewPath := filepath.Join(workspace, "claim-renew.json")
	writeActivationInput(t, renewPath, map[string]any{"claim_id": claim["id"], "expected_version": 1, "ttl_seconds": 120})
	renewed := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "renew", "--actor", "actor-builder", "--input", renewPath, "--idempotency-key", "claim-renew-1")
	if renewed.err != nil || renewed.stderr != "" {
		t.Fatalf("renew = %#v, want stdout-only success", renewed)
	}
	renewedClaim := claimFromResponse(t, renewed.stdout)
	if got := stringifyClaimValue(renewedClaim["version"]); got != "2" {
		t.Fatalf("renewed version = %#v, want 2", got)
	}
	beforeStaleRenew := snapshotWorkspacePersistence(t, workspace)
	staleRenewPath := filepath.Join(workspace, "claim-renew-stale.json")
	writeActivationInput(t, staleRenewPath, map[string]any{"claim_id": claim["id"], "expected_version": 1, "ttl_seconds": 120})
	staleRenew := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "renew", "--actor", "actor-builder", "--input", staleRenewPath, "--idempotency-key", "claim-renew-stale")
	if staleRenew.err == nil || staleRenew.stdout != "" || !strings.Contains(staleRenew.stderr, `"code":"claim-stale"`) {
		t.Fatalf("stale renew = %#v, want typed stale rejection", staleRenew)
	}
	assertClaimStateUnchanged(t, workspace, beforeStaleRenew)
	replay := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "renew", "--actor", "actor-builder", "--input", renewPath, "--idempotency-key", "claim-renew-1")
	if replay.err != nil || replay.stderr != "" || claimResponseReplayed(t, replay.stdout) != true {
		t.Fatalf("renew replay = %#v, want exact idempotent replay", replay)
	}

	releasePath := filepath.Join(workspace, "claim-release.json")
	writeActivationInput(t, releasePath, map[string]any{"claim_id": claim["id"], "expected_version": 2, "reason": "completed"})
	beforeStaleRelease := snapshotWorkspacePersistence(t, workspace)
	staleReleasePath := filepath.Join(workspace, "claim-release-stale.json")
	writeActivationInput(t, staleReleasePath, map[string]any{"claim_id": claim["id"], "expected_version": 1, "reason": "completed"})
	staleRelease := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", staleReleasePath, "--idempotency-key", "claim-release-stale")
	if staleRelease.err == nil || staleRelease.stdout != "" || !strings.Contains(staleRelease.stderr, `"code":"claim-stale"`) {
		t.Fatalf("stale release = %#v, want typed stale rejection", staleRelease)
	}
	assertClaimStateUnchanged(t, workspace, beforeStaleRelease)
	beforeRelease := journalEventCount(t, workspace)
	released := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", releasePath, "--idempotency-key", "claim-release-1")
	if released.err != nil || released.stderr != "" {
		t.Fatalf("release = %#v, want stdout-only success", released)
	}
	releasedClaim := claimFromResponse(t, released.stdout)
	if got := stringifyClaimValue(releasedClaim["version"]); got != "3" || releasedClaim["status"] != "released" || releasedClaim["release_reason"] != "completed" {
		t.Fatalf("released claim = %#v, want released version 3 completed", releasedClaim)
	}
	releaseReplay := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", "actor-builder", "--input", releasePath, "--idempotency-key", "claim-release-1")
	if releaseReplay.err != nil || releaseReplay.stderr != "" || !claimResponseReplayed(t, releaseReplay.stdout) {
		t.Fatalf("release replay = %#v, want exact idempotent replay", releaseReplay)
	}
	if after := journalEventCount(t, workspace); after != beforeRelease+1 {
		t.Fatalf("release/replay events = %d, want %d", after, beforeRelease+1)
	}
}

func claimEnvelopeForTest(t *testing.T) (string, string, string) {
	t.Helper()
	workspace, values := cloneReusableWorkspaceFixture(t, "claim-envelope", func() (string, []string) {
		workspace, issueID, envelopeID := buildClaimEnvelopeForTest(t)
		return workspace, []string{issueID, envelopeID}
	})
	return workspace, values[0], values[1]
}

func buildClaimEnvelopeForTest(t *testing.T) (string, string, string) {
	t.Helper()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	envelopeInput := envelopeCreateInput(t, workspace, issueID, registration)
	grant := envelopeInput["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "claim.release", "claim.renew", "envelope.activate"}
	grant["actor_selector"] = map[string]any{"actor_id": "actor-builder", "required_capabilities": []string{"claim.acquire", "claim.release", "claim.renew", "envelope.activate"}}
	envelopePath := filepath.Join(workspace, "generic-claim-envelope.json")
	writeActivationInput(t, envelopePath, envelopeInput)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", envelopePath, "--idempotency-key", "generic-claim-envelope-create")
	if created.err != nil {
		t.Fatalf("create claim envelope: %v; stderr=%q", created.err, created.stderr)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
	activatePath := filepath.Join(workspace, "generic-claim-envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder", "--input", activatePath, "--idempotency-key", "generic-claim-envelope-activate")
	if activated.err != nil {
		t.Fatalf("activate claim envelope: %v; stderr=%q", activated.err, activated.stderr)
	}
	return workspace, issueID, envelopeID
}

func acquireGenericClaimForTest(t *testing.T, workspace, issueID, envelopeID, key string) (map[string]any, string) {
	return acquireGenericClaimForTestWithTTL(t, workspace, issueID, envelopeID, key, 60)
}

func acquireGenericClaimForTestWithTTL(t *testing.T, workspace, issueID, envelopeID, key string, ttlSeconds int) (map[string]any, string) {
	return acquireGenericClaimForTestWithMode(t, workspace, issueID, envelopeID, key, ttlSeconds, "mutation")
}

func acquireGenericClaimForTestWithMode(t *testing.T, workspace, issueID, envelopeID, key string, ttlSeconds int, mode string) (map[string]any, string) {
	t.Helper()
	path := filepath.Join(workspace, key+".json")
	worktreeID := testEnvelopeRepository(t, workspace)["worktree_id"]
	writeActivationInput(t, path, map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "worktree_id": worktreeID, "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "mode": mode, "ttl_seconds": ttlSeconds})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", "actor-builder", "--input", path, "--idempotency-key", key)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("acquire = %#v, want stdout-only success", result)
	}
	return claimFromResponse(t, result.stdout), path
}

func claimFromResponse(t *testing.T, output string) map[string]any {
	t.Helper()
	return decodeCommandJSONStructure(t, output).(map[string]any)["result"].(map[string]any)["claim"].(map[string]any)
}

func claimResponseReplayed(t *testing.T, output string) bool {
	t.Helper()
	value := decodeCommandJSONStructure(t, output).(map[string]any)["result"].(map[string]any)["replayed"]
	replayed, ok := value.(bool)
	if !ok {
		t.Fatalf("replayed = %#v, want bool", value)
	}
	return replayed
}

func TestLegacyIssueClaimSpellingsAreNotAliases(t *testing.T) {
	t.Parallel()
	for _, legacy := range []string{"issue claim", "issue unclaim"} {
		for _, leaf := range liveExecutableLeafCommands(newRootCommandAt(t.TempDir())) {
			if leaf == legacy {
				t.Fatalf("legacy leaf %q remains registered", legacy)
			}
		}
	}
}

func TestClaimResponseCarriesOnlyContractFields(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := claimEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "claim-contract-fields")
	encoded, err := json.Marshal(claim)
	if err != nil || !strings.Contains(string(encoded), `"worktree_id":"worktree-1"`) {
		t.Fatalf("claim JSON = %q, %v; want contract binding", encoded, err)
	}
}
