package db_test

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"database/sql"
	"testing"
)

func TestHistoricalTaskCompletionKindGuards(t *testing.T) {
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
			name: "rejects task completion on epic",
			query: `
				INSERT INTO task_completions (id, issue_id, attempt_number, status, completed_at)
				VALUES (?, ?, ?, ?, ?)
			`,
			args:    []any{"TC-invalid-epic", "FID-guard-epic", 1, "dev_complete", "2026-06-07T14:00:00Z"},
			wantErr: "task_completions: issue_id must reference a task or subtask",
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

	taskCompletionCases := []struct {
		name    string
		id      string
		issueID string
	}{
		{name: "task", id: "TC-valid-task", issueID: "FID-guard-task"},
		{name: "subtask", id: "TC-valid-subtask", issueID: "FID-guard-subtask"},
	}
	for _, tt := range taskCompletionCases {
		t.Run("accepts task completion on "+tt.name, func(t *testing.T) {
			if _, err := conn.Exec(`
				INSERT INTO task_completions (id, issue_id, attempt_number, status, completed_at)
				VALUES (?, ?, ?, ?, ?)
			`, tt.id, tt.issueID, 1, "dev_complete", "2026-06-07T14:05:00Z"); err != nil {
				t.Fatalf("insert task_completions(%q) error = %v", tt.id, err)
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
