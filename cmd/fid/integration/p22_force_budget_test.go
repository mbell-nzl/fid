//go:build integration

package integration

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func p22ForceBudgetFixture(t *testing.T, maximum int) *p22FIDRecoveryAdapter {
	t.Helper()
	return p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, p22FinalActions,
		p22ExecutionFixture{Argv: []string{"/bin/true"}, ExtraBudgets: []any{map[string]any{"name": "interventions", "unit": "interventions", "maximum": maximum, "consumed": nil, "remaining": nil, "scope": "envelope", "warning_threshold": nil, "exhaustion_action": "budget-exhausted"}}})
}

func p22ForceBudgetInput(t *testing.T, a *p22FIDRecoveryAdapter, key string) map[string]any {
	t.Helper()
	p22AcceptedCall(t, a, key+"-observe", "run.observe", map[string]any{"liveness": map[string]any{"executor_id": a.executorID, "state": "live", "source_ref": "process:force-budget", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
	p22AcceptedCall(t, a, key+"-inspect", "run.inspect", nil)
	run := p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
	return map[string]any{"authority": map[string]any{"mode": "force", "actor_id": a.actor, "reason": "authorized takeover", "reference": "HUM-force-budget", "grant_ref": map[string]any{"envelope_id": p21Path(t, run, "result", "run", "envelope_id"), "revision": 1, "ordinal": 0}}, "executor": map[string]any{"id": "executor-" + key, "source": "activated-harness"}, "ttl_seconds": 3600}
}

func p22AssertInterventions(t *testing.T, response map[string]any, consumed int) {
	t.Helper()
	for _, raw := range p21Path(t, response, "result", "run", "budgets").([]any) {
		budget := raw.(map[string]any)
		if budget["name"] == "interventions" {
			if budget["consumed"] != float64(consumed) {
				t.Fatalf("interventions=%#v, want %d", budget, consumed)
			}
			return
		}
	}
	t.Fatal("missing interventions budget")
}

func TestP22ForceRecoveryBudgetChainBoundaryReplayAndPortability(t *testing.T) {
	t.Parallel()
	a := p22ForceBudgetFixture(t, 2)
	var first map[string]any
	for i := 1; i <= 2; i++ {
		key := fmt.Sprintf("force-%d", i)
		input := p22ForceBudgetInput(t, a, key)
		response := p22AcceptedCall(t, a, key, "run.recover", input)
		a.executorID = "executor-" + key
		p22AssertInterventions(t, response, i)
		if i == 1 {
			first = response
		}
	}
	replay := p22AcceptedCall(t, a, "force-1", "run.recover", nil)
	p22AssertInterventions(t, replay, 1)
	if !reflect.DeepEqual(first["event"], replay["event"]) || p22ResultPath(replay, "result.replayed") != true {
		t.Fatal("retry lost original Event/snapshot")
	}
	conn := p22FaultDB(t, a.workspace)
	var samples, stops, links int
	p22FaultOK(t, conn.QueryRow(`SELECT count(*) FROM budget_samples WHERE budget_name='interventions'`).Scan(&samples))
	p22FaultOK(t, conn.QueryRow(`SELECT count(*) FROM stops WHERE code='budget-exhausted' AND subject_id=?`, a.runID).Scan(&stops))
	p22FaultOK(t, conn.QueryRow(`SELECT count(*) FROM attention_stop_links`).Scan(&links))
	if samples != 2 || stops != 1 || links != 1 {
		t.Fatalf("samples/stops/links=%d/%d/%d", samples, stops, links)
	}
	var chained int
	p22FaultOK(t, conn.QueryRow(`SELECT count(*) FROM budget_samples b JOIN budget_samples p ON b.previous_sample_id=p.id JOIN run_recoveries r ON b.source_sample_id=r.id WHERE b.budget_name='interventions' AND b.source='fid' AND b.value=1 AND b.measurement_ref_json IS NULL AND b.run_id=r.successor_run_id AND b.run_id<>p.run_id`).Scan(&chained))
	if chained != 1 {
		t.Fatal("successor did not retain immutable intervention chain")
	}
	p22FaultOK(t, conn.Close())
	input := p22ForceBudgetInput(t, a, "blocked-force")
	before := p22FaultProjection(t, a.workspace)
	events := p21JournalEventCount(t, a.workspace)
	denied := a.Execute(t, p22ConformanceCase{Name: "blocked-force", Operation: "run.recover", Input: input})
	if denied.ExitCode == 0 {
		t.Fatal("exhausted successor admitted another force recovery")
	}
	if events != p21JournalEventCount(t, a.workspace) || !reflect.DeepEqual(before, p22FaultProjection(t, a.workspace)) {
		t.Fatal("denied force recovery mutated state")
	}
	p22PortableThreeWay(t, a, nil)
	replay = p22AcceptedCall(t, a, "force-1", "run.recover", nil)
	p22AssertInterventions(t, replay, 1)
}
