package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestJournalInspectCommandUsesDirectLegacyJSONWithoutRevert(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Inspect journal event",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)
	_, transitionStderr := runFIDCommand(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "in_progress",
		"--actor", "agent:developer",
	)
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": todo → in_progress") {
		t.Fatalf("transition stderr = %q, want confirmation", transitionStderr)
	}

	journalPath := mustJournalPath(t, workspace)
	event := mustReadLastJournalEvent(t, journalPath)
	journalFile := filepath.ToSlash(filepath.Join(config.LedgerDir(gittest.ResolvedPath(t, workspace)), config.JournalDirName, filepath.Base(journalPath)))
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
	})

	stdout, stderr := runFIDCommand(t, workspace, "journal", "inspect", event.EventID)
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if !strings.HasSuffix(stdout, "\n") {
		t.Fatalf("stdout = %q, want one trailing newline", stdout)
	}
	if got := decodeCommandJSONStructure(t, stdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("JSON structure = %#v, want %#v", got, want)
	}
}

func TestJournalInspectCommandDoesNotExposeRevertOrResponseEnvelope(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Inspect non-eligible event",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	event := mustReadCommandEventByType(t, mustJournalPath(t, workspace), "issue.created")
	stdout, stderr := runFIDCommand(t, workspace, "journal", "inspect", event.EventID)
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	payload := decodeCommandJSONStructure(t, stdout).(map[string]any)
	if got := payload["event_type"]; got != "issue.created" {
		t.Fatalf("event_type = %#v, want %q", got, "issue.created")
	}

	for _, forbidden := range []string{"revert", "contract_version", "event"} {
		if _, ok := payload[forbidden]; ok {
			t.Fatalf("journal inspect JSON includes %q: %#v", forbidden, payload)
		}
	}

	if issueID == "" {
		t.Fatal("issue id = empty, want created issue id")
	}
}

func TestJournalInspectCommandRejectsMissingEventIDs(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init")

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "journal", "inspect", "EVT-20260606-999999")
	if err == nil {
		t.Fatal("Execute() error = nil, want missing event rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}

	journalDir := filepath.ToSlash(filepath.Join(config.LedgerDir(gittest.ResolvedPath(t, workspace)), config.JournalDirName))
	want := `error: journal event "EVT-20260606-999999" not found` + "\n" +
		"  check the event id and confirm it exists under " + journalDir + "/\n"
	if stderr != want {
		t.Fatalf("stderr = %q, want %q", stderr, want)
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
