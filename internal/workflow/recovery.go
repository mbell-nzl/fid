package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
)

const (
	runRecoverAction      = "run.recover"
	runForceRecoverAction = "run.force-recover"
)

type RecoveryGrantRef struct {
	EnvelopeID string `json:"envelope_id"`
	Revision   int    `json:"revision"`
	Ordinal    int    `json:"ordinal"`
}

type RecoveryAuthority struct {
	Mode      string            `json:"mode"`
	ActorID   string            `json:"actor_id"`
	Reason    string            `json:"reason"`
	Reference *string           `json:"reference"`
	GrantRef  *RecoveryGrantRef `json:"grant_ref"`
}

type RunRecoveryRequest struct {
	PriorRunID           string            `json:"prior_run_id"`
	ExpectedRunVersion   int               `json:"expected_run_version"`
	ExpectedClaimVersion int               `json:"expected_claim_version"`
	ObservationID        string            `json:"observation_id"`
	ReconciliationDigest string            `json:"reconciliation_digest"`
	AdoptedPaths         []string          `json:"adopted_paths"`
	Authority            RecoveryAuthority `json:"authority"`
	Executor             Executor          `json:"executor"`
	TTLSeconds           int               `json:"ttl_seconds"`
}

type RecoveryClaim struct {
	ID            string          `json:"id"`
	IssueID       string          `json:"issue_id"`
	EnvelopeID    string          `json:"envelope_id"`
	Revision      int             `json:"envelope_revision"`
	ActorID       string          `json:"actor_id"`
	WorktreeID    string          `json:"worktree_id"`
	ScopeJSON     string          `json:"-"`
	Scope         json.RawMessage `json:"scope"`
	Mode          string          `json:"mode"`
	Version       int             `json:"version"`
	Status        string          `json:"status"`
	AcquiredAt    string          `json:"acquired_at"`
	ExpiresAt     string          `json:"expires_at"`
	ReleasedAt    *string         `json:"released_at"`
	ReleaseReason *string         `json:"release_reason"`
}

type Recovery struct {
	ID               string            `json:"id"`
	PriorRunID       string            `json:"prior_run_id"`
	SuccessorRunID   string            `json:"successor_run_id"`
	PriorClaimID     string            `json:"prior_claim_id"`
	SuccessorClaimID string            `json:"successor_claim_id"`
	ObservationID    string            `json:"observation_id"`
	Reconciliation   Reconciliation    `json:"reconciliation"`
	AdoptedPaths     []string          `json:"adopted_paths"`
	Authority        RecoveryAuthority `json:"authority"`
	EventID          string            `json:"event_id"`
}

type RunRecoveryResult struct {
	Recovery       Recovery      `json:"recovery"`
	Run            Run           `json:"run"`
	Claim          RecoveryClaim `json:"claim"`
	IdempotencyKey string        `json:"idempotency_key"`
	Replayed       bool          `json:"replayed"`
	EventID        string
}

type ExecutorLivenessUnknownError struct {
	RunID, ExecutorID, ObservationID string
}

func (err *ExecutorLivenessUnknownError) Error() string { return "executor liveness is unknown" }

type ExecutorStillLiveError struct {
	RunID, ExecutorID, ObservationID string
}

func (err *ExecutorStillLiveError) Error() string { return "executor is still live" }

type ObservationNotFoundError struct{ ID string }

func (err *ObservationNotFoundError) Error() string {
	return fmt.Sprintf("observation %q was not found", err.ID)
}

type RecoveryAlreadyLinkedError struct {
	RunID, SuccessorRunID string
}

func (err *RecoveryAlreadyLinkedError) Error() string { return "recovery already linked" }

