package db

import (
	"path/filepath"
	"testing"
)

func TestMigration0023KeepsLifecycleHistoryAndOneOpenQuestion(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(22); err != nil {
		t.Fatal(err)
	}
	statements := []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-a','task','Work','task-default','todo','developer','2026-11-14T00:00:00Z','2026-11-14T00:00:00Z')`,
		`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at,ended_at,end_reason,kept,park_reason) VALUES('HLD-1','FID-a','one','/w/one','agent:builder','2026-11-14T07:00:00Z','2026-11-14T08:00:00Z','2026-11-14T07:05:00Z','interrupted',1,'Usage limit')`,
		`INSERT INTO person_decisions VALUES('PDN-1','takeover','human:mark','task','FID-a','Continue here','FID-a','2026-11-14T07:00:00Z')`,
	}
	for _, s := range statements {
		if _, err := conn.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
	if err := migrator.Migrate(23); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := conn.QueryRow(`SELECT count(*) FROM holds WHERE kept=1 AND park_reason='Usage limit' AND last_heartbeat_at=since`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("upgraded hold: %d, %v", n, err)
	}
	if _, err := conn.Exec(`UPDATE holds SET end_reason='waiting',park_reason=NULL WHERE id='HLD-1'`); err != nil {
		t.Fatal(err)
	}
	insert := `INSERT INTO questions(id,task_id,hold_id,text,asker,asked_at) VALUES(?,'FID-a','HLD-1','Include timeouts?','agent:builder','2026-11-14T07:05:00Z')`
	if _, err := conn.Exec(insert, "QST-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(insert, "QST-2"); err == nil {
		t.Fatal("two open questions accepted")
	}
	if err := migrator.Migrate(22); err != nil {
		t.Fatal(err)
	}
	if err := conn.QueryRow(`SELECT count(*) FROM holds WHERE kept=1 AND end_reason='interrupted' AND park_reason='Include timeouts?'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("downgraded hold: %d, %v", n, err)
	}
	if err := conn.QueryRow(`SELECT count(*) FROM person_decisions WHERE id='PDN-1'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("downgraded decision: %d, %v", n, err)
	}
}
