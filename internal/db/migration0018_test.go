package db

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
)

func TestLatestSchemaHasNoV1JudgmentTables(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(ActiveMigrationVersion); err != nil {
		t.Fatalf("migrate latest: %v", err)
	}
	for _, name := range []string{"task_completions", "ac_verifications", "requirements", "decisions", "judgment_concerns"} {
		if _, exists := schemaObjectSQL(t, conn, "table", name); exists {
			t.Errorf("removed V1 table %q survives latest schema", name)
		}
	}
	if _, exists := schemaObjectSQL(t, conn, "table", "concerns"); !exists {
		t.Error("renamed generic concerns table is missing")
	}
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('concerns') WHERE name IN ('subject_type','subject_id','status')`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("generic concerns columns present = %d, want 3", count)
	}
}

func TestMigration0018DownRestoresEmptyV1Tables(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(18); err != nil {
		t.Fatalf("migrate 18: %v", err)
	}
	if err := migrator.Migrate(17); err != nil {
		t.Fatalf("migrate down to 17: %v", err)
	}
	for _, name := range []string{"task_completions", "ac_verifications", "requirements", "decisions", "judgment_concerns"} {
		if _, exists := schemaObjectSQL(t, conn, "table", name); !exists {
			t.Errorf("restored V1 table %q is missing", name)
		}
	}
}

func TestMigration0018DownMatchesFreshSchema17(t *testing.T) {
	t.Parallel()
	dir := testMigrationsDir(t)
	fresh, freshDB := newTestMigrator(t, filepath.Join(t.TempDir(), "fresh.sqlite"), dir)
	defer closeMigrator(t, fresh, freshDB)
	if err := fresh.Migrate(17); err != nil {
		t.Fatal(err)
	}
	round, roundDB := newTestMigrator(t, filepath.Join(t.TempDir(), "round.sqlite"), dir)
	defer closeMigrator(t, round, roundDB)
	if err := round.Migrate(18); err != nil {
		t.Fatal(err)
	}
	if err := round.Migrate(17); err != nil {
		t.Fatal(err)
	}
	want := migration0018SchemaObjects(t, freshDB)
	got := migration0018SchemaObjects(t, roundDB)
	if len(got) != len(want) {
		t.Errorf("schema object count = %d, want %d", len(got), len(want))
	}
	for key, expected := range want {
		actual, ok := got[key]
		if !ok {
			t.Errorf("missing schema object %q", key)
			continue
		}
		// SQLite quotes names touched by ALTER TABLE RENAME.
		normalize := func(s string) string {
			return strings.NewReplacer(`"judgment_concerns"`, "judgment_concerns", `"agent_runs"`, "agent_runs").Replace(s)
		}
		if normalize(actual) != normalize(expected) {
			t.Errorf("schema object %q differs:\ngot  %s\nwant %s", key, actual, expected)
		}
	}
}

func migration0018SchemaObjects(t *testing.T, conn *sql.DB) map[string]string {
	t.Helper()
	rows, err := conn.Query(`SELECT type, name, tbl_name, sql FROM sqlite_master WHERE name NOT LIKE 'sqlite_%' AND name != 'schema_migrations' ORDER BY type, name`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	objects := map[string]string{}
	for rows.Next() {
		var kind, name, table, definition string
		if err := rows.Scan(&kind, &name, &table, &definition); err != nil {
			t.Fatal(err)
		}
		objects[kind+":"+name+":"+table] = definition
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return objects
}

func TestMigration0018DownPreservesAgentRunAndFullRollback(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, path, testMigrationsDir(t))
	if err := migrator.Migrate(18); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO agent_runs(id,role,model) VALUES('RUN-18','developer','test')`); err != nil {
		t.Fatal(err)
	}
	if err := migrator.Migrate(17); err != nil {
		t.Fatal(err)
	}
	var role, model string
	var completion sql.NullString
	if err := conn.QueryRow(`SELECT role,model,completion_id FROM agent_runs WHERE id='RUN-18'`).Scan(&role, &model, &completion); err != nil {
		t.Fatal(err)
	}
	if role != "developer" || model != "test" || completion.Valid {
		t.Errorf("restored run = %q %q %v", role, model, completion)
	}
	closeMigrator(t, migrator, conn)
	if err := MigrateDown(path, testMigrationsDir(t)); err != nil {
		t.Fatalf("full rollback: %v", err)
	}
}

func TestMigration0018RoundTripIntegrity(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	for _, version := range []uint{17, 18, 17, 18} {
		if err := migrator.Migrate(version); err != nil {
			t.Fatalf("migrate %d: %v", version, err)
		}
	}
	assertForeignKeyCheckClean(t, conn)
	assertIntegrityCheckClean(t, conn)
}
