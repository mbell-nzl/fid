package review

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type ConcernRecordRequest struct {
	ID         string        `json:"id"`
	Subject    Subject       `json:"subject"`
	Scope      Scope         `json:"scope"`
	Severity   string        `json:"severity"`
	Summary    string        `json:"summary"`
	Reasoning  string        `json:"reasoning"`
	References []SourceInput `json:"references"`
}

type ConcernDisposeRequest struct {
	ConcernID           string  `json:"concern_id"`
	Action              string  `json:"action"`
	IssueID             *string `json:"issue_id"`
	SupersedesConcernID *string `json:"supersedes_concern_id"`
	SupersedesFindingID *string `json:"supersedes_finding_id"`
	Reason              string  `json:"reason"`
}

type JudgmentOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
}

type ConcernMutationResult struct {
	Concern        map[string]any `json:"concern"`
	IdempotencyKey string         `json:"idempotency_key"`
	Replayed       bool           `json:"replayed"`
	EventID        string
}

type ConcernNotFoundError struct{ ID string }

func (err *ConcernNotFoundError) Error() string {
	return fmt.Sprintf("Concern %q was not found", err.ID)
}

type ConcernAlreadyDisposedError struct{ ID, DisposedAt string }

func (err *ConcernAlreadyDisposedError) Error() string {
	return fmt.Sprintf("Concern %q is already disposed", err.ID)
}

type ConcernIssueNotFoundError struct{ ID string }

func (err *ConcernIssueNotFoundError) Error() string {
	return fmt.Sprintf("Issue %q was not found", err.ID)
}

