package links

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func List(options ListOptions) ([]ListItem, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return nil, fmt.Errorf("link list: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return nil, fmt.Errorf("link list: migrations dir is required")
	}
	entityID := strings.TrimSpace(options.EntityID)
	if entityID == "" {
		return nil, fmt.Errorf("link list: entity id is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}

	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer closeLinksConn(conn)

	return queryLinkList(conn, entityID)
}

func queryLinkList(conn *sql.DB, entityID string) ([]ListItem, error) {
	rows, err := conn.Query(`
		WITH entity_links AS (
			SELECT
				'outgoing' AS direction,
				relation,
				to_type AS entity_type,
				to_id AS entity_id,
				COALESCE(created_at, '') AS created_at
			FROM links
			WHERE from_id = ?
			UNION ALL
			SELECT
				'incoming' AS direction,
				relation,
				from_type AS entity_type,
				from_id AS entity_id,
				COALESCE(created_at, '') AS created_at
			FROM links
			WHERE to_id = ?
		)
		SELECT
			el.direction,
			el.relation,
			el.entity_type,
			el.entity_id,
			CASE
				WHEN el.entity_type = 'issue' THEN issue.title
				WHEN el.entity_type = 'requirement' THEN requirement.summary
				ELSE NULL
			END AS entity_title
		FROM entity_links el
		LEFT JOIN issues issue
			ON el.entity_type = 'issue' AND issue.id = el.entity_id
		LEFT JOIN requirements requirement
			ON el.entity_type = 'requirement' AND requirement.id = el.entity_id
		ORDER BY
			CASE el.direction WHEN 'outgoing' THEN 0 ELSE 1 END,
			el.created_at,
			el.relation,
			el.entity_type,
			el.entity_id
	`, entityID, entityID)
	if err != nil {
		return nil, fmt.Errorf("query link list for %q: %w", entityID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	items := []ListItem{}
	for rows.Next() {
		var (
			item        ListItem
			entityTitle sql.NullString
		)
		if err := rows.Scan(&item.Direction, &item.Relation, &item.EntityType, &item.EntityID, &entityTitle); err != nil {
			return nil, fmt.Errorf("scan link list row for %q: %w", entityID, err)
		}
		if entityTitle.Valid {
			item.EntityTitle = listStringPointer(entityTitle.String)
		}
		if !isListRelation(item.Relation) {
			continue
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate link list rows for %q: %w", entityID, err)
	}
	return items, nil
}

func listStringPointer(value string) *string {
	return &value
}

func isListRelation(relation string) bool {
	switch relation {
	case "blocks", "relates_to", "duplicates", "supersedes", "implements", "caused_by", "discovered_from":
		return true
	default:
		return false
	}
}
