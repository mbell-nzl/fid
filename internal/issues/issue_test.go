package issues

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestGetReturnsNotFoundErrorForMissingIssue(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	_, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            "FID-missing",
	})
	if err == nil {
		t.Fatal("Get() error = nil, want not found")
	}

	var notFoundErr *NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("Get() error = %T, want *NotFoundError", err)
	}
	if notFoundErr.IssueID != "FID-missing" {
		t.Fatalf("NotFoundError.IssueID = %q, want %q", notFoundErr.IssueID, "FID-missing")
	}
	if notFoundErr.Error() != `issue "FID-missing" not found` {
		t.Fatalf("NotFoundError.Error() = %q, want %q", notFoundErr.Error(), `issue "FID-missing" not found`)
	}
}

func TestCreateAndGetValidateInputGuards(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		call    func() error
		wantErr string
	}{
		{
			name: "create requires workspace root",
			call: func() error {
				_, err := Create(CreateOptions{
					MigrationsDir: testIssueMigrationsDir(t),
					Kind:          "task",
					Title:         "Implement coverage",
					ActorRole:     "developer",
				})
				return err
			},
			wantErr: "create issue: workspace root is required",
		},
		{
			name: "create requires migrations dir",
			call: func() error {
				_, err := Create(CreateOptions{
					WorkspaceRoot: t.TempDir(),
					Kind:          "task",
					Title:         "Implement coverage",
					ActorRole:     "developer",
				})
				return err
			},
			wantErr: "create issue: migrations dir is required",
		},
		{
			name: "create requires title",
			call: func() error {
				workspace := newIssueWorkspace(t)
				_, err := Create(CreateOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Kind:          "task",
					Title:         "  ",
					ActorRole:     "developer",
				})
				return err
			},
			wantErr: "create issue: title is required",
		},
		{
			name: "get requires workspace root",
			call: func() error {
				_, err := Get(GetOptions{
					MigrationsDir: testIssueMigrationsDir(t),
					ID:            "FID-123456",
				})
				return err
			},
			wantErr: "get issue: workspace root is required",
		},
		{
			name: "get requires migrations dir",
			call: func() error {
				_, err := Get(GetOptions{
					WorkspaceRoot: t.TempDir(),
					ID:            "FID-123456",
				})
				return err
			},
			wantErr: "get issue: migrations dir is required",
		},
		{
			name: "get requires issue id",
			call: func() error {
				workspace := newIssueWorkspace(t)
				_, err := Get(GetOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					ID:            "   ",
				})
				return err
			},
			wantErr: "get issue: issue id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatalf("error = nil, want %q", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCreateRejectsUnsupportedKind(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	_, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "bug",
		Title:         "Unsupported kind",
		ActorRole:     "developer",
	})
	if err == nil {
		t.Fatal("Create() error = nil, want unsupported kind")
	}

	var unsupportedErr *UnsupportedKindError
	if !errors.As(err, &unsupportedErr) {
		t.Fatalf("Create() error = %T, want *UnsupportedKindError", err)
	}
	if unsupportedErr.Kind != "bug" {
		t.Fatalf("UnsupportedKindError.Kind = %q, want %q", unsupportedErr.Kind, "bug")
	}
	if unsupportedErr.Error() != `kind "bug" must be epic, task, or subtask` {
		t.Fatalf("UnsupportedKindError.Error() = %q", unsupportedErr.Error())
	}
}

