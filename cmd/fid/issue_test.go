package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"gopkg.in/yaml.v3"
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
		TriageStatus       *string  `json:"triage_status"`
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
	if got.TriageStatus != nil {
		t.Fatalf("triage_status = %#v, want nil", got.TriageStatus)
	}
	if got.AcceptanceCriteria == nil || len(got.AcceptanceCriteria) != 0 {
		t.Fatalf("acceptance_criteria = %#v, want empty array", got.AcceptanceCriteria)
	}
}

func TestIssueCreateEpicAcceptsSchedulingAndReturnsDefaultReadShape(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Phase 4 authoring", "--scheduling", "sequential", "--by", "tl")
	issueID := parseCreatedIssueID(t, stderr)

	if !strings.Contains(stderr, fmt.Sprintf("created issue %s (epic) %q", issueID, "Phase 4 authoring")) {
		t.Fatalf("stderr = %q, want create confirmation", stderr)
	}

	stdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

	var got struct {
		ID             string   `json:"id"`
		Kind           string   `json:"kind"`
		Title          string   `json:"title"`
		Status         string   `json:"status"`
		StatusCategory string   `json:"status_category"`
		WorkflowID     string   `json:"workflow_id"`
		ParentID       *string  `json:"parent_id"`
		OwnerRole      *string  `json:"owner_role"`
		Labels         []string `json:"labels"`
		InInbox        bool     `json:"in_inbox"`
	}
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.ID != issueID {
		t.Fatalf("id = %q, want %q", got.ID, issueID)
	}
	if got.Kind != "epic" {
		t.Fatalf("kind = %q, want epic", got.Kind)
	}
	if got.Title != "Phase 4 authoring" {
		t.Fatalf("title = %q, want %q", got.Title, "Phase 4 authoring")
	}
	if got.Status != "todo" {
		t.Fatalf("status = %q, want todo", got.Status)
	}
	if got.StatusCategory != "unstarted" {
		t.Fatalf("status_category = %q, want unstarted", got.StatusCategory)
	}
	if got.WorkflowID != "epic-default" {
		t.Fatalf("workflow_id = %q, want epic-default", got.WorkflowID)
	}
	if got.ParentID != nil {
		t.Fatalf("parent_id = %#v, want nil", got.ParentID)
	}
	if got.OwnerRole == nil || *got.OwnerRole != "tl" {
		t.Fatalf("owner_role = %#v, want tl", got.OwnerRole)
	}
	if got.Labels == nil || len(got.Labels) != 0 {
		t.Fatalf("labels = %#v, want empty array", got.Labels)
	}
	if got.InInbox {
		t.Fatal("in_inbox = true, want false")
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var scheduling string
	if err := conn.QueryRow(`SELECT scheduling FROM epic_metadata WHERE issue_id = ?`, issueID).Scan(&scheduling); err != nil {
		t.Fatalf("query epic_metadata: %v", err)
	}
	if scheduling != "sequential" {
		t.Fatalf("epic_metadata scheduling = %q, want sequential", scheduling)
	}
}

func TestIssueCreateSupportsHierarchyAndAuthoringFlags(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	if _, err := conn.Exec(`INSERT INTO labels (name) VALUES (?), (?)`, "core", "phase-4"); err != nil {
		t.Fatalf("insert labels: %v", err)
	}

	_, epicStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Phase 4 authoring", "--scheduling", "sequential", "--by", "tl")
	epicID := parseCreatedIssueID(t, epicStderr)

	_, taskStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Implement authoring domain",
		"--epic", epicID,
		"--priority", "1",
		"--sequence", "1",
		"--spec-md", "project-docs/docs/phases/phase-4/tasks/task-2.md",
		"--est-tokens", "2100",
		"--due-date", "2026-06-30",
		"--label", "phase-4",
		"--label", "core",
		"--inbox",
		"--owner-role", "qa_lead",
		"--owner-actor", "human:mark",
		"--by", "tl",
	)
	taskID := parseCreatedIssueID(t, taskStderr)

	_, unattachedStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Inbox capture", "--by", "developer")
	unattachedTaskID := parseCreatedIssueID(t, unattachedStderr)

	_, subtaskStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "subtask", "--title", "Cover rebuild parity", "--parent", taskID, "--by", "developer")
	subtaskID := parseCreatedIssueID(t, subtaskStderr)

	type issueJSON struct {
		ID             string   `json:"id"`
		Kind           string   `json:"kind"`
		Title          string   `json:"title"`
		Status         string   `json:"status"`
		StatusCategory string   `json:"status_category"`
		WorkflowID     string   `json:"workflow_id"`
		ParentID       *string  `json:"parent_id"`
		Priority       *int     `json:"priority"`
		Sequence       *int     `json:"sequence"`
		SpecMDPath     *string  `json:"spec_md_path"`
		EstTokens      *int     `json:"est_tokens"`
		OwnerRole      *string  `json:"owner_role"`
		OwnerActor     *string  `json:"owner_actor"`
		Labels         []string `json:"labels"`
		InInbox        bool     `json:"in_inbox"`
	}

	readIssue := func(issueID string) issueJSON {
		t.Helper()

		stdout, stderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
		if stderr != "" {
			t.Fatalf("issue get %s stderr = %q, want empty", issueID, stderr)
		}

		var got issueJSON
		if err := json.Unmarshal([]byte(stdout), &got); err != nil {
			t.Fatalf("Unmarshal(%s) error = %v", issueID, err)
		}
		return got
	}

	epic := readIssue(epicID)
	if epic.Kind != "epic" || epic.WorkflowID != "epic-default" || epic.Status != "todo" || epic.StatusCategory != "unstarted" {
		t.Fatalf("epic read-back = %+v, want epic defaults", epic)
	}
	if epic.ParentID != nil {
		t.Fatalf("epic parent_id = %#v, want nil", epic.ParentID)
	}
	if epic.OwnerRole == nil || *epic.OwnerRole != "tl" {
		t.Fatalf("epic owner_role = %#v, want tl", epic.OwnerRole)
	}

	task := readIssue(taskID)
	if task.Kind != "task" || task.WorkflowID != "task-default" {
		t.Fatalf("task read-back = %+v, want task defaults", task)
	}
	if task.ParentID == nil || *task.ParentID != epicID {
		t.Fatalf("task parent_id = %#v, want %q", task.ParentID, epicID)
	}
	if task.Priority == nil || *task.Priority != 1 {
		t.Fatalf("task priority = %#v, want 1", task.Priority)
	}
	if task.Sequence == nil || *task.Sequence != 1 {
		t.Fatalf("task sequence = %#v, want 1", task.Sequence)
	}
	if task.SpecMDPath == nil || *task.SpecMDPath != "project-docs/docs/phases/phase-4/tasks/task-2.md" {
		t.Fatalf("task spec_md_path = %#v", task.SpecMDPath)
	}
	if task.EstTokens == nil || *task.EstTokens != 2100 {
		t.Fatalf("task est_tokens = %#v, want 2100", task.EstTokens)
	}
	if task.OwnerRole == nil || *task.OwnerRole != "qa_lead" {
		t.Fatalf("task owner_role = %#v, want qa_lead", task.OwnerRole)
	}
	if task.OwnerActor == nil || *task.OwnerActor != "human:mark" {
		t.Fatalf("task owner_actor = %#v, want human:mark", task.OwnerActor)
	}
	if !task.InInbox {
		t.Fatal("task in_inbox = false, want true")
	}
	if want := []string{"core", "phase-4"}; !reflect.DeepEqual(task.Labels, want) {
		t.Fatalf("task labels = %#v, want %#v", task.Labels, want)
	}

	unattachedTask := readIssue(unattachedTaskID)
	if unattachedTask.ParentID != nil {
		t.Fatalf("unattached task parent_id = %#v, want nil", unattachedTask.ParentID)
	}
	if unattachedTask.OwnerRole == nil || *unattachedTask.OwnerRole != "developer" {
		t.Fatalf("unattached task owner_role = %#v, want developer", unattachedTask.OwnerRole)
	}

	subtask := readIssue(subtaskID)
	if subtask.Kind != "subtask" || subtask.WorkflowID != "subtask-default" {
		t.Fatalf("subtask read-back = %+v, want subtask defaults", subtask)
	}
	if subtask.ParentID == nil || *subtask.ParentID != taskID {
		t.Fatalf("subtask parent_id = %#v, want %q", subtask.ParentID, taskID)
	}
	if subtask.OwnerRole == nil || *subtask.OwnerRole != "developer" {
		t.Fatalf("subtask owner_role = %#v, want developer", subtask.OwnerRole)
	}

	var dueDate string
	if err := conn.QueryRow(`SELECT date(due_date) FROM issues WHERE id = ?`, taskID).Scan(&dueDate); err != nil {
		t.Fatalf("query task due_date: %v", err)
	}
	if dueDate != "2026-06-30" {
		t.Fatalf("task due_date = %q, want 2026-06-30", dueDate)
	}
}