// RecoverRun atomically fences a predecessor and creates its one linked successor.
func RecoverRun(options RunOptions, request RunRecoveryRequest) (RunRecoveryResult, error) {
	if err := validateRunOptions(options); err != nil || !validRunRecoveryRequest(request) {
		return RunRecoveryResult{}, invalidRun("run recovery input is invalid")
	}
	var result RunRecoveryResult
	action := runRecoverAction
	if request.Authority.Mode == "force" {
		action = runForceRecoverAction
	}
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid " + action}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openWorkflowDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer closeWorkflowDB(conn)
		projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
		if err != nil {
			return fmt.Errorf("run recovery: read project config: %w", err)
		}
		digest, err := runRequestDigest(action, projectConfig.ProjectID, options.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: action + ":" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			if err := json.Unmarshal([]byte(replay.ResultJSON), &result); err != nil {
				return fmt.Errorf("decode Run replay result: %w", err)
			}
			result.IdempotencyKey, result.Replayed, result.EventID = options.IdempotencyKey, true, replay.EventID
			return nil
		}
		prior, err := loadRun(conn, request.PriorRunID)
		if err != nil {
			return err
		}
		stops, err := ApplicableStops(conn, runRecoverAction, Subject{Type: "run", ID: prior.ID}, Subject{Type: "envelope", ID: prior.EnvelopeID})
		if err != nil {
			return err
		}
		if len(stops) != 0 {
			return &StopBlockedError{Action: runRecoverAction, Stops: stops}
		}
		if prior.Version != request.ExpectedRunVersion {
			return &RunVersionConflictError{RunID: prior.ID, ExpectedVersion: request.ExpectedRunVersion, ActualVersion: prior.Version}
		}
		if successor, err := existingRecoverySuccessor(conn, prior.ID); err != nil {
			return err
		} else if successor != "" {
			return &RecoveryAlreadyLinkedError{RunID: prior.ID, SuccessorRunID: successor}
		}
		priorClaim, err := loadRecoveryClaim(conn, prior.ClaimID)
		if err != nil {
			return err
		}
		if priorClaim.Version != request.ExpectedClaimVersion {
			return invalidRun("run recovery claim version is invalid")
		}
		observation, err := recoveryObservationAtWorkspace(options.WorkspaceRoot, conn, request.ObservationID, prior)
		if err != nil {
			return err
		}
		if observation.Liveness.State == "unknown" {
			return &ExecutorLivenessUnknownError{RunID: prior.ID, ExecutorID: prior.Executor.ID, ObservationID: observation.ID}
		}
		if err := validateRecoveryAuthority(options, prior, request.Authority, observation); err != nil {
			return err
		}
		if prior.Status == "succeeded" || prior.Status == "abandoned" || (prior.Status != "active" && prior.Status != "failed" && prior.Status != "interrupted") {
			return invalidRun("run recovery predecessor is not recoverable")
		}
		if request.Authority.Mode == "ordinary" && (observation.Liveness.State != "terminated" || priorClaim.Status == "active") {
			if observation.Liveness.State == "live" {
				return &ExecutorStillLiveError{RunID: prior.ID, ExecutorID: prior.Executor.ID, ObservationID: observation.ID}
			}
			return invalidRun("ordinary recovery requires a released or expired claim")
		}
		current, err := project.ObserveRepository(options.WorkspaceRoot)
		if err != nil {
			return fmt.Errorf("run recovery: observe repository: %w", err)
		}
		approvedPaths, err := runClaimPaths(conn, prior.ClaimID)
		if err != nil {
			return err
		}
		checkpoint, err := latestRunCheckpoint(conn, prior.ID)
		if err != nil {
			return err
		}
		base, checkpointID := prior.StartRepository, (*string)(nil)
		if checkpoint != nil {
			base, checkpointID = checkpoint.Reconciliation.Repository, &checkpoint.ID
		}
		reconciliation, err := reconcileRunRepository(options.WorkspaceRoot, base, current, checkpointID, approvedPaths)
		if err != nil {
			return err
		}
		if err := validateReconciliationAdoption(reconciliation, request.AdoptedPaths, request.ReconciliationDigest, request.Authority.Mode == "force"); err != nil {
			return err
		}
		if err := validateClaimModeAdoption(priorClaim, len(request.AdoptedPaths) != 0); err != nil {
			return err
		}
		if err := validateRecoveryExecutor(conn, request.Executor); err != nil {
			return err
		}
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		now := runNow(options.Now)
		successorClaimID, err := uniqueRecoveryID(conn, "CLM", "claims")
		if err != nil {
			return err
		}
		successorRunID, err := uniqueRecoveryID(conn, "RUN", "runs")
		if err != nil {
			return err
		}
		recoveryID, err := uniqueRecoveryID(conn, "REC", "run_recoveries")
		if err != nil {
			return err
		}
		successorClaim := RecoveryClaim{ID: successorClaimID, IssueID: priorClaim.IssueID, EnvelopeID: priorClaim.EnvelopeID, Revision: priorClaim.Revision, ActorID: options.ActorID, WorktreeID: reconciliation.Repository.WorktreeID, ScopeJSON: priorClaim.ScopeJSON, Scope: append(json.RawMessage(nil), priorClaim.Scope...), Mode: priorClaim.Mode, Version: 1, Status: "active", AcquiredAt: now.Format(time.RFC3339Nano), ExpiresAt: now.Add(time.Duration(request.TTLSeconds) * time.Second).Format(time.RFC3339Nano)}
		successor := Run{ID: successorRunID, IssueID: prior.IssueID, EnvelopeID: prior.EnvelopeID, EnvelopeRevision: prior.EnvelopeRevision, ClaimID: successorClaim.ID, ClaimVersion: successorClaim.Version, Version: 1, Executor: request.Executor, ObligationID: prior.ObligationID, ActorID: options.ActorID, StartRepository: reconciliation.Repository, Status: "active", Budgets: prior.Budgets, StartedAt: now.Format(time.RFC3339Nano), PredecessorRunID: &prior.ID}
		recovery := Recovery{ID: recoveryID, PriorRunID: prior.ID, SuccessorRunID: successor.ID, PriorClaimID: priorClaim.ID, SuccessorClaimID: successorClaim.ID, ObservationID: observation.ID, Reconciliation: reconciliation, AdoptedPaths: append([]string{}, request.AdoptedPaths...), Authority: request.Authority}
		budgetEffects, err := recoveryInterventionEffects(conn, &successor, recovery)
		if err != nil {
			return err
		}
		effects := append(recoveryEffects(prior, priorClaim, successor, successorClaim, recovery, now), budgetEffects...)
		sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: prior.ID, Type: "run"})
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "run.recovered", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: prior.ID, Type: "run", Sequence: sequence}, CommandName: action, Repository: envelopeRepository(reconciliation.Repository), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			for index := range prepared.Effects {
				if prepared.Effects[index].Values["created_event_id"] == "<event-id>" {
					prepared.Effects[index].Values["created_event_id"] = prepared.EventID
				}
			}
			recovery.EventID = prepared.EventID
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"recovery": recovery, "run": successor, "claim": successorClaim, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: recovery.ID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		loadedRecovery, err := loadRecovery(conn, outcome.EntityID)
		if err != nil {
			return err
		}
		loadedRun, err := loadRun(conn, successor.ID)
		if err != nil {
			return err
		}
		loadedClaim, err := loadRecoveryClaim(conn, successorClaim.ID)
		if err != nil {
			return err
		}
		result = RunRecoveryResult{Recovery: loadedRecovery, Run: loadedRun, Claim: loadedClaim, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func validRunRecoveryRequest(request RunRecoveryRequest) bool {
	if strings.TrimSpace(request.PriorRunID) == "" || request.ExpectedRunVersion < 1 || request.ExpectedClaimVersion < 1 || strings.TrimSpace(request.ObservationID) == "" || strings.TrimSpace(request.ReconciliationDigest) == "" || request.TTLSeconds < 1 || request.TTLSeconds > 86400 || strings.TrimSpace(request.Executor.ID) == "" || request.Executor.Source != "activated-harness" || request.Authority.ActorID == "" || request.Authority.Reason == "" || (request.Authority.Mode != "ordinary" && request.Authority.Mode != "force") {
		return false
	}
	if request.Authority.Mode == "ordinary" {
		return request.Authority.Reference == nil && request.Authority.GrantRef == nil
	}
	return request.Authority.Reference != nil && strings.TrimSpace(*request.Authority.Reference) != "" && request.Authority.GrantRef != nil
}

func validateRecoveryAuthority(options RunOptions, prior Run, authority RecoveryAuthority, observation RunObservation) error {
	if authority.ActorID != options.ActorID {
		return invalidRun("recovery authority actor is invalid")
	}
	if authority.Mode == "ordinary" {
		return AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runRecoverAction, Capability: runRecoverAction, IssueID: prior.IssueID, EnvelopeID: prior.EnvelopeID, Revision: prior.EnvelopeRevision})
	}
	if observation.Liveness.State != "live" {
		return invalidRun("force recovery requires a live observation")
	}
	actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
	if err != nil {
		return err
	}
	if actor.ActorType != "human" {
		return invalidRun("force recovery requires a Human actor")
	}
	if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runRecoverAction, Capability: runRecoverAction, IssueID: prior.IssueID, EnvelopeID: prior.EnvelopeID, Revision: prior.EnvelopeRevision}); err != nil {
		return err
	}
	if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runForceRecoverAction, Capability: runForceRecoverAction, IssueID: prior.IssueID, EnvelopeID: prior.EnvelopeID, Revision: prior.EnvelopeRevision}); err != nil {
		return err
	}
	envelope, err := GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, prior.EnvelopeID)
	if err != nil {
		return err
	}
	ref := authority.GrantRef
	if ref.EnvelopeID != envelope.ID || ref.Revision != envelope.CurrentRevision || ref.Ordinal < 0 || ref.Ordinal >= len(envelope.AuthorityGrants) {
		return invalidRun("force recovery grant reference is invalid")
	}
	grant := envelope.AuthorityGrants[ref.Ordinal]
	if grant.ActorSelector.ActorID != options.ActorID || grant.Scope.IssueID != prior.IssueID || !containsRecoveryAction(grant.ActorSelector.RequiredCapabilities, runForceRecoverAction) || !containsRecoveryAction(grant.Actions, runForceRecoverAction) {
		return invalidRun("force recovery grant reference is stale")
	}
	return nil
}

