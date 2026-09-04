package issues

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
)

func TestListFiltersByKind(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	epic, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "epic",
		Title:         "Phase 4 discovery",
		Actor:         testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(),
		Scheduling: "sequential",
	})
	if err != nil {
		t.Fatalf("Create(epic) error = %v", err)
	}

	task, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Implement issue list",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		ParentID: stringPointer(epic.ID),
		Sequence: intPointer(1),
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	if _, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "subtask",
		Title:         "Add filter coverage",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		ParentID: stringPointer(task.ID),
	}); err != nil {
		t.Fatalf("Create(subtask) error = %v", err)
	}

	got, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("List() len = %d, want 1", len(got))
	}
	if got[0].ID != task.ID {
		t.Fatalf("List() id = %q, want %q", got[0].ID, task.ID)
	}
	if got[0].Kind != "task" {
		t.Fatalf("List() kind = %q, want task", got[0].Kind)
	}
}

func TestListFiltersAndOrdering(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	seedIssueListFixture(t, conn)

	t.Run("filters", func(t *testing.T) {
		tests := []struct {
			name    string
			options ListOptions
			wantIDs []string
		}{
			{
				name: "status",
				options: ListOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Status:        "todo",
				},
				wantIDs: []string{"FID-epic01", "FID-task01"},
			},
			{
				name: "status category",
				options: ListOptions{
					WorkspaceRoot:  workspace,
					MigrationsDir:  testIssueMigrationsDir(t),
					StatusCategory: "started",
				},
				wantIDs: []string{"FID-task02"},
			},
			{
				name: "epic",
				options: ListOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					EpicID:        "FID-epic01",
				},
				wantIDs: []string{"FID-task01", "FID-task02", "FID-task03"},
			},
			{
				name: "label",
				options: ListOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Label:         "core",
				},
				wantIDs: []string{"FID-task01", "FID-task03"},
			},
			{
				name: "priority",
				options: ListOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Priorities:    []int{1},
				},
				wantIDs: []string{"FID-task01", "FID-task04"},
			},
			{
				name: "open concerns",
				options: ListOptions{
					WorkspaceRoot:   workspace,
					MigrationsDir:   testIssueMigrationsDir(t),
					HasOpenConcerns: boolPointer(true),
				},
				wantIDs: []string{"FID-task01"},
			},
			{
				name: "empty result",
				options: ListOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Label:         "missing",
				},
				wantIDs: []string{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := List(tt.options)
				if err != nil {
					t.Fatalf("List() error = %v", err)
				}
				if gotIDs := issueListIDs(got); !equalStrings(gotIDs, tt.wantIDs) {
					t.Fatalf("List() ids = %#v, want %#v", gotIDs, tt.wantIDs)
				}
			})
		}
	})

	t.Run("ordering", func(t *testing.T) {
		tests := []struct {
			name    string
			order   string
			wantIDs []string
		}{
			{
				name:    "sequence",
				order:   "sequence",
				wantIDs: []string{"FID-task02", "FID-task01", "FID-task03", "FID-task04"},
			},
			{
				name:    "priority",
				order:   "priority",
				wantIDs: []string{"FID-task02", "FID-task01", "FID-task04", "FID-task03"},
			},
			{
				name:    "created",
				order:   "created",
				wantIDs: []string{"FID-task01", "FID-task04", "FID-task02", "FID-task03"},
			},
			{
				name:    "updated",
				order:   "updated",
				wantIDs: []string{"FID-task02", "FID-task03", "FID-task01", "FID-task04"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := List(ListOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Kind:          "task",
					Order:         tt.order,
				})
				if err != nil {
					t.Fatalf("List() error = %v", err)
				}
				if gotIDs := issueListIDs(got); !equalStrings(gotIDs, tt.wantIDs) {
					t.Fatalf("List() ids = %#v, want %#v", gotIDs, tt.wantIDs)
				}
			})
		}
	})
}

func TestParsePriorityFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr string
		wantNil bool
	}{
		{
			name:    "blank input returns nil",
			input:   "   ",
			wantNil: true,
		},
		{
			name:  "comma separated priorities are parsed",
			input: "1, 3,5",
			want:  []int{1, 3, 5},
		},
		{
			name:    "embedded empty value is rejected",
			input:   "1, ,5",
			wantErr: "priority filter contains an empty value",
		},
		{
			name:    "non integer value is rejected",
			input:   "1,high,5",
			wantErr: `priority filter "high" is not an integer`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePriorityFilter(tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("ParsePriorityFilter() error = nil, want %q", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("ParsePriorityFilter() error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePriorityFilter() error = %v", err)
			}
			if tt.wantNil {
				if got != nil {
					t.Fatalf("ParsePriorityFilter() = %#v, want nil", got)
				}
				return
			}
			if !equalInts(got, tt.want) {
				t.Fatalf("ParsePriorityFilter() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestIssueListOrderClause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		order      string
		wantClause string
		wantErr    string
	}{
		{
			name:       "blank order defaults to created",
			order:      "   ",
			wantClause: "ie.created_at ASC, ie.id ASC",
		},
		{
			name:       "updated order",
			order:      "updated",
			wantClause: "ie.updated_at ASC, ie.id ASC",
		},
		{
			name:       "priority order",
			order:      "priority",
			wantClause: "ie.priority IS NULL ASC, ie.priority ASC, ie.created_at ASC, ie.id ASC",
		},
		{
			name:       "sequence order",
			order:      "sequence",
			wantClause: "ie.sequence IS NULL ASC, ie.sequence ASC, ie.created_at ASC, ie.id ASC",
		},
		{
			name:    "unsupported order is rejected with supported list",
			order:   "rank",
			wantErr: `unsupported issue list order "rank"; supported orders: sequence, priority, created, updated`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := issueListOrderClause(tt.order)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("issueListOrderClause() error = nil, want %q", tt.wantErr)
				}
				var unsupportedErr *UnsupportedListOrderError
				if !errors.As(err, &unsupportedErr) {
					t.Fatalf("issueListOrderClause() error = %T, want *UnsupportedListOrderError", err)
				}
				if unsupportedErr.Order != "rank" {
					t.Fatalf("UnsupportedListOrderError.Order = %q, want rank", unsupportedErr.Order)
				}
				if !equalStrings(unsupportedErr.SupportedOrder, []string{"sequence", "priority", "created", "updated"}) {
					t.Fatalf("UnsupportedListOrderError.SupportedOrder = %#v", unsupportedErr.SupportedOrder)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("issueListOrderClause() error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("issueListOrderClause() error = %v", err)
			}
			if got != tt.wantClause {
				t.Fatalf("issueListOrderClause() = %q, want %q", got, tt.wantClause)
			}
		})
	}
}

func seedIssueListFixture(t *testing.T, conn *sql.DB) {
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
			query: `
				INSERT INTO task_completions (
					id, issue_id, attempt_number, developer_role, model, status, gate_results, human_intervention, completed_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"TC-open01",
				"FID-task01",
				1,
				"developer",
				"gpt-5",
				"dev_complete",
				`["go test ./..."]`,
				false,
				"2026-06-09T09:05:00Z",
				"TC-addressed01",
				"FID-task03",
				1,
				"developer",
				"gpt-5",
				"done",
				`["go test ./..."]`,
				false,
				"2026-06-09T09:06:00Z",
			},
		},
		{
			query: `
				INSERT INTO concerns (
					id, issue_id, raised_in_completion_id, raised_by_role, summary, status, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"CON-open01",
				"FID-task01",
				"TC-open01",
				"developer",
				"Needs triage",
				"open",
				"2026-06-09T09:05:00Z",
				"2026-06-09T09:05:00Z",
				"CON-addressed01",
				"FID-task03",
				"TC-addressed01",
				"developer",
				"Already addressed",
				"addressed",
				"2026-06-09T09:06:00Z",
				"2026-06-09T09:06:00Z",
			},
		},
	}

	for _, statement := range statements {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed issue list fixture exec error = %v\nquery=%s", err, statement.query)
		}
	}
}

func equalInts(got, want []int) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func issueListIDs(items []ListItem) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func boolPointer(value bool) *bool {
	return &value
}
