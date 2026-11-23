package db

import (
	"path/filepath"
	"testing"
)

func TestReviewClaimsMigrationAndVersions(t *testing.T) {
	t.Parallel()
	if ActiveMigrationVersion != 39 || ProjectionSchemaVersion != 21 {
		t.Fatalf("work decisions require migration 39 / cache 21; got %d / %d", ActiveMigrationVersion, ProjectionSchemaVersion)
	}
	m, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, m, conn)
	if err := m.Migrate(35); err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,created_at,updated_at) VALUES('FID-a','task','Task','task-default','in_review','2026-11-22T00:00:00Z','2026-11-22T00:00:00Z')`,
		`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at) VALUES('HLD-a','FID-a','builder','/w/build','agent:builder','2026-11-22T00:00:00Z','2026-11-22T01:00:00Z')`,
		`INSERT INTO policy_snapshots VALUES(1,'built-in defaults',NULL,'digest','{}')`,
		`INSERT INTO handoffs VALUES('HOF-a','FID-a','HLD-a',1,'{}')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.Migrate(36); err != nil {
		t.Fatal(err)
	}
	insert := `INSERT INTO review_claims VALUES(?, 'FID-a','HOF-a','primary',?,'/w/review','agent:reviewer','claude/opus-5.5','2026-11-22T00:00:00Z','2026-11-22T00:00:00Z','2026-11-22T01:00:00Z',NULL,'{}')`
	if _, err := conn.Exec(insert, "RCL-one", "reviewer-one"); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(insert, "RCL-two", "reviewer-two"); err == nil {
		t.Fatal("two reviewers reserved the same task")
	}
	if err := m.Migrate(35); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := conn.QueryRow(`SELECT count(*) FROM holds`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("migration changed build history: %d/%v", count, err)
	}
}