func TestCreateSupportsEpicTaskAndSubtaskAuthoring(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	seedIssueLabels(t, workspace, "core", "phase-4")

	epicCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "epic",
		Title:         "Phase 4 authoring",
		ActorRole:     "tl",
		Scheduling:    "sequential",
	})
	if err != nil {
		t.Fatalf("Create(epic) error = %v", err)
	}

	epic, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            epicCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get(epic) error = %v", err)
	}
	if epic.Kind != "epic" {
		t.Fatalf("Get(epic) kind = %q, want epic", epic.Kind)
	}
	if epic.WorkflowID != "epic-default" {
		t.Fatalf("Get(epic) workflow_id = %q, want epic-default", epic.WorkflowID)
	}
	if epic.Status != "todo" {
		t.Fatalf("Get(epic) status = %q, want todo", epic.Status)
	}
	if epic.StatusCategory != "unstarted" {
		t.Fatalf("Get(epic) status_category = %q, want unstarted", epic.StatusCategory)
	}
	if epic.ParentID != nil {
		t.Fatalf("Get(epic) parent_id = %#v, want nil", epic.ParentID)
	}
	if epic.OwnerRole == nil || *epic.OwnerRole != "tl" {
		t.Fatalf("Get(epic) owner_role = %#v, want tl", epic.OwnerRole)
	}

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Implement authoring domain",
		ActorRole:     "tl",
		ParentID:      stringPointer(epicCreated.ID),
		Priority:      intPointer(1),
		Sequence:      intPointer(1),
		SpecMDPath:    stringPointer("project-docs/docs/phases/phase-4/tasks/task-2.md"),
		EstTokens:     intPointer(2100),
		DueDate:       stringPointer("2026-06-30"),
		Labels:        []string{"phase-4", "core"},
		InInbox:       true,
		OwnerRole:     stringPointer("qa_lead"),
		OwnerActor:    stringPointer("human:mark"),
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	task, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get(task) error = %v", err)
	}
	if task.Kind != "task" {
		t.Fatalf("Get(task) kind = %q, want task", task.Kind)
	}
	if task.WorkflowID != "task-default" {
		t.Fatalf("Get(task) workflow_id = %q, want task-default", task.WorkflowID)
	}
	if task.ParentID == nil || *task.ParentID != epicCreated.ID {
		t.Fatalf("Get(task) parent_id = %#v, want %q", task.ParentID, epicCreated.ID)
	}
	if task.Priority == nil || *task.Priority != 1 {
		t.Fatalf("Get(task) priority = %#v, want 1", task.Priority)
	}
	if task.Sequence == nil || *task.Sequence != 1 {
		t.Fatalf("Get(task) sequence = %#v, want 1", task.Sequence)
	}
	if task.SpecMDPath == nil || *task.SpecMDPath != "project-docs/docs/phases/phase-4/tasks/task-2.md" {
		t.Fatalf("Get(task) spec_md_path = %#v", task.SpecMDPath)
	}
	if task.EstTokens == nil || *task.EstTokens != 2100 {
		t.Fatalf("Get(task) est_tokens = %#v, want 2100", task.EstTokens)
	}
	if task.OwnerRole == nil || *task.OwnerRole != "qa_lead" {
		t.Fatalf("Get(task) owner_role = %#v, want qa_lead", task.OwnerRole)
	}
	if task.OwnerActor == nil || *task.OwnerActor != "human:mark" {
		t.Fatalf("Get(task) owner_actor = %#v, want human:mark", task.OwnerActor)
	}
	if !task.InInbox {
		t.Fatal("Get(task) in_inbox = false, want true")
	}
	if want := []string{"core", "phase-4"}; !equalStrings(task.Labels, want) {
		t.Fatalf("Get(task) labels = %#v, want %#v", task.Labels, want)
	}

	subtaskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "subtask",
		Title:         "Cover rebuild parity",
		ActorRole:     "developer",
		ParentID:      stringPointer(taskCreated.ID),
	})
	if err != nil {
		t.Fatalf("Create(subtask) error = %v", err)
	}

	subtask, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            subtaskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get(subtask) error = %v", err)
	}
	if subtask.Kind != "subtask" {
		t.Fatalf("Get(subtask) kind = %q, want subtask", subtask.Kind)
	}
	if subtask.WorkflowID != "subtask-default" {
		t.Fatalf("Get(subtask) workflow_id = %q, want subtask-default", subtask.WorkflowID)
	}
	if subtask.ParentID == nil || *subtask.ParentID != taskCreated.ID {
		t.Fatalf("Get(subtask) parent_id = %#v, want %q", subtask.ParentID, taskCreated.ID)
	}
	if subtask.OwnerRole == nil || *subtask.OwnerRole != "developer" {
		t.Fatalf("Get(subtask) owner_role = %#v, want developer", subtask.OwnerRole)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var scheduling string
	if err := conn.QueryRow(`SELECT scheduling FROM epic_metadata WHERE issue_id = ?`, epicCreated.ID).Scan(&scheduling); err != nil {
		t.Fatalf("query epic_metadata: %v", err)
	}
	if scheduling != "sequential" {
		t.Fatalf("epic_metadata scheduling = %q, want sequential", scheduling)
	}

	var dueDate string
	if err := conn.QueryRow(`SELECT date(due_date) FROM issues WHERE id = ?`, taskCreated.ID).Scan(&dueDate); err != nil {
		t.Fatalf("query task due_date: %v", err)
	}
	if dueDate != "2026-06-30" {
		t.Fatalf("task due_date = %q, want 2026-06-30", dueDate)
	}

	stored := readIssueStoredEvents(t, journalPathForWorkspace(t, workspace))
	if len(stored) != 5 {
		t.Fatalf("stored events = %d, want 5", len(stored))
	}
	for _, event := range stored[2:] {
		if event.EventType != "issue.created" {
			t.Fatalf("stored event type = %q, want issue.created", event.EventType)
		}
	}
}

