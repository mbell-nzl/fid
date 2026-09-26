package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

// EvaluationReason is a stable, ordered failure reason for a guarded action.
type EvaluationReason struct {
	Code        string `json:"code"`
	Detail      string `json:"detail"`
	SubjectID   string `json:"subject_id"`
	SubjectType string `json:"subject_type"`
}

type ExecutionTransitionRequest struct {
	IssueID               string
	ActorID               string
	ExpectedIssueRevision *int
	ActualIssueRevision   int
	WorkflowEdgeValid     bool
	TerminalSuccess       bool
}

type ExecutionTransitionEvaluation struct {
	Bound      bool
	EnvelopeID string
	Reasons    []EvaluationReason
}

// ApplicableStops returns unresolved Stops whose exact subject and declared
// action match this mutation. Callers use the same predicate before any
// durable effect is prepared.
func ApplicableStops(conn *sql.DB, action string, subjects ...Subject) ([]Stop, error) {
	if conn == nil || !blockedActionSet[action] {
		return nil, fmt.Errorf("Stop predicate input is invalid")
	}
	items := []Stop{}
	for _, subject := range subjects {
		if !validSubject(subject) {
			return nil, fmt.Errorf("Stop predicate subject is invalid")
		}
		rows, err := conn.Query(`
			SELECT s.id
			FROM stops s
			LEFT JOIN stop_resolutions r ON r.stop_id = s.id
			WHERE s.subject_type = ? AND s.subject_id = ? AND r.stop_id IS NULL
			  AND EXISTS (SELECT 1 FROM json_each(s.blocked_actions_json) WHERE value = ?)
			ORDER BY s.raised_at,
			  CASE WHEN s.code = 'protected-decision-required' THEN COALESCE((
			    SELECT p.ordinal FROM envelope_protected_decisions p
			    WHERE p.envelope_id = s.subject_id
			      AND p.revision = (SELECT current_revision FROM work_envelopes WHERE id = s.subject_id)
			      AND p.decision_class = json_extract(s.details_json, '$.decision_class')
			  ), 2147483647) ELSE 0 END,
			  s.id
		`, subject.Type, subject.ID, action)
		if err != nil {
			return nil, fmt.Errorf("evaluate Stops: %w", err)
		}
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				_ = rows.Close()
				return nil, err
			}
			stop, err := loadStop(conn, id)
			if err != nil {
				_ = rows.Close()
				return nil, err
			}
			items = append(items, stop)
		}
		if err := rows.Close(); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func stopEvaluationReasons(conn *sql.DB, action string, subjects ...Subject) ([]EvaluationReason, error) {
	stops, err := ApplicableStops(conn, action, subjects...)
	if err != nil {
		return nil, err
	}
	reasons := make([]EvaluationReason, 0, len(stops))
	for _, stop := range stops {
		reasons = append(reasons, evaluationReason(stop.Code, action, stop.ID, "stop"))
	}
	return reasons, nil
}

// ObligationEvaluationInput carries durable Evidence presence into the shared
// action-independent Obligation evaluator.
type ObligationEvaluationInput struct {
	Action      string
	Obligations []Obligation
	Evidence    map[string]bool
}

// ObligationEvidenceState carries the independent Evidence dimensions used by
// pure evaluation and guarded transitions.
type ObligationEvidenceState struct {
	ObligationID                 string
	EvidenceID                   string
	EvidenceType                 string
	ReceiptID                    string
	ReceiptSuccessful            bool
	ReviewID                     string
	ReviewVerdict                string
	ReviewIndependent            bool
	HumanOutcomeID               string
	HumanOutcomeResult           string
	OutcomeSuperseded            bool
	UnresolvedBlockingFindingIDs []string
	Validity                     string
	Freshness                    string
	ArtifactAvailability         string
	ArtifactIDs                  []string
}

// EvaluateObligations returns the same ordered result for each guarded action.
// Later Evidence work supplies the remaining Receipt, validity, freshness, and
// Artifact dimensions through this one boundary.
func EvaluateObligations(input ObligationEvaluationInput) []EvaluationReason {
	states := make([]ObligationEvidenceState, 0, len(input.Obligations))
	for _, obligation := range input.Obligations {
		states = append(states, ObligationEvidenceState{ObligationID: obligation.ID, Validity: "valid", Freshness: "fresh", ArtifactAvailability: "available", ReceiptSuccessful: input.Evidence[obligation.ID]})
	}
	return EvaluateEvidenceStates(states)
}

