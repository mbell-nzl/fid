package db

import (
	"path/filepath"
	"testing"
)

func TestReworkMigrationPreservesHoldHistory(t *testing.T) {
	t.Parallel()
	m, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsFS(t))
	defer closeMigrator(t, m, conn)
	if err := m.Migrate(31); err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-a','task','Work','task-default','in_review','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`,
		`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,model,since,expires_at,ended_at,end_reason,kept) VALUES('HLD-1','FID-a','one','/w','agent:builder','codex/gpt-6.1-sol','2026-11-15T00:00:00Z','2026-11-15T01:00:00Z','2026-11-15T00:10:00Z','review',1)`,
	} {
		if _, err := conn.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.Up(); err != nil {
		t.Fatal(err)
	}
	var untouched int
	if err := conn.QueryRow(`SELECT count(*) FROM holds WHERE rework_at IS NULL AND model='codex/gpt-6.1-sol' AND kept=1 AND ended_at='2026-11-15T00:10:00Z' AND last_heartbeat_at=since`).Scan(&untouched); err != nil || untouched != 1 {
		t.Fatalf("migration changed hold history: %d/%v", untouched, err)
	}
	if _, err := conn.Exec(`UPDATE holds SET rework_at='2026-11-15T00:11:00Z'`); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`UPDATE holds SET rework_at='2026-11-15T00:09:00Z'`); err == nil {
		t.Fatal("rework predating hand-off was accepted")
	}
	if err := m.Migrate(31); err != nil {
		t.Fatal(err)
	}
	var model, reason string
	if err := conn.QueryRow(`SELECT model,end_reason FROM holds`).Scan(&model, &reason); err != nil || model != "codex/gpt-6.1-sol" || reason != "review" {
		t.Fatalf("rollback changed original hold: %q/%q/%v", model, reason, err)
	}
}
