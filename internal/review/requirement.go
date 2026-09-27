package review

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type RequirementReferenceRequest struct {
	ID     string `json:"id"`
	Source Source `json:"source"`
	Status string `json:"status"`
}

type RequirementLinkRequest struct {
	RequirementID string `json:"requirement_id"`
	IssueID       string `json:"issue_id"`
	Relation      string `json:"relation"`
}

type RequirementMutationResult struct {
	Requirement    map[string]any `json:"requirement,omitempty"`
	Link           map[string]any `json:"link,omitempty"`
	IdempotencyKey string         `json:"idempotency_key"`
	Replayed       bool           `json:"replayed"`
	EventID        string         `json:"-"`
}

type RequirementNotFoundError struct{ ID string }

func (err *RequirementNotFoundError) Error() string {
	return fmt.Sprintf("requirement %q was not found", err.ID)
}

type RequirementListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Status        string
	IssueID       string
	Limit         int
}

func RecordRequirementReference(options JudgmentOptions, request RequirementReferenceRequest) (RequirementMutationResult, error) {
	if !validJudgmentOptions(options) || strings.TrimSpace(request.ID) == "" || !oneOf(request.Status, "active", "withdrawn") || request.Source.SourceType != "project-file" || request.Source.Path == nil || request.Source.ExternalID != nil || request.Source.DisplayLabel == "" || request.Source.Digest == "" {
		return RequirementMutationResult{}, fmt.Errorf("requirement reference input is invalid")
	}
	var result RequirementMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid requirement record-reference"}, func(lock *journal.WriteLock) error {
		conn, projectID, actor, err := openJudgmentMutation(lock, options)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		digest, err := judgmentDigest("requirement.record-reference", projectID, options.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "requirement.record-reference:" + projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return decodeRequirementResult(replay, &result)
		}
		if err := verifyDecisionSource(options.WorkspaceRoot, request.Source); err != nil {
			return err
		}
		refID, err := journal.GenerateID("RREF", func(candidate string) (bool, error) { return rowExists(conn, "requirement_refs", candidate) })
		if err != nil {
			return err
		}
		effects, err := requirementReferenceEffects(conn, request.ID, refID, request.Status)
		if err != nil {
			return err
		}
		sourceJSON, err := journal.MarshalCanonicalJSON(request.Source)
		if err != nil {
			return err
		}
		effects[len(effects)-1].Values["source_json"] = string(sourceJSON)
		effects[len(effects)-1].Values["superseded_by_id"] = nil
		effects[len(effects)-1].Values["created_at"] = "<occurred-at>"
		effects[len(effects)-1].Values["created_event_id"] = "<event-id>"
		repository, err := observedRepository(options.WorkspaceRoot, projectID)
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "requirement.referenced", ActorV2: judgmentActor(actor), Aggregate: journal.Aggregate{ID: refID, Sequence: 1, Type: "requirement"}, CommandName: "requirement.record-reference", Repository: journalRepository(repository), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			setJudgmentEventIdentity(prepared.Effects, prepared)
			value := map[string]any{"id": refID, "source": request.Source, "status": request.Status, "superseded_by_id": nil, "created_at": prepared.OccurredAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"), "event_id": prepared.EventID}
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"requirement": value, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: refID, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		return decodeRequirementResult(outcome, &result)
	})
	return result, err
}

func decodeRequirementResult(outcome journal.IdempotencyOutcome, result *RequirementMutationResult) error {
	if err := json.Unmarshal([]byte(outcome.ResultJSON), result); err != nil {
		return err
	}
	result.Replayed, result.EventID = outcome.Replayed, outcome.EventID
	return nil
}

func LinkRequirement(options JudgmentOptions, request RequirementLinkRequest) (RequirementMutationResult, error) {
	if !validJudgmentOptions(options) || request.RequirementID == "" || request.IssueID == "" || !oneOf(request.Relation, "implements", "verifies", "relates-to") {
		return RequirementMutationResult{}, fmt.Errorf("requirement link input is invalid")
	}
	var result RequirementMutationResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid requirement link"}, func(lock *journal.WriteLock) error {
		conn, projectID, actor, err := openJudgmentMutation(lock, options)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		digest, err := judgmentDigest("requirement.link", projectID, options.ActorID, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "requirement.link:" + projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if replay, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			return decodeRequirementResult(replay, &result)
		}
		if err := validateRequirementLinkTargets(conn, request); err != nil {
			return err
		}
		repository, err := observedRepository(options.WorkspaceRoot, projectID)
		if err != nil {
			return err
		}
		key := request.RequirementID + ":" + request.IssueID + ":" + request.Relation
		effects := []journal.Effect{{Op: "insert", Table: "requirement_issue_links", PK: map[string]any{"requirement_id": request.RequirementID, "issue_id": request.IssueID, "relation": request.Relation}, Values: map[string]any{"created_at": "<occurred-at>", "created_event_id": "<event-id>"}}}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "requirement.linked", ActorV2: judgmentActor(actor), Aggregate: journal.Aggregate{ID: key, Sequence: 1, Type: "requirement"}, CommandName: "requirement.link", Repository: journalRepository(repository), Effects: effects}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			setJudgmentEventIdentity(prepared.Effects, prepared)
			value := map[string]any{"requirement_id": request.RequirementID, "issue_id": request.IssueID, "relation": request.Relation, "created_at": prepared.OccurredAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"), "event_id": prepared.EventID}
			payload, err := journal.MarshalCanonicalJSON(map[string]any{"link": value, "idempotency_key": options.IdempotencyKey, "replayed": false})
			return journal.IdempotencyMutation{EntityID: key, ResultJSON: string(payload)}, err
		})
		if err != nil {
			return err
		}
		return decodeRequirementResult(outcome, &result)
	})
	return result, err
}

