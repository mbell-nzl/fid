package db

import (
	"path/filepath"
	"testing"
)

func TestClosureMigrationPreservesAnsweredQuestion(t *testing.T) {
	t.Parallel()
	m, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsFS(t))
	defer closeMigrator(t, m, conn)
	if err := m.Migrate(30); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-a','task','Task','task-default','todo','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`,
		`INSERT INTO holds(id,task_id,actor,worktree_id,worktree_root,since,expires_at) VALUES('HLD-1','FID-a','agent:builder','one','/w','2026-11-15T00:00:00Z','2026-11-15T01:00:00Z')`,
		`INSERT INTO person_decisions VALUES('PDN-answer','answer','human:mark','task','FID-a','Yes','QST-1','2026-11-15T00:02:00Z')`,
		`INSERT INTO questions VALUES('QST-1','FID-a','HLD-1','Proceed?','agent:builder','2026-11-15T00:01:00Z','2026-11-15T00:02:00Z','answered','PDN-answer')`,
	} {
		if _, err := conn.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.Up(); err != nil {
		t.Fatal(err)
	}
	for _, version := range []uint{31, 30} {
		if version == 30 {
			if err := m.Migrate(version); err != nil {
				t.Fatal(err)
			}
		}
		var text, reason string
		if err := conn.QueryRow(`SELECT q.text,p.reason FROM questions q JOIN person_decisions p ON p.id=q.decision_id`).Scan(&text, &reason); err != nil || text != "Proceed?" || reason != "Yes" {
			t.Fatalf("v%d: %q/%q/%v", version, text, reason, err)
		}
		rows, err := conn.Query(`PRAGMA foreign_key_check`)
		if err != nil {
			t.Fatal(err)
		}
		if rows.Next() {
			t.Fatal("foreign key violation after migration")
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
