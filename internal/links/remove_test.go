package links

import (
	"database/sql"
	"errors"
	"strings"
	"testing"
)

func TestRemoveDeletesDeclaredAndInverseRowsAndProjectsActivity(t *testing.T) {
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

	result, err := Remove(RemoveOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        firstIssueID,
		Relation:      "blocks",
		ToID:          secondIssueID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if got := countLinksJournalEvents(t, workspace); got != beforeEvents+1 {
		t.Fatalf("journal events after Remove() = %d, want %d", got, beforeEvents+1)
	}
	if result.Declared != (Link{FromType: "issue", FromID: firstIssueID, Relation: "blocks", ToType: "issue", ToID: secondIssueID}) {
		t.Fatalf("declared link = %#v", result.Declared)
	}
	if result.Inverse != (Link{FromType: "issue", FromID: secondIssueID, Relation: "blocked_by", ToType: "issue", ToID: firstIssueID}) {
		t.Fatalf("inverse link = %#v", result.Inverse)
	}

	conn := openLinksDB(t, workspace)
	defer closeLinksDB(t, conn)

	if rows := loadLinkRows(t, conn); len(rows) != 0 {
		t.Fatalf("links rows after Remove() = %#v, want none", rows)
	}

	eventName, actor, details := loadSingleActivityByRef(t, conn, "issue", firstIssueID)
	if eventName != "link_removed" {
		t.Fatalf("activity_log event = %q, want %q", eventName, "link_removed")
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want %q", actor, "developer")
	}
	for _, want := range []string{
		`"relation":"blocks"`,
		`"to_id":"` + secondIssueID + `"`,
		`"inverse_relation":"blocked_by"`,
	} {
		if !strings.Contains(details, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details, want)
		}
	}
}

func TestRemoveRejectsMissingLinkBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newLinksWorkspace(t)
	firstIssueID := createLinksIssue(t, workspace, "First issue")
	secondIssueID := createLinksIssue(t, workspace, "Second issue")
	beforeEvents := countLinksJournalEvents(t, workspace)

	_, err := Remove(RemoveOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        firstIssueID,
		Relation:      "blocks",
		ToID:          secondIssueID,
		ActorRole:     "developer",
	})
	if err == nil {
		t.Fatal("Remove() error = nil, want rejection")
	}

	var target *MissingLinkError
	if !errors.As(err, &target) {
		t.Fatalf("error = %T, want *MissingLinkError", err)
	}
	if got := countLinksJournalEvents(t, workspace); got != beforeEvents {
		t.Fatalf("journal events after rejection = %d, want %d", got, beforeEvents)
	}
}

func TestRemoveSupportsRequirementInverse(t *testing.T) {
	t.Parallel()

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
		t.Fatalf("seed Add(implements) error = %v", err)
	}

	result, err := Remove(RemoveOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        issueID,
		Relation:      "implements",
		ToID:          requirementID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Remove(implements) error = %v", err)
	}
	if result.Inverse.Relation != "implemented_by" {
		t.Fatalf("inverse relation = %q, want %q", result.Inverse.Relation, "implemented_by")
	}

	conn := openLinksDB(t, workspace)
	defer closeLinksDB(t, conn)

	if rows := loadLinkRows(t, conn); len(rows) != 0 {
		t.Fatalf("links rows after Remove(implements) = %#v, want none", rows)
	}
}

func loadSingleActivityByRef(t *testing.T, conn *sql.DB, refType, refID string) (string, string, string) {
	t.Helper()

	var (
		event   string
		actor   string
		details string
	)
	if err := conn.QueryRow(`
		SELECT event, COALESCE(actor, ''), COALESCE(details, '')
		FROM activity_log
		WHERE ref_type = ? AND ref_id = ?
		ORDER BY occurred_at DESC, id DESC
		LIMIT 1
	`, refType, refID).Scan(&event, &actor, &details); err != nil {
		t.Fatalf("query activity_log row: %v", err)
	}
	return event, actor, details
}
