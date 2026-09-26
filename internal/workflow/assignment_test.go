package workflow

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestValidateAssignmentRequestRejectsAndAcceptsScopeMatrix(t *testing.T) {
	t.Parallel()
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)

	tests := []struct {
		name    string
		request AssignmentRequest
		wantErr string
	}{
		{
			name:    "unknown workflow",
			request: AssignmentRequest{Name: "missing-flow", Kind: "task"},
			wantErr: `"missing-flow"`,
		},
		{
			name:    "kind incompatible with applies to",
			request: AssignmentRequest{Name: "epic-default", Kind: "task"},
			wantErr: `"epic"`,
		},
		{
			name:    "project requires any workflow",
			request: AssignmentRequest{Name: "task-default", Project: true},
			wantErr: `"task"`,
		},
		{
			name:    "both scopes",
			request: AssignmentRequest{Name: "task-default", Kind: "task", Project: true},
			wantErr: "both kind and project",
		},
		{
			name:    "neither scope",
			request: AssignmentRequest{Name: "task-default"},
			wantErr: "kind or project",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAssignmentRequest(conn, tt.request)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("ValidateAssignmentRequest() error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}

	for _, request := range []AssignmentRequest{
		{Name: "task-default", Kind: "task"},
		{Name: "any-flow", Project: true},
	} {
		if request.Name == "any-flow" {
			if _, err := conn.Exec(`INSERT INTO workflows (id, name, applies_to, active) VALUES ('any-flow', 'any-flow', NULL, 1)`); err != nil {
				t.Fatalf("insert any workflow: %v", err)
			}
		}
		if err := ValidateAssignmentRequest(conn, request); err != nil {
			t.Fatalf("ValidateAssignmentRequest(%#v) error = %v", request, err)
		}
	}
}

func TestValidateAssignmentRequestDoesNotWrite(t *testing.T) {
	t.Parallel()
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	before := workflowTableSnapshot(t, conn)
	journalPath := filepath.Join(t.TempDir(), "events.jsonl")
	if err := os.WriteFile(journalPath, []byte("journal before\n"), 0o600); err != nil {
		t.Fatalf("write journal fixture: %v", err)
	}
	beforeJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal before: %v", err)
	}
	if _, err := conn.Exec(`PRAGMA query_only = ON`); err != nil {
		t.Fatalf("set query_only: %v", err)
	}
	if err := ValidateAssignmentRequest(conn, AssignmentRequest{Name: "task-default", Kind: "task"}); err != nil {
		t.Fatalf("ValidateAssignmentRequest() error = %v", err)
	}
	if after := workflowTableSnapshot(t, conn); !reflect.DeepEqual(after, before) {
		t.Fatalf("workflow tables after assignment validation = %#v, want %#v", after, before)
	}
	afterJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal after: %v", err)
	}
	if !reflect.DeepEqual(afterJournal, beforeJournal) {
		t.Fatalf("journal bytes after assignment validation = %q, want %q", afterJournal, beforeJournal)
	}
}
