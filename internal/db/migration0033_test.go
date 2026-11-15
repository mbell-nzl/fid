package db

import (
	"path/filepath"
	"testing"
)

func TestApprovalMigrationPreservesQuestionsAndFindingDecisions(t *testing.T) {
	t.Parallel()
	m, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, m, conn)
	if err := m.Migrate(32); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-a','task','Task','task-default','in_review','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`,
		`INSERT INTO holds(id,task_id,actor,worktree_id,worktree_root,since,expires_at) VALUES('HLD-1','FID-a','agent:builder','one','/w','2026-11-15T00:00:00Z','2026-11-15T01:00:00Z')`,
		`INSERT INTO policy_snapshots VALUES(1,'built-in defaults',NULL,'digest','{}')`,
		`INSERT INTO handoffs VALUES('HOF-1','FID-a','HLD-1',1,'{}')`,
		`INSERT INTO review_rounds VALUES('RND-1','FID-a','HOF-1','{}')`,
		`INSERT INTO round_findings VALUES('FND-1','RND-1','{}')`,
		`INSERT INTO person_decisions VALUES('PDN-answer','answer','human:mark','task','FID-a','Yes','QST-1','2026-11-15T00:02:00Z')`,
		`INSERT INTO questions VALUES('QST-1','FID-a','HLD-1','Proceed?','agent:builder','2026-11-15T00:01:00Z','2026-11-15T00:02:00Z','answered','PDN-answer')`,
		`INSERT INTO person_decisions VALUES('PDN-close','close-finding','human:mark','task','FID-a','Accepted risk','FND-1','2026-11-15T00:03:00Z')`,
		`INSERT INTO finding_closures VALUES('FND-1',NULL,'PDN-close','closed','Accepted risk')`,
	} {
		if _, err := conn.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
	for _, version := range []uint{33, 32} {
		if err := m.Migrate(version); err != nil {
			t.Fatal(err)
		}
		var answer, closure string
		if err := conn.QueryRow(`SELECT p.reason FROM questions q JOIN person_decisions p ON p.id=q.decision_id`).Scan(&answer); err != nil || answer != "Yes" {
			t.Fatal(answer, err)
		}
		if err := conn.QueryRow(`SELECT p.reason FROM finding_closures c JOIN person_decisions p ON p.id=c.decision_id`).Scan(&closure); err != nil || closure != "Accepted risk" {
			t.Fatal(closure, err)
		}
		rows, err := conn.Query(`PRAGMA foreign_key_check`)
		if err != nil {
			t.Fatal(err)
		}
		if rows.Next() {
			t.Fatal("foreign key violation after approval migration")
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
