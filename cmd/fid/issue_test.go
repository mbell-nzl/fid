package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestIssueCreateThenGetReturnsJSONByDefaultAndWithFormatFlag(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Implement tracer bullet", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	defaultStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID)
	jsonStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")

	if defaultStdout != jsonStdout {
		t.Fatalf("default stdout = %q, want same as --format json %q", defaultStdout, jsonStdout)
	}

	var got struct {
		ID                 string   `json:"id"`
		Kind               string   `json:"kind"`
		Title              string   `json:"title"`
		Status             string   `json:"status"`
		StatusCategory     string   `json:"status_category"`
		WorkflowID         string   `json:"workflow_id"`
		OwnerRole          string   `json:"owner_role"`
		Labels             []string `json:"labels"`
		InInbox            bool     `json:"in_inbox"`
		ConcernsRaised     bool     `json:"concerns_raised"`
		OpenConcernsCount  int      `json:"open_concerns_count"`
		AcceptanceCriteria []any    `json:"acceptance_criteria"`
	}
	if err := json.Unmarshal([]byte(defaultStdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.ID != issueID {
		t.Fatalf("id = %q, want %q", got.ID, issueID)
	}
	if got.Kind != "task" {
		t.Fatalf("kind = %q, want %q", got.Kind, "task")
	}
	if got.Title != "Implement tracer bullet" {
		t.Fatalf("title = %q, want %q", got.Title, "Implement tracer bullet")
	}
	if got.Status != "todo" {
		t.Fatalf("status = %q, want %q", got.Status, "todo")
	}
	if got.StatusCategory != "unstarted" {
		t.Fatalf("status_category = %q, want %q", got.StatusCategory, "unstarted")
	}
	if got.WorkflowID != "task-default" {
		t.Fatalf("workflow_id = %q, want %q", got.WorkflowID, "task-default")
	}
	if got.OwnerRole != "developer" {
		t.Fatalf("owner_role = %q, want %q", got.OwnerRole, "developer")
	}
	if got.Labels == nil || len(got.Labels) != 0 {
		t.Fatalf("labels = %#v, want empty array", got.Labels)
	}
	if got.InInbox {
		t.Fatal("in_inbox = true, want false")
	}
	if got.ConcernsRaised {
		t.Fatal("concerns_raised = true, want false")
	}
	if got.OpenConcernsCount != 0 {
		t.Fatalf("open_concerns_count = %d, want 0", got.OpenConcernsCount)
	}
	if got.AcceptanceCriteria == nil || len(got.AcceptanceCriteria) != 0 {
		t.Fatalf("acceptance_criteria = %#v, want empty array", got.AcceptanceCriteria)
	}
}

func TestIssueCreateAppendsOneEventAndProjectsOneIssue(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	journalPath := mustJournalPath(t, workspace)
	beforePayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}
	beforeLines := splitJournalLines(beforePayload)

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Append one event", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	afterPayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after create error = %v", journalPath, err)
	}
	afterLines := splitJournalLines(afterPayload)
	if len(afterLines) != len(beforeLines)+1 {
		t.Fatalf("journal lines = %d, want %d", len(afterLines), len(beforeLines)+1)
	}

	var event journal.Event
	if err := json.Unmarshal([]byte(afterLines[len(afterLines)-1]), &event); err != nil {
		t.Fatalf("Unmarshal() event error = %v", err)
	}
	if event.EventType != "issue.created" {
		t.Fatalf("event_type = %q, want %q", event.EventType, "issue.created")
	}
	if event.EventID == "" || event.EventSchema != journal.EventSchemaV1 || event.Actor.Role != "developer" || len(event.Command.Argv) == 0 || len(event.Effects) != 1 {
		t.Fatalf("event envelope = %+v, want populated v1 envelope with one effect", event)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var projectedTitle string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, issueID).Scan(&projectedTitle); err != nil {
		t.Fatalf("query created issue: %v", err)
	}
	if projectedTitle != "Append one event" {
		t.Fatalf("projected title = %q, want %q", projectedTitle, "Append one event")
	}
}

func TestIssueCreateAndRebuildRoundTripReturnsIdenticalJSON(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Rebuild parity", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	beforeStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runFIDCommand(t, workspace, "rebuild", "--by", "developer")

	afterStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID)
	if beforeStdout != afterStdout {
		t.Fatalf("issue JSON after rebuild = %q, want %q", afterStdout, beforeStdout)
	}
}

func TestIssueCreateAndRebuildValidateRoleExistence(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	t.Run("issue create missing by", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "create", "--kind", "task", "--title", "Missing role")
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

	t.Run("issue create unknown role", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "create", "--kind", "task", "--title", "Unknown role", "--by", "intern")
		if err == nil {
			t.Fatal("Execute() error = nil, want unknown role rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
			t.Fatalf("stderr = %q, want unknown role error", stderr)
		}
	})

	t.Run("rebuild missing by", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "rebuild")
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

	t.Run("rebuild unknown role", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "rebuild", "--by", "intern")
		if err == nil {
			t.Fatal("Execute() error = nil, want unknown role rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
			t.Fatalf("stderr = %q, want unknown role error", stderr)
		}
	})
}

func TestIssueGetRejectsUnsupportedFormats(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Unsupported format", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	for _, formatName := range []string{"toon", "markdown", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "get", issueID, "--format", formatName)
			if err == nil {
				t.Fatal("Execute() error = nil, want unsupported format rejection")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, `not yet supported`) {
				t.Fatalf("stderr = %q, want not yet supported error", stderr)
			}
		})
	}
}

func runFIDCommand(t *testing.T, workspace string, args ...string) (string, string) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	defer mustChdir(t, cwd)
	mustChdir(t, workspace)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(%v) error = %v", args, err)
	}

	return stdout.String(), stderr.String()
}

func runFIDCommandExpectError(t *testing.T, workspace string, args ...string) (string, string, error) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	defer mustChdir(t, cwd)
	mustChdir(t, workspace)

	err = cmd.Execute()
	return stdout.String(), stderr.String(), err
}

func parseCreatedIssueID(t *testing.T, stderr string) string {
	t.Helper()

	fields := strings.Fields(stderr)
	for index := range fields {
		if fields[index] == "issue" && index+1 < len(fields) {
			return fields[index+1]
		}
	}
	t.Fatalf("stderr = %q, want created issue line with id", stderr)
	return ""
}
