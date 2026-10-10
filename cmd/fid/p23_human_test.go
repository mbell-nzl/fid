package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func TestPhase23HumanLeafTextMatchesRecordedFacts(t *testing.T) {
	t.Parallel()
	source := map[string]any{"display_label": "contract section", "path": "docs/contract.md", "digest": "sha256:abc"}
	items := map[string]map[string]any{
		"review":        {"id": "REV-one", "issue_id": "ISS-one", "envelope_id": "ENV-one", "envelope_revision": 2, "run_id": "RUN-one", "actor_id": "actor-reviewer", "profile": "review-default", "verdict": "changes-required", "independence": "distinct-actor", "repository": map[string]any{"project_id": "PRJ-one", "worktree_id": "WT-one", "head_oid": "abcdef0123456789abcdef0123456789abcdef01ab", "tree_oid": "def", "base_oid": "base", "capture_source": "git", "dirty": false, "dirty_digest": nil}, "findings": []any{map[string]any{"id": "FND-one", "review_id": "REV-one", "severity": "major", "blocking": true, "state": "addressed", "summary": "RECORDED REVIEW PROSE"}}},
		"finding":       {"id": "FND-one", "review_id": "REV-one", "severity": "major", "blocking": true, "state": "addressed", "summary": "RECORDED FINDING PROSE"},
		"concern":       {"id": "CON-one", "subject": map[string]any{"type": "issue", "id": "ISS-one"}, "severity": "major", "status": "disposed", "actor_id": "actor-raiser", "disposition": map[string]any{"action": "carried-forward", "actor_id": "actor-disposer", "issue_id": "ISS-next", "reason": "RECORDED CONCERN PROSE"}, "references": []any{source}, "reasoning": "RECORDED CONCERN PROSE"},
		"deviation":     {"id": "DEV-one", "subject": map[string]any{"type": "issue", "id": "ISS-one"}, "disposition": "pending-decision", "authority": map[string]any{"actor_id": "actor-human"}, "references": []any{source}, "reasoning": "RECORDED DEVIATION PROSE"},
		"decision":      {"id": "DREF-one", "class": "security", "outcome": "approved", "authority": map[string]any{"actor_id": "actor-human", "required_capabilities": []any{"decision.record-reference"}}, "source": source},
		"human-outcome": {"id": "OUT-one", "actor_id": "actor-human", "class": "security", "result": "approved", "provenance": map[string]any{"activation_id": "ACT-one", "activation_version": 2, "profile_ref": "profile:human"}, "references": []any{source}, "reasoning": "RECORDED OUTCOME PROSE"},
		"requirement":   {"id": "RREF-one", "requirement_id": "REQ-one", "status": "active", "source": source},
	}
	bodies := map[string]string{
		"review":        "Reviews\n  Review REV-one\n    Issue: ISS-one; Envelope: ENV-one revision 2\n    Run: RUN-one; Reviewer: actor-reviewer; Profile: review-default\n    Verdict: changes-required; Independence: distinct-actor\n    Repository: HEAD abcdef012345; clean\n    Findings by severity: blocker 0; major 1; minor 0; nit 0\n  Finding FND-one\n    Review: REV-one; Severity: major; Blocking: true; Disposition: addressed\n    Addressed is a Builder assertion, not a resolution.\n    Summary: RECORDED REVIEW PROSE",
		"finding":       "Findings\n  Finding FND-one\n    Review: REV-one; Severity: major; Blocking: true; Disposition: addressed\n    Addressed is a Builder assertion, not a resolution.\n    Summary: RECORDED FINDING PROSE",
		"concern":       "Concerns\n  These records express judgment, not deterministic proof.\n  Concern CON-one\n    Subject: issue:ISS-one\n    Severity: major; Status: disposed; Raised by: actor-raiser\n    Disposition: carried-forward; Disposed by: actor-disposer\n    Work now owned by Issue: ISS-next\n    Reference: contract section; path docs/contract.md; digest sha256:abc\n    Reasoning: RECORDED CONCERN PROSE\n    Disposition reason: RECORDED CONCERN PROSE",
		"deviation":     "Deviations\n  These records express judgment, not deterministic proof.\n  Deviation DEV-one\n    Subject: issue:ISS-one\n    Disposition: pending-decision; Authority Actor: actor-human\n    Reference: contract section; path docs/contract.md; digest sha256:abc\n    Reasoning: RECORDED DEVIATION PROSE",
		"decision":      "Decision References\n  Decision Reference DREF-one\n    Class: security; Outcome: approved; Authority Actor: actor-human\n    Required capabilities: decision.record-reference\n    Reference: contract section; path docs/contract.md; digest sha256:abc",
		"human-outcome": "Human Outcomes\n  These Outcomes express authorized Human judgment, not deterministic proof.\n  Human Outcome OUT-one\n    Actor: actor-human; Class: security; Result: approved\n    Activation: ACT-one revision 2; Profile: profile:human\n    Reference: contract section; path docs/contract.md; digest sha256:abc\n    Reasoning: RECORDED OUTCOME PROSE",
		"requirement":   "Requirement References\n  Requirement REQ-one\n    Reference ID: RREF-one; Status: active\n    Reference: contract section; path docs/contract.md; digest sha256:abc",
	}
	operations := []struct{ operation, kind, purpose string }{
		{"review.record", "review", "Record an independent Review of a Run and its bound Issue."},
		{"review.get", "review", "Inspect independent Reviews and their Findings."},
		{"review.list", "review", "Inspect independent Reviews and their Findings."},
		{"finding.address", "finding", "Record a Builder assertion that a Finding was addressed."},
		{"finding.list", "finding", "Inspect Finding dispositions and current blocking state."},
		{"concern.record", "concern", "Inspect a raised concern and its disposition as judgment."},
		{"concern.dispose", "concern", "Inspect a raised concern and its disposition as judgment."},
		{"concern.get", "concern", "Inspect a raised concern and its disposition as judgment."},
		{"concern.list", "concern", "Inspect a raised concern and its disposition as judgment."},
		{"deviation.record", "deviation", "Inspect a recorded deviation as judgment."},
		{"deviation.get", "deviation", "Inspect a recorded deviation as judgment."},
		{"deviation.list", "deviation", "Inspect a recorded deviation as judgment."},
		{"decision.record-reference", "decision", "Inspect Decision References and their authority."},
		{"decision.list", "decision", "Inspect Decision References and their authority."},
		{"human-outcome.record", "human-outcome", "Inspect authorized Human judgments and their provenance."},
		{"human-outcome.get", "human-outcome", "Inspect authorized Human judgments and their provenance."},
		{"human-outcome.list", "human-outcome", "Inspect authorized Human judgments and their provenance."},
		{"requirement.record-reference", "requirement", "Inspect Requirement References and Issue links."},
		{"requirement.link", "requirement", "Inspect Requirement References and Issue links."},
		{"requirement.get", "requirement", "Inspect Requirement References and Issue links."},
		{"requirement.list", "requirement", "Inspect Requirement References and Issue links."},
	}
	for _, tc := range operations {
		t.Run(tc.operation, func(t *testing.T) {
			t.Parallel()
			item := items[tc.kind]
			if tc.operation == "requirement.link" {
				item = map[string]any{"requirement_id": "REQ-one", "relation": "implements", "issue_id": "ISS-one"}
			}
			encoded, err := json.Marshal(item)
			if err != nil {
				t.Fatal(err)
			}
			var jsonFact map[string]any
			if err := json.Unmarshal(encoded, &jsonFact); err != nil {
				t.Fatal(err)
			}
			paths := map[string][]string{
				"review":        {"id", "issue_id", "envelope_id", "envelope_revision", "run_id", "actor_id", "profile", "verdict", "independence", "repository.project_id", "repository.worktree_id", "repository.head_oid", "repository.tree_oid", "repository.base_oid", "repository.capture_source", "repository.dirty", "repository.dirty_digest", "findings"},
				"finding":       {"id", "review_id", "severity", "blocking", "state"},
				"concern":       {"id", "subject.type", "subject.id", "severity", "status", "actor_id", "disposition.action", "disposition.actor_id", "disposition.issue_id", "references"},
				"deviation":     {"id", "subject.type", "subject.id", "disposition", "authority.actor_id", "references"},
				"decision":      {"id", "class", "outcome", "authority.actor_id", "authority.required_capabilities", "source.display_label", "source.path", "source.digest"},
				"human-outcome": {"id", "actor_id", "class", "result", "provenance.activation_id", "provenance.activation_version", "provenance.profile_ref", "references"},
				"requirement":   {"id", "requirement_id", "status", "source.display_label", "source.path", "source.digest"},
			}
			if tc.operation == "requirement.link" {
				paths["requirement"] = []string{"requirement_id", "relation", "issue_id"}
			}
			for _, path := range paths[tc.kind] {
				value := any(jsonFact)
				for _, part := range strings.Split(path, ".") {
					fact, ok := value.(map[string]any)
					if !ok {
						t.Fatalf("%s Human fact %s has no JSON object", tc.operation, path)
					}
					value, ok = fact[part]
					if !ok {
						t.Fatalf("%s Human fact %s is absent from JSON", tc.operation, path)
					}
				}
			}
			var output bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOut(&output)
			if err := writeP23Human(cmd, tc.operation, jsonFact); err != nil {
				t.Fatal(err)
			}
			want := "Purpose: " + tc.purpose + "\n\n" + bodies[tc.kind] + "\n"
			if tc.operation == "requirement.link" {
				want = "Purpose: " + tc.purpose + "\n\nIssue links\n  REQ-one implements ISS-one\n"
			}
			if tc.operation == "requirement.get" {
				want = "Purpose: " + tc.purpose + "\n\n" + bodies[tc.kind] + "\n\nIssue links\n  None.\n"
			}
			if output.String() != want {
				t.Fatalf("Human text:\n%s\nwant:\n%s", output.String(), want)
			}
			if strings.Contains(output.String(), "PROJECT_FILE_CONTENT_SENTINEL") {
				t.Fatal("Human text reproduced referenced file content")
			}
			if tc.kind == "concern" || tc.kind == "deviation" || tc.kind == "human-outcome" {
				if !strings.Contains(output.String(), "not deterministic proof.") {
					t.Fatal("judgment view omitted non-proof qualifier")
				}
			}
		})
	}
}