func TestIssueCreateRejectsInvalidFlagsAndConstraintViolations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    func(t *testing.T, workspace string) []string
		wantErr string
		verify  func(t *testing.T, workspace string)
	}{
		{
			name: "task rejects subtask parent flag",
			args: func(t *testing.T, workspace string) []string {
				_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Parent task", "--by", "developer")
				taskID := parseCreatedIssueID(t, stderr)
				return []string{"issue", "create", "--kind", "task", "--title", "Wrong flag", "--parent", taskID, "--by", "developer"}
			},
			wantErr: "create issue: --parent is only valid for subtasks",
		},
		{
			name: "subtask rejects task epic flag",
			args: func(t *testing.T, workspace string) []string {
				_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Sequential parent", "--scheduling", "sequential", "--by", "tl")
				epicID := parseCreatedIssueID(t, stderr)
				return []string{"issue", "create", "--kind", "subtask", "--title", "Wrong flag", "--epic", epicID, "--by", "developer"}
			},
			wantErr: "create issue: --epic is only valid for tasks",
		},
		{
			name: "task under sequential epic requires sequence",
			args: func(t *testing.T, workspace string) []string {
				_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Sequential parent", "--scheduling", "sequential", "--by", "tl")
				epicID := parseCreatedIssueID(t, stderr)
				return []string{"issue", "create", "--kind", "task", "--title", "Missing sequence", "--epic", epicID, "--by", "developer"}
			},
			wantErr: "issues: task under sequential epic requires sequence",
		},
		{
			name: "task under prioritised epic must not set sequence",
			args: func(t *testing.T, workspace string) []string {
				_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Prioritised parent", "--scheduling", "prioritised", "--by", "tl")
				epicID := parseCreatedIssueID(t, stderr)
				return []string{"issue", "create", "--kind", "task", "--title", "Prioritised sequence", "--epic", epicID, "--sequence", "1", "--by", "developer"}
			},
			wantErr: "issues: task under prioritised epic must not have sequence",
		},
		{
			name: "unattached task must not set sequence",
			args: func(t *testing.T, workspace string) []string {
				return []string{"issue", "create", "--kind", "task", "--title", "Unattached sequence", "--sequence", "1", "--by", "developer"}
			},
			wantErr: "issues: unattached task must not have sequence",
		},
		{
			name: "task under epic without metadata is rejected",
			args: func(t *testing.T, workspace string) []string {
				conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
				if err != nil {
					t.Fatalf("Open() error = %v", err)
				}
				defer closeCommandDB(t, conn)
				if _, err := conn.Exec(`
					INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)
				`,
					"FID-nometa",
					"epic",
					"Metadata missing epic",
					"todo",
					"epic-default",
					"tl",
					"2026-06-07T14:00:00Z",
					"2026-06-07T14:00:00Z",
				); err != nil {
					t.Fatalf("insert metadata-less epic: %v", err)
				}
				return []string{"issue", "create", "--kind", "task", "--title", "No metadata", "--epic", "FID-nometa", "--sequence", "1", "--by", "developer"}
			},
			wantErr: "issues: task parent epic must have epic_metadata",
		},
		{
			name: "unknown label rejected without auto-create",
			args: func(t *testing.T, workspace string) []string {
				seedCommandLabels(t, workspace, "known")
				return []string{"issue", "create", "--kind", "task", "--title", "Unknown label", "--label", "known", "--label", "missing", "--by", "developer"}
			},
			wantErr: `create issue: label "missing" does not exist`,
			verify: func(t *testing.T, workspace string) {
				conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
				if err != nil {
					t.Fatalf("Open() error = %v", err)
				}
				defer closeCommandDB(t, conn)

				var missingLabelCount int
				if err := conn.QueryRow(`SELECT COUNT(*) FROM labels WHERE name = ?`, "missing").Scan(&missingLabelCount); err != nil {
					t.Fatalf("count missing label rows: %v", err)
				}
				if missingLabelCount != 0 {
					t.Fatalf("missing label rows = %d, want 0", missingLabelCount)
				}

				var unknownLabelIssueCount int
				if err := conn.QueryRow(`SELECT COUNT(*) FROM issues WHERE title = ?`, "Unknown label").Scan(&unknownLabelIssueCount); err != nil {
					t.Fatalf("count unknown label issues: %v", err)
				}
				if unknownLabelIssueCount != 0 {
					t.Fatalf("unknown label issues = %d, want 0", unknownLabelIssueCount)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)

			stdout, stderr, err := runFIDCommandExpectError(t, workspace, tt.args(t, workspace)...)
			if err == nil {
				t.Fatalf("Execute() error = nil, want rejection")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, tt.wantErr) {
				t.Fatalf("stderr = %q, want substring %q", stderr, tt.wantErr)
			}
			if tt.verify != nil {
				tt.verify(t, workspace)
			}
		})
	}
}

