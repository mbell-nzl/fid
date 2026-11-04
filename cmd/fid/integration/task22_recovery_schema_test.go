//go:build integration

package integration

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestRecoverySchemaFreshBinaryInitAndRebuild(t *testing.T) {
	t.Parallel()
	workspace := newGitWorkspace(t)
	binary := integrationFIDBinary(t)
	closureRun(t, workspace, binary, "--json", "init", "--prefix", "P22")

	database, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("open recovery projection: %v", err)
	}
	defer func() { _ = database.Close() }()
	for _, table := range []string{"run_runtime", "run_observations", "run_recoveries", "budget_samples", "stops", "attention_items", "activation_resources", "activation_recoveries"} {
		var count int
		if err := database.QueryRow(`SELECT COUNT(*) FROM sqlite_schema WHERE type = 'table' AND name = ?`, table).Scan(&count); err != nil {
			t.Fatalf("inspect %s: %v", table, err)
		}
		if count != 1 {
			t.Fatalf("recovery table %q count = %d, want 1", table, count)
		}
	}
	assertRecoverySchemaRetiredAuthorityAbsent(t, database)

	rebuild, _ := closureRun(t, workspace, binary, "--json", "rebuild")
	var result struct {
		Result struct {
			SchemaVersion int `json:"schema_version"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(rebuild), &result); err != nil {
		t.Fatalf("decode rebuild result: %v", err)
	}
	if result.Result.SchemaVersion != fiddb.ActiveMigrationVersion {
		t.Fatalf("rebuild schema_version = %d, want active schema %d", result.Result.SchemaVersion, fiddb.ActiveMigrationVersion)
	}
	assertRecoverySchemaRetiredAuthorityAbsent(t, database)
}

func assertRecoverySchemaRetiredAuthorityAbsent(t *testing.T, database *sql.DB) {
	t.Helper()
	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM sqlite_schema WHERE type = 'table' AND name = 'escalations'`).Scan(&count); err != nil {
		t.Fatalf("inspect retired escalation authority: %v", err)
	}
	if count != 0 {
		t.Fatalf("retired escalation authority table count = %d, want 0", count)
	}
}
