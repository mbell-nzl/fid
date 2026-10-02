package workflow

import (
	"errors"
	"testing"
	"time"
)

func TestObserveRunBudgetAppendsExternalSampleAndBoundaryEffectsExactlyOnce(t *testing.T) {
	t.Parallel()
	options, recovery := recoveryFixture(t)
	observeBudgetAuthorize(t, options.WorkspaceRoot, "input-tokens", 5)
	value := int64(5)
	request := RunBudgetObserveRequest{
		RunID: recovery.PriorRunID, ExpectedVersion: recovery.ExpectedRunVersion,
		BudgetName: "external", SourceSampleID: "harness-sample-1", SampleKind: "baseline",
		Value: &value, ObservedAt: time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano),
	}

	first, err := ObserveRunBudget(options, request)
	if err != nil {
		t.Fatalf("ObserveRunBudget() error = %v", err)
	}
	if first.Replayed || first.Sample.ID == "" || first.Sample.EventID != first.EventID || first.Budget.Remaining == nil || *first.Budget.Remaining != 0 || len(first.Stops) != 1 || len(first.Attention) != 1 {
		t.Fatalf("ObserveRunBudget() = %#v, want durable exhausted sample with one Stop and Attention", first)
	}
	replayed, err := ObserveRunBudget(options, request)
	if err != nil || !replayed.Replayed || replayed.Sample.ID != first.Sample.ID || replayed.EventID != first.EventID || len(replayed.Stops) != 1 || replayed.Stops[0].ID != first.Stops[0].ID {
		t.Fatalf("ObserveRunBudget(replay) = (%#v, %v), want original sample and Stop", replayed, err)
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	var samples, stops int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM budget_samples`).Scan(&samples); err != nil {
		t.Fatalf("count budget samples: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM stops WHERE subject_id = ?`, recovery.PriorRunID).Scan(&stops); err != nil {
		t.Fatalf("count budget Stops: %v", err)
	}
	if samples != 1 || stops != 1 {
		t.Fatalf("sample/Stop counts = %d/%d, want 1/1 after replay", samples, stops)
	}
}

func TestObserveRunBudgetRejectsFidOwnedUnitAndConflictWithoutEffects(t *testing.T) {
	t.Parallel()
	options, recovery := recoveryFixture(t)
	observeBudgetAuthorize(t, options.WorkspaceRoot, "attempts", 2)
	value := int64(1)
	request := RunBudgetObserveRequest{RunID: recovery.PriorRunID, ExpectedVersion: recovery.ExpectedRunVersion, BudgetName: "external", SourceSampleID: "harness-sample-1", SampleKind: "baseline", Value: &value, ObservedAt: time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}
	before := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
	if _, err := ObserveRunBudget(options, request); err == nil {
		t.Fatal("ObserveRunBudget(fid-owned) error = nil, want rejection")
	}
	if after := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot)); string(after) != string(before) {
		t.Fatal("fid-owned budget rejection appended effects")
	}

	observeBudgetAuthorize(t, options.WorkspaceRoot, "input-tokens", 10)
	if _, err := ObserveRunBudget(options, request); err != nil {
		t.Fatalf("ObserveRunBudget(first) error = %v", err)
	}
	changed := request
	changed.ExpectedVersion++
	changed.Value = int64Pointer(2)
	changed.PreviousSampleID = nil
	changed.SourceSampleID = request.SourceSampleID
	changedKey := options
	changedKey.IdempotencyKey = "observe-budget-conflict"
	before = readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
	if _, err := ObserveRunBudget(changedKey, changed); err == nil {
		t.Fatal("ObserveRunBudget(conflicting sample) error = nil, want rejection")
	} else {
		var conflict *BudgetSampleConflictError
		if !errors.As(err, &conflict) {
			t.Fatalf("ObserveRunBudget(conflicting sample) error = %v, want BudgetSampleConflictError", err)
		}
	}
	if after := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot)); string(after) != string(before) {
		t.Fatal("conflicting sample rejection appended effects")
	}
}

func TestObserveRunBudgetReplayKeepsHistoricalDerivedState(t *testing.T) {
	t.Parallel()
	options, recovery := recoveryFixture(t)
	observeBudgetAuthorize(t, options.WorkspaceRoot, "input-tokens", 10)
	baseline := int64(2)
	firstRequest := RunBudgetObserveRequest{RunID: recovery.PriorRunID, ExpectedVersion: recovery.ExpectedRunVersion, BudgetName: "external", SourceSampleID: "historical-baseline", SampleKind: "baseline", Value: &baseline, ObservedAt: time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}
	first, err := ObserveRunBudget(options, firstRequest)
	if err != nil || first.Budget.State != "available" || first.Budget.Remaining == nil || *first.Budget.Remaining != 8 {
		t.Fatalf("ObserveRunBudget(first) = (%#v, %v), want available historical state", first, err)
	}
	charged := int64(8)
	secondOptions := options
	secondOptions.IdempotencyKey = "observe-budget-later"
	secondRequest := RunBudgetObserveRequest{RunID: recovery.PriorRunID, ExpectedVersion: recovery.ExpectedRunVersion + 1, BudgetName: "external", SourceSampleID: "later-charge", SampleKind: "delta", PreviousSampleID: &first.Sample.ID, Value: &charged, ObservedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	second, err := ObserveRunBudget(secondOptions, secondRequest)
	if err != nil || second.Budget.State != "exhausted" {
		t.Fatalf("ObserveRunBudget(second) = (%#v, %v), want later exhausted state", second, err)
	}
	replayed, err := ObserveRunBudget(options, firstRequest)
	if err != nil || !replayed.Replayed || replayed.Budget.State != "available" || replayed.Budget.Remaining == nil || *replayed.Budget.Remaining != 8 {
		t.Fatalf("ObserveRunBudget(replay) = (%#v, %v), want original available snapshot", replayed, err)
	}
}

func observeBudgetAuthorize(t *testing.T, workspace, unit string, maximum int64) {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	budgetJSON := string(mustRunJSON(Budget{Name: "external", Scope: "envelope", Unit: unit, Maximum: maximum, ExhaustionAction: "budget-exhausted"}))
	statements := []struct {
		query string
		args  []any
	}{
		{`UPDATE activation_actor_profiles SET producer = 'activated-harness', capabilities_json = '["run.observe-budget"]' WHERE actor_id = 'agent:recovery'`, nil},
		{`UPDATE envelope_authority_grants SET actor_selector_json = '{"actor_id":"agent:recovery","required_capabilities":["run.observe-budget"]}', actions_json = '["run.observe-budget"]' WHERE envelope_id = 'ENV-recovery' AND ordinal = 0`, nil},
		{`DELETE FROM envelope_budgets WHERE envelope_id = 'ENV-recovery' AND revision = 1`, nil},
		{`INSERT INTO envelope_budgets (envelope_id, revision, name, budget_json) VALUES ('ENV-recovery', 1, 'external', ?)`, []any{budgetJSON}},
		{`DELETE FROM run_budget_bindings WHERE run_id = 'RUN-prior'`, nil},
		{`INSERT INTO run_budget_bindings (run_id, name, budget_json) VALUES ('RUN-prior', 'external', ?)`, []any{budgetJSON}},
	}
	for index, statement := range statements {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("authorize budget observation statement %d: %v", index, err)
		}
	}
}
