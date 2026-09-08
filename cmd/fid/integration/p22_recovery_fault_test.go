//go:build integration

package integration

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
	"modernc.org/sqlite"
)

type p22RecoveryCrashInput struct {
	Mode    string
	Options workflow.RunOptions
	Request workflow.RunRecoveryRequest
}

func TestP22RecoverySurvivesActualResponseLossAndAppendProjectionFault(t *testing.T) {
	for _, mode := range []string{"before-append", "after-append", "lost-response"} {
		t.Run(mode, func(t *testing.T) {
			a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, p22FinalActions)
			p22AcceptedCall(t, a, "fault-baseline", "run.observe-budget", map[string]any{"value": 0})
			p22JourneyProof(t, a, "fault-retained")
			a.p22Observe(t, "terminated")
			a.p22Expire(t)
			fixture := p22ConformanceCase{Name: "crashed-recovery", Operation: "run.recover"}
			_, raw := a.p22Args(t, fixture)
			encoded, err := json.Marshal(raw)
			p22FaultOK(t, err)
			input := p22RecoveryCrashInput{Mode: mode, Options: workflow.RunOptions{WorkspaceRoot: a.workspace, MigrationsDir: filepath.Join(integrationRepositoryRoot(), "migrations"), ActorID: a.actor, IdempotencyKey: "p22-corpus-crashed-recovery"}}
			p22FaultOK(t, json.Unmarshal(encoded, &input.Request))
			encoded, err = json.Marshal(input)
			p22FaultOK(t, err)
			path := filepath.Join(t.TempDir(), "crash.json")
			p22FaultOK(t, os.WriteFile(path, encoded, 0o600))
			before := p22FaultProjection(t, a.workspace)
			beforeEvents := p21JournalEventCount(t, a.workspace)
			cmd := exec.Command(os.Args[0], "-test.run=^TestP22RecoveryCrashProcess$", "-test.count=1")
			cmd.Env = append(os.Environ(), "FID_P22_RECOVERY_CRASH="+path, "GOCOVERDIR="+t.TempDir())
			output, err := cmd.CombinedOutput()
			var exited *exec.ExitError
			if !errors.As(err, &exited) || exited.ExitCode() != 23 || len(output) != 0 {
				t.Fatalf("actual interrupted process: %v, %q", err, output)
			}
			afterFault := p22FaultProjection(t, a.workspace)
			if mode != "lost-response" && !reflect.DeepEqual(before, afterFault) {
				t.Fatal("failed transaction partially projected")
			}
			wantEvents := beforeEvents + 1
			if mode == "before-append" {
				wantEvents = beforeEvents
			}
			if p21JournalEventCount(t, a.workspace) != wantEvents {
				t.Fatal("fault occurred on wrong side of durable append")
			}
			conn := p22FaultDB(t, a.workspace)
			_, err = conn.Exec(`DROP TRIGGER IF EXISTS p22_recovery_fault`)
			p22FaultOK(t, err)
			p22FaultOK(t, conn.Close())
			// A fresh fid process catches up first. Retry then consumes the exact
			// original key, never a caller-chosen replacement successor identity.
			p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
			result := p22AcceptedCall(t, a, fixture.Name, fixture.Operation, nil)
			if p22ResultPath(result, "result.replayed") != (mode != "before-append") {
				t.Fatal("incorrect recovery replay disposition")
			}
			successor := p21PathString(t, result, "result", "run", "id")
			replay := p22AcceptedCall(t, a, fixture.Name, fixture.Operation, nil)
			if p22ResultPath(replay, "result.run.id") != successor || p21JournalEventCount(t, a.workspace) != beforeEvents+1 {
				t.Fatal("same key duplicated recovery")
			}
			after := p22FaultProjection(t, a.workspace)
			for table, rows := range before {
				switch table {
				case "runs", "run_runtime", "run_recoveries", "claims", "run_budget_bindings", "idempotency_results", "applied_events", "projection_meta":
					continue
				}
				if !reflect.DeepEqual(rows, after[table]) {
					t.Fatalf("recovery changed preserved %s context", table)
				}
			}
			if len(after["run_recoveries"]) != 1 || len(after["runs"]) != len(before["runs"])+1 || len(after["claims"]) != len(before["claims"])+1 {
				t.Fatal("missing or duplicate successor/Claim/link")
			}
			if len(before["envelope_obligations"]) == 0 || len(before["evidence_records"]) == 0 || len(before["budget_samples"]) == 0 {
				t.Fatal("preservation fixture lacks actual proof and budget context")
			}
			p22RejectSecondRecovery(t, a, input.Request, after, beforeEvents+1)
		})
	}
}