func validateRequirementLinkTargets(conn *sql.DB, request RequirementLinkRequest) error {
	var marker int
	if err := conn.QueryRow(`SELECT 1 FROM requirement_refs WHERE requirement_id=? LIMIT 1`, request.RequirementID).Scan(&marker); err != nil {
		if err == sql.ErrNoRows {
			return &RequirementNotFoundError{ID: request.RequirementID}
		}
		return err
	}
	if err := conn.QueryRow(`SELECT 1 FROM issues WHERE id=?`, request.IssueID).Scan(&marker); err != nil {
		return err
	}
	if err := conn.QueryRow(`SELECT 1 FROM requirement_issue_links WHERE requirement_id=? AND issue_id=? AND relation=?`, request.RequirementID, request.IssueID, request.Relation).Scan(&marker); err == nil {
		return fmt.Errorf("requirement link already exists")
	} else if err != sql.ErrNoRows {
		return err
	}
	return nil
}

func GetRequirement(options RequirementListOptions, id string) (map[string]any, []map[string]any, error) {
	conn, err := openRequirementRead(options)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = conn.Close() }()
	var refID, requirementID string
	err = conn.QueryRow(`SELECT id, requirement_id FROM requirement_refs WHERE id=? OR requirement_id=? ORDER BY CASE WHEN id=? THEN 0 ELSE 1 END, CASE WHEN status='active' THEN 0 ELSE 1 END, created_at DESC, id DESC LIMIT 1`, id, id, id).Scan(&refID, &requirementID)
	if err == sql.ErrNoRows {
		return nil, nil, &RequirementNotFoundError{ID: id}
	}
	if err != nil {
		return nil, nil, err
	}
	value, err := loadRequirementRef(conn, refID)
	if err != nil {
		return nil, nil, err
	}
	links, err := loadRequirementLinks(conn, requirementID)
	return value, links, err
}

func ListRequirements(options RequirementListOptions) ([]map[string]any, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}
	if options.Limit < 1 || options.Limit > 10000 || (options.Status != "" && !oneOf(options.Status, "active", "superseded", "withdrawn")) {
		return nil, fmt.Errorf("requirement list input is invalid")
	}
	conn, err := openRequirementRead(options)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	rows, err := conn.Query(`SELECT id FROM requirement_refs r WHERE (?='' OR status=?) AND (?='' OR EXISTS (SELECT 1 FROM requirement_issue_links l WHERE l.requirement_id=r.requirement_id AND l.issue_id=?)) ORDER BY created_at,id LIMIT ?`, options.Status, options.Status, options.IssueID, options.IssueID, options.Limit+1)
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
		item, err := loadRequirementRef(conn, id)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if len(items) > options.Limit {
		return nil, fmt.Errorf("requirement list exceeds limit %d", options.Limit)
	}
	return items, rows.Err()
}

func openRequirementRead(options RequirementListOptions) (*sql.DB, error) {
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	return fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
}

func loadRequirementRef(conn *sql.DB, id string) (map[string]any, error) {
	var sourceJSON, status, createdAt, eventID string
	var supersededBy sql.NullString
	if err := conn.QueryRow(`SELECT source_json,status,superseded_by_id,created_at,created_event_id FROM requirement_refs WHERE id=?`, id).Scan(&sourceJSON, &status, &supersededBy, &createdAt, &eventID); err != nil {
		return nil, err
	}
	var source any
	if err := json.Unmarshal([]byte(sourceJSON), &source); err != nil {
		return nil, err
	}
	return map[string]any{"id": id, "source": source, "status": status, "superseded_by_id": nullableValue(supersededBy), "created_at": createdAt, "event_id": eventID}, nil
}

func loadRequirementLinks(conn *sql.DB, requirementID string) ([]map[string]any, error) {
	rows, err := conn.Query(`SELECT requirement_id,issue_id,relation,created_at,created_event_id FROM requirement_issue_links WHERE requirement_id=? ORDER BY issue_id,relation`, requirementID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := []map[string]any{}
	for rows.Next() {
		var requirementID, issueID, relation, createdAt, eventID string
		if err := rows.Scan(&requirementID, &issueID, &relation, &createdAt, &eventID); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{"requirement_id": requirementID, "issue_id": issueID, "relation": relation, "created_at": createdAt, "event_id": eventID})
	}
	return items, rows.Err()
}

func requirementReferenceEffects(conn *sql.DB, requirementID, refID, status string) ([]journal.Effect, error) {
	var oldID string
	err := conn.QueryRow(`SELECT id FROM requirement_refs WHERE requirement_id=? AND status='active'`, requirementID).Scan(&oldID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("find active requirement reference: %w", err)
	}
	effects := make([]journal.Effect, 0, 2)
	if err == nil {
		effects = append(effects, journal.Effect{Op: "update", Table: "requirement_refs", PK: map[string]any{"id": oldID}, Preconditions: map[string]any{"status": "active", "superseded_by_id": nil}, Values: map[string]any{"status": "superseded", "superseded_by_id": refID}})
	}
	effects = append(effects, journal.Effect{Op: "insert", Table: "requirement_refs", PK: map[string]any{"id": refID}, Values: map[string]any{"requirement_id": requirementID, "status": status}})
	return effects, nil
}
