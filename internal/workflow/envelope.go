package workflow

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
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
	Source                Source   `json:"source"`
}

type Obligation struct {
	ObligationInput
	DerivedStatus string `json:"derived_status"`
	Reasons       []any  `json:"reasons"`
}

type EnvelopeCreateRequest struct {
	Outcome    string   `json:"outcome"`
	NonGoals   []string `json:"non_goals"`
	RiskResult struct {
		Source       Source `json:"source"`
		ResultDigest string `json:"result_digest"`
	} `json:"risk_result"`
	RepositoryScope Repository        `json:"repository_scope"`
	Sources         []Source          `json:"sources"`
	PrimaryIssueID  string            `json:"primary_issue_id"`
	RelatedIssueIDs []string          `json:"related_issue_ids"`
	AuthorityGrants []AuthorityGrant  `json:"authority_grants"`
	Obligations     []ObligationInput `json:"obligations"`
	Budgets         []Budget          `json:"budgets"`
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
	ID              string           `json:"id"`
	CurrentRevision int              `json:"current_revision"`
	Status          string           `json:"status"`
	Outcome         string           `json:"outcome"`
	NonGoals        []string         `json:"non_goals"`
	RiskResult      json.RawMessage  `json:"risk_result"`
	RepositoryScope Repository       `json:"repository_scope"`
	Sources         []Source         `json:"sources"`
	PrimaryIssueID  string           `json:"primary_issue_id"`
	RelatedIssueIDs []string         `json:"related_issue_ids"`
	AuthorityGrants []AuthorityGrant `json:"authority_grants"`
	Obligations     []Obligation     `json:"obligations"`
	Budgets         []Budget         `json:"budgets"`
	CreatedActorID  string           `json:"created_actor_id"`
	ApprovedActorID *string          `json:"approved_actor_id"`
	CreatedAt       string           `json:"created_at"`
	ActivatedAt     *string          `json:"activated_at"`
	CompletedAt     *string          `json:"completed_at"`
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
		requestDigest, err := envelopeRequestDigest("envelope.activate", project.ProjectID, actor.ActorID, request)
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "envelope.activated", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: envelope.ID, Sequence: 2, Type: "envelope"}, CommandName: "envelope.activate", Repository: envelopeRepository(envelope.RepositoryScope), Effects: []journal.Effect{
			{Op: "update", Table: "work_envelopes", PK: map[string]any{"id": envelope.ID}, Values: map[string]any{"status": "active"}},
			{Op: "update", Table: "envelope_revisions", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1}, Values: map[string]any{"status": "active", "activated_at": activatedAt}},
		}}
		resultJSON, err := envelopeActivateResultJSON(envelope, activatedAt, options.IdempotencyKey)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, journal.IdempotencyRequest{OperationScope: "envelope.activate:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}, event, func(journal.Event) (journal.IdempotencyMutation, error) {
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
			loaded, err := loadEnvelope(conn, outcome.EntityID)
			if err != nil {
				return err
			}
			evaluation, err := currentEnvelopeEvaluation(conn, loaded, request.Repository)
			if err != nil {
				return err
			}
			result = EnvelopeMutationResult{Envelope: loaded, Evaluation: evaluation, IdempotencyKey: options.IdempotencyKey, Replayed: true, EventID: outcome.EventID}
			return nil
		}

		envelope, err := loadEnvelope(conn, request.EnvelopeID)
		if err != nil {
			return err
		}
		reasons, evaluation, err := evaluateEnvelopeCompletion(options, conn, envelope, request)
		if err != nil {
			return err
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

func evaluateEnvelopeCompletion(options EnvelopeOptions, conn *sql.DB, envelope Envelope, request EnvelopeCompleteRequest) ([]EvaluationReason, ObligationEvaluation, error) {
	reasons := []EvaluationReason{}
	if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: "envelope.complete", Capability: "envelope.complete", IssueID: envelope.PrimaryIssueID, EnvelopeID: envelope.ID, Revision: envelope.CurrentRevision}); err != nil {
		reasons = append(reasons, evaluationReason("authority-unavailable", "envelope.complete", options.ActorID, "actor"))
	}
	if envelope.Status != "active" {
		reasons = append(reasons, evaluationReason("envelope-inactive", envelope.Status, envelope.ID, "envelope"))
	}
	evaluation, err := currentEnvelopeEvaluation(conn, envelope, request.Repository)
	if err != nil {
		return nil, ObligationEvaluation{}, err
	}
	reasons = append(reasons, evaluation.Reasons...)
	terminal, err := issueTerminalSuccess(conn, envelope.PrimaryIssueID)
	if err != nil {
		return nil, ObligationEvaluation{}, err
	}
	if !terminal {
		reasons = append(reasons, evaluationReason("issue-not-terminal-success", "completed workflow category", envelope.PrimaryIssueID, "issue"))
	}
	return sortEvaluationReasons(reasons), evaluation, nil
}

