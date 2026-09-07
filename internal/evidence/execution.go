package evidence

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

// GateRunRequest is the closed gate.run input after CLI decoding.
type GateRunRequest struct {
	RegistrationID      string     `json:"registration_id"`
	RegistrationVersion int        `json:"registration_version"`
	ObligationID        string     `json:"obligation_id"`
	RunID               string     `json:"run_id"`
	Repository          Repository `json:"repository"`
}

// Attempt is the durable pre-spawn execution fact.
type Attempt struct {
	ID                   string     `json:"id"`
	RegistrationID       string     `json:"registration_id"`
	RegistrationVersion  int        `json:"registration_version"`
	ObligationID         string     `json:"obligation_id"`
	RunID                string     `json:"run_id"`
	ActorID              string     `json:"actor_id"`
	StartRepository      Repository `json:"start_repository"`
	ResolvedPolicyDigest string     `json:"resolved_policy_digest"`
	StartedAt            string     `json:"started_at"`
}

// Receipt is the immutable terminal observation for one Attempt.
type Receipt struct {
	ID                  string              `json:"id"`
	AttemptID           string              `json:"attempt_id"`
	RegistrationID      string              `json:"registration_id"`
	RegistrationVersion int                 `json:"registration_version"`
	ObligationID        string              `json:"obligation_id"`
	RunID               string              `json:"run_id"`
	ActorID             string              `json:"actor_id"`
	TrustSource         string              `json:"trust_source"`
	Outcome             string              `json:"outcome"`
	ExitCode            *int                `json:"exit_code"`
	Signal              *string             `json:"signal"`
	TimedOut            bool                `json:"timed_out"`
	StartedAt           string              `json:"started_at"`
	EndedAt             string              `json:"ended_at"`
	TerminalRepository  Repository          `json:"terminal_repository"`
	Toolchain           map[string]string   `json:"toolchain"`
	OutputExcerpt       string              `json:"output_excerpt"`
	OutputDigest        string              `json:"output_digest"`
	OutputBytes         int64               `json:"output_bytes"`
	OutputTruncated     bool                `json:"output_truncated"`
	RedactionApplied    bool                `json:"redaction_applied"`
	Artifacts           []Artifact          `json:"artifacts"`
	RejectedArtifacts   []ArtifactRejection `json:"rejected_artifacts"`
}

type GateRunOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
	Request        GateRunRequest
	runLive        func(context.Context, string, GateDefinition, string) (LiveResult, error)
}

type GateRunResult struct {
	Attempt         Attempt
	Receipt         *Receipt
	Denial          *BudgetDenial
	Live            LiveResult
	EventID         string
	EventOccurredAt string
	IdempotencyKey  string
	Replayed        bool
	definition      GateDefinition
}

// BudgetDenial is the durable pre-spawn rejection produced by a required
// unknown or exhausted budget chain.
type BudgetDenial struct {
	Budget          workflow.Budget `json:"budget"`
	StopIDs         []string        `json:"stop_ids"`
	AttentionIDs    []string        `json:"attention_ids"`
	EventID         string          `json:"event_id"`
	EventOccurredAt string          `json:"event_occurred_at"`
}

type BudgetAdmissionDeniedError struct {
	Code   string
	Denial BudgetDenial
}

func (err *BudgetAdmissionDeniedError) Error() string { return err.Code }

type gateRunBinding struct {
	registration Registration
	actor        storedActor
	projectID    string
}

type fidBudgetCharge struct {
	BudgetName     string
	Source         string
	SourceSampleID string
	SampleKind     string
	Value          *int64
	ObservedAt     string
}

type budgetChain struct {
	Budget           workflow.Budget
	EnvelopeID       string
	EnvelopeRevision int
	Samples          []workflow.BudgetSample
	RuntimeVersion   int
}

// fidBudgetCharges derives automatic charges only from facts fid observed while
// executing a registered gate. Caller input never contributes a charge.
func fidBudgetCharges(budgets []workflow.Budget, attempt Attempt, receipt *Receipt) ([]fidBudgetCharge, error) {
	charges := make([]fidBudgetCharge, 0, len(budgets))
	for _, budget := range budgets {
		var value *int64
		sourceID, observedAt := attempt.ID, attempt.StartedAt
		switch budget.Unit {
		case "attempts":
			if receipt == nil {
				unit := int64(1)
				value = &unit
			}
		case "failures":
			if receipt != nil && receipt.Outcome != "succeeded" {
				unit := int64(1)
				value, sourceID, observedAt = &unit, receipt.ID, receipt.EndedAt
			}
		case "milliseconds":
			if receipt != nil {
				started, startErr := time.Parse(time.RFC3339Nano, receipt.StartedAt)
				ended, endErr := time.Parse(time.RFC3339Nano, receipt.EndedAt)
				if startErr != nil || endErr != nil || ended.Before(started) {
					return nil, fmt.Errorf("derive fid budget milliseconds: invalid Receipt timing")
				}
				unit := ended.Sub(started).Milliseconds()
				value, sourceID, observedAt = &unit, receipt.ID, receipt.EndedAt
			}
		case "tool-failures":
			if receipt != nil && receipt.Outcome == "spawn-failed" {
				unit := int64(1)
				value, sourceID, observedAt = &unit, receipt.ID, receipt.EndedAt
			}
		}
		if value != nil {
			charges = append(charges, fidBudgetCharge{BudgetName: budget.Name, Source: "fid", SourceSampleID: sourceID, SampleKind: "delta", Value: value, ObservedAt: observedAt})
		}
	}
	return charges, nil
}

