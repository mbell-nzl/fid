package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func resumeNextActions(conn *sql.DB, issueID string, envelope Envelope, terminal bool, evaluation ObligationEvaluation, claim *ResumeClaim, run *Run, repository Repository) ([]NextAction, error) {
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
	judgment, err := resumeJudgmentActions(conn, issueID, envelope, claim, run, repository)
	if err != nil {
		return nil, err
	}
	actions = append(actions, judgment...)
	return actions, nil
}

func resumeJudgmentActions(conn *sql.DB, issueID string, envelope Envelope, claim *ResumeClaim, run *Run, repository Repository) ([]NextAction, error) {
	actions := []NextAction{}
	for _, operation := range []string{"concern.record", "deviation.record", "decision.record-reference", "human-outcome.record"} {
		action := NextAction{Operation: operation, Subject: Subject{Type: "issue", ID: issueID}, RequiredCapabilities: []string{operation}, Reasons: []EvaluationReason{}}
		authorized, err := resumeAnyActorAuthority(conn, envelope, operation)
		if err != nil {
			return nil, err
		}
		if !authorized {
			action.Reasons = append(action.Reasons, evaluationReason("authority-unavailable", operation, issueID, "issue"))
		}
		actions = append(actions, withSortedReasons(action))
	}
	if run != nil && run.Status == "active" && claim != nil && claim.Mode == "review" {
		action := NextAction{Operation: "review.record", Subject: Subject{Type: "run", ID: run.ID}, RequiredCapabilities: []string{"review.record"}, Reasons: []EvaluationReason{}}
		if !sameRepository(run.StartRepository, repository) {
			action.Reasons = append(action.Reasons, evaluationReason("evidence-stale", "review repository differs", run.ID, "run"))
		}
		if claim.EffectiveStatus != "active" {
			action.Reasons = append(action.Reasons, evaluationReason("claim-expired", claim.ID, claim.ID, "claim"))
		}
		var count int
		if err := conn.QueryRow(`SELECT COUNT(*) FROM review_records WHERE run_id=?`, run.ID).Scan(&count); err != nil {
			return nil, err
		}
		if count != 0 {
			action.Reasons = append(action.Reasons, evaluationReason("review-not-independent", "Run already reviewed", run.ID, "run"))
		}
		conflict, err := FindReviewIndependenceConflict(conn, run.ActorID, issueID, envelope.ID, envelope.CurrentRevision)
		if err != nil {
			return nil, err
		}
		if conflict != nil {
			action.Reasons = append(action.Reasons, evaluationReason("review-not-independent", conflict.SubjectType, conflict.SubjectID, conflict.SubjectType))
		}
		authorized, err := resumeActionAuthority(conn, envelope, run.ActorID, "review.record")
		if err != nil {
			return nil, err
		}
		if !authorized {
			action.Reasons = append(action.Reasons, evaluationReason("authority-unavailable", "review.record", run.ActorID, "actor"))
		}
		actions = append(actions, withSortedReasons(action))
	}
	if claim != nil && claim.Mode == "mutation" && claim.EffectiveStatus == "active" {
		address, err := resumeFindingAddressActions(conn, issueID, envelope, claim)
		if err != nil {
			return nil, err
		}
		actions = append(actions, address...)
	}
	dispose, err := resumeConcernDisposeActions(conn, issueID, envelope)
	if err != nil {
		return nil, err
	}
	actions = append(actions, dispose...)
	return actions, nil
}

func resumeFindingAddressActions(conn *sql.DB, issueID string, envelope Envelope, claim *ResumeClaim) ([]NextAction, error) {
	rows, err := conn.Query(`SELECT f.id FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE r.issue_id=?
		AND NOT EXISTS(SELECT 1 FROM finding_dispositions d WHERE d.finding_id=f.id) ORDER BY f.id LIMIT ?`, issueID, resumeJudgmentLimit+1)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	actions := []NextAction{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		action := NextAction{Operation: "finding.address", Subject: Subject{Type: "finding", ID: id}, RequiredCapabilities: []string{"finding.address"}, Reasons: []EvaluationReason{}}
		authorized, err := resumeActionAuthority(conn, envelope, claim.ActorID, "finding.address")
		if err != nil {
			return nil, err
		}
		if !authorized {
			action.Reasons = append(action.Reasons, evaluationReason("authority-unavailable", "finding.address", claim.ActorID, "actor"))
		}
		actions = append(actions, withSortedReasons(action))
	}
	return actions, rows.Err()
}

func resumeConcernDisposeActions(conn *sql.DB, issueID string, envelope Envelope) ([]NextAction, error) {
	rows, err := conn.Query(`SELECT c.id FROM concerns c WHERE c.status='open' AND (`+resumeJudgmentSubjectFilter("c")+`) ORDER BY c.id LIMIT ?`,
		issueID, issueID, issueID, issueID, issueID, issueID, resumeJudgmentLimit+1)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	actions := []NextAction{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		action := NextAction{Operation: "concern.dispose", Subject: Subject{Type: "concern", ID: id}, RequiredCapabilities: []string{"concern.dispose"}, Reasons: []EvaluationReason{}}
		authorized, err := resumeAnyActorAuthority(conn, envelope, "concern.dispose")
		if err != nil {
			return nil, err
		}
		if !authorized {
			action.Reasons = append(action.Reasons, evaluationReason("authority-unavailable", "concern.dispose", id, "concern"))
		}
		actions = append(actions, withSortedReasons(action))
	}
	return actions, rows.Err()
}

func resumeAnyActorAuthority(conn *sql.DB, envelope Envelope, operation string) (bool, error) {
	var incomplete int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM kit_activations WHERE status IN ('prepared','staged','committed')`).Scan(&incomplete); err != nil {
		return false, err
	}
	if incomplete != 0 {
		return false, nil
	}
	rows, err := conn.Query(`SELECT p.actor_id,p.actor_type,p.capabilities_json FROM activation_actor_profiles p JOIN kit_activations a ON a.id=p.activation_id WHERE a.status='finalized' ORDER BY p.actor_id`)
	if err != nil {
		return false, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var actorID, actorType, capabilitiesJSON string
		if err := rows.Scan(&actorID, &actorType, &capabilitiesJSON); err != nil {
			return false, err
		}
		var capabilities []string
		if err := json.Unmarshal([]byte(capabilitiesJSON), &capabilities); err != nil {
			return false, err
		}
		if !containsResumeAction(capabilities, operation) {
			continue
		}
		if operation == "human-outcome.record" {
			if actorType == "human" {
				return true, nil
			}
			continue
		}
		allowed, err := resumeActionAuthority(conn, envelope, actorID, operation)
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}
	return false, rows.Err()
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
