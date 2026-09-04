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

func Remove(options RemoveOptions) (RemoveResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return RemoveResult{}, fmt.Errorf("link remove: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return RemoveResult{}, fmt.Errorf("link remove: migrations dir is required")
	}
	if strings.TrimSpace(options.FromID) == "" || strings.TrimSpace(options.ToID) == "" {
		return RemoveResult{}, fmt.Errorf("link remove: from and to ids are required")
	}
	actor, err := standaloneActorForOptions(options.Actor)
	if err != nil {
		return RemoveResult{}, err
	}

	fromType, err := endpointType(options.FromType, "from-type", "link remove")
	if err != nil {
		return RemoveResult{}, err
	}
	toType, err := endpointType(options.ToType, "to-type", "link remove")
	if err != nil {
		return RemoveResult{}, err
	}
	fromID := strings.TrimSpace(options.FromID)
	relation := strings.TrimSpace(options.Relation)
	toID := strings.TrimSpace(options.ToID)
	if _, ok := inverseRelation(relation); !ok {
		return RemoveResult{}, &UnknownRelationError{Operation: "link remove", Relation: relation}
	}
	if fromType == toType && fromID == toID {
		return RemoveResult{}, &SelfLinkError{Operation: "link remove", FromID: fromID, Relation: relation, ToID: toID}
	}

	commandArgv := linkRemoveArgv(fromID, fromType, relation, toID, toType, options.Actor, options.IdempotencyKey)
	var result RemoveResult
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
		digest, err := linkRequestDigest("link.remove", projectConfig.ProjectID, actor, map[string]any{"from_id": fromID, "from_type": fromType, "relation": relation, "to_id": toID, "to_type": toType})
		if err != nil {
			return err
		}
		idempotencyRequest := journal.IdempotencyRequest{OperationScope: "link.remove:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
			return err
		} else if found {
			decoded, err := removeResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
			if err != nil {
				return err
			}
			result = decoded
			return nil
		}

		declared, inverse, err := plannedRemoval(conn, fromType, fromID, relation, toType, toID)
		if err != nil {
			return err
		}
		if err := assertEndpointExists(conn, "link remove", declared.FromType, declared.FromID, relation, declared.ToType, declared.ToID); err != nil {
			return err
		}
		if err := assertEndpointExists(conn, "link remove", declared.ToType, declared.ToID, relation, declared.FromType, declared.FromID); err != nil {
			return err
		}
		storedDeclared, err := loadStoredLink(conn, declared)
		if err != nil {
			return err
		}
		storedInverse, err := loadStoredLink(conn, inverse)
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
		baseEffects := []journal.Effect{
			linkDeleteEffect(storedDeclared),
			linkDeleteEffect(storedInverse),
		}
		baseEffects = append(baseEffects, revisionEffects...)
		sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: declared.FromID + ":" + declared.Relation + ":" + declared.ToID, Type: "link"})
		if err != nil {
			return err
		}
		event, err := journal.PrepareEvent(lock, journal.Event{EventSchema: journal.EventSchemaV2, EventType: "link.removed", ActorV2: actor.eventActor(), Aggregate: journal.Aggregate{ID: declared.FromID + ":" + declared.Relation + ":" + declared.ToID, Sequence: sequence, Type: "link"}, CommandName: "link.remove", Effects: baseEffects})
		if err != nil {
			return err
		}
		event.Effects = append(baseEffects, linkRemoveActivityEffect(activityID, event.EventID, actor.id, event.OccurredAt.UTC().Format(time.RFC3339), declared, inverse))
		seed := RemoveResult{Declared: declared, Inverse: inverse, IssueRevisions: resultIssueRevisions(observedRevisions), OccurredAt: event.OccurredAt.UTC().Format(time.RFC3339Nano), IdempotencyKey: options.IdempotencyKey}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			encoded, err := removeResultJSON(seed)
			if err != nil {
				return journal.IdempotencyMutation{}, err
			}
			return journal.IdempotencyMutation{EntityID: declared.FromID + ":" + declared.Relation + ":" + declared.ToID, ResultJSON: encoded}, nil
		})
		if err != nil {
			return err
		}
		decoded, err := removeResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		if err != nil {
			return err
		}
		result = decoded
		return nil
	})
	if err != nil {
		return RemoveResult{}, err
	}

	return result, nil
}

type storedLink struct {
	Link
	CreatedAt string
}

func plannedRemoval(conn *sql.DB, fromType, fromID, relation, toType, toID string) (Link, Link, error) {
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
	if !relationAllowedForTypes(declared.Relation, declared.FromType, declared.ToType) {
		return Link{}, Link{}, &IncompatibleRelationError{
			Operation: "link remove",
			Relation:  relation,
			FromType:  declared.FromType,
			FromID:    declared.FromID,
			ToType:    declared.ToType,
			ToID:      declared.ToID,
		}
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

func loadStoredLink(conn *sql.DB, link Link) (storedLink, error) {
	var createdAt string
	err := conn.QueryRow(`
		SELECT COALESCE(created_at, '')
		FROM links
		WHERE from_type = ? AND from_id = ? AND relation = ? AND to_type = ? AND to_id = ?
		LIMIT 1
	`, link.FromType, link.FromID, link.Relation, link.ToType, link.ToID).Scan(&createdAt)
	if err == nil {
		return storedLink{Link: link, CreatedAt: createdAt}, nil
	}
	if err == sql.ErrNoRows {
		return storedLink{}, &MissingLinkError{
			Operation: "link remove",
			FromType:  link.FromType,
			FromID:    link.FromID,
			Relation:  link.Relation,
			ToType:    link.ToType,
			ToID:      link.ToID,
		}
	}
	return storedLink{}, fmt.Errorf("link remove: query existing link %s %s %s: %w", link.FromID, link.Relation, link.ToID, err)
}

func linkRemoveArgv(fromID, fromType, relation, toID, toType, actor, key string) []string {
	return []string{"fid", "link", "remove", fromID, relation, toID, "--from-type", fromType, "--to-type", toType, "--actor", actor, "--idempotency-key", key}
}

func linkDeleteEffect(link storedLink) journal.Effect {
	return journal.Effect{
		Op:    "delete",
		Table: "links",
		PK: map[string]any{
			"from_type": link.FromType,
			"from_id":   link.FromID,
			"relation":  link.Relation,
			"to_type":   link.ToType,
			"to_id":     link.ToID,
		},
		Preconditions: map[string]any{
			"created_at": link.CreatedAt,
		},
	}
}

func linkRemoveActivityEffect(activityID, eventID, actorRole, occurredAt string, declared, inverse Link) journal.Effect {
	return journal.Effect{
		Op:    "activity",
		Table: "activity_log",
		PK:    map[string]any{"id": activityID},
		Values: map[string]any{
			"journal_event_id": eventID,
			"ref_type":         declared.FromType,
			"ref_id":           declared.FromID,
			"event":            "link-removed",
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

func linkActivityIDExists(conn *sql.DB, activityID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM activity_log WHERE id = ? LIMIT 1`, activityID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query activity id %q: %w", activityID, err)
}
