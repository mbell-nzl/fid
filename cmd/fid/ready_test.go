package main

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestReadyHumanSuccessUsesStdout(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)
	seedReadyCommandRows(t, conn, []readyCommandRow{{id: "FID-human-ready", kind: "task", title: "Human ready", status: "todo", workflowID: "task-default", ownerRole: "developer", createdAt: "2026-06-14T10:00:00Z", updatedAt: "2026-06-14T10:00:00Z"}})
	result := executeCommandForTest(workspace, "ready")
	if result.err != nil || result.stderr != "" || !strings.Contains(result.stdout, "FID-human-ready") || strings.Contains(result.stdout, ".fid/") {
		t.Fatalf("human ready = %#v", result)
	}
}

func TestReadyEpicScopesSelectionAndBlockerStates(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	seedReadyCommandRows(t, conn, []readyCommandRow{
		{
			id:         "FID-epic-seq-ready",
			kind:       "epic",
			title:      "Sequential ready epic",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T11:00:00Z",
			updatedAt:  "2026-06-14T11:00:00Z",
			scheduling: stringPointer("sequential"),
		},
		{
			id:         "FID-seq-ready-current",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq-ready"),
			title:      "Current sequential task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(1),
			priority:   intPointer(2),
			createdAt:  "2026-06-14T11:01:00Z",
			updatedAt:  "2026-06-14T11:01:00Z",
		},
		{
			id:         "FID-seq-ready-later",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq-ready"),
			title:      "Later sequential task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(2),
			priority:   intPointer(0),
			createdAt:  "2026-06-14T11:02:00Z",
			updatedAt:  "2026-06-14T11:02:00Z",
		},
		{
			id:         "FID-epic-seq-blocked",
			kind:       "epic",
			title:      "Sequential blocked epic",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T11:03:00Z",
			updatedAt:  "2026-06-14T11:03:00Z",
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
			createdAt:  "2026-06-14T11:04:00Z",
			updatedAt:  "2026-06-14T11:04:00Z",
		},
		{
			id:         "FID-seq-blocked-later",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq-blocked"),
			title:      "Blocked sequential later task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(2),
			createdAt:  "2026-06-14T11:05:00Z",
			updatedAt:  "2026-06-14T11:05:00Z",
		},
		{
			id:         "FID-epic-pri",
			kind:       "epic",
			title:      "Prioritised epic",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T11:06:00Z",
			updatedAt:  "2026-06-14T11:06:00Z",
			scheduling: stringPointer("prioritised"),
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
			createdAt:  "2026-06-14T11:07:00Z",
			updatedAt:  "2026-06-14T11:07:00Z",
		},
		{
			id:         "FID-pri-ready",
			kind:       "task",
			parentID:   stringPointer("FID-epic-pri"),
			title:      "Ready prioritised winner",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(1),
			createdAt:  "2026-06-14T11:08:00Z",
			updatedAt:  "2026-06-14T11:08:00Z",
		},
		{
			id:         "FID-pri-lower",
			kind:       "task",
			parentID:   stringPointer("FID-epic-pri"),
			title:      "Lower priority task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(2),
			createdAt:  "2026-06-14T11:09:00Z",
			updatedAt:  "2026-06-14T11:09:00Z",
		},
	})

	tests := []struct {
		name          string
		epicID        string
		wantTaskID    string
		wantNull      bool
		wantNoteParts []string
	}{
		{
			name:       "sequential actionable current task",
			epicID:     "FID-epic-seq-ready",
			wantTaskID: "FID-seq-ready-current",
		},
		{
			name:          "sequential started current task returns empty note",
			epicID:        "FID-epic-seq-blocked",
			wantNull:      true,
			wantNoteParts: []string{"note:", "FID-epic-seq-blocked", "FID-seq-blocked-current", "in_progress"},
		},
		{
			name:       "prioritised epic uses highest priority actionable task",
			epicID:     "FID-epic-pri",
			wantTaskID: "FID-pri-ready",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "--json", "ready", "--epic", tt.epicID)
			got := decodeCanonicalResult(t, stdout)

			task, hasTask := got["task"]
			if !hasTask {
				t.Fatalf("stdout = %#v, want task key", got)
			}

			if tt.wantNull {
				if task != nil {
					t.Fatalf("stdout task = %#v, want nil", task)
				}
				if stderr != "" {
					t.Fatalf("stderr = %q, want empty", stderr)
				}
				return
			}

			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			taskMap, ok := task.(map[string]any)
			if !ok {
				t.Fatalf("stdout task = %#v, want object", task)
			}
			if gotID := taskMap["id"]; gotID != tt.wantTaskID {
				t.Fatalf("stdout task id = %#v, want %q", gotID, tt.wantTaskID)
			}
		})
	}
}

