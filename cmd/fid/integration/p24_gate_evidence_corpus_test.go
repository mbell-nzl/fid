//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var p24GateEvidenceFields = map[string]string{"gate.run": "registration_version", "evidence.link": "obligation_id", "envelope.complete": "expected_revision"}

func TestP24GateEvidenceAndEnvelopeKeyedLeafCorpus(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, p22FinalActions)
	p24StartKeyedCapture(a.workspace)
	p22AcceptedCall(t, a, "p24-gate-evidence-budget", "run.observe-budget", map[string]any{"value": 0})
	input := map[string]any{"run_id": a.runID, "registration_id": a.registrationID, "registration_version": a.registrationVersion, "obligation_id": "OBL-p22-corpus", "expected_repository": map[string]any{"head_oid": a.repository["head_oid"], "dirty": a.repository["dirty"]}}
	baseline := p23RunInput(t, a, "gate", "run", a.actor, "p24-gate-evidence", input)
	if baseline.ExitCode != 0 || baseline.Stderr != "" {
		t.Fatalf("gate.run baseline = %d stderr=%q stdout=%q", baseline.ExitCode, baseline.Stderr, baseline.Stdout)
	}
	terminal := p24GateTerminal(t, baseline.Stdout)
	response := terminal["response"].(map[string]any)
	receiptID := p21PathString(t, response, "result", "receipt", "id")
	gateEvents := p21JournalEventCount(t, a.workspace)
	canonical := p24ReverseJSON(input)
	canonicalArgs := []string{"--project", a.workspace, "--json", "gate", "run", "--actor", a.actor, "--input", p21Input(t, a.workspace, "p24-gate-canonical", input), "--idempotency-key", "p24-gate-evidence"}
	canonicalArgs = p24ReplaceInput(t, a.workspace, canonicalArgs, "gate-reordered", canonical)
	replayStdout, replayStderr, replayErr := closureRunCommand(a.workspace, a.binary, canonicalArgs...)
	if replayErr != nil || replayStderr != "" {
		t.Fatalf("gate.run reordered replay = %v stderr=%q", replayErr, replayStderr)
	}
	replayedTerminal := p24GateTerminal(t, replayStdout)
	if want := strings.Replace(p24GateTerminalResponseBytes(t, baseline.Stdout), `"replayed":false`, `"replayed":true`, 1); want == p24GateTerminalResponseBytes(t, baseline.Stdout) || p24GateTerminalResponseBytes(t, replayStdout) != want {
		t.Fatal("gate.run terminal response bytes changed beyond replayed")
	}
	wantResponse := p23Clone(response)
	wantResponse["result"].(map[string]any)["replayed"] = true
	if got := replayedTerminal["response"]; !reflect.DeepEqual(got, wantResponse) {
		t.Fatalf("gate.run replay response = %#v, want %#v", got, wantResponse)
	}
	changed := p23Clone(input)
	changed["registration_version"] = p23NumberValue(changed["registration_version"]) + 1
	conflict := p23RunInput(t, a, "gate", "run", a.actor, "p24-gate-evidence", changed)
	assertP22ProcessFailure(t, "gate.run", "idempotency-conflict", 4, conflict)
	if got := p21Path(t, p23Envelope(t, conflict.Stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{"registration_version"}) {
		t.Fatalf("gate.run differing_fields = %#v", got)
	}
	if after := p21JournalEventCount(t, a.workspace); after != gateEvents {
		t.Fatalf("gate.run replay/conflict appended Events: %d -> %d", gateEvents, after)
	}
	linked := p23RunInput(t, a, "evidence", "link", a.actor, "p24-evidence-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": receiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}})
	if linked.ExitCode != 0 || linked.Stderr != "" {
		t.Fatalf("evidence.link = %d %q", linked.ExitCode, linked.Stderr)
	}
	for i, state := range []string{"in_progress", "dev_complete", "qa_review", "qa_pass"} {
		p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "issue", "transition", a.issueID, "--to", state, "--actor", "agent:"+a.actor, "--expected-revision", strconv.Itoa(i+1), "--idempotency-key", "p24-gate-issue-"+state)
	}
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "p24-gate-evidence")
	complete := p23RunInput(t, a, "envelope", "complete", a.actor, "p24-envelope-complete", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1})
	if complete.ExitCode != 0 || complete.Stderr != "" {
		t.Fatalf("envelope.complete = %d %q", complete.ExitCode, complete.Stderr)
	}
	captured := p24EndKeyedCapture(a.workspace)
	for operation, name := range p24GateEvidenceFields {
		if operation == "gate.run" {
			continue // The streamed Gate row is exercised above.
		}
		entry, ok := captured[operation]
		if !ok {
			t.Fatalf("missing %s in gate/evidence capture: %v", operation, p24CapturedNames(captured))
		}
		p24ExerciseCapturedLeaf(t, a.workspace, a.binary, operation, name, entry)
	}
	p24AssertFamilyConflict(t, a.workspace, a.binary, "evidence.link", captured["evidence.link"], []string{"obligation_id", "receipt_id"})
	p24AssertFamilyConflict(t, a.workspace, a.binary, "evidence.link", captured["evidence.link"], []string{"--actor"})
	p24AssertRebuildMapParity(t, a.workspace, a.binary)
}

func p24GateTerminal(t *testing.T, stdout string) map[string]any {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSpace(stdout), "\n") {
		var frame map[string]any
		if err := json.Unmarshal([]byte(line), &frame); err != nil {
			t.Fatal(err)
		}
		if frame["stream"] == "terminal" {
			return frame["data"].(map[string]any)
		}
	}
	t.Fatalf("missing gate terminal frame: %q", stdout)
	return nil
}

func p24GateTerminalResponseBytes(t *testing.T, stdout string) string {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSpace(stdout), "\n") {
		var frame struct {
			Stream string `json:"stream"`
			Data   struct {
				Response json.RawMessage `json:"response"`
			} `json:"data"`
		}
		if err := json.Unmarshal([]byte(line), &frame); err != nil {
			t.Fatal(err)
		}
		if frame.Stream == "terminal" {
			return string(frame.Data.Response)
		}
	}
	t.Fatal("missing gate terminal response")
	return ""
}
