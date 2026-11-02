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
	runCloseAction   = "run.close"
	runAbandonAction = "run.abandon"
)

type RunCloseRequest struct {
	RunID              string                     `json:"run_id"`
	ExpectedVersion    int                        `json:"expected_version"`
	ObservationID      string                     `json:"observation_id"`
	ExpectedRepository project.ExpectedRepository `json:"expected_repository"`
}

type RunAbandonRequest struct {
	RunID                string                     `json:"run_id"`
	ExpectedVersion      int                        `json:"expected_version"`
	ExpectedClaimVersion int                        `json:"expected_claim_version"`
	ObservationID        string                     `json:"observation_id"`
	Authority            RecoveryAuthority          `json:"authority"`
	ExpectedRepository   project.ExpectedRepository `json:"expected_repository"`
}

type RunAbandonResult struct {
	Run            Run           `json:"run"`
	Claim          RecoveryClaim `json:"claim"`
	IdempotencyKey string        `json:"idempotency_key"`
	Replayed       bool          `json:"replayed"`
	EventID        string
}

// CloseRun derives a terminal Run state from trusted executor and Evidence facts.
func CloseRun(options RunOptions, request RunCloseRequest) (RunMutationResult, error) {
	if err := validateRunOptions(options); err != nil || !validRunCloseRequest(request) {
		return RunMutationResult{}, invalidRun("run close input is invalid")
	}
	return mutateRun(options, runCloseAction, request, func(conn *sql.DB, _ string, nowTime time.Time) (journal.Event, Run, error) {
		run, err := loadRun(conn, request.RunID)
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		if run.Status != "active" {
			return journal.Event{}, Run{}, &RunInactiveError{RunID: run.ID, Status: run.Status}
		}
		if run.Version != request.ExpectedVersion {
			return journal.Event{}, Run{}, &RunVersionConflictError{RunID: run.ID, ExpectedVersion: request.ExpectedVersion, ActualVersion: run.Version}
		}
		if err := requireTerminalObservation(options.WorkspaceRoot, conn, request.ObservationID, run); err != nil {
			return journal.Event{}, Run{}, err
		}
		if err := requireNoOutstandingAttempts(conn, run.ID); err != nil {
			return journal.Event{}, Run{}, err
		}
		if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runCloseAction, Capability: runCloseAction, IssueID: run.IssueID, EnvelopeID: run.EnvelopeID, Revision: run.EnvelopeRevision}); err != nil {
			return journal.Event{}, Run{}, err
		}
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		envelope, err := loadEnvelope(conn, run.EnvelopeID)
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		observed, err := observeRunTerminalRepository(options.WorkspaceRoot, request.ExpectedRepository)
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		reasons, err := evaluateObligations(conn, envelope.Obligations, observed)
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		status := "failed"
		if len(reasons) == 0 {
			status = "succeeded"
		}
		endedAt := nowTime.Format(time.RFC3339Nano)
		run.Status, run.TerminalRepository, run.EndedAt, run.Version = status, &observed, &endedAt, run.Version+1
		event, err := runEvent(options.WorkspaceRoot, actor, run, "run.closed", runCloseAction, []journal.Effect{
			{Op: "update", Table: "runs", PK: map[string]any{"id": run.ID}, Values: map[string]any{"status": run.Status, "terminal_repository_json": string(mustRunJSON(observed)), "ended_at": endedAt}, Preconditions: map[string]any{"status": "active"}},
			{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": run.ID}, Values: map[string]any{"version": run.Version}, Preconditions: map[string]any{"version": run.Version - 1}},
		})
		return event, run, err
	})
}