// This subprocess exits before returning any recovery result to its caller.
// A stateful SQLite scalar in a test-only trigger passes dry-run and fails the
// actual projection, so the after-append case uses a genuinely durable Event.
func TestP22RecoveryCrashProcess(t *testing.T) {
	path := os.Getenv("FID_P22_RECOVERY_CRASH")
	if path == "" {
		return
	}
	encoded, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var input p22RecoveryCrashInput
	p22FaultOK(t, json.Unmarshal(encoded, &input))
	var calls atomic.Int32
	if input.Mode != "lost-response" {
		failAt := int32(2)
		if input.Mode == "before-append" {
			failAt = 1
		}
		p22FaultOK(t, sqlite.RegisterScalarFunction("p22_recovery_interruption", 0, func(*sqlite.FunctionContext, []driver.Value) (driver.Value, error) {
			if calls.Add(1) == failAt {
				return nil, fmt.Errorf("injected recovery interruption")
			}
			return int64(0), nil
		}))
		conn := p22FaultDB(t, input.Options.WorkspaceRoot)
		_, err := conn.Exec(`CREATE TRIGGER p22_recovery_fault BEFORE INSERT ON run_recoveries BEGIN SELECT p22_recovery_interruption(); END`)
		p22FaultOK(t, err)
		p22FaultOK(t, conn.Close())
	}
	_, err = workflow.RecoverRun(input.Options, input.Request)
	if input.Mode == "lost-response" {
		p22FaultOK(t, err)
	} else {
		wantCalls := int32(2)
		if input.Mode == "before-append" {
			wantCalls = 1
		}
		if err == nil || !strings.Contains(err.Error(), "injected recovery interruption") || calls.Load() != wantCalls {
			t.Fatalf("wrong fault point: calls=%d, %v", calls.Load(), err)
		}
	}
	os.Exit(23)
}

func p22RejectSecondRecovery(t *testing.T, a *p22FIDRecoveryAdapter, original workflow.RunRecoveryRequest, before map[string][]string, events int) {
	t.Helper()
	for _, kind := range []string{"changed-key", "changed-digest"} {
		request := original
		key := "p22-corpus-crashed-recovery"
		if kind == "changed-key" {
			key = "second-recovery"
			prior := p21Run(t, a.binary, a.workspace, "--json", "run", "get", original.PriorRunID)
			request.ExpectedRunVersion = p22Number(p21Path(t, prior, "result", "run", "version"))
		} else {
			request.TTLSeconds++
		}
		if kind == "changed-digest" {
			_, err := workflow.RecoverRun(workflow.RunOptions{WorkspaceRoot: a.workspace, MigrationsDir: filepath.Join(integrationRepositoryRoot(), "migrations"), ActorID: a.actor, IdempotencyKey: key}, request)
			var conflict *journal.IdempotencyConflictError
			if !errors.As(err, &conflict) {
				t.Fatalf("changed request digest = %v, want idempotency conflict", err)
			}
			if p21JournalEventCount(t, a.workspace) != events || !reflect.DeepEqual(before, p22FaultProjection(t, a.workspace)) {
				t.Fatal("changed digest mutated recovery")
			}
			continue
		}
		encoded, err := json.Marshal(request)
		p22FaultOK(t, err)
		path := filepath.Join(t.TempDir(), "request.json")
		p22FaultOK(t, os.WriteFile(path, encoded, 0o600))
		response := a.p22Run(t, "run", "recover", "--actor", a.actor, "--input", path, "--idempotency-key", key)
		assertP22ProcessFailure(t, "run.recover", "recovery-already-linked", 4, response)
		if p21JournalEventCount(t, a.workspace) != events || !reflect.DeepEqual(before, p22FaultProjection(t, a.workspace)) {
			t.Fatal("competing replay changed recovery facts")
		}
	}
}

func p22FaultDB(t *testing.T, root string) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", filepath.Join(root, config.WorkspaceDirName, config.DatabaseFileName))
	p22FaultOK(t, err)
	return conn
}

func p22FaultProjection(t *testing.T, root string) map[string][]string {
	t.Helper()
	conn := p22FaultDB(t, root)
	defer conn.Close()
	names, err := conn.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	p22FaultOK(t, err)
	var tables []string
	for names.Next() {
		var name string
		p22FaultOK(t, names.Scan(&name))
		tables = append(tables, name)
	}
	p22FaultOK(t, names.Err())
	p22FaultOK(t, names.Close())
	result := map[string][]string{}
	for _, table := range tables {
		rows, err := conn.Query(`SELECT * FROM "` + table + `"`)
		p22FaultOK(t, err)
		columns, err := rows.Columns()
		p22FaultOK(t, err)
		result[table] = []string{}
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			p22FaultOK(t, rows.Scan(pointers...))
			encoded, err := json.Marshal(values)
			p22FaultOK(t, err)
			result[table] = append(result[table], string(encoded))
		}
		p22FaultOK(t, rows.Err())
		p22FaultOK(t, rows.Close())
		sort.Strings(result[table])
	}
	return result
}

func p22FaultOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