func TestIssueCreateUsesOneJournalEventPerCommandAndRebuildsRepresentativeHierarchy(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	seedCommandLabels(t, workspace, "core", "phase-4")

	journalPath := mustJournalPath(t, workspace)
	beforePayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}
	beforeLines := splitJournalLines(beforePayload)

	_, epicStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Phase 4 authoring", "--scheduling", "sequential", "--by", "tl")
	epicID := parseCreatedIssueID(t, epicStderr)
	afterEpicLines := splitJournalLines(mustReadCommandFile(t, journalPath))
	if len(afterEpicLines) != len(beforeLines)+1 {
		t.Fatalf("journal lines after epic = %d, want %d", len(afterEpicLines), len(beforeLines)+1)
	}

	_, taskStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Implement authoring domain",
		"--epic", epicID,
		"--priority", "1",
		"--sequence", "1",
		"--spec-md", "project-docs/docs/phases/phase-4/tasks/task-2.md",
		"--est-tokens", "2100",
		"--due-date", "2026-06-30",
		"--label", "phase-4",
		"--label", "core",
		"--inbox",
		"--owner-role", "qa_lead",
		"--owner-actor", "human:mark",
		"--by", "tl",
	)
	taskID := parseCreatedIssueID(t, taskStderr)
	afterTaskLines := splitJournalLines(mustReadCommandFile(t, journalPath))
	if len(afterTaskLines) != len(afterEpicLines)+1 {
		t.Fatalf("journal lines after task = %d, want %d", len(afterTaskLines), len(afterEpicLines)+1)
	}

	_, subtaskStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "subtask", "--title", "Cover rebuild parity", "--parent", taskID, "--by", "developer")
	subtaskID := parseCreatedIssueID(t, subtaskStderr)
	afterSubtaskLines := splitJournalLines(mustReadCommandFile(t, journalPath))
	if len(afterSubtaskLines) != len(afterTaskLines)+1 {
		t.Fatalf("journal lines after subtask = %d, want %d", len(afterSubtaskLines), len(afterTaskLines)+1)
	}

	var latestEvent journal.Event
	if err := json.Unmarshal([]byte(afterSubtaskLines[len(afterSubtaskLines)-1]), &latestEvent); err != nil {
		t.Fatalf("Unmarshal(last event) error = %v", err)
	}
	if latestEvent.EventType != "issue.created" {
		t.Fatalf("last event type = %q, want issue.created", latestEvent.EventType)
	}

	beforeRebuild := map[string]string{
		epicID:    mustGetIssueJSON(t, workspace, epicID),
		taskID:    mustGetIssueJSON(t, workspace, taskID),
		subtaskID: mustGetIssueJSON(t, workspace, subtaskID),
	}

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")

	for issueID, want := range beforeRebuild {
		if got := mustGetIssueJSON(t, workspace, issueID); got != want {
			t.Fatalf("issue %s JSON after rebuild = %q, want %q", issueID, got, want)
		}
	}
}

func TestIssueGetReturnsEnrichedJSONAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Surface enriched issue state", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	appendCommandIssueEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "issue.enriched_state_seeded",
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
				Table: "issues",
				PK: map[string]any{
					"id": "FID-epiccmd",
				},
				Values: map[string]any{
					"kind":        "epic",
					"title":       "CLI enriched epic",
					"workflow_id": "epic-default",
					"status":      "todo",
					"created_at":  "2026-06-07T14:00:00Z",
					"updated_at":  "2026-06-07T14:00:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "task_completions",
				PK: map[string]any{
					"id": "TC-CMD-001",
				},
				Values: map[string]any{
					"issue_id":           issueID,
					"attempt_number":     1,
					"developer_role":     "developer",
					"model":              "gpt-5",
					"status":             "dev_complete",
					"gate_results":       `["go test ./..."]`,
					"human_intervention": false,
					"completed_at":       "2026-06-07T14:01:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-CMD-001",
				},
				Values: map[string]any{
					"issue_id":                issueID,
					"raised_in_completion_id": "TC-CMD-001",
					"raised_by_role":          "developer",
					"summary":                 "CLI concern",
					"body_md_path":            "project-docs/concerns/CON-CMD-001.md",
					"status":                  "open",
					"created_at":              "2026-06-07T14:01:30Z",
					"updated_at":              "2026-06-07T14:01:30Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-CMD-001",
				},
				Values: map[string]any{
					"epic_id":       "FID-epiccmd",
					"issue_id":      issueID,
					"triager_role":  "qa_lead",
					"source":        "task_completion",
					"source_ref_id": "TC-CMD-001",
					"verdict":       "RESUME",
					"reasoning_md":  "project-docs/triage/TRG-CMD-001.md",
					"created_at":    "2026-06-07T14:02:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-CMD-002",
				},
				Values: map[string]any{
					"epic_id":       "FID-epiccmd",
					"issue_id":      issueID,
					"triager_role":  "qa_lead",
					"source":        "task_completion",
					"source_ref_id": "TC-CMD-001",
					"verdict":       "BLOCK",
					"reasoning_md":  "project-docs/triage/TRG-CMD-002.md",
					"created_at":    "2026-06-07T14:03:00Z",
				},
			},
		},
	})

	beforeStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")

	var before struct {
		ConcernsRaised    bool    `json:"concerns_raised"`
		OpenConcernsCount int     `json:"open_concerns_count"`
		TriageStatus      *string `json:"triage_status"`
	}
	if err := json.Unmarshal([]byte(beforeStdout), &before); err != nil {
		t.Fatalf("Unmarshal(before) error = %v", err)
	}
	if !before.ConcernsRaised {
		t.Fatal("before concerns_raised = false, want true")
	}
	if before.OpenConcernsCount != 1 {
		t.Fatalf("before open_concerns_count = %d, want 1", before.OpenConcernsCount)
	}
	if before.TriageStatus == nil || *before.TriageStatus != "blocked" {
		t.Fatalf("before triage_status = %#v, want blocked", before.TriageStatus)
	}

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")

	afterStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if beforeStdout != afterStdout {
		t.Fatalf("issue JSON after rebuild = %q, want %q", afterStdout, beforeStdout)
	}
}

