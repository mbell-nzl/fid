package workflow

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

const (
	runStartAction  = "run.start"
	runFinishAction = "run.finish"
)

type Executor struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type RunStartRequest struct {
	IssueID          string     `json:"issue_id"`
	EnvelopeID       string     `json:"envelope_id"`
	EnvelopeRevision int        `json:"envelope_revision"`
	ClaimID          string     `json:"claim_id"`
	ClaimVersion     int        `json:"claim_version"`
	Executor         Executor   `json:"executor"`
	ObligationID     *string    `json:"obligation_id"`
	Repository       Repository `json:"repository"`
}

type RunFinishRequest struct {
	RunID          string     `json:"run_id"`
	ExpectedStatus string     `json:"expected_status"`
	Outcome        string     `json:"outcome"`
	Repository     Repository `json:"repository"`
}

type RunOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
	Now            time.Time
}

type Run struct {
	ID                 string      `json:"id"`
	IssueID            string      `json:"issue_id"`
	EnvelopeID         string      `json:"envelope_id"`
	EnvelopeRevision   int         `json:"envelope_revision"`
	ClaimID            string      `json:"claim_id"`
	ClaimVersion       int         `json:"claim_version"`
	Version            int         `json:"version"`
	Executor           Executor    `json:"executor"`
	ObligationID       *string     `json:"obligation_id"`
	ActorID            string      `json:"actor_id"`
	StartRepository    Repository  `json:"start_repository"`
	TerminalRepository *Repository `json:"terminal_repository"`
	Status             string      `json:"status"`
	Budgets            []Budget    `json:"budgets"`
	StartedAt          string      `json:"started_at"`
	EndedAt            *string     `json:"ended_at"`
}

type RunMutationResult struct {
	Run            Run    `json:"run"`
	IdempotencyKey string `json:"idempotency_key"`
	Replayed       bool   `json:"replayed"`
	EventID        string
}

type RunListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	EnvelopeID    string
	Status        string
	Limit         int
}

type RunNotFoundError struct{ ID string }

func (err *RunNotFoundError) Error() string { return fmt.Sprintf("run %q was not found", err.ID) }

type RunInactiveError struct {
	RunID  string
	Status string
}

func (err *RunInactiveError) Error() string {
	return fmt.Sprintf("run %q is %s", err.RunID, err.Status)
}

type RunInvalidError struct{ Reason string }

func (err *RunInvalidError) Error() string { return err.Reason }

type RunVersionConflictError struct {
	RunID           string
	ExpectedVersion int
	ActualVersion   int
}

func (err *RunVersionConflictError) Error() string { return "run version conflict" }

type DuplicateRunExecutorError struct{ ExecutorID string }

func (err *DuplicateRunExecutorError) Error() string { return "run executor already exists" }

