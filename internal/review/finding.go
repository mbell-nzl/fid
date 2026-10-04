package review

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type AddressRequest struct {
	FindingID            string `json:"finding_id"`
	ClaimID              string `json:"claim_id"`
	ExpectedClaimVersion int    `json:"expected_claim_version"`
	Reason               string `json:"reason"`
}

type AddressOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
	Now            time.Time
}

type AddressResult struct {
	Finding        map[string]any `json:"finding"`
	IdempotencyKey string         `json:"idempotency_key"`
	Replayed       bool           `json:"replayed"`
	EventID        string
}

type ListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	ReviewID      string
	State         string
	Severity      string
	Limit         int
}

type FindingNotFoundError struct{ ID string }

func (err *FindingNotFoundError) Error() string {
	return fmt.Sprintf("Finding %q was not found", err.ID)
}

type FindingAlreadyDisposedError struct {
	ID    string
	State string
}

func (err *FindingAlreadyDisposedError) Error() string {
	return fmt.Sprintf("Finding %q is already %s", err.ID, err.State)
}

type FindingTargetInvalidError struct{ ID string }

func (err *FindingTargetInvalidError) Error() string {
	return fmt.Sprintf("Finding %q is not a valid disposition target", err.ID)
}

type ResultLimitExceededError struct {
	Limit    int
	Observed int
}

func (err *ResultLimitExceededError) Error() string { return "Finding result limit exceeded" }

