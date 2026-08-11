package db

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
)

func TestIssuesTableExcludesRetiredPhase20Fields(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrationsDir := testMigrationsDir(t)
	migrator, conn := newTestMigrator(t, dbPath, migrationsDir)
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Migrate(11); err != nil {
		t.Fatalf("Migrate(11) error = %v", err)
	}

	assertIssueColumnsAbsent(t, conn, "issues", "in_inbox", "est_tokens", "due_date")
	assertIssueColumnsAbsent(t, conn, "issues_enriched", "in_inbox", "est_tokens", "due_date")

	var issuesSQL string
	if err := conn.QueryRow(`SELECT sql FROM sqlite_schema WHERE type = 'table' AND name = 'issues'`).Scan(&issuesSQL); err != nil {
		t.Fatalf("read issues schema: %v", err)
	}
	if strings.Contains(issuesSQL, "est_tokens") {
		t.Fatalf("issues schema unexpectedly retains est_tokens CHECK: %s", issuesSQL)
	}

	if err := migrator.Migrate(10); err != nil {
		t.Fatalf("Migrate(10) error = %v", err)
	}
	assertIssueColumnsPresent(t, conn, "issues", "in_inbox", "est_tokens", "due_date")

	if _, err := conn.Exec(`
		INSERT INTO issues (id, kind, title, workflow_id, status, est_tokens)
		VALUES ('FID-retired-check', 'epic', 'retired check', 'epic-default', 'todo', 1)
	`); err == nil {
		t.Fatal("insert epic with est_tokens error = nil, want restored CHECK rejection")
	}
}

func assertIssueColumnsAbsent(t *testing.T, conn *sql.DB, table string, names ...string) {
	t.Helper()
	columns := issueColumns(t, conn, table)
	for _, name := range names {
		if columns[name] {
			t.Fatalf("%s unexpectedly contains retired column %q", table, name)
		}
	}
}

func assertIssueColumnsPresent(t *testing.T, conn *sql.DB, table string, names ...string) {
	t.Helper()
	columns := issueColumns(t, conn, table)
	for _, name := range names {
		if !columns[name] {
			t.Fatalf("%s does not contain restored column %q", table, name)
		}
	}
}

func issueColumns(t *testing.T, conn *sql.DB, table string) map[string]bool {
	t.Helper()
	rows, err := conn.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		t.Fatalf("read %s columns: %v", table, err)
	}
	defer closeRows(t, rows)

	columns := make(map[string]bool)
	for rows.Next() {
		var (
			columnID int
			name     string
			typeName string
			notNull  int
			defaultV any
			primary  int
		)
		if err := rows.Scan(&columnID, &name, &typeName, &notNull, &defaultV, &primary); err != nil {
			t.Fatalf("scan %s column: %v", table, err)
		}
		columns[name] = true
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate %s columns: %v", table, err)
	}

	return columns
}
