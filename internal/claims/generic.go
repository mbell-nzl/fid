package claims

import (
	"database/sql"
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

const (
	claimAcquireAction = "claim.acquire"
	claimExpireAction  = "claim.expire"
	claimRenewAction   = "claim.renew"
	claimReleaseAction = "claim.release"
)

type Scope struct {
	Paths []string `json:"paths"`
}

type AcquireRequest struct {
	IssueID          string `json:"issue_id"`
	EnvelopeID       string `json:"envelope_id"`
	EnvelopeRevision int    `json:"envelope_revision"`
	WorktreeID       string `json:"-"`
	Scope            Scope  `json:"scope"`
	Mode             string `json:"mode,omitempty"`
	TTLSeconds       int    `json:"ttl_seconds"`
}

type RenewRequest struct {
	ClaimID         string `json:"claim_id"`
	ExpectedVersion int    `json:"expected_version"`
	TTLSeconds      int    `json:"ttl_seconds"`
}

type ExpireRequest struct {
	ClaimID         string `json:"claim_id"`
	ExpectedVersion int    `json:"expected_version"`
}

type ReleaseRequest struct {
	ClaimID         string `json:"claim_id"`
	ExpectedVersion int    `json:"expected_version"`
	Reason          string `json:"reason"`
}

type MutationOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
	Now            time.Time
}

type GenericClaim struct {
	ID               string  `json:"id"`
	IssueID          string  `json:"issue_id"`
	EnvelopeID       string  `json:"envelope_id"`
	EnvelopeRevision int     `json:"envelope_revision"`
	ActorID          string  `json:"actor_id"`
	WorktreeID       string  `json:"worktree_id"`
	Scope            Scope   `json:"scope"`
	Mode             string  `json:"mode"`
	Version          int     `json:"version"`
	Status           string  `json:"status"`
	AcquiredAt       string  `json:"acquired_at"`
	ExpiresAt        string  `json:"expires_at"`
	ReleasedAt       *string `json:"released_at"`
	ReleaseReason    *string `json:"release_reason"`
}

type MutationResult struct {
	Claim          GenericClaim `json:"claim"`
	IdempotencyKey string       `json:"idempotency_key"`
	Replayed       bool         `json:"replayed"`
	EventID        string
}

type GenericListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	EnvelopeID    string
	Status        string
	Limit         int
	Now           time.Time
}

type ClaimConflictError struct {
	IssueID         string
	ClaimID         *string
	ExpectedVersion *int
	ActualVersion   *int
	ExpiresAt       *string
}

func (err *ClaimConflictError) Error() string { return "claim conflict" }

type ClaimExpiredError ClaimConflictError

func (err *ClaimExpiredError) Error() string { return "claim expired" }

type ClaimStaleError ClaimConflictError

func (err *ClaimStaleError) Error() string { return "claim stale" }

type ClaimNotFoundError struct{ ID string }

func (err *ClaimNotFoundError) Error() string { return fmt.Sprintf("claim %q was not found", err.ID) }

type GenericClaimInvalidError struct{ Reason string }

func (err *GenericClaimInvalidError) Error() string { return err.Reason }

func closeClaimsConn(conn *sql.DB) {
	_ = conn.Close()
}

func closeClaimsRows(rows *sql.Rows) {
	_ = rows.Close()
}