func TestIssueGetReturnsMarkdownForEmptyAndPopulatedIssues(t *testing.T) {
	t.Run("empty issue omits empty sections", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		defer closeCommandDB(t, conn)

		if _, err := conn.Exec(`
			INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
			"FID-emptymd",
			"task",
			"Render empty issue",
			"todo",
			"task-default",
			"developer",
			"2026-06-07T14:00:00Z",
			"2026-06-07T14:00:00Z",
		); err != nil {
			t.Fatalf("insert empty markdown issue: %v", err)
		}

		stdout, stderr := runFIDCommand(t, workspace, "issue", "get", "FID-emptymd", "--format", "markdown")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}

		want := "" +
			"# FID-emptymd: Render empty issue\n\n" +
			"**Kind:** task | **Status:** todo | **Priority:** (none)\n" +
			"**Owner:** developer\n\n" +
			"---\n" +
			"Created: 2026-06-07T14:00:00Z | Updated: 2026-06-07T14:00:00Z\n"
		if stdout != want {
			t.Fatalf("markdown stdout = %q, want %q", stdout, want)
		}
		if strings.Contains(stdout, "## Acceptance Criteria") || strings.Contains(stdout, "## Labels") || strings.Contains(stdout, "## Links") {
			t.Fatalf("markdown stdout = %q, want empty sections omitted", stdout)
		}
	})

	t.Run("populated issue matches markdown layout", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		defer closeCommandDB(t, conn)

		seedMarkdownPopulatedIssue(t, conn)

		stdout, stderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--format", "markdown")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}

		want := "" +
			"# FID-a1b2c3: Implement journal append with flush\n\n" +
			"**Kind:** task | **Status:** in_progress | **Priority:** P1\n" +
			"**Epic:** FID-e1f2g3 (Core persistence layer)\n" +
			"**Owner:** developer | **Sequence:** 3\n" +
			"**Spec:** project-docs/docs/phases/phase-1/tasks/FID-a1b2c3.md\n\n" +
			"## Acceptance Criteria\n\n" +
			"- [x] AC-d4e5f6: Journal file is flushed to disk before returning success\n" +
			"  - developer: pass (2026-05-27T10:15:00Z)\n" +
			"- [ ] AC-e5f6g7: Partial write leaves journal intact\n" +
			"  - developer: (not verified)\n\n" +
			"## Labels\n\n" +
			"core, journal\n\n" +
			"## Links\n\n" +
			"- blocks FID-7c3e9f (Schema migration runner)\n" +
			"- implements PRD-ISSUE-001\n\n" +
			"---\n" +
			"Created: 2026-05-25T09:00:00Z | Updated: 2026-05-27T10:20:00Z\n"
		if stdout != want {
			t.Fatalf("markdown stdout = %q, want %q", stdout, want)
		}
	})
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
	if event.EventID == "" || event.EventSchema != journal.EventSchemaV1 || event.Actor.Role != "developer" || len(event.Command.Argv) == 0 {
		t.Fatalf("event envelope = %+v, want populated v1 envelope", event)
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

	if got := countActivityLogRows(t, conn); got != 1 {
		t.Fatalf("activity_log rows after create = %d, want 1", got)
	}

	var (
		activityCount int
		activityEvent string
		actor         string
		details       sql.NullString
	)
	if err := conn.QueryRow(`
		SELECT COUNT(*), event, COALESCE(actor, ''), details
		FROM activity_log
		WHERE journal_event_id = ?
		GROUP BY event, actor, details
	`, event.EventID).Scan(&activityCount, &activityEvent, &actor, &details); err != nil {
		t.Fatalf("query activity_log row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if activityEvent != "issue_created" {
		t.Fatalf("activity_log event = %q, want issue_created", activityEvent)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	for _, want := range []string{`"kind":"task"`, `"title":"Append one event"`} {
		if !strings.Contains(details.String, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details.String, want)
		}
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

	runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")

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

	t.Run("issue create denied role leaves journal unchanged", func(t *testing.T) {
		journalPath := mustJournalPath(t, workspace)
		linesBefore := splitJournalLines(mustReadCommandFile(t, journalPath))

		stdout, stderr, err := runFIDCommandExpectError(
			t,
			workspace,
			"issue", "create",
			"--kind", "epic",
			"--title", "Denied epic create",
			"--scheduling", "sequential",
			"--by", "developer",
		)
		if err == nil {
			t.Fatal("Execute() error = nil, want permission rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		for _, want := range []string{
			`error: role "developer" not permitted for action "issue.create.epic"`,
			"check `fid policy show` for permitted role/action combinations",
		} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("stderr = %q, want substring %q", stderr, want)
			}
		}

		linesAfter := splitJournalLines(mustReadCommandFile(t, journalPath))
		if len(linesAfter) != len(linesBefore) {
			t.Fatalf("journal lines after rejection = %d, want unchanged %d", len(linesAfter), len(linesBefore))
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

func TestIssueGetStopsOnJournalCatchUpErrors(t *testing.T) {
	t.Run("malformed trailing line", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Catch-up corruption", "--by", "developer")
		issueID := parseCreatedIssueID(t, stderr)

		stdout, stderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
		if stderr != "" {
			t.Fatalf("stderr before corruption = %q, want empty", stderr)
		}
		if !strings.Contains(stdout, issueID) {
			t.Fatalf("stdout before corruption = %q, want issue payload", stdout)
		}

		appendRawJournalLine(t, mustJournalPath(t, workspace), `{"event_id":"EVT-20260607-999999"`)

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "get", issueID, "--format", "json")
		if err == nil {
			t.Fatal("Execute() error = nil, want corruption error")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		for _, want := range []string{
			"error: malformed journal entry at .fid/journal/",
			"parse error: unexpected end of JSON input",
			"read stopped - manual repair required",
		} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("stderr = %q, want substring %q", stderr, want)
			}
		}
	})

	t.Run("precondition conflict", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Catch-up conflict", "--by", "developer")
		issueID := parseCreatedIssueID(t, stderr)

		stdout, stderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
		if stderr != "" {
			t.Fatalf("stderr before conflict = %q, want empty", stderr)
		}
		if !strings.Contains(stdout, issueID) {
			t.Fatalf("stdout before conflict = %q, want issue payload", stdout)
		}

		appendCommandIssueEvent(t, workspace, journal.Event{
			EventSchema: journal.EventSchemaV1,
			EventType:   "issue.transitioned",
			Actor: journal.Actor{
				Role: "tl",
			},
			Command: journal.Command{
				Argv:       []string{"fid", "issue", "transition", issueID, "--to", "done", "--by", "tl"},
				CWD:        ".",
				FIDVersion: "0.1.0",
			},
			Effects: []journal.Effect{
				{
					Op:    "update",
					Table: "issues",
					PK: map[string]any{
						"id": issueID,
					},
					Preconditions: map[string]any{
						"status": "in_progress",
					},
					Values: map[string]any{
						"status": "done",
					},
				},
			},
		})

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "get", issueID, "--format", "json")
		if err == nil {
			t.Fatal("Execute() error = nil, want precondition conflict")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		for _, want := range []string{
			"error: replay precondition conflict at .fid/journal/",
			"event: EVT-",
			"issue.transitioned",
			"entity: issues pk={id:",
			"expected: {status:in_progress}",
			"actual: {status:todo}",
			"read stopped - manual repair required",
		} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("stderr = %q, want substring %q", stderr, want)
			}
		}
	})
}

func TestIssueGetReturnsYAMLForEmptyAndPopulatedIssues(t *testing.T) {
	t.Run("empty issue", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "YAML empty issue", "--by", "developer")
		issueID := parseCreatedIssueID(t, stderr)

		jsonStdout, jsonStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
		yamlStdout, yamlStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "yaml")

		if jsonStderr != "" {
			t.Fatalf("json stderr = %q, want empty", jsonStderr)
		}
		if yamlStderr != "" {
			t.Fatalf("yaml stderr = %q, want empty", yamlStderr)
		}

		want := decodeCommandJSONStructure(t, jsonStdout)
		got := decodeCommandYAMLStructure(t, yamlStdout)
		if !commandStructuresEqual(want, got) {
			t.Fatalf("YAML structure = %#v, want %#v", got, want)
		}
	})

	t.Run("populated issue", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		defer closeCommandDB(t, conn)

		seedMarkdownPopulatedIssue(t, conn)

		jsonStdout, jsonStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--format", "json")
		yamlStdout, yamlStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--format", "yaml")

		if jsonStderr != "" {
			t.Fatalf("json stderr = %q, want empty", jsonStderr)
		}
		if yamlStderr != "" {
			t.Fatalf("yaml stderr = %q, want empty", yamlStderr)
		}

		want := decodeCommandJSONStructure(t, jsonStdout)
		got := decodeCommandYAMLStructure(t, yamlStdout)
		if !commandStructuresEqual(want, got) {
			t.Fatalf("YAML structure = %#v, want %#v", got, want)
		}
	})
}

func TestIssueGetDispatchesStructuredFormatsAndRejectsUnknownFormats(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Unsupported format", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	for _, formatName := range []string{"toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			jsonStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
			formattedStdout, formattedStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", formatName)

			if formattedStderr != "" {
				t.Fatalf("stderr = %q, want empty", formattedStderr)
			}

			want := decodeCommandJSONStructure(t, jsonStdout)
			var got any
			switch formatName {
			case "toon":
				got = decodeCommandTOONStructure(t, formattedStdout)
			case "yaml":
				got = decodeCommandYAMLStructure(t, formattedStdout)
			default:
				t.Fatalf("unexpected format %q", formatName)
			}

			if !commandStructuresEqual(want, got) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, got, want)
			}
		})
	}

	t.Run("xml", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "get", issueID, "--format", "xml")
		if err == nil {
			t.Fatal("Execute() error = nil, want unsupported format rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `unsupported format "xml"; known formats: json, toon, markdown, yaml`) {
			t.Fatalf("stderr = %q, want actionable unknown format error", stderr)
		}
	})
}

func seedMarkdownPopulatedIssue(t *testing.T, conn *sql.DB) {
	t.Helper()

	statements := []struct {
		query string
		args  []any
	}{
		{
			query: `
				INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-e1f2g3",
				"epic",
				"Core persistence layer",
				"todo",
				"epic-default",
				"developer",
				"2026-05-25T08:00:00Z",
				"2026-05-25T08:00:00Z",
			},
		},
		{
			query: `
				INSERT INTO epic_metadata (issue_id, scheduling)
				VALUES (?, ?)
			`,
			args: []any{
				"FID-e1f2g3",
				"sequential",
			},
		},
		{
			query: `
				INSERT INTO issues (id, kind, parent_id, title, workflow_id, status, priority, sequence, spec_md_path, owner_role, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-a1b2c3",
				"task",
				"FID-e1f2g3",
				"Implement journal append with flush",
				"task-default",
				"in_progress",
				1,
				3,
				"project-docs/docs/phases/phase-1/tasks/FID-a1b2c3.md",
				"developer",
				"2026-05-25T09:00:00Z",
				"2026-05-27T10:20:00Z",
			},
		},
		{
			query: `
				INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-7c3e9f",
				"task",
				"Schema migration runner",
				"todo",
				"task-default",
				"developer",
				"2026-05-25T09:05:00Z",
				"2026-05-25T09:05:00Z",
			},
		},
		{
			query: `
				INSERT INTO requirements (id, type, summary, source_doc, source_section, status)
				VALUES (?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"PRD-ISSUE-001",
				"PRD",
				"Implement journal append with flush",
				"01_PRD.md",
				"ISSUE",
				"active",
			},
		},
		{
			query: `INSERT INTO labels (name) VALUES (?), (?)`,
			args:  []any{"core", "journal"},
		},
		{
			query: `INSERT INTO issue_labels (issue_id, label) VALUES (?, ?), (?, ?)`,
			args:  []any{"FID-a1b2c3", "core", "FID-a1b2c3", "journal"},
		},
		{
			query: `
				INSERT INTO acceptance_criteria (id, issue_id, ordinal, text)
				VALUES (?, ?, ?, ?), (?, ?, ?, ?)
			`,
			args: []any{
				"AC-d4e5f6",
				"FID-a1b2c3",
				1,
				"Journal file is flushed to disk before returning success",
				"AC-e5f6g7",
				"FID-a1b2c3",
				2,
				"Partial write leaves journal intact",
			},
		},
		{
			query: `
				INSERT INTO ac_verifications (id, ac_id, verifier_role, status, verified_at)
				VALUES (?, ?, ?, ?, ?)
			`,
			args: []any{
				"ACV-d4e5f6-dev",
				"AC-d4e5f6",
				"developer",
				"pass",
				"2026-05-27T10:15:00Z",
			},
		},
		{
			query: `
				INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
				VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"issue",
				"FID-a1b2c3",
				"blocks",
				"issue",
				"FID-7c3e9f",
				"2026-05-27T10:18:00Z",
				"issue",
				"FID-a1b2c3",
				"implements",
				"requirement",
				"PRD-ISSUE-001",
				"2026-05-27T10:19:00Z",
			},
		},
	}

	for _, statement := range statements {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("Exec(%q) error = %v", statement.query, err)
		}
	}
}

