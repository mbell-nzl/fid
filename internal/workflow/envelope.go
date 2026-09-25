package workflow

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

type Repository struct {
	BaseOID       *string `json:"base_oid"`
	CaptureSource string  `json:"capture_source"`
	Dirty         bool    `json:"dirty"`
	DirtyDigest   *string `json:"dirty_digest"`
	HeadOID       string  `json:"head_oid"`
	ProjectID     string  `json:"project_id"`
	TreeOID       string  `json:"tree_oid"`
	WorktreeID    string  `json:"worktree_id"`
}

type Source struct {
	Digest       string     `json:"digest"`
	DisplayLabel string     `json:"display_label"`
	ExternalID   *string    `json:"external_id"`
	Path         *string    `json:"path"`
	Repository   Repository `json:"repository"`
	Section      *string    `json:"section"`
	SourceType   string     `json:"source_type"`
}

type AuthorityGrant struct {
	Actions       []string `json:"actions"`
	ActorSelector struct {
		ActorID              string   `json:"actor_id"`
		RequiredCapabilities []string `json:"required_capabilities"`
	} `json:"actor_selector"`
	Scope struct {
		IssueID string `json:"issue_id"`
	} `json:"scope"`
}

type Budget struct {
	Name             string `json:"name"`
	Scope            string `json:"scope"`
	Unit             string `json:"unit"`
	Maximum          int64  `json:"maximum"`
	WarningThreshold *int64 `json:"warning_threshold"`
	ExhaustionAction string `json:"exhaustion_action"`
	Consumed         *int64 `json:"consumed"`
	Remaining        *int64 `json:"remaining"`
	State            string `json:"state,omitempty"`
}

