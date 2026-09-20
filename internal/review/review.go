// Package review owns immutable Review Record admission.
package review

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type RecordRequest struct {
	RunID              string         `json:"run_id"`
	ExpectedRunVersion int            `json:"expected_run_version"`
	Profile            string         `json:"profile"`
	Repository         Repository     `json:"repository"`
	Coverage           Coverage       `json:"coverage"`
	Verdict            string         `json:"verdict"`
	Findings           []FindingInput `json:"findings"`
	ResolvesFindingIDs []string       `json:"resolves_finding_ids"`
	EvidenceRefs       []Source       `json:"evidence_refs"`
	ReceiptIDs         []string       `json:"receipt_ids"`
}

type RecordOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
}

type Repository struct {
	BaseOID       *string `json:"base_oid"`
	CaptureSource string  `json:"capture_source"`
	Dirty         bool    `json:"dirty"`
	DirtyDigest   *string `json:"dirty_digest"`
	HeadOID       string  `json:"head_oid"`
	ProjectID     string  `json:"project_id"`
	TreeOID       string  `json:"tree_oid"`
	WorktreeID    string  `json:"worktree_id"`
}

type Source struct {
	SourceType   string     `json:"source_type"`
	DisplayLabel string     `json:"display_label"`
	Path         *string    `json:"path"`
	Section      *string    `json:"section"`
	ExternalID   *string    `json:"external_id"`
	Digest       string     `json:"digest"`
	Repository   Repository `json:"repository"`
}

type CoverageItem struct {
	Source   Source `json:"source"`
	Assessed bool   `json:"assessed"`
}

type CoverageExclusion struct {
	Source Source `json:"source"`
	Reason string `json:"reason"`
}

type Coverage struct {
	Items      []CoverageItem      `json:"items"`
	Exclusions []CoverageExclusion `json:"exclusions"`
}

type Subject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Scope struct {
	Paths []string `json:"paths"`
}

type FindingInput struct {
	ID       string  `json:"id"`
	Subject  Subject `json:"subject"`
	Scope    Scope   `json:"scope"`
	Severity string  `json:"severity"`
	Summary  string  `json:"summary"`
	Blocking bool    `json:"blocking"`
}

type NotIndependentError struct {
	ReviewActorID      string
	ConflictingType    string
	ConflictingID      string
	ConflictingActorID string
}

func (err *NotIndependentError) Error() string { return "review actor is not independent" }

type VerdictInvalidError struct {
	Verdict string
	Reason  string
}

func (err *VerdictInvalidError) Error() string { return "review verdict is invalid" }

type RecordResult struct {
	Review         map[string]any `json:"review"`
	IdempotencyKey string         `json:"idempotency_key"`
	Replayed       bool           `json:"replayed"`
	EventID        string
}

type RepositoryMismatchError struct {
	Expected Repository
	Observed Repository
}

func (err *RepositoryMismatchError) Error() string {
	return "review repository does not match current repository"
}

type ReviewNotFoundError struct{ ID string }

func (err *ReviewNotFoundError) Error() string { return fmt.Sprintf("review %q was not found", err.ID) }

type ReviewConflictError struct{ RunID string }

func (err *ReviewConflictError) Error() string { return "review already exists for run" }

// Get returns the immutable Review projection and its derived Finding state.
func Get(workspaceRoot, migrationsDir, reviewID string) (map[string]any, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return nil, err
	}
	conn, err := fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return nil, fmt.Errorf("open review projection: %w", err)
	}
	defer func() { _ = conn.Close() }()
	return loadReview(conn, reviewID)
}

