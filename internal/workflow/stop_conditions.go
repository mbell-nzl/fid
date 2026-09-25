package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/journal"
)

func qualifyingExecutorTermination(conn *sql.DB, event journal.Event, condition Condition) error {
	if condition.Subject.Type != "run" || event.Aggregate.Type != "run" || event.Aggregate.ID != condition.Subject.ID {
		return fmt.Errorf("executor termination fact belongs to another Run")
	}
	var executorID string
	if err := conn.QueryRow(`SELECT executor_id FROM run_runtime WHERE run_id = ?`, condition.Subject.ID).Scan(&executorID); err != nil {
		return err
	}
	for _, effect := range event.Effects {
		if effect.Op != "insert" || effect.Table != "run_observations" || effect.Values["run_id"] != condition.Subject.ID || effect.Values["kind"] != "liveness" {
			continue
		}
		var observation RunObservation
		raw, ok := effect.Values["observation_json"].(string)
		if !ok || json.Unmarshal([]byte(raw), &observation) != nil || observation.ID != effect.PK["id"] || observation.RunID != condition.Subject.ID || observation.Kind != "liveness" || observation.Liveness.ExecutorID != executorID || observation.Liveness.State != "terminated" {
			return fmt.Errorf("executor observation does not prove termination")
		}
		current, err := latestLiveness(conn, condition.Subject.ID)
		if err != nil {
			return err
		}
		if current == nil || current.ExecutorID != executorID || current.State != "terminated" {
			return fmt.Errorf("current executor state is not terminated")
		}
		return nil
	}
	return fmt.Errorf("executor termination fact has no matching observation")
}

func qualifyingActivationCondition(conn *sql.DB, event journal.Event, condition Condition) error {
	if condition.Subject.Type != "activation" || event.Aggregate.Type != "activation" || event.Aggregate.ID != condition.Subject.ID {
		return fmt.Errorf("activation fact belongs to another subject")
	}
	activation, err := loadActivation(conn, condition.Subject.ID)
	if err != nil {
		return err
	}
	if condition.Type == "activation-recovered" && activation.Status == "aborted" {
		if event.EventType != "activation.recovery-completed" || activation.Recovery == nil || activation.Recovery.CompletedEventID == nil || *activation.Recovery.CompletedEventID != event.EventID || activation.Recovery.Direction != "restore-prior" {
			return fmt.Errorf("activation fact does not prove completed restoration")
		}
		return requirePriorResourceAgreement(activation)
	}
	if activation.Status != "finalized" {
		return fmt.Errorf("activation is not the current finalized selection")
	}
	if err := requireNewResourceAgreement(conn, activation.ID); err != nil {
		return err
	}
	if condition.Type != "capability-available" {
		return nil
	}
	var backend string
	if err := conn.QueryRow(`SELECT new_identity FROM activation_resources WHERE activation_id = ? AND component = 'backend'`, activation.ID).Scan(&backend); err != nil {
		return err
	}
	if backend != activation.BackendIdentity || condition.CapabilityID == nil {
		return fmt.Errorf("capability backend or named capability does not match")
	}
	for _, capability := range capabilityMatrix {
		if capability.capabilityID == *condition.CapabilityID && capability.available {
			return nil
		}
	}
	return fmt.Errorf("named capability is not available in the current binary")
}

func qualifyingRecoveryCondition(conn *sql.DB, event journal.Event, condition Condition) error {
	if condition.Subject.Type != "run" || event.Aggregate.Type != "run" || event.Aggregate.ID != condition.Subject.ID {
		return fmt.Errorf("recovery fact does not name the predecessor Run")
	}
	var recoveryID string
	return conn.QueryRow(`SELECT id FROM run_recoveries WHERE prior_run_id = ? AND created_event_id = ?`, condition.Subject.ID, event.EventID).Scan(&recoveryID)
}

func qualifyingBudgetCondition(conn *sql.DB, event journal.Event, condition Condition) error {
	if condition.BudgetName == nil || event.Aggregate.Type != "run" {
		return fmt.Errorf("budget fact has no named budget or Run binding")
	}
	envelopeID, revision, err := stopConditionEnvelope(conn, condition.Subject)
	if err != nil {
		return err
	}
	matched := false
	for _, effect := range event.Effects {
		if effect.Op != "insert" || effect.Table != "budget_samples" || effect.Values["envelope_id"] != envelopeID || effect.Values["budget_name"] != *condition.BudgetName || effect.Values["run_id"] != event.Aggregate.ID {
			continue
		}
		var sampleID string
		err := conn.QueryRow(`SELECT id FROM budget_samples WHERE id = ? AND envelope_id = ? AND envelope_revision = ? AND budget_name = ? AND created_event_id = ? AND run_id = ?`, effect.PK["id"], envelopeID, revision, *condition.BudgetName, event.EventID, event.Aggregate.ID).Scan(&sampleID)
		if err != nil {
			return err
		}
		matched = true
	}
	if !matched {
		return fmt.Errorf("budget fact does not belong to the named Envelope budget chain")
	}
	budget, err := loadRunBudget(conn, envelopeID, revision, *condition.BudgetName)
	if err != nil {
		return err
	}
	samples, err := loadRunBudgetSamples(conn, envelopeID, revision, *condition.BudgetName)
	if err != nil {
		return err
	}
	current, err := deriveBudget(budget, samples)
	if err != nil {
		return err
	}
	if current.Consumed == nil || current.Remaining == nil || *current.Remaining <= 0 {
		return fmt.Errorf("current named budget is unknown or exhausted")
	}
	return nil
}