func containsRecoveryAction(actions []string, action string) bool {
	for _, value := range actions {
		if value == action {
			return true
		}
	}
	return false
}

func recoveryObservation(conn *sql.DB, id string, run Run) (RunObservation, error) {
	var found int
	if err := conn.QueryRow(`SELECT 1 FROM run_observations WHERE id = ?`, id).Scan(&found); err == sql.ErrNoRows {
		return RunObservation{}, &ObservationNotFoundError{ID: id}
	} else if err != nil {
		return RunObservation{}, fmt.Errorf("find recovery observation: %w", err)
	}
	observation, err := loadRunObservation(conn, id)
	if err != nil {
		return RunObservation{}, err
	}
	if observation.Kind != "liveness" || observation.RunID != run.ID || observation.Liveness.ExecutorID != run.Executor.ID {
		return RunObservation{}, invalidRun("recovery observation binding is invalid")
	}
	return observation, nil
}

func recoveryObservationAtWorkspace(workspaceRoot string, conn *sql.DB, id string, run Run) (RunObservation, error) {
	observation, err := recoveryObservation(conn, id, run)
	if err != nil {
		return RunObservation{}, err
	}
	if importedObservation(workspaceRoot, conn, id) {
		return RunObservation{}, invalidRun("recovery observation is imported history")
	}
	return observation, nil
}

