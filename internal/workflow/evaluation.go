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

// StopBlockedError reports an unresolved Stop that applies to a guarded
// operation. It carries the immutable facts so command adapters do not infer
// scope from an error string.
type StopBlockedError struct {
	Action string
	Stops  []Stop
}

func (err *StopBlockedError) Error() string { return "guarded action is blocked by unresolved Stop" }

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
	ObligationID         string
	EvidenceID           string
	ReceiptID            string
	ReceiptSuccessful    bool
	Validity             string
	Freshness            string
	ArtifactAvailability string
	ArtifactIDs          []string
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
	for _, state := range states {
		if state.EvidenceID == "" && !state.ReceiptSuccessful {
			reasons = append(reasons, evaluationReason("obligation-evidence-absent", "qualifying Evidence is absent", state.ObligationID, "obligation"))
			continue
		}
		if !state.ReceiptSuccessful {
			reasons = append(reasons, evaluationReason("receipt-unsuccessful", "Receipt outcome is not succeeded", state.ReceiptID, "receipt"))
		}
		if state.Validity != "valid" {
			reasons = append(reasons, evaluationReason("evidence-invalid", "Evidence validity is invalid", state.EvidenceID, "evidence"))
		}
		if state.Freshness != "fresh" {
			reasons = append(reasons, evaluationReason("evidence-stale", "repository tree differs from Evidence target", state.EvidenceID, "evidence"))
		}
		if state.ArtifactAvailability != "available" {
			for _, artifactID := range state.ArtifactIDs {
				reasons = append(reasons, evaluationReason("artifact-unavailable", "required Artifact is unavailable", artifactID, "artifact"))
			}
		}
	}
	return sortEvaluationReasons(reasons)
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
	return json.Unmarshal([]byte(repositoryJSON), &repository) == nil && sameRepository(repository, envelope.RepositoryScope)
}

func evaluateObligations(conn *sql.DB, obligations []Obligation, current Repository) ([]EvaluationReason, error) {
	states := make([]ObligationEvidenceState, 0, len(obligations))
	for _, obligation := range obligations {
		loaded, err := loadObligationEvidenceStates(conn, obligation.ID, current)
		if err != nil {
			return nil, err
		}
		states = append(states, loaded...)
	}
	return EvaluateEvidenceStates(states), nil
}

func loadObligationEvidenceStates(conn *sql.DB, obligationID string, current Repository) ([]ObligationEvidenceState, error) {
	rows, err := conn.Query(`
		SELECT e.id, e.receipt_id, e.validity, e.target_repository_json, r.terminal_status
		FROM evidence_records e
		JOIN execution_receipts r ON r.id = e.receipt_id
		WHERE e.obligation_id = ?
		ORDER BY e.id
	`, obligationID)
	if err != nil {
		return nil, fmt.Errorf("evaluate obligation %q: %w", obligationID, err)
	}
	defer func() { _ = rows.Close() }()

	states := []ObligationEvidenceState{}
	for rows.Next() {
		var evidenceID, receiptID, validity, targetJSON, outcome string
		if err := rows.Scan(&evidenceID, &receiptID, &validity, &targetJSON, &outcome); err != nil {
			return nil, err
		}
		var target Repository
		if err := json.Unmarshal([]byte(targetJSON), &target); err != nil {
			return nil, fmt.Errorf("decode Evidence target repository: %w", err)
		}
		state := ObligationEvidenceState{
			ObligationID: obligationID, EvidenceID: evidenceID, ReceiptID: receiptID,
			ReceiptSuccessful: outcome == "succeeded", Validity: validity,
			Freshness: evidenceFreshness(target, current), ArtifactAvailability: "available",
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

func evidenceFreshness(target, current Repository) string {
	if sameRepository(target, current) {
		return "fresh"
	}
	return "stale"
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
		"obligation-evidence-absent": 8, "receipt-unsuccessful": 9, "evidence-invalid": 10,
		"evidence-stale": 11, "artifact-unavailable": 12, "issue-not-terminal-success": 13,
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
