//go:build integration

package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV02aWorkflowAuthoringExitScenarioAndRebuildParity(t *testing.T) {
	binary := buildClosureFIDBinary(t)
	workspace := t.TempDir()
	closureRunFID(t, binary, workspace, "init")
	journalPath := closureJournalPath(t, workspace)

	assertReservedRulesAreNotSeeded(t, binary, workspace)
	assertReservedWorkflowRulesRejected(t, binary, workspace, journalPath)
	controlID := closureCreateIssue(t, binary, workspace, "task", "Unassigned control", nil, nil, nil)

	closureRunFID(t, binary, workspace,
		"workflow", "define", "--name", "task-strict", "--applies-to", "task",
		"--states", "todo,in_progress,dev_complete,done,cancelled",
		"--categories", "unstarted,started,started,completed,cancelled",
		"--transitions", "todo->in_progress,in_progress->dev_complete,dev_complete->done",
		"--require", "in_progress->dev_complete:concerns_addressed", "--by", "tl",
	)
	for _, formatName := range []string{"json", "toon", "markdown", "yaml"} {
		stdout, _ := closureRunFID(t, binary, workspace, "workflow", "list", "--format", formatName)
		if !strings.Contains(stdout, "task-strict") {
			t.Fatalf("workflow list %s = %q, want task-strict", formatName, stdout)
		}
	}
	closureRunFID(t, binary, workspace, "workflow", "assign", "--name", "task-strict", "--kind", "task", "--by", "tl")
	strictID := closureCreateIssue(t, binary, workspace, "task", "Assigned strict task", nil, nil, nil)
	assertIssueWorkflow(t, binary, workspace, controlID, "task-default", "todo")
	assertIssueWorkflow(t, binary, workspace, strictID, "task-strict", "todo")

	closureRunFID(t, binary, workspace, "issue", "claim", strictID, "--by", "developer")
	closureRunFID(t, binary, workspace, "issue", "transition", strictID, "--to", "in_progress", "--by", "developer")
	completionID := closureCompleteIssue(t, binary, workspace, strictID)
	_, concernStderr := closureRunFID(t, binary, workspace, "concern", "raise", "--on", strictID, "--by", "developer", "--in-completion", completionID, "--summary", "Must resolve before development completes")
	concernID := closureIDAfterWord(t, concernStderr, "concern")
	beforeBlockedTransition := closureReadFile(t, journalPath)
	_, blockedStderr, blockedErr := closureRunFIDExpectError(t, binary, workspace, "issue", "transition", strictID, "--to", "dev_complete", "--by", "developer")
	if blockedErr == nil || !strings.Contains(blockedStderr, "concerns_addressed") || !strings.Contains(blockedStderr, concernID) {
		t.Fatalf("blocked transition = (%v, %q), want BlockedTransitionError naming %s and concerns_addressed", blockedErr, blockedStderr, concernID)
	}
	assertClosureJournalBytesEqual(t, journalPath, beforeBlockedTransition, "blocked concerns-addressed transition")
	addressingCompletionID := closureCompleteIssue(t, binary, workspace, strictID)
	closureRunFID(t, binary, workspace, "concern", "address", concernID, "--in-completion", addressingCompletionID, "--by", "developer")
	closureRunFID(t, binary, workspace, "issue", "transition", strictID, "--to", "dev_complete", "--by", "developer")
	assertIssueWorkflow(t, binary, workspace, strictID, "task-strict", "dev_complete")
	assertIssueWorkflow(t, binary, workspace, controlID, "task-default", "todo")

	assertCustomStateTransitionBoundary(t, binary, workspace, journalPath)
	assertWholeProjectionRebuildParity(t, binary, workspace, journalPath)
}

func assertReservedRulesAreNotSeeded(t *testing.T, binary, workspace string) {
	t.Helper()
	stdout, _ := closureRunFID(t, binary, workspace, "sql", `
		SELECT workflow_id, from_state, to_state, rule
		FROM workflow_transition_requires
		WHERE rule IN ('required_gates_current', 'has_completion')
		ORDER BY workflow_id, from_state, to_state, rule
	`, "--format", "json")
	if strings.TrimSpace(stdout) != "[]" {
		t.Fatalf("seeded reserved rules = %q, want []", stdout)
	}
}

