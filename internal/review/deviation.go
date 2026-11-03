package review

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type DeviationAuthority struct {
	ActorID   string  `json:"actor_id"`
	Reference *string `json:"reference"`
}

type DeviationRecordRequest struct {
	ID                    string             `json:"id"`
	Subject               Subject            `json:"subject"`
	Expected              string             `json:"expected"`
	Observed              string             `json:"observed"`
	Authority             DeviationAuthority `json:"authority"`
	Disposition           string             `json:"disposition"`
	Reasoning             string             `json:"reasoning"`
	References            []SourceInput      `json:"references"`
	SupersedesDeviationID *string            `json:"supersedes_deviation_id"`
}

type DeviationMutationResult struct {
	Deviation      map[string]any `json:"deviation"`
	IdempotencyKey string         `json:"idempotency_key"`
	Replayed       bool           `json:"replayed"`
	EventID        string
}

type DeviationListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Subject       string
	Disposition   string
	Limit         int
}

type DeviationNotFoundError struct{ ID string }

func (err *DeviationNotFoundError) Error() string {
	return fmt.Sprintf("Deviation %q was not found", err.ID)
}

func RecordDeviation(options JudgmentOptions, request DeviationRecordRequest) (DeviationMutationResult, error) {
	if err := validateDeviationRecord(options, request); err != nil {
		return DeviationMutationResult{}, err
	}
	var result DeviationMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid deviation record"}, func(lock *journal.WriteLock) error {
		conn, projectID, actor, err := openJudgmentMutation(lock, options)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"operation": "deviation.record", "project_id": projectID, "actor_id": options.ActorID, "request": request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "deviation.record:" + projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return replayDeviation(replay, options.IdempotencyKey, &result)
		}
		if err := authorizeJudgmentSubject(options, conn, actor, "deviation.record", request.Subject); err != nil {
			return err
		}
		if exists, err := rowExists(conn, "deviations", request.ID); err != nil {
			return err
		} else if exists {
			return fmt.Errorf("deviation %q already exists", request.ID)
		}
		if request.SupersedesDeviationID != nil {
			if *request.SupersedesDeviationID == request.ID {
				return fmt.Errorf("deviation cannot supersede itself")
			}
			if exists, err := rowExists(conn, "deviations", *request.SupersedesDeviationID); err != nil {
				return err
			} else if !exists {
				return &DeviationNotFoundError{ID: *request.SupersedesDeviationID}
			}
		}
		repository, err := observedRepository(options.WorkspaceRoot, projectID)
		if err != nil {
			return err
		}
		resolved, err := resolveSources(options.WorkspaceRoot, request.References, repository, "references")
		if err != nil {
			return err
		}
		references, err := journal.MarshalCanonicalJSON(resolved)
		if err != nil {
			return err
		}
		effect := journal.Effect{Op: "insert", Table: "deviations", PK: map[string]any{"id": request.ID}, Values: map[string]any{"subject_type": request.Subject.Type, "subject_id": request.Subject.ID, "actor_id": options.ActorID, "expected": request.Expected, "observed": request.Observed, "authority_actor_id": request.Authority.ActorID, "authority_reference": request.Authority.Reference, "disposition": request.Disposition, "reasoning": request.Reasoning, "references_json": string(references), "supersedes_deviation_id": request.SupersedesDeviationID, "created_at": "<occurred-at>", "created_event_id": "<event-id>"}}
		effects := []journal.Effect{effect}
		activity, err := judgmentSubjectActivityEffect(conn, request.Subject, "deviation-recorded", options.ActorID, map[string]any{"deviation_id": request.ID})
		if err != nil {
			return err
		}
		effects = append(effects, activity)
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "deviation.recorded", ActorV2: judgmentActor(actor), Aggregate: journal.Aggregate{ID: request.ID, Sequence: 1, Type: "deviation"}, CommandName: "deviation.record", Repository: journalRepository(repository), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			setJudgmentEventIdentity(prepared.Effects, prepared)
			value := deviationRecordValue(request, resolved, options.ActorID, prepared)
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"deviation": value, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: request.ID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		if outcome.Replayed {
			return replayDeviation(outcome, options.IdempotencyKey, &result)
		}
		value, err := loadDeviation(conn, request.ID)
		if err != nil {
			return err
		}
		result = DeviationMutationResult{Deviation: value, IdempotencyKey: options.IdempotencyKey, EventID: outcome.EventID}
		return nil
	})
	return result, err
}

