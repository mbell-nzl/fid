package journal

import (
	"path/filepath"
	"runtime"
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
