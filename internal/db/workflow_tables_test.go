package db

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestMigrateUpSeedsTaskDefaultWorkflowStates(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	want := map[string]string{"todo": "unstarted", "in_progress": "started", "in_review": "started", "done": "completed", "cancelled": "cancelled"}

	got := make(map[string]string, len(want))
	rows, err := conn.Query(`
		SELECT state, category
		FROM workflow_states
		WHERE workflow_id = ?
	`, "task-default")
	if err != nil {
		t.Fatalf("query task-default workflow states: %v", err)
	}
	defer closeRows(t, rows)

	for rows.Next() {
		var state, category string
		if err := rows.Scan(&state, &category); err != nil {
			t.Fatalf("scan workflow state: %v", err)
		}
		got[state] = category
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate workflow states: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("task-default state count = %d, want %d; got %#v", len(got), len(want), got)
	}
	for state, wantCategory := range want {
		if got[state] != wantCategory {
			t.Fatalf("task-default category for %q = %q, want %q", state, got[state], wantCategory)
		}
	}
}

func TestMigration0007CreatesWorkflowSchema(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowTablesDB(t)
	defer closeDB(t, conn)

	tables := []struct {
		name    string
		columns []string
		primary []string
	}{
		{"workflows", []string{"id", "name", "applies_to", "active", "created_at"}, []string{"id"}},
		{"workflow_states", []string{"workflow_id", "state", "category", "ordinal"}, []string{"workflow_id", "state"}},
		{"workflow_transitions", []string{"workflow_id", "from_state", "to_state"}, []string{"workflow_id", "from_state", "to_state"}},
		{"workflow_transition_requires", []string{"workflow_id", "from_state", "to_state", "rule", "params_json"}, []string{}},
	}
	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			assertTableExists(t, conn, table.name)
			assertTableColumns(t, conn, table.name, table.columns)
			assertPrimaryKeyColumns(t, conn, table.name, table.primary)
		})
	}

	assertCompositeTransitionRequirementForeignKey(t, conn)
	assertNoPrimaryKeyOrUniqueConstraint(t, conn, "workflow_transition_requires")
}

func TestMigration0007DownAndUpRestoresWorkflowTables(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() error = %v", err)
	}
	if err := migrator.Migrate(6); err != nil {
		t.Fatalf("Migrate(6) error = %v", err)
	}
	for _, table := range workflowTableNames {
		assertTableDoesNotExist(t, conn, table)
	}

	if err := migrator.Up(); err != nil {
		t.Fatalf("second Up() error = %v", err)
	}
	for _, table := range workflowTableNames {
		assertTableExists(t, conn, table)
	}
}

var workflowTableNames = []string{
	"workflows",
	"workflow_states",
	"workflow_transitions",
	"workflow_transition_requires",
}

func migratedWorkflowTablesDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}
	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return conn
}

func assertPrimaryKeyColumns(t *testing.T, conn *sql.DB, table string, want []string) {
	t.Helper()

	rows, err := conn.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		t.Fatalf("PRAGMA table_info(%q): %v", table, err)
	}
	defer closeRows(t, rows)

	gotByPosition := make(map[int]string)
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, dataType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatalf("scan table info for %q: %v", table, err)
		}
		if primaryKey > 0 {
			gotByPosition[primaryKey] = name
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate table info for %q: %v", table, err)
	}

	got := make([]string, 0, len(gotByPosition))
	for position := 1; position <= len(gotByPosition); position++ {
		got = append(got, gotByPosition[position])
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("primary key columns for %q = %#v, want %#v", table, got, want)
	}
}

func assertCompositeTransitionRequirementForeignKey(t *testing.T, conn *sql.DB) {
	t.Helper()

	rows, err := conn.Query(`PRAGMA foreign_key_list(workflow_transition_requires)`)
	if err != nil {
		t.Fatalf("PRAGMA foreign_key_list(workflow_transition_requires): %v", err)
	}
	defer closeRows(t, rows)

	type foreignKeyColumn struct{ from, to string }
	columns := make(map[int][]foreignKeyColumn)
	tables := make(map[int]string)
	for rows.Next() {
		var id, sequence int
		var table, from, to, onUpdate, onDelete, match string
		if err := rows.Scan(&id, &sequence, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			t.Fatalf("scan workflow requirement foreign key: %v", err)
		}
		tables[id] = table
		columns[id] = append(columns[id], foreignKeyColumn{from: from, to: to})
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate workflow requirement foreign keys: %v", err)
	}

	want := []foreignKeyColumn{{"workflow_id", "workflow_id"}, {"from_state", "from_state"}, {"to_state", "to_state"}}
	for id, got := range columns {
		if tables[id] == "workflow_transitions" && reflect.DeepEqual(got, want) {
			return
		}
	}
	t.Fatalf("workflow_transition_requires foreign keys = %#v to %#v, want composite workflow_transitions key %#v", columns, tables, want)
}

func assertNoPrimaryKeyOrUniqueConstraint(t *testing.T, conn *sql.DB, table string) {
	t.Helper()

	var schemaSQL string
	if err := conn.QueryRow(`SELECT sql FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&schemaSQL); err != nil {
		t.Fatalf("load schema SQL for %q: %v", table, err)
	}
	normalized := strings.ToUpper(strings.Join(strings.Fields(schemaSQL), " "))
	if strings.Contains(normalized, "PRIMARY KEY") || strings.Contains(normalized, "UNIQUE") {
		t.Fatalf("schema for %q unexpectedly has primary/unique constraint: %s", table, normalized)
	}
}