func runFIDCommand(t *testing.T, workspace string, args ...string) (string, string) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(%v) error = %v", args, err)
	}

	return stdout.String(), stderr.String()
}

func runFIDCommandExpectError(t *testing.T, workspace string, args ...string) (string, string, error) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	err := cmd.Execute()
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

func mustReadCommandFile(t *testing.T, path string) []byte {
	t.Helper()

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	return payload
}

func mustGetIssueJSON(t *testing.T, workspace, issueID string) string {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if stderr != "" {
		t.Fatalf("issue get %s stderr = %q, want empty", issueID, stderr)
	}
	return stdout
}

func seedCommandLabels(t *testing.T, workspace string, labels ...string) {
	t.Helper()

	effects := make([]journal.Effect, 0, len(labels))
	for _, label := range labels {
		effects = append(effects, journal.Effect{
			Op:    "insert",
			Table: "labels",
			PK:    map[string]any{"name": label},
			Values: map[string]any{
				"description": label + " label",
			},
		})
	}

	appendCommandIssueEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "label.seeded",
		Actor: journal.Actor{
			Role: "tl",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: effects,
	})
}

func appendCommandIssueEvent(t *testing.T, workspace string, event journal.Event) {
	t.Helper()

	if err := journal.WithLock(workspace, journal.LockOwner{
		Command:   "fid maintenance seed",
		ActorRole: event.Actor.Role,
	}, func(lock *journal.WriteLock) error {
		_, err := journal.AppendEvent(lock, event)
		return err
	}); err != nil {
		t.Fatalf("AppendEvent() error = %v", err)
	}
}

