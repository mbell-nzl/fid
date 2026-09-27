package links

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func Add(options AddOptions) (AddResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return AddResult{}, fmt.Errorf("link add: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return AddResult{}, fmt.Errorf("link add: migrations dir is required")
	}
	if strings.TrimSpace(options.FromID) == "" || strings.TrimSpace(options.ToID) == "" {
		return AddResult{}, fmt.Errorf("link add: from and to ids are required")
	}
	actor, err := standaloneActorForOptions(options.Actor)
	if err != nil {
		return AddResult{}, err
	}

	fromType, err := endpointType(options.FromType, "from-type", "link add")
	if err != nil {
		return AddResult{}, err
	}
	toType, err := endpointType(options.ToType, "to-type", "link add")
	if err != nil {
		return AddResult{}, err
	}
	declaredRelation := strings.TrimSpace(options.Relation)
	if _, ok := inverseRelation(declaredRelation); !ok {
		return AddResult{}, &UnknownRelationError{Operation: "link add", Relation: declaredRelation}
	}

	fromID := strings.TrimSpace(options.FromID)
	toID := strings.TrimSpace(options.ToID)
	if fromType == "issue" && toType == "issue" && fromID == toID && isDependencyRelation(declaredRelation) {
		return AddResult{}, &DependencyCycleError{Operation: "link add", Path: []string{fromID, toID}}
	}
	if fromType == toType && fromID == toID {
		return AddResult{}, &SelfLinkError{Operation: "link add", FromID: fromID, Relation: declaredRelation, ToID: toID}
	}

	var result AddResult
	commandArgv := linkAddArgv(fromID, fromType, declaredRelation, toID, toType, options.Actor, options.IdempotencyKey)
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command: strings.Join(commandArgv, " "),
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeLinksConn(conn)
		projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
		if err != nil {
			return err
		}
		digest, err := linkRequestDigest("link.add", projectConfig.ProjectID, actor, map[string]any{"from_id": fromID, "from_type": fromType, "relation": declaredRelation, "to_id": toID, "to_type": toType})
		if err != nil {
			return err
		}
		idempotencyRequest := journal.IdempotencyRequest{OperationScope: "link.add:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
			return err
		} else if found {
			decoded, err := addResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
			if err != nil {
				return err
			}
			result = decoded
			return nil
		}

		declared, inverse, err := plannedLinks(conn, fromType, fromID, declaredRelation, toType, toID)
		if err != nil {
			return err
		}

		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return linkActivityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		observedRevisions, err := linkIssueRevisions(conn, declared)
		if err != nil {
			return err
		}
		revisionEffects, err := issues.BuildIssueRevisionEffects(observedRevisions)
		if err != nil {
			return err
		}
		sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: declared.FromID + ":" + declared.Relation + ":" + declared.ToID, Type: "link"})
		if err != nil {
			return err
		}
		baseEffects := append([]journal.Effect{linkInsertEffect(declared, ""), linkInsertEffect(inverse, "")}, revisionEffects...)
		event, err := journal.PrepareEvent(lock, journal.Event{EventSchema: journal.EventSchemaV2, EventType: "link.added", ActorV2: actor.eventActor(), Aggregate: journal.Aggregate{ID: declared.FromID + ":" + declared.Relation + ":" + declared.ToID, Sequence: sequence, Type: "link"}, CommandName: "link.add", Effects: baseEffects})
		if err != nil {
			return err
		}
		occurredAt := event.OccurredAt.UTC().Format(time.RFC3339)
		createdAt := event.OccurredAt.UTC().Format("2006-01-02T15:04:05Z07:00")
		event.Effects = append([]journal.Effect{linkInsertEffect(declared, createdAt), linkInsertEffect(inverse, createdAt)}, revisionEffects...)
		event.Effects = append(event.Effects, linkAddActivityEffect(activityID, event.EventID, actor.id, occurredAt, declared, inverse))
		seed := AddResult{Declared: declared, Inverse: inverse, IssueRevisions: resultIssueRevisions(observedRevisions), OccurredAt: event.OccurredAt.UTC().Format(time.RFC3339Nano), IdempotencyKey: options.IdempotencyKey}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			encoded, err := addResultJSON(seed)
			if err != nil {
				return journal.IdempotencyMutation{}, err
			}
			return journal.IdempotencyMutation{EntityID: declared.FromID + ":" + declared.Relation + ":" + declared.ToID, ResultJSON: encoded}, nil
		})
		if err != nil {
			return err
		}
		decoded, err := addResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		if err != nil {
			return err
		}
		result = decoded
		return nil
	})
	if err != nil {
		return AddResult{}, err
	}

	return result, nil
}

