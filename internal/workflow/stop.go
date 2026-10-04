package workflow

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Subject identifies the sole resource a Stop or Attention can block.
type Subject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Condition describes the durable fact required to resolve a Stop.
type Condition struct {
	Type          string  `json:"type"`
	Subject       Subject `json:"subject"`
	BudgetName    *string `json:"budget_name"`
	CapabilityID  *string `json:"capability_id"`
	DecisionClass *string `json:"decision_class,omitempty"`
	OutcomeClass  *string `json:"outcome_class,omitempty"`
	FindingID     *string `json:"finding_id,omitempty"`
}

type Stop struct {
	ID                  string          `json:"id"`
	Subject             Subject         `json:"subject"`
	Code                string          `json:"code"`
	BlockedActions      []string        `json:"blocked_actions"`
	Details             json.RawMessage `json:"details"`
	ResolutionCondition Condition       `json:"resolution_condition"`
	ResolverCapability  string          `json:"resolver_capability"`
	ActorID             string          `json:"actor_id"`
	RaisedAt            string          `json:"raised_at"`
	Resolution          *StopResolution `json:"resolution"`
}

type StopResolution struct {
	ActorID     string `json:"actor_id"`
	FactEventID string `json:"fact_event_id"`
	ResolvedAt  string `json:"resolved_at"`
	EventID     string `json:"event_id"`
}

type Attention struct {
	ID                string   `json:"id"`
	SemanticKey       string   `json:"semantic_key"`
	Subject           Subject  `json:"subject"`
	Code              string   `json:"code"`
	Status            string   `json:"status"`
	Version           int      `json:"version"`
	StopIDs           []string `json:"stop_ids"`
	RaisedAt          string   `json:"raised_at"`
	ResolutionEventID *string  `json:"resolution_event_id"`
}

type StopOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
}

type StopRaiseRequest struct {
	Subject             Subject         `json:"subject"`
	Code                string          `json:"code"`
	BlockedActions      []string        `json:"blocked_actions"`
	Details             json.RawMessage `json:"details"`
	ResolutionCondition Condition       `json:"resolution_condition"`
	ResolverCapability  string          `json:"resolver_capability"`
}

type StopResolveRequest struct {
	StopID      string `json:"stop_id"`
	FactEventID string `json:"fact_event_id"`
}

type AttentionResolveRequest struct {
	AttentionID     string   `json:"attention_id"`
	ExpectedVersion int      `json:"expected_version"`
	FactEventIDs    []string `json:"fact_event_ids"`
}

type StopRaiseResult struct {
	Stop           Stop      `json:"stop"`
	Attention      Attention `json:"attention"`
	IdempotencyKey string    `json:"idempotency_key"`
	Replayed       bool      `json:"replayed"`
	EventID        string    `json:"event_id"`
}

type StopResolveResult struct {
	Stop           Stop   `json:"stop"`
	IdempotencyKey string `json:"idempotency_key"`
	Replayed       bool   `json:"replayed"`
	EventID        string `json:"event_id"`
}

type AttentionResolveResult struct {
	Attention      Attention `json:"attention"`
	IdempotencyKey string    `json:"idempotency_key"`
	Replayed       bool      `json:"replayed"`
	EventID        string    `json:"event_id"`
}

// StopBlockedError reports an unresolved Stop that applies to a guarded
// operation. It carries the immutable facts so command adapters do not infer
// scope from an error string.
type StopBlockedError struct {
	Action       string
	Stops        []Stop
	AttentionIDs []string
	Budget       *Budget
}

func (err *StopBlockedError) Error() string { return "guarded action is blocked by unresolved Stop" }

