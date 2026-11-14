package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
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

	seedIssueListCommandFixture(t, workspace)

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

func seedIssueListCommandFixture(t *testing.T, workspace string) {
	t.Helper()
	rows := []struct {
		table            string
		columns, primary []string
		values           []any
	}{
		{table: "issues", columns: strings.Fields("id kind title status workflow_id owner_role created_at updated_at"), primary: strings.Fields("id"), values: []any{
			"FID-epic01",
			"epic",
			"Phase 4 discovery",
			"todo",
			"epic-default",
			"tl",
			"2026-06-09T09:00:00Z",
			"2026-06-09T09:00:00Z",
		}},
		{table: "epic_metadata", columns: strings.Fields("issue_id scheduling"), primary: strings.Fields("issue_id"), values: []any{"FID-epic01", "sequential"}},
		{table: "labels", columns: strings.Fields("name"), primary: strings.Fields("name"), values: []any{"backend", "core", "docs"}},
		{table: "issues", columns: strings.Fields("id kind parent_id title status workflow_id priority sequence owner_role created_at updated_at"), primary: strings.Fields("id"), values: []any{
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
		}},
		{table: "issues", columns: strings.Fields("id kind parent_id title status workflow_id priority sequence owner_role created_at updated_at"), primary: strings.Fields("id"), values: []any{
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
		}},
		{table: "issues", columns: strings.Fields("id kind parent_id title status workflow_id priority sequence owner_role created_at updated_at"), primary: strings.Fields("id"), values: []any{
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
		}},
		{table: "issues", columns: strings.Fields("id kind title status workflow_id priority owner_role created_at updated_at"), primary: strings.Fields("id"), values: []any{
			"FID-task04",
			"task",
			"Backlog task",
			"done",
			"task-default",
			1,
			"developer",
			"2026-06-09T09:01:00Z",
			"2026-06-09T09:03:00Z",
		}},
		{table: "issues", columns: strings.Fields("id kind parent_id title status workflow_id owner_role created_at updated_at"), primary: strings.Fields("id"), values: []any{
			"FID-sub01",
			"subtask",
			"FID-task02",
			"Cancelled subtask",
			"cancelled",
			"subtask-default",
			"developer",
			"2026-06-09T09:04:00Z",
			"2026-06-09T09:04:00Z",
		}},
		{table: "issue_labels", columns: strings.Fields("issue_id label"), primary: strings.Fields("issue_id label"), values: []any{
			"FID-task01", "backend",
			"FID-task01", "core",
			"FID-task03", "core",
			"FID-sub01", "docs",
		}},
	}
	effects := []journal.Effect{}
	for _, row := range rows {
		for offset := 0; offset < len(row.values); offset += len(row.columns) {
			effect := journal.Effect{Op: "insert", Table: row.table, PK: map[string]any{}, Values: map[string]any{}}
			for i, column := range row.columns {
				primary := false
				for _, key := range row.primary {
					primary = primary || key == column
				}
				if primary {
					effect.PK[column] = row.values[offset+i]
				} else {
					effect.Values[column] = row.values[offset+i]
				}
			}
			effects = append(effects, effect)
		}
	}
	appendCommandIssueEvent(t, workspace, journal.Event{EventSchema: journal.EventSchemaV1, EventType: "issue.read-fixture", Actor: journal.Actor{Role: "developer"}, Command: journal.Command{Argv: []string{"fid", "test"}, CWD: ".", FIDVersion: "test"}, Effects: effects})
	if _, err := journal.EnsureProjectionFresh(workspace, issueMigrationsDir()); err != nil {
		t.Fatal(err)
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