func TestCreateRejectsInvalidHierarchyAndUnknownLabels(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	seedIssueLabels(t, workspace, "known")

	epicCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "epic",
		Title:         "Parent epic",
		ActorRole:     "tl",
		Scheduling:    "prioritised",
	})
	if err != nil {
		t.Fatalf("Create(parent epic) error = %v", err)
	}

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Parent task",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(parent task) error = %v", err)
	}

	subtaskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "subtask",
		Title:         "Existing subtask",
		ActorRole:     "developer",
		ParentID:      stringPointer(taskCreated.ID),
	})
	if err != nil {
		t.Fatalf("Create(parent subtask) error = %v", err)
	}

	tests := []struct {
		name    string
		options CreateOptions
		wantErr string
	}{
		{
			name: "epic cannot have parent",
			options: CreateOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				Kind:          "epic",
				Title:         "Invalid epic",
				ActorRole:     "tl",
				ParentID:      stringPointer(epicCreated.ID),
				Scheduling:    "sequential",
			},
			wantErr: "issues: epic parent_id must be NULL",
		},
		{
			name: "subtask requires task parent",
			options: CreateOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				Kind:          "subtask",
				Title:         "Missing subtask parent",
				ActorRole:     "developer",
			},
			wantErr: "issues: subtask parent_id must reference a task",
		},
		{
			name: "task parent must be epic",
			options: CreateOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				Kind:          "task",
				Title:         "Task under subtask",
				ActorRole:     "developer",
				ParentID:      stringPointer(subtaskCreated.ID),
			},
			wantErr: "issues: task parent_id must reference an epic or be NULL",
		},
		{
			name: "subtask parent cannot be epic",
			options: CreateOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				Kind:          "subtask",
				Title:         "Subtask under epic",
				ActorRole:     "developer",
				ParentID:      stringPointer(epicCreated.ID),
			},
			wantErr: "issues: subtask parent_id must reference a task",
		},
		{
			name: "non-epic metadata rejected",
			options: CreateOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				Kind:          "task",
				Title:         "Task with scheduling",
				ActorRole:     "developer",
				Scheduling:    "sequential",
			},
			wantErr: "create issue: scheduling is only valid for epics",
		},
		{
			name: "unknown labels rejected",
			options: CreateOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				Kind:          "task",
				Title:         "Unknown label",
				ActorRole:     "developer",
				Labels:        []string{"known", "missing"},
			},
			wantErr: `create issue: label "missing" does not exist`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Create(tt.options)
			if err == nil {
				t.Fatalf("Create() error = nil, want %q", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("Create() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCreateRejectsProjectionFailingWriteBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	epicCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "epic",
		Title:         "Prioritised parent",
		ActorRole:     "tl",
		Scheduling:    "prioritised",
	})
	if err != nil {
		t.Fatalf("Create(epic) error = %v", err)
	}

	validTask, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Valid sibling",
		ActorRole:     "developer",
		ParentID:      stringPointer(epicCreated.ID),
	})
	if err != nil {
		t.Fatalf("Create(valid task) error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	_, err = Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Invalid sequenced task",
		ActorRole:     "developer",
		ParentID:      stringPointer(epicCreated.ID),
		Sequence:      intPointer(1),
	})
	if err == nil {
		t.Fatal("Create(invalid task) error = nil, want projection failure")
	}
	if !strings.Contains(err.Error(), "issues: task under prioritised epic must not have sequence") {
		t.Fatalf("Create(invalid task) error = %q, want prioritised-epic sequence rejection", err.Error())
	}

	if got := countJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected create = %d, want unchanged %d", got, linesBefore)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            validTask.ID,
	})
	if err != nil {
		t.Fatalf("Get(valid task after rejected create) error = %v", err)
	}
	if got.ID != validTask.ID {
		t.Fatalf("Get(valid task after rejected create) id = %q, want %q", got.ID, validTask.ID)
	}
}