func TestP24ConcernListSummary(t *testing.T) {
	t.Parallel()
	text := strings.Repeat("界", 81)
	got := strings.Join(p23HumanLines("concern.list", map[string]any{"id": "CON-one", "summary": text}), "\n")
	want := "    Summary: " + strings.Repeat("界", 80) + "…"
	if !strings.Contains(got, want) {
		t.Fatalf("concern list missing cut summary %q: %s", want, got)
	}
}

// The leaf set follows writeP23Human/writeP23MutationHuman callers in concern,
// deviation, human_outcome, review and finding, plus resumeJudgmentHumanLines.
var p24ReasoningViews = map[string]string{
	"concern.record": "concern", "concern.dispose": "concern", "concern.get": "concern", "concern.list": "concern",
	"deviation.record": "deviation", "deviation.get": "deviation", "deviation.list": "deviation",
	"human-outcome.record": "human-outcome", "human-outcome.get": "human-outcome", "human-outcome.list": "human-outcome",
	"review.record": "review", "review.get": "review", "review.list": "review",
	"finding.address": "finding", "finding.list": "finding",
	"resume.concern": "concern", "resume.deviation": "deviation", "resume.human-outcome": "human-outcome",
	"resume.review": "review", "resume.finding": "finding",
}

func p24JudgmentOutput(view, kind string, items []map[string]any) string {
	if strings.HasPrefix(view, "resume.") {
		resume := workflow.Resume{}
		switch kind {
		case "concern":
			resume.Concerns = items
		case "deviation":
			resume.Deviations = items
		case "human-outcome":
			resume.HumanOutcomes = items
		case "review":
			resume.Reviews = items
		case "finding":
			resume.Findings = items
		}
		return strings.Join(resumeJudgmentHumanLines(resume), "\n")
	}
	if len(items) != 1 || (view != "review.record" && view != "human-outcome.record") {
		return strings.Join(p23HumanLines(view, items...), "\n")
	}
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)
	if err := writeP23MutationHuman(cmd, view, items[0], nil, nil, nil, nil); err != nil {
		panic(err)
	}
	return out.String()
}

