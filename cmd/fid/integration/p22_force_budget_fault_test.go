//go:build integration

package integration

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestP22ForceRecoveryBudgetSurvivesProcessFaults(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"before-append", "after-append", "lost-response"} {
		t.Run(mode, func(t *testing.T) {
			a := p22ForceBudgetFixture(t, 1)
			force := p22ForceBudgetInput(t, a, "fault-force")
			fixture := p22ConformanceCase{Name: "crashed-recovery", Operation: "run.recover", Input: force}
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
			p22AssertInterventions(t, result, 1)
			p22AssertInterventions(t, replay, 1)
			for table, want := range map[string]int{"budget_samples": 1, "stops": 1, "attention_items": 1, "attention_stop_links": 1, "run_recoveries": 1} {
				if len(after[table]) != want {
					t.Fatalf("%s count = %d, want %d", table, len(after[table]), want)
				}
			}
			p21Run(t, a.binary, a.workspace, "--json", "rebuild")
			if !reflect.DeepEqual(after, p22FaultProjection(t, a.workspace)) {
				t.Fatal("rebuild changed committed force accounting")
			}

		})
	}
}
