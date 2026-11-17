package activity

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestListReturnsChronologicalEntriesWithStructuredDetails(t *testing.T) {
	t.Parallel()

	workspace := newActivityWorkspace(t)
	primaryIssueID := createActivityIssue(t, workspace, "Primary issue")
	secondaryIssueID := createActivityIssue(t, workspace, "Secondary issue")
	firstEntry := newActivitySeedIdentity(t, workspace, 100, 33*time.Hour)
	secondEntry := newActivitySeedIdentity(t, workspace, 101, 35*time.Hour)
	thirdEntry := newActivitySeedIdentity(t, workspace, 102, 35*time.Hour)

	secondEntry.eventID = appendActivityEvent(t, workspace, activitySeedEvent{
		activityID: secondEntry.activityID,
		refID:      secondaryIssueID,
		eventName:  "comment_added",
		actor:      "qa_reviewer",
		occurredAt: secondEntry.occurredAt,
		details:    map[string]any{"text": "Needs evidence"},
	})
	firstEntry.eventID = appendActivityEvent(t, workspace, activitySeedEvent{
		activityID: firstEntry.activityID,
		refID:      primaryIssueID,
		eventName:  "status_changed",
		actor:      "developer",
		occurredAt: firstEntry.occurredAt,
		details:    map[string]any{"from": "todo", "to": "in_progress"},
	})
	thirdEntry.eventID = appendActivityEvent(t, workspace, activitySeedEvent{
		activityID: thirdEntry.activityID,
		refID:      primaryIssueID,
		eventName:  "issue_updated",
		actor:      "developer",
		occurredAt: thirdEntry.occurredAt,
		details:    map[string]any{"field": "title", "to": "Refined title"},
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
			ID:             firstEntry.activityID,
			JournalEventID: firstEntry.eventID,
			RefType:        "issue",
			RefID:          primaryIssueID,
			Event:          "status-changed",
			Actor:          "developer",
			OccurredAt:     firstEntry.occurredAt,
			Details:        map[string]any{"from": "todo", "to": "in_progress"},
		},
		{
			ID:             secondEntry.activityID,
			JournalEventID: secondEntry.eventID,
			RefType:        "issue",
			RefID:          secondaryIssueID,
			Event:          "comment-added",
			Actor:          "qa_reviewer",
			OccurredAt:     secondEntry.occurredAt,
			Details:        map[string]any{"text": "Needs evidence"},
		},
		{
			ID:             thirdEntry.activityID,
			JournalEventID: thirdEntry.eventID,
			RefType:        "issue",
			RefID:          primaryIssueID,
			Event:          "issue-updated",
			Actor:          "developer",
			OccurredAt:     thirdEntry.occurredAt,
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
	firstEntry := newActivitySeedIdentity(t, workspace, 99, 23*time.Hour+59*time.Minute+59*time.Second)
	secondEntry := newActivitySeedIdentity(t, workspace, 100, 24*time.Hour)
	thirdEntry := newActivitySeedIdentity(t, workspace, 101, 32*time.Hour)
	secondOccurredAt, err := time.Parse(time.RFC3339, secondEntry.occurredAt)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", secondEntry.occurredAt, err)
	}

	appendActivityEvent(t, workspace, activitySeedEvent{
		activityID: firstEntry.activityID,
		refID:      primaryIssueID,
		eventName:  "status_changed",
		actor:      "developer",
		occurredAt: firstEntry.occurredAt,
		details:    map[string]any{"from": "todo", "to": "in_progress"},
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		activityID: secondEntry.activityID,
		refID:      primaryIssueID,
		eventName:  "status_changed",
		actor:      "developer",
		occurredAt: secondEntry.occurredAt,
		details:    map[string]any{"from": "in_progress", "to": "qa_review"},
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		activityID: thirdEntry.activityID,
		refID:      secondaryIssueID,
		eventName:  "issue_updated",
		actor:      "qa_reviewer",
		occurredAt: thirdEntry.occurredAt,
		details:    map[string]any{"field": "title", "to": "Evidence review"},
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
			wantIDs: []string{firstEntry.activityID, secondEntry.activityID},
		},
		{
			name: "since boundary includes exact match",
			options: ListOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testActivityMigrationsDir(t),
				Since:         secondEntry.occurredAt,
			},
			wantIDs: []string{secondEntry.activityID, thirdEntry.activityID},
		},
		{
			name: "since date only",
			options: ListOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testActivityMigrationsDir(t),
				Since:         secondOccurredAt.Format("2006-01-02"),
			},
			wantIDs: []string{secondEntry.activityID, thirdEntry.activityID},
		},
		{
			name: "event filter",
			options: ListOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testActivityMigrationsDir(t),
				Event:         "issue-updated",
			},
			wantIDs: []string{thirdEntry.activityID},
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
	activityID string
	refID      string
	eventName  string
	actor      string
	occurredAt string
	details    map[string]any
}

