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
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/project"
)

const runCheckpointAction = "run.checkpoint"

type RunCheckpointRequest struct {
	RunID                string `json:"run_id"`
	ExpectedVersion      int    `json:"expected_version"`
	ReconciliationDigest string `json:"reconciliation_digest"`
}

type RunCheckpointRecord struct {
	ID             string         `json:"id"`
	RunID          string         `json:"run_id"`
	Kind           string         `json:"kind"`
	ActorID        string         `json:"actor_id"`
	ObservedAt     string         `json:"observed_at"`
	Reconciliation Reconciliation `json:"reconciliation"`
}

type RunCheckpointResult struct {
	Checkpoint     RunCheckpointRecord `json:"checkpoint"`
	Run            Run                 `json:"run"`
	IdempotencyKey string              `json:"idempotency_key"`
	Replayed       bool                `json:"replayed"`
	EventID        string
}

func RunCheckpoint(options RunOptions, request RunCheckpointRequest) (RunCheckpointResult, error) {
	if err := validateRunOptions(options); err != nil || !validRunCheckpointRequest(request) {
		return RunCheckpointResult{}, invalidRun("run checkpoint input is invalid")
	}
	var result RunCheckpointResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid " + runCheckpointAction}, func(lock *journal.WriteLock) error {
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
			return fmt.Errorf("run checkpoint: read project config: %w", err)
		}
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor_id": options.ActorID, "operation": runCheckpointAction, "project_id": projectConfig.ProjectID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: runCheckpointAction + ":" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
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
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		if actor.Producer != "activated-harness" {
			return &policy.EnvelopeAuthorityDeniedError{ActorID: options.ActorID, Action: runCheckpointAction, Capability: runCheckpointAction}
		}
		if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runCheckpointAction, Capability: runCheckpointAction, IssueID: run.IssueID, EnvelopeID: run.EnvelopeID, Revision: run.EnvelopeRevision}); err != nil {
			return err
		}
		checkpoint, err := latestRunCheckpoint(conn, run.ID)
		if err != nil {
			return err
		}
		base, checkpointID := run.StartRepository, (*string)(nil)
		if checkpoint != nil {
			base, checkpointID = checkpoint.Reconciliation.Repository, &checkpoint.ID
		}
		approvedPaths, err := runClaimPaths(conn, run.ClaimID)
		if err != nil {
			return err
		}
		current, err := project.ObserveRepository(options.WorkspaceRoot)
		if err != nil {
			return fmt.Errorf("run checkpoint: observe repository: %w", err)
		}
		reconciliation, err := ReconcileRepository(base, current, checkpointID, approvedPaths)
		if err != nil {
			return err
		}
		if reconciliation.Digest != request.ReconciliationDigest {
			return recoveryConflict(nil, "reconciliation digest is stale or changed")
		}
		checkpointIDValue, err := journal.GenerateID("OBS", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM run_observations WHERE id = ?`, candidate).Scan(&found)
			return err == nil, nil
		})
		if err != nil {
			return err
		}
		checkpoint = &RunCheckpointRecord{ID: checkpointIDValue, RunID: run.ID, Kind: "checkpoint", ActorID: actor.ActorID, ObservedAt: runNow(options.Now).Format(time.RFC3339Nano), Reconciliation: reconciliation}
		payload, err := journal.MarshalCanonicalJSON(checkpoint)
		if err != nil {
			return err
		}
		run.Version++
		event, err := runEvent(options.WorkspaceRoot, actor, run, "run.checkpointed", runCheckpointAction, []journal.Effect{
			{Op: "insert", Table: "run_observations", PK: map[string]any{"id": checkpoint.ID}, Values: map[string]any{"run_id": checkpoint.RunID, "kind": checkpoint.Kind, "actor_id": checkpoint.ActorID, "observed_at": checkpoint.ObservedAt, "observation_json": string(payload), "created_event_id": "<event-id>"}},
			{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": run.ID}, Values: map[string]any{"version": run.Version}, Preconditions: map[string]any{"version": run.Version - 1}},
		})
		if err != nil {
			return err
		}
		resultJSON, err := journal.MarshalCanonicalJSON(map[string]any{"checkpoint": checkpoint, "run": run, "idempotency_key": options.IdempotencyKey, "replayed": false})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			prepared.Effects[0].Values["created_event_id"] = prepared.EventID
			return journal.IdempotencyMutation{EntityID: checkpoint.ID, ResultJSON: string(resultJSON)}, nil
		})
		if err != nil {
			return err
		}
		loaded, err := loadRun(conn, run.ID)
		if err != nil {
			return err
		}
		result = RunCheckpointResult{Checkpoint: *checkpoint, Run: loaded, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func validRunCheckpointRequest(request RunCheckpointRequest) bool {
	return strings.TrimSpace(request.RunID) != "" && request.ExpectedVersion > 0 && strings.HasPrefix(request.ReconciliationDigest, "sha256:")
}

func loadRunCheckpoint(conn *sql.DB, id string) (RunCheckpointRecord, error) {
	var checkpoint RunCheckpointRecord
	var raw string
	err := conn.QueryRow(`SELECT id, run_id, kind, actor_id, observed_at, observation_json FROM run_observations WHERE id = ? AND kind = 'checkpoint'`, id).Scan(&checkpoint.ID, &checkpoint.RunID, &checkpoint.Kind, &checkpoint.ActorID, &checkpoint.ObservedAt, &raw)
	if err != nil {
		return RunCheckpointRecord{}, fmt.Errorf("load run checkpoint: %w", err)
	}
	if err := json.Unmarshal([]byte(raw), &checkpoint); err != nil {
		return RunCheckpointRecord{}, fmt.Errorf("decode run checkpoint: %w", err)
	}
	return checkpoint, nil
}

func latestRunCheckpoint(conn *sql.DB, runID string) (*RunCheckpointRecord, error) {
	var id string
	err := conn.QueryRow(`SELECT id FROM run_observations WHERE run_id = ? AND kind = 'checkpoint' ORDER BY observed_at DESC, id DESC LIMIT 1`, runID).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load latest run checkpoint: %w", err)
	}
	checkpoint, err := loadRunCheckpoint(conn, id)
	if err != nil {
		return nil, err
	}
	return &checkpoint, nil
}

func runClaimPaths(conn *sql.DB, claimID string) ([]string, error) {
	var scope struct {
		Paths []string `json:"paths"`
	}
	var raw string
	if err := conn.QueryRow(`SELECT scope_json FROM claims WHERE id = ?`, claimID).Scan(&raw); err != nil {
		return nil, fmt.Errorf("load Run Claim paths: %w", err)
	}
	if err := json.Unmarshal([]byte(raw), &scope); err != nil {
		return nil, fmt.Errorf("decode Run Claim paths: %w", err)
	}
	return scope.Paths, nil
}
