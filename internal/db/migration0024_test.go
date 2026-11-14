package db

import (
	"path/filepath"
	"testing"
)

func TestMigration0024PreservesAllowedStatesAndConstrainsTaskLanes(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(23); err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('E','epic','Epic','epic-default','in_progress','tl','2026-11-14T00:00:00Z','2026-11-14T00:00:00Z')`,
		`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('T','task','Task','task-default','todo','developer','2026-11-14T00:00:00Z','2026-11-14T00:00:00Z')`,
		`INSERT INTO issues(id,kind,parent_id,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('S','subtask','T','Subtask','subtask-default','todo','developer','2026-11-14T00:00:00Z','2026-11-14T00:00:00Z')`,
	} {
		if _, err := conn.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	if err := migrator.Migrate(24); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := conn.QueryRow(`SELECT count(*) FROM issues WHERE lane IS NULL`).Scan(&n); err != nil || n != 3 {
		t.Fatalf("upgraded rows = %d, %v", n, err)
	}
	for _, row := range []struct {
		kind  string
		lane  any
		valid bool
	}{
		{"T", "routine", true}, {"T", "standard", true}, {"T", "assured", true}, {"T", nil, true}, {"E", "standard", false}, {"S", "routine", false}, {"T", "urgent", false}, {"T", "Standard", false},
	} {
		_, err := conn.Exec(`UPDATE issues SET lane=? WHERE id=?`, row.lane, row.kind)
		if (err == nil) != row.valid {
			t.Fatalf("lane %v on %s = %v", row.lane, row.kind, err)
		}
	}
	for _, row := range []struct {
		id, state string
		valid     bool
	}{
		{"T", "in_review", true}, {"T", "dev_complete", false}, {"T", "qa_pass", false}, {"E", "done", true}, {"E", "wrapped", false}, {"E", "closed", false}, {"S", "in_progress", true}, {"S", "done", true},
	} {
		_, err := conn.Exec(`UPDATE issues SET status=? WHERE id=?`, row.state, row.id)
		if (err == nil) != row.valid {
			t.Fatalf("state %s on %s = %v", row.state, row.id, err)
		}
	}
	// Allowed shared states can be downgraded without a converter; new states
	// deliberately have no mapping back into the retired graph.
	if _, err := conn.Exec(`UPDATE issues SET status='todo',lane=NULL WHERE id IN ('T','E')`); err != nil {
		t.Fatal(err)
	}
	if err := migrator.Migrate(23); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`SELECT lane FROM issues`); err == nil {
		t.Fatal("downgrade retained lane column")
	}
	if err := conn.QueryRow(`SELECT count(*) FROM workflow_states WHERE workflow_id='task-default' AND state IN ('dev_complete','qa_pass')`).Scan(&n); err != nil || n != 2 {
		t.Fatalf("restored graph = %d, %v", n, err)
	}
}

func TestMigration0024DoesNotConvertNewStatesOnDowngrade(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ kind, state string }{{"task", "in_review"}, {"epic", "done"}} {
		t.Run(row.kind, func(t *testing.T) {
			t.Parallel()
			migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
			defer closeMigrator(t, migrator, conn)
			if err := migrator.Migrate(24); err != nil {
				t.Fatal(err)
			}
			if _, err := conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('new',?,'New state',?,?,'developer','2026-11-14T00:00:00Z','2026-11-14T00:00:00Z')`, row.kind, row.kind+"-default", row.state); err != nil {
				t.Fatal(err)
			}
			if err := migrator.Migrate(23); err == nil {
				t.Fatal("downgrade silently accepted a state absent from the old graph")
			}
			var status string
			if err := conn.QueryRow(`SELECT status FROM issues WHERE id='new'`).Scan(&status); err != nil || status != row.state {
				t.Fatalf("downgrade converted status to %q: %v", status, err)
			}
		})
	}
}
