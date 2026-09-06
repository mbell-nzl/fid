package workflow

import (
	"fmt"
	"math"
)

// BudgetSample is one append-only observed measurement in a budget chain.
// A nil Value is an explicitly unknown observation.
type BudgetSample struct {
	ID               string  `json:"id"`
	Source           string  `json:"source"`
	SourceSampleID   string  `json:"source_sample_id"`
	SampleKind       string  `json:"sample_kind"`
	PreviousSampleID *string `json:"previous_sample_id"`
	Value            *int64  `json:"value"`
}

// BudgetAdmission is the state derived while admitting an execution fact.
type BudgetAdmission struct {
	Admitted  bool
	Budget    Budget
	Exhausted bool
	Warning   bool
}

type BudgetSampleConflictError struct {
	BudgetName               string
	SourceSampleID           string
	ExpectedPreviousSampleID *string
	ActualPreviousSampleID   *string
}

func (err *BudgetSampleConflictError) Error() string { return "budget sample conflict" }

// AdmitBudgetSample derives the result of appending one observed sample. The
// durable appliers use the same rules before preparing their Event effects.
func AdmitBudgetSample(budget Budget, samples []BudgetSample, sample BudgetSample) (BudgetAdmission, error) {
	if err := validateBudgetSample(budget, sample); err != nil {
		return BudgetAdmission{}, fmt.Errorf("budget sample is invalid")
	}
	if sample.ID == "" {
		sample.ID = sample.SourceSampleID
	}
	for _, prior := range samples {
		if prior.Source == sample.Source && prior.SourceSampleID == sample.SourceSampleID {
			if sameBudgetSample(prior, sample) {
				derived, err := deriveBudget(budget, samples)
				if err != nil {
					return BudgetAdmission{}, err
				}
				return budgetAdmission(derived), nil
			}
			return BudgetAdmission{}, &BudgetSampleConflictError{BudgetName: budget.Name, SourceSampleID: sample.SourceSampleID, ExpectedPreviousSampleID: prior.PreviousSampleID, ActualPreviousSampleID: sample.PreviousSampleID}
		}
	}
	if len(samples) == 0 {
		if sample.PreviousSampleID != nil {
			return BudgetAdmission{}, &BudgetSampleConflictError{BudgetName: budget.Name, SourceSampleID: sample.SourceSampleID, ExpectedPreviousSampleID: nil, ActualPreviousSampleID: sample.PreviousSampleID}
		}
	} else {
		expected := samples[len(samples)-1].ID
		if sample.PreviousSampleID == nil || *sample.PreviousSampleID != expected {
			return BudgetAdmission{}, &BudgetSampleConflictError{BudgetName: budget.Name, SourceSampleID: sample.SourceSampleID, ExpectedPreviousSampleID: &expected, ActualPreviousSampleID: sample.PreviousSampleID}
		}
	}
	derived, err := deriveBudget(budget, append(append([]BudgetSample{}, samples...), sample))
	if err != nil {
		return BudgetAdmission{}, err
	}
	return budgetAdmission(derived), nil
}

func validateBudgetSample(budget Budget, sample BudgetSample) error {
	if budget.Maximum < 1 || sample.SourceSampleID == "" || (sample.Source != "fid" && sample.Source != "activated-harness") || (sample.SampleKind != "delta" && sample.SampleKind != "baseline") || (sample.Value != nil && *sample.Value < 0) {
		return fmt.Errorf("invalid")
	}
	if sample.Source == "fid" && !fidOwnedBudgetUnit(budget.Unit) {
		return fmt.Errorf("fid source may only charge owned units")
	}
	return nil
}

func fidOwnedBudgetUnit(unit string) bool {
	return unit == "attempts" || unit == "failures" || unit == "milliseconds" || unit == "tool-failures" || unit == "interventions"
}

func validBudgetUnit(unit string) bool {
	switch unit {
	case "attempts", "failures", "milliseconds", "input-tokens", "output-tokens", "micro-usd", "tool-failures", "diff-bytes", "interventions":
		return true
	default:
		return false
	}
}

func deriveBudget(budget Budget, samples []BudgetSample) (Budget, error) {
	known := fidOwnedBudgetUnit(budget.Unit)
	var consumed int64
	for _, sample := range samples {
		if sample.Value == nil {
			known = false
			continue
		}
		if sample.SampleKind == "baseline" {
			if known && *sample.Value < consumed {
				return Budget{}, &BudgetSampleConflictError{BudgetName: budget.Name, SourceSampleID: sample.SourceSampleID}
			}
			consumed, known = *sample.Value, true
			continue
		}
		if !known {
			continue
		}
		if *sample.Value > math.MaxInt64-consumed {
			return Budget{}, fmt.Errorf("budget consumption overflow")
		}
		consumed += *sample.Value
	}
	if !known {
		budget.Consumed, budget.Remaining = nil, nil
		return budget, nil
	}
	budget.Consumed = &consumed
	remaining := budget.Maximum - consumed
	if remaining < 0 {
		remaining = 0
	}
	budget.Remaining = &remaining
	return budget, nil
}

func budgetAdmission(budget Budget) BudgetAdmission {
	exhausted := budget.Remaining != nil && *budget.Remaining == 0
	warning := budget.WarningThreshold != nil && budget.Remaining != nil && *budget.Remaining <= *budget.WarningThreshold && !exhausted
	return BudgetAdmission{Admitted: true, Budget: budget, Exhausted: exhausted, Warning: warning}
}

func sameBudgetSample(left, right BudgetSample) bool {
	leftID, rightID := left.ID, right.ID
	if leftID == "" {
		leftID = left.SourceSampleID
	}
	if rightID == "" {
		rightID = right.SourceSampleID
	}
	return leftID == rightID && left.Source == right.Source && left.SourceSampleID == right.SourceSampleID && left.SampleKind == right.SampleKind && sameStringPointer(left.PreviousSampleID, right.PreviousSampleID) && sameInt64Pointer(left.Value, right.Value)
}

func sameInt64Pointer(left, right *int64) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}
