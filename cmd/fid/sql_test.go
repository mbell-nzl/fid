package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
	internalsqlquery "github.com/mbell-nzl/fid/internal/sqlquery"
)

func TestSQLCommandAcceptsReadStatementsAsDirectLegacyJSON(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"sql", "SELECT id, active FROM workflows ORDER BY id LIMIT 2",
	)
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	wantRows := []internalsqlquery.Row{
		{"id": "epic-default", "active": int64(1)},
		{"id": "subtask-default", "active": int64(1)},
	}
	want := expectedCommandStructure(t, wantRows)
	if got := decodeCommandJSONStructure(t, stdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("sql json = %#v, want %#v", got, want)
	}

	for _, query := range []string{
		"EXPLAIN QUERY PLAN SELECT id FROM workflows ORDER BY id",
		"PRAGMA table_info(workflows)",
		"WITH selected AS (SELECT id FROM workflows WHERE id = 'task-default') SELECT id FROM selected",
	} {
		t.Run(query, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "sql", query)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}
			rows := decodeCommandJSONStructure(t, stdout).([]any)
			if len(rows) == 0 {
				t.Fatalf("stdout rows = %d, want at least one row", len(rows))
			}
		})
	}
}

func TestSQLCommandNamesAStoppedJournalByAbsolutePath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		fixture   string
		firstLine string
	}{
		{"corrupt-tail.jsonl", "error: malformed journal entry at %s:2\n"},
		{"conflicting-precondition.jsonl", "error: replay precondition conflict at %s:2\n"},
	}
	for _, tc := range tests {
		t.Run(tc.fixture, func(t *testing.T) {
			t.Parallel()
			workspace := gittest.NewWorkspace(t)
			runInitCommand(t, workspace, "init")
			// A later date than init's journal, so the projection catches up
			// through it.
			writeCommandFixture(t, workspace, tc.fixture, "9999-12-01.jsonl")

			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "sql", "SELECT id FROM issues")
			if err == nil {
				t.Fatal("Execute() error = nil, want stopped journal")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			want := fmt.Sprintf(tc.firstLine, rebuildJournalFile(t, workspace, "9999-12-01.jsonl"))
			if got := strings.SplitAfter(stderr, "\n")[0]; got != want {
				t.Fatalf("stderr first line = %q, want %q\nstderr=%q", got, want, stderr)
			}
		})
	}
}

func TestSQLCommandRejectsMutationsWithDocumentedError(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	tests := []struct {
		name         string
		query        string
		wantRejected string
	}{
		{name: "insert", query: "INSERT INTO issues(id) VALUES('FID-123')", wantRejected: "INSERT INTO issues(id) VALUES('FID-123')"},
		{name: "update", query: "UPDATE issues SET title = 'changed'", wantRejected: "UPDATE issues SET title = 'changed'"},
		{name: "delete", query: "DELETE FROM issues", wantRejected: "DELETE FROM issues"},
		{name: "drop", query: "DROP TABLE issues", wantRejected: "DROP TABLE issues"},
		{name: "with delete", query: "WITH doomed AS (SELECT id FROM issues) DELETE FROM issues", wantRejected: "DELETE FROM issues"},
		{name: "multiple statements", query: "SELECT id FROM issues; DROP TABLE issues", wantRejected: "DROP TABLE issues"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "sql", tt.query)
			if err == nil {
				t.Fatalf("Execute() error = nil, want rejection for %q", tt.query)
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}

			wantMessage := fmt.Sprintf(
				"error: fid sql accepts only SELECT, EXPLAIN, PRAGMA query, and WITH statements\n  rejected: %s\n",
				tt.wantRejected,
			)
			if stderr != wantMessage {
				t.Fatalf("stderr = %q, want %q", stderr, wantMessage)
			}
		})
	}
}

func TestSQLCommandRequiresSingleQueryArgument(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "sql")
	if err == nil {
		t.Fatal("Execute() error = nil, want cobra arg validation error")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if !strings.Contains(err.Error(), "accepts 1 arg(s), received 0") {
		t.Fatalf("error = %q, want cobra arg validation error", err.Error())
	}
}
