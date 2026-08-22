package journal

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestFindEventReturnsStoredEventLocationAndHash(t *testing.T) {
	journalDir := filepath.Join(t.TempDir(), ".fid", "journal")
	mustMkdirAll(t, journalDir)

	first := writerRoleEvent("EVT-20260605-000001")
	target := Event{
		EventID:     "EVT-20260606-000002",
		EventSchema: EventSchemaV1,
		EventType:   "issue.transitioned",
		OccurredAt:  time.Date(2026, time.June, 6, 10, 30, 0, 0, time.UTC),
		Actor: Actor{
			Role: "developer",
		},
		Command: Command{
			Argv:       []string{"fid", "issue", "transition", "FID-a1b2c3", "--to", "in_progress", "--by", "developer"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "update",
				Table: "issues",
				PK:    map[string]any{"id": "FID-a1b2c3"},
				Values: map[string]any{
					"status": "in_progress",
				},
				Prior: map[string]any{
					"status": "todo",
				},
				Preconditions: map[string]any{
					"status":      "todo",
					"workflow_id": "task-default",
				},
			},
		},
	}
	last := writerRoleEvent("EVT-20260606-000003")

	writeWriterEvents(t, filepath.Join(journalDir, "2026-06-05.jsonl"), first)
	writeWriterEvents(t, filepath.Join(journalDir, "2026-06-06.jsonl"), target, last)

	got, err := FindEvent(journalDir, target.EventID)
	if err != nil {
		t.Fatalf("FindEvent() error = %v", err)
	}

	if !reflect.DeepEqual(got.Event, target) {
		t.Fatalf("FindEvent().Event = %#v, want %#v", got.Event, target)
	}
	if got.JournalFile != ".fid/journal/2026-06-06.jsonl" {
		t.Fatalf("FindEvent().JournalFile = %q, want %q", got.JournalFile, ".fid/journal/2026-06-06.jsonl")
	}
	if got.Line != 1 {
		t.Fatalf("FindEvent().Line = %d, want %d", got.Line, 1)
	}

	wantHash, err := HashEvent(target)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}
	if got.Hash != wantHash {
		t.Fatalf("FindEvent().Hash = %q, want %q", got.Hash, wantHash)
	}
}

func TestFindEventReturnsTypedNotFoundError(t *testing.T) {
	journalDir := filepath.Join(t.TempDir(), ".fid", "journal")
	mustMkdirAll(t, journalDir)
	writeWriterEvents(t, filepath.Join(journalDir, "2026-06-05.jsonl"), writerRoleEvent("EVT-20260605-000001"))

	_, err := FindEvent(journalDir, "EVT-20260605-999999")
	if err == nil {
		t.Fatal("FindEvent() error = nil, want not found")
	}

	var target *EventNotFoundError
	if !errors.As(err, &target) {
		t.Fatalf("FindEvent() error = %T, want *EventNotFoundError", err)
	}
	if target.EventID != "EVT-20260605-999999" {
		t.Fatalf("EventNotFoundError.EventID = %q, want %q", target.EventID, "EVT-20260605-999999")
	}
}