type ObligationInput struct {
	ID             string   `json:"id"`
	Type           string   `json:"type"`
	Applicability  string   `json:"applicability"`
	GuardedActions []string `json:"guarded_actions"`
	Subject        struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"subject"`
	Registration struct {
		ID      string `json:"id"`
		Version int    `json:"version"`
	} `json:"registration"`
	AcceptedEvidenceTypes []string `json:"accepted_evidence_types"`
	AcceptedTrustSources  []string `json:"accepted_trust_sources"`
	FreshnessTarget       string   `json:"freshness_target"`
	ArtifactPolicy        string   `json:"artifact_policy"`
	Independence          string   `json:"independence"`
	RequiredLane          *string  `json:"required_lane"`
	ReviewProfile         *string  `json:"review_profile,omitempty"`
	OutcomeClass          *string  `json:"outcome_class,omitempty"`
	Source                Source   `json:"source"`
}

// MarshalJSON preserves the wire contract's nullable registration field while
// retaining the Phase 21 value-type field for existing callers.
func (input ObligationInput) MarshalJSON() ([]byte, error) {
	type rawObligationInput ObligationInput
	payload, err := json.Marshal(rawObligationInput(input))
	if err != nil || input.Registration.ID != "" || input.Registration.Version != 0 {
		return payload, err
	}
	fields := map[string]json.RawMessage{}
	if err := json.Unmarshal(payload, &fields); err != nil {
		return nil, err
	}
	fields["registration"] = json.RawMessage("null")
	return json.Marshal(fields)
}

type Obligation struct {
	ObligationInput
	DerivedStatus string `json:"derived_status"`
	Reasons       []any  `json:"reasons"`
}

func (obligation Obligation) MarshalJSON() ([]byte, error) {
	payload, err := json.Marshal(obligation.ObligationInput)
	if err != nil {
		return nil, err
	}
	fields := map[string]json.RawMessage{}
	if err := json.Unmarshal(payload, &fields); err != nil {
		return nil, err
	}
	status, err := json.Marshal(obligation.DerivedStatus)
	if err != nil {
		return nil, err
	}
	reasons, err := json.Marshal(obligation.Reasons)
	if err != nil {
		return nil, err
	}
	fields["derived_status"] = status
	fields["reasons"] = reasons
	return json.Marshal(fields)
}

type EnvelopeCreateRequest struct {
	Outcome    string   `json:"outcome"`
	NonGoals   []string `json:"non_goals"`
	RiskResult struct {
		Source       Source `json:"source"`
		ResultDigest string `json:"result_digest"`
	} `json:"risk_result"`
	RepositoryScope    Repository               `json:"repository_scope"`
	Sources            []Source                 `json:"sources"`
	PrimaryIssueID     string                   `json:"primary_issue_id"`
	RelatedIssueIDs    []string                 `json:"related_issue_ids"`
	AuthorityGrants    []AuthorityGrant         `json:"authority_grants"`
	Obligations        []ObligationInput        `json:"obligations"`
	ProtectedDecisions []ProtectedDecisionInput `json:"protected_decisions"`
	Budgets            []Budget                 `json:"budgets"`
}

type ProtectedDecisionInput struct {
	Ordinal              int      `json:"ordinal"`
	DecisionClass        string   `json:"decision_class"`
	RequiredCapabilities []string `json:"required_capabilities"`
	AffectedActions      []string `json:"affected_actions"`
}

type ProtectedDecision struct {
	ProtectedDecisionInput
	Status         string  `json:"status"`
	DecisionRefID  *string `json:"decision_ref_id"`
	HumanOutcomeID *string `json:"human_outcome_id"`
}

type EnvelopeActivateRequest struct {
	EnvelopeID       string `json:"envelope_id"`
	ExpectedRevision int    `json:"expected_revision"`
}

type EnvelopeCompleteRequest struct {
	EnvelopeID       string     `json:"envelope_id"`
	ExpectedRevision int        `json:"expected_revision"`
	Repository       Repository `json:"repository"`
}

type ObligationEvaluation struct {
	ObligationID         string             `json:"obligation_id"`
	Status               string             `json:"status"`
	Validity             string             `json:"validity"`
	Freshness            string             `json:"freshness"`
	ArtifactAvailability string             `json:"artifact_availability"`
	Reasons              []EvaluationReason `json:"reasons"`
}

type Envelope struct {
	ID                 string              `json:"id"`
	CurrentRevision    int                 `json:"current_revision"`
	Status             string              `json:"status"`
	Outcome            string              `json:"outcome"`
	NonGoals           []string            `json:"non_goals"`
	RiskResult         json.RawMessage     `json:"risk_result"`
	RepositoryScope    Repository          `json:"repository_scope"`
	Sources            []Source            `json:"sources"`
	PrimaryIssueID     string              `json:"primary_issue_id"`
	RelatedIssueIDs    []string            `json:"related_issue_ids"`
	AuthorityGrants    []AuthorityGrant    `json:"authority_grants"`
	Obligations        []Obligation        `json:"obligations"`
	ProtectedDecisions []ProtectedDecision `json:"protected_decisions"`
	Budgets            []Budget            `json:"budgets"`
	CreatedActorID     string              `json:"created_actor_id"`
	ApprovedActorID    *string             `json:"approved_actor_id"`
	CreatedAt          string              `json:"created_at"`
	ActivatedAt        *string             `json:"activated_at"`
	CompletedAt        *string             `json:"completed_at"`
}

type EnvelopeMutationResult struct {
	Envelope       Envelope `json:"envelope"`
	Evaluation     any      `json:"evaluation,omitempty"`
	IdempotencyKey string   `json:"idempotency_key"`
	Replayed       bool     `json:"replayed"`
	EventID        string
}

type EnvelopeOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
}

type EnvelopeNotFoundError struct{ ID string }

func (err *EnvelopeNotFoundError) Error() string {
	return fmt.Sprintf("envelope %q was not found", err.ID)
}

type EnvelopeInvalidError struct{ Reason string }

func (err *EnvelopeInvalidError) Error() string { return err.Reason }

type EnvelopeCompletionError struct {
	EnvelopeID string
	Reasons    []EvaluationReason
}

func (err *EnvelopeCompletionError) Error() string {
	return fmt.Sprintf("complete envelope %q is not permitted", err.EnvelopeID)
}

type EnvelopeAuthorityError struct {
	ActorID    string
	Capability string
}

func (err *EnvelopeAuthorityError) Error() string {
	return fmt.Sprintf("actor %q lacks %s", err.ActorID, err.Capability)
}

func CreateEnvelope(options EnvelopeOptions, request EnvelopeCreateRequest) (EnvelopeMutationResult, error) {
	if err := validateEnvelopeOptions(options); err != nil {
		return EnvelopeMutationResult{}, err
	}
	if err := validateEnvelopeCreate(request); err != nil {
		return EnvelopeMutationResult{}, err
	}
	var result EnvelopeMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid envelope create"}, func(lock *journal.WriteLock) error {
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
			return fmt.Errorf("create envelope: read project config: %w", err)
		}
		actor, err := envelopeActor(conn, options.ActorID, "envelope.create")
		if err != nil {
			return err
		}
		if err := validateEnvelopeBindings(conn, project.ProjectID, request); err != nil {
			return err
		}
		envelopeID, err := journal.GenerateID("ENV", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM work_envelopes WHERE id = ?`, candidate).Scan(&found)
			return err == nil, nil
		})
		if err != nil {
			return err
		}
		createdAt := time.Now().UTC().Format(time.RFC3339Nano)
		envelope := envelopeFromCreate(envelopeID, actor.ActorID, createdAt, request)
		if err := deriveInitialEnvelopeBudgets(&envelope); err != nil {
			return err
		}
		requestDigest, err := envelopeRequestDigest("envelope.create", project.ProjectID, actor.ActorID, request)
		if err != nil {
			return err
		}
		resultJSON, err := envelopeResultJSON(envelope, options.IdempotencyKey)
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "envelope.created", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: envelopeID, Sequence: 1, Type: "envelope"}, CommandName: "envelope.create", Repository: envelopeRepository(request.RepositoryScope), Effects: envelopeCreateEffects(envelope, request)}
		outcome, err := journal.CommitIdempotent(lock, conn, journal.IdempotencyRequest{OperationScope: "envelope.create:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: envelopeID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		if outcome.Replayed {
			return loadEnvelopeMutationReplay(outcome, options.IdempotencyKey, &result)
		}
		loaded, err := loadEnvelope(conn, outcome.EntityID)
		if err != nil {
			return err
		}
		result = EnvelopeMutationResult{Envelope: loaded, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func GetEnvelope(workspaceRoot, migrationsDir, id string) (Envelope, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Envelope{}, err
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return Envelope{}, err
	}
	defer closeWorkflowDB(conn)
	return loadEnvelope(conn, strings.TrimSpace(id))
}

func ActivateEnvelope(options EnvelopeOptions, request EnvelopeActivateRequest) (EnvelopeMutationResult, error) {
	if err := validateEnvelopeOptions(options); err != nil || strings.TrimSpace(request.EnvelopeID) == "" || request.ExpectedRevision != 1 {
		return EnvelopeMutationResult{}, &EnvelopeInvalidError{Reason: "envelope activation input is invalid"}
	}
	var result EnvelopeMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid envelope activate"}, func(lock *journal.WriteLock) error {
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
			return err
		}
		actor, err := envelopeActor(conn, options.ActorID, "envelope.activate")
		if err != nil {
			return err
		}
		requestDigest, err := envelopeRequestDigest("envelope.activate", project.ProjectID, actor.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "envelope.activate:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return loadEnvelopeMutationReplay(outcome, options.IdempotencyKey, &result)
		}
		envelope, err := loadEnvelope(conn, request.EnvelopeID)
		if err != nil {
			return err
		}
		if envelope.Status != "draft" {
			return &EnvelopeInvalidError{Reason: "envelope is not draft"}
		}
		var activeClaim int
		if err := conn.QueryRow(`SELECT COUNT(*) FROM claims WHERE issue_id = ? AND status = 'active'`, envelope.PrimaryIssueID).Scan(&activeClaim); err != nil {
			return fmt.Errorf("activate envelope: read active claims: %w", err)
		}
		if activeClaim != 0 {
			return &EnvelopeInvalidError{Reason: "active mutation claim blocks envelope activation"}
		}
		activatedAt := time.Now().UTC().Format(time.RFC3339Nano)
		stopEffects, err := protectedDecisionStopEffects(conn, envelope, options.ActorID)
		if err != nil {
			return err
		}
		// A fresh git observation is unnecessary here: the only Obligation type
		// that can be raised at activation time is human-outcome
		// (freshness_target "none", artifact_policy "metadata-sufficient"), so
		// its evaluation never consults the repository. Reusing the Envelope's
		// already-loaded RepositoryScope avoids a new git dependency on
		// envelope.activate, which is not otherwise repository-bound.
		readinessEffects, _, _, err := humanOutcomeReadinessEffectsForEnvelope(conn, envelope, options.ActorID, envelope.RepositoryScope, nil)
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "envelope.activated", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: envelope.ID, Sequence: 2, Type: "envelope"}, CommandName: "envelope.activate", Repository: envelopeRepository(envelope.RepositoryScope), Effects: []journal.Effect{
			{Op: "update", Table: "work_envelopes", PK: map[string]any{"id": envelope.ID}, Values: map[string]any{"status": "active"}},
			{Op: "update", Table: "envelope_revisions", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1}, Values: map[string]any{"status": "active", "activated_at": activatedAt}},
		}}
		event.Effects = append(event.Effects, stopEffects...)
		event.Effects = append(event.Effects, readinessEffects...)
		resultJSON, err := envelopeActivateResultJSON(envelope, activatedAt, options.IdempotencyKey)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			setEnvelopeEventIdentity(prepared.Effects, prepared)
			return journal.IdempotencyMutation{EntityID: envelope.ID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		loaded, err := loadEnvelope(conn, outcome.EntityID)
		if err != nil {
			return err
		}
		result = EnvelopeMutationResult{Envelope: loaded, Evaluation: envelopeEvaluation(loaded), IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func CompleteEnvelope(options EnvelopeOptions, request EnvelopeCompleteRequest) (EnvelopeMutationResult, error) {
	if err := validateEnvelopeOptions(options); err != nil || strings.TrimSpace(request.EnvelopeID) == "" || request.ExpectedRevision != 1 || !validEnvelopeRepository(request.Repository) {
		return EnvelopeMutationResult{}, &EnvelopeInvalidError{Reason: "envelope completion input is invalid"}
	}
	var result EnvelopeMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid envelope complete"}, func(lock *journal.WriteLock) error {
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
			return fmt.Errorf("complete envelope: read project config: %w", err)
		}
		requestDigest, err := envelopeRequestDigest("envelope.complete", project.ProjectID, options.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "envelope.complete:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return loadEnvelopeMutationReplay(outcome, options.IdempotencyKey, &result)
		}

		envelope, err := loadEnvelope(conn, request.EnvelopeID)
		if err != nil {
			return err
		}
		currentRepository, err := resumeObservedRepository(options.WorkspaceRoot, options.WorkspaceRoot)
		if err != nil {
			return fmt.Errorf("complete envelope: observe repository: %w", err)
		}
		request.Repository = currentRepository
		reasons, evaluation, err := evaluateEnvelopeCompletion(options, conn, envelope, request)
		if err != nil {
			return err
		}
		if len(reasons) > 0 && (reasons[0].Code == "protected-decision-required" || reasons[0].Code == "human-outcome-required") {
			stops, err := ApplicableStops(conn, "envelope.complete", Subject{Type: "envelope", ID: envelope.ID})
			if err != nil {
				return err
			}
			for _, stop := range stops {
				if stop.Code == reasons[0].Code {
					blocked, err := NewStopBlockedError(conn, "envelope.complete", []Stop{stop})
					if err != nil {
						return err
					}
					return blocked
				}
			}
		}
		if len(reasons) > 0 {
			return &EnvelopeCompletionError{EnvelopeID: envelope.ID, Reasons: reasons}
		}
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		completedAt := time.Now().UTC().Format(time.RFC3339Nano)
		sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "journal"), journal.Aggregate{ID: envelope.ID, Type: "envelope"})
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "envelope.completed", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: envelope.ID, Sequence: sequence, Type: "envelope"}, CommandName: "envelope.complete", Repository: envelopeRepository(request.Repository), Effects: []journal.Effect{
			{Op: "update", Table: "work_envelopes", PK: map[string]any{"id": envelope.ID}, Values: map[string]any{"status": "completed"}, Preconditions: map[string]any{"status": "active"}},
			{Op: "update", Table: "envelope_revisions", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1}, Values: map[string]any{"status": "completed", "completed_at": completedAt}, Preconditions: map[string]any{"status": "active"}},
		}}
		resultJSON, err := envelopeCompleteResultJSON(envelope, completedAt, evaluation, options.IdempotencyKey)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: envelope.ID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		loaded, err := loadEnvelope(conn, outcome.EntityID)
		if err != nil {
			return err
		}
		result = EnvelopeMutationResult{Envelope: loaded, Evaluation: evaluation, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func evaluateEnvelopeCompletion(options EnvelopeOptions, conn *sql.DB, envelope Envelope, request EnvelopeCompleteRequest) ([]EvaluationReason, []ObligationEvaluation, error) {
	reasons := []EvaluationReason{}
	if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: "envelope.complete", Capability: "envelope.complete", IssueID: envelope.PrimaryIssueID, EnvelopeID: envelope.ID, Revision: envelope.CurrentRevision}); err != nil {
		reasons = append(reasons, evaluationReason("authority-unavailable", "envelope.complete", options.ActorID, "actor"))
	}
	if envelope.Status != "active" {
		reasons = append(reasons, evaluationReason("envelope-inactive", envelope.Status, envelope.ID, "envelope"))
	}
	stopReasons, err := stopEvaluationReasons(conn, "envelope.complete", Subject{Type: "envelope", ID: envelope.ID})
	if err != nil {
		return nil, nil, err
	}
	reasons = append(reasons, stopReasons...)
	evaluation, err := currentEnvelopeEvaluations(conn, envelope, request.Repository)
	if err != nil {
		return nil, nil, err
	}
	for _, item := range evaluation {
		reasons = append(reasons, item.Reasons...)
	}
	terminal, err := issueTerminalSuccess(conn, envelope.PrimaryIssueID)
	if err != nil {
		return nil, nil, err
	}
	if !terminal {
		reasons = append(reasons, evaluationReason("issue-not-terminal-success", "completed workflow category", envelope.PrimaryIssueID, "issue"))
	}
	return sortEvaluationReasons(reasons), evaluation, nil
}

