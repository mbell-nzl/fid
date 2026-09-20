package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
			ORDER BY s.raised_at, s.id
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
		obligationReasons, err := evaluateObligations(conn, envelope.Obligations, repository)
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
			state.OutcomeSuperseded = governing.ID != humanOutcomeID && governing.Result != humanOutcomeResult
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
		"obligation-evidence-absent": 8, "human-outcome-absent": 9, "finding-unresolved": 10,
		"review-not-independent": 11, "outcome-superseded": 12, "receipt-unsuccessful": 13,
		"evidence-invalid": 14, "evidence-stale": 15, "artifact-unavailable": 16, "issue-not-terminal-success": 17,
		"executor-liveness-unknown": 8, "recovery-state-conflict": 8, "budget-exhausted": 8,
		"budget-data-unavailable": 8, "activation-incomplete": 8, "capability-unavailable": 8, "operator-stop": 8,
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
