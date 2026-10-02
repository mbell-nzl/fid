package workflow

import (
	"errors"
	"math"
	"testing"
)

func TestFinalAllowedAttemptChargesAndRaisesStopAtomically(t *testing.T) {
	t.Parallel()
	maximum := int64(1)
	budget := Budget{Name: "attempt-budget", Scope: "envelope", Unit: "attempts", Maximum: maximum, ExhaustionAction: "budget-exhausted"}

	admission, err := AdmitBudgetSample(budget, nil, BudgetSample{
		Source:         "fid",
		SourceSampleID: "ATT-final",
		SampleKind:     "delta",
		Value:          int64Pointer(1),
	})
	if err != nil {
		t.Fatalf("AdmitBudgetSample() error = %v", err)
	}
	if !admission.Admitted || admission.Budget.Consumed == nil || *admission.Budget.Consumed != 1 || admission.Budget.Remaining == nil || *admission.Budget.Remaining != 0 || !admission.Exhausted {
		t.Fatalf("admission = %#v, want admitted exhausted final attempt with zero remaining", admission)
	}
}

func int64Pointer(value int64) *int64 { return &value }

func TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines(t *testing.T) {
	t.Parallel()
	budget := Budget{Name: "external", Scope: "envelope", Unit: "input-tokens", Maximum: 100, ExhaustionAction: "budget-exhausted"}
	unknown := BudgetSample{Source: "activated-harness", SourceSampleID: "sample-unknown", SampleKind: "delta", Value: nil}
	first, err := AdmitBudgetSample(budget, nil, unknown)
	if err != nil || first.Budget.Consumed != nil || first.Budget.Remaining != nil {
		t.Fatalf("unknown admission = %#v, %v; want null totals", first, err)
	}
	unknown.ID = unknown.SourceSampleID
	baselineID := "sample-unknown"
	baseline := BudgetSample{Source: "activated-harness", SourceSampleID: "sample-baseline", SampleKind: "baseline", PreviousSampleID: &baselineID, Value: int64Pointer(40)}
	second, err := AdmitBudgetSample(budget, []BudgetSample{unknown}, baseline)
	if err != nil || second.Budget.Consumed == nil || *second.Budget.Consumed != 40 || second.Budget.Remaining == nil || *second.Budget.Remaining != 60 {
		t.Fatalf("baseline admission = %#v, %v; want 40 consumed and 60 remaining", second, err)
	}
	duplicate, err := AdmitBudgetSample(budget, []BudgetSample{unknown, baseline}, baseline)
	if err != nil || duplicate.Budget.Consumed == nil || *duplicate.Budget.Consumed != 40 {
		t.Fatalf("duplicate admission = %#v, %v; want idempotent 40", duplicate, err)
	}
	different := baseline
	different.Value = int64Pointer(41)
	if _, err := AdmitBudgetSample(budget, []BudgetSample{unknown, baseline}, different); !errors.As(err, new(*BudgetSampleConflictError)) {
		t.Fatalf("different same-key admission error = %v, want BudgetSampleConflictError", err)
	}
	staleID := "sample-unknown"
	stale := BudgetSample{Source: "activated-harness", SourceSampleID: "sample-stale", SampleKind: "baseline", PreviousSampleID: &staleID, Value: int64Pointer(39)}
	if _, err := AdmitBudgetSample(budget, []BudgetSample{unknown, baseline}, stale); !errors.As(err, new(*BudgetSampleConflictError)) {
		t.Fatalf("stale baseline admission error = %v, want BudgetSampleConflictError", err)
	}
}

func TestBudgetSamplesRejectOverflow(t *testing.T) {
	t.Parallel()
	budget := Budget{Name: "attempt-budget", Scope: "envelope", Unit: "attempts", Maximum: math.MaxInt64, ExhaustionAction: "budget-exhausted"}
	first := BudgetSample{ID: "one", Source: "fid", SourceSampleID: "one", SampleKind: "delta", Value: int64Pointer(math.MaxInt64)}
	previous := "one"
	second := BudgetSample{Source: "fid", SourceSampleID: "two", SampleKind: "delta", PreviousSampleID: &previous, Value: int64Pointer(1)}
	if _, err := AdmitBudgetSample(budget, []BudgetSample{first}, second); err == nil {
		t.Fatal("overflow admission error = nil, want rejection")
	}
}

