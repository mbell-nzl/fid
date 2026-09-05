package db

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestMigration0015RecoverySchemaConstraintsAndRoundTrip(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() {
		if conn != nil {
			closeDB(t, conn)
		}
	}()
	assertRecoverySchema(t, conn)
	assertRecoveryConstraints(t, conn)
	assertActiveMigration(t, conn)

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() before down error = %v", err)
	}
	conn = nil
	if err := MigrateDown(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateDown() error = %v", err)
	}
	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("second MigrateUp() error = %v", err)
	}

	conn, err = Open(dbPath)
	if err != nil {
		t.Fatalf("Open() after round trip error = %v", err)
	}
	assertRecoverySchema(t, conn)
	assertRecoveryTablesEmpty(t, conn)
	assertActiveMigration(t, conn)
}

func assertActiveMigration(t *testing.T, conn *sql.DB) {
	t.Helper()
	var version int
	if err := conn.QueryRow(`SELECT version FROM schema_migrations`).Scan(&version); err != nil {
		t.Fatalf("read active schema migration: %v", err)
	}
	if version != ActiveMigrationVersion {
		t.Fatalf("active schema migration = %d, want %d", version, ActiveMigrationVersion)
	}
}

func assertRecoverySchema(t *testing.T, conn *sql.DB) {
	t.Helper()
	for _, table := range []string{
		"run_runtime", "run_observations", "run_recoveries", "budget_samples",
		"stops", "stop_resolutions", "attention_items", "attention_stop_links",
		"activation_resources", "activation_recoveries",
	} {
		assertTableExists(t, conn, table)
	}

	var notNull int
	if err := conn.QueryRow(`SELECT "notnull" FROM pragma_table_info('run_recoveries') WHERE name = 'observation_id'`).Scan(&notNull); err != nil {
		t.Fatalf("inspect run_recoveries.observation_id: %v", err)
	}
	if notNull != 1 {
		t.Fatalf("run_recoveries.observation_id notnull = %d, want 1", notNull)
	}
}

func assertRecoveryConstraints(t *testing.T, conn *sql.DB) {
	t.Helper()
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable foreign keys for NOT NULL proof: %v", err)
	}
	defer func() { _, _ = conn.Exec(`PRAGMA foreign_keys = ON`) }()
	if _, err := conn.Exec(`
		INSERT INTO run_recoveries (
			id, prior_run_id, successor_run_id, prior_claim_id, successor_claim_id,
			observation_id, reconciliation_json, adopted_paths_json, authority_json,
			reason, created_event_id
		) VALUES (?, ?, ?, ?, ?, NULL, ?, ?, ?, ?, ?)
	`, "REC-null-observation", "RUN-prior", "RUN-successor", "CLM-prior", "CLM-successor", "{}", "[]", "{}", "test", "EVT-recovery-null"); err == nil {
		t.Fatal("run_recoveries accepted NULL observation_id")
	}

	for _, code := range recoveryStopCodes {
		if _, err := conn.Exec(`
			INSERT INTO stops (
				id, subject_type, subject_id, code, blocked_actions_json, details_json,
				resolution_condition_json, resolver_capability, actor_id, raised_at, created_event_id
			) VALUES (?, 'run', 'RUN-1', ?, '[]', '{}', '{}', 'stop.resolve', 'actor-1', '2026-09-05T00:00:00Z', ?)
		`, "STOP-"+code, code, "EVT-stop-"+code); err != nil {
			t.Fatalf("stops accepts code %q: %v", code, err)
		}
		if _, err := conn.Exec(`
			INSERT INTO attention_items (
				id, semantic_key, subject_type, subject_id, code, status, version, raised_at
			) VALUES (?, ?, 'run', 'RUN-1', ?, 'open', 1, '2026-09-05T00:00:00Z')
		`, "ATT-"+code, "attention-"+code, code); err != nil {
			t.Fatalf("attention_items accepts code %q: %v", code, err)
		}
	}
	for _, statement := range []string{
		`INSERT INTO stops (id, subject_type, subject_id, code, blocked_actions_json, details_json, resolution_condition_json, resolver_capability, actor_id, raised_at, created_event_id) VALUES ('STOP-invalid', 'run', 'RUN-1', 'invalid', '[]', '{}', '{}', 'stop.resolve', 'actor-1', '2026-09-05T00:00:00Z', 'EVT-stop-invalid')`,
		`INSERT INTO attention_items (id, semantic_key, subject_type, subject_id, code, status, version, raised_at) VALUES ('ATT-invalid', 'attention-invalid', 'run', 'RUN-1', 'invalid', 'open', 1, '2026-09-05T00:00:00Z')`,
	} {
		if _, err := conn.Exec(statement); err == nil {
			t.Fatalf("recovery code constraint accepted invalid code: %s", statement)
		}
	}
}

func assertRecoveryTablesEmpty(t *testing.T, conn *sql.DB) {
	t.Helper()
	for _, table := range []string{
		"run_runtime", "run_observations", "run_recoveries", "budget_samples",
		"stops", "stop_resolutions", "attention_items", "attention_stop_links",
		"activation_resources", "activation_recoveries",
	} {
		var count int
		if err := conn.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("%s rows after down/up = %d, want 0", table, count)
		}
	}
}

var recoveryStopCodes = []string{
	"executor-liveness-unknown",
	"recovery-state-conflict",
	"budget-exhausted",
	"budget-data-unavailable",
	"activation-incomplete",
	"capability-unavailable",
	"operator-stop",
}