func validateRecoveryExecutor(conn *sql.DB, executor Executor) error {
	var found string
	err := conn.QueryRow(`SELECT run_id FROM run_runtime WHERE executor_id = ?`, executor.ID).Scan(&found)
	if err == nil {
		return &DuplicateRunExecutorError{ExecutorID: executor.ID}
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("load recovery executor: %w", err)
	}
	return nil
}

func uniqueRecoveryID(conn *sql.DB, prefix, table string) (string, error) {
	return journal.GenerateID(prefix, func(candidate string) (bool, error) {
		var found int
		err := conn.QueryRow(`SELECT 1 FROM `+table+` WHERE id = ?`, candidate).Scan(&found)
		if err == sql.ErrNoRows {
			return false, nil
		}
		return err == nil, err
	})
}

func existingRecoverySuccessor(conn *sql.DB, priorID string) (string, error) {
	var successor string
	err := conn.QueryRow(`SELECT successor_run_id FROM run_recoveries WHERE prior_run_id = ?`, priorID).Scan(&successor)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return successor, err
}

func recoveryEffects(prior Run, priorClaim RecoveryClaim, successor Run, successorClaim RecoveryClaim, recovery Recovery, now time.Time) []journal.Effect {
	effects := []journal.Effect{}
	if prior.Status == "active" {
		endedAt := now.Format(time.RFC3339Nano)
		effects = append(effects, journal.Effect{Op: "update", Table: "runs", PK: map[string]any{"id": prior.ID}, Values: map[string]any{"status": "interrupted", "terminal_repository_json": string(mustRunJSON(recovery.Reconciliation.Repository)), "ended_at": endedAt}, Preconditions: map[string]any{"status": "active"}})
		effects = append(effects, journal.Effect{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": prior.ID}, Values: map[string]any{"version": prior.Version + 1}, Preconditions: map[string]any{"version": prior.Version}})
	}
	if priorClaim.Status == "active" {
		releasedAt := now.Format(time.RFC3339Nano)
		effects = append(effects, journal.Effect{Op: "update", Table: "claims", PK: map[string]any{"id": priorClaim.ID}, Values: map[string]any{"version": priorClaim.Version + 1, "status": "superseded", "released_at": releasedAt, "release_reason": "recovery-superseded"}, Preconditions: map[string]any{"status": "active", "version": priorClaim.Version}})
	}
	effects = append(effects, journal.Effect{Op: "insert", Table: "claims", PK: map[string]any{"id": successorClaim.ID}, Values: recoveryClaimValues(successorClaim)})
	effects = append(effects, runStartEffects(successor)...)
	reconciliationJSON, _ := journal.MarshalCanonicalJSON(recovery.Reconciliation)
	adoptedJSON, _ := journal.MarshalCanonicalJSON(recovery.AdoptedPaths)
	authorityJSON, _ := journal.MarshalCanonicalJSON(recovery.Authority)
	effects = append(effects, journal.Effect{Op: "insert", Table: "run_recoveries", PK: map[string]any{"id": recovery.ID}, Values: map[string]any{"prior_run_id": recovery.PriorRunID, "successor_run_id": recovery.SuccessorRunID, "prior_claim_id": recovery.PriorClaimID, "successor_claim_id": recovery.SuccessorClaimID, "observation_id": recovery.ObservationID, "reconciliation_json": string(reconciliationJSON), "adopted_paths_json": string(adoptedJSON), "authority_json": string(authorityJSON), "reason": recovery.Authority.Reason, "created_event_id": "<event-id>"}})
	return effects
}

func recoveryClaimValues(claim RecoveryClaim) map[string]any {
	return map[string]any{"issue_id": claim.IssueID, "envelope_id": claim.EnvelopeID, "envelope_revision": claim.Revision, "actor_id": claim.ActorID, "worktree_id": claim.WorktreeID, "scope_json": claim.ScopeJSON, "mode": claim.Mode, "version": claim.Version, "status": claim.Status, "acquired_at": claim.AcquiredAt, "expires_at": claim.ExpiresAt, "released_at": nil, "release_reason": nil}
}

func loadRecoveryClaim(conn *sql.DB, id string) (RecoveryClaim, error) {
	var claim RecoveryClaim
	var releasedAt, reason sql.NullString
	err := conn.QueryRow(`SELECT id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at, released_at, release_reason FROM claims WHERE id = ?`, id).Scan(&claim.ID, &claim.IssueID, &claim.EnvelopeID, &claim.Revision, &claim.ActorID, &claim.WorktreeID, &claim.ScopeJSON, &claim.Mode, &claim.Version, &claim.Status, &claim.AcquiredAt, &claim.ExpiresAt, &releasedAt, &reason)
	if err != nil {
		return RecoveryClaim{}, err
	}
	if !json.Valid([]byte(claim.ScopeJSON)) {
		return RecoveryClaim{}, fmt.Errorf("decode recovery Claim scope: invalid JSON")
	}
	claim.Scope = append(json.RawMessage(nil), claim.ScopeJSON...)
	if releasedAt.Valid {
		claim.ReleasedAt = &releasedAt.String
	}
	if reason.Valid {
		claim.ReleaseReason = &reason.String
	}
	return claim, nil
}

func loadRecovery(conn *sql.DB, id string) (Recovery, error) {
	var recovery Recovery
	var reconciliation, adopted, authority string
	err := conn.QueryRow(`SELECT id, prior_run_id, successor_run_id, prior_claim_id, successor_claim_id, observation_id, reconciliation_json, adopted_paths_json, authority_json, created_event_id FROM run_recoveries WHERE id = ?`, id).Scan(&recovery.ID, &recovery.PriorRunID, &recovery.SuccessorRunID, &recovery.PriorClaimID, &recovery.SuccessorClaimID, &recovery.ObservationID, &reconciliation, &adopted, &authority, &recovery.EventID)
	if err != nil {
		return Recovery{}, err
	}
	if err := json.Unmarshal([]byte(reconciliation), &recovery.Reconciliation); err != nil {
		return Recovery{}, err
	}
	if err := json.Unmarshal([]byte(adopted), &recovery.AdoptedPaths); err != nil {
		return Recovery{}, err
	}
	if err := json.Unmarshal([]byte(authority), &recovery.Authority); err != nil {
		return Recovery{}, err
	}
	return recovery, nil
}
