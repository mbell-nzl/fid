//go:build integration

package integration

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestP24HumanProtectedDecisionFailure(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{},
		[]string{"claim.acquire", "envelope.activate", "envelope.create", "gate.register", "issue.transition", "run.start"},
		p22ExecutionFixture{ProtectedDecisions: []any{map[string]any{"ordinal": 0, "decision_class": "security", "required_capabilities": []string{"decision.record-reference"}, "affected_actions": []string{"issue.transition-terminal-success"}}}})
	for revision, state := range []string{"in_progress", "dev_complete", "qa_review"} {
		assertP23DecisionCommandSuccess(t, state, p23IssueTransition(t, a, state, revision+1, "p24-"+state))
	}
	denial := p23IssueTransition(t, a, "qa_pass", 4, "p24-denial")
	assertP22ProcessFailure(t, "issue.transition", "protected-decision-required", 5, denial)
	replay := p23IssueTransition(t, a, "qa_pass", 4, "p24-denial")
	assertP22ProcessFailure(t, "issue.transition", "protected-decision-required", 5, replay)
	want := p24ExpectedFailureText(t, p23Envelope(t, replay.Stderr), "issue.transition", false)
	stdout, stderr, err := closureRunCommand(a.workspace, a.binary, "--project", a.workspace, "issue", "transition", a.issueID, "--to", "qa_pass", "--actor", "agent:"+a.actor, "--expected-revision", "4", "--idempotency-key", "p24-denial")
	if err == nil || processExitCode(t, err) != 5 || stdout != "" || stderr != want {
		t.Fatalf("protected Decision Human refusal: stdout=%q stderr=%q error=%v", stdout, stderr, err)
	}
}

func TestP24HumanGateTerminalFailure(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, p22FinalActions, p22ExecutionFixture{Argv: []string{"/bin/false"}})
	p22AcceptedCall(t, a, "p24-terminal-budget", "run.observe-budget", map[string]any{"value": 0})
	_, input := a.p22BuildArgs(t, p22ConformanceCase{Operation: "gate.run"})
	args := []string{"--project", a.workspace, "gate", "run", "--actor", a.actor, "--input", p21Input(t, a.workspace, "p24-terminal", input), "--idempotency-key", "p24-terminal"}
	jsonOut, jsonErr, jsonExit := closureRunCommand(a.workspace, a.binary, append([]string{"--json"}, args...)...)
	if jsonExit == nil || processExitCode(t, jsonExit) != 8 || jsonErr != "" {
		t.Fatalf("JSON gate: %v %q %q", jsonExit, jsonOut, jsonErr)
	}
	response := p24GateTerminal(t, jsonOut)["response"].(map[string]any)
	if got := p21PathString(t, response, "error", "code"); got != "gate-unsuccessful" {
		t.Fatalf("code=%s", got)
	}
	humanOut, humanErr, humanExit := closureRunCommand(a.workspace, a.binary, args...)
	if humanExit == nil || processExitCode(t, humanExit) != 8 || humanOut == "" || humanErr != p24ExpectedFailureText(t, response, "gate.run", true) {
		t.Fatalf("Human gate: %v stdout=%q stderr=%q", humanExit, humanOut, humanErr)
	}
}

