package links

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

const actionLinkRemove = "link.remove"

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
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionLinkRemove); err != nil {
		return RemoveResult{}, err
	}

	fromID := strings.TrimSpace(options.FromID)
	relation := strings.TrimSpace(options.Relation)
	toID := strings.TrimSpace(options.ToID)
	if _, ok := inverseRelation(relation); !ok {
		return RemoveResult{}, &UnknownRelationError{Operation: "link remove", Relation: relation}
	}
	if fromID == toID {
		return RemoveResult{}, &SelfLinkError{Operation: "link remove", FromID: fromID, Relation: relation, ToID: toID}
	}

	commandArgv := linkRemoveArgv(fromID, relation, toID, options.ActorRole)
	var result RemoveResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeLinksConn(conn)

		declared, inverse, err := plannedRemoval(fromID, relation, toID)
		if err != nil {
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

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "link.removed",
			Actor:       journal.Actor{Role: options.ActorRole},
			Command: journal.Command{
				Argv:       commandArgv,
				CWD:        ".",
				FIDVersion: fidVersion,
			},
		})
		if err != nil {
			return err
		}

		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return linkActivityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		occurredAt := event.OccurredAt.UTC().Format(time.RFC3339)
		event.Effects = []journal.Effect{
			linkDeleteEffect(storedDeclared),
			linkDeleteEffect(storedInverse),
			linkRemoveActivityEffect(activityID, event.EventID, options.ActorRole, occurredAt, declared, inverse),
		}

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}
		if err := projectLinksEvent(options.WorkspaceRoot, conn, appended); err != nil {
			return err
		}

		result = RemoveResult{Declared: declared, Inverse: inverse}
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

func plannedRemoval(fromID, relation, toID string) (Link, Link, error) {
	declared := Link{
		FromType: endpointTypeForID(fromID),
		FromID:   fromID,
		Relation: relation,
		ToType:   endpointTypeForID(toID),
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
			FromID:    link.FromID,
			Relation:  link.Relation,
			ToID:      link.ToID,
		}
	}
	return storedLink{}, fmt.Errorf("link remove: query existing link %s %s %s: %w", link.FromID, link.Relation, link.ToID, err)
}

func linkRemoveArgv(fromID, relation, toID, actorRole string) []string {
	return []string{"fid", "link", "remove", fromID, relation, toID, "--by", actorRole}
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
			"event":            "link_removed",
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
