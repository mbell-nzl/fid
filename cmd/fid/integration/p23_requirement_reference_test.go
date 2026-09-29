//go:build integration

package integration

import (
	"encoding/json"
	"testing"
)

func TestP23RequirementReferenceAndLinkRealProcess(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "decision.record-reference", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.observe-budget", "run.start", "stop.resolve"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, p22ExecutionFixture{})
	a.task13ReplayCorpus = true
	source := p23DecisionSource(a, true)
	seed := map[string]any{"id": "PRD-TRACE-001", "source": source, "status": "active"}
	first := p23RunInput(t, a, "requirement", "record-reference", a.actor, "task12-first", seed)
	if first.ExitCode != 0 || first.Stderr != "" {
		t.Fatalf("first reference = %#v", first)
	}
	p23AssertActivitySubject(t, a.workspace, first, "requirement-referenced", "requirement", "PRD-TRACE-001", 1)
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
	p23AssertActivitySubject(t, a.workspace, second, "requirement-referenced", "requirement", "PRD-TRACE-001", 1)
	var secondEnvelope map[string]any
	if err := json.Unmarshal([]byte(second.Stdout), &secondEnvelope); err != nil {
		t.Fatal(err)
	}
	newID := p21PathString(t, secondEnvelope, "result", "requirement", "id")
	if oldID == newID {
		t.Fatal("re-recording reused the old row")
	}
	beforeMismatch := p21JournalEventCount(t, a.workspace)
	mismatch := p23RunInput(t, a, "requirement", "record-reference", a.actor, "task12-mismatch", map[string]any{"id": "PRD-TRACE-001", "source": p23DecisionSource(a, false), "status": "active"})
	assertP22ProcessFailure(t, "requirement.record-reference", "source-digest-mismatch", 4, mismatch)
	if mismatch.ExitCode != 4 || mismatch.Stdout != "" || !p23HasJSON(mismatch.Stderr, `"code":"source-digest-mismatch"`) || p21JournalEventCount(t, a.workspace) != beforeMismatch {
		t.Fatalf("digest mismatch = %#v; events %d -> %d", mismatch, beforeMismatch, p21JournalEventCount(t, a.workspace))
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
	active := a.p22Run(t, "requirement", "list", "--status", "active")
	if active.ExitCode != 0 || active.Stderr != "" {
		t.Fatalf("active references = %#v", active)
	}
	var activeEnvelope map[string]any
	if err := json.Unmarshal([]byte(active.Stdout), &activeEnvelope); err != nil {
		t.Fatal(err)
	}
	items := p21Path(t, activeEnvelope, "result", "items").([]any)
	if len(items) != 1 || items[0].(map[string]any)["id"] != newID {
		t.Fatalf("active references = %#v, want only %s", items, newID)
	}
	for _, relation := range []string{"implements", "verifies", "relates-to"} {
		linked := p23RunInput(t, a, "requirement", "link", a.actor, "task12-link-"+relation, map[string]any{"requirement_id": "PRD-TRACE-001", "issue_id": a.issueID, "relation": relation})
		if linked.ExitCode != 0 || linked.Stderr != "" {
			t.Fatalf("link %s = %#v", relation, linked)
		}
		p23AssertActivitySubject(t, a.workspace, linked, "requirement-linked", "requirement", "PRD-TRACE-001", 1)
		var envelope map[string]any
		if err := json.Unmarshal([]byte(linked.Stdout), &envelope); err != nil {
			t.Fatal(err)
		}
		if p21PathString(t, envelope, "result", "link", "relation") != relation {
			t.Fatalf("link relation = %s, want %s", linked.Stdout, relation)
		}
	}
	readback := a.p22Run(t, "requirement", "get", "PRD-TRACE-001")
	if readback.ExitCode != 0 || readback.Stderr != "" {
		t.Fatalf("linked requirement readback = %#v", readback)
	}
	var readEnvelope map[string]any
	if err := json.Unmarshal([]byte(readback.Stdout), &readEnvelope); err != nil {
		t.Fatal(err)
	}
	storedLinks := p21Path(t, readEnvelope, "result", "links").([]any)
	if len(storedLinks) != 3 {
		t.Fatalf("stored traceability links = %#v, want all three relations", storedLinks)
	}
	seen := map[string]bool{}
	for _, link := range storedLinks {
		value := link.(map[string]any)
		if value["requirement_id"] != "PRD-TRACE-001" || value["issue_id"] != a.issueID {
			t.Fatalf("stored traceability link = %#v", value)
		}
		seen[value["relation"].(string)] = true
	}
	for _, relation := range []string{"implements", "verifies", "relates-to"} {
		if !seen[relation] {
			t.Fatalf("stored links omit %s: %#v", relation, storedLinks)
		}
	}
	beforeUnknown := p21JournalEventCount(t, a.workspace)
	unknownID := p23RunInput(t, a, "requirement", "link", a.actor, "task12-unknown-id", map[string]any{"requirement_id": "PRD-UNKNOWN", "issue_id": a.issueID, "relation": "verifies"})
	assertP22ProcessFailure(t, "requirement.link", "requirement-not-found", 3, unknownID)
	unknownRelation := p23RunInput(t, a, "requirement", "link", a.actor, "task12-unknown-relation", map[string]any{"requirement_id": "PRD-TRACE-001", "issue_id": a.issueID, "relation": "unknown"})
	if unknownID.ExitCode == 0 || !p23HasJSON(unknownID.Stderr, `"code":"requirement-not-found"`) || unknownRelation.ExitCode == 0 || p21JournalEventCount(t, a.workspace) != beforeUnknown {
		t.Fatalf("unknown links = %#v, %#v; events %d -> %d", unknownID, unknownRelation, beforeUnknown, p21JournalEventCount(t, a.workspace))
	}
	for _, args := range [][]string{
		{"link", "add", a.issueID, "implements", "PRD-TRACE-001", "--from-type", "issue", "--to-type", "requirement", "--actor", "agent:developer", "--idempotency-key", "task12-legacy-endpoint"},
		{"link", "add", a.issueID, "implements", "FID-unknown", "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "task12-legacy-implements"},
		{"link", "add", a.issueID, "implemented_by", "FID-unknown", "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "task12-legacy-relation"},
	} {
		before := p21JournalEventCount(t, a.workspace)
		rejected := a.p22Run(t, args...)
		if rejected.ExitCode == 0 || rejected.Stdout != "" || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("removed link form %v = %#v; events %d -> %d", args, rejected, before, p21JournalEventCount(t, a.workspace))
		}
	}
}
