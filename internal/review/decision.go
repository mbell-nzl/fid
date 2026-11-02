package review

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type DecisionAuthority struct {
	ActorID              string   `json:"actor_id"`
	RequiredCapabilities []string `json:"required_capabilities"`
}

type DecisionScope struct {
	IssueID string `json:"issue_id"`
}

type DecisionReferenceRequest struct {
	ID        string            `json:"id"`
	Class     string            `json:"class"`
	Source    SourceInput       `json:"source"`
	Authority DecisionAuthority `json:"authority"`
	Outcome   string            `json:"outcome"`
	Scope     DecisionScope     `json:"scope"`
}

type DecisionMutationResult struct {
	Decision           map[string]any               `json:"decision"`
	ProtectedDecisions []workflow.ProtectedDecision `json:"protected_decisions"`
	Stops              []workflow.Stop              `json:"stops"`
	Attention          []workflow.Attention         `json:"attention"`
	IdempotencyKey     string                       `json:"idempotency_key"`
	Replayed           bool                         `json:"replayed"`
	EventID            string
}

type DecisionListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Class         string
	Outcome       string
	IssueID       string
	Limit         int
}

type SourceDigestMismatchError = project.SourceDigestMismatchError

func RecordDecisionReference(options JudgmentOptions, request DecisionReferenceRequest) (DecisionMutationResult, error) {
	if err := validateDecisionReference(options, request); err != nil {
		return DecisionMutationResult{}, err
	}
	var result DecisionMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid decision record-reference"}, func(lock *journal.WriteLock) error {
		conn, projectID, actor, err := openJudgmentMutation(lock, options)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"operation": "decision.record-reference", "project_id": projectID, "actor_id": options.ActorID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "decision.record-reference:" + projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return replayDecision(replay, options.IdempotencyKey, &result)
		}
		if request.Authority.ActorID != options.ActorID || !contains(actor.Capabilities, "decision.record-reference") {
			return &policy.EnvelopeAuthorityDeniedError{ActorID: options.ActorID, Action: "decision.record-reference", Capability: "decision.record-reference"}
		}
		for _, capability := range request.Authority.RequiredCapabilities {
			if !contains(actor.Capabilities, capability) {
				return &policy.EnvelopeAuthorityDeniedError{ActorID: options.ActorID, Action: "decision.record-reference", Capability: capability}
			}
		}
		if err := authorizeJudgmentSubject(options, conn, actor, "decision.record-reference", Subject{Type: "issue", ID: request.Scope.IssueID}); err != nil {
			return err
		}
		repository, err := observedRepository(options.WorkspaceRoot, projectID)
		if err != nil {
			return err
		}
		resolved, err := resolveSources(options.WorkspaceRoot, []SourceInput{request.Source}, repository, "source")
		if err != nil {
			return err
		}
		source := resolved[0]
		if exists, err := rowExists(conn, "decision_refs", request.ID); err != nil {
			return err
		} else if exists {
			return fmt.Errorf("decision reference %q already exists", request.ID)
		}
		fact := workflow.PendingDecisionReferenceFact(request.ID, request.Class, request.Outcome, request.Scope.IssueID, request.Authority.RequiredCapabilities)
		decisions, stops, attention, resolutionEffects, err := plannedJudgmentResolutions(conn, request.Scope.IssueID, fact, options.ActorID, false)
		if err != nil {
			return err
		}
		sourceJSON, _ := journal.MarshalCanonicalJSON(source)
		capabilitiesJSON, _ := journal.MarshalCanonicalJSON(request.Authority.RequiredCapabilities)
		effects := []journal.Effect{{Op: "insert", Table: "decision_refs", PK: map[string]any{"id": request.ID}, Values: map[string]any{
			"class": request.Class, "source_json": string(sourceJSON), "authority_actor_id": request.Authority.ActorID,
			"authority_capabilities_json": string(capabilitiesJSON), "outcome": request.Outcome,
			"scope_issue_id": request.Scope.IssueID, "created_at": "<occurred-at>", "created_event_id": "<event-id>",
		}}}
		effects = append(effects, resolutionEffects...)
		activity, err := reviewSubjectActivityEffect(conn, "decision", request.ID, "decision-referenced", options.ActorID, map[string]any{"decision_id": request.ID})
		if err != nil {
			return err
		}
		effects = append(effects, activity)
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "decision.referenced", ActorV2: judgmentActor(actor), Aggregate: journal.Aggregate{ID: request.ID, Sequence: 1, Type: "decision"}, CommandName: "decision.record-reference", Repository: journalRepository(repository), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			setJudgmentEventIdentity(prepared.Effects, prepared)
			value := decisionReferenceValue(request, source, prepared)
			finalizeDecisionResolutionValues(stops, attention, prepared)
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"decision": value, "protected_decisions": decisions, "stops": stops, "attention": attention, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: request.ID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		return replayDecision(outcome, options.IdempotencyKey, &result)
	})
	return result, err
}

