package workflow

import (
	"database/sql"
	"encoding/json"
	"errors"
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

// errNotQualifying marks a fact that legitimately fails a Condition's
// qualifying-fact predicate, as distinct from an infrastructure failure.
var errNotQualifying = errors.New("fact does not qualify")

func notQualifying(format string, args ...any) error {
	return fmt.Errorf("%w: %s", errNotQualifying, fmt.Sprintf(format, args...))
}

// JudgmentFact is the durable content of one judgment Event: a
// human-outcome.recorded Event's Human Outcome or a decision.referenced
// Event's Decision Reference. It is the single input to the qualifying-fact
// predicates of the human-outcome-recorded and decision-satisfied Conditions
// (DEC-250 R1). Manual `stop resolve` builds it from the committed Event;
// the recording Event builds it, marked Pending, from the request it is about
// to commit. Both then call the same predicate, so the two paths cannot
// diverge.
type JudgmentFact struct {
	EventType string
	ID        string
	Class     string
	// Result is the Outcome's result, or the Reference's outcome.
	Result       string
	ScopeIssueID string
	// ActivationID is the Outcome's provenance Activation.
	ActivationID string
	// RequiredCapabilities is the Reference's recorded authority.
	RequiredCapabilities []string
	// EventID names the committed Event. It is empty while Pending.
	EventID string
	// Pending marks a fact whose Event has not committed. It is by
	// construction the newest fact of its kind and is absent from the
	// projection.
	Pending bool
}

func (fact JudgmentFact) approved() bool { return fact.Result == "approved" }

// PendingHumanOutcomeFact describes the Human Outcome an Event is about to
// commit.
func PendingHumanOutcomeFact(id, class, result, scopeIssueID, activationID string) JudgmentFact {
	return JudgmentFact{EventType: "human-outcome.recorded", ID: id, Class: class, Result: result, ScopeIssueID: scopeIssueID, ActivationID: activationID, Pending: true}
}

// PendingDecisionReferenceFact describes the Decision Reference an Event is
// about to commit.
func PendingDecisionReferenceFact(id, class, outcome, scopeIssueID string, requiredCapabilities []string) JudgmentFact {
	return JudgmentFact{EventType: "decision.referenced", ID: id, Class: class, Result: outcome, ScopeIssueID: scopeIssueID, RequiredCapabilities: requiredCapabilities, Pending: true}
}

// QualifyingJudgmentStops returns, in Stop ID order, every open Stop whose
// recorded Condition the fact qualifies under DEC-250 R2 or R3, whatever code
// raised the Stop. A Stop carrying any other Condition is never returned.
func QualifyingJudgmentStops(conn *sql.DB, fact JudgmentFact) ([]Stop, error) {
	rows, err := conn.Query(`
		SELECT s.id FROM stops s
		LEFT JOIN stop_resolutions r ON r.stop_id = s.id
		WHERE r.stop_id IS NULL AND json_extract(s.resolution_condition_json, '$.type') IN ('human-outcome-recorded', 'decision-satisfied')
		ORDER BY s.id`)
	if err != nil {
		return nil, fmt.Errorf("list open judgment Stops: %w", err)
	}
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	stops := []Stop{}
	for _, id := range ids {
		stop, err := loadStop(conn, id)
		if err != nil {
			return nil, err
		}
		if err := qualifyingJudgmentPredicate(conn, stop.ResolutionCondition, fact); err != nil {
			if errors.Is(err, errNotQualifying) {
				continue
			}
			return nil, err
		}
		stops = append(stops, stop)
	}
	return stops, nil
}

// qualifyingJudgmentCondition is manual `stop resolve`'s entry: it reads the
// committed Event's judgment row and applies the shared predicate.
func qualifyingJudgmentCondition(conn *sql.DB, event journal.Event, condition Condition) error {
	fact, err := judgmentFactFromEvent(event)
	if err != nil {
		return err
	}
	return qualifyingJudgmentPredicate(conn, condition, fact)
}

func qualifyingJudgmentPredicate(conn *sql.DB, condition Condition, fact JudgmentFact) error {
	switch condition.Type {
	case "human-outcome-recorded":
		return qualifyingHumanOutcomeCondition(conn, condition, fact)
	case "decision-satisfied":
		return qualifyingDecisionCondition(conn, condition, fact)
	}
	return notQualifying("condition %q is not resolved by a judgment fact", condition.Type)
}

func judgmentFactFromEvent(event journal.Event) (JudgmentFact, error) {
	table, ok := map[string]string{"human-outcome.recorded": "human_outcomes", "decision.referenced": "decision_refs"}[event.EventType]
	if !ok {
		return JudgmentFact{}, fmt.Errorf("event %q is not a judgment Event", event.EventType)
	}
	for _, effect := range event.Effects {
		if effect.Op != "insert" || effect.Table != table {
			continue
		}
		fact := JudgmentFact{EventType: event.EventType, EventID: event.EventID}
		fact.ID, _ = effect.PK["id"].(string)
		fact.Class, _ = effect.Values["class"].(string)
		fact.ScopeIssueID, _ = effect.Values["scope_issue_id"].(string)
		if event.EventType == "human-outcome.recorded" {
			fact.Result, _ = effect.Values["result"].(string)
			fact.ActivationID, _ = effect.Values["activation_id"].(string)
		} else {
			fact.Result, _ = effect.Values["outcome"].(string)
			raw, _ := effect.Values["authority_capabilities_json"].(string)
			if err := json.Unmarshal([]byte(raw), &fact.RequiredCapabilities); err != nil {
				return JudgmentFact{}, fmt.Errorf("decode Decision Reference authority: %w", err)
			}
		}
		if fact.ID == "" || fact.Class == "" || fact.Result == "" || fact.ScopeIssueID == "" {
			return JudgmentFact{}, fmt.Errorf("%s Event carries an incomplete %s row", event.EventType, table)
		}
		return fact, nil
	}
	return JudgmentFact{}, fmt.Errorf("%s Event carries no %s row", event.EventType, table)
}

// qualifyingHumanOutcomeCondition is DEC-250 R2, the single predicate for the
// human-outcome-recorded Condition. The fact must be a Human Outcome of the
// Condition's class with result approved, not superseded under contracts.md
// section 7, and inside the subject's scope.
func qualifyingHumanOutcomeCondition(conn *sql.DB, condition Condition, fact JudgmentFact) error {
	if condition.OutcomeClass == nil || *condition.OutcomeClass == "" {
		return fmt.Errorf("human-outcome-recorded condition names no outcome class")
	}
	if fact.EventType != "human-outcome.recorded" {
		return notQualifying("Event %q is not a Human Outcome", fact.EventType)
	}
	if fact.Class != *condition.OutcomeClass || !fact.approved() {
		return notQualifying("Human Outcome class or result does not match the condition")
	}
	if !fact.Pending {
		governing, err := ResolveGoverningHumanOutcome(conn, fact.Class, fact.ScopeIssueID)
		if err != nil {
			return err
		}
		if HumanOutcomeSuperseded(governing, fact.ID, fact.Result) {
			return notQualifying("Human Outcome %q is superseded", fact.ID)
		}
	}
	return humanOutcomeInSubjectScope(conn, condition.Subject, fact)
}

func humanOutcomeInSubjectScope(conn *sql.DB, subject Subject, fact JudgmentFact) error {
	switch subject.Type {
	case "envelope":
		var found int
		err := conn.QueryRow(`SELECT 1 FROM work_envelopes e JOIN envelope_issue_selection s ON s.envelope_id = e.id AND s.revision = e.current_revision WHERE e.id = ? AND s.issue_id = ?`, subject.ID, fact.ScopeIssueID).Scan(&found)
		if errors.Is(err, sql.ErrNoRows) {
			return notQualifying("Outcome scope Issue is not selected in the Envelope's current revision")
		}
		return err
	case "run":
		var issueID string
		err := conn.QueryRow(`SELECT issue_id FROM runs WHERE id = ?`, subject.ID).Scan(&issueID)
		if errors.Is(err, sql.ErrNoRows) {
			return notQualifying("Run %q does not exist", subject.ID)
		}
		if err != nil {
			return err
		}
		if issueID != fact.ScopeIssueID {
			return notQualifying("Outcome scope Issue is not the Run's Issue")
		}
		return nil
	case "activation":
		if fact.ActivationID != subject.ID {
			return notQualifying("Outcome provenance Activation is not the condition's Activation")
		}
		return nil
	}
	return notQualifying("condition subject type %q has no Outcome scope", subject.Type)
}

// qualifyingDecisionCondition is DEC-250 R3, the single predicate for the
// decision-satisfied Condition. The subject resolves to an Envelope (a Run's
// Envelope for a run subject); that Envelope's current revision must declare a
// protected Decision of the class, which must be satisfied under contracts.md
// section 9, and the fact must be the pair-completing Event: the later-created
// of the bound Decision Reference and Human Outcome.
func qualifyingDecisionCondition(conn *sql.DB, condition Condition, fact JudgmentFact) error {
	if condition.DecisionClass == nil || *condition.DecisionClass == "" {
		return fmt.Errorf("decision-satisfied condition names no Decision class")
	}
	if fact.EventType != "decision.referenced" && fact.EventType != "human-outcome.recorded" {
		return notQualifying("Event %q is not a judgment Event", fact.EventType)
	}
	if condition.Subject.Type != "envelope" && condition.Subject.Type != "run" {
		return notQualifying("condition subject type %q has no Envelope scope", condition.Subject.Type)
	}
	envelopeID, _, err := stopConditionEnvelope(conn, condition.Subject)
	if errors.Is(err, sql.ErrNoRows) {
		return notQualifying("condition subject %q does not exist", condition.Subject.ID)
	}
	if err != nil {
		return err
	}
	envelope, err := loadEnvelope(conn, envelopeID)
	if err != nil {
		return err
	}
	for _, declared := range envelope.ProtectedDecisions {
		if declared.DecisionClass != *condition.DecisionClass {
			continue
		}
		var pending *JudgmentFact
		if fact.Pending {
			pending = &fact
		}
		binding, err := deriveProtectedDecision(conn, envelope.PrimaryIssueID, &declared, pending)
		if err != nil {
			return err
		}
		if declared.Status != "satisfied" {
			return notQualifying("protected Decision %q is not satisfied", declared.DecisionClass)
		}
		if !isPairCompleting(fact, binding) {
			return notQualifying("Event is not the one that completed the Decision pair")
		}
		return nil
	}
	return notQualifying("Envelope %q declares no protected Decision of class %q", envelopeID, *condition.DecisionClass)
}

// isPairCompleting reports whether the fact is whichever of the bound Decision
// Reference's and Human Outcome's creating Events committed later. A pending
// fact that is itself bound is the later one by construction.
func isPairCompleting(fact JudgmentFact, binding protectedDecisionBinding) bool {
	if fact.Pending {
		return (fact.EventType == "decision.referenced" && binding.RefPending) || (fact.EventType == "human-outcome.recorded" && binding.OutcomePending)
	}
	return fact.EventID != "" && fact.EventID == max(binding.RefEventID, binding.OutcomeEventID)
}
