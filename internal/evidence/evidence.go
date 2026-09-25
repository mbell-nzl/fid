package evidence

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

// Evidence is one immutable typed link between an Obligation and its subject.
type Evidence struct {
	ID                  string     `json:"id"`
	ObligationID        string     `json:"obligation_id"`
	EvidenceType        string     `json:"evidence_type"`
	ReceiptID           *string    `json:"receipt_id"`
	ReviewID            *string    `json:"review_id"`
	HumanOutcomeID      *string    `json:"human_outcome_id"`
	Producer            string     `json:"producer"`
	TrustSource         string     `json:"trust_source"`
	TargetRepository    Repository `json:"target_repository"`
	Validity            string     `json:"validity"`
	RequiredArtifactIDs []string   `json:"required_artifact_ids"`
	CreatedAt           string     `json:"created_at"`
}

// Evaluation is the public, independently serialized Evidence result.
type Evaluation struct {
	ObligationID         string                      `json:"obligation_id"`
	Status               string                      `json:"status"`
	Validity             string                      `json:"validity"`
	Freshness            string                      `json:"freshness"`
	ArtifactAvailability string                      `json:"artifact_availability"`
	Reasons              []workflow.EvaluationReason `json:"reasons"`
}

type EvidenceLinkRequest struct {
	ObligationID        string     `json:"obligation_id"`
	EvidenceType        string     `json:"evidence_type"`
	ReceiptID           *string    `json:"receipt_id"`
	ReviewID            *string    `json:"review_id"`
	HumanOutcomeID      *string    `json:"human_outcome_id"`
	TargetRepository    Repository `json:"target_repository"`
	RequiredArtifactIDs []string   `json:"required_artifact_ids"`
}

type EvidenceEvaluateRequest struct {
	ObligationID     *string    `json:"obligation_id"`
	TargetRepository Repository `json:"target_repository"`
}

type EvidenceLinkOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
	Request        EvidenceLinkRequest
}

type EvidenceLinkResult struct {
	Evidence        Evidence   `json:"evidence"`
	Evaluation      Evaluation `json:"evaluation"`
	IdempotencyKey  string     `json:"idempotency_key"`
	Replayed        bool       `json:"replayed"`
	EventID         string
	EventOccurredAt string
}

type EvidenceNotFoundError struct{ ID string }

func (err *EvidenceNotFoundError) Error() string {
	return fmt.Sprintf("Evidence %q was not found", err.ID)
}

type EvidenceBindingError struct{ Reason string }

func (err *EvidenceBindingError) Error() string { return err.Reason }

