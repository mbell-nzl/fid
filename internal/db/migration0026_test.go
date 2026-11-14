package db

import (
	"path/filepath"
	"testing"
)

func TestMigration0026PreservesHistoryAndReviewReservation(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(25); err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,lane,owner_role,created_at,updated_at) VALUES('FID-a','task','Work','task-default','in_progress','standard','developer','2026-11-14T00:00:00Z','2026-11-14T00:00:00Z')`,
		`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at,ended_at,end_reason,kept) VALUES('HLD-1','FID-a','one','/w/one','agent:builder','2026-11-14T07:00:00Z','2026-11-14T08:00:00Z','2026-11-14T07:05:00Z','waiting',1)`,
		`INSERT INTO questions(id,task_id,hold_id,text,asker,asked_at) VALUES('QST-1','FID-a','HLD-1','Include timeouts?','agent:builder','2026-11-14T07:05:00Z')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	if err := migrator.Migrate(26); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := conn.QueryRow(`SELECT count(*) FROM questions q JOIN holds h ON h.id=q.hold_id WHERE h.kept=1 AND h.end_reason='waiting'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("upgraded history=%d/%v", count, err)
	}
	if _, err := conn.Exec(`UPDATE holds SET end_reason='review' WHERE id='HLD-1'`); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at) VALUES('HLD-2','FID-a','two','/w/two','agent:other','2026-11-14T07:10:00Z','2026-11-14T08:10:00Z')`); err == nil {
		t.Fatal("review task accepted a second reservation")
	}
	if err := migrator.Migrate(25); err != nil {
		t.Fatal(err)
	}
	if err := conn.QueryRow(`SELECT count(*) FROM questions q JOIN holds h ON h.id=q.hold_id WHERE h.kept=1 AND h.end_reason='interrupted' AND h.park_reason='In review'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("downgraded reservation=%d/%v", count, err)
	}
}
