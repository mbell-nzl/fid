package db

import (
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestMigration0017CreatesExactlyTheAcceptedObjects(t *testing.T) {
	t.Parallel()
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Migrate(16); err != nil {
		t.Fatalf("Migrate(16) error = %v", err)
	}
	before := snapshotDatabaseState(t, conn)
	if err := migrator.Migrate(17); err != nil {
		t.Fatalf("Migrate(17) error = %v", err)
	}

	got := operationalKernelSchemaObjects(t, before, conn)
	want := map[string]struct{}{}
	for _, object := range []string{
		"table:review_records", "table:review_findings", "table:finding_dispositions",
		"table:judgment_concerns", "table:deviations", "table:decision_refs",
		"table:human_outcomes", "table:envelope_protected_decisions",
		"table:requirement_refs", "table:requirement_issue_links",
		"index:review_records_by_issue", "index:review_records_by_envelope",
		"index:review_findings_by_review", "index:finding_dispositions_one_terminal",
		"index:finding_dispositions_by_finding", "index:judgment_concerns_by_subject",
		"index:deviations_by_subject", "index:decision_refs_by_class",
		"index:human_outcomes_by_class", "index:requirement_refs_one_active",
		"index:requirement_refs_by_requirement", "index:requirement_issue_links_by_issue",
		"index:envelope_obligations_one_review_profile",
		"index:envelope_obligations_one_outcome_class",
		"index:envelope_obligations_one_gate_evidence", "index:evidence_records_one_receipt",
		"index:evidence_records_one_review", "index:evidence_records_one_outcome",
		"trigger:requirement_issue_links_ref_guard_insert",
	} {
		want[object] = struct{}{}
	}
	if gotKeys := schemaObjectKeys(got); !reflect.DeepEqual(gotKeys, schemaObjectKeysFromSet(want)) {
		t.Fatalf("0017 new schema object keys = %#v, want %#v", gotKeys, schemaObjectKeysFromSet(want))
	}

	for _, table := range []string{"evidence_records", "envelope_obligations"} {
		if _, exists := schemaObjectSQL(t, conn, "table", table); !exists {
			t.Fatalf("0017 did not retain rebuilt table %q", table)
		}
	}
	assertForeignKeyCheckClean(t, conn)
	assertIntegrityCheckClean(t, conn)
	assertMigration0017MatchesAcceptedDDL(t)
}

func TestMigration0017ObligationRegistrationDiscriminator(t *testing.T) {
	t.Parallel()
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(17); err != nil {
		t.Fatalf("Migrate(17) error = %v", err)
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable foreign keys: %v", err)
	}

	insert := func(id, obligationType string, registrationID any, registrationVersion any) error {
		_, err := conn.Exec(`
			INSERT INTO envelope_obligations
			(id, envelope_id, revision, obligation_type, registration_id, registration_version, obligation_json)
			VALUES (?, 'ENV-0017', 1, ?, ?, ?, '{}')
		`, id, obligationType, registrationID, registrationVersion)
		return err
	}
	for _, test := range []struct {
		name                string
		obligationType      string
		registrationID      any
		registrationVersion any
		wantRejected        bool
	}{
		{"gate without registration", "gate-evidence", nil, nil, true},
		{"review with registration", "review", "REG-0017", 1, true},
		{"gate with registration", "gate-evidence", "REG-0017", 1, false},
		{"human outcome without registration", "human-outcome", nil, nil, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := insert("OBL-"+strings.ReplaceAll(test.name, " ", "-"), test.obligationType, test.registrationID, test.registrationVersion)
			if (err != nil) != test.wantRejected {
				t.Fatalf("insert %s error = %v, want rejected=%t", test.name, err, test.wantRejected)
			}
		})
	}
}

func TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks(t *testing.T) {
	t.Parallel()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	output, err := exec.Command(
		"git", "-C", repoRoot, "diff", "--no-ext-diff", "--unified=0",
		"0195897a347dfe1872db0f7be3800c5c3ccdefcf", "--",
		"migrations/0013_p21_operational_kernel.up.sql",
		"migrations/0015_p22_recovery.up.sql",
	).Output()
	if err != nil {
		t.Fatalf("diff accepted baseline: %v", err)
	}
	const want = `diff --git a/migrations/0013_p21_operational_kernel.up.sql b/migrations/0013_p21_operational_kernel.up.sql
index 01899a8..d034272 100644
--- a/migrations/0013_p21_operational_kernel.up.sql
+++ b/migrations/0013_p21_operational_kernel.up.sql
@@ -155,3 +155,3 @@ CREATE TABLE envelope_obligations (
-  obligation_type TEXT NOT NULL CHECK (obligation_type = 'gate-evidence'),
-  registration_id TEXT NOT NULL,
-  registration_version INTEGER NOT NULL,
+  obligation_type TEXT NOT NULL CHECK (obligation_type IN ('gate-evidence','review','human-outcome')),
+  registration_id TEXT,
+  registration_version INTEGER,
@@ -158,0 +159,4 @@ CREATE TABLE envelope_obligations (
+  CHECK ((obligation_type = 'gate-evidence'
+          AND registration_id IS NOT NULL AND registration_version IS NOT NULL)
+      OR (obligation_type IN ('review','human-outcome')
+          AND registration_id IS NULL AND registration_version IS NULL)),
@@ -181 +185 @@ CREATE TABLE work_claims (
-  mode TEXT NOT NULL CHECK (mode = 'mutation'),
+  mode TEXT NOT NULL CHECK (mode IN ('mutation','review')),
diff --git a/migrations/0015_p22_recovery.up.sql b/migrations/0015_p22_recovery.up.sql
index 6bcffc4..715a384 100644
--- a/migrations/0015_p22_recovery.up.sql
+++ b/migrations/0015_p22_recovery.up.sql
@@ -61 +61,2 @@ CREATE TABLE stops (
-    'operator-stop'
+    'operator-stop','review-rework-exhausted','protected-decision-required',
+    'human-outcome-required'
@@ -87 +88,2 @@ CREATE TABLE attention_items (
-    'operator-stop'
+    'operator-stop','review-rework-exhausted','protected-decision-required',
+    'human-outcome-required'
`
	if string(output) != want {
		t.Fatalf("Phase 23 baseline diff =\n%s\nwant exactly seven accepted hunks:\n%s", output, want)
	}
}

func TestMigration0017PreservesReceiptEvidenceAndArtifactLinksInPlace(t *testing.T) {
	t.Parallel()
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(16); err != nil {
		t.Fatalf("Migrate(16) error = %v", err)
	}
	seedSchema16ReceiptEvidence(t, conn)

	var beforeEvidence string
	if err := conn.QueryRow(`SELECT id || '|' || obligation_id || '|' || receipt_id || '|' || evidence_type || '|' || producer || '|' || trust_source || '|' || target_repository_json || '|' || validity || '|' || created_at FROM evidence_records`).Scan(&beforeEvidence); err != nil {
		t.Fatalf("snapshot schema-16 evidence: %v", err)
	}
	if err := migrator.Migrate(17); err != nil {
		t.Fatalf("Migrate(17) error = %v", err)
	}

	var afterEvidence string
	var reviewID, outcomeID any
	if err := conn.QueryRow(`SELECT id || '|' || obligation_id || '|' || receipt_id || '|' || evidence_type || '|' || producer || '|' || trust_source || '|' || target_repository_json || '|' || validity || '|' || created_at, review_id, human_outcome_id FROM evidence_records`).Scan(&afterEvidence, &reviewID, &outcomeID); err != nil {
		t.Fatalf("snapshot schema-17 evidence: %v", err)
	}
	if afterEvidence != beforeEvidence || reviewID != nil || outcomeID != nil {
		t.Fatalf("migrated evidence = (%q, %v, %v), want (%q, <nil>, <nil>)", afterEvidence, reviewID, outcomeID, beforeEvidence)
	}
	var requiredMaterial int
	if err := conn.QueryRow(`SELECT required_material FROM evidence_artifact_links WHERE evidence_id = 'EVD-0017' AND artifact_id = 'ART-0017'`).Scan(&requiredMaterial); err != nil {
		t.Fatalf("query restored evidence artifact link: %v", err)
	}
	if requiredMaterial != 1 {
		t.Fatalf("restored link required_material = %d, want 1", requiredMaterial)
	}
	assertForeignKeyCheckClean(t, conn)
	assertSchemaMigrationCleanAt17(t, conn)
}

