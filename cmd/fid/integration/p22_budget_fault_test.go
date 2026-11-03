//go:build integration

package integration

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/mbell-nzl/fid/internal/workflow"
	"modernc.org/sqlite"
)

type p22BudgetCrashInput struct {
	Kind, Mode string
	Options    workflow.RunOptions
	Sample     workflow.RunBudgetObserveRequest
	Gate       evidence.GateRunOptions
}

func TestP22BudgetBoundarySurvivesAppendAndProjectionFaults(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"external", "attempt", "receipt"} {
		for _, mode := range []string{"before-append", "after-append", "normal"} {
			t.Run(kind+"/"+mode, func(t *testing.T) { p22BudgetBoundaryCase(t, kind, mode) })
		}
	}
}

func p22BudgetBoundaryCase(t *testing.T, kind, mode string) {
	t.Helper()
	a := p22BudgetFaultFixture(t, kind)
	input := p22BudgetCrashInput{Kind: kind, Mode: mode, Options: workflow.RunOptions{WorkspaceRoot: a.workspace, MigrationsDir: filepath.Join(integrationRepositoryRoot(), "migrations"), ActorID: a.actor, IdempotencyKey: "budget-boundary"}}
	input.Gate = p22BudgetGateOptions(t, a, "budget-boundary")
	if kind == "external" {
		_, raw := a.p22Args(t, p22ConformanceCase{Name: "budget-boundary", Operation: "run.observe-budget", Input: map[string]any{"value": 5}})
		encoded, err := json.Marshal(raw)
		p22FaultOK(t, err)
		p22FaultOK(t, json.Unmarshal(encoded, &input.Sample))
	}
	before, events := p22FaultProjection(t, a.workspace), p21JournalEventCount(t, a.workspace)
	spawns := p22BudgetSpawnCount(t, a)
	if mode == "normal" {
		p22BudgetExecute(t, input)
	} else {
		p22BudgetCrash(t, input)
		after := p22FaultProjection(t, a.workspace)
		for table, rows := range before {
			if kind == "receipt" {
				switch table {
				case "execution_attempts", "budget_samples", "run_runtime", "idempotency_results", "applied_events", "projection_meta":
					continue
				}
			}
			if !reflect.DeepEqual(rows, after[table]) {
				t.Fatalf("fault partially changed %s", table)
			}
		}
		extra := 0
		if kind == "receipt" {
			extra++
		}
		if mode == "after-append" {
			extra++
		}
		if p21JournalEventCount(t, a.workspace) != events+extra {
			t.Fatal("fault occurred on wrong side of append")
		}
		if len(after["execution_receipts"]) != len(before["execution_receipts"]) || len(after["stops"]) != 0 || len(after["attention_items"]) != 0 || len(after["attention_stop_links"]) != 0 {
			t.Fatal("partial Receipt or boundary facts")
		}
		wantSpawns := spawns
		if kind == "receipt" {
			wantSpawns++
		}
		if p22BudgetSpawnCount(t, a) != wantSpawns {
			t.Fatal("wrong spawn count at fault")
		}
		conn := p22FaultDB(t, a.workspace)
		_, err := conn.Exec(`DROP TRIGGER p22_budget_fault`)
		p22FaultOK(t, err)
		p22FaultOK(t, conn.Close())
		p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
		if mode == "before-append" && kind != "receipt" {
			p22BudgetExecute(t, input)
		}
	}
	complete := mode == "normal" || kind == "external" || (kind == "attempt" && mode == "before-append") || (kind == "receipt" && mode == "after-append")
	// A pre-append Receipt failure retains an abandoned Attempt. It cannot invent
	// a Receipt or failure charge on replay. The other paths completed a boundary.
	boundary := !(kind == "receipt" && mode == "before-append")
	p22AssertBudgetTransaction(t, a, before, kind, complete, boundary)
	wantSpawns := spawns
	if kind == "receipt" || (kind == "attempt" && complete) {
		wantSpawns++
	}
	if p22BudgetSpawnCount(t, a) != wantSpawns {
		t.Fatal("boundary spawned an unexpected process count")
	}
	stable := p22FaultProjection(t, a.workspace)
	stableEvents, stableSpawns := p21JournalEventCount(t, a.workspace), p22BudgetSpawnCount(t, a)
	if kind == "external" {
		replay, err := workflow.ObserveRunBudget(input.Options, input.Sample)
		p22FaultOK(t, err)
		if !replay.Replayed {
			t.Fatal("sample retry did not replay")
		}
	} else {
		replay, err := evidence.RunGate(context.Background(), input.Gate)
		if complete {
			p22FaultOK(t, err)
		} else if err == nil || !strings.Contains(err.Error(), "has no observed terminal Receipt") {
			t.Fatalf("abandoned replay = %v", err)
		}
		if !replay.Replayed || (replay.Receipt != nil) != complete {
			t.Fatal("Attempt/Receipt replay disposition changed")
		}
	}
	if !reflect.DeepEqual(stable, p22FaultProjection(t, a.workspace)) || stableEvents != p21JournalEventCount(t, a.workspace) || stableSpawns != p22BudgetSpawnCount(t, a) {
		t.Fatal("exact replay duplicated facts or spawned")
	}
	if boundary {
		p22AssertNextBudgetGateDenied(t, a, stableSpawns)
	}
	// Full rebuild reconstructs the same durable identities and transaction state.
	stable = p22FaultProjection(t, a.workspace)
	p21Run(t, a.binary, a.workspace, "--json", "rebuild")
	rebuilt := p22FaultProjection(t, a.workspace)
	for table, rows := range stable {
		if table == "projection_meta" {
			continue
		}
		if !reflect.DeepEqual(rows, rebuilt[table]) {
			t.Fatalf("rebuild changed %s", table)
		}
	}
}