func StartRun(options RunOptions, request RunStartRequest) (RunMutationResult, error) {
	if err := validateRunOptions(options); err != nil || !validRunStartRequest(request) {
		return RunMutationResult{}, invalidRun("run start input is invalid")
	}
	return mutateRun(options, runStartAction, request, func(conn *sql.DB, projectID string, now time.Time) (journal.Event, Run, error) {
		if err := validateRunStartBinding(conn, options, projectID, request, now); err != nil {
			return journal.Event{}, Run{}, err
		}
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		id, err := journal.GenerateID("RUN", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM runs WHERE id = ?`, candidate).Scan(&found)
			return err == nil, nil
		})
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		envelope, err := loadEnvelope(conn, request.EnvelopeID)
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		run := Run{ID: id, IssueID: request.IssueID, EnvelopeID: request.EnvelopeID, EnvelopeRevision: request.EnvelopeRevision, ClaimID: request.ClaimID, ClaimVersion: request.ClaimVersion, Version: 1, Executor: request.Executor, ObligationID: request.ObligationID, ActorID: options.ActorID, StartRepository: request.Repository, Status: "active", Budgets: envelope.Budgets, StartedAt: now.Format(time.RFC3339Nano)}
		event, err := runEvent(options.WorkspaceRoot, actor, run, "run.started", runStartAction, runStartEffects(run))
		return event, run, err
	})
}

func FinishRun(options RunOptions, request RunFinishRequest) (RunMutationResult, error) {
	if err := validateRunOptions(options); err != nil || !validRunFinishRequest(request) {
		return RunMutationResult{}, invalidRun("run finish input is invalid")
	}
	return mutateRun(options, runFinishAction, request, func(conn *sql.DB, _ string, now time.Time) (journal.Event, Run, error) {
		run, err := loadRun(conn, request.RunID)
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		if run.Status != "active" {
			return journal.Event{}, Run{}, &RunInactiveError{RunID: run.ID, Status: run.Status}
		}
		if run.ActorID != options.ActorID || !sameRepositoryLocation(run.StartRepository, request.Repository) {
			return journal.Event{}, Run{}, invalidRun("run finish binding is invalid")
		}
		if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runFinishAction, Capability: runFinishAction, IssueID: run.IssueID, EnvelopeID: run.EnvelopeID, Revision: run.EnvelopeRevision}); err != nil {
			return journal.Event{}, Run{}, err
		}
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return journal.Event{}, Run{}, err
		}
		endedAt := now.Format(time.RFC3339Nano)
		run.Status, run.TerminalRepository, run.EndedAt = request.Outcome, &request.Repository, &endedAt
		event, err := runEvent(options.WorkspaceRoot, actor, run, "run.finished", runFinishAction, []journal.Effect{{
			Op: "update", Table: "runs", PK: map[string]any{"id": run.ID},
			Values:        map[string]any{"status": run.Status, "terminal_repository_json": string(mustRunJSON(request.Repository)), "ended_at": endedAt},
			Preconditions: map[string]any{"status": "active"},
		}})
		return event, run, err
	})
}

func mutateRun(options RunOptions, operation string, request any, build func(*sql.DB, string, time.Time) (journal.Event, Run, error)) (RunMutationResult, error) {
	var result RunMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid " + operation}, func(lock *journal.WriteLock) error {
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
			return fmt.Errorf("%s: read project config: %w", operation, err)
		}
		digest, err := runRequestDigest(operation, project.ProjectID, options.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: operation + ":" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			run, err := loadRun(conn, replay.EntityID)
			if err != nil {
				return err
			}
			result = RunMutationResult{Run: run, IdempotencyKey: options.IdempotencyKey, Replayed: true, EventID: replay.EventID}
			return nil
		}
		event, run, err := build(conn, project.ProjectID, runNow(options.Now))
		if err != nil {
			return err
		}
		resultJSON, err := journal.MarshalCanonicalJSON(map[string]any{"run": run, "idempotency_key": options.IdempotencyKey, "replayed": false})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: run.ID, ResultJSON: string(resultJSON)}, nil
		})
		if err != nil {
			return err
		}
		loaded, err := loadRun(conn, outcome.EntityID)
		if err != nil {
			return err
		}
		result = RunMutationResult{Run: loaded, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func validateRunStartBinding(conn *sql.DB, options RunOptions, projectID string, request RunStartRequest, now time.Time) error {
	if !validEnvelopeRepository(request.Repository) || request.Repository.ProjectID != projectID {
		return invalidRun("run repository is invalid")
	}
	var claim struct {
		IssueID, EnvelopeID, ActorID, WorktreeID, Status, ExpiresAt string
		EnvelopeRevision, Version                                   int
	}
	err := conn.QueryRow(`SELECT issue_id, envelope_id, envelope_revision, actor_id, worktree_id, version, status, expires_at FROM claims WHERE id = ?`, request.ClaimID).Scan(&claim.IssueID, &claim.EnvelopeID, &claim.EnvelopeRevision, &claim.ActorID, &claim.WorktreeID, &claim.Version, &claim.Status, &claim.ExpiresAt)
	if err == sql.ErrNoRows {
		return invalidRun("run claim was not found")
	}
	if err != nil {
		return fmt.Errorf("load run claim: %w", err)
	}
	expiresAt, err := time.Parse(time.RFC3339Nano, claim.ExpiresAt)
	if err != nil || !now.Before(expiresAt.UTC()) || claim.Status != "active" || claim.IssueID != request.IssueID || claim.EnvelopeID != request.EnvelopeID || claim.EnvelopeRevision != request.EnvelopeRevision || claim.Version != request.ClaimVersion || claim.ActorID != options.ActorID || claim.WorktreeID != request.Repository.WorktreeID {
		return invalidRun("run claim binding is invalid")
	}
	envelope, err := loadEnvelope(conn, request.EnvelopeID)
	if err != nil {
		return err
	}
	if envelope.Status != "active" || envelope.PrimaryIssueID != request.IssueID || envelope.CurrentRevision != request.EnvelopeRevision || !sameRepository(envelope.RepositoryScope, request.Repository) {
		return invalidRun("run envelope binding is invalid")
	}
	if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: runStartAction, Capability: runStartAction, IssueID: request.IssueID, EnvelopeID: request.EnvelopeID, Revision: request.EnvelopeRevision}); err != nil {
		return err
	}
	if request.ObligationID != nil {
		var found int
		if err := conn.QueryRow(`SELECT 1 FROM envelope_obligations WHERE id = ? AND envelope_id = ? AND revision = ?`, *request.ObligationID, request.EnvelopeID, request.EnvelopeRevision).Scan(&found); err != nil {
			if err == sql.ErrNoRows {
				return invalidRun("run obligation binding is invalid")
			}
			return fmt.Errorf("load run obligation: %w", err)
		}
	}
	var existingRun string
	if err := conn.QueryRow(`SELECT run_id FROM run_runtime WHERE executor_id = ?`, request.Executor.ID).Scan(&existingRun); err == nil {
		return &DuplicateRunExecutorError{ExecutorID: request.Executor.ID}
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("load run executor: %w", err)
	}
	var id string
	err = conn.QueryRow(`SELECT id FROM runs WHERE claim_id = ? AND status = 'active'`, request.ClaimID).Scan(&id)
	if err == nil {
		return invalidRun("claim already has an active run")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("load active run: %w", err)
	}
	return nil
}

func runStartEffects(run Run) []journal.Effect {
	effects := []journal.Effect{{Op: "insert", Table: "runs", PK: map[string]any{"id": run.ID}, Values: map[string]any{
		"issue_id": run.IssueID, "envelope_id": run.EnvelopeID, "envelope_revision": run.EnvelopeRevision, "claim_id": run.ClaimID, "claim_version": run.ClaimVersion, "actor_id": run.ActorID,
		"start_repository_json": string(mustRunJSON(run.StartRepository)), "terminal_repository_json": nil, "status": run.Status, "started_at": run.StartedAt, "ended_at": nil,
	}}, {Op: "insert", Table: "run_runtime", PK: map[string]any{"run_id": run.ID}, Values: map[string]any{"version": run.Version, "executor_id": run.Executor.ID, "obligation_id": run.ObligationID}}}
	for _, budget := range run.Budgets {
		effects = append(effects, journal.Effect{Op: "insert", Table: "run_budget_bindings", PK: map[string]any{"run_id": run.ID, "name": budget.Name}, Values: map[string]any{"budget_json": string(mustRunJSON(budget))}})
	}
	return effects
}

func runEvent(workspaceRoot string, actor ActorProfile, run Run, eventType, command string, effects []journal.Effect) (journal.Event, error) {
	next, err := journal.NextAggregateSequence(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: run.ID, Type: "run"})
	if err != nil {
		return journal.Event{}, err
	}
	repository := run.StartRepository
	if run.TerminalRepository != nil {
		repository = *run.TerminalRepository
	}
	return journal.Event{EventSchema: journal.EventSchemaV2, EventType: eventType, ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: run.ID, Sequence: next, Type: "run"}, CommandName: command, Repository: envelopeRepository(repository), Effects: effects}, nil
}

func validateRunOptions(options RunOptions) error {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.IdempotencyKey) == "" {
		return invalidRun("run options are incomplete")
	}
	return nil
}

func validRunStartRequest(request RunStartRequest) bool {
	return strings.TrimSpace(request.IssueID) != "" && strings.TrimSpace(request.EnvelopeID) != "" && request.EnvelopeRevision == 1 && strings.TrimSpace(request.ClaimID) != "" && request.ClaimVersion > 0 && strings.TrimSpace(request.Executor.ID) != "" && request.Executor.Source == "activated-harness"
}

func validRunFinishRequest(request RunFinishRequest) bool {
	return strings.TrimSpace(request.RunID) != "" && request.ExpectedStatus == "active" && (request.Outcome == "abandoned" || request.Outcome == "failed" || request.Outcome == "succeeded") && validEnvelopeRepository(request.Repository)
}

func runRequestDigest(operation, projectID, actorID string, request any) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"actor_id": actorID, "operation": operation, "project_id": projectID, "request": request})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func runNow(now time.Time) time.Time {
	if now.IsZero() {
		return time.Now().UTC()
	}
	return now.UTC()
}

func mustRunJSON(value any) []byte {
	payload, err := journal.MarshalCanonicalJSON(value)
	if err != nil {
		panic(fmt.Sprintf("marshal Run value: %v", err))
	}
	return payload
}

func invalidRun(reason string) error { return &RunInvalidError{Reason: reason} }
