package issues

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestReadySelectsLowestSequenceActionableTaskForSequentialEpic(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	seedReadyRows(t, conn, []readyRow{
		{
			id:         "FID-epic-seq",
			kind:       "epic",
			title:      "Phase 6 ready queue",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T10:00:00Z",
			updatedAt:  "2026-06-14T10:00:00Z",
			scheduling: stringPointer("sequential"),
		},
		{
			id:         "FID-task-01",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq"),
			title:      "Complete prep",
			status:     "done",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(1),
			priority:   intPointer(2),
			estTokens:  intPointer(1200),
			createdAt:  "2026-06-14T10:01:00Z",
			updatedAt:  "2026-06-14T10:01:00Z",
		},
		{
			id:         "FID-task-02",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq"),
			title:      "Implement engine",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(2),
			priority:   intPointer(1),
			estTokens:  intPointer(2400),
			createdAt:  "2026-06-14T10:02:00Z",
			updatedAt:  "2026-06-14T10:02:00Z",
		},
		{
			id:         "FID-task-03",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq"),
			title:      "Wire CLI",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(3),
			priority:   intPointer(0),
			estTokens:  intPointer(1800),
			createdAt:  "2026-06-14T10:03:00Z",
			updatedAt:  "2026-06-14T10:03:00Z",
		},
	})

	got, err := Ready(ReadyOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		EpicID:        "FID-epic-seq",
	})
	if err != nil {
		t.Fatalf("Ready() error = %v", err)
	}
	if got.Task == nil {
		t.Fatal("Ready().Task = nil, want selected task")
	}
	if got.Explanation != nil {
		t.Fatalf("Ready().Explanation = %#v, want nil", got.Explanation)
	}
	if got.Task.ID != "FID-task-02" {
		t.Fatalf("Ready().Task.ID = %q, want %q", got.Task.ID, "FID-task-02")
	}
	if got.Task.Title != "Implement engine" {
		t.Fatalf("Ready().Task.Title = %q, want %q", got.Task.Title, "Implement engine")
	}
	if got.Task.Status != "todo" {
		t.Fatalf("Ready().Task.Status = %q, want %q", got.Task.Status, "todo")
	}
	if got.Task.EpicID == nil || *got.Task.EpicID != "FID-epic-seq" {
		t.Fatalf("Ready().Task.EpicID = %#v, want FID-epic-seq", got.Task.EpicID)
	}
	if got.Task.EpicTitle == nil || *got.Task.EpicTitle != "Phase 6 ready queue" {
		t.Fatalf("Ready().Task.EpicTitle = %#v, want Phase 6 ready queue", got.Task.EpicTitle)
	}
	if got.Task.Sequence == nil || *got.Task.Sequence != 2 {
		t.Fatalf("Ready().Task.Sequence = %#v, want 2", got.Task.Sequence)
	}
	if got.Task.Priority == nil || *got.Task.Priority != 1 {
		t.Fatalf("Ready().Task.Priority = %#v, want 1", got.Task.Priority)
	}
	if got.Task.EstTokens == nil || *got.Task.EstTokens != 2400 {
		t.Fatalf("Ready().Task.EstTokens = %#v, want 2400", got.Task.EstTokens)
	}
}

