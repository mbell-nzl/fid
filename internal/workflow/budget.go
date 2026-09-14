package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

// BudgetSample is one append-only observed measurement in a budget chain.
// A nil Value is an explicitly unknown observation.
type BudgetSample struct {
	ID               string  `json:"id"`
	EnvelopeID       string  `json:"envelope_id"`
	EnvelopeRevision int     `json:"envelope_revision"`
	BudgetName       string  `json:"budget_name"`
	RunID            string  `json:"run_id"`
	ActorID          string  `json:"actor_id"`
	Source           string  `json:"source"`
	SourceSampleID   string  `json:"source_sample_id"`
	SampleKind       string  `json:"sample_kind"`
	PreviousSampleID *string `json:"previous_sample_id"`
	Value            *int64  `json:"value"`
	MeasurementRef   *Source `json:"measurement_ref"`
	ObservedAt       string  `json:"observed_at"`
	EventID          string  `json:"event_id"`
}

// BudgetAdmission is the state derived while admitting an execution fact.
type BudgetAdmission struct {
	Admitted  bool
	Budget    Budget
	Exhausted bool
	Warning   bool
}

// BudgetAdmissionBlock identifies the durable reason execution cannot begin.
type BudgetAdmissionBlock struct {
	Code   string
	Budget Budget
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
		budget.State = "unknown"
		return budget, nil
	}
	budget.Consumed = &consumed
	remaining := budget.Maximum - consumed
	if remaining < 0 {
		remaining = 0
	}
	budget.Remaining = &remaining
	budget.State = budgetState(budget)
	return budget, nil
}

func budgetState(budget Budget) string {
	if budget.Consumed == nil || budget.Remaining == nil {
		return "unknown"
	}
	if *budget.Remaining == 0 {
		return "exhausted"
	}
	if budget.WarningThreshold != nil && *budget.Remaining <= *budget.WarningThreshold {
		return "warning"
	}
	return "available"
}

// budgetPolicyValue removes derived read fields before an immutable budget is
// persisted into policy or an Event Effect.
func budgetPolicyValue(budget Budget) Budget {
	budget.Consumed = nil
	budget.Remaining = nil
	budget.State = ""
	return budget
}

// DeriveBudget returns the current budget state for an already durable sample
// chain. It does not append or normalize a sample.
func DeriveBudget(budget Budget, samples []BudgetSample) (Budget, error) {
	return deriveBudget(budget, samples)
}

func budgetAdmission(budget Budget) BudgetAdmission {
	exhausted := budget.Remaining != nil && *budget.Remaining == 0
	warning := budget.WarningThreshold != nil && budget.Remaining != nil && *budget.Remaining <= *budget.WarningThreshold && !exhausted
	return BudgetAdmission{Admitted: true, Budget: budget, Exhausted: exhausted, Warning: warning}
}

// BudgetAdmissionBlocks returns each unknown or exhausted budget that blocks
// a new execution. Warning budgets remain admissible.
func BudgetAdmissionBlocks(budgets []Budget) []BudgetAdmissionBlock {
	blocks := make([]BudgetAdmissionBlock, 0, len(budgets))
	for _, budget := range budgets {
		code := ""
		switch {
		case budget.Consumed == nil || budget.Remaining == nil:
			code = "budget-data-unavailable"
		case *budget.Remaining == 0:
			code = "budget-exhausted"
		}
		if code != "" {
			blocks = append(blocks, BudgetAdmissionBlock{Code: code, Budget: budget})
		}
	}
	return blocks
}

// BudgetStopEffects prepares the Stop/Attention effects that must share the
// Event which denied execution or crossed a budget exhaustion boundary.
func BudgetStopEffects(conn *sql.DB, subject Subject, envelopeID, actorID string, blocks []BudgetAdmissionBlock, raisedAt string) ([]journal.Effect, []Stop, []Attention, error) {
	if conn == nil || !validSubject(subject) || envelopeID == "" || actorID == "" {
		return nil, nil, nil, fmt.Errorf("budget Stop effects are invalid")
	}
	if raisedAt == "" {
		raisedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	effects := []journal.Effect{}
	stops := make([]Stop, 0, len(blocks))
	attention := make([]Attention, 0, len(blocks))
	preparedAttention := map[string]Attention{}
	for _, block := range blocks {
		if block.Code != "budget-exhausted" && block.Code != "budget-data-unavailable" {
			return nil, nil, nil, fmt.Errorf("budget Stop code is invalid")
		}
		condition := Condition{Type: "work-ended", Subject: subject}
		if block.Code == "budget-data-unavailable" {
			name := block.Budget.Name
			condition = Condition{Type: "budget-observed", Subject: subject, BudgetName: &name}
		}
		details, err := json.Marshal(map[string]any{
			"envelope_id": envelopeID,
			"budget_name": block.Budget.Name,
			"unit":        block.Budget.Unit,
			"maximum":     block.Budget.Maximum,
			"consumed":    block.Budget.Consumed,
			"remaining":   block.Budget.Remaining,
		})
		if err != nil {
			return nil, nil, nil, err
		}
		request := StopRaiseRequest{Subject: subject, Code: block.Code, BlockedActions: []string{"gate.run", "run.recover", "run.start"}, Details: details, ResolutionCondition: condition, ResolverCapability: "stop.resolve"}
		if err := validateStopRaise(request); err != nil {
			return nil, nil, nil, err
		}
		stopID, err := newStopID(conn)
		if err != nil {
			return nil, nil, nil, err
		}
		key, err := attentionSemanticKey(subject, block.Code, condition)
		if err != nil {
			return nil, nil, nil, err
		}
		item, found, err := loadOpenAttentionByKey(conn, key)
		if err != nil {
			return nil, nil, nil, err
		}
		if prepared, exists := preparedAttention[key]; exists {
			item, found = prepared, true
		}
		if !found {
			item.ID, err = newAttentionID(conn)
			if err != nil {
				return nil, nil, nil, err
			}
			item = Attention{ID: item.ID, SemanticKey: key, Subject: subject, Code: block.Code, Status: "open", Version: 1}
		}
		stop := Stop{ID: stopID, Subject: subject, Code: block.Code, BlockedActions: request.BlockedActions, Details: details, ResolutionCondition: condition, ResolverCapability: request.ResolverCapability, ActorID: actorID, RaisedAt: raisedAt}
		if !found {
			item.RaisedAt = raisedAt
		}
		item.StopIDs = append(item.StopIDs, stopID)
		sort.Strings(item.StopIDs)
		preparedAttention[key] = item
		stopEffects := stopRaisedEventForDomain(stop, item, !found, "<event-id>").Effects
		effects = append(effects, stopEffects...)
		stops, attention = append(stops, stop), append(attention, item)
	}
	return effects, stops, attention, nil
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
