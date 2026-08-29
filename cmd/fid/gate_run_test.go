package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame(t *testing.T) {
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "gate-run-claim")
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	runPath := filepath.Join(workspace, "gate-run-start.json")
	writeActivationInput(t, runPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace,
		"--project", workspace, "--json", "run", "start", "--actor", "actor-builder",
		"--input", runPath, "--idempotency-key", "gate-run-start",
	)
	if started.err != nil {
		t.Fatalf("start Run: %v; stderr=%q", started.err, started.stderr)
	}
	runID := runFromResponse(t, started.stdout)["id"].(string)

	envelope := decodeCommandJSONStructure(t, executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "get", envelopeID,
	).stdout).(map[string]any)["result"].(map[string]any)["envelope"].(map[string]any)
	obligation := envelope["obligations"].([]any)[0].(map[string]any)
	registration := obligation["registration"].(map[string]any)
	inputPath := filepath.Join(workspace, "gate-run.json")
	writeActivationInput(t, inputPath, map[string]any{
		"registration_id": registration["id"], "registration_version": registration["version"],
		"obligation_id": obligation["id"], "run_id": runID, "repository": repository,
	})

	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "gate", "run", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "gate-run-success",
	)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("gate run = %#v (%v), want live stdout-only success", result, result.err)
	}
	lines := strings.Split(strings.TrimSpace(result.stdout), "\n")
	if len(lines) == 0 {
		t.Fatal("gate run emitted no live frames")
	}
	var terminal struct {
		Stream string `json:"stream"`
		Data   struct {
			AttemptID string `json:"attempt_id"`
			ReceiptID string `json:"receipt_id"`
			Outcome   string `json:"outcome"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &terminal); err != nil {
		t.Fatalf("decode terminal frame: %v; frame=%q", err, lines[len(lines)-1])
	}
	if terminal.Stream != "terminal" || terminal.Data.AttemptID == "" || terminal.Data.ReceiptID == "" || terminal.Data.Outcome != "succeeded" {
		t.Fatalf("terminal = %#v, want successful Attempt and Receipt", terminal)
	}
	journal := string(mustReadCommandFile(t, mustJournalPath(t, workspace)))
	if strings.Index(journal, `"event_type":"gate.execution-started"`) >= strings.Index(journal, `"event_type":"gate.receipt-recorded"`) {
		t.Fatalf("gate Event order = %q, want Attempt before Receipt", journal)
	}
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open gate database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var attempts, receipts, evidenceRecords int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts`).Scan(&attempts); err != nil {
		t.Fatalf("count Attempts: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_receipts`).Scan(&receipts); err != nil {
		t.Fatalf("count Receipts: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM evidence_records`).Scan(&evidenceRecords); err != nil {
		t.Fatalf("count Evidence: %v", err)
	}
	if attempts != 1 || receipts != 1 || evidenceRecords != 0 {
		t.Fatalf("Attempt/Receipt/Evidence counts = %d/%d/%d, want 1/1/0", attempts, receipts, evidenceRecords)
	}
	replayed := executeCommandForTest(workspace,
		"--project", workspace, "--json", "gate", "run", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "gate-run-success",
	)
	if replayed.err != nil || replayed.stderr != "" {
		t.Fatalf("replayed gate run = %#v, want one terminal success", replayed)
	}
	replayedLines := strings.Split(strings.TrimSpace(replayed.stdout), "\n")
	if len(replayedLines) != 1 {
		t.Fatalf("replayed gate run frames = %d, want terminal frame only", len(replayedLines))
	}
	var replayTerminal struct {
		Stream string `json:"stream"`
		Data   struct {
			Replayed bool `json:"replayed"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(replayedLines[0]), &replayTerminal); err != nil || replayTerminal.Stream != "terminal" || !replayTerminal.Data.Replayed {
		t.Fatalf("replayed terminal = %#v (%v), want replayed terminal", replayTerminal, err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts`).Scan(&attempts); err != nil {
		t.Fatalf("count replayed Attempts: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_receipts`).Scan(&receipts); err != nil {
		t.Fatalf("count replayed Receipts: %v", err)
	}
	if attempts != 1 || receipts != 1 {
		t.Fatalf("replayed Attempt/Receipt counts = %d/%d, want 1/1", attempts, receipts)
	}
}

func TestGateRunRecordsExactNonSuccessTerminalOutcomes(t *testing.T) {
	cases := []struct {
		name         string
		modify       func(map[string]any)
		wantOutcome  string
		wantExitCode any
		wantSignal   bool
		wantTimedOut bool
	}{
		{"failed", func(definition map[string]any) { definition["argv"] = []string{"/bin/false"} }, "failed", float64(1), false, false},
		{"spawn failed", func(definition map[string]any) { definition["argv"] = []string{"/missing/fid-gate"} }, "spawn-failed", nil, false, false},
		{"timed out", func(definition map[string]any) {
			definition["argv"] = []string{"/bin/sleep", "1"}
			definition["timeout_ms"] = 1
		}, "timed-out", nil, true, true},
		{"signaled", func(definition map[string]any) {
			definition["argv"] = []string{os.Args[0], "-test.run=^TestGateRunHelperProcess$", "--"}
			definition["environment"] = []any{map[string]any{"name": "FID_GATE_HELPER", "source": "literal", "value": "signal"}}
		}, "signaled", nil, true, false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, testCase.modify)
			terminal, result := executeGateRunForTest(t, workspace, issueID, envelopeID, "gate-run-"+strings.ReplaceAll(testCase.name, " ", "-"))
			if result.err == nil || result.stderr != "" {
				t.Fatalf("gate run %s = %#v, want terminal failure only", testCase.name, result)
			}
			if exitCode := cli.ExitCode(result.err); exitCode != 8 {
				t.Fatalf("gate run %s exit class = %d, want 8", testCase.name, exitCode)
			}
			if terminal.Stream != "terminal" || terminal.Data.ReceiptID == "" || terminal.Data.Outcome != testCase.wantOutcome || terminal.Data.ExitCode != testCase.wantExitCode || (terminal.Data.Signal != nil) != testCase.wantSignal || terminal.Data.TimedOut != testCase.wantTimedOut {
				t.Fatalf("%s terminal = %#v, want outcome=%q exit=%#v signal=%t timed_out=%t", testCase.name, terminal, testCase.wantOutcome, testCase.wantExitCode, testCase.wantSignal, testCase.wantTimedOut)
			}
		})
	}
}

func TestGateRunHelperProcess(t *testing.T) {
	if os.Getenv("FID_GATE_HELPER") != "signal" {
		return
	}
	_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

type gateTerminalFrame struct {
	Stream string `json:"stream"`
	Data   struct {
		AttemptID string `json:"attempt_id"`
		ReceiptID string `json:"receipt_id"`
		Outcome   string `json:"outcome"`
		ExitCode  any    `json:"exit_code"`
		Signal    any    `json:"signal"`
		TimedOut  bool   `json:"timed_out"`
	} `json:"data"`
}

func executeGateRunForTest(t *testing.T, workspace, issueID, envelopeID, key string) (gateTerminalFrame, commandExecutionResult) {
	t.Helper()
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, key+"-claim")
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	runPath := filepath.Join(workspace, key+"-run.json")
	writeActivationInput(t, runPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", runPath, "--idempotency-key", key+"-start")
	if started.err != nil {
		t.Fatalf("start Run: %v; stderr=%q", started.err, started.stderr)
	}
	envelopeResult := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "get", envelopeID)
	if envelopeResult.err != nil {
		t.Fatalf("get envelope: %v; stderr=%q", envelopeResult.err, envelopeResult.stderr)
	}
	obligation := decodeCommandJSONStructure(t, envelopeResult.stdout).(map[string]any)["result"].(map[string]any)["envelope"].(map[string]any)["obligations"].([]any)[0].(map[string]any)
	registration := obligation["registration"].(map[string]any)
	inputPath := filepath.Join(workspace, key+".json")
	writeActivationInput(t, inputPath, map[string]any{"registration_id": registration["id"], "registration_version": registration["version"], "obligation_id": obligation["id"], "run_id": runFromResponse(t, started.stdout)["id"], "repository": repository})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", key)
	lines := strings.Split(strings.TrimSpace(result.stdout), "\n")
	if len(lines) == 0 || lines[0] == "" {
		t.Fatalf("gate run output = %#v, want terminal frame", result)
	}
	var terminal gateTerminalFrame
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &terminal); err != nil {
		t.Fatalf("decode terminal frame: %v; frame=%q", err, lines[len(lines)-1])
	}
	return terminal, result
}