func plannedLinks(conn *sql.DB, fromType, fromID, relation, toType, toID string) (Link, Link, error) {
	resolvedFromType, err := resolveEndpointType(conn, fromType, fromID)
	if err != nil {
		return Link{}, Link{}, err
	}
	resolvedToType, err := resolveEndpointType(conn, toType, toID)
	if err != nil {
		return Link{}, Link{}, err
	}
	declared := Link{
		FromType: resolvedFromType,
		FromID:   fromID,
		Relation: relation,
		ToType:   resolvedToType,
		ToID:     toID,
	}
	if err := assertEndpointExists(conn, "link add", declared.FromType, declared.FromID, relation, declared.ToType, declared.ToID); err != nil {
		return Link{}, Link{}, err
	}
	if err := assertEndpointExists(conn, "link add", declared.ToType, declared.ToID, relation, declared.FromType, declared.FromID); err != nil {
		return Link{}, Link{}, err
	}
	if !relationAllowedForTypes(declared.Relation, declared.FromType, declared.ToType) {
		return Link{}, Link{}, &IncompatibleRelationError{
			Operation: "link add",
			Relation:  relation,
			FromType:  declared.FromType,
			FromID:    declared.FromID,
			ToType:    declared.ToType,
			ToID:      declared.ToID,
		}
	}
	if exists, err := linkExists(conn, declared); err != nil {
		return Link{}, Link{}, err
	} else if exists {
		return Link{}, Link{}, &DuplicateLinkError{Operation: "link add", FromID: declared.FromID, Relation: declared.Relation, ToID: declared.ToID}
	}
	if path, err := dependencyCyclePath(conn, declared); err != nil {
		return Link{}, Link{}, err
	} else if path != nil {
		return Link{}, Link{}, &DependencyCycleError{Operation: "link add", Path: path}
	}

	inverseRelationName, _ := inverseRelation(relation)
	inverse := Link{
		FromType: declared.ToType,
		FromID:   declared.ToID,
		Relation: inverseRelationName,
		ToType:   declared.FromType,
		ToID:     declared.FromID,
	}
	return declared, inverse, nil
}

func resolveEndpointType(conn *sql.DB, suppliedType, endpointID string) (string, error) {
	if suppliedType != "" {
		return suppliedType, nil
	}
	matches := []string{}
	for _, candidate := range []string{"issue"} {
		exists, err := endpointExists(conn, candidate, endpointID)
		if err != nil {
			return "", err
		}
		if exists {
			matches = append(matches, candidate)
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	return "", &InvalidEndpointTypeError{Operation: "link", Argument: "endpoint type", Value: suppliedType}
}

func linkAddArgv(fromID, fromType, relation, toID, toType, actor, key string) []string {
	return []string{"fid", "link", "add", fromID, relation, toID, "--from-type", fromType, "--to-type", toType, "--actor", actor, "--idempotency-key", key}
}

func linkInsertEffect(link Link, createdAt string) journal.Effect {
	return journal.Effect{
		Op:    "insert",
		Table: "links",
		PK: map[string]any{
			"from_type": link.FromType,
			"from_id":   link.FromID,
			"relation":  link.Relation,
			"to_type":   link.ToType,
			"to_id":     link.ToID,
		},
		Values: map[string]any{
			"created_at": createdAt,
		},
	}
}

func linkAddActivityEffect(activityID, eventID, actorRole, occurredAt string, declared, inverse Link) journal.Effect {
	return journal.Effect{
		Op:    "activity",
		Table: "activity_log",
		PK:    map[string]any{"id": activityID},
		Values: map[string]any{
			"journal_event_id": eventID,
			"ref_type":         declared.FromType,
			"ref_id":           declared.FromID,
			"event":            "link-added",
			"details": map[string]any{
				"from_type":         declared.FromType,
				"from_id":           declared.FromID,
				"relation":          declared.Relation,
				"to_type":           declared.ToType,
				"to_id":             declared.ToID,
				"inverse_relation":  inverse.Relation,
				"inverse_from_type": inverse.FromType,
				"inverse_from_id":   inverse.FromID,
				"inverse_to_type":   inverse.ToType,
				"inverse_to_id":     inverse.ToID,
			},
			"actor":       actorRole,
			"occurred_at": occurredAt,
		},
	}
}

func assertEndpointExists(conn *sql.DB, operation, endpointType, endpointID, relation, otherType, otherID string) error {
	exists, err := endpointExists(conn, endpointType, endpointID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return &EndpointNotFoundError{
		Operation:    operation,
		EndpointType: endpointType,
		EndpointID:   endpointID,
		Relation:     relation,
		OtherType:    otherType,
		OtherID:      otherID,
	}
}

func endpointExists(conn *sql.DB, endpointType, endpointID string) (bool, error) {
	query := ""
	switch endpointType {
	case "issue":
		query = `SELECT 1 FROM issues WHERE id = ? LIMIT 1`
	default:
		return false, fmt.Errorf("link add: unsupported endpoint type %q", endpointType)
	}

	var marker int
	err := conn.QueryRow(query, endpointID).Scan(&marker)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("link add: query %s %q: %w", endpointType, endpointID, err)
}

func linkExists(conn *sql.DB, link Link) (bool, error) {
	var marker int
	err := conn.QueryRow(`
		SELECT 1
		FROM links
		WHERE from_type = ? AND from_id = ? AND relation = ? AND to_type = ? AND to_id = ?
		LIMIT 1
	`, link.FromType, link.FromID, link.Relation, link.ToType, link.ToID).Scan(&marker)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("link add: query existing link %s %s %s: %w", link.FromID, link.Relation, link.ToID, err)
}

func closeLinksConn(conn *sql.DB) {
	_ = conn.Close()
}
