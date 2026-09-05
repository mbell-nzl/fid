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
)

const runObserveAction = "run.observe"

type Liveness struct {
	ExecutorID string `json:"executor_id"`
	State      string `json:"state"`
	SourceRef  string `json:"source_ref"`
	ObservedAt string `json:"observed_at"`
}

type RunObserveRequest struct {
	RunID           string   `json:"run_id"`
	ExpectedVersion int      `json:"expected_version"`
	Liveness        Liveness `json:"liveness"`
}

type RunObservation struct {
	ID         string   `json:"id"`
	RunID      string   `json:"run_id"`
	Kind       string   `json:"kind"`
	ActorID    string   `json:"actor_id"`
	ObservedAt string   `json:"observed_at"`
	Liveness   Liveness `json:"liveness"`
}

type RunObservationResult struct {
	Observation    RunObservation `json:"observation"`
	Run            Run            `json:"run"`
	IdempotencyKey string         `json:"idempotency_key"`
	Replayed       bool           `json:"replayed"`
	EventID        string
}

func ObserveRun(options RunOptions, request RunObserveRequest) (RunObservationResult, error) {
	if err := validateRunOptions(options); err != nil || !validRunObserveRequest(request) {
		return RunObservationResult{}, invalidRun("run observe input is invalid")
	}
	var result RunObservationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid " + runObserveAction}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openWorkflowDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer closeWorkflowDB(conn)
		project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
		if err != nil {
			return fmt.Errorf("run observe: read project config: %w", err)
		}
		digest, err := runRequestDigest(runObserveAction, project.ProjectID, options.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: runObserveAction + ":" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			observation, err := loadRunObservation(conn, replay.EntityID)
			if err != nil {
				return err
			}
			run, err := loadRun(conn, observation.RunID)
			if err != nil {
				return err
			}
			result = RunObservationResult{Observation: observation, Run: run, IdempotencyKey: options.IdempotencyKey, Replayed: true, EventID: replay.EventID}
			return nil
		}
		now := runNow(options.Now)
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
		if run.Executor.ID != request.Liveness.ExecutorID {
			return invalidRun("run observation executor binding is invalid")
		}
		observedAt, err := time.Parse(time.RFC3339Nano, request.Liveness.ObservedAt)
		if err != nil || observedAt.After(now) {
			return invalidRun("run observation time is invalid")
		}
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		if actor.Producer != "activated-harness" {
			return &policy.EnvelopeAuthorityDeniedError{ActorID: options.ActorID, Action: runObserveAction, Capability: runObserveAction}
		}
		if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runObserveAction, Capability: runObserveAction, IssueID: run.IssueID, EnvelopeID: run.EnvelopeID, Revision: run.EnvelopeRevision}); err != nil {
			return err
		}
		if prior, err := latestLiveness(conn, run.ID); err != nil {
			return err
		} else if prior != nil && prior.State == "terminated" && request.Liveness.State != "terminated" {
			return invalidRun("terminated executor cannot become live")
		}
		observationID, err := journal.GenerateID("OBS", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM run_observations WHERE id = ?`, candidate).Scan(&found)
			return err == nil, nil
		})
		if err != nil {
			return err
		}
		observation := RunObservation{ID: observationID, RunID: run.ID, Kind: "liveness", ActorID: actor.ActorID, ObservedAt: request.Liveness.ObservedAt, Liveness: request.Liveness}
		payload, err := journal.MarshalCanonicalJSON(observation)
		if err != nil {
			return err
		}
		run.Version++
		event, err := runEvent(options.WorkspaceRoot, actor, run, "run.observed", runObserveAction, []journal.Effect{
			{Op: "insert", Table: "run_observations", PK: map[string]any{"id": observation.ID}, Values: map[string]any{"run_id": observation.RunID, "kind": observation.Kind, "actor_id": observation.ActorID, "observed_at": observation.ObservedAt, "observation_json": string(payload), "created_event_id": "<event-id>"}},
			{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": run.ID}, Values: map[string]any{"version": run.Version}, Preconditions: map[string]any{"version": run.Version - 1}},
		})
		if err != nil {
			return err
		}
		resultJSON, err := journal.MarshalCanonicalJSON(map[string]any{"observation": observation, "run": run, "idempotency_key": options.IdempotencyKey, "replayed": false})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			prepared.Effects[0].Values["created_event_id"] = prepared.EventID
			return journal.IdempotencyMutation{EntityID: observation.ID, ResultJSON: string(resultJSON)}, nil
		})
		if err != nil {
			return err
		}
		loaded, err := loadRun(conn, run.ID)
		if err != nil {
			return err
		}
		result = RunObservationResult{Observation: observation, Run: loaded, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func validRunObserveRequest(request RunObserveRequest) bool {
	if strings.TrimSpace(request.RunID) == "" || request.ExpectedVersion < 1 || strings.TrimSpace(request.Liveness.ExecutorID) == "" || strings.TrimSpace(request.Liveness.SourceRef) == "" || (request.Liveness.State != "live" && request.Liveness.State != "terminated" && request.Liveness.State != "unknown") {
		return false
	}
	_, err := time.Parse(time.RFC3339Nano, request.Liveness.ObservedAt)
	return err == nil
}

func loadRunObservation(conn *sql.DB, id string) (RunObservation, error) {
	var observation RunObservation
	var raw string
	err := conn.QueryRow(`SELECT id, run_id, kind, actor_id, observed_at, observation_json FROM run_observations WHERE id = ?`, id).Scan(&observation.ID, &observation.RunID, &observation.Kind, &observation.ActorID, &observation.ObservedAt, &raw)
	if err != nil {
		return RunObservation{}, fmt.Errorf("load run observation: %w", err)
	}
	if err := json.Unmarshal([]byte(raw), &observation); err != nil {
		return RunObservation{}, fmt.Errorf("decode run observation: %w", err)
	}
	return observation, nil
}

func latestLiveness(conn *sql.DB, runID string) (*Liveness, error) {
	var raw string
	err := conn.QueryRow(`SELECT observation_json FROM run_observations WHERE run_id = ? AND kind = 'liveness' ORDER BY observed_at DESC, id DESC LIMIT 1`, runID).Scan(&raw)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load run liveness: %w", err)
	}
	var observation RunObservation
	if err := json.Unmarshal([]byte(raw), &observation); err != nil {
		return nil, fmt.Errorf("decode run liveness: %w", err)
	}
	return &observation.Liveness, nil
}
