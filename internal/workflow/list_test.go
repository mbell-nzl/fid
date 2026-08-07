package workflow

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestListFiltersCountsAssignmentsAndOrdering(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	seedListedWorkflows(t, conn)

	items, err := List(ListOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Kind: "qa_review", Active: true})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if got, want := workflowListIDs(items), []string{"any-active", "qa-active"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("filtered IDs = %v, want %v", got, want)
	}
	if items[0].AppliesTo != nil || items[0].States != 2 || items[0].Transitions != 1 {
		t.Fatalf("any-active = %#v, want null applicability and projected counts", items[0])
	}
	if got, want := items[0].AssignedTo, []Assignment{{Scope: "kind", Kind: "subtask"}, {Scope: "kind", Kind: "task"}, {Scope: "project"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("any-active assignments = %#v, want %#v", got, want)
	}
	if items[1].AssignedTo == nil || len(items[1].AssignedTo) != 0 {
		t.Fatalf("qa-active assignments = %#v, want empty array", items[1].AssignedTo)
	}

	unfiltered, err := List(ListOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Kind: "qa_review"})
	if err != nil {
		t.Fatalf("List(unfiltered) error = %v", err)
	}
	if got, want := workflowListIDs(unfiltered), []string{"any-active", "qa-active", "qa-inactive"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unfiltered IDs = %v, want %v", got, want)
	}
}

func TestListReturnsEmptySliceForNoMatch(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	items, err := List(ListOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Kind: "missing"})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if items == nil || len(items) != 0 {
		t.Fatalf("items = %#v, want empty slice", items)
	}
}

func seedListedWorkflows(t *testing.T, conn *sql.DB) {
	t.Helper()
	for _, workflow := range []struct {
		id, appliesTo string
		active        bool
	}{
		{"any-active", "", true},
		{"qa-active", "qa_review", true},
		{"qa-inactive", "qa_review", false},
		{"other-active", "epic", true},
	} {
		var appliesTo any
		if workflow.appliesTo != "" {
			appliesTo = workflow.appliesTo
		}
		if _, err := conn.Exec(`INSERT INTO workflows (id, name, applies_to, active) VALUES (?, ?, ?, ?)`, workflow.id, workflow.id, appliesTo, workflow.active); err != nil {
			t.Fatalf("insert workflow %q: %v", workflow.id, err)
		}
	}
	for _, state := range []string{"open", "closed"} {
		if _, err := conn.Exec(`INSERT INTO workflow_states (workflow_id, state, category, ordinal) VALUES ('any-active', ?, 'unstarted', ?)`, state, len(state)); err != nil {
			t.Fatalf("insert state %q: %v", state, err)
		}
	}
	if _, err := conn.Exec(`INSERT INTO workflow_transitions (workflow_id, from_state, to_state) VALUES ('any-active', 'open', 'closed')`); err != nil {
		t.Fatalf("insert transition: %v", err)
	}
	for _, assignment := range []struct{ scope, kind string }{{"kind", "task"}, {"kind", "subtask"}, {"project", ""}} {
		var kind any
		if assignment.kind != "" {
			kind = assignment.kind
		}
		if _, err := conn.Exec(`INSERT INTO workflow_assignments (scope, kind, workflow_id) VALUES (?, ?, 'any-active')`, assignment.scope, kind); err != nil {
			t.Fatalf("insert %s assignment: %v", assignment.scope, err)
		}
	}
}

func workflowListIDs(items []ListItem) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}
