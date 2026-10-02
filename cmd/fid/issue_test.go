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
)

func TestIssueGetHumanSuccessUsesStdout(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Human read", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, stderr)
	result := executeCommandForTest(workspace, "issue", "get", issueID)
	if result.err != nil || result.stderr != "" || !strings.Contains(result.stdout, issueID) || strings.Contains(result.stdout, ".fid/") {
		t.Fatalf("human read = %#v", result)
	}
}

func TestIssueCreateEpicAcceptsSchedulingAndReturnsDefaultReadShape(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Phase 4 authoring", "--scheduling", "sequential", "--actor", "agent:tl")
	issueID := parseCreatedIssueID(t, stderr)

	if !strings.Contains(stderr, fmt.Sprintf("created issue %s (epic) %q", issueID, "Phase 4 authoring")) {
		t.Fatalf("stderr = %q, want create confirmation", stderr)
	}

	stdout := mustGetIssueJSON(t, workspace, issueID)

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
	if got.OwnerRole != nil {
		t.Fatalf("canonical owner_role = %#v, want field absent", got.OwnerRole)
	}
	if gotRole := mustIssueOwnerRole(t, workspace, issueID); gotRole != "tl" {
		t.Fatalf("stored owner_role = %q, want tl", gotRole)
	}
	if got.Labels == nil || len(got.Labels) != 0 {
		t.Fatalf("labels = %#v, want empty array", got.Labels)
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
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	if _, err := conn.Exec(`INSERT INTO labels (name) VALUES (?), (?)`, "core", "phase-4"); err != nil {
		t.Fatalf("insert labels: %v", err)
	}

	_, epicStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Phase 4 authoring", "--scheduling", "sequential", "--actor", "agent:tl")
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
		"--label", "phase-4",
		"--label", "core",
		"--actor", "agent:tl",
	)
	taskID := parseCreatedIssueID(t, taskStderr)

	_, unattachedStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Inbox capture", "--actor", "agent:developer")
	unattachedTaskID := parseCreatedIssueID(t, unattachedStderr)

	_, subtaskStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "subtask", "--title", "Cover rebuild parity", "--parent", taskID, "--actor", "agent:developer")
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
		OwnerRole      *string  `json:"owner_role"`
		OwnerActor     *string  `json:"owner_actor"`
		Labels         []string `json:"labels"`
	}

	readIssue := func(issueID string) issueJSON {
		t.Helper()

		stdout := mustGetIssueJSON(t, workspace, issueID)

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
	if epic.OwnerRole != nil {
		t.Fatalf("epic canonical owner_role = %#v, want field absent", epic.OwnerRole)
	}
	if gotRole := mustIssueOwnerRole(t, workspace, epicID); gotRole != "tl" {
		t.Fatalf("epic stored owner_role = %q, want tl", gotRole)
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
	if task.OwnerRole != nil || task.OwnerActor != nil {
		t.Fatalf("task canonical owner fields = (%#v, %#v), want fields absent", task.OwnerRole, task.OwnerActor)
	}
	if gotRole := mustIssueOwnerRole(t, workspace, taskID); gotRole != "developer" {
		t.Fatalf("task stored owner_role = %q, want developer", gotRole)
	}
	if want := []string{"core", "phase-4"}; !reflect.DeepEqual(task.Labels, want) {
		t.Fatalf("task labels = %#v, want %#v", task.Labels, want)
	}

	unattachedTask := readIssue(unattachedTaskID)
	if unattachedTask.ParentID != nil {
		t.Fatalf("unattached task parent_id = %#v, want nil", unattachedTask.ParentID)
	}
	if unattachedTask.OwnerRole != nil {
		t.Fatalf("unattached task canonical owner_role = %#v, want field absent", unattachedTask.OwnerRole)
	}
	if gotRole := mustIssueOwnerRole(t, workspace, unattachedTaskID); gotRole != "developer" {
		t.Fatalf("unattached task stored owner_role = %q, want developer", gotRole)
	}

	subtask := readIssue(subtaskID)
	if subtask.Kind != "subtask" || subtask.WorkflowID != "subtask-default" {
		t.Fatalf("subtask read-back = %+v, want subtask defaults", subtask)
	}
	if subtask.ParentID == nil || *subtask.ParentID != taskID {
		t.Fatalf("subtask parent_id = %#v, want %q", subtask.ParentID, taskID)
	}
	if subtask.OwnerRole != nil {
		t.Fatalf("subtask canonical owner_role = %#v, want field absent", subtask.OwnerRole)
	}
	if gotRole := mustIssueOwnerRole(t, workspace, subtaskID); gotRole != "developer" {
		t.Fatalf("subtask stored owner_role = %q, want developer", gotRole)
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
				_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Parent task", "--actor", "agent:developer")
				taskID := parseCreatedIssueID(t, stderr)
				return []string{"issue", "create", "--kind", "task", "--title", "Wrong flag", "--parent", taskID, "--actor", "agent:developer"}
			},
			wantErr: "create issue: --parent is only valid for subtasks",
		},
		{
			name: "subtask rejects task epic flag",
			args: func(t *testing.T, workspace string) []string {
				_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Sequential parent", "--scheduling", "sequential", "--actor", "agent:tl")
				epicID := parseCreatedIssueID(t, stderr)
				return []string{"issue", "create", "--kind", "subtask", "--title", "Wrong flag", "--epic", epicID, "--actor", "agent:developer"}
			},
			wantErr: "create issue: --epic is only valid for tasks",
		},
		{
			name: "task under sequential epic requires sequence",
			args: func(t *testing.T, workspace string) []string {
				_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Sequential parent", "--scheduling", "sequential", "--actor", "agent:tl")
				epicID := parseCreatedIssueID(t, stderr)
				return []string{"issue", "create", "--kind", "task", "--title", "Missing sequence", "--epic", epicID, "--actor", "agent:developer"}
			},
			wantErr: "issues: task under sequential epic requires sequence",
		},
		{
			name: "task under prioritised epic must not set sequence",
			args: func(t *testing.T, workspace string) []string {
				_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Prioritised parent", "--scheduling", "prioritised", "--actor", "agent:tl")
				epicID := parseCreatedIssueID(t, stderr)
				return []string{"issue", "create", "--kind", "task", "--title", "Prioritised sequence", "--epic", epicID, "--sequence", "1", "--actor", "agent:developer"}
			},
			wantErr: "issues: task under prioritised epic must not have sequence",
		},
		{
			name: "unattached task must not set sequence",
			args: func(t *testing.T, workspace string) []string {
				return []string{"issue", "create", "--kind", "task", "--title", "Unattached sequence", "--sequence", "1", "--actor", "agent:developer"}
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
				return []string{"issue", "create", "--kind", "task", "--title", "No metadata", "--epic", "FID-nometa", "--sequence", "1", "--actor", "agent:developer"}
			},
			wantErr: "issues: task parent epic must have epic_metadata",
		},
		{
			name: "unknown label rejected without auto-create",
			args: func(t *testing.T, workspace string) []string {
				seedCommandLabels(t, workspace, "known")
				return []string{"issue", "create", "--kind", "task", "--title", "Unknown label", "--label", "known", "--label", "missing", "--actor", "agent:developer"}
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
	t.Parallel()
	workspace := clonePristineWorkspace(t)

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

	_, epicStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Phase 4 authoring", "--scheduling", "sequential", "--actor", "agent:tl")
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
		"--label", "phase-4",
		"--label", "core",
		"--actor", "agent:tl",
	)
	taskID := parseCreatedIssueID(t, taskStderr)
	afterTaskLines := splitJournalLines(mustReadCommandFile(t, journalPath))
	if len(afterTaskLines) != len(afterEpicLines)+1 {
		t.Fatalf("journal lines after task = %d, want %d", len(afterTaskLines), len(afterEpicLines)+1)
	}

	_, subtaskStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "subtask", "--title", "Cover rebuild parity", "--parent", taskID, "--actor", "agent:developer")
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

	runFIDCommand(t, workspace, "rebuild")

	for issueID, want := range beforeRebuild {
		if got := mustGetIssueJSON(t, workspace, issueID); got != want {
			t.Fatalf("issue %s JSON after rebuild = %q, want %q", issueID, got, want)
		}
	}
}

