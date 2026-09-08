//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestP22RecoveryUsesClaimScopeAndExactHumanOutsideAdoption(t *testing.T) {
	for _, mode := range []string{"ordinary-forged-scope", "force-outside", "force-protected"} {
		t.Run(mode, func(t *testing.T) {
			a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, p22FinalActions)
			state := "live"
			if mode == "ordinary-forged-scope" {
				state = "terminated"
			}
			a.p22Observe(t, state)
			if mode == "ordinary-forged-scope" {
				a.p22Expire(t)
			}
			name := "outside.txt"
			if mode == "force-protected" {
				name = ".fid/protected-adoption"
			}
			p22Write(t, a.workspace, name, "precise outside work")
			inspection := a.p22Run(t, "run", "inspect", a.runID)
			if inspection.ExitCode != 0 {
				t.Fatalf("inspect: %#v", inspection)
			}
			a.p22Advance(t, "run.inspect", inspection.Stdout)
			run := p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
			input := map[string]any{"adopted_paths": []string{name}}
			if mode == "ordinary-forged-scope" {
				var baseline workflow.Repository
				encoded, err := json.Marshal(a.repository)
				p22FaultOK(t, err)
				p22FaultOK(t, json.Unmarshal(encoded, &baseline))
				current, err := project.ObserveRepository(a.workspace)
				p22FaultOK(t, err)
				forged, err := workflow.ReconcileRepository(baseline, current, nil, []string{name})
				p22FaultOK(t, err)
				input["reconciliation_digest"] = forged.Digest
			} else {
				input["authority"] = map[string]any{"mode": "force", "actor_id": a.actor, "reason": "adopt exact identified outside work", "reference": "HUM-explicit-outside", "grant_ref": map[string]any{"envelope_id": p21PathString(t, run, "result", "run", "envelope_id"), "revision": 1, "ordinal": 0}}
			}
			before := p22FaultProjection(t, a.workspace)
			events := p21JournalEventCount(t, a.workspace)
			response := a.Execute(t, p22ConformanceCase{Name: "exact-scope", Operation: "run.recover", Input: input})
			if mode != "force-outside" {
				p22AssertWireFailure(t, "run.recover", "recovery-state-conflict", response)
				if p21JournalEventCount(t, a.workspace) != events || !reflect.DeepEqual(before, p22FaultProjection(t, a.workspace)) {
					t.Fatal("rejected scope changed facts")
				}
				return
			}
			if response.ExitCode != 0 {
				t.Fatalf("exact Human outside adoption: %#v", response)
			}
			var result map[string]any
			p22FaultOK(t, json.Unmarshal([]byte(response.Stdout), &result))
			if p22ResultPath(result, "result.recovery.reconciliation.paths.0.classification") != "outside-scope" || p22ResultPath(result, "result.recovery.authority.reference") != "HUM-explicit-outside" {
				t.Fatal("force adoption lost classification or Human provenance")
			}
			conn := p22FaultDB(t, a.workspace)
			var scope string
			p22FaultOK(t, conn.QueryRow(`SELECT scope_json FROM claims WHERE id = ?`, a.claimID).Scan(&scope))
			p22FaultOK(t, conn.Close())
			if scope != `{"paths":["base.txt"]}` || p21JournalEventCount(t, a.workspace) != events+1 {
				t.Fatal("force adoption widened immutable Claim or duplicated recovery")
			}
		})
	}
}
