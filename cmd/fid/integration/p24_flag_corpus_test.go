//go:build integration

package integration

import (
	"testing"
)

var p24FlagFields = map[string]string{
	"label.add": "--colour", "issue.ac.add": "--text", "issue.update": "--title", "issue.label.add": "label", "issue.label.remove": "label", "issue.transition": "--to", "link.add": "relation", "link.remove": "relation",
}

func leanFlagFixture(t *testing.T, binary, workspace string) string {
	t.Helper()
	closureRun(t, workspace, "git", "init")
	p21Run(t, binary, workspace, "init", "--prefix", "FID", "--json")
	parent := p21Run(t, binary, workspace, "--json", "issue", "create", "--kind", "task", "--title", "Parent", "--actor", "agent:developer", "--idempotency-key", "flag-parent")
	child := p21Run(t, binary, workspace, "--json", "issue", "create", "--kind", "subtask", "--parent", p21PathString(t, parent, "result", "issue", "id"), "--title", "Flag corpus", "--actor", "agent:developer", "--idempotency-key", "flag-child")
	return p21PathString(t, child, "result", "issue", "id")
}

func TestP24FlagKeyedLeafCorpus(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	issueID := leanFlagFixture(t, binary, workspace)
	p24StartKeyedCapture(workspace)
	base := []string{"--project", workspace, "--json"}
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, workspace, append(append([]string(nil), base...), args...)...)
	}
	call("label", "add", "p24-label", "--colour", "red", "--actor", "agent:developer", "--idempotency-key", "p24-label-add")
	call("label", "add", "p24-other-label", "--colour", "blue", "--actor", "agent:developer", "--idempotency-key", "p24-other-label-add")
	call("issue", "transition", issueID, "--to", "in_progress", "--expected-revision", "1", "--actor", "agent:developer", "--idempotency-key", "p24-issue-transition")
	call("issue", "ac", "add", issueID, "--text", "first criterion", "--ordinal", "1", "--actor", "agent:developer", "--idempotency-key", "p24-ac-add")
	call("issue", "update", issueID, "--title", "Phase 24 updated", "--expected-revision", "3", "--actor", "agent:developer", "--idempotency-key", "p24-issue-update")
	call("issue", "label", "add", issueID, "p24-label", "--actor", "agent:developer", "--idempotency-key", "p24-issue-label-add")
	call("issue", "label", "remove", issueID, "p24-label", "--actor", "agent:developer", "--idempotency-key", "p24-issue-label-remove")
	other := call("issue", "create", "--kind", "task", "--title", "Phase 24 link peer", "--actor", "agent:developer", "--idempotency-key", "p24-link-peer")
	peerID := p21PathString(t, other, "result", "issue", "id")
	call("link", "add", issueID, "blocks", peerID, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "p24-link-add")
	call("link", "remove", issueID, "blocks", peerID, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "p24-link-remove")
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
	assertLeanIntegrationJournal(t, workspace)
	p24AssertRebuildMapParity(t, workspace, binary)
}
