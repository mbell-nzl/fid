//go:build integration

package integration

import (
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV02aWorkflowAuthoringExitScenarioAndRebuildParity(t *testing.T) {
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	closureRunFID(t, binary, workspace, "init", "--prefix", "FID")
	journalPath := closureJournalPath(t, workspace)

	before := closureReadFile(t, journalPath)
	_, stderr, err := closureRunFIDExpectError(t, binary, workspace,
		"workflow", "define", "--name", "rejected-rule", "--applies-to", "task",
		"--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done",
		"--require", "todo->done:concerns_addressed", "--by", "tl",
	)
	if err == nil || !strings.Contains(stderr, "unknown flag: --require") {
		t.Fatalf("removed requirement rejection = (%v, %q), want unregistered --require", err, stderr)
	}
	assertClosureJournalBytesEqual(t, journalPath, before, "unregistered workflow requirement")

	closureRunFID(t, binary, workspace,
		"workflow", "define", "--name", "task-strict", "--applies-to", "task",
		"--states", "todo,in_progress,done,cancelled",
		"--categories", "unstarted,started,completed,cancelled",
		"--transitions", "todo->in_progress,in_progress->done", "--by", "tl",
	)
	closureRunFID(t, binary, workspace, "workflow", "assign", "--name", "task-strict", "--kind", "task", "--by", "tl")
	controlID := closureCreateIssue(t, binary, workspace, "task", "Assigned strict task", nil, nil, nil)
	assertIssueWorkflow(t, binary, workspace, controlID, "task-strict", "todo")
	assertWholeProjectionRebuildParity(t, binary, workspace, journalPath)
}

func TestMaintenanceCheckReportsIncompatibleWorkflowAssignment(t *testing.T) {
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	closureRunFID(t, binary, workspace, "init", "--prefix", "FID")

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO workflow_assignments (scope, kind, workflow_id) VALUES ('kind', 'epic', 'task-default')`); err != nil {
		_ = conn.Close()
		t.Fatalf("seed incompatible assignment: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close projection: %v", err)
	}

	stdout, stderr, err := closureRunFIDExpectError(t, binary, workspace, "maintenance", "check", "--warn-as-error", "--format", "json")
	if err == nil || stderr != "" || !strings.Contains(stdout, `"check":"workflow_assignment"`) || !strings.Contains(stdout, `"entity_id":"kind:epic"`) {
		t.Fatalf("maintenance check = (%v, %q, %q), want non-zero workflow_assignment finding", err, stdout, stderr)
	}
}

func assertIssueWorkflow(t *testing.T, binary, workspace, issueID, workflowID, status string) {
	t.Helper()
	stdout, _ := closureRunFID(t, binary, workspace, "--json", "issue", "get", issueID)
	for _, value := range []string{`"workflow_id":"` + workflowID + `"`, `"status":"` + status + `"`} {
		if !strings.Contains(stdout, value) {
			t.Fatalf("issue get %s = %q, want %s", issueID, stdout, value)
		}
	}
}

func assertWholeProjectionRebuildParity(t *testing.T, binary, workspace, journalPath string) {
	t.Helper()
	beforeTables := closureProjectedTableSnapshots(t, binary, workspace)
	beforeSchema := closureSchemaSnapshot(t, binary, workspace)
	beforeJournal := closureJournalFileBytes(t, workspace)
	assertClosureForeignKeysClean(t, binary, workspace)
	closureRunFID(t, binary, workspace, "rebuild", "--by", "tl")
	assertClosureJournalBytesEqual(t, journalPath, beforeJournal[filepath.Base(journalPath)], "rebuild")
	assertClosureJournalFilesEqual(t, closureJournalFileBytes(t, workspace), beforeJournal, "rebuild")
	assertClosureProjectedTablesEqual(t, closureProjectedTableSnapshots(t, binary, workspace), beforeTables)
	if got := closureSchemaSnapshot(t, binary, workspace); got != beforeSchema {
		t.Fatal("sqlite_schema after rebuild differs from before rebuild")
	}
	assertClosureForeignKeysClean(t, binary, workspace)
}
