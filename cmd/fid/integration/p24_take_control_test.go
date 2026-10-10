//go:build integration

package integration

import (
	"crypto/sha256"
	"fmt"
	"slices"
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
	issue := p21Run(t, binary, workspace, "--json", "issue", "create", "--kind", "task", "--title", "Human takes control", "--actor", "agent:builder", "--idempotency-key", "start-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")
	actions := []string{"claim.acquire", "claim.release", "decision.record-reference", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.start"}
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
