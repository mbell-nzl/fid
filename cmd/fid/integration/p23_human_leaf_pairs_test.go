//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func p23GenerateLeafPairs(t *testing.T, binary, root string) []p23HumanLeafPair {
	t.Helper()
	workspace := p22HumanDirectory(t, root, "leaf-contract")
	actions := []string{
		"claim.acquire", "claim.release", "concern.record", "concern.dispose",
		"decision.record-reference", "deviation.record", "envelope.activate", "envelope.create",
		"finding.address", "gate.register", "human-outcome.record", "review.record",
		"run.close", "run.observe", "run.start",
	}
	a := p22NewFIDRecoveryAdapterAt(t, binary, workspace, "human", 3600, []string{}, actions)
	return p23GenerateLeafPairsForAdapter(t, a)
}

func p23GenerateLeafPairsForAdapter(t *testing.T, a *p22FIDRecoveryAdapter) []p23HumanLeafPair {
	t.Helper()
	workspace, binary := a.workspace, a.binary
	pairs := []p23HumanLeafPair{}
	addMutation := func(operation, actor, key string, input map[string]any) {
		t.Helper()
		path := p21Input(t, workspace, "leaf-"+key, input)
		parts := strings.Split(operation, ".")
		leaf := []string{parts[0], parts[1], "--input", path, "--actor", actor, "--idempotency-key", key}
		pairs = append(pairs, p23LeafPair(t, binary, workspace, operation, leaf, true))
	}
	addRead := func(operation string, leaf ...string) {
		t.Helper()
		pairs = append(pairs, p23LeafPair(t, binary, workspace, operation, leaf, false))
	}

	addMutation("requirement.record-reference", a.actor, "leaf-requirement", map[string]any{
		"id": "REQ-leaf", "source": p23DecisionSource(a, true), "status": "active",
	})
	addMutation("requirement.link", a.actor, "leaf-link", map[string]any{
		"requirement_id": "REQ-leaf", "issue_id": a.issueID, "relation": "implements",
	})
	addMutation("concern.record", a.actor, "leaf-concern", map[string]any{
		"id": "CON-leaf", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"scope": map[string]any{"paths": []string{}}, "severity": "major",
		"summary": "PRIVATE_CONCERN_SUMMARY", "reasoning": "PRIVATE_CONCERN_REASONING", "references": []any{p23DecisionSource(a, true)},
	})
	addMutation("concern.dispose", a.actor, "leaf-dispose", map[string]any{
		"concern_id": "CON-leaf", "action": "carried-forward", "issue_id": a.issueID,
		"supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "PRIVATE_DISPOSITION_REASON",
	})
	addMutation("deviation.record", a.actor, "leaf-deviation", map[string]any{
		"id": "DEV-leaf", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"expected": "PRIVATE_EXPECTED", "observed": "PRIVATE_OBSERVED",
		"authority":   map[string]any{"actor_id": a.actor, "reference": nil},
		"disposition": "pending-decision", "reasoning": "PRIVATE_DEVIATION_REASONING",
		"references": []any{p23DecisionSource(a, true)}, "supersedes_deviation_id": nil,
	})
	addMutation("decision.record-reference", a.actor, "leaf-decision", map[string]any{
		"id": "DREF-leaf", "class": "security", "source": p23DecisionSource(a, true),
		"authority": map[string]any{"actor_id": a.actor, "required_capabilities": []string{"decision.record-reference"}},
		"outcome":   "approved", "scope": map[string]any{"issue_id": a.issueID},
	})
	addMutation("human-outcome.record", a.actor, "leaf-outcome", map[string]any{
		"id": "OUT-leaf", "class": "security", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"result": "approved", "scope": map[string]any{"issue_id": a.issueID},
		"reasoning": "PRIVATE_OUTCOME_REASONING", "references": []any{p23DecisionSource(a, true)},
	})

	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "leaf-builder")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "leaf-review-claim")
	run := p23StartRun(t, a, "actor-denied", claim, "leaf-review-run")
	addMutation("review.record", "actor-denied", "leaf-review", p23ReviewInput(a,
		p21PathString(t, run, "id"), p23Number(t, run["version"]), "changes-required",
		[]any{p23Finding(a.issueID, "FND-leaf", "major", true)}, nil))
	reviewID := p23LeafFact(t, pairs[len(pairs)-1].JSONFacts, "review")["id"].(string)
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, run, "id"), p21PathString(t, claim, "id"), "leaf-review")
	mutationClaim := p23AcquireClaim(t, a, a.actor, "mutation", "leaf-address-claim")
	addMutation("finding.address", a.actor, "leaf-address", map[string]any{
		"finding_id": "FND-leaf", "claim_id": mutationClaim["id"],
		"expected_claim_version": mutationClaim["version"], "reason": "PRIVATE_FINDING_REASON",
	})
	p23ReleaseClaim(t, a, a.actor, mutationClaim, "leaf-address")

	addRead("review.get", "review", "get", reviewID)
	addRead("review.list", "review", "list", "--issue", a.issueID)
	addRead("finding.list", "finding", "list", "--issue", a.issueID)
	addRead("concern.list", "concern", "list", "--subject", "issue:"+a.issueID)
	addRead("deviation.list", "deviation", "list", "--subject", "issue:"+a.issueID)
	addRead("decision.list", "decision", "list", "--issue", a.issueID)
	addRead("human-outcome.list", "human-outcome", "list", "--issue", a.issueID)
	addRead("requirement.get", "requirement", "get", "REQ-leaf")
	addRead("requirement.list", "requirement", "list")
	return pairs
}

