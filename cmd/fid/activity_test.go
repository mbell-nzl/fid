package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	internalactivity "github.com/mbell-nzl/fid/internal/activity"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestActivityListAppliesFilters(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	primaryIssueID := createActivityCommandIssue(t, workspace, "Primary issue")
	secondaryIssueID := createActivityCommandIssue(t, workspace, "Secondary issue")

	seedActivityCommandEntry(t, workspace, activityCommandSeed{
		eventID:    "EVT-20260704-000099",
		activityID: "ACT-20260704-000099",
		refID:      primaryIssueID,
		eventName:  "status_changed",
		actor:      "developer",
		occurredAt: "2026-07-04T23:59:59Z",
		details:    map[string]any{"from": "todo", "to": "in_progress"},
	})
	seedActivityCommandEntry(t, workspace, activityCommandSeed{
		eventID:    "EVT-20260705-000100",
		activityID: "ACT-20260705-000100",
		refID:      primaryIssueID,
		eventName:  "status_changed",
		actor:      "developer",
		occurredAt: "2026-07-05T00:00:00Z",
		details:    map[string]any{"from": "in_progress", "to": "qa_review"},
	})
	seedActivityCommandEntry(t, workspace, activityCommandSeed{
		eventID:    "EVT-20260705-000101",
		activityID: "ACT-20260705-000101",
		refID:      secondaryIssueID,
		eventName:  "issue_updated",
		actor:      "qa_reviewer",
		occurredAt: "2026-07-05T08:00:00Z",
		details:    map[string]any{"field": "title", "to": "Evidence review"},
	})

	tests := []struct {
		name    string
		args    []string
		wantIDs []string
	}{
		{
			name:    "on filter",
			args:    []string{"activity", "list", "--on", "issue/" + primaryIssueID, "--format", "json"},
			wantIDs: []string{"ACT-20260704-000099", "ACT-20260705-000100"},
		},
		{
			name:    "since filter",
			args:    []string{"activity", "list", "--since", "2026-07-05T00:00:00Z", "--format", "json"},
			wantIDs: []string{"ACT-20260705-000100", "ACT-20260705-000101"},
		},
		{
			name:    "event filter",
			args:    []string{"activity", "list", "--event", "issue_updated", "--format", "json"},
			wantIDs: []string{"ACT-20260705-000101"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, tt.args...)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			rows := decodeCommandJSONStructure(t, stdout).([]any)
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

func TestActivityListDispatchesFormats(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	issueID := createActivityCommandIssue(t, workspace, "Primary issue")
	seedActivityCommandEntry(t, workspace, activityCommandSeed{
		eventID:    "EVT-20260705-000100",
		activityID: "ACT-20260705-000100",
		refID:      issueID,
		eventName:  "status_changed",
		actor:      "developer",
		occurredAt: "2026-07-05T09:00:00Z",
		details:    map[string]any{"from": "todo", "to": "in_progress"},
	})

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "activity", "list", "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	want := expectedCommandStructure(t, []internalactivity.Entry{
		{
			ID:             "ACT-20260705-000100",
			JournalEventID: "EVT-20260705-000100",
			RefType:        "issue",
			RefID:          issueID,
			Event:          "status_changed",
			Actor:          "developer",
			OccurredAt:     "2026-07-05T09:00:00Z",
			Details:        map[string]any{"from": "todo", "to": "in_progress"},
		},
	})
	if got := decodeCommandJSONStructure(t, jsonStdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("activity list json = %#v, want %#v", got, want)
	}

	for _, formatName := range []string{"toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "activity", "list", "--format", formatName)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			var got any
			switch formatName {
			case "toon":
				got = decodeCommandTOONStructure(t, stdout)
			case "yaml":
				got = decodeCommandYAMLStructure(t, stdout)
			}

			if !commandStructuresEqual(want, got) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, got, want)
			}
		})
	}

	markdownStdout, markdownStderr := runFIDCommand(t, workspace, "activity", "list", "--format", "markdown")
	if markdownStderr != "" {
		t.Fatalf("markdown stderr = %q, want empty", markdownStderr)
	}

	wantMarkdown := "" +
		"# Activity\n\n" +
		"| Occurred At | Event | Ref | Actor | ID | Details |\n" +
		"| --- | --- | --- | --- | --- | --- |\n" +
		fmt.Sprintf("| 2026-07-05T09:00:00Z | status_changed | issue/%s | developer | ACT-20260705-000100 | {\"from\":\"todo\",\"to\":\"in_progress\"} |\n", issueID)
	if markdownStdout != wantMarkdown {
		t.Fatalf("markdown stdout = %q, want %q", markdownStdout, wantMarkdown)
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

	return created.ID
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
