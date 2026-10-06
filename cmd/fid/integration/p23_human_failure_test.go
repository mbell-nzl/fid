//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestP23StopFailureHumanAndJSONParity(t *testing.T) {
	t.Parallel()
	fixture := p23ReviewExhaustionFixture(t, integrationFIDBinary(t), t.TempDir())
	args := append([]string(nil), fixture.Steps[0].Args...)
	for i, arg := range args {
		if arg == "--json" {
			args = append(args[:i], args[i+1:]...)
			break
		}
	}
	stdout, stderr, err := closureRunCommand(fixture.Project, integrationFIDBinary(t), args...)
	want := "error: review rework budget is exhausted\nnext: an activated Human records an approved Human Outcome of class review-rework-extension for this Envelope's Issue; see fid stop list --subject envelope:" + fixture.IDs["envelope_id"] + "\n"
	if err == nil || processExitCode(t, err) != 4 || stdout != "" || stderr != want {
		t.Fatalf("claim acquire Human denial = (exit %v, stdout %q, stderr %q), want exit 4 and two exact lines", err, stdout, stderr)
	}
	if strings.Contains(stderr, "guidance:") {
		t.Fatalf("unexpected guidance: %q", stderr)
	}
}

type p23StopFailureCase struct {
	name, sourceCode, code, operation, workspace string
	args                                         []string
}

// Every row runs the installed command tree twice against the same durable
// fixture. The code inventory below makes a newly rendered code need a row.
func TestP23StopFailureCodeMatrix(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	a := p22NewFIDRecoveryAdapter(t)
	p22AcceptedCall(t, a, "human-gate-budget-baseline", "run.observe-budget", nil)
	envelopeID := p23EnvelopeID(t, a)
	envelopeSubject := map[string]any{"type": "envelope", "id": envelopeID}
	raised := p23RunInput(t, a, "stop", "raise", a.actor, "human-envelope-stop", map[string]any{
		"subject": envelopeSubject, "code": "operator-stop", "blocked_actions": []string{"claim.acquire", "run.start"},
		"details":              map[string]any{"reason": "operator review", "reference": "policy:human-failure"},
		"resolution_condition": p23HumanOutcomeCondition(envelopeSubject, "security"), "resolver_capability": "stop.resolve",
	})
	p23AssertReviewSuccess(t, "raise envelope Stop", raised)
	stopID := p21PathString(t, p23Envelope(t, raised.Stdout), "result", "stop", "id")
	attentionID := p21PathString(t, p23Envelope(t, raised.Stdout), "result", "attention", "id")
	p22AcceptedCall(t, a, "human-gate-stop", "stop.raise", nil)
	mutation := func(group, leaf, key string, input map[string]any) []string {
		return []string{"--project", a.workspace, group, leaf, "--actor", a.actor,
			"--input", p21Input(t, a.workspace, "human-"+key, input), "--idempotency-key", key}
	}
	claimInput := map[string]any{"issue_id": a.issueID, "envelope_id": envelopeID, "envelope_revision": 1, "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": "mutation", "ttl_seconds": 60}
	runInput := map[string]any{"issue_id": a.issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": a.claimID, "claim_version": a.claimVersion, "executor": map[string]any{"id": "executor-human-failure", "source": "activated-harness"}, "obligation_id": nil, "expected_repository": p24ExpectedRepository(a.repository)}
	_, gateInput := a.p22BuildArgs(t, p22ConformanceCase{Operation: "gate.run"})
	stopResolve := func(key, id, event string) []string {
		return mutation("stop", "resolve", key, map[string]any{"stop_id": id, "fact_event_id": event})
	}
	attentionResolve := func(key, id string, version int, events []string) []string {
		return mutation("attention", "resolve", key, map[string]any{"attention_id": id, "expected_version": version, "fact_event_ids": events})
	}
	cases := []p23StopFailureCase{
		{"claim/open-stop", "CodeStopUnresolved", "stop-unresolved", "claim.acquire", a.workspace, mutation("claim", "acquire", "stop-claim", claimInput)},
		{"run/open-stop", "CodeStopUnresolved", "stop-unresolved", "run.start", a.workspace, mutation("run", "start", "stop-run", runInput)},
		{"gate/open-stop", "CodeStopUnresolved", "stop-unresolved", "gate.run", a.workspace, mutation("gate", "run", "stop-gate", gateInput)},
		{"stop/not-found", "CodeStopNotFound", "stop-not-found", "stop.resolve", a.workspace, stopResolve("stop-missing", "STP-missing", "EVT-missing")},
		{"stop/resolution-invalid", "CodeResolutionInvalid", "resolution-invalid", "stop.resolve", a.workspace, stopResolve("stop-invalid-fact", stopID, "EVT-missing")},
		{"stop/invalid-argument", "CodeInvalidArgument", "invalid-argument", "stop.resolve", a.workspace, stopResolve("stop-invalid-input", "", "EVT-missing")},
		{"attention/not-found", "CodeAttentionNotFound", "attention-not-found", "attention.resolve", a.workspace, attentionResolve("attention-missing", "ATT-missing", 1, []string{"EVT-missing"})},
		{"attention/version-conflict", "CodeAttentionVersionConflict", "attention-version-conflict", "attention.resolve", a.workspace, attentionResolve("attention-stale", attentionID, 999, []string{"EVT-missing"})},
		{"attention/resolution-invalid", "CodeResolutionInvalid", "resolution-invalid", "attention.resolve", a.workspace, attentionResolve("attention-invalid-fact", attentionID, 1, []string{"EVT-missing"})},
		{"attention/invalid-argument", "CodeInvalidArgument", "invalid-argument", "attention.resolve", a.workspace, attentionResolve("attention-invalid-input", attentionID, 0, []string{"EVT-missing"})},
		{"attention-list/invalid-argument", "CodeInvalidArgument", "invalid-argument", "attention.list", a.workspace, []string{"--project", a.workspace, "attention", "list", "--subject", "issue:ISS-1"}},
		{"stop-list/invalid-argument", "CodeInvalidArgument", "invalid-argument", "stop.list", a.workspace, []string{"--project", a.workspace, "stop", "list", "--subject", "issue:ISS-1"}},
		{"stop-raise/invalid-argument", "CodeInvalidArgument", "invalid-argument", "stop.raise", a.workspace, mutation("stop", "raise", "stop-raise-invalid", map[string]any{})},
		{"stop-raise/idempotency-conflict", "CodeIdempotencyConflict", "idempotency-conflict", "stop.raise", a.workspace, mutation("stop", "raise", "human-envelope-stop", map[string]any{
			"subject": envelopeSubject, "code": "operator-stop", "blocked_actions": []string{"claim.acquire", "run.start"},
			"details":              map[string]any{"reason": "changed operator review", "reference": "policy:human-failure"},
			"resolution_condition": p23HumanOutcomeCondition(envelopeSubject, "security"), "resolver_capability": "stop.resolve",
		})},
	}
	rework := p23StopFailureReworkFixture(t, binary)
	cases = append(cases,
		p23StopFailureCase{"claim/rework-exhausted", "CodeReviewReworkExhausted", "review-rework-exhausted", "claim.acquire", rework.workspace, rework.claimArgs},
		p23StopFailureCase{"run/rework-exhausted", "CodeReviewReworkExhausted", "review-rework-exhausted", "run.start", rework.workspace, rework.runArgs},
	)
	p23AssertStopFailureCodeInventory(t, cases)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) { p23AssertStopFailurePair(t, binary, tc) })
	}
}

