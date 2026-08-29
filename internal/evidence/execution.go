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
	Artifacts           []any               `json:"artifacts"`
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
	Live            LiveResult
	EventID         string
	EventOccurredAt string
	IdempotencyKey  string
	Replayed        bool
	definition      GateDefinition
}

type gateRunBinding struct {
	registration Registration
	actor        storedActor
	projectID    string
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
			receipt, err := loadReceipt(conn, outcome.EntityID)
			if err != nil {
				return err
			}
			occurredAt, err := gateEventOccurredAt(options.WorkspaceRoot, outcome.EventID)
			if err != nil {
				return err
			}
			prepared = GateRunResult{Attempt: Attempt{ID: receipt.AttemptID}, Receipt: &receipt, Live: LiveResult{Outcome: "succeeded", ExitCode: receipt.ExitCode, Signal: receipt.Signal, TimedOut: receipt.TimedOut, OutputBytes: receipt.OutputBytes, OutputTruncated: receipt.OutputTruncated, RedactionApplied: receipt.RedactionApplied, RejectedArtifacts: receipt.RejectedArtifacts, executionID: receipt.AttemptID}, EventID: outcome.EventID, EventOccurredAt: occurredAt, IdempotencyKey: options.IdempotencyKey, Replayed: true}
			return nil
		}

		if _, _, cancel, err := executorCommand(ctx, options.WorkspaceRoot, binding.registration.Policy); err != nil {
			return err
		} else {
			cancel()
		}
		attempt, err := appendAttempt(lock, conn, binding, options.Request)
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
		receipt, eventID, err := appendReceipt(lock, conn, binding, prepared.Attempt, options, requestDigest, live)
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

func appendAttempt(lock *journal.WriteLock, conn *sql.DB, binding gateRunBinding, request GateRunRequest) (Attempt, error) {
	id, err := journal.GenerateID("ATT", func(candidate string) (bool, error) {
		var count int
		err := conn.QueryRow(`SELECT COUNT(*) FROM execution_attempts WHERE id = ?`, candidate).Scan(&count)
		return count > 0, err
	})
	if err != nil {
		return Attempt{}, err
	}
	now := time.Now().UTC()
	repositoryJSON, err := journal.MarshalCanonicalJSON(request.Repository)
	if err != nil {
		return Attempt{}, err
	}
	attempt := Attempt{ID: id, RegistrationID: binding.registration.ID, RegistrationVersion: binding.registration.Version, ObligationID: request.ObligationID, RunID: request.RunID, ActorID: binding.actor.ActorID, StartRepository: request.Repository, ResolvedPolicyDigest: binding.registration.DefinitionDigest, StartedAt: now.Format(time.RFC3339Nano)}
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "gate.execution-started", ActorV2: journal.EventActor{ActorID: binding.actor.ActorID, ActorType: binding.actor.ActorType, Producer: binding.actor.Producer, ProvenanceRef: binding.actor.ProvenanceRef, SessionRef: binding.actor.SessionRef}, Aggregate: journal.Aggregate{ID: attempt.ID, Sequence: 1, Type: "attempt"}, CommandName: "gate.run", CorrelationID: &attempt.RunID, Repository: repositoryEventValue(request.Repository), Effects: []journal.Effect{{Op: "insert", Table: "execution_attempts", PK: map[string]any{"id": attempt.ID}, Values: map[string]any{"registration_id": attempt.RegistrationID, "registration_version": attempt.RegistrationVersion, "obligation_id": attempt.ObligationID, "run_id": attempt.RunID, "actor_id": attempt.ActorID, "start_repository_json": string(repositoryJSON), "resolved_policy_digest": attempt.ResolvedPolicyDigest, "producer": "fid", "started_at": attempt.StartedAt}}}}
	if _, err := journal.CommitEvent(lock, conn, event); err != nil {
		return Attempt{}, err
	}
	return attempt, nil
}