func currentEnvelopeEvaluation(conn *sql.DB, envelope Envelope, repository Repository) (ObligationEvaluation, error) {
	evaluations, err := currentEnvelopeEvaluations(conn, envelope, repository)
	if err != nil {
		return ObligationEvaluation{}, err
	}
	return evaluations[0], nil
}

func currentEnvelopeEvaluations(conn *sql.DB, envelope Envelope, repository Repository) ([]ObligationEvaluation, error) {
	states := map[string][]ObligationEvidenceState{}
	for _, obligation := range envelope.Obligations {
		loaded, err := loadObligationEvidenceStates(conn, obligation.ID, repository)
		if err != nil {
			return nil, err
		}
		states[obligation.ID] = loaded
	}
	return EvaluateEnvelopeObligations(envelope.Obligations, states), nil
}

func issueTerminalSuccess(conn *sql.DB, issueID string) (bool, error) {
	var workflowID, status string
	if err := conn.QueryRow(`SELECT workflow_id, status FROM issues WHERE id = ?`, issueID).Scan(&workflowID, &status); err != nil {
		return false, fmt.Errorf("load completion Issue: %w", err)
	}
	registry, err := Load(conn)
	if err != nil {
		return false, err
	}
	category, found := registry.CategoryForState(workflowID, status)
	return found && category == "completed", nil
}