func appendRawJournalLine(t *testing.T, path, line string) {
	t.Helper()

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("OpenFile(%q) error = %v", path, err)
	}
	defer closeCommandFile(t, file)

	if _, err := file.WriteString(line); err != nil {
		t.Fatalf("WriteString(%q) error = %v", path, err)
	}
	if err := file.Sync(); err != nil {
		t.Fatalf("Sync(%q) error = %v", path, err)
	}
}

func decodeCommandJSONStructure(t *testing.T, value string) any {
	t.Helper()

	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.UseNumber()

	var result any
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("Decode JSON structure error = %v", err)
	}
	return result
}

func decodeCommandTOONStructure(t *testing.T, value string) any {
	t.Helper()

	parser := commandTOONParser{input: strings.TrimSpace(value)}
	result, err := parser.parseValue()
	if err != nil {
		t.Fatalf("parse TOON structure error = %v", err)
	}
	parser.skipWhitespace()
	if parser.position != len(parser.input) {
		t.Fatalf("parse TOON structure stopped at %d of %d", parser.position, len(parser.input))
	}
	return result
}

func decodeCommandYAMLStructure(t *testing.T, value string) any {
	t.Helper()

	var result any
	if err := yaml.Unmarshal([]byte(value), &result); err != nil {
		t.Fatalf("Unmarshal YAML structure error = %v", err)
	}
	return result
}

func commandStructuresEqual(left any, right any) bool {
	switch leftValue := left.(type) {
	case map[string]any:
		rightValue, ok := right.(map[string]any)
		if !ok || len(leftValue) != len(rightValue) {
			return false
		}
		for key, nestedLeft := range leftValue {
			nestedRight, ok := rightValue[key]
			if !ok || !commandStructuresEqual(nestedLeft, nestedRight) {
				return false
			}
		}
		return true
	case []any:
		rightValue, ok := right.([]any)
		if !ok || len(leftValue) != len(rightValue) {
			return false
		}
		for index := range leftValue {
			if !commandStructuresEqual(leftValue[index], rightValue[index]) {
				return false
			}
		}
		return true
	case json.Number:
		return commandNumericStructuresEqual(leftValue, right)
	default:
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
	}
}

