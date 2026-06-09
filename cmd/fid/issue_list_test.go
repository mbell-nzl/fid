package main

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestIssueListFiltersAndOrdering(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
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
			args:    []string{"issue", "list", "--kind", "task", "--status", "todo", "--format", "json"},
			wantIDs: []string{"FID-task01"},
		},
		{
			name:    "status category",
			args:    []string{"issue", "list", "--status-category", "started", "--format", "json"},
			wantIDs: []string{"FID-task02"},
		},
		{
			name: "combined filters",
			args: []string{
				"issue", "list",
				"--epic", "FID-epic01",
				"--label", "core",
				"--priority", "1",
				"--inbox",
				"--has-concerns",
				"--format", "json",
			},
			wantIDs: []string{"FID-task01"},
		},
		{
			name:    "deterministic order",
			args:    []string{"issue", "list", "--kind", "task", "--order", "priority", "--format", "json"},
			wantIDs: []string{"FID-task02", "FID-task01", "FID-task04", "FID-task03"},
		},
		{
			name:    "empty result",
			args:    []string{"issue", "list", "--kind", "task", "--status", "cancelled", "--format", "json"},
			wantIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, tt.args...)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			rows := decodeCommandJSONStructure(t, stdout).([]any)
			if gotIDs := commandIssueListIDs(rows); !equalIssueIDs(gotIDs, tt.wantIDs) {
				t.Fatalf("stdout ids = %#v, want %#v", gotIDs, tt.wantIDs)
			}
		})
	}
}

func TestIssueListDispatchesFormatsAndRejectsUnknownFormats(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	seedIssueListCommandFixture(t, conn)

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "issue", "list", "--kind", "task", "--order", "sequence", "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	want := decodeCommandJSONStructure(t, jsonStdout)

	for _, formatName := range []string{"toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "issue", "list", "--kind", "task", "--order", "sequence", "--format", formatName)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			var got any
			switch formatName {
			case "toon":
				got = decodeCommandTOONStructure(t, stdout)
			case "yaml":
				got = decodeCommandYAMLStructure(t, stdout)
			}

			if !commandStructuresEqual(want, got) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, got, want)
			}
		})
	}

	t.Run("markdown", func(t *testing.T) {
		stdout, stderr := runFIDCommand(t, workspace, "issue", "list", "--kind", "task", "--order", "sequence", "--format", "markdown")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}

		want := "" +
			"# Issues\n\n" +
			"- FID-task02 | task | in_progress | started | epic FID-epic01 | P0 | concerns: 0\n" +
			"  Started task\n" +
			"- FID-task01 | task | todo | unstarted | epic FID-epic01 | P1 | inbox | concerns: 1\n" +
			"  Inbox task\n" +
			"- FID-task03 | task | done | completed | epic FID-epic01 | P2 | concerns: 0\n" +
			"  Completed task\n" +
			"- FID-task04 | task | qa_pass | completed | no epic | P1 | concerns: 0\n" +
			"  Backlog task\n"
		if stdout != want {
			t.Fatalf("markdown stdout = %q, want %q", stdout, want)
		}
	})

	t.Run("xml", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "list", "--format", "xml")
		if err == nil {
			t.Fatal("Execute() error = nil, want unsupported format rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `unsupported format "xml"; known formats: json, toon, markdown, yaml`) {
			t.Fatalf("stderr = %q, want actionable unknown format error", stderr)
		}
	})
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
					id, kind, parent_id, title, status, workflow_id, priority, sequence, in_inbox, owner_role, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
				true,
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
