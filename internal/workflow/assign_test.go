package workflow

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestAssignEmitsScopeAwareReplacementEffects(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	for _, fixture := range []struct {
		name      string
		appliesTo string
	}{
		{name: "second-task", appliesTo: "task"},
		{name: "first-project", appliesTo: "any"},
		{name: "second-project", appliesTo: "any"},
	} {
		if _, err := Define(DefineOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: migrationsDir,
			ActorRole:     "tl",
			Request: DefinitionRequest{
				Name:        fixture.name,
				AppliesTo:   fixture.appliesTo,
				States:      []string{"queued", "done", "cancelled"},
				Categories:  []string{"unstarted", "completed", "cancelled"},
				Transitions: []Transition{{From: "queued", To: "done"}},
			},
		}); err != nil {
			t.Fatalf("Define(%q) error = %v", fixture.name, err)
		}
	}

	for _, fixture := range []struct {
		name    string
		first   AssignmentRequest
		second  AssignmentRequest
		firstID string
		lastID  string
		delete  map[string]any
	}{
		{
			name:    "kind",
			first:   AssignmentRequest{Name: "task-default", Kind: "task"},
			second:  AssignmentRequest{Name: "second-task", Kind: "task"},
			firstID: "task-default",
			lastID:  "second-task",
			delete:  map[string]any{"scope": "kind", "kind": "task"},
		},
		{
			name:    "project",
			first:   AssignmentRequest{Name: "first-project", Project: true},
			second:  AssignmentRequest{Name: "second-project", Project: true},
			firstID: "first-project",
			lastID:  "second-project",
			delete:  map[string]any{"scope": "project"},
		},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			journalPath := workflowJournalPath(t, workspace)
			beforeLines := strings.Split(strings.TrimSpace(string(readWorkflowTestFile(t, journalPath))), "\n")
			if _, err := Assign(AssignOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorRole: "tl", Request: fixture.first}); err != nil {
				t.Fatalf("Assign(first) error = %v", err)
			}
			first := assignmentEventAt(t, journalPath, len(beforeLines))
			if first.EventType != "workflow.assigned" || len(first.Effects) != 1 || first.Effects[0].Op != "insert" {
				t.Fatalf("first effects = %#v, want one insert", first.Effects)
			}
			if got := first.Effects[0].Values["workflow_id"]; got != fixture.firstID {
				t.Fatalf("first workflow_id = %#v, want %q", got, fixture.firstID)
			}
			if _, err := Assign(AssignOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorRole: "tl", Request: fixture.second}); err != nil {
				t.Fatalf("Assign(second) error = %v", err)
			}
			second := assignmentEventAt(t, journalPath, len(beforeLines)+1)
			if len(second.Effects) != 2 || second.Effects[0].Op != "delete" || second.Effects[1].Op != "insert" {
				t.Fatalf("second effects = %#v, want delete then insert", second.Effects)
			}
			if got := second.Effects[0].PK; !mapsEqual(got, fixture.delete) {
				t.Fatalf("delete PK = %#v, want %#v", got, fixture.delete)
			}
			if got := second.Effects[1].Values["workflow_id"]; got != fixture.lastID {
				t.Fatalf("second workflow_id = %#v, want %q", got, fixture.lastID)
			}
			if _, err := journal.RebuildWorkspace(workspace, migrationsDir); err != nil {
				t.Fatalf("RebuildWorkspace() error = %v", err)
			}
			rebuilt, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
			if err != nil {
				t.Fatalf("open rebuilt projection: %v", err)
			}
			defer closeRegistryDB(t, rebuilt)
			var replayed string
			query := `SELECT workflow_id FROM workflow_assignments WHERE scope = ?`
			args := []any{"project"}
			if !fixture.second.Project {
				query += ` AND kind = ?`
				args = []any{"kind", fixture.second.Kind}
			}
			if err := rebuilt.QueryRow(query, args...).Scan(&replayed); err != nil {
				t.Fatalf("query replayed assignment: %v", err)
			}
			if replayed != fixture.lastID {
				t.Fatalf("replayed assignment workflow_id = %q, want %q", replayed, fixture.lastID)
			}
		})
	}
}

func assignmentEventAt(t *testing.T, journalPath string, index int) journal.Event {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(string(readWorkflowTestFile(t, journalPath))), "\n")
	var event journal.Event
	if err := json.Unmarshal([]byte(lines[index]), &event); err != nil {
		t.Fatalf("decode event %d: %v", index, err)
	}
	return event
}

func mapsEqual(left, right map[string]any) bool {
	if len(left) != len(right) {
		return false
	}
	for key, want := range right {
		if left[key] != want {
			return false
		}
	}
	return true
}
