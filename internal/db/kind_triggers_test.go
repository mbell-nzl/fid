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

func TestIssueParentKindTriggerRejectsInvalidHierarchiesAndAllowsValidRows(t *testing.T) {
	t.Parallel()

	conn := openKindTriggerDB(t)
	defer closeKindTriggerDB(t, conn)

	insertIssueRow(t, conn, issueRow{
		id:         "FID-parent-epic",
		kind:       "epic",
		workflowID: "epic-default",
		status:     "todo",
		title:      "Parent epic",
	})
	insertEpicMetadataRow(t, conn, "FID-parent-epic", "prioritised")
	insertIssueRow(t, conn, issueRow{
		id:         "FID-parent-task",
		kind:       "task",
		workflowID: "task-default",
		status:     "todo",
		title:      "Parent task",
	})
	insertIssueRow(t, conn, issueRow{
		id:         "FID-parent-subtask",
		kind:       "subtask",
		workflowID: "subtask-default",
		status:     "todo",
		title:      "Parent subtask",
		parentID:   sql.NullString{String: "FID-parent-task", Valid: true},
	})

	tests := []struct {
		name    string
		row     issueRow
		wantErr string
	}{
		{
			name: "rejects epic with parent",
			row: issueRow{
				id:         "FID-invalid-epic",
				kind:       "epic",
				workflowID: "epic-default",
				status:     "todo",
				title:      "Invalid epic",
				parentID:   sql.NullString{String: "FID-parent-epic", Valid: true},
			},
			wantErr: "issues: epic parent_id must be NULL",
		},
		{
			name: "rejects task with subtask parent",
			row: issueRow{
				id:         "FID-invalid-task",
				kind:       "task",
				workflowID: "task-default",
				status:     "todo",
				title:      "Invalid task",
				parentID:   sql.NullString{String: "FID-parent-subtask", Valid: true},
			},
			wantErr: "issues: task parent_id must reference an epic or be NULL",
		},
		{
			name: "rejects subtask without parent",
			row: issueRow{
				id:         "FID-invalid-subtask-nil",
				kind:       "subtask",
				workflowID: "subtask-default",
				status:     "todo",
				title:      "Invalid subtask without parent",
			},
			wantErr: "issues: subtask parent_id must reference a task",
		},
		{
			name: "rejects subtask with epic parent",
			row: issueRow{
				id:         "FID-invalid-subtask-epic",
				kind:       "subtask",
				workflowID: "subtask-default",
				status:     "todo",
				title:      "Invalid subtask with epic parent",
				parentID:   sql.NullString{String: "FID-parent-epic", Valid: true},
			},
			wantErr: "issues: subtask parent_id must reference a task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := execInsertIssue(conn, tt.row)
			assertSQLErrorContains(t, err, tt.wantErr)
		})
	}

	validRows := []issueRow{
		{
			id:         "FID-valid-epic",
			kind:       "epic",
			workflowID: "epic-default",
			status:     "todo",
			title:      "Valid epic",
		},
		{
			id:         "FID-valid-task-standalone",
			kind:       "task",
			workflowID: "task-default",
			status:     "todo",
			title:      "Valid standalone task",
		},
		{
			id:         "FID-valid-task-child",
			kind:       "task",
			workflowID: "task-default",
			status:     "todo",
			title:      "Valid task under epic",
			parentID:   sql.NullString{String: "FID-parent-epic", Valid: true},
		},
		{
			id:         "FID-valid-subtask",
			kind:       "subtask",
			workflowID: "subtask-default",
			status:     "todo",
			title:      "Valid subtask under task",
			parentID:   sql.NullString{String: "FID-parent-task", Valid: true},
		},
	}

	for _, row := range validRows {
		t.Run("accepts "+row.id, func(t *testing.T) {
			if err := execInsertIssue(conn, row); err != nil {
				t.Fatalf("execInsertIssue(%q) error = %v", row.id, err)
			}
		})
	}
}

