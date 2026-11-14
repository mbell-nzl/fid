//go:build integration

package integration

import (
	"testing"
)

var p24FoundationFields = map[string]string{
	"issue.create": "--title", "activation.prepare": "kit_identity",
	"activation.resource": "evidence_ref", "activation.component": "content",
	"activation.commit": "expected_version", "activation.finalize": "expected_version",
	"gate.register": "definition_version", "envelope.create": "outcome",
	"envelope.activate": "expected_revision", "claim.acquire": "ttl_seconds",
	"run.start": "claim_version",
}

func TestP24FoundationKeyedLeafConflictAndCanonicalRetryCorpus(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p24StartKeyedCapture(workspace)
	_ = p22NewFIDRecoveryAdapterAt(t, binary, workspace, "agent", 3600, []string{}, nil)
	captured := p24EndKeyedCapture(workspace)
	for operation, name := range p24FoundationFields {
		entry, ok := captured[operation]
		if !ok {
			t.Fatalf("fixture did not capture keyed %s; captured = %v", operation, p24CapturedNames(captured))
		}
		p24ExerciseCapturedLeaf(t, workspace, binary, operation, name, entry)
	}
	for operation, second := range map[string]string{"activation.prepare": "resources", "envelope.create": "non_goals", "run.start": "executor", "claim.acquire": "mode", "gate.register": "definition_id"} {
		entry := captured[operation]
		p24AssertFamilyConflict(t, workspace, binary, operation, entry, []string{p24FoundationFields[operation], second})
		p24AssertFamilyConflict(t, workspace, binary, operation, entry, []string{"--actor"})
	}
	p24AssertFamilyConflict(t, workspace, binary, "envelope.create", captured["envelope.create"], []string{"protected_decisions"})
	p24AssertRebuildMapParity(t, workspace, binary)
}

var p24PositionalNames = map[string][]string{"check.execute": {"names"},
	"task.answer":   {"task", "text"},
	"task.takeover": {"task"},
	"issue.update":  {"id"}, "issue.transition": {"id"}, "issue.ac.add": {"issue-id"},
	"issue.label.add": {"issue-id", "label"}, "issue.label.remove": {"issue-id", "label"},
	"label.add": {"name"}, "link.add": {"from-id", "relation", "to-id"}, "link.remove": {"from-id", "relation", "to-id"},
}
