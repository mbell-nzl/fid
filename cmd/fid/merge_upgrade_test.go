package main

import (
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestCaughtUpVersion9ProjectionGainsMergeCIStorage(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runFIDCommand(t, workspace, "init", "--prefix", "FID")
	runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Preserved task", "--actor", "human:mark")
	before := executeCommandForTest(workspace, "--json", "issue", "list")
	if before.err != nil {
		t.Fatal(before.err)
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	// Model an existing caught-up main projection: the checkpoint already
	// matches the journal, but its schema lacks the new observation tables.
	for _, sql := range []string{`DROP TABLE ci_results`, `DROP TABLE merge_observations`, `UPDATE schema_migrations SET version=27`, `UPDATE projection_meta SET projection_schema_version=9`} {
		if _, err := conn.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	after := executeCommandForTest(workspace, "--json", "issue", "list")
	if after.err != nil || after.stdout != before.stdout {
		t.Fatalf("upgrade changed tasks: %#v", after)
	}
	conn, err = fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, table := range []string{"merge_observations", "ci_results"} {
		var count int
		if err := conn.QueryRow("SELECT count(*) FROM " + table).Scan(&count); err != nil || count != 0 {
			t.Fatalf("upgraded %s=%d/%v", table, count, err)
		}
	}
	var cache, migration int
	if err := conn.QueryRow(`SELECT projection_schema_version FROM projection_meta`).Scan(&cache); err != nil {
		t.Fatal(err)
	}
	if err := conn.QueryRow(`SELECT version FROM schema_migrations`).Scan(&migration); err != nil {
		t.Fatal(err)
	}
	if cache != fiddb.ProjectionSchemaVersion || migration != fiddb.ActiveMigrationVersion {
		t.Fatalf("upgraded cache/migration=%d/%d", cache, migration)
	}
}