// Record currently establishes the durable independence admission boundary.
// Subsequent Task 4 slices append the immutable Review Event after this check.
func Record(options RecordOptions, request RecordRequest) (RecordResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.IdempotencyKey) == "" || strings.TrimSpace(request.RunID) == "" {
		return RecordResult{}, fmt.Errorf("review record input is incomplete")
	}
	var result RecordResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid review record"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return fmt.Errorf("open review projection: %w", err)
		}
		defer func() { _ = conn.Close() }()
		projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
		if err != nil {
			return fmt.Errorf("review record: read project config: %w", err)
		}
		requestDigest, err := reviewRequestDigest(projectConfig.ProjectID, options.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "review.record:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return recordReplay(outcome, options.IdempotencyKey, &result)
		}
		binding, err := loadRecordBinding(conn, request.RunID)
		if err != nil {
			return err
		}
		if err := checkIndependence(conn, options.ActorID, request.RunID); err != nil {
			return err
		}
		if binding.Version != request.ExpectedRunVersion || binding.Status != "active" || binding.ClaimMode != "review" || binding.ClaimActorID != options.ActorID {
			return fmt.Errorf("review record Run binding is invalid")
		}
		if err := validateCoverage(request.Coverage); err != nil {
			return err
		}
		if err := validateVerdictFindingAgreement(request.Verdict, request.Findings, request.Coverage); err != nil {
			return err
		}
		if len(request.ResolvesFindingIDs) > 0 && request.Verdict != "pass" && request.Verdict != "changes-required" {
			return &VerdictInvalidError{Verdict: request.Verdict, Reason: "resolving Findings requires pass or changes-required"}
		}
		if err := workflow.AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, workflow.EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: "review.record", Capability: "review.record", IssueID: binding.IssueID, EnvelopeID: binding.EnvelopeID, Revision: binding.EnvelopeRevision}); err != nil {
			return err
		}
		observed, err := observedRepository(options.WorkspaceRoot, projectConfig.ProjectID)
		if err != nil {
			return err
		}
		if !sameRepository(request.Repository, observed) {
			return &RepositoryMismatchError{Expected: request.Repository, Observed: observed}
		}
		var existing string
		if err := conn.QueryRow(`SELECT id FROM review_records WHERE run_id = ?`, request.RunID).Scan(&existing); err == nil {
			return &ReviewConflictError{RunID: request.RunID}
		} else if err != sql.ErrNoRows {
			return fmt.Errorf("review record: read existing Review: %w", err)
		}
		for _, receiptID := range request.ReceiptIDs {
			var receiptRunID string
			if err := conn.QueryRow(`SELECT a.run_id FROM execution_receipts r JOIN execution_attempts a ON a.id = r.attempt_id WHERE r.id = ?`, receiptID).Scan(&receiptRunID); err != nil || receiptRunID != request.RunID {
				return fmt.Errorf("review record receipt %q does not belong to Run", receiptID)
			}
		}
		actor, err := workflow.ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		reviewID, err := journal.GenerateID("REV", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM review_records WHERE id = ?`, candidate).Scan(&found)
			return err == nil, nil
		})
		if err != nil {
			return err
		}
		resolutionEffects, err := resolveFindingEffects(conn, request.ResolvesFindingIDs, binding, options.ActorID, reviewID)
		if err != nil {
			return err
		}
		next, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: reviewID, Type: "review"})
		if err != nil {
			return err
		}
		review, err := recordReviewValue(reviewID, binding, options.ActorID, request, observed, "<occurred-at>", "<event-id>")
		if err != nil {
			return err
		}
		effects := append(reviewRecordEffects(review, "<event-id>"), resolutionEffects...)
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "review.recorded", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: reviewID, Sequence: next, Type: "review"}, CommandName: "review.record", Repository: journalRepository(observed), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			setReviewEventIdentity(review, prepared)
			resultJSON, err := journal.MarshalCanonicalJSON(map[string]any{"review": review, "budget": nil, "stops": []any{}, "attention": []any{}, "idempotency_key": options.IdempotencyKey, "replayed": false})
			if err != nil {
				return journal.IdempotencyMutation{}, err
			}
			return journal.IdempotencyMutation{EntityID: reviewID, ResultJSON: string(resultJSON)}, nil
		})
		if err != nil {
			return err
		}
		if outcome.Replayed {
			return recordReplay(outcome, options.IdempotencyKey, &result)
		}
		result = RecordResult{Review: review, IdempotencyKey: options.IdempotencyKey, Replayed: false, EventID: outcome.EventID}
		return nil
	})
	if err != nil {
		return RecordResult{}, err
	}
	return result, nil
}

func resolveFindingEffects(conn *sql.DB, findingIDs []string, binding recordBinding, actorID, reviewID string) ([]journal.Effect, error) {
	seen := make(map[string]struct{}, len(findingIDs))
	effects := make([]journal.Effect, 0, len(findingIDs))
	for _, findingID := range findingIDs {
		if _, found := seen[findingID]; found {
			return nil, fmt.Errorf("review resolution Finding %q is duplicated", findingID)
		}
		seen[findingID] = struct{}{}
		var issueID, envelopeID, creatorReviewID string
		var revision int
		err := conn.QueryRow(`SELECT r.issue_id, r.envelope_id, r.envelope_revision, f.review_id FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE f.id=?`, findingID).Scan(&issueID, &envelopeID, &revision, &creatorReviewID)
		if err != nil || issueID != binding.IssueID || envelopeID != binding.EnvelopeID || revision != binding.EnvelopeRevision || creatorReviewID == reviewID {
			return nil, fmt.Errorf("review resolution Finding %q is invalid", findingID)
		}
		var terminal int
		if err := conn.QueryRow(`SELECT 1 FROM finding_dispositions WHERE finding_id=? AND action IN ('resolved','inapplicable')`, findingID).Scan(&terminal); err != sql.ErrNoRows {
			if err == nil {
				return nil, fmt.Errorf("review resolution Finding %q is terminal", findingID)
			}
			return nil, fmt.Errorf("review resolution Finding %q: %w", findingID, err)
		}
		dispositionID, err := journal.GenerateID("FDP", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM finding_dispositions WHERE id=?`, candidate).Scan(&found)
			return err == nil, nil
		})
		if err != nil {
			return nil, err
		}
		effects = append(effects, journal.Effect{Op: "insert", Table: "finding_dispositions", PK: map[string]any{"id": dispositionID}, Values: map[string]any{"finding_id": findingID, "action": "resolved", "actor_id": actorID, "claim_id": nil, "qualifying_review_id": reviewID, "reason": "resolved by qualifying Review", "disposed_at": "<occurred-at>", "created_event_id": "<event-id>"}})
	}
	return effects, nil
}