func TestReadyReturnsStructuredBlockerWhenSequentialCurrentTaskIsStarted(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	seedReadyRows(t, conn, []readyRow{
		{
			id:         "FID-epic-seq",
			kind:       "epic",
			title:      "Phase 6 ready queue",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T11:00:00Z",
			updatedAt:  "2026-06-14T11:00:00Z",
			scheduling: stringPointer("sequential"),
		},
		{
			id:         "FID-task-01",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq"),
			title:      "Complete prep",
			status:     "done",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(1),
			createdAt:  "2026-06-14T11:01:00Z",
			updatedAt:  "2026-06-14T11:01:00Z",
		},
		{
			id:         "FID-task-02",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq"),
			title:      "Implement engine",
			status:     "in_progress",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(2),
			createdAt:  "2026-06-14T11:02:00Z",
			updatedAt:  "2026-06-14T11:02:00Z",
		},
		{
			id:         "FID-task-03",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq"),
			title:      "Wire CLI",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(3),
			createdAt:  "2026-06-14T11:03:00Z",
			updatedAt:  "2026-06-14T11:03:00Z",
		},
	})

	got, err := Ready(ReadyOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		EpicID:        "FID-epic-seq",
	})
	if err != nil {
		t.Fatalf("Ready() error = %v", err)
	}
	if got.Task != nil {
		t.Fatalf("Ready().Task = %#v, want nil", got.Task)
	}
	if got.Explanation == nil {
		t.Fatal("Ready().Explanation = nil, want blocker explanation")
	}
	if got.Explanation.Code != "sequential_current_not_actionable" {
		t.Fatalf("Ready().Explanation.Code = %q, want %q", got.Explanation.Code, "sequential_current_not_actionable")
	}
	if got.Explanation.Scope != "epic" {
		t.Fatalf("Ready().Explanation.Scope = %q, want epic", got.Explanation.Scope)
	}
	if got.Explanation.EpicID == nil || *got.Explanation.EpicID != "FID-epic-seq" {
		t.Fatalf("Ready().Explanation.EpicID = %#v, want FID-epic-seq", got.Explanation.EpicID)
	}
	if got.Explanation.BlockingTaskID == nil || *got.Explanation.BlockingTaskID != "FID-task-02" {
		t.Fatalf("Ready().Explanation.BlockingTaskID = %#v, want FID-task-02", got.Explanation.BlockingTaskID)
	}
	if got.Explanation.BlockingTaskTitle == nil || *got.Explanation.BlockingTaskTitle != "Implement engine" {
		t.Fatalf("Ready().Explanation.BlockingTaskTitle = %#v, want Implement engine", got.Explanation.BlockingTaskTitle)
	}
	if got.Explanation.BlockingTaskStatus == nil || *got.Explanation.BlockingTaskStatus != "in_progress" {
		t.Fatalf("Ready().Explanation.BlockingTaskStatus = %#v, want in_progress", got.Explanation.BlockingTaskStatus)
	}
}

func TestReadySelectsHighestPriorityActionableTaskForPrioritisedEpic(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	seedReadyRows(t, conn, []readyRow{
		{
			id:         "FID-epic-pri",
			kind:       "epic",
			title:      "Phase 6 polish",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T12:00:00Z",
			updatedAt:  "2026-06-14T12:00:00Z",
			scheduling: stringPointer("prioritised"),
		},
		{
			id:         "FID-task-01",
			kind:       "task",
			parentID:   stringPointer("FID-epic-pri"),
			title:      "Started higher priority",
			status:     "in_progress",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(0),
			createdAt:  "2026-06-14T12:01:00Z",
			updatedAt:  "2026-06-14T12:01:00Z",
		},
		{
			id:         "FID-task-02",
			kind:       "task",
			parentID:   stringPointer("FID-epic-pri"),
			title:      "Ready winner",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(1),
			estTokens:  intPointer(1300),
			createdAt:  "2026-06-14T12:02:00Z",
			updatedAt:  "2026-06-14T12:02:00Z",
		},
		{
			id:         "FID-task-03",
			kind:       "task",
			parentID:   stringPointer("FID-epic-pri"),
			title:      "Lower priority",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(2),
			createdAt:  "2026-06-14T12:03:00Z",
			updatedAt:  "2026-06-14T12:03:00Z",
		},
	})

	got, err := Ready(ReadyOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		EpicID:        "FID-epic-pri",
	})
	if err != nil {
		t.Fatalf("Ready() error = %v", err)
	}
	if got.Task == nil {
		t.Fatal("Ready().Task = nil, want selected task")
	}
	if got.Task.ID != "FID-task-02" {
		t.Fatalf("Ready().Task.ID = %q, want %q", got.Task.ID, "FID-task-02")
	}
	if got.Explanation != nil {
		t.Fatalf("Ready().Explanation = %#v, want nil", got.Explanation)
	}
}