func assertMigration0017MatchesAcceptedDDL(t *testing.T) {
	t.Helper()
	for _, test := range []struct{ filename, heading, nextHeading string }{
		{"0017_p23_review_judgment.up.sql", "## 2. Migration 0017 Up", "## 3. Migration 0017 Down"},
		{"0017_p23_review_judgment.down.sql", "## 3. Migration 0017 Down", "## 4. Migration 0018 Up - V1 Judgment Removal"},
	} {
		t.Run(test.filename, func(t *testing.T) {
			got, err := os.ReadFile(filepath.Join(testMigrationsDir(t), test.filename))
			if err != nil {
				t.Fatalf("read migration: %v", err)
			}
			if want := acceptedDDLSection(t, test.heading, test.nextHeading); string(got) != want {
				t.Fatalf("%s differs from accepted Phase 23 DDL", test.filename)
			}
		})
	}
}

func acceptedDDLSection(t *testing.T, heading, nextHeading string) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	schemaPath := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "project-docs", "docs", "phases", "phase-23", "schema.md"))
	source, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read accepted schema: %v", err)
	}
	section := string(source)
	start := strings.Index(section, heading)
	end := strings.Index(section, nextHeading)
	if start < 0 || end < 0 || end <= start {
		t.Fatalf("locate accepted DDL section %q", heading)
	}
	section = section[start:end]
	codeStart := strings.Index(section, "```sql\n")
	codeEnd := strings.LastIndex(section, "\n```")
	if codeStart < 0 || codeEnd < 0 || codeEnd <= codeStart {
		t.Fatalf("locate accepted DDL code fence %q", heading)
	}
	return section[codeStart+len("```sql\n") : codeEnd+1]
}

