package main

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestReadyDefaultsToFrozenJSONShapeWithoutNote(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	seedReadyCommandRows(t, conn, []readyCommandRow{
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
			createdAt:  "2026-06-14T10:03:00Z",
			updatedAt:  "2026-06-14T10:03:00Z",
		},
	})

	defaultStdout, defaultStderr := runFIDCommand(t, workspace, "ready")
	if defaultStderr != "" {
		t.Fatalf("default stderr = %q, want empty", defaultStderr)
	}

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "ready", "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	if defaultStdout != jsonStdout {
		t.Fatalf("default stdout = %q, want json stdout %q", defaultStdout, jsonStdout)
	}

	got := decodeCommandJSONStructure(t, defaultStdout)
	want := map[string]any{
		"task": map[string]any{
			"id":         "FID-task-02",
			"title":      "Implement engine",
			"status":     "todo",
			"epic_id":    "FID-epic-seq",
			"epic_title": "Phase 6 ready queue",
			"sequence":   "2",
			"priority":   "1",
			"est_tokens": "2400",
		},
	}

	if !commandStructuresEqual(want, got) {
		t.Fatalf("stdout structure = %#v, want %#v", got, want)
	}
}

func TestReadyEpicScopesSelectionAndBlockerStates(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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
			stdout, stderr := runFIDCommand(t, workspace, "ready", "--epic", tt.epicID, "--format", "json")
			got := decodeCommandJSONStructure(t, stdout).(map[string]any)

			task, hasTask := got["task"]
			if !hasTask {
				t.Fatalf("stdout = %#v, want task key", got)
			}

			if tt.wantNull {
				if task != nil {
					t.Fatalf("stdout task = %#v, want nil", task)
				}
				for _, part := range tt.wantNoteParts {
					if !strings.Contains(stderr, part) {
						t.Fatalf("stderr = %q, want substring %q", stderr, part)
					}
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
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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

	stdout, stderr := runFIDCommand(t, workspace, "ready", "--format", "json")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	got := decodeCommandJSONStructure(t, stdout).(map[string]any)
	task := got["task"].(map[string]any)
	if task["id"] != "FID-unattached-top" {
		t.Fatalf("stdout task id = %#v, want %q", task["id"], "FID-unattached-top")
	}
}

func TestReadyFormatsRenderSelectedAndEmptyStates(t *testing.T) {
	selectedWorkspace := t.TempDir()
	runInitCommand(t, selectedWorkspace, "init")

	selectedConn, err := fiddb.Open(filepath.Join(selectedWorkspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open(selected) error = %v", err)
	}
	defer closeCommandDB(t, selectedConn)

	seedReadyCommandRows(t, selectedConn, []readyCommandRow{
		{
			id:         "FID-task-01",
			kind:       "task",
			title:      "Selected task",
			status:     "todo",
			workflowID: "task-default",
			ownerRole:  "developer",
			priority:   intPointer(1),
			estTokens:  intPointer(900),
			createdAt:  "2026-06-14T13:00:00Z",
			updatedAt:  "2026-06-14T13:00:00Z",
		},
	})

	emptyWorkspace := t.TempDir()
	runInitCommand(t, emptyWorkspace, "init")

	emptyConn, err := fiddb.Open(filepath.Join(emptyWorkspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open(empty) error = %v", err)
	}
	defer closeCommandDB(t, emptyConn)

	seedReadyCommandRows(t, emptyConn, []readyCommandRow{
		{
			id:         "FID-task-02",
			kind:       "task",
			title:      "Started task",
			status:     "in_progress",
			workflowID: "task-default",
			ownerRole:  "developer",
			createdAt:  "2026-06-14T13:01:00Z",
			updatedAt:  "2026-06-14T13:01:00Z",
		},
	})

	selectedJSON, selectedJSONStderr := runFIDCommand(t, selectedWorkspace, "ready", "--format", "json")
	if selectedJSONStderr != "" {
		t.Fatalf("selected json stderr = %q, want empty", selectedJSONStderr)
	}
	wantSelected := decodeCommandJSONStructure(t, selectedJSON)

	for _, formatName := range []string{"json", "toon", "yaml"} {
		t.Run("selected "+formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, selectedWorkspace, "ready", "--format", formatName)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			var got any
			switch formatName {
			case "json":
				got = decodeCommandJSONStructure(t, stdout)
			case "toon":
				got = decodeCommandTOONStructure(t, stdout)
			case "yaml":
				got = decodeCommandYAMLStructure(t, stdout)
			}

			if !commandStructuresEqual(wantSelected, got) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, got, wantSelected)
			}
		})
	}

	t.Run("selected markdown", func(t *testing.T) {
		stdout, stderr := runFIDCommand(t, selectedWorkspace, "ready", "--format", "markdown")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}

		want := "" +
			"# Ready Task\n\n" +
			"- ID: FID-task-01\n" +
			"- Title: Selected task\n" +
			"- Status: todo\n" +
			"- Epic: (none)\n" +
			"- Sequence: (none)\n" +
			"- Priority: P1\n" +
			"- Est. tokens: 900\n"
		if stdout != want {
			t.Fatalf("markdown stdout = %q, want %q", stdout, want)
		}
	})

	for _, formatName := range []string{"json", "toon", "yaml"} {
		t.Run("empty "+formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, emptyWorkspace, "ready", "--format", formatName)
			if !strings.Contains(stderr, "note:") {
				t.Fatalf("stderr = %q, want note", stderr)
			}

			var got any
			switch formatName {
			case "json":
				if strings.TrimSpace(stdout) != `{"task":null}` {
					t.Fatalf("json stdout = %q, want exact null payload", stdout)
				}
				got = decodeCommandJSONStructure(t, stdout)
			case "toon":
				got = decodeCommandTOONStructure(t, stdout)
			case "yaml":
				got = decodeCommandYAMLStructure(t, stdout)
			}

			want := map[string]any{"task": nil}
			if !commandStructuresEqual(want, got) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, got, want)
			}
		})
	}

	t.Run("empty markdown", func(t *testing.T) {
		stdout, stderr := runFIDCommand(t, emptyWorkspace, "ready", "--format", "markdown")
		if !strings.Contains(stderr, "note:") {
			t.Fatalf("stderr = %q, want note", stderr)
		}

		want := "# Ready Task\n\n(no ready task)\n"
		if stdout != want {
			t.Fatalf("markdown stdout = %q, want %q", stdout, want)
		}
	})
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
	estTokens  *int
	createdAt  string
	updatedAt  string
	scheduling *string
}

func seedReadyCommandRows(t *testing.T, conn *sql.DB, rows []readyCommandRow) {
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

func stringPointer(value string) *string {
	result := value
	return &result
}

func intPointer(value int) *int {
	result := value
	return &result
}
