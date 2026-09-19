package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
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

	if err := migrator.Migrate(13); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migrate(13) error = %v", err)
	}
	afterUp := snapshotDatabaseState(t, conn)

	assertOperationalKernelInventory(t, before, conn)
	assertOperationalKernelConstraints(t, conn)

	if err := migrator.Migrate(12); err != nil {
		t.Fatalf("Migrate(12) down error = %v", err)
	}
	if afterDown := snapshotDatabaseState(t, conn); !reflect.DeepEqual(afterDown, before) {
		t.Fatalf("state after 0013 down = %#v, want pre-0013 state %#v", afterDown, before)
	}

	if err := migrator.Migrate(13); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("second Migrate(13) error = %v", err)
	}
	if afterSecondUp := snapshotDatabaseState(t, conn); !reflect.DeepEqual(afterSecondUp, afterUp) {
		t.Fatalf("state after second 0013 up = %#v, want first post-up state %#v", afterSecondUp, afterUp)
	}
}

func TestMigration0013CancelledReceiptTerminalConstraint(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Migrate(13); err != nil {
		t.Fatalf("Migrate(13) error = %v", err)
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable foreign keys for isolated Receipt constraint: %v", err)
	}
	insert := func(id string, exitCode any, signal any, timedOut int) error {
		_, err := conn.Exec(`
			INSERT INTO execution_receipts
			(id, attempt_id, trust_source, terminal_status, exit_code, signal, timed_out, started_at, ended_at,
			 terminal_repository_json, toolchain_json, output_excerpt, output_digest, output_bytes, output_truncated,
			 redaction_applied, rejected_artifacts_json)
			VALUES (?, ?, 'fid-executed', 'cancelled', ?, ?, ?, '2026-09-07T00:00:00Z', '2026-09-07T00:00:01Z',
			        '{}', '{}', '', 'sha256:test', 0, 0, 0, '[]')
		`, id, "ATT-"+id, exitCode, signal, timedOut)
		return err
	}
	for _, receipt := range []struct {
		id       string
		exitCode any
		signal   any
	}{
		{id: "RCP-cancelled-exit", exitCode: 1},
		{id: "RCP-cancelled-signal", signal: "terminated"},
		{id: "RCP-cancelled-empty"},
	} {
		if err := insert(receipt.id, receipt.exitCode, receipt.signal, 0); err != nil {
			t.Fatalf("insert valid cancelled Receipt %q: %v", receipt.id, err)
		}
	}
	if err := insert("RCP-cancelled-timeout", nil, nil, 1); err == nil {
		t.Fatal("cancelled Receipt with timed_out = 1 succeeded, want CHECK rejection")
	}
	if err := insert("RCP-cancelled-both", 1, "terminated", 0); err == nil {
		t.Fatal("cancelled Receipt with exit code and signal succeeded, want CHECK rejection")
	}
}

