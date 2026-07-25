package db

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestMigrateUpSeedsTaskDefaultWorkflowStates(t *testing.T) {
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

	want := map[string]string{
		"todo":         "unstarted",
		"in_progress":  "started",
		"dev_complete": "started",
		"qa_review":    "started",
		"qa_block":     "started",
		"failed_gate":  "started",
		"qa_pass":      "completed",
		"done":         "completed",
		"cancelled":    "cancelled",
	}
	got := make(map[string]string, len(want))
	rows, err := conn.Query(`
		SELECT state, category
		FROM workflow_states
		WHERE workflow_id = ?
	`, "task-default")
	if err != nil {
		t.Fatalf("query task-default workflow states: %v", err)
	}
	defer closeRows(t, rows)

	for rows.Next() {
		var state, category string
		if err := rows.Scan(&state, &category); err != nil {
			t.Fatalf("scan workflow state: %v", err)
		}
		got[state] = category
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate workflow states: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("task-default state count = %d, want %d; got %#v", len(got), len(want), got)
	}
	for state, wantCategory := range want {
		if got[state] != wantCategory {
			t.Fatalf("task-default category for %q = %q, want %q", state, got[state], wantCategory)
		}
	}
}

func TestMigration0007CreatesWorkflowSchema(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowTablesDB(t)
	defer closeDB(t, conn)

	tables := []struct {
		name    string
		columns []string
		primary []string
	}{
		{"workflows", []string{"id", "name", "applies_to", "active", "created_at"}, []string{"id"}},
		{"workflow_states", []string{"workflow_id", "state", "category", "ordinal"}, []string{"workflow_id", "state"}},
		{"workflow_transitions", []string{"workflow_id", "from_state", "to_state"}, []string{"workflow_id", "from_state", "to_state"}},
		{"workflow_transition_requires", []string{"workflow_id", "from_state", "to_state", "rule", "params_json"}, []string{}},
	}
	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			assertTableExists(t, conn, table.name)
			assertTableColumns(t, conn, table.name, table.columns)
			assertPrimaryKeyColumns(t, conn, table.name, table.primary)
		})
	}

	assertCompositeTransitionRequirementForeignKey(t, conn)
	assertNoPrimaryKeyOrUniqueConstraint(t, conn, "workflow_transition_requires")
}

func TestMigration0007SeedsFrozenV01WorkflowDefinitions(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowTablesDB(t)
	defer closeDB(t, conn)

	definitions := workflow.FrozenDefinitions()
	assertSeededWorkflows(t, conn, definitions)
	assertSeededStates(t, conn, definitions)
	assertSeededTransitions(t, conn, definitions)
	assertSeededRequirements(t, conn, definitions)
}

func TestMigration0007ReservesGateObligationShapeWithoutSeedingReservedRules(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowTablesDB(t)
	defer closeDB(t, conn)

	for _, rule := range []string{workflow.RuleRequiredGatesCurrent, workflow.RuleHasCompletion} {
		var count int
		if err := conn.QueryRow(`
			SELECT COUNT(*)
			FROM workflow_transition_requires
			WHERE rule = ?
		`, rule).Scan(&count); err != nil {
			t.Fatalf("count seeded %q rules: %v", rule, err)
		}
		if count != 0 {
			t.Fatalf("seeded %q rule count = %d, want 0", rule, count)
		}
	}

	const paramsJSON = `{"gates":[{"id":"gate-test","version":1}]}`
	if _, err := conn.Exec(`
		INSERT INTO workflow_transition_requires (workflow_id, from_state, to_state, rule, params_json)
		VALUES (?, ?, ?, ?, ?)
	`, "task-default", "todo", "in_progress", workflow.RuleRequiredGatesCurrent, paramsJSON); err != nil {
		t.Fatalf("insert gate-obligation params_json: %v", err)
	}

	var loaded string
	if err := conn.QueryRow(`
		SELECT params_json
		FROM workflow_transition_requires
		WHERE workflow_id = ? AND from_state = ? AND to_state = ? AND rule = ?
	`, "task-default", "todo", "in_progress", workflow.RuleRequiredGatesCurrent).Scan(&loaded); err != nil {
		t.Fatalf("load gate-obligation params_json: %v", err)
	}
	if loaded != paramsJSON {
		t.Fatalf("loaded params_json = %q, want %q", loaded, paramsJSON)
	}

	var parsed struct {
		Gates []struct {
			ID      string `json:"id"`
			Version int    `json:"version"`
		} `json:"gates"`
	}
	if err := json.Unmarshal([]byte(loaded), &parsed); err != nil {
		t.Fatalf("parse loaded params_json: %v", err)
	}
	if len(parsed.Gates) != 1 || parsed.Gates[0].ID != "gate-test" || parsed.Gates[0].Version != 1 {
		t.Fatalf("parsed gate obligation = %#v, want gate-test version 1", parsed.Gates)
	}

	roundTripped, err := json.Marshal(parsed)
	if err != nil {
		t.Fatalf("marshal parsed gate obligation: %v", err)
	}
	if string(roundTripped) != paramsJSON {
		t.Fatalf("round-tripped params_json = %q, want %q", roundTripped, paramsJSON)
	}
}

