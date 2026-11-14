package db

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"path/filepath"
	"strings"
	"testing"
)

func TestHistoricalMigration0008PreservesClaimReferenceTriggers(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, path, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(7); err != nil {
		t.Fatal(err)
	}
	seedMigration0008Fixture(t, conn)
	if err := migrator.Migrate(8); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Exec(`INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at) VALUES (?, ?, ?, ?)`, "issue", "FID-task-0001", "developer", "2026-07-19T02:04:00Z"); err != nil {
		t.Fatalf("insert valid claim: %v", err)
	}
	if _, err := conn.Exec(`UPDATE claims SET ref_id = ? WHERE ref_type = ? AND ref_id = ?`, "FID-missing", "issue", "FID-task-0001"); err == nil || !strings.Contains(err.Error(), "claims: ref_type/ref_id must resolve to an existing row") {
		t.Fatalf("update claim dangling ref error = %v, want ref-integrity rejection", err)
	}
}