// LinkEvidence persists a link only when the Receipt was observed by fid and
// its complete durable chain matches the selected Obligation.
func LinkEvidence(options EvidenceLinkOptions) (EvidenceLinkResult, error) {
	if err := validateEvidenceLinkOptions(options); err != nil {
		return EvidenceLinkResult{}, err
	}
	var result EvidenceLinkResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid evidence link"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openEvidenceDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
		if err != nil {
			return fmt.Errorf("link Evidence: read project config: %w", err)
		}
		binding, err := loadEvidenceBinding(conn, options.Request)
		if err != nil {
			return err
		}
		if err := workflow.AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, workflow.EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: "evidence.link", Capability: "evidence.link", IssueID: binding.IssueID, EnvelopeID: binding.EnvelopeID, Revision: binding.Revision}); err != nil {
			return err
		}
		_, actor, err := loadActiveActor(conn, options.ActorID)
		if err != nil {
			return err
		}
		if !hasCapability(actor.Capabilities, "evidence.link") {
			return &AuthorityDeniedError{ActorID: options.ActorID}
		}
		requestDigest, err := evidenceRequestDigest(project.ProjectID, options.ActorID, options.Request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "evidence.link:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			if err := json.Unmarshal([]byte(outcome.ResultJSON), &result); err != nil {
				return fmt.Errorf("link Evidence: decode replay: %w", err)
			}
			result.IdempotencyKey, result.Replayed, result.EventID = options.IdempotencyKey, true, outcome.EventID
			return nil
		}

		evidenceID, err := nextEvidenceID(conn)
		if err != nil {
			return err
		}
		createdAt := time.Now().UTC().Format(time.RFC3339Nano)
		evidence := Evidence{
			ID:                  evidenceID,
			ObligationID:        options.Request.ObligationID,
			EvidenceType:        binding.EvidenceType,
			ReceiptID:           binding.ReceiptID,
			ReviewID:            binding.ReviewID,
			HumanOutcomeID:      binding.HumanOutcomeID,
			Producer:            binding.Producer,
			TrustSource:         binding.TrustSource,
			TargetRepository:    options.Request.TargetRepository,
			Validity:            "valid",
			RequiredArtifactIDs: append([]string{}, options.Request.RequiredArtifactIDs...),
			CreatedAt:           createdAt,
		}
		evaluation := evaluationFromBinding(evidence, binding, options.Request.TargetRepository)
		result = EvidenceLinkResult{Evidence: evidence, Evaluation: evaluation, IdempotencyKey: options.IdempotencyKey}
		resultJSON, err := journal.MarshalCanonicalJSON(result)
		if err != nil {
			return err
		}
		targetJSON, err := journal.MarshalCanonicalJSON(evidence.TargetRepository)
		if err != nil {
			return err
		}
		effects := []journal.Effect{{Op: "insert", Table: "evidence_records", PK: map[string]any{"id": evidence.ID}, Values: map[string]any{"obligation_id": evidence.ObligationID, "receipt_id": evidence.ReceiptID, "review_id": evidence.ReviewID, "human_outcome_id": evidence.HumanOutcomeID, "evidence_type": evidence.EvidenceType, "producer": evidence.Producer, "trust_source": evidence.TrustSource, "target_repository_json": string(targetJSON), "validity": evidence.Validity, "created_at": evidence.CreatedAt}}}
		for _, artifactID := range evidence.RequiredArtifactIDs {
			effects = append(effects, journal.Effect{Op: "insert", Table: "evidence_artifact_links", PK: map[string]any{"evidence_id": evidence.ID, "artifact_id": artifactID}, Values: map[string]any{"required_material": true}})
		}
		readinessEffects, _, _, err := workflow.HumanOutcomeReadinessStopEffects(conn, binding.EnvelopeID, binding.Revision, actor.ActorID, workflowRepository(options.Request.TargetRepository), map[string]string{options.Request.ObligationID: evaluation.Status})
		if err != nil {
			return err
		}
		effects = append(effects, readinessEffects...)
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "evidence.linked", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: evidence.ID, Sequence: 1, Type: "evidence"}, CommandName: "evidence.link", CorrelationID: binding.CorrelationID, Repository: repositoryEventValue(evidence.TargetRepository), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			substituteEvidenceEventPlaceholders(prepared.Effects, prepared)
			return journal.IdempotencyMutation{EntityID: evidence.ID, ResultJSON: string(resultJSON)}, nil
		})
		if err != nil {
			return err
		}
		result.EventID = outcome.EventID
		return nil
	})
	if err != nil {
		return EvidenceLinkResult{}, err
	}
	occurredAt, err := gateEventOccurredAt(options.WorkspaceRoot, result.EventID)
	if err != nil {
		return EvidenceLinkResult{}, err
	}
	result.EventOccurredAt = occurredAt
	return result, nil
}

// EvaluateEvidence reads durable Evidence dimensions without an Event or Effect.
func EvaluateEvidence(workspaceRoot, migrationsDir string, request EvidenceEvaluateRequest) ([]Evaluation, error) {
	if strings.TrimSpace(workspaceRoot) == "" || strings.TrimSpace(migrationsDir) == "" || (request.ObligationID != nil && strings.TrimSpace(*request.ObligationID) == "") || !validEvidenceRepository(request.TargetRepository) {
		return nil, &EvidenceBindingError{Reason: "Evidence evaluation input is invalid"}
	}
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return nil, err
	}
	conn, err := openEvidenceDB(workspaceRoot)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	return evaluateEvidence(conn, request)
}

func GetEvidence(workspaceRoot, migrationsDir, id string) (Evidence, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Evidence{}, err
	}
	conn, err := openEvidenceDB(workspaceRoot)
	if err != nil {
		return Evidence{}, err
	}
	defer func() { _ = conn.Close() }()
	evidence, err := loadEvidence(conn, id)
	if err == sql.ErrNoRows {
		return Evidence{}, &EvidenceNotFoundError{ID: id}
	}
	return evidence, err
}