func TestIssueSequenceSchedulingTriggerRejectsInvalidTaskStatesAndAllowsValidRows(t *testing.T) {
	t.Parallel()

	conn := openKindTriggerDB(t)
	defer closeKindTriggerDB(t, conn)

	insertIssueRow(t, conn, issueRow{
		id:         "FID-seq-epic",
		kind:       "epic",
		workflowID: "epic-default",
		status:     "todo",
		title:      "Sequential epic",
	})
	insertEpicMetadataRow(t, conn, "FID-seq-epic", "sequential")
	insertIssueRow(t, conn, issueRow{
		id:         "FID-pri-epic",
		kind:       "epic",
		workflowID: "epic-default",
		status:     "todo",
		title:      "Prioritised epic",
	})
	insertEpicMetadataRow(t, conn, "FID-pri-epic", "prioritised")
	insertIssueRow(t, conn, issueRow{
		id:         "FID-metadata-missing-epic",
		kind:       "epic",
		workflowID: "epic-default",
		status:     "todo",
		title:      "Epic without metadata",
	})

	tests := []struct {
		name    string
		row     issueRow
		wantErr string
	}{
		{
			name: "rejects sequential epic task without sequence",
			row: issueRow{
				id:         "FID-seq-task-missing",
				kind:       "task",
				parentID:   sql.NullString{String: "FID-seq-epic", Valid: true},
				workflowID: "task-default",
				status:     "todo",
				title:      "Missing sequence",
			},
			wantErr: "issues: task under sequential epic requires sequence",
		},
		{
			name: "rejects prioritised epic task with sequence",
			row: issueRow{
				id:         "FID-pri-task-with-seq",
				kind:       "task",
				parentID:   sql.NullString{String: "FID-pri-epic", Valid: true},
				workflowID: "task-default",
				status:     "todo",
				title:      "Priority task with sequence",
				sequence:   sql.NullInt64{Int64: 1, Valid: true},
			},
			wantErr: "issues: task under prioritised epic must not have sequence",
		},
		{
			name: "rejects unattached task with sequence",
			row: issueRow{
				id:         "FID-standalone-with-seq",
				kind:       "task",
				workflowID: "task-default",
				status:     "todo",
				title:      "Standalone task with sequence",
				sequence:   sql.NullInt64{Int64: 1, Valid: true},
			},
			wantErr: "issues: unattached task must not have sequence",
		},
		{
			name: "rejects task under epic missing metadata",
			row: issueRow{
				id:         "FID-metadata-missing-task",
				kind:       "task",
				parentID:   sql.NullString{String: "FID-metadata-missing-epic", Valid: true},
				workflowID: "task-default",
				status:     "todo",
				title:      "Task under metadata-less epic",
			},
			wantErr: "issues: task parent epic must have epic_metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := execInsertIssue(conn, tt.row)
			assertSQLErrorContains(t, err, tt.wantErr)
		})
	}

	validRows := []issueRow{
		{
			id:         "FID-seq-task-valid",
			kind:       "task",
			parentID:   sql.NullString{String: "FID-seq-epic", Valid: true},
			workflowID: "task-default",
			status:     "todo",
			title:      "Sequenced task",
			sequence:   sql.NullInt64{Int64: 1, Valid: true},
		},
		{
			id:         "FID-pri-task-valid",
			kind:       "task",
			parentID:   sql.NullString{String: "FID-pri-epic", Valid: true},
			workflowID: "task-default",
			status:     "todo",
			title:      "Prioritised task",
		},
	}

	for _, row := range validRows {
		t.Run("accepts "+row.id, func(t *testing.T) {
			if err := execInsertIssue(conn, row); err != nil {
				t.Fatalf("execInsertIssue(%q) error = %v", row.id, err)
			}
		})
	}

	err := execInsertIssue(conn, issueRow{
		id:         "FID-seq-task-duplicate",
		kind:       "task",
		parentID:   sql.NullString{String: "FID-seq-epic", Valid: true},
		workflowID: "task-default",
		status:     "todo",
		title:      "Duplicate sequence",
		sequence:   sql.NullInt64{Int64: 1, Valid: true},
	})
	assertSQLErrorContains(t, err, "UNIQUE constraint failed: issues.parent_id, issues.sequence")
}