func TestDeriveBudgetPreservesUnknownAndExhaustedChains(t *testing.T) {
	t.Parallel()
	budget := Budget{Name: "attempts", Scope: "envelope", Unit: "attempts", Maximum: 1, ExhaustionAction: "budget-exhausted"}
	unknown, err := DeriveBudget(budget, []BudgetSample{{ID: "unknown", Source: "activated-harness", SourceSampleID: "unknown", SampleKind: "delta"}})
	if err != nil || unknown.Consumed != nil || unknown.Remaining != nil {
		t.Fatalf("DeriveBudget() unknown = %#v, %v; want null totals", unknown, err)
	}
	exhausted, err := DeriveBudget(budget, []BudgetSample{{ID: "attempt", Source: "fid", SourceSampleID: "ATT-1", SampleKind: "delta", Value: int64Pointer(1)}})
	if err != nil || exhausted.Remaining == nil || *exhausted.Remaining != 0 {
		t.Fatalf("DeriveBudget() exhausted = %#v, %v; want zero remaining", exhausted, err)
	}
}

func TestDeriveBudgetProjectsEveryAcceptedOutputState(t *testing.T) {
	t.Parallel()
	maximum := int64(10)
	warning := int64(3)
	budget := Budget{Name: "tokens", Scope: "envelope", Unit: "input-tokens", Maximum: maximum, WarningThreshold: &warning, ExhaustionAction: "budget-exhausted"}

	tests := []struct {
		name    string
		samples []BudgetSample
		want    string
	}{
		{name: "unknown", samples: []BudgetSample{{Source: "activated-harness", SourceSampleID: "unknown", SampleKind: "delta"}}, want: "unknown"},
		{name: "available", samples: []BudgetSample{{Source: "activated-harness", SourceSampleID: "available", SampleKind: "baseline", Value: int64Pointer(2)}}, want: "available"},
		{name: "warning", samples: []BudgetSample{{Source: "activated-harness", SourceSampleID: "warning", SampleKind: "baseline", Value: int64Pointer(7)}}, want: "warning"},
		{name: "exhausted", samples: []BudgetSample{{Source: "activated-harness", SourceSampleID: "exhausted", SampleKind: "baseline", Value: int64Pointer(10)}}, want: "exhausted"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			derived, err := DeriveBudget(budget, test.samples)
			if err != nil {
				t.Fatalf("DeriveBudget() error = %v", err)
			}
			if derived.State != test.want {
				t.Fatalf("DeriveBudget() state = %q, want %q", derived.State, test.want)
			}
		})
	}
}

func TestBudgetAdmissionBlockAndEffectsUseExactBudgetStopShape(t *testing.T) {
	t.Parallel()
	workspace, _ := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)

	unknown := Budget{Name: "input", Scope: "envelope", Unit: "input-tokens", Maximum: 10}
	exhausted := Budget{Name: "attempts", Scope: "envelope", Unit: "attempts", Maximum: 1, Consumed: int64Pointer(1), Remaining: int64Pointer(0)}
	blocks := BudgetAdmissionBlocks([]Budget{unknown, exhausted})
	if len(blocks) != 2 || blocks[0].Code != "budget-data-unavailable" || blocks[1].Code != "budget-exhausted" {
		t.Fatalf("BudgetAdmissionBlocks() = %#v, want unknown and exhausted blocks", blocks)
	}
	effects, stops, attention, err := BudgetStopEffects(conn, Subject{Type: "run", ID: "RUN-budget"}, "ENV-budget", "actor:budget", blocks, "2026-09-06T00:00:00Z")
	if err != nil {
		t.Fatalf("BudgetStopEffects() error = %v", err)
	}
	if len(stops) != 2 || len(attention) != 2 || len(effects) != 6 {
		t.Fatalf("BudgetStopEffects() = %d effects, %d Stops, %d Attention items; want 6/2/2", len(effects), len(stops), len(attention))
	}
	for _, stop := range stops {
		if stop.Subject != (Subject{Type: "run", ID: "RUN-budget"}) || len(stop.BlockedActions) != 3 || stop.BlockedActions[0] != "gate.run" || stop.BlockedActions[1] != "run.recover" || stop.BlockedActions[2] != "run.start" {
			t.Fatalf("budget Stop = %#v, want run-scoped exact blocked actions", stop)
		}
	}
	if effects[0].Table != "stops" || effects[1].Table != "attention_items" || effects[2].Table != "attention_stop_links" {
		t.Fatalf("first budget effects = %#v, want Stop, Attention, link order", effects[:3])
	}
}

