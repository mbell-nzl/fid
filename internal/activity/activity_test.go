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
	fiddb "github.com/mbell-nzl/fid/internal/db"
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

	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        secondEntry.eventID,
		activityID:     secondEntry.activityID,
		refID:          secondaryIssueID,
		eventName:      "comment_added",
		actor:          "qa_reviewer",
		occurredAt:     secondEntry.occurredAt,
		details:        map[string]any{"text": "Needs evidence"},
		journalEventID: secondEntry.eventID,
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        firstEntry.eventID,
		activityID:     firstEntry.activityID,
		refID:          primaryIssueID,
		eventName:      "status_changed",
		actor:          "developer",
		occurredAt:     firstEntry.occurredAt,
		details:        map[string]any{"from": "todo", "to": "in_progress"},
		journalEventID: firstEntry.eventID,
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        thirdEntry.eventID,
		activityID:     thirdEntry.activityID,
		refID:          primaryIssueID,
		eventName:      "issue_updated",
		actor:          "developer",
		occurredAt:     thirdEntry.occurredAt,
		details:        map[string]any{"field": "title", "to": "Refined title"},
		journalEventID: thirdEntry.eventID,
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
			Event:          "status_changed",
			Actor:          "developer",
			OccurredAt:     firstEntry.occurredAt,
			Details:        map[string]any{"from": "todo", "to": "in_progress"},
		},
		{
			ID:             secondEntry.activityID,
			JournalEventID: secondEntry.eventID,
			RefType:        "issue",
			RefID:          secondaryIssueID,
			Event:          "comment_added",
			Actor:          "qa_reviewer",
			OccurredAt:     secondEntry.occurredAt,
			Details:        map[string]any{"text": "Needs evidence"},
		},
		{
			ID:             thirdEntry.activityID,
			JournalEventID: thirdEntry.eventID,
			RefType:        "issue",
			RefID:          primaryIssueID,
			Event:          "issue_updated",
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
		eventID:        firstEntry.eventID,
		activityID:     firstEntry.activityID,
		refID:          primaryIssueID,
		eventName:      "status_changed",
		actor:          "developer",
		occurredAt:     firstEntry.occurredAt,
		details:        map[string]any{"from": "todo", "to": "in_progress"},
		journalEventID: firstEntry.eventID,
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        secondEntry.eventID,
		activityID:     secondEntry.activityID,
		refID:          primaryIssueID,
		eventName:      "status_changed",
		actor:          "developer",
		occurredAt:     secondEntry.occurredAt,
		details:        map[string]any{"from": "in_progress", "to": "qa_review"},
		journalEventID: secondEntry.eventID,
	})
	appendActivityEvent(t, workspace, activitySeedEvent{
		eventID:        thirdEntry.eventID,
		activityID:     thirdEntry.activityID,
		refID:          secondaryIssueID,
		eventName:      "issue_updated",
		actor:          "qa_reviewer",
		occurredAt:     thirdEntry.occurredAt,
		details:        map[string]any{"field": "title", "to": "Evidence review"},
		journalEventID: thirdEntry.eventID,
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
				Event:         "issue_updated",
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
	eventID        string
	activityID     string
	refID          string
	eventName      string
	actor          string
	occurredAt     string
	details        map[string]any
	journalEventID string
}

type activitySeedIdentity struct {
	eventID    string
	activityID string
	occurredAt string
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

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeActivityDB(t, conn)

	if _, err := conn.Exec(`
		DELETE FROM activity_log
		WHERE ref_type = 'issue' AND ref_id = ? AND event = 'issue_created'
	`, result.ID); err != nil {
		t.Fatalf("delete create activity row: %v", err)
	}

	return result.ID
}

func newActivitySeedIdentity(t *testing.T, workspace string, sequence int, offset time.Duration) activitySeedIdentity {
	t.Helper()

	occurredAt := latestActivityJournalDate(t, workspace).Add(offset).UTC()

	return activitySeedIdentity{
		eventID:    fmt.Sprintf("EVT-%s-%06d", occurredAt.Format("20060102"), sequence),
		activityID: fmt.Sprintf("ACT-%s-%06d", occurredAt.Format("20060102"), sequence),
		occurredAt: occurredAt.Format(time.RFC3339),
	}
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

func latestActivityJournalDate(t *testing.T, workspace string) time.Time {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(workspace, ".fid", "journal"))
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
