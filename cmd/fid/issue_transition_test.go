package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestIssueTransitionChangesStateAndPrintsConfirmation(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI transition", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	_, transitionStderr := runFIDCommand(t, workspace, "issue", "transition", issueID, "--to", "in_progress", "--by", "developer")
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": todo → in_progress") {
		t.Fatalf("stderr = %q, want transition confirmation", transitionStderr)
	}

	stdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

	var got struct {
		Status         string `json:"status"`
		StatusCategory string `json:"status_category"`
	}
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Status != "in_progress" {
		t.Fatalf("status = %q, want in_progress", got.Status)
	}
	if got.StatusCategory != "started" {
		t.Fatalf("status_category = %q, want started", got.StatusCategory)
	}
}

func TestIssueTransitionRejectsInvalidEdgeAndLeavesStateUnchanged(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI invalid transition", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	beforeStdout, beforeStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if beforeStderr != "" {
		t.Fatalf("issue get before stderr = %q, want empty", beforeStderr)
	}

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "transition", issueID, "--to", "done", "--by", "tl")
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid transition rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"error: transition issue",
		`"todo"`,
		`"done"`,
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %q, want substring %q", stderr, want)
		}
	}

	afterStdout, afterStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if afterStderr != "" {
		t.Fatalf("issue get after stderr = %q, want empty", afterStderr)
	}
	if afterStdout != beforeStdout {
		t.Fatalf("issue JSON after rejected transition = %q, want %q", afterStdout, beforeStdout)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore {
		t.Fatalf("journal lines after rejected transition = %d, want %d", linesAfter, linesBefore)
	}
}

func TestIssueTransitionRendersRuleFailures(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI rule failure", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	appendCommandIssueEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "acceptance_criteria.seeded",
		Actor: journal.Actor{
			Role: "qa_lead",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []journal.Effect{
			{
				Op:    "insert",
				Table: "acceptance_criteria",
				PK:    map[string]any{"id": "AC-cli-001"},
				Values: map[string]any{
					"issue_id": issueID,
					"ordinal":  1,
					"text":     "CLI AC",
				},
			},
		},
	})

	runFIDCommand(t, workspace, "issue", "transition", issueID, "--to", "in_progress", "--by", "developer")

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "transition", issueID, "--to", "dev_complete", "--by", "developer")
	if err == nil {
		t.Fatal("Execute() error = nil, want rule rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"all_acs_verified",
		"has_completion_report",
		"AC-cli-001",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %q, want substring %q", stderr, want)
		}
	}
}

func TestIssueTransitionRequiresTargetStateAndKnownByRole(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI role validation", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	t.Run("missing to", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "transition", issueID, "--by", "developer")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --to rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: transition issue: target state is required") {
			t.Fatalf("stderr = %q, want missing --to error", stderr)
		}
	})

	t.Run("missing by", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "transition", issueID, "--to", "in_progress")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --by rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: --by is required") {
			t.Fatalf("stderr = %q, want missing --by error", stderr)
		}
	})

	t.Run("unknown by", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "transition", issueID, "--to", "in_progress", "--by", "intern")
		if err == nil {
			t.Fatal("Execute() error = nil, want unknown --by rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
			t.Fatalf("stderr = %q, want unknown role error", stderr)
		}
	})
}

func TestIssueTransitionUsesStateSpecificPermissions(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI state permissions", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	_, startStderr := runFIDCommand(t, workspace, "issue", "transition", issueID, "--to", "in_progress", "--by", "developer")
	if !strings.Contains(startStderr, "transitioned "+issueID+": todo → in_progress") {
		t.Fatalf("in_progress stderr = %q, want start confirmation", startStderr)
	}

	_, completeStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete", issueID,
		"--by", "developer",
		"--notes-md", "project-docs/docs/phases/phase-10/tasks/task-3-done.md",
		"--tokens", "1200",
	)
	if !strings.Contains(completeStderr, "completed "+issueID+" (attempt 1)") {
		t.Fatalf("complete stderr = %q, want completion confirmation", completeStderr)
	}

	journalPath := mustJournalPath(t, workspace)
	linesBeforeDenied := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "transition", issueID, "--to", "dev_complete", "--by", "qa_reviewer")
	if err == nil {
		t.Fatal("qa_reviewer dev_complete error = nil, want permission denial")
	}
	if stdout != "" {
		t.Fatalf("qa_reviewer stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		`error: role "qa_reviewer" not permitted for action "issue.transition.dev_complete"`,
		"check `fid policy show` for permitted role/action combinations",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("qa_reviewer stderr = %q, want substring %q", stderr, want)
		}
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforeDenied {
		t.Fatalf("journal lines after qa_reviewer denial = %d, want %d", got, linesBeforeDenied)
	}

	allowedStdout, allowedStderr := runFIDCommand(t, workspace, "issue", "transition", issueID, "--to", "dev_complete", "--by", "developer")
	if allowedStdout != "" {
		t.Fatalf("developer stdout = %q, want empty", allowedStdout)
	}
	if !strings.Contains(allowedStderr, "transitioned "+issueID+": in_progress → dev_complete") {
		t.Fatalf("developer stderr = %q, want dev_complete confirmation", allowedStderr)
	}
}