func checkIndependence(conn *sql.DB, actorID, runID string) error {
	var issueID, envelopeID, claimID string
	var revision int
	err := conn.QueryRow(`SELECT issue_id, envelope_id, envelope_revision, claim_id FROM runs WHERE id = ?`, runID).Scan(&issueID, &envelopeID, &revision, &claimID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("review Run %q was not found", runID)
	}
	if err != nil {
		return fmt.Errorf("load review Run: %w", err)
	}
	var conflictingID string
	err = conn.QueryRow(`SELECT id FROM claims WHERE issue_id = ? AND envelope_id = ? AND envelope_revision = ? AND mode = 'mutation' AND actor_id = ? ORDER BY acquired_at, id LIMIT 1`, issueID, envelopeID, revision, actorID).Scan(&conflictingID)
	if err == nil {
		return &NotIndependentError{ReviewActorID: actorID, ConflictingType: "claim", ConflictingID: conflictingID, ConflictingActorID: actorID}
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("load mutating Claim history: %w", err)
	}
	err = conn.QueryRow(`SELECT r.id FROM runs r JOIN claims c ON c.id = r.claim_id WHERE r.issue_id = ? AND r.envelope_id = ? AND r.envelope_revision = ? AND c.mode = 'mutation' AND r.actor_id = ? ORDER BY r.started_at, r.id LIMIT 1`, issueID, envelopeID, revision, actorID).Scan(&conflictingID)
	if err == nil {
		return &NotIndependentError{ReviewActorID: actorID, ConflictingType: "run", ConflictingID: conflictingID, ConflictingActorID: actorID}
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("load mutating Run history: %w", err)
	}
	return nil
}

func validateVerdictFindingAgreement(verdict string, findings []FindingInput, coverage Coverage) error {
	blocking := false
	for _, finding := range findings {
		if (finding.Severity == "blocker" || finding.Severity == "major") != finding.Blocking {
			return &VerdictInvalidError{Verdict: verdict, Reason: "finding severity and blocking disagree"}
		}
		blocking = blocking || finding.Blocking
	}
	switch verdict {
	case "pass":
		if blocking {
			return &VerdictInvalidError{Verdict: verdict, Reason: "pass cannot carry a blocking finding"}
		}
	case "changes-required":
		if !blocking {
			return &VerdictInvalidError{Verdict: verdict, Reason: "changes-required requires a blocking finding"}
		}
	case "blocked":
		if len(findings) == 0 && len(coverage.Exclusions) == 0 {
			return &VerdictInvalidError{Verdict: verdict, Reason: "blocked requires a finding or exclusion"}
		}
	default:
		return &VerdictInvalidError{Verdict: verdict, Reason: "verdict is not accepted"}
	}
	return nil
}