// NewStopBlockedError enriches a non-committing denial with the immutable
// Attention identities and, for review exhaustion, the derived budget state.
func NewStopBlockedError(conn *sql.DB, action string, stops []Stop) (*StopBlockedError, error) {
	blocked := &StopBlockedError{Action: action, Stops: stops, AttentionIDs: []string{}}
	for _, stop := range stops {
		rows, err := conn.Query(`SELECT attention_id FROM attention_stop_links WHERE stop_id=? ORDER BY attention_id`, stop.ID)
		if err != nil {
			return nil, fmt.Errorf("load blocking Attention: %w", err)
		}
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				_ = rows.Close()
				return nil, err
			}
			blocked.AttentionIDs = append(blocked.AttentionIDs, id)
		}
		if err := rows.Close(); err != nil {
			return nil, err
		}
		if stop.Code == "review-rework-exhausted" && blocked.Budget == nil {
			var details struct {
				EnvelopeID string `json:"envelope_id"`
				BudgetName string `json:"budget_name"`
			}
			if err := json.Unmarshal(stop.Details, &details); err != nil {
				return nil, fmt.Errorf("decode review exhaustion details: %w", err)
			}
			budget, err := loadRunBudget(conn, details.EnvelopeID, 1, details.BudgetName)
			if err != nil {
				return nil, fmt.Errorf("load review exhaustion budget: %w", err)
			}
			samples, err := loadRunBudgetSamples(conn, details.EnvelopeID, 1, details.BudgetName)
			if err != nil {
				return nil, err
			}
			derived, err := DeriveBudget(budget, samples)
			if err != nil {
				return nil, err
			}
			blocked.Budget = &derived
		}
	}
	return blocked, nil
}

var stopCodes = map[string]bool{
	"executor-liveness-unknown": true, "recovery-state-conflict": true,
	"budget-exhausted": true, "budget-data-unavailable": true,
	"review-rework-exhausted":     true,
	"protected-decision-required": true,
	"human-outcome-required":      true,
	"activation-incomplete":       true, "capability-unavailable": true, "operator-stop": true,
}

var conditionTypes = map[string]bool{
	"executor-terminated": true, "run-recovered": true, "activation-recovered": true,
	"capability-available": true, "budget-observed": true, "work-ended": true,
	"decision-satisfied": true, "human-outcome-recorded": true,
	"review-recorded": true, "finding-resolved": true,
}

var blockedActionSet = map[string]bool{
	"claim.acquire": true, "run.start": true, "run.recover": true, "gate.run": true,
	"issue.transition-terminal-success": true, "envelope.complete": true,
}

func RaiseStop(options StopOptions, request StopRaiseRequest) (StopRaiseResult, error) {
	if err := validateStopOptions(options); err != nil {
		return StopRaiseResult{}, err
	}
	if err := validateStopRaise(request); err != nil {
		return StopRaiseResult{}, err
	}
	var result StopRaiseResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid stop raise"}, func(lock *journal.WriteLock) error {
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
		if err := authorizeStop(conn, options, request.Subject, "stop.raise", ""); err != nil {
			return err
		}
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"operation": "stop.raise", "project_id": project.ProjectID, "actor_id": options.ActorID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idem := journal.IdempotencyRequest{OperationScope: "stop.raise:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idem); err != nil {
			return err
		} else if found {
			return loadStopRaiseResult(replay, &result)
		}
		stopID, err := newStopID(conn)
		if err != nil {
			return err
		}
		key, err := attentionSemanticKey(request.Subject, request.Code, request.ResolutionCondition)
		if err != nil {
			return err
		}
		attention, found, err := loadOpenAttentionByKey(conn, key)
		if err != nil {
			return err
		}
		if !found {
			attention.ID, err = newAttentionID(conn)
			if err != nil {
				return err
			}
			attention = Attention{ID: attention.ID, SemanticKey: key, Subject: request.Subject, Code: request.Code, Status: "open", Version: 1}
		}
		now := time.Now().UTC().Format(time.RFC3339Nano)
		stop := Stop{ID: stopID, Subject: request.Subject, Code: request.Code, BlockedActions: append([]string(nil), request.BlockedActions...), Details: append(json.RawMessage(nil), request.Details...), ResolutionCondition: request.ResolutionCondition, ResolverCapability: request.ResolverCapability, ActorID: options.ActorID, RaisedAt: now}
		if !found {
			attention.RaisedAt = now
		}
		attention.StopIDs = append(attention.StopIDs, stop.ID)
		sort.Strings(attention.StopIDs)
		event := stopRaisedEventForDomain(stop, attention, !found, "")
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		event = stopRaisedEventForDomain(stop, attention, !found, prepared.EventID)
		outcome, err := journal.CommitIdempotent(lock, conn, idem, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"stop": stop, "attention": attention, "idempotency_key": options.IdempotencyKey, "replayed": false})
			if err != nil {
				return journal.IdempotencyMutation{}, err
			}
			return journal.IdempotencyMutation{EntityID: stop.ID, ResultJSON: string(payload)}, nil
		})
		if err != nil {
			return err
		}
		return loadStopRaiseResult(outcome, &result)
	})
	return result, err
}

