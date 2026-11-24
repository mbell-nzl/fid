package db

import (
	"path/filepath"
	"testing"
)

func TestMigration0028AddsMergeAndCIStorageToExistingProjection(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsFS(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(27); err != nil {
		t.Fatal(err)
	}
	if err := migrator.Up(); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"merge_observations", "ci_results"} {
		var count int
		if err := conn.QueryRow(`SELECT count(*) FROM sqlite_schema WHERE type='table' AND name=?`, table).Scan(&count); err != nil || count != 1 {
			t.Fatalf("upgraded %s=%d/%v", table, count, err)
		}
	}
}
