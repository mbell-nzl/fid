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

// ObligationEvaluationInput carries durable Evidence presence into the shared
// action-independent Obligation evaluator.
type ObligationEvaluationInput struct {
	Action      string
	Obligations []Obligation
	Evidence    map[string]bool
}

// EvaluateObligations returns the same ordered result for each guarded action.
// Later Evidence work supplies the remaining Receipt, validity, freshness, and
// Artifact dimensions through this one boundary.
func EvaluateObligations(input ObligationEvaluationInput) []EvaluationReason {
	reasons := make([]EvaluationReason, 0, len(input.Obligations))
	for _, obligation := range input.Obligations {
		if !input.Evidence[obligation.ID] {
			reasons = append(reasons, evaluationReason("obligation-evidence-absent", "qualifying Evidence is absent", obligation.ID, "obligation"))
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
		obligationReasons, err := evaluateObligations(conn, envelope.Obligations)
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
		FROM work_claims
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

func evaluateObligations(conn *sql.DB, obligations []Obligation) ([]EvaluationReason, error) {
	evidence := make(map[string]bool, len(obligations))
	for _, obligation := range obligations {
		var count int
		if err := conn.QueryRow(`SELECT COUNT(*) FROM evidence_records WHERE obligation_id = ?`, obligation.ID).Scan(&count); err != nil {
			return nil, fmt.Errorf("evaluate obligation %q: %w", obligation.ID, err)
		}
		evidence[obligation.ID] = count > 0
	}
	return EvaluateObligations(ObligationEvaluationInput{Action: "issue.transition-terminal-success", Obligations: obligations, Evidence: evidence}), nil
}

func evaluationReason(code, detail, subjectID, subjectType string) EvaluationReason {
	return EvaluationReason{Code: code, Detail: detail, SubjectID: subjectID, SubjectType: subjectType}
}

func sortEvaluationReasons(reasons []EvaluationReason) []EvaluationReason {
	ordered := append([]EvaluationReason(nil), reasons...)
	precedence := map[string]int{
		"authority-unavailable": 1, "envelope-inactive": 2, "claim-missing": 3, "claim-expired": 4,
		"run-inactive": 5, "issue-revision-conflict": 6, "workflow-transition-invalid": 7,
		"obligation-evidence-absent": 8, "receipt-unsuccessful": 9, "evidence-invalid": 10,
		"evidence-stale": 11, "artifact-unavailable": 12, "issue-not-terminal-success": 13,
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