func schemaObjectKeys(objects map[string]string) []string {
	keys := make([]string, 0, len(objects))
	for key := range objects {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func schemaObjectKeysFromSet(objects map[string]struct{}) []string {
	keys := make([]string, 0, len(objects))
	for key := range objects {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func schemaObjectSQL(t *testing.T, conn interface{ QueryRow(string, ...any) *sql.Row }, objectType, name string) (string, bool) {
	t.Helper()
	var statement string
	err := conn.QueryRow(`SELECT sql FROM sqlite_schema WHERE type = ? AND name = ?`, objectType, name).Scan(&statement)
	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		t.Fatalf("query schema object %s:%s: %v", objectType, name, err)
	}
	return statement, true
}

func assertIntegrityCheckClean(t *testing.T, conn interface{ QueryRow(string, ...any) *sql.Row }) {
	t.Helper()
	var result string
	if err := conn.QueryRow(`PRAGMA integrity_check`).Scan(&result); err != nil {
		t.Fatalf("PRAGMA integrity_check: %v", err)
	}
	if result != "ok" {
		t.Fatalf("PRAGMA integrity_check = %q, want ok", result)
	}
}

func assertSchemaMigrationCleanAt17(t *testing.T, conn interface{ QueryRow(string, ...any) *sql.Row }) {
	t.Helper()
	var version, dirty int
	if err := conn.QueryRow(`SELECT version, dirty FROM schema_migrations`).Scan(&version, &dirty); err != nil {
		t.Fatalf("schema_migrations: %v", err)
	}
	if version != 17 || dirty != 0 {
		t.Fatalf("schema_migrations = %d|%d, want 17|0", version, dirty)
	}
}

func seedSchema16ReceiptEvidence(t *testing.T, conn interface {
	Exec(string, ...any) (sql.Result, error)
}) {
	t.Helper()
	statements := []string{
		`PRAGMA foreign_keys = OFF`,
		`INSERT INTO issues (id, kind, title, workflow_id, status) VALUES ('ISS-0017', 'task', 'migration fixture', 'task-default', 'todo')`,
		`INSERT INTO kit_activations (id, kit_identity, backend_identity, phase, status, direction, created_at, updated_at) VALUES ('ACT-0017', 'kit', 'backend', 16, 'finalized', 'activate', '2026-09-20T00:00:00Z', '2026-09-20T00:00:00Z')`,
		`INSERT INTO activation_gate_definitions (activation_id, definition_id, definition_version, definition_digest, definition_json) VALUES ('ACT-0017', 'DEF-0017', 1, 'sha256:def', '{}')`,
		`INSERT INTO gate_registrations (id, version, activation_id, definition_id, definition_version, definition_digest, binding_digest, project_id, worktree_id, scope_json, policy_json, authorization_provenance, created_event_id, created_at) VALUES ('REG-0017', 1, 'ACT-0017', 'DEF-0017', 1, 'sha256:def', 'sha256:binding', 'project', 'worktree', '{}', '{}', 'fixture', 'EVT-registration', '2026-09-20T00:00:00Z')`,
		`INSERT INTO work_envelopes (id, current_revision, status, created_at) VALUES ('ENV-0017', 1, 'active', '2026-09-20T00:00:00Z')`,
		`INSERT INTO envelope_revisions (envelope_id, revision, outcome, non_goals_json, risk_result_json, repository_scope_json, status, created_actor_id, created_at, activated_at) VALUES ('ENV-0017', 1, 'outcome', '[]', '{}', '{}', 'active', 'actor', '2026-09-20T00:00:00Z', '2026-09-20T00:00:00Z')`,
		`INSERT INTO envelope_obligations (id, envelope_id, revision, obligation_type, registration_id, registration_version, obligation_json) VALUES ('OBL-0017', 'ENV-0017', 1, 'gate-evidence', 'REG-0017', 1, '{}')`,
		`INSERT INTO claims (id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at) VALUES ('CLM-0017', 'ISS-0017', 'ENV-0017', 1, 'actor', 'worktree', '{}', 'mutation', 1, 'active', '2026-09-20T00:00:00Z', '2026-09-21T00:00:00Z')`,
		`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, status, started_at) VALUES ('RUN-0017', 'ISS-0017', 'ENV-0017', 1, 'CLM-0017', 1, 'actor', '{}', 'active', '2026-09-20T00:00:00Z')`,
		`INSERT INTO execution_attempts (id, registration_id, registration_version, obligation_id, run_id, actor_id, start_repository_json, resolved_policy_digest, producer, started_at) VALUES ('ATT-0017', 'REG-0017', 1, 'OBL-0017', 'RUN-0017', 'actor', '{}', 'sha256:policy', 'fid', '2026-09-20T00:00:00Z')`,
		`INSERT INTO execution_receipts (id, attempt_id, trust_source, terminal_status, exit_code, timed_out, started_at, ended_at, terminal_repository_json, toolchain_json, output_excerpt, output_digest, output_bytes, output_truncated, redaction_applied, rejected_artifacts_json) VALUES ('RCP-0017', 'ATT-0017', 'fid-executed', 'succeeded', 0, 0, '2026-09-20T00:00:00Z', '2026-09-20T00:00:01Z', '{}', '{}', '', 'sha256:output', 0, 0, 0, '[]')`,
		`INSERT INTO referenced_artifacts (id, digest, media_type, size_bytes, retention_class, availability, location_policy, created_at) VALUES ('ART-0017', 'sha256:artifact', 'text/plain', 1, 'evidence', 'available', 'content-addressed', '2026-09-20T00:00:00Z')`,
		`INSERT INTO evidence_records (id, obligation_id, receipt_id, evidence_type, producer, trust_source, target_repository_json, validity, created_at) VALUES ('EVD-0017', 'OBL-0017', 'RCP-0017', 'receipt', 'fid', 'fid-executed', '{}', 'valid', '2026-09-20T00:00:00Z')`,
		`INSERT INTO evidence_artifact_links (evidence_id, artifact_id, required_material) VALUES ('EVD-0017', 'ART-0017', 1)`,
		`PRAGMA foreign_keys = ON`,
	}
	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("seed schema-16 receipt Evidence: %s: %v", statement, err)
		}
	}
}