func stopConditionEnvelope(conn *sql.DB, subject Subject) (string, int, error) {
	var id string
	var revision int
	var err error
	switch subject.Type {
	case "run":
		err = conn.QueryRow(`SELECT envelope_id, envelope_revision FROM runs WHERE id = ?`, subject.ID).Scan(&id, &revision)
	case "envelope":
		err = conn.QueryRow(`SELECT id, current_revision FROM work_envelopes WHERE id = ?`, subject.ID).Scan(&id, &revision)
	default:
		err = fmt.Errorf("condition subject has no Envelope")
	}
	return id, revision, err
}

func qualifyingWorkEnded(conn *sql.DB, event journal.Event, condition Condition) error {
	if event.EventType == "envelope.completed" {
		if condition.Subject.Type != "envelope" || event.Aggregate.Type != "envelope" || event.Aggregate.ID != condition.Subject.ID {
			return fmt.Errorf("completed Envelope fact belongs to another subject")
		}
		var status string
		if err := conn.QueryRow(`SELECT status FROM work_envelopes WHERE id = ?`, condition.Subject.ID).Scan(&status); err != nil {
			return err
		}
		if status != "completed" {
			return fmt.Errorf("Envelope is not currently completed")
		}
		return requireNoActiveRun(conn, condition.Subject.ID)
	}
	if event.Aggregate.Type != "run" {
		return fmt.Errorf("work-ended fact is not a Run fact")
	}
	var envelopeID, status string
	if err := conn.QueryRow(`SELECT envelope_id, status FROM runs WHERE id = ?`, event.Aggregate.ID).Scan(&envelopeID, &status); err != nil {
		return err
	}
	if status != "succeeded" && status != "failed" && status != "abandoned" {
		return fmt.Errorf("Run is not currently terminal")
	}
	if condition.Subject.Type == "run" && condition.Subject.ID == event.Aggregate.ID {
		return nil
	}
	if condition.Subject.Type == "envelope" && condition.Subject.ID == envelopeID {
		return requireNoActiveRun(conn, envelopeID)
	}
	return fmt.Errorf("terminal Run belongs to another condition subject")
}

// qualifyingReviewRecordedCondition proves a review.recorded Event belongs to
// the Condition's own subject (its recording Run or its Envelope); the type
// carries no further discriminant (contracts.md section 11).
func qualifyingReviewRecordedCondition(conn *sql.DB, event journal.Event, condition Condition) error {
	if !reviewRecordEffectMatchesSubject(event, condition.Subject) {
		return fmt.Errorf("review record fact belongs to another subject")
	}
	return nil
}

// qualifyingFindingResolvedCondition additionally requires the review.recorded
// Event to carry a 'resolved' disposition for the named Finding; 'addressed'
// and 'inapplicable' dispositions never qualify.
func qualifyingFindingResolvedCondition(conn *sql.DB, event journal.Event, condition Condition) error {
	if condition.FindingID == nil {
		return fmt.Errorf("finding-resolved condition has no named Finding")
	}
	if !reviewRecordEffectMatchesSubject(event, condition.Subject) {
		return fmt.Errorf("review record fact belongs to another subject")
	}
	for _, effect := range event.Effects {
		if effect.Op == "insert" && effect.Table == "finding_dispositions" && effect.Values["action"] == "resolved" && effect.Values["finding_id"] == *condition.FindingID {
			return nil
		}
	}
	return fmt.Errorf("review record fact does not resolve the named Finding")
}

func reviewRecordEffectMatchesSubject(event journal.Event, subject Subject) bool {
	for _, effect := range event.Effects {
		if effect.Op != "insert" || effect.Table != "review_records" {
			continue
		}
		switch subject.Type {
		case "run":
			return effect.Values["run_id"] == subject.ID
		case "envelope":
			return effect.Values["envelope_id"] == subject.ID
		}
	}
	return false
}

func requireNoActiveRun(conn *sql.DB, envelopeID string) error {
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM runs WHERE envelope_id = ? AND status = 'active'`, envelopeID).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("Envelope still has an active Run")
	}
	return nil
}
