package workflow

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/db"
)

func TestRegistryEqualsV01Definitions(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)

	registry, err := Load(conn)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	frozen := FrozenDefinitions()
	if len(registry.workflows) != len(frozen) {
		t.Fatalf("loaded workflow count = %d, want %d frozen defaults", len(registry.workflows), len(frozen))
	}
	for workflowID, want := range frozen {
		got, ok := registry.Workflow(workflowID)
		if !ok {
			t.Fatalf("loaded registry is missing frozen workflow %q", workflowID)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("loaded frozen workflow %q = %#v, want %#v", workflowID, got, want)
		}
	}
}

func TestLoadReleasesWorkflowRowsBeforeNestedQueries(t *testing.T) {
	t.Parallel()
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	conn.SetMaxOpenConns(1)

	result := make(chan error, 1)
	go func() {
		_, err := Load(conn)
		result <- err
	}()

	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Load() did not return with one open database connection")
	}
}

func TestLoadSurfacesReservedRulesAsUnimplementedWithoutWriting(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)

	for _, rule := range []struct {
		name       string
		paramsJSON string
	}{
		{name: RuleRequiredGatesCurrent, paramsJSON: `{"gates":[{"id":"gate-test","version":1}]}`},
		{name: RuleHasCompletion, paramsJSON: `{"completion":"required"}`},
	} {
		if _, err := conn.Exec(`
			INSERT INTO workflow_transition_requires (workflow_id, from_state, to_state, rule, params_json)
			VALUES (?, ?, ?, ?, ?)
		`, "task-default", "todo", "in_progress", rule.name, rule.paramsJSON); err != nil {
			t.Fatalf("insert %q rule: %v", rule.name, err)
		}
	}
	if _, err := conn.Exec(`PRAGMA query_only = ON`); err != nil {
		t.Fatalf("set query_only: %v", err)
	}

	registry, err := Load(conn)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	rules, ok := registry.RulesForEdge("task-default", "todo", "in_progress")
	if !ok {
		t.Fatal("loaded task-default todo -> in_progress edge is missing")
	}
	want := []RuleReference{
		{Name: RuleRequiredGatesCurrent, Unimplemented: true},
		{Name: RuleHasCompletion, Unimplemented: true},
	}
	if !reflect.DeepEqual(rules, want) {
		t.Fatalf("reserved rules = %#v, want %#v", rules, want)
	}
}

func TestLoadImplementsProjectDecisionRuleWithoutRetiredSeed(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)

	var retired int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM workflows WHERE id = 'escalation-default'`).Scan(&retired); err != nil || retired != 0 {
		t.Fatalf("retired seed count = %d, error = %v", retired, err)
	}
	if _, err := conn.Exec(`INSERT INTO workflow_transition_requires (workflow_id, from_state, to_state, rule) VALUES ('task-default', 'todo', 'in_progress', 'linked_to_decision')`); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Exec(`PRAGMA query_only = ON`); err != nil {
		t.Fatalf("set query_only: %v", err)
	}

	registry, err := Load(conn)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	rules, ok := registry.RulesForEdge("task-default", "todo", "in_progress")
	if !ok {
		t.Fatal("loaded project Decision rule edge is missing")
	}
	want := []RuleReference{{Name: RuleLinkedToDecision}}
	if !reflect.DeepEqual(rules, want) {
		t.Fatalf("project Decision rule = %#v, want %#v", rules, want)
	}
}

func migratedWorkflowRegistryDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := db.MigrateUp(dbPath, workflowMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}
	conn, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return conn
}

func closeRegistryDB(t *testing.T, conn *sql.DB) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func workflowMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}