func TestMigration0014RetiresV1AuthorityExecutionWithoutReconstruction(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Migrate(13); err != nil {
		t.Fatalf("Migrate(13) error = %v", err)
	}

	for _, table := range []string{"roles", "claims", "gate_runs", "work_claims", "task_completions", "agent_runs", "ac_verifications", "concerns"} {
		if !tableExists(t, conn, table) {
			t.Fatalf("pre-0014 table %q is absent", table)
		}
	}

	if err := migrator.Migrate(14); err != nil {
		t.Fatalf("Migrate(14) error = %v", err)
	}

	for _, table := range []string{"roles", "gate_runs", "work_claims"} {
		if tableExists(t, conn, table) {
			t.Fatalf("post-0014 removed table %q is present", table)
		}
	}
	if !tableExists(t, conn, "claims") {
		t.Fatal("post-0014 generic claims table is absent")
	}
	for _, table := range []string{"task_completions", "agent_runs", "ac_verifications", "concerns"} {
		if !tableExists(t, conn, table) {
			t.Fatalf("post-0014 frozen table %q is absent", table)
		}
	}

	var claimColumns string
	if err := conn.QueryRow(`SELECT group_concat(name, ',') FROM pragma_table_info('claims')`).Scan(&claimColumns); err != nil {
		t.Fatalf("inspect post-0014 claims columns: %v", err)
	}
	if !strings.Contains(claimColumns, "envelope_id") || strings.Contains(claimColumns, "ref_type") {
		t.Fatalf("post-0014 claims columns = %q, want generic claim shape only", claimColumns)
	}

	rows, err := conn.Query(`SELECT workflow_id, rule FROM workflow_transition_requires ORDER BY workflow_id, rowid`)
	if err != nil {
		t.Fatalf("query post-0014 transition requirements: %v", err)
	}
	defer closeRows(t, rows)
	var rules []string
	for rows.Next() {
		var workflowID, rule string
		if err := rows.Scan(&workflowID, &rule); err != nil {
			t.Fatalf("scan post-0014 transition requirement: %v", err)
		}
		rules = append(rules, workflowID+":"+rule)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate post-0014 transition requirements: %v", err)
	}
	if want := []string{"epic-default:all_child_tasks_done_or_cancelled", "escalation-default:linked_to_decision"}; !reflect.DeepEqual(rules, want) {
		t.Fatalf("post-0014 transition requirements = %v, want %v", rules, want)
	}

	if err := migrator.Migrate(13); err != nil {
		t.Fatalf("Migrate(13) after 0014 error = %v", err)
	}
	for _, table := range []string{"roles", "claims", "gate_runs", "work_claims"} {
		if !tableExists(t, conn, table) {
			t.Fatalf("empty-state 0014 down table %q is absent", table)
		}
	}
	if err := migrator.Migrate(14); err != nil {
		t.Fatalf("second Migrate(14) error = %v", err)
	}
}

func TestOperationalKernelInventoryComparisonRejectsSchemaDrift(t *testing.T) {
	want := map[string]string{
		"table:expected": "CREATE TABLE expected (id TEXT PRIMARY KEY)",
	}
	tests := []struct {
		name  string
		apply func(map[string]string)
	}{
		{
			name: "extra table",
			apply: func(got map[string]string) {
				got["table:unexpected"] = "CREATE TABLE unexpected (id TEXT PRIMARY KEY)"
			},
		},
		{
			name: "extra index",
			apply: func(got map[string]string) {
				got["index:unexpected"] = "CREATE INDEX unexpected ON expected(id)"
			},
		},
		{
			name: "extra check",
			apply: func(got map[string]string) {
				got["table:expected"] = "CREATE TABLE expected (id TEXT PRIMARY KEY, CHECK (length(id) > 0))"
			},
		},
		{
			name: "extra foreign key",
			apply: func(got map[string]string) {
				got["table:expected"] = "CREATE TABLE expected (id TEXT PRIMARY KEY, parent_id TEXT REFERENCES parents(id))"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := map[string]string{
				"table:expected": "CREATE TABLE expected (id TEXT PRIMARY KEY)",
			}
			tt.apply(got)

			if err := compareOperationalKernelInventory(want, got); err == nil {
				t.Fatalf("exact inventory accepted %s", tt.name)
			}
		})
	}
}

func assertOperationalKernelInventory(t *testing.T, before databaseState, conn *sql.DB) {
	t.Helper()

	got := operationalKernelSchemaObjects(t, before, conn)
	if err := compareOperationalKernelInventory(expectedOperationalKernelSchemaObjects(), got); err != nil {
		t.Fatal(err)
	}
}

func compareOperationalKernelInventory(want, got map[string]string) error {
	if reflect.DeepEqual(want, got) {
		return nil
	}
	return fmt.Errorf("0013 schema inventory = %#v, want %#v", got, want)
}

