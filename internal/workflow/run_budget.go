package workflow

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	fidproject "github.com/mbell-nzl/fid/internal/project"
)

type RunBudgetObserveRequest struct {
	RunID            string                  `json:"run_id"`
	ExpectedVersion  int                     `json:"expected_version"`
	BudgetName       string                  `json:"budget_name"`
	SourceSampleID   string                  `json:"source_sample_id"`
	SampleKind       string                  `json:"sample_kind"`
	PreviousSampleID *string                 `json:"previous_sample_id"`
	Value            *int64                  `json:"value"`
	MeasurementRef   *fidproject.SourceInput `json:"measurement_ref"`
	ObservedAt       string                  `json:"observed_at"`
}

type RunBudgetObserveResult struct {
	Sample         BudgetSample `json:"sample"`
	Budget         Budget       `json:"budget"`
	Stops          []Stop       `json:"stops"`
	Attention      []Attention  `json:"attention"`
	IdempotencyKey string       `json:"idempotency_key"`
	Replayed       bool         `json:"replayed"`
	EventID        string
}

// ObserveRunBudget records an activated-harness measurement and derives its
// current counter from the immutable sample chain.
func ObserveRunBudget(options RunOptions, request RunBudgetObserveRequest) (RunBudgetObserveResult, error) {
	if err := validateRunOptions(options); err != nil || !validRunBudgetObserveRequest(request) {
		return RunBudgetObserveResult{}, invalidRun("run observe-budget input is invalid")
	}
	var result RunBudgetObserveResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid run.observe-budget"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openWorkflowDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer closeWorkflowDB(conn)
		project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
		if err != nil {
			return err
		}
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor_id": options.ActorID, "operation": "run.observe-budget", "project_id": project.ProjectID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idem := journal.IdempotencyRequest{OperationScope: "run.observe-budget:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idem); err != nil {
			return err
		} else if found {
			return loadBudgetObserveReplay(conn, replay, options.IdempotencyKey, &result)
		}
		var measurement *Source
		var observed *fidproject.RecordedRepository
		if request.MeasurementRef != nil {
			repository, err := fidproject.ObserveRecordedRepository(options.WorkspaceRoot, project.ProjectID)
			if err != nil {
				return err
			}
			digest, err := fidproject.ResolveSourceInput(options.WorkspaceRoot, *request.MeasurementRef)
			if err != nil {
				return err
			}
			input := request.MeasurementRef
			measurement = &Source{SourceType: input.SourceType, DisplayLabel: input.DisplayLabel, Path: input.Path, Section: input.Section, ExternalID: input.ExternalID, Digest: digest, Repository: Repository(repository)}
			if !validEnvelopeSource(*measurement) {
				return invalidRun("run observe-budget measurement source is invalid")
			}
			observed = &repository
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
		if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: "run.observe-budget", Capability: "run.observe-budget", IssueID: run.IssueID, EnvelopeID: run.EnvelopeID, Revision: run.EnvelopeRevision}); err != nil {
			return err
		}
		budget, err := loadRunBudget(conn, run.EnvelopeID, run.EnvelopeRevision, request.BudgetName)
		if err != nil {
			return err
		}
		if fidOwnedBudgetUnit(budget.Unit) {
			return invalidRun("external budget observations cannot charge fid-owned units")
		}
		samples, err := loadRunBudgetSamples(conn, run.EnvelopeID, run.EnvelopeRevision, request.BudgetName)
		if err != nil {
			return err
		}
		sample := BudgetSample{EnvelopeID: run.EnvelopeID, EnvelopeRevision: run.EnvelopeRevision, BudgetName: request.BudgetName, RunID: run.ID, ActorID: options.ActorID, Source: "activated-harness", SourceSampleID: request.SourceSampleID, SampleKind: request.SampleKind, PreviousSampleID: request.PreviousSampleID, Value: request.Value, MeasurementRef: measurement, ObservedAt: request.ObservedAt}
		admission, err := AdmitBudgetSample(budget, samples, sample)
		if err != nil {
			return err
		}
		sample.ID, err = journal.GenerateID("BSM", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM budget_samples WHERE id = ?`, candidate).Scan(&found)
			return err == nil, nil
		})
		if err != nil {
			return err
		}
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		if actor.Producer != "activated-harness" {
			return &policy.EnvelopeAuthorityDeniedError{ActorID: options.ActorID, Action: "run.observe-budget", Capability: "run.observe-budget"}
		}
		before, err := DeriveBudget(budget, samples)
		if err != nil {
			return err
		}
		blocks := budgetBoundaryBlocks(before, admission.Budget)
		stopEffects, stops, attention, err := BudgetStopEffects(conn, Subject{Type: "run", ID: run.ID}, run.EnvelopeID, actor.ActorID, blocks, request.ObservedAt)
		if err != nil {
			return err
		}
		effects := []journal.Effect{{Op: "insert", Table: "budget_samples", PK: map[string]any{"id": sample.ID}, Values: budgetSampleValues(sample)}}
		effects = append(effects, stopEffects...)
		effects = append(effects, journal.Effect{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": run.ID}, Values: map[string]any{"version": run.Version + 1}, Preconditions: map[string]any{"version": run.Version}})
		event, err := runEvent(options.WorkspaceRoot, actor, run, "run.budget-observed", "run.observe-budget", effects)
		if err != nil {
			return err
		}
		if observed != nil {
			event.Repository = envelopeRepository(Repository(*observed))
		}
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		sample.EventID = prepared.EventID
		resultJSON, err := journal.MarshalCanonicalJSON(map[string]any{"sample": sample, "budget": admission.Budget, "stops": stops, "attention": attention, "idempotency_key": options.IdempotencyKey, "replayed": false})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idem, prepared, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			fillBudgetCreatedEventIDs(prepared.Effects, prepared.EventID)
			return journal.IdempotencyMutation{EntityID: sample.ID, ResultJSON: string(resultJSON)}, nil
		})
		if err != nil {
			return err
		}
		result = RunBudgetObserveResult{Sample: sample, Budget: admission.Budget, Stops: stops, Attention: attention, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func validRunBudgetObserveRequest(request RunBudgetObserveRequest) bool {
	if strings.TrimSpace(request.RunID) == "" || request.ExpectedVersion < 1 || strings.TrimSpace(request.BudgetName) == "" || strings.TrimSpace(request.SourceSampleID) == "" || strings.TrimSpace(request.ObservedAt) == "" {
		return false
	}
	if request.SampleKind != "delta" && request.SampleKind != "baseline" {
		return false
	}
	if request.SampleKind == "baseline" && request.PreviousSampleID != nil {
		return false
	}
	if request.Value != nil && *request.Value < 0 {
		return false
	}
	if _, err := time.Parse(time.RFC3339Nano, request.ObservedAt); err != nil {
		return false
	}
	return true // SourceInput is checked after replay without reading on retries.
}

func loadRunBudget(conn *sql.DB, envelopeID string, revision int, name string) (Budget, error) {
	var budget Budget
	var raw string
	err := conn.QueryRow(`SELECT budget_json FROM envelope_budgets WHERE envelope_id = ? AND revision = ? AND name = ?`, envelopeID, revision, name).Scan(&raw)
	if err != nil {
		return Budget{}, err
	}
	if err := json.Unmarshal([]byte(raw), &budget); err != nil {
		return Budget{}, err
	}
	return budget, nil
}

func budgetSampleValues(sample BudgetSample) map[string]any {
	var measurement any
	if sample.MeasurementRef != nil {
		measurement = string(mustRunJSON(*sample.MeasurementRef))
	}
	return map[string]any{"envelope_id": sample.EnvelopeID, "envelope_revision": sample.EnvelopeRevision, "budget_name": sample.BudgetName, "run_id": sample.RunID, "actor_id": sample.ActorID, "source": sample.Source, "source_sample_id": sample.SourceSampleID, "sample_kind": sample.SampleKind, "previous_sample_id": sample.PreviousSampleID, "value": sample.Value, "measurement_ref_json": measurement, "observed_at": sample.ObservedAt, "created_event_id": "<event-id>"}
}

func budgetBoundaryBlocks(before, after Budget) []BudgetAdmissionBlock {
	if before.Remaining != nil && *before.Remaining == 0 {
		return nil
	}
	return BudgetAdmissionBlocks([]Budget{after})
}

func fillBudgetCreatedEventIDs(effects []journal.Effect, eventID string) {
	for index := range effects {
		if effects[index].Values["created_event_id"] == "<event-id>" {
			effects[index].Values["created_event_id"] = eventID
		}
	}
}

func loadBudgetObserveReplay(conn *sql.DB, replay journal.IdempotencyOutcome, key string, result *RunBudgetObserveResult) error {
	var stored RunBudgetObserveResult
	if err := json.Unmarshal([]byte(replay.ResultJSON), &stored); err != nil {
		return err
	}
	stored.IdempotencyKey, stored.Replayed, stored.EventID = key, true, replay.EventID
	*result = stored
	return nil
}
