package journal

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestProjectEventAppliesEffectsAndUpdatesProjectionMeta(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeProjectDB(t, conn)

	event := Event{
		EventID:     "EVT-20260605-000001",
		EventSchema: EventSchemaV1,
		EventType:   "project.initialized",
		OccurredAt:  time.Date(2026, time.June, 5, 14, 23, 11, 0, time.UTC),
		Actor: Actor{
			Role: "human",
		},
		Command: Command{
			Argv:       []string{"fid", "init"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "insert",
				Table: "roles",
				PK: map[string]any{
					"name": "developer",
				},
				Values: map[string]any{
					"description": "Developer",
					"active":      true,
					"is_human":    false,
					"created_at":  "2026-06-05T14:23:11Z",
				},
			},
		},
	}
	eventHash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}

	err = ProjectEvent(conn, ProjectionInput{
		Event:     event,
		EventHash: eventHash,
		File:      "2026-06-05.jsonl",
		Line:      1,
		Offset:    123,
	})
	if err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}

	var description string
	var active bool
	if err := conn.QueryRow(`
		SELECT description, active
		FROM roles
		WHERE name = ?
	`, "developer").Scan(&description, &active); err != nil {
		t.Fatalf("query projected role: %v", err)
	}

	if description != "Developer" {
		t.Fatalf("projected description = %q, want %q", description, "Developer")
	}
	if !active {
		t.Fatal("projected role active = false, want true")
	}

	meta, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta() error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta() found = false, want true")
	}
	if meta.JournalFile != "2026-06-05.jsonl" {
		t.Fatalf("projection meta journal_file = %q, want %q", meta.JournalFile, "2026-06-05.jsonl")
	}
	if meta.JournalOffset != 123 {
		t.Fatalf("projection meta journal_offset = %d, want %d", meta.JournalOffset, 123)
	}
	if meta.LastEventID != event.EventID {
		t.Fatalf("projection meta last_event_id = %q, want %q", meta.LastEventID, event.EventID)
	}
	if meta.LastEventHash != eventHash {
		t.Fatalf("projection meta last_event_hash = %q, want %q", meta.LastEventHash, eventHash)
	}
	if meta.ReplayStatus != fiddb.ReplayStatusClean {
		t.Fatalf("projection meta replay_status = %q, want %q", meta.ReplayStatus, fiddb.ReplayStatusClean)
	}
}

func TestProjectEventValidatesRequiredInputs(t *testing.T) {
	t.Parallel()

	event := testV1Event()

	err := ProjectEvent(nil, ProjectionInput{
		Event:     event,
		EventHash: "sha256:test",
	})
	if err == nil {
		t.Fatal("ProjectEvent(nil) error = nil, want nil connection rejection")
	}
	if err.Error() != "project event: connection is nil" {
		t.Fatalf("ProjectEvent(nil) error = %q, want %q", err.Error(), "project event: connection is nil")
	}

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}
	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeProjectDB(t, conn)

	err = ProjectEvent(conn, ProjectionInput{
		Event:     Event{EventSchema: EventSchemaV1},
		EventHash: "sha256:test",
	})
	if err == nil {
		t.Fatal("ProjectEvent(missing id) error = nil, want event_id validation")
	}
	if err.Error() != "project event: event_id is required" {
		t.Fatalf("ProjectEvent(missing id) error = %q, want %q", err.Error(), "project event: event_id is required")
	}

	err = ProjectEvent(conn, ProjectionInput{
		Event: event,
	})
	if err == nil {
		t.Fatal("ProjectEvent(missing hash) error = nil, want event hash validation")
	}
	if err.Error() != `project event "EVT-20260529-000001": event hash is required` {
		t.Fatalf("ProjectEvent(missing hash) error = %q", err.Error())
	}
}