// EvaluateEvidenceStates returns every independent Evidence failure in the
// accepted reason order. A missing Evidence record is distinct from an
// unsuccessful Receipt and from each persisted Evidence dimension.
func EvaluateEvidenceStates(states []ObligationEvidenceState) []EvaluationReason {
	reasons := make([]EvaluationReason, 0, len(states)*3)
	seen := map[EvaluationReason]bool{}
	appendReason := func(reason EvaluationReason) {
		if !seen[reason] {
			reasons = append(reasons, reason)
			seen[reason] = true
		}
	}
	for _, state := range states {
		if state.EvidenceID == "" && !state.ReceiptSuccessful {
			appendReason(evaluationReason("obligation-evidence-absent", "qualifying Evidence is absent", state.ObligationID, "obligation"))
			continue
		}
		if state.EvidenceType == "human-outcome" && state.OutcomeSuperseded {
			appendReason(evaluationReason("outcome-superseded", "a later Human Outcome has a different result", state.HumanOutcomeID, "human-outcome"))
		}
		if state.EvidenceType != "review" && state.EvidenceType != "human-outcome" && !state.ReceiptSuccessful {
			appendReason(evaluationReason("receipt-unsuccessful", "Receipt outcome is not succeeded", state.ReceiptID, "receipt"))
		}
		for _, findingID := range state.UnresolvedBlockingFindingIDs {
			appendReason(evaluationReason("finding-unresolved", "blocking Finding remains unresolved", findingID, "finding"))
		}
		if state.EvidenceType == "review" && !state.ReviewIndependent {
			appendReason(evaluationReason("review-not-independent", "Review actor is not independent", state.ReviewID, "review"))
		}
		if state.Validity != "valid" {
			appendReason(evaluationReason("evidence-invalid", "Evidence validity is invalid", state.EvidenceID, "evidence"))
		}
		if state.Freshness != "fresh" {
			appendReason(evaluationReason("evidence-stale", "repository tree differs from Evidence target", state.EvidenceID, "evidence"))
		}
		if state.ArtifactAvailability != "available" {
			for _, artifactID := range state.ArtifactIDs {
				appendReason(evaluationReason("artifact-unavailable", "required Artifact is unavailable", artifactID, "artifact"))
			}
		}
	}
	return sortEvaluationReasons(reasons)
}

// EvaluateEnvelopeObligations serializes the independently evaluated state of
// every Obligation in stable ID order. Metadata-only Obligations deliberately
// have no Artifact availability requirement.
func EvaluateEnvelopeObligations(obligations []Obligation, statesByObligation map[string][]ObligationEvidenceState) []ObligationEvaluation {
	ordered := append([]Obligation(nil), obligations...)
	sort.Slice(ordered, func(left, right int) bool { return ordered[left].ID < ordered[right].ID })
	evaluations := make([]ObligationEvaluation, 0, len(ordered))
	for _, obligation := range ordered {
		states := append([]ObligationEvidenceState(nil), statesByObligation[obligation.ID]...)
		if len(states) == 0 {
			states = []ObligationEvidenceState{{ObligationID: obligation.ID}}
		}
		metadataSufficient := obligation.ArtifactPolicy == "metadata-sufficient"
		if metadataSufficient {
			for index := range states {
				states[index].ArtifactAvailability = "not-applicable"
			}
		}
		reasons := EvaluateEvidenceStates(states)
		if obligation.Type == "human-outcome" && evidenceAbsent(states) {
			reasons = []EvaluationReason{evaluationReason("human-outcome-absent", "qualifying Human Outcome is absent", obligation.ID, "obligation")}
		}
		evaluation := ObligationEvaluation{ObligationID: obligation.ID, Status: "unsatisfied", Validity: "invalid", Freshness: "stale", ArtifactAvailability: "missing", Reasons: reasons}
		if obligation.FreshnessTarget == "none" {
			evaluation.Freshness = "fresh"
		}
		if metadataSufficient {
			evaluation.ArtifactAvailability = "not-applicable"
		}
		for _, state := range states {
			if state.Validity == "valid" {
				evaluation.Validity = "valid"
			}
			if state.Freshness == "fresh" {
				evaluation.Freshness = "fresh"
			}
			if !metadataSufficient && state.ArtifactAvailability == "available" {
				evaluation.ArtifactAvailability = "available"
			}
		}
		if len(reasons) == 0 {
			evaluation.Status, evaluation.Validity, evaluation.Freshness = "satisfied", "valid", "fresh"
			if !metadataSufficient {
				evaluation.ArtifactAvailability = "available"
			}
		}
		evaluations = append(evaluations, evaluation)
	}
	return evaluations
}