// RunGate appends an Attempt before invoking the process and records exactly
// one Receipt only after RunLive returned an observed terminal process result.
func RunGate(ctx context.Context, options GateRunOptions) (GateRunResult, error) {
	if err := validateGateRunOptions(options); err != nil {
		return GateRunResult{}, err
	}
	if options.runLive == nil {
		options.runLive = RunLive
	}

	var prepared GateRunResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid gate run"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openEvidenceDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()

		binding, requestDigest, err := loadGateRunBinding(conn, options)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "gate.run:" + binding.projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return replayGateRunResult(conn, options, outcome, &prepared)
		}
		chains, err := loadBudgetChains(conn, options.Request.RunID)
		if err != nil {
			return err
		}
		blocks, err := budgetAdmissionBlocks(chains)
		if err != nil {
			return err
		}
		if len(blocks) != 0 {
			denial, err := appendExecutionBlocked(lock, conn, binding, options, requestDigest, chains, blocks)
			if err != nil {
				return err
			}
			prepared = GateRunResult{Denial: &denial, EventID: denial.EventID, EventOccurredAt: denial.EventOccurredAt, IdempotencyKey: options.IdempotencyKey}
			return &BudgetAdmissionDeniedError{Code: blocks[0].Code, Denial: denial}
		}
		stops, err := workflow.ApplicableStops(conn, "gate.run", workflow.Subject{Type: "run", ID: options.Request.RunID})
		if err != nil {
			return err
		}
		if len(stops) != 0 {
			return &workflow.StopBlockedError{Action: "gate.run", Stops: stops}
		}

		if _, _, cancel, err := executorCommand(ctx, options.WorkspaceRoot, binding.registration.Policy); err != nil {
			return err
		} else {
			cancel()
		}
		attempt, err := appendAttempt(lock, conn, binding, options, requestDigest)
		if err != nil {
			return err
		}
		prepared = GateRunResult{Attempt: attempt, IdempotencyKey: options.IdempotencyKey, definition: binding.registration.Policy}
		return nil
	})
	if err != nil || prepared.Replayed {
		return prepared, err
	}

	live, err := options.runLive(ctx, options.WorkspaceRoot, prepared.definition, prepared.Attempt.ID)
	if err != nil {
		return GateRunResult{Attempt: prepared.Attempt, Live: live, IdempotencyKey: options.IdempotencyKey}, err
	}
	prepared.Live = live

	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid gate run receipt"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openEvidenceDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		binding, requestDigest, err := loadGateRunBinding(conn, options)
		if err != nil {
			return err
		}
		receipt, eventID, err := appendReceipt(lock, conn, binding, prepared.Attempt, options, requestDigest, &prepared.Live)
		if err != nil {
			return err
		}
		occurredAt, err := gateEventOccurredAt(options.WorkspaceRoot, eventID)
		if err != nil {
			return err
		}
		prepared.Receipt, prepared.EventID, prepared.EventOccurredAt = &receipt, eventID, occurredAt
		return nil
	})
	if err != nil {
		return GateRunResult{Attempt: prepared.Attempt, Live: live, IdempotencyKey: options.IdempotencyKey}, err
	}
	return prepared, nil
}

func gateEventOccurredAt(workspaceRoot, eventID string) (string, error) {
	found, err := journal.FindEvent(filepath.Join(workspaceRoot, config.WorkspaceDirName, "journal"), eventID)
	if err != nil {
		return "", fmt.Errorf("run gate: read terminal Event %q: %w", eventID, err)
	}
	return found.Event.OccurredAt.UTC().Format(time.RFC3339Nano), nil
}

func validateGateRunOptions(options GateRunOptions) error {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" {
		return fmt.Errorf("run gate: workspace root, migrations directory, and actor are required")
	}
	request := options.Request
	if strings.TrimSpace(request.RegistrationID) == "" || request.RegistrationVersion < 1 || strings.TrimSpace(request.ObligationID) == "" || strings.TrimSpace(request.RunID) == "" {
		return fmt.Errorf("run gate: registration, obligation, and Run are required")
	}
	return nil
}

