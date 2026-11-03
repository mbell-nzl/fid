package main

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestIssueListFiltersAndOrdering(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	seedIssueListCommandFixture(t, conn)

	tests := []struct {
		name    string
		args    []string
		wantIDs []string
	}{
		{
			name:    "kind and status",
			args:    []string{"--json", "issue", "list", "--kind", "task", "--status", "todo"},
			wantIDs: []string{"FID-task01"},
		},
		{
			name:    "status category",
			args:    []string{"--json", "issue", "list", "--status-category", "started"},
			wantIDs: []string{"FID-task02"},
		},
		{
			name: "combined filters",
			args: []string{
				"--json", "issue", "list",
				"--epic", "FID-epic01",
				"--label", "core",
				"--priority", "1",
				"--has-concerns",
			},
			wantIDs: []string{"FID-task01"},
		},
		{
			name:    "deterministic order",
			args:    []string{"--json", "issue", "list", "--kind", "task", "--order", "priority"},
			wantIDs: []string{"FID-task02", "FID-task01", "FID-task04", "FID-task03"},
		},
		{
			name:    "empty result",
			args:    []string{"--json", "issue", "list", "--kind", "task", "--status", "cancelled"},
			wantIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, tt.args...)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			rows := decodeCanonicalResult(t, stdout)["issues"].([]any)
			if gotIDs := commandIssueListIDs(rows); !equalIssueIDs(gotIDs, tt.wantIDs) {
				t.Fatalf("stdout ids = %#v, want %#v", gotIDs, tt.wantIDs)
			}
		})
	}
}

func TestIssueListRejectsRetiredFormatFlag(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	if result := executeCommandForTest(workspace, "issue", "list", "--format", "json"); result.err == nil {
		t.Fatal("--format was accepted")
	}
}

func seedIssueListCommandFixture(t *testing.T, conn *sql.DB) {
	t.Helper()

	statements := []struct {
		query string
		args  []any
	}{
		{
			query: `
				INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-epic01",
				"epic",
				"Phase 4 discovery",
				"todo",
				"epic-default",
				"tl",
				"2026-06-09T09:00:00Z",
				"2026-06-09T09:00:00Z",
			},
		},
		{
			query: `INSERT INTO epic_metadata (issue_id, scheduling) VALUES (?, ?)`,
			args:  []any{"FID-epic01", "sequential"},
		},
		{
			query: `INSERT INTO labels (name) VALUES (?), (?), (?)`,
			args:  []any{"backend", "core", "docs"},
		},
		{
			query: `
				INSERT INTO issues (
					id, kind, parent_id, title, status, workflow_id, priority, sequence, owner_role, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-task01",
				"task",
				"FID-epic01",
				"Inbox task",
				"todo",
				"task-default",
				1,
				2,
				"developer",
				"2026-06-09T09:01:00Z",
				"2026-06-09T09:03:00Z",
			},
		},
		{
			query: `
				INSERT INTO issues (
					id, kind, parent_id, title, status, workflow_id, priority, sequence, owner_role, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-task02",
				"task",
				"FID-epic01",
				"Started task",
				"in_progress",
				"task-default",
				0,
				1,
				"developer",
				"2026-06-09T09:02:00Z",
				"2026-06-09T09:01:00Z",
			},
		},
		{
			query: `
				INSERT INTO issues (
					id, kind, parent_id, title, status, workflow_id, priority, sequence, owner_role, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-task03",
				"task",
				"FID-epic01",
				"Completed task",
				"done",
				"task-default",
				2,
				3,
				"developer",
				"2026-06-09T09:03:00Z",
				"2026-06-09T09:02:00Z",
			},
		},
		{
			query: `
				INSERT INTO issues (
					id, kind, title, status, workflow_id, priority, owner_role, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-task04",
				"task",
				"Backlog task",
				"qa_pass",
				"task-default",
				1,
				"developer",
				"2026-06-09T09:01:00Z",
				"2026-06-09T09:03:00Z",
			},
		},
		{
			query: `
				INSERT INTO issues (
					id, kind, parent_id, title, status, workflow_id, owner_role, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-sub01",
				"subtask",
				"FID-task02",
				"Cancelled subtask",
				"cancelled",
				"subtask-default",
				"developer",
				"2026-06-09T09:04:00Z",
				"2026-06-09T09:04:00Z",
			},
		},
		{
			query: `
				INSERT INTO issue_labels (issue_id, label)
				VALUES (?, ?), (?, ?), (?, ?), (?, ?)
			`,
			args: []any{
				"FID-task01", "backend",
				"FID-task01", "core",
				"FID-task03", "core",
				"FID-sub01", "docs",
			},
		},
		{
			query: `INSERT INTO concerns
				(id, subject_type, subject_id, actor_id, scope_json, severity, summary, reasoning, references_json,
				 status, created_at, created_event_id)
				VALUES (?, 'issue', ?, 'developer', '{"paths":[]}', 'major', ?, 'Needs review', '[]',
				        'open', ?, ?)`,
			args: []any{"CON-open01", "FID-task01", "Needs triage", "2026-06-09T09:05:00Z", "EVT-open01"},
		},
		{
			query: `INSERT INTO concerns
				(id, subject_type, subject_id, actor_id, scope_json, severity, summary, reasoning, references_json,
				 status, disposition_action, disposition_reason, disposed_actor_id, disposed_at,
				 disposed_event_id, created_at, created_event_id)
				VALUES (?, 'issue', ?, 'developer', '{"paths":[]}', 'major', ?, 'Needs review', '[]',
				        'disposed', 'rejected', 'Addressed', 'qa_lead', ?, ?, ?, ?)`,
			args: []any{"CON-addressed01", "FID-task03", "Already addressed", "2026-06-09T09:06:00Z",
				"EVT-dispose01", "2026-06-09T09:06:00Z", "EVT-addressed01"},
		},
	}

	for _, statement := range statements {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed issue list command fixture exec error = %v\nquery=%s", err, statement.query)
		}
	}
}

func commandIssueListIDs(rows []any) []string {
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		item, ok := row.(map[string]any)
		if !ok {
			return nil
		}
		id, _ := item["id"].(string)
		ids = append(ids, id)
	}
	return ids
}

func equalIssueIDs(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