func TestReadyProjectWideIncludesOnlyActionableUnattachedTasks(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	seedReadyRows(t, conn, []readyRow{
		{
			id:         "FID-task-01",
			kind:       "task",
			title:      "Eligible unattached task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(1),
			estTokens:  intPointer(800),
			createdAt:  "2026-06-14T13:00:00Z",
			updatedAt:  "2026-06-14T13:00:00Z",
		},
		{
			id:         "FID-task-02",
			kind:       "task",
			title:      "Started unattached task",
			status:     "in_progress",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(0),
			createdAt:  "2026-06-14T13:01:00Z",
			updatedAt:  "2026-06-14T13:01:00Z",
		},
		{
			id:         "FID-task-03",
			kind:       "task",
			title:      "Completed unattached task",
			status:     "done",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(0),
			createdAt:  "2026-06-14T13:02:00Z",
			updatedAt:  "2026-06-14T13:02:00Z",
		},
		{
			id:         "FID-task-04",
			kind:       "task",
			title:      "Cancelled unattached task",
			status:     "cancelled",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(0),
			createdAt:  "2026-06-14T13:03:00Z",
			updatedAt:  "2026-06-14T13:03:00Z",
		},
	})

	got, err := Ready(ReadyOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("Ready() error = %v", err)
	}
	if got.Task == nil {
		t.Fatal("Ready().Task = nil, want unattached task")
	}
	if got.Task.ID != "FID-task-01" {
		t.Fatalf("Ready().Task.ID = %q, want %q", got.Task.ID, "FID-task-01")
	}
	if got.Task.EpicID != nil {
		t.Fatalf("Ready().Task.EpicID = %#v, want nil", got.Task.EpicID)
	}
	if got.Task.EpicTitle != nil {
		t.Fatalf("Ready().Task.EpicTitle = %#v, want nil", got.Task.EpicTitle)
	}
}

func TestReadyProjectWideSelectsDeterministicTopTaskAcrossMixedSources(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	seedReadyRows(t, conn, []readyRow{
		{
			id:         "FID-epic-seq-ready",
			kind:       "epic",
			title:      "Sequential ready epic",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T14:00:00Z",
			updatedAt:  "2026-06-14T14:00:00Z",
			scheduling: stringPointer("sequential"),
		},
		{
			id:         "FID-seq-current",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq-ready"),
			title:      "Current sequential task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(1),
			priority:   intPointer(2),
			createdAt:  "2026-06-14T14:01:00Z",
			updatedAt:  "2026-06-14T14:01:00Z",
		},
		{
			id:         "FID-seq-later",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq-ready"),
			title:      "Later sequential task must not contribute",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(2),
			priority:   intPointer(0),
			createdAt:  "2026-06-14T14:02:00Z",
			updatedAt:  "2026-06-14T14:02:00Z",
		},
		{
			id:         "FID-epic-seq-blocked",
			kind:       "epic",
			title:      "Sequential blocked epic",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T14:03:00Z",
			updatedAt:  "2026-06-14T14:03:00Z",
			scheduling: stringPointer("sequential"),
		},
		{
			id:         "FID-seq-blocked-current",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq-blocked"),
			title:      "Blocked sequential current task",
			status:     "in_progress",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(1),
			priority:   intPointer(5),
			createdAt:  "2026-06-14T14:04:00Z",
			updatedAt:  "2026-06-14T14:04:00Z",
		},
		{
			id:         "FID-seq-blocked-later",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq-blocked"),
			title:      "Blocked epic later task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(2),
			priority:   intPointer(0),
			createdAt:  "2026-06-14T14:05:00Z",
			updatedAt:  "2026-06-14T14:05:00Z",
		},
		{
			id:         "FID-epic-pri",
			kind:       "epic",
			title:      "Prioritised epic",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T14:06:00Z",
			updatedAt:  "2026-06-14T14:06:00Z",
			scheduling: stringPointer("prioritised"),
		},
		{
			id:         "FID-pri-best",
			kind:       "task",
			parentID:   stringPointer("FID-epic-pri"),
			title:      "Prioritised winner candidate",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(1),
			createdAt:  "2026-06-14T14:07:00Z",
			updatedAt:  "2026-06-14T14:07:00Z",
		},
		{
			id:         "FID-pri-started",
			kind:       "task",
			parentID:   stringPointer("FID-epic-pri"),
			title:      "Started prioritised task",
			status:     "in_progress",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(0),
			createdAt:  "2026-06-14T14:08:00Z",
			updatedAt:  "2026-06-14T14:08:00Z",
		},
		{
			id:         "FID-unattached-top",
			kind:       "task",
			title:      "Unattached deterministic winner",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(1),
			createdAt:  "2026-06-14T14:00:30Z",
			updatedAt:  "2026-06-14T14:00:30Z",
		},
	})

	got, err := Ready(ReadyOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("Ready() error = %v", err)
	}
	if got.Task == nil {
		t.Fatal("Ready().Task = nil, want selected task")
	}
	if got.Task.ID != "FID-unattached-top" {
		t.Fatalf("Ready().Task.ID = %q, want %q", got.Task.ID, "FID-unattached-top")
	}
}