func appendReceipt(lock *journal.WriteLock, conn *sql.DB, binding gateRunBinding, attempt Attempt, options GateRunOptions, requestDigest string, live LiveResult) (Receipt, string, error) {
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
	output := liveOutput(live)
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
	receipt := Receipt{ID: id, AttemptID: attempt.ID, RegistrationID: attempt.RegistrationID, RegistrationVersion: attempt.RegistrationVersion, ObligationID: attempt.ObligationID, RunID: attempt.RunID, ActorID: attempt.ActorID, TrustSource: "fid-executed", Outcome: live.Outcome, ExitCode: live.ExitCode, Signal: live.Signal, TimedOut: live.TimedOut, StartedAt: attempt.StartedAt, EndedAt: endedAt.Format(time.RFC3339Nano), TerminalRepository: options.Request.Repository, Toolchain: toolchain, OutputExcerpt: excerpt, OutputDigest: "sha256:" + hex.EncodeToString(digest[:]), OutputBytes: live.OutputBytes, OutputTruncated: live.OutputTruncated, RedactionApplied: live.RedactionApplied, Artifacts: []any{}, RejectedArtifacts: append([]ArtifactRejection(nil), live.RejectedArtifacts...)}
	event := receiptEvent(binding.actor, receipt, string(repositoryJSON), string(toolchainJSON), options.Request.Repository)
	if receipt.Outcome != "succeeded" {
		appended, err := journal.CommitEvent(lock, conn, event)
		if err != nil {
			return Receipt{}, "", err
		}
		return receipt, appended.EventID, nil
	}
	resultJSON, err := journal.MarshalCanonicalJSON(map[string]any{"receipt": receipt, "idempotency_key": options.IdempotencyKey, "replayed": false})
	if err != nil {
		return Receipt{}, "", err
	}
	idempotency := journal.IdempotencyRequest{OperationScope: "gate.run:" + binding.projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
	outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(journal.Event) (journal.IdempotencyMutation, error) {
		return journal.IdempotencyMutation{EntityID: receipt.ID, ResultJSON: string(resultJSON)}, nil
	})
	if err != nil {
		return Receipt{}, "", err
	}
	return receipt, outcome.EventID, nil
}

func receiptEvent(actor storedActor, receipt Receipt, repositoryJSON, toolchainJSON string, repository Repository) journal.Event {
	return journal.Event{EventSchema: journal.EventSchemaV2, EventType: "gate.receipt-recorded", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: receipt.ID, Sequence: 1, Type: "receipt"}, CommandName: "gate.run", CorrelationID: &receipt.RunID, Repository: repositoryEventValue(repository), Effects: []journal.Effect{{Op: "insert", Table: "execution_receipts", PK: map[string]any{"id": receipt.ID}, Values: map[string]any{"attempt_id": receipt.AttemptID, "trust_source": receipt.TrustSource, "terminal_status": receipt.Outcome, "exit_code": receipt.ExitCode, "signal": receipt.Signal, "timed_out": receipt.TimedOut, "started_at": receipt.StartedAt, "ended_at": receipt.EndedAt, "terminal_repository_json": repositoryJSON, "toolchain_json": toolchainJSON, "output_excerpt": receipt.OutputExcerpt, "output_digest": receipt.OutputDigest, "output_bytes": receipt.OutputBytes, "output_truncated": receipt.OutputTruncated, "redaction_applied": receipt.RedactionApplied}}}}
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

func loadReceipt(conn *sql.DB, id string) (Receipt, error) {
	var receipt Receipt
	var repositoryJSON, toolchainJSON string
	err := conn.QueryRow(`SELECT r.id, r.attempt_id, a.registration_id, a.registration_version, a.obligation_id, a.run_id, a.actor_id, r.trust_source, r.terminal_status, r.exit_code, r.signal, r.timed_out, r.started_at, r.ended_at, r.terminal_repository_json, r.toolchain_json, r.output_excerpt, r.output_digest, r.output_bytes, r.output_truncated, r.redaction_applied FROM execution_receipts r JOIN execution_attempts a ON a.id = r.attempt_id WHERE r.id = ?`, id).Scan(&receipt.ID, &receipt.AttemptID, &receipt.RegistrationID, &receipt.RegistrationVersion, &receipt.ObligationID, &receipt.RunID, &receipt.ActorID, &receipt.TrustSource, &receipt.Outcome, &receipt.ExitCode, &receipt.Signal, &receipt.TimedOut, &receipt.StartedAt, &receipt.EndedAt, &repositoryJSON, &toolchainJSON, &receipt.OutputExcerpt, &receipt.OutputDigest, &receipt.OutputBytes, &receipt.OutputTruncated, &receipt.RedactionApplied)
	if err == sql.ErrNoRows {
		return Receipt{}, fmt.Errorf("load receipt %q: not found", id)
	}
	if err != nil {
		return Receipt{}, fmt.Errorf("load receipt %q: %w", id, err)
	}
	if err := json.Unmarshal([]byte(repositoryJSON), &receipt.TerminalRepository); err != nil {
		return Receipt{}, fmt.Errorf("load receipt repository: %w", err)
	}
	if err := json.Unmarshal([]byte(toolchainJSON), &receipt.Toolchain); err != nil {
		return Receipt{}, fmt.Errorf("load receipt toolchain: %w", err)
	}
	receipt.Artifacts = []any{}
	return receipt, nil
}