func validateCoverage(coverage Coverage) error {
	if len(coverage.Items) == 0 || len(coverage.Items) > 64 || len(coverage.Exclusions) > 64 {
		return fmt.Errorf("review coverage cardinality is invalid")
	}
	digests := make(map[string]struct{}, len(coverage.Items)+len(coverage.Exclusions))
	for _, item := range coverage.Items {
		if strings.TrimSpace(item.Source.Digest) == "" {
			return fmt.Errorf("review coverage Source digest is required")
		}
		if _, found := digests[item.Source.Digest]; found {
			return fmt.Errorf("review coverage Source digest is duplicated")
		}
		digests[item.Source.Digest] = struct{}{}
	}
	for _, exclusion := range coverage.Exclusions {
		if strings.TrimSpace(exclusion.Source.Digest) == "" || strings.TrimSpace(exclusion.Reason) == "" {
			return fmt.Errorf("review coverage exclusion is incomplete")
		}
		if _, found := digests[exclusion.Source.Digest]; found {
			return fmt.Errorf("review coverage Source digest is duplicated")
		}
		digests[exclusion.Source.Digest] = struct{}{}
	}
	return nil
}

type recordBinding struct {
	IssueID, EnvelopeID, ClaimMode, ClaimActorID, Status string
	EnvelopeRevision, Version                            int
}

func loadRecordBinding(conn *sql.DB, runID string) (recordBinding, error) {
	var binding recordBinding
	err := conn.QueryRow(`SELECT r.issue_id, r.envelope_id, r.envelope_revision, r.status, rt.version, c.mode, c.actor_id FROM runs r JOIN run_runtime rt ON rt.run_id=r.id JOIN claims c ON c.id=r.claim_id WHERE r.id=?`, runID).Scan(&binding.IssueID, &binding.EnvelopeID, &binding.EnvelopeRevision, &binding.Status, &binding.Version, &binding.ClaimMode, &binding.ClaimActorID)
	if err == sql.ErrNoRows {
		return recordBinding{}, fmt.Errorf("review Run %q was not found", runID)
	}
	if err != nil {
		return recordBinding{}, fmt.Errorf("load review Run: %w", err)
	}
	return binding, nil
}

func observedRepository(root, projectID string) (Repository, error) {
	observation, err := project.ObserveRepository(root)
	if err != nil {
		return Repository{}, fmt.Errorf("review record: observe repository: %w", err)
	}
	base := observation.HeadOID
	return Repository{BaseOID: &base, CaptureSource: "git", Dirty: observation.Dirty, DirtyDigest: observation.DirtyDigest, HeadOID: observation.HeadOID, ProjectID: projectID, TreeOID: observation.TreeOID, WorktreeID: observation.WorktreeID}, nil
}

func sameRepository(left, right Repository) bool {
	return left.CaptureSource == right.CaptureSource && left.Dirty == right.Dirty && sameString(left.BaseOID, right.BaseOID) && sameString(left.DirtyDigest, right.DirtyDigest) && left.HeadOID == right.HeadOID && left.ProjectID == right.ProjectID && left.TreeOID == right.TreeOID && left.WorktreeID == right.WorktreeID
}

func sameString(left, right *string) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func journalRepository(repository Repository) *journal.Repository {
	base := ""
	if repository.BaseOID != nil {
		base = *repository.BaseOID
	}
	return &journal.Repository{BaseOID: base, CaptureSource: repository.CaptureSource, Dirty: repository.Dirty, DirtyDigest: repository.DirtyDigest, HeadOID: repository.HeadOID, ProjectID: repository.ProjectID, TreeOID: repository.TreeOID, WorktreeID: repository.WorktreeID}
}

func reviewRequestDigest(projectID, actorID string, request RecordRequest) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"operation": "review.record", "project_id": projectID, "actor_id": actorID, "request": request})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func recordReplay(outcome journal.IdempotencyOutcome, key string, result *RecordResult) error {
	var decoded struct {
		Review map[string]any `json:"review"`
	}
	if err := json.Unmarshal([]byte(outcome.ResultJSON), &decoded); err != nil {
		return fmt.Errorf("decode Review replay result: %w", err)
	}
	*result = RecordResult{Review: decoded.Review, IdempotencyKey: key, Replayed: true, EventID: outcome.EventID}
	return nil
}

