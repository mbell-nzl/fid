package main

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

type readyCommandRow struct {
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

func journalEventCount(t *testing.T, workspace string) int {
	t.Helper()
	count := 0
	for _, path := range mustJournalPaths(t, workspace) {
		count += len(splitJournalLines(mustReadCommandFile(t, path)))
	}
	return count
}

func stringPointer(value string) *string {
	result := value
	return &result
}

func intPointer(value int) *int {
	result := value
	return &result
}

func expectedCommandStructure(t *testing.T, value any) any {
	t.Helper()

	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal(expected structure) error = %v", err)
	}

	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	decoder.UseNumber()

	var result any
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("Decode(expected structure) error = %v", err)
	}
	return result
}

func seedReadyCommandRows(t *testing.T, conn *sql.DB, rows []readyCommandRow) {
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

func mustReadLastJournalEvent(t *testing.T, journalPath string) journal.Event {
	t.Helper()
	lines := splitJournalLines(mustReadCommandFile(t, journalPath))
	if len(lines) == 0 {
		t.Fatal("journal lines = 0, want at least one line")
	}
	var event journal.Event
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &event); err != nil {
		t.Fatalf("Unmarshal(last journal event) error = %v", err)
	}
	return event
}

func p24BindingPersistenceSnapshot(t *testing.T, workspace string) workspacePersistenceSnapshot {
	t.Helper()
	return snapshotWorkspacePersistence(t, workspace)
}