func TestReadyProjectWideUsesMixedSourceSelectionOrder(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	seedReadyCommandRows(t, conn, []readyCommandRow{
		{
			id:         "FID-epic-seq",
			kind:       "epic",
			title:      "Sequential epic",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T12:00:00Z",
			updatedAt:  "2026-06-14T12:00:00Z",
			scheduling: stringPointer("sequential"),
		},
		{
			id:         "FID-seq-current",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq"),
			title:      "Sequential current task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(1),
			priority:   intPointer(2),
			createdAt:  "2026-06-14T12:01:00Z",
			updatedAt:  "2026-06-14T12:01:00Z",
		},
		{
			id:         "FID-seq-later",
			kind:       "task",
			parentID:   stringPointer("FID-epic-seq"),
			title:      "Sequential later task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			sequence:   intPointer(2),
			priority:   intPointer(0),
			createdAt:  "2026-06-14T12:02:00Z",
			updatedAt:  "2026-06-14T12:02:00Z",
		},
		{
			id:         "FID-epic-pri",
			kind:       "epic",
			title:      "Prioritised epic",
			status:     "todo",
			workflowID: "epic-default",
			ownerRole:  "tl",
			createdAt:  "2026-06-14T12:03:00Z",
			updatedAt:  "2026-06-14T12:03:00Z",
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
			createdAt:  "2026-06-14T12:04:00Z",
			updatedAt:  "2026-06-14T12:04:00Z",
		},
		{
			id:         "FID-unattached-top",
			kind:       "task",
			title:      "Unattached deterministic winner",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(1),
			createdAt:  "2026-06-14T12:00:30Z",
			updatedAt:  "2026-06-14T12:00:30Z",
		},
		{
			id:         "FID-unattached-no-priority",
			kind:       "task",
			title:      "Nil priority loses",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			createdAt:  "2026-06-14T11:59:00Z",
			updatedAt:  "2026-06-14T11:59:00Z",
		},
	})

	stdout, stderr := runFIDCommand(t, workspace, "--json", "ready")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	got := decodeCanonicalResult(t, stdout)
	task := got["task"].(map[string]any)
	if task["id"] != "FID-unattached-top" {
		t.Fatalf("stdout task id = %#v, want %q", task["id"], "FID-unattached-top")
	}
}

func TestReadyRejectsRetiredFormatFlag(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	if result := executeCommandForTest(workspace, "ready", "--format", "json"); result.err == nil {
		t.Fatal("--format was accepted")
	}
}

type readyCommandRow struct {
	id         string
	kind       string
	parentID   *string
	title      string
	status     string
	workflowID string
	ownerRole  string
	sequence   *int
	priority   *int
	createdAt  string
	updatedAt  string
	scheduling *string
}

func seedReadyCommandRows(t *testing.T, conn *sql.DB, rows []readyCommandRow) {
	t.Helper()

	for _, row := range rows {
		if _, err := conn.Exec(`
			INSERT INTO issues (
				id, kind, parent_id, title, status, workflow_id, owner_role, sequence, priority, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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

func stringPointer(value string) *string {
	result := value
	return &result
}

func intPointer(value int) *int {
	result := value
	return &result
}
