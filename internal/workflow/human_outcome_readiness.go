package workflow

import (
	"database/sql"

	"github.com/mbell-nzl/fid/internal/journal"
)

// onlyUnsatisfiedHumanOutcomeObligation returns the sole unsatisfied
// Obligation when it is the only one remaining and its type is human-outcome
// (DEC-224, contracts.md section 10). Every other combination returns false,
// including when the sole unsatisfied Obligation is not a human-outcome type.
func onlyUnsatisfiedHumanOutcomeObligation(obligations []Obligation, evaluations []ObligationEvaluation) (Obligation, bool) {
	unsatisfiedCount := 0
	var target Obligation
	found := false
	for _, evaluation := range evaluations {
		if evaluation.Status == "satisfied" {
			continue
		}
		unsatisfiedCount++
		for _, obligation := range obligations {
			if obligation.ID == evaluation.ObligationID {
				target, found = obligation, true
			}
		}
	}
	if unsatisfiedCount != 1 || !found || target.Type != "human-outcome" || target.OutcomeClass == nil || *target.OutcomeClass == "" {
		return Obligation{}, false
	}
	return target, true
}

// humanOutcomeReadinessEffectsForEnvelope raises a human-outcome-required Stop,
// atomically with the caller's own committing Event, exactly when the given
// Envelope's human-outcome Obligation is the sole unsatisfied Obligation. It
// operates on an already-loaded Envelope so callers mid-mutation (whose own
// projected status has not yet changed) evaluate the correct in-flight state.
// At most one unresolved Stop exists per Obligation: when one is already open
// for the exact same subject, code, and Condition, this raises nothing.
//
// overrides forces one or more Obligation IDs to the given evaluation status,
// keyed by Obligation ID. It exists because Effects are computed before the
// caller's own committing Event reaches the projection, so a live re-query
// inside the same transaction cannot yet see that Event's own Effect (for
// example, the Evidence row evidence.link is about to insert, or the
// Artifact availability flip artifact.verify is about to apply). Callers pass
// the known post-Effect status for exactly the Obligation(s) their own
// mutation affects; every other Obligation is read from committed state.
func humanOutcomeReadinessEffectsForEnvelope(conn *sql.DB, envelope Envelope, actorID string, repository Repository, overrides map[string]string) ([]journal.Effect, []Stop, []Attention, error) {
	evaluations, err := currentEnvelopeEvaluations(conn, envelope, repository)
	if err != nil {
		return nil, nil, nil, err
	}
	for index := range evaluations {
		if forced, ok := overrides[evaluations[index].ObligationID]; ok {
			evaluations[index].Status = forced
		}
	}
	obligation, ok := onlyUnsatisfiedHumanOutcomeObligation(envelope.Obligations, evaluations)
	if !ok {
		return nil, nil, nil, nil
	}
	class := *obligation.OutcomeClass
	request := StopRaiseRequest{
		Subject:             Subject{Type: "envelope", ID: envelope.ID},
		Code:                "human-outcome-required",
		BlockedActions:      obligation.GuardedActions,
		ResolutionCondition: Condition{Type: "human-outcome-recorded", Subject: Subject{Type: "envelope", ID: envelope.ID}, OutcomeClass: &class},
		ResolverCapability:  "stop.resolve",
	}
	request.Details = mustEnvelopeJSON(map[string]any{"envelope_id": envelope.ID, "obligation_id": obligation.ID, "outcome_class": class})
	if err := validateStopRaise(request); err != nil {
		return nil, nil, nil, err
	}
	key, err := attentionSemanticKey(request.Subject, request.Code, request.ResolutionCondition)
	if err != nil {
		return nil, nil, nil, err
	}
	if _, found, err := loadOpenAttentionByKey(conn, key); err != nil {
		return nil, nil, nil, err
	} else if found {
		return nil, nil, nil, nil
	}
	stopID, err := newStopID(conn)
	if err != nil {
		return nil, nil, nil, err
	}
	attentionID, err := newAttentionID(conn)
	if err != nil {
		return nil, nil, nil, err
	}
	attention := Attention{ID: attentionID, SemanticKey: key, Subject: request.Subject, Code: request.Code, Status: "open", Version: 1, StopIDs: []string{stopID}, RaisedAt: "<occurred-at>"}
	stop := Stop{ID: stopID, Subject: request.Subject, Code: request.Code, BlockedActions: request.BlockedActions, Details: request.Details, ResolutionCondition: request.ResolutionCondition, ResolverCapability: request.ResolverCapability, ActorID: actorID, RaisedAt: "<occurred-at>"}
	raised := stopRaisedEventForDomain(stop, attention, true, "<event-id>")
	return raised.Effects, []Stop{stop}, []Attention{attention}, nil
}

// HumanOutcomeReadinessStopEffects is the cross-package entry point used by
// the evidence.linked, review.recorded, and artifact.availability-observed
// producers (DEC-224). It reloads the Envelope so producers outside this
// package observe committed state, and requires the Envelope still be active
// at the expected revision. See humanOutcomeReadinessEffectsForEnvelope for
// what overrides is for; pass nil when the caller's own mutation does not
// affect any Obligation's evaluation ahead of the projection.
func HumanOutcomeReadinessStopEffects(conn *sql.DB, envelopeID string, revision int, actorID string, repository Repository, overrides map[string]string) ([]journal.Effect, []Stop, []Attention, error) {
	envelope, err := loadEnvelope(conn, envelopeID)
	if err != nil {
		return nil, nil, nil, err
	}
	if envelope.CurrentRevision != revision || envelope.Status != "active" {
		return nil, nil, nil, nil
	}
	return humanOutcomeReadinessEffectsForEnvelope(conn, envelope, actorID, repository, overrides)
}

// ArtifactAvailabilityOverrides returns the Obligation-status overrides that
// apply once the named Artifact becomes available: an Obligation whose
// current evaluation's sole unsatisfied reason is exactly artifact-unavailable
// for this Artifact is forced to satisfied, because artifact.verify's own
// referenced_artifacts availability update has not yet reached the projection
// within the same committing Event.
func ArtifactAvailabilityOverrides(conn *sql.DB, envelopeID string, revision int, artifactID string, repository Repository) (map[string]string, error) {
	envelope, err := loadEnvelope(conn, envelopeID)
	if err != nil {
		return nil, err
	}
	if envelope.CurrentRevision != revision || envelope.Status != "active" {
		return nil, nil
	}
	evaluations, err := currentEnvelopeEvaluations(conn, envelope, repository)
	if err != nil {
		return nil, err
	}
	overrides := map[string]string{}
	for _, evaluation := range evaluations {
		if len(evaluation.Reasons) == 1 && evaluation.Reasons[0].Code == "artifact-unavailable" && evaluation.Reasons[0].SubjectID == artifactID {
			overrides[evaluation.ObligationID] = "satisfied"
		}
	}
	return overrides, nil
}