func TestMigration0007DownAndUpRestoresWorkflowTables(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Up(); err != nil {
		t.Fatalf("Up() error = %v", err)
	}
	if err := migrator.Migrate(6); err != nil {
		t.Fatalf("Migrate(6) error = %v", err)
	}
	for _, table := range workflowTableNames {
		assertTableDoesNotExist(t, conn, table)
	}

	if err := migrator.Up(); err != nil {
		t.Fatalf("second Up() error = %v", err)
	}
	for _, table := range workflowTableNames {
		assertTableExists(t, conn, table)
	}
}

var workflowTableNames = []string{
	"workflows",
	"workflow_states",
	"workflow_transitions",
	"workflow_transition_requires",
}

type workflowStateTuple struct {
	workflowID string
	state      string
	category   string
}

type workflowTransitionTuple struct {
	workflowID string
	fromState  string
	toState    string
}

type workflowRequirementTuple struct {
	workflowID string
	fromState  string
	toState    string
	rule       string
	paramsJSON string
}

func migratedWorkflowTablesDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}
	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return conn
}

func assertPrimaryKeyColumns(t *testing.T, conn *sql.DB, table string, want []string) {
	t.Helper()

	rows, err := conn.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		t.Fatalf("PRAGMA table_info(%q): %v", table, err)
	}
	defer closeRows(t, rows)

	gotByPosition := make(map[int]string)
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, dataType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatalf("scan table info for %q: %v", table, err)
		}
		if primaryKey > 0 {
			gotByPosition[primaryKey] = name
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate table info for %q: %v", table, err)
	}

	got := make([]string, 0, len(gotByPosition))
	for position := 1; position <= len(gotByPosition); position++ {
		got = append(got, gotByPosition[position])
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("primary key columns for %q = %#v, want %#v", table, got, want)
	}
}

func assertCompositeTransitionRequirementForeignKey(t *testing.T, conn *sql.DB) {
	t.Helper()

	rows, err := conn.Query(`PRAGMA foreign_key_list(workflow_transition_requires)`)
	if err != nil {
		t.Fatalf("PRAGMA foreign_key_list(workflow_transition_requires): %v", err)
	}
	defer closeRows(t, rows)

	type foreignKeyColumn struct{ from, to string }
	columns := make(map[int][]foreignKeyColumn)
	tables := make(map[int]string)
	for rows.Next() {
		var id, sequence int
		var table, from, to, onUpdate, onDelete, match string
		if err := rows.Scan(&id, &sequence, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			t.Fatalf("scan workflow requirement foreign key: %v", err)
		}
		tables[id] = table
		columns[id] = append(columns[id], foreignKeyColumn{from: from, to: to})
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate workflow requirement foreign keys: %v", err)
	}

	want := []foreignKeyColumn{{"workflow_id", "workflow_id"}, {"from_state", "from_state"}, {"to_state", "to_state"}}
	for id, got := range columns {
		if tables[id] == "workflow_transitions" && reflect.DeepEqual(got, want) {
			return
		}
	}
	t.Fatalf("workflow_transition_requires foreign keys = %#v to %#v, want composite workflow_transitions key %#v", columns, tables, want)
}

