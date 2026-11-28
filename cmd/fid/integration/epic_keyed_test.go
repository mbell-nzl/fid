//go:build integration

package integration

import "testing"

var p24EpicFields = map[string]string{"epic.approve": "--reason", "epic.validation.add": "--claim", "epic.validation.update": "--claim", "epic.validation.remove": "id", "epic.validation.pass": "--reason", "epic.validation.verify": "epic"}

func TestP24EpicKeyedLeafCorpus(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepoWithArtifacts(t, true, true)
	p24StartKeyedCapture(root)
	verified := p24EpicCommands(t, root, binary, "test")
	c := loopContract{t: t, decl: readLoopDeclaration(t)}
	c.observe("epic-validation-verify", verified, 0)
	captured := p24EndKeyedCapture(root)
	for operation, field := range p24EpicFields {
		entry, ok := captured[operation]
		if !ok {
			t.Fatalf("missing executable %s", operation)
		}
		p24ExerciseCapturedLeaf(t, root, binary, operation, field, entry)
	}
	p24AssertRebuildMapParity(t, root, binary)
}

func p24EpicCommands(t *testing.T, root, binary, check string) map[string]any {
	t.Helper()
	created := p21Run(t, binary, root, "--json", "issue", "create", "--kind", "epic", "--scheduling", "prioritised", "--title", "Keyed epic", "--actor", "human:mark", "--idempotency-key", "keyed-epic")
	epic := p21PathString(t, created, "result", "issue", "id")
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, root, append([]string{"--json"}, args...)...)
	}
	call("epic", "validation", "add", epic, "--claim", "Review claim", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "add-review")
	call("epic", "validation", "update", epic, "V1", "--claim", "Updated review claim", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "update")
	call("epic", "validation", "pass", epic, "V1", "--actor", "human:mark", "--reason", "Reviewed", "--idempotency-key", "pass")
	call("epic", "validation", "remove", epic, "V1", "--actor", "agent:builder", "--idempotency-key", "remove")
	call("epic", "validation", "add", epic, "--claim", "Check claim", "--verifier", "check:"+check, "--actor", "agent:builder", "--idempotency-key", "add-check")
	verified := call("epic", "validation", "verify", epic, "--actor", "agent:builder", "--idempotency-key", "verify")
	if out, stderr, err := approvalCorpusCommand(t, root, binary, "--json", "epic", "approve", epic, "--actor", "human:mark", "--reason", "Approve plan", "--idempotency-key", "approve"); err != nil {
		t.Fatalf("approve=%v %s %s", err, out, stderr)
	}
	return verified
}
