package links

import (
	"database/sql"
	"fmt"
)

func dependencyCyclePath(conn *sql.DB, link Link) ([]string, error) {
	fromID, toID, ok := normalizedBlocksEndpoints(link)
	if !ok {
		return nil, nil
	}

	path, found, err := blocksPath(conn, toID, fromID)
	if err != nil || !found {
		return nil, err
	}
	return append([]string{fromID}, path...), nil
}

func normalizedBlocksEndpoints(link Link) (string, string, bool) {
	if link.FromType != "issue" || link.ToType != "issue" {
		return "", "", false
	}
	switch link.Relation {
	case "blocks":
		return link.FromID, link.ToID, true
	case "blocked_by":
		return link.ToID, link.FromID, true
	default:
		return "", "", false
	}
}

func isDependencyRelation(relation string) bool {
	return relation == "blocks" || relation == "blocked_by"
}

func blocksPath(conn *sql.DB, fromID, targetID string) ([]string, bool, error) {
	type pathNode struct {
		id   string
		path []string
	}
	queue := []pathNode{{id: fromID, path: []string{fromID}}}
	visited := map[string]bool{fromID: true}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		neighbors, err := blocksNeighbors(conn, current.id)
		if err != nil {
			return nil, false, err
		}
		for _, neighbor := range neighbors {
			path := append(append([]string{}, current.path...), neighbor)
			if neighbor == targetID {
				return path, true, nil
			}
			if visited[neighbor] {
				continue
			}
			visited[neighbor] = true
			queue = append(queue, pathNode{id: neighbor, path: path})
		}
	}
	return nil, false, nil
}

func blocksNeighbors(conn *sql.DB, issueID string) ([]string, error) {
	rows, err := conn.Query(`
		SELECT to_id
		FROM links
		WHERE from_type = 'issue' AND from_id = ?
			AND relation = 'blocks' AND to_type = 'issue'
		ORDER BY to_id
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query dependency edges for %q: %w", issueID, err)
	}
	defer func() { _ = rows.Close() }()

	neighbors := []string{}
	for rows.Next() {
		var neighbor string
		if err := rows.Scan(&neighbor); err != nil {
			return nil, fmt.Errorf("scan dependency edge for %q: %w", issueID, err)
		}
		neighbors = append(neighbors, neighbor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dependency edges for %q: %w", issueID, err)
	}
	return neighbors, nil
}
