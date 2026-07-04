package issues

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRaiseConcernProjectsConcernRaisedActivityInSameEvent(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	issueID, completionID := seedConcernRaiseIssueAndCompletion(t, workspace)
	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	activityBefore := countIssueActivityRows(t, conn)
	closeIssueTestDB(t, conn)

	result, err := RaiseConcern(RaiseConcernOptions{
		WorkspaceRoot:        workspace,
		MigrationsDir:        testIssueMigrationsDir(t),
		IssueID:              issueID,
		ActorRole:            "developer",
		RaisedInCompletionID: completionID,
		Summary:              "Potential race in rebuild parity",
	})
	if err != nil {
		t.Fatalf("RaiseConcern() error = %v", err)
	}

	if got := countJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after RaiseConcern() = %d, want %d", got, linesBefore+1)
	}

	stored := readIssueStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "issue.concern_raised" {
		t.Fatalf("last event type = %q, want issue.concern_raised", last.EventType)
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
		if got := effect.Values["event"]; got != "concern_raised" {
			t.Fatalf("activity event = %#v, want concern_raised", got)
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
		if got := details["concern_id"]; got != result.ID {
			t.Fatalf("activity details concern_id = %#v, want %q", got, result.ID)
		}
		if got := details["summary"]; got != "Potential race in rebuild parity" {
			t.Fatalf("activity details summary = %#v, want concern summary", got)
		}
	}
	if !activityFound {
		t.Fatal("missing concern_raised activity effect")
	}

	conn = openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	if got := countIssueActivityRows(t, conn); got != activityBefore+1 {
		t.Fatalf("activity_log rows after RaiseConcern() = %d, want %d", got, activityBefore+1)
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
		t.Fatalf("query concern_raised activity row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if eventName != "concern_raised" {
		t.Fatalf("activity_log event = %q, want concern_raised", eventName)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	for _, want := range []string{
		`"concern_id":"` + result.ID + `"`,
		`"summary":"Potential race in rebuild parity"`,
	} {
		if !strings.Contains(details, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details, want)
		}
	}
}
