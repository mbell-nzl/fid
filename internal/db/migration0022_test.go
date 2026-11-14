package db

import (
	"path/filepath"
	"testing"
)

func TestMigration0022KeepsHistoryAndReservesParkedWork(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(21); err != nil {
		t.Fatal(err)
	}
	statements := []string{
		`INSERT INTO issues (id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES ('FID-a','task','A','task-default','in_progress','developer','2026-11-14T00:00:00Z','2026-11-14T00:00:00Z')`,
		`INSERT INTO issues (id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES ('FID-b','task','B','task-default','todo','developer','2026-11-14T00:00:00Z','2026-11-14T00:00:00Z')`,
		`INSERT INTO holds (id,task_id,worktree_id,worktree_root,actor,since,expires_at,ended_at,end_reason) VALUES ('HLD-1','FID-a','one','/w/one','agent:builder','2026-11-14T07:00:00Z','2026-11-14T08:00:00Z','2026-11-14T08:00:00Z','replaced')`,
		`INSERT INTO holds (id,task_id,worktree_id,worktree_root,actor,since,expires_at,replaces_hold_id) VALUES ('HLD-2','FID-a','one','/w/one','agent:builder','2026-11-14T08:00:00Z','2026-11-14T09:00:00Z','HLD-1')`,
	}
	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	if err := migrator.Migrate(22); err != nil {
		t.Fatal(err)
	}
	var initialized int
	if err := conn.QueryRow(`SELECT count(*) FROM holds WHERE last_heartbeat_at=since`).Scan(&initialized); err != nil || initialized != 2 {
		t.Fatalf("initial heartbeat count = %d, %v", initialized, err)
	}
	if _, err := conn.Exec(`UPDATE holds SET ended_at='2026-11-14T08:05:00Z',end_reason='interrupted',park_reason='Usage limit',kept=1 WHERE id='HLD-2'`); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ name, task, worktree string }{{"same task", "FID-a", "two"}, {"same worktree", "FID-b", "one"}} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := conn.Exec(`INSERT INTO holds (id,task_id,worktree_id,worktree_root,actor,since,expires_at) VALUES ('HLD-3',?,?,'/w/new','agent:builder','2026-11-14T08:06:00Z','2026-11-14T09:06:00Z')`, test.task, test.worktree); err == nil {
				t.Fatal("parked reservation allowed conflicting hold")
			}
		})
	}
	if _, err := conn.Exec(`UPDATE holds SET kept=0 WHERE id='HLD-2'`); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO holds (id,task_id,worktree_id,worktree_root,actor,since,expires_at,replaces_hold_id) VALUES ('HLD-3','FID-a','one','/w/one','human:mark','2026-11-14T08:06:00Z','2026-11-14T09:06:00Z','HLD-2')`); err != nil {
		t.Fatal(err)
	}
	if err := migrator.Migrate(21); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := conn.QueryRow(`SELECT count(*) FROM holds`).Scan(&count); err != nil || count != 3 {
		t.Fatalf("downgrade lost hold history: %d, %v", count, err)
	}
}
