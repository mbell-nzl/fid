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

func TestRevertVerdictKeepsSpecificStrategiesForSupportedEventTypes(t *testing.T) {
	tests := []struct {
		name             string
		event            Event
		wantEligible     bool
		wantMode         string
		wantRestoreKey   string
		wantPrecondition string
		wantCorrection   string
	}{
		{
			name:             "issue update without field effects uses guarded fallback",
			event:            Event{EventType: "issue.updated"},
			wantEligible:     true,
			wantMode:         "inverse",
			wantRestoreKey:   "prior_state",
			wantPrecondition: "current_state",
		},
		{
			name:             "concern address without field effects uses guarded fallback",
			event:            Event{EventType: "issue.concern_addressed"},
			wantEligible:     true,
			wantMode:         "inverse",
			wantRestoreKey:   "prior_state",
			wantPrecondition: "current_state",
		},
		{
			name: "link add removes only present rows",
			event: Event{EventType: "link.added", Effects: []Effect{{
				Table: "links",
				PK:    map[string]any{"from_id": "FID-a1b2c3", "to_id": "FID-d4e5f6"},
			}}},
			wantEligible:     true,
			wantMode:         "inverse",
			wantRestoreKey:   "links_absent",
			wantPrecondition: "links_present",
		},
		{
			name: "link removal restores only absent rows",
			event: Event{EventType: "link.removed", Effects: []Effect{{
				Table: "links",
				PK:    map[string]any{"from_id": "FID-a1b2c3", "to_id": "FID-d4e5f6"},
			}}},
			wantEligible:     true,
			wantMode:         "inverse",
			wantRestoreKey:   "links_present",
			wantPrecondition: "links_absent",
		},
		{
			name:         "acceptance criterion addition retracts fact",
			event:        Event{EventType: "issue.ac_added"},
			wantEligible: true,
			wantMode:     "retraction",
		},
		{
			name:         "concern raise retracts fact",
			event:        Event{EventType: "issue.concern_raised"},
			wantEligible: true,
			wantMode:     "retraction",
		},
		{
			name:         "gate result retracts fact",
			event:        Event{EventType: "dev.gate.recorded"},
			wantEligible: true,
			wantMode:     "retraction",
		},
		{
			name:           "workspace initialization stays non-revertible",
			event:          Event{EventType: "project.initialized"},
			wantCorrection: "create a fresh workspace if the project bootstrap was wrong",
		},
		{
			name:           "label creation stays non-revertible",
			event:          Event{EventType: "label.added"},
			wantCorrection: "stop applying the label, or create the corrected label and move usage forward",
		},
		{
			name:           "completion stays non-revertible",
			event:          Event{EventType: "issue.completed"},
			wantCorrection: "move the issue forward with the appropriate workflow command instead of removing the completion history",
		},
		{
			name:           "run telemetry stays non-revertible",
			event:          Event{EventType: "run.recorded"},
			wantCorrection: "record a correcting run entry or leave the telemetry history intact",
		},
		{
			name:           "requirement sync stays non-revertible",
			event:          Event{EventType: "requirement.synced"},
			wantCorrection: "rerun fid requirement sync from the corrected source document",
		},
		{
			name:           "developer finish stays non-revertible",
			event:          Event{EventType: "dev.finished"},
			wantCorrection: "use the issue and dev commands to move workflow state forward instead of reverting the composite event",
		},
		{
			name:           "developer block stays non-revertible",
			event:          Event{EventType: "dev.blocked"},
			wantCorrection: "use the issue and dev commands to move workflow state forward instead of reverting the composite event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RevertVerdict(tt.event)
			if got.Eligible != tt.wantEligible {
				t.Fatalf("eligible = %t, want %t", got.Eligible, tt.wantEligible)
			}
			if got.Mode != tt.wantMode {
				t.Fatalf("mode = %q, want %q", got.Mode, tt.wantMode)
			}
			if tt.wantRestoreKey != "" && got.Restores[tt.wantRestoreKey] == nil {
				t.Fatalf("restores = %#v, want %q entry", got.Restores, tt.wantRestoreKey)
			}
			if tt.wantPrecondition != "" && got.Precondition[tt.wantPrecondition] == nil {
				t.Fatalf("precondition = %#v, want %q entry", got.Precondition, tt.wantPrecondition)
			}
			if got.Correction != tt.wantCorrection {
				t.Fatalf("correction = %q, want %q", got.Correction, tt.wantCorrection)
			}
			if got.Note == "" {
				t.Fatal("note = empty, want an operator explanation")
			}
		})
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
