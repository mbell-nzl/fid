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

// Evidence is the immutable fid-created link between one Obligation and one
// exact Receipt.
type Evidence struct {
	ID                  string     `json:"id"`
	ObligationID        string     `json:"obligation_id"`
	ReceiptID           string     `json:"receipt_id"`
	EvidenceType        string     `json:"evidence_type"`
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
	ReceiptID           string     `json:"receipt_id"`
	TargetRepository    Repository `json:"target_repository"`
	RequiredArtifactIDs []string   `json:"required_artifact_ids"`
}

type EvidenceEvaluateRequest struct {
	ObligationID     string     `json:"obligation_id"`
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
		evidence := Evidence{ID: evidenceID, ObligationID: options.Request.ObligationID, ReceiptID: options.Request.ReceiptID, EvidenceType: "receipt", Producer: "fid", TrustSource: "fid-executed", TargetRepository: options.Request.TargetRepository, Validity: "valid", RequiredArtifactIDs: append([]string(nil), options.Request.RequiredArtifactIDs...), CreatedAt: createdAt}
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
		effects := []journal.Effect{{Op: "insert", Table: "evidence_records", PK: map[string]any{"id": evidence.ID}, Values: map[string]any{"obligation_id": evidence.ObligationID, "receipt_id": evidence.ReceiptID, "evidence_type": evidence.EvidenceType, "producer": evidence.Producer, "trust_source": evidence.TrustSource, "target_repository_json": string(targetJSON), "validity": evidence.Validity, "created_at": evidence.CreatedAt}}}
		for _, artifactID := range evidence.RequiredArtifactIDs {
			effects = append(effects, journal.Effect{Op: "insert", Table: "evidence_artifact_links", PK: map[string]any{"evidence_id": evidence.ID, "artifact_id": artifactID}, Values: map[string]any{"required_material": true}})
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "evidence.linked", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: evidence.ID, Sequence: 1, Type: "evidence"}, CommandName: "evidence.link", CorrelationID: &binding.RunID, Repository: repositoryEventValue(evidence.TargetRepository), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(journal.Event) (journal.IdempotencyMutation, error) {
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
func EvaluateEvidence(workspaceRoot, migrationsDir string, request EvidenceEvaluateRequest) (Evaluation, error) {
	if strings.TrimSpace(workspaceRoot) == "" || strings.TrimSpace(migrationsDir) == "" || strings.TrimSpace(request.ObligationID) == "" || !validEvidenceRepository(request.TargetRepository) {
		return Evaluation{}, &EvidenceBindingError{Reason: "Evidence evaluation input is invalid"}
	}
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Evaluation{}, err
	}
	conn, err := openEvidenceDB(workspaceRoot)
	if err != nil {
		return Evaluation{}, err
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
	EnvelopeID string
	Revision   int
	IssueID    string
	RunID      string
	ReceiptID  string
	Outcome    string
	Repository Repository
	Artifacts  map[string]string
}

func validateEvidenceLinkOptions(options EvidenceLinkOptions) error {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.IdempotencyKey) == "" || strings.TrimSpace(options.Request.ObligationID) == "" || strings.TrimSpace(options.Request.ReceiptID) == "" || !validEvidenceRepository(options.Request.TargetRepository) || !sortedUnique(options.Request.RequiredArtifactIDs) {
		return &EvidenceBindingError{Reason: "Evidence link input is invalid"}
	}
	return nil
}

func loadEvidenceBinding(conn *sql.DB, request EvidenceLinkRequest) (evidenceBinding, error) {
	var binding evidenceBinding
	var obligationJSON, receiptRepositoryJSON, trustSource, producer string
	err := conn.QueryRow(`SELECT o.envelope_id, o.revision, o.obligation_json, r.id, r.terminal_status, r.terminal_repository_json, r.trust_source, a.producer, a.registration_id, a.registration_version, a.obligation_id, a.run_id, run.issue_id, run.envelope_id, run.envelope_revision FROM envelope_obligations o JOIN execution_receipts r ON r.id = ? JOIN execution_attempts a ON a.id = r.attempt_id JOIN runs run ON run.id = a.run_id WHERE o.id = ?`, request.ReceiptID, request.ObligationID).Scan(&binding.EnvelopeID, &binding.Revision, &obligationJSON, &binding.ReceiptID, &binding.Outcome, &receiptRepositoryJSON, &trustSource, &producer, new(string), new(int), new(string), &binding.RunID, &binding.IssueID, new(string), new(int))
	if err == sql.ErrNoRows {
		return evidenceBinding{}, &EvidenceBindingError{Reason: "Receipt and Obligation binding is not exact"}
	}
	if err != nil {
		return evidenceBinding{}, fmt.Errorf("load Evidence binding: %w", err)
	}
	var obligation struct {
		Subject struct {
			ID string `json:"id"`
		} `json:"subject"`
		Registration struct {
			ID      string `json:"id"`
			Version int    `json:"version"`
		} `json:"registration"`
	}
	if err := json.Unmarshal([]byte(obligationJSON), &obligation); err != nil {
		return evidenceBinding{}, fmt.Errorf("decode Evidence Obligation: %w", err)
	}
	if err := json.Unmarshal([]byte(receiptRepositoryJSON), &binding.Repository); err != nil {
		return evidenceBinding{}, fmt.Errorf("decode Receipt repository: %w", err)
	}
	var registrationID, attemptObligationID, runEnvelopeID string
	var registrationVersion, runRevision int
	err = conn.QueryRow(`SELECT a.registration_id, a.registration_version, a.obligation_id, run.envelope_id, run.envelope_revision FROM execution_receipts r JOIN execution_attempts a ON a.id = r.attempt_id JOIN runs run ON run.id = a.run_id WHERE r.id = ?`, request.ReceiptID).Scan(&registrationID, &registrationVersion, &attemptObligationID, &runEnvelopeID, &runRevision)
	if err != nil {
		return evidenceBinding{}, fmt.Errorf("load Receipt chain: %w", err)
	}
	if trustSource != "fid-executed" || producer != "fid" || registrationID != obligation.Registration.ID || registrationVersion != obligation.Registration.Version || attemptObligationID != request.ObligationID || binding.EnvelopeID != runEnvelopeID || binding.Revision != runRevision || obligation.Subject.ID != binding.IssueID || !sameEvidenceRepository(binding.Repository, request.TargetRepository) {
		return evidenceBinding{}, &EvidenceBindingError{Reason: "Receipt and Obligation binding is not exact"}
	}
	binding.Artifacts = map[string]string{}
	rows, err := conn.Query(`SELECT artifact_id, availability FROM receipt_artifacts ra JOIN referenced_artifacts a ON a.id = ra.artifact_id WHERE ra.receipt_id = ? ORDER BY artifact_id`, request.ReceiptID)
	if err != nil {
		return evidenceBinding{}, fmt.Errorf("load Receipt Artifacts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id, availability string
		if err := rows.Scan(&id, &availability); err != nil {
			return evidenceBinding{}, err
		}
		binding.Artifacts[id] = availability
	}
	if err := rows.Err(); err != nil {
		return evidenceBinding{}, err
	}
	for _, artifactID := range request.RequiredArtifactIDs {
		if _, found := binding.Artifacts[artifactID]; !found {
			return evidenceBinding{}, &EvidenceBindingError{Reason: "required Artifact is not linked to Receipt"}
		}
	}
	return binding, nil
}

func evaluateEvidence(conn *sql.DB, request EvidenceEvaluateRequest) (Evaluation, error) {
	rows, err := conn.Query(`SELECT e.id, e.receipt_id, e.validity, r.terminal_status, r.terminal_repository_json FROM evidence_records e JOIN execution_receipts r ON r.id = e.receipt_id WHERE e.obligation_id = ? ORDER BY e.id`, request.ObligationID)
	if err != nil {
		return Evaluation{}, fmt.Errorf("evaluate Evidence: %w", err)
	}
	defer func() { _ = rows.Close() }()
	states := []workflow.ObligationEvidenceState{}
	validity, freshness, availability := "invalid", "stale", "missing"
	for rows.Next() {
		var evidenceID, receiptID, recordValidity, outcome, repositoryJSON string
		if err := rows.Scan(&evidenceID, &receiptID, &recordValidity, &outcome, &repositoryJSON); err != nil {
			return Evaluation{}, err
		}
		var target Repository
		if err := json.Unmarshal([]byte(repositoryJSON), &target); err != nil {
			return Evaluation{}, fmt.Errorf("decode Evidence Receipt repository: %w", err)
		}
		state := workflow.ObligationEvidenceState{ObligationID: request.ObligationID, EvidenceID: evidenceID, ReceiptID: receiptID, ReceiptSuccessful: outcome == "succeeded", Validity: recordValidity, Freshness: freshnessForRepository(target, request.TargetRepository), ArtifactAvailability: "available"}
		artifactRows, err := conn.Query(`SELECT l.artifact_id, a.availability FROM evidence_artifact_links l JOIN referenced_artifacts a ON a.id = l.artifact_id WHERE l.evidence_id = ? AND l.required_material = 1 ORDER BY l.artifact_id`, evidenceID)
		if err != nil {
			return Evaluation{}, err
		}
		for artifactRows.Next() {
			var artifactID, artifactState string
			if err := artifactRows.Scan(&artifactID, &artifactState); err != nil {
				_ = artifactRows.Close()
				return Evaluation{}, err
			}
			state.ArtifactIDs = append(state.ArtifactIDs, artifactID)
			if artifactState != "available" {
				state.ArtifactAvailability = "missing"
			}
		}
		if err := artifactRows.Close(); err != nil {
			return Evaluation{}, err
		}
		states = append(states, state)
		if recordValidity == "valid" {
			validity = "valid"
		}
		if state.Freshness == "fresh" {
			freshness = "fresh"
		}
		if state.ArtifactAvailability == "available" {
			availability = "available"
		}
	}
	if err := rows.Err(); err != nil {
		return Evaluation{}, err
	}
	if len(states) == 0 {
		states = append(states, workflow.ObligationEvidenceState{ObligationID: request.ObligationID})
	}
	reasons := workflow.EvaluateEvidenceStates(states)
	status := "unsatisfied"
	if len(reasons) == 0 {
		status, validity, freshness, availability = "satisfied", "valid", "fresh", "available"
	}
	return Evaluation{ObligationID: request.ObligationID, Status: status, Validity: validity, Freshness: freshness, ArtifactAvailability: availability, Reasons: reasons}, nil
}

func evaluationFromBinding(evidence Evidence, binding evidenceBinding, target Repository) Evaluation {
	availability := "available"
	artifactIDs := append([]string(nil), evidence.RequiredArtifactIDs...)
	for _, id := range artifactIDs {
		if binding.Artifacts[id] != "available" {
			availability = "missing"
		}
	}
	state := workflow.ObligationEvidenceState{ObligationID: evidence.ObligationID, EvidenceID: evidence.ID, ReceiptID: evidence.ReceiptID, ReceiptSuccessful: binding.Outcome == "succeeded", Validity: evidence.Validity, Freshness: freshnessForRepository(binding.Repository, target), ArtifactAvailability: availability, ArtifactIDs: artifactIDs}
	reasons := workflow.EvaluateEvidenceStates([]workflow.ObligationEvidenceState{state})
	status := "unsatisfied"
	if len(reasons) == 0 {
		status = "satisfied"
	}
	return Evaluation{ObligationID: evidence.ObligationID, Status: status, Validity: evidence.Validity, Freshness: state.Freshness, ArtifactAvailability: availability, Reasons: reasons}
}

func loadEvidence(conn *sql.DB, id string) (Evidence, error) {
	var evidence Evidence
	var repositoryJSON string
	err := conn.QueryRow(`SELECT id, obligation_id, receipt_id, evidence_type, producer, trust_source, target_repository_json, validity, created_at FROM evidence_records WHERE id = ?`, id).Scan(&evidence.ID, &evidence.ObligationID, &evidence.ReceiptID, &evidence.EvidenceType, &evidence.Producer, &evidence.TrustSource, &repositoryJSON, &evidence.Validity, &evidence.CreatedAt)
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
