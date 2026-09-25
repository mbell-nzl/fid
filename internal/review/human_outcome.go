package review

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

// HumanOutcomeRequest records authorized Human judgment (contracts.md section 10).
type HumanOutcomeRequest struct {
	ID         string        `json:"id"`
	Class      string        `json:"class"`
	Subject    Subject       `json:"subject"`
	Result     string        `json:"result"`
	Scope      DecisionScope `json:"scope"`
	Reasoning  string        `json:"reasoning"`
	References []Source      `json:"references"`
}

type HumanOutcomeMutationResult struct {
	Outcome            map[string]any               `json:"outcome"`
	ProtectedDecisions []workflow.ProtectedDecision `json:"protected_decisions"`
	Stops              []workflow.Stop              `json:"stops"`
	Attention          []workflow.Attention         `json:"attention"`
	IdempotencyKey     string                       `json:"idempotency_key"`
	Replayed           bool                         `json:"replayed"`
	EventID            string
}

type HumanOutcomeListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Class         string
	Result        string
	IssueID       string
	Limit         int
}

// AuthorityMissingError reports that active Kit Activation provenance could
// not establish authority for a Phase 23 leaf: a non-human actor type, an
// unfinalized or absent Activation, a standalone attributed actor that was
// never staged into an Activation profile, or a missing capability
// (PRD-AUTH-007/008).
type AuthorityMissingError struct {
	Operation string
	ActorID   string
}

func (err *AuthorityMissingError) Error() string {
	return fmt.Sprintf("actor %q lacks %q authority", err.ActorID, err.Operation)
}

type actorProvenance struct {
	ActivationID      string
	ActivationVersion int
	ProfileRef        string
}

// RecordHumanOutcome requires a human Actor type in a finalized Kit Activation
// profile carrying the human-outcome.record capability. It never presents the
// recorded judgment as deterministic proof and resolves only the protected
// Decisions and human-outcome-required Stops this exact Outcome qualifies.
func RecordHumanOutcome(options JudgmentOptions, request HumanOutcomeRequest) (HumanOutcomeMutationResult, error) {
	if err := validateHumanOutcomeRequest(options, request); err != nil {
		return HumanOutcomeMutationResult{}, err
	}
	var result HumanOutcomeMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid human-outcome record"}, func(lock *journal.WriteLock) error {
		conn, projectID, actor, err := openJudgmentMutation(lock, options)
		if err != nil {
			var unavailable *workflow.TrustedActorUnavailableError
			if errors.As(err, &unavailable) {
				return &AuthorityMissingError{Operation: "human-outcome.record", ActorID: options.ActorID}
			}
			return err
		}
		defer func() { _ = conn.Close() }()
		if actor.ActorType != "human" || !contains(actor.Capabilities, "human-outcome.record") {
			return &AuthorityMissingError{Operation: "human-outcome.record", ActorID: options.ActorID}
		}
		digest, err := judgmentDigest("human-outcome.record", projectID, options.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "human-outcome.record:" + projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return replayHumanOutcome(replay, options.IdempotencyKey, &result)
		}
		if exists, err := rowExists(conn, "human_outcomes", request.ID); err != nil {
			return err
		} else if exists {
			return fmt.Errorf("human outcome %q already exists", request.ID)
		}
		provenance, err := loadActorProvenance(conn, options.ActorID)
		if err != nil {
			return fmt.Errorf("load Human Outcome provenance: %w", err)
		}
		decisions, stops, attention, resolutionEffects, err := plannedHumanOutcomeResolutions(conn, request, options.ActorID)
		if err != nil {
			return err
		}
		repository, err := observedRepository(options.WorkspaceRoot, projectID)
		if err != nil {
			return err
		}
		referencesJSON, err := journal.MarshalCanonicalJSON(request.References)
		if err != nil {
			return err
		}
		effects := []journal.Effect{{Op: "insert", Table: "human_outcomes", PK: map[string]any{"id": request.ID}, Values: map[string]any{
			"actor_id": options.ActorID, "activation_id": provenance.ActivationID, "activation_version": provenance.ActivationVersion, "profile_ref": provenance.ProfileRef,
			"class": request.Class, "subject_type": request.Subject.Type, "subject_id": request.Subject.ID, "result": request.Result,
			"scope_issue_id": request.Scope.IssueID, "reasoning": request.Reasoning, "references_json": string(referencesJSON),
			"created_at": "<occurred-at>", "created_event_id": "<event-id>",
		}}}
		effects = append(effects, resolutionEffects...)
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "human-outcome.recorded", ActorV2: judgmentActor(actor), Aggregate: journal.Aggregate{ID: request.ID, Sequence: 1, Type: "human-outcome"}, CommandName: "human-outcome.record", Repository: journalRepository(repository), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			setJudgmentEventIdentity(prepared.Effects, prepared)
			value := humanOutcomeValue(request, provenance, prepared)
			finalizeDecisionResolutionValues(stops, attention, prepared)
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"outcome": value, "protected_decisions": decisions, "stops": stops, "attention": attention, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: request.ID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		return replayHumanOutcome(outcome, options.IdempotencyKey, &result)
	})
	return result, err
}

