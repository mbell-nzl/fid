package workflow

import (
	"reflect"
	"testing"
)

func TestForceRecoveryChargesInterventionBudgetsExactlyOnce(t *testing.T) {
	for _, mode := range []string{"ordinary", "force"} {
		t.Run(mode, func(t *testing.T) {
			options, request := recoveryFixture(t)
			if mode == "force" {
				request.Authority = exactRecoveryForce(t, options)
			}
			seedRecoveryInterventionBudget(t, options.WorkspaceRoot, "interventions", 1)
			first, err := RecoverRun(options, request)
			mustRecoveryTest(t, err)
			want := int64(0)
			if mode == "force" {
				want = 1
			}
			if len(first.Run.Budgets) != 1 || first.Run.Budgets[0].Consumed == nil || *first.Run.Budgets[0].Consumed != want {
				t.Fatalf("%s recovery budgets = %#v, want consumed %d", mode, first.Run.Budgets, want)
			}
			replayed, err := RecoverRun(options, request)
			mustRecoveryTest(t, err)
			if !replayed.Replayed || !reflect.DeepEqual(first.Run, replayed.Run) || first.EventID != replayed.EventID {
				t.Fatal("recovery retry changed its original snapshot")
			}
			conn, err := openWorkflowDB(options.WorkspaceRoot)
			mustRecoveryTest(t, err)
			defer closeWorkflowDB(conn)
			var samples, stops, links int64
			mustRecoveryTest(t, conn.QueryRow(`SELECT COUNT(*) FROM budget_samples`).Scan(&samples))
			mustRecoveryTest(t, conn.QueryRow(`SELECT COUNT(*) FROM stops WHERE subject_id=? AND code='budget-exhausted'`, first.Run.ID).Scan(&stops))
			mustRecoveryTest(t, conn.QueryRow(`SELECT COUNT(*) FROM attention_stop_links`).Scan(&links))
			if samples != want || stops != want || links != want {
				t.Fatalf("samples/Stops/links = %d/%d/%d, want %d each", samples, stops, links, want)
			}
			if mode == "force" {
				var exact int
				mustRecoveryTest(t, conn.QueryRow(`SELECT COUNT(*) FROM budget_samples WHERE source='fid' AND source_sample_id=? AND run_id=? AND actor_id=? AND value=1 AND sample_kind='delta' AND previous_sample_id IS NULL AND measurement_ref_json IS NULL AND observed_at=? AND created_event_id=?`, first.Recovery.ID, first.Run.ID, options.ActorID, first.Run.StartedAt, first.EventID).Scan(&exact))
				if exact != 1 {
					t.Fatal("force sample lost Recovery/Run/Actor/time/Event binding")
				}
			}
		})
	}
}

func seedRecoveryInterventionBudget(t *testing.T, workspace, name string, maximum int64) {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	mustRecoveryTest(t, err)
	defer closeWorkflowDB(conn)
	budget := Budget{Name: name, Scope: "envelope", Unit: "interventions", Maximum: maximum, ExhaustionAction: "budget-exhausted"}
	raw := string(mustRunJSON(budgetPolicyValue(budget)))
	_, err = conn.Exec(`INSERT INTO envelope_budgets (envelope_id,revision,name,budget_json) VALUES ('ENV-recovery',1,?,?)`, name, raw)
	mustRecoveryTest(t, err)
	_, err = conn.Exec(`INSERT INTO run_budget_bindings (run_id,name,budget_json) VALUES ('RUN-prior',?,?)`, name, raw)
	mustRecoveryTest(t, err)
}

func TestForceRecoveryChargesEveryConfiguredInterventionBudget(t *testing.T) {
	options, request := recoveryFixture(t)
	request.Authority = exactRecoveryForce(t, options)
	seedRecoveryInterventionBudget(t, options.WorkspaceRoot, "interventions-a", 1)
	seedRecoveryInterventionBudget(t, options.WorkspaceRoot, "interventions-b", 1)
	result, err := RecoverRun(options, request)
	mustRecoveryTest(t, err)
	for _, budget := range result.Run.Budgets {
		if budget.Consumed == nil || *budget.Consumed != 1 || budget.State != "exhausted" {
			t.Fatalf("uncharged configured budget: %#v", budget)
		}
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	mustRecoveryTest(t, err)
	defer closeWorkflowDB(conn)
	var samples, stops, attention, links int
	mustRecoveryTest(t, conn.QueryRow(`SELECT COUNT(*) FROM budget_samples`).Scan(&samples))
	mustRecoveryTest(t, conn.QueryRow(`SELECT COUNT(*) FROM stops`).Scan(&stops))
	mustRecoveryTest(t, conn.QueryRow(`SELECT COUNT(*) FROM attention_items`).Scan(&attention))
	mustRecoveryTest(t, conn.QueryRow(`SELECT COUNT(*) FROM attention_stop_links`).Scan(&links))
	if samples != 2 || stops != 2 || attention != 1 || links != 2 {
		t.Fatalf("configured budget transaction = %d/%d/%d/%d, want 2/2/1/2", samples, stops, attention, links)
	}
}
