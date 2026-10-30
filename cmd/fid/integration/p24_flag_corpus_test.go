//go:build integration

package integration

import (
	"testing"
	"time"
)

var p24FlagFields = map[string]string{
	"label.add": "--colour", "issue.ac.add": "--text",
	"issue.update": "--title", "issue.label.add": "label",
	"issue.label.remove": "label", "issue.transition": "--to",
	"link.add": "relation", "link.remove": "relation",
	"claim.renew": "ttl_seconds", "claim.release": "expected_version",
}

func TestP24FlagAndClaimKeyedLeafCorpus(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	a := p22NewFIDRecoveryAdapterAt(t, binary, workspace, "agent", 3600, []string{}, p22FinalActions)
	p24StartKeyedCapture(workspace)
	base := []string{"--project", workspace, "--json"}
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, workspace, append(append([]string(nil), base...), args...)...)
	}
	call("label", "add", "p24-label", "--colour", "red", "--actor", "agent:developer", "--idempotency-key", "p24-label-add")
	call("label", "add", "p24-other-label", "--colour", "blue", "--actor", "agent:developer", "--idempotency-key", "p24-other-label-add")
	call("issue", "transition", a.issueID, "--to", "in_progress", "--expected-revision", "1", "--actor", "agent:"+a.actor, "--idempotency-key", "p24-issue-transition")
	call("issue", "ac", "add", a.issueID, "--text", "first criterion", "--ordinal", "1", "--actor", "agent:developer", "--idempotency-key", "p24-ac-add")
	call("issue", "update", a.issueID, "--title", "Phase 24 updated", "--expected-revision", "3", "--actor", "agent:developer", "--idempotency-key", "p24-issue-update")
	call("issue", "label", "add", a.issueID, "p24-label", "--actor", "agent:developer", "--idempotency-key", "p24-issue-label-add")
	call("issue", "label", "remove", a.issueID, "p24-label", "--actor", "agent:developer", "--idempotency-key", "p24-issue-label-remove")
	other := call("issue", "create", "--kind", "task", "--title", "Phase 24 link peer", "--actor", "agent:developer", "--idempotency-key", "p24-link-peer")
	peerID := p21PathString(t, other, "result", "issue", "id")
	call("link", "add", a.issueID, "blocks", peerID, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "p24-link-add")
	call("link", "remove", a.issueID, "blocks", peerID, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "p24-link-remove")
	renew := p23RunInput(t, a, "claim", "renew", a.actor, "p24-claim-renew", map[string]any{"claim_id": a.claimID, "expected_version": 1, "ttl_seconds": 3600})
	if renew.ExitCode != 0 || renew.Stderr != "" {
		t.Fatalf("claim.renew = %d %q", renew.ExitCode, renew.Stderr)
	}
	run := call("run", "get", a.runID)
	observed := p23RunInput(t, a, "run", "observe", a.actor, "p24-flag-observe", map[string]any{"run_id": a.runID, "expected_version": p21Path(t, run, "result", "run", "version"), "liveness": map[string]any{"executor_id": p21PathString(t, run, "result", "run", "executor", "id"), "state": "terminated", "source_ref": "process:p24", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
	p23AssertReviewSuccess(t, "observe Run", observed)
	observation := p23Envelope(t, observed.Stdout)
	closed := p23RunInput(t, a, "run", "close", a.actor, "p24-flag-close", map[string]any{"run_id": a.runID, "expected_version": p21Path(t, observation, "result", "run", "version"), "observation_id": p21PathString(t, observation, "result", "observation", "id"), "expected_repository": p24ExpectedRepository(a.repository)})
	p23AssertReviewSuccess(t, "close Run", closed)
	released := p23RunInput(t, a, "claim", "release", a.actor, "p24-flag-release", map[string]any{"claim_id": a.claimID, "expected_version": 2, "reason": "completed"})
	p23AssertReviewSuccess(t, "release Claim", released)
	captured := p24EndKeyedCapture(workspace)
	for operation, name := range p24FlagFields {
		entry, found := captured[operation]
		if !found {
			t.Fatalf("missing %s in flag/claim capture: %v", operation, p24CapturedNames(captured))
		}
		p24ExerciseCapturedLeaf(t, workspace, binary, operation, name, entry)
	}
	for operation, second := range map[string]string{"issue.update": "--expected-revision", "label.add": "--description", "link.add": "to-id"} {
		entry := captured[operation]
		p24AssertFamilyConflict(t, workspace, binary, operation, entry, []string{p24FlagFields[operation], second})
		p24AssertFamilyConflict(t, workspace, binary, operation, entry, []string{"--actor"})
	}
	p24AssertRebuildMapParity(t, workspace, binary)
}
