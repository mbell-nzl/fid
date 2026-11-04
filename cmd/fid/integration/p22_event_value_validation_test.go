//go:build integration

package integration

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type p22EventCorruption struct {
	name, family string
	mutate       func(*journal.Event)
}

// These cases preserve the independent final-review-2 counterexamples. Every
// case gets a genuine producer prefix and a separate projection for each mode.
func TestP22CurrentEventClosedValuesAndConditionalEffects(t *testing.T) {
	t.Parallel()
	a := p22CompleteOneIssueRecoveryJourney(t)
	binary := a.binary
	activation, id := p22ActivationBoundaryFixture(t, binary, "committed")
	state := p21ActivationStatus(t, activation, id)
	p21Run(t, binary, activation, "--json", "activation", "recover", "--input",
		p21Input(t, activation, "closed-select", map[string]any{"activation_id": id, "expected_version": state.Version, "expected_status": state.Status, "direction": "complete-forward", "reason": "closed Event control"}),
		"--idempotency-key", "closed-select")
	sources := map[string]string{"run.observed": a.workspace, "run.recovered": a.workspace, "activation.recovery-selected": activation}
	for family, source := range sources {
		t.Run("control/"+family, func(t *testing.T) {
			prefix, event := p22CurrentEventPrefix(t, source, family)
			for _, mode := range []string{"live", "replay", "rebuild"} {
				t.Run(mode, func(t *testing.T) { p22CheckCurrentEventCase(t, binary, source, prefix, event, mode, true) })
			}
		})
	}
	// The review's Claim-removal case was a no-op on an already expired Claim.
	prefix, recovered := p22CurrentEventPrefix(t, a.workspace, "run.recovered")
	for _, effect := range recovered.Effects {
		if effect.Table == "claims" && effect.Op == "update" {
			t.Fatal("expired-Claim control unexpectedly disposes a Claim")
		}
	}
	t.Run("control/expired-claim-no-op", func(t *testing.T) {
		p22CheckCurrentEventCase(t, binary, a.workspace, prefix, recovered, "rebuild", true)
	})
	cases := p22ClosedEventCorruptions()
	if len(cases) != 23 {
		t.Fatal("the 23 independent counterexamples must stay explicit")
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			source := sources[test.family]
			prefix, event := p22CurrentEventPrefix(t, source, test.family)
			before, err := json.Marshal(event)
			p22FaultOK(t, err)
			test.mutate(&event)
			after, err := json.Marshal(event)
			p22FaultOK(t, err)
			if bytes.Equal(before, after) {
				t.Fatal("corruption was a no-op")
			}
			for _, mode := range []string{"live", "replay", "rebuild"} {
				t.Run(mode, func(t *testing.T) { p22CheckCurrentEventCase(t, binary, source, prefix, event, mode, false) })
			}
		})
	}
}

