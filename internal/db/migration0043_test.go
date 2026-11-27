package db

import (
	"path/filepath"
	"testing"
)

func TestUsageMigration(t *testing.T) {
	t.Parallel()
	m, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsFS(t))
	defer closeMigrator(t, m, conn)
	if err := m.Migrate(42); err != nil {
		t.Fatal(err)
	}
	if err := m.Migrate(43); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := conn.QueryRow(`SELECT count(*) FROM usage_sessions`).Scan(&n); err != nil || n != 0 {
		t.Fatalf("new table: %d %v", n, err)
	}
	assertForeignKeysContain(t, conn, "usage_sessions", "task_id->issues.id")
	if err := m.Migrate(42); err != nil {
		t.Fatal(err)
	}
	if err := conn.QueryRow(`SELECT count(*) FROM usage_sessions`).Scan(&n); err == nil {
		t.Fatal("down left usage table")
	}
	if err := m.Migrate(43); err != nil {
		t.Fatal(err)
	}
}