func TestMaintenanceCheckReportsIncompatibleWorkflowAssignment(t *testing.T) {
	binary := buildClosureFIDBinary(t)
	workspace := t.TempDir()
	closureRunFID(t, binary, workspace, "init")

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	if _, err := conn.Exec(`
		INSERT INTO workflow_assignments (scope, kind, workflow_id)
		VALUES ('kind', 'epic', 'task-default')
	`); err != nil {
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

func assertReservedWorkflowRulesRejected(t *testing.T, binary, workspace, journalPath string) {
	t.Helper()
	for _, fixture := range []struct {
		rule    string
		version string
	}{
		{rule: "required_gates_current", version: "v0.2b"},
		{rule: "has_completion", version: "v0.2c"},
	} {
		before := closureReadFile(t, journalPath)
		_, stderr, err := closureRunFIDExpectError(t, binary, workspace,
			"workflow", "define", "--name", "reserved-"+strings.ReplaceAll(fixture.rule, "_", "-"), "--applies-to", "task",
			"--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done",
			"--require", "todo->done:"+fixture.rule, "--by", "tl",
		)
		if err == nil || !strings.Contains(stderr, fixture.rule) || !strings.Contains(stderr, fixture.version) {
			t.Fatalf("reserved %s rejection = (%v, %q), want version %s", fixture.rule, err, stderr, fixture.version)
		}
		assertClosureJournalBytesEqual(t, journalPath, before, "reserved workflow rule "+fixture.rule)
	}
}

func assertCustomStateTransitionBoundary(t *testing.T, binary, workspace, journalPath string) {
	t.Helper()
	closureRunFID(t, binary, workspace,
		"workflow", "define", "--name", "custom-state-flow", "--applies-to", "task",
		"--states", "todo,backlog,done,cancelled", "--categories", "unstarted,started,completed,cancelled",
		"--transitions", "todo->backlog,backlog->done", "--by", "tl",
	)
	closureRunFID(t, binary, workspace, "workflow", "assign", "--name", "custom-state-flow", "--kind", "task", "--by", "tl")
	customStateID := closureCreateIssue(t, binary, workspace, "task", "Custom state boundary", nil, nil, nil)
	before := closureReadFile(t, journalPath)
	_, stderr, err := closureRunFIDExpectError(t, binary, workspace, "issue", "transition", customStateID, "--to", "backlog", "--by", "developer")
	if err == nil || !strings.Contains(stderr, `authorize action "issue.transition.backlog": no permissions matrix entry`) {
		t.Fatalf("custom state transition = (%v, %q), want documented no-matrix-entry error", err, stderr)
	}
	assertClosureJournalBytesEqual(t, journalPath, before, "custom state transition boundary")

	closureRunFID(t, binary, workspace, "workflow", "assign", "--name", "task-strict", "--kind", "task", "--by", "tl")
	compatibleID := closureCreateIssue(t, binary, workspace, "task", "v0.1 state control", nil, nil, nil)
	closureRunFID(t, binary, workspace, "issue", "claim", compatibleID, "--by", "developer")
	closureRunFID(t, binary, workspace, "issue", "transition", compatibleID, "--to", "in_progress", "--by", "developer")
	assertIssueWorkflow(t, binary, workspace, compatibleID, "task-strict", "in_progress")
}

func assertWholeProjectionRebuildParity(t *testing.T, binary, workspace, journalPath string) {
	t.Helper()
	beforeTables := closureProjectedTableSnapshots(t, binary, workspace)
	beforeSchema := closureSchemaSnapshot(t, binary, workspace)
	beforeJournal := closureJournalFileBytes(t, workspace)
	assertClosureForeignKeysClean(t, binary, workspace)

	closureRunFID(t, binary, workspace, "rebuild", "--by", "tl")
	assertClosureJournalBytesEqual(t, journalPath, beforeJournal[filepath.Base(journalPath)], "first rebuild")
	assertClosureJournalFilesEqual(t, closureJournalFileBytes(t, workspace), beforeJournal, "first rebuild")
	assertClosureProjectedTablesEqual(t, closureProjectedTableSnapshots(t, binary, workspace), beforeTables)
	if got := closureSchemaSnapshot(t, binary, workspace); got != beforeSchema {
		t.Fatal("sqlite_schema after first rebuild differs from before rebuild")
	}
	assertClosureForeignKeysClean(t, binary, workspace)

	closureRunFID(t, binary, workspace, "rebuild", "--by", "tl")
	assertClosureJournalFilesEqual(t, closureJournalFileBytes(t, workspace), beforeJournal, "second rebuild")
	assertClosureProjectedTablesEqual(t, closureProjectedTableSnapshots(t, binary, workspace), beforeTables)
	if got := closureSchemaSnapshot(t, binary, workspace); got != beforeSchema {
		t.Fatal("sqlite_schema after second rebuild differs from before rebuild")
	}
	assertClosureForeignKeysClean(t, binary, workspace)
}

func assertIssueWorkflow(t *testing.T, binary, workspace, issueID, workflowID, status string) {
	t.Helper()
	stdout, _ := closureRunFID(t, binary, workspace, "issue", "get", issueID, "--format", "json")
	for _, value := range []string{`"workflow_id":"` + workflowID + `"`, `"status":"` + status + `"`} {
		if !bytes.Contains([]byte(stdout), []byte(value)) {
			t.Fatalf("issue get %s = %q, want %s", issueID, stdout, value)
		}
	}
}