func p23LeafPair(t *testing.T, binary, workspace, operation string, leaf []string, mutation bool) p23HumanLeafPair {
	t.Helper()
	humanArgs := append([]string{"--project", workspace}, leaf...)
	jsonArgs := append([]string{"--project", workspace, "--json"}, leaf...)
	jsonText, jsonErrText, jsonErr := closureRunCommand(workspace, binary, jsonArgs...)
	if jsonErr != nil || jsonErrText != "" {
		t.Fatalf("%s JSON: %v %s", operation, jsonErr, jsonErrText)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(jsonText), &envelope); err != nil {
		t.Fatalf("%s JSON decode: %v", operation, err)
	}
	if envelope["operation"] != operation {
		t.Fatalf("%s JSON operation = %v", operation, envelope["operation"])
	}
	facts, ok := envelope["result"].(map[string]any)
	if !ok {
		t.Fatalf("%s JSON result = %#v", operation, envelope["result"])
	}
	human, humanErrText, humanErr := closureRunCommand(workspace, binary, humanArgs...)
	if humanErr != nil || humanErrText != "" {
		t.Fatalf("%s Human: %v %s", operation, humanErr, humanErrText)
	}
	expected := p23ExpectedLeafText(t, operation, facts)
	if human != expected {
		t.Fatalf("%s Human:\n%s\nwant from JSON facts:\n%s", operation, human, expected)
	}
	if strings.Contains(human, "PRIVATE_") {
		t.Fatalf("%s leaked project prose: %s", operation, human)
	}
	p23CheckJudgmentQualifiers(t, operation, human)
	if mutation {
		replayText, replayErrText, replayErr := closureRunCommand(workspace, binary, jsonArgs...)
		if replayErr != nil || replayErrText != "" {
			t.Fatalf("%s JSON replay: %v %s", operation, replayErr, replayErrText)
		}
		var replay map[string]any
		if err := json.Unmarshal([]byte(replayText), &replay); err != nil {
			t.Fatal(err)
		}
		replayed := replay["result"].(map[string]any)
		facts["replayed"] = true
		if !reflect.DeepEqual(facts, replayed) {
			t.Fatalf("%s Human/JSON replay facts drifted: %#v %#v", operation, facts, replayed)
		}
	}
	return p23HumanLeafPair{Operation: operation, CWD: workspace, HumanArgs: humanArgs,
		JSONArgs: jsonArgs, ExpectedHuman: expected, JSONFacts: facts}
}

func p23CheckJudgmentQualifiers(t *testing.T, operation, human string) {
	t.Helper()
	kind := strings.Split(operation, ".")[0]
	switch kind {
	case "concern", "deviation", "human-outcome":
		if !strings.Contains(human, "not deterministic proof.") {
			t.Fatalf("%s omits judgment non-proof qualifier", operation)
		}
	case "finding":
		if !strings.Contains(human, "Addressed is a Builder assertion, not a resolution.") {
			t.Fatalf("%s omits addressed qualifier", operation)
		}
	}
}

func p23LeafFact(t *testing.T, facts map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := facts[key].(map[string]any)
	if !ok {
		t.Fatalf("JSON result %s = %#v", key, facts[key])
	}
	return value
}

func p23LeafItems(t *testing.T, facts map[string]any, key string) []any {
	t.Helper()
	value, ok := facts[key].([]any)
	if !ok {
		t.Fatalf("JSON result %s = %#v", key, facts[key])
	}
	return value
}

func p23LeafSource(value any) string {
	source, _ := value.(map[string]any)
	return fmt.Sprintf("    Reference: %v; path %v; digest %v", source["display_label"], source["path"], source["digest"])
}

