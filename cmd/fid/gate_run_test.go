package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "gate-run-claim")
	repository := runRepositoryForTest(t, workspace)
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
			AttemptID string         `json:"attempt_id"`
			ReceiptID string         `json:"receipt_id"`
			Outcome   string         `json:"outcome"`
			Response  map[string]any `json:"response"`
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
	if err := assertCompleteReceiptBinding(terminal.Data.AttemptID, terminal.Data.ReceiptID, readImmutableReceiptFields(t, conn, terminal.Data.ReceiptID)); err != nil {
		t.Fatal(err)
	}
	assertReceiptMutationRejectedAndImmutable(t, conn, workspace, terminal.Data.ReceiptID, terminal.Data.Response)
	if _, err := conn.Exec(`INSERT INTO evidence_records (id, obligation_id, receipt_id, evidence_type, producer, trust_source, target_repository_json, validity, created_at) VALUES (?, ?, ?, 'receipt', 'caller', 'fid-executed', ?, 'valid', ?)`, "EVD-caller-reported", obligation["id"], terminal.Data.ReceiptID, `{"project_id":"project-1"}`, time.Now().UTC().Format(time.RFC3339Nano)); err == nil {
		t.Fatal("caller-reported Evidence insert succeeded, want rejection")
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM evidence_records`).Scan(&evidenceRecords); err != nil {
		t.Fatalf("count Evidence after caller proof rejection: %v", err)
	}
	if evidenceRecords != 0 {
		t.Fatalf("Evidence count after caller proof rejection = %d, want 0", evidenceRecords)
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

func TestGateRunAttemptAppendFailureDoesNotSpawn(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/touch", "attempt-append-spawned"}
	})
	inputPath, key := gateRunInputForTest(t, workspace, issueID, envelopeID, "gate-run-append-failure")
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open append-failure database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	if _, err := conn.Exec(`CREATE TRIGGER reject_attempt_append BEFORE INSERT ON execution_attempts BEGIN SELECT RAISE(ABORT, 'injected Attempt append failure'); END`); err != nil {
		t.Fatalf("inject Attempt append failure: %v", err)
	}
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", key)
	if result.err == nil || cli.ExitCode(result.err) != 2 || result.stdout != "" {
		t.Fatalf("append-failure gate run = %#v, want handled pre-spawn input failure", result)
	}
	if _, err := os.Stat(filepath.Join(workspace, "attempt-append-spawned")); !os.IsNotExist(err) {
		t.Fatalf("append-failure spawn marker stat = %v, want no spawned process side effect", err)
	}
	var attempts int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts`).Scan(&attempts); err != nil {
		t.Fatalf("count append-failure Attempts: %v", err)
	}
	if attempts != 0 {
		t.Fatalf("append-failure Attempts = %d, want 0", attempts)
	}
}