func GetStop(workspaceRoot, migrationsDir, id string) (Stop, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Stop{}, err
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		if err == sql.ErrNoRows {
			return Stop{}, &StopNotFoundError{ID: id}
		}
		return Stop{}, err
	}
	defer closeWorkflowDB(conn)
	return loadStop(conn, id)
}

// GetStopFromProjection reads a Stop from an already locked projection.
func GetStopFromProjection(conn *sql.DB, id string) (Stop, error) {
	return loadStop(conn, id)
}

// AttentionForStopFromProjection reads the Attention linked to a Stop from an
// already locked projection.
func AttentionForStopFromProjection(conn *sql.DB, stopID string) (Attention, error) {
	return loadAttentionForStop(conn, stopID)
}

func ResolveStop(options StopOptions, request StopResolveRequest) (StopResolveResult, error) {
	if err := validateStopOptions(options); err != nil || strings.TrimSpace(request.StopID) == "" || strings.TrimSpace(request.FactEventID) == "" {
		return StopResolveResult{}, fmt.Errorf("stop resolve input is invalid")
	}
	var result StopResolveResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid stop resolve"}, func(lock *journal.WriteLock) error {
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
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"operation": "stop.resolve", "project_id": project.ProjectID, "actor_id": options.ActorID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idem := journal.IdempotencyRequest{OperationScope: "stop.resolve:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idem); err != nil {
			return err
		} else if found {
			return loadStopResolveResult(replay, &result)
		}
		stop, err := loadStop(conn, request.StopID)
		if err != nil {
			return err
		}
		if stop.Resolution != nil {
			return fmt.Errorf("stop %q is already resolved", stop.ID)
		}
		if err := authorizeStop(conn, options, stop.Subject, "stop.resolve", stop.ResolverCapability); err != nil {
			return err
		}
		if err := qualifyingFact(options.WorkspaceRoot, conn, stop, request.FactEventID); err != nil {
			return &ResolutionInvalidError{Subject: stop.Subject, FactEventIDs: []string{request.FactEventID}, Reasons: []EvaluationReason{evaluationReason("resolution-invalid", err.Error(), stop.ID, "stop")}}
		}
		sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: stop.ID, Type: "stop"})
		if err != nil {
			return err
		}
		now := time.Now().UTC().Format(time.RFC3339Nano)
		event := stopResolvedEvent(stop, options.ActorID, request.FactEventID, now, "", sequence)
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		event = stopResolvedEvent(stop, options.ActorID, request.FactEventID, now, prepared.EventID, sequence)
		outcome, err := journal.CommitIdempotent(lock, conn, idem, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			stop.Resolution = &StopResolution{ActorID: options.ActorID, FactEventID: request.FactEventID, ResolvedAt: now, EventID: prepared.EventID}
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"stop": stop, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: stop.ID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		return loadStopResolveResult(outcome, &result)
	})
	return result, err
}