type evidenceBinding struct {
	EnvelopeID                   string
	Revision                     int
	IssueID                      string
	CorrelationID                *string
	EvidenceType                 string
	ReceiptID                    *string
	ReviewID                     *string
	HumanOutcomeID               *string
	Producer                     string
	TrustSource                  string
	ReceiptOutcome               string
	ReviewVerdict                string
	ReviewIndependent            bool
	HumanOutcomeClass            string
	UnresolvedBlockingFindingIDs []string
	Repository                   Repository
	Artifacts                    map[string]string
	Obligation                   workflow.ObligationInput
}

func validateEvidenceLinkOptions(options EvidenceLinkOptions) error {
	request := options.Request
	subjectCount := nonNilCount(request.ReceiptID, request.ReviewID, request.HumanOutcomeID)
	subjectMatchesType := (request.EvidenceType == "receipt" && nonBlankPointer(request.ReceiptID)) || (request.EvidenceType == "review" && nonBlankPointer(request.ReviewID)) || (request.EvidenceType == "human-outcome" && nonBlankPointer(request.HumanOutcomeID))
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.IdempotencyKey) == "" || strings.TrimSpace(request.ObligationID) == "" || subjectCount != 1 || !subjectMatchesType || !validEvidenceRepository(request.TargetRepository) || !sortedUnique(request.RequiredArtifactIDs) {
		return &EvidenceBindingError{Reason: "Evidence link input is invalid"}
	}
	return nil
}

func nonNilCount(values ...*string) int {
	count := 0
	for _, value := range values {
		if value != nil {
			count++
		}
	}
	return count
}

func nonBlankPointer(value *string) bool {
	return value != nil && strings.TrimSpace(*value) != ""
}

func loadEvidenceBinding(conn *sql.DB, request EvidenceLinkRequest) (evidenceBinding, error) {
	binding, err := loadObligationBinding(conn, request.ObligationID)
	if err != nil {
		return evidenceBinding{}, err
	}
	switch request.EvidenceType {
	case "receipt":
		err = loadReceiptEvidenceBinding(conn, request, &binding)
	case "review":
		err = loadReviewEvidenceBinding(conn, request, &binding)
	case "human-outcome":
		err = loadHumanOutcomeEvidenceBinding(conn, request, &binding)
	default:
		return evidenceBinding{}, &EvidenceBindingError{Reason: "Evidence link input is invalid"}
	}
	if err != nil {
		return evidenceBinding{}, err
	}
	return binding, nil
}

func loadObligationBinding(conn *sql.DB, obligationID string) (evidenceBinding, error) {
	var binding evidenceBinding
	var obligationJSON string
	err := conn.QueryRow(`SELECT envelope_id,revision,obligation_json FROM envelope_obligations WHERE id=?`, obligationID).Scan(&binding.EnvelopeID, &binding.Revision, &obligationJSON)
	if err == sql.ErrNoRows {
		return binding, &EvidenceBindingError{Reason: "Evidence Obligation was not found"}
	}
	if err != nil {
		return binding, fmt.Errorf("load Evidence Obligation: %w", err)
	}
	if err := json.Unmarshal([]byte(obligationJSON), &binding.Obligation); err != nil {
		return binding, fmt.Errorf("decode Evidence Obligation: %w", err)
	}
	binding.IssueID = binding.Obligation.Subject.ID
	return binding, nil
}

