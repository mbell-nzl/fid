//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestP22AbandonUsesExactHumanForceAndTypedObservationFailures(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"ordinary", "force"} {
		cases := []string{"valid", "unknown", "missing-observation", "null-observation", "nonexistent-observation", "forged-target"}
		if mode == "ordinary" {
			cases = append(cases, "live")
		} else {
			cases = append(cases, "missing-reference")
		}
		for _, kind := range cases {
			t.Run(mode+"/"+kind, func(t *testing.T) {
				actorType, state := "agent", "terminated"
				if mode == "force" {
					actorType, state = "human", "live"
				}
				if kind == "unknown" || kind == "live" {
					state = kind
				}
				a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), actorType, 3600, []string{}, p22FinalActions)
				a.p22Observe(t, state)
				authority := map[string]any{"mode": mode, "actor_id": a.actor, "reason": "intentional cessation", "reference": nil, "grant_ref": nil}
				if mode == "force" {
					run := p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
					authority["reference"] = "HUM-close-fixture"
					authority["grant_ref"] = map[string]any{"envelope_id": p21Path(t, run, "result", "run", "envelope_id"), "revision": 1, "ordinal": 0}
				}
				input := map[string]any{"run_id": a.runID, "expected_version": a.runVersion, "expected_claim_version": a.claimVersion, "observation_id": a.observation, "authority": authority, "repository": a.repository}
				code, exit := "", 0
				switch kind {
				case "unknown":
					code, exit = "executor-liveness-unknown", 4
				case "live":
					code, exit = "executor-still-live", 4
				case "missing-observation":
					delete(input, "observation_id")
					code, exit = "invalid-argument", 2
				case "null-observation":
					input["observation_id"] = nil
					code, exit = "invalid-argument", 2
				case "nonexistent-observation":
					input["observation_id"] = "OBS-absent"
					code, exit = "observation-not-found", 3
				case "missing-reference":
					authority["reference"] = nil
					code, exit = "invalid-argument", 2
				case "forged-target":
					p22Write(t, a.workspace, "base.txt", "changed terminal target\n")
					code, exit = "recovery-state-conflict", 4
				}
				before, events := p22FaultProjection(t, a.workspace), p21JournalEventCount(t, a.workspace)
				response := a.p22RunInput(t, "run", "abandon", a.actor, input)
				if code != "" {
					assertP22ProcessFailure(t, "run.abandon", code, exit, response)
					if !reflect.DeepEqual(before, p22FaultProjection(t, a.workspace)) || events != p21JournalEventCount(t, a.workspace) {
						t.Fatal("rejected abandon changed facts")
					}
					return
				}
				if response.ExitCode != 0 || response.Stderr != "" {
					t.Fatalf("abandon: %#v", response)
				}
				var result map[string]any
				p22FaultOK(t, json.Unmarshal([]byte(response.Stdout), &result))
				for path, want := range map[string]any{"result.run.status": "abandoned", "result.claim.status": "released", "result.claim.release_reason": "abandoned", "result.claim.version": float64(a.claimVersion + 1), "result.run.version": float64(a.runVersion + 1), "result.run.actor_id": a.actor} {
					if !reflect.DeepEqual(p22ResultPath(result, path), want) {
						t.Fatalf("abandon %s differs", path)
					}
				}
				after := p22FaultProjection(t, a.workspace)
				for table, rows := range before {
					switch table {
					case "runs", "run_runtime", "claims", "idempotency_results", "applied_events", "projection_meta":
						continue
					}
					if !reflect.DeepEqual(rows, after[table]) {
						t.Fatalf("abandon changed %s proof/context", table)
					}
				}
				if len(after["runs"]) != len(before["runs"]) || len(after["claims"]) != len(before["claims"]) || events+1 != p21JournalEventCount(t, a.workspace) {
					t.Fatal("abandon created successor identity or duplicate Event")
				}
				replay := a.p22RunInput(t, "run", "abandon", a.actor, input)
				var again map[string]any
				p22FaultOK(t, json.Unmarshal([]byte(replay.Stdout), &again))
				if replay.ExitCode != 0 || replay.Stderr != "" || p22ResultPath(again, "result.replayed") != true {
					t.Fatal("abandon replay failed")
				}
				again["result"].(map[string]any)["replayed"] = false
				if !reflect.DeepEqual(result, again) || !reflect.DeepEqual(after, p22FaultProjection(t, a.workspace)) || events+1 != p21JournalEventCount(t, a.workspace) {
					t.Fatal("abandon replay changed facts/provenance")
				}
			})
		}
	}
}