func p22ClosedEventCorruptions() []p22EventCorruption {
	obs := func(name string, mutate func(map[string]any)) p22EventCorruption {
		return p22EventCorruption{name, "run.observed", p22MutateEventJSON("run_observations", "observation_json", mutate)}
	}
	recovery := func(name string, column string, value any) p22EventCorruption {
		return p22EventCorruption{name, "run.recovered", p22SetEventValue("run_recoveries", column, value)}
	}
	stop := func(name, column string, value any) p22EventCorruption {
		return p22EventCorruption{name, "activation.recovery-selected", p22SetEventValue("stops", column, value)}
	}
	result := func(name string, mutate func(map[string]any)) p22EventCorruption {
		return p22EventCorruption{name, "run.observed", p22MutateEventJSON("idempotency_results", "result_json", mutate)}
	}
	return []p22EventCorruption{
		obs("observation-extra-liveness-field", func(v map[string]any) { v["liveness"].(map[string]any)["qa_unknown"] = true }),
		obs("observation-unknown-liveness-state", func(v map[string]any) { v["liveness"].(map[string]any)["state"] = "invented-state" }),
		obs("observation-null-liveness", func(v map[string]any) { v["liveness"] = nil }),
		obs("observation-wrong-executor", func(v map[string]any) { v["liveness"].(map[string]any)["executor_id"] = "executor-not-this-run" }),
		obs("observation-nested-actor-mismatch", func(v map[string]any) { v["actor_id"] = "actor-not-observer" }),
		{"observation-bad-time", "run.observed", p22SetEventValue("run_observations", "observed_at", "not-a-timestamp")},
		{"recovery-missing-predecessor-terminalization", "run.recovered", func(e *journal.Event) {
			kept := []journal.Effect{}
			for _, effect := range e.Effects {
				if effect.Op == "update" && (effect.Table == "runs" || effect.Table == "run_runtime") {
					continue
				}
				kept = append(kept, effect)
			}
			e.Effects = kept
		}},
		{"recovery-unknown-authority-field", "run.recovered", p22MutateEventJSON("run_recoveries", "authority_json", func(v map[string]any) { v["qa_unknown"] = true })},
		recovery("recovery-malformed-reconciliation", "reconciliation_json", "{}"),
		recovery("recovery-adopted-paths-wrong-type", "adopted_paths_json", "\""+"not-an-array"+"\""),
		recovery("recovery-protected-adoption", "adopted_paths_json", "[\".fid/fid.sqlite\"]"),
		{"stop-unknown-detail-field", "activation.recovery-selected", p22MutateEventJSON("stops", "details_json", func(v map[string]any) { v["qa_unknown"] = true })},
		{"stop-invalid-condition-type", "activation.recovery-selected", p22MutateEventJSON("stops", "resolution_condition_json", func(v map[string]any) { v["type"] = "invented-condition" })},
		stop("stop-invalid-blocked-action", "blocked_actions_json", "[\"invented.action\"]"),
		stop("stop-missing-condition", "resolution_condition_json", "{}"),
		{"activation-impossible-recovery-direction", "activation.recovery-selected", p22SetEventValue("activation_recoveries", "direction", "restore-prior")},
		result("result-payload-identity", func(v map[string]any) { v["event_id"] = "EVT-not-this-event" }),
		{"observation-duplicate-json-key", "run.observed", func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table == "run_observations" {
					raw := e.Effects[i].Values["observation_json"].(string)
					e.Effects[i].Values["observation_json"] = "{\"actor_id\":\"duplicate-wrong-actor\"," + raw[1:]
				}
			}
		}},
		obs("observation-missing-source-ref", func(v map[string]any) { delete(v["liveness"].(map[string]any), "source_ref") }),
		obs("observation-inconsistent-time", func(v map[string]any) { v["observed_at"] = "2026-01-01T00:00:00Z" }),
		result("result-wrong-observation-id", func(v map[string]any) { v["observation"].(map[string]any)["id"] = "OBS-aaaaaa" }),
		result("result-wrong-run-id", func(v map[string]any) { v["run"].(map[string]any)["id"] = "RUN-aaaaaa" }),
		{"recovery-missing-budget-bindings", "run.recovered", func(e *journal.Event) {
			kept := []journal.Effect{}
			for _, effect := range e.Effects {
				if effect.Table != "run_budget_bindings" {
					kept = append(kept, effect)
				}
			}
			e.Effects = kept
		}},
	}
}

func p22MutateEventJSON(table, column string, mutate func(map[string]any)) func(*journal.Event) {
	return func(event *journal.Event) {
		for i := range event.Effects {
			effect := &event.Effects[i]
			if effect.Table != table {
				continue
			}
			var value map[string]any
			if err := json.Unmarshal([]byte(effect.Values[column].(string)), &value); err != nil {
				panic(err)
			}
			mutate(value)
			raw, err := json.Marshal(value)
			if err != nil {
				panic(err)
			}
			effect.Values[column] = string(raw)
			return
		}
	}
}
func p22SetEventValue(table, column string, value any) func(*journal.Event) {
	return func(event *journal.Event) {
		for i := range event.Effects {
			if event.Effects[i].Table == table {
				event.Effects[i].Values[column] = value
				return
			}
		}
	}
}
func p22CurrentEventPrefix(t *testing.T, source, family string) ([]byte, journal.Event) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(config.LedgerDir(source), "journal"))
	p22FaultOK(t, err)
	var prefix []byte
	for _, entry := range entries {
		raw, err := os.ReadFile(filepath.Join(config.LedgerDir(source), "journal", entry.Name()))
		p22FaultOK(t, err)
		for _, line := range bytes.Split(bytes.TrimSpace(raw), []byte{'\n'}) {
			var event journal.Event
			p22FaultOK(t, json.Unmarshal(line, &event))
			if event.EventType == family {
				return prefix, event
			}
			prefix = append(append(prefix, line...), '\n')
		}
	}
	t.Fatalf("missing genuine %s", family)
	return nil, journal.Event{}
}