func TestBudgetDenialReportsCommittedStopWithoutStartFrame(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	inputPath, _ := gateRunInputForTest(t, workspace, issueID, envelopeID, "budget-denial")
	for _, key := range []string{"budget-denial-one", "budget-denial-two", "budget-denial-three"} {
		result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", key)
		if result.err != nil {
			t.Fatalf("admitted budget run %q = %#v", key, result)
		}
	}

	denied := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "budget-denial-four")
	if denied.err == nil || denied.stderr != "" || cli.ExitCode(denied.err) != 4 {
		t.Fatalf("budget denial = %#v (%v), want frame-only committed denial with exit 4", denied, denied.err)
	}
	lines := strings.Split(strings.TrimSpace(denied.stdout), "\n")
	if len(lines) != 1 {
		t.Fatalf("budget denial frames = %d, want only one fid-error terminal", len(lines))
	}
	var terminal gateTerminalFrame
	if err := json.Unmarshal([]byte(lines[0]), &terminal); err != nil {
		t.Fatalf("decode budget denial terminal: %v", err)
	}
	if terminal.Stream != "terminal" || terminal.Data.AttemptID != "" || terminal.Data.ReceiptID != "" || terminal.Data.Outcome != "fid-error" {
		t.Fatalf("budget denial terminal = %#v, want no Attempt or Receipt", terminal)
	}
	errorValue, ok := terminal.Data.Response["error"].(map[string]any)
	if !ok || errorValue["code"] != "budget-exhausted" {
		t.Fatalf("budget denial response = %#v, want budget-exhausted error", terminal.Data.Response)
	}
	details, ok := errorValue["details"].(map[string]any)
	if !ok || details["event"] == nil || len(details["stop_ids"].([]any)) == 0 || len(details["attention_ids"].([]any)) == 0 {
		t.Fatalf("budget denial details = %#v, want committed Event, Stop, and Attention IDs", details)
	}
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open budget-denial database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var attempts int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts`).Scan(&attempts); err != nil {
		t.Fatalf("count budget-denial Attempts: %v", err)
	}
	if attempts != 3 {
		t.Fatalf("budget denial Attempt count = %d, want no fourth Attempt", attempts)
	}
}

func TestGateRunReceiptAppendFailureEmitsFrameOnlyFidError(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/echo", "started"}
	})
	inputPath, key := gateRunInputForTest(t, workspace, issueID, envelopeID, "gate-run-receipt-failure")
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open receipt-failure database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	if _, err := conn.Exec(`CREATE TRIGGER reject_receipt_append BEFORE INSERT ON execution_receipts BEGIN SELECT RAISE(ABORT, 'injected Receipt append failure'); END`); err != nil {
		t.Fatalf("inject Receipt append failure: %v", err)
	}
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", key)
	if result.err == nil || result.stderr != "" || cli.ExitCode(result.err) != 7 {
		t.Fatalf("receipt-failure gate run = %#v, want frame-only fid-error with exit 7", result)
	}
	lines := strings.Split(strings.TrimSpace(result.stdout), "\n")
	if len(lines) != 2 {
		t.Fatalf("receipt-failure frames = %d, want one data frame and one terminal", len(lines))
	}
	var terminal gateTerminalFrame
	if err := json.Unmarshal([]byte(lines[1]), &terminal); err != nil {
		t.Fatalf("decode fid-error terminal: %v", err)
	}
	if terminal.Stream != "terminal" || terminal.Data.ReceiptID != "" || terminal.Data.Outcome != "fid-error" || terminal.Data.ExitCode != nil || terminal.Data.Signal != nil || terminal.Data.TimedOut || terminal.Data.Replayed {
		t.Fatalf("fid-error terminal = %#v, want null Receipt and exact frame-only outcome", terminal)
	}
	var attempts, receipts, idempotency int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts`).Scan(&attempts); err != nil {
		t.Fatalf("count fid-error Attempts: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_receipts`).Scan(&receipts); err != nil {
		t.Fatalf("count fid-error Receipts: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM idempotency_results WHERE operation_scope LIKE 'gate.run:%'`).Scan(&idempotency); err != nil {
		t.Fatalf("count fid-error idempotency results: %v", err)
	}
	if attempts != 1 || receipts != 0 || idempotency != 1 {
		t.Fatalf("fid-error Attempt/Receipt/idempotency counts = %d/%d/%d, want 1/0/1", attempts, receipts, idempotency)
	}
}

func TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/echo", "started"}
	})
	inputPath, key := gateRunInputForTest(t, workspace, issueID, envelopeID, "killed-fid")
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open killed-fid database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	if _, err := conn.Exec(`CREATE TRIGGER lose_terminal_observation BEFORE INSERT ON execution_receipts BEGIN SELECT RAISE(ABORT, 'injected fid termination before Receipt'); END`); err != nil {
		t.Fatalf("inject terminal observation loss: %v", err)
	}

	first := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", key)
	if first.err == nil || cli.ExitCode(first.err) != 7 {
		t.Fatalf("initial interrupted gate run = %#v, want frame-only fid error", first)
	}
	retry := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", key)
	if retry.err == nil || cli.ExitCode(retry.err) != 2 || retry.stdout != "" {
		t.Fatalf("same-key retry after lost terminal observation = %#v, want rejection without spawn", retry)
	}
	var attempts, receipts int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts`).Scan(&attempts); err != nil {
		t.Fatalf("count interrupted Attempts: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_receipts`).Scan(&receipts); err != nil {
		t.Fatalf("count interrupted Receipts: %v", err)
	}
	if attempts != 1 || receipts != 0 {
		t.Fatalf("interrupted Attempt/Receipt counts after same-key retry = %d/%d, want 1/0", attempts, receipts)
	}
}

