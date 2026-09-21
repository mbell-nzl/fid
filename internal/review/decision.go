package review

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
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
	Source    Source            `json:"source"`
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

type SourceDigestMismatchError struct {
	DisplayLabel   string
	Path           *string
	ExpectedDigest string
	ObservedDigest *string
}

func (err *SourceDigestMismatchError) Error() string {
	return fmt.Sprintf("Source %q digest does not match the observed project file", err.DisplayLabel)
}

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
		digest, err := judgmentDigest("decision.record-reference", projectID, options.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "decision.record-reference:" + projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return replayDecision(replay, options.IdempotencyKey, &result)
		}
		if err := verifyDecisionSource(options.WorkspaceRoot, request.Source); err != nil {
			return err
		}
		if exists, err := rowExists(conn, "decision_refs", request.ID); err != nil {
			return err
		} else if exists {
			return fmt.Errorf("decision reference %q already exists", request.ID)
		}
		envelope, decisions, err := plannedDecisionStates(conn, request)
		if err != nil {
			return err
		}
		stops, attention, resolutionEffects, err := plannedDecisionResolutions(conn, envelope, decisions, request.Class, options.ActorID)
		if err != nil {
			return err
		}
		sourceJSON, _ := journal.MarshalCanonicalJSON(request.Source)
		capabilitiesJSON, _ := journal.MarshalCanonicalJSON(request.Authority.RequiredCapabilities)
		effects := []journal.Effect{{Op: "insert", Table: "decision_refs", PK: map[string]any{"id": request.ID}, Values: map[string]any{
			"class": request.Class, "source_json": string(sourceJSON), "authority_actor_id": request.Authority.ActorID,
			"authority_capabilities_json": string(capabilitiesJSON), "outcome": request.Outcome,
			"scope_issue_id": request.Scope.IssueID, "created_at": "<occurred-at>", "created_event_id": "<event-id>",
		}}}
		effects = append(effects, resolutionEffects...)
		repository, err := observedRepository(options.WorkspaceRoot, projectID)
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "decision.referenced", ActorV2: judgmentActor(actor), Aggregate: journal.Aggregate{ID: request.ID, Sequence: 1, Type: "decision"}, CommandName: "decision.record-reference", Repository: journalRepository(repository), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			setJudgmentEventIdentity(prepared.Effects, prepared)
			value := decisionReferenceValue(request, prepared)
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
	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
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
	if request.Source.SourceType != "project-file" || request.Source.Path == nil || request.Source.ExternalID != nil || request.Source.DisplayLabel == "" || request.Source.Digest == "" {
		return fmt.Errorf("decision Source is invalid")
	}
	return nil
}

func verifyDecisionSource(root string, source Source) error {
	if source.Path == nil || filepath.IsAbs(*source.Path) || filepath.Clean(*source.Path) != *source.Path {
		return &SourceDigestMismatchError{DisplayLabel: source.DisplayLabel, Path: source.Path, ExpectedDigest: source.Digest}
	}
	path := filepath.Join(root, filepath.FromSlash(*source.Path))
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return &SourceDigestMismatchError{DisplayLabel: source.DisplayLabel, Path: source.Path, ExpectedDigest: source.Digest}
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return &SourceDigestMismatchError{DisplayLabel: source.DisplayLabel, Path: source.Path, ExpectedDigest: source.Digest}
	}
	sum := sha256.Sum256(payload)
	observed := "sha256:" + hex.EncodeToString(sum[:])
	if observed != source.Digest {
		return &SourceDigestMismatchError{DisplayLabel: source.DisplayLabel, Path: source.Path, ExpectedDigest: source.Digest, ObservedDigest: &observed}
	}
	return nil
}

func plannedDecisionStates(conn *sql.DB, request DecisionReferenceRequest) (workflow.Envelope, []workflow.ProtectedDecision, error) {
	envelope, decisions, err := workflow.ProtectedDecisionsForIssue(conn, request.Scope.IssueID)
	if err != nil {
		return workflow.Envelope{}, nil, err
	}
	for index := range decisions {
		decision := &decisions[index]
		if decision.DecisionClass != request.Class || request.Outcome != "approved" || !equalStrings(decision.RequiredCapabilities, request.Authority.RequiredCapabilities) {
			continue
		}
		decision.DecisionRefID = &request.ID
		if decision.HumanOutcomeID != nil {
			decision.Status = "satisfied"
		}
	}
	return envelope, decisions, nil
}

func plannedDecisionResolutions(conn *sql.DB, envelope workflow.Envelope, decisions []workflow.ProtectedDecision, class, actorID string) ([]workflow.Stop, []workflow.Attention, []journal.Effect, error) {
	satisfied := false
	for _, decision := range decisions {
		if decision.DecisionClass == class && decision.Status == "satisfied" {
			satisfied = true
		}
	}
	if !satisfied {
		return []workflow.Stop{}, []workflow.Attention{}, nil, nil
	}
	rows, err := conn.Query(`SELECT id FROM stops WHERE subject_type='envelope' AND subject_id=? AND code='protected-decision-required' AND id NOT IN (SELECT stop_id FROM stop_resolutions) ORDER BY id`, envelope.ID)
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
		if stop.ResolutionCondition.DecisionClass == nil || *stop.ResolutionCondition.DecisionClass != class {
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

func decisionReferenceValue(request DecisionReferenceRequest, event journal.Event) map[string]any {
	return map[string]any{"id": request.ID, "class": request.Class, "source": request.Source, "authority": request.Authority, "outcome": request.Outcome, "scope": request.Scope, "created_at": event.OccurredAt.UTC().Format(time.RFC3339Nano), "event_id": event.EventID}
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

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func stringPointerValue(value string) *string { return &value }