func ResolveAttention(options StopOptions, request AttentionResolveRequest) (AttentionResolveResult, error) {
	if err := validateStopOptions(options); err != nil || strings.TrimSpace(request.AttentionID) == "" || request.ExpectedVersion < 1 || !sortedUniqueIDs(request.FactEventIDs) {
		return AttentionResolveResult{}, fmt.Errorf("attention resolve input is invalid")
	}
	var result AttentionResolveResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid attention resolve"}, func(lock *journal.WriteLock) error {
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
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"operation": "attention.resolve", "project_id": project.ProjectID, "actor_id": options.ActorID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idem := journal.IdempotencyRequest{OperationScope: "attention.resolve:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idem); err != nil {
			return err
		} else if found {
			return loadAttentionResolveResult(replay, &result)
		}
		attention, err := loadAttentionByQuery(conn, `SELECT id, semantic_key, subject_type, subject_id, code, status, version, raised_at, resolution_event_id FROM attention_items WHERE id = ?`, request.AttentionID)
		if err != nil {
			return err
		}
		if attention.Version != request.ExpectedVersion {
			return &AttentionVersionConflictError{AttentionID: attention.ID, ExpectedVersion: request.ExpectedVersion, ActualVersion: attention.Version}
		}
		if attention.Status != "open" {
			return &ResolutionInvalidError{Subject: attention.Subject, FactEventIDs: request.FactEventIDs, Reasons: []EvaluationReason{evaluationReason("attention-not-open", attention.Status, attention.ID, "attention")}}
		}
		if err := authorizeStop(conn, options, attention.Subject, "attention.resolve", ""); err != nil {
			return err
		}
		if err := qualifyingDispositions(options.WorkspaceRoot, conn, attention, request.FactEventIDs); err != nil {
			return &ResolutionInvalidError{Subject: attention.Subject, FactEventIDs: request.FactEventIDs, Reasons: []EvaluationReason{evaluationReason("resolution-invalid", err.Error(), attention.ID, "attention")}}
		}
		event := attentionResolvedEvent(attention, options.ActorID, "")
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		event = attentionResolvedEvent(attention, options.ActorID, prepared.EventID)
		outcome, err := journal.CommitIdempotent(lock, conn, idem, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			attention.Status = "resolved"
			attention.Version++
			attention.ResolutionEventID = &prepared.EventID
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"attention": attention, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: attention.ID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		return loadAttentionResolveResult(outcome, &result)
	})
	return result, err
}

func validateStopOptions(options StopOptions) error {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.IdempotencyKey) == "" {
		return fmt.Errorf("stop options require workspace, migrations, actor, and idempotency key")
	}
	return nil
}

func validateStopRaise(request StopRaiseRequest) error {
	if !validSubject(request.Subject) || !stopCodes[request.Code] || !validBlockedActions(request.BlockedActions) || strings.TrimSpace(request.ResolverCapability) == "" {
		return fmt.Errorf("stop raise input is invalid")
	}
	if err := validateCondition(request.ResolutionCondition, request.Subject); err != nil {
		return err
	}
	var details map[string]any
	if json.Unmarshal(request.Details, &details) != nil || !validStopDetails(request.Code, details) {
		return fmt.Errorf("stop details are invalid for %q", request.Code)
	}
	return nil
}

func validSubject(subject Subject) bool {
	return (subject.Type == "run" || subject.Type == "envelope" || subject.Type == "activation") && strings.TrimSpace(subject.ID) != ""
}

func validBlockedActions(actions []string) bool {
	if len(actions) == 0 {
		return false
	}
	last := ""
	for _, action := range actions {
		if !blockedActionSet[action] || action <= last {
			return false
		}
		last = action
	}
	return true
}

func validateCondition(condition Condition, subject Subject) error {
	if !conditionTypes[condition.Type] || condition.Subject != subject {
		return fmt.Errorf("stop resolution condition is invalid")
	}
	if condition.Type == "decision-satisfied" && subject.Type == "activation" {
		// An activation has no Envelope scope, so this Condition could never
		// qualify (contracts.md section 11, DEC-250 R3).
		return fmt.Errorf("stop resolution condition is invalid")
	}
	if (condition.Type == "budget-observed") != (condition.BudgetName != nil && *condition.BudgetName != "") {
		return fmt.Errorf("stop resolution condition budget binding is invalid")
	}
	if (condition.Type == "capability-available") != (condition.CapabilityID != nil && *condition.CapabilityID != "") {
		return fmt.Errorf("stop resolution condition capability binding is invalid")
	}
	if (condition.Type == "decision-satisfied") != (condition.DecisionClass != nil && *condition.DecisionClass != "") {
		return fmt.Errorf("stop resolution condition Decision binding is invalid")
	}
	if (condition.Type == "human-outcome-recorded") != (condition.OutcomeClass != nil && *condition.OutcomeClass != "") {
		return fmt.Errorf("stop resolution condition Human Outcome binding is invalid")
	}
	if (condition.Type == "finding-resolved") != (condition.FindingID != nil && *condition.FindingID != "") {
		return fmt.Errorf("stop resolution condition Finding binding is invalid")
	}
	return nil
}

