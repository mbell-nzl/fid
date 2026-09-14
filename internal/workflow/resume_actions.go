package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func resumeNextActions(conn *sql.DB, envelope Envelope, terminal bool, evaluation ObligationEvaluation, claim *ResumeClaim, run *Run, repository Repository) ([]NextAction, error) {
	if envelope.Status == "completed" {
		return []NextAction{}, nil
	}
	actions := []NextAction{}
	if run != nil {
		actions = append(actions, NextAction{Operation: "run.inspect", Subject: Subject{Type: "run", ID: run.ID}, RequiredCapabilities: []string{}, Reasons: []EvaluationReason{}})
		observe, err := resumeRunObserveAction(conn, envelope, run)
		if err != nil {
			return nil, err
		}
		actions = append(actions, observe)
	}
	action, found, err := resumeGateAction(conn, envelope, run, repository)
	if err != nil {
		return nil, err
	}
	if found {
		actions = append(actions, action)
	}
	if claim != nil {
		renew, err := resumeClaimRenewAction(conn, envelope, claim)
		if err != nil {
			return nil, err
		}
		actions = append(actions, renew)
		action, err := resumeClaimReleaseAction(conn, envelope, claim)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func resumeRunObserveAction(conn *sql.DB, envelope Envelope, run *Run) (NextAction, error) {
	action := NextAction{Operation: "run.observe", Subject: Subject{Type: "run", ID: run.ID}, RequiredCapabilities: []string{"run.observe"}, Reasons: []EvaluationReason{}}
	if run.Status != "active" {
		action.Reasons = append(action.Reasons, evaluationReason("run-inactive", run.Status, run.ID, "run"))
	}
	authorized, err := resumeActionAuthority(conn, envelope, run.ActorID, "run.observe")
	if err != nil {
		return NextAction{}, err
	}
	if !authorized {
		action.Reasons = append(action.Reasons, evaluationReason("authority-unavailable", "run.observe", run.ActorID, "actor"))
	}
	return withSortedReasons(action), nil
}

func resumeClaimRenewAction(conn *sql.DB, envelope Envelope, claim *ResumeClaim) (NextAction, error) {
	action := NextAction{Operation: "claim.renew", Subject: Subject{Type: "claim", ID: claim.ID}, RequiredCapabilities: []string{"claim.renew"}, Reasons: []EvaluationReason{}}
	if claim.EffectiveStatus == "expired" {
		action.Reasons = append(action.Reasons, evaluationReason("claim-expired", claim.ID, claim.ID, "claim"))
	}
	authorized, err := resumeActionAuthority(conn, envelope, claim.ActorID, "claim.renew")
	if err != nil {
		return NextAction{}, err
	}
	if !authorized {
		action.Reasons = append(action.Reasons, evaluationReason("authority-unavailable", "claim.renew", claim.ActorID, "actor"))
	}
	return withSortedReasons(action), nil
}

func resumeClaimReleaseAction(conn *sql.DB, envelope Envelope, claim *ResumeClaim) (NextAction, error) {
	action := NextAction{Operation: "claim.release", Subject: Subject{Type: "claim", ID: claim.ID}, RequiredCapabilities: []string{"claim.release"}, Reasons: []EvaluationReason{}}
	if claim.EffectiveStatus == "expired" {
		action.Reasons = append(action.Reasons, evaluationReason("claim-expired", claim.ID, claim.ID, "claim"))
	}
	authorized, err := resumeActionAuthority(conn, envelope, claim.ActorID, "claim.release")
	if err != nil {
		return NextAction{}, err
	}
	if !authorized {
		action.Reasons = append(action.Reasons, evaluationReason("authority-unavailable", "claim.release", claim.ActorID, "actor"))
	}
	return withSortedReasons(action), nil
}

func resumeGateAction(conn *sql.DB, envelope Envelope, run *Run, repository Repository) (NextAction, bool, error) {
	if run == nil || !sameRepository(run.StartRepository, repository) {
		return NextAction{}, false, nil
	}
	action := NextAction{Operation: "gate.run", Subject: Subject{Type: "envelope", ID: envelope.ID}, RequiredCapabilities: []string{"gate.run"}, Reasons: []EvaluationReason{}}
	action.Subject = Subject{Type: "run", ID: run.ID}
	if run.Status != "active" {
		action.Reasons = append(action.Reasons, evaluationReason("run-inactive", run.Status, run.ID, "run"))
	}
	for _, block := range BudgetAdmissionBlocks(run.Budgets) {
		action.Reasons = append(action.Reasons, evaluationReason(block.Code, block.Budget.Name, envelope.ID, "envelope"))
	}
	subjects := []Subject{{Type: "run", ID: run.ID}}
	activationIDs, err := resumeOpenActivationIDs(conn)
	if err != nil {
		return NextAction{}, false, err
	}
	for _, activationID := range activationIDs {
		subjects = append(subjects, Subject{Type: "activation", ID: activationID})
	}
	stops, err := ApplicableStops(conn, "gate.run", subjects...)
	if err != nil {
		return NextAction{}, false, fmt.Errorf("resume gate Stops: %w", err)
	}
	for _, stop := range stops {
		action.Reasons = append(action.Reasons, evaluationReason(stop.Code, "gate.run", stop.ID, "stop"))
	}
	authorized, err := resumeActionAuthority(conn, envelope, run.ActorID, "gate.run")
	if err != nil {
		return NextAction{}, false, err
	}
	if !authorized {
		action.Reasons = append(action.Reasons, evaluationReason("authority-unavailable", "gate.run", run.ActorID, "actor"))
	}
	return withSortedReasons(action), true, nil
}

func resumeActionAuthority(conn *sql.DB, envelope Envelope, actorID, action string) (bool, error) {
	var incomplete int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM kit_activations WHERE status IN ('prepared', 'staged', 'committed')`).Scan(&incomplete); err != nil {
		return false, fmt.Errorf("resume gate incomplete activation: %w", err)
	}
	if incomplete != 0 {
		return false, nil
	}
	var activationID string
	if err := conn.QueryRow(`SELECT id FROM kit_activations WHERE status = 'finalized' LIMIT 1`).Scan(&activationID); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("resume gate activation: %w", err)
	}
	var producer, capabilitiesJSON string
	if err := conn.QueryRow(`SELECT producer, capabilities_json FROM activation_actor_profiles WHERE activation_id = ? AND actor_id = ?`, activationID, actorID).Scan(&producer, &capabilitiesJSON); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("resume gate actor: %w", err)
	}
	var capabilities []string
	if err := json.Unmarshal([]byte(capabilitiesJSON), &capabilities); err != nil {
		return false, fmt.Errorf("resume gate capabilities: %w", err)
	}
	if !containsResumeAction(capabilities, action) {
		return false, nil
	}
	if action == "run.observe" && producer != "activated-harness" {
		return false, nil
	}
	rows, err := conn.Query(`SELECT actor_selector_json, actions_json FROM envelope_authority_grants WHERE envelope_id = ? AND revision = ? ORDER BY ordinal`, envelope.ID, envelope.CurrentRevision)
	if err != nil {
		return false, fmt.Errorf("resume gate grants: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var selectorJSON, actionsJSON string
		var selector struct {
			ActorID              string   `json:"actor_id"`
			RequiredCapabilities []string `json:"required_capabilities"`
		}
		var actions []string
		if err := rows.Scan(&selectorJSON, &actionsJSON); err != nil {
			return false, err
		}
		if err := json.Unmarshal([]byte(selectorJSON), &selector); err != nil {
			return false, err
		}
		if err := json.Unmarshal([]byte(actionsJSON), &actions); err != nil {
			return false, err
		}
		if selector.ActorID == actorID && containsResumeAction(actions, action) && containsResumeAction(selector.RequiredCapabilities, action) {
			return true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	return false, nil
}

func containsResumeAction(actions []string, wanted string) bool {
	for _, action := range actions {
		if action == wanted {
			return true
		}
	}
	return false
}

func withSortedReasons(action NextAction) NextAction {
	action.Reasons = sortEvaluationReasons(action.Reasons)
	return action
}