func p23LeafReferenceLines(value any) []string {
	refs, _ := value.([]any)
	lines := []string{}
	for _, raw := range refs {
		lines = append(lines, p23LeafSource(raw))
	}
	return lines
}

func p23LeafFactLines(t *testing.T, kind string, item map[string]any) []string {
	t.Helper()
	id := item["id"]
	if kind == "requirement" {
		if id == nil {
			return []string{fmt.Sprintf("  %v %v %v", item["requirement_id"], item["relation"], item["issue_id"])}
		}
		lines := []string{fmt.Sprintf("  Requirement %v", item["requirement_id"])}
		line := fmt.Sprintf("    Reference ID: %v; Status: %v", id, item["status"])
		if item["superseded_by_id"] != nil {
			line += fmt.Sprintf("; Superseded by: %v", item["superseded_by_id"])
		}
		return append(lines, line, p23LeafSource(item["source"]))
	}
	names := map[string]string{"review": "Review", "finding": "Finding", "concern": "Concern", "deviation": "Deviation", "decision": "Decision Reference", "human-outcome": "Human Outcome", "requirement": "Requirement"}
	lines := []string{fmt.Sprintf("  %s %v", names[kind], id)}
	switch kind {
	case "review":
		repo := item["repository"].(map[string]any)
		findings := item["findings"].([]any)
		counts := map[string]int{}
		for _, raw := range findings {
			counts[raw.(map[string]any)["severity"].(string)]++
		}
		lines = append(lines,
			fmt.Sprintf("    Issue: %v; Envelope: %v revision %v", item["issue_id"], item["envelope_id"], item["envelope_revision"]),
			fmt.Sprintf("    Run: %v; Reviewer: %v; Profile: %v", item["run_id"], item["actor_id"], item["profile"]),
			fmt.Sprintf("    Verdict: %v; Independence: %v", item["verdict"], item["independence"]),
			fmt.Sprintf("    Repository: HEAD %.12s; %s", repo["head_oid"], map[bool]string{false: "clean", true: "dirty"}[repo["dirty"].(bool)]),
			fmt.Sprintf("    Findings by severity: blocker %d; major %d; minor %d; nit %d", counts["blocker"], counts["major"], counts["minor"], counts["nit"]))
		for _, raw := range findings {
			lines = append(lines, p23LeafFactLines(t, "finding", raw.(map[string]any))...)
		}
	case "finding":
		lines = append(lines, fmt.Sprintf("    Review: %v; Severity: %v; Blocking: %v; Disposition: %v", item["review_id"], item["severity"], item["blocking"], item["state"]))
		if item["state"] == "addressed" {
			lines = append(lines, "    Addressed is a Builder assertion, not a resolution.")
		}
	case "concern":
		subject := item["subject"].(map[string]any)
		lines = append(lines, fmt.Sprintf("    Subject: %v:%v", subject["type"], subject["id"]),
			fmt.Sprintf("    Severity: %v; Status: %v; Raised by: %v", item["severity"], item["status"], item["actor_id"]))
		if disposition, ok := item["disposition"].(map[string]any); ok {
			lines = append(lines, fmt.Sprintf("    Disposition: %v; Disposed by: %v", disposition["action"], disposition["actor_id"]))
			if disposition["action"] == "carried-forward" {
				lines = append(lines, fmt.Sprintf("    Work now owned by Issue: %v", disposition["issue_id"]))
			}
		}
		lines = append(lines, p23LeafReferenceLines(item["references"])...)
		lines = append(lines, "    This Concern is recorded judgment, not deterministic proof.")
	case "deviation":
		subject := item["subject"].(map[string]any)
		authority := item["authority"].(map[string]any)
		lines = append(lines, fmt.Sprintf("    Subject: %v:%v", subject["type"], subject["id"]),
			fmt.Sprintf("    Disposition: %v; Authority Actor: %v", item["disposition"], authority["actor_id"]))
		lines = append(lines, p23LeafReferenceLines(item["references"])...)
		lines = append(lines, "    This Deviation is recorded judgment, not deterministic proof.")
	case "decision":
		authority := item["authority"].(map[string]any)
		capabilities, _ := authority["required_capabilities"].([]any)
		values := []string{}
		for _, capability := range capabilities {
			values = append(values, fmt.Sprint(capability))
		}
		lines = append(lines, fmt.Sprintf("    Class: %v; Outcome: %v; Authority Actor: %v", item["class"], item["outcome"], authority["actor_id"]),
			"    Required capabilities: "+strings.Join(values, ", "), p23LeafSource(item["source"]))
	case "human-outcome":
		provenance := item["provenance"].(map[string]any)
		lines = append(lines, fmt.Sprintf("    Actor: %v; Class: %v; Result: %v", item["actor_id"], item["class"], item["result"]),
			fmt.Sprintf("    Activation: %v revision %v; Profile: %v", provenance["activation_id"], provenance["activation_version"], provenance["profile_ref"]),
			"    This Outcome is authorized Human judgment, not deterministic proof.")
		lines = append(lines, p23LeafReferenceLines(item["references"])...)
	}
	return lines
}