func evaluateProtectedDecisions(conn *sql.DB, envelope Envelope, action string) ([]EvaluationReason, error) {
	reasons := []EvaluationReason{}
	for index := range envelope.ProtectedDecisions {
		decision := &envelope.ProtectedDecisions[index]
		if err := resolveProtectedDecision(conn, envelope.PrimaryIssueID, decision); err != nil {
			return nil, err
		}
		if decision.Status == "satisfied" || !containsAction(decision.AffectedActions, action) {
			continue
		}
		reasons = append(reasons, evaluationReason(
			"protected-decision-required",
			"protected Decision requires approved reference and Human Outcome",
			decision.DecisionClass,
			"decision",
		))
	}
	return sortEvaluationReasons(reasons), nil
}

// ProtectedDecisionsForIssue derives the current status of every protected
// Decision on the Issue's active or draft Envelope revision.
func ProtectedDecisionsForIssue(conn *sql.DB, issueID string) (Envelope, []ProtectedDecision, error) {
	return ProtectedDecisionsForIssueWithFact(conn, issueID, nil)
}

// ProtectedDecisionsForIssueWithFact is ProtectedDecisionsForIssue with the
// judgment fact of a not-yet-committed Event applied, so a producer reports the
// statuses its own Event will leave behind (DEC-250 R1).
func ProtectedDecisionsForIssueWithFact(conn *sql.DB, issueID string, pending *JudgmentFact) (Envelope, []ProtectedDecision, error) {
	var envelopeID string
	err := conn.QueryRow(`
		SELECT e.id
		FROM work_envelopes e
		JOIN envelope_issue_selection s ON s.envelope_id=e.id AND s.revision=e.current_revision
		WHERE s.issue_id=? AND e.status IN ('active','draft')
		ORDER BY CASE e.status WHEN 'active' THEN 0 ELSE 1 END,e.id
		LIMIT 1`, issueID).Scan(&envelopeID)
	if err != nil {
		return Envelope{}, nil, fmt.Errorf("load protected Decision Envelope: %w", err)
	}
	envelope, err := loadEnvelope(conn, envelopeID)
	if err != nil {
		return Envelope{}, nil, err
	}
	for index := range envelope.ProtectedDecisions {
		if _, err := deriveProtectedDecision(conn, envelope.PrimaryIssueID, &envelope.ProtectedDecisions[index], pending); err != nil {
			return Envelope{}, nil, err
		}
	}
	return envelope, envelope.ProtectedDecisions, nil
}

// ProtectedDecisionBlockForIssue reloads the immutable Stop that produced an
// execution-bound Issue denial. Command adapters use it to render the same
// protected Decision facts as other guarded actions.
func ProtectedDecisionBlockForIssue(workspaceRoot, migrationsDir, issueID, action string) (*StopBlockedError, error) {
	return stopBlockForIssue(workspaceRoot, migrationsDir, issueID, action, "protected-decision-required")
}

// HumanOutcomeRequiredBlockForIssue reloads the immutable Stop that produced
// an execution-bound Issue denial for the human-outcome-required precedence
// case, mirroring ProtectedDecisionBlockForIssue.
func HumanOutcomeRequiredBlockForIssue(workspaceRoot, migrationsDir, issueID, action string) (*StopBlockedError, error) {
	return stopBlockForIssue(workspaceRoot, migrationsDir, issueID, action, "human-outcome-required")
}