func AcquireGeneric(options MutationOptions, request AcquireRequest) (MutationResult, error) {
	if request.Mode == "" {
		request.Mode = "mutation"
	}
	if err := validateMutationOptions(options); err != nil || !validAcquireRequest(request) {
		return MutationResult{}, invalidGenericClaim("claim acquire input is invalid")
	}
	return mutateGenericClaim(options, claimAcquireAction, request, func(conn *sql.DB, _ string, now time.Time) (journal.Event, GenericClaim, error) {
		projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		observation, err := project.ObserveRecordedRepository(options.WorkspaceRoot, projectConfig.ProjectID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		repository := workflow.Repository(observation)
		envelope, err := validateAcquireAuthority(options, request, repository)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		blocked, err := claimAcquireStopBlock(conn, request)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		if blocked != nil {
			return journal.Event{}, GenericClaim{}, blocked
		}
		current, err := loadGenericActiveForIssue(conn, request.IssueID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		if current != nil {
			if current.expiredAt(now) {
				return journal.Event{}, GenericClaim{}, expiredGenericClaim(*current)
			}
			return journal.Event{}, GenericClaim{}, conflictGenericClaim(current)
		}
		prior, err := loadUnresolvedRunClaimForIssue(conn, request.IssueID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		if prior != nil {
			return journal.Event{}, GenericClaim{}, conflictGenericClaim(prior)
		}
		actor, err := workflow.ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		claimID, err := journal.GenerateID("CLM", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM claims WHERE id = ?`, candidate).Scan(&found)
			return err == nil, nil
		})
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		claim := GenericClaim{ID: claimID, IssueID: request.IssueID, EnvelopeID: request.EnvelopeID, EnvelopeRevision: request.EnvelopeRevision, ActorID: options.ActorID, WorktreeID: observation.WorktreeID, Scope: request.Scope, Mode: request.Mode, Version: 1, Status: "active", AcquiredAt: now.Format(time.RFC3339Nano), ExpiresAt: now.Add(time.Duration(request.TTLSeconds) * time.Second).Format(time.RFC3339Nano)}
		event, err := genericClaimEvent(options.WorkspaceRoot, actor, envelope.RepositoryScope, claimID, 1, "claim.acquired", claimAcquireAction, []journal.Effect{{
			Op: "insert", Table: "claims", PK: map[string]any{"id": claimID}, Values: claimValues(claim),
		}})
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		return event, claim, nil
	})
}

func claimAcquireStopBlock(conn *sql.DB, request AcquireRequest) (*workflow.StopBlockedError, error) {
	stops, err := workflow.ApplicableStops(conn, claimAcquireAction, workflow.Subject{Type: "envelope", ID: request.EnvelopeID})
	if err != nil || len(stops) == 0 {
		return nil, err
	}
	return workflow.NewStopBlockedError(conn, claimAcquireAction, stops)
}

func RenewGeneric(options MutationOptions, request RenewRequest) (MutationResult, error) {
	if err := validateMutationOptions(options); err != nil || !validRenewRequest(request) {
		return MutationResult{}, invalidGenericClaim("claim renew input is invalid")
	}
	return mutateGenericClaim(options, claimRenewAction, request, func(conn *sql.DB, _ string, now time.Time) (journal.Event, GenericClaim, error) {
		claim, err := loadGenericClaim(conn, request.ClaimID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		if err := validateStoredClaimAuthority(options, claim, claimRenewAction); err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		if claim.expiredAt(now) {
			return journal.Event{}, GenericClaim{}, expiredGenericClaim(claim)
		}
		if claim.Status != "active" || claim.Version != request.ExpectedVersion || claim.ActorID != options.ActorID {
			return journal.Event{}, GenericClaim{}, staleGenericClaim(claim, request.ExpectedVersion)
		}
		actor, err := workflow.ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		claim.Version++
		claim.ExpiresAt = now.Add(time.Duration(request.TTLSeconds) * time.Second).Format(time.RFC3339Nano)
		event, err := genericClaimEvent(options.WorkspaceRoot, actor, mustClaimEnvelope(options, claim), claim.ID, claim.Version, "claim.renewed", claimRenewAction, []journal.Effect{{
			Op: "update", Table: "claims", PK: map[string]any{"id": claim.ID}, Values: map[string]any{"version": claim.Version, "expires_at": claim.ExpiresAt}, Preconditions: map[string]any{"version": claim.Version - 1, "status": "active", "actor_id": options.ActorID},
		}})
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		return event, claim, nil
	})
}

// ExpireGeneric records an explicit lease expiry. It never changes the stored
// deadline, so expiry cannot later be revived by renewal.
func ExpireGeneric(options MutationOptions, request ExpireRequest) (MutationResult, error) {
	if err := validateMutationOptions(options); err != nil || !validExpireRequest(request) {
		return MutationResult{}, invalidGenericClaim("claim expire input is invalid")
	}
	return mutateGenericClaim(options, claimExpireAction, request, func(conn *sql.DB, _ string, now time.Time) (journal.Event, GenericClaim, error) {
		claim, err := loadGenericClaim(conn, request.ClaimID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		if err := validateStoredClaimAuthority(options, claim, claimExpireAction); err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		if claim.Status != "active" || claim.Version != request.ExpectedVersion || claim.ActorID != options.ActorID {
			return journal.Event{}, GenericClaim{}, staleGenericClaim(claim, request.ExpectedVersion)
		}
		actor, err := workflow.ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		claim.Version++
		claim.Status = "expired"
		event, err := genericClaimEvent(options.WorkspaceRoot, actor, mustClaimEnvelope(options, claim), claim.ID, claim.Version, "claim.expired", claimExpireAction, []journal.Effect{{
			Op: "update", Table: "claims", PK: map[string]any{"id": claim.ID},
			Values:        map[string]any{"version": claim.Version, "status": claim.Status},
			Preconditions: map[string]any{"version": claim.Version - 1, "status": "active", "actor_id": options.ActorID},
		}})
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		return event, claim, nil
	})
}

func ReleaseGeneric(options MutationOptions, request ReleaseRequest) (MutationResult, error) {
	if err := validateMutationOptions(options); err != nil || !validReleaseRequest(request) {
		return MutationResult{}, invalidGenericClaim("claim release input is invalid")
	}
	return mutateGenericClaim(options, claimReleaseAction, request, func(conn *sql.DB, _ string, now time.Time) (journal.Event, GenericClaim, error) {
		claim, err := loadGenericClaim(conn, request.ClaimID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		if err := validateStoredClaimAuthority(options, claim, claimReleaseAction); err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		if claim.expiredAt(now) {
			return journal.Event{}, GenericClaim{}, expiredGenericClaim(claim)
		}
		if claim.Status != "active" || claim.Version != request.ExpectedVersion || claim.ActorID != options.ActorID {
			return journal.Event{}, GenericClaim{}, staleGenericClaim(claim, request.ExpectedVersion)
		}
		actor, err := workflow.ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		claim.Version++
		claim.Status = "released"
		releasedAt := now.Format(time.RFC3339Nano)
		claim.ReleasedAt, claim.ReleaseReason = &releasedAt, &request.Reason
		event, err := genericClaimEvent(options.WorkspaceRoot, actor, mustClaimEnvelope(options, claim), claim.ID, claim.Version, "claim.released", claimReleaseAction, []journal.Effect{{
			Op: "update", Table: "claims", PK: map[string]any{"id": claim.ID}, Values: map[string]any{"version": claim.Version, "status": claim.Status, "released_at": releasedAt, "release_reason": request.Reason}, Preconditions: map[string]any{"version": claim.Version - 1, "status": "active", "actor_id": options.ActorID},
		}})
		if err != nil {
			return journal.Event{}, GenericClaim{}, err
		}
		return event, claim, nil
	})
}

func mutateGenericClaim(options MutationOptions, operation string, request any, build func(*sql.DB, string, time.Time) (journal.Event, GenericClaim, error)) (MutationResult, error) {
	var result MutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid " + operation}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openGenericClaimsDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer closeClaimsConn(conn)
		project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
		if err != nil {
			return fmt.Errorf("claim mutation: read project config: %w", err)
		}
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor_id": options.ActorID, "operation": operation, "project_id": project.ProjectID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: operation + ":" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			var stored struct {
				Claim GenericClaim `json:"claim"`
			}
			if err := json.Unmarshal([]byte(replay.ResultJSON), &stored); err != nil {
				return fmt.Errorf("decode Claim replay: %w", err)
			}
			result = MutationResult{Claim: stored.Claim, IdempotencyKey: options.IdempotencyKey, Replayed: true, EventID: replay.EventID}
			return nil
		}
		event, claim, err := build(conn, project.ProjectID, genericClaimNow(options.Now))
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"claim": claim, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: claim.ID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		claim, err = loadGenericClaim(conn, outcome.EntityID)
		if err != nil {
			return err
		}
		result = MutationResult{Claim: claim, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func validateAcquireAuthority(options MutationOptions, request AcquireRequest, observed workflow.Repository) (workflow.Envelope, error) {
	envelope, err := workflow.GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, request.EnvelopeID)
	if err != nil {
		return workflow.Envelope{}, err
	}
	if envelope.PrimaryIssueID != request.IssueID || envelope.CurrentRevision != request.EnvelopeRevision || !workflow.SameWorktree(options.WorkspaceRoot, envelope.RepositoryScope, observed) {
		return workflow.Envelope{}, invalidGenericClaim("claim acquire binding is invalid")
	}
	if err := workflow.AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, workflow.EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: claimAcquireAction, Capability: claimAcquireAction, IssueID: request.IssueID, EnvelopeID: request.EnvelopeID, Revision: request.EnvelopeRevision}); err != nil {
		return workflow.Envelope{}, err
	}
	return envelope, nil
}

func validateStoredClaimAuthority(options MutationOptions, claim GenericClaim, action string) error {
	envelope, err := workflow.GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, claim.EnvelopeID)
	if err != nil {
		return err
	}
	if envelope.RepositoryScope.WorktreeID != claim.WorktreeID {
		return invalidGenericClaim("claim worktree binding is invalid")
	}
	return workflow.AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, workflow.EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: action, Capability: action, IssueID: claim.IssueID, EnvelopeID: claim.EnvelopeID, Revision: claim.EnvelopeRevision})
}

func mustClaimEnvelope(options MutationOptions, claim GenericClaim) workflow.Repository {
	envelope, err := workflow.GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, claim.EnvelopeID)
	if err != nil {
		panic(err)
	}
	return envelope.RepositoryScope
}

func genericClaimEvent(workspaceRoot string, actor workflow.ActorProfile, repository workflow.Repository, claimID string, sequence int, eventType, command string, effects []journal.Effect) (journal.Event, error) {
	next, err := journal.NextAggregateSequence(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: claimID, Type: "claim"})
	if err != nil {
		return journal.Event{}, err
	}
	_ = sequence // Claim version and aggregate history diverge after atomic recovery creates the Claim under a Run Event.
	base := ""
	if repository.BaseOID != nil {
		base = *repository.BaseOID
	}
	return journal.Event{EventSchema: journal.EventSchemaV2, EventType: eventType, ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: claimID, Sequence: next, Type: "claim"}, CommandName: command, Repository: &journal.Repository{BaseOID: base, CaptureSource: repository.CaptureSource, Dirty: repository.Dirty, DirtyDigest: repository.DirtyDigest, HeadOID: repository.HeadOID, ProjectID: repository.ProjectID, TreeOID: repository.TreeOID, WorktreeID: repository.WorktreeID}, Effects: effects}, nil
}

func loadGenericClaim(conn *sql.DB, id string) (GenericClaim, error) {
	var claim GenericClaim
	var scopeJSON string
	var releasedAt, reason sql.NullString
	err := conn.QueryRow(`SELECT id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at, released_at, release_reason FROM claims WHERE id = ?`, id).Scan(&claim.ID, &claim.IssueID, &claim.EnvelopeID, &claim.EnvelopeRevision, &claim.ActorID, &claim.WorktreeID, &scopeJSON, &claim.Mode, &claim.Version, &claim.Status, &claim.AcquiredAt, &claim.ExpiresAt, &releasedAt, &reason)
	if err == sql.ErrNoRows {
		return GenericClaim{}, &ClaimNotFoundError{ID: id}
	}
	if err != nil {
		return GenericClaim{}, fmt.Errorf("load claim: %w", err)
	}
	if err := json.Unmarshal([]byte(scopeJSON), &claim.Scope); err != nil {
		return GenericClaim{}, fmt.Errorf("decode claim scope: %w", err)
	}
	if releasedAt.Valid {
		claim.ReleasedAt = &releasedAt.String
	}
	if reason.Valid {
		claim.ReleaseReason = &reason.String
	}
	return claim, nil
}

func loadGenericActiveForIssue(conn *sql.DB, issueID string) (*GenericClaim, error) {
	var id string
	err := conn.QueryRow(`SELECT id FROM claims WHERE issue_id = ? AND status = 'active'`, issueID).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load active claim: %w", err)
	}
	claim, err := loadGenericClaim(conn, id)
	return &claim, err
}

func loadUnresolvedRunClaimForIssue(conn *sql.DB, issueID string) (*GenericClaim, error) {
	var claimID string
	err := conn.QueryRow(`SELECT claim_id FROM runs WHERE issue_id = ? AND status = 'active' ORDER BY started_at, id LIMIT 1`, issueID).Scan(&claimID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load unresolved run: %w", err)
	}
	claim, err := loadGenericClaim(conn, claimID)
	if err != nil {
		return nil, err
	}
	return &claim, nil
}

func claimValues(claim GenericClaim) map[string]any {
	scope, err := journal.MarshalCanonicalJSON(claim.Scope)
	if err != nil {
		panic(err)
	}
	return map[string]any{"issue_id": claim.IssueID, "envelope_id": claim.EnvelopeID, "envelope_revision": claim.EnvelopeRevision, "actor_id": claim.ActorID, "worktree_id": claim.WorktreeID, "scope_json": string(scope), "mode": claim.Mode, "version": claim.Version, "status": claim.Status, "acquired_at": claim.AcquiredAt, "expires_at": claim.ExpiresAt, "released_at": nil, "release_reason": nil}
}

func (claim GenericClaim) expiredAt(now time.Time) bool {
	expiresAt, err := time.Parse(time.RFC3339Nano, claim.ExpiresAt)
	return err != nil || !now.Before(expiresAt.UTC())
}

func conflictGenericClaim(claim *GenericClaim) error {
	return &ClaimConflictError{IssueID: claim.IssueID, ClaimID: &claim.ID, ActualVersion: &claim.Version, ExpiresAt: &claim.ExpiresAt}
}

func expiredGenericClaim(claim GenericClaim) error {
	return &ClaimExpiredError{IssueID: claim.IssueID, ClaimID: &claim.ID, ActualVersion: &claim.Version, ExpiresAt: &claim.ExpiresAt}
}

func staleGenericClaim(claim GenericClaim, expected int) error {
	return &ClaimStaleError{IssueID: claim.IssueID, ClaimID: &claim.ID, ExpectedVersion: &expected, ActualVersion: &claim.Version, ExpiresAt: &claim.ExpiresAt}
}

func validateMutationOptions(options MutationOptions) error {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.IdempotencyKey) == "" {
		return invalidGenericClaim("claim mutation options are incomplete")
	}
	return nil
}

func validAcquireRequest(request AcquireRequest) bool {
	if strings.TrimSpace(request.IssueID) == "" || strings.TrimSpace(request.EnvelopeID) == "" || request.EnvelopeRevision != 1 || (request.Mode != "mutation" && request.Mode != "review") || request.TTLSeconds < 1 || request.TTLSeconds > 86400 || len(request.Scope.Paths) == 0 {
		return false
	}
	previous := ""
	for _, path := range request.Scope.Paths {
		if strings.TrimSpace(path) == "" || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "../") || strings.Contains(path, "/../") || previous >= path {
			return false
		}
		previous = path
	}
	return true
}

func validRenewRequest(request RenewRequest) bool {
	return strings.TrimSpace(request.ClaimID) != "" && request.ExpectedVersion > 0 && request.TTLSeconds >= 1 && request.TTLSeconds <= 86400
}

func validExpireRequest(request ExpireRequest) bool {
	return strings.TrimSpace(request.ClaimID) != "" && request.ExpectedVersion > 0
}

func validReleaseRequest(request ReleaseRequest) bool {
	return strings.TrimSpace(request.ClaimID) != "" && request.ExpectedVersion > 0 && (request.Reason == "completed" || request.Reason == "failed" || request.Reason == "operator-released")
}

func genericClaimNow(now time.Time) time.Time {
	if now.IsZero() {
		return time.Now().UTC()
	}
	return now.UTC()
}

func invalidGenericClaim(reason string) error { return &GenericClaimInvalidError{Reason: reason} }

func openGenericClaimsDB(workspaceRoot string) (*sql.DB, error) {
	return fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
}

func ListGeneric(options GenericListOptions) ([]GenericClaim, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || options.Limit < 1 || options.Limit > 10000 || (options.Status != "" && options.Status != "active" && options.Status != "expired" && options.Status != "released") {
		return nil, invalidGenericClaim("claim list input is invalid")
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := openGenericClaimsDB(options.WorkspaceRoot)
	if err != nil {
		return nil, err
	}
	defer closeClaimsConn(conn)
	query := `SELECT id FROM claims WHERE 1 = 1`
	args := []any{}
	if options.IssueID != "" {
		query += ` AND issue_id = ?`
		args = append(args, options.IssueID)
	}
	if options.EnvelopeID != "" {
		query += ` AND envelope_id = ?`
		args = append(args, options.EnvelopeID)
	}
	query += ` ORDER BY issue_id, id LIMIT ?`
	args = append(args, options.Limit)
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list claims: %w", err)
	}
	defer closeClaimsRows(rows)
	claims := []GenericClaim{}
	now := genericClaimNow(options.Now)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		claim, err := loadGenericClaim(conn, id)
		if err != nil {
			return nil, err
		}
		if claim.Status == "active" && claim.expiredAt(now) {
			claim.Status = "expired"
		}
		if options.Status == "" || options.Status == claim.Status {
			claims = append(claims, claim)
		}
	}
	return claims, rows.Err()
}
