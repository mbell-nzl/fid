package journal

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAppendEventRejectsV2AggregateSequenceGap(t *testing.T) {
	workspace := t.TempDir()
	if err := os.MkdirAll(filepath.Join(workspace, ".fid", "journal"), 0o755); err != nil {
		t.Fatalf("MkdirAll(journal) error = %v", err)
	}

	err := WithLock(workspace, LockOwner{Command: "test Event-v2 sequence"}, func(lock *WriteLock) error {
		_, appendErr := AppendEvent(lock, eventV2ForSequence(2))
		return appendErr
	})
	if err == nil {
		t.Fatal("AppendEvent() error = nil, want aggregate sequence gap rejection")
	}
	if !strings.Contains(err.Error(), "aggregate activation/ACT-1 sequence 2, want 1") {
		t.Fatalf("AppendEvent() error = %q, want aggregate sequence rejection", err)
	}
}

func TestReplayJournalRejectsOutOfOrderV2AggregateSequence(t *testing.T) {
	workspace := newReplayWorkspace(t)
	journalPath := filepath.Join(workspace, ".fid", "journal", "2026-08-22.jsonl")
	first := eventV2ForSequence(1)
	first.EventID = "EVT-20260822-000001"
	first.OccurredAt = time.Date(2026, time.August, 22, 0, 0, 0, 0, time.UTC)
	second := eventV2ForSequence(3)
	second.EventID = "EVT-20260822-000002"
	second.OccurredAt = time.Date(2026, time.August, 22, 0, 1, 0, 0, time.UTC)
	writeReplayEvents(t, journalPath, first, second)

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)
	_, err := ReplayJournal(conn, filepath.Join(workspace, ".fid", "journal"), nil)
	var corrupt *CorruptJournalError
	if !errors.As(err, &corrupt) {
		t.Fatalf("ReplayJournal() error = %T %v, want *CorruptJournalError", err, err)
	}
	if corrupt.Line != 2 || !strings.Contains(corrupt.Error(), "aggregate activation/ACT-1 sequence 3, want 2") {
		t.Fatalf("ReplayJournal() error = %v, want second-line aggregate sequence rejection", corrupt)
	}
}

func TestEventV2RejectsUnknownTopLevelKey(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(`{"actor":{"actor_id":"fid-bootstrap","actor_type":"service","producer":"fid","provenance_ref":"fid:init","session_ref":null},"aggregate":{"id":"project-1","sequence":1,"type":"project"},"causation_id":null,"command":"project.init","correlation_id":null,"effects":[],"event_id":"EVT-20260822-000001","event_schema":2,"event_type":"project.initialized","idempotency":null,"occurred_at":"2026-08-22T00:00:00Z","repository":null,"unknown":true}`), &event)
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("Unmarshal(Event-v2 with unknown key) error = %v, want unknown-field rejection", err)
	}
}

func eventV2ForSequence(sequence int) Event {
	event := testIdempotentActivationEvent()
	event.EventID = "EVT-20260822-000001"
	event.OccurredAt = time.Date(2026, time.August, 22, 0, 0, 0, 0, time.UTC)
	event.Aggregate.Sequence = sequence
	event = currentProjectionIdempotency(event, "ACT-1")
	event.Effects[len(event.Effects)-1].Values["result_json"] = testActivationResultJSON("projection-fixture", false)
	return event
}
