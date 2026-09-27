//go:build integration

package integration

import (
	"encoding/json"
	"testing"
)

func TestP23RequirementReferenceAndLinkRealProcess(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "decision.record-reference", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "requirement.link", "requirement.record-reference", "run.close", "run.observe", "run.observe-budget", "run.start", "stop.resolve"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, p22ExecutionFixture{})
	source := p23DecisionSource(a, true)
	seed := map[string]any{"id": "PRD-TRACE-001", "source": source, "status": "active"}
	first := p23RunInput(t, a, "requirement", "record-reference", a.actor, "task12-first", seed)
	if first.ExitCode != 0 || first.Stderr != "" {
		t.Fatalf("first reference = %#v", first)
	}
	var firstEnvelope map[string]any
	if err := json.Unmarshal([]byte(first.Stdout), &firstEnvelope); err != nil {
		t.Fatal(err)
	}
	oldID := p21PathString(t, firstEnvelope, "result", "requirement", "id")
	secondSource := p23DecisionSource(a, true)
	secondSource["display_label"] = "Second reference to the same file"
	second := p23RunInput(t, a, "requirement", "record-reference", a.actor, "task12-second", map[string]any{"id": "PRD-TRACE-001", "source": secondSource, "status": "active"})
	if second.ExitCode != 0 || second.Stderr != "" {
		t.Fatalf("second reference = %#v", second)
	}
	var secondEnvelope map[string]any
	if err := json.Unmarshal([]byte(second.Stdout), &secondEnvelope); err != nil {
		t.Fatal(err)
	}
	newID := p21PathString(t, secondEnvelope, "result", "requirement", "id")
	if oldID == newID {
		t.Fatal("re-recording reused the old row")
	}
	old := a.p22Run(t, "requirement", "get", oldID)
	if old.ExitCode != 0 || old.Stderr != "" {
		t.Fatalf("old reference = %#v", old)
	}
	var oldEnvelope map[string]any
	if err := json.Unmarshal([]byte(old.Stdout), &oldEnvelope); err != nil {
		t.Fatal(err)
	}
	if p21PathString(t, oldEnvelope, "result", "requirement", "status") != "superseded" || p21PathString(t, oldEnvelope, "result", "requirement", "superseded_by_id") != newID {
		t.Fatalf("prior reference lost supersession pointer: %s", old.Stdout)
	}
	linked := p23RunInput(t, a, "requirement", "link", a.actor, "task12-link", map[string]any{"requirement_id": "PRD-TRACE-001", "issue_id": a.issueID, "relation": "verifies"})
	if linked.ExitCode != 0 || linked.Stderr != "" {
		t.Fatalf("link = %#v", linked)
	}
}
