package issues

import (
	"database/sql"
	"testing"
)

type readyRow struct {
	id         string
	kind       string
	parentID   *string
	title      string
	status     string
	workflowID string
	ownerRole  string
	sequence   *int
	priority   *int
	createdAt  string
	updatedAt  string
	scheduling *string
}

func seedReadyRows(t *testing.T, conn *sql.DB, rows []readyRow) {
	t.Helper()

	for _, row := range rows {
		if _, err := conn.Exec(`
			INSERT INTO issues (
				id, kind, parent_id, title, status, workflow_id, owner_role, sequence, priority, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			row.id,
			row.kind,
			row.parentID,
			row.title,
			row.status,
			row.workflowID,
			row.ownerRole,
			row.sequence,
			row.priority,
			row.createdAt,
			row.updatedAt,
		); err != nil {
			t.Fatalf("seed ready issue %q: %v", row.id, err)
		}
		if row.scheduling != nil {
			if _, err := conn.Exec(`INSERT INTO epic_metadata (issue_id, scheduling) VALUES (?, ?)`, row.id, *row.scheduling); err != nil {
				t.Fatalf("seed ready epic_metadata %q: %v", row.id, err)
			}
		}
	}
}
