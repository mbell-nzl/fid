package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	migratesqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestMigrateUpCreatesPhaseTwoV01TablesAndConstraints(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	tests := []struct {
		table              string
		columns            []string
		checkSQLFragments  []string
		foreignKeyContains []string
	}{
		{
			table:              "epic_metadata",
			columns:            []string{"issue_id", "scheduling", "brief_md_path", "target_start", "target_end"},
			foreignKeyContains: []string{"issue_id->issues.id"},
		},
		{
			table:              "acceptance_criteria",
			columns:            []string{"id", "issue_id", "ordinal", "text"},
			foreignKeyContains: []string{"issue_id->issues.id"},
		},
		{
			table:              "ac_verifications",
			columns:            []string{"id", "ac_id", "verifier_role", "status", "completion_id", "verified_at"},
			foreignKeyContains: []string{"ac_id->acceptance_criteria.id", "completion_id->task_completions.id"},
		},
		{
			table:              "gate_runs",
			columns:            []string{"id", "issue_id", "attempt_number", "agent_run_id", "command", "exit_code", "ran_at"},
			checkSQLFragments:  []string{"CHECK (attempt_number > 0)"},
			foreignKeyContains: []string{"issue_id->issues.id", "agent_run_id->agent_runs.id"},
		},
		{
			table:              "task_completions",
			columns:            []string{"id", "issue_id", "attempt_number", "status", "gate_results", "human_intervention", "completed_at"},
			foreignKeyContains: []string{"issue_id->issues.id"},
		},
		{
			table:              "concerns",
			columns:            []string{"id", "issue_id", "raised_in_completion_id", "status", "escalated_to_esc_id", "triaged_in_event_id"},
			foreignKeyContains: []string{"issue_id->issues.id", "raised_in_completion_id->task_completions.id", "escalated_to_esc_id->escalations.id", "triaged_in_event_id->triage_events.id"},
		},
		{
			table:   "requirements",
			columns: []string{"id", "type", "summary", "source_doc", "source_section", "last_synced_hash", "status"},
		},
		{
			table:              "agent_runs",
			columns:            []string{"id", "issue_id", "completion_id", "role", "model", "token_breakdown_json", "provider_cost_usd"},
			foreignKeyContains: []string{"issue_id->issues.id", "completion_id->task_completions.id"},
		},
		{
			table:   "labels",
			columns: []string{"name", "colour", "description"},
		},
		{
			table:              "issue_labels",
			columns:            []string{"issue_id", "label"},
			foreignKeyContains: []string{"issue_id->issues.id", "label->labels.name"},
		},
		{
			table:   "links",
			columns: []string{"from_type", "from_id", "relation", "to_type", "to_id", "created_at"},
		},
		{
			table:   "claims",
			columns: []string{"ref_type", "ref_id", "claimed_by", "claimed_at", "ttl_seconds", "note"},
		},
		{
			table:   "activity_log",
			columns: []string{"id", "journal_event_id", "ref_type", "ref_id", "event", "details", "actor", "occurred_at"},
		},
		{
			table:              "decisions",
			columns:            []string{"id", "epic_id", "author_role", "summary", "body_md_path", "area", "status", "created_at"},
			foreignKeyContains: []string{"epic_id->issues.id"},
		},
		{
			table:              "escalations",
			columns:            []string{"id", "status", "trigger_class", "raised_by_role", "raised_at", "resolved_at", "resolution_dec"},
			foreignKeyContains: []string{"resolution_dec->decisions.id"},
		},
		{
			table:              "triage_events",
			columns:            []string{"id", "epic_id", "issue_id", "triager_role", "source", "source_ref_id", "verdict", "esc_id", "reasoning_md", "created_at"},
			foreignKeyContains: []string{"epic_id->issues.id", "issue_id->issues.id", "esc_id->escalations.id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.table, func(t *testing.T) {
			assertTableExists(t, conn, tt.table)
			assertTableColumns(t, conn, tt.table, tt.columns)
			if len(tt.checkSQLFragments) > 0 {
				assertSchemaContains(t, conn, tt.table, tt.checkSQLFragments...)
			}
			if len(tt.foreignKeyContains) > 0 {
				assertForeignKeysContain(t, conn, tt.table, tt.foreignKeyContains...)
			}
		})
	}
}