func p22CheckCurrentEventCase(t *testing.T, binary, source string, prefix []byte, event journal.Event, mode string, accepted bool) {
	t.Helper()
	if mode == "rebuild" {
		p22CheckCurrentRebuild(t, binary, source, prefix, event, accepted)
		return
	}
	workspace := t.TempDir()
	p22FaultOK(t, os.MkdirAll(config.LedgerDir(workspace), 0700))
	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
	p22FaultOK(t, fiddb.MigrateUp(dbPath, filepath.Join(integrationRepositoryRoot(), "migrations")))
	conn, err := fiddb.Open(dbPath)
	p22FaultOK(t, err)
	defer conn.Close()
	dir := t.TempDir()
	path := filepath.Join(dir, "2026-09-09.jsonl")
	p22FaultOK(t, os.WriteFile(path, prefix, 0600))
	_, err = journal.ReplayJournal(conn, dir, nil)
	p22FaultOK(t, err)
	before := p22FaultProjection(t, workspace)
	if mode == "live" {
		hash, hashErr := journal.HashEvent(event)
		p22FaultOK(t, hashErr)
		err = journal.ProjectEvent(conn, journal.ProjectionInput{Event: event, EventHash: hash})
	} else {
		raw, marshalErr := json.Marshal(event)
		p22FaultOK(t, marshalErr)
		body := append(append(append([]byte{}, prefix...), raw...), '\n')
		p22FaultOK(t, os.WriteFile(path, body, 0600))
		_, err = journal.ReplayJournal(conn, dir, &journal.ReplayOptions{StartFile: filepath.Base(path), StartOffset: int64(len(prefix))})
	}
	if accepted {
		p22FaultOK(t, err)
		return
	}
	if err == nil {
		t.Fatal("malformed current Event accepted")
	}
	if !reflect.DeepEqual(before, p22FaultProjection(t, workspace)) {
		t.Fatal("rejected Event changed projection")
	}
}
func p22CheckCurrentRebuild(t *testing.T, binary, source string, prefix []byte, event journal.Event, accepted bool) {
	t.Helper()
	workspace := t.TempDir()
	copyParityFixture(t, source, workspace)
	dir := filepath.Join(config.LedgerDir(workspace), "journal")
	entries, err := os.ReadDir(dir)
	p22FaultOK(t, err)
	for _, entry := range entries {
		p22FaultOK(t, os.Remove(filepath.Join(dir, entry.Name())))
	}
	raw, err := json.Marshal(event)
	p22FaultOK(t, err)
	body := append(append(append([]byte{}, prefix...), raw...), '\n')
	p22FaultOK(t, os.WriteFile(filepath.Join(dir, "2026-09-09.jsonl"), body, 0600))
	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
	before, err := os.ReadFile(dbPath)
	p22FaultOK(t, err)
	command := exec.Command(binary, "--project", workspace, "--json", "rebuild")
	command.Dir = workspace
	output, err := command.CombinedOutput()
	if accepted {
		if err != nil {
			t.Fatalf("valid rebuild: %v %s", err, output)
		}
		return
	}
	if err == nil {
		t.Fatalf("malformed rebuild accepted: %s", output)
	}
	after, readErr := os.ReadFile(dbPath)
	p22FaultOK(t, readErr)
	if sha256.Sum256(before) != sha256.Sum256(after) {
		t.Fatal("rejected rebuild replaced projection")
	}
}

func TestP22CurrentProducerEmitsRequiredClosedFields(t *testing.T) {
	t.Parallel()
	a := p22CompleteOneIssueRecoveryJourney(t)
	b, _ := p22PortableSource(t, "deduplicated-attention")
	p22CheckProducerRequiredFields(t, []string{a.workspace, b.workspace})
}

func p22CheckProducerRequiredFields(t *testing.T, workspaces []string) {
	checkedCondition, checkedClaim, checkedStop, checkedAttention := false, false, false, false
	for _, workspace := range workspaces {
		files, err := os.ReadDir(filepath.Join(config.LedgerDir(workspace), "journal"))
		p22FaultOK(t, err)
		require := func(label string, value map[string]any, keys ...string) {
			for _, key := range keys {
				if _, ok := value[key]; !ok {
					t.Errorf("%s omits required field %s", label, key)
				}
			}
		}
		for _, file := range files {
			raw, err := os.ReadFile(filepath.Join(config.LedgerDir(workspace), "journal", file.Name()))
			p22FaultOK(t, err)
			for _, line := range bytes.Split(bytes.TrimSpace(raw), []byte{'\n'}) {
				var event journal.Event
				p22FaultOK(t, json.Unmarshal(line, &event))
				for _, effect := range event.Effects {
					if effect.Table == "stops" && effect.Op == "insert" {
						var condition map[string]any
						p22FaultOK(t, json.Unmarshal([]byte(effect.Values["resolution_condition_json"].(string)), &condition))
						require("stored Condition", condition, "budget_name", "capability_id")
						checkedCondition = true
					}
					if effect.Table != "idempotency_results" {
						continue
					}
					var result map[string]any
					p22FaultOK(t, json.Unmarshal([]byte(effect.Values["result_json"].(string)), &result))
					if event.EventType == "run.recovered" {
						claim := result["claim"].(map[string]any)
						require("Recovery Claim", claim, "scope", "mode")
						checkedClaim = true
					}
					if stop, ok := result["stop"].(map[string]any); ok {
						require("Stop", stop, "resolution")
						require("Stop Condition", stop["resolution_condition"].(map[string]any), "budget_name", "capability_id")
						checkedStop = true
					}
					if attention, ok := result["attention"].(map[string]any); ok {
						require("Attention", attention, "resolution_event_id")
						checkedAttention = true
					}
				}
			}
		}
	}
	if !checkedCondition || !checkedClaim || !checkedStop || !checkedAttention {
		t.Fatal("producer fixtures did not exercise every required shape")
	}
}