func operationalKernelSchemaObjects(t *testing.T, before databaseState, conn *sql.DB) map[string]string {
	t.Helper()

	beforeKeys := make(map[string]struct{}, len(before.schema))
	for _, object := range before.schema {
		parts := strings.SplitN(object, ":", 3)
		beforeKeys[parts[0]+":"+parts[1]] = struct{}{}
	}

	rows, err := conn.Query(`
		SELECT type, name, sql
		FROM sqlite_schema
		WHERE name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'
		ORDER BY type, name
	`)
	if err != nil {
		t.Fatalf("query post-0013 schema objects: %v", err)
	}
	defer closeRows(t, rows)

	got := make(map[string]string)
	for rows.Next() {
		var objectType, name, sqlText string
		if err := rows.Scan(&objectType, &name, &sqlText); err != nil {
			t.Fatalf("scan post-0013 schema object: %v", err)
		}
		key := objectType + ":" + name
		if _, existed := beforeKeys[key]; !existed {
			got[key] = normalizeSchemaSQL(sqlText)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate post-0013 schema objects: %v", err)
	}
	return got
}

func normalizeSchemaSQL(sqlText string) string {
	return strings.TrimSuffix(strings.Join(strings.Fields(sqlText), " "), ";")
}

func expectedOperationalKernelSchemaObjects() map[string]string {
	return normalizedSchemaObjects(map[string]string{
		"table:system_capabilities": `CREATE TABLE system_capabilities (
  capability_id TEXT PRIMARY KEY,
  schema_min INTEGER NOT NULL CHECK (schema_min > 0),
  schema_max INTEGER NOT NULL CHECK (schema_max >= schema_min),
  available INTEGER NOT NULL CHECK (available IN (0, 1)),
  limits_json TEXT NOT NULL CHECK (json_valid(limits_json)),
  trust_source TEXT NOT NULL CHECK (trust_source IN (
    'active-kit-activation','fid-executed','fid-git-observation','fid-journal',
    'phase-22-unavailable','phase-23-unavailable','phase-24-unavailable','phase-25-unavailable'
  ))
)`,
		"table:idempotency_results": `CREATE TABLE idempotency_results (
  operation_scope TEXT NOT NULL,
  idempotency_key TEXT NOT NULL CHECK (length(idempotency_key) BETWEEN 1 AND 128),
  request_digest TEXT NOT NULL,
  result_entity_id TEXT NOT NULL,
  result_event_id TEXT NOT NULL,
  aggregate_sequence INTEGER NOT NULL CHECK (aggregate_sequence > 0),
  result_json TEXT NOT NULL CHECK (json_valid(result_json)),
  created_at DATETIME NOT NULL,
  PRIMARY KEY (operation_scope, idempotency_key)
)`,
		"table:kit_activations": `CREATE TABLE kit_activations (
  id TEXT PRIMARY KEY,
  prior_activation_id TEXT REFERENCES kit_activations(id),
  kit_identity TEXT NOT NULL,
  backend_identity TEXT NOT NULL,
  phase INTEGER NOT NULL CHECK (phase > 0),
  status TEXT NOT NULL CHECK (status IN ('prepared','staged','committed','finalized','aborted','superseded')),
  point_of_no_return INTEGER NOT NULL DEFAULT 0 CHECK (point_of_no_return IN (0, 1)),
  direction TEXT NOT NULL CHECK (direction = 'activate'),
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
)`,
		"index:kit_activations_one_nonterminal": `CREATE UNIQUE INDEX kit_activations_one_nonterminal
  ON kit_activations((1)) WHERE status IN ('prepared','staged','committed')`,
		"index:kit_activations_one_active": `CREATE UNIQUE INDEX kit_activations_one_active
  ON kit_activations((1)) WHERE status = 'finalized'`,
		"table:activation_components": `CREATE TABLE activation_components (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  component_type TEXT NOT NULL CHECK (component_type IN ('authority-profile','gate-definitions')),
  component_id TEXT NOT NULL,
  component_version INTEGER NOT NULL CHECK (component_version > 0),
  content_digest TEXT NOT NULL,
  expected_prior_digest TEXT,
  status TEXT NOT NULL CHECK (status IN ('staged','committed','finalized','aborted')),
  evidence_json TEXT NOT NULL CHECK (json_valid(evidence_json)),
  error_code TEXT,
  PRIMARY KEY (activation_id, component_type, component_id, component_version)
)`,
		"table:activation_actor_profiles": `CREATE TABLE activation_actor_profiles (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  actor_id TEXT NOT NULL,
  actor_type TEXT NOT NULL CHECK (actor_type IN ('agent','human','service')),
  producer TEXT NOT NULL,
  provenance_ref TEXT NOT NULL,
  session_ref TEXT,
  capabilities_json TEXT NOT NULL CHECK (json_valid(capabilities_json)),
  PRIMARY KEY (activation_id, actor_id)
)`,
		"table:activation_gate_definitions": `CREATE TABLE activation_gate_definitions (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  definition_id TEXT NOT NULL,
  definition_version INTEGER NOT NULL CHECK (definition_version > 0),
  definition_digest TEXT NOT NULL,
  definition_json TEXT NOT NULL CHECK (json_valid(definition_json)),
  PRIMARY KEY (activation_id, definition_id, definition_version),
  UNIQUE (activation_id, definition_digest)
)`,
		"table:gate_registrations": `CREATE TABLE gate_registrations (
  id TEXT NOT NULL,
  version INTEGER NOT NULL CHECK (version > 0),
  activation_id TEXT NOT NULL,
  definition_id TEXT NOT NULL,
  definition_version INTEGER NOT NULL,
  definition_digest TEXT NOT NULL,
  binding_digest TEXT NOT NULL,
  project_id TEXT NOT NULL,
  worktree_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  policy_json TEXT NOT NULL CHECK (json_valid(policy_json)),
  authorization_provenance TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  PRIMARY KEY (id, version),
  UNIQUE (definition_digest, binding_digest),
  FOREIGN KEY (activation_id, definition_id, definition_version)
    REFERENCES activation_gate_definitions(activation_id, definition_id, definition_version)
)`,
		"table:work_envelopes": `CREATE TABLE work_envelopes (
  id TEXT PRIMARY KEY,
  current_revision INTEGER NOT NULL CHECK (current_revision = 1),
  status TEXT NOT NULL CHECK (status IN ('draft','active','completed')),
  created_at DATETIME NOT NULL
)`,
		"table:envelope_revisions": `CREATE TABLE envelope_revisions (
  envelope_id TEXT NOT NULL REFERENCES work_envelopes(id),
  revision INTEGER NOT NULL CHECK (revision = 1),
  prior_revision INTEGER CHECK (prior_revision IS NULL),
  outcome TEXT NOT NULL,
  non_goals_json TEXT NOT NULL CHECK (json_valid(non_goals_json)),
  risk_result_json TEXT NOT NULL CHECK (json_valid(risk_result_json)),
  repository_scope_json TEXT NOT NULL CHECK (json_valid(repository_scope_json)),
  status TEXT NOT NULL CHECK (status IN ('draft','active','completed')),
  created_actor_id TEXT NOT NULL,
  approved_actor_id TEXT,
  created_at DATETIME NOT NULL,
  activated_at DATETIME,
  completed_at DATETIME,
  PRIMARY KEY (envelope_id, revision)
)`,
		"table:envelope_sources": `CREATE TABLE envelope_sources (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  ordinal INTEGER NOT NULL CHECK (ordinal >= 0),
  source_json TEXT NOT NULL CHECK (json_valid(source_json)),
  PRIMARY KEY (envelope_id, revision, ordinal),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
)`,
		"table:envelope_issue_selection": `CREATE TABLE envelope_issue_selection (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  primary_issue INTEGER NOT NULL CHECK (primary_issue = 1),
  PRIMARY KEY (envelope_id, revision, issue_id),
  UNIQUE (envelope_id, revision, primary_issue),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
)`,
		"table:envelope_authority_grants": `CREATE TABLE envelope_authority_grants (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  ordinal INTEGER NOT NULL CHECK (ordinal >= 0),
  actor_selector_json TEXT NOT NULL CHECK (json_valid(actor_selector_json)),
  actions_json TEXT NOT NULL CHECK (json_valid(actions_json)),
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  PRIMARY KEY (envelope_id, revision, ordinal),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
)`,
		"table:envelope_obligations": `CREATE TABLE envelope_obligations (
  id TEXT PRIMARY KEY,
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
	  obligation_type TEXT NOT NULL CHECK (obligation_type IN ('gate-evidence','review','human-outcome')),
	  registration_id TEXT,
	  registration_version INTEGER,
	  obligation_json TEXT NOT NULL CHECK (json_valid(obligation_json)),
	  CHECK ((obligation_type = 'gate-evidence'
	          AND registration_id IS NOT NULL AND registration_version IS NOT NULL)
	      OR (obligation_type IN ('review','human-outcome')
	          AND registration_id IS NULL AND registration_version IS NULL)),
  UNIQUE (envelope_id, revision, id),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision),
  FOREIGN KEY (registration_id, registration_version) REFERENCES gate_registrations(id, version)
)`,
		"table:envelope_budgets": `CREATE TABLE envelope_budgets (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  name TEXT NOT NULL,
  budget_json TEXT NOT NULL CHECK (json_valid(budget_json)),
  PRIMARY KEY (envelope_id, revision, name),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
)`,
		"table:work_claims": `CREATE TABLE work_claims (
  id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  worktree_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
	  mode TEXT NOT NULL CHECK (mode IN ('mutation','review')),
  version INTEGER NOT NULL CHECK (version > 0),
  status TEXT NOT NULL CHECK (status IN ('active','released','expired','superseded')),
  acquired_at DATETIME NOT NULL,
  expires_at DATETIME NOT NULL,
  released_at DATETIME,
  release_reason TEXT CHECK (release_reason IS NULL OR release_reason IN (
    'completed','failed','operator-released','interrupted','recovery-superseded','abandoned'
  )),
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
)`,
		"index:work_claims_one_active_mutation": `CREATE UNIQUE INDEX work_claims_one_active_mutation
  ON work_claims(issue_id) WHERE status = 'active'`,
		"table:runs": `CREATE TABLE runs (
  id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  claim_id TEXT NOT NULL REFERENCES work_claims(id),
  claim_version INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  start_repository_json TEXT NOT NULL CHECK (json_valid(start_repository_json)),
  terminal_repository_json TEXT CHECK (terminal_repository_json IS NULL OR json_valid(terminal_repository_json)),
  status TEXT NOT NULL CHECK (status IN ('active','succeeded','failed','abandoned','interrupted')),
  started_at DATETIME NOT NULL,
  ended_at DATETIME,
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
)`,
		"index:runs_one_active_per_claim": `CREATE UNIQUE INDEX runs_one_active_per_claim
  ON runs(claim_id) WHERE status = 'active'`,
		"table:run_budget_bindings": `CREATE TABLE run_budget_bindings (
  run_id TEXT NOT NULL REFERENCES runs(id),
  name TEXT NOT NULL,
  budget_json TEXT NOT NULL CHECK (json_valid(budget_json)),
  PRIMARY KEY (run_id, name)
)`,
		"table:execution_attempts": `CREATE TABLE execution_attempts (
  id TEXT PRIMARY KEY,
  registration_id TEXT NOT NULL,
  registration_version INTEGER NOT NULL,
  obligation_id TEXT NOT NULL REFERENCES envelope_obligations(id),
  run_id TEXT NOT NULL REFERENCES runs(id),
  actor_id TEXT NOT NULL,
  start_repository_json TEXT NOT NULL CHECK (json_valid(start_repository_json)),
  resolved_policy_digest TEXT NOT NULL,
  producer TEXT NOT NULL CHECK (producer = 'fid'),
  started_at DATETIME NOT NULL,
  FOREIGN KEY (registration_id, registration_version) REFERENCES gate_registrations(id, version)
)`,
		"table:referenced_artifacts": `CREATE TABLE referenced_artifacts (
  id TEXT PRIMARY KEY,
  digest TEXT NOT NULL UNIQUE,
  media_type TEXT NOT NULL,
  size_bytes INTEGER NOT NULL CHECK (size_bytes >= 0),
  retention_class TEXT NOT NULL CHECK (retention_class = 'evidence'),
  availability TEXT NOT NULL CHECK (availability IN ('available','missing')),
  location_policy TEXT NOT NULL CHECK (location_policy = 'content-addressed'),
  created_at DATETIME NOT NULL,
  expires_at DATETIME
)`,
		"table:execution_receipts": `CREATE TABLE execution_receipts (
  id TEXT PRIMARY KEY,
  attempt_id TEXT NOT NULL UNIQUE REFERENCES execution_attempts(id),
  trust_source TEXT NOT NULL CHECK (trust_source = 'fid-executed'),
  terminal_status TEXT NOT NULL CHECK (terminal_status IN ('succeeded','failed','timed-out','signaled','spawn-failed','cancelled')),
  exit_code INTEGER,
  signal TEXT,
  timed_out INTEGER NOT NULL CHECK (timed_out IN (0, 1)),
  started_at DATETIME NOT NULL,
  ended_at DATETIME NOT NULL,
  terminal_repository_json TEXT NOT NULL CHECK (json_valid(terminal_repository_json)),
  toolchain_json TEXT NOT NULL CHECK (json_valid(toolchain_json)),
  output_excerpt TEXT NOT NULL,
  output_digest TEXT NOT NULL,
  output_bytes INTEGER NOT NULL CHECK (output_bytes >= 0),
  output_truncated INTEGER NOT NULL CHECK (output_truncated IN (0, 1)),
  redaction_applied INTEGER NOT NULL CHECK (redaction_applied IN (0, 1)),
  rejected_artifacts_json TEXT NOT NULL
    CHECK (json_valid(rejected_artifacts_json) AND json_type(rejected_artifacts_json) = 'array'),
  CHECK ((terminal_status = 'spawn-failed' AND exit_code IS NULL AND signal IS NULL AND timed_out = 0)
      OR (terminal_status = 'timed-out' AND timed_out = 1)
      OR (terminal_status = 'signaled' AND exit_code IS NULL AND signal IS NOT NULL AND timed_out = 0)
      OR (terminal_status IN ('succeeded','failed') AND exit_code IS NOT NULL AND signal IS NULL AND timed_out = 0)
      OR (terminal_status = 'cancelled' AND timed_out = 0 AND NOT (exit_code IS NOT NULL AND signal IS NOT NULL)))
)`,
		"table:receipt_artifacts": `CREATE TABLE receipt_artifacts (
  receipt_id TEXT NOT NULL REFERENCES execution_receipts(id),
  artifact_id TEXT NOT NULL REFERENCES referenced_artifacts(id),
  declared_path TEXT NOT NULL,
  PRIMARY KEY (receipt_id, artifact_id)
)`,
		"table:evidence_records": `CREATE TABLE evidence_records (
  id TEXT PRIMARY KEY,
  obligation_id TEXT NOT NULL REFERENCES envelope_obligations(id),
  receipt_id TEXT NOT NULL REFERENCES execution_receipts(id),
  evidence_type TEXT NOT NULL CHECK (evidence_type = 'receipt'),
  producer TEXT NOT NULL CHECK (producer = 'fid'),
  trust_source TEXT NOT NULL CHECK (trust_source = 'fid-executed'),
  target_repository_json TEXT NOT NULL CHECK (json_valid(target_repository_json)),
  validity TEXT NOT NULL CHECK (validity IN ('valid','invalid')),
  created_at DATETIME NOT NULL,
  UNIQUE (obligation_id, receipt_id)
)`,
		"table:evidence_artifact_links": `CREATE TABLE evidence_artifact_links (
  evidence_id TEXT NOT NULL REFERENCES evidence_records(id),
  artifact_id TEXT NOT NULL REFERENCES referenced_artifacts(id),
  required_material INTEGER NOT NULL CHECK (required_material IN (0, 1)),
  PRIMARY KEY (evidence_id, artifact_id)
)`,
	})
}

func normalizedSchemaObjects(objects map[string]string) map[string]string {
	for key, sqlText := range objects {
		objects[key] = normalizeSchemaSQL(sqlText)
	}
	return objects
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
