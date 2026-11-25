package db

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestReviewRenderMigrationPreservesReceipts(t *testing.T) {
	t.Parallel()
	m, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsFS(t))
	defer closeMigrator(t, m, conn)
	if err := m.Migrate(41); err != nil {
		t.Fatal(err)
	}
	for _, stmt := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,created_at,updated_at) VALUES('FID-a','task','Task','task-default','in_review','2026-11-25T00:00:00Z','2026-11-25T00:00:00Z')`,
		`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at) VALUES('HLD-a','FID-a','builder','/w','agent:builder','2026-11-25T00:00:00Z','2026-11-25T01:00:00Z')`,
		`INSERT INTO policy_snapshots VALUES(1,'built-in defaults',NULL,'digest','{}')`,
		`INSERT INTO handoffs VALUES('HOF-a','FID-a','HLD-a',1,'{}')`,
		`INSERT INTO review_posts VALUES('RND-old','FID-a','HOF-a','pr','old-url',1,'old-digest','2026-11-25T00:01:00Z')`,
	} {
		if _, err := conn.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}
	const insert = `INSERT INTO review_posts VALUES(?,'FID-a','HOF-a','pr',?,?, 'digest','2026-11-25T00:02:00Z')`
	if _, err := conn.Exec(insert, "RND-new", "new-url", 2); err == nil {
		t.Fatal("schema 41 admitted render version 2")
	}
	if err := m.Migrate(42); err != nil {
		t.Fatal(err)
	}
	var row string
	if err := conn.QueryRow(`SELECT json_array(id,task_id,handoff_id,pr,url,render_version,digest,posted_at) FROM review_posts`).Scan(&row); err != nil || row != `["RND-old","FID-a","HOF-a","pr","old-url",1,"old-digest","2026-11-25T00:01:00Z"]` {
		t.Fatalf("migration changed receipt: %s / %v", row, err)
	}
	for _, version := range []int{1, 2, 0, 3, 99} {
		_, err := conn.Exec(insert, "RND-new", "new-url", version)
		if version == 1 || version == 2 {
			if err != nil {
				t.Fatalf("render version %d: %v", version, err)
			}
			if _, err := conn.Exec(insert, "RND-duplicate", "new-url", version); err == nil {
				t.Fatal("migration lost unique comment URL")
			}
			if _, err := conn.Exec(`DELETE FROM review_posts WHERE id='RND-new'`); err != nil {
				t.Fatal(err)
			}
		} else if err == nil || !strings.Contains(err.Error(), "CHECK constraint failed") {
			t.Fatalf("unknown render version %d refusal: %v", version, err)
		}
	}
	assertForeignKeysContain(t, conn, "review_posts", "task_id->issues.id", "handoff_id->handoffs.id")
	if err := m.Migrate(41); err != nil {
		t.Fatal(err)
	}
	var restored string
	if err := conn.QueryRow(`SELECT json_array(id,task_id,handoff_id,pr,url,render_version,digest,posted_at) FROM review_posts`).Scan(&restored); err != nil || restored != row {
		t.Fatalf("downgrade changed receipt: %s / %v", restored, err)
	}
	if _, err := conn.Exec(insert, "RND-new", "new-url", 2); err == nil {
		t.Fatal("downgrade admitted render version 2")
	}
	if err := m.Migrate(42); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(insert, "RND-new", "new-url", 2); err != nil {
		t.Fatal(err)
	}
	if err := m.Migrate(41); err == nil || !strings.Contains(err.Error(), "CHECK constraint failed") {
		t.Fatalf("downgrade with version-2 receipts: %v", err)
	}
	var count int
	if err := conn.QueryRow(`SELECT count(*) FROM review_posts`).Scan(&count); err != nil || count != 2 {
		t.Fatalf("refused downgrade lost receipts: %d / %v", count, err)
	}
}
