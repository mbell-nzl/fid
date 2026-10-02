//go:build integration

package integration

import (
	"crypto/sha256"
	"encoding/json"
	"testing"
)

func TestP23RequirementNewShapeJournalRoundTrip(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "decision.record-reference", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.observe-budget", "run.start", "stop.resolve"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, p22ExecutionFixture{})
	request := map[string]any{"id": "PRD-ROUNDTRIP-001", "source": p23DecisionSource(a, true), "status": "active"}
	first := p23RunInput(t, a, "requirement", "record-reference", a.actor, "task20-roundtrip-first", request)
	if first.ExitCode != 0 || first.Stderr != "" {
		t.Fatalf("first reference: %#v", first)
	}
	var firstResult map[string]any
	if err := json.Unmarshal([]byte(first.Stdout), &firstResult); err != nil {
		t.Fatal(err)
	}
	firstID := p21PathString(t, firstResult, "result", "requirement", "id")
	secondSource := p23DecisionSource(a, true)
	secondSource["display_label"] = "newer source label"
	second := p23RunInput(t, a, "requirement", "record-reference", a.actor, "task20-roundtrip-second", map[string]any{"id": "PRD-ROUNDTRIP-001", "source": secondSource, "status": "active"})
	if second.ExitCode != 0 || second.Stderr != "" {
		t.Fatalf("superseding reference: %#v", second)
	}
	var secondResult map[string]any
	if err := json.Unmarshal([]byte(second.Stdout), &secondResult); err != nil {
		t.Fatal(err)
	}
	secondID := p21PathString(t, secondResult, "result", "requirement", "id")
	if firstID == secondID {
		t.Fatal("superseding reference reused prior row ID")
	}
	for _, relation := range []string{"implements", "verifies"} {
		linked := p23RunInput(t, a, "requirement", "link", a.actor, "task20-roundtrip-"+relation, map[string]any{"requirement_id": "PRD-ROUNDTRIP-001", "issue_id": a.issueID, "relation": relation})
		if linked.ExitCode != 0 || linked.Stderr != "" {
			t.Fatalf("link %s: %#v", relation, linked)
		}
	}
	reads := [][]string{{"requirement", "get", "PRD-ROUNDTRIP-001"}, {"requirement", "get", firstID}, {"requirement", "get", secondID}, {"requirement", "list"}, {"requirement", "list", "--issue", a.issueID}}
	before := make([]string, len(reads))
	for i, args := range reads {
		result := a.p22Run(t, args...)
		if result.ExitCode != 0 || result.Stderr != "" {
			t.Fatalf("read before rebuild %v: %#v", args, result)
		}
		before[i] = result.Stdout
	}
	digest := sha256.Sum256(p21JournalBytes(t, a.workspace))
	rebuild := a.p22Run(t, "rebuild")
	if rebuild.ExitCode != 0 || rebuild.Stderr != "" {
		t.Fatalf("new-shape rebuild: %#v", rebuild)
	}
	if got := sha256.Sum256(p21JournalBytes(t, a.workspace)); got != digest {
		t.Fatalf("rebuild changed source journal digest: %x -> %x", digest, got)
	}
	for i, args := range reads {
		result := a.p22Run(t, args...)
		if result.ExitCode != 0 || result.Stderr != "" || result.Stdout != before[i] {
			t.Fatalf("read after rebuild %v changed: before %q, after %#v", args, before[i], result)
		}
	}
	events := p21JournalEventCount(t, a.workspace)
	retry := p23RunInput(t, a, "requirement", "record-reference", a.actor, "task20-roundtrip-first", request)
	if retry.ExitCode != 0 || retry.Stderr != "" || p21JournalEventCount(t, a.workspace) != events {
		t.Fatalf("idempotent retry changed journal or failed: %#v", retry)
	}
	var replay map[string]any
	if err := json.Unmarshal([]byte(retry.Stdout), &replay); err != nil {
		t.Fatal(err)
	}
	if p21PathString(t, replay, "result", "requirement", "requirement_id") != "PRD-ROUNDTRIP-001" || p21PathString(t, replay, "result", "requirement", "id") != firstID || p21Path(t, replay, "result", "replayed") != true {
		t.Fatalf("idempotent retry lost stable identity: %s", retry.Stdout)
	}
	if got := sha256.Sum256(p21JournalBytes(t, a.workspace)); got != digest {
		t.Fatalf("idempotent retry changed source journal digest: %x -> %x", digest, got)
	}
}
