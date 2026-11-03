package links

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
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

	stored := readLinkStoredEvents(t, filepath.Join(config.LedgerDir(workspace), "journal"))
	last := stored[len(stored)-1]
	eventName, actor, details := loadSingleActivityByEvent(t, conn, last.EventID)
	if eventName != "link-removed" {
		t.Fatalf("activity_log event = %q, want %q", eventName, "link-removed")
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
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

func TestRemoveRejectsRemovedRequirementInverse(t *testing.T) {
	t.Parallel()

	workspace := newLinksWorkspace(t)
	issueID := createLinksIssue(t, workspace, "Implements issue")
	before := countLinksJournalEvents(t, workspace)
	_, err := Remove(RemoveOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        issueID,
		Relation:      "implements",
		ToType:        "requirement",
		ToID:          "PRD-LINK-002",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	var invalid *InvalidEndpointTypeError
	if !errors.As(err, &invalid) || countLinksJournalEvents(t, workspace) != before {
		t.Fatalf("removed requirement inverse = %v; Event count %d -> %d", err, before, countLinksJournalEvents(t, workspace))
	}

	conn := openLinksDB(t, workspace)
	defer closeLinksDB(t, conn)

	if rows := loadLinkRows(t, conn); len(rows) != 0 {
		t.Fatalf("links rows after Remove(implements) = %#v, want none", rows)
	}
}
