package main

import (
	"encoding/json"
	"path/filepath"
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
	t.Helper()
	path := filepath.Join(workspace, key+".json")
	writeActivationInput(t, path, map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "worktree_id": "worktree-1", "scope": map[string]any{"paths": []string{"cmd/fid/claim.go"}}, "ttl_seconds": ttlSeconds})
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