func assertNoPrimaryKeyOrUniqueConstraint(t *testing.T, conn *sql.DB, table string) {
	t.Helper()

	var schemaSQL string
	if err := conn.QueryRow(`SELECT sql FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&schemaSQL); err != nil {
		t.Fatalf("load schema SQL for %q: %v", table, err)
	}
	normalized := strings.ToUpper(strings.Join(strings.Fields(schemaSQL), " "))
	if strings.Contains(normalized, "PRIMARY KEY") || strings.Contains(normalized, "UNIQUE") {
		t.Fatalf("schema for %q unexpectedly has primary/unique constraint: %s", table, normalized)
	}
}

func assertSeededWorkflows(t *testing.T, conn *sql.DB, definitions map[string]workflow.Definition) {
	t.Helper()

	rows, err := conn.Query(`SELECT id, applies_to, active FROM workflows ORDER BY id`)
	if err != nil {
		t.Fatalf("query workflows: %v", err)
	}
	defer closeRows(t, rows)

	got := make(map[string]struct {
		appliesTo string
		active    bool
	})
	for rows.Next() {
		var id, appliesTo string
		var active bool
		if err := rows.Scan(&id, &appliesTo, &active); err != nil {
			t.Fatalf("scan workflow: %v", err)
		}
		got[id] = struct {
			appliesTo string
			active    bool
		}{appliesTo, active}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate workflows: %v", err)
	}

	wantAppliesTo := map[string]string{"epic-default": "epic", "task-default": "task", "subtask-default": "subtask"}
	if len(got) != len(definitions) {
		t.Fatalf("seeded workflow count = %d, want %d; got %#v", len(got), len(definitions), got)
	}
	for workflowID := range definitions {
		row, ok := got[workflowID]
		if !ok || row.appliesTo != wantAppliesTo[workflowID] || !row.active {
			t.Fatalf("seeded workflow %q = %#v, want applies_to %q and active", workflowID, row, wantAppliesTo[workflowID])
		}
	}
}

func assertSeededStates(t *testing.T, conn *sql.DB, definitions map[string]workflow.Definition) {
	t.Helper()

	rows, err := conn.Query(`SELECT workflow_id, state, category FROM workflow_states ORDER BY workflow_id, state`)
	if err != nil {
		t.Fatalf("query workflow states: %v", err)
	}
	defer closeRows(t, rows)

	var got []workflowStateTuple
	for rows.Next() {
		var tuple workflowStateTuple
		if err := rows.Scan(&tuple.workflowID, &tuple.state, &tuple.category); err != nil {
			t.Fatalf("scan workflow state: %v", err)
		}
		got = append(got, tuple)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate workflow states: %v", err)
	}

	var want []workflowStateTuple
	for workflowID, definition := range definitions {
		for state, category := range definition.States {
			want = append(want, workflowStateTuple{workflowID, state, category})
		}
	}
	sort.Slice(want, func(i, j int) bool { return stateTupleKey(want[i]) < stateTupleKey(want[j]) })
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("seeded workflow states = %#v, want %#v", got, want)
	}
}

func assertSeededTransitions(t *testing.T, conn *sql.DB, definitions map[string]workflow.Definition) {
	t.Helper()

	rows, err := conn.Query(`SELECT workflow_id, from_state, to_state FROM workflow_transitions ORDER BY workflow_id, from_state, to_state`)
	if err != nil {
		t.Fatalf("query workflow transitions: %v", err)
	}
	defer closeRows(t, rows)

	var got []workflowTransitionTuple
	for rows.Next() {
		var tuple workflowTransitionTuple
		if err := rows.Scan(&tuple.workflowID, &tuple.fromState, &tuple.toState); err != nil {
			t.Fatalf("scan workflow transition: %v", err)
		}
		got = append(got, tuple)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate workflow transitions: %v", err)
	}

	var want []workflowTransitionTuple
	for workflowID, definition := range definitions {
		for fromState, destinations := range definition.Transitions {
			for toState := range destinations {
				want = append(want, workflowTransitionTuple{workflowID, fromState, toState})
			}
		}
		for state := range definition.States {
			if state != "cancelled" {
				want = append(want, workflowTransitionTuple{workflowID, state, "cancelled"})
			}
		}
	}
	sort.Slice(want, func(i, j int) bool { return transitionTupleKey(want[i]) < transitionTupleKey(want[j]) })
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("seeded workflow transitions = %#v, want %#v", got, want)
	}
}

func assertSeededRequirements(t *testing.T, conn *sql.DB, definitions map[string]workflow.Definition) {
	t.Helper()

	rows, err := conn.Query(`
		SELECT workflow_id, from_state, to_state, rule, COALESCE(params_json, '')
		FROM workflow_transition_requires
		ORDER BY workflow_id, from_state, to_state, rule
	`)
	if err != nil {
		t.Fatalf("query workflow requirements: %v", err)
	}
	defer closeRows(t, rows)

	var got []workflowRequirementTuple
	for rows.Next() {
		var tuple workflowRequirementTuple
		if err := rows.Scan(&tuple.workflowID, &tuple.fromState, &tuple.toState, &tuple.rule, &tuple.paramsJSON); err != nil {
			t.Fatalf("scan workflow requirement: %v", err)
		}
		got = append(got, tuple)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate workflow requirements: %v", err)
	}

	var want []workflowRequirementTuple
	for workflowID, definition := range definitions {
		for fromState, destinations := range definition.Transitions {
			for toState, rules := range destinations {
				for _, rule := range rules {
					paramsJSON := ""
					if len(rule.RequiredRoles) > 0 {
						encoded, err := json.Marshal(struct {
							Roles []string `json:"roles"`
						}{rule.RequiredRoles})
						if err != nil {
							t.Fatalf("marshal roles for %q: %v", rule.Name, err)
						}
						paramsJSON = string(encoded)
					}
					want = append(want, workflowRequirementTuple{workflowID, fromState, toState, rule.Name, paramsJSON})
				}
			}
		}
	}
	sort.Slice(want, func(i, j int) bool { return requirementTupleKey(want[i]) < requirementTupleKey(want[j]) })
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("seeded workflow requirements = %#v, want %#v", got, want)
	}
}

func stateTupleKey(tuple workflowStateTuple) string {
	return tuple.workflowID + "\x00" + tuple.state
}

func transitionTupleKey(tuple workflowTransitionTuple) string {
	return tuple.workflowID + "\x00" + tuple.fromState + "\x00" + tuple.toState
}

func requirementTupleKey(tuple workflowRequirementTuple) string {
	return tuple.workflowID + "\x00" + tuple.fromState + "\x00" + tuple.toState + "\x00" + tuple.rule
}