func Address(options AddressOptions, request AddressRequest) (AddressResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.IdempotencyKey) == "" || strings.TrimSpace(request.FindingID) == "" || strings.TrimSpace(request.ClaimID) == "" || request.ExpectedClaimVersion < 1 || strings.TrimSpace(request.Reason) == "" || len(request.Reason) > 4096 {
		return AddressResult{}, fmt.Errorf("finding address input is invalid")
	}
	var result AddressResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid finding address"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return fmt.Errorf("open Finding projection: %w", err)
		}
		defer func() { _ = conn.Close() }()
		projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
		if err != nil {
			return fmt.Errorf("finding address: read project config: %w", err)
		}
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"operation": "finding.address", "project_id": projectConfig.ProjectID, "actor_id": options.ActorID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "finding.address:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return replayAddress(replay, options.IdempotencyKey, &result)
		}
		binding, err := loadFindingBinding(conn, request.FindingID)
		if err != nil {
			return err
		}
		var terminalState string
		if err := conn.QueryRow(`SELECT action FROM finding_dispositions WHERE finding_id=? AND action IN ('resolved','inapplicable')`, request.FindingID).Scan(&terminalState); err == nil {
			return &FindingAlreadyDisposedError{ID: request.FindingID, State: terminalState}
		} else if err != sql.ErrNoRows {
			return fmt.Errorf("load terminal Finding state: %w", err)
		}
		claim, err := loadAddressClaim(conn, request.ClaimID)
		if err != nil {
			return err
		}
		now := options.Now
		if now.IsZero() {
			now = time.Now().UTC()
		}
		if claim.Mode != "mutation" {
			return &workflow.ClaimModeInvalidError{ClaimID: claim.ID, Mode: claim.Mode, RequiredMode: "mutation"}
		}
		if claim.Status != "active" || claim.Version != request.ExpectedClaimVersion || claim.ActorID != options.ActorID || claim.IssueID != binding.IssueID || claim.EnvelopeID != binding.EnvelopeID || claim.EnvelopeRevision != binding.EnvelopeRevision || !now.Before(claim.ExpiresAt) {
			return staleAddressClaim(claim, request.ExpectedClaimVersion)
		}
		if err := workflow.AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, workflow.EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: "finding.address", Capability: "finding.address", IssueID: binding.IssueID, EnvelopeID: binding.EnvelopeID, Revision: binding.EnvelopeRevision}); err != nil {
			return err
		}
		actor, err := workflow.ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		dispositionID, err := nextDispositionID(conn)
		if err != nil {
			return err
		}
		next, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: request.FindingID, Type: "finding"})
		if err != nil {
			return err
		}
		envelope, err := workflow.GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, binding.EnvelopeID)
		if err != nil {
			return err
		}
		effect := journal.Effect{Op: "insert", Table: "finding_dispositions", PK: map[string]any{"id": dispositionID}, Values: map[string]any{"finding_id": request.FindingID, "action": "addressed", "actor_id": options.ActorID, "claim_id": request.ClaimID, "qualifying_review_id": nil, "reason": request.Reason, "disposed_at": "<occurred-at>", "created_event_id": "<event-id>"}}
		activity, err := reviewActivityEffect(conn, binding.IssueID, "finding-addressed", options.ActorID, map[string]any{"finding_id": request.FindingID})
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "finding.addressed", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: request.FindingID, Sequence: next, Type: "finding"}, CommandName: "finding.address", Repository: workflowJournalRepository(envelope.RepositoryScope), Effects: []journal.Effect{effect, activity}}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			for index := range prepared.Effects {
				setEffectIdentity(&prepared.Effects[index], prepared)
			}
			finding, err := loadFinding(conn, request.FindingID)
			if err != nil {
				return journal.IdempotencyMutation{}, err
			}
			disposition := map[string]any{"id": dispositionID, "action": "addressed", "actor_id": options.ActorID, "qualifying_review_id": nil, "reason": request.Reason, "disposed_at": prepared.OccurredAt.UTC().Format(time.RFC3339Nano), "event_id": prepared.EventID}
			finding["dispositions"] = append(finding["dispositions"].([]any), disposition)
			finding["state"] = "addressed"
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"finding": finding, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: request.FindingID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		finding, err := loadFinding(conn, request.FindingID)
		if err != nil {
			return err
		}
		result = AddressResult{Finding: finding, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func List(options ListOptions) ([]map[string]any, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}
	if options.Limit < 1 || options.Limit > 10000 || (options.State != "" && options.State != "open" && options.State != "addressed" && options.State != "resolved" && options.State != "inapplicable") || (options.Severity != "" && options.Severity != "blocker" && options.Severity != "major" && options.Severity != "minor" && options.Severity != "nit") {
		return nil, fmt.Errorf("finding list input is invalid")
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return nil, fmt.Errorf("open Finding projection: %w", err)
	}
	defer func() { _ = conn.Close() }()
	rows, err := conn.Query(`SELECT f.id FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE (?='' OR r.issue_id=?) AND (?='' OR f.review_id=?) AND (?='' OR f.severity=?) AND (?='' OR COALESCE((SELECT action FROM finding_dispositions d WHERE d.finding_id=f.id ORDER BY disposed_at DESC,id DESC LIMIT 1),'open')=?) ORDER BY f.id LIMIT ?`, options.IssueID, options.IssueID, options.ReviewID, options.ReviewID, options.Severity, options.Severity, options.State, options.State, options.Limit+1)
	if err != nil {
		return nil, fmt.Errorf("list Findings: %w", err)
	}
	defer func() { _ = rows.Close() }()
	items := []map[string]any{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		finding, err := loadFinding(conn, id)
		if err != nil {
			return nil, err
		}
		items = append(items, finding)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) > options.Limit {
		return nil, &ResultLimitExceededError{Limit: options.Limit, Observed: len(items)}
	}
	return items, nil
}

type findingBinding struct {
	IssueID, EnvelopeID string
	EnvelopeRevision    int
}

type addressClaim struct {
	ID, IssueID, EnvelopeID, ActorID, Mode, Status string
	EnvelopeRevision, Version                      int
	ExpiresAt                                      time.Time
}

