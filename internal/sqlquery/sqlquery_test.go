package sqlquery

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
)

func TestClassifyAllowsReadStatements(t *testing.T) {
	t.Parallel()

	queries := []string{
		"SELECT id FROM issues",
		" explain query plan select id from issues ",
		"/* leading block */\nSELECT 1",
		"-- leading line\nPRAGMA table_info(issues)",
		"WITH seeded AS (SELECT 1 AS n) SELECT n FROM seeded",
		"SELECT id FROM issues;",
	}

	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			t.Parallel()

			if err := Classify(query); err != nil {
				t.Fatalf("Classify(%q) error = %v", query, err)
			}
		})
	}
}

func TestClassifyRejectsMutationsAndMultiStatementInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		query        string
		wantRejected string
	}{
		{name: "insert", query: "INSERT INTO issues(id) VALUES('FID-123')", wantRejected: "INSERT INTO issues(id) VALUES('FID-123')"},
		{name: "update", query: "UPDATE issues SET title = 'changed'", wantRejected: "UPDATE issues SET title = 'changed'"},
		{name: "delete", query: "DELETE FROM issues", wantRejected: "DELETE FROM issues"},
		{name: "drop", query: "DROP TABLE issues", wantRejected: "DROP TABLE issues"},
		{name: "create", query: "CREATE TABLE temp_items(id TEXT)", wantRejected: "CREATE TABLE temp_items(id TEXT)"},
		{name: "replace", query: "REPLACE INTO issues(id) VALUES('FID-123')", wantRejected: "REPLACE INTO issues(id) VALUES('FID-123')"},
		{name: "with delete", query: "WITH doomed AS (SELECT id FROM issues) DELETE FROM issues", wantRejected: "DELETE FROM issues"},
		{name: "multiple statements", query: "SELECT id FROM issues; DROP TABLE issues", wantRejected: "DROP TABLE issues"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Classify(tt.query)
			if err == nil {
				t.Fatalf("Classify(%q) error = nil, want rejection", tt.query)
			}

			rejectedErr, ok := err.(*RejectionError)
			if !ok {
				t.Fatalf("Classify(%q) error = %T, want *RejectionError", tt.query, err)
			}
			if rejectedErr.Statement != tt.wantRejected {
				t.Fatalf("rejected statement = %q, want %q", rejectedErr.Statement, tt.wantRejected)
			}
		})
	}
}

func TestRunReturnsStructuredRows(t *testing.T) {
	t.Parallel()

	workspace := newSQLQueryWorkspace(t)

	result, err := Run(Options{
		WorkspaceRoot: workspace,
		MigrationsDir: testSQLQueryMigrationsDir(t),
		Query:         "SELECT id, active, applies_to FROM workflows ORDER BY id LIMIT 2",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	wantColumns := []string{"id", "active", "applies_to"}
	if !reflect.DeepEqual(result.Columns, wantColumns) {
		t.Fatalf("Run().Columns = %#v, want %#v", result.Columns, wantColumns)
	}

	wantRows := []Row{
		{"id": "epic-default", "active": int64(1), "applies_to": "epic"},
		{"id": "subtask-default", "active": int64(1), "applies_to": "subtask"},
	}
	if !reflect.DeepEqual(result.Rows, wantRows) {
		t.Fatalf("Run().Rows = %#v, want %#v", result.Rows, wantRows)
	}
}

func TestReadOnlyConnectionRejectsMutationWhenClassifierIsBypassed(t *testing.T) {
	t.Parallel()

	workspace := newSQLQueryWorkspace(t)

	conn, err := openReadOnly(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("openReadOnly() error = %v", err)
	}
	defer closeSQLQueryConn(conn)

	if _, err := conn.Exec(`INSERT INTO labels(name) VALUES('observer')`); err == nil {
		t.Fatal("Exec(INSERT) error = nil, want read-only failure")
	} else if !strings.Contains(strings.ToLower(err.Error()), "readonly") {
		t.Fatalf("Exec(INSERT) error = %v, want readonly failure", err)
	}
}

func newSQLQueryWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testSQLQueryMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}
	return workspace
}

func testSQLQueryMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}
