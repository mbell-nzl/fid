package workflow

import (
	"errors"
	"math"
	"testing"
)

func TestFinalAllowedAttemptChargesAndRaisesStopAtomically(t *testing.T) {
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
	budget := Budget{Name: "attempt-budget", Scope: "envelope", Unit: "attempts", Maximum: math.MaxInt64, ExhaustionAction: "budget-exhausted"}
	first := BudgetSample{ID: "one", Source: "fid", SourceSampleID: "one", SampleKind: "delta", Value: int64Pointer(math.MaxInt64)}
	previous := "one"
	second := BudgetSample{Source: "fid", SourceSampleID: "two", SampleKind: "delta", PreviousSampleID: &previous, Value: int64Pointer(1)}
	if _, err := AdmitBudgetSample(budget, []BudgetSample{first}, second); err == nil {
		t.Fatal("overflow admission error = nil, want rejection")
	}
}