func TestIssueKindGuardTriggersRejectWrongKindsAndAllowValidRows(t *testing.T) {
	t.Parallel()

	conn := openSchema17DB(t, testKindTriggerMigrationsDir(t))
	defer closeKindTriggerDB(t, conn)

	insertIssueRow(t, conn, issueRow{
		id:         "FID-guard-epic",
		kind:       "epic",
		workflowID: "epic-default",
		status:     "todo",
		title:      "Guard epic",
	})
	insertIssueRow(t, conn, issueRow{
		id:         "FID-guard-task",
		kind:       "task",
		workflowID: "task-default",
		status:     "todo",
		title:      "Guard task",
	})
	insertIssueRow(t, conn, issueRow{
		id:         "FID-guard-subtask",
		kind:       "subtask",
		workflowID: "subtask-default",
		status:     "todo",
		title:      "Guard subtask",
		parentID:   sql.NullString{String: "FID-guard-task", Valid: true},
	})

	tests := []struct {
		name    string
		query   string
		args    []any
		wantErr string
	}{
		{
			name: "rejects acceptance criteria on epic",
			query: `
				INSERT INTO acceptance_criteria (id, issue_id, ordinal, text)
				VALUES (?, ?, ?, ?)
			`,
			args:    []any{"AC-invalid-epic", "FID-guard-epic", 1, "Epic AC"},
			wantErr: "acceptance_criteria: issue_id must reference a task or subtask",
		},

		{
			name: "rejects epic metadata on task",
			query: `
				INSERT INTO epic_metadata (issue_id, scheduling)
				VALUES (?, ?)
			`,
			args:    []any{"FID-guard-task", "sequential"},
			wantErr: "epic_metadata: issue_id must reference an epic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conn.Exec(tt.query, tt.args...)
			assertSQLErrorContains(t, err, tt.wantErr)
		})
	}

	acceptanceCriteriaCases := []struct {
		name    string
		id      string
		issueID string
	}{
		{name: "task", id: "AC-valid-task", issueID: "FID-guard-task"},
		{name: "subtask", id: "AC-valid-subtask", issueID: "FID-guard-subtask"},
	}
	for _, tt := range acceptanceCriteriaCases {
		t.Run("accepts acceptance criteria on "+tt.name, func(t *testing.T) {
			if _, err := conn.Exec(`
				INSERT INTO acceptance_criteria (id, issue_id, ordinal, text)
				VALUES (?, ?, ?, ?)
			`, tt.id, tt.issueID, 1, "Valid acceptance criterion"); err != nil {
				t.Fatalf("insert acceptance_criteria(%q) error = %v", tt.id, err)
			}
		})
	}

	if _, err := conn.Exec(`
		INSERT INTO epic_metadata (issue_id, scheduling)
		VALUES (?, ?)
	`, "FID-guard-epic", "prioritised"); err != nil {
		t.Fatalf("insert epic_metadata(valid epic) error = %v", err)
	}
}