func loadReceiptEvidenceBinding(conn *sql.DB, request EvidenceLinkRequest, binding *evidenceBinding) error {
	var repositoryJSON, trustSource, producer, registrationID, attemptObligationID, runEnvelopeID, runID, issueID string
	var registrationVersion, runRevision int
	err := conn.QueryRow(`SELECT r.terminal_status,r.terminal_repository_json,r.trust_source,a.producer,a.registration_id,a.registration_version,a.obligation_id,a.run_id,run.issue_id,run.envelope_id,run.envelope_revision FROM execution_receipts r JOIN execution_attempts a ON a.id=r.attempt_id JOIN runs run ON run.id=a.run_id WHERE r.id=?`, *request.ReceiptID).Scan(&binding.ReceiptOutcome, &repositoryJSON, &trustSource, &producer, &registrationID, &registrationVersion, &attemptObligationID, &runID, &issueID, &runEnvelopeID, &runRevision)
	if err == sql.ErrNoRows {
		return &EvidenceBindingError{Reason: "Receipt and Obligation binding is not exact"}
	}
	if err != nil {
		return fmt.Errorf("load Receipt chain: %w", err)
	}
	if err := json.Unmarshal([]byte(repositoryJSON), &binding.Repository); err != nil {
		return fmt.Errorf("decode Receipt repository: %w", err)
	}
	obligation := binding.Obligation
	if trustSource != "fid-executed" || producer != "fid" || binding.EnvelopeID != runEnvelopeID || binding.Revision != runRevision || binding.IssueID != issueID || !sameEvidenceRepository(binding.Repository, request.TargetRepository) {
		return &EvidenceBindingError{Reason: "Receipt and Obligation binding is not exact"}
	}
	binding.EvidenceType, binding.ReceiptID = "receipt", request.ReceiptID
	binding.Producer, binding.TrustSource, binding.CorrelationID = "fid", "fid-executed", &runID
	if err := classifyBindingAdmission(request, binding); err != nil {
		return err
	}
	if registrationID != obligation.Registration.ID || registrationVersion != obligation.Registration.Version || attemptObligationID != request.ObligationID {
		return &EvidenceBindingError{Reason: "Receipt and Obligation binding is not exact"}
	}
	binding.Artifacts = map[string]string{}
	rows, err := conn.Query(`SELECT ra.artifact_id,a.availability FROM receipt_artifacts ra JOIN referenced_artifacts a ON a.id=ra.artifact_id WHERE ra.receipt_id=? ORDER BY ra.artifact_id`, *request.ReceiptID)
	if err != nil {
		return fmt.Errorf("load Receipt Artifacts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id, availability string
		if err := rows.Scan(&id, &availability); err != nil {
			return err
		}
		binding.Artifacts[id] = availability
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, artifactID := range request.RequiredArtifactIDs {
		if _, found := binding.Artifacts[artifactID]; !found {
			return &EvidenceBindingError{Reason: "required Artifact is not linked to Receipt"}
		}
	}
	return nil
}

func loadReviewEvidenceBinding(conn *sql.DB, request EvidenceLinkRequest, binding *evidenceBinding) error {
	var runID, issueID, envelopeID, actorID, profile, repositoryJSON string
	var revision int
	err := conn.QueryRow(`SELECT run_id,issue_id,envelope_id,envelope_revision,actor_id,profile,repository_json,verdict FROM review_records WHERE id=?`, *request.ReviewID).Scan(&runID, &issueID, &envelopeID, &revision, &actorID, &profile, &repositoryJSON, &binding.ReviewVerdict)
	if err == sql.ErrNoRows {
		return &ReviewNotFoundError{ID: *request.ReviewID}
	}
	if err != nil {
		return fmt.Errorf("load Review Evidence binding: %w", err)
	}
	if err := json.Unmarshal([]byte(repositoryJSON), &binding.Repository); err != nil {
		return fmt.Errorf("decode Review repository: %w", err)
	}
	if issueID != binding.IssueID || envelopeID != binding.EnvelopeID || revision != binding.Revision || !sameEvidenceRepository(binding.Repository, request.TargetRepository) {
		return &EvidenceBindingError{Reason: "Review and Obligation binding is not exact"}
	}
	blocking, err := blockingFindingIDs(conn, `f.review_id=?`, *request.ReviewID)
	if err != nil {
		return err
	}
	if err := validateReviewEvidenceAdmission(request.ObligationID, *request.ReviewID, binding.ReviewVerdict, blocking); err != nil {
		return err
	}
	conflict, err := workflow.FindReviewIndependenceConflict(conn, actorID, issueID, envelopeID, revision)
	if err != nil {
		return err
	}
	if conflict != nil {
		return &ReviewNotIndependentError{ReviewActorID: actorID, ConflictingType: conflict.SubjectType, ConflictingID: conflict.SubjectID, ConflictingActorID: conflict.ActorID}
	}
	binding.EvidenceType, binding.ReviewID = "review", request.ReviewID
	binding.Producer, binding.TrustSource, binding.CorrelationID = "activated-actor", "activated-actor", &runID
	if err := classifyBindingAdmission(request, binding); err != nil {
		return err
	}
	if binding.Obligation.ReviewProfile == nil || profile != *binding.Obligation.ReviewProfile {
		return &EvidenceBindingError{Reason: "Review and Obligation binding is not exact"}
	}
	unresolved, err := blockingFindingIDs(conn, `r.issue_id=? AND r.envelope_id=? AND r.envelope_revision=? AND NOT EXISTS (SELECT 1 FROM finding_dispositions d WHERE d.finding_id=f.id AND d.action IN ('resolved','inapplicable'))`, issueID, envelopeID, revision)
	if err != nil {
		return err
	}
	binding.ReviewIndependent = true
	binding.UnresolvedBlockingFindingIDs = unresolved
	binding.Artifacts = map[string]string{}
	return nil
}

func loadHumanOutcomeEvidenceBinding(conn *sql.DB, request EvidenceLinkRequest, binding *evidenceBinding) error {
	var class, subjectType, subjectID, scopeIssueID string
	err := conn.QueryRow(`SELECT class,subject_type,subject_id,result,scope_issue_id FROM human_outcomes WHERE id=?`, *request.HumanOutcomeID).Scan(&class, &subjectType, &subjectID, &binding.ReceiptOutcome, &scopeIssueID)
	if err == sql.ErrNoRows {
		return &OutcomeNotFoundError{ID: *request.HumanOutcomeID}
	}
	if err != nil {
		return fmt.Errorf("load Human Outcome Evidence binding: %w", err)
	}
	if subjectType != "issue" || subjectID != binding.IssueID || scopeIssueID != binding.IssueID {
		return &EvidenceBindingError{Reason: "Human Outcome and Obligation binding is not exact"}
	}
	binding.EvidenceType, binding.HumanOutcomeID = "human-outcome", request.HumanOutcomeID
	binding.Producer, binding.TrustSource = "human", "activated-human"
	binding.Repository, binding.Artifacts, binding.HumanOutcomeClass = request.TargetRepository, map[string]string{}, class
	return classifyBindingAdmission(request, binding)
}

func classifyBindingAdmission(request EvidenceLinkRequest, binding *evidenceBinding) error {
	return classifyEvidenceAdmission(evidenceAdmission{
		ObligationID:           request.ObligationID,
		ObligationType:         binding.Obligation.Type,
		AcceptedEvidenceTypes:  binding.Obligation.AcceptedEvidenceTypes,
		AcceptedTrustSources:   binding.Obligation.AcceptedTrustSources,
		EvidenceType:           binding.EvidenceType,
		TrustSource:            binding.TrustSource,
		HumanOutcomeID:         pointerValue(binding.HumanOutcomeID),
		HumanOutcomeResult:     binding.ReceiptOutcome,
		HumanOutcomeClass:      binding.HumanOutcomeClass,
		ObligationOutcomeClass: pointerValue(binding.Obligation.OutcomeClass),
	})
}

func blockingFindingIDs(conn *sql.DB, where string, arguments ...any) ([]string, error) {
	rows, err := conn.Query(`SELECT f.id FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE f.blocking=1 AND `+where+` ORDER BY f.id`, arguments...)
	if err != nil {
		return nil, fmt.Errorf("load blocking Findings: %w", err)
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

func pointerValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func evaluateEvidence(conn *sql.DB, request EvidenceEvaluateRequest) ([]Evaluation, error) {
	obligations, err := loadEvaluationObligations(conn, request.ObligationID)
	if err != nil {
		return nil, err
	}
	statesByObligation := make(map[string][]workflow.ObligationEvidenceState, len(obligations))
	for _, obligation := range obligations {
		states, err := loadEvaluationStates(conn, obligation, request.TargetRepository)
		if err != nil {
			return nil, err
		}
		statesByObligation[obligation.ID] = states
	}
	shared := workflow.EvaluateEnvelopeObligations(obligations, statesByObligation)
	evaluations := make([]Evaluation, 0, len(shared))
	for _, evaluation := range shared {
		evaluations = append(evaluations, evaluationFromWorkflow(evaluation))
	}
	return evaluations, nil
}

func evaluationFromBinding(evidence Evidence, binding evidenceBinding, target Repository) Evaluation {
	availability := "available"
	artifactIDs := append([]string(nil), evidence.RequiredArtifactIDs...)
	for _, id := range artifactIDs {
		if binding.Artifacts[id] != "available" {
			availability = "missing"
		}
	}
	state := workflow.ObligationEvidenceState{
		ObligationID:                 evidence.ObligationID,
		EvidenceID:                   evidence.ID,
		EvidenceType:                 evidence.EvidenceType,
		ReceiptID:                    pointerValue(evidence.ReceiptID),
		ReviewID:                     pointerValue(evidence.ReviewID),
		ReceiptSuccessful:            binding.ReceiptOutcome == "succeeded",
		ReviewVerdict:                binding.ReviewVerdict,
		ReviewIndependent:            binding.ReviewIndependent,
		UnresolvedBlockingFindingIDs: append([]string(nil), binding.UnresolvedBlockingFindingIDs...),
		Validity:                     evidence.Validity,
		Freshness:                    freshnessForTarget(binding.Obligation.FreshnessTarget, binding.Repository, target),
		ArtifactAvailability:         availability,
		ArtifactIDs:                  artifactIDs,
	}
	obligation := workflow.Obligation{ObligationInput: binding.Obligation}
	result := workflow.EvaluateEnvelopeObligations([]workflow.Obligation{obligation}, map[string][]workflow.ObligationEvidenceState{evidence.ObligationID: {state}})
	return evaluationFromWorkflow(result[0])
}

func loadEvaluationObligations(conn *sql.DB, obligationID *string) ([]workflow.Obligation, error) {
	query := `SELECT o.obligation_json FROM envelope_obligations o JOIN work_envelopes e ON e.id=o.envelope_id AND e.current_revision=o.revision WHERE e.status='active' ORDER BY o.id`
	arguments := []any{}
	if obligationID != nil {
		query = `SELECT obligation_json FROM envelope_obligations WHERE id=? ORDER BY id`
		arguments = append(arguments, *obligationID)
	}
	rows, err := conn.Query(query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("load Evidence evaluation Obligations: %w", err)
	}
	defer func() { _ = rows.Close() }()
	obligations := []workflow.Obligation{}
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return nil, err
		}
		var input workflow.ObligationInput
		if err := json.Unmarshal([]byte(payload), &input); err != nil {
			return nil, fmt.Errorf("decode Evidence evaluation Obligation: %w", err)
		}
		obligations = append(obligations, workflow.Obligation{ObligationInput: input})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if obligationID != nil && len(obligations) == 0 {
		return nil, &EvidenceBindingError{Reason: "Evidence Obligation was not found"}
	}
	return obligations, nil
}

func loadEvaluationStates(conn *sql.DB, obligation workflow.Obligation, target Repository) ([]workflow.ObligationEvidenceState, error) {
	var envelopeID string
	var revision int
	if err := conn.QueryRow(`SELECT envelope_id,revision FROM envelope_obligations WHERE id=?`, obligation.ID).Scan(&envelopeID, &revision); err != nil {
		return nil, fmt.Errorf("load Evidence evaluation Obligation binding: %w", err)
	}
	rows, err := conn.Query(`
		SELECT e.id,e.evidence_type,e.receipt_id,e.review_id,e.human_outcome_id,e.validity,e.target_repository_json,
		       r.terminal_status,rv.verdict,ho.result,ho.class,ho.scope_issue_id
		FROM evidence_records e
		LEFT JOIN execution_receipts r ON r.id=e.receipt_id
		LEFT JOIN review_records rv ON rv.id=e.review_id
		LEFT JOIN human_outcomes ho ON ho.id=e.human_outcome_id
		WHERE e.obligation_id=?
		ORDER BY e.id`, obligation.ID)
	if err != nil {
		return nil, fmt.Errorf("evaluate Evidence: %w", err)
	}
	defer func() { _ = rows.Close() }()
	states := []workflow.ObligationEvidenceState{}
	for rows.Next() {
		var evidenceID, evidenceType, validity, repositoryJSON string
		var receiptID, reviewID, humanOutcomeID, receiptOutcome, reviewVerdict, humanOutcomeResult, outcomeClass, scopeIssueID sql.NullString
		if err := rows.Scan(&evidenceID, &evidenceType, &receiptID, &reviewID, &humanOutcomeID, &validity, &repositoryJSON, &receiptOutcome, &reviewVerdict, &humanOutcomeResult, &outcomeClass, &scopeIssueID); err != nil {
			return nil, err
		}
		var evidenceTarget Repository
		if err := json.Unmarshal([]byte(repositoryJSON), &evidenceTarget); err != nil {
			return nil, fmt.Errorf("decode Evidence target repository: %w", err)
		}
		state := workflow.ObligationEvidenceState{
			ObligationID:         obligation.ID,
			EvidenceID:           evidenceID,
			EvidenceType:         evidenceType,
			ReceiptID:            receiptID.String,
			ReviewID:             reviewID.String,
			ReceiptSuccessful:    receiptOutcome.String == "succeeded",
			ReviewVerdict:        reviewVerdict.String,
			ReviewIndependent:    evidenceType != "review" || reviewID.Valid,
			HumanOutcomeID:       humanOutcomeID.String,
			HumanOutcomeResult:   humanOutcomeResult.String,
			Validity:             validity,
			Freshness:            freshnessForTarget(obligation.FreshnessTarget, evidenceTarget, target),
			ArtifactAvailability: "available",
		}
		if evidenceType == "human-outcome" {
			governing, err := workflow.ResolveGoverningHumanOutcome(conn, outcomeClass.String, scopeIssueID.String)
			if err != nil {
				return nil, err
			}
			state.OutcomeSuperseded = governing.ID != humanOutcomeID.String && governing.Result != humanOutcomeResult.String
		}
		if evidenceType == "review" {
			blocking, err := blockingFindingIDs(conn, `r.issue_id=? AND r.envelope_id=? AND r.envelope_revision=? AND NOT EXISTS (SELECT 1 FROM finding_dispositions d WHERE d.finding_id=f.id AND d.action IN ('resolved','inapplicable'))`, obligation.Subject.ID, envelopeID, revision)
			if err != nil {
				return nil, err
			}
			state.UnresolvedBlockingFindingIDs = blocking
		}
		if err := loadStateArtifacts(conn, &state); err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, rows.Err()
}

func loadStateArtifacts(conn *sql.DB, state *workflow.ObligationEvidenceState) error {
	rows, err := conn.Query(`SELECT l.artifact_id,a.availability FROM evidence_artifact_links l JOIN referenced_artifacts a ON a.id=l.artifact_id WHERE l.evidence_id=? AND l.required_material=1 ORDER BY l.artifact_id`, state.EvidenceID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id, availability string
		if err := rows.Scan(&id, &availability); err != nil {
			return err
		}
		state.ArtifactIDs = append(state.ArtifactIDs, id)
		if availability != "available" {
			state.ArtifactAvailability = "missing"
		}
	}
	return rows.Err()
}

func evaluationFromWorkflow(input workflow.ObligationEvaluation) Evaluation {
	return Evaluation{ObligationID: input.ObligationID, Status: input.Status, Validity: input.Validity, Freshness: input.Freshness, ArtifactAvailability: input.ArtifactAvailability, Reasons: input.Reasons}
}

func loadEvidence(conn *sql.DB, id string) (Evidence, error) {
	evidence := Evidence{RequiredArtifactIDs: []string{}}
	var repositoryJSON string
	err := conn.QueryRow(`SELECT id,obligation_id,evidence_type,receipt_id,review_id,human_outcome_id,producer,trust_source,target_repository_json,validity,created_at FROM evidence_records WHERE id=?`, id).Scan(&evidence.ID, &evidence.ObligationID, &evidence.EvidenceType, &evidence.ReceiptID, &evidence.ReviewID, &evidence.HumanOutcomeID, &evidence.Producer, &evidence.TrustSource, &repositoryJSON, &evidence.Validity, &evidence.CreatedAt)
	if err != nil {
		return Evidence{}, err
	}
	if err := json.Unmarshal([]byte(repositoryJSON), &evidence.TargetRepository); err != nil {
		return Evidence{}, err
	}
	rows, err := conn.Query(`SELECT artifact_id FROM evidence_artifact_links WHERE evidence_id = ? AND required_material = 1 ORDER BY artifact_id`, id)
	if err != nil {
		return Evidence{}, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var artifactID string
		if err := rows.Scan(&artifactID); err != nil {
			return Evidence{}, err
		}
		evidence.RequiredArtifactIDs = append(evidence.RequiredArtifactIDs, artifactID)
	}
	return evidence, rows.Err()
}

func nextEvidenceID(conn *sql.DB) (string, error) {
	return journal.GenerateID("EVD", func(candidate string) (bool, error) {
		var count int
		err := conn.QueryRow(`SELECT COUNT(*) FROM evidence_records WHERE id = ?`, candidate).Scan(&count)
		return count > 0, err
	})
}

func evidenceRequestDigest(projectID, actorID string, request EvidenceLinkRequest) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"actor_id": actorID, "project_id": projectID, "request": request})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func validEvidenceRepository(repository Repository) bool {
	return repository.CaptureSource == "git" && strings.TrimSpace(repository.ProjectID) != "" && strings.TrimSpace(repository.WorktreeID) != "" && validOID(repository.HeadOID) && validOID(repository.TreeOID) && (repository.BaseOID == nil || validOID(*repository.BaseOID)) && (!repository.Dirty || repository.DirtyDigest != nil) && (repository.Dirty || repository.DirtyDigest == nil)
}

func sameEvidenceRepository(left, right Repository) bool {
	return left.CaptureSource == right.CaptureSource && left.Dirty == right.Dirty && sameString(left.BaseOID, right.BaseOID) && sameString(left.DirtyDigest, right.DirtyDigest) && left.HeadOID == right.HeadOID && left.ProjectID == right.ProjectID && left.TreeOID == right.TreeOID && left.WorktreeID == right.WorktreeID
}

func sameString(left, right *string) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func freshnessForRepository(receipt, target Repository) string {
	if receipt.ProjectID == target.ProjectID && receipt.WorktreeID == target.WorktreeID && receipt.TreeOID == target.TreeOID {
		return "fresh"
	}
	return "stale"
}

func freshnessForTarget(freshnessTarget string, evidence, target Repository) string {
	if freshnessTarget == "none" {
		return "fresh"
	}
	return freshnessForRepository(evidence, target)
}

// workflowRepository converts the evidence package's own Repository shape to
// the workflow package's identical wire shape, so the shared human-outcome
// readiness helper can be reused across package boundaries without a type alias.
func workflowRepository(repository Repository) workflow.Repository {
	return workflow.Repository{BaseOID: repository.BaseOID, CaptureSource: repository.CaptureSource, Dirty: repository.Dirty, DirtyDigest: repository.DirtyDigest, HeadOID: repository.HeadOID, ProjectID: repository.ProjectID, TreeOID: repository.TreeOID, WorktreeID: repository.WorktreeID}
}

// substituteEvidenceEventPlaceholders resolves the "<event-id>"/"<occurred-at>"
// convention used by optional readiness-Stop Effects appended to this Event,
// mirroring the equivalent per-package helper in internal/review and
// internal/workflow. Pre-existing evidence_records/evidence_artifact_links
// Effects never use these placeholders, so this is a no-op for them.
func substituteEvidenceEventPlaceholders(effects []journal.Effect, event journal.Event) {
	occurredAt := event.OccurredAt.UTC().Format(time.RFC3339Nano)
	for index := range effects {
		for key, value := range effects[index].Values {
			switch value {
			case "<event-id>":
				effects[index].Values[key] = event.EventID
			case "<occurred-at>":
				effects[index].Values[key] = occurredAt
			}
		}
	}
}

func sortedUnique(values []string) bool {
	if !sort.StringsAreSorted(values) {
		return false
	}
	for index := 1; index < len(values); index++ {
		if values[index] == values[index-1] || strings.TrimSpace(values[index]) == "" {
			return false
		}
	}
	return len(values) == 0 || strings.TrimSpace(values[0]) != ""
}