func validStopDetails(code string, details map[string]any) bool {
	keys := map[string][]string{"executor-liveness-unknown": {"run_id", "executor_id", "observation_id"}, "recovery-state-conflict": {"run_id", "observation_digest", "paths", "reasons"}, "budget-exhausted": {"envelope_id", "budget_name", "unit", "maximum", "consumed", "remaining"}, "budget-data-unavailable": {"envelope_id", "budget_name", "unit", "maximum", "consumed", "remaining"}, "review-rework-exhausted": {"envelope_id", "budget_name", "unit", "maximum", "consumed", "remaining", "last_review_id"}, "protected-decision-required": {"envelope_id", "decision_class", "required_capabilities", "affected_actions"}, "human-outcome-required": {"envelope_id", "obligation_id", "outcome_class"}, "activation-incomplete": {"activation_id", "point_of_no_return", "components"}, "capability-unavailable": {"capability_id", "available", "owner_phase"}, "operator-stop": {"reason", "reference"}}
	want, ok := keys[code]
	if !ok || len(details) != len(want) {
		return false
	}
	for _, key := range want {
		if _, ok := details[key]; !ok {
			return false
		}
	}
	return true
}

func attentionSemanticKey(subject Subject, code string, condition Condition) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"subject": subject, "code": code, "resolution_condition": condition})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func newStopID(conn *sql.DB) (string, error) {
	return journal.GenerateID("STP", func(id string) (bool, error) {
		var found int
		err := conn.QueryRow(`SELECT 1 FROM stops WHERE id = ?`, id).Scan(&found)
		return err == nil, nil
	})
}
func newAttentionID(conn *sql.DB) (string, error) {
	return journal.GenerateID("ATT", func(id string) (bool, error) {
		var found int
		err := conn.QueryRow(`SELECT 1 FROM attention_items WHERE id = ?`, id).Scan(&found)
		return err == nil, nil
	})
}

func stopRaisedEventForDomain(stop Stop, attention Attention, createAttention bool, eventID string) journal.Event {
	blocked, _ := journal.MarshalCanonicalJSON(stop.BlockedActions)
	condition, _ := journal.MarshalCanonicalJSON(stop.ResolutionCondition)
	effects := []journal.Effect{{Op: "insert", Table: "stops", PK: map[string]any{"id": stop.ID}, Values: map[string]any{"subject_type": stop.Subject.Type, "subject_id": stop.Subject.ID, "code": stop.Code, "blocked_actions_json": string(blocked), "details_json": string(stop.Details), "resolution_condition_json": string(condition), "resolver_capability": stop.ResolverCapability, "actor_id": stop.ActorID, "raised_at": stop.RaisedAt, "created_event_id": eventID}}}
	if createAttention {
		effects = append(effects, journal.Effect{Op: "insert", Table: "attention_items", PK: map[string]any{"id": attention.ID}, Values: map[string]any{"semantic_key": attention.SemanticKey, "subject_type": attention.Subject.Type, "subject_id": attention.Subject.ID, "code": attention.Code, "status": "open", "version": 1, "raised_at": attention.RaisedAt, "resolution_event_id": nil}})
	}
	effects = append(effects, journal.Effect{Op: "insert", Table: "attention_stop_links", PK: map[string]any{"attention_id": attention.ID, "stop_id": stop.ID}})
	return journal.Event{EventSchema: journal.EventSchemaV2, EventType: "stop.raised", ActorV2: journal.EventActor{ActorID: stop.ActorID, ActorType: "agent", Producer: "fid", ProvenanceRef: "workflow.stop"}, Aggregate: journal.Aggregate{ID: stop.ID, Type: "stop", Sequence: 1}, CommandName: "stop.raise", Effects: effects}
}