func TestEpicReferenceKindGuardTriggersRejectWrongKindsAndAllowValidRows(t *testing.T) {
	t.Parallel()

	conn := openSchema17DB(t, testKindTriggerMigrationsDir(t))
	defer closeKindTriggerDB(t, conn)

	insertIssueRow(t, conn, issueRow{
		id:         "FID-triage-epic",
		kind:       "epic",
		workflowID: "epic-default",
		status:     "todo",
		title:      "Triage epic",
	})
	insertIssueRow(t, conn, issueRow{
		id:         "FID-triage-task",
		kind:       "task",
		workflowID: "task-default",
		status:     "todo",
		title:      "Triage task",
	})

	_, err := conn.Exec(`
		INSERT INTO triage_events (
			id, epic_id, issue_id, triager_role, source, verdict, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "TRG-invalid", "FID-triage-task", "FID-triage-task", "qa_lead", "task_completion", "BLOCK", "2026-06-07T14:10:00Z")
	assertSQLErrorContains(t, err, "triage_events: epic_id must reference an epic")

	_, err = conn.Exec(`
		INSERT INTO decisions (
			id, epic_id, author_role, summary, body_md_path, status, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "DEC-invalid", "FID-triage-task", "tl", "Invalid decision", "project-docs/decisions/DEC-invalid.md", "open", "2026-06-07T14:11:00Z")
	assertSQLErrorContains(t, err, "decisions: epic_id must reference an epic when set")

	if _, err := conn.Exec(`
		INSERT INTO triage_events (
			id, epic_id, issue_id, triager_role, source, verdict, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "TRG-valid", "FID-triage-epic", "FID-triage-task", "qa_lead", "task_completion", "RESUME", "2026-06-07T14:12:00Z"); err != nil {
		t.Fatalf("insert triage_events(valid epic) error = %v", err)
	}

	if _, err := conn.Exec(`
		INSERT INTO decisions (
			id, epic_id, author_role, summary, body_md_path, status, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "DEC-valid-epic", "FID-triage-epic", "tl", "Valid decision", "project-docs/decisions/DEC-valid-epic.md", "open", "2026-06-07T14:13:00Z"); err != nil {
		t.Fatalf("insert decisions(valid epic) error = %v", err)
	}

	if _, err := conn.Exec(`
		INSERT INTO decisions (
			id, epic_id, author_role, summary, body_md_path, status, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "DEC-valid-null", nil, "tl", "Decision without epic", "project-docs/decisions/DEC-valid-null.md", "open", "2026-06-07T14:14:00Z"); err != nil {
		t.Fatalf("insert decisions(NULL epic_id) error = %v", err)
	}
}

func TestIssueKindImmutabilityTriggerRejectsKindChangesAndAllowsOtherUpdates(t *testing.T) {
	t.Parallel()

	conn := openKindTriggerDB(t)
	defer closeKindTriggerDB(t, conn)

	insertIssueRow(t, conn, issueRow{
		id:         "FID-kind-lock",
		kind:       "task",
		workflowID: "task-default",
		status:     "todo",
		title:      "Immutable task",
	})

	_, err := conn.Exec(`
		UPDATE issues
		SET kind = ?, workflow_id = ?, status = ?, updated_at = ?
		WHERE id = ?
	`, "subtask", "subtask-default", "todo", "2026-06-07T14:20:00Z", "FID-kind-lock")
	assertSQLErrorContains(t, err, "issues: kind is immutable")

	if _, err := conn.Exec(`
		UPDATE issues
		SET title = ?, updated_at = ?
		WHERE id = ?
	`, "Renamed immutable task", "2026-06-07T14:21:00Z", "FID-kind-lock"); err != nil {
		t.Fatalf("update issues(non-kind columns) error = %v", err)
	}
}

func TestValidJournalProjectedIssueStillAppliesWithKindTriggersActive(t *testing.T) {
	t.Parallel()

	conn := openKindTriggerDB(t)
	defer closeKindTriggerDB(t, conn)

	event := journal.Event{
		EventID:     "EVT-20260607-000601",
		EventSchema: journal.EventSchemaV1,
		EventType:   "phase2.kind-trigger.control",
		OccurredAt:  time.Date(2026, time.June, 7, 14, 30, 0, 0, time.UTC),
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
					"id": "FID-journal-epic",
				},
				Values: map[string]any{
					"kind":        "epic",
					"title":       "Journal epic",
					"workflow_id": "epic-default",
					"status":      "todo",
					"created_at":  "2026-06-07T14:30:00Z",
					"updated_at":  "2026-06-07T14:30:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "epic_metadata",
				PK: map[string]any{
					"issue_id": "FID-journal-epic",
				},
				Values: map[string]any{
					"scheduling": "sequential",
				},
			},
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-journal-task",
				},
				Values: map[string]any{
					"kind":        "task",
					"parent_id":   "FID-journal-epic",
					"sequence":    1,
					"title":       "Journal task",
					"workflow_id": "task-default",
					"status":      "todo",
					"created_at":  "2026-06-07T14:31:00Z",
					"updated_at":  "2026-06-07T14:31:00Z",
				},
			},
		},
	}

	projectDBEvent(t, conn, event)

	var count int
	if err := conn.QueryRow(`
		SELECT COUNT(*)
		FROM issues
		WHERE id IN (?, ?)
	`, "FID-journal-epic", "FID-journal-task").Scan(&count); err != nil {
		t.Fatalf("query projected journal issues: %v", err)
	}
	if count != 2 {
		t.Fatalf("projected issue count = %d, want %d", count, 2)
	}
}

type issueRow struct {
	id         string
	kind       string
	parentID   sql.NullString
	sequence   sql.NullInt64
	title      string
	workflowID string
	status     string
}

func openKindTriggerDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testKindTriggerMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return conn
}

func insertIssueRow(t *testing.T, conn *sql.DB, row issueRow) {
	t.Helper()

	if err := execInsertIssue(conn, row); err != nil {
		t.Fatalf("execInsertIssue(%q) error = %v", row.id, err)
	}
}

func insertEpicMetadataRow(t *testing.T, conn *sql.DB, issueID, scheduling string) {
	t.Helper()

	if _, err := conn.Exec(`
		INSERT INTO epic_metadata (issue_id, scheduling)
		VALUES (?, ?)
	`, issueID, scheduling); err != nil {
		t.Fatalf("insert epic_metadata(%q) error = %v", issueID, err)
	}
}

func execInsertIssue(conn *sql.DB, row issueRow) error {
	_, err := conn.Exec(`
		INSERT INTO issues (id, kind, parent_id, sequence, title, workflow_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		row.id,
		row.kind,
		nullStringValue(row.parentID),
		nullInt64Value(row.sequence),
		row.title,
		row.workflowID,
		row.status,
		"2026-06-07T13:59:00Z",
		"2026-06-07T13:59:00Z",
	)
	return err
}

func assertSQLErrorContains(t *testing.T, err error, want string) {
	t.Helper()

	if err == nil {
		t.Fatalf("error = nil, want substring %q", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error = %q, want substring %q", err.Error(), want)
	}
}

func testKindTriggerMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func closeKindTriggerDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func nullStringValue(value sql.NullString) any {
	if !value.Valid {
		return nil
	}
	return value.String
}

func nullInt64Value(value sql.NullInt64) any {
	if !value.Valid {
		return nil
	}
	return value.Int64
}
