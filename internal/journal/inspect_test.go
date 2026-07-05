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

func TestRevertVerdictClassifiesKnownEventTypes(t *testing.T) {
	transition := RevertVerdict(Event{
		EventType: "issue.transitioned",
		Effects: []Effect{{
			Op:    "update",
			Table: "issues",
			PK:    map[string]any{"id": "FID-a1b2c3"},
			Values: map[string]any{
				"status":     "in_progress",
				"updated_at": "2026-06-06T10:30:00Z",
			},
			Prior: map[string]any{
				"status":     "todo",
				"updated_at": "2026-06-06T10:00:00Z",
			},
		}},
	})
	if !transition.Eligible || transition.Mode != "inverse" {
		t.Fatalf("transition verdict = %#v, want eligible inverse", transition)
	}
	if got := transition.Restores["status"]; got != "todo" {
		t.Fatalf("transition restores[status] = %#v, want %q", got, "todo")
	}
	if got := transition.Precondition["status"]; got != "in_progress" {
		t.Fatalf("transition precondition[status] = %#v, want %q", got, "in_progress")
	}

	acVerified := RevertVerdict(Event{EventType: "issue.ac_verified"})
	if !acVerified.Eligible || acVerified.Mode != "retraction" {
		t.Fatalf("ac verified verdict = %#v, want eligible retraction", acVerified)
	}
	if acVerified.Note == "" {
		t.Fatalf("ac verified note = empty, want explanation")
	}

	created := RevertVerdict(Event{EventType: "issue.created"})
	if created.Eligible {
		t.Fatalf("created verdict = %#v, want not eligible", created)
	}
	if created.Correction == "" {
		t.Fatalf("created correction = empty, want forward path")
	}

	claim := RevertVerdict(Event{EventType: "issue.claim.acquired"})
	if claim.Eligible {
		t.Fatalf("claim verdict = %#v, want not eligible", claim)
	}
	if claim.Correction == "" {
		t.Fatalf("claim correction = empty, want forward path")
	}

	composite := RevertVerdict(Event{EventType: "dev.started"})
	if composite.Eligible {
		t.Fatalf("composite verdict = %#v, want not eligible", composite)
	}
	if composite.Correction == "" {
		t.Fatalf("composite correction = empty, want forward path")
	}

	labelAdded := RevertVerdict(Event{
		EventType: "issue.label.added",
		Effects: []Effect{{
			Op:    "insert",
			Table: "issue_labels",
			PK: map[string]any{
				"issue_id": "FID-a1b2c3",
				"label":    "blocked",
			},
		}},
	})
	if !labelAdded.Eligible || labelAdded.Mode != "inverse" {
		t.Fatalf("label added verdict = %#v, want eligible inverse", labelAdded)
	}
	if got := labelAdded.Restores["labels_absent"]; got == nil {
		t.Fatalf("label added restores = %#v, want labels_absent entry", labelAdded.Restores)
	}
	if got := labelAdded.Precondition["labels_present"]; got == nil {
		t.Fatalf("label added precondition = %#v, want labels_present entry", labelAdded.Precondition)
	}

	labelRemoved := RevertVerdict(Event{
		EventType: "issue.label.removed",
		Effects: []Effect{{
			Op:    "delete",
			Table: "issue_labels",
			PK: map[string]any{
				"issue_id": "FID-a1b2c3",
				"label":    "blocked",
			},
		}},
	})
	if !labelRemoved.Eligible || labelRemoved.Mode != "inverse" {
		t.Fatalf("label removed verdict = %#v, want eligible inverse", labelRemoved)
	}
	if got := labelRemoved.Restores["labels_present"]; got == nil {
		t.Fatalf("label removed restores = %#v, want labels_present entry", labelRemoved.Restores)
	}
	if got := labelRemoved.Precondition["labels_absent"]; got == nil {
		t.Fatalf("label removed precondition = %#v, want labels_absent entry", labelRemoved.Precondition)
	}
}

func TestRevertVerdictDefaultsUnknownTypesToSafeNotEligible(t *testing.T) {
	got := RevertVerdict(Event{EventType: "mystery.event"})

	if got.Eligible {
		t.Fatalf("unknown verdict = %#v, want not eligible", got)
	}
	if got.Correction == "" {
		t.Fatalf("unknown correction = empty, want safe default")
	}
	if got.Note == "" {
		t.Fatalf("unknown note = empty, want safe default note")
	}
}