func TestReadyProjectWideUsesIDAsFinalTieBreak(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	seedReadyRows(t, conn, []readyRow{
		{
			id:         "FID-task-b",
			kind:       "task",
			title:      "Later ID",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(1),
			createdAt:  "2026-06-14T15:00:00Z",
			updatedAt:  "2026-06-14T15:00:00Z",
		},
		{
			id:         "FID-task-a",
			kind:       "task",
			title:      "Earlier ID",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(1),
			createdAt:  "2026-06-14T15:00:00Z",
			updatedAt:  "2026-06-14T15:00:00Z",
		},
	})

	got, err := Ready(ReadyOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("Ready() error = %v", err)
	}
	if got.Task == nil {
		t.Fatal("Ready().Task = nil, want selected task")
	}
	if got.Task.ID != "FID-task-a" {
		t.Fatalf("Ready().Task.ID = %q, want %q", got.Task.ID, "FID-task-a")
	}
}

func TestReadyProjectWideReturnsStructuredEmptyExplanationWhenNothingIsReady(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	seedReadyRows(t, conn, []readyRow{
		{
			id:         "FID-task-01",
			kind:       "task",
			title:      "Completed task",
			status:     "done",
			workflowID: "task-default",
			ownerRole:  "developer",
			createdAt:  "2026-06-14T16:00:00Z",
			updatedAt:  "2026-06-14T16:00:00Z",
		},
		{
			id:         "FID-task-02",
			kind:       "task",
			title:      "Cancelled task",
			status:     "cancelled",
			workflowID: "task-default",
			ownerRole:  "developer",
			createdAt:  "2026-06-14T16:01:00Z",
			updatedAt:  "2026-06-14T16:01:00Z",
		},
		{
			id:         "FID-task-03",
			kind:       "task",
			title:      "Started task",
			status:     "in_progress",
			workflowID: "task-default",
			ownerRole:  "developer",
			createdAt:  "2026-06-14T16:02:00Z",
			updatedAt:  "2026-06-14T16:02:00Z",
		},
	})

	got, err := Ready(ReadyOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("Ready() error = %v", err)
	}
	if got.Task != nil {
		t.Fatalf("Ready().Task = %#v, want nil", got.Task)
	}
	if got.Explanation == nil {
		t.Fatal("Ready().Explanation = nil, want empty explanation")
	}
	if got.Explanation.Code != "no_ready_tasks" {
		t.Fatalf("Ready().Explanation.Code = %q, want %q", got.Explanation.Code, "no_ready_tasks")
	}
	if got.Explanation.Scope != "project" {
		t.Fatalf("Ready().Explanation.Scope = %q, want project", got.Explanation.Scope)
	}
}

type readyRow struct {
	id         string
	kind       string
	parentID   *string
	title      string
	status     string
	workflowID string
	ownerRole  string
	sequence   *int
	priority   *int
	estTokens  *int
	createdAt  string
	updatedAt  string
	scheduling *string
}

func seedReadyRows(t *testing.T, conn *sql.DB, rows []readyRow) {
	t.Helper()

	for _, row := range rows {
		if _, err := conn.Exec(`
			INSERT INTO issues (
				id, kind, parent_id, title, status, workflow_id, owner_role, sequence, priority, est_tokens, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			row.id,
			row.kind,
			row.parentID,
			row.title,
			row.status,
			row.workflowID,
			row.ownerRole,
			row.sequence,
			row.priority,
			row.estTokens,
			row.createdAt,
			row.updatedAt,
		); err != nil {
			t.Fatalf("seed ready issue %q: %v", row.id, err)
		}
		if row.scheduling != nil {
			if _, err := conn.Exec(`INSERT INTO epic_metadata (issue_id, scheduling) VALUES (?, ?)`, row.id, *row.scheduling); err != nil {
				t.Fatalf("seed ready epic_metadata %q: %v", row.id, err)
			}
		}
	}
}