// AbandonRun records intentional cessation without creating Evidence or a successor Run.
func AbandonRun(options RunOptions, request RunAbandonRequest) (RunAbandonResult, error) {
	if err := validateRunOptions(options); err != nil || !validRunAbandonRequest(request) {
		return RunAbandonResult{}, invalidRun("run abandon input is invalid")
	}
	var result RunAbandonResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid " + runAbandonAction}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openWorkflowDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer closeWorkflowDB(conn)
		projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.ConfigFileName))
		if err != nil {
			return fmt.Errorf("run abandon: read project config: %w", err)
		}
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor_id": options.ActorID, "operation": runAbandonAction, "project_id": projectConfig.ProjectID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: runAbandonAction + ":" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			if err := json.Unmarshal([]byte(replay.ResultJSON), &result); err != nil {
				return fmt.Errorf("decode Run replay result: %w", err)
			}
			result.IdempotencyKey, result.Replayed, result.EventID = options.IdempotencyKey, true, replay.EventID
			return nil
		}

		run, err := loadRun(conn, request.RunID)
		if err != nil {
			return err
		}
		if run.Status != "active" {
			return &RunInactiveError{RunID: run.ID, Status: run.Status}
		}
		if run.Version != request.ExpectedVersion {
			return &RunVersionConflictError{RunID: run.ID, ExpectedVersion: request.ExpectedVersion, ActualVersion: run.Version}
		}
		observation, err := recoveryObservationAtWorkspace(options.WorkspaceRoot, conn, request.ObservationID, run)
		if err != nil {
			return err
		}
		if observation.Liveness.State == "unknown" {
			return &ExecutorLivenessUnknownError{RunID: run.ID, ExecutorID: run.Executor.ID, ObservationID: observation.ID}
		}
		if err := validateAbandonAuthority(options, run, request.Authority, observation); err != nil {
			return err
		}
		if request.Authority.Mode == "ordinary" && observation.Liveness.State != "terminated" {
			return &ExecutorStillLiveError{RunID: run.ID, ExecutorID: run.Executor.ID, ObservationID: observation.ID}
		}
		claim, err := loadRecoveryClaim(conn, run.ClaimID)
		if err != nil {
			return err
		}
		if claim.Version != request.ExpectedClaimVersion || claim.Status != "active" {
			return invalidRun("run abandon claim is invalid")
		}
		observed, err := observeRunTerminalRepository(options.WorkspaceRoot, request.ExpectedRepository)
		if err != nil {
			return err
		}
		if err := validateClaimModeAdoption(claim, observed.Dirty); err != nil {
			return err
		}
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		nowTime := runNow(options.Now)
		endedAt := nowTime.Format(time.RFC3339Nano)
		run.Status, run.TerminalRepository, run.EndedAt, run.Version = "abandoned", &observed, &endedAt, run.Version+1
		releasedAt := endedAt
		claim.Version, claim.Status, claim.ReleasedAt = claim.Version+1, "released", &releasedAt
		reason := "abandoned"
		claim.ReleaseReason = &reason
		event, err := runEvent(options.WorkspaceRoot, actor, run, "run.abandoned", runAbandonAction, []journal.Effect{
			{Op: "update", Table: "runs", PK: map[string]any{"id": run.ID}, Values: map[string]any{"status": run.Status, "terminal_repository_json": string(mustRunJSON(observed)), "ended_at": endedAt}, Preconditions: map[string]any{"status": "active"}},
			{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": run.ID}, Values: map[string]any{"version": run.Version}, Preconditions: map[string]any{"version": run.Version - 1}},
			{Op: "update", Table: "claims", PK: map[string]any{"id": claim.ID}, Values: map[string]any{"version": claim.Version, "status": claim.Status, "released_at": releasedAt, "release_reason": reason}, Preconditions: map[string]any{"status": "active", "version": claim.Version - 1}},
		})
		if err != nil {
			return err
		}
		resultJSON, err := journal.MarshalCanonicalJSON(map[string]any{"run": run, "claim": claim, "idempotency_key": options.IdempotencyKey, "replayed": false})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: run.ID, ResultJSON: string(resultJSON)}, nil
		})
		if err != nil {
			return err
		}
		loadedRun, err := loadRun(conn, outcome.EntityID)
		if err != nil {
			return err
		}
		loadedClaim, err := loadRecoveryClaim(conn, loadedRun.ClaimID)
		if err != nil {
			return err
		}
		result = RunAbandonResult{Run: loadedRun, Claim: loadedClaim, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func validRunCloseRequest(request RunCloseRequest) bool {
	return strings.TrimSpace(request.RunID) != "" && request.ExpectedVersion > 0 && strings.TrimSpace(request.ObservationID) != "" && validEnvelopeOID(request.ExpectedRepository.HeadOID)
}

func validRunAbandonRequest(request RunAbandonRequest) bool {
	if strings.TrimSpace(request.RunID) == "" || request.ExpectedVersion < 1 || request.ExpectedClaimVersion < 1 || strings.TrimSpace(request.ObservationID) == "" || !validEnvelopeOID(request.ExpectedRepository.HeadOID) || strings.TrimSpace(request.Authority.ActorID) == "" || strings.TrimSpace(request.Authority.Reason) == "" {
		return false
	}
	if request.Authority.Mode == "ordinary" {
		return request.Authority.Reference == nil && request.Authority.GrantRef == nil
	}
	return request.Authority.Mode == "force" && request.Authority.Reference != nil && strings.TrimSpace(*request.Authority.Reference) != "" && request.Authority.GrantRef != nil
}

func requireTerminalObservation(workspaceRoot string, conn *sql.DB, observationID string, run Run) error {
	observation, err := recoveryObservationAtWorkspace(workspaceRoot, conn, observationID, run)
	if err != nil {
		return err
	}
	if observation.Liveness.State == "unknown" {
		return &ExecutorLivenessUnknownError{RunID: run.ID, ExecutorID: run.Executor.ID, ObservationID: observation.ID}
	}
	if observation.Liveness.State != "terminated" {
		return &ExecutorStillLiveError{RunID: run.ID, ExecutorID: run.Executor.ID, ObservationID: observation.ID}
	}
	return nil
}

func requireNoOutstandingAttempts(conn *sql.DB, runID string) error {
	var found int
	err := conn.QueryRow(`SELECT 1 FROM execution_attempts a LEFT JOIN execution_receipts r ON r.attempt_id = a.id WHERE a.run_id = ? AND r.id IS NULL LIMIT 1`, runID).Scan(&found)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("load outstanding Run Attempts: %w", err)
	}
	return invalidRun("run close has outstanding Attempts")
}

