package db

import (
	"path/filepath"
	"strings"
	"testing"
)

// Migration 20 adds the holds table (RESET 1.4a). The database backs up
// fid next's lock: at most one open hold per task and per worktree. A
// lapse is computed from the clock, so only an ended hold stops counting.
func TestMigration0020HoldsAllowOneOpenHoldPerTaskAndWorktree(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(19); err != nil {
		t.Fatal(err)
	}
	before := migration0018SchemaObjects(t, conn)
	if err := migrator.Migrate(20); err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		`INSERT INTO issues (id, kind, title, workflow_id, status, owner_role, created_at, updated_at) VALUES ('FID-a', 'task', 'A', 'task-default', 'todo', 'developer', '2026-11-10T00:00:00Z', '2026-11-10T00:00:00Z')`,
		`INSERT INTO issues (id, kind, title, workflow_id, status, owner_role, created_at, updated_at) VALUES ('FID-b', 'task', 'B', 'task-default', 'todo', 'developer', '2026-11-10T00:00:00Z', '2026-11-10T00:00:00Z')`,
		`INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at) VALUES ('HLD-1', 'FID-a', 'sha256:one', '/w/one', 'agent:one', '2026-11-10T09:00:00Z', '2026-11-10T10:00:00Z')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("seed %q: %v", statement, err)
		}
	}
	for _, test := range []struct {
		name, statement, refusal string
	}{
		{"a second open hold on the task", `INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at) VALUES ('HLD-2', 'FID-a', 'sha256:two', '/w/two', 'agent:two', '2026-11-10T09:00:00Z', '2026-11-10T10:00:00Z')`, "UNIQUE"},
		{"a second open hold in the worktree", `INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at) VALUES ('HLD-2', 'FID-b', 'sha256:one', '/w/one', 'agent:one', '2026-11-10T09:00:00Z', '2026-11-10T10:00:00Z')`, "UNIQUE"},
		{"a hold on a missing task", `INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at) VALUES ('HLD-2', 'FID-z', 'sha256:two', '/w/two', 'agent:two', '2026-11-10T09:00:00Z', '2026-11-10T10:00:00Z')`, "FOREIGN KEY"},
		{"an end without a reason", `UPDATE holds SET ended_at = '2026-11-10T09:30:00Z' WHERE id = 'HLD-1'`, "CHECK"},
		{"an unknown end reason", `UPDATE holds SET ended_at = '2026-11-10T09:30:00Z', end_reason = 'parked' WHERE id = 'HLD-1'`, "CHECK"},
		{"a hold that expires before it starts", `INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at) VALUES ('HLD-2', 'FID-b', 'sha256:two', '/w/two', 'agent:two', '2026-11-10T09:00:00Z', '2026-11-10T09:00:00Z')`, "CHECK"},
	} {
		if _, err := conn.Exec(test.statement); err == nil || !strings.Contains(err.Error(), test.refusal) {
			t.Errorf("%s = %v, want a %s refusal", test.name, err, test.refusal)
		}
	}
	for _, statement := range []string{
		`UPDATE holds SET ended_at = '2026-11-10T10:05:00Z', end_reason = 'replaced' WHERE id = 'HLD-1'`,
		`INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at, replaces_hold_id) VALUES ('HLD-2', 'FID-a', 'sha256:one', '/w/one', 'agent:one', '2026-11-10T10:05:00Z', '2026-11-10T11:05:00Z', 'HLD-1')`,
		`INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at) VALUES ('HLD-3', 'FID-b', 'sha256:two', '/w/two', 'agent:two', '2026-11-10T10:05:00Z', '2026-11-10T11:05:00Z')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("after the first hold ended, %q: %v", statement, err)
		}
	}
	if _, err := conn.Exec(`INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at, replaces_hold_id, ended_at, end_reason) VALUES ('HLD-4', 'FID-b', 'sha256:three', '/w/three', 'agent:three', '2026-11-10T10:05:00Z', '2026-11-10T11:05:00Z', 'HLD-1', '2026-11-10T10:06:00Z', 'replaced')`); err == nil {
		t.Error("a second hold replacing HLD-1 was accepted")
	}

	if err := migrator.Migrate(19); err != nil {
		t.Fatal(err)
	}
	restored := migration0018SchemaObjects(t, conn)
	if len(restored) != len(before) {
		t.Errorf("down schema object count = %d, want %d", len(restored), len(before))
	}
	for name, sql := range before {
		if restored[name] != sql {
			t.Errorf("down schema object %s differs:\ngot %s\nwant %s", name, restored[name], sql)
		}
	}
}