func TestCreatePreservesAuthoringFieldsAcrossRebuild(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	seedIssueLabels(t, workspace, "phase-4", "ready")

	epicCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "epic",
		Title:         "Rebuild parent",
		ActorRole:     "tl",
		Scheduling:    "sequential",
	})
	if err != nil {
		t.Fatalf("Create(epic) error = %v", err)
	}

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Rebuild task",
		ActorRole:     "developer",
		ParentID:      stringPointer(epicCreated.ID),
		Priority:      intPointer(2),
		Sequence:      intPointer(1),
		SpecMDPath:    stringPointer("project-docs/specs/02_TRD.md"),
		EstTokens:     intPointer(1440),
		DueDate:       stringPointer("2026-07-01"),
		Labels:        []string{"ready", "phase-4"},
		InInbox:       true,
		OwnerActor:    stringPointer("agent"),
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	before, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get(before rebuild) error = %v", err)
	}

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	if _, err := journal.RebuildWorkspace(workspace, testIssueMigrationsDir(t), "developer"); err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}

	after, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get(after rebuild) error = %v", err)
	}

	assertIssuesEqual(t, after, before)

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var dueDate string
	if err := conn.QueryRow(`SELECT date(due_date) FROM issues WHERE id = ?`, taskCreated.ID).Scan(&dueDate); err != nil {
		t.Fatalf("query rebuilt due_date: %v", err)
	}
	if dueDate != "2026-07-01" {
		t.Fatalf("rebuilt due_date = %q, want 2026-07-01", dueDate)
	}

	var issueID, scheduling string
	if err := conn.QueryRow(`SELECT issue_id, scheduling FROM epic_metadata WHERE issue_id = ?`, epicCreated.ID).Scan(&issueID, &scheduling); err != nil {
		t.Fatalf("query rebuilt epic_metadata: %v", err)
	}
	if issueID != epicCreated.ID || scheduling != "sequential" {
		t.Fatalf("rebuilt epic_metadata = (%q, %q), want (%q, %q)", issueID, scheduling, epicCreated.ID, "sequential")
	}
}

