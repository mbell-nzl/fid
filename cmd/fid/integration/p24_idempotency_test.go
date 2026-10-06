//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestP24RunStartConflictNamesChangedClaimVersion(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapter(t)
	before := p21JournalEventCount(t, a.workspace)
	input := map[string]any{
		"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a),
		"envelope_revision": 1, "claim_id": a.claimID, "claim_version": 2,
		"executor":      map[string]any{"id": "executor-corpus-prior", "source": "activated-harness"},
		"obligation_id": nil, "expected_repository": p24ExpectedRepository(a.repository),
	}
	response := p23RunInput(t, a, "run", "start", a.actor, "p22-corpus-run-start", input)
	assertP22ProcessFailure(t, "run.start", "idempotency-conflict", 4, response)
	fields := p21Path(t, p23Envelope(t, response.Stderr), "error", "details", "differing_fields")
	if !reflect.DeepEqual(fields, []any{"claim_version"}) {
		t.Fatalf("differing_fields = %#v, want [claim_version]", fields)
	}
	if got := p21JournalEventCount(t, a.workspace); got != before {
		t.Fatalf("conflict appended Event: %d -> %d", before, got)
	}
}

func TestP24IssueCreateConflictNamesChangedTitle(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapter(t)
	before := p21JournalEventCount(t, a.workspace)
	stdout, stderr, err := closureRunCommand(a.workspace, a.binary, "--project", a.workspace, "--json", "issue", "create", "--kind", "task", "--title", "different title", "--actor", "agent:developer", "--idempotency-key", "p22-corpus-issue")
	if err == nil || processExitCode(t, err) != 4 || stdout != "" {
		t.Fatalf("issue.create conflict = %v stdout=%q stderr=%q", err, stdout, stderr)
	}
	if got := p21Path(t, p23Envelope(t, stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{"--title"}) {
		t.Fatalf("issue.create differing_fields = %#v", got)
	}
	if got := p21JournalEventCount(t, a.workspace); got != before {
		t.Fatalf("issue.create conflict appended Event: %d -> %d", before, got)
	}
}

func TestP24StopRaiseConflictNamesChangedDetails(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapter(t)
	_, input := a.p22BuildArgs(t, p22ConformanceCase{Operation: "stop.raise"})
	first := p23RunInput(t, a, "stop", "raise", a.actor, "p24-stop-conflict", input)
	if first.ExitCode != 0 {
		t.Fatalf("stop.raise first = %d stderr=%q", first.ExitCode, first.Stderr)
	}
	before := p21JournalEventCount(t, a.workspace)
	changed := p23Clone(input)
	changed["details"].(map[string]any)["reason"] = "another reason"
	response := p23RunInput(t, a, "stop", "raise", a.actor, "p24-stop-conflict", changed)
	assertP22ProcessFailure(t, "stop.raise", "idempotency-conflict", 4, response)
	if got := p21Path(t, p23Envelope(t, response.Stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{"details"}) {
		t.Fatalf("stop.raise differing_fields = %#v", got)
	}
	if got := p21JournalEventCount(t, a.workspace); got != before {
		t.Fatalf("stop.raise conflict appended Event: %d -> %d", before, got)
	}
}

func TestP24ClaimAcquireConflictNamesChangedLease(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapter(t)
	before := p21JournalEventCount(t, a.workspace)
	input := map[string]any{
		"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a),
		"envelope_revision": 1, "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": "mutation",
		"ttl_seconds": 61,
	}
	response := p23RunInput(t, a, "claim", "acquire", a.actor, "p22-corpus-claim-acquire", input)
	assertP22ProcessFailure(t, "claim.acquire", "idempotency-conflict", 4, response)
	if got := p21Path(t, p23Envelope(t, response.Stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{"ttl_seconds"}) {
		t.Fatalf("differing_fields = %#v, want [ttl_seconds]", got)
	}
	if got := p21JournalEventCount(t, a.workspace); got != before {
		t.Fatalf("claim conflict appended Event: %d -> %d", before, got)
	}
}

func TestP24GateRunPersistsNamedFieldsThroughReceipt(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterWithLease(t, "agent", 3600)
	p22AcceptedCall(t, a, "p24-gate-budget-baseline", "run.observe-budget", map[string]any{"value": 0})
	_, input := a.p22BuildArgs(t, p22ConformanceCase{Operation: "gate.run"})
	response := p23RunInput(t, a, "gate", "run", a.actor, "p24-gate-fields", input)
	if response.ExitCode != 0 {
		t.Fatalf("gate.run = %d stdout=%q stderr=%q", response.ExitCode, response.Stdout, response.Stderr)
	}
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var encoded string
	if err := conn.QueryRow(`SELECT request_field_digests_json FROM idempotency_results WHERE operation_scope LIKE 'gate.run:%' AND idempotency_key = ?`, "p24-gate-fields").Scan(&encoded); err != nil {
		t.Fatal(err)
	}
	var fields map[string]string
	if err := json.Unmarshal([]byte(encoded), &fields); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"run_id", "registration_id", "registration_version", "obligation_id", "repository", "--actor"} {
		if !strings.HasPrefix(fields[name], "sha256:") {
			t.Fatalf("gate.run field %s = %q; all fields = %v", name, fields[name], fields)
		}
	}
	before := p21JournalEventCount(t, a.workspace)
	conflict := p23RunInput(t, a, "gate", "run", "actor-denied", "p24-gate-fields", input)
	assertP22ProcessFailure(t, "gate.run", "idempotency-conflict", 4, conflict)
	if got := p21Path(t, p23Envelope(t, conflict.Stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{"--actor"}) {
		t.Fatalf("gate.run differing_fields = %#v", got)
	}
	if got := p21JournalEventCount(t, a.workspace); got != before {
		t.Fatalf("gate.run conflict appended Event: %d -> %d", before, got)
	}
}

func TestP24GateRunCommittedDenialPersistsNamedFields(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapter(t)
	_, input := a.p22BuildArgs(t, p22ConformanceCase{Operation: "gate.run"})
	response := p23RunInput(t, a, "gate", "run", a.actor, "p24-gate-denial-fields", input)
	if response.ExitCode != 4 || !strings.Contains(response.Stdout, "budget-data-unavailable") {
		t.Fatalf("gate.run denial = %d stdout=%q stderr=%q", response.ExitCode, response.Stdout, response.Stderr)
	}
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var encoded string
	if err := conn.QueryRow(`SELECT request_field_digests_json FROM idempotency_results WHERE operation_scope LIKE 'gate.run:%' AND idempotency_key = ?`, "p24-gate-denial-fields").Scan(&encoded); err != nil {
		t.Fatal(err)
	}
	var fields map[string]string
	if err := json.Unmarshal([]byte(encoded), &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 6 || !strings.HasPrefix(fields["--actor"], "sha256:") || !strings.HasPrefix(fields["run_id"], "sha256:") {
		t.Fatalf("gate.run denial fields = %v", fields)
	}
	before := p21JournalEventCount(t, a.workspace)
	replay := p23RunInput(t, a, "gate", "run", a.actor, "p24-gate-denial-fields", input)
	if replay.ExitCode != 4 || !strings.Contains(replay.Stdout, `"replayed":true`) {
		t.Fatalf("gate.run denial replay = %d stdout=%q stderr=%q", replay.ExitCode, replay.Stdout, replay.Stderr)
	}
	if got := p21JournalEventCount(t, a.workspace); got != before {
		t.Fatalf("gate.run denial replay appended Event: %d -> %d", before, got)
	}
}

func TestP24ActivationPrepareConflictNamesChangedKit(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapter(t)
	before := p21JournalEventCount(t, a.workspace)
	input := map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "different-kit", "phase": 21,
		"expected_prior_activation_id": nil, "expected_prior_version": nil, "resources": p21ActivationResources()}
	path := p21Input(t, a.workspace, "p24-activation-conflict", input)
	stdout, stderr, err := closureRunCommand(a.workspace, a.binary, "--project", a.workspace, "--json", "activation", "prepare", "--input", path, "--idempotency-key", "p22-corpus-activation-prepare")
	if err == nil || processExitCode(t, err) != 4 || stdout != "" {
		t.Fatalf("activation.prepare conflict = %v stdout=%q stderr=%q, want exit 4", err, stdout, stderr)
	}
	got := p21Path(t, p23Envelope(t, stderr), "error", "details", "differing_fields")
	if !reflect.DeepEqual(got, []any{"kit_identity"}) {
		t.Fatalf("differing_fields = %#v, want [kit_identity]", got)
	}
	if after := p21JournalEventCount(t, a.workspace); after != before {
		t.Fatalf("activation conflict appended Event: %d -> %d", before, after)
	}
}

func TestP24ReviewOptionalPresenceAndEmptyException(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation")
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "p24-review-optional-setup")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "p24-review-optional")
	run := p23StartRun(t, a, "actor-denied", claim, "p24-review-optional")
	input := p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "pass", []any{}, nil)
	delete(input, "inapplicable_findings")
	first := p23ReviewRecord(t, a, "actor-denied", "p24-review-optional", input)
	p23AssertReviewSuccess(t, "first review", first)
	before := p21JournalEventCount(t, a.workspace)
	changed := p23Clone(input)
	changed["inapplicable_findings"] = []any{map[string]any{"finding_id": "FND-missing", "reason": "not applicable"}}
	conflict := p23ReviewRecord(t, a, "actor-denied", "p24-review-optional", changed)
	assertP22ProcessFailure(t, "review.record", "idempotency-conflict", 4, conflict)
	if got := p21Path(t, p23Envelope(t, conflict.Stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{"inapplicable_findings"}) {
		t.Fatalf("differing_fields = %#v, want [inapplicable_findings]", got)
	}
	empty := p23Clone(input)
	empty["inapplicable_findings"] = []any{}
	replay := p23ReviewRecord(t, a, "actor-denied", "p24-review-optional", empty)
	p23AssertReviewSuccess(t, "empty optional replay", replay)
	if got := p21Path(t, p23Envelope(t, replay.Stdout), "result", "replayed"); got != true {
		t.Fatalf("empty optional replay = %#v, want true", got)
	}
	if got := p21JournalEventCount(t, a.workspace); got != before {
		t.Fatalf("review conflict/replay appended Event: %d -> %d", before, got)
	}
}

func TestP24ObservingJudgmentLeavesNameChangedPublicField(t *testing.T) {
	t.Parallel()
	pairs := p23GenerateLeafPairs(t, integrationFIDBinary(t), t.TempDir())
	wantField := map[string]string{
		"concern.record": "id", "concern.dispose": "reason",
		"deviation.record": "id", "human-outcome.record": "id",
		"decision.record-reference": "id", "requirement.record-reference": "id",
		"requirement.link": "requirement_id", "finding.address": "reason",
	}
	seen := map[string]bool{}
	for _, pair := range pairs {
		field, included := wantField[pair.Operation]
		if !included {
			continue
		}
		seen[pair.Operation] = true
		request := p23LeafRequest(t, pair)
		request[field] = request[field].(string) + "-changed"
		path := p21Input(t, pair.CWD, "p24-"+strings.ReplaceAll(pair.Operation, ".", "-"), request)
		args := append([]string(nil), pair.JSONArgs...)
		for i, arg := range args {
			if arg == "--input" {
				args[i+1] = path
			}
		}
		before := p21JournalEventCount(t, pair.CWD)
		_, stderr, err := closureRunCommand(pair.CWD, integrationFIDBinary(t), args...)
		if err == nil || processExitCode(t, err) != 4 {
			t.Fatalf("%s changed field = %v %s, want exit 4", pair.Operation, err, stderr)
		}
		var envelope map[string]any
		if err := json.Unmarshal([]byte(stderr), &envelope); err != nil {
			t.Fatal(err)
		}
		failure := envelope["error"].(map[string]any)
		if failure["code"] != "idempotency-conflict" {
			t.Fatalf("%s changed field code = %v", pair.Operation, failure["code"])
		}
		got := failure["details"].(map[string]any)["differing_fields"]
		if !reflect.DeepEqual(got, []any{field}) {
			t.Fatalf("%s differing_fields = %#v, want [%s]", pair.Operation, got, field)
		}
		if after := p21JournalEventCount(t, pair.CWD); after != before {
			t.Fatalf("%s conflict appended Event: %d -> %d", pair.Operation, before, after)
		}
	}
	if len(seen) != len(wantField) {
		t.Fatalf("observing leaf fixture coverage = %v, want %v", seen, wantField)
	}
}

func TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	pairs := p23GenerateLeafPairs(t, binary, t.TempDir())
	observing := map[string]bool{
		"concern.record": true, "concern.dispose": true, "deviation.record": true,
		"human-outcome.record": true, "decision.record-reference": true,
		"requirement.record-reference": true, "requirement.link": true,
	}
	baseline := map[string]string{}
	var workspace string
	for _, pair := range pairs {
		if !observing[pair.Operation] {
			continue
		}
		workspace = pair.CWD
		stdout, stderr, err := closureRunCommand(workspace, binary, pair.JSONArgs...)
		if err != nil || stderr != "" {
			t.Fatalf("%s baseline retry = %v %q", pair.Operation, err, stderr)
		}
		_ = p23Envelope(t, stdout)
		baseline[pair.Operation] = stdout
	}
	if len(baseline) != len(observing) {
		t.Fatalf("observing fixture coverage = %d, want %d", len(baseline), len(observing))
	}
	if err := os.WriteFile(filepath.Join(workspace, "base.txt"), []byte("changed source\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	p22Git(t, workspace, "add", "base.txt")
	p22Git(t, workspace, "commit", "-m", "change observed source")
	before := p21JournalEventCount(t, workspace)
	for _, pair := range pairs {
		want, included := baseline[pair.Operation]
		if !included {
			continue
		}
		stdout, stderr, err := closureRunCommand(workspace, binary, pair.JSONArgs...)
		if err != nil || stderr != "" {
			t.Fatalf("%s changed-source retry = %v %q", pair.Operation, err, stderr)
		}
		if got := stdout; got != want {
			t.Fatalf("%s changed-source replay = %#v, want %#v", pair.Operation, got, want)
		}
		if after := p21JournalEventCount(t, workspace); after != before {
			t.Fatalf("%s replay appended Event: %d -> %d", pair.Operation, before, after)
		}
	}
}

func TestP24ObservingJudgmentReplayNeverObserves(t *testing.T) {
	t.Parallel()
	sourceBearing := map[string]bool{
		"concern.record": true, "deviation.record": true, "human-outcome.record": true,
		"decision.record-reference": true, "requirement.record-reference": true,
	}
	for _, operation := range []string{
		"concern.record", "concern.dispose", "deviation.record", "human-outcome.record",
		"decision.record-reference", "requirement.record-reference", "requirement.link",
	} {
		t.Run(operation, func(t *testing.T) {
			t.Parallel()
			binary := integrationFIDBinary(t)
			pairs := p23GenerateLeafPairs(t, binary, t.TempDir())
			byOperation := map[string]p23HumanLeafPair{}
			for _, pair := range pairs {
				byOperation[pair.Operation] = pair
			}
			pair, found := byOperation[operation]
			if !found {
				t.Fatalf("missing observing leaf %s", operation)
			}
			workspace := pair.CWD
			baseline := p24RetryJudgment(t, binary, pair)
			control := func(marker string) []string {
				t.Helper()
				request := p23LeafRequest(t, pair)
				switch operation {
				case "concern.dispose":
					newID := "CON-p24-" + marker
					p24InvokeJudgment(t, binary, byOperation["concern.record"], map[string]any{"id": newID}, "setup-"+marker, true)
					request["concern_id"] = newID
				case "requirement.link":
					if marker == "preflight" {
						request["relation"] = "verifies"
					} else {
						request["relation"] = "relates-to"
					}
				default:
					request["id"] = request["id"].(string) + "-" + marker
				}
				return p24JudgmentArgs(t, pair, request, marker)
			}
			preflight := control("preflight")
			if stdout, stderr, err := closureRunCommand(workspace, binary, preflight...); err != nil || stderr != "" {
				t.Fatalf("%s observable new-key control = %v stdout=%q stderr=%q", operation, err, stdout, stderr)
			}
			blocked := control("blocked")
			p22MakeRepositoryUnsafe(t, &p22FIDRecoveryAdapter{workspace: workspace}, "unmerged")
			before := p21JournalEventCount(t, workspace)
			if got := p24RetryJudgment(t, binary, pair); !reflect.DeepEqual(got, baseline) {
				t.Fatalf("%s unmerged replay = %#v, want %#v", operation, got, baseline)
			}
			if stdout, stderr, err := closureRunCommand(workspace, binary, blocked...); err == nil || stdout != "" || stderr == "" {
				t.Fatalf("%s unmerged new-key control = %v stdout=%q stderr=%q, want refusal", operation, err, stdout, stderr)
			}
			if after := p21JournalEventCount(t, workspace); after != before {
				t.Fatalf("%s unmerged replay/control appended Event: %d -> %d", operation, before, after)
			}
			request := p23LeafRequest(t, pair)
			_, hasSource := request["source"]
			if sourceBearing[operation] != (hasSource || request["references"] != nil) {
				t.Fatalf("%s project-file Source enumeration disagrees with fixture", operation)
			}
			if sourceBearing[operation] {
				deleted := control("deleted")
				if err := os.Remove(filepath.Join(workspace, "base.txt")); err != nil {
					t.Fatal(err)
				}
				if got := p24RetryJudgment(t, binary, pair); !reflect.DeepEqual(got, baseline) {
					t.Fatalf("%s deleted-source replay = %#v, want %#v", operation, got, baseline)
				}
				if after := p21JournalEventCount(t, workspace); after != before {
					t.Fatalf("%s deleted-source retry appended Event: %d -> %d", operation, before, after)
				}
				stdout, stderr, err := closureRunCommand(workspace, binary, deleted...)
				if err == nil || stdout != "" || stderr == "" {
					t.Fatalf("%s deleted-source new-key control = %v stdout=%q stderr=%q, want refusal", operation, err, stdout, stderr)
				}
				if code, ok := p21Path(t, p23Envelope(t, stderr), "error", "code").(string); !ok || code == "" {
					t.Fatalf("%s deleted-source new-key control lacks typed error code", operation)
				}
				if after := p21JournalEventCount(t, workspace); after != before {
					t.Fatalf("%s deleted-source control appended Event: %d -> %d", operation, before, after)
				}
			}
		})
	}
}

func p24RetryJudgment(t *testing.T, binary string, pair p23HumanLeafPair) string {
	t.Helper()
	stdout, stderr, err := closureRunCommand(pair.CWD, binary, pair.JSONArgs...)
	if err != nil || stderr != "" {
		t.Fatalf("%s replay = %v stderr=%q", pair.Operation, err, stderr)
	}
	_ = p23Envelope(t, stdout)
	return stdout
}

func p24JudgmentArgs(t *testing.T, pair p23HumanLeafPair, request map[string]any, marker string) []string {
	t.Helper()
	args := append([]string(nil), pair.JSONArgs...)
	for i, arg := range args {
		switch arg {
		case "--input":
			args[i+1] = p21Input(t, pair.CWD, "p24-"+strings.ReplaceAll(pair.Operation, ".", "-")+"-"+marker, request)
		case "--idempotency-key":
			args[i+1] += "-p24-" + marker
		}
	}
	return args
}

func p24InvokeJudgment(t *testing.T, binary string, pair p23HumanLeafPair, changes map[string]any, marker string, wantSuccess bool) {
	t.Helper()
	request := p23LeafRequest(t, pair)
	for field, value := range changes {
		request[field] = value
	}
	args := p24JudgmentArgs(t, pair, request, marker)
	stdout, stderr, err := closureRunCommand(pair.CWD, binary, args...)
	if wantSuccess && (err != nil || stderr != "") {
		t.Fatalf("%s setup = %v stdout=%q stderr=%q", pair.Operation, err, stdout, stderr)
	}
}
