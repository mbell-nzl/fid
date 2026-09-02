package issues

import (
	"path/filepath"
	"testing"
)

func TestGetListAndReadyUseProjectedWorkflowCategories(t *testing.T) {
	workspace := newIssueWorkspace(t)
	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Projected category consumer",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	if _, err := conn.Exec(`
		UPDATE workflow_states
		SET category = ?
		WHERE workflow_id = ? AND state = ?
	`, "started", "task-default", "todo"); err != nil {
		closeIssueTestDB(t, conn)
		t.Fatalf("update projected workflow category: %v", err)
	}
	closeIssueTestDB(t, conn)

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.StatusCategory != "started" {
		t.Fatalf("Get().StatusCategory = %q, want projected category started", got.StatusCategory)
	}

	listed, err := List(ListOptions{
		WorkspaceRoot:  workspace,
		MigrationsDir:  testIssueMigrationsDir(t),
		StatusCategory: "started",
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(listed) != 1 || listed[0].ID != created.ID || listed[0].StatusCategory != "started" {
		t.Fatalf("List() = %#v, want projected-category issue %q", listed, created.ID)
	}

	ready, err := Ready(ReadyOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("Ready() error = %v", err)
	}
	if ready.Task != nil {
		t.Fatalf("Ready().Task = %#v, want no task after projected category change", ready.Task)
	}
}
