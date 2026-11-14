package db

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"unicode"
)

func TestLatestSchemaContainsOnlyRetainedTables(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Up(); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"acceptance_criteria", "activity_log", "agent_runs", "check_files", "check_processes",
		"check_requests", "check_results", "check_starts", "epic_metadata", "handoffs", "holds",
		"idempotency_results", "issue_labels", "issues", "labels", "links", "person_decisions",
		"policy_snapshots", "projection_meta", "questions", "referenced_artifacts", "schema_migrations",
		"triage_events", "workflow_states", "workflow_transition_requires", "workflow_transitions", "workflows",
	}
	rows, err := conn.Query(`SELECT name FROM sqlite_schema WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		t.Fatal(err)
	}
	defer closeRows(t, rows)
	var got []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		got = append(got, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("latest tables = %v, want only %v", got, want)
	}
	assertIssueColumnsAbsent(t, conn, "issues_enriched", "concerns_raised", "open_concerns_count")
	assertIssueColumnsPresent(t, conn, "issues_enriched", "triage_status", "lane", "revision")
}

func TestMigration0027RoundTripPreservesRetainedState(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(26); err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,lane) VALUES('FID-keep','task','Keep','task-default','in_progress','assured')`,
		`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at,ended_at,end_reason,kept) VALUES('HLD-keep','FID-keep','one','/w/one','agent:builder','2026-11-15T07:00:00Z','2026-11-15T08:00:00Z','2026-11-15T07:05:00Z','review',1)`,
		`INSERT INTO referenced_artifacts(id,digest,media_type,size_bytes,retention_class,availability,location_policy,created_at) VALUES('ART-keep','sha256:keep','text/plain',1,'evidence','available','content-addressed','2026-11-15T07:00:00Z')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	before := snapshotDatabaseState(t, conn)
	if err := migrator.Migrate(27); err != nil {
		t.Fatal(err)
	}
	after := snapshotDatabaseState(t, conn)
	assertNoDroppedTableReferences(t, before, after)
	for table, rows := range after.rows {
		if !reflect.DeepEqual(rows, before.rows[table]) {
			t.Fatalf("%s changed: %v, want %v", table, rows, before.rows[table])
		}
	}
	var violations int
	if err := conn.QueryRow(`SELECT count(*) FROM pragma_foreign_key_check`).Scan(&violations); err != nil || violations != 0 {
		t.Fatalf("foreign key violations = %d: %v", violations, err)
	}
	if err := migrator.Migrate(26); err != nil {
		t.Fatal(err)
	}
	if got := snapshotDatabaseState(t, conn); !reflect.DeepEqual(got, before) {
		t.Fatal("downgrade did not restore schema 26 and retained rows exactly")
	}
	if err := migrator.Migrate(27); err != nil {
		t.Fatal(err)
	}
	if got := snapshotDatabaseState(t, conn); !reflect.DeepEqual(got, after) {
		t.Fatal("second upgrade differs from first")
	}
}

func assertNoDroppedTableReferences(t *testing.T, before, after databaseState) {
	t.Helper()
	for table := range before.rows {
		if _, retained := after.rows[table]; retained {
			continue
		}
		for _, object := range after.schema {
			for _, word := range strings.FieldsFunc(object, func(r rune) bool { return !unicode.IsLetter(r) && r != '_' }) {
				if strings.EqualFold(word, table) {
					t.Fatalf("retained schema references dropped table %s: %s", table, object)
				}
			}
		}
	}
}

func TestMigration0027DropsPopulatedKernelAndKeepsSharedArtifacts(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(16); err != nil {
		t.Fatal(err)
	}
	seedSchema16ReceiptEvidence(t, conn)
	if err := migrator.Migrate(27); err != nil {
		t.Fatal(err)
	}
	assertTableDoesNotExist(t, conn, "evidence_records")
	assertTableDoesNotExist(t, conn, "claims")
	var id string
	if err := conn.QueryRow(`SELECT id FROM referenced_artifacts`).Scan(&id); err != nil || id != "ART-0017" {
		t.Fatalf("shared Artifact = %s, %v", id, err)
	}
	assertIntegrityCheckClean(t, conn)
	assertCurrentActivityGuard(t, conn)
}

func assertCurrentActivityGuard(t *testing.T, conn *sql.DB) {
	t.Helper()
	for _, refType := range []string{"issue", "requirement", "decision"} {
		_, err := conn.Exec(`INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event)
		 VALUES(?, 'EVT-guard', ?, 'ISS-0017', 'updated')`, "ACT-"+refType, refType)
		if refType == "issue" {
			if err != nil {
				t.Fatal(err)
			}
		} else if err == nil || !strings.Contains(err.Error(), "activity_log: ref_type/ref_id") {
			t.Fatalf("activity %s error = %v, want reference refusal", refType, err)
		}
	}
}