func loadReview(conn *sql.DB, reviewID string) (map[string]any, error) {
	var review map[string]any
	var repositoryJSON, coverageJSON, evidenceJSON, receiptJSON []byte
	var revision int
	var runID, issueID, envelopeID, actorID, profile, verdict, independence, recordedAt, eventID string
	err := conn.QueryRow(`SELECT run_id, issue_id, envelope_id, envelope_revision, actor_id, profile, repository_json, coverage_json, verdict, evidence_refs_json, receipt_ids_json, independence, recorded_at, created_event_id FROM review_records WHERE id=?`, reviewID).Scan(&runID, &issueID, &envelopeID, &revision, &actorID, &profile, &repositoryJSON, &coverageJSON, &verdict, &evidenceJSON, &receiptJSON, &independence, &recordedAt, &eventID)
	if err == sql.ErrNoRows {
		return nil, &ReviewNotFoundError{ID: reviewID}
	}
	if err != nil {
		return nil, fmt.Errorf("load Review: %w", err)
	}
	var repository, coverage, evidenceRefs, receiptIDs any
	if err := json.Unmarshal(repositoryJSON, &repository); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(coverageJSON, &coverage); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(evidenceJSON, &evidenceRefs); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(receiptJSON, &receiptIDs); err != nil {
		return nil, err
	}
	findings, err := loadReviewFindings(conn, reviewID)
	if err != nil {
		return nil, err
	}
	resolvedFindingIDs, err := loadResolvedFindingIDs(conn, reviewID)
	if err != nil {
		return nil, err
	}
	review = map[string]any{"id": reviewID, "run_id": runID, "issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": revision, "actor_id": actorID, "profile": profile, "repository": repository, "coverage": coverage, "verdict": verdict, "evidence_refs": evidenceRefs, "receipt_ids": receiptIDs, "findings": findings, "resolved_finding_ids": resolvedFindingIDs, "independence": independence, "recorded_at": recordedAt, "event_id": eventID}
	return review, nil
}

func loadReviewFindings(conn *sql.DB, reviewID string) ([]any, error) {
	rows, err := conn.Query(`SELECT id, subject_type, subject_id, scope_json, severity, summary, blocking, created_at FROM review_findings WHERE review_id=? ORDER BY id`, reviewID)
	if err != nil {
		return nil, fmt.Errorf("load Review Findings: %w", err)
	}
	defer func() { _ = rows.Close() }()
	findings := []any{}
	for rows.Next() {
		var id, subjectType, subjectID, scopeJSON, severity, summary, createdAt string
		var blocking bool
		if err := rows.Scan(&id, &subjectType, &subjectID, &scopeJSON, &severity, &summary, &blocking, &createdAt); err != nil {
			return nil, err
		}
		var scope any
		if err := json.Unmarshal([]byte(scopeJSON), &scope); err != nil {
			return nil, err
		}
		dispositions, err := loadFindingDispositions(conn, id)
		if err != nil {
			return nil, err
		}
		findings = append(findings, map[string]any{"id": id, "subject": map[string]any{"type": subjectType, "id": subjectID}, "scope": scope, "severity": severity, "summary": summary, "blocking": blocking, "review_id": reviewID, "created_at": createdAt, "dispositions": dispositions, "state": findingState(dispositions)})
	}
	return findings, rows.Err()
}