func TestCreateAndGetRoundTrip(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Implement coverage remediation",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID == "" {
		t.Fatal("Create() ID = empty, want generated id")
	}
	if created.Kind != "task" {
		t.Fatalf("Create() kind = %q, want %q", created.Kind, "task")
	}
	if created.Title != "Implement coverage remediation" {
		t.Fatalf("Create() title = %q, want %q", created.Title, "Implement coverage remediation")
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("Get() id = %q, want %q", got.ID, created.ID)
	}
	if got.Status != "todo" {
		t.Fatalf("Get() status = %q, want %q", got.Status, "todo")
	}
	if got.StatusCategory != "unstarted" {
		t.Fatalf("Get() status_category = %q, want %q", got.StatusCategory, "unstarted")
	}
	if got.OwnerRole == nil || *got.OwnerRole != "developer" {
		t.Fatalf("Get() owner_role = %#v, want developer", got.OwnerRole)
	}
	if got.TriageStatus != nil {
		t.Fatalf("Get() triage_status = %#v, want nil", got.TriageStatus)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	exists, err := issueIDExists(conn, created.ID)
	if err != nil {
		t.Fatalf("issueIDExists(existing) error = %v", err)
	}
	if !exists {
		t.Fatal("issueIDExists(existing) = false, want true")
	}

	exists, err = issueIDExists(conn, "FID-does-not-exist")
	if err != nil {
		t.Fatalf("issueIDExists(missing) error = %v", err)
	}
	if exists {
		t.Fatal("issueIDExists(missing) = true, want false")
	}

	journalPath := journalPathForWorkspace(t, workspace)
	if got := countJournalLines(journalPath); got != 2 {
		t.Fatalf("countJournalLines(%q) = %d, want 2", journalPath, got)
	}
}

