package db

import (
	"path/filepath"
	"testing"
)

func TestShadowJudgeMigration(t *testing.T) {
	t.Parallel()
	m, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsFS(t))
	defer closeMigrator(t, m, conn)
	for _, version := range []uint{43, 44, 43, 44} {
		if err := m.Migrate(version); err != nil {
			t.Fatal(err)
		}
	}
	for _, table := range []string{"merge_holds", "merge_calls"} {
		var count int
		if err := conn.QueryRow("SELECT count(*) FROM " + table).Scan(&count); err != nil || count != 0 {
			t.Fatalf("%s: %d %v", table, count, err)
		}
	}
	assertForeignKeysContain(t, conn, "merge_holds", "task_id->issues.id", "handoff_id->handoffs.id")
}
