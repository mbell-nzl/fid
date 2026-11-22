// Package workrecord reads task-bound person decisions and creation attribution.
package workrecord

import (
	"database/sql"
	"fmt"
)

type Queryer interface {
	Query(string, ...any) (*sql.Rows, error)
}

// Unaccepted includes every open agent-created task or subtask, irrespective
// of creation date. Creation attribution is derived from canonical Events.
func Unaccepted(conn Queryer) (map[string]string, error) {
	rows, err := conn.Query(`SELECT i.id,c.actor FROM issues i JOIN issue_creators c ON c.issue_id=i.id
 JOIN workflow_states s ON s.workflow_id=i.workflow_id AND s.state=i.status
 WHERE i.kind IN ('task','subtask') AND s.category NOT IN ('completed','cancelled')
 AND c.actor LIKE 'agent:%' AND NOT EXISTS(SELECT 1 FROM person_decisions d WHERE d.kind='accept' AND d.subject_id=i.id) ORDER BY i.id`)
	if err != nil {
		return nil, fmt.Errorf("read unaccepted work: %w", err)
	}
	defer func() { _ = rows.Close() }()
	result := map[string]string{}
	for rows.Next() {
		var id, actor string
		if err := rows.Scan(&id, &actor); err != nil {
			return nil, err
		}
		result[id] = actor
	}
	return result, rows.Err()
}