func TestIssueGetReturnsEnrichedJSONAndRebuildParity(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Surface enriched issue state", "--actor", "agent:developer")
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
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-CMD-001",
				},
				Values: map[string]any{
					"subject_type": "issue", "subject_id": issueID,
					"actor_id": "developer", "scope_json": `{"paths":[]}`,
					"severity": "major", "summary": "CLI concern", "reasoning": "Needs review",
					"references_json": "[]", "status": "open",
					"created_at": "2026-06-07T14:01:30Z", "created_event_id": "EVT-CON-CMD-001",
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
					"source":        "manual",
					"source_ref_id": "REV-CMD-001",
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
					"source":        "manual",
					"source_ref_id": "REV-CMD-001",
					"verdict":       "BLOCK",
					"reasoning_md":  "project-docs/triage/TRG-CMD-002.md",
					"created_at":    "2026-06-07T14:03:00Z",
				},
			},
		},
	})

	beforeStdout := mustGetIssueJSON(t, workspace, issueID)
	beforeConcerns := mustGetConcernListJSONForIssue(t, workspace, issueID)
	if len(beforeConcerns) != 1 || beforeConcerns[0].Status != "open" {
		t.Fatalf("before concerns = %#v, want one open concern", beforeConcerns)
	}
	if got := mustIssueTriageStatus(t, workspace, issueID); got == nil || *got != "blocked" {
		t.Fatalf("before triage_status = %#v, want blocked", got)
	}

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runFIDCommand(t, workspace, "rebuild")

	afterStdout := mustGetIssueJSON(t, workspace, issueID)
	if beforeStdout != afterStdout {
		t.Fatalf("issue JSON after rebuild = %q, want %q", afterStdout, beforeStdout)
	}
	if afterConcerns := mustGetConcernListJSONForIssue(t, workspace, issueID); !reflect.DeepEqual(afterConcerns, beforeConcerns) {
		t.Fatalf("concerns after rebuild = %#v, want %#v", afterConcerns, beforeConcerns)
	}
	if got := mustIssueTriageStatus(t, workspace, issueID); got == nil || *got != "blocked" {
		t.Fatalf("after triage_status = %#v, want blocked", got)
	}
}

