package activity

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestListReturnsChronologicalEntriesWithStructuredDetails(t *testing.T) {
	t.Parallel()

	workspace := newActivityWorkspace(t)
	primaryIssueID := createActivityIssue(t, workspace, "Primary issue")
	secondaryIssueID := createActivityIssue(t, workspace, "Secondary issue")

	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        "EVT-20260705-000101",
		activityID:     "ACT-20260705-000101",
		refID:          secondaryIssueID,
		eventName:      "comment_added",
		actor:          "qa_reviewer",
		occurredAt:     "2026-07-05T09:02:00Z",
		details:        map[string]any{"text": "Needs evidence"},
		journalEventID: "EVT-20260705-000101",
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        "EVT-20260705-000100",
		activityID:     "ACT-20260705-000100",
		refID:          primaryIssueID,
		eventName:      "status_changed",
		actor:          "developer",
		occurredAt:     "2026-07-05T09:00:00Z",
		details:        map[string]any{"from": "todo", "to": "in_progress"},
		journalEventID: "EVT-20260705-000100",
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        "EVT-20260705-000102",
		activityID:     "ACT-20260705-000102",
		refID:          primaryIssueID,
		eventName:      "issue_updated",
		actor:          "developer",
		occurredAt:     "2026-07-05T09:02:00Z",
		details:        map[string]any{"field": "title", "to": "Refined title"},
		journalEventID: "EVT-20260705-000102",
	})

	got, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testActivityMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := []Entry{
		{
			ID:             "ACT-20260705-000100",
			JournalEventID: "EVT-20260705-000100",
			RefType:        "issue",
			RefID:          primaryIssueID,
			Event:          "status_changed",
			Actor:          "developer",
			OccurredAt:     "2026-07-05T09:00:00Z",
			Details:        map[string]any{"from": "todo", "to": "in_progress"},
		},
		{
			ID:             "ACT-20260705-000101",
			JournalEventID: "EVT-20260705-000101",
			RefType:        "issue",
			RefID:          secondaryIssueID,
			Event:          "comment_added",
			Actor:          "qa_reviewer",
			OccurredAt:     "2026-07-05T09:02:00Z",
			Details:        map[string]any{"text": "Needs evidence"},
		},
		{
			ID:             "ACT-20260705-000102",
			JournalEventID: "EVT-20260705-000102",
			RefType:        "issue",
			RefID:          primaryIssueID,
			Event:          "issue_updated",
			Actor:          "developer",
			OccurredAt:     "2026-07-05T09:02:00Z",
			Details:        map[string]any{"field": "title", "to": "Refined title"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %#v, want %#v", got, want)
	}
}

func TestListFiltersByOnSinceAndEvent(t *testing.T) {
	t.Parallel()

	workspace := newActivityWorkspace(t)
	primaryIssueID := createActivityIssue(t, workspace, "Primary issue")
	secondaryIssueID := createActivityIssue(t, workspace, "Secondary issue")

	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        "EVT-20260704-000099",
		activityID:     "ACT-20260704-000099",
		refID:          primaryIssueID,
		eventName:      "status_changed",
		actor:          "developer",
		occurredAt:     "2026-07-04T23:59:59Z",
		details:        map[string]any{"from": "todo", "to": "in_progress"},
		journalEventID: "EVT-20260704-000099",
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        "EVT-20260705-000100",
		activityID:     "ACT-20260705-000100",
		refID:          primaryIssueID,
		eventName:      "status_changed",
		actor:          "developer",
		occurredAt:     "2026-07-05T00:00:00Z",
		details:        map[string]any{"from": "in_progress", "to": "qa_review"},
		journalEventID: "EVT-20260705-000100",
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        "EVT-20260705-000101",
		activityID:     "ACT-20260705-000101",
		refID:          secondaryIssueID,
		eventName:      "issue_updated",
		actor:          "qa_reviewer",
		occurredAt:     "2026-07-05T08:00:00Z",
		details:        map[string]any{"field": "title", "to": "Evidence review"},
		journalEventID: "EVT-20260705-000101",
	})

	tests := []struct {
		name    string
		options ListOptions
		wantIDs []string
	}{
		{
			name: "on filter",
			options: ListOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testActivityMigrationsDir(t),
				On:            "issue/" + primaryIssueID,
			},
			wantIDs: []string{"ACT-20260704-000099", "ACT-20260705-000100"},
		},
		{
			name: "since boundary includes exact match",
			options: ListOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testActivityMigrationsDir(t),
				Since:         "2026-07-05T00:00:00Z",
			},
			wantIDs: []string{"ACT-20260705-000100", "ACT-20260705-000101"},
		},
		{
			name: "since date only",
			options: ListOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testActivityMigrationsDir(t),
				Since:         "2026-07-05",
			},
			wantIDs: []string{"ACT-20260705-000100", "ACT-20260705-000101"},
		},
		{
			name: "event filter",
			options: ListOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testActivityMigrationsDir(t),
				Event:         "issue_updated",
			},
			wantIDs: []string{"ACT-20260705-000101"},
		},
		{
			name: "empty result is non nil slice",
			options: ListOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testActivityMigrationsDir(t),
				On:            "issue/FID-missing",
			},
			wantIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.options)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}
			if got == nil {
				t.Fatal("List() returned nil slice, want empty slice")
			}

			gotIDs := make([]string, 0, len(got))
			for _, entry := range got {
				gotIDs = append(gotIDs, entry.ID)
			}
			if !reflect.DeepEqual(gotIDs, tt.wantIDs) {
				t.Fatalf("List() ids = %#v, want %#v", gotIDs, tt.wantIDs)
			}
		})
	}
}

type activitySeedEvent struct {
	eventID        string
	activityID     string
	refID          string
	eventName      string
	actor          string
	occurredAt     string
	details        map[string]any
	journalEventID string
}

func newActivityWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testActivityMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}

	return workspace
}

func createActivityIssue(t *testing.T, workspace, title string) string {
	t.Helper()

	result, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testActivityMigrationsDir(t),
		Kind:          "task",
		Title:         title,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	return result.ID
}

func appendActivityEvent(t *testing.T, workspace string, seed activitySeedEvent) {
	t.Helper()

	occurredAt, err := time.Parse(time.RFC3339, seed.occurredAt)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", seed.occurredAt, err)
	}

	event := journal.Event{
		EventID:     seed.eventID,
		EventSchema: journal.EventSchemaV1,
		EventType:   "activity.seeded",
		OccurredAt:  occurredAt,
		Actor: journal.Actor{
			Role: seed.actor,
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []journal.Effect{
			{
				Op:    "activity",
				Table: "activity_log",
				PK: map[string]any{
					"id": seed.activityID,
				},
				Values: map[string]any{
					"journal_event_id": seed.journalEventID,
					"ref_type":         "issue",
					"ref_id":           seed.refID,
					"event":            seed.eventName,
					"actor":            seed.actor,
					"occurred_at":      seed.occurredAt,
					"details":          seed.details,
				},
			},
		},
	}

	if err := journal.WithLock(workspace, journal.LockOwner{
		Command:   "fid maintenance seed",
		ActorRole: seed.actor,
	}, func(lock *journal.WriteLock) error {
		_, err := journal.AppendEvent(lock, event)
		return err
	}); err != nil {
		t.Fatalf("AppendEvent() error = %v", err)
	}
}

func testActivityMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}
