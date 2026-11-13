//go:build integration

package integration

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// Setup creates the accepted stuck state. Its values are never available to
// the acting caller: the journey must rediscover them through Human output.
type p24TakeControlStart struct {
	workspace, binary, issueID string
	blockedArgs                []string
}

func p24NewTakeControlStart(t *testing.T, binary, root string) p24TakeControlStart {
	t.Helper()
	workspace := p22HumanDirectory(t, root, "take-control")
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p24@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 24 fixture")
	p22Write(t, workspace, "decision.md", "Project security decision awaiting Human approval.\n")
	p22Git(t, workspace, "add", "decision.md")
	p22Git(t, workspace, "commit", "-qm", "take-control start")
	initialized := p21Run(t, binary, workspace, "--json", "init", "--prefix", "P24")
	p22Git(t, workspace, "add", ".fid.toml")
	p22Git(t, workspace, "commit", "-qm", "project identity")
	issue := p21Run(t, binary, workspace, "--json", "issue", "create", "--kind", "task", "--title", "Human takes control", "--actor", "agent:builder", "--idempotency-key", "start-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")
	actions := []string{"claim.acquire", "claim.release", "concern.record", "decision.record-reference", "deviation.record", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.start"}
	activationID := p21Activation(t, binary, workspace, "prepare", "start-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p24", "phase": 21})
	profiles := []any{}
	grants := []any{}
	for _, actor := range []struct{ id, kind string }{{"builder", "agent"}, {"human", "human"}} {
		profiles = append(profiles, map[string]any{"actor_id": actor.id, "actor_type": actor.kind, "producer": "activated-harness", "provenance_ref": "activation:p24", "session_ref": nil, "capabilities": actions})
		grants = append(grants, map[string]any{"actions": actions, "actor_selector": map[string]any{"actor_id": actor.id, "required_capabilities": actions}, "scope": map[string]any{"issue_id": issueID}})
	}
	digest := p21Component(t, binary, workspace, activationID, "authority-profile", "take-control-authority", profiles)
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "take-control-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-take-control", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	digests := []string{digest, gateDigest}
	slices.Sort(digests)
	p21Activation(t, binary, workspace, "commit", "start-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	p21Activation(t, binary, workspace, "finalize", "start-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})
	source := map[string]any{"source_type": "project-file", "display_label": "Security decision", "path": "decision.md", "section": nil, "external_id": nil}
	guarded := []string{"envelope.complete", "issue.transition-terminal-success"}
	obligation := map[string]any{"id": "OBL-security", "type": "human-outcome", "applicability": "required", "guarded_actions": guarded, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "outcome_class": "security", "accepted_evidence_types": []string{"human-outcome"}, "accepted_trust_sources": []string{"activated-human"}, "freshness_target": "none", "artifact_policy": "metadata-sufficient", "independence": "human", "required_lane": nil, "source": source}
	input := map[string]any{"outcome": "Human takes control", "non_goals": []string{}, "risk_result": map[string]any{"source": source, "result_digest": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("security risk")))}, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{}, "authority_grants": grants, "obligations": []any{obligation}, "budgets": []any{}, "protected_decisions": []any{map[string]any{"ordinal": 0, "decision_class": "security", "required_capabilities": []string{"decision.record-reference"}, "affected_actions": guarded}}}
	// Inputs live outside the worktree, including during setup.
	created := p21Run(t, binary, workspace, "--json", "envelope", "create", "--actor", "human", "--input", p21Input(t, t.TempDir(), "envelope", input), "--idempotency-key", "start-envelope")
	envelopeID := p21PathString(t, created, "result", "envelope", "id")
	p21Run(t, binary, workspace, "--json", "envelope", "activate", "--actor", "human", "--input", p21Input(t, t.TempDir(), "activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1}), "--idempotency-key", "start-activate")
	claim := p21Run(t, binary, workspace, "--json", "claim", "acquire", "--actor", "builder", "--input", p21Input(t, t.TempDir(), "claim", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "scope": map[string]any{"paths": []string{"decision.md"}}, "mode": "mutation", "ttl_seconds": 86400}), "--idempotency-key", "start-claim")
	repository := p22CorpusRepository(t, workspace, p21PathString(t, initialized, "result", "project", "id"))
	p21Run(t, binary, workspace, "--json", "run", "start", "--actor", "builder", "--input", p21Input(t, t.TempDir(), "run", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": p21PathString(t, claim, "result", "claim", "id"), "claim_version": 1, "executor": map[string]any{"id": "builder-executor", "source": "activated-harness"}, "obligation_id": nil, "expected_repository": p24ExpectedRepository(repository)}), "--idempotency-key", "start-run")
	for revision, state := range []string{"in_progress", "dev_complete", "qa_review"} {
		out, stderr, exit := p24GetRunExit(t, binary, workspace, "issue", "transition", issueID, "--to", state, "--actor", "agent:builder", "--expected-revision", fmt.Sprint(revision+1), "--idempotency-key", "start-"+state)
		if exit != 0 || stderr != "" || out == "" {
			t.Fatalf("start transition: exit=%d stdout=%q stderr=%q", exit, out, stderr)
		}
	}
	return p24TakeControlStart{workspace: workspace, binary: binary, issueID: issueID, blockedArgs: []string{"issue", "transition", issueID, "--to", "qa_pass", "--actor", "agent:builder", "--expected-revision", "4", "--idempotency-key", "terminal-success"}}
}

func TestP24TakeControlTerminalStateProvenance(t *testing.T) {
	t.Parallel()
	start := p24NewTakeControlStart(t, integrationFIDBinary(t), t.TempDir())
	// The first Builder refusal establishes the stuck start, rather than
	// supplying the acting Human with a previously unseen workflow state.
	before := p21JournalEventCount(t, start.workspace)
	out, stderr, exit := p24GetRunExit(t, start.binary, start.workspace, start.blockedArgs...)
	if exit != 5 || out != "" || !strings.Contains(stderr, "reason: protected-decision-required") || !strings.Contains(stderr, "blocked by:") || p21JournalEventCount(t, start.workspace) != before {
		t.Fatalf("protected start: exit=%d stdout=%q stderr=%q", exit, out, stderr)
	}
	t.Logf("stuck Builder request: exit=%d\n%s", exit, stderr)
	// Collect the state-discovery surfaces and all help advertised by the
	// live command tree. No JSON-mode output, source, SQL or .fid read supplies
	// the Human's --to value. Setup IDs are not used for these reads.
	rootHelp := p24TakeControlHumanRead(t, start, "--help")
	issueHelp := p24TakeControlHumanRead(t, start, "issue", "--help")
	issues := p24TakeControlHumanRead(t, start, "issue", "list")
	issueID := strings.TrimSpace(strings.Split(issues, " | ")[0])
	if issueID == "" || !strings.Contains(rootHelp, "resume") || !strings.Contains(issueHelp, "list") {
		t.Fatal("cannot discover Issue and resume from help/output")
	}
	t.Logf("provenance field=issue-id class=fid-output source=%q", strings.Split(issues, "\n")[0])
	sources := []string{rootHelp, issueHelp, issues,
		p24TakeControlHumanRead(t, start, "resume", issueID, "--repository", start.workspace),
		p24TakeControlHumanRead(t, start, "workflow", "list"),
		p24TakeControlHumanRead(t, start, "stop", "list"),
		p24TakeControlHumanRead(t, start, "attention", "list"),
	}
	p24TakeControlHelpSources(t, start, nil, &sources)
	position := sources[3]
	if !strings.Contains(position, "  Revision: 4\n  Next states: cancelled (cancelled), qa_block (started), qa_pass (completed)\n") {
		t.Fatal("resume must supply both terminal target and current revision")
	}
	for _, source := range sources {
		for _, line := range strings.Split(source, "\n") {
			if strings.Contains(line, "qa_pass") {
				t.Logf("terminal target provenance class=fid-output/fid-help source=%q", line)
				return
			}
		}
	}
	t.Fatal("BLOCKED provenance: issue.transition --to requires terminal-success state qa_pass; no earlier Human fid output or help supplies it, and DEC-278 does not permit setup workflow states or caller-authored state values")
}

// This pairing is outside the Human journey; its JSON is never a value source.
func TestP24TakeControlResumePair(t *testing.T) {
	t.Parallel()
	start := p24NewTakeControlStart(t, integrationFIDBinary(t), t.TempDir())
	human := p24TakeControlHumanRead(t, start, "resume", start.issueID, "--repository", start.workspace)
	paired := p21Run(t, start.binary, start.workspace, "--json", "resume", start.issueID, "--repository", start.workspace)
	issue := p21Path(t, paired, "result", "resume", "issue").(map[string]any)
	if !p23SameKeys(issue, "id status status_category revision next_states") || issue["id"] != start.issueID || issue["revision"] != float64(4) || issue["status"] != "qa_review" || issue["status_category"] != "started" {
		t.Fatalf("start position = %#v", issue)
	}
	states := issue["next_states"].([]any)
	var labels []string
	for _, raw := range states {
		next := raw.(map[string]any)
		if !p23SameKeys(next, "state category") {
			t.Fatalf("NextState shape = %#v", next)
		}
		labels = append(labels, fmt.Sprintf("%s (%s)", next["state"], next["category"]))
	}
	want := "cancelled (cancelled), qa_block (started), qa_pass (completed)"
	if strings.Join(labels, ", ") != want || !strings.Contains(human, "  Revision: 4\n  Next states: "+want+"\n") {
		t.Fatalf("start Human/JSON pairing: %s / %#v", human, issue)
	}
	p21Run(t, start.binary, start.workspace, "--json", "issue", "transition", start.issueID, "--to", "cancelled", "--reason", "Test cancellation", "--actor", "human:mark", "--expected-revision", "4", "--idempotency-key", "pair-cancelled")
	terminal := p21Run(t, start.binary, start.workspace, "--json", "resume", start.issueID, "--repository", start.workspace)
	states, ok := p21Path(t, terminal, "result", "resume", "issue", "next_states").([]any)
	if !ok || len(states) != 0 {
		t.Fatalf("no-edge JSON must be [], got %#v", p21Path(t, terminal, "result", "resume", "issue"))
	}
	human = p24TakeControlHumanRead(t, start, "resume", start.issueID, "--repository", start.workspace)
	if !strings.Contains(human, "  Revision: 5\n  Next states: none.\n") {
		t.Fatalf("no-edge Human position: %s", human)
	}
}

type p24Journey struct {
	start   p24TakeControlStart
	sources []string
	values  map[string]p24Provenance
}

type p24Provenance struct{ class, source string }

func TestP24JourneyProvenanceKeepsExactSourceLine(t *testing.T) {
	t.Parallel()
	line := "  Envelope: ENV-abc (revision 1, active)"
	j := p24Journey{sources: []string{"Issue P24-111\n  Revision: 4\n" + line + "\n"}, values: map[string]p24Provenance{}}
	j.value(t, "1", "fid-output", line)
	if j.values["1"].source != line {
		t.Fatalf("provenance source = %q, want %q", j.values["1"].source, line)
	}
}

func (j *p24Journey) value(t *testing.T, value, class, source string) string {
	t.Helper()
	if !slices.Contains([]string{"fid-output", "fid-help", "caller-authored", "git", "setup"}, class) || source == "" {
		t.Fatalf("unclassified value %q", value)
	}
	if class == "fid-output" || class == "fid-help" {
		observed := false
		for _, output := range j.sources {
			observed = observed || slices.Contains(strings.Split(output, "\n"), source)
		}
		if !observed || strings.Contains(source, "\n") || !strings.Contains(source, value) {
			t.Fatalf("unobserved source for %q: %q", value, source)
		}
	}
	if class == "setup" && value != j.start.workspace && value != "human" && value != "agent:builder" {
		t.Fatal("forbidden setup value")
	}
	j.values[value] = p24Provenance{class, source}
	t.Logf("provenance value=%q class=%s source=%q", value, class, source)
	return value
}

func (j *p24Journey) read(t *testing.T, args ...string) string {
	t.Helper()
	out := p24TakeControlHumanRead(t, j.start, args...)
	j.sources = append(j.sources, out)
	return out
}

func (j *p24Journey) command(t *testing.T, parts ...string) []string {
	t.Helper()
	parent := parts[:len(parts)-1]
	help := j.read(t, append(append([]string{}, parent...), "--help")...)
	name := parts[len(parts)-1]
	if !regexp.MustCompile(`(?m)^  ` + regexp.QuoteMeta(name) + `\s`).MatchString(help) {
		t.Fatalf("command %v not advertised: %s", parts, help)
	}
	line := regexp.MustCompile(`(?m)^  ` + regexp.QuoteMeta(name) + `\s[^\n]*`).FindString(help)
	j.sources = append(j.sources, line)
	j.value(t, name, "fid-help", line)
	t.Logf("command provenance command=%q class=fid-help source=%q", strings.Join(parts, " "), j.values[name].source)
	return parts
}

func (j *p24Journey) input(t *testing.T, parts []string, key string, replacements map[string]any) string {
	t.Helper()
	help := j.read(t, append(append([]string{}, parts...), "--help")...)
	example := j.read(t, append(append([]string{}, parts...), "--example")...)
	var input map[string]any
	if err := json.Unmarshal([]byte(example), &input); err != nil {
		t.Fatal(err)
	}
	originals := map[string]string{}
	var templateValues func(any, string)
	templateValues = func(value any, path string) {
		switch v := value.(type) {
		case map[string]any:
			for name, member := range v {
				templateValues(member, path+"."+name)
			}
		case []any:
			for i, member := range v {
				templateValues(member, fmt.Sprintf("%s[%d]", path, i))
			}
		default:
			encoded, err := json.Marshal(v)
			if err != nil {
				t.Fatal(err)
			}
			originals[path] = string(encoded)
		}
	}
	templateValues(input, "input")
	for path, value := range replacements {
		fields := strings.Split(path, ".")
		object := input
		for _, field := range fields[:len(fields)-1] {
			object = object[field].(map[string]any)
		}
		object[fields[len(fields)-1]] = value
	}
	var record func(any, string)
	record = func(value any, path string) {
		switch v := value.(type) {
		case map[string]any:
			for name, member := range v {
				record(member, path+"."+name)
			}
		case []any:
			for i, member := range v {
				record(member, fmt.Sprintf("%s[%d]", path, i))
			}
		case []string:
			for i, member := range v {
				record(member, fmt.Sprintf("%s[%d]", path, i))
			}
		default:
			text := fmt.Sprint(v)
			provenance, ok := j.values[text]
			if !ok {
				encoded, err := json.Marshal(v)
				if err != nil {
					t.Fatal(err)
				}
				if v != nil && strings.Contains(text, "<") || originals[path] != string(encoded) {
					t.Fatalf("unclassified input %s=%v; help: %s", path, v, help)
				}
				for _, line := range strings.Split(example, "\n") {
					if strings.Contains(line, string(encoded)) {
						provenance = p24Provenance{"fid-help", line}
						break
					}
				}
			}
			t.Logf("input provenance field=%s value=%v class=%s source=%q", path, v, provenance.class, provenance.source)
		}
	}
	record(input, "input")
	path := p21Input(t, t.TempDir(), key, input)
	j.value(t, path, "caller-authored", "Caller wrote the --example substitutions outside the worktree: "+path)
	return path
}

func (j *p24Journey) mutate(t *testing.T, parts []string, actor, key, path string) string {
	t.Helper()
	j.value(t, actor, "setup", "Acting Actor ID: "+actor)
	j.value(t, key, "caller-authored", "Caller chose idempotency key: "+key)
	args := append(append([]string{}, parts...), "--actor", actor, "--input", path, "--idempotency-key", key)
	out := j.read(t, args...)
	return out
}

func p24JourneyLine(t *testing.T, text, prefix string) string {
	t.Helper()
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	t.Fatalf("missing source line %q in %s", prefix, text)
	return ""
}

func p24CompleteTakeControl(t *testing.T, start p24TakeControlStart) p24TakeControlStart {
	t.Helper()
	// Setup's refused request establishes the precondition; its arguments are
	// never supplied to the acting caller as provenance.
	before := p21JournalEventCount(t, start.workspace)
	out, refusal, exit := p24GetRunExit(t, start.binary, start.workspace, start.blockedArgs...)
	if exit != 5 || out != "" || !strings.Contains(refusal, "reason: protected-decision-required") || !strings.Contains(refusal, "blocked by:") || p21JournalEventCount(t, start.workspace) != before {
		t.Fatalf("stuck start: %d %s %s", exit, out, refusal)
	}
	j := p24Journey{start: start, sources: []string{refusal}, values: map[string]p24Provenance{}}
	j.value(t, start.workspace, "setup", "Project directory: "+start.workspace)
	j.command(t, "issue")
	issues := j.read(t, j.command(t, "issue", "list")...)
	issueID := strings.TrimSpace(strings.Split(issues, " | ")[0])
	j.value(t, issueID, "fid-output", strings.Split(issues, "\n")[0])
	resumeArgs := append(j.command(t, "resume"), issueID, "--repository", start.workspace)
	resume := j.read(t, resumeArgs...)
	envelopeLine := p24JourneyLine(t, resume, "  Envelope: ")
	envelopeID := strings.Fields(envelopeLine)[1]
	envelopeRevision := strings.TrimSuffix(strings.Fields(envelopeLine)[3], ",")
	j.value(t, envelopeID, "fid-output", envelopeLine)
	j.value(t, envelopeRevision, "fid-output", envelopeLine)
	stops := j.read(t, j.command(t, "stop", "list")...)
	classLine := p24JourneyLine(t, stops, "  Decision class: ")
	class := strings.TrimPrefix(classLine, "  Decision class: ")
	j.value(t, class, "fid-output", classLine)
	capabilityLine := p24JourneyLine(t, stops, "  Required capabilities: ")
	capability := strings.TrimPrefix(capabilityLine, "  Required capabilities: ")
	j.value(t, capability, "fid-output", capabilityLine)
	j.value(t, "human", "setup", "Acting Actor ID: human")
	file := j.value(t, "approval.md", "caller-authored", "Caller writes approval.md under the project root")
	if err := os.WriteFile(filepath.Join(start.workspace, file), []byte("Human reviewed and approved the security decision.\n"), 0600); err != nil {
		t.Fatal(err)
	}
	display := j.value(t, "Human security approval", "caller-authored", "Caller chose the approval display label")
	decisionID := j.value(t, "DREF-take-control", "caller-authored", "Caller chose a new Decision Reference ID")
	decision := j.command(t, "decision", "record-reference")
	path := j.input(t, decision, "decision", map[string]any{"id": decisionID, "class": class, "source.path": file, "source.display_label": display, "authority.actor_id": "human", "authority.required_capabilities": []string{capability}, "scope.issue_id": issueID})
	j.mutate(t, decision, "human", "journey-decision", path)
	outcomeID := j.value(t, "OUT-take-control", "caller-authored", "Caller chose a new Human Outcome ID")
	reasoning := j.value(t, "Security decision approved after Human review", "caller-authored", "Caller wrote the approval reasoning")
	outcome := j.command(t, "human-outcome", "record")
	path = j.input(t, outcome, "outcome", map[string]any{"id": outcomeID, "class": class, "subject.id": issueID, "scope.issue_id": issueID, "reasoning": reasoning})
	j.mutate(t, outcome, "human", "journey-outcome", path)
	resolved := j.read(t, j.command(t, "stop", "list")...)
	if !strings.Contains(resolved, "Status: resolved") || strings.Contains(resolved, "Status: open") {
		t.Fatalf("Stops not resolved: %s", resolved)
	}
	resume = j.read(t, resumeArgs...)
	obligationLine := p24JourneyLine(t, resume, "  Obligation ")
	obligationID := strings.TrimSuffix(strings.Fields(obligationLine)[1], ":")
	j.value(t, obligationID, "fid-output", obligationLine)
	link := j.command(t, "evidence", "link")
	help := j.read(t, append(append([]string{}, link...), "--help")...)
	j.value(t, "human-outcome", "fid-help", p24JourneyLine(t, help, "  input.evidence_type: "))
	path = j.input(t, link, "evidence", map[string]any{"obligation_id": obligationID, "evidence_type": "human-outcome", "human_outcome_id": outcomeID})
	j.mutate(t, link, "human", "journey-evidence", path)
	// This read follows the last mutation. Both transition preconditions come
	// exclusively from its Human lines, never from setup or JSON.
	resume = j.read(t, resumeArgs...)
	nextLine := p24JourneyLine(t, resume, "  Next states: ")
	match := regexp.MustCompile(`([a-z][a-z0-9_]*) \(completed\)`).FindStringSubmatch(nextLine)
	if len(match) != 2 {
		t.Fatalf("no completed edge in %s", nextLine)
	}
	j.sources = append(j.sources, nextLine)
	target := j.value(t, match[1], "fid-output", nextLine)
	revisionLine := p24JourneyLine(t, resume, "  Revision: ")
	j.sources = append(j.sources, revisionLine)
	revision := j.value(t, strings.TrimPrefix(revisionLine, "  Revision: "), "fid-output", revisionLine)
	t.Logf("transition --to provenance class=fid-output source=%q; --expected-revision provenance class=fid-output source=%q", nextLine, revisionLine)
	transition := j.command(t, "issue", "transition")
	j.value(t, "agent:builder", "setup", "Acting Actor ID: agent:builder")
	j.value(t, "journey-terminal", "caller-authored", "Caller chose terminal transition key")
	j.read(t, append(transition, issueID, "--to", target, "--expected-revision", revision, "--actor", "agent:builder", "--idempotency-key", "journey-terminal")...)
	complete := j.command(t, "envelope", "complete")
	expected, err := strconv.Atoi(envelopeRevision)
	if err != nil {
		t.Fatal(err)
	}
	path = j.input(t, complete, "complete", map[string]any{"envelope_id": envelopeID, "expected_revision": expected})
	finished := j.mutate(t, complete, "human", "journey-complete", path)
	if finished != "completed envelope "+envelopeID+"\n" {
		t.Fatalf("unexpected completion: %s", finished)
	}
	t.Logf("envelope.complete success: %s", finished)
	return start
}

func TestP24HumanTakeControlJourney(t *testing.T) {
	t.Parallel()
	p24CompleteTakeControl(t, p24NewTakeControlStart(t, integrationFIDBinary(t), t.TempDir()))
}

func p24TakeControlHumanRead(t *testing.T, start p24TakeControlStart, args ...string) string {
	t.Helper()
	for _, arg := range args {
		if arg == "--json" || (arg == "sql" || arg == "journal") && !slices.Contains(args, "--help") {
			t.Fatal("forbidden journey source")
		}
	}
	out, stderr, exit := p24GetRunExit(t, start.binary, start.workspace, args...)
	if exit != 0 || stderr != "" {
		t.Fatalf("Human read %v: exit=%d stdout=%q stderr=%q", args, exit, out, stderr)
	}
	if !slices.Contains(args, "--help") {
		t.Logf("Human command=%q exit=%d\n%s", strings.Join(args, " "), exit, out)
	}
	return out
}

func p24TakeControlHelpSources(t *testing.T, start p24TakeControlStart, path []string, sources *[]string) {
	t.Helper()
	help := p24TakeControlHumanRead(t, start, append(append([]string{}, path...), "--help")...)
	*sources = append(*sources, help)
	available := false
	for _, line := range strings.Split(help, "\n") {
		if line == "Available Commands:" {
			available = true
			continue
		}
		if line == "" {
			available = false
		}
		if !available {
			continue
		}
		name := strings.Fields(line)[0]
		if name == "help" || name == "completion" {
			continue
		}
		t.Logf("command provenance class=fid-help source=%q", line)
		p24TakeControlHelpSources(t, start, append(append([]string{}, path...), name), sources)
	}
}
