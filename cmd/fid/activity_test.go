package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	internalactivity "github.com/mbell-nzl/fid/internal/activity"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestActivityListAppliesFilters(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	primaryIssueID := createActivityCommandIssue(t, workspace, "Primary issue")
	secondaryIssueID := createActivityCommandIssue(t, workspace, "Secondary issue")

	firstEntry := newActivityCommandSeedIdentity(t, workspace, 99, 23*time.Hour+59*time.Minute+59*time.Second)
	seedActivityCommandEntry(t, workspace, activityCommandSeed{
		eventID:    firstEntry.eventID,
		activityID: firstEntry.activityID,
		refID:      primaryIssueID,
		eventName:  "status_changed",
		actor:      "developer",
		occurredAt: firstEntry.occurredAt,
		details:    map[string]any{"from": "todo", "to": "in_progress"},
	})
	secondEntry := newActivityCommandSeedIdentity(t, workspace, 100, 24*time.Hour)
	seedActivityCommandEntry(t, workspace, activityCommandSeed{
		eventID:    secondEntry.eventID,
		activityID: secondEntry.activityID,
		refID:      primaryIssueID,
		eventName:  "status_changed",
		actor:      "developer",
		occurredAt: secondEntry.occurredAt,
		details:    map[string]any{"from": "in_progress", "to": "qa_review"},
	})
	thirdEntry := newActivityCommandSeedIdentity(t, workspace, 101, 32*time.Hour)
	seedActivityCommandEntry(t, workspace, activityCommandSeed{
		eventID:    thirdEntry.eventID,
		activityID: thirdEntry.activityID,
		refID:      secondaryIssueID,
		eventName:  "issue_updated",
		actor:      "qa_reviewer",
		occurredAt: thirdEntry.occurredAt,
		details:    map[string]any{"field": "title", "to": "Evidence review"},
	})

	tests := []struct {
		name    string
		args    []string
		wantIDs []string
	}{
		{
			name:    "on filter",
			args:    []string{"--json", "activity", "list", "--on", "issue/" + primaryIssueID},
			wantIDs: []string{firstEntry.activityID, secondEntry.activityID},
		},
		{
			name:    "since filter",
			args:    []string{"--json", "activity", "list", "--since", secondEntry.occurredAt},
			wantIDs: []string{secondEntry.activityID, thirdEntry.activityID},
		},
		{
			name:    "event filter",
			args:    []string{"--json", "activity", "list", "--event", "issue-updated"},
			wantIDs: []string{thirdEntry.activityID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, tt.args...)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			envelope := decodeCommandJSONStructure(t, stdout).(map[string]any)
			result := envelope["result"].(map[string]any)
			rows := result["items"].([]any)
			gotIDs := make([]string, 0, len(rows))
			for _, row := range rows {
				gotIDs = append(gotIDs, row.(map[string]any)["id"].(string))
			}
			if !reflect.DeepEqual(gotIDs, tt.wantIDs) {
				t.Fatalf("stdout ids = %#v, want %#v", gotIDs, tt.wantIDs)
			}
		})
	}
}

func TestActivityListUsesEnvelopeAndNormalizedActivityEvent(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	issueID := createActivityCommandIssue(t, workspace, "Primary issue")
	entry := newActivityCommandSeedIdentity(t, workspace, 100, 33*time.Hour)
	seedActivityCommandEntry(t, workspace, activityCommandSeed{
		eventID:    entry.eventID,
		activityID: entry.activityID,
		refID:      issueID,
		eventName:  "status_changed",
		actor:      "developer",
		occurredAt: entry.occurredAt,
		details:    map[string]any{"from": "todo", "to": "in_progress"},
	})

	stdout, stderr := runFIDCommand(t, workspace, "--json", "activity", "list")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	want := expectedCommandStructure(t, []internalactivity.Entry{
		{
			ID:             entry.activityID,
			JournalEventID: entry.eventID,
			RefType:        "issue",
			RefID:          issueID,
			Event:          "status-changed",
			Actor:          "developer",
			OccurredAt:     entry.occurredAt,
			Details:        map[string]any{"from": "todo", "to": "in_progress"},
		},
	})
	envelope := decodeCommandJSONStructure(t, stdout).(map[string]any)
	if got := envelope["result"].(map[string]any)["items"]; !reflect.DeepEqual(got, want) {
		t.Fatalf("activity list json = %#v, want %#v", got, want)
	}
}

type activityCommandSeed struct {
	eventID    string
	activityID string
	refID      string
	eventName  string
	actor      string
	occurredAt string
	details    map[string]any
}

type activityCommandSeedIdentity struct {
	eventID    string
	activityID string
	occurredAt string
}

func createActivityCommandIssue(t *testing.T, workspace, title string) string {
	t.Helper()

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: issueMigrationsDir(),
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
	defer closeCommandDB(t, conn)

	if _, err := conn.Exec(`
		DELETE FROM activity_log
		WHERE ref_type = 'issue' AND ref_id = ? AND event = 'issue-created'
	`, created.ID); err != nil {
		t.Fatalf("delete create activity row: %v", err)
	}

	return created.ID
}

func newActivityCommandSeedIdentity(t *testing.T, workspace string, sequence int, offset time.Duration) activityCommandSeedIdentity {
	t.Helper()

	baseDate := latestActivityCommandJournalDate(t, workspace)
	occurredAt := baseDate.Add(offset).UTC()

	return activityCommandSeedIdentity{
		eventID:    fmt.Sprintf("EVT-%s-%06d", occurredAt.Format("20060102"), sequence),
		activityID: fmt.Sprintf("ACT-%s-%06d", occurredAt.Format("20060102"), sequence),
		occurredAt: occurredAt.Format(time.RFC3339),
	}
}

func latestActivityCommandJournalDate(t *testing.T, workspace string) time.Time {
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

func seedActivityCommandEntry(t *testing.T, workspace string, seed activityCommandSeed) {
	t.Helper()

	occurredAt, err := time.Parse(time.RFC3339, seed.occurredAt)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", seed.occurredAt, err)
	}

	appendCommandIssueEvent(t, workspace, journal.Event{
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
					"journal_event_id": seed.eventID,
					"ref_type":         "issue",
					"ref_id":           seed.refID,
					"event":            seed.eventName,
					"actor":            seed.actor,
					"occurred_at":      seed.occurredAt,
					"details":          seed.details,
				},
			},
		},
	})
}
