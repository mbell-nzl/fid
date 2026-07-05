package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestJournalInspectCommandReportsKnownEventAcrossFormats(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Inspect journal event",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)
	claimIssueForTransitionCLI(t, workspace, issueID, "developer")
	_, transitionStderr := runFIDCommand(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "in_progress",
		"--by", "developer",
	)
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": todo → in_progress") {
		t.Fatalf("transition stderr = %q, want confirmation", transitionStderr)
	}

	journalPath := mustJournalPath(t, workspace)
	event := mustReadLastJournalEvent(t, journalPath)
	journalFile := ".fid/journal/" + filepath.Base(journalPath)
	want := expectedCommandStructure(t, journalInspectOutput{
		EventID:     event.EventID,
		EventSchema: event.EventSchema,
		EventType:   event.EventType,
		JournalFile: journalFile,
		Line:        len(splitJournalLines(mustReadCommandFile(t, journalPath))),
		Hash:        mustCommandEventHash(t, event),
		Actor:       event.Actor,
		Command:     event.Command,
		Effects:     event.Effects,
		Revert:      journal.RevertVerdict(event),
	})

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "journal", "inspect", event.EventID, "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}
	if got := decodeCommandJSONStructure(t, jsonStdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("json structure = %#v, want %#v", got, want)
	}

	for _, formatName := range []string{"toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "journal", "inspect", event.EventID, "--format", formatName)
			if stderr != "" {
				t.Fatalf("%s stderr = %q, want empty", formatName, stderr)
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

	markdownStdout, markdownStderr := runFIDCommand(t, workspace, "journal", "inspect", event.EventID, "--format", "markdown")
	if markdownStderr != "" {
		t.Fatalf("markdown stderr = %q, want empty", markdownStderr)
	}
	for _, wantSubstring := range []string{
		"# Journal Event " + event.EventID,
		event.EventType,
		journalFile + ":",
		"## Revert",
		"reverts this transition to the prior status",
	} {
		if !strings.Contains(markdownStdout, wantSubstring) {
			t.Fatalf("markdown stdout = %q, want substring %q", markdownStdout, wantSubstring)
		}
	}
}

func TestJournalInspectCommandReportsCorrectionForNonEligibleEvents(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Inspect non-eligible event",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	event := mustReadCommandEventByType(t, mustJournalPath(t, workspace), "issue.created")
	stdout, stderr := runFIDCommand(t, workspace, "journal", "inspect", event.EventID, "--format", "json")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	payload := decodeCommandJSONStructure(t, stdout).(map[string]any)
	if got := payload["event_type"]; got != "issue.created" {
		t.Fatalf("event_type = %#v, want %q", got, "issue.created")
	}

	revertPayload, ok := payload["revert"].(map[string]any)
	if !ok {
		t.Fatalf("revert = %#v, want object", payload["revert"])
	}
	if got := revertPayload["eligible"]; got != false {
		t.Fatalf("revert.eligible = %#v, want false", got)
	}
	if got := revertPayload["correction"]; got != "cancel the issue with fid issue transition <id> --to cancelled --by <role>" {
		t.Fatalf("revert.correction = %#v, want cancellation guidance", got)
	}
	if got := revertPayload["note"]; !strings.Contains(got.(string), "append-only history") {
		t.Fatalf("revert.note = %#v, want append-only explanation", got)
	}

	if issueID == "" {
		t.Fatal("issue id = empty, want created issue id")
	}
}

func TestJournalInspectCommandRejectsMissingEventIDs(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "journal", "inspect", "EVT-20260606-999999")
	if err == nil {
		t.Fatal("Execute() error = nil, want missing event rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}

	for _, wantSubstring := range []string{
		`error: journal event "EVT-20260606-999999" not found`,
		"check the event id",
	} {
		if !strings.Contains(stderr, wantSubstring) {
			t.Fatalf("stderr = %q, want substring %q", stderr, wantSubstring)
		}
	}
}

func mustReadCommandEventByType(t *testing.T, journalPath, eventType string) journal.Event {
	t.Helper()

	lines := splitJournalLines(mustReadCommandFile(t, journalPath))
	for _, line := range lines {
		var event journal.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("Unmarshal(journal event) error = %v", err)
		}
		if event.EventType == eventType {
			return event
		}
	}

	t.Fatalf("journal %q has no event type %q", journalPath, eventType)
	return journal.Event{}
}

func mustCommandEventHash(t *testing.T, event journal.Event) string {
	t.Helper()

	hash, err := journal.HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent(%s) error = %v", event.EventID, err)
	}
	return hash
}