func TestMigration0002DownLeavesPhaseOneSchemaAndUpRestoresPhaseTwoTables(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrationsDir := testMigrationsDir(t)

	migrator, conn := newTestMigrator(t, dbPath, migrationsDir)
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Up() error = %v", err)
	}

	if err := migrator.Migrate(1); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migrate(1) error = %v", err)
	}

	assertTableExists(t, conn, "projection_meta")
	assertTableExists(t, conn, "roles")
	assertTableExists(t, conn, "issues")
	assertTableDoesNotExist(t, conn, "epic_metadata")
	assertTableDoesNotExist(t, conn, "concerns")
	assertTableDoesNotExist(t, conn, "triage_events")

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("second Up() error = %v", err)
	}

	assertTableExists(t, conn, "epic_metadata")
	assertTableExists(t, conn, "concerns")
	assertTableExists(t, conn, "triage_events")
}

func TestMigration0006DownRemovesSequenceSchedulingBackstopOnly(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrationsDir := testMigrationsDir(t)

	migrator, conn := newTestMigrator(t, dbPath, migrationsDir)
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Up() error = %v", err)
	}

	if err := execIssueInsert(conn, map[string]any{
		"id":          "FID-parent-epic",
		"kind":        "epic",
		"title":       "Parent epic",
		"workflow_id": "epic-default",
		"status":      "todo",
	}); err != nil {
		t.Fatalf("insert epic before down migration: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO epic_metadata (issue_id, scheduling) VALUES (?, ?)`, "FID-parent-epic", "sequential"); err != nil {
		t.Fatalf("insert epic_metadata before down migration: %v", err)
	}

	if err := execIssueInsert(conn, map[string]any{
		"id":          "FID-seq-invalid-before-down",
		"kind":        "task",
		"parent_id":   "FID-parent-epic",
		"title":       "Missing sequence",
		"workflow_id": "task-default",
		"status":      "todo",
	}); err == nil || !strings.Contains(err.Error(), "issues: task under sequential epic requires sequence") {
		t.Fatalf("sequence backstop before down = %v, want sequence trigger failure", err)
	}

	if err := migrator.Migrate(5); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migrate(5) error = %v", err)
	}

	if err := execIssueInsert(conn, map[string]any{
		"id":          "FID-child-task",
		"kind":        "task",
		"parent_id":   "FID-parent-epic",
		"title":       "Missing sequence allowed after down",
		"workflow_id": "task-default",
		"status":      "todo",
	}); err != nil {
		t.Fatalf("insert task after down migration: %v", err)
	}

	if err := execIssueInsert(conn, map[string]any{
		"id":          "FID-parent-task",
		"kind":        "task",
		"title":       "Task parent",
		"workflow_id": "task-default",
		"status":      "todo",
	}); err != nil {
		t.Fatalf("insert task parent after down migration: %v", err)
	}

	err := execIssueInsert(conn, map[string]any{
		"id":          "FID-invalid-child-task",
		"kind":        "task",
		"parent_id":   "FID-parent-task",
		"title":       "Task under task",
		"workflow_id": "task-default",
		"status":      "todo",
	})
	if err == nil || !strings.Contains(err.Error(), "issues: task parent_id must reference an epic or be NULL") {
		t.Fatalf("kind trigger after down = %v, want parent kind trigger failure", err)
	}
}

func TestDecideProjectionFreshnessRequiresPhase19SchemaVersion(t *testing.T) {
	t.Parallel()

	if ProjectionSchemaVersion != 5 {
		t.Fatalf("ProjectionSchemaVersion = %d, want %d", ProjectionSchemaVersion, 5)
	}

	stale := DecideProjectionFreshness(ProjectionFreshnessInput{
		MetaExists:                    true,
		ReplayStatus:                  ReplayStatusClean,
		StoredProjectionSchemaVersion: 4,
		CurrentProjectionSchema:       ProjectionSchemaVersion,
		CheckpointExists:              true,
		CheckpointMatches:             true,
	})
	if stale != ProjectionRebuild {
		t.Fatalf("stale freshness decision = %q, want %q", stale, ProjectionRebuild)
	}

	fresh := DecideProjectionFreshness(ProjectionFreshnessInput{
		MetaExists:                    true,
		ReplayStatus:                  ReplayStatusClean,
		StoredProjectionSchemaVersion: ProjectionSchemaVersion,
		CurrentProjectionSchema:       ProjectionSchemaVersion,
		CheckpointExists:              true,
		CheckpointMatches:             true,
	})
	if fresh != ProjectionFresh {
		t.Fatalf("freshness decision = %q, want %q", fresh, ProjectionFresh)
	}
}