func ListHumanOutcomes(options HumanOutcomeListOptions) ([]map[string]any, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}
	if options.Limit < 1 || options.Limit > 10000 || (options.Result != "" && !oneOf(options.Result, "approved", "rejected", "changes-required", "acknowledged")) {
		return nil, fmt.Errorf("human outcome list input is invalid")
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	rows, err := conn.Query(`SELECT id FROM human_outcomes WHERE (?='' OR class=?) AND (?='' OR result=?) AND (?='' OR scope_issue_id=?) ORDER BY id LIMIT ?`, options.Class, options.Class, options.Result, options.Result, options.IssueID, options.IssueID, options.Limit+1)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := []map[string]any{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		value, err := loadHumanOutcome(conn, id)
		if err != nil {
			return nil, err
		}
		items = append(items, value)
	}
	if len(items) > options.Limit {
		return nil, &ResultLimitExceededError{Limit: options.Limit, Observed: len(items)}
	}
	return items, rows.Err()
}

func validateHumanOutcomeRequest(options JudgmentOptions, request HumanOutcomeRequest) error {
	if !validJudgmentOptions(options) || strings.TrimSpace(request.ID) == "" || strings.TrimSpace(request.Class) == "" || !validJudgmentSubject(request.Subject) || !oneOf(request.Result, "approved", "rejected", "changes-required", "acknowledged") || strings.TrimSpace(request.Scope.IssueID) == "" || !validText(request.Reasoning) || len(request.References) > 32 {
		return fmt.Errorf("human outcome input is invalid")
	}
	return nil
}

func loadActorProvenance(conn *sql.DB, actorID string) (actorProvenance, error) {
	var provenance actorProvenance
	err := conn.QueryRow(`SELECT a.id,a.version,p.provenance_ref FROM kit_activations a JOIN activation_actor_profiles p ON p.activation_id=a.id WHERE a.status='finalized' AND p.actor_id=? LIMIT 1`, actorID).Scan(&provenance.ActivationID, &provenance.ActivationVersion, &provenance.ProfileRef)
	return provenance, err
}

