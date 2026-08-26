package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestWorkflowAssignJSONPrecedenceAndScopedReplacement(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	defineWorkflowForTest(t, workspace, "project-flow", "any")
	defineWorkflowForTest(t, workspace, "task-flow", "task")
	if result := executeCommandForTest(workspace, "--json", "workflow", "assign", "--name", "project-flow", "--project-wide", "--by", "tl"); result.err != nil {
		t.Fatalf("project assignment = %#v", result)
	}
	result := executeCommandForTest(workspace, "--json", "workflow", "assign", "--name", "task-flow", "--kind", "task", "--by", "tl")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("kind assignment = %#v", result)
	}
	var envelope struct {
		Result struct {
			Assignment struct {
				ScopeType  string `json:"scope_type"`
				ScopeValue string `json:"scope_value"`
				WorkflowID string `json:"workflow_id"`
			} `json:"assignment"`
		} `json:"result"`
		Event *struct {
			ID         string `json:"id"`
			OccurredAt string `json:"occurred_at"`
		} `json:"event"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil || envelope.Result.Assignment.ScopeType != "kind" || envelope.Result.Assignment.ScopeValue != "task" || envelope.Result.Assignment.WorkflowID != "task-flow" || envelope.Event == nil || envelope.Event.ID == "" || envelope.Event.OccurredAt == "" {
		t.Fatalf("assignment response = %+v, err = %v", envelope, err)
	}
	tail := mustReadLastJournalEvent(t, mustJournalPath(t, workspace))
	if tail.EventID != envelope.Event.ID || tail.OccurredAt.UTC().Format(time.RFC3339Nano) != envelope.Event.OccurredAt {
		t.Fatalf("assignment event = %+v, journal tail = %#v; want matching identity", envelope.Event, tail)
	}
	conn := openWorkflowCommandDB(t, workspace)
	defer closeCommandDB(t, conn)
	resolved, err := workflow.ResolveIssueWorkflow(conn, "task")
	if err != nil || resolved.ID != "task-flow" {
		t.Fatalf("ResolveIssueWorkflow = %+v, %v; want kind assignment precedence", resolved, err)
	}
	if result := executeCommandForTest(workspace, "--json", "workflow", "assign", "--name", "project-flow", "--project-wide", "--by", "tl"); result.err != nil {
		t.Fatalf("replacement assignment = %#v", result)
	}
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM workflow_assignments WHERE scope = 'project'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("project assignment count = %d, err = %v; want one", count, err)
	}
}

func TestWorkflowAssignRequiresExactlyOneScopeAndMapsUnknownWorkflow(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	journal := mustJournalPath(t, workspace)
	before := mustReadCommandFile(t, journal)
	for _, args := range [][]string{{"--json", "workflow", "assign", "--name", "missing", "--by", "tl"}, {"--json", "workflow", "assign", "--name", "missing", "--kind", "task", "--project-wide", "--by", "tl"}, {"--json", "workflow", "assign", "--name", "missing", "--kind", "qa_review", "--by", "tl"}} {
		result := executeCommandForTest(workspace, args...)
		if result.err == nil || result.stdout != "" {
			t.Fatalf("scope rejection = %#v", result)
		}
		assertWorkflowInvalidEnvelope(t, result.stderr, "workflow.assign")
		if after := mustReadCommandFile(t, journal); string(after) != string(before) {
			t.Fatal("rejected assignment appended an event")
		}
	}
	result := executeCommandForTest(workspace, "--json", "workflow", "assign", "--name", "missing", "--kind", "task", "--by", "tl")
	var failure struct {
		Error struct {
			Code    string `json:"code"`
			Subject struct {
				ID string `json:"id"`
			} `json:"subject"`
		} `json:"error"`
	}
	if result.err == nil || json.Unmarshal([]byte(result.stderr), &failure) != nil || failure.Error.Code != "workflow-not-found" || failure.Error.Subject.ID != "missing" {
		t.Fatalf("unknown workflow = %#v, failure = %+v", result, failure)
	}
	if after := mustReadCommandFile(t, journal); string(after) != string(before) {
		t.Fatal("unknown workflow assignment appended an event")
	}
}

func TestWorkflowAssignCombinesRootProjectSelectorAndProjectWideScope(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	defineWorkflowForTest(t, workspace, "root-project-flow", "any")

	assign := newWorkflowAssignCommand()
	if assign.Flags().Lookup("project") != nil {
		t.Fatal("workflow assign registers a local --project flag, want only root --project <path>")
	}
	result := executeCommandForTest(t.TempDir(), "--project", workspace, "--json", "workflow", "assign", "--name", "root-project-flow", "--project-wide", "--by", "tl")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("root project plus project-wide assignment = %#v", result)
	}
	var envelope struct {
		Result struct {
			Assignment struct {
				ScopeType  string `json:"scope_type"`
				ScopeValue any    `json:"scope_value"`
				WorkflowID string `json:"workflow_id"`
			} `json:"assignment"`
		} `json:"result"`
		Event any `json:"event"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil || envelope.Result.Assignment.ScopeType != "project" || envelope.Result.Assignment.ScopeValue != nil || envelope.Result.Assignment.WorkflowID != "root-project-flow" || envelope.Event == nil {
		t.Fatalf("root project assignment envelope = %+v, err = %v", envelope, err)
	}
}

func defineWorkflowForTest(t *testing.T, workspace, name, appliesTo string) {
	t.Helper()
	result := executeCommandForTest(workspace, "--json", "workflow", "define", "--name", name, "--applies-to", appliesTo, "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--by", "tl")
	if result.err != nil {
		t.Fatalf("define %s = %#v", name, result)
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
