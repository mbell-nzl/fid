package db

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFinalSchemaHasNoEscalationsOrEscalationReferences(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := MigrateUp(path, testMigrationsDir(t)); err != nil {
		t.Fatal(err)
	}
	conn, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer closeDB(t, conn)
	assertTableDoesNotExist(t, conn, "escalations")
	for _, table := range []string{"triage_events", "concerns", "activity_log", "runs", "stops", "attention_items"} {
		assertTableExists(t, conn, table)
	}
	rows, err := conn.Query(`SELECT name, COALESCE(sql, '') FROM sqlite_schema ORDER BY name`)
	if err != nil {
		t.Fatal(err)
	}
	defer closeRows(t, rows)
	for rows.Next() {
		var name, ddl string
		if err := rows.Scan(&name, &ddl); err != nil {
			t.Fatal(err)
		}
		for _, retired := range []string{"escalations", "esc_id", "escalated_to_esc_id", "'escalation'"} {
			if strings.Contains(ddl, retired) {
				t.Errorf("latest schema %s retains %s: %s", name, retired, ddl)
			}
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if err := rows.Close(); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"workflows", "workflow_states", "workflow_transitions", "workflow_transition_requires"} {
		column := "workflow_id"
		if table == "workflows" {
			column = "id"
		}
		var count int
		if err := conn.QueryRow("SELECT COUNT(*) FROM " + table + " WHERE " + column + " = 'escalation-default'").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Errorf("%s retains %d retired seed rows", table, count)
		}
	}
	if ActiveMigrationVersion != 20 {
		t.Fatalf("active migration = %d, want 20", ActiveMigrationVersion)
	}
	if _, err := conn.Exec(`INSERT INTO activity_log (id, journal_event_id, ref_type, ref_id, event, occurred_at) VALUES ('ACT-retired', 'EVT-retired', 'escalation', 'ESC-retired', 'raised', '2026-09-08T00:00:00Z')`); err == nil || !strings.Contains(err.Error(), "activity_log: ref_type/ref_id") {
		t.Fatalf("latest activity accepts retired reference or wrong rejection: %v", err)
	}
}

func TestFinalSchemaResetsThroughZeroWithoutRetiredAuthority(t *testing.T) {
	t.Parallel()
	path, migrations := filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t)
	if err := MigrateUp(path, migrations); err != nil {
		t.Fatal(err)
	}
	if err := MigrateDown(path, migrations); err != nil {
		t.Fatal(err)
	}
	if err := MigrateUp(path, migrations); err != nil {
		t.Fatal(err)
	}
	conn, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer closeDB(t, conn)
	assertTableDoesNotExist(t, conn, "escalations")
	assertTableExists(t, conn, "triage_events")
	assertTableExists(t, conn, "concerns")
	assertTableExists(t, conn, "stops")
}
