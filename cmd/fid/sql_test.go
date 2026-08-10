package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	internalsqlquery "github.com/mbell-nzl/fid/internal/sqlquery"
)

func TestSQLCommandAcceptsReadStatementsAndFormats(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	jsonStdout, jsonStderr := runFIDCommand(
		t,
		workspace,
		"sql", "SELECT name, active FROM roles ORDER BY name LIMIT 2",
		"--format", "json",
	)
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	wantRows := []internalsqlquery.Row{
		{"name": "developer", "active": int64(1)},
		{"name": "human", "active": int64(1)},
	}
	want := expectedCommandStructure(t, wantRows)
	if got := decodeCommandJSONStructure(t, jsonStdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("sql json = %#v, want %#v", got, want)
	}

	for _, formatName := range []string{"toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(
				t,
				workspace,
				"sql", "SELECT name, active FROM roles ORDER BY name LIMIT 2",
				"--format", formatName,
			)
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

	markdownStdout, markdownStderr := runFIDCommand(
		t,
		workspace,
		"sql", "SELECT name, active FROM roles ORDER BY name LIMIT 2",
		"--format", "markdown",
	)
	if markdownStderr != "" {
		t.Fatalf("markdown stderr = %q, want empty", markdownStderr)
	}
	wantMarkdown := "" +
		"# SQL Query Results\n\n" +
		"| name | active |\n" +
		"| --- | --- |\n" +
		"| developer | 1 |\n" +
		"| human | 1 |\n"
	if markdownStdout != wantMarkdown {
		t.Fatalf("markdown stdout = %q, want %q", markdownStdout, wantMarkdown)
	}

	for _, query := range []string{
		"EXPLAIN QUERY PLAN SELECT name FROM roles ORDER BY name",
		"PRAGMA table_info(roles)",
		"WITH selected AS (SELECT name FROM roles WHERE name = 'developer') SELECT name FROM selected",
	} {
		t.Run(query, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "sql", query, "--format", "json")
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
