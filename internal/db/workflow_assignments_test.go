package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
)

func TestMigration0009CreatesWorkflowAssignmentsWithScopedUniqueness(t *testing.T) {
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

	insertAssignment := func(scope, kind, workflowID string) error {
		_, err := conn.Exec(`
			INSERT INTO workflow_assignments (scope, kind, workflow_id)
			VALUES (?, ?, ?)
		`, scope, kind, workflowID)
		return err
	}

	if err := insertAssignment("kind", "task", "task-default"); err != nil {
		t.Fatalf("insert first kind assignment: %v", err)
	}
	if err := insertAssignment("kind", "task", "epic-default"); err == nil {
		t.Fatal("second kind assignment for task succeeded, want uniqueness error")
	}
	if err := insertAssignment("kind", "epic", "epic-default"); err != nil {
		t.Fatalf("insert distinct kind assignment: %v", err)
	}
}

func TestMigration0009WorkflowAssignmentsSchemaAndConstraints(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowAssignmentsDB(t)
	defer closeDB(t, conn)

	assertTableColumns(t, conn, "workflow_assignments", []string{"scope", "kind", "workflow_id", "assigned_at"})
	assertSchemaContains(t, conn, "workflow_assignments",
		"scope       TEXT NOT NULL",
		"workflow_id TEXT NOT NULL REFERENCES workflows(id)",
		"CHECK (scope IN ('kind', 'project'))",
		"CHECK ((scope = 'kind' AND kind IS NOT NULL)",
		"OR (scope = 'project' AND kind IS NULL))",
	)
	assertForeignKeysContain(t, conn, "workflow_assignments", "workflow_id->workflows.id")
	assertPartialIndexes(t, conn, map[string]string{
		"idx_workflow_assignments_kind":    "CREATE UNIQUE INDEX idx_workflow_assignments_kind\n  ON workflow_assignments(kind) WHERE scope = 'kind'",
		"idx_workflow_assignments_project": "CREATE UNIQUE INDEX idx_workflow_assignments_project\n  ON workflow_assignments(scope) WHERE scope = 'project'",
	})

	tests := []struct {
		name    string
		scope   string
		kind    any
		wantErr bool
	}{
		{name: "kind requires kind", scope: "kind", kind: nil, wantErr: true},
		{name: "project rejects kind", scope: "project", kind: "task", wantErr: true},
		{name: "kind with kind accepted", scope: "kind", kind: "subtask"},
		{name: "project without kind accepted", scope: "project", kind: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conn.Exec(`
				INSERT INTO workflow_assignments (scope, kind, workflow_id)
				VALUES (?, ?, 'task-default')
			`, tt.scope, tt.kind)
			if (err != nil) != tt.wantErr {
				t.Fatalf("insert error = %v, want error = %t", err, tt.wantErr)
			}
		})
	}

	if _, err := conn.Exec(`
		INSERT INTO workflow_assignments (scope, kind, workflow_id)
		VALUES ('project', NULL, 'epic-default')
	`); err == nil {
		t.Fatal("second project assignment succeeded, want uniqueness error")
	}
}

func TestMigrations0009And0010RoundTripWithoutResidue(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Migrate(8); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migrate(8) error = %v", err)
	}
	before := snapshotDatabaseState(t, conn)

	if err := migrator.Migrate(10); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migrate(10) error = %v", err)
	}
	afterUp := snapshotDatabaseState(t, conn)
	if !stateContains(afterUp.schema, "workflow_assignments") || !stateContains(afterUp.rows["workflows"], "escalation-default") {
		t.Fatalf("post-up state = %#v, want assignments table and escalation-default", afterUp)
	}

	if err := migrator.Migrate(8); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migrate(8) after up error = %v", err)
	}
	afterDown := snapshotDatabaseState(t, conn)
	if !reflect.DeepEqual(afterDown, before) {
		t.Fatalf("state after down = %#v, want pre-0009 state %#v", afterDown, before)
	}

	if err := migrator.Migrate(10); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("second Migrate(10) error = %v", err)
	}
	if afterSecondUp := snapshotDatabaseState(t, conn); !reflect.DeepEqual(afterSecondUp, afterUp) {
		t.Fatalf("state after second up = %#v, want first post-up state %#v", afterSecondUp, afterUp)
	}
}

func stateContains(values []string, want string) bool {
	for _, value := range values {
		if strings.Contains(value, want) {
			return true
		}
	}
	return false
}

func migratedWorkflowAssignmentsDB(t *testing.T) *sql.DB {
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

func assertPartialIndexes(t *testing.T, conn *sql.DB, want map[string]string) {
	t.Helper()

	rows, err := conn.Query(`
		SELECT name, sql
		FROM sqlite_schema
		WHERE type = 'index' AND tbl_name = 'workflow_assignments'
		ORDER BY name
	`)
	if err != nil {
		t.Fatalf("query workflow assignment indexes: %v", err)
	}
	defer closeRows(t, rows)

	got := make(map[string]string)
	for rows.Next() {
		var name, sqlText string
		if err := rows.Scan(&name, &sqlText); err != nil {
			t.Fatalf("scan workflow assignment index: %v", err)
		}
		got[name] = sqlText
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate workflow assignment indexes: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("workflow assignment indexes = %#v, want %#v", got, want)
	}
}

type databaseState struct {
	schema []string
	rows   map[string][]string
}

func snapshotDatabaseState(t *testing.T, conn *sql.DB) databaseState {
	t.Helper()

	state := databaseState{rows: make(map[string][]string)}
	objects, err := conn.Query(`
		SELECT type || ':' || name || ':' || COALESCE(sql, '')
		FROM sqlite_schema
		WHERE name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'
		ORDER BY type, name
	`)
	if err != nil {
		t.Fatalf("query sqlite_schema: %v", err)
	}
	defer closeRows(t, objects)
	for objects.Next() {
		var object string
		if err := objects.Scan(&object); err != nil {
			t.Fatalf("scan sqlite_schema: %v", err)
		}
		state.schema = append(state.schema, object)
	}
	if err := objects.Err(); err != nil {
		t.Fatalf("iterate sqlite_schema: %v", err)
	}

	tables, err := conn.Query(`
		SELECT name FROM sqlite_schema
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'
		ORDER BY name
	`)
	if err != nil {
		t.Fatalf("query tables: %v", err)
	}
	defer closeRows(t, tables)
	for tables.Next() {
		var table string
		if err := tables.Scan(&table); err != nil {
			t.Fatalf("scan table name: %v", err)
		}
		state.rows[table] = snapshotTableRows(t, conn, table)
	}
	if err := tables.Err(); err != nil {
		t.Fatalf("iterate table names: %v", err)
	}
	return state
}

func snapshotTableRows(t *testing.T, conn *sql.DB, table string) []string {
	t.Helper()

	rows, err := conn.Query(fmt.Sprintf(`SELECT * FROM %q ORDER BY rowid`, table))
	if err != nil {
		t.Fatalf("query rows from %q: %v", table, err)
	}
	defer closeRows(t, rows)

	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("read columns from %q: %v", table, err)
	}
	result := make([]string, 0)
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(values))
		for index := range values {
			pointers[index] = &values[index]
		}
		if err := rows.Scan(pointers...); err != nil {
			t.Fatalf("scan rows from %q: %v", table, err)
		}
		result = append(result, fmt.Sprint(values...))
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate rows from %q: %v", table, err)
	}
	sort.Strings(result)
	return result
}