func TestReviewReworkStopRequiresApprovedOutcome(t *testing.T) {
	t.Parallel()
	workspace, _ := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer closeWorkflowDB(conn)
	block := BudgetAdmissionBlock{Code: "review-rework-exhausted", Budget: Budget{Name: "review-cycles", Unit: "review-cycles", Maximum: 1}, LastReviewID: "REV-1"}
	_, stops, _, err := BudgetStopEffects(conn, Subject{Type: "envelope", ID: "ENV-1"}, "ENV-1", "actor-reviewer", []BudgetAdmissionBlock{block}, "2026-09-20T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if len(stops) != 1 {
		t.Fatalf("Stops = %#v, want one", stops)
	}
	condition := stops[0].ResolutionCondition
	if condition.Type != "human-outcome-recorded" || condition.Subject != (Subject{Type: "envelope", ID: "ENV-1"}) || condition.OutcomeClass == nil || *condition.OutcomeClass != "review-rework-extension" || condition.BudgetName != nil || condition.CapabilityID != nil {
		t.Fatalf("rework resolution condition = %#v", condition)
	}
}

func TestGenericBudgetAdmissionExcludesOnlyReworkAction(t *testing.T) {
	t.Parallel()
	budgets := []Budget{
		{Name: "review-cycles", Unit: "review-cycles", ExhaustionAction: "review-rework-exhausted", Consumed: int64Pointer(1), Remaining: int64Pointer(0)},
		{Name: "attempts", Unit: "attempts", ExhaustionAction: "budget-exhausted", Consumed: int64Pointer(1), Remaining: int64Pointer(0)},
		{Name: "tokens", Unit: "input-tokens", ExhaustionAction: "budget-data-unavailable"},
	}
	blocks := GenericBudgetAdmissionBlocks(budgets)
	if len(blocks) != 2 || blocks[0].Code != "budget-exhausted" || blocks[0].Budget.Name != "attempts" || blocks[1].Code != "budget-data-unavailable" || blocks[1].Budget.Name != "tokens" {
		t.Fatalf("generic admission blocks = %#v, want attempts and tokens only", blocks)
	}
}

func TestResumeGateReasonsExcludeReworkBudget(t *testing.T) {
	t.Parallel()
	workspace, _ := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer closeWorkflowDB(conn)
	repository := Repository{ProjectID: "P-1", WorktreeID: "W-1", HeadOID: "head", TreeOID: "tree"}
	run := &Run{ID: "RUN-1", ActorID: "actor-builder", Status: "active", StartRepository: repository, Budgets: []Budget{
		{Name: "review-cycles", Unit: "review-cycles", ExhaustionAction: "review-rework-exhausted", Consumed: int64Pointer(1), Remaining: int64Pointer(0)},
		{Name: "attempts", Unit: "attempts", ExhaustionAction: "budget-exhausted", Consumed: int64Pointer(1), Remaining: int64Pointer(0)},
	}}
	action, present, err := resumeGateAction(conn, Envelope{ID: "ENV-1"}, run, repository)
	if err != nil || !present {
		t.Fatalf("resume gate action = %#v, present %t, err %v", action, present, err)
	}
	var attempts int
	for _, reason := range action.Reasons {
		if reason.Code == "budget-exhausted" && reason.Detail == "attempts" {
			attempts++
		}
		if reason.Detail == "review-cycles" || reason.Code == "budget-data-unavailable" {
			t.Fatalf("review budget leaked into resume: %#v", action.Reasons)
		}
	}
	if attempts != 1 {
		t.Fatalf("attempts budget reason count = %d, want one: %#v", attempts, action.Reasons)
	}
}