func assertTableColumns(t *testing.T, conn *sql.DB, table string, want []string) {
	t.Helper()

	rows, err := conn.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		t.Fatalf("PRAGMA table_info(%q) error = %v", table, err)
	}
	defer closeRows(t, rows)

	got := make(map[string]struct{}, len(want))
	for rows.Next() {
		var (
			cid      int
			name     string
			dataType string
			notNull  int
			defaultV sql.NullString
			primaryK int
		)
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultV, &primaryK); err != nil {
			t.Fatalf("scan PRAGMA table_info(%q): %v", table, err)
		}
		got[name] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate PRAGMA table_info(%q): %v", table, err)
	}

	for _, column := range want {
		if _, ok := got[column]; !ok {
			t.Fatalf("table %q missing column %q; got columns %v", table, column, mapKeys(got))
		}
	}
}

func assertSchemaContains(t *testing.T, conn *sql.DB, table string, fragments ...string) {
	t.Helper()

	var sqlText string
	if err := conn.QueryRow(`
		SELECT sql
		FROM sqlite_master
		WHERE type = 'table' AND name = ?
	`, table).Scan(&sqlText); err != nil {
		t.Fatalf("query sqlite_master SQL for %q: %v", table, err)
	}

	normalized := strings.Join(strings.Fields(sqlText), " ")
	for _, fragment := range fragments {
		if !strings.Contains(normalized, strings.Join(strings.Fields(fragment), " ")) {
			t.Fatalf("schema for %q missing fragment %q; sql=%s", table, fragment, normalized)
		}
	}
}

func assertForeignKeysContain(t *testing.T, conn *sql.DB, table string, want ...string) {
	t.Helper()

	rows, err := conn.Query(`PRAGMA foreign_key_list(` + table + `)`)
	if err != nil {
		t.Fatalf("PRAGMA foreign_key_list(%q) error = %v", table, err)
	}
	defer closeRows(t, rows)

	var got []string
	for rows.Next() {
		var (
			id       int
			seq      int
			refTable string
			from     string
			to       string
			onUpdate string
			onDelete string
			match    string
		)
		if err := rows.Scan(&id, &seq, &refTable, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			t.Fatalf("scan PRAGMA foreign_key_list(%q): %v", table, err)
		}
		got = append(got, fmt.Sprintf("%s->%s.%s", from, refTable, to))
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate PRAGMA foreign_key_list(%q): %v", table, err)
	}

	for _, foreignKey := range want {
		if !containsString(got, foreignKey) {
			t.Fatalf("table %q missing foreign key %q; got %v", table, foreignKey, got)
		}
	}
}

func newTestMigrator(t *testing.T, dbPath, migrationsDir string) (*migrate.Migrate, *sql.DB) {
	t.Helper()

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	driver, err := migratesqlite.WithInstance(conn, &migratesqlite.Config{})
	if err != nil {
		closeDB(t, conn)
		t.Fatalf("WithInstance() error = %v", err)
	}

	absoluteMigrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		closeDB(t, conn)
		t.Fatalf("filepath.Abs(%q) error = %v", migrationsDir, err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://"+absoluteMigrationsDir, driverName, driver)
	if err != nil {
		closeDB(t, conn)
		t.Fatalf("NewWithDatabaseInstance() error = %v", err)
	}

	return migrator, conn
}

func closeMigrator(t *testing.T, migrator *migrate.Migrate, conn *sql.DB) {
	t.Helper()

	sourceErr, databaseErr := migrator.Close()
	if sourceErr != nil || databaseErr != nil {
		t.Fatalf("migrator.Close() errors = %v / %v", sourceErr, databaseErr)
	}
	closeDB(t, conn)
}

func closeRows(t *testing.T, rows *sql.Rows) {
	t.Helper()

	if err := rows.Close(); err != nil {
		t.Fatalf("Rows.Close() error = %v", err)
	}
}

func mapKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func execIssueInsert(conn *sql.DB, values map[string]any) error {
	_, err := conn.Exec(`
		INSERT INTO issues (id, kind, parent_id, sequence, title, workflow_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		values["id"],
		values["kind"],
		values["parent_id"],
		values["sequence"],
		values["title"],
		values["workflow_id"],
		values["status"],
		"2026-06-07T15:00:00Z",
		"2026-06-07T15:00:00Z",
	)
	return err
}