// These retained fixtures have a closed code set. Each expected line below
// comes from section 9 and the JSON failure, not the production renderer.
func p24StopFailureText(t *testing.T, payload map[string]any, operation string) string {
	t.Helper()
	f := payload["error"].(map[string]any)
	d := f["details"].(map[string]any)
	code := f["code"].(string)
	join := func(key, separator string) string {
		values, _ := d[key].([]any)
		parts := []string{}
		for _, value := range values {
			parts = append(parts, value.(string))
		}
		return strings.Join(parts, separator)
	}
	detail, next := "", ""
	switch code {
	case "stop-not-found", "attention-not-found":
		detail = "id=" + d["id"].(string)
	case "attention-version-conflict":
		detail = fmt.Sprintf("actual_version=%v, attention_id=%s, expected_version=%v", d["actual_version"], d["attention_id"], d["expected_version"])
	case "stop-unresolved", "resolution-invalid":
		parts := []string{}
		for _, raw := range d["reasons"].([]any) {
			parts = append(parts, raw.(map[string]any)["code"].(string))
		}
		detail = strings.Join(parts, ",")
		if code == "stop-unresolved" {
			next = "read each Stop's qualifying resolution: fid stop list --status open"
		}
	case "invalid-argument":
		detail = "argument=" + d["argument"].(string) + ", allowed=" + join("allowed", "|")
		if operation != "attention.list" && operation != "stop.list" {
			leaf := "fid " + strings.ReplaceAll(operation, ".", " ")
			next = "see " + leaf + " --help; template: " + leaf + " --example"
		}
	case "idempotency-conflict":
		detail = "differing fields: " + join("differing_fields", ",")
		next = "resend the original input, or use a new --idempotency-key"
	case "review-rework-exhausted":
		detail = "attention_ids=" + join("attention_ids", ",") + ", event=null, replayed=false, stop_ids=" + join("stop_ids", ",")
		next = "an activated Human records an approved Human Outcome of class review-rework-extension: fid human-outcome record (template: fid human-outcome record --example)"
	default:
		t.Fatalf("missing retained expectation for %s", code)
	}
	reason := "reason: " + code
	if detail != "" {
		reason += "; " + detail
	}
	result := "error: " + f["message"].(string) + "\n" + reason + "\n"
	ids := join("stop_ids", ",")
	attention := join("attention_ids", ",")
	if ids != "" && attention != "" {
		ids += ","
	}
	ids += attention
	if ids != "" {
		result += "blocked by: " + ids + "\n"
	}
	if next != "" {
		result += "next: " + next + "\n"
	}
	return result
}

func TestP24GenericFallbackSensitiveProcess(t *testing.T) {
	t.Parallel()
	binary, workspace := integrationFIDBinary(t), newGitWorkspace(t)
	p21Run(t, binary, workspace, "init", "--prefix", "P24", "--json")
	activationID := p21Activation(t, binary, workspace, "prepare", "p24-sensitive-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p24", "phase": 21})
	const sentinel = "P24-GENERIC-SECRET-SENTINEL"
	content := []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": []string{}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-p24-sensitive", "definition_version": 1, "environment": []any{map[string]any{"name": "P24_SECRET", "source": "literal", "value": sentinel}}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}}
	p21Component(t, binary, workspace, activationID, "gate-definitions", "p24-sensitive-first", content)
	version := p21ActivationStatus(t, workspace, activationID).Version
	path := p21Input(t, workspace, "p24-sensitive-second", map[string]any{"activation_id": activationID, "component_type": "gate-definitions", "component_id": "p24-sensitive-second", "component_version": 1, "expected_version": version, "expected_prior_digest": nil, "content": content})
	args := []string{"activation", "component", "--input", path, "--idempotency-key", "p24-sensitive-second"}
	before := p21JournalEventCount(t, workspace)
	jsonOut, jsonErr, jsonExit := closureRunCommand(workspace, binary, append([]string{"--json"}, args...)...)
	if jsonExit == nil || processExitCode(t, jsonExit) != 2 || jsonOut != "" {
		t.Fatalf("duplicate component: %v %q %q", jsonExit, jsonOut, jsonErr)
	}
	f := p23Envelope(t, jsonErr)["error"].(map[string]any)
	wantDetails := map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted activation input"}}
	message := f["message"].(string)
	if f["code"] != "invalid-argument" || f["category"] != "invalid-input" || f["retryable"] != false || f["subject"] != nil || !reflect.DeepEqual(f["details"], wantDetails) || !strings.Contains(message, "idempotency_results/insert") || !strings.Contains(message, "array must be sorted and unique") {
		t.Fatalf("did not reach generic domain fallback: %#v", f)
	}
	humanOut, humanErr, humanExit := closureRunCommand(workspace, binary, args...)
	want := "error: " + message + "\nreason: invalid-argument; argument=input, allowed=accepted activation input\nnext: see fid activation component --help; template: fid activation component --example\n"
	if humanExit == nil || processExitCode(t, humanExit) != 2 || humanOut != "" || humanErr != want {
		t.Fatalf("Human duplicate: %v %q %q want %q", humanExit, humanOut, humanErr, want)
	}
	if strings.Contains(jsonOut+jsonErr+humanOut+humanErr, sentinel) {
		t.Fatal("sensitive descriptor value leaked")
	}
	if after := p21JournalEventCount(t, workspace); after != before {
		t.Fatalf("refusal appended Events: %d -> %d", before, after)
	}
}