func loadGateRunBinding(conn *sql.DB, options GateRunOptions) (gateRunBinding, string, error) {
	project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
	if err != nil {
		return gateRunBinding{}, "", fmt.Errorf("run gate: read project config: %w", err)
	}
	registration, err := loadRegistration(conn, options.Request.RegistrationID, options.Request.RegistrationVersion)
	if err != nil {
		return gateRunBinding{}, "", err
	}
	if err := validateRepositoryBinding(RegisterRequest{WorktreeID: registration.WorktreeID, Repository: options.Request.Repository}, project.ProjectID); err != nil || registration.ProjectID != project.ProjectID {
		return gateRunBinding{}, "", fmt.Errorf("run gate: repository binding is invalid")
	}
	_, actor, err := loadActiveActor(conn, options.ActorID)
	if err != nil {
		return gateRunBinding{}, "", err
	}
	if !hasCapability(actor.Capabilities, "gate.run") {
		return gateRunBinding{}, "", &AuthorityDeniedError{ActorID: options.ActorID}
	}
	if err := validateGateRunReferences(conn, registration, actor.ActorID, options.Request); err != nil {
		return gateRunBinding{}, "", err
	}
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"actor_id": options.ActorID, "project_id": project.ProjectID, "request": options.Request})
	if err != nil {
		return gateRunBinding{}, "", fmt.Errorf("run gate: canonicalize request: %w", err)
	}
	sum := sha256.Sum256(payload)
	return gateRunBinding{registration: registration, actor: actor, projectID: project.ProjectID}, "sha256:" + hex.EncodeToString(sum[:]), nil
}

func validateGateRunReferences(conn *sql.DB, registration Registration, actorID string, request GateRunRequest) error {
	var envelopeID string
	var revision int
	err := conn.QueryRow(`SELECT envelope_id, revision FROM envelope_obligations WHERE id = ? AND registration_id = ? AND registration_version = ?`, request.ObligationID, registration.ID, registration.Version).Scan(&envelopeID, &revision)
	if err == sql.ErrNoRows {
		return fmt.Errorf("run gate: obligation %q is not bound to registration", request.ObligationID)
	}
	if err != nil {
		return fmt.Errorf("run gate: load obligation: %w", err)
	}
	var issueID, runActor, status, repositoryJSON string
	err = conn.QueryRow(`SELECT issue_id, actor_id, status, start_repository_json FROM runs WHERE id = ? AND envelope_id = ? AND envelope_revision = ?`, request.RunID, envelopeID, revision).Scan(&issueID, &runActor, &status, &repositoryJSON)
	if err == sql.ErrNoRows {
		return fmt.Errorf("run gate: Run %q is not bound to obligation", request.RunID)
	}
	if err != nil {
		return fmt.Errorf("run gate: load Run: %w", err)
	}
	if status != "active" || runActor != actorID || issueID != registration.Scope.IssueID {
		return fmt.Errorf("run gate: active Run binding is invalid")
	}
	var runRepository Repository
	if err := json.Unmarshal([]byte(repositoryJSON), &runRepository); err != nil || !sameGateRepository(runRepository, request.Repository) {
		return fmt.Errorf("run gate: repository does not match active Run")
	}
	var actionsJSON, selectorJSON, scopeJSON string
	err = conn.QueryRow(`SELECT actions_json, actor_selector_json, scope_json FROM envelope_authority_grants WHERE envelope_id = ? AND revision = ?`, envelopeID, revision).Scan(&actionsJSON, &selectorJSON, &scopeJSON)
	if err != nil {
		return fmt.Errorf("run gate: load authority grant: %w", err)
	}
	var actions []string
	var selector struct {
		ActorID              string   `json:"actor_id"`
		RequiredCapabilities []string `json:"required_capabilities"`
	}
	var scope GateScope
	if json.Unmarshal([]byte(actionsJSON), &actions) != nil || json.Unmarshal([]byte(selectorJSON), &selector) != nil || json.Unmarshal([]byte(scopeJSON), &scope) != nil || selector.ActorID != actorID || scope.IssueID != issueID || !hasCapability(actions, "gate.run") || !hasCapability(selector.RequiredCapabilities, "gate.run") {
		return fmt.Errorf("run gate: authority grant does not permit gate.run")
	}
	return nil
}

func sameGateRepository(left, right Repository) bool {
	return sameStringPointer(left.BaseOID, right.BaseOID) && left.CaptureSource == right.CaptureSource && left.Dirty == right.Dirty && sameStringPointer(left.DirtyDigest, right.DirtyDigest) && left.HeadOID == right.HeadOID && left.ProjectID == right.ProjectID && left.TreeOID == right.TreeOID && left.WorktreeID == right.WorktreeID
}