func currentEnvelopeEvaluation(conn *sql.DB, envelope Envelope, repository Repository) (ObligationEvaluation, error) {
	states, err := loadObligationEvidenceStates(conn, envelope.Obligations[0].ID, repository)
	if err != nil {
		return ObligationEvaluation{}, err
	}
	reasons := EvaluateEvidenceStates(states)
	evaluation := ObligationEvaluation{ObligationID: envelope.Obligations[0].ID, Status: "unsatisfied", Validity: "invalid", Freshness: "stale", ArtifactAvailability: "missing", Reasons: reasons}
	for _, state := range states {
		if state.Validity == "valid" {
			evaluation.Validity = "valid"
		}
		if state.Freshness == "fresh" {
			evaluation.Freshness = "fresh"
		}
		if state.ArtifactAvailability == "available" {
			evaluation.ArtifactAvailability = "available"
		}
	}
	if len(reasons) == 0 {
		evaluation.Status, evaluation.Validity, evaluation.Freshness, evaluation.ArtifactAvailability = "satisfied", "valid", "fresh", "available"
	}
	return evaluation, nil
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
	if strings.TrimSpace(request.Outcome) == "" || len(request.Outcome) > 2048 || len(request.NonGoals) > 64 || len(request.Sources) == 0 || len(request.Sources) > 64 || strings.TrimSpace(request.PrimaryIssueID) == "" || len(request.RelatedIssueIDs) != 0 || len(request.AuthorityGrants) == 0 || len(request.AuthorityGrants) > 64 || len(request.Obligations) != 1 || len(request.Budgets) > 64 {
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
	obligation := request.Obligations[0]
	if obligation.ID == "" || obligation.Type != "gate-evidence" || obligation.Applicability != "required" || obligation.Subject.Type != "issue" || obligation.Subject.ID != request.PrimaryIssueID || obligation.Registration.ID == "" || obligation.Registration.Version < 1 || strings.Join(obligation.GuardedActions, ",") != "envelope.complete,issue.transition-terminal-success" || strings.Join(obligation.AcceptedEvidenceTypes, ",") != "receipt" || strings.Join(obligation.AcceptedTrustSources, ",") != "fid-executed" || obligation.FreshnessTarget != "repository-tree" || obligation.ArtifactPolicy != "required-available" || obligation.Independence != "fid-executed" || obligation.RequiredLane != nil || !validEnvelopeSource(obligation.Source) {
		return &EnvelopeInvalidError{Reason: "gate-evidence obligation is not the exact supported value"}
	}
	previous = ""
	for _, budget := range request.Budgets {
		if budget.Name == "" || previous >= budget.Name || budget.Scope != "envelope" || (budget.Unit != "attempts" && budget.Unit != "milliseconds" && budget.Unit != "tokens") || budget.Maximum < 1 || (budget.WarningThreshold != nil && (*budget.WarningThreshold < 1 || *budget.WarningThreshold > budget.Maximum)) || budget.ExhaustionAction != "unavailable" || budget.Consumed != nil || budget.Remaining != nil {
			return &EnvelopeInvalidError{Reason: "budget values must remain immutable null counters"}
		}
		previous = budget.Name
	}
	return nil
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
	for _, source := range append(append([]Source{request.RiskResult.Source}, request.Sources...), request.Obligations[0].Source) {
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
	obligation := request.Obligations[0]
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
	obligation := Obligation{ObligationInput: request.Obligations[0], DerivedStatus: "unsatisfied", Reasons: []any{}}
	return Envelope{ID: id, CurrentRevision: 1, Status: "draft", Outcome: request.Outcome, NonGoals: request.NonGoals, RiskResult: mustEnvelopeJSON(request.RiskResult), RepositoryScope: request.RepositoryScope, Sources: request.Sources, PrimaryIssueID: request.PrimaryIssueID, RelatedIssueIDs: []string{}, AuthorityGrants: request.AuthorityGrants, Obligations: []Obligation{obligation}, Budgets: request.Budgets, CreatedActorID: actorID, ApprovedActorID: nil, CreatedAt: createdAt, ActivatedAt: nil, CompletedAt: nil}
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
	obligation := request.Obligations[0]
	effects = append(effects, journal.Effect{Op: "insert", Table: "envelope_obligations", PK: map[string]any{"id": obligation.ID}, Values: map[string]any{"envelope_id": envelope.ID, "revision": 1, "obligation_type": obligation.Type, "registration_id": obligation.Registration.ID, "registration_version": obligation.Registration.Version, "obligation_json": string(mustEnvelopeJSON(obligation))}})
	for _, budget := range request.Budgets {
		effects = append(effects, journal.Effect{Op: "insert", Table: "envelope_budgets", PK: map[string]any{"envelope_id": envelope.ID, "revision": 1, "name": budget.Name}, Values: map[string]any{"budget_json": string(mustEnvelopeJSON(budget))}})
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
		envelope.Budgets = append(envelope.Budgets, budget)
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

func envelopeActivateResultJSON(envelope Envelope, activatedAt, key string) (string, error) {
	envelope.Status, envelope.ActivatedAt = "active", &activatedAt
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"envelope": envelope, "evaluation": envelopeEvaluation(envelope), "idempotency_key": key, "replayed": false})
	return string(payload), err
}

func envelopeCompleteResultJSON(envelope Envelope, completedAt string, evaluation ObligationEvaluation, key string) (string, error) {
	envelope.Status, envelope.CompletedAt = "completed", &completedAt
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"envelope": envelope, "evaluation": evaluation, "idempotency_key": key, "replayed": false})
	return string(payload), err
}

func envelopeEvaluation(envelope Envelope) map[string]any {
	return map[string]any{"obligation_id": envelope.Obligations[0].ID, "status": "unsatisfied", "validity": "invalid", "freshness": "stale", "artifact_availability": "missing", "reasons": []any{}}
}

func envelopeRepository(repository Repository) *journal.Repository {
	base := ""
	if repository.BaseOID != nil {
		base = *repository.BaseOID
	}
	return &journal.Repository{BaseOID: base, CaptureSource: repository.CaptureSource, Dirty: repository.Dirty, DirtyDigest: repository.DirtyDigest, HeadOID: repository.HeadOID, ProjectID: repository.ProjectID, TreeOID: repository.TreeOID, WorktreeID: repository.WorktreeID}
}