func validateEnvelopeOptions(options EnvelopeOptions) error {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.IdempotencyKey) == "" {
		return &EnvelopeInvalidError{Reason: "envelope options are incomplete"}
	}
	return nil
}

func validateEnvelopeCreate(request EnvelopeCreateRequest) error {
	if strings.TrimSpace(request.Outcome) == "" || len(request.Outcome) > 2048 || len(request.NonGoals) > 64 || len(request.Sources) == 0 || len(request.Sources) > 64 || strings.TrimSpace(request.PrimaryIssueID) == "" || len(request.RelatedIssueIDs) != 0 || len(request.AuthorityGrants) == 0 || len(request.AuthorityGrants) > 64 || len(request.Obligations) < 1 || len(request.Obligations) > 16 || len(request.ProtectedDecisions) > 16 || len(request.Budgets) > 64 {
		return &EnvelopeInvalidError{Reason: "envelope input exceeds the one-revision boundary"}
	}
	if !validEnvelopeRepository(request.RepositoryScope) || !validEnvelopeSource(request.RiskResult.Source) || !validEnvelopeDigest(request.RiskResult.ResultDigest) {
		return &EnvelopeInvalidError{Reason: "envelope repository or risk result is invalid"}
	}
	for _, source := range request.Sources {
		if !validEnvelopeSource(source) {
			return &EnvelopeInvalidError{Reason: "envelope source is invalid"}
		}
	}
	previous := ""
	for _, grant := range request.AuthorityGrants {
		if grant.ActorSelector.ActorID == "" || grant.Scope.IssueID != request.PrimaryIssueID || !sortedUniqueEnvelopeActions(grant.Actions) || !sortedUniqueCapabilities(grant.ActorSelector.RequiredCapabilities) || previous >= grant.ActorSelector.ActorID {
			return &EnvelopeInvalidError{Reason: "authority grants must be sorted and exact"}
		}
		previous = grant.ActorSelector.ActorID
	}
	if err := validateEnvelopeObligations(request); err != nil {
		return err
	}
	if err := validateProtectedDecisions(request); err != nil {
		return err
	}
	previous = ""
	hasReviewObligation := false
	for _, obligation := range request.Obligations {
		hasReviewObligation = hasReviewObligation || obligation.Type == "review"
	}
	for _, budget := range request.Budgets {
		reviewCycles := budget.Unit == "review-cycles"
		reviewExhaustion := budget.ExhaustionAction == "review-rework-exhausted"
		if budget.Name == "" || previous >= budget.Name || budget.Scope != "envelope" || !validBudgetUnit(budget.Unit) || budget.Maximum < 1 || (budget.WarningThreshold != nil && (*budget.WarningThreshold < 0 || *budget.WarningThreshold > budget.Maximum)) || reviewCycles != reviewExhaustion || (!reviewCycles && budget.ExhaustionAction != "budget-exhausted") || (reviewCycles && !hasReviewObligation) || budget.Consumed != nil || budget.Remaining != nil || budget.State != "" {
			return &EnvelopeInvalidError{Reason: "budget values must remain immutable null counters"}
		}
		previous = budget.Name
	}
	return nil
}

func validateEnvelopeObligations(request EnvelopeCreateRequest) error {
	ids := map[string]bool{}
	terminalTypes := map[string]bool{}
	reviewProfiles := map[string]bool{}
	outcomeClasses := map[string]bool{}
	gateEvidenceCount := 0
	for _, obligation := range request.Obligations {
		if obligation.ID == "" || ids[obligation.ID] || obligation.Applicability != "required" || obligation.Subject.Type != "issue" || obligation.Subject.ID != request.PrimaryIssueID || obligation.RequiredLane != nil || !validEnvelopeSource(obligation.Source) || !validGuardedActions(obligation.GuardedActions) {
			return &EnvelopeInvalidError{Reason: "obligation is not an accepted value"}
		}
		ids[obligation.ID] = true
		if containsAction(obligation.GuardedActions, "issue.transition-terminal-success") {
			if terminalTypes[obligation.Type] {
				return &EnvelopeInvalidError{Reason: "only one obligation type may guard terminal success"}
			}
			terminalTypes[obligation.Type] = true
		}
		switch obligation.Type {
		case "gate-evidence":
			gateEvidenceCount++
			if gateEvidenceCount > 1 || obligation.Registration.ID == "" || obligation.Registration.Version < 1 || obligation.ReviewProfile != nil || obligation.OutcomeClass != nil || strings.Join(obligation.AcceptedEvidenceTypes, ",") != "receipt" || strings.Join(obligation.AcceptedTrustSources, ",") != "fid-executed" || obligation.FreshnessTarget != "repository-tree" || obligation.ArtifactPolicy != "required-available" || obligation.Independence != "fid-executed" {
				return &EnvelopeInvalidError{Reason: "gate-evidence obligation is not the exact supported value"}
			}
		case "review":
			if obligation.Registration.ID != "" || obligation.Registration.Version != 0 || obligation.ReviewProfile == nil || strings.TrimSpace(*obligation.ReviewProfile) == "" || reviewProfiles[*obligation.ReviewProfile] || obligation.OutcomeClass != nil || strings.Join(obligation.AcceptedEvidenceTypes, ",") != "review" || strings.Join(obligation.AcceptedTrustSources, ",") != "activated-actor" || obligation.FreshnessTarget != "repository-tree" || obligation.ArtifactPolicy != "metadata-sufficient" || obligation.Independence != "distinct-actor" {
				return &EnvelopeInvalidError{Reason: "review obligation is not the exact supported value"}
			}
			reviewProfiles[*obligation.ReviewProfile] = true
		case "human-outcome":
			if obligation.Registration.ID != "" || obligation.Registration.Version != 0 || obligation.ReviewProfile != nil || obligation.OutcomeClass == nil || strings.TrimSpace(*obligation.OutcomeClass) == "" || outcomeClasses[*obligation.OutcomeClass] || strings.Join(obligation.AcceptedEvidenceTypes, ",") != "human-outcome" || strings.Join(obligation.AcceptedTrustSources, ",") != "activated-human" || obligation.FreshnessTarget != "none" || obligation.ArtifactPolicy != "metadata-sufficient" || obligation.Independence != "human" {
				return &EnvelopeInvalidError{Reason: "human-outcome obligation is not the exact supported value"}
			}
			outcomeClasses[*obligation.OutcomeClass] = true
		default:
			return &EnvelopeInvalidError{Reason: "obligation type is not supported"}
		}
	}
	return nil
}

