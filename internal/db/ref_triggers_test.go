package db_test

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
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
	insertRefTriggerRequirement(t, conn, "REQ-ref-001")

	t.Run("claims reject dangling issue ref", func(t *testing.T) {
		_, err := conn.Exec(`
			INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-missing-issue", "developer", "2026-06-07T15:00:00Z", 1800, "missing issue")
		assertRefSQLErrorContains(t, err, "claims: ref_type/ref_id must resolve to an existing row")
	})

	t.Run("claims accept resolvable issue ref", func(t *testing.T) {
		if _, err := conn.Exec(`
			INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-ref-task", "developer", "2026-06-07T15:01:00Z", 1800, "task claim"); err != nil {
			t.Fatalf("insert claims(valid) error = %v", err)
		}
	})

	t.Run("links reject dangling from endpoint", func(t *testing.T) {
		_, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-missing-issue", "implements", "requirement", "REQ-ref-001", "2026-06-07T15:02:00Z")
		assertRefSQLErrorContains(t, err, "links: from endpoint must resolve to an existing row")
	})

	t.Run("links reject dangling to endpoint", func(t *testing.T) {
		_, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-ref-task", "implements", "requirement", "REQ-missing", "2026-06-07T15:03:00Z")
		assertRefSQLErrorContains(t, err, "links: to endpoint must resolve to an existing row")
	})

	t.Run("links accept resolvable endpoints", func(t *testing.T) {
		if _, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-ref-task", "relates_to", "requirement", "REQ-ref-001", "2026-06-07T15:04:00Z"); err != nil {
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
	insertRefTriggerRequirement(t, conn, "REQ-link-001")
	insertRefTriggerDecision(t, conn, "DEC-link-001", "FID-link-epic")
	insertRefTriggerEscalation(t, conn, "ESC-link-001")

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

	t.Run("accepts implements issue to requirement", func(t *testing.T) {
		if _, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-link-task-a", "implements", "requirement", "REQ-link-001", "2026-06-07T15:13:00Z"); err != nil {
			t.Fatalf("insert links(implements) error = %v", err)
		}
	})

	t.Run("accepts blocks issue to issue", func(t *testing.T) {
		if _, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-link-task-a", "blocks", "issue", "FID-link-task-b", "2026-06-07T15:14:00Z"); err != nil {
			t.Fatalf("insert links(blocks) error = %v", err)
		}
	})

	t.Run("accepts relates_to issue to requirement", func(t *testing.T) {
		if _, err := conn.Exec(`
			INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", "FID-link-task-b", "relates_to", "requirement", "REQ-link-001", "2026-06-07T15:15:00Z"); err != nil {
			t.Fatalf("insert links(relates_to) error = %v", err)
		}
	})
}

func TestClaimScopeTriggersRejectUnsupportedRefTypeAndAllowIssueRef(t *testing.T) {
	t.Parallel()

	conn := openRefTriggerDB(t)
	defer closeRefTriggerDB(t, conn)

	insertRefTriggerIssue(t, conn, issueRow{
		id:         "FID-claim-task",
		kind:       "task",
		workflowID: "task-default",
		status:     "todo",
		title:      "Claim task",
	})

	_, err := conn.Exec(`
		INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "qa_review", "QA-001", "developer", "2026-06-07T15:20:00Z", 1800, "unsupported type")
	assertRefSQLErrorContains(t, err, "claims: ref_type must be issue in v0.1")

	if _, err := conn.Exec(`
		INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "issue", "FID-claim-task", "developer", "2026-06-07T15:21:00Z", 1800, "supported type"); err != nil {
		t.Fatalf("insert claims(issue) error = %v", err)
	}
}

func TestValidJournalProjectedRowsStillApplyWithRefTriggersActive(t *testing.T) {
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
				Table: "requirements",
				PK: map[string]any{
					"id": "REQ-journal-ref-001",
				},
				Values: map[string]any{
					"type":           "PRD",
					"summary":        "Journal requirement",
					"source_doc":     "01_PRD.md",
					"source_section": "OPS",
					"status":         "active",
				},
			},
			{
				Op:    "insert",
				Table: "claims",
				PK: map[string]any{
					"ref_type": "issue",
					"ref_id":   "FID-journal-ref-task",
				},
				Values: map[string]any{
					"claimed_by":  "developer",
					"claimed_at":  "2026-06-07T15:31:00Z",
					"ttl_seconds": 1800,
					"note":        "journal claim",
				},
			},
			{
				Op:    "insert",
				Table: "links",
				PK: map[string]any{
					"from_type": "issue",
					"from_id":   "FID-journal-ref-task",
					"relation":  "implements",
					"to_type":   "requirement",
					"to_id":     "REQ-journal-ref-001",
				},
				Values: map[string]any{
					"created_at": "2026-06-07T15:32:00Z",
				},
			},
			{
				Op:    "insert",
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

	assertRefRowCount(t, conn, "claims", 1)
	assertRefRowCount(t, conn, "links", 1)
	assertRefRowCount(t, conn, "activity_log", 1)
}

func openRefTriggerDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testRefTriggerMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	return conn
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

func insertRefTriggerRequirement(t *testing.T, conn *sql.DB, requirementID string) {
	t.Helper()

	if _, err := conn.Exec(`
		INSERT INTO requirements (id, type, summary, source_doc, source_section, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`, requirementID, "PRD", "Requirement "+requirementID, "01_PRD.md", "OPS", "active"); err != nil {
		t.Fatalf("insert requirements(%q) error = %v", requirementID, err)
	}
}

func insertRefTriggerDecision(t *testing.T, conn *sql.DB, decisionID, epicID string) {
	t.Helper()

	if _, err := conn.Exec(`
		INSERT INTO decisions (id, epic_id, author_role, summary, body_md_path, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, decisionID, epicID, "tl", "Decision "+decisionID, "project-docs/decisions/"+decisionID+".md", "open", "2026-06-07T15:08:00Z"); err != nil {
		t.Fatalf("insert decisions(%q) error = %v", decisionID, err)
	}
}

func insertRefTriggerEscalation(t *testing.T, conn *sql.DB, escalationID string) {
	t.Helper()

	if _, err := conn.Exec(`
		INSERT INTO escalations (id, status, trigger_class, raised_by_role, raised_at)
		VALUES (?, ?, ?, ?, ?)
	`, escalationID, "open", "policy", "qa_lead", "2026-06-07T15:09:00Z"); err != nil {
		t.Fatalf("insert escalations(%q) error = %v", escalationID, err)
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
