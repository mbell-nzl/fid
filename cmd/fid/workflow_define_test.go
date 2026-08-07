package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestWorkflowDefineProjectsDeclaredStatesTransitionsAndRules(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	journalPath := mustJournalPath(t, workspace)
	journalBefore := mustReadCommandFile(t, journalPath)
	conn := openWorkflowCommandDB(t, workspace)
	defer closeCommandDB(t, conn)
	activityBefore := countActivityLogRows(t, conn)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"workflow", "define",
		"--name", "dev-flow-strict",
		"--applies-to", "task",
		"--states", "todo,in_progress,done,cancelled",
		"--categories", "unstarted,started,completed,cancelled",
		"--transitions", "todo->in_progress,in_progress->done",
		"--require", `in_progress->done:all_acs_verified{"roles":["developer"]}`,
		"--require", "in_progress->done:has_completion_report",
		"--by", "tl",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if want := "defined workflow dev-flow-strict for task\n  states: 4\n  transitions: 5\n  requirements: 2\n"; stderr != want {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}

	journalAfter := mustReadCommandFile(t, journalPath)
	if got, want := len(splitJournalLines(journalAfter)), len(splitJournalLines(journalBefore))+1; got != want {
		t.Fatalf("journal lines = %d, want %d", got, want)
	}
	event := mustReadLastJournalEvent(t, journalPath)
	if event.EventType != "workflow.defined" {
		t.Fatalf("event type = %q, want workflow.defined", event.EventType)
	}
	if event.EventSchema != journal.CurrentEventSchema {
		t.Fatalf("event schema = %d, want %d", event.EventSchema, journal.CurrentEventSchema)
	}
	if got, want := len(event.Effects), 12; got != want {
		t.Fatalf("event effects = %d, want %d", got, want)
	}

	assertWorkflowDefinitionProjection(t, conn)
	registry, err := workflow.Load(conn)
	if err != nil {
		t.Fatalf("workflow.Load() error = %v", err)
	}
	if _, ok := registry.Workflow("dev-flow-strict"); !ok {
		t.Fatal("workflow.Load() registry missing dev-flow-strict")
	}
	if activityAfter := countActivityLogRows(t, conn); activityAfter != activityBefore {
		t.Fatalf("activity_log rows = %d, want unchanged %d", activityAfter, activityBefore)
	}
}