func loadResolvedFindingIDs(conn *sql.DB, reviewID string) ([]any, error) {
	rows, err := conn.Query(`SELECT finding_id FROM finding_dispositions WHERE qualifying_review_id=? AND action='resolved' ORDER BY finding_id`, reviewID)
	if err != nil {
		return nil, fmt.Errorf("load resolved Finding IDs: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ids := []any{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func loadFindingDispositions(conn *sql.DB, findingID string) ([]any, error) {
	rows, err := conn.Query(`SELECT id, action, actor_id, qualifying_review_id, reason, disposed_at, created_event_id FROM finding_dispositions WHERE finding_id=? ORDER BY disposed_at, id`, findingID)
	if err != nil {
		return nil, fmt.Errorf("load Finding dispositions: %w", err)
	}
	defer func() { _ = rows.Close() }()
	dispositions := []any{}
	for rows.Next() {
		var id, action, actorID, reason, disposedAt, eventID string
		var qualifyingReviewID *string
		if err := rows.Scan(&id, &action, &actorID, &qualifyingReviewID, &reason, &disposedAt, &eventID); err != nil {
			return nil, err
		}
		dispositions = append(dispositions, map[string]any{"id": id, "action": action, "actor_id": actorID, "qualifying_review_id": qualifyingReviewID, "reason": reason, "disposed_at": disposedAt, "event_id": eventID})
	}
	return dispositions, rows.Err()
}

func findingState(dispositions []any) string {
	if len(dispositions) == 0 {
		return "open"
	}
	return dispositions[len(dispositions)-1].(map[string]any)["action"].(string)
}

func recordReviewValue(id string, binding recordBinding, actorID string, request RecordRequest, repository Repository, recordedAt, eventID string) (map[string]any, error) {
	value, err := reviewJSONValue(map[string]any{"id": id, "run_id": request.RunID, "issue_id": binding.IssueID, "envelope_id": binding.EnvelopeID, "envelope_revision": binding.EnvelopeRevision, "actor_id": actorID, "profile": request.Profile, "repository": repository, "coverage": request.Coverage, "verdict": request.Verdict, "evidence_refs": request.EvidenceRefs, "receipt_ids": request.ReceiptIDs, "findings": request.Findings, "resolved_finding_ids": request.ResolvesFindingIDs, "independence": "distinct-actor", "recorded_at": recordedAt, "event_id": eventID})
	if err != nil {
		return nil, err
	}
	if value["evidence_refs"] == nil {
		value["evidence_refs"] = []any{}
	}
	if value["receipt_ids"] == nil {
		value["receipt_ids"] = []any{}
	}
	if value["findings"] == nil {
		value["findings"] = []any{}
	}
	if value["resolved_finding_ids"] == nil {
		value["resolved_finding_ids"] = []any{}
	}
	for _, raw := range value["findings"].([]any) {
		finding := raw.(map[string]any)
		finding["review_id"] = id
		finding["created_at"] = recordedAt
		finding["dispositions"] = []any{}
		finding["state"] = "open"
	}
	return value, nil
}

func reviewJSONValue(input any) (map[string]any, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var value map[string]any
	if err := json.Unmarshal(payload, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func reviewRecordEffects(review map[string]any, eventID string) []journal.Effect {
	marshal := func(value any) string {
		payload, err := journal.MarshalCanonicalJSON(value)
		if err != nil {
			panic(err)
		}
		return string(payload)
	}
	effects := []journal.Effect{{Op: "insert", Table: "review_records", PK: map[string]any{"id": review["id"]}, Values: map[string]any{"run_id": review["run_id"], "issue_id": review["issue_id"], "envelope_id": review["envelope_id"], "envelope_revision": review["envelope_revision"], "actor_id": review["actor_id"], "profile": review["profile"], "repository_json": marshal(review["repository"]), "coverage_json": marshal(review["coverage"]), "verdict": review["verdict"], "evidence_refs_json": marshal(review["evidence_refs"]), "receipt_ids_json": marshal(review["receipt_ids"]), "independence": review["independence"], "recorded_at": review["recorded_at"], "created_event_id": eventID}}}
	for _, raw := range review["findings"].([]any) {
		finding := raw.(map[string]any)
		effects = append(effects, journal.Effect{Op: "insert", Table: "review_findings", PK: map[string]any{"id": finding["id"]}, Values: map[string]any{"review_id": review["id"], "subject_type": finding["subject"].(map[string]any)["type"], "subject_id": finding["subject"].(map[string]any)["id"], "scope_json": marshal(finding["scope"]), "severity": finding["severity"], "summary": finding["summary"], "blocking": finding["blocking"], "created_at": review["recorded_at"], "created_event_id": eventID}})
	}
	return effects
}

func setReviewEventIdentity(review map[string]any, event journal.Event) {
	recordedAt := event.OccurredAt.UTC().Format(time.RFC3339Nano)
	review["recorded_at"], review["event_id"] = recordedAt, event.EventID
	for _, raw := range review["findings"].([]any) {
		raw.(map[string]any)["created_at"] = recordedAt
	}
	for index := range event.Effects {
		for key, value := range event.Effects[index].Values {
			if value == "<event-id>" {
				event.Effects[index].Values[key] = event.EventID
			}
			if value == "<occurred-at>" {
				event.Effects[index].Values[key] = recordedAt
			}
		}
	}
}
