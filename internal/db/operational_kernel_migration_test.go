package db

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/golang-migrate/migrate/v4"
)

func TestMigration0013OperationalKernelInventoryAndRoundTrip(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Migrate(12); err != nil {
		t.Fatalf("Migrate(12) error = %v", err)
	}
	before := snapshotDatabaseState(t, conn)

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Up() error = %v", err)
	}
	afterUp := snapshotDatabaseState(t, conn)

	assertOperationalKernelInventory(t, conn)
	assertOperationalKernelConstraints(t, conn)

	if err := migrator.Migrate(12); err != nil {
		t.Fatalf("Migrate(12) down error = %v", err)
	}
	if afterDown := snapshotDatabaseState(t, conn); !reflect.DeepEqual(afterDown, before) {
		t.Fatalf("state after 0013 down = %#v, want pre-0013 state %#v", afterDown, before)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("second Up() error = %v", err)
	}
	if afterSecondUp := snapshotDatabaseState(t, conn); !reflect.DeepEqual(afterSecondUp, afterUp) {
		t.Fatalf("state after second 0013 up = %#v, want first post-up state %#v", afterSecondUp, afterUp)
	}
}

func assertOperationalKernelInventory(t *testing.T, conn *sql.DB) {
	t.Helper()

	wantTables := []string{
		"activation_actor_profiles",
		"activation_components",
		"activation_gate_definitions",
		"envelope_authority_grants",
		"envelope_budgets",
		"envelope_issue_selection",
		"envelope_obligations",
		"envelope_revisions",
		"envelope_sources",
		"evidence_artifact_links",
		"evidence_records",
		"execution_attempts",
		"execution_receipts",
		"gate_registrations",
		"idempotency_results",
		"kit_activations",
		"receipt_artifacts",
		"referenced_artifacts",
		"run_budget_bindings",
		"runs",
		"system_capabilities",
		"work_claims",
		"work_envelopes",
	}
	gotTables := schemaObjectNames(t, conn, "table")
	for _, table := range wantTables {
		if !containsString(gotTables, table) {
			t.Fatalf("0013 schema missing table %q; got %v", table, gotTables)
		}
	}

	assertSchemaContains(t, conn, "system_capabilities",
		"CHECK (schema_min > 0)",
		"CHECK (schema_max >= schema_min)",
		"CHECK (available IN (0, 1))",
		"CHECK (json_valid(limits_json))",
	)
	assertSchemaContains(t, conn, "kit_activations",
		"CHECK (status IN ('prepared','staged','committed','finalized','aborted'))",
		"CHECK (point_of_no_return IN (0, 1))",
		"CHECK (direction = 'activate')",
	)
	assertSchemaContains(t, conn, "envelope_issue_selection", "CHECK (primary_issue = 1)")
	assertSchemaContains(t, conn, "work_claims",
		"CHECK (mode = 'mutation')",
		"CHECK (status IN ('active','released','expired'))",
	)
	assertSchemaContains(t, conn, "execution_receipts",
		"CHECK (terminal_status IN ('succeeded','failed','timed-out','signaled','spawn-failed'))",
		"terminal_status = 'spawn-failed' AND exit_code IS NULL AND signal IS NULL AND timed_out = 0",
	)
	assertSchemaContains(t, conn, "evidence_records",
		"CHECK (evidence_type = 'receipt')",
		"CHECK (trust_source = 'fid-executed')",
	)

	assertForeignKeysContain(t, conn, "activation_components", "activation_id->kit_activations.id")
	assertForeignKeysContain(t, conn, "gate_registrations", "activation_id->activation_gate_definitions.activation_id")
	assertForeignKeysContain(t, conn, "envelope_issue_selection",
		"issue_id->issues.id",
		"envelope_id->envelope_revisions.envelope_id",
	)
	assertForeignKeysContain(t, conn, "work_claims",
		"issue_id->issues.id",
		"envelope_id->envelope_revisions.envelope_id",
	)
	assertForeignKeysContain(t, conn, "execution_receipts", "attempt_id->execution_attempts.id")
	assertForeignKeysContain(t, conn, "evidence_artifact_links",
		"evidence_id->evidence_records.id",
		"artifact_id->referenced_artifacts.id",
	)

	wantIndexes := []string{
		"kit_activations_one_active",
		"kit_activations_one_nonterminal",
		"runs_one_active_per_claim",
		"work_claims_one_active_mutation",
	}
	gotIndexes := schemaObjectNames(t, conn, "index")
	for _, index := range wantIndexes {
		if !containsString(gotIndexes, index) {
			t.Fatalf("0013 schema missing index %q; got %v", index, gotIndexes)
		}
	}
}

