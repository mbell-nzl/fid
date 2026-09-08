package maintenance

import (
	"database/sql"
	"fmt"
)

func recoveryIntegrityFindings(db *sql.DB) ([]Finding, error) {
	checks := []struct {
		query, check, message string
	}{
		{
			query: `SELECT recovery.id FROM run_recoveries AS recovery
LEFT JOIN runs AS prior ON prior.id = recovery.prior_run_id
LEFT JOIN runs AS successor ON successor.id = recovery.successor_run_id
LEFT JOIN run_observations AS observation ON observation.id = recovery.observation_id
LEFT JOIN claims AS successor_claim ON successor_claim.id = recovery.successor_claim_id
WHERE prior.claim_id IS NOT recovery.prior_claim_id
   OR successor.claim_id IS NOT recovery.successor_claim_id
   OR observation.run_id IS NOT recovery.prior_run_id
   OR successor.issue_id IS NOT prior.issue_id
   OR successor.envelope_id IS NOT prior.envelope_id
   OR successor.envelope_revision IS NOT prior.envelope_revision
   OR successor_claim.issue_id IS NOT successor.issue_id
   OR successor_claim.envelope_id IS NOT successor.envelope_id
   OR successor_claim.envelope_revision IS NOT successor.envelope_revision
ORDER BY recovery.id`,
			check: "recovery-link-integrity", message: "recovery predecessor, successor, claim, or observation does not agree",
		},
		{
			query: `SELECT sample.id FROM budget_samples AS sample
JOIN runs AS run ON run.id = sample.run_id
WHERE sample.envelope_id IS NOT run.envelope_id
   OR sample.envelope_revision IS NOT run.envelope_revision
   OR NOT EXISTS (
		SELECT 1 FROM run_budget_bindings AS binding
		WHERE binding.run_id = sample.run_id AND binding.name = sample.budget_name
	)
ORDER BY sample.id`,
			check: "budget-binding-integrity", message: "budget sample does not bind the Run envelope budget",
		},
		{
			query: `SELECT recovery.activation_id FROM activation_recoveries AS recovery
JOIN kit_activations AS activation ON activation.id = recovery.activation_id
WHERE (recovery.completed_event_id IS NULL AND activation.status NOT IN ('prepared', 'staged', 'committed'))
   OR (recovery.completed_event_id IS NOT NULL AND (
		(recovery.direction = 'restore-prior' AND activation.status <> 'aborted')
		OR (recovery.direction = 'complete-forward' AND activation.status NOT IN ('finalized', 'superseded'))
	))
ORDER BY recovery.activation_id`,
			check: "activation-recovery-integrity", message: "activation recovery completion does not agree with activation status",
		},
	}
	findings := make([]Finding, 0)
	for _, check := range checks {
		rows, err := db.Query(check.query)
		if err != nil {
			return nil, fmt.Errorf("check %s: query rows: %w", check.check, err)
		}
		items, err := simpleFindings(rows, check.check, check.message)
		closeErr := rows.Close()
		if err != nil {
			return nil, err
		}
		if closeErr != nil {
			return nil, fmt.Errorf("check %s: close rows: %w", check.check, closeErr)
		}
		findings = append(findings, items...)
	}
	return findings, nil
}
