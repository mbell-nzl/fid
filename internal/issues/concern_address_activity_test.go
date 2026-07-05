package issues

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestAddressConcernProjectsConcernAddressedActivityInSameEvent(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	concernID, issueID, addressCompletionID := seedAddressConcernScenario(t, workspace)
	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	activityBefore := countIssueActivityRows(t, conn)
	closeIssueTestDB(t, conn)

	result, err := AddressConcern(AddressConcernOptions{
		WorkspaceRoot:           workspace,
		MigrationsDir:           testIssueMigrationsDir(t),
		ConcernID:               concernID,
		ActorRole:               "developer",
		AddressedInCompletionID: addressCompletionID,
	})
	if err != nil {
		t.Fatalf("AddressConcern() error = %v", err)
	}
	if result.IssueID != issueID {
		t.Fatalf("AddressConcern().IssueID = %q, want %q", result.IssueID, issueID)
	}

	if got := countJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after AddressConcern() = %d, want %d", got, linesBefore+1)
	}

	stored := readIssueStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "issue.concern_addressed" {
		t.Fatalf("last event type = %q, want issue.concern_addressed", last.EventType)
	}
	if len(last.Effects) != 2 {
		t.Fatalf("effect count = %d, want 2", len(last.Effects))
	}

	var activityFound bool
	for index := range last.Effects {
		effect := last.Effects[index]
		if effect.Table != "activity_log" || effect.Op != "activity" {
			continue
		}
		activityFound = true
		if got := effect.Values["event"]; got != "concern_addressed" {
			t.Fatalf("activity event = %#v, want concern_addressed", got)
		}
		if got := effect.Values["ref_type"]; got != "issue" {
			t.Fatalf("activity ref_type = %#v, want issue", got)
		}
		if got := effect.Values["ref_id"]; got != issueID {
			t.Fatalf("activity ref_id = %#v, want %q", got, issueID)
		}
		details, ok := effect.Values["details"].(map[string]any)
		if !ok {
			t.Fatalf("activity details = %#v, want map", effect.Values["details"])
		}
		if got := details["concern_id"]; got != concernID {
			t.Fatalf("activity details concern_id = %#v, want %q", got, concernID)
		}
	}
	if !activityFound {
		t.Fatal("missing concern_addressed activity effect")
	}

	conn = openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	if got := countIssueActivityRows(t, conn); got != activityBefore+1 {
		t.Fatalf("activity_log rows after AddressConcern() = %d, want %d", got, activityBefore+1)
	}

	var (
		activityCount int
		eventName     string
		actor         string
		details       string
	)
	if err := conn.QueryRow(`
		SELECT COUNT(*), event, COALESCE(actor, ''), COALESCE(details, '')
		FROM activity_log
		WHERE journal_event_id = ?
		GROUP BY event, actor, details
	`, last.EventID).Scan(&activityCount, &eventName, &actor, &details); err != nil {
		t.Fatalf("query concern_addressed activity row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if eventName != "concern_addressed" {
		t.Fatalf("activity_log event = %q, want concern_addressed", eventName)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	if !strings.Contains(details, `"concern_id":"`+concernID+`"`) {
		t.Fatalf("activity_log details = %q, want concern_id %q", details, concernID)
	}
}

func TestAddressConcernRejectsWithoutActivityAppend(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	concernID, _, addressCompletionID := seedAddressConcernScenario(t, workspace)

	if _, err := AddressConcern(AddressConcernOptions{
		WorkspaceRoot:           workspace,
		MigrationsDir:           testIssueMigrationsDir(t),
		ConcernID:               concernID,
		ActorRole:               "developer",
		AddressedInCompletionID: addressCompletionID,
	}); err != nil {
		t.Fatalf("AddressConcern(seed addressed) error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	activityBefore := countIssueActivityRows(t, conn)
	closeIssueTestDB(t, conn)

	_, err := AddressConcern(AddressConcernOptions{
		WorkspaceRoot:           workspace,
		MigrationsDir:           testIssueMigrationsDir(t),
		ConcernID:               concernID,
		ActorRole:               "developer",
		AddressedInCompletionID: addressCompletionID,
	})
	if err == nil {
		t.Fatal("AddressConcern(already addressed) error = nil, want rejection")
	}

	if !strings.Contains(err.Error(), "only open concerns can be addressed") {
		t.Fatalf("AddressConcern(already addressed) error = %q, want addressed rejection", err.Error())
	}
	if got := countJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected AddressConcern() = %d, want unchanged %d", got, linesBefore)
	}

	conn = openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)
	if got := countIssueActivityRows(t, conn); got != activityBefore {
		t.Fatalf("activity_log rows after rejected AddressConcern() = %d, want unchanged %d", got, activityBefore)
	}
}
