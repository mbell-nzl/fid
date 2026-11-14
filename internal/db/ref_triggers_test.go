package db_test

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestPolymorphicRefTriggersRejectDanglingRefsAndAllowResolvableRows(t *testing.T) {
	t.Parallel()

	conn := openRefTriggerDB(t)
	defer closeRefTriggerDB(t, conn)

	insertRefTriggerIssue(t, conn, issueRow{
		id:         "FID-ref-task",
		kind:       "task",
		workflowID: "task-default",
		status:     "todo",
		title:      "Reference task",
	})
	insertRefTriggerIssue(t, conn, issueRow{id: "FID-ref-peer", kind: "task", workflowID: "task-default", status: "todo", title: "Reference peer"})

	t.Run("links reject dangling from endpoint", func(t *testing.T) {
		_, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-missing-issue", "blocks", "issue", "FID-ref-peer", "2026-06-07T15:02:00Z")
		assertRefSQLErrorContains(t, err, "links: from endpoint must resolve to an existing row")
	})

	t.Run("links reject dangling to endpoint", func(t *testing.T) {
		_, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-ref-task", "blocks", "issue", "FID-missing-peer", "2026-06-07T15:03:00Z")
		assertRefSQLErrorContains(t, err, "links: to endpoint must resolve to an existing row")
	})

	t.Run("links accept resolvable endpoints", func(t *testing.T) {
		if _, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-ref-task", "relates_to", "issue", "FID-ref-peer", "2026-06-07T15:04:00Z"); err != nil {
			t.Fatalf("insert links(valid) error = %v", err)
		}
	})

	t.Run("activity log rejects dangling ref", func(t *testing.T) {
		_, err := conn.Exec(`
			INSERT INTO activity_log (id, journal_event_id, ref_type, ref_id, event, details, actor, occurred_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, "ACT-missing", "EVT-20260607-000701", "issue", "FID-missing-issue", "status_changed", `{"from":"todo","to":"doing"}`, "developer", "2026-06-07T15:05:00Z")
		assertRefSQLErrorContains(t, err, "activity_log: ref_type/ref_id must resolve to an existing row")
	})

	t.Run("activity log accepts resolvable ref", func(t *testing.T) {
		if _, err := conn.Exec(`
			INSERT INTO activity_log (id, journal_event_id, ref_type, ref_id, event, details, actor, occurred_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, "ACT-valid", "EVT-20260607-000702", "issue", "FID-ref-task", "status_changed", `{"from":"todo","to":"doing"}`, "developer", "2026-06-07T15:06:00Z"); err != nil {
			t.Fatalf("insert activity_log(valid) error = %v", err)
		}
	})
}

func TestLinkScopeAndRelationTriggersRejectInvalidCombosAndAllowValidRows(t *testing.T) {
	t.Parallel()

	conn := openRefTriggerDB(t)
	defer closeRefTriggerDB(t, conn)

	insertRefTriggerIssue(t, conn, issueRow{
		id:         "FID-link-epic",
		kind:       "epic",
		workflowID: "epic-default",
		status:     "todo",
		title:      "Link epic",
	})
	insertEpicMetadataRow(t, conn, "FID-link-epic", "prioritised")
	insertRefTriggerIssue(t, conn, issueRow{
		id:         "FID-link-task-a",
		kind:       "task",
		parentID:   sql.NullString{String: "FID-link-epic", Valid: true},
		workflowID: "task-default",
		status:     "todo",
		title:      "Link task A",
	})
	insertRefTriggerIssue(t, conn, issueRow{
		id:         "FID-link-task-b",
		kind:       "task",
		parentID:   sql.NullString{String: "FID-link-epic", Valid: true},
		workflowID: "task-default",
		status:     "todo",
		title:      "Link task B",
	})

	t.Run("rejects decision endpoint", func(t *testing.T) {
		_, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "decision", "DEC-link-001", "supersedes", "issue", "FID-link-task-a", "2026-06-07T15:10:00Z")
		assertRefSQLErrorContains(t, err, "links: endpoint types must be issue or requirement in v0.1")
	})

	t.Run("rejects escalation endpoint", func(t *testing.T) {
		_, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-link-task-a", "relates_to", "escalation", "ESC-link-001", "2026-06-07T15:11:00Z")
		assertRefSQLErrorContains(t, err, "links: endpoint types must be issue or requirement in v0.1")
	})

	t.Run("rejects incompatible relation", func(t *testing.T) {
		_, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-link-task-a", "implements", "issue", "FID-link-task-b", "2026-06-07T15:12:00Z")
		assertRefSQLErrorContains(t, err, "links: relation/endpoint combination is not allowed in v0.1")
	})

	t.Run("accepts blocks issue to issue", func(t *testing.T) {
		if _, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-link-task-a", "blocks", "issue", "FID-link-task-b", "2026-06-07T15:14:00Z"); err != nil {
			t.Fatalf("insert links(blocks) error = %v", err)
		}
	})

}