func stopBlockForIssue(workspaceRoot, migrationsDir, issueID, action, code string) (*StopBlockedError, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return nil, err
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return nil, err
	}
	defer closeWorkflowDB(conn)

	envelope, _, err := ProtectedDecisionsForIssue(conn, issueID)
	if err != nil {
		return nil, err
	}
	stops, err := ApplicableStops(conn, action, Subject{Type: "envelope", ID: envelope.ID})
	if err != nil {
		return nil, err
	}
	for _, stop := range stops {
		if stop.Code == code {
			return NewStopBlockedError(conn, action, []Stop{stop})
		}
	}
	return nil, fmt.Errorf("Stop %q is unavailable for issue %q", code, issueID)
}

func resolveProtectedDecision(conn *sql.DB, issueID string, decision *ProtectedDecision) error {
	_, err := deriveProtectedDecision(conn, issueID, decision, nil)
	return err
}

// protectedDecisionBinding records which Events created the Decision Reference
// and Human Outcome bound to a protected Decision, so a Stop predicate can name
// the pair-completing Event. A half that comes from the pending fact has no
// committed Event yet and is by construction the later one.
type protectedDecisionBinding struct {
	RefEventID, OutcomeEventID string
	RefPending, OutcomePending bool
}

// deriveProtectedDecision is the one section 9 derivation, used by every
// guarded action and every Decision-satisfied Stop predicate. pending is the
// judgment fact of an Event that has not committed yet: it competes as the
// newest fact of its kind, exactly as it will once committed.
func deriveProtectedDecision(conn *sql.DB, issueID string, decision *ProtectedDecision, pending *JudgmentFact) (protectedDecisionBinding, error) {
	required, err := journal.MarshalCanonicalJSON(decision.RequiredCapabilities)
	if err != nil {
		return protectedDecisionBinding{}, fmt.Errorf("encode protected Decision authority: %w", err)
	}
	var binding protectedDecisionBinding
	var referenceID, outcomeID string
	err = conn.QueryRow(`
		SELECT id, created_event_id FROM decision_refs
		WHERE class = ? AND scope_issue_id = ? AND outcome = 'approved'
		  AND authority_capabilities_json = ?
		ORDER BY created_event_id DESC LIMIT 1`,
		decision.DecisionClass, issueID, string(required)).Scan(&referenceID, &binding.RefEventID)
	if err != nil && err != sql.ErrNoRows {
		return protectedDecisionBinding{}, fmt.Errorf("resolve protected Decision reference: %w", err)
	}
	referenceFound := err == nil
	outcomeErr := conn.QueryRow(`
		SELECT id, created_event_id FROM human_outcomes
		WHERE class = ? AND scope_issue_id = ? AND result = 'approved'
		ORDER BY created_event_id DESC LIMIT 1`, decision.DecisionClass, issueID).Scan(&outcomeID, &binding.OutcomeEventID)
	if outcomeErr != nil && outcomeErr != sql.ErrNoRows {
		return protectedDecisionBinding{}, fmt.Errorf("resolve protected Decision Human Outcome: %w", outcomeErr)
	}
	outcomeFound := outcomeErr == nil
	if pending != nil && pending.approved() && pending.Class == decision.DecisionClass && pending.ScopeIssueID == issueID {
		switch {
		case pending.EventType == "decision.referenced" && slices.Equal(pending.RequiredCapabilities, decision.RequiredCapabilities):
			referenceID, referenceFound, binding.RefEventID, binding.RefPending = pending.ID, true, "", true
		case pending.EventType == "human-outcome.recorded":
			outcomeID, outcomeFound, binding.OutcomeEventID, binding.OutcomePending = pending.ID, true, "", true
		}
	}
	decision.DecisionRefID = nil
	decision.HumanOutcomeID = nil
	decision.Status = "unsatisfied"
	if referenceFound {
		decision.DecisionRefID = &referenceID
	}
	if outcomeFound {
		decision.HumanOutcomeID = &outcomeID
	}
	if referenceFound && outcomeFound {
		decision.Status = "satisfied"
	}
	return binding, nil
}