func p22BudgetFaultFixture(t *testing.T, kind string) *p22FIDRecoveryAdapter {
	t.Helper()
	exit := "0"
	cfg := p22ExecutionFixture{}
	if kind == "receipt" {
		exit = "1"
		cfg.ExtraBudgets = []any{map[string]any{"name": "failures", "scope": "envelope", "unit": "failures", "maximum": 1, "warning_threshold": nil, "exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil}}
	}
	cfg.Argv = []string{"/bin/sh", "-c", "printf x >> .fid/fq5-spawn; exit " + exit}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, p22FinalActions, cfg)
	p22AcceptedCall(t, a, "budget-initial", "run.observe-budget", map[string]any{"value": 0})
	if kind == "external" {
		p22AcceptedCall(t, a, "budget-below", "run.observe-budget", map[string]any{"value": 5})
	}
	if kind == "attempt" {
		for i := 0; i < 2; i++ {
			result, err := evidence.RunGate(context.Background(), p22BudgetGateOptions(t, a, fmt.Sprintf("below-%d", i)))
			p22FaultOK(t, err)
			if result.Receipt == nil || result.Receipt.Outcome != "succeeded" {
				t.Fatal("below-boundary process did not finish")
			}
		}
	}
	rows := p22FaultProjection(t, a.workspace)
	if len(rows["stops"]) != 0 || len(rows["attention_items"]) != 0 {
		t.Fatal("below boundary created exhaustion facts")
	}
	return a
}

func p22BudgetGateOptions(t *testing.T, a *p22FIDRecoveryAdapter, key string) evidence.GateRunOptions {
	t.Helper()
	_, raw := a.p22BuildArgs(t, p22ConformanceCase{Operation: "gate.run"})
	encoded, err := json.Marshal(raw)
	p22FaultOK(t, err)
	options := evidence.GateRunOptions{WorkspaceRoot: a.workspace, MigrationsDir: filepath.Join(integrationRepositoryRoot(), "migrations"), ActorID: a.actor, IdempotencyKey: key}
	p22FaultOK(t, json.Unmarshal(encoded, &options.Request))
	return options
}

func p22BudgetExecute(t *testing.T, input p22BudgetCrashInput) {
	t.Helper()
	if input.Kind == "external" {
		_, err := workflow.ObserveRunBudget(input.Options, input.Sample)
		p22FaultOK(t, err)
		return
	}
	result, err := evidence.RunGate(context.Background(), input.Gate)
	p22FaultOK(t, err)
	if result.Receipt == nil {
		t.Fatal("final allowed process lost its Receipt")
	}
}

func p22BudgetCrash(t *testing.T, input p22BudgetCrashInput) {
	t.Helper()
	encoded, err := json.Marshal(input)
	p22FaultOK(t, err)
	path := filepath.Join(t.TempDir(), "budget.json")
	p22FaultOK(t, os.WriteFile(path, encoded, 0o600))
	cmd := exec.Command(os.Args[0], "-test.run=^TestP22BudgetCrashProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), "FID_P22_BUDGET_CRASH="+path, "GOCOVERDIR="+t.TempDir())
	output, err := cmd.CombinedOutput()
	var exited *exec.ExitError
	if !errors.As(err, &exited) || exited.ExitCode() != 23 || len(output) != 0 {
		t.Fatalf("budget interruption: %v %q", err, output)
	}
}

func TestP22BudgetCrashProcess(t *testing.T) {
	t.Parallel()
	path := os.Getenv("FID_P22_BUDGET_CRASH")
	if path == "" {
		return
	}
	encoded, err := os.ReadFile(path)
	p22FaultOK(t, err)
	var input p22BudgetCrashInput
	p22FaultOK(t, json.Unmarshal(encoded, &input))
	failAt := int32(2)
	if input.Mode == "before-append" {
		failAt = 1
	}
	var calls atomic.Int32
	p22FaultOK(t, sqlite.RegisterScalarFunction("p22_budget_interruption", 0, func(*sqlite.FunctionContext, []driver.Value) (driver.Value, error) {
		if calls.Add(1) == failAt {
			return nil, fmt.Errorf("injected budget interruption")
		}
		return int64(0), nil
	}))
	name := map[string]string{"external": "tokens", "attempt": "attempts", "receipt": "failures"}[input.Kind]
	conn := p22FaultDB(t, input.Options.WorkspaceRoot)
	_, err = conn.Exec(`CREATE TRIGGER p22_budget_fault BEFORE INSERT ON budget_samples WHEN NEW.budget_name = '` + name + `' BEGIN SELECT p22_budget_interruption(); END`)
	p22FaultOK(t, err)
	p22FaultOK(t, conn.Close())
	if input.Kind == "external" {
		_, err = workflow.ObserveRunBudget(input.Options, input.Sample)
	} else {
		_, err = evidence.RunGate(context.Background(), input.Gate)
	}
	if err == nil || !strings.Contains(err.Error(), "injected budget interruption") || calls.Load() != failAt {
		t.Fatalf("wrong budget fault: calls=%d %v", calls.Load(), err)
	}
	os.Exit(23)
}

func p22BudgetSpawnCount(t *testing.T, a *p22FIDRecoveryAdapter) int {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(config.LedgerDir(a.workspace), "fq5-spawn"))
	if errors.Is(err, os.ErrNotExist) {
		return 0
	}
	p22FaultOK(t, err)
	return len(data)
}

func p22AssertBudgetTransaction(t *testing.T, a *p22FIDRecoveryAdapter, before map[string][]string, kind string, complete, boundary bool) {
	t.Helper()
	after := p22FaultProjection(t, a.workspace)
	types := p21JournalEventTypes(t, a.workspace)
	if kind == "external" {
		if types[len(types)-1] != "run.budget-observed" {
			t.Fatal("missing final sample Event")
		}
	} else if complete {
		if !reflect.DeepEqual(types[len(types)-2:], []string{"gate.execution-started", "gate.receipt-recorded"}) {
			t.Fatal("final allowed Attempt/Receipt order changed")
		}
	} else if types[len(types)-1] != "gate.execution-started" {
		t.Fatal("abandoned Attempt gained a terminal Event")
	}
	samples := len(before["budget_samples"]) + 1
	if kind == "receipt" && complete {
		samples++
	}
	if len(after["budget_samples"]) != samples {
		t.Fatal("partial or duplicate budget samples")
	}
	attempts, receipts := len(before["execution_attempts"]), len(before["execution_receipts"])
	if kind != "external" {
		attempts++
		if complete {
			receipts++
		}
	}
	if len(after["execution_attempts"]) != attempts || len(after["execution_receipts"]) != receipts {
		t.Fatal("wrong final Attempt/Receipt counts")
	}
	wanted := 0
	if boundary {
		wanted = 1
	}
	for _, table := range []string{"stops", "attention_items", "attention_stop_links"} {
		if len(after[table]) != wanted {
			t.Fatalf("%s count=%d want=%d", table, len(after[table]), wanted)
		}
	}
	conn := p22FaultDB(t, a.workspace)
	defer conn.Close()
	name := map[string]string{"external": "tokens", "attempt": "attempts", "receipt": "failures"}[kind]
	var total int
	p22FaultOK(t, conn.QueryRow(`SELECT COALESCE(SUM(value),0) FROM budget_samples WHERE budget_name=?`, name).Scan(&total))
	want := map[string]int{"external": 10, "attempt": 3, "receipt": 1}[kind]
	if !boundary {
		want = 0
	}
	if total != want {
		t.Fatalf("budget total=%d want=%d", total, want)
	}
	if boundary {
		var bound int
		p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM stops s JOIN budget_samples b ON b.created_event_id=s.created_event_id JOIN attention_stop_links l ON l.stop_id=s.id JOIN attention_items a ON a.id=l.attention_id WHERE b.budget_name=? AND s.code='budget-exhausted' AND a.status='open' AND s.subject_id=?`, name, a.runID).Scan(&bound))
		if bound != 1 {
			t.Fatal("sample/Stop/Attention do not form one boundary transaction")
		}
	}
}

func p22AssertNextBudgetGateDenied(t *testing.T, a *p22FIDRecoveryAdapter, spawns int) {
	t.Helper()
	before := p22FaultProjection(t, a.workspace)
	_, raw := a.p22BuildArgs(t, p22ConformanceCase{Operation: "gate.run"})
	response := a.p22RunInput(t, "gate", "run", a.actor, raw)
	if response.ExitCode != 4 || response.Stderr != "" {
		t.Fatalf("budget denial: %#v", response)
	}
	var frame map[string]any
	p22FaultOK(t, json.Unmarshal([]byte(response.Stdout), &frame))
	for path, want := range map[string]any{"frame_schema": "fid.live.v1", "stream": "terminal", "sequence": float64(1), "data.attempt_id": nil, "data.receipt_id": nil, "data.outcome": "fid-error", "data.replayed": false, "data.response.ok": false, "data.response.operation": "gate.run", "data.response.error.code": "budget-exhausted", "data.response.error.category": "precondition", "data.response.error.details.budget.remaining": float64(0), "data.response.error.details.stop_ids.#": float64(1), "data.response.error.details.attention_ids.#": float64(1)} {
		if !reflect.DeepEqual(p22ResultPath(frame, path), want) {
			t.Fatalf("denial %s = %#v want %#v", path, p22ResultPath(frame, path), want)
		}
	}
	after := p22FaultProjection(t, a.workspace)
	for _, table := range []string{"execution_attempts", "execution_receipts", "budget_samples", "attention_items"} {
		if !reflect.DeepEqual(before[table], after[table]) {
			t.Fatalf("denied execution changed %s", table)
		}
	}
	for _, table := range []string{"stops", "attention_stop_links"} {
		if len(after[table]) != len(before[table])+1 {
			t.Fatalf("denial did not append one %s fact", table)
		}
		for _, old := range before[table] {
			found := false
			for _, row := range after[table] {
				if row == old {
					found = true
				}
			}
			if !found {
				t.Fatalf("denial rewrote existing %s", table)
			}
		}
	}
	if p22BudgetSpawnCount(t, a) != spawns {
		t.Fatal("denied execution spawned")
	}
	events := p21JournalEventCount(t, a.workspace)
	replay := a.p22RunInput(t, "gate", "run", a.actor, raw)
	var replayFrame map[string]any
	p22FaultOK(t, json.Unmarshal([]byte(replay.Stdout), &replayFrame))
	if replay.ExitCode != 4 || replay.Stderr != "" || p22ResultPath(replayFrame, "data.replayed") != true || p22ResultPath(replayFrame, "data.response.error.details.event.id") != p22ResultPath(frame, "data.response.error.details.event.id") {
		t.Fatal("denial replay lost original identity")
	}
	if !reflect.DeepEqual(after, p22FaultProjection(t, a.workspace)) || p21JournalEventCount(t, a.workspace) != events || p22BudgetSpawnCount(t, a) != spawns {
		t.Fatal("denial replay duplicated facts or spawned")
	}
}
