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
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeRegistryDB(t, conn)
	for _, id := range []string{"second-task", "first-project", "second-project"} {
		if _, err := conn.Exec(`INSERT INTO workflows (id, name, applies_to, active) VALUES (?, ?, NULL, 1)`, id, id); err != nil {
			t.Fatalf("insert workflow %q: %v", id, err)
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