func validateProtectedDecisions(request EnvelopeCreateRequest) error {
	classes := map[string]bool{}
	ordinals := map[int]bool{}
	guarded := map[string]bool{}
	for _, obligation := range request.Obligations {
		for _, action := range obligation.GuardedActions {
			guarded[action] = true
		}
	}
	for _, decision := range request.ProtectedDecisions {
		if decision.Ordinal < 0 || ordinals[decision.Ordinal] || strings.TrimSpace(decision.DecisionClass) == "" || classes[decision.DecisionClass] || len(decision.RequiredCapabilities) == 0 || len(decision.RequiredCapabilities) > 16 || !sortedUniqueCapabilities(decision.RequiredCapabilities) || !validGuardedActions(decision.AffectedActions) {
			return &EnvelopeInvalidError{Reason: "protected decision is not an accepted value"}
		}
		for _, action := range decision.AffectedActions {
			if !guarded[action] {
				return &EnvelopeInvalidError{Reason: "protected decision affects an unguarded action"}
			}
		}
		ordinals[decision.Ordinal], classes[decision.DecisionClass] = true, true
	}
	return nil
}

func validGuardedActions(actions []string) bool {
	if len(actions) == 0 || len(actions) > 2 || !sort.StringsAreSorted(actions) {
		return false
	}
	for index, action := range actions {
		if (action != "envelope.complete" && action != "issue.transition-terminal-success") || (index > 0 && action == actions[index-1]) {
			return false
		}
	}
	return true
}

func containsAction(actions []string, want string) bool {
	for _, action := range actions {
		if action == want {
			return true
		}
	}
	return false
}

func validEnvelopeRepository(repository Repository) bool {
	return repository.CaptureSource == "git" && repository.ProjectID != "" && repository.WorktreeID != "" && validEnvelopeOID(repository.HeadOID) && validEnvelopeOID(repository.TreeOID) && (repository.BaseOID == nil || validEnvelopeOID(*repository.BaseOID)) && ((repository.Dirty && repository.DirtyDigest != nil && validEnvelopeDigest(*repository.DirtyDigest)) || (!repository.Dirty && repository.DirtyDigest == nil))
}

func validEnvelopeSource(source Source) bool {
	return source.DisplayLabel != "" && len(source.DisplayLabel) <= 128 && validEnvelopeDigest(source.Digest) && validEnvelopeRepository(source.Repository) && ((source.SourceType == "project-file" && source.Path != nil && source.ExternalID == nil) || (source.SourceType == "external-reference" && source.Path == nil && source.ExternalID != nil))
}

func validEnvelopeOID(value string) bool {
	return len(value) == 40 && strings.Trim(value, "0123456789abcdef") == ""
}

func validEnvelopeDigest(value string) bool {
	return strings.HasPrefix(value, "sha256:") && len(value) == 71 && strings.TrimPrefix(value, "sha256:") != ""
}

func sortedUniqueEnvelopeActions(actions []string) bool {
	if len(actions) == 0 {
		return false
	}
	previous := ""
	for _, action := range actions {
		if !sortedUniqueCapabilities([]string{action}) || previous >= action {
			return false
		}
		previous = action
	}
	return true
}

func validateEnvelopeBindings(conn *sql.DB, projectID string, request EnvelopeCreateRequest) error {
	if request.RepositoryScope.ProjectID != projectID {
		return &EnvelopeInvalidError{Reason: "envelope repository project does not match"}
	}
	sources := append([]Source{request.RiskResult.Source}, request.Sources...)
	for _, obligation := range request.Obligations {
		sources = append(sources, obligation.Source)
	}
	for _, source := range sources {
		if source.Repository.ProjectID != projectID || source.Repository.WorktreeID != request.RepositoryScope.WorktreeID {
			return &EnvelopeInvalidError{Reason: "envelope source repository does not match"}
		}
	}
	var issueID string
	if err := conn.QueryRow(`SELECT id FROM issues WHERE id = ?`, request.PrimaryIssueID).Scan(&issueID); err == sql.ErrNoRows {
		return &EnvelopeInvalidError{Reason: "primary issue was not found"}
	} else if err != nil {
		return fmt.Errorf("validate envelope issue: %w", err)
	}
	for _, obligation := range request.Obligations {
		if obligation.Type != "gate-evidence" {
			continue
		}
		var registrationProject, registrationWorktree, scopeJSON string
		err := conn.QueryRow(`SELECT project_id, worktree_id, scope_json FROM gate_registrations WHERE id = ? AND version = ?`, obligation.Registration.ID, obligation.Registration.Version).Scan(&registrationProject, &registrationWorktree, &scopeJSON)
		if err == sql.ErrNoRows {
			return &EnvelopeInvalidError{Reason: "prior gate registration was not found"}
		}
		if err != nil {
			return fmt.Errorf("validate envelope registration: %w", err)
		}
		var scope struct {
			IssueID string `json:"issue_id"`
		}
		if err := json.Unmarshal([]byte(scopeJSON), &scope); err != nil {
			return fmt.Errorf("decode envelope registration scope: %w", err)
		}
		if registrationProject != projectID || registrationWorktree != request.RepositoryScope.WorktreeID || scope.IssueID != request.PrimaryIssueID {
			return &EnvelopeInvalidError{Reason: "prior registration binding does not match envelope"}
		}
	}
	return nil
}