func TestGetReturnsEnrichedConcernFields(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Surface open concerns",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	event := journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "issue.concerned",
		Actor: journal.Actor{
			Role: "developer",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []journal.Effect{
			{
				Op:    "insert",
				Table: "task_completions",
				PK: map[string]any{
					"id": "TC-000201",
				},
				Values: map[string]any{
					"issue_id":           created.ID,
					"attempt_number":     1,
					"developer_role":     "developer",
					"model":              "gpt-5",
					"status":             "dev_complete",
					"gate_results":       `["go test ./..."]`,
					"human_intervention": false,
					"completed_at":       "2026-06-07T12:00:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-000201",
				},
				Values: map[string]any{
					"issue_id":                created.ID,
					"raised_in_completion_id": "TC-000201",
					"raised_by_role":          "developer",
					"summary":                 "Need triage before closing",
					"body_md_path":            "project-docs/concerns/CON-000201.md",
					"status":                  "open",
					"created_at":              "2026-06-07T12:01:00Z",
					"updated_at":              "2026-06-07T12:01:00Z",
				},
			},
		},
	}

	if err := journal.WithLock(workspace, journal.LockOwner{
		Command:   "fid maintenance seed",
		ActorRole: "developer",
	}, func(lock *journal.WriteLock) error {
		_, err := journal.AppendEvent(lock, event)
		return err
	}); err != nil {
		t.Fatalf("AppendEvent() error = %v", err)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if !got.ConcernsRaised {
		t.Fatal("Get() concerns_raised = false, want true")
	}
	if got.OpenConcernsCount != 1 {
		t.Fatalf("Get() open_concerns_count = %d, want 1", got.OpenConcernsCount)
	}
	if got.TriageStatus != nil {
		t.Fatalf("Get() triage_status = %#v, want nil", got.TriageStatus)
	}
}

func TestGetTreatsAddressedConcernsAsNotOpen(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Ignore addressed concerns",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	appendIssueJournalEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "issue.concern_triaged",
		Actor: journal.Actor{
			Role: "developer",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []journal.Effect{
			{
				Op:    "insert",
				Table: "task_completions",
				PK: map[string]any{
					"id": "TC-000301",
				},
				Values: map[string]any{
					"issue_id":           created.ID,
					"attempt_number":     1,
					"developer_role":     "developer",
					"model":              "gpt-5",
					"status":             "dev_complete",
					"gate_results":       `["go test ./..."]`,
					"human_intervention": false,
					"completed_at":       "2026-06-07T12:10:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "task_completions",
				PK: map[string]any{
					"id": "TC-000302",
				},
				Values: map[string]any{
					"issue_id":           created.ID,
					"attempt_number":     2,
					"developer_role":     "developer",
					"model":              "gpt-5",
					"status":             "qa_pass",
					"gate_results":       `["go test ./..."]`,
					"human_intervention": false,
					"completed_at":       "2026-06-07T12:20:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-000301",
				},
				Values: map[string]any{
					"issue_id":                   created.ID,
					"raised_in_completion_id":    "TC-000301",
					"raised_by_role":             "developer",
					"summary":                    "Resolved concern",
					"body_md_path":               "project-docs/concerns/CON-000301.md",
					"status":                     "addressed",
					"addressed_in_completion_id": "TC-000302",
					"created_at":                 "2026-06-07T12:11:00Z",
					"updated_at":                 "2026-06-07T12:21:00Z",
				},
			},
		},
	})

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.ConcernsRaised {
		t.Fatal("Get() concerns_raised = true, want false")
	}
	if got.OpenConcernsCount != 0 {
		t.Fatalf("Get() open_concerns_count = %d, want 0", got.OpenConcernsCount)
	}
}

func TestGetReturnsLatestTriageStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		eventIDs   []string
		verdicts   []string
		createdAts []string
		want       string
	}{
		{
			name:       "latest resume wins",
			eventIDs:   []string{"TRG-000401", "TRG-000402"},
			verdicts:   []string{"BLOCK", "RESUME"},
			createdAts: []string{"2026-06-07T12:30:00Z", "2026-06-07T12:40:00Z"},
			want:       "resumed",
		},
		{
			name:       "latest block wins",
			eventIDs:   []string{"TRG-000403", "TRG-000404"},
			verdicts:   []string{"RESUME", "BLOCK"},
			createdAts: []string{"2026-06-07T12:30:00Z", "2026-06-07T12:40:00Z"},
			want:       "blocked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := newIssueWorkspace(t)

			created, err := Create(CreateOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				Kind:          "task",
				Title:         "Map latest triage verdict",
				ActorRole:     "developer",
			})
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}

			appendIssueJournalEvent(t, workspace, journal.Event{
				EventSchema: journal.EventSchemaV1,
				EventType:   "issue.triaged",
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
							"id": "FID-epictriage",
						},
						Values: map[string]any{
							"kind":        "epic",
							"title":       "Phase 2 triage epic",
							"workflow_id": "epic-default",
							"status":      "todo",
							"created_at":  "2026-06-07T12:25:00Z",
							"updated_at":  "2026-06-07T12:25:00Z",
						},
					},
					{
						Op:    "insert",
						Table: "triage_events",
						PK: map[string]any{
							"id": tt.eventIDs[0],
						},
						Values: map[string]any{
							"epic_id":       "FID-epictriage",
							"issue_id":      created.ID,
							"triager_role":  "qa_lead",
							"source":        "task_completion",
							"source_ref_id": "TC-triage-1",
							"verdict":       tt.verdicts[0],
							"reasoning_md":  "project-docs/triage/TRG-1.md",
							"created_at":    tt.createdAts[0],
						},
					},
					{
						Op:    "insert",
						Table: "triage_events",
						PK: map[string]any{
							"id": tt.eventIDs[1],
						},
						Values: map[string]any{
							"epic_id":       "FID-epictriage",
							"issue_id":      created.ID,
							"triager_role":  "qa_lead",
							"source":        "task_completion",
							"source_ref_id": "TC-triage-2",
							"verdict":       tt.verdicts[1],
							"reasoning_md":  "project-docs/triage/TRG-2.md",
							"created_at":    tt.createdAts[1],
						},
					},
				},
			})

			got, err := Get(GetOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				ID:            created.ID,
			})
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}

			if got.TriageStatus == nil || *got.TriageStatus != tt.want {
				t.Fatalf("Get() triage_status = %#v, want %q", got.TriageStatus, tt.want)
			}
		})
	}
}