func RecordConcern(options JudgmentOptions, request ConcernRecordRequest) (ConcernMutationResult, error) {
	if err := validateConcernRecord(options, request); err != nil {
		return ConcernMutationResult{}, err
	}
	var result ConcernMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid concern record"}, func(lock *journal.WriteLock) error {
		conn, projectID, actor, err := openJudgmentMutation(lock, options)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"operation": "concern.record", "project_id": projectID, "actor_id": options.ActorID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "concern.record:" + projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return replayConcern(replay, options.IdempotencyKey, &result)
		}
		if err := authorizeJudgmentSubject(options, conn, actor, "concern.record", request.Subject); err != nil {
			return err
		}
		if exists, err := rowExists(conn, "concerns", request.ID); err != nil {
			return err
		} else if exists {
			return fmt.Errorf("concern %q already exists", request.ID)
		}
		repository, err := observedRepository(options.WorkspaceRoot, projectID)
		if err != nil {
			return err
		}
		references, err := resolveSources(options.WorkspaceRoot, request.References, repository, "references")
		if err != nil {
			return err
		}
		event := concernRecordEvent(request, references, options.ActorID, actor, repository)
		activity, err := judgmentSubjectActivityEffect(conn, request.Subject, "concern-recorded", options.ActorID, map[string]any{"concern_id": request.ID})
		if err != nil {
			return err
		}
		event.Effects = append(event.Effects, activity)
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			concern := concernRecordValue(request, references, options.ActorID, prepared)
			setJudgmentEventIdentity(prepared.Effects, prepared)
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"concern": concern, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: request.ID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		if outcome.Replayed {
			return replayConcern(outcome, options.IdempotencyKey, &result)
		}
		value, err := loadConcern(conn, request.ID)
		if err != nil {
			return err
		}
		result = ConcernMutationResult{Concern: value, IdempotencyKey: options.IdempotencyKey, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func DisposeConcern(options JudgmentOptions, request ConcernDisposeRequest) (ConcernMutationResult, error) {
	if err := validateConcernDispose(options, request); err != nil {
		return ConcernMutationResult{}, err
	}
	var result ConcernMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid concern dispose"}, func(lock *journal.WriteLock) error {
		conn, projectID, actor, err := openJudgmentMutation(lock, options)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"operation": "concern.dispose", "project_id": projectID, "actor_id": options.ActorID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "concern.dispose:" + projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return replayConcern(replay, options.IdempotencyKey, &result)
		}
		concern, err := loadConcern(conn, request.ConcernID)
		if err != nil {
			return err
		}
		if concern["status"] == "disposed" {
			disposition := concern["disposition"].(map[string]any)
			return &ConcernAlreadyDisposedError{ID: request.ConcernID, DisposedAt: disposition["disposed_at"].(string)}
		}
		if err := validateConcernDispositionTargets(conn, request); err != nil {
			return err
		}
		if err := authorizeJudgmentSubject(options, conn, actor, "concern.dispose", subjectFromValue(concern)); err != nil {
			return err
		}
		repository, err := observedRepository(options.WorkspaceRoot, projectID)
		if err != nil {
			return err
		}
		next, err := journal.NextAggregateSequence(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.JournalDirName), journal.Aggregate{ID: request.ConcernID, Type: "concern"})
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "concern.disposed", ActorV2: judgmentActor(actor), Aggregate: journal.Aggregate{ID: request.ConcernID, Sequence: next, Type: "concern"}, CommandName: "concern.dispose", Repository: journalRepository(repository), Effects: []journal.Effect{concernDisposeEffect(request)}}
		subject := subjectFromValue(concern)
		activity, err := judgmentSubjectActivityEffect(conn, subject, "concern-disposed", options.ActorID, map[string]any{"concern_id": request.ConcernID})
		if err != nil {
			return err
		}
		event.Effects = append(event.Effects, activity)
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			setJudgmentEventIdentity(prepared.Effects, prepared)
			value := disposedConcernValue(concern, request, options.ActorID, prepared)
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"concern": value, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: request.ConcernID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		if outcome.Replayed {
			return replayConcern(outcome, options.IdempotencyKey, &result)
		}
		value, err := loadConcern(conn, request.ConcernID)
		if err != nil {
			return err
		}
		result = ConcernMutationResult{Concern: value, IdempotencyKey: options.IdempotencyKey, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func validateConcernRecord(options JudgmentOptions, request ConcernRecordRequest) error {
	if !validJudgmentOptions(options) || strings.TrimSpace(request.ID) == "" || !validJudgmentSubject(request.Subject) || !validScope(request.Scope) || !oneOf(request.Severity, "blocker", "major", "minor", "nit") || !validText(request.Summary) || !validText(request.Reasoning) || len(request.References) > 32 {
		return fmt.Errorf("concern record input is invalid")
	}
	return nil
}

func validateConcernDispose(options JudgmentOptions, request ConcernDisposeRequest) error {
	if !validJudgmentOptions(options) || strings.TrimSpace(request.ConcernID) == "" || !validText(request.Reason) {
		return fmt.Errorf("concern dispose input is invalid")
	}
	return nil
}

func validateConcernDispositionTargets(conn *sql.DB, request ConcernDisposeRequest) error {
	switch request.Action {
	case "carried-forward":
		if request.IssueID == nil || request.SupersedesConcernID != nil || request.SupersedesFindingID != nil {
			return fmt.Errorf("carried-forward disposition requires one Issue")
		}
		if ok, err := rowExists(conn, "issues", *request.IssueID); err != nil {
			return err
		} else if !ok {
			return &ConcernIssueNotFoundError{ID: *request.IssueID}
		}
	case "rejected":
		if request.IssueID != nil || request.SupersedesConcernID != nil || request.SupersedesFindingID != nil {
			return fmt.Errorf("rejected disposition cannot bind a target")
		}
	case "superseded":
		if request.IssueID != nil || (request.SupersedesConcernID == nil) == (request.SupersedesFindingID == nil) {
			return fmt.Errorf("superseded disposition requires one Concern or Finding")
		}
		if request.SupersedesConcernID != nil {
			if *request.SupersedesConcernID == request.ConcernID {
				return fmt.Errorf("concern cannot supersede itself")
			}
			if ok, err := rowExists(conn, "concerns", *request.SupersedesConcernID); err != nil {
				return err
			} else if !ok {
				return &ConcernNotFoundError{ID: *request.SupersedesConcernID}
			}
		} else if ok, err := rowExists(conn, "review_findings", *request.SupersedesFindingID); err != nil {
			return err
		} else if !ok {
			return &FindingNotFoundError{ID: *request.SupersedesFindingID}
		}
	default:
		return fmt.Errorf("unknown Concern disposition action")
	}
	return nil
}

func concernRecordEvent(request ConcernRecordRequest, sources []Source, actorID string, actor workflow.ActorProfile, repository Repository) journal.Event {
	scope, _ := journal.MarshalCanonicalJSON(request.Scope)
	references, _ := journal.MarshalCanonicalJSON(sources)
	effect := journal.Effect{Op: "insert", Table: "concerns", PK: map[string]any{"id": request.ID}, Values: map[string]any{"subject_type": request.Subject.Type, "subject_id": request.Subject.ID, "actor_id": actorID, "scope_json": string(scope), "severity": request.Severity, "summary": request.Summary, "reasoning": request.Reasoning, "references_json": string(references), "status": "open", "disposition_action": nil, "disposition_issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "disposition_reason": nil, "disposed_actor_id": nil, "disposed_at": nil, "disposed_event_id": nil, "created_at": "<occurred-at>", "created_event_id": "<event-id>"}}
	return journal.Event{EventSchema: journal.EventSchemaV2, EventType: "concern.recorded", ActorV2: judgmentActor(actor), Aggregate: journal.Aggregate{ID: request.ID, Sequence: 1, Type: "concern"}, CommandName: "concern.record", Repository: journalRepository(repository), Effects: []journal.Effect{effect}}
}

func judgmentActivityIssueID(conn *sql.DB, subject Subject) (string, error) {
	if subject.Type == "issue" {
		return subject.ID, nil
	}
	binding, found, err := judgmentEnvelopeBinding(conn, subject)
	if err != nil {
		return "", fmt.Errorf("resolve activity Issue: %w", err)
	}
	if !found {
		return "", fmt.Errorf("resolve activity Issue for %s %s: missing subject", subject.Type, subject.ID)
	}
	return binding.IssueID, nil
}

func reviewActivityEffect(conn *sql.DB, issueID, name, actorID string, details map[string]any) (journal.Effect, error) {
	return reviewSubjectActivityEffect(conn, "issue", issueID, name, actorID, details)
}

func judgmentSubjectActivityEffect(conn *sql.DB, subject Subject, name, actorID string, details map[string]any) (journal.Effect, error) {
	if subject.Type == "requirement" {
		return reviewSubjectActivityEffect(conn, "requirement", subject.ID, name, actorID, details)
	}
	issueID, err := judgmentActivityIssueID(conn, subject)
	if err != nil {
		return journal.Effect{}, err
	}
	return reviewActivityEffect(conn, issueID, name, actorID, details)
}

func reviewSubjectActivityEffect(conn *sql.DB, subjectType, subjectID, name, actorID string, details map[string]any) (journal.Effect, error) {
	id, err := journal.GenerateID("ACT", func(candidate string) (bool, error) { return rowExists(conn, "activity_log", candidate) })
	if err != nil {
		return journal.Effect{}, err
	}
	raw, err := journal.MarshalCanonicalJSON(details)
	if err != nil {
		return journal.Effect{}, err
	}
	return journal.Effect{Op: "activity", Table: "activity_log", PK: map[string]any{"id": id}, Values: map[string]any{
		"journal_event_id": "<event-id>", "ref_type": subjectType, "ref_id": subjectID,
		"event": name, "details": string(raw), "actor": actorID, "occurred_at": "<occurred-at>",
	}}, nil
}

func concernDisposeEffect(request ConcernDisposeRequest) journal.Effect {
	return journal.Effect{Op: "update", Table: "concerns", PK: map[string]any{"id": request.ConcernID}, Preconditions: map[string]any{"status": "open"}, Values: map[string]any{"status": "disposed", "disposition_action": request.Action, "disposition_issue_id": request.IssueID, "supersedes_concern_id": request.SupersedesConcernID, "supersedes_finding_id": request.SupersedesFindingID, "disposition_reason": request.Reason, "disposed_actor_id": "<actor-id>", "disposed_at": "<occurred-at>", "disposed_event_id": "<event-id>"}}
}

func concernRecordValue(request ConcernRecordRequest, references []Source, actorID string, event journal.Event) map[string]any {
	return map[string]any{"id": request.ID, "subject": map[string]any{"type": request.Subject.Type, "id": request.Subject.ID}, "actor_id": actorID, "scope": map[string]any{"paths": request.Scope.Paths}, "severity": request.Severity, "summary": request.Summary, "reasoning": request.Reasoning, "references": references, "status": "open", "disposition": nil, "created_at": event.OccurredAt.UTC().Format(time.RFC3339Nano), "event_id": event.EventID}
}

func disposedConcernValue(value map[string]any, request ConcernDisposeRequest, actorID string, event journal.Event) map[string]any {
	copyValue := map[string]any{}
	for key, item := range value {
		copyValue[key] = item
	}
	copyValue["status"] = "disposed"
	copyValue["disposition"] = map[string]any{"action": request.Action, "issue_id": request.IssueID, "supersedes_concern_id": request.SupersedesConcernID, "supersedes_finding_id": request.SupersedesFindingID, "reason": request.Reason, "actor_id": actorID, "disposed_at": event.OccurredAt.UTC().Format(time.RFC3339Nano), "event_id": event.EventID}
	return copyValue
}

func loadConcern(conn *sql.DB, id string) (map[string]any, error) {
	var subjectType, subjectID, actorID, scopeJSON, severity, summary, reasoning, refsJSON, status, createdAt, eventID string
	var action, issueID, concernID, findingID, reason, disposedActor, disposedAt, disposedEvent sql.NullString
	err := conn.QueryRow(`SELECT subject_type,subject_id,actor_id,scope_json,severity,summary,reasoning,references_json,status,disposition_action,disposition_issue_id,supersedes_concern_id,supersedes_finding_id,disposition_reason,disposed_actor_id,disposed_at,disposed_event_id,created_at,created_event_id FROM concerns WHERE id=?`, id).Scan(&subjectType, &subjectID, &actorID, &scopeJSON, &severity, &summary, &reasoning, &refsJSON, &status, &action, &issueID, &concernID, &findingID, &reason, &disposedActor, &disposedAt, &disposedEvent, &createdAt, &eventID)
	if err == sql.ErrNoRows {
		return nil, &ConcernNotFoundError{ID: id}
	}
	if err != nil {
		return nil, err
	}
	var scope, references any
	if err := json.Unmarshal([]byte(scopeJSON), &scope); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(refsJSON), &references); err != nil {
		return nil, err
	}
	var disposition any
	if action.Valid {
		disposition = map[string]any{"action": action.String, "issue_id": nullableValue(issueID), "supersedes_concern_id": nullableValue(concernID), "supersedes_finding_id": nullableValue(findingID), "reason": reason.String, "actor_id": disposedActor.String, "disposed_at": disposedAt.String, "event_id": disposedEvent.String}
	}
	return map[string]any{"id": id, "subject": map[string]any{"type": subjectType, "id": subjectID}, "actor_id": actorID, "scope": scope, "severity": severity, "summary": summary, "reasoning": reasoning, "references": references, "status": status, "disposition": disposition, "created_at": createdAt, "event_id": eventID}, nil
}

func replayConcern(outcome journal.IdempotencyOutcome, key string, result *ConcernMutationResult) error {
	var decoded struct {
		Concern map[string]any `json:"concern"`
	}
	if err := json.Unmarshal([]byte(outcome.ResultJSON), &decoded); err != nil {
		return err
	}
	*result = ConcernMutationResult{Concern: decoded.Concern, IdempotencyKey: key, Replayed: true, EventID: outcome.EventID}
	return nil
}