func evidenceAbsent(states []ObligationEvidenceState) bool {
	for _, state := range states {
		if state.EvidenceID != "" || state.ReceiptSuccessful {
			return false
		}
	}
	return true
}

// EvaluateExecutionTransition evaluates the generic Envelope-bound path. A
// non-bound Issue keeps using its retained standalone workflow behavior.
func EvaluateExecutionTransition(workspaceRoot, migrationsDir string, request ExecutionTransitionRequest) (ExecutionTransitionEvaluation, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return ExecutionTransitionEvaluation{}, err
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return ExecutionTransitionEvaluation{}, err
	}
	defer closeWorkflowDB(conn)

	return evaluateExecutionTransition(conn, workspaceRoot, migrationsDir, request)
}

func evaluateExecutionTransition(conn *sql.DB, workspaceRoot, migrationsDir string, request ExecutionTransitionRequest) (ExecutionTransitionEvaluation, error) {
	var envelopeID string
	err := conn.QueryRow(`
		SELECT e.id
		FROM work_envelopes e
		JOIN envelope_issue_selection s ON s.envelope_id = e.id AND s.revision = e.current_revision
		WHERE s.issue_id = ?
	`, request.IssueID).Scan(&envelopeID)
	if err == sql.ErrNoRows {
		return ExecutionTransitionEvaluation{}, nil
	}
	if err != nil {
		return ExecutionTransitionEvaluation{}, fmt.Errorf("evaluate execution transition envelope: %w", err)
	}

	envelope, err := loadEnvelope(conn, envelopeID)
	if err != nil {
		return ExecutionTransitionEvaluation{}, err
	}
	reasons := []EvaluationReason{}
	if err := AuthorizeEnvelopeGrant(workspaceRoot, migrationsDir, EnvelopeAuthorityRequest{ActorID: request.ActorID, Action: "issue.transition", Capability: "issue.transition", IssueID: request.IssueID, EnvelopeID: envelope.ID, Revision: envelope.CurrentRevision}); err != nil {
		reasons = append(reasons, evaluationReason("authority-unavailable", "issue.transition", request.ActorID, "actor"))
	}
	if envelope.Status != "active" {
		reasons = append(reasons, evaluationReason("envelope-inactive", envelope.Status, envelope.ID, "envelope"))
	}

	claim, claimFound, claimExpired, err := loadExecutionClaim(conn, request.IssueID, envelope)
	if err != nil {
		return ExecutionTransitionEvaluation{}, err
	}
	if !claimFound || claim.ActorID != request.ActorID {
		reasons = append(reasons, evaluationReason("claim-missing", "active mutation claim", request.IssueID, "issue"))
	} else if claimExpired {
		reasons = append(reasons, evaluationReason("claim-expired", claim.ID, claim.ID, "claim"))
	}
	if !executionRunActive(conn, claim, claimFound && !claimExpired && claim.ActorID == request.ActorID, envelope) {
		reasons = append(reasons, evaluationReason("run-inactive", "active run", request.IssueID, "issue"))
	}
	if request.ExpectedIssueRevision == nil || *request.ExpectedIssueRevision != request.ActualIssueRevision {
		reasons = append(reasons, evaluationReason("issue-revision-conflict", "expected current Issue revision", request.IssueID, "issue"))
	}
	if !request.WorkflowEdgeValid {
		reasons = append(reasons, evaluationReason("workflow-transition-invalid", "workflow edge", request.IssueID, "issue"))
	}
	if request.TerminalSuccess {
		repository, err := resumeObservedRepository(workspaceRoot, workspaceRoot)
		if err != nil {
			return ExecutionTransitionEvaluation{}, fmt.Errorf("observe terminal transition repository: %w", err)
		}
		stopReasons, err := stopEvaluationReasons(conn, "issue.transition-terminal-success", Subject{Type: "envelope", ID: envelope.ID})
		if err != nil {
			return ExecutionTransitionEvaluation{}, err
		}
		reasons = append(reasons, stopReasons...)
		obligationReasons, err := evaluateObligations(conn, obligationsGuarding(envelope.Obligations, "issue.transition-terminal-success"), repository)
		if err != nil {
			return ExecutionTransitionEvaluation{}, err
		}
		reasons = append(reasons, obligationReasons...)
	}
	return ExecutionTransitionEvaluation{Bound: true, EnvelopeID: envelope.ID, Reasons: sortEvaluationReasons(reasons)}, nil
}

