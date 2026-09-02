package db

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/golang-migrate/migrate/v4"
)

func TestMigration0010SeedsEscalationDefaultExactlyAndPreservesV01Tuples(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Migrate(8); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migrate(8) error = %v", err)
	}
	defaultsBefore := workflowTableTuples(t, conn, []string{"epic-default", "task-default", "subtask-default"})

	if err := migrator.Migrate(13); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migrate(13) error = %v", err)
	}
	gotSeed := workflowTableTuples(t, conn, []string{"escalation-default"})
	wantSeed := map[string][]string{
		"workflows":                    {"escalation-default|Escalation default|escalation|1|<NULL>"},
		"workflow_states":              {"escalation-default|open|started|1", "escalation-default|resolved|completed|2"},
		"workflow_transitions":         {"escalation-default|open|resolved"},
		"workflow_transition_requires": {"escalation-default|open|resolved|linked_to_decision|<NULL>"},
	}
	if !reflect.DeepEqual(gotSeed, wantSeed) {
		t.Fatalf("escalation-default tuples = %#v, want %#v", gotSeed, wantSeed)
	}

	defaultsAfter := workflowTableTuples(t, conn, []string{"epic-default", "task-default", "subtask-default"})
	if !reflect.DeepEqual(defaultsAfter, defaultsBefore) {
		t.Fatalf("v0.1 default tuples after 0010 = %#v, want unchanged tuples %#v", defaultsAfter, defaultsBefore)
	}
}

func workflowTableTuples(t *testing.T, conn *sql.DB, workflowIDs []string) map[string][]string {
	t.Helper()

	if len(workflowIDs) == 0 {
		t.Fatal("workflow IDs are required")
	}
	placeholders := "?"
	for range workflowIDs[1:] {
		placeholders += ", ?"
	}
	queries := map[string]string{
		"workflows": `
			SELECT id || '|' || name || '|' || COALESCE(applies_to, '<NULL>') || '|' || active || '|' || COALESCE(created_at, '<NULL>')
			FROM workflows WHERE id IN (` + placeholders + `) ORDER BY id`,
		"workflow_states": `
			SELECT workflow_id || '|' || state || '|' || category || '|' || ordinal
			FROM workflow_states WHERE workflow_id IN (` + placeholders + `) ORDER BY workflow_id, ordinal, state`,
		"workflow_transitions": `
			SELECT workflow_id || '|' || from_state || '|' || to_state
			FROM workflow_transitions WHERE workflow_id IN (` + placeholders + `) ORDER BY workflow_id, from_state, to_state`,
		"workflow_transition_requires": `
			SELECT workflow_id || '|' || from_state || '|' || to_state || '|' || rule || '|' || COALESCE(params_json, '<NULL>')
			FROM workflow_transition_requires WHERE workflow_id IN (` + placeholders + `) ORDER BY workflow_id, from_state, to_state, rowid`,
	}

	result := make(map[string][]string, len(queries))
	for table, query := range queries {
		rows, err := conn.Query(query, stringsToAny(workflowIDs)...)
		if err != nil {
			t.Fatalf("query %s tuples: %v", table, err)
		}
		values := make([]string, 0)
		for rows.Next() {
			var value string
			if err := rows.Scan(&value); err != nil {
				_ = rows.Close()
				t.Fatalf("scan %s tuple: %v", table, err)
			}
			values = append(values, value)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			t.Fatalf("iterate %s tuples: %v", table, err)
		}
		closeRows(t, rows)
		result[table] = values
	}
	return result
}

func stringsToAny(values []string) []any {
	result := make([]any, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}
