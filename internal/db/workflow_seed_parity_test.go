package db_test

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestMigration0007SeedsFrozenV01WorkflowDefinitions(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowSeedDB(t)
	defer closeWorkflowSeedDB(t, conn)

	definitions := workflow.FrozenDefinitions()
	assertSeededWorkflows(t, conn, definitions)
	assertSeededStates(t, conn, definitions)
	assertSeededTransitions(t, conn, definitions)
	assertSeededRequirements(t, conn, definitions)
}

func TestMigration0007ReservesGateObligationShapeWithoutSeedingReservedRules(t *testing.T) {
	t.Parallel()

	conn := migratedWorkflowSeedDB(t)
	defer closeWorkflowSeedDB(t, conn)

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

func migratedWorkflowSeedDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testWorkflowSeedMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}
	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return conn
}

func testWorkflowSeedMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func closeWorkflowSeedDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func closeWorkflowSeedRows(t *testing.T, rows *sql.Rows) {
	t.Helper()

	if err := rows.Close(); err != nil {
		t.Fatalf("Rows.Close() error = %v", err)
	}
}

func assertSeededWorkflows(t *testing.T, conn *sql.DB, definitions map[string]workflow.Definition) {
	t.Helper()

	rows, err := conn.Query(`SELECT id, applies_to, active FROM workflows ORDER BY id`)
	if err != nil {
		t.Fatalf("query workflows: %v", err)
	}
	defer closeWorkflowSeedRows(t, rows)

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
		if _, ok := definitions[id]; !ok {
			continue
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
	defer closeWorkflowSeedRows(t, rows)

	var got []workflowStateTuple
	for rows.Next() {
		var tuple workflowStateTuple
		if err := rows.Scan(&tuple.workflowID, &tuple.state, &tuple.category); err != nil {
			t.Fatalf("scan workflow state: %v", err)
		}
		if _, ok := definitions[tuple.workflowID]; !ok {
			continue
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
	defer closeWorkflowSeedRows(t, rows)

	var got []workflowTransitionTuple
	for rows.Next() {
		var tuple workflowTransitionTuple
		if err := rows.Scan(&tuple.workflowID, &tuple.fromState, &tuple.toState); err != nil {
			t.Fatalf("scan workflow transition: %v", err)
		}
		if _, ok := definitions[tuple.workflowID]; !ok {
			continue
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
	defer closeWorkflowSeedRows(t, rows)

	var got []workflowRequirementTuple
	for rows.Next() {
		var tuple workflowRequirementTuple
		if err := rows.Scan(&tuple.workflowID, &tuple.fromState, &tuple.toState, &tuple.rule, &tuple.paramsJSON); err != nil {
			t.Fatalf("scan workflow requirement: %v", err)
		}
		if _, ok := definitions[tuple.workflowID]; !ok {
			continue
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