func TestValidJournalProjectedLinkAndActivityRowsStillApplyWithRefTriggersActive(t *testing.T) {
	t.Parallel()

	conn := openRefTriggerDB(t)
	defer closeRefTriggerDB(t, conn)

	event := journal.Event{
		EventID:     "EVT-20260607-000801",
		EventSchema: journal.EventSchemaV1,
		EventType:   "phase2.ref-trigger.control",
		OccurredAt:  time.Date(2026, time.June, 7, 15, 30, 0, 0, time.UTC),
		Actor: journal.Actor{
			Role: "developer",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "replay"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []journal.Effect{
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-journal-ref-task",
				},
				Values: map[string]any{
					"kind":        "task",
					"title":       "Journal ref task",
					"workflow_id": "task-default",
					"status":      "todo",
					"created_at":  "2026-06-07T15:30:00Z",
					"updated_at":  "2026-06-07T15:30:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-journal-ref-peer",
				},
				Values: map[string]any{"kind": "task", "title": "Journal peer", "workflow_id": "task-default", "status": "todo", "created_at": "2026-06-07T15:30:00Z", "updated_at": "2026-06-07T15:30:00Z"},
			},
			{
				Op:    "insert",
				Table: "links",
				PK: map[string]any{
					"from_type": "issue",
					"from_id":   "FID-journal-ref-task",
					"relation":  "blocks",
					"to_type":   "issue",
					"to_id":     "FID-journal-ref-peer",
				},
				Values: map[string]any{
					"created_at": "2026-06-07T15:32:00Z",
				},
			},
			{
				Op:    "activity",
				Table: "activity_log",
				PK: map[string]any{
					"id": "ACT-journal-ref-001",
				},
				Values: map[string]any{
					"journal_event_id": "EVT-20260607-000801",
					"ref_type":         "issue",
					"ref_id":           "FID-journal-ref-task",
					"event":            "status_changed",
					"details":          `{"from":"todo","to":"doing"}`,
					"actor":            "developer",
					"occurred_at":      "2026-06-07T15:33:00Z",
				},
			},
		},
	}

	projectRefDBEvent(t, conn, event)

	assertRefRowCount(t, conn, "links", 1)
	assertRefRowCount(t, conn, "activity_log", 1)
}

func openRefTriggerDB(t *testing.T) *sql.DB {
	t.Helper()

	return openSchema17DB(t, testRefTriggerMigrationsDir(t))
}

func testRefTriggerMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func closeRefTriggerDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func assertRefSQLErrorContains(t *testing.T, err error, want string) {
	t.Helper()

	if err == nil {
		t.Fatalf("error = nil, want substring %q", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error = %q, want substring %q", err.Error(), want)
	}
}

func insertRefTriggerIssue(t *testing.T, conn *sql.DB, row issueRow) {
	t.Helper()

	if err := execInsertIssue(conn, row); err != nil {
		t.Fatalf("execInsertIssue(%q) error = %v", row.id, err)
	}
}

func projectRefDBEvent(t *testing.T, conn *sql.DB, event journal.Event) {
	t.Helper()

	eventHash, err := journal.HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}

	if err := journal.ProjectEvent(conn, journal.ProjectionInput{
		Event:     event,
		EventHash: eventHash,
		File:      "2026-06-07.jsonl",
		Line:      1,
		Offset:    4096,
	}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}
}

func assertRefRowCount(t *testing.T, conn *sql.DB, table string, want int) {
	t.Helper()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&count); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if count != want {
		t.Fatalf("%s row count = %d, want %d", table, count, want)
	}
}
