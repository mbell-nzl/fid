package workflow

import (
	"database/sql"
	"fmt"

	"github.com/mbell-nzl/fid/internal/journal"
)

// The successor inherits the Envelope chain. Its creation and the intervention
// are one transaction, so the boundary Stop guards the newly active Run.
func recoveryInterventionEffects(conn *sql.DB, successor *Run, recovery Recovery) ([]journal.Effect, error) {
	if recovery.Authority.Mode != "force" {
		return nil, nil
	}
	successor.Budgets = append([]Budget{}, successor.Budgets...)
	effects := []journal.Effect{}
	blocks := []BudgetAdmissionBlock{}
	for index, budget := range successor.Budgets {
		if budget.Unit != "interventions" {
			continue
		}
		samples, err := loadRunBudgetSamples(conn, successor.EnvelopeID, successor.EnvelopeRevision, budget.Name)
		if err != nil {
			return nil, err
		}
		id, err := uniqueRecoveryID(conn, "BSM", "budget_samples")
		if err != nil {
			return nil, err
		}
		value := int64(1)
		sample := BudgetSample{ID: id, EnvelopeID: successor.EnvelopeID, EnvelopeRevision: successor.EnvelopeRevision, BudgetName: budget.Name, RunID: successor.ID, ActorID: successor.ActorID, Source: "fid", SourceSampleID: recovery.ID, SampleKind: "delta", Value: &value, ObservedAt: successor.StartedAt}
		if len(samples) != 0 {
			previous := samples[len(samples)-1].ID
			sample.PreviousSampleID = &previous
		}
		admission, err := AdmitBudgetSample(budget, samples, sample)
		if err != nil {
			return nil, fmt.Errorf("admit recovery intervention %q: %w", budget.Name, err)
		}
		blocks = append(blocks, budgetBoundaryBlocks(budget, admission.Budget)...)
		successor.Budgets[index] = admission.Budget
		effects = append(effects, journal.Effect{Op: "insert", Table: "budget_samples", PK: map[string]any{"id": id}, Values: budgetSampleValues(sample)})
	}
	stopEffects, _, _, err := BudgetStopEffects(conn, Subject{Type: "run", ID: successor.ID}, successor.EnvelopeID, successor.ActorID, blocks, successor.StartedAt)
	if err != nil {
		return nil, err
	}
	return append(effects, stopEffects...), nil
}