// plannedHumanOutcomeResolutions computes the in-memory protected-Decision
// status this Outcome contributes to, and every open Stop (protected-Decision
// and human-outcome-required) it resolves - both atomically, in the same
// Event, the same pattern task 9 used for Decision References.
func plannedHumanOutcomeResolutions(conn *sql.DB, request HumanOutcomeRequest, actorID string) ([]workflow.ProtectedDecision, []workflow.Stop, []workflow.Attention, []journal.Effect, error) {
	envelope, decisions, err := workflow.ProtectedDecisionsForIssue(conn, request.Scope.IssueID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []workflow.ProtectedDecision{}, []workflow.Stop{}, []workflow.Attention{}, nil, nil
		}
		return nil, nil, nil, nil, err
	}
	for index := range decisions {
		decision := &decisions[index]
		if decision.DecisionClass != request.Class || request.Result != "approved" {
			continue
		}
		decision.HumanOutcomeID = &request.ID
		if decision.DecisionRefID != nil {
			decision.Status = "satisfied"
		}
	}
	decisionStops, decisionAttention, decisionEffects, err := plannedDecisionResolutions(conn, envelope, decisions, request.Class, actorID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	readinessStops, readinessAttention, readinessEffects, err := plannedHumanOutcomeReadinessResolutions(conn, envelope.ID, request, actorID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	stops := append(decisionStops, readinessStops...)
	attention := append(decisionAttention, readinessAttention...)
	effects := append(decisionEffects, readinessEffects...)
	return decisions, stops, attention, effects, nil
}

// plannedHumanOutcomeReadinessResolutions resolves every open
// human-outcome-required Stop on the Envelope whose Condition names this
// exact class, when and only when the recorded result is approved
// (contracts.md section 11).
func plannedHumanOutcomeReadinessResolutions(conn *sql.DB, envelopeID string, request HumanOutcomeRequest, actorID string) ([]workflow.Stop, []workflow.Attention, []journal.Effect, error) {
	if request.Result != "approved" {
		return []workflow.Stop{}, []workflow.Attention{}, nil, nil
	}
	rows, err := conn.Query(`SELECT id FROM stops WHERE subject_type='envelope' AND subject_id=? AND code='human-outcome-required' AND id NOT IN (SELECT stop_id FROM stop_resolutions) ORDER BY id`, envelopeID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer func() { _ = rows.Close() }()
	stops, attention, effects := []workflow.Stop{}, []workflow.Attention{}, []journal.Effect{}
	seenAttention := map[string]bool{}
	for rows.Next() {
		var stopID string
		if err := rows.Scan(&stopID); err != nil {
			return nil, nil, nil, err
		}
		stop, err := workflow.GetStopFromProjection(conn, stopID)
		if err != nil {
			return nil, nil, nil, err
		}
		if stop.ResolutionCondition.OutcomeClass == nil || *stop.ResolutionCondition.OutcomeClass != request.Class {
			continue
		}
		stop.Resolution = &workflow.StopResolution{ActorID: actorID, FactEventID: "<event-id>", ResolvedAt: "<occurred-at>", EventID: "<event-id>"}
		stops = append(stops, stop)
		effects = append(effects, journal.Effect{Op: "insert", Table: "stop_resolutions", PK: map[string]any{"stop_id": stop.ID}, Values: map[string]any{"actor_id": actorID, "fact_event_id": "<event-id>", "resolved_at": "<occurred-at>", "created_event_id": "<event-id>"}})
		item, err := workflow.AttentionForStopFromProjection(conn, stop.ID)
		if err != nil {
			return nil, nil, nil, err
		}
		if !seenAttention[item.ID] {
			item.Status, item.Version, item.ResolutionEventID = "resolved", item.Version+1, stringPointerValue("<event-id>")
			attention = append(attention, item)
			effects = append(effects, journal.Effect{Op: "update", Table: "attention_items", PK: map[string]any{"id": item.ID}, Values: map[string]any{"status": "resolved", "version": item.Version, "resolution_event_id": "<event-id>"}, Preconditions: map[string]any{"status": "open", "version": item.Version - 1}})
			seenAttention[item.ID] = true
		}
	}
	return stops, attention, effects, rows.Err()
}

func humanOutcomeValue(request HumanOutcomeRequest, provenance actorProvenance, event journal.Event) map[string]any {
	return map[string]any{
		"id": request.ID, "actor_id": event.ActorV2.ActorID,
		"provenance": map[string]any{"activation_id": provenance.ActivationID, "activation_version": provenance.ActivationVersion, "profile_ref": provenance.ProfileRef},
		"class":      request.Class, "subject": map[string]any{"type": request.Subject.Type, "id": request.Subject.ID},
		"result": request.Result, "scope": map[string]any{"issue_id": request.Scope.IssueID}, "reasoning": request.Reasoning, "references": request.References,
		"created_at": event.OccurredAt.UTC().Format(time.RFC3339Nano), "event_id": event.EventID,
	}
}

func loadHumanOutcome(conn *sql.DB, id string) (map[string]any, error) {
	var actorID, activationID, profileRef, class, subjectType, subjectID, result, scopeIssueID, reasoning, refsJSON, createdAt, eventID string
	var activationVersion int
	err := conn.QueryRow(`SELECT actor_id,activation_id,activation_version,profile_ref,class,subject_type,subject_id,result,scope_issue_id,reasoning,references_json,created_at,created_event_id FROM human_outcomes WHERE id=?`, id).Scan(&actorID, &activationID, &activationVersion, &profileRef, &class, &subjectType, &subjectID, &result, &scopeIssueID, &reasoning, &refsJSON, &createdAt, &eventID)
	if err != nil {
		return nil, err
	}
	var references any
	if err := json.Unmarshal([]byte(refsJSON), &references); err != nil {
		return nil, err
	}
	return map[string]any{
		"id": id, "actor_id": actorID,
		"provenance": map[string]any{"activation_id": activationID, "activation_version": activationVersion, "profile_ref": profileRef},
		"class":      class, "subject": map[string]any{"type": subjectType, "id": subjectID}, "result": result,
		"scope": map[string]any{"issue_id": scopeIssueID}, "reasoning": reasoning, "references": references,
		"created_at": createdAt, "event_id": eventID,
	}, nil
}

func replayHumanOutcome(outcome journal.IdempotencyOutcome, key string, result *HumanOutcomeMutationResult) error {
	var decoded HumanOutcomeMutationResult
	if err := json.Unmarshal([]byte(outcome.ResultJSON), &decoded); err != nil {
		return err
	}
	decoded.IdempotencyKey, decoded.Replayed, decoded.EventID = key, outcome.Replayed, outcome.EventID
	*result = decoded
	return nil
}
