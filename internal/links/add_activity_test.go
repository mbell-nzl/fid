package links

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestAddProjectsLinkAddedActivityForIssueAndRequirementLinks(t *testing.T) {
	t.Parallel()

	t.Run("issue to issue", func(t *testing.T) {
		workspace := newLinksWorkspace(t)
		firstIssueID := createLinksIssue(t, workspace, "First issue")
		secondIssueID := createLinksIssue(t, workspace, "Second issue")
		beforeEvents := countLinksJournalEvents(t, workspace)

		conn := openLinksDB(t, workspace)
		activityBefore := countLinkActivityRows(t, conn)
		closeLinksDB(t, conn)

		result, err := Add(AddOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testLinksMigrationsDir(t),
			FromID:        firstIssueID,
			Relation:      "blocks",
			ToID:          secondIssueID,
			ActorRole:     "developer",
		})
		if err != nil {
			t.Fatalf("Add(blocks) error = %v", err)
		}

		if got := countLinksJournalEvents(t, workspace); got != beforeEvents+1 {
			t.Fatalf("journal events after Add(blocks) = %d, want %d", got, beforeEvents+1)
		}

		stored := readLinkStoredEvents(t, filepath.Join(workspace, ".fid", "journal"))
		last := stored[len(stored)-1]
		if last.EventType != "link.added" {
			t.Fatalf("last event type = %q, want link.added", last.EventType)
		}
		if len(last.Effects) != 3 {
			t.Fatalf("effect count = %d, want 3", len(last.Effects))
		}

		var activityFound bool
		for index := range last.Effects {
			effect := last.Effects[index]
			if effect.Table != "activity_log" || effect.Op != "activity" {
				continue
			}
			activityFound = true
			if got := effect.Values["event"]; got != "link_added" {
				t.Fatalf("activity event = %#v, want link_added", got)
			}
			if got := effect.Values["ref_type"]; got != "issue" {
				t.Fatalf("activity ref_type = %#v, want issue", got)
			}
			if got := effect.Values["ref_id"]; got != firstIssueID {
				t.Fatalf("activity ref_id = %#v, want %q", got, firstIssueID)
			}
		}
		if !activityFound {
			t.Fatal("missing link_added activity effect")
		}

		conn = openLinksDB(t, workspace)
		defer closeLinksDB(t, conn)

		if got := countLinkActivityRows(t, conn); got != activityBefore+1 {
			t.Fatalf("activity_log rows after Add(blocks) = %d, want %d", got, activityBefore+1)
		}

		eventName, actor, details := loadSingleActivityByEvent(t, conn, last.EventID)
		if eventName != "link_added" {
			t.Fatalf("activity_log event = %q, want link_added", eventName)
		}
		if actor != "developer" {
			t.Fatalf("activity_log actor = %q, want developer", actor)
		}
		for _, want := range []string{
			`"from_type":"issue"`,
			`"from_id":"` + firstIssueID + `"`,
			`"relation":"blocks"`,
			`"to_type":"issue"`,
			`"to_id":"` + secondIssueID + `"`,
			`"inverse_relation":"blocked_by"`,
			`"inverse_from_id":"` + secondIssueID + `"`,
			`"inverse_to_id":"` + firstIssueID + `"`,
		} {
			if !strings.Contains(details, want) {
				t.Fatalf("activity_log details = %q, want substring %q", details, want)
			}
		}
		if result.Inverse.Relation != "blocked_by" {
			t.Fatalf("inverse relation = %q, want blocked_by", result.Inverse.Relation)
		}
	})

	t.Run("issue implements requirement", func(t *testing.T) {
		workspace := newLinksWorkspace(t)
		issueID := createLinksIssue(t, workspace, "Implements issue")
		requirementID := syncLinksRequirement(t, workspace, "PRD-LINK-002")

		if _, err := Add(AddOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testLinksMigrationsDir(t),
			FromID:        issueID,
			Relation:      "implements",
			ToID:          requirementID,
			ActorRole:     "developer",
		}); err != nil {
			t.Fatalf("Add(implements) error = %v", err)
		}

		stored := readLinkStoredEvents(t, filepath.Join(workspace, ".fid", "journal"))
		last := stored[len(stored)-1]

		conn := openLinksDB(t, workspace)
		defer closeLinksDB(t, conn)

		eventName, actor, details := loadSingleActivityByEvent(t, conn, last.EventID)
		if eventName != "link_added" {
			t.Fatalf("activity_log event = %q, want link_added", eventName)
		}
		if actor != "developer" {
			t.Fatalf("activity_log actor = %q, want developer", actor)
		}
		for _, want := range []string{
			`"relation":"implements"`,
			`"to_type":"requirement"`,
			`"to_id":"` + requirementID + `"`,
			`"inverse_relation":"implemented_by"`,
		} {
			if !strings.Contains(details, want) {
				t.Fatalf("activity_log details = %q, want substring %q", details, want)
			}
		}
	})
}

func TestAddRejectsDuplicateWithoutActivityAppend(t *testing.T) {
	t.Parallel()

	workspace := newLinksWorkspace(t)
	firstIssueID := createLinksIssue(t, workspace, "First issue")
	secondIssueID := createLinksIssue(t, workspace, "Second issue")
	if _, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        firstIssueID,
		Relation:      "blocks",
		ToID:          secondIssueID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("seed Add(blocks) error = %v", err)
	}

	beforeEvents := countLinksJournalEvents(t, workspace)
	conn := openLinksDB(t, workspace)
	activityBefore := countLinkActivityRows(t, conn)
	closeLinksDB(t, conn)

	_, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        firstIssueID,
		Relation:      "blocks",
		ToID:          secondIssueID,
		ActorRole:     "developer",
	})
	if err == nil {
		t.Fatal("Add(duplicate) error = nil, want rejection")
	}

	if got := countLinksJournalEvents(t, workspace); got != beforeEvents {
		t.Fatalf("journal events after duplicate Add() = %d, want unchanged %d", got, beforeEvents)
	}

	conn = openLinksDB(t, workspace)
	defer closeLinksDB(t, conn)
	if got := countLinkActivityRows(t, conn); got != activityBefore {
		t.Fatalf("activity_log rows after duplicate Add() = %d, want unchanged %d", got, activityBefore)
	}
}

func countLinkActivityRows(t *testing.T, conn *sql.DB) int {
	t.Helper()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&count); err != nil {
		t.Fatalf("count activity_log rows: %v", err)
	}
	return count
}

func loadSingleActivityByEvent(t *testing.T, conn *sql.DB, eventID string) (string, string, string) {
	t.Helper()

	var (
		event   string
		actor   string
		details string
	)
	if err := conn.QueryRow(`
		SELECT event, COALESCE(actor, ''), COALESCE(details, '')
		FROM activity_log
		WHERE journal_event_id = ?
	`, eventID).Scan(&event, &actor, &details); err != nil {
		t.Fatalf("query activity_log row for %s: %v", eventID, err)
	}
	return event, actor, details
}

func readLinkStoredEvents(t *testing.T, journalDir string) []journal.Event {
	t.Helper()

	entries, err := os.ReadDir(journalDir)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", journalDir, err)
	}

	lines := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		payload, err := os.ReadFile(filepath.Join(journalDir, entry.Name()))
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", entry.Name(), err)
		}
		trimmed := strings.TrimSpace(string(payload))
		if trimmed == "" {
			continue
		}
		lines = append(lines, strings.Split(trimmed, "\n")...)
	}

	events := make([]journal.Event, 0, len(lines))
	for _, line := range lines {
		var event journal.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		events = append(events, event)
	}
	return events
}