func envelopeActor(conn *sql.DB, actorID, capability string) (ActorProfile, error) {
	var activationID, capabilitiesJSON string
	actor := ActorProfile{}
	err := conn.QueryRow(`SELECT id FROM kit_activations WHERE status = 'finalized' LIMIT 1`).Scan(&activationID)
	if err == sql.ErrNoRows {
		return ActorProfile{}, &EnvelopeAuthorityError{ActorID: actorID, Capability: capability}
	}
	if err != nil {
		return ActorProfile{}, fmt.Errorf("load finalized envelope activation: %w", err)
	}
	err = conn.QueryRow(`SELECT actor_id, actor_type, producer, provenance_ref, session_ref, capabilities_json FROM activation_actor_profiles WHERE activation_id = ? AND actor_id = ?`, activationID, actorID).Scan(&actor.ActorID, &actor.ActorType, &actor.Producer, &actor.ProvenanceRef, &actor.SessionRef, &capabilitiesJSON)
	if err == sql.ErrNoRows {
		return ActorProfile{}, &EnvelopeAuthorityError{ActorID: actorID, Capability: capability}
	}
	if err != nil {
		return ActorProfile{}, fmt.Errorf("load envelope actor: %w", err)
	}
	if err := json.Unmarshal([]byte(capabilitiesJSON), &actor.Capabilities); err != nil {
		return ActorProfile{}, fmt.Errorf("decode envelope actor capabilities: %w", err)
	}
	if !hasEnvelopeCapability(actor.Capabilities, capability) {
		return ActorProfile{}, &EnvelopeAuthorityError{ActorID: actorID, Capability: capability}
	}
	return actor, nil
}

func hasEnvelopeCapability(capabilities []string, capability string) bool {
	for _, candidate := range capabilities {
		if candidate == capability {
			return true
		}
	}
	return false
}

func envelopeFromCreate(id, actorID, createdAt string, request EnvelopeCreateRequest) Envelope {
	obligations := make([]Obligation, 0, len(request.Obligations))
	for _, input := range request.Obligations {
		obligations = append(obligations, Obligation{ObligationInput: input, DerivedStatus: "unsatisfied", Reasons: []any{}})
	}
	decisions := make([]ProtectedDecision, 0, len(request.ProtectedDecisions))
	for _, input := range request.ProtectedDecisions {
		decisions = append(decisions, ProtectedDecision{ProtectedDecisionInput: input, Status: "unsatisfied"})
	}
	return Envelope{ID: id, CurrentRevision: 1, Status: "draft", Outcome: request.Outcome, NonGoals: request.NonGoals, RiskResult: mustEnvelopeJSON(request.RiskResult), RepositoryScope: request.RepositoryScope, Sources: request.Sources, PrimaryIssueID: request.PrimaryIssueID, RelatedIssueIDs: []string{}, AuthorityGrants: request.AuthorityGrants, Obligations: obligations, ProtectedDecisions: decisions, Budgets: request.Budgets, CreatedActorID: actorID, ApprovedActorID: nil, CreatedAt: createdAt, ActivatedAt: nil, CompletedAt: nil}
}

func deriveInitialEnvelopeBudgets(envelope *Envelope) error {
	derived := make([]Budget, 0, len(envelope.Budgets))
	for _, budget := range envelope.Budgets {
		value, err := deriveBudget(budget, nil)
		if err != nil {
			return fmt.Errorf("derive initial envelope budget %q: %w", budget.Name, err)
		}
		derived = append(derived, value)
	}
	envelope.Budgets = derived
	return nil
}

func mustEnvelopeJSON(value any) json.RawMessage {
	payload, err := journal.MarshalCanonicalJSON(value)
	if err != nil {
		panic(fmt.Sprintf("marshal validated envelope value: %v", err))
	}
	return payload
}

func envelopeCreateEffects(envelope Envelope, request EnvelopeCreateRequest) []journal.Effect {
	riskJSON := mustEnvelopeJSON(request.RiskResult)
	repositoryJSON := mustEnvelopeJSON(request.RepositoryScope)
	effects := []journal.Effect{
		{Op: "insert", Table: "work_envelopes", PK: map[string]any{"id": envelope.ID}, Values: map[string]any{"current_revision": 1, "status": "draft", "created_at": envelope.CreatedAt}},
		{Op: "insert", Table: "envelope_revisions", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1}, Values: map[string]any{"prior_revision": nil, "outcome": request.Outcome, "non_goals_json": string(mustEnvelopeJSON(request.NonGoals)), "risk_result_json": string(riskJSON), "repository_scope_json": string(repositoryJSON), "status": "draft", "created_actor_id": envelope.CreatedActorID, "approved_actor_id": nil, "created_at": envelope.CreatedAt, "activated_at": nil, "completed_at": nil}},
		{Op: "insert", Table: "envelope_issue_selection", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1, "issue_id": request.PrimaryIssueID}, Values: map[string]any{"primary_issue": 1}},
	}
	for ordinal, source := range request.Sources {
		effects = append(effects, journal.Effect{Op: "insert", Table: "envelope_sources", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1, "ordinal": ordinal}, Values: map[string]any{"source_json": string(mustEnvelopeJSON(source))}})
	}
	for ordinal, grant := range request.AuthorityGrants {
		effects = append(effects, journal.Effect{Op: "insert", Table: "envelope_authority_grants", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1, "ordinal": ordinal}, Values: map[string]any{"actor_selector_json": string(mustEnvelopeJSON(grant.ActorSelector)), "actions_json": string(mustEnvelopeJSON(grant.Actions)), "scope_json": string(mustEnvelopeJSON(grant.Scope))}})
	}
	for _, obligation := range request.Obligations {
		registrationID, registrationVersion := any(nil), any(nil)
		if obligation.Type == "gate-evidence" {
			registrationID, registrationVersion = obligation.Registration.ID, obligation.Registration.Version
		}
		effects = append(effects, journal.Effect{Op: "insert", Table: "envelope_obligations", PK: map[string]any{"id": obligation.ID}, Values: map[string]any{"envelope_id": envelope.ID, "revision": 1, "obligation_type": obligation.Type, "registration_id": registrationID, "registration_version": registrationVersion, "review_profile": obligation.ReviewProfile, "outcome_class": obligation.OutcomeClass, "obligation_json": string(mustEnvelopeJSON(obligation))}})
	}
	for _, decision := range request.ProtectedDecisions {
		effects = append(effects, journal.Effect{Op: "insert", Table: "envelope_protected_decisions", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1, "ordinal": decision.Ordinal}, Values: map[string]any{"decision_class": decision.DecisionClass, "required_capabilities_json": string(mustEnvelopeJSON(decision.RequiredCapabilities)), "affected_actions_json": string(mustEnvelopeJSON(decision.AffectedActions))}})
	}
	for _, budget := range request.Budgets {
		policyBudget := budgetPolicyValue(budget)
		effects = append(effects, journal.Effect{Op: "insert", Table: "envelope_budgets", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1, "name": budget.Name}, Values: map[string]any{"budget_json": string(mustEnvelopeJSON(policyBudget))}})
	}
	return effects
}