func validateAbandonAuthority(options RunOptions, run Run, authority RecoveryAuthority, observation RunObservation) error {
	if authority.ActorID != options.ActorID {
		return invalidRun("abandon authority actor is invalid")
	}
	if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runAbandonAction, Capability: runAbandonAction, IssueID: run.IssueID, EnvelopeID: run.EnvelopeID, Revision: run.EnvelopeRevision}); err != nil {
		return err
	}
	if authority.Mode == "ordinary" {
		return nil
	}
	if observation.Liveness.State != "live" {
		return invalidRun("force abandon requires a live observation")
	}
	actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
	if err != nil {
		return err
	}
	if actor.ActorType != "human" {
		return invalidRun("force abandon requires a Human actor")
	}
	if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runForceRecoverAction, Capability: runForceRecoverAction, IssueID: run.IssueID, EnvelopeID: run.EnvelopeID, Revision: run.EnvelopeRevision}); err != nil {
		return err
	}
	envelope, err := GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, run.EnvelopeID)
	if err != nil {
		return err
	}
	ref := authority.GrantRef
	if ref.EnvelopeID != envelope.ID || ref.Revision != envelope.CurrentRevision || ref.Ordinal < 0 || ref.Ordinal >= len(envelope.AuthorityGrants) {
		return invalidRun("force abandon grant reference is invalid")
	}
	grant := envelope.AuthorityGrants[ref.Ordinal]
	if grant.ActorSelector.ActorID != options.ActorID || grant.Scope.IssueID != run.IssueID || !containsRecoveryAction(grant.Actions, runForceRecoverAction) || !containsRecoveryAction(grant.ActorSelector.RequiredCapabilities, runForceRecoverAction) {
		return invalidRun("force abandon grant reference is stale")
	}
	return nil
}

func observeRunTerminalRepository(workspaceRoot string, expected project.ExpectedRepository) (Repository, error) {
	current, err := resumeObservedRepository(workspaceRoot, workspaceRoot)
	if err != nil {
		return Repository{}, recoveryConflict(nil, err.Error())
	}
	if !expected.MatchesRecordedRepository(project.RecordedRepository(current)) {
		return Repository{}, &RunRepositoryMismatchError{Expected: expected, Observed: current}
	}
	return current, nil
}