func stopResolvedEvent(stop Stop, actorID, factEventID, resolvedAt, eventID string, sequence int) journal.Event {
	return journal.Event{EventSchema: journal.EventSchemaV2, EventType: "stop.resolved", ActorV2: journal.EventActor{ActorID: actorID, ActorType: "agent", Producer: "fid", ProvenanceRef: "workflow.stop"}, Aggregate: journal.Aggregate{ID: stop.ID, Type: "stop", Sequence: sequence}, CommandName: "stop.resolve", Effects: []journal.Effect{{Op: "insert", Table: "stop_resolutions", PK: map[string]any{"stop_id": stop.ID}, Values: map[string]any{"actor_id": actorID, "fact_event_id": factEventID, "resolved_at": resolvedAt, "created_event_id": eventID}}}}
}

func attentionResolvedEvent(attention Attention, actorID, eventID string) journal.Event {
	return journal.Event{EventSchema: journal.EventSchemaV2, EventType: "attention.resolved", ActorV2: journal.EventActor{ActorID: actorID, ActorType: "agent", Producer: "fid", ProvenanceRef: "workflow.attention"}, Aggregate: journal.Aggregate{ID: attention.ID, Type: "attention", Sequence: 1}, CommandName: "attention.resolve", Effects: []journal.Effect{{Op: "update", Table: "attention_items", PK: map[string]any{"id": attention.ID}, Values: map[string]any{"status": "resolved", "version": attention.Version + 1, "resolution_event_id": eventID}, Preconditions: map[string]any{"status": "open", "version": attention.Version}}}}
}