func assertOperationalKernelConstraints(t *testing.T, conn *sql.DB) {
	t.Helper()

	if _, err := conn.Exec(`
		INSERT INTO system_capabilities
		(capability_id, schema_min, schema_max, available, limits_json, trust_source)
		VALUES ('workflow.state.v1', 1, 1, 1, '{}', 'fid-journal')
	`); err != nil {
		t.Fatalf("insert valid system capability: %v", err)
	}
	if _, err := conn.Exec(`
		INSERT INTO system_capabilities
		(capability_id, schema_min, schema_max, available, limits_json, trust_source)
		VALUES ('invalid-capability', 1, 1, 2, '{}', 'fid-journal')
	`); err == nil {
		t.Fatal("invalid capability availability insert succeeded, want CHECK rejection")
	}

	if _, err := conn.Exec(`
		INSERT INTO idempotency_results
		(operation_scope, idempotency_key, request_digest, result_entity_id, result_event_id, aggregate_sequence, result_json, created_at)
		VALUES ('issue.create:project-1', 'key-1', 'sha256:request', 'ISS-1', 'EVT-1', 1, '{}', '2026-08-23T00:00:00Z')
	`); err != nil {
		t.Fatalf("insert valid idempotency result: %v", err)
	}
	if _, err := conn.Exec(`
		INSERT INTO idempotency_results
		(operation_scope, idempotency_key, request_digest, result_entity_id, result_event_id, aggregate_sequence, result_json, created_at)
		VALUES ('issue.create:project-1', '', 'sha256:request', 'ISS-2', 'EVT-2', 1, '{}', '2026-08-23T00:00:00Z')
	`); err == nil {
		t.Fatal("empty idempotency key insert succeeded, want CHECK rejection")
	}

	if _, err := conn.Exec(`
		INSERT INTO kit_activations
		(id, kit_identity, backend_identity, phase, status, direction, created_at, updated_at)
		VALUES ('ACT-1', 'kit-1', 'fid-1', 21, 'prepared', 'activate', '2026-08-23T00:00:00Z', '2026-08-23T00:00:00Z')
	`); err != nil {
		t.Fatalf("insert valid activation: %v", err)
	}
	if _, err := conn.Exec(`
		INSERT INTO kit_activations
		(id, kit_identity, backend_identity, phase, status, direction, created_at, updated_at)
		VALUES ('ACT-2', 'kit-2', 'fid-1', 21, 'staged', 'activate', '2026-08-23T00:00:00Z', '2026-08-23T00:00:00Z')
	`); err == nil {
		t.Fatal("second nonterminal activation insert succeeded, want partial index rejection")
	}
}

func schemaObjectNames(t *testing.T, conn *sql.DB, objectType string) []string {
	t.Helper()

	rows, err := conn.Query(`
		SELECT name
		FROM sqlite_schema
		WHERE type = ? AND name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'
	`, objectType)
	if err != nil {
		t.Fatalf("query %s names: %v", objectType, err)
	}
	defer closeRows(t, rows)

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan %s name: %v", objectType, err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate %s names: %v", objectType, err)
	}
	sort.Strings(names)
	return names
}
