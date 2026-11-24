package db

import (
	"path/filepath"
	"testing"
)

func TestReviewMigrationsPreserveExistingHolds(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsFS(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(28); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-a','task','Work','task-default','todo','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO holds(id,task_id,actor,worktree_id,worktree_root,since,expires_at) VALUES('HLD-1','FID-a','agent:legacy','one','/w','2026-11-15T00:00:00Z','2026-11-15T01:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if err := migrator.Up(); err != nil {
		t.Fatal(err)
	}
	var model string
	if err := conn.QueryRow(`SELECT model FROM holds WHERE id='HLD-1'`).Scan(&model); err != nil || model != "" {
		t.Fatalf("legacy model=%q/%v", model, err)
	}
	for _, table := range []string{"review_rounds", "round_findings"} {
		var n int
		if err := conn.QueryRow(`SELECT count(*) FROM ` + table).Scan(&n); err != nil || n != 0 {
			t.Fatalf("table %s=%d/%v", table, n, err)
		}
	}
	if err := migrator.Migrate(28); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := conn.QueryRow(`SELECT count(*) FROM holds WHERE id='HLD-1'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("down lost hold=%d/%v", n, err)
	}
}