type executionClaim struct {
	ID        string
	ActorID   string
	ExpiresAt string
}

func loadExecutionClaim(conn *sql.DB, issueID string, envelope Envelope) (executionClaim, bool, bool, error) {
	var claim executionClaim
	err := conn.QueryRow(`
		SELECT id, actor_id, expires_at
		FROM claims
		WHERE issue_id = ? AND envelope_id = ? AND envelope_revision = ? AND worktree_id = ? AND status = 'active'
		ORDER BY id
		LIMIT 1
	`, issueID, envelope.ID, envelope.CurrentRevision, envelope.RepositoryScope.WorktreeID).Scan(&claim.ID, &claim.ActorID, &claim.ExpiresAt)
	if err == sql.ErrNoRows {
		return executionClaim{}, false, false, nil
	}
	if err != nil {
		return executionClaim{}, false, false, fmt.Errorf("evaluate execution transition claim: %w", err)
	}
	expiresAt, err := time.Parse(time.RFC3339Nano, claim.ExpiresAt)
	if err != nil {
		return executionClaim{}, false, false, fmt.Errorf("evaluate execution transition claim expiry: %w", err)
	}
	return claim, true, !time.Now().UTC().Before(expiresAt.UTC()), nil
}

func executionRunActive(conn *sql.DB, claim executionClaim, usableClaim bool, envelope Envelope) bool {
	if !usableClaim {
		return false
	}
	var repositoryJSON string
	err := conn.QueryRow(`
		SELECT start_repository_json
		FROM runs
		WHERE claim_id = ? AND envelope_id = ? AND envelope_revision = ? AND status = 'active'
	`, claim.ID, envelope.ID, envelope.CurrentRevision).Scan(&repositoryJSON)
	if err != nil {
		return false
	}
	var repository Repository
	// A recovered Run has a newly observed repository target but remains bound
	// to the Envelope's approved project and worktree. Comparing its full target
	// to the original Envelope snapshot would reject the very reconciliation
	// that admitted the successor Run; content may change, identity may not.
	return json.Unmarshal([]byte(repositoryJSON), &repository) == nil &&
		repository.ProjectID == envelope.RepositoryScope.ProjectID &&
		repository.WorktreeID == envelope.RepositoryScope.WorktreeID
}

// obligationsGuarding returns the Obligations applicable to a guarded action: a
// required Obligation whose guarded_actions contains that action (DEC-249). An
// Envelope in which none guards the action puts no Obligation requirement on it.
// Reads never use this filter; they report every Obligation.
func obligationsGuarding(obligations []Obligation, action string) []Obligation {
	guarding := make([]Obligation, 0, len(obligations))
	for _, obligation := range obligations {
		if obligation.Applicability == "required" && containsAction(obligation.GuardedActions, action) {
			guarding = append(guarding, obligation)
		}
	}
	return guarding
}

func evaluateObligations(conn *sql.DB, obligations []Obligation, current Repository) ([]EvaluationReason, error) {
	states := map[string][]ObligationEvidenceState{}
	for _, obligation := range obligations {
		loaded, err := loadObligationEvidenceStates(conn, obligation.ID, current)
		if err != nil {
			return nil, err
		}
		states[obligation.ID] = loaded
	}
	evaluations := EvaluateEnvelopeObligations(obligations, states)
	reasons := []EvaluationReason{}
	for _, evaluation := range evaluations {
		reasons = append(reasons, evaluation.Reasons...)
	}
	return sortEvaluationReasons(reasons), nil
}