type p23StopFailureRework struct {
	workspace          string
	claimArgs, runArgs []string
}

func p23StopFailureReworkFixture(t *testing.T, binary string) p23StopFailureRework {
	t.Helper()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, binary, t.TempDir(), "agent", 3600, []string{}, actions, "mutation", p22ExecutionFixture{ReviewObligation: true, ReviewCyclesMaximum: 1})
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "human-failure-builder")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "human-failure-review")
	run := p23StartRun(t, a, "actor-denied", claim, "human-failure-review")
	recorded := p23ReviewRecord(t, a, "actor-denied", "human-failure-review", p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-human-failure", "major", true)}, nil))
	p23AssertReviewSuccess(t, "exhausting Review", recorded)
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, run, "id"), p21PathString(t, claim, "id"), "human-failure-review")
	mutation := func(group, leaf, key string, input map[string]any) []string {
		return []string{"--project", a.workspace, group, leaf, "--actor", a.actor,
			"--input", p21Input(t, a.workspace, "human-"+key, input), "--idempotency-key", key}
	}
	envelopeID := p23EnvelopeID(t, a)
	claimInput := map[string]any{"issue_id": a.issueID, "envelope_id": envelopeID, "envelope_revision": 1, "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": "mutation", "ttl_seconds": 60}
	runInput := map[string]any{"issue_id": a.issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": "CLM-would-be-rework", "claim_version": 1, "executor": map[string]any{"id": "executor-would-be-rework", "source": "activated-harness"}, "obligation_id": nil, "expected_repository": p24ExpectedRepository(a.repository)}
	return p23StopFailureRework{a.workspace, mutation("claim", "acquire", "rework-claim", claimInput), mutation("run", "start", "rework-run", runInput)}
}

