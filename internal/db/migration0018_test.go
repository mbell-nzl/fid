package db

import (
	"path/filepath"
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