func ListDecisionReferences(options DecisionListOptions) ([]map[string]any, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}
	if options.Limit < 1 || options.Limit > 10000 || (options.Outcome != "" && !oneOf(options.Outcome, "approved", "rejected", "deferred")) {
		return nil, fmt.Errorf("decision list input is invalid")
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	rows, err := conn.Query(`SELECT id FROM decision_refs WHERE (?='' OR class=?) AND (?='' OR outcome=?) AND (?='' OR scope_issue_id=?) ORDER BY id LIMIT ?`, options.Class, options.Class, options.Outcome, options.Outcome, options.IssueID, options.IssueID, options.Limit+1)
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
		value, err := loadDecisionReference(conn, id)
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

func validateDecisionReference(options JudgmentOptions, request DecisionReferenceRequest) error {
	if !validJudgmentOptions(options) || strings.TrimSpace(request.ID) == "" || strings.TrimSpace(request.Class) == "" || request.Authority.ActorID == "" || request.Scope.IssueID == "" || !oneOf(request.Outcome, "approved", "rejected", "deferred") || len(request.Authority.RequiredCapabilities) == 0 || len(request.Authority.RequiredCapabilities) > 16 || !sort.StringsAreSorted(request.Authority.RequiredCapabilities) || !uniqueStrings(request.Authority.RequiredCapabilities) {
		return fmt.Errorf("decision reference input is invalid")
	}
	if request.Source.SourceType != "project-file" || request.Source.Path == nil || request.Source.ExternalID != nil || request.Source.DisplayLabel == "" {
		return fmt.Errorf("decision Source is invalid")
	}
	return nil
}

// plannedJudgmentResolutions computes, for a judgment Event that has not
// committed yet, the protected-Decision statuses it leaves behind and the
// resolution Effects for every open Stop it qualifies. Resolution is keyed on
// each Stop's recorded Condition through the one shared predicate
// (workflow.QualifyingJudgmentStops, DEC-250 R1), never on the code that raised
// the Stop, so a producer's Stop and an operator-stop resolve identically.
func plannedJudgmentResolutions(conn *sql.DB, issueID string, fact workflow.JudgmentFact, actorID string, tolerateNoEnvelope bool) ([]workflow.ProtectedDecision, []workflow.Stop, []workflow.Attention, []journal.Effect, error) {
	_, decisions, err := workflow.ProtectedDecisionsForIssueWithFact(conn, issueID, &fact)
	if err != nil {
		if !tolerateNoEnvelope || !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil, nil, err
		}
		decisions = []workflow.ProtectedDecision{}
	}
	qualified, err := workflow.QualifyingJudgmentStops(conn, fact)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	stops, attention, effects, err := resolutionsForStops(conn, qualified, actorID)
	return decisions, stops, attention, effects, err
}

// resolutionsForStops builds the stop_resolutions and attention_items Effects
// that resolve the given Stops, and their Attention, in the caller's own Event.
func resolutionsForStops(conn *sql.DB, qualified []workflow.Stop, actorID string) ([]workflow.Stop, []workflow.Attention, []journal.Effect, error) {
	stops, attention, effects := []workflow.Stop{}, []workflow.Attention{}, []journal.Effect{}
	seenAttention := map[string]bool{}
	for _, stop := range qualified {
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
	return stops, attention, effects, nil
}

func finalizeDecisionResolutionValues(stops []workflow.Stop, attention []workflow.Attention, event journal.Event) {
	stamp := event.OccurredAt.UTC().Format(time.RFC3339Nano)
	for index := range stops {
		stops[index].Resolution.FactEventID = event.EventID
		stops[index].Resolution.ResolvedAt = stamp
		stops[index].Resolution.EventID = event.EventID
	}
	for index := range attention {
		attention[index].ResolutionEventID = stringPointerValue(event.EventID)
	}
}

func decisionReferenceValue(request DecisionReferenceRequest, source Source, event journal.Event) map[string]any {
	return map[string]any{"id": request.ID, "class": request.Class, "source": source, "authority": request.Authority, "outcome": request.Outcome, "scope": request.Scope, "created_at": event.OccurredAt.UTC().Format(time.RFC3339Nano), "event_id": event.EventID}
}

func loadDecisionReference(conn *sql.DB, id string) (map[string]any, error) {
	var class, sourceJSON, actorID, capabilitiesJSON, outcome, issueID, createdAt, eventID string
	if err := conn.QueryRow(`SELECT class,source_json,authority_actor_id,authority_capabilities_json,outcome,scope_issue_id,created_at,created_event_id FROM decision_refs WHERE id=?`, id).Scan(&class, &sourceJSON, &actorID, &capabilitiesJSON, &outcome, &issueID, &createdAt, &eventID); err != nil {
		return nil, err
	}
	var source any
	var capabilities any
	if json.Unmarshal([]byte(sourceJSON), &source) != nil || json.Unmarshal([]byte(capabilitiesJSON), &capabilities) != nil {
		return nil, fmt.Errorf("decode Decision Reference")
	}
	return map[string]any{"id": id, "class": class, "source": source, "authority": map[string]any{"actor_id": actorID, "required_capabilities": capabilities}, "outcome": outcome, "scope": map[string]any{"issue_id": issueID}, "created_at": createdAt, "event_id": eventID}, nil
}

func replayDecision(outcome journal.IdempotencyOutcome, key string, result *DecisionMutationResult) error {
	var decoded DecisionMutationResult
	if err := json.Unmarshal([]byte(outcome.ResultJSON), &decoded); err != nil {
		return err
	}
	decoded.IdempotencyKey, decoded.Replayed, decoded.EventID = key, outcome.Replayed, outcome.EventID
	*result = decoded
	return nil
}

func stringPointerValue(value string) *string { return &value }