func TestIssueGetRejectsRetiredFormatFlags(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Format retirement", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, stderr)
	for _, flag := range []string{"--format", "--section"} {
		args := []string{"issue", "get", issueID, flag, "json"}
		if flag == "--section" {
			args[4] = "acceptance"
		}
		if result := executeCommandForTest(workspace, args...); result.err == nil {
			t.Fatalf("%s was accepted", flag)
		}
	}
}

func TestIssueCreateAppendsOneEventAndProjectsOneIssue(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	journalPath := mustJournalPath(t, workspace)
	beforePayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}
	beforeLines := splitJournalLines(beforePayload)

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Append one event", "--actor", "agent:developer")
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
	if event.EventID == "" || event.EventSchema != journal.EventSchemaV2 || event.ActorV2.ActorID != "developer" || event.ActorV2.ActorType != "agent" || event.ActorV2.Producer != "fid-standalone" || event.ActorV2.ProvenanceRef != "project-local-attribution" || event.Aggregate.ID != issueID || event.Aggregate.Sequence != 1 || event.CommandName != "issue.create" || event.Idempotency == nil {
		t.Fatalf("event envelope = %+v, want populated v2 standalone Issue Event", event)
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
	if activityEvent != "issue-created" {
		t.Fatalf("activity_log event = %q, want issue-created", activityEvent)
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
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Rebuild parity", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, stderr)

	beforeStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runFIDCommand(t, workspace, "rebuild")

	afterStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID)
	if beforeStdout != afterStdout {
		t.Fatalf("issue JSON after rebuild = %q, want %q", afterStdout, beforeStdout)
	}
}

func TestIssueCreateStandaloneActorIsAttributionOnly(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	t.Run("issue create missing actor", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "create", "--kind", "task", "--title", "Missing role")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --actor rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: --actor is required") {
			t.Fatalf("stderr = %q, want missing --actor error", stderr)
		}
	})

	t.Run("unregistered actor is recorded as attribution", func(t *testing.T) {
		_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Unregistered actor", "--actor", "agent:intern")
		issueID := parseCreatedIssueID(t, stderr)
		tail := lastCommandJournalEvent(t, mustJournalPath(t, workspace))
		if tail.ActorV2.ActorID != "intern" || tail.ActorV2.ActorType != "agent" || tail.Aggregate.ID != issueID {
			t.Fatalf("Event Actor/aggregate = (%#v, %#v), want agent intern attribution for %q", tail.ActorV2, tail.Aggregate, issueID)
		}
	})

	t.Run("actor does not select standalone create authority", func(t *testing.T) {
		_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Attributed epic create", "--scheduling", "sequential", "--actor", "agent:developer")
		if issueID := parseCreatedIssueID(t, stderr); issueID == "" {
			t.Fatal("attributed epic create did not return an Issue ID")
		}
	})

}