func TestP24JudgmentSectionDisclaimers(t *testing.T) {
	t.Parallel()
	for view, kind := range p24ReasoningViews {
		if kind != "concern" && kind != "deviation" && kind != "human-outcome" {
			continue
		}
		for _, count := range []int{0, 1, 3} {
			t.Run(fmt.Sprintf("%s/%d", view, count), func(t *testing.T) {
				items := []map[string]any{}
				for i := 0; i < count; i++ {
					items = append(items, map[string]any{"id": fmt.Sprintf("record-%d", i)})
				}
				got := p24JudgmentOutput(view, kind, items)
				want := 1
				if count == 0 {
					want = 0
					if !strings.Contains(strings.ToLower(got), "none") {
						t.Fatal("empty section missing none")
					}
				}
				if strings.Count(got, "not deterministic proof.") != want {
					t.Fatalf("disclaimer count, want %d: %s", want, got)
				}
				if count > 0 && strings.Index(got, "not deterministic proof.") > strings.Index(got, "record-0") {
					t.Fatalf("disclaimer follows item: %s", got)
				}
				for _, old := range []string{"This Concern is", "This Deviation is", "This Outcome is"} {
					if strings.Contains(got, old) {
						t.Fatalf("per-item disclaimer: %s", got)
					}
				}
			})
		}
	}
}