func appendAttempt(lock *journal.WriteLock, conn *sql.DB, binding gateRunBinding, options GateRunOptions, requestDigest string) (Attempt, error) {
	id, err := journal.GenerateID("ATT", func(candidate string) (bool, error) {
		var count int
		err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts WHERE id = ?`, candidate).Scan(&count)
		return count > 0, err
	})
	if err != nil {
		return Attempt{}, err
	}
	now := time.Now().UTC()
	repositoryJSON, err := journal.MarshalCanonicalJSON(options.Request.Repository)
	if err != nil {
		return Attempt{}, err
	}
	attempt := Attempt{ID: id, RegistrationID: binding.registration.ID, RegistrationVersion: binding.registration.Version, ObligationID: options.Request.ObligationID, RunID: options.Request.RunID, ActorID: binding.actor.ActorID, StartRepository: options.Request.Repository, ResolvedPolicyDigest: binding.registration.DefinitionDigest, StartedAt: now.Format(time.RFC3339Nano)}
	chains, err := loadBudgetChains(conn, attempt.RunID)
	if err != nil {
		return Attempt{}, err
	}
	charges, err := fidBudgetCharges(chainBudgets(chains), attempt, nil)
	if err != nil {
		return Attempt{}, err
	}
	pendingResult, err := journal.MarshalCanonicalJSON(map[string]any{"attempt": attempt, "idempotency_key": options.IdempotencyKey, "receipt": nil, "replayed": false})
	if err != nil {
		return Attempt{}, fmt.Errorf("run gate: canonicalize pending Attempt result: %w", err)
	}
	effects := []journal.Effect{
		{Op: "insert", Table: "execution_attempts", PK: map[string]any{"id": attempt.ID}, Values: map[string]any{"registration_id": attempt.RegistrationID, "registration_version": attempt.RegistrationVersion, "obligation_id": attempt.ObligationID, "run_id": attempt.RunID, "actor_id": attempt.ActorID, "start_repository_json": string(repositoryJSON), "resolved_policy_digest": attempt.ResolvedPolicyDigest, "producer": "fid", "started_at": attempt.StartedAt}},
		{Op: "insert", Table: "idempotency_results", PK: map[string]any{"operation_scope": "gate.run:" + binding.projectID, "idempotency_key": options.IdempotencyKey}, Values: map[string]any{"request_digest": requestDigest, "result_entity_id": attempt.ID, "result_event_id": "<event-id>", "aggregate_sequence": 1, "result_json": string(pendingResult), "created_at": "<occurred-at>"}},
	}
	budgetEffects, err := appendFIDBudgetEffects(conn, chains, charges, attempt.RunID, binding.actor.ActorID)
	if err != nil {
		return Attempt{}, err
	}
	effects = append(effects, budgetEffects...)
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "gate.execution-started", ActorV2: journal.EventActor{ActorID: binding.actor.ActorID, ActorType: binding.actor.ActorType, Producer: binding.actor.Producer, ProvenanceRef: binding.actor.ProvenanceRef, SessionRef: binding.actor.SessionRef}, Aggregate: journal.Aggregate{ID: attempt.ID, Sequence: 1, Type: "attempt"}, CommandName: "gate.run", CorrelationID: &attempt.RunID, Idempotency: &journal.EventIdempotency{Key: options.IdempotencyKey, RequestDigest: requestDigest, Scope: "gate.run:" + binding.projectID}, Repository: repositoryEventValue(options.Request.Repository), Effects: effects}
	prepared, err := journal.PrepareEvent(lock, event)
	if err != nil {
		return Attempt{}, err
	}
	fillCreatedEventIDs(prepared.Effects, prepared.EventID)
	fillAttemptIdempotencyEventValues(prepared.Effects, prepared.EventID, prepared.OccurredAt)
	if _, err := journal.CommitEvent(lock, conn, prepared); err != nil {
		return Attempt{}, err
	}
	return attempt, nil
}

func chainBudgets(chains []budgetChain) []workflow.Budget {
	budgets := make([]workflow.Budget, 0, len(chains))
	for _, chain := range chains {
		budgets = append(budgets, chain.Budget)
	}
	return budgets
}

func budgetBoundaryBlocks(before, after workflow.Budget) []workflow.BudgetAdmissionBlock {
	if before.Remaining != nil && *before.Remaining == 0 {
		return nil
	}
	return workflow.BudgetAdmissionBlocks([]workflow.Budget{after})
}

func budgetAdmissionBlocks(chains []budgetChain) ([]workflow.BudgetAdmissionBlock, error) {
	budgets := make([]workflow.Budget, 0, len(chains))
	for _, chain := range chains {
		budget, err := workflow.DeriveBudget(chain.Budget, chain.Samples)
		if err != nil {
			return nil, fmt.Errorf("derive budget admission %q: %w", chain.Budget.Name, err)
		}
		budgets = append(budgets, budget)
	}
	return workflow.BudgetAdmissionBlocks(budgets), nil
}

func loadBudgetChains(conn *sql.DB, runID string) ([]budgetChain, error) {
	var envelopeID string
	var revision, runtimeVersion int
	if err := conn.QueryRow(`SELECT r.envelope_id, r.envelope_revision, rr.version FROM runs r JOIN run_runtime rr ON rr.run_id = r.id WHERE r.id = ?`, runID).Scan(&envelopeID, &revision, &runtimeVersion); err != nil {
		return nil, fmt.Errorf("load budget chain Run: %w", err)
	}
	rows, err := conn.Query(`SELECT name, budget_json FROM run_budget_bindings WHERE run_id = ? ORDER BY name`, runID)
	if err != nil {
		return nil, fmt.Errorf("load budget bindings: %w", err)
	}
	defer func() { _ = rows.Close() }()
	chains := []budgetChain{}
	for rows.Next() {
		var name, raw string
		var budget workflow.Budget
		if err := rows.Scan(&name, &raw); err != nil || json.Unmarshal([]byte(raw), &budget) != nil || budget.Name != name {
			return nil, fmt.Errorf("decode Run budget")
		}
		samples, err := loadBudgetSamples(conn, envelopeID, revision, name)
		if err != nil {
			return nil, err
		}
		chains = append(chains, budgetChain{Budget: budget, EnvelopeID: envelopeID, EnvelopeRevision: revision, Samples: samples, RuntimeVersion: runtimeVersion})
	}
	return chains, rows.Err()
}

func loadBudgetSamples(conn *sql.DB, envelopeID string, revision int, name string) ([]workflow.BudgetSample, error) {
	rows, err := conn.Query(`SELECT id, source, source_sample_id, sample_kind, previous_sample_id, value FROM budget_samples WHERE envelope_id = ? AND envelope_revision = ? AND budget_name = ? ORDER BY observed_at, id`, envelopeID, revision, name)
	if err != nil {
		return nil, fmt.Errorf("load budget samples: %w", err)
	}
	defer func() { _ = rows.Close() }()
	samples := []workflow.BudgetSample{}
	for rows.Next() {
		var sample workflow.BudgetSample
		if err := rows.Scan(&sample.ID, &sample.Source, &sample.SourceSampleID, &sample.SampleKind, &sample.PreviousSampleID, &sample.Value); err != nil {
			return nil, fmt.Errorf("read budget sample: %w", err)
		}
		samples = append(samples, sample)
	}
	return samples, rows.Err()
}

func appendFIDBudgetEffects(conn *sql.DB, chains []budgetChain, charges []fidBudgetCharge, runID, actorID string) ([]journal.Effect, error) {
	effects := []journal.Effect{}
	charged := false
	for _, charge := range charges {
		for index := range chains {
			chain := &chains[index]
			if chain.Budget.Name != charge.BudgetName {
				continue
			}
			id, err := journal.GenerateID("BSM", func(candidate string) (bool, error) {
				var count int
				err := conn.QueryRow(`SELECT COUNT(*) FROM budget_samples WHERE id = ?`, candidate).Scan(&count)
				return count > 0, err
			})
			if err != nil {
				return nil, err
			}
			sample := workflow.BudgetSample{ID: id, Source: charge.Source, SourceSampleID: charge.SourceSampleID, SampleKind: charge.SampleKind, Value: charge.Value}
			if len(chain.Samples) > 0 {
				previous := chain.Samples[len(chain.Samples)-1].ID
				sample.PreviousSampleID = &previous
			}
			before, err := workflow.DeriveBudget(chain.Budget, chain.Samples)
			if err != nil {
				return nil, fmt.Errorf("derive fid budget %q: %w", chain.Budget.Name, err)
			}
			admission, err := workflow.AdmitBudgetSample(chain.Budget, chain.Samples, sample)
			if err != nil {
				return nil, fmt.Errorf("admit fid budget sample %q: %w", chain.Budget.Name, err)
			}
			boundary := budgetBoundaryBlocks(before, admission.Budget)
			chain.Samples = append(chain.Samples, sample)
			effects = append(effects, journal.Effect{Op: "insert", Table: "budget_samples", PK: map[string]any{"id": sample.ID}, Values: map[string]any{"envelope_id": chain.EnvelopeID, "envelope_revision": chain.EnvelopeRevision, "budget_name": chain.Budget.Name, "run_id": runID, "actor_id": actorID, "source": sample.Source, "source_sample_id": sample.SourceSampleID, "sample_kind": sample.SampleKind, "previous_sample_id": sample.PreviousSampleID, "value": sample.Value, "measurement_ref_json": nil, "observed_at": charge.ObservedAt, "created_event_id": "<event-id>"}})
			if len(boundary) != 0 {
				stopEffects, _, _, err := workflow.BudgetStopEffects(conn, workflow.Subject{Type: "run", ID: runID}, chain.EnvelopeID, actorID, boundary, charge.ObservedAt)
				if err != nil {
					return nil, err
				}
				effects = append(effects, stopEffects...)
			}
			charged = true
			break
		}
	}
	if charged && len(chains) > 0 {
		effects = append(effects, journal.Effect{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": runID}, Values: map[string]any{"version": chains[0].RuntimeVersion + 1}, Preconditions: map[string]any{"version": chains[0].RuntimeVersion}})
	}
	return effects, nil
}

func appendExecutionBlocked(lock *journal.WriteLock, conn *sql.DB, binding gateRunBinding, options GateRunOptions, requestDigest string, chains []budgetChain, blocks []workflow.BudgetAdmissionBlock) (BudgetDenial, error) {
	if len(chains) == 0 || len(blocks) == 0 {
		return BudgetDenial{}, fmt.Errorf("execution budget denial is invalid")
	}
	effects, stops, attention, err := workflow.BudgetStopEffects(conn, workflow.Subject{Type: "run", ID: options.Request.RunID}, chains[0].EnvelopeID, binding.actor.ActorID, blocks, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return BudgetDenial{}, err
	}
	effects = append(effects, journal.Effect{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": options.Request.RunID}, Values: map[string]any{"version": chains[0].RuntimeVersion + 1}, Preconditions: map[string]any{"version": chains[0].RuntimeVersion}})
	sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: options.Request.RunID, Type: "run"})
	if err != nil {
		return BudgetDenial{}, err
	}
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "run.execution-blocked", ActorV2: journal.EventActor{ActorID: binding.actor.ActorID, ActorType: binding.actor.ActorType, Producer: binding.actor.Producer, ProvenanceRef: binding.actor.ProvenanceRef, SessionRef: binding.actor.SessionRef}, Aggregate: journal.Aggregate{ID: options.Request.RunID, Sequence: sequence, Type: "run"}, CommandName: "gate.run", CorrelationID: &options.Request.RunID, Repository: repositoryEventValue(options.Request.Repository), Effects: effects}
	prepared, err := journal.PrepareEvent(lock, event)
	if err != nil {
		return BudgetDenial{}, err
	}
	fillCreatedEventIDs(prepared.Effects, prepared.EventID)
	denial := BudgetDenial{Budget: blocks[0].Budget, EventID: prepared.EventID, EventOccurredAt: prepared.OccurredAt.UTC().Format(time.RFC3339Nano)}
	for _, stop := range stops {
		denial.StopIDs = append(denial.StopIDs, stop.ID)
	}
	for _, item := range attention {
		denial.AttentionIDs = append(denial.AttentionIDs, item.ID)
	}
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"denial": denial, "idempotency_key": options.IdempotencyKey, "replayed": false})
	if err != nil {
		return BudgetDenial{}, err
	}
	idempotency := journal.IdempotencyRequest{OperationScope: "gate.run:" + binding.projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
	if _, err := journal.CommitIdempotent(lock, conn, idempotency, prepared, func(journal.Event) (journal.IdempotencyMutation, error) {
		return journal.IdempotencyMutation{EntityID: denial.StopIDs[0], ResultJSON: string(payload)}, nil
	}); err != nil {
		return BudgetDenial{}, err
	}
	return denial, nil
}

func replayGateRunResult(conn *sql.DB, options GateRunOptions, outcome journal.IdempotencyOutcome, result *GateRunResult) error {
	var stored struct {
		Attempt *Attempt      `json:"attempt"`
		Denial  *BudgetDenial `json:"denial"`
		Receipt *Receipt      `json:"receipt"`
	}
	if err := json.Unmarshal([]byte(outcome.ResultJSON), &stored); err != nil {
		return fmt.Errorf("decode gate replay result: %w", err)
	}
	if stored.Denial != nil {
		stored.Denial.EventID = outcome.EventID
		stored.Denial.EventOccurredAt, _ = gateEventOccurredAt(options.WorkspaceRoot, outcome.EventID)
		*result = GateRunResult{Denial: stored.Denial, EventID: outcome.EventID, EventOccurredAt: stored.Denial.EventOccurredAt, IdempotencyKey: options.IdempotencyKey, Replayed: true}
		return &BudgetAdmissionDeniedError{Code: budgetDenialCode(stored.Denial), Denial: *stored.Denial}
	}
	if stored.Receipt == nil {
		if stored.Attempt == nil {
			return fmt.Errorf("decode gate replay result: missing Attempt and Receipt")
		}
		*result = GateRunResult{Attempt: *stored.Attempt, EventID: outcome.EventID, IdempotencyKey: options.IdempotencyKey, Replayed: true}
		return fmt.Errorf("run gate: Attempt %q has no observed terminal Receipt", stored.Attempt.ID)
	}
	receipt, err := loadReceipt(conn, outcome.EntityID)
	if err != nil {
		return err
	}
	occurredAt, err := gateEventOccurredAt(options.WorkspaceRoot, outcome.EventID)
	if err != nil {
		return err
	}
	*result = GateRunResult{Attempt: Attempt{ID: receipt.AttemptID}, Receipt: &receipt, Live: LiveResult{Outcome: receipt.Outcome, ExitCode: receipt.ExitCode, Signal: receipt.Signal, TimedOut: receipt.TimedOut, OutputBytes: receipt.OutputBytes, OutputTruncated: receipt.OutputTruncated, RedactionApplied: receipt.RedactionApplied, RejectedArtifacts: receipt.RejectedArtifacts, executionID: receipt.AttemptID}, EventID: outcome.EventID, EventOccurredAt: occurredAt, IdempotencyKey: options.IdempotencyKey, Replayed: true}
	return nil
}

func budgetDenialCode(denial *BudgetDenial) string {
	if denial != nil && denial.Budget.Remaining != nil && *denial.Budget.Remaining == 0 {
		return "budget-exhausted"
	}
	return "budget-data-unavailable"
}

func fillCreatedEventIDs(effects []journal.Effect, eventID string) {
	for index := range effects {
		if effects[index].Values["created_event_id"] == "<event-id>" {
			effects[index].Values["created_event_id"] = eventID
		}
	}
}

func appendReceipt(lock *journal.WriteLock, conn *sql.DB, binding gateRunBinding, attempt Attempt, options GateRunOptions, requestDigest string, live *LiveResult) (Receipt, string, error) {
	if live == nil {
		return Receipt{}, "", fmt.Errorf("run gate: live result is required")
	}
	if live.Outcome == "" {
		return Receipt{}, "", fmt.Errorf("run gate: process termination was not observed")
	}
	id, err := journal.GenerateID("RCP", func(candidate string) (bool, error) {
		var count int
		err := conn.QueryRow(`SELECT COUNT(*) FROM execution_receipts WHERE id = ?`, candidate).Scan(&count)
		return count > 0, err
	})
	if err != nil {
		return Receipt{}, "", err
	}
	endedAt := time.Now().UTC()
	output := liveOutput(*live)
	digest := sha256.Sum256([]byte(output))
	excerpt := output
	if int64(len(excerpt)) > binding.registration.Policy.ExcerptBytes {
		excerpt = truncateUTF8(excerpt, int(binding.registration.Policy.ExcerptBytes))
	}
	repositoryJSON, err := journal.MarshalCanonicalJSON(options.Request.Repository)
	if err != nil {
		return Receipt{}, "", err
	}
	toolchain := map[string]string{"executable": binding.registration.Policy.Argv[0], "goarch": runtime.GOARCH, "goos": runtime.GOOS}
	toolchainJSON, err := journal.MarshalCanonicalJSON(toolchain)
	if err != nil {
		return Receipt{}, "", err
	}
	retained, rejected, err := retainArtifacts(options.WorkspaceRoot, binding.registration.Policy)
	if err != nil {
		return Receipt{}, "", err
	}
	for _, rejection := range rejected {
		if !containsArtifactRejection(live.RejectedArtifacts, rejection) {
			live.RejectedArtifacts = append(live.RejectedArtifacts, rejection)
			live.appendControl(map[string]any{"kind": "artifact-rejected", "path": rejection.Path, "reason": rejection.Reason})
		}
	}
	artifacts, effects, err := receiptArtifactEffects(conn, retained, endedAt)
	if err != nil {
		return Receipt{}, "", err
	}
	rejectedArtifacts := append([]ArtifactRejection{}, live.RejectedArtifacts...)
	sortArtifactRejections(rejectedArtifacts)
	live.RejectedArtifacts = append([]ArtifactRejection{}, rejectedArtifacts...)
	rejectedJSON, err := journal.MarshalCanonicalJSON(rejectedArtifacts)
	if err != nil {
		return Receipt{}, "", fmt.Errorf("run gate: canonicalize rejected Artifacts: %w", err)
	}
	receipt := Receipt{ID: id, AttemptID: attempt.ID, RegistrationID: attempt.RegistrationID, RegistrationVersion: attempt.RegistrationVersion, ObligationID: attempt.ObligationID, RunID: attempt.RunID, ActorID: attempt.ActorID, TrustSource: "fid-executed", Outcome: live.Outcome, ExitCode: live.ExitCode, Signal: live.Signal, TimedOut: live.TimedOut, StartedAt: attempt.StartedAt, EndedAt: endedAt.Format(time.RFC3339Nano), TerminalRepository: options.Request.Repository, Toolchain: toolchain, OutputExcerpt: excerpt, OutputDigest: "sha256:" + hex.EncodeToString(digest[:]), OutputBytes: live.OutputBytes, OutputTruncated: live.OutputTruncated, RedactionApplied: live.RedactionApplied, Artifacts: artifacts, RejectedArtifacts: rejectedArtifacts}
	chains, err := loadBudgetChains(conn, receipt.RunID)
	if err != nil {
		return Receipt{}, "", err
	}
	charges, err := fidBudgetCharges(chainBudgets(chains), attempt, &receipt)
	if err != nil {
		return Receipt{}, "", err
	}
	budgetEffects, err := appendFIDBudgetEffects(conn, chains, charges, receipt.RunID, binding.actor.ActorID)
	if err != nil {
		return Receipt{}, "", err
	}
	effects = append(effects, budgetEffects...)
	event := receiptEvent(binding.actor, receipt, string(repositoryJSON), string(toolchainJSON), string(rejectedJSON), options.Request.Repository, effects)
	resultJSON, err := journal.MarshalCanonicalJSON(map[string]any{"receipt": receipt, "idempotency_key": options.IdempotencyKey, "replayed": false})
	if err != nil {
		return Receipt{}, "", err
	}
	event.Idempotency = &journal.EventIdempotency{Key: options.IdempotencyKey, RequestDigest: requestDigest, Scope: "gate.run:" + binding.projectID}
	event.Effects = append(event.Effects, journal.Effect{Op: "update", Table: "idempotency_results", PK: map[string]any{"operation_scope": "gate.run:" + binding.projectID, "idempotency_key": options.IdempotencyKey}, Values: map[string]any{"result_entity_id": receipt.ID, "result_event_id": "<event-id>", "aggregate_sequence": 1, "result_json": string(resultJSON)}, Preconditions: map[string]any{"request_digest": requestDigest}})
	prepared, err := journal.PrepareEvent(lock, event)
	if err != nil {
		return Receipt{}, "", err
	}
	fillCreatedEventIDs(prepared.Effects, prepared.EventID)
	fillReceiptIdempotencyEventValues(prepared.Effects, prepared.EventID)
	appended, err := journal.CommitEvent(lock, conn, prepared)
	if err != nil {
		return Receipt{}, "", err
	}
	return receipt, appended.EventID, nil
}

func fillAttemptIdempotencyEventValues(effects []journal.Effect, eventID string, occurredAt time.Time) {
	for index := range effects {
		if effects[index].Table != "idempotency_results" || effects[index].Values["result_event_id"] != "<event-id>" {
			continue
		}
		effects[index].Values["result_event_id"] = eventID
		effects[index].Values["created_at"] = occurredAt.UTC().Format(time.RFC3339Nano)
	}
}

func fillReceiptIdempotencyEventValues(effects []journal.Effect, eventID string) {
	for index := range effects {
		if effects[index].Table == "idempotency_results" && effects[index].Values["result_event_id"] == "<event-id>" {
			effects[index].Values["result_event_id"] = eventID
		}
	}
}

func receiptEvent(actor storedActor, receipt Receipt, repositoryJSON, toolchainJSON, rejectedJSON string, repository Repository, artifactEffects []journal.Effect) journal.Event {
	effects := append([]journal.Effect{}, artifactEffects...)
	effects = append(effects, journal.Effect{Op: "insert", Table: "execution_receipts", PK: map[string]any{"id": receipt.ID}, Values: map[string]any{"attempt_id": receipt.AttemptID, "trust_source": receipt.TrustSource, "terminal_status": receipt.Outcome, "exit_code": receipt.ExitCode, "signal": receipt.Signal, "timed_out": receipt.TimedOut, "started_at": receipt.StartedAt, "ended_at": receipt.EndedAt, "terminal_repository_json": repositoryJSON, "toolchain_json": toolchainJSON, "output_excerpt": receipt.OutputExcerpt, "output_digest": receipt.OutputDigest, "output_bytes": receipt.OutputBytes, "output_truncated": receipt.OutputTruncated, "redaction_applied": receipt.RedactionApplied, "rejected_artifacts_json": rejectedJSON}})
	for _, artifact := range receipt.Artifacts {
		effects = append(effects, journal.Effect{Op: "insert", Table: "receipt_artifacts", PK: map[string]any{"receipt_id": receipt.ID, "artifact_id": artifact.ID}, Values: map[string]any{"declared_path": artifact.DeclaredPath}})
	}
	return journal.Event{EventSchema: journal.EventSchemaV2, EventType: "gate.receipt-recorded", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: receipt.ID, Sequence: 1, Type: "receipt"}, CommandName: "gate.run", CorrelationID: &receipt.RunID, Repository: repositoryEventValue(repository), Effects: effects}
}

func liveOutput(live LiveResult) string {
	var output strings.Builder
	for _, frame := range live.Frames {
		if frame.Stream == "stdout" || frame.Stream == "stderr" {
			output.WriteString(frame.Data["text"].(string))
		}
	}
	return output.String()
}
