//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/evidence"
)

func TestP22CloseRunUsesObservedRepositoryAndQualifyingProof(t *testing.T) {
	for _, kind := range []string{"qualifying", "missing", "stale", "failed", "outstanding", "forged-old-target", "live", "unknown", "missing-observation", "null-observation", "stale-version", "nil-base", "missing-base"} {
		t.Run(kind, func(t *testing.T) {
			a := p22CloseProofFixture(t, kind)
			state := "terminated"
			if kind == "live" || kind == "unknown" {
				state = kind
			}
			a.p22Observe(t, state)
			if kind == "stale" || kind == "forged-old-target" {
				p22Write(t, a.workspace, "base.txt", "changed proof target\n")
				if kind == "stale" {
					a.repository = p22CorpusRepository(t, a.workspace, a.projectID)
				}
			}
			if kind == "nil-base" {
				a.repository["base_oid"] = nil
			}
			if kind == "missing-base" {
				a.repository["base_oid"] = strings.Repeat("0", 40)
			}
			input := map[string]any{"run_id": a.runID, "expected_version": a.runVersion, "observation_id": a.observation, "repository": a.repository}
			if kind == "missing-observation" {
				delete(input, "observation_id")
			}
			if kind == "null-observation" {
				input["observation_id"] = nil
			}
			if kind == "stale-version" {
				input["expected_version"] = 99
			}
			before, events := p22FaultProjection(t, a.workspace), p21JournalEventCount(t, a.workspace)
			response := a.p22RunInput(t, "run", "close", a.actor, input)
			wantCode, wantExit := "", 0
			switch kind {
			case "outstanding", "missing-observation", "null-observation":
				wantCode, wantExit = "invalid-argument", 2
			case "forged-old-target", "missing-base":
				wantCode, wantExit = "recovery-state-conflict", 4
			case "live":
				wantCode, wantExit = "executor-still-live", 4
			case "unknown":
				wantCode, wantExit = "executor-liveness-unknown", 4
			case "stale-version":
				wantCode, wantExit = "run-version-conflict", 4
			}
			if wantCode != "" {
				assertP22ProcessFailure(t, "run.close", wantCode, wantExit, response)
				if events != p21JournalEventCount(t, a.workspace) || !reflect.DeepEqual(before, p22FaultProjection(t, a.workspace)) {
					t.Fatal("rejected close changed facts")
				}
				return
			}
			if response.ExitCode != 0 || response.Stderr != "" {
				t.Fatalf("close: %#v", response)
			}
			var result map[string]any
			p22FaultOK(t, json.Unmarshal([]byte(response.Stdout), &result))
			status := "failed"
			if kind == "qualifying" {
				status = "succeeded"
			}
			if p22ResultPath(result, "result.run.status") != status || p22ResultPath(result, "result.run.version") != float64(a.runVersion+1) || p22ResultPath(result, "result.run.ended_at") == nil {
				t.Fatalf("close did not derive %s: %#v", status, result)
			}
			after := p22FaultProjection(t, a.workspace)
			for table, rows := range before {
				switch table {
				case "runs", "run_runtime", "idempotency_results", "applied_events", "projection_meta":
					continue
				}
				if !reflect.DeepEqual(rows, after[table]) {
					t.Fatalf("close changed %s", table)
				}
			}
			if p21JournalEventCount(t, a.workspace) != events+1 {
				t.Fatal("close appended wrong Event count")
			}
			replay := a.p22RunInput(t, "run", "close", a.actor, input)
			var again map[string]any
			p22FaultOK(t, json.Unmarshal([]byte(replay.Stdout), &again))
			if replay.ExitCode != 0 || replay.Stderr != "" || p22ResultPath(again, "result.replayed") != true || !reflect.DeepEqual(p22ResultPath(again, "result.run"), p22ResultPath(result, "result.run")) || p22ResultPath(again, "event.id") != p22ResultPath(result, "event.id") {
				t.Fatal("close replay changed terminal proof")
			}
			if !reflect.DeepEqual(after, p22FaultProjection(t, a.workspace)) || p21JournalEventCount(t, a.workspace) != events+1 {
				t.Fatal("close replay changed facts")
			}
		})
	}
}

func p22CloseProofFixture(t *testing.T, kind string) *p22FIDRecoveryAdapter {
	t.Helper()
	argv := []string{"/bin/true"}
	if kind == "failed" {
		argv = []string{"/bin/false"}
	}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, p22FinalActions, p22ExecutionFixture{Argv: argv})
	p22AcceptedCall(t, a, "close-baseline", "run.observe-budget", map[string]any{"value": 0})
	if kind == "missing" {
		return a
	}
	if kind == "outstanding" {
		conn := p22FaultDB(t, a.workspace)
		_, err := conn.Exec(`CREATE TRIGGER p22_close_no_receipt BEFORE INSERT ON execution_receipts BEGIN SELECT RAISE(ABORT,'close outstanding Attempt fixture'); END`)
		p22FaultOK(t, err)
		p22FaultOK(t, conn.Close())
	}
	gate, err := evidence.RunGate(context.Background(), p22BudgetGateOptions(t, a, "close-proof-gate"))
	if kind == "outstanding" {
		if err == nil || !strings.Contains(err.Error(), "close outstanding Attempt fixture") || gate.Receipt != nil {
			t.Fatalf("outstanding fixture: %#v %v", gate, err)
		}
		conn := p22FaultDB(t, a.workspace)
		_, err = conn.Exec(`DROP TRIGGER p22_close_no_receipt`)
		p22FaultOK(t, err)
		p22FaultOK(t, conn.Close())
	} else {
		p22FaultOK(t, err)
		if gate.Receipt == nil {
			t.Fatal("proof fixture has no Receipt")
		}
		expected := "succeeded"
		if kind == "failed" {
			expected = "failed"
		}
		if gate.Receipt.Outcome != expected {
			t.Fatal("proof fixture did not observe expected process outcome")
		}
		p21Run(t, a.binary, a.workspace, "--json", "evidence", "link", "--actor", a.actor, "--input", p21Input(t, a.workspace, "close-proof-link", map[string]any{"obligation_id": "OBL-p22-corpus", "receipt_id": gate.Receipt.ID, "target_repository": a.repository, "required_artifact_ids": []string{}}), "--idempotency-key", "close-proof-link")
	}
	current := p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
	a.runVersion = p22Number(p21Path(t, current, "result", "run", "version"))
	return a
}