func ListDeviations(options DeviationListOptions) ([]map[string]any, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}
	if options.Limit < 1 || options.Limit > 10000 || (options.Disposition != "" && !oneOf(options.Disposition, "accepted", "rejected", "pending-decision")) {
		return nil, fmt.Errorf("deviation list input is invalid")
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	subjectType, subjectID, err := splitSubjectFilter(options.Subject)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(`SELECT id FROM deviations WHERE (?='' OR (subject_type=? AND subject_id=?)) AND (?='' OR disposition=?) ORDER BY created_at,id LIMIT ?`, options.Subject, subjectType, subjectID, options.Disposition, options.Disposition, options.Limit+1)
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
		value, err := loadDeviation(conn, id)
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

func validateDeviationRecord(options JudgmentOptions, request DeviationRecordRequest) error {
	if !validJudgmentOptions(options) || strings.TrimSpace(request.ID) == "" || !validJudgmentSubject(request.Subject) || !validText(request.Expected) || !validText(request.Observed) || strings.TrimSpace(request.Authority.ActorID) == "" || !oneOf(request.Disposition, "accepted", "rejected", "pending-decision") || !validText(request.Reasoning) || len(request.References) > 32 {
		return fmt.Errorf("deviation record input is invalid")
	}
	return nil
}

func deviationRecordValue(request DeviationRecordRequest, references []Source, actorID string, event journal.Event) map[string]any {
	return map[string]any{"id": request.ID, "subject": map[string]any{"type": request.Subject.Type, "id": request.Subject.ID}, "actor_id": actorID, "expected": request.Expected, "observed": request.Observed, "authority": map[string]any{"actor_id": request.Authority.ActorID, "reference": request.Authority.Reference}, "disposition": request.Disposition, "reasoning": request.Reasoning, "references": references, "supersedes_deviation_id": request.SupersedesDeviationID, "created_at": event.OccurredAt.UTC().Format(time.RFC3339Nano), "event_id": event.EventID}
}

func loadDeviation(conn *sql.DB, id string) (map[string]any, error) {
	var subjectType, subjectID, actorID, expected, observed, authorityActor, disposition, reasoning, referencesJSON, createdAt, eventID string
	var authorityReference, supersedes sql.NullString
	err := conn.QueryRow(`SELECT subject_type,subject_id,actor_id,expected,observed,authority_actor_id,authority_reference,disposition,reasoning,references_json,supersedes_deviation_id,created_at,created_event_id FROM deviations WHERE id=?`, id).Scan(&subjectType, &subjectID, &actorID, &expected, &observed, &authorityActor, &authorityReference, &disposition, &reasoning, &referencesJSON, &supersedes, &createdAt, &eventID)
	if err == sql.ErrNoRows {
		return nil, &DeviationNotFoundError{ID: id}
	}
	if err != nil {
		return nil, err
	}
	var references any
	if err := json.Unmarshal([]byte(referencesJSON), &references); err != nil {
		return nil, err
	}
	return map[string]any{"id": id, "subject": map[string]any{"type": subjectType, "id": subjectID}, "actor_id": actorID, "expected": expected, "observed": observed, "authority": map[string]any{"actor_id": authorityActor, "reference": nullableValue(authorityReference)}, "disposition": disposition, "reasoning": reasoning, "references": references, "supersedes_deviation_id": nullableValue(supersedes), "created_at": createdAt, "event_id": eventID}, nil
}

func replayDeviation(outcome journal.IdempotencyOutcome, key string, result *DeviationMutationResult) error {
	var decoded struct {
		Deviation map[string]any `json:"deviation"`
	}
	if err := json.Unmarshal([]byte(outcome.ResultJSON), &decoded); err != nil {
		return err
	}
	*result = DeviationMutationResult{Deviation: decoded.Deviation, IdempotencyKey: key, Replayed: true, EventID: outcome.EventID}
	return nil
}

func splitSubjectFilter(value string) (string, string, error) {
	if value == "" {
		return "", "", nil
	}
	typeName, id, found := strings.Cut(value, ":")
	if !found || !validJudgmentSubject(Subject{Type: typeName, ID: id}) {
		return "", "", fmt.Errorf("invalid judgment subject filter")
	}
	return typeName, id, nil
}

// GetDeviation returns the same record as the corresponding list item.
func GetDeviation(workspaceRoot, migrationsDir, id string) (map[string]any, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return nil, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, fmt.Errorf("open deviation projection: %w", err)
	}
	defer func() { _ = conn.Close() }()
	return loadDeviation(conn, id)
}
