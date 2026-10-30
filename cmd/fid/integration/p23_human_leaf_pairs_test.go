//go:build integration

package integration

import (
	"crypto/sha256"
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
	source := p24SeedJudgmentSource(t, a)
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
		"id": "REQ-leaf", "source": source, "status": "active",
	})
	addMutation("requirement.link", a.actor, "leaf-link", map[string]any{
		"requirement_id": "REQ-leaf", "issue_id": a.issueID, "relation": "implements",
	})
	addMutation("concern.record", a.actor, "leaf-concern", map[string]any{
		"id": "CON-leaf", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"scope": map[string]any{"paths": []string{}}, "severity": "major",
		"summary": "RECORDED_CONCERN_SUMMARY", "reasoning": "RECORDED_CONCERN_REASONING", "references": []any{source},
	})
	addMutation("concern.dispose", a.actor, "leaf-dispose", map[string]any{
		"concern_id": "CON-leaf", "action": "carried-forward", "issue_id": a.issueID,
		"supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "RECORDED_DISPOSITION_REASON",
	})
	addMutation("deviation.record", a.actor, "leaf-deviation", map[string]any{
		"id": "DEV-leaf", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"expected": "RECORDED_EXPECTED", "observed": "RECORDED_OBSERVED",
		"authority":   map[string]any{"actor_id": a.actor, "reference": nil},
		"disposition": "pending-decision", "reasoning": "RECORDED_DEVIATION_REASONING",
		"references": []any{source}, "supersedes_deviation_id": nil,
	})
	addMutation("decision.record-reference", a.actor, "leaf-decision", map[string]any{
		"id": "DREF-leaf", "class": "security", "source": source,
		"authority": map[string]any{"actor_id": a.actor, "required_capabilities": []string{"decision.record-reference"}},
		"outcome":   "approved", "scope": map[string]any{"issue_id": a.issueID},
	})
	addMutation("human-outcome.record", a.actor, "leaf-outcome", map[string]any{
		"id": "OUT-leaf", "class": "security", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"result": "approved", "scope": map[string]any{"issue_id": a.issueID},
		"reasoning": "RECORDED_OUTCOME_REASONING", "references": []any{source},
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
		"expected_claim_version": mutationClaim["version"], "reason": "RECORDED_FINDING_REASON",
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
	p24CheckProjectContentSentinel(t, operation, human)
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

func p23LeafFactLines(t *testing.T, kind string, item map[string]any, list bool) []string {
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
			lines = append(lines, p23LeafFactLines(t, "finding", raw.(map[string]any), list)...)
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
	case "deviation":
		subject := item["subject"].(map[string]any)
		authority := item["authority"].(map[string]any)
		lines = append(lines, fmt.Sprintf("    Subject: %v:%v", subject["type"], subject["id"]),
			fmt.Sprintf("    Disposition: %v; Authority Actor: %v", item["disposition"], authority["actor_id"]))
		lines = append(lines, p23LeafReferenceLines(item["references"])...)
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
			fmt.Sprintf("    Activation: %v revision %v; Profile: %v", provenance["activation_id"], provenance["activation_version"], provenance["profile_ref"]))
		lines = append(lines, p23LeafReferenceLines(item["references"])...)
	}
	return append(lines, p24ExpectedRecordedLines(kind, item, list)...)
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
	if len(items) > 0 {
		switch kind {
		case "concern", "deviation":
			lines = append(lines, "  These records express judgment, not deterministic proof.")
		case "human-outcome":
			lines = append(lines, "  These Outcomes express authorized Human judgment, not deterministic proof.")
		}
	}
	if len(items) == 0 {
		lines = append(lines, "  None.")
	}
	for _, raw := range items {
		lines = append(lines, p23LeafFactLines(t, kind, raw.(map[string]any), strings.HasSuffix(operation, ".list"))...)
	}
	if operation == "requirement.get" {
		lines = append(lines, "", "Issue links")
		links := p23LeafItems(t, facts, "links")
		if len(links) == 0 {
			lines = append(lines, "  None.")
		}
		for _, raw := range links {
			lines = append(lines, p23LeafFactLines(t, kind, raw.(map[string]any), strings.HasSuffix(operation, ".list"))...)
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

const p24ProjectContentSentinel = "PROJECT_FILE_CONTENT_SENTINEL_P24_JUDGMENT"

func p24SeedJudgmentSource(t *testing.T, a *p22FIDRecoveryAdapter) map[string]any {
	t.Helper()
	const path = "project-docs/p24-judgment-sentinel.json"
	content := p24ProjectContentSentinel + "\n"
	p22Write(t, a.workspace, path, content)
	// The existing fixture ignores *.json. Referenced file bytes are checked
	// without changing a repository tree already bound to Gate Evidence.
	return map[string]any{"source_type": "project-file", "display_label": "Judgment source", "path": path, "section": nil, "external_id": nil, "digest": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(content)))}
}

func p24CheckProjectContentSentinel(t *testing.T, view, human string) {
	t.Helper()
	if strings.Contains(human, p24ProjectContentSentinel) {
		t.Fatalf("%s leaked referenced project-file content", view)
	}
}

// Independent expected text preserves the exact JSON-fact comparison introduced
// in b71ed1ba (Phase 23 task 14), converged under DEC-272/274.
func p24ExpectedRecordedLines(kind string, item map[string]any, list bool) []string {
	fields := []struct {
		label string
		value any
	}{}
	switch kind {
	case "concern":
		fields = append(fields, struct {
			label string
			value any
		}{"Summary", item["summary"]}, struct {
			label string
			value any
		}{"Reasoning", item["reasoning"]})
		if disposition, ok := item["disposition"].(map[string]any); ok {
			fields = append(fields, struct {
				label string
				value any
			}{"Disposition reason", disposition["reason"]})
		}
	case "finding":
		fields = append(fields, struct {
			label string
			value any
		}{"Summary", item["summary"]})
		if dispositions, ok := item["dispositions"].([]any); ok {
			for _, raw := range dispositions {
				fields = append(fields, struct {
					label string
					value any
				}{"Disposition reason", raw.(map[string]any)["reason"]})
			}
		}
	case "deviation":
		fields = append(fields, struct {
			label string
			value any
		}{"Expected", item["expected"]}, struct {
			label string
			value any
		}{"Observed", item["observed"]}, struct {
			label string
			value any
		}{"Reasoning", item["reasoning"]})
	case "human-outcome":
		fields = append(fields, struct {
			label string
			value any
		}{"Reasoning", item["reasoning"]})
	}
	lines := []string{}
	for _, field := range fields {
		text, ok := field.value.(string)
		if !ok {
			continue
		}
		text = strings.NewReplacer("\x1b", `\u001b`, "\r", `\r`, "\a", `\u0007`, "\x00", `\u0000`).Replace(text)
		if list {
			text = strings.Join(strings.Fields(text), " ")
			scalars := []rune(text)
			if len(scalars) > 80 {
				text = string(scalars[:80]) + "…"
			}
		}
		lines = append(lines, "    "+field.label+": "+strings.ReplaceAll(text, "\n", "\n      "))
	}
	return lines
}

func TestP24JudgmentDetailAndResumeSentinel(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	pairs := p23GenerateLeafPairs(t, binary, t.TempDir())
	workspace := pairs[0].CWD
	var issueID string
	for _, pair := range pairs {
		if pair.Operation == "concern.record" {
			issueID = pair.JSONFacts["concern"].(map[string]any)["subject"].(map[string]any)["id"].(string)
		}
	}
	if issueID == "" {
		t.Fatal("missing fixture Issue")
	}
	views := []struct {
		args     []string
		recorded []string
	}{
		{[]string{"concern", "get", "CON-leaf"}, []string{"RECORDED_CONCERN_SUMMARY", "RECORDED_CONCERN_REASONING", "RECORDED_DISPOSITION_REASON"}},
		{[]string{"deviation", "get", "DEV-leaf"}, []string{"RECORDED_EXPECTED", "RECORDED_OBSERVED", "RECORDED_DEVIATION_REASONING"}},
		{[]string{"human-outcome", "get", "OUT-leaf"}, []string{"RECORDED_OUTCOME_REASONING"}},
		{[]string{"resume", issueID, "--repository", workspace}, []string{"RECORDED_CONCERN_SUMMARY", "RECORDED_CONCERN_REASONING", "RECORDED_DISPOSITION_REASON", "RECORDED_EXPECTED", "RECORDED_OBSERVED", "RECORDED_DEVIATION_REASONING", "RECORDED_OUTCOME_REASONING", "RECORDED_FINDING_REASON"}},
	}
	for _, view := range views {
		stdout, stderr, err := closureRunCommand(workspace, binary, append([]string{"--project", workspace}, view.args...)...)
		if err != nil || stderr != "" {
			t.Fatalf("%v: %v %q", view.args, err, stderr)
		}
		p24CheckProjectContentSentinel(t, strings.Join(view.args, " "), stdout)
		for _, marker := range view.recorded {
			if !strings.Contains(stdout, marker) {
				t.Fatalf("%v missing %s: %s", view.args, marker, stdout)
			}
		}
		if !strings.Contains(stdout, "Reference: Judgment source; path project-docs/p24-judgment-sentinel.json; digest sha256:") {
			t.Fatalf("%v missing referenced metadata", view.args)
		}
		want := 1
		if view.args[0] == "resume" {
			want = 3
		}
		if strings.Count(stdout, "not deterministic proof.") != want {
			t.Fatalf("%v disclaimer count: %s", view.args, stdout)
		}
	}
}

func TestP24SentinelPreservesRepositoryBinding(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, nil)
	before := p22CorpusRepository(t, a.workspace, a.projectID)
	p24SeedJudgmentSource(t, a)
	after := p22CorpusRepository(t, a.workspace, a.projectID)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("seeding referenced content changed the bound repository: before=%v after=%v", before, after)
	}
}