func TestGateRunRecordsExactNonSuccessTerminalOutcomes(t *testing.T) {
	t.Parallel()
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
			definition["argv"] = []string{os.Args[0], "-test.run=^TestGateRunHelperProcess$", "--"}
			definition["timeout_ms"] = 5_000
			definition["environment"] = []any{map[string]any{"name": "FID_GATE_HELPER", "source": "literal", "value": "timeout"}}
		}, "timed-out", nil, true, true},
		{"signaled", func(definition map[string]any) {
			definition["argv"] = []string{os.Args[0], "-test.run=^TestGateRunHelperProcess$", "--"}
			definition["timeout_ms"] = 10_000
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
			replayed := executeCommandForTest(workspace,
				"--project", workspace, "--json", "gate", "run", "--actor", "actor-builder",
				"--input", filepath.Join(workspace, "gate-run-"+strings.ReplaceAll(testCase.name, " ", "-")+".json"), "--idempotency-key", "gate-run-"+strings.ReplaceAll(testCase.name, " ", "-"),
			)
			if replayed.err == nil || replayed.stderr != "" || cli.ExitCode(replayed.err) != 8 {
				t.Fatalf("replayed gate run %s = %#v (%v), want terminal failure with exit 8", testCase.name, replayed, replayed.err)
			}
			replayedLines := strings.Split(strings.TrimSpace(replayed.stdout), "\n")
			if len(replayedLines) != 1 {
				t.Fatalf("replayed gate run %s frames = %d, want terminal frame only", testCase.name, len(replayedLines))
			}
			var replayTerminal gateTerminalFrame
			if err := json.Unmarshal([]byte(replayedLines[0]), &replayTerminal); err != nil {
				t.Fatalf("decode replayed %s terminal: %v", testCase.name, err)
			}
			if replayTerminal.Stream != "terminal" || !replayTerminal.Data.Replayed || replayTerminal.Data.ReceiptID != terminal.Data.ReceiptID || replayTerminal.Data.Outcome != testCase.wantOutcome || replayTerminal.Data.ExitCode != testCase.wantExitCode || (replayTerminal.Data.Signal != nil) != testCase.wantSignal || replayTerminal.Data.TimedOut != testCase.wantTimedOut || !reflect.DeepEqual(replayTerminal.Data.Response, terminal.Data.Response) {
				t.Fatalf("replayed %s terminal = %#v, want original Receipt and exact replayed terminal", testCase.name, replayTerminal)
			}
			conn, err := fiddb.Open(filepath.Join(workspace, ".fid", config.DatabaseFileName))
			if err != nil {
				t.Fatalf("open replay database: %v", err)
			}
			defer func() { _ = conn.Close() }()
			var attempts, receipts, idempotency int
			if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts`).Scan(&attempts); err != nil {
				t.Fatalf("count replayed Attempts: %v", err)
			}
			if err := conn.QueryRow(`SELECT COUNT(*) FROM execution_receipts`).Scan(&receipts); err != nil {
				t.Fatalf("count replayed Receipts: %v", err)
			}
			if err := conn.QueryRow(`SELECT COUNT(*) FROM idempotency_results WHERE operation_scope LIKE 'gate.run:%' AND result_entity_id = ?`, terminal.Data.ReceiptID).Scan(&idempotency); err != nil {
				t.Fatalf("count replayed idempotency records: %v", err)
			}
			if attempts != 1 || receipts != 1 || idempotency != 1 {
				t.Fatalf("replayed %s Attempt/Receipt/idempotency counts = %d/%d/%d, want 1/1/1", testCase.name, attempts, receipts, idempotency)
			}
		})
	}
}

func TestGateRunHelperProcess(t *testing.T) {
	switch os.Getenv("FID_GATE_HELPER") {
	case "signal":
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
		select {}
	case "timeout":
		_, _ = os.Stdout.WriteString("ready\n")
		select {}
	}
}

type gateTerminalFrame struct {
	Stream string `json:"stream"`
	Data   struct {
		AttemptID string         `json:"attempt_id"`
		ReceiptID string         `json:"receipt_id"`
		Outcome   string         `json:"outcome"`
		ExitCode  any            `json:"exit_code"`
		Signal    any            `json:"signal"`
		TimedOut  bool           `json:"timed_out"`
		Replayed  bool           `json:"replayed"`
		Response  map[string]any `json:"response"`
	} `json:"data"`
}

type immutableReceiptFields struct {
	ID                  string
	AttemptID           string
	RegistrationID      string
	RegistrationVersion int
	ObligationID        string
	RunID               string
	ActorID             string
	TrustSource         string
	Outcome             string
	ExitCode            *int
	Signal              *string
	TimedOut            bool
	StartedAt           string
	EndedAt             string
	TerminalRepository  string
	Toolchain           string
	OutputExcerpt       string
	OutputDigest        string
	OutputBytes         int64
	OutputTruncated     bool
	RedactionApplied    bool
}

func assertCompleteReceiptBinding(attemptID, receiptID string, receipt immutableReceiptFields) error {
	if receipt.ID != receiptID || receipt.AttemptID != attemptID {
		return fmt.Errorf("Receipt identity binding = %#v, want receipt %q bound to Attempt %q", receipt, receiptID, attemptID)
	}
	if receipt.RegistrationID == "" || receipt.RegistrationVersion < 1 || receipt.ObligationID == "" || receipt.RunID == "" || receipt.ActorID == "" || receipt.TrustSource != "fid-executed" {
		return fmt.Errorf("Receipt execution binding is incomplete: %#v", receipt)
	}
	if receipt.Outcome != "succeeded" || receipt.ExitCode == nil || *receipt.ExitCode != 0 || receipt.Signal != nil || receipt.TimedOut || receipt.StartedAt == "" || receipt.EndedAt == "" {
		return fmt.Errorf("Receipt terminal fields = %#v, want successful observed terminal", receipt)
	}
	if receipt.TerminalRepository == "" || receipt.Toolchain == "" || receipt.OutputDigest == "" || receipt.OutputBytes < 0 {
		return fmt.Errorf("Receipt terminal binding is incomplete: %#v", receipt)
	}
	return nil
}

func assertReceiptMutationRejectedAndImmutable(t *testing.T, conn *sql.DB, workspace, receiptID string, response map[string]any) {
	t.Helper()
	before := readImmutableReceiptFields(t, conn, receiptID)
	event, ok := response["event"].(map[string]any)
	if !ok {
		t.Fatalf("successful terminal response event = %#v, want EventRef", response["event"])
	}
	eventID, ok := event["id"].(string)
	if !ok {
		t.Fatalf("successful terminal response Event ID = %#v, want string", event["id"])
	}
	source, err := journal.FindEvent(filepath.Join(workspace, ".fid", "journal"), eventID)
	if err != nil {
		t.Fatalf("read Receipt Event: %v", err)
	}
	mutation := source.Event
	mutation.EventID = "EVT-20991231-999999"
	mutation.EventType = "gate.receipt-mutated"
	mutation.OccurredAt = time.Now().UTC()
	mutation.Aggregate.Sequence++
	mutation.Idempotency = nil
	mutation.Effects = []journal.Effect{{Op: "update", Table: "execution_receipts", PK: map[string]any{"id": receiptID}, Values: map[string]any{"terminal_status": "failed"}}}
	if err := journal.DryRunProjectEvent(conn, mutation); err == nil {
		t.Fatal("Receipt mutation projection succeeded, want immutable Receipt rejection")
	}
	after := readImmutableReceiptFields(t, conn, receiptID)
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("Receipt immutable fields after rejected mutation = %#v, want %#v", after, before)
	}
}

func readImmutableReceiptFields(t *testing.T, conn *sql.DB, receiptID string) immutableReceiptFields {
	t.Helper()
	var receipt immutableReceiptFields
	err := conn.QueryRow(`SELECT r.id, r.attempt_id, a.registration_id, a.registration_version, a.obligation_id, a.run_id, a.actor_id, r.trust_source, r.terminal_status, r.exit_code, r.signal, r.timed_out, r.started_at, r.ended_at, r.terminal_repository_json, r.toolchain_json, r.output_excerpt, r.output_digest, r.output_bytes, r.output_truncated, r.redaction_applied FROM execution_receipts r JOIN execution_attempts a ON a.id = r.attempt_id WHERE r.id = ?`, receiptID).Scan(&receipt.ID, &receipt.AttemptID, &receipt.RegistrationID, &receipt.RegistrationVersion, &receipt.ObligationID, &receipt.RunID, &receipt.ActorID, &receipt.TrustSource, &receipt.Outcome, &receipt.ExitCode, &receipt.Signal, &receipt.TimedOut, &receipt.StartedAt, &receipt.EndedAt, &receipt.TerminalRepository, &receipt.Toolchain, &receipt.OutputExcerpt, &receipt.OutputDigest, &receipt.OutputBytes, &receipt.OutputTruncated, &receipt.RedactionApplied)
	if err != nil {
		t.Fatalf("read immutable Receipt fields: %v", err)
	}
	return receipt
}

func gateRunInputForTest(t *testing.T, workspace, issueID, envelopeID, key string) (string, string) {
	t.Helper()
	repository := runRepositoryForTest(t, workspace)
	binding := envelopeBindingForWorkspace(t, workspace)
	runID := activeRunIDForGateTest(t, workspace, issueID, envelopeID, key, repository, binding)
	inputPath := filepath.Join(workspace, key+".json")
	writeActivationInput(t, inputPath, map[string]any{"registration_id": binding.registrationID, "registration_version": binding.registrationVersion, "obligation_id": binding.obligationID, "run_id": runID, "repository": repository})
	return inputPath, key
}

func executeGateRunForTest(t *testing.T, workspace, issueID, envelopeID, key string) (gateTerminalFrame, commandExecutionResult) {
	t.Helper()
	repository := runRepositoryForTest(t, workspace)
	binding := envelopeBindingForWorkspace(t, workspace)
	runID := activeRunIDForGateTest(t, workspace, issueID, envelopeID, key, repository, binding)
	inputPath := filepath.Join(workspace, key+".json")
	writeActivationInput(t, inputPath, map[string]any{"registration_id": binding.registrationID, "registration_version": binding.registrationVersion, "obligation_id": binding.obligationID, "run_id": runID, "repository": repository})
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

func activeRunIDForGateTest(t *testing.T, workspace, issueID, envelopeID, key string, repository map[string]any, binding envelopeBindingForTest) string {
	t.Helper()
	if binding.runID != "" {
		return binding.runID
	}
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, key+"-claim")
	runPath := filepath.Join(workspace, key+"-run.json")
	writeActivationInput(t, runPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", runPath, "--idempotency-key", key+"-start")
	if started.err != nil {
		t.Fatalf("start Run: %v; stderr=%q", started.err, started.stderr)
	}
	return runFromResponse(t, started.stdout)["id"].(string)
}