func p23ExpectedLeafText(t *testing.T, operation string, facts map[string]any) string {
	t.Helper()
	kind := strings.Split(operation, ".")[0]
	purpose := map[string]string{
		"review.record":                "Record an independent Review of a Run and its bound Issue.",
		"review.get":                   "Inspect independent Reviews and their Findings.",
		"review.list":                  "Inspect independent Reviews and their Findings.",
		"finding.address":              "Record a Builder assertion that a Finding was addressed.",
		"finding.list":                 "Inspect Finding dispositions and current blocking state.",
		"concern.record":               "Inspect a raised concern and its disposition as judgment.",
		"concern.dispose":              "Inspect a raised concern and its disposition as judgment.",
		"concern.list":                 "Inspect a raised concern and its disposition as judgment.",
		"deviation.record":             "Inspect a recorded deviation as judgment.",
		"deviation.list":               "Inspect a recorded deviation as judgment.",
		"decision.record-reference":    "Inspect Decision References and their authority.",
		"decision.list":                "Inspect Decision References and their authority.",
		"human-outcome.record":         "Inspect authorized Human judgments and their provenance.",
		"human-outcome.list":           "Inspect authorized Human judgments and their provenance.",
		"requirement.record-reference": "Inspect Requirement References and Issue links.",
		"requirement.link":             "Inspect Requirement References and Issue links.",
		"requirement.get":              "Inspect Requirement References and Issue links.",
		"requirement.list":             "Inspect Requirement References and Issue links.",
	}
	headings := map[string]string{"review": "Reviews", "finding": "Findings", "concern": "Concerns", "deviation": "Deviations", "decision": "Decision References", "human-outcome": "Human Outcomes", "requirement": "Requirement References"}
	if operation == "requirement.link" {
		headings["requirement"] = "Issue links"
	}
	lines := []string{"Purpose: " + purpose[operation], "", headings[kind]}
	items := []any{}
	switch operation {
	case "review.get", "review.record":
		items = append(items, p23LeafFact(t, facts, "review"))
	case "finding.address":
		items = append(items, p23LeafFact(t, facts, "finding"))
	case "concern.record", "concern.dispose":
		items = append(items, p23LeafFact(t, facts, "concern"))
	case "deviation.record":
		items = append(items, p23LeafFact(t, facts, "deviation"))
	case "decision.record-reference":
		items = append(items, p23LeafFact(t, facts, "decision"))
	case "human-outcome.record":
		items = append(items, p23LeafFact(t, facts, "outcome"))
	case "requirement.record-reference":
		items = append(items, p23LeafFact(t, facts, "requirement"))
	case "requirement.link":
		items = append(items, p23LeafFact(t, facts, "link"))
	case "requirement.get":
		items = append(items, p23LeafFact(t, facts, "requirement"))
	default:
		items = p23LeafItems(t, facts, "items")
	}
	if len(items) == 0 {
		lines = append(lines, "  None.")
	}
	for _, raw := range items {
		lines = append(lines, p23LeafFactLines(t, kind, raw.(map[string]any))...)
	}
	if operation == "requirement.get" {
		lines = append(lines, "", "Issue links")
		links := p23LeafItems(t, facts, "links")
		if len(links) == 0 {
			lines = append(lines, "  None.")
		}
		for _, raw := range links {
			lines = append(lines, p23LeafFactLines(t, kind, raw.(map[string]any))...)
		}
	}
	if budget, ok := facts["budget"].(map[string]any); ok {
		remaining, consumed := "unknown", "unknown"
		if budget["remaining"] != nil {
			remaining = fmt.Sprint(budget["remaining"])
		}
		if budget["consumed"] != nil {
			consumed = fmt.Sprint(budget["consumed"])
		}
		line := fmt.Sprintf("Budget %v: remaining %s %v; consumed %s; maximum %v",
			budget["name"], remaining, budget["unit"], consumed, budget["maximum"])
		if budget["state"] != nil && budget["state"] != "" {
			line += "; state " + fmt.Sprint(budget["state"])
		}
		lines = append(lines, "", "Review budget", line)
	}
	return strings.Join(lines, "\n") + "\n"
}
