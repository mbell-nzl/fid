package links

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

const (
	actionLinkAdd = "link.add"
	fidVersion    = "0.1.0"
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
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionLinkAdd); err != nil {
		return AddResult{}, err
	}

	declaredRelation := strings.TrimSpace(options.Relation)
	if _, ok := inverseRelation(declaredRelation); !ok {
		return AddResult{}, &UnknownRelationError{Operation: "link add", Relation: declaredRelation}
	}

	fromID := strings.TrimSpace(options.FromID)
	toID := strings.TrimSpace(options.ToID)
	if fromID == toID {
		return AddResult{}, &SelfLinkError{Operation: "link add", FromID: fromID, Relation: declaredRelation, ToID: toID}
	}

	var result AddResult
	commandArgv := linkAddArgv(fromID, declaredRelation, toID, options.ActorRole)
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

		declared, inverse, err := plannedLinks(conn, fromID, declaredRelation, toID)
		if err != nil {
			return err
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "link.added",
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
		createdAt := event.OccurredAt.UTC().Format("2006-01-02T15:04:05Z07:00")
		event.Effects = []journal.Effect{
			linkInsertEffect(declared, createdAt),
			linkInsertEffect(inverse, createdAt),
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

		result = AddResult{Declared: declared, Inverse: inverse}
		return nil
	})
	if err != nil {
		return AddResult{}, err
	}

	return result, nil
}

func plannedLinks(conn *sql.DB, fromID, relation, toID string) (Link, Link, error) {
	declared := Link{
		FromType: endpointTypeForID(fromID),
		FromID:   fromID,
		Relation: relation,
		ToType:   endpointTypeForID(toID),
		ToID:     toID,
	}
	if err := assertEndpointExists(conn, declared.FromType, declared.FromID, relation, declared.ToType, declared.ToID); err != nil {
		return Link{}, Link{}, err
	}
	if err := assertEndpointExists(conn, declared.ToType, declared.ToID, relation, declared.FromType, declared.FromID); err != nil {
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

func linkAddArgv(fromID, relation, toID, actorRole string) []string {
	return []string{"fid", "link", "add", fromID, relation, toID, "--by", actorRole}
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

func assertEndpointExists(conn *sql.DB, endpointType, endpointID, relation, otherType, otherID string) error {
	exists, err := endpointExists(conn, endpointType, endpointID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return &EndpointNotFoundError{
		Operation:    "link add",
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
	case "requirement":
		query = `SELECT 1 FROM requirements WHERE id = ? LIMIT 1`
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

func projectLinksEvent(workspaceRoot string, conn *sql.DB, event journal.Event) error {
	eventHash, err := journal.HashEvent(event)
	if err != nil {
		return fmt.Errorf("hash link event %q: %w", event.EventID, err)
	}

	journalFile := event.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return fmt.Errorf("stat link journal %q: %w", journalPath, err)
	}

	return journal.ProjectEvent(conn, journal.ProjectionInput{
		Event:     event,
		EventHash: eventHash,
		File:      journalFile,
		Line:      countLinkJournalLines(journalPath),
		Offset:    info.Size(),
	})
}

func countLinkJournalLines(path string) int {
	payload, err := os.ReadFile(path)
	if err != nil {
		return 1
	}
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func closeLinksConn(conn *sql.DB) {
	_ = conn.Close()
}
