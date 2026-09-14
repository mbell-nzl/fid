package journal

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestRecoverySchemaFreshInitAndReplay(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() { _ = conn.Close() }()
	if !sqliteTableExists(t, conn, "run_recoveries") {
		t.Fatal("run_recoveries table absent after fresh init")
	}

	event := unknownRecoveryEvent()
	eventHash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}
	err = ProjectEvent(conn, ProjectionInput{Event: event, EventHash: eventHash})
	assertUnknownRecoveryEvent(t, err)
	assertInvalidStopCodeRejectsLiveAndReplay(t, conn)

	workspace := newReplayWorkspace(t)
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-05.jsonl"), event)
	replayDBPath := filepath.Join(t.TempDir(), "replay.sqlite")
	if err := fiddb.MigrateUp(replayDBPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp(replay) error = %v", err)
	}
	replayConn, err := fiddb.Open(replayDBPath)
	if err != nil {
		t.Fatalf("Open(replay) error = %v", err)
	}
	defer func() { _ = replayConn.Close() }()
	_, err = ReplayJournal(replayConn, filepath.Join(workspace, ".fid", "journal"), nil)
	assertUnknownRecoveryEvent(t, err)
}

func assertInvalidStopCodeRejectsLiveAndReplay(t *testing.T, conn *sql.DB) {
	t.Helper()
	event := invalidStopCodeEvent()
	eventHash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent(invalid stop) error = %v", err)
	}
	err = ProjectEvent(conn, ProjectionInput{Event: event, EventHash: eventHash})
	if err == nil || !strings.Contains(err.Error(), "CHECK constraint failed") {
		t.Fatalf("ProjectEvent(invalid stop) error = %v, want Stop code CHECK rejection", err)
	}

	workspace := newReplayWorkspace(t)
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-05.jsonl"), event)
	replayDBPath := filepath.Join(t.TempDir(), "invalid-stop-replay.sqlite")
	if err := fiddb.MigrateUp(replayDBPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp(invalid stop replay) error = %v", err)
	}
	replayConn, err := fiddb.Open(replayDBPath)
	if err != nil {
		t.Fatalf("Open(invalid stop replay) error = %v", err)
	}
	defer func() { _ = replayConn.Close() }()
	_, err = ReplayJournal(replayConn, filepath.Join(workspace, ".fid", "journal"), nil)
	var corrupt *CorruptJournalError
	if !errors.As(err, &corrupt) || !strings.Contains(corrupt.Error(), "CHECK constraint failed") {
		t.Fatalf("ReplayJournal(invalid stop) error = %v, want corrupt Stop code CHECK rejection", err)
	}
}

func sqliteTableExists(t *testing.T, conn *sql.DB, table string) bool {
	t.Helper()
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&count); err != nil {
		t.Fatalf("query sqlite_master for %q: %v", table, err)
	}
	return count == 1
}

func assertUnknownRecoveryEvent(t *testing.T, err error) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), "unsupported effect op \"unknown\"") {
		t.Fatalf("recovery Event error = %v, want unknown Event effect rejection", err)
	}
	var corrupt *CorruptJournalError
	if errors.As(err, &corrupt) && !strings.Contains(corrupt.Error(), "unsupported effect op \"unknown\"") {
		t.Fatalf("ReplayJournal() error = %v, want unknown Event effect rejection", err)
	}
}

func unknownRecoveryEvent() Event {
	return Event{
		EventID:     "EVT-20260905-000001",
		EventSchema: EventSchemaV2,
		EventType:   "unknown.recovery-event",
		OccurredAt:  time.Date(2026, time.September, 5, 0, 0, 0, 0, time.UTC),
		ActorV2: EventActor{
			ActorID:       "actor-recovery",
			ActorType:     "agent",
			Producer:      "fid",
			ProvenanceRef: "test:recovery",
		},
		Aggregate:   Aggregate{ID: "RUN-prior", Sequence: 1, Type: "run"},
		CommandName: "run.observe-budget",
		Effects: []Effect{{
			Op:    "unknown",
			Table: "run_recoveries",
			PK:    map[string]any{"id": "REC-1"},
		}},
	}
}

func invalidStopCodeEvent() Event {
	return Event{
		EventID:     "EVT-20260905-000002",
		EventSchema: EventSchemaV1,
		EventType:   "project.initialized",
		OccurredAt:  time.Date(2026, time.September, 5, 0, 1, 0, 0, time.UTC),
		Actor:       Actor{Role: "developer"},
		Command:     Command{Argv: []string{"fid", "test"}, CWD: ".", FIDVersion: "0.1.0"},
		Effects: []Effect{{
			Op:    "insert",
			Table: "stops",
			PK:    map[string]any{"id": "STOP-invalid"},
			Values: map[string]any{
				"subject_type":              "run",
				"subject_id":                "RUN-1",
				"code":                      "invalid",
				"blocked_actions_json":      "[]",
				"details_json":              "{}",
				"resolution_condition_json": "{}",
				"resolver_capability":       "stop.resolve",
				"actor_id":                  "actor-1",
				"raised_at":                 "2026-09-05T00:01:00Z",
				"created_event_id":          "EVT-stop-invalid",
			},
		}},
	}
}