func qualifyingDispositions(workspace string, conn *sql.DB, attention Attention, eventIDs []string) error {
	rows, err := conn.Query(`SELECT l.stop_id, r.fact_event_id, r.created_event_id FROM attention_stop_links l LEFT JOIN stop_resolutions r ON r.stop_id = l.stop_id WHERE l.attention_id = ? ORDER BY r.created_event_id`, attention.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	expected := make([]string, 0, len(attention.StopIDs))
	seenDispositions := make(map[string]bool, len(attention.StopIDs))
	stops := make([]Stop, 0, len(attention.StopIDs))
	facts := make([]string, 0, len(attention.StopIDs))
	for rows.Next() {
		var stopID string
		var fact, disposition sql.NullString
		if err := rows.Scan(&stopID, &fact, &disposition); err != nil {
			return err
		}
		if !fact.Valid || !disposition.Valid {
			return fmt.Errorf("attention %q has an unresolved Stop", attention.ID)
		}
		stop, err := loadStop(conn, stopID)
		if err != nil {
			return err
		}
		if !seenDispositions[disposition.String] {
			expected = append(expected, disposition.String)
			seenDispositions[disposition.String] = true
		}
		stops, facts = append(stops, stop), append(facts, fact.String)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	sort.Strings(expected)
	if len(expected) != len(eventIDs) {
		return fmt.Errorf("attention disposition Event set is not exact")
	}
	for i := range expected {
		if expected[i] != eventIDs[i] {
			return fmt.Errorf("attention disposition Event set is not exact")
		}
	}
	for i, stop := range stops {
		if err := qualifyingFact(workspace, conn, stop, facts[i]); err != nil {
			return err
		}
	}
	return nil
}

func qualifyingFact(workspace string, conn *sql.DB, stop Stop, factEventID string) error {
	found, err := journal.FindEvent(filepath.Join(workspace, config.WorkspaceDirName, "journal"), factEventID)
	if err != nil {
		return fmt.Errorf("resolve Stop %q: qualifying fact: %w", stop.ID, err)
	}
	event := found.Event
	if !eventMatchesCondition(event, stop.ResolutionCondition) {
		return fmt.Errorf("resolve Stop %q: fact Event %q does not qualify", stop.ID, factEventID)
	}
	switch stop.ResolutionCondition.Type {
	case "executor-terminated":
		return qualifyingExecutorTermination(conn, event, stop.ResolutionCondition)
	case "activation-recovered", "capability-available":
		return qualifyingActivationCondition(conn, event, stop.ResolutionCondition)
	case "run-recovered":
		return qualifyingRecoveryCondition(conn, event, stop.ResolutionCondition)
	case "budget-observed":
		return qualifyingBudgetCondition(conn, event, stop.ResolutionCondition)
	case "work-ended":
		return qualifyingWorkEnded(conn, event, stop.ResolutionCondition)
	case "review-recorded":
		return qualifyingReviewRecordedCondition(conn, event, stop.ResolutionCondition)
	case "finding-resolved":
		return qualifyingFindingResolvedCondition(conn, event, stop.ResolutionCondition)
	case "human-outcome-recorded", "decision-satisfied":
		return qualifyingJudgmentCondition(conn, event, stop.ResolutionCondition)
	}
	return fmt.Errorf("unsupported Stop resolution condition")
}

func eventMatchesCondition(event journal.Event, condition Condition) bool {
	want := map[string][]string{"executor-terminated": {"run.observed"}, "run-recovered": {"run.recovered"}, "activation-recovered": {"activation.recovery-completed", "activation.finalized"}, "capability-available": {"activation.finalized", "activation.recovery-completed"}, "budget-observed": {"run.budget-observed"}, "work-ended": {"run.closed", "run.abandoned", "envelope.completed"}, "review-recorded": {"review.recorded"}, "finding-resolved": {"review.recorded"}, "human-outcome-recorded": {"human-outcome.recorded"}, "decision-satisfied": {"decision.referenced", "human-outcome.recorded"}}
	for _, eventType := range want[condition.Type] {
		if event.EventType == eventType {
			return true
		}
	}
	return false
}

func loadOpenAttentionByKey(conn *sql.DB, key string) (Attention, bool, error) {
	attention, err := loadAttentionByQuery(conn, `SELECT id, semantic_key, subject_type, subject_id, code, status, version, raised_at, resolution_event_id FROM attention_items WHERE semantic_key = ? AND status = 'open'`, key)
	var notFound *AttentionNotFoundError
	if errors.As(err, &notFound) {
		return Attention{}, false, nil
	}
	return attention, err == nil, err
}

func loadStopRaiseResult(outcome journal.IdempotencyOutcome, result *StopRaiseResult) error {
	if err := json.Unmarshal([]byte(outcome.ResultJSON), result); err != nil {
		return fmt.Errorf("decode Stop raise snapshot: %w", err)
	}
	result.Replayed, result.EventID = outcome.Replayed, outcome.EventID
	return nil
}

func loadStopResolveResult(outcome journal.IdempotencyOutcome, result *StopResolveResult) error {
	if err := json.Unmarshal([]byte(outcome.ResultJSON), result); err != nil {
		return fmt.Errorf("decode Stop resolution snapshot: %w", err)
	}
	result.Replayed, result.EventID = outcome.Replayed, outcome.EventID
	return nil
}

func loadAttentionResolveResult(outcome journal.IdempotencyOutcome, result *AttentionResolveResult) error {
	if err := json.Unmarshal([]byte(outcome.ResultJSON), result); err != nil {
		return fmt.Errorf("decode Attention resolution snapshot: %w", err)
	}
	result.Replayed, result.EventID = outcome.Replayed, outcome.EventID
	return nil
}

func loadStop(conn *sql.DB, id string) (Stop, error) {
	var stop Stop
	var blocked, details, condition string
	err := conn.QueryRow(`SELECT id, subject_type, subject_id, code, blocked_actions_json, details_json, resolution_condition_json, resolver_capability, actor_id, raised_at FROM stops WHERE id = ?`, id).Scan(&stop.ID, &stop.Subject.Type, &stop.Subject.ID, &stop.Code, &blocked, &details, &condition, &stop.ResolverCapability, &stop.ActorID, &stop.RaisedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Stop{}, &StopNotFoundError{ID: id}
		}
		return Stop{}, err
	}
	if err := json.Unmarshal([]byte(blocked), &stop.BlockedActions); err != nil {
		return Stop{}, err
	}
	stop.Details = json.RawMessage(details)
	if err := json.Unmarshal([]byte(condition), &stop.ResolutionCondition); err != nil {
		return Stop{}, err
	}
	var resolution StopResolution
	err = conn.QueryRow(`SELECT actor_id, fact_event_id, resolved_at, created_event_id FROM stop_resolutions WHERE stop_id = ?`, id).Scan(&resolution.ActorID, &resolution.FactEventID, &resolution.ResolvedAt, &resolution.EventID)
	if err == nil {
		stop.Resolution = &resolution
	} else if err != sql.ErrNoRows {
		return Stop{}, err
	}
	return stop, nil
}

func loadAttentionForStop(conn *sql.DB, stopID string) (Attention, error) {
	return loadAttentionByQuery(conn, `SELECT a.id, a.semantic_key, a.subject_type, a.subject_id, a.code, a.status, a.version, a.raised_at, a.resolution_event_id FROM attention_items a JOIN attention_stop_links l ON l.attention_id = a.id WHERE l.stop_id = ?`, stopID)
}
func loadAttentionByQuery(conn *sql.DB, query string, argument string) (Attention, error) {
	var item Attention
	var resolution sql.NullString
	err := conn.QueryRow(query, argument).Scan(&item.ID, &item.SemanticKey, &item.Subject.Type, &item.Subject.ID, &item.Code, &item.Status, &item.Version, &item.RaisedAt, &resolution)
	if err != nil {
		if err == sql.ErrNoRows {
			return Attention{}, &AttentionNotFoundError{ID: argument}
		}
		return Attention{}, err
	}
	if resolution.Valid {
		item.ResolutionEventID = &resolution.String
	}
	rows, err := conn.Query(`SELECT stop_id FROM attention_stop_links WHERE attention_id = ? ORDER BY stop_id`, item.ID)
	if err != nil {
		return Attention{}, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return Attention{}, err
		}
		item.StopIDs = append(item.StopIDs, id)
	}
	return item, rows.Err()
}