func loadEnvelope(conn *sql.DB, id string) (Envelope, error) {
	var envelope Envelope
	var nonGoals, riskResult, repositoryScope string
	err := conn.QueryRow(`SELECT e.id, e.current_revision, e.status, r.outcome, r.non_goals_json, r.risk_result_json, r.repository_scope_json, r.created_actor_id, r.approved_actor_id, r.created_at, r.activated_at, r.completed_at FROM work_envelopes e JOIN envelope_revisions r ON r.envelope_id = e.id AND r.revision = e.current_revision WHERE e.id = ?`, id).Scan(&envelope.ID, &envelope.CurrentRevision, &envelope.Status, &envelope.Outcome, &nonGoals, &riskResult, &repositoryScope, &envelope.CreatedActorID, &envelope.ApprovedActorID, &envelope.CreatedAt, &envelope.ActivatedAt, &envelope.CompletedAt)
	if err == sql.ErrNoRows {
		return Envelope{}, &EnvelopeNotFoundError{ID: id}
	}
	if err != nil {
		return Envelope{}, fmt.Errorf("load envelope: %w", err)
	}
	if err := json.Unmarshal([]byte(nonGoals), &envelope.NonGoals); err != nil || json.Unmarshal([]byte(riskResult), &envelope.RiskResult) != nil || json.Unmarshal([]byte(repositoryScope), &envelope.RepositoryScope) != nil {
		return Envelope{}, fmt.Errorf("decode envelope values")
	}
	if err := loadEnvelopeSources(conn, &envelope); err != nil {
		return Envelope{}, err
	}
	if err := conn.QueryRow(`SELECT issue_id FROM envelope_issue_selection WHERE envelope_id = ? AND revision = 1`, id).Scan(&envelope.PrimaryIssueID); err != nil {
		return Envelope{}, fmt.Errorf("load envelope issue: %w", err)
	}
	if err := loadEnvelopeGrants(conn, &envelope); err != nil {
		return Envelope{}, err
	}
	if err := loadEnvelopeObligations(conn, &envelope); err != nil {
		return Envelope{}, err
	}
	if err := loadEnvelopeProtectedDecisions(conn, &envelope); err != nil {
		return Envelope{}, err
	}
	if err := loadEnvelopeBudgets(conn, &envelope); err != nil {
		return Envelope{}, err
	}
	envelope.RelatedIssueIDs = []string{}
	return envelope, nil
}