func TestWorkflowDefineAuthorizesBeforeParsingAndDoesNotAppendOnDenial(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	journalPath := mustJournalPath(t, workspace)
	before := mustReadCommandFile(t, journalPath)

	stdout, stderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"workflow", "define",
		"--name", "denied-flow",
		"--applies-to", "not-a-kind",
		"--states", "todo,cancelled",
		"--categories", "unstarted,cancelled",
		"--transitions", "todo->cancelled",
		"--by", "developer",
	)
	if err == nil {
		t.Fatal("workflow define error = nil, want permission denial")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if want := "error: role \"developer\" not permitted for action \"workflow.define\"\n  check `fid policy show` for permitted role/action combinations\n"; stderr != want {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
	if after := mustReadCommandFile(t, journalPath); !bytes.Equal(after, before) {
		t.Fatal("denied workflow define changed journal bytes")
	}
}

func TestWorkflowDefinePreservesJSONRequirementParametersAfterTheEdgeColon(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"workflow", "define",
		"--name", "review-flow",
		"--applies-to", "qa_review",
		"--states", "open,resolved",
		"--categories", "unstarted,completed",
		"--transitions", "open->resolved",
		"--require", `open->resolved:has_completion_report["review:done"]`,
		"--by", "tl",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if want := "defined workflow review-flow for qa_review\n  states: 2\n  transitions: 1\n  requirements: 1\n"; stderr != want {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}

	conn := openWorkflowCommandDB(t, workspace)
	defer closeCommandDB(t, conn)
	var params string
	if err := conn.QueryRow(`SELECT params_json FROM workflow_transition_requires WHERE workflow_id = 'review-flow'`).Scan(&params); err != nil {
		t.Fatalf("query requirement params: %v", err)
	}
	if want := `["review:done"]`; params != want {
		t.Fatalf("params_json = %q, want %q", params, want)
	}
}

func TestWorkflowDefineRejectsWithoutMutatingJournalOrWorkflowTables(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{"invalid applies to", []string{"--name", "reject-flow", "--applies-to", "invalid-kind", "--states", "todo,cancelled", "--categories", "unstarted,cancelled", "--transitions", "todo->cancelled"}, `"invalid-kind"`},
		{"invalid slug", []string{"--name", "Bad Flow", "--applies-to", "task", "--states", "todo,cancelled", "--categories", "unstarted,cancelled", "--transitions", "todo->cancelled"}, `"Bad Flow"`},
		{"empty states", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "", "--categories", "", "--transitions", ""}, "must not be empty"},
		{"category count mismatch", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,cancelled", "--transitions", "todo->done"}, "3 states, 2 categories"},
		{"unknown category", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,cancelled", "--categories", "waiting,cancelled", "--transitions", "todo->cancelled"}, `"waiting"`},
		{"invalid state", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "Todo,cancelled", "--categories", "unstarted,cancelled", "--transitions", "Todo->cancelled"}, `"Todo"`},
		{"duplicate state", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,todo,cancelled", "--categories", "unstarted,started,cancelled", "--transitions", "todo->cancelled"}, `"todo"`},
		{"unknown transition state", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,cancelled", "--categories", "unstarted,cancelled", "--transitions", "todo->missing"}, `"todo" -> "missing"`},
		{"duplicate transition", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done,todo->done"}, `"todo" -> "done"`},
		{"unknown rule", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--require", "todo->done:invented_rule"}, `"invented_rule"`},
		{"reserved rule", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--require", "todo->done:required_gates_current"}, "v0.2b"},
		{"invalid rule json", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--require", "todo->done:has_completion_report{"}, "invalid JSON"},
		{"all acs requires roles", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--require", `todo->done:all_acs_verified{"roles":[]}`}, "non-empty roles"},
		{"missing cancelled state", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,done", "--categories", "unstarted,completed", "--transitions", "todo->done"}, "PRD-WF-013"},
		{"multiple cancelled states", []string{"--name", "reject-flow", "--applies-to", "task", "--states", "todo,cancelled,abandoned", "--categories", "unstarted,cancelled,cancelled", "--transitions", "todo->cancelled"}, `"cancelled" and "abandoned"`},
	}

	seenErrors := map[string]bool{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := t.TempDir()
			runInitCommand(t, workspace, "init")
			journalPath := mustJournalPath(t, workspace)
			beforeJournal := mustReadCommandFile(t, journalPath)
			conn := openWorkflowCommandDB(t, workspace)
			beforeTables := workflowCommandTableSnapshot(t, conn)
			closeCommandDB(t, conn)

			args := append([]string{"workflow", "define"}, tt.args...)
			args = append(args, "--by", "tl")
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, args...)
			if err == nil {
				t.Fatal("workflow define error = nil, want rejection")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, tt.wantErr) {
				t.Fatalf("stderr = %q, want %q", stderr, tt.wantErr)
			}
			if seenErrors[stderr] {
				t.Fatalf("stderr = %q, want distinct rejection message", stderr)
			}
			seenErrors[stderr] = true
			if afterJournal := mustReadCommandFile(t, journalPath); !bytes.Equal(afterJournal, beforeJournal) {
				t.Fatal("rejected workflow define changed journal bytes")
			}
			conn = openWorkflowCommandDB(t, workspace)
			afterTables := workflowCommandTableSnapshot(t, conn)
			closeCommandDB(t, conn)
			if !reflect.DeepEqual(afterTables, beforeTables) {
				t.Fatalf("workflow tables after rejection = %#v, want %#v", afterTables, beforeTables)
			}
		})
	}
}

func TestWorkflowDefineNormalizesArrowsAndAcceptsAppliesToVocabulary(t *testing.T) {
	for _, arrow := range []string{"->", "-->", "→", "⟶"} {
		t.Run("arrow "+arrow, func(t *testing.T) {
			workspace := t.TempDir()
			runInitCommand(t, workspace, "init")
			runFIDCommand(t, workspace, "workflow", "define", "--name", "arrow-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", " todo "+arrow+" done ", "--by", "tl")
			conn := openWorkflowCommandDB(t, workspace)
			defer closeCommandDB(t, conn)
			assertWorkflowRows(t, conn, "workflow_transitions", `SELECT from_state || '->' || to_state FROM workflow_transitions WHERE workflow_id = 'arrow-flow' ORDER BY rowid`, []string{"todo->done", "todo->cancelled", "done->cancelled"})
		})
	}

	for _, appliesTo := range []string{"epic", "task", "subtask", "qa_review", "escalation", "any"} {
		t.Run("applies to "+appliesTo, func(t *testing.T) {
			workspace := t.TempDir()
			runInitCommand(t, workspace, "init")
			states, categories, transitions := "todo,done,cancelled", "unstarted,completed,cancelled", "todo->done"
			if appliesTo == "qa_review" || appliesTo == "escalation" {
				states, categories, transitions = "open,resolved", "unstarted,completed", "open->resolved"
			}
			runFIDCommand(t, workspace, "workflow", "define", "--name", "applies-flow", "--applies-to", appliesTo, "--states", states, "--categories", categories, "--transitions", transitions, "--by", "tl")
			conn := openWorkflowCommandDB(t, workspace)
			defer closeCommandDB(t, conn)
			var value sql.NullString
			if err := conn.QueryRow(`SELECT applies_to FROM workflows WHERE id = 'applies-flow'`).Scan(&value); err != nil {
				t.Fatalf("query applies_to: %v", err)
			}
			if appliesTo == "any" {
				if value.Valid {
					t.Fatalf("applies_to = %q, want NULL", value.String)
				}
			} else if !value.Valid || value.String != appliesTo {
				t.Fatalf("applies_to = %#v, want %q", value, appliesTo)
			}
		})
	}
}

func openWorkflowCommandDB(t *testing.T, workspace string) *sql.DB {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open workflow command database: %v", err)
	}
	return conn
}

func assertWorkflowDefinitionProjection(t *testing.T, conn *sql.DB) {
	t.Helper()

	var (
		name      string
		appliesTo string
		active    bool
	)
	if err := conn.QueryRow(`SELECT name, applies_to, active FROM workflows WHERE id = 'dev-flow-strict'`).Scan(&name, &appliesTo, &active); err != nil {
		t.Fatalf("query workflow: %v", err)
	}
	if name != "dev-flow-strict" || appliesTo != "task" || !active {
		t.Fatalf("workflow row = (%q, %q, %t), want (dev-flow-strict, task, true)", name, appliesTo, active)
	}

	assertWorkflowRows(t, conn, "workflow_states", `SELECT state || ':' || category || ':' || ordinal FROM workflow_states WHERE workflow_id = 'dev-flow-strict' ORDER BY ordinal`, []string{
		"todo:unstarted:1", "in_progress:started:2", "done:completed:3", "cancelled:cancelled:4",
	})
	assertWorkflowRows(t, conn, "workflow_transitions", `SELECT from_state || '->' || to_state FROM workflow_transitions WHERE workflow_id = 'dev-flow-strict' ORDER BY rowid`, []string{
		"todo->in_progress", "in_progress->done", "todo->cancelled", "in_progress->cancelled", "done->cancelled",
	})
	assertWorkflowRows(t, conn, "workflow_transition_requires", `SELECT from_state || '->' || to_state || ':' || rule || ':' || COALESCE(params_json, 'NULL') FROM workflow_transition_requires WHERE workflow_id = 'dev-flow-strict' ORDER BY rowid`, []string{
		`in_progress->done:all_acs_verified:{"roles":["developer"]}`,
		"in_progress->done:has_completion_report:NULL",
	})
}

func assertWorkflowRows(t *testing.T, conn *sql.DB, table, query string, want []string) {
	t.Helper()
	rows, err := conn.Query(query)
	if err != nil {
		t.Fatalf("query %s: %v", table, err)
	}
	defer closeRows(t, rows)

	var got []string
	for rows.Next() {
		var row string
		if err := rows.Scan(&row); err != nil {
			t.Fatalf("scan %s: %v", table, err)
		}
		got = append(got, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate %s: %v", table, err)
	}
	if len(got) != len(want) {
		t.Fatalf("%s rows = %#v, want %#v", table, got, want)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("%s row %d = %q, want %q", table, index, got[index], want[index])
		}
	}
}

func workflowCommandTableSnapshot(t *testing.T, conn *sql.DB) map[string][]string {
	t.Helper()
	snapshot := make(map[string][]string)
	for _, table := range []string{"workflows", "workflow_states", "workflow_transitions", "workflow_transition_requires"} {
		rows, err := conn.Query("SELECT * FROM " + table + " ORDER BY rowid")
		if err != nil {
			t.Fatalf("query %s: %v", table, err)
		}
		columns, err := rows.Columns()
		if err != nil {
			closeRows(t, rows)
			t.Fatalf("columns %s: %v", table, err)
		}
		for rows.Next() {
			values := make([]any, len(columns))
			pointers := make([]any, len(columns))
			for index := range values {
				pointers[index] = &values[index]
			}
			if err := rows.Scan(pointers...); err != nil {
				closeRows(t, rows)
				t.Fatalf("scan %s: %v", table, err)
			}
			snapshot[table] = append(snapshot[table], fmt.Sprint(values...))
		}
		if err := rows.Err(); err != nil {
			closeRows(t, rows)
			t.Fatalf("iterate %s: %v", table, err)
		}
		closeRows(t, rows)
	}
	return snapshot
}
