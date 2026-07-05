package claims

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseProjectsClaimReleasedActivityInSameEvent(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Release claim target")

	if _, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Acquire(initial) error = %v", err)
	}

	journalPath := claimsJournalPathForWorkspace(t, workspace)
	linesBefore := countClaimsJournalLines(t, journalPath)

	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	activityBefore := countClaimActivityRows(t, conn)
	closeClaimsDB(t, conn)

	result, err := Release(ReleaseOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Release() error = %v", err)
	}
	if result.RefType != "issue" {
		t.Fatalf("Release().RefType = %q, want issue", result.RefType)
	}
	if result.RefID != created.ID {
		t.Fatalf("Release().RefID = %q, want %q", result.RefID, created.ID)
	}
	if result.ReleasedBy != "developer" {
		t.Fatalf("Release().ReleasedBy = %q, want developer", result.ReleasedBy)
	}

	if got := countClaimsJournalLines(t, journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after release = %d, want %d", got, linesBefore+1)
	}

	stored := readClaimStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "issue.claim.released" {
		t.Fatalf("last event type = %q, want issue.claim.released", last.EventType)
	}

	var (
		deleteEffects   int
		activityEffects int
	)
	for index := range last.Effects {
		effect := last.Effects[index]
		switch {
		case effect.Table == "claims" && effect.Op == "delete":
			deleteEffects++
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffects++
			if got := effect.Values["event"]; got != "claim_released" {
				t.Fatalf("activity event = %#v, want claim_released", got)
			}
			if got := effect.Values["ref_type"]; got != "issue" {
				t.Fatalf("activity ref_type = %#v, want issue", got)
			}
			if got := effect.Values["ref_id"]; got != created.ID {
				t.Fatalf("activity ref_id = %#v, want %q", got, created.ID)
			}
			if got := effect.Values["actor"]; got != "developer" {
				t.Fatalf("activity actor = %#v, want developer", got)
			}
			details, ok := effect.Values["details"].(map[string]any)
			if !ok {
				t.Fatalf("activity details = %#v, want map", effect.Values["details"])
			}
			if got := details["released_by"]; got != "developer" {
				t.Fatalf("activity details released_by = %#v, want developer", got)
			}
		default:
			t.Fatalf("unexpected effect = %#v", effect)
		}
	}
	if deleteEffects != 1 || activityEffects != 1 {
		t.Fatalf("effect counts = %d deletes / %d activity, want 1 / 1", deleteEffects, activityEffects)
	}

	conn = openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)

	if got := countClaimActivityRows(t, conn); got != activityBefore+1 {
		t.Fatalf("activity_log rows after release = %d, want %d", got, activityBefore+1)
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
		t.Fatalf("query release activity row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if eventName != "claim_released" {
		t.Fatalf("activity_log event = %q, want claim_released", eventName)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	if !strings.Contains(details, `"released_by":"developer"`) {
		t.Fatalf("activity_log details = %q, want released_by developer", details)
	}
}

func TestReleaseRejectsNonHolderWithoutActivityAppend(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Release wrong holder target")

	if _, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Acquire(initial) error = %v", err)
	}

	journalPath := claimsJournalPathForWorkspace(t, workspace)
	linesBefore := countClaimsJournalLines(t, journalPath)

	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	activityBefore := countClaimActivityRows(t, conn)
	closeClaimsDB(t, conn)

	_, err := Release(ReleaseOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "qa_reviewer",
	})
	if err == nil {
		t.Fatal("Release(non-holder) error = nil, want rejection")
	}

	var notHeldErr *ClaimNotHeldError
	if !errors.As(err, &notHeldErr) {
		t.Fatalf("Release(non-holder) error = %T, want *ClaimNotHeldError", err)
	}
	if got := countClaimsJournalLines(t, journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected release = %d, want unchanged %d", got, linesBefore)
	}

	conn = openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)
	if got := countClaimActivityRows(t, conn); got != activityBefore {
		t.Fatalf("activity_log rows after rejected release = %d, want unchanged %d", got, activityBefore)
	}
}

func countClaimActivityRows(t *testing.T, conn *sql.DB) int {
	t.Helper()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&count); err != nil {
		t.Fatalf("count activity_log rows: %v", err)
	}
	return count
}