func loadEnvelopeSources(conn *sql.DB, envelope *Envelope) error {
	rows, err := conn.Query(`SELECT source_json FROM envelope_sources WHERE envelope_id = ? AND revision = 1 ORDER BY ordinal`, envelope.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	envelope.Sources = []Source{}
	for rows.Next() {
		var raw string
		var source Source
		if err := rows.Scan(&raw); err != nil || json.Unmarshal([]byte(raw), &source) != nil {
			return fmt.Errorf("decode envelope source")
		}
		envelope.Sources = append(envelope.Sources, source)
	}
	return rows.Err()
}

func loadEnvelopeGrants(conn *sql.DB, envelope *Envelope) error {
	rows, err := conn.Query(`SELECT actor_selector_json, actions_json, scope_json FROM envelope_authority_grants WHERE envelope_id = ? AND revision = 1 ORDER BY ordinal`, envelope.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	envelope.AuthorityGrants = []AuthorityGrant{}
	for rows.Next() {
		var actor, actions, scope string
		var grant AuthorityGrant
		if err := rows.Scan(&actor, &actions, &scope); err != nil || json.Unmarshal([]byte(actor), &grant.ActorSelector) != nil || json.Unmarshal([]byte(actions), &grant.Actions) != nil || json.Unmarshal([]byte(scope), &grant.Scope) != nil {
			return fmt.Errorf("decode envelope authority grant")
		}
		envelope.AuthorityGrants = append(envelope.AuthorityGrants, grant)
	}
	return rows.Err()
}

func loadEnvelopeObligations(conn *sql.DB, envelope *Envelope) error {
	rows, err := conn.Query(`SELECT obligation_json FROM envelope_obligations WHERE envelope_id = ? AND revision = 1 ORDER BY id`, envelope.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	envelope.Obligations = []Obligation{}
	for rows.Next() {
		var raw string
		var input ObligationInput
		if err := rows.Scan(&raw); err != nil || json.Unmarshal([]byte(raw), &input) != nil {
			return fmt.Errorf("decode envelope obligation")
		}
		envelope.Obligations = append(envelope.Obligations, Obligation{ObligationInput: input, DerivedStatus: "unsatisfied", Reasons: []any{}})
	}
	return rows.Err()
}

func loadEnvelopeProtectedDecisions(conn *sql.DB, envelope *Envelope) error {
	rows, err := conn.Query(`SELECT ordinal, decision_class, required_capabilities_json, affected_actions_json FROM envelope_protected_decisions WHERE envelope_id = ? AND revision = ? ORDER BY ordinal`, envelope.ID, envelope.CurrentRevision)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	envelope.ProtectedDecisions = []ProtectedDecision{}
	for rows.Next() {
		var decision ProtectedDecision
		var capabilities, actions string
		if err := rows.Scan(&decision.Ordinal, &decision.DecisionClass, &capabilities, &actions); err != nil || json.Unmarshal([]byte(capabilities), &decision.RequiredCapabilities) != nil || json.Unmarshal([]byte(actions), &decision.AffectedActions) != nil {
			return fmt.Errorf("decode envelope protected decision")
		}
		decision.Status = "unsatisfied"
		if err := resolveProtectedDecision(conn, envelope.PrimaryIssueID, &decision); err != nil {
			return err
		}
		envelope.ProtectedDecisions = append(envelope.ProtectedDecisions, decision)
	}
	return rows.Err()
}

func protectedDecisionStopEffects(conn *sql.DB, envelope Envelope, actorID string) ([]journal.Effect, error) {
	effects := []journal.Effect{}
	for _, decision := range envelope.ProtectedDecisions {
		if decision.Status == "satisfied" {
			continue
		}
		class := decision.DecisionClass
		request := StopRaiseRequest{
			Subject:        Subject{Type: "envelope", ID: envelope.ID},
			Code:           "protected-decision-required",
			BlockedActions: decision.AffectedActions,
			ResolutionCondition: Condition{Type: "decision-satisfied", Subject: Subject{Type: "envelope", ID: envelope.ID},
				DecisionClass: &class},
			ResolverCapability: "stop.resolve",
		}
		request.Details = mustEnvelopeJSON(map[string]any{"envelope_id": envelope.ID, "decision_class": class, "required_capabilities": decision.RequiredCapabilities, "affected_actions": decision.AffectedActions})
		if err := validateStopRaise(request); err != nil {
			return nil, err
		}
		stopID, err := newStopID(conn)
		if err != nil {
			return nil, err
		}
		key, err := attentionSemanticKey(request.Subject, request.Code, request.ResolutionCondition)
		if err != nil {
			return nil, err
		}
		attention, found, err := loadOpenAttentionByKey(conn, key)
		if err != nil {
			return nil, err
		}
		if !found {
			attentionID, err := newAttentionID(conn)
			if err != nil {
				return nil, err
			}
			attention = Attention{ID: attentionID, SemanticKey: key, Subject: request.Subject, Code: request.Code, Status: "open", Version: 1, StopIDs: []string{stopID}, RaisedAt: "<occurred-at>"}
		}
		stop := Stop{ID: stopID, Subject: request.Subject, Code: request.Code, BlockedActions: request.BlockedActions, Details: request.Details, ResolutionCondition: request.ResolutionCondition, ResolverCapability: request.ResolverCapability, ActorID: actorID, RaisedAt: "<occurred-at>"}
		raised := stopRaisedEventForDomain(stop, attention, !found, "<event-id>")
		effects = append(effects, raised.Effects...)
	}
	return effects, nil
}

func setEnvelopeEventIdentity(effects []journal.Effect, event journal.Event) {
	stamp := event.OccurredAt.UTC().Format(time.RFC3339Nano)
	for index := range effects {
		for key, value := range effects[index].Values {
			switch value {
			case "<event-id>":
				effects[index].Values[key] = event.EventID
			case "<occurred-at>":
				effects[index].Values[key] = stamp
			}
		}
	}
}

func loadEnvelopeBudgets(conn *sql.DB, envelope *Envelope) error {
	rows, err := conn.Query(`SELECT budget_json FROM envelope_budgets WHERE envelope_id = ? AND revision = 1 ORDER BY name`, envelope.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	envelope.Budgets = []Budget{}
	for rows.Next() {
		var raw string
		var budget Budget
		if err := rows.Scan(&raw); err != nil || json.Unmarshal([]byte(raw), &budget) != nil {
			return fmt.Errorf("decode envelope budget")
		}
		samples, err := loadRunBudgetSamples(conn, envelope.ID, envelope.CurrentRevision, budget.Name)
		if err != nil {
			return err
		}
		derived, err := DeriveBudget(budget, samples)
		if err != nil {
			return fmt.Errorf("derive envelope budget %q: %w", budget.Name, err)
		}
		envelope.Budgets = append(envelope.Budgets, derived)
	}
	return rows.Err()
}

func envelopeRequestDigest(operation, projectID, actorID string, value any) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"actor_id": actorID, "operation": operation, "project_id": projectID, "request": value})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func envelopeResultJSON(envelope Envelope, key string) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"envelope": envelope, "idempotency_key": key, "replayed": false})
	return string(payload), err
}

func loadEnvelopeMutationReplay(outcome journal.IdempotencyOutcome, key string, result *EnvelopeMutationResult) error {
	var stored EnvelopeMutationResult
	if err := json.Unmarshal([]byte(outcome.ResultJSON), &stored); err != nil {
		return fmt.Errorf("decode envelope replay result: %w", err)
	}
	stored.IdempotencyKey, stored.Replayed, stored.EventID = key, true, outcome.EventID
	*result = stored
	return nil
}

func envelopeActivateResultJSON(envelope Envelope, activatedAt, key string) (string, error) {
	envelope.Status, envelope.ActivatedAt = "active", &activatedAt
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"envelope": envelope, "evaluation": envelopeEvaluation(envelope), "idempotency_key": key, "replayed": false})
	return string(payload), err
}

func envelopeCompleteResultJSON(envelope Envelope, completedAt string, evaluation []ObligationEvaluation, key string) (string, error) {
	envelope.Status, envelope.CompletedAt = "completed", &completedAt
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"envelope": envelope, "evaluation": evaluation, "idempotency_key": key, "replayed": false})
	return string(payload), err
}

func envelopeEvaluation(envelope Envelope) []ObligationEvaluation {
	return EvaluateEnvelopeObligations(envelope.Obligations, nil)
}

func envelopeRepository(repository Repository) *journal.Repository {
	base := ""
	if repository.BaseOID != nil {
		base = *repository.BaseOID
	}
	return &journal.Repository{BaseOID: base, CaptureSource: repository.CaptureSource, Dirty: repository.Dirty, DirtyDigest: repository.DirtyDigest, HeadOID: repository.HeadOID, ProjectID: repository.ProjectID, TreeOID: repository.TreeOID, WorktreeID: repository.WorktreeID}
}