func TestProjectEventAppliesUpdateDeleteAndReplayStatus(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeProjectDB(t, conn)

	insert := projectRoleEvent("EVT-20260605-000001", "insert", map[string]any{
		"description": "Developer",
		"active":      true,
	}, nil)
	insertHash, err := HashEvent(insert)
	if err != nil {
		t.Fatalf("HashEvent(insert) error = %v", err)
	}
	if err := ProjectEvent(conn, ProjectionInput{
		Event:      insert,
		EventHash:  insertHash,
		File:       "2026-06-05.jsonl",
		Line:       1,
		Offset:     100,
		ReplayMode: fiddb.ReplayStatusRebuilding,
	}); err != nil {
		t.Fatalf("ProjectEvent(insert) error = %v", err)
	}

	update := projectRoleEvent("EVT-20260605-000002", "update", map[string]any{
		"description": "Developer Updated",
		"active":      false,
	}, map[string]any{
		"active": true,
	})
	updateHash, err := HashEvent(update)
	if err != nil {
		t.Fatalf("HashEvent(update) error = %v", err)
	}
	if err := ProjectEvent(conn, ProjectionInput{
		Event:      update,
		EventHash:  updateHash,
		File:       "2026-06-05.jsonl",
		Line:       2,
		Offset:     200,
		ReplayMode: fiddb.ReplayStatusCatchingUp,
	}); err != nil {
		t.Fatalf("ProjectEvent(update) error = %v", err)
	}

	var description string
	var active bool
	if err := conn.QueryRow(`SELECT description, active FROM roles WHERE name = ?`, "developer").Scan(&description, &active); err != nil {
		t.Fatalf("query updated role: %v", err)
	}
	if description != "Developer Updated" {
		t.Fatalf("updated description = %q, want %q", description, "Developer Updated")
	}
	if active {
		t.Fatal("updated active = true, want false")
	}

	deleteEvent := projectRoleEvent("EVT-20260605-000003", "delete", nil, map[string]any{
		"active": false,
	})
	deleteHash, err := HashEvent(deleteEvent)
	if err != nil {
		t.Fatalf("HashEvent(delete) error = %v", err)
	}
	if err := ProjectEvent(conn, ProjectionInput{
		Event:     deleteEvent,
		EventHash: deleteHash,
		File:      "2026-06-05.jsonl",
		Line:      3,
		Offset:    300,
	}); err != nil {
		t.Fatalf("ProjectEvent(delete) error = %v", err)
	}

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM roles WHERE name = ?`, "developer").Scan(&count); err != nil {
		t.Fatalf("count deleted role: %v", err)
	}
	if count != 0 {
		t.Fatalf("deleted role count = %d, want %d", count, 0)
	}

	meta, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta() error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta() found = false, want true")
	}
	if meta.ReplayStatus != fiddb.ReplayStatusClean {
		t.Fatalf("final replay status = %q, want %q", meta.ReplayStatus, fiddb.ReplayStatusClean)
	}
}

func TestProjectEventHonoursSequenceSchedulingBackstop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		event   Event
		wantErr string
	}{
		{
			name: "accepts sequential epic task with sequence when metadata is projected first",
			event: Event{
				EventID:     "EVT-20260607-000710",
				EventSchema: EventSchemaV1,
				EventType:   "issue.created",
				OccurredAt:  time.Date(2026, time.June, 7, 15, 10, 0, 0, time.UTC),
				Actor: Actor{
					Role: "developer",
				},
				Command: Command{
					Argv:       []string{"fid", "maintenance", "replay"},
					CWD:        ".",
					FIDVersion: "0.1.0",
				},
				Effects: []Effect{
					projectIssueEffect("FID-project-epic", map[string]any{
						"kind":        "epic",
						"title":       "Projected epic",
						"workflow_id": "epic-default",
						"status":      "todo",
						"created_at":  "2026-06-07T15:10:00Z",
						"updated_at":  "2026-06-07T15:10:00Z",
					}),
					{
						Op:    "insert",
						Table: "epic_metadata",
						PK: map[string]any{
							"issue_id": "FID-project-epic",
						},
						Values: map[string]any{
							"scheduling": "sequential",
						},
					},
					projectIssueEffect("FID-project-task", map[string]any{
						"kind":        "task",
						"parent_id":   "FID-project-epic",
						"sequence":    1,
						"title":       "Projected task",
						"workflow_id": "task-default",
						"status":      "todo",
						"created_at":  "2026-06-07T15:11:00Z",
						"updated_at":  "2026-06-07T15:11:00Z",
					}),
				},
			},
		},
		{
			name: "rejects sequential epic task without sequence during projection",
			event: Event{
				EventID:     "EVT-20260607-000711",
				EventSchema: EventSchemaV1,
				EventType:   "issue.created",
				OccurredAt:  time.Date(2026, time.June, 7, 15, 11, 0, 0, time.UTC),
				Actor: Actor{
					Role: "developer",
				},
				Command: Command{
					Argv:       []string{"fid", "maintenance", "replay"},
					CWD:        ".",
					FIDVersion: "0.1.0",
				},
				Effects: []Effect{
					projectIssueEffect("FID-project-epic-invalid", map[string]any{
						"kind":        "epic",
						"title":       "Projected epic invalid",
						"workflow_id": "epic-default",
						"status":      "todo",
						"created_at":  "2026-06-07T15:12:00Z",
						"updated_at":  "2026-06-07T15:12:00Z",
					}),
					{
						Op:    "insert",
						Table: "epic_metadata",
						PK: map[string]any{
							"issue_id": "FID-project-epic-invalid",
						},
						Values: map[string]any{
							"scheduling": "sequential",
						},
					},
					projectIssueEffect("FID-project-task-invalid", map[string]any{
						"kind":        "task",
						"parent_id":   "FID-project-epic-invalid",
						"title":       "Projected task invalid",
						"workflow_id": "task-default",
						"status":      "todo",
						"created_at":  "2026-06-07T15:13:00Z",
						"updated_at":  "2026-06-07T15:13:00Z",
					}),
				},
			},
			wantErr: "issues: task under sequential epic requires sequence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
			if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
				t.Fatalf("MigrateUp() error = %v", err)
			}

			conn, err := fiddb.Open(dbPath)
			if err != nil {
				t.Fatalf("Open() error = %v", err)
			}
			defer closeProjectDB(t, conn)

			eventHash, err := HashEvent(tt.event)
			if err != nil {
				t.Fatalf("HashEvent() error = %v", err)
			}

			err = ProjectEvent(conn, ProjectionInput{
				Event:     tt.event,
				EventHash: eventHash,
				File:      "2026-06-07.jsonl",
				Line:      1,
				Offset:    2048,
			})
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("ProjectEvent() error = nil, want substring %q", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("ProjectEvent() error = %q, want substring %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ProjectEvent() error = %v", err)
			}

			var sequence int
			if err := conn.QueryRow(`SELECT sequence FROM issues WHERE id = ?`, "FID-project-task").Scan(&sequence); err != nil {
				t.Fatalf("query projected task sequence: %v", err)
			}
			if sequence != 1 {
				t.Fatalf("projected task sequence = %d, want %d", sequence, 1)
			}
		})
	}
}

func TestDryRunProjectEventRollsBackEffectsAndLeavesProjectionMetaUnchanged(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeProjectDB(t, conn)

	baseEvent := projectRoleEvent("EVT-20260613-000001", "insert", map[string]any{
		"description": "Developer",
		"active":      true,
	}, nil)
	baseHash, err := HashEvent(baseEvent)
	if err != nil {
		t.Fatalf("HashEvent(base) error = %v", err)
	}
	if err := ProjectEvent(conn, ProjectionInput{
		Event:     baseEvent,
		EventHash: baseHash,
		File:      "2026-06-13.jsonl",
		Line:      1,
		Offset:    100,
	}); err != nil {
		t.Fatalf("ProjectEvent(base) error = %v", err)
	}

	metaBefore, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta(before) error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta(before) found = false, want true")
	}

	err = DryRunProjectEvent(conn, Event{
		EventType: "issue.created",
		Actor: Actor{
			Role: "developer",
		},
		Command: Command{
			Argv:       []string{"fid", "issue", "create"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "insert",
				Table: "roles",
				PK: map[string]any{
					"name": "qa_lead",
				},
				Values: map[string]any{
					"description": "QA Lead",
					"active":      true,
					"is_human":    true,
					"created_at":  "2026-06-13T10:00:00Z",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("DryRunProjectEvent() error = %v", err)
	}

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM roles WHERE name = ?`, "qa_lead").Scan(&count); err != nil {
		t.Fatalf("count dry-run role: %v", err)
	}
	if count != 0 {
		t.Fatalf("dry-run inserted roles = %d, want 0", count)
	}

	metaAfter, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta(after) error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta(after) found = false, want true")
	}
	if metaAfter != metaBefore {
		t.Fatalf("projection_meta after dry run = %#v, want unchanged %#v", metaAfter, metaBefore)
	}
}

