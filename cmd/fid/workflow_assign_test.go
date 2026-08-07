package main

import (
	"database/sql"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestWorkflowAssignProjectsOneKindAssignment(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	runFIDCommand(t, workspace,
		"workflow", "define",
		"--name", "assigned-task-flow",
		"--applies-to", "task",
		"--states", "queued,done,cancelled",
		"--categories", "unstarted,completed,cancelled",
		"--transitions", "queued->done",
		"--by", "tl",
	)

	journalPath := mustJournalPath(t, workspace)
	beforeLines := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	stdout, stderr := runFIDCommand(t, workspace, "workflow", "assign", "--name", "assigned-task-flow", "--kind", "task", "--by", "tl")
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if want := "assigned workflow assigned-task-flow to kind task\n"; stderr != want {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != beforeLines+1 {
		t.Fatalf("journal lines = %d, want %d", got, beforeLines+1)
	}
	event := mustReadLastJournalEvent(t, journalPath)
	if event.EventType != "workflow.assigned" || event.EventSchema != journal.CurrentEventSchema || len(event.Effects) != 1 {
		t.Fatalf("event = %#v, want one current workflow.assigned event", event)
	}

	conn := openWorkflowCommandDB(t, workspace)
	defer closeCommandDB(t, conn)
	var (
		kind       string
		workflowID string
	)
	if err := conn.QueryRow(`SELECT kind, workflow_id FROM workflow_assignments WHERE scope = 'kind'`).Scan(&kind, &workflowID); err != nil {
		t.Fatalf("query assignment: %v", err)
	}
	if kind != "task" || workflowID != "assigned-task-flow" {
		t.Fatalf("assignment = (%q, %q), want (task, assigned-task-flow)", kind, workflowID)
	}
	var activityCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&activityCount); err != nil {
		t.Fatalf("count activity_log: %v", err)
	}
	if activityCount != 0 {
		t.Fatalf("activity_log rows = %d, want 0", activityCount)
	}
}

func TestWorkflowAssignRejectsWithoutMutation(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	journalPath := mustJournalPath(t, workspace)

	seen := make(map[string]struct{})
	for _, fixture := range []struct {
		args []string
		want string
	}{
		{[]string{"workflow", "assign", "--name", "missing", "--kind", "task", "--by", "tl"}, "error: workflow \"missing\" does not exist\n"},
		{[]string{"workflow", "assign", "--name", "task-default", "--project", "--by", "tl"}, "error: workflow \"task-default\" applies_to \"task\" cannot be assigned to project\n"},
		{[]string{"workflow", "assign", "--name", "task-default", "--kind", "task", "--project", "--by", "tl"}, "error: workflow assignment specifies both kind and project scopes\n"},
		{[]string{"workflow", "assign", "--name", "task-default", "--by", "tl"}, "error: workflow assignment requires a kind or project scope\n"},
		{[]string{"workflow", "assign", "--name", "task-default", "--kind", "epic", "--by", "tl"}, "error: workflow \"task-default\" applies_to \"task\" is incompatible with kind \"epic\"\n"},
		{[]string{"workflow", "assign", "--name", "task-default", "--kind", "task", "--by", "developer"}, "error: role \"developer\" not permitted for action \"workflow.assign\"\n  check `fid policy show` for permitted role/action combinations\n"},
	} {
		args := fixture.args
		before := mustReadCommandFile(t, journalPath)
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, args...)
		if err == nil || stdout != "" || stderr != fixture.want {
			t.Fatalf("assign %v = (%q, %q, %v), want rejection", args, stdout, stderr, err)
		}
		if _, duplicate := seen[stderr]; duplicate {
			t.Fatalf("assign %v stderr = %q, want a value-naming error distinct from every prior rejection", args, stderr)
		}
		seen[stderr] = struct{}{}
		if after := mustReadCommandFile(t, journalPath); string(after) != string(before) {
			t.Fatalf("assign %v changed journal", args)
		}
		conn := openWorkflowCommandDB(t, workspace)
		var count int
		if err := conn.QueryRow(`SELECT COUNT(*) FROM workflow_assignments`).Scan(&count); err != nil {
			closeCommandDB(t, conn)
			t.Fatalf("count assignments: %v", err)
		}
		closeCommandDB(t, conn)
		if count != 0 {
			t.Fatalf("assign %v left %d assignment rows", args, count)
		}
	}
	if len(seen) != 6 {
		t.Fatalf("distinct rejection messages = %d, want 6", len(seen))
	}
}

var _ sql.NullString
