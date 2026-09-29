package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPhase23HumanLeafTextMatchesRecordedFacts(t *testing.T) {
	t.Parallel()
	source := map[string]any{"display_label": "contract section", "path": "docs/contract.md", "digest": "sha256:abc"}
	items := map[string]map[string]any{
		"review":        {"id": "REV-one", "issue_id": "ISS-one", "envelope_id": "ENV-one", "envelope_revision": 2, "run_id": "RUN-one", "actor_id": "actor-reviewer", "profile": "review-default", "verdict": "changes-required", "independence": "distinct-actor", "repository": map[string]any{"project_id": "PRJ-one", "worktree_id": "WT-one", "head_oid": "abc", "tree_oid": "def", "base_oid": "base", "capture_source": "git", "dirty": false, "dirty_digest": nil}, "findings": []any{map[string]any{"id": "FND-one", "review_id": "REV-one", "severity": "major", "blocking": true, "state": "addressed", "summary": "PRIVATE REVIEW PROSE"}}},
		"finding":       {"id": "FND-one", "review_id": "REV-one", "severity": "major", "blocking": true, "state": "addressed", "summary": "PRIVATE FINDING PROSE"},
		"concern":       {"id": "CON-one", "subject": map[string]any{"type": "issue", "id": "ISS-one"}, "severity": "major", "status": "disposed", "actor_id": "actor-raiser", "disposition": map[string]any{"action": "carried-forward", "actor_id": "actor-disposer", "issue_id": "ISS-next", "reason": "PRIVATE CONCERN PROSE"}, "references": []any{source}, "reasoning": "PRIVATE CONCERN PROSE"},
		"deviation":     {"id": "DEV-one", "subject": map[string]any{"type": "issue", "id": "ISS-one"}, "disposition": "pending-decision", "authority": map[string]any{"actor_id": "actor-human"}, "references": []any{source}, "reasoning": "PRIVATE DEVIATION PROSE"},
		"decision":      {"id": "DREF-one", "class": "security", "outcome": "approved", "authority": map[string]any{"actor_id": "actor-human", "required_capabilities": []any{"decision.record-reference"}}, "source": source},
		"human-outcome": {"id": "OUT-one", "actor_id": "actor-human", "class": "security", "result": "approved", "provenance": map[string]any{"activation_id": "ACT-one", "activation_version": 2, "profile_ref": "profile:human"}, "references": []any{source}, "reasoning": "PRIVATE OUTCOME PROSE"},
		"requirement":   {"id": "RREF-one", "status": "active", "source": source},
	}
	bodies := map[string]string{
		"review":        "Reviews\n  Review REV-one\n    Issue: ISS-one; Envelope: ENV-one revision 2\n    Run: RUN-one; Reviewer: actor-reviewer; Profile: review-default\n    Verdict: changes-required; Independence: distinct-actor\n    Repository: project PRJ-one; worktree WT-one; HEAD abc; tree def\n    Repository capture: git; base base; dirty false; dirty digest <nil>\n    Findings by severity: blocker 0; major 1; minor 0; nit 0\n  Finding FND-one\n    Review: REV-one; Severity: major; Blocking: true; Disposition: addressed\n    Addressed is a Builder assertion, not a resolution.",
		"finding":       "Findings\n  Finding FND-one\n    Review: REV-one; Severity: major; Blocking: true; Disposition: addressed\n    Addressed is a Builder assertion, not a resolution.",
		"concern":       "Concerns\n  Concern CON-one\n    Subject: issue:ISS-one\n    Severity: major; Status: disposed; Raised by: actor-raiser\n    Disposition: carried-forward; Disposed by: actor-disposer\n    Work now owned by Issue: ISS-next\n    Reference: contract section; path docs/contract.md; digest sha256:abc\n    This Concern is recorded judgment, not deterministic proof.",
		"deviation":     "Deviations\n  Deviation DEV-one\n    Subject: issue:ISS-one\n    Disposition: pending-decision; Authority Actor: actor-human\n    Reference: contract section; path docs/contract.md; digest sha256:abc\n    This Deviation is recorded judgment, not deterministic proof.",
		"decision":      "Decision References\n  Decision Reference DREF-one\n    Class: security; Outcome: approved; Authority Actor: actor-human\n    Required capabilities: decision.record-reference\n    Reference: contract section; path docs/contract.md; digest sha256:abc",
		"human-outcome": "Human Outcomes\n  Human Outcome OUT-one\n    Actor: actor-human; Class: security; Result: approved\n    Activation: ACT-one revision 2; Profile: profile:human\n    This Outcome is authorized Human judgment, not deterministic proof.\n    Reference: contract section; path docs/contract.md; digest sha256:abc",
		"requirement":   "Requirement References and links\n  Requirement RREF-one\n    Status: active\n    Reference: contract section; path docs/contract.md; digest sha256:abc",
	}
	operations := []struct{ operation, kind, purpose string }{
		{"review.record", "review", "Record an independent Review of a Run and its bound Issue."},
		{"review.get", "review", "Inspect independent Reviews and their Findings."},
		{"review.list", "review", "Inspect independent Reviews and their Findings."},
		{"finding.address", "finding", "Record a Builder assertion that a Finding was addressed."},
		{"finding.list", "finding", "Inspect Finding dispositions and current blocking state."},
		{"concern.record", "concern", "Inspect a raised concern and its disposition as judgment."},
		{"concern.dispose", "concern", "Inspect a raised concern and its disposition as judgment."},
		{"concern.list", "concern", "Inspect a raised concern and its disposition as judgment."},
		{"deviation.record", "deviation", "Inspect a recorded deviation as judgment."},
		{"deviation.list", "deviation", "Inspect a recorded deviation as judgment."},
		{"decision.record-reference", "decision", "Inspect Decision References and their authority."},
		{"decision.list", "decision", "Inspect Decision References and their authority."},
		{"human-outcome.record", "human-outcome", "Inspect authorized Human judgments and their provenance."},
		{"human-outcome.list", "human-outcome", "Inspect authorized Human judgments and their provenance."},
		{"requirement.record-reference", "requirement", "Inspect Requirement References and Issue links."},
		{"requirement.link", "requirement", "Inspect Requirement References and Issue links."},
		{"requirement.get", "requirement", "Inspect Requirement References and Issue links."},
		{"requirement.list", "requirement", "Inspect Requirement References and Issue links."},
	}
	for _, tc := range operations {
		t.Run(tc.operation, func(t *testing.T) {
			t.Parallel()
			var output bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOut(&output)
			if err := writeP23Human(cmd, tc.operation, items[tc.kind]); err != nil {
				t.Fatal(err)
			}
			want := "Purpose: " + tc.purpose + "\n\n" + bodies[tc.kind] + "\n"
			if output.String() != want {
				t.Fatalf("Human text:\n%s\nwant:\n%s", output.String(), want)
			}
			if strings.Contains(output.String(), "PRIVATE") {
				t.Fatal("Human text reproduced project prose")
			}
		})
	}
}