func TestIssueGetStopsOnJournalCatchUpErrors(t *testing.T) {
	t.Parallel()
	t.Run("malformed trailing line", func(t *testing.T) {
		workspace := clonePristineWorkspace(t)

		_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Catch-up corruption", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, stderr)

		stdout := mustGetIssueJSON(t, workspace, issueID)
		if !strings.Contains(stdout, issueID) {
			t.Fatalf("stdout before corruption = %q, want issue payload", stdout)
		}

		appendRawJournalLine(t, mustJournalPath(t, workspace), `{"event_id":"EVT-20260607-999999"`)

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "--json", "issue", "get", issueID)
		if err == nil {
			t.Fatal("Execute() error = nil, want corruption error")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `"code":"journal-replay-stopped"`) {
			t.Fatalf("stderr = %q, want journal-replay-stopped envelope", stderr)
		}
	})

	t.Run("precondition conflict", func(t *testing.T) {
		workspace := clonePristineWorkspace(t)

		_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Catch-up conflict", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, stderr)

		stdout := mustGetIssueJSON(t, workspace, issueID)
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
				Argv:       []string{"fid", "issue", "transition", issueID, "--to", "done", "--actor", "agent:tl"},
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

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "--json", "issue", "get", issueID)
		if err == nil {
			t.Fatal("Execute() error = nil, want precondition conflict")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `"code":"journal-replay-stopped"`) {
			t.Fatalf("stderr = %q, want journal-replay-stopped envelope", stderr)
		}
	})
}

func TestIssueGetRejectsAllFormatValues(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Format rejection", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, stderr)
	for _, value := range []string{"json", "toon", "yaml", "markdown", "xml"} {
		if result := executeCommandForTest(workspace, "issue", "get", issueID, "--format", value); result.err == nil {
			t.Fatalf("--format %s was accepted", value)
		}
	}
}

func runFIDCommand(t *testing.T, workspace string, args ...string) (string, string) {
	t.Helper()
	args = withIssueMutationTestKey(args)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(%v) error = %v", args, err)
	}

	if len(args) >= 2 && args[0] == "issue" && (args[1] == "create" || args[1] == "update" || args[1] == "transition") {
		// Pre-convergence callers shared this helper expecting mutation
		// confirmations on stderr. Canonical mutations now write Human success to
		// stdout; focused transport tests exercise the real streams directly.
		return "", stdout.String() + stderr.String()
	}
	return stdout.String(), stderr.String()
}

func runFIDCommandExpectError(t *testing.T, workspace string, args ...string) (string, string, error) {
	t.Helper()
	args = withIssueMutationTestKey(args)

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

	result := executeCommandForTest(workspace, "--json", "issue", "get", issueID)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("issue get %s = %#v, want JSON stdout success", issueID, result)
	}
	payload, err := json.Marshal(decodeReadEnvelope(t, result.stdout, "issue.get").Result["issue"])
	if err != nil {
		t.Fatalf("Marshal(issue.get result) error = %v", err)
	}
	return string(payload)
}

type issueConcernListItem struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func mustGetConcernListJSONForIssue(t *testing.T, workspace, issueID string) []issueConcernListItem {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "--json", "concern", "list", "--subject", "issue:"+issueID)
	if stderr != "" {
		t.Fatalf("concern list for %s stderr = %q, want empty", issueID, stderr)
	}
	var envelope struct {
		Result struct {
			Items []issueConcernListItem `json:"items"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(concern list for %s) error = %v", issueID, err)
	}
	return envelope.Result.Items
}

func mustIssueTriageStatus(t *testing.T, workspace, issueID string) *string {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var status sql.NullString
	if err := conn.QueryRow(`SELECT triage_status FROM issues_enriched WHERE id = ?`, issueID).Scan(&status); err != nil {
		t.Fatalf("query triage_status for %s: %v", issueID, err)
	}
	if !status.Valid {
		return nil
	}
	return &status.String
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
		Command: "fid maintenance seed",
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

func decodeCanonicalResult(t *testing.T, value string) map[string]any {
	t.Helper()

	envelope, ok := decodeCommandJSONStructure(t, value).(map[string]any)
	if !ok {
		t.Fatalf("canonical response = %#v, want envelope object", value)
	}
	result, ok := envelope["result"].(map[string]any)
	if !ok {
		t.Fatalf("canonical response result = %#v, want result object", envelope["result"])
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