func TestGetMapsPopulatedNullableFields(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	_, err := conn.Exec(`
		INSERT INTO issues (id, kind, title, status, workflow_id, owner_role)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		"FID-parent",
		"epic",
		"Parent issue",
		"todo",
		"epic-default",
		"developer",
	)
	if err != nil {
		t.Fatalf("insert parent issue: %v", err)
	}

	_, err = conn.Exec(`
		INSERT INTO epic_metadata (issue_id, scheduling)
		VALUES (?, ?)
	`,
		"FID-parent",
		"sequential",
	)
	if err != nil {
		t.Fatalf("insert parent epic metadata: %v", err)
	}

	_, err = conn.Exec(`
		INSERT INTO issues (
			id, kind, title, status, workflow_id, parent_id, priority, sequence,
			spec_md_path, est_tokens, owner_role, owner_actor, in_inbox, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		"FID-populated",
		"task",
		"Populated values",
		"qa_pass",
		"task-default",
		"FID-parent",
		5,
		42,
		"docs/spec.md",
		99,
		"qa_lead",
		"human",
		true,
		"2026-06-07T10:00:00Z",
		"2026-06-07T11:00:00Z",
	)
	if err != nil {
		t.Fatalf("insert populated issue: %v", err)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            "FID-populated",
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.StatusCategory != "completed" {
		t.Fatalf("Get() status_category = %q, want %q", got.StatusCategory, "completed")
	}
	if got.ParentID == nil || *got.ParentID != "FID-parent" {
		t.Fatalf("Get() parent_id = %#v, want FID-parent", got.ParentID)
	}
	if got.Priority == nil || *got.Priority != 5 {
		t.Fatalf("Get() priority = %#v, want 5", got.Priority)
	}
	if got.Sequence == nil || *got.Sequence != 42 {
		t.Fatalf("Get() sequence = %#v, want 42", got.Sequence)
	}
	if got.SpecMDPath == nil || *got.SpecMDPath != "docs/spec.md" {
		t.Fatalf("Get() spec_md_path = %#v, want docs/spec.md", got.SpecMDPath)
	}
	if got.EstTokens == nil || *got.EstTokens != 99 {
		t.Fatalf("Get() est_tokens = %#v, want 99", got.EstTokens)
	}
	if got.OwnerRole == nil || *got.OwnerRole != "qa_lead" {
		t.Fatalf("Get() owner_role = %#v, want qa_lead", got.OwnerRole)
	}
	if got.OwnerActor == nil || *got.OwnerActor != "human" {
		t.Fatalf("Get() owner_actor = %#v, want human", got.OwnerActor)
	}
	if got.CreatedAt == nil || *got.CreatedAt != "2026-06-07T10:00:00Z" {
		t.Fatalf("Get() created_at = %#v, want 2026-06-07T10:00:00Z", got.CreatedAt)
	}
	if got.UpdatedAt == nil || *got.UpdatedAt != "2026-06-07T11:00:00Z" {
		t.Fatalf("Get() updated_at = %#v, want 2026-06-07T11:00:00Z", got.UpdatedAt)
	}
}

func TestIssueHelpersCoverStatusAndNullHandling(t *testing.T) {
	t.Parallel()

	t.Run("status category mapping", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			status string
			want   string
		}{
			{status: "todo", want: "unstarted"},
			{status: "in_progress", want: "started"},
			{status: "dev_complete", want: "started"},
			{status: "qa_review", want: "started"},
			{status: "qa_block", want: "started"},
			{status: "failed_gate", want: "started"},
			{status: "wrapped", want: "started"},
			{status: "qa_pass", want: "completed"},
			{status: "done", want: "completed"},
			{status: "closed", want: "completed"},
			{status: "cancelled", want: "cancelled"},
			{status: "unknown", want: "unstarted"},
		}

		for _, tt := range tests {
			if got := statusCategory(tt.status); got != tt.want {
				t.Fatalf("statusCategory(%q) = %q, want %q", tt.status, got, tt.want)
			}
		}
	})

	t.Run("null pointer helpers", func(t *testing.T) {
		t.Parallel()

		if nullStringPointer(sql.NullString{}) != nil {
			t.Fatal("nullStringPointer(invalid) = non-nil, want nil")
		}
		if nullIntPointer(sql.NullInt64{}) != nil {
			t.Fatal("nullIntPointer(invalid) = non-nil, want nil")
		}

		stringValue := nullStringPointer(sql.NullString{String: "developer", Valid: true})
		if stringValue == nil || *stringValue != "developer" {
			t.Fatalf("nullStringPointer(valid) = %#v, want developer", stringValue)
		}

		intValue := nullIntPointer(sql.NullInt64{Int64: 13, Valid: true})
		if intValue == nil || *intValue != 13 {
			t.Fatalf("nullIntPointer(valid) = %#v, want 13", intValue)
		}
	})

	t.Run("count journal lines", func(t *testing.T) {
		t.Parallel()

		if got := countJournalLines(filepath.Join(t.TempDir(), "missing.jsonl")); got != 1 {
			t.Fatalf("countJournalLines(missing) = %d, want 1", got)
		}

		path := filepath.Join(t.TempDir(), "journal.jsonl")
		if err := os.WriteFile(path, []byte("line1\nline2\n"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if got := countJournalLines(path); got != 2 {
			t.Fatalf("countJournalLines(%q) = %d, want 2", path, got)
		}
	})
}

func testIssueMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func newIssueWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}

	return workspace
}

func openIssueDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", dbPath, err)
	}

	return conn
}

func closeIssueTestDB(t *testing.T, conn interface{ Close() error }) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func journalPathForWorkspace(t *testing.T, workspace string) string {
	t.Helper()

	journalDir := filepath.Join(workspace, ".fid", "journal")
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", journalDir, err)
	}
	if len(entries) != 1 {
		t.Fatalf("journal entries = %d, want 1", len(entries))
	}
	if entries[0].IsDir() {
		t.Fatalf("journal entry %q is a directory, want file", entries[0].Name())
	}

	return filepath.Join(journalDir, entries[0].Name())
}

func appendIssueJournalEvent(t *testing.T, workspace string, event journal.Event) {
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

func seedIssueLabels(t *testing.T, workspace string, labels ...string) {
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

	appendIssueJournalEvent(t, workspace, journal.Event{
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

func readIssueStoredEvents(t *testing.T, path string) []journal.Event {
	t.Helper()

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}

	lines := strings.Split(strings.TrimSpace(string(payload)), "\n")
	events := make([]journal.Event, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var event journal.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("UnmarshalCanonicalJSON() error = %v", err)
		}
		events = append(events, event)
	}

	return events
}

func assertIssuesEqual(t *testing.T, got, want Issue) {
	t.Helper()

	if got.ID != want.ID ||
		got.Kind != want.Kind ||
		got.Title != want.Title ||
		got.Status != want.Status ||
		got.StatusCategory != want.StatusCategory ||
		got.WorkflowID != want.WorkflowID ||
		got.InInbox != want.InInbox ||
		got.ConcernsRaised != want.ConcernsRaised ||
		got.OpenConcernsCount != want.OpenConcernsCount {
		t.Fatalf("issue mismatch = %#v, want %#v", got, want)
	}
	if !equalStringPointers(got.ParentID, want.ParentID) ||
		!equalIntPointers(got.Priority, want.Priority) ||
		!equalIntPointers(got.Sequence, want.Sequence) ||
		!equalStringPointers(got.SpecMDPath, want.SpecMDPath) ||
		!equalIntPointers(got.EstTokens, want.EstTokens) ||
		!equalStringPointers(got.OwnerRole, want.OwnerRole) ||
		!equalStringPointers(got.OwnerActor, want.OwnerActor) ||
		!equalStringPointers(got.TriageStatus, want.TriageStatus) ||
		!equalStringPointers(got.CreatedAt, want.CreatedAt) ||
		!equalStringPointers(got.UpdatedAt, want.UpdatedAt) ||
		!equalStrings(got.Labels, want.Labels) {
		t.Fatalf("issue pointer/slice mismatch = %#v, want %#v", got, want)
	}
}

func equalStringPointers(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func equalIntPointers(left, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func stringPointer(value string) *string {
	return &value
}

func intPointer(value int) *int {
	return &value
}