func closeCommandFile(t *testing.T, file *os.File) {
	t.Helper()

	if err := file.Close(); err != nil {
		t.Fatalf("Close(%q) error = %v", file.Name(), err)
	}
}

func commandNumericStructuresEqual(left json.Number, right any) bool {
	switch rightValue := right.(type) {
	case json.Number:
		return left.String() == rightValue.String()
	case int:
		return left.String() == strconv.Itoa(rightValue)
	case int8:
		return left.String() == strconv.FormatInt(int64(rightValue), 10)
	case int16:
		return left.String() == strconv.FormatInt(int64(rightValue), 10)
	case int32:
		return left.String() == strconv.FormatInt(int64(rightValue), 10)
	case int64:
		return left.String() == strconv.FormatInt(rightValue, 10)
	case uint:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint8:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint16:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint32:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint64:
		return left.String() == strconv.FormatUint(rightValue, 10)
	case float32:
		return left.String() == strconv.FormatFloat(float64(rightValue), 'f', -1, 32)
	case float64:
		return left.String() == strconv.FormatFloat(rightValue, 'f', -1, 64)
	default:
		return false
	}
}

type commandTOONParser struct {
	input    string
	position int
}

func (parser *commandTOONParser) parseValue() (any, error) {
	parser.skipWhitespace()
	if parser.position >= len(parser.input) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	switch parser.input[parser.position] {
	case '{':
		return parser.parseObject()
	case '[':
		return parser.parseArray()
	case '"':
		return parser.parseString()
	default:
		return parser.parseScalar()
	}
}

func (parser *commandTOONParser) parseObject() (map[string]any, error) {
	parser.position++
	result := map[string]any{}

	parser.skipWhitespace()
	if parser.consumeIf('}') {
		return result, nil
	}

	for {
		key, err := parser.parseKey()
		if err != nil {
			return nil, err
		}
		if !parser.consumeIf(':') {
			return nil, fmt.Errorf("expected ':' after key %q", key)
		}
		value, err := parser.parseValue()
		if err != nil {
			return nil, err
		}
		result[key] = value

		parser.skipWhitespace()
		if parser.consumeIf('}') {
			return result, nil
		}
		if !parser.consumeIf(',') {
			return nil, fmt.Errorf("expected ',' or '}' after key %q", key)
		}
	}
}

func (parser *commandTOONParser) parseArray() ([]any, error) {
	parser.position++
	var result []any

	parser.skipWhitespace()
	if parser.consumeIf(']') {
		return result, nil
	}

	for {
		value, err := parser.parseValue()
		if err != nil {
			return nil, err
		}
		result = append(result, value)

		parser.skipWhitespace()
		if parser.consumeIf(']') {
			return result, nil
		}
		if !parser.consumeIf(',') {
			return nil, fmt.Errorf("expected ',' or ']' in array")
		}
	}
}

func (parser *commandTOONParser) parseKey() (string, error) {
	parser.skipWhitespace()
	if parser.position >= len(parser.input) {
		return "", fmt.Errorf("unexpected end of input")
	}
	if parser.input[parser.position] == '"' {
		return parser.parseString()
	}

	start := parser.position
	for parser.position < len(parser.input) {
		switch parser.input[parser.position] {
		case ':', ' ', '\n', '\t', '\r':
			return parser.input[start:parser.position], nil
		default:
			parser.position++
		}
	}
	return parser.input[start:parser.position], nil
}

func (parser *commandTOONParser) parseString() (string, error) {
	start := parser.position
	parser.position++
	escaped := false
	for parser.position < len(parser.input) {
		current := parser.input[parser.position]
		if escaped {
			escaped = false
			parser.position++
			continue
		}
		if current == '\\' {
			escaped = true
			parser.position++
			continue
		}
		if current == '"' {
			parser.position++
			var result string
			if err := json.Unmarshal([]byte(parser.input[start:parser.position]), &result); err != nil {
				return "", err
			}
			return result, nil
		}
		parser.position++
	}
	return "", fmt.Errorf("unterminated string")
}

func (parser *commandTOONParser) parseScalar() (any, error) {
	start := parser.position
	for parser.position < len(parser.input) {
		switch parser.input[parser.position] {
		case ',', '}', ']', ' ', '\n', '\t', '\r':
			token := parser.input[start:parser.position]
			return parseCommandTOONScalar(token)
		default:
			parser.position++
		}
	}
	return parseCommandTOONScalar(parser.input[start:parser.position])
}

func (parser *commandTOONParser) skipWhitespace() {
	for parser.position < len(parser.input) {
		switch parser.input[parser.position] {
		case ' ', '\n', '\t', '\r':
			parser.position++
		default:
			return
		}
	}
}

func (parser *commandTOONParser) consumeIf(expected byte) bool {
	parser.skipWhitespace()
	if parser.position < len(parser.input) && parser.input[parser.position] == expected {
		parser.position++
		return true
	}
	return false
}

func parseCommandTOONScalar(token string) (any, error) {
	switch token {
	case "null":
		return nil, nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	}

	if _, err := strconv.ParseInt(token, 10, 64); err == nil {
		return json.Number(token), nil
	}
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return json.Number(token), nil
	}
	return token, nil
}