func TestP24RecordedReasoningEveryFieldAndView(t *testing.T) {
	t.Parallel()
	fields := map[string][]struct{ path, label string }{
		"concern":       {{"summary", "Summary"}, {"reasoning", "Reasoning"}, {"disposition.reason", "Disposition reason"}},
		"finding":       {{"summary", "Summary"}, {"dispositions.reason", "Disposition reason"}},
		"deviation":     {{"expected", "Expected"}, {"observed", "Observed"}, {"reasoning", "Reasoning"}},
		"human-outcome": {{"reasoning", "Reasoning"}},
	}
	cases := []struct{ name, input, list, detail string }{
		{"80", strings.Repeat("界", 80), strings.Repeat("界", 80), strings.Repeat("界", 80)},
		{"81", strings.Repeat("界", 81), strings.Repeat("界", 80) + "…", strings.Repeat("界", 81)},
		{"multiline", "  first  line\nsecond   line\n", "first line second line", "  first  line\n      second   line\n      "},
		{"whitespace", " \tfirst\t \n  second\u2003line ", "first second line", " \\tfirst\\t \n        second\u2003line "},
		{"controls", "a\x1b\r\a\x00z", `a\u001b\r\u0007\u0000z`, `a\u001b\r\u0007\u0000z`},
	}
	for view, kind := range p24ReasoningViews {
		fieldKind := kind
		if kind == "review" {
			fieldKind = "finding"
		}
		for _, field := range fields[fieldKind] {
			for _, tc := range cases {
				t.Run(view+"/"+field.path+"/"+tc.name, func(t *testing.T) {
					item := map[string]any{"id": "record-one"}
					switch field.path {
					case "disposition.reason":
						item["disposition"] = map[string]any{"reason": tc.input}
					case "dispositions.reason":
						item["dispositions"] = []any{map[string]any{"reason": tc.input}, map[string]any{"reason": "second disposition"}}
					default:
						item[field.path] = tc.input
					}
					if kind == "review" {
						item = map[string]any{"id": "review-one", "findings": []any{item}}
					}
					got := p24JudgmentOutput(view, kind, []map[string]any{item})
					want := tc.detail
					if strings.HasSuffix(view, ".list") || strings.HasPrefix(view, "resume.") {
						want = tc.list
					}
					if !strings.Contains(got, "    "+field.label+": "+want) {
						t.Fatalf("missing recorded field %q: %s", want, got)
					}
					if field.path == "dispositions.reason" && !strings.Contains(got, "Disposition reason: second disposition") {
						t.Fatalf("missing later disposition: %s", got)
					}
					if strings.ContainsAny(got, "\x1b\r\a\x00") {
						t.Fatalf("raw control: %q", got)
					}
				})
			}
		}
	}
}