func authorizeStop(conn *sql.DB, options StopOptions, subject Subject, action, resolver string) error {
	if subject.Type == "activation" {
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		if !hasEnvelopeCapability(actor.Capabilities, action) || (resolver != "" && !hasEnvelopeCapability(actor.Capabilities, resolver)) {
			return &EnvelopeAuthorityError{ActorID: options.ActorID, Capability: action}
		}
		return nil
	}
	var envelopeID, issueID string
	var revision int
	var err error
	if subject.Type == "run" {
		err = conn.QueryRow(`SELECT envelope_id, issue_id, envelope_revision FROM runs WHERE id = ?`, subject.ID).Scan(&envelopeID, &issueID, &revision)
	} else {
		err = conn.QueryRow(`SELECT e.id, s.issue_id, e.current_revision FROM work_envelopes e JOIN envelope_issue_selection s ON s.envelope_id = e.id AND s.revision = e.current_revision AND s.primary_issue = 1 WHERE e.id = ?`, subject.ID).Scan(&envelopeID, &issueID, &revision)
	}
	if err != nil {
		return err
	}
	if err := AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: action, Capability: action, IssueID: issueID, EnvelopeID: envelopeID, Revision: revision}); err != nil {
		return err
	}
	if resolver != "" {
		actor, err := ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
		if err != nil {
			return err
		}
		if !hasEnvelopeCapability(actor.Capabilities, resolver) {
			return &EnvelopeAuthorityError{ActorID: options.ActorID, Capability: resolver}
		}
	}
	return nil
}

func sortedUniqueIDs(ids []string) bool {
	if len(ids) == 0 {
		return false
	}
	last := ""
	for _, id := range ids {
		if id == "" || id <= last {
			return false
		}
		last = id
	}
	return true
}