func TestDryRunProjectEventRollsBackOnProjectionFailure(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeProjectDB(t, conn)

	baseEvent := projectRoleEvent("EVT-20260613-000001", "insert", map[string]any{
		"description": "Developer",
		"active":      true,
	}, nil)
	baseHash, err := HashEvent(baseEvent)
	if err != nil {
		t.Fatalf("HashEvent(base) error = %v", err)
	}
	if err := ProjectEvent(conn, ProjectionInput{
		Event:     baseEvent,
		EventHash: baseHash,
		File:      "2026-06-13.jsonl",
		Line:      1,
		Offset:    100,
	}); err != nil {
		t.Fatalf("ProjectEvent(base) error = %v", err)
	}

	metaBefore, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta(before) error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta(before) found = false, want true")
	}

	err = DryRunProjectEvent(conn, Event{
		EventType: "issue.created",
		Actor: Actor{
			Role: "developer",
		},
		Command: Command{
			Argv:       []string{"fid", "issue", "create"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			projectIssueEffect("FID-priority-epic", map[string]any{
				"kind":        "epic",
				"title":       "Prioritised epic",
				"workflow_id": "epic-default",
				"status":      "todo",
				"created_at":  "2026-06-13T10:00:00Z",
				"updated_at":  "2026-06-13T10:00:00Z",
			}),
			{
				Op:    "insert",
				Table: "epic_metadata",
				PK: map[string]any{
					"issue_id": "FID-priority-epic",
				},
				Values: map[string]any{
					"scheduling": "prioritised",
				},
			},
			projectIssueEffect("FID-invalid-task", map[string]any{
				"kind":        "task",
				"parent_id":   "FID-priority-epic",
				"sequence":    1,
				"title":       "Invalid task",
				"workflow_id": "task-default",
				"status":      "todo",
				"created_at":  "2026-06-13T10:01:00Z",
				"updated_at":  "2026-06-13T10:01:00Z",
			}),
		},
	})
	if err == nil {
		t.Fatal("DryRunProjectEvent() error = nil, want projection rejection")
	}
	if !strings.Contains(err.Error(), "issues: task under prioritised epic must not have sequence") {
		t.Fatalf("DryRunProjectEvent() error = %q, want prioritised-epic sequence rejection", err.Error())
	}

	var issueCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM issues WHERE id IN (?, ?)`, "FID-priority-epic", "FID-invalid-task").Scan(&issueCount); err != nil {
		t.Fatalf("count dry-run issues: %v", err)
	}
	if issueCount != 0 {
		t.Fatalf("dry-run projected issues = %d, want 0", issueCount)
	}

	metaAfter, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta(after) error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta(after) found = false, want true")
	}
	if metaAfter != metaBefore {
		t.Fatalf("projection_meta after failed dry run = %#v, want unchanged %#v", metaAfter, metaBefore)
	}
}

func TestProjectHelpersValidateEffectsAndResults(t *testing.T) {
	t.Parallel()

	t.Run("precondition conflict formats details", func(t *testing.T) {
		t.Parallel()

		err := &PreconditionConflictError{
			EventID:   "EVT-20260605-000010",
			EventType: "issue.updated",
			Table:     "issues",
			PK:        map[string]any{"id": "FID-123456"},
			Expected:  map[string]any{"status": "todo"},
			Actual:    map[string]any{"status": "done"},
			File:      "2026-06-05.jsonl",
			Line:      2,
		}
		if !stringsContainsAll(err.Error(), []string{"2026-06-05.jsonl:2", "EVT-20260605-000010", "issue.updated", "issues"}) {
			t.Fatalf("PreconditionConflictError.Error() = %q, want structured details", err.Error())
		}
	})

	t.Run("ensure single row affected", func(t *testing.T) {
		t.Parallel()

		if err := ensureSingleRowAffected(fakeSQLResult{rowsAffected: 1}, "EVT-1", "update", "roles"); err != nil {
			t.Fatalf("ensureSingleRowAffected(one row) error = %v", err)
		}

		err := ensureSingleRowAffected(fakeSQLResult{rowsAffected: 0}, "EVT-1", "delete", "roles")
		if err == nil {
			t.Fatal("ensureSingleRowAffected(zero rows) error = nil, want failure")
		}

		err = ensureSingleRowAffected(fakeSQLResult{rowsAffectedErr: errors.New("boom")}, "EVT-1", "delete", "roles")
		if err == nil {
			t.Fatal("ensureSingleRowAffected(rows error) error = nil, want failure")
		}
	})

	t.Run("identifier and value helpers", func(t *testing.T) {
		t.Parallel()

		if _, err := validatedIdentifier("roles"); err != nil {
			t.Fatalf("validatedIdentifier(valid) error = %v", err)
		}
		if _, err := validatedIdentifier("bad-name"); err == nil {
			t.Fatal("validatedIdentifier(invalid) error = nil, want rejection")
		}

		if !valuesEqual([]byte("abc"), "abc") {
			t.Fatal("valuesEqual([]byte(\"abc\"), \"abc\") = false, want true")
		}
		if !valuesEqual(true, int64(1)) {
			t.Fatal("valuesEqual(true, int64(1)) = false, want true")
		}
		if !valuesEqual(int(3), float64(3)) {
			t.Fatal("valuesEqual(int(3), float64(3)) = false, want true")
		}
		if valuesEqual("todo", "done") {
			t.Fatal("valuesEqual(\"todo\", \"done\") = true, want false")
		}

		tests := []struct {
			value any
			want  float64
		}{
			{value: int8(1), want: 1},
			{value: int16(2), want: 2},
			{value: int32(3), want: 3},
			{value: uint(4), want: 4},
			{value: uint8(5), want: 5},
			{value: uint16(6), want: 6},
			{value: uint32(7), want: 7},
			{value: uint64(8), want: 8},
			{value: float32(9), want: 9},
		}
		for _, tt := range tests {
			got, ok := asFloat(tt.value)
			if !ok || got != tt.want {
				t.Fatalf("asFloat(%T(%v)) = (%v, %t), want (%v, true)", tt.value, tt.value, got, ok, tt.want)
			}
		}
		if _, ok := asFloat("nope"); ok {
			t.Fatal("asFloat(string) ok = true, want false")
		}
	})
}

func testProjectMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func closeProjectDB(t *testing.T, conn interface{ Close() error }) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func projectRoleEvent(eventID, op string, values, preconditions map[string]any) Event {
	return Event{
		EventID:     eventID,
		EventSchema: EventSchemaV1,
		EventType:   "project.initialized",
		OccurredAt:  time.Date(2026, time.June, 5, 14, 23, 11, 0, time.UTC),
		Actor: Actor{
			Role: "human",
		},
		Command: Command{
			Argv:       []string{"fid", "init"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:            op,
				Table:         "roles",
				PK:            map[string]any{"name": "developer"},
				Values:        values,
				Preconditions: preconditions,
			},
		},
	}
}

func projectIssueEffect(id string, values map[string]any) Effect {
	return Effect{
		Op:    "insert",
		Table: "issues",
		PK: map[string]any{
			"id": id,
		},
		Values: values,
	}
}

type fakeSQLResult struct {
	rowsAffected    int64
	rowsAffectedErr error
}

func (result fakeSQLResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (result fakeSQLResult) RowsAffected() (int64, error) {
	if result.rowsAffectedErr != nil {
		return 0, result.rowsAffectedErr
	}
	return result.rowsAffected, nil
}

func stringsContainsAll(value string, fragments []string) bool {
	for _, fragment := range fragments {
		if !strings.Contains(value, fragment) {
			return false
		}
	}
	return true
}
