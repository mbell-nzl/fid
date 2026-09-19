package db

import (
	"path/filepath"
	"testing"
)

func TestMigration0017CreatesExactlyTheAcceptedObjects(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Migrate(17); err != nil {
		t.Fatalf("Migrate(17) error = %v", err)
	}

	for _, table := range []string{
		"review_records",
		"review_findings",
		"finding_dispositions",
		"judgment_concerns",
		"deviations",
		"decision_refs",
		"human_outcomes",
		"envelope_protected_decisions",
		"requirement_refs",
		"requirement_issue_links",
	} {
		assertTableExists(t, conn, table)
	}
}