func p23AssertStopFailureCodeInventory(t *testing.T, cases []p23StopFailureCase) {
	t.Helper()
	source, err := os.ReadFile(filepath.Join(integrationRepositoryRoot(), "cmd", "fid", "stop.go"))
	if err != nil {
		t.Fatal(err)
	}
	part := strings.SplitN(string(source), "func renderStopError(", 2)
	if len(part) != 2 {
		t.Fatal("renderStopError missing")
	}
	body := strings.SplitN(part[1], "func writeStopFailure(", 2)[0]
	pattern := regexp.MustCompile(`cli\.(Code[A-Za-z]+)`)
	want := map[string]bool{}
	for _, match := range pattern.FindAllStringSubmatch(body, -1) {
		want[match[1]] = true
	}
	got := map[string]bool{}
	for _, tc := range cases {
		got[tc.sourceCode] = true
	}
	if len(got) != len(want) {
		t.Fatalf("Stop failure code inventory: source=%v table=%v", want, got)
	}
	for code := range want {
		if !got[code] {
			t.Fatalf("Stop failure code %s has no built-binary row", code)
		}
	}
}

func p23AssertStopFailurePair(t *testing.T, binary string, tc p23StopFailureCase) {
	t.Helper()
	before := p21JournalEventCount(t, tc.workspace)
	jsonArgs := append([]string{"--json"}, tc.args...)
	jsonOut, jsonErr, jsonExit := closureRunCommand(tc.workspace, binary, jsonArgs...)
	if jsonExit == nil {
		t.Fatalf("%s JSON succeeded: %q", tc.operation, jsonOut)
	}
	jsonStatus := processExitCode(t, jsonExit)
	var envelope struct {
		ContractVersion string                         `json:"contract_version"`
		OK              bool                           `json:"ok"`
		Operation       string                         `json:"operation"`
		Error           struct{ Code, Message string } `json:"error"`
	}
	if err := json.Unmarshal([]byte(jsonErr), &envelope); err != nil {
		t.Fatalf("JSON failure %q: %v", jsonErr, err)
	}
	if jsonOut != "" || envelope.ContractVersion != "fid-cli.v1" || envelope.OK || envelope.Operation != tc.operation || envelope.Error.Code != tc.code || !strings.HasSuffix(jsonErr, "\n") {
		t.Fatalf("JSON %s = exit %d stdout %q stderr %q", tc.name, jsonStatus, jsonOut, jsonErr)
	}
	humanOut, humanErr, humanExit := closureRunCommand(tc.workspace, binary, tc.args...)
	if humanExit == nil {
		t.Fatalf("%s Human succeeded: %q", tc.operation, humanOut)
	}
	humanStatus := processExitCode(t, humanExit)
	wantHuman := "error: " + envelope.Error.Message + "\n"
	if tc.code == "review-rework-exhausted" {
		var payload map[string]any
		if err := json.Unmarshal([]byte(jsonErr), &payload); err != nil {
			t.Fatal(err)
		}
		id := p21Path(t, payload, "error", "subject", "id").(string)
		wantHuman += "next: an activated Human records an approved Human Outcome of class review-rework-extension for this Envelope's Issue; see fid stop list --subject envelope:" + id + "\n"
	}
	if humanStatus != jsonStatus || humanOut != "" || humanErr != wantHuman {
		t.Fatalf("%s Human = (exit %d, stdout %q, stderr %q), JSON = (exit %d, stderr %q)", tc.name, humanStatus, humanOut, humanErr, jsonStatus, jsonErr)
	}
	jsonOutAgain, jsonErrAgain, jsonExitAgain := closureRunCommand(tc.workspace, binary, jsonArgs...)
	if jsonExitAgain == nil || processExitCode(t, jsonExitAgain) != jsonStatus || jsonOutAgain != jsonOut || jsonErrAgain != jsonErr {
		t.Fatalf("%s JSON failure changed after Human invocation: before %q, after %q", tc.name, jsonErr, jsonErrAgain)
	}
	if after := p21JournalEventCount(t, tc.workspace); after != before {
		t.Fatalf("%s changed Event count %d -> %d", tc.name, before, after)
	}
}