func loadObligationEvidenceStates(conn *sql.DB, obligationID string, current Repository) ([]ObligationEvidenceState, error) {
	rows, err := conn.Query(`
		SELECT e.id, e.evidence_type, COALESCE(e.receipt_id,''), COALESCE(e.review_id,''), COALESCE(e.human_outcome_id,''), e.validity, e.target_repository_json,
		       COALESCE(r.terminal_status,''), COALESCE(rv.verdict,''),
		       COALESCE(rv.issue_id,''), COALESCE(rv.envelope_id,''), COALESCE(rv.envelope_revision,0),
		       COALESCE(ho.result,''), COALESCE(ho.class,''), COALESCE(ho.scope_issue_id,'')
		FROM evidence_records e
		LEFT JOIN execution_receipts r ON r.id = e.receipt_id
		LEFT JOIN review_records rv ON rv.id = e.review_id
		LEFT JOIN human_outcomes ho ON ho.id = e.human_outcome_id
		WHERE e.obligation_id = ?
		ORDER BY e.id
	`, obligationID)
	if err != nil {
		return nil, fmt.Errorf("evaluate obligation %q: %w", obligationID, err)
	}
	defer func() { _ = rows.Close() }()

	states := []ObligationEvidenceState{}
	for rows.Next() {
		var evidenceID, evidenceType, receiptID, reviewID, humanOutcomeID, validity, targetJSON, outcome, verdict, issueID, envelopeID string
		var humanOutcomeResult, outcomeClass, scopeIssueID string
		var revision int
		if err := rows.Scan(&evidenceID, &evidenceType, &receiptID, &reviewID, &humanOutcomeID, &validity, &targetJSON, &outcome, &verdict, &issueID, &envelopeID, &revision, &humanOutcomeResult, &outcomeClass, &scopeIssueID); err != nil {
			return nil, err
		}
		var target Repository
		if err := json.Unmarshal([]byte(targetJSON), &target); err != nil {
			return nil, fmt.Errorf("decode Evidence target repository: %w", err)
		}
		state := ObligationEvidenceState{
			ObligationID: obligationID, EvidenceID: evidenceID, EvidenceType: evidenceType, ReceiptID: receiptID, ReviewID: reviewID, ReviewVerdict: verdict,
			HumanOutcomeID: humanOutcomeID, HumanOutcomeResult: humanOutcomeResult,
			ReceiptSuccessful: outcome == "succeeded", Validity: validity,
			Freshness: evidenceFreshnessForType(evidenceType, target, current), ArtifactAvailability: "available",
		}
		if evidenceType == "human-outcome" {
			governing, err := ResolveGoverningHumanOutcome(conn, outcomeClass, scopeIssueID)
			if err != nil {
				return nil, err
			}
			state.OutcomeSuperseded = HumanOutcomeSuperseded(governing, humanOutcomeID, humanOutcomeResult)
		}
		if evidenceType == "review" {
			state.ReviewIndependent = true
			state.UnresolvedBlockingFindingIDs, err = unresolvedBlockingFindingIDs(conn, issueID, envelopeID, revision)
			if err != nil {
				return nil, err
			}
		}
		artifactRows, err := conn.Query(`
			SELECT l.artifact_id, a.availability
			FROM evidence_artifact_links l
			JOIN referenced_artifacts a ON a.id = l.artifact_id
			WHERE l.evidence_id = ? AND l.required_material = 1
			ORDER BY l.artifact_id
		`, evidenceID)
		if err != nil {
			return nil, fmt.Errorf("evaluate obligation artifacts: %w", err)
		}
		for artifactRows.Next() {
			var artifactID, availability string
			if err := artifactRows.Scan(&artifactID, &availability); err != nil {
				_ = artifactRows.Close()
				return nil, err
			}
			state.ArtifactIDs = append(state.ArtifactIDs, artifactID)
			if availability != "available" {
				state.ArtifactAvailability = "missing"
			}
		}
		if err := artifactRows.Close(); err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(states) == 0 {
		states = append(states, ObligationEvidenceState{ObligationID: obligationID})
	}
	return states, nil
}

// ReviewIndependenceConflict identifies the first mutating-side history that
// conflicts with a Review actor. Claims take precedence over Runs.
type ReviewIndependenceConflict struct {
	ReviewActorID string
	SubjectType   string
	SubjectID     string
	ActorID       string
}

// FindReviewIndependenceConflict applies the shared durable-history rule used
// both when a Review is recorded and when it is later linked as Evidence.
func FindReviewIndependenceConflict(conn *sql.DB, actorID, issueID, envelopeID string, revision int) (*ReviewIndependenceConflict, error) {
	var conflictingID string
	err := conn.QueryRow(`SELECT id FROM claims WHERE issue_id=? AND envelope_id=? AND envelope_revision=? AND mode='mutation' AND actor_id=? ORDER BY acquired_at,id LIMIT 1`, issueID, envelopeID, revision, actorID).Scan(&conflictingID)
	if err == nil {
		return &ReviewIndependenceConflict{ReviewActorID: actorID, SubjectType: "claim", SubjectID: conflictingID, ActorID: actorID}, nil
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("load mutating Claim history: %w", err)
	}
	err = conn.QueryRow(`SELECT r.id FROM runs r JOIN claims c ON c.id=r.claim_id WHERE r.issue_id=? AND r.envelope_id=? AND r.envelope_revision=? AND c.mode='mutation' AND r.actor_id=? ORDER BY r.started_at,r.id LIMIT 1`, issueID, envelopeID, revision, actorID).Scan(&conflictingID)
	if err == nil {
		return &ReviewIndependenceConflict{ReviewActorID: actorID, SubjectType: "run", SubjectID: conflictingID, ActorID: actorID}, nil
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("load mutating Run history: %w", err)
	}
	return nil, nil
}

func unresolvedBlockingFindingIDs(conn *sql.DB, issueID, envelopeID string, revision int) ([]string, error) {
	rows, err := conn.Query(`SELECT f.id FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE r.issue_id=? AND r.envelope_id=? AND r.envelope_revision=? AND f.blocking=1 AND NOT EXISTS (SELECT 1 FROM finding_dispositions d WHERE d.finding_id=f.id AND d.action IN ('resolved','inapplicable')) ORDER BY f.id`, issueID, envelopeID, revision)
	if err != nil {
		return nil, fmt.Errorf("load unresolved blocking Findings: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func evidenceFreshness(target, current Repository) string {
	if sameRepository(target, current) {
		return "fresh"
	}
	return "stale"
}

func evidenceFreshnessForType(evidenceType string, target, current Repository) string {
	if evidenceType == "human-outcome" {
		return "fresh"
	}
	return evidenceFreshness(target, current)
}

func evaluationReason(code, detail, subjectID, subjectType string) EvaluationReason {
	return EvaluationReason{Code: code, Detail: detail, SubjectID: subjectID, SubjectType: subjectType}
}

func sortEvaluationReasons(reasons []EvaluationReason) []EvaluationReason {
	ordered := make([]EvaluationReason, len(reasons))
	copy(ordered, reasons)
	precedence := map[string]int{
		"authority-unavailable": 1, "envelope-inactive": 2, "claim-missing": 3, "claim-expired": 4,
		"run-inactive": 5, "issue-revision-conflict": 6, "workflow-transition-invalid": 7,
		"protected-decision-required": 8,
		"human-outcome-required":      9,
		"obligation-evidence-absent":  10, "human-outcome-absent": 11, "finding-unresolved": 12,
		"review-not-independent": 13, "outcome-superseded": 14, "receipt-unsuccessful": 15,
		"evidence-invalid": 16, "evidence-stale": 17, "artifact-unavailable": 18, "issue-not-terminal-success": 19,
		"executor-liveness-unknown": 7, "recovery-state-conflict": 7, "budget-exhausted": 7,
		"budget-data-unavailable": 7, "activation-incomplete": 7, "capability-unavailable": 7, "operator-stop": 7,
	}
	sort.Slice(ordered, func(left, right int) bool {
		leftRank, rightRank := precedence[ordered[left].Code], precedence[ordered[right].Code]
		if leftRank != rightRank {
			return leftRank < rightRank
		}
		if ordered[left].SubjectType != ordered[right].SubjectType {
			return ordered[left].SubjectType < ordered[right].SubjectType
		}
		if ordered[left].SubjectID != ordered[right].SubjectID {
			return ordered[left].SubjectID < ordered[right].SubjectID
		}
		return ordered[left].Detail < ordered[right].Detail
	})
	return ordered
}
