//go:build integration

package integration

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestP22ResumeAndGateAgreeDuringIncompleteActivation(t *testing.T) {
	adapter := p22NewFIDRecoveryAdapter(t)
	p22ExecuteHumanSetup(t, adapter, p22ConformanceCase{Name: "activation-resume-budget", Operation: "run.observe-budget", Input: map[string]any{"value": 0}})
	run := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "get", adapter.runID)
	issueID := p21PathString(t, run, "result", "run", "issue_id")
	before := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "resume", issueID)
	if action := p22GateAction(t, p21Path(t, before, "result", "resume").(map[string]any)); len(action["reasons"].([]any)) != 0 {
		t.Fatalf("complete activation did not qualify gate action: %#v", action)
	}
	prior := p21ActivationStatus(t, adapter.workspace, adapter.activation["id"].(string))
	resources := p21ActivationResources()
	for _, item := range resources {
		resource := item.(map[string]any)
		resource["old_identity"] = resource["new_identity"]
		resource["new_identity"] = resource["component"].(string) + "-successor"
	}
	prepared := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "activation", "prepare", "--input", p21Input(t, adapter.workspace, "activation-resume-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "resume-upgrade", "phase": 21, "expected_prior_activation_id": prior.ID, "expected_prior_version": prior.Version, "resources": resources}), "--idempotency-key", "activation-resume-prepare")
	activationID := p21PathString(t, prepared, "result", "activation", "id")
	// Recovery selection produces the accepted activation-incomplete Stop and
	// Attention, in addition to the activation's own incomplete trust guard.
	p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "activation", "recover", "--input", p21Input(t, adapter.workspace, "activation-resume-recover", map[string]any{"activation_id": activationID, "expected_version": 1, "expected_status": "prepared", "direction": "restore-prior", "reason": "recover interrupted upgrade"}), "--idempotency-key", "activation-resume-recover")
	journalBefore := p21JournalBytes(t, adapter.workspace)
	resumeEnvelope := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "resume", issueID)
	resume := p21Path(t, resumeEnvelope, "result", "resume").(map[string]any)
	action := p22GateAction(t, resume)
	if len(action["reasons"].([]any)) == 0 {
		t.Fatalf("incomplete activation advertised gate as executable: %#v", action)
	}
	stops := resume["stops"].([]any)
	attention := resume["attention"].([]any)
	if len(stops) != 1 || len(attention) != 1 || p22ResultPath(stops[0].(map[string]any), "subject.id") != activationID || p22ResultPath(stops[0].(map[string]any), "code") != "activation-incomplete" {
		t.Fatalf("resume omitted activation Stop/Attention: %#v", resume)
	}
	human, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, "--project", adapter.workspace, "resume", issueID)
	if err != nil || stderr != "" || !strings.Contains(human, activationID) || !strings.Contains(human, "activation-incomplete") {
		t.Fatalf("Human resume hides activation blocker: %q %q %v", human, stderr, err)
	}
	gate := adapter.Execute(t, p22ConformanceCase{Name: "activation-resume-gate", Operation: "gate.run"})
	if gate.ExitCode == 0 || strings.Contains(gate.Stdout, `"kind":"start"`) {
		t.Fatalf("incomplete activation actually executed a gate: %#v", gate)
	}
	var failure map[string]any
	if err := json.Unmarshal([]byte(gate.Stderr), &failure); err != nil {
		t.Fatalf("gate did not return typed rejection: %#v (%v)", gate, err)
	}
	if code := p22ResultPath(failure, "error.code"); code != "authority-unavailable" {
		t.Fatalf("gate error = %v, want retained authority-unavailable guard", code)
	}
	foundAuthorityReason := false
	for _, reason := range action["reasons"].([]any) {
		foundAuthorityReason = foundAuthorityReason || reason.(map[string]any)["code"] == "authority-unavailable"
	}
	if !foundAuthorityReason {
		t.Fatalf("resume omitted the actual gate authority guard: %#v", action)
	}
	details := p22ResultPath(failure, "error.details").(map[string]any)
	p22AssertKeys(t, details, "actor_id operation required_actions required_capabilities")
	if details["operation"] != "gate.run" || details["actor_id"] != nil {
		t.Fatalf("gate authority failure details drifted: %#v", details)
	}
	if after := p21JournalBytes(t, adapter.workspace); string(after) != string(journalBefore) {
		t.Fatal("inspection or rejected gate appended an Event")
	}
}

func p22GateAction(t *testing.T, resume map[string]any) map[string]any {
	t.Helper()
	for _, item := range resume["next_actions"].([]any) {
		action := item.(map[string]any)
		if action["operation"] == "gate.run" {
			return action
		}
	}
	t.Fatal("resume has no Gate action")
	return nil
}