type activitySeedIdentity struct {
	eventID    string
	activityID string
	occurredAt string
}

func newActivityWorkspace(t *testing.T) string {
	t.Helper()

	workspace := gittest.NewWorkspace(t)
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
		WorkspaceRoot:  workspace,
		MigrationsDir:  testActivityMigrationsDir(t),
		Kind:           "task",
		Title:          title,
		Actor:          "agent:developer",
		IdempotencyKey: "activity-create-" + title,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeActivityDB(t, conn)

	if _, err := conn.Exec(`
		DELETE FROM activity_log
		WHERE ref_type = 'issue' AND ref_id = ? AND event = 'issue-created'
	`, result.ID); err != nil {
		t.Fatalf("delete create activity row: %v", err)
	}

	return result.ID
}

func newActivitySeedIdentity(t *testing.T, workspace string, sequence int, offset time.Duration) activitySeedIdentity {
	t.Helper()

	occurredAt := latestActivityJournalDate(t, workspace).Add(offset).UTC()

	// The Event ID is the writers' next one when the entry is appended.
	return activitySeedIdentity{
		activityID: fmt.Sprintf("ACT-%s-%06d", occurredAt.Format("20060102"), sequence),
		occurredAt: occurredAt.Format(time.RFC3339),
	}
}

// appendActivityEvent appends one seeded entry under the Event ID the
// writers allocate next, as they do, and returns that ID.
func appendActivityEvent(t *testing.T, workspace string, seed activitySeedEvent) string {
	t.Helper()

	occurredAt, err := time.Parse(time.RFC3339, seed.occurredAt)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", seed.occurredAt, err)
	}

	event := journal.Event{
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
					"journal_event_id": nil,
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

	var eventID string
	if err := journal.WithLock(workspace, journal.LockOwner{
		Command: "fid maintenance seed",
	}, func(lock *journal.WriteLock) error {
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		eventID = prepared.EventID
		prepared.Effects[0].Values["journal_event_id"] = eventID
		_, err = journal.AppendEvent(lock, prepared)
		return err
	}); err != nil {
		t.Fatalf("AppendEvent() error = %v", err)
	}
	return eventID
}

func testActivityMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func latestActivityJournalDate(t *testing.T, workspace string) time.Time {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(config.LedgerDir(workspace), "journal"))
	if err != nil {
		t.Fatalf("ReadDir(journal) error = %v", err)
	}

	var latest string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".jsonl" {
			continue
		}

		dateText := name[:len(name)-len(".jsonl")]
		if _, err := time.Parse("2006-01-02", dateText); err != nil {
			continue
		}
		if dateText > latest {
			latest = dateText
		}
	}

	if latest == "" {
		t.Fatal("latest journal date = empty")
	}

	date, err := time.Parse("2006-01-02", latest)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", latest, err)
	}
	return date.UTC()
}

func closeActivityDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