func loadFindingBinding(conn *sql.DB, id string) (findingBinding, error) {
	var binding findingBinding
	err := conn.QueryRow(`SELECT r.issue_id,r.envelope_id,r.envelope_revision FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE f.id=?`, id).Scan(&binding.IssueID, &binding.EnvelopeID, &binding.EnvelopeRevision)
	if err == sql.ErrNoRows {
		return binding, &FindingNotFoundError{ID: id}
	}
	return binding, err
}

func loadAddressClaim(conn *sql.DB, id string) (addressClaim, error) {
	var claim addressClaim
	var expires string
	err := conn.QueryRow(`SELECT id,issue_id,envelope_id,envelope_revision,actor_id,mode,version,status,expires_at FROM claims WHERE id=?`, id).Scan(&claim.ID, &claim.IssueID, &claim.EnvelopeID, &claim.EnvelopeRevision, &claim.ActorID, &claim.Mode, &claim.Version, &claim.Status, &expires)
	if err == sql.ErrNoRows {
		return claim, &claims.ClaimNotFoundError{ID: id}
	}
	if err != nil {
		return claim, err
	}
	claim.ExpiresAt, err = time.Parse(time.RFC3339Nano, expires)
	return claim, err
}

func staleAddressClaim(claim addressClaim, expected int) error {
	actual, expiry := claim.Version, claim.ExpiresAt.Format(time.RFC3339Nano)
	return (*claims.ClaimStaleError)(&claims.ClaimConflictError{IssueID: claim.IssueID, ClaimID: &claim.ID, ExpectedVersion: &expected, ActualVersion: &actual, ExpiresAt: &expiry})
}

func loadFinding(conn *sql.DB, id string) (map[string]any, error) {
	var reviewID, subjectType, subjectID, scopeJSON, severity, summary, createdAt string
	var blocking bool
	err := conn.QueryRow(`SELECT review_id,subject_type,subject_id,scope_json,severity,summary,blocking,created_at FROM review_findings WHERE id=?`, id).Scan(&reviewID, &subjectType, &subjectID, &scopeJSON, &severity, &summary, &blocking, &createdAt)
	if err == sql.ErrNoRows {
		return nil, &FindingNotFoundError{ID: id}
	}
	if err != nil {
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
	return map[string]any{"id": id, "subject": map[string]any{"type": subjectType, "id": subjectID}, "scope": scope, "severity": severity, "summary": summary, "blocking": blocking, "review_id": reviewID, "created_at": createdAt, "dispositions": dispositions, "state": findingState(dispositions)}, nil
}

func nextDispositionID(conn *sql.DB) (string, error) {
	return journal.GenerateID("FDP", func(candidate string) (bool, error) {
		var found int
		err := conn.QueryRow(`SELECT 1 FROM finding_dispositions WHERE id=?`, candidate).Scan(&found)
		return err == nil, nil
	})
}

func replayAddress(outcome journal.IdempotencyOutcome, key string, result *AddressResult) error {
	var decoded struct {
		Finding map[string]any `json:"finding"`
	}
	if err := json.Unmarshal([]byte(outcome.ResultJSON), &decoded); err != nil {
		return err
	}
	*result = AddressResult{Finding: decoded.Finding, IdempotencyKey: key, Replayed: true, EventID: outcome.EventID}
	return nil
}

func setEffectIdentity(effect *journal.Effect, event journal.Event) {
	for key, value := range effect.Values {
		if value == "<event-id>" {
			effect.Values[key] = event.EventID
		}
		if value == "<occurred-at>" {
			effect.Values[key] = event.OccurredAt.UTC().Format(time.RFC3339Nano)
		}
	}
}

func workflowJournalRepository(repository workflow.Repository) *journal.Repository {
	base := ""
	if repository.BaseOID != nil {
		base = *repository.BaseOID
	}
	return &journal.Repository{BaseOID: base, CaptureSource: repository.CaptureSource, Dirty: repository.Dirty, DirtyDigest: repository.DirtyDigest, HeadOID: repository.HeadOID, ProjectID: repository.ProjectID, TreeOID: repository.TreeOID, WorktreeID: repository.WorktreeID}
}
