package links

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

func TestListReturnsOutgoingAndIncomingLinksWithResolvedTitlesInStableOrder(t *testing.T) {
	t.Parallel()

	workspace := newLinksWorkspace(t)
	primaryIssueID := createLinksIssue(t, workspace, "Primary issue")
	outgoingIssueID := createLinksIssue(t, workspace, "Schema migration runner")
	incomingIssueID := createLinksIssue(t, workspace, "Add index on concerns(issue_id, status)")
	requirementID := syncLinksRequirement(t, workspace, "PRD-LINK-002")

	if _, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        primaryIssueID,
		Relation:      "blocks",
		ToID:          outgoingIssueID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("seed Add(blocks) error = %v", err)
	}
	if _, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        primaryIssueID,
		Relation:      "implements",
		ToID:          requirementID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("seed Add(implements) error = %v", err)
	}
	if _, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        incomingIssueID,
		Relation:      "discovered_from",
		ToID:          primaryIssueID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("seed Add(discovered_from) error = %v", err)
	}

	got, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		EntityID:      primaryIssueID,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := []ListItem{
		{
			Direction:   "outgoing",
			Relation:    "blocks",
			EntityType:  "issue",
			EntityID:    outgoingIssueID,
			EntityTitle: stringPointer("Schema migration runner"),
		},
		{
			Direction:   "outgoing",
			Relation:    "implements",
			EntityType:  "requirement",
			EntityID:    requirementID,
			EntityTitle: stringPointer("The system shall support links from issues to requirements."),
		},
		{
			Direction:   "incoming",
			Relation:    "discovered_from",
			EntityType:  "issue",
			EntityID:    incomingIssueID,
			EntityTitle: stringPointer("Add index on concerns(issue_id, status)"),
		},
	}
	if len(got) != len(want) {
		t.Fatalf("len(List()) = %d, want %d (%#v)", len(got), len(want), got)
	}
	for index := range want {
		if !equalListItem(got[index], want[index]) {
			t.Fatalf("List()[%d] = %#v, want %#v", index, got[index], want[index])
		}
	}
}

func TestListReturnsEmptySliceForEntityWithNoLinks(t *testing.T) {
	t.Parallel()

	workspace := newLinksWorkspace(t)
	issueID := createLinksIssue(t, workspace, "Unlinked issue")

	got, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		EntityID:      issueID,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("List() = %#v, want empty slice", got)
	}
}

func TestQueryLinkListHonorsAtAndOverResultLimitWithoutTruncation(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeLinksDB(t, conn)
	for _, statement := range []string{
		`CREATE TABLE links (from_type TEXT, from_id TEXT, relation TEXT, to_type TEXT, to_id TEXT, created_at TEXT)`,
		`CREATE TABLE issues (id TEXT, title TEXT)`,
		`CREATE TABLE requirements (id TEXT, summary TEXT)`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("Exec(%q) error = %v", statement, err)
		}
	}

	insert, err := conn.Prepare(`
		INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
		VALUES ('issue', 'FID-primary', 'relates_to', 'issue', ?, '2026-08-15T00:00:00Z')
	`)
	if err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}
	defer func() { _ = insert.Close() }()
	for index := 0; index < linkListLimit; index++ {
		if _, err := insert.Exec(linkListItemID(index)); err != nil {
			t.Fatalf("insert at-limit row %d: %v", index, err)
		}
	}

	atLimit, err := queryLinkList(conn, "issue", "FID-primary")
	if err != nil {
		t.Fatalf("queryLinkList(at limit) error = %v", err)
	}
	if len(atLimit) != linkListLimit {
		t.Fatalf("len(queryLinkList at limit) = %d, want %d", len(atLimit), linkListLimit)
	}

	if _, err := insert.Exec(linkListItemID(linkListLimit)); err != nil {
		t.Fatalf("insert over-limit row: %v", err)
	}
	_, err = queryLinkList(conn, "issue", "FID-primary")
	var limit *ResultLimitExceededError
	if !errors.As(err, &limit) {
		t.Fatalf("queryLinkList(over limit) error = %T %v, want *ResultLimitExceededError", err, err)
	}
	if limit.Collection != "links" || limit.Limit != linkListLimit || limit.Observed != linkListLimit+1 {
		t.Fatalf("limit error = %#v, want links/%d/%d", limit, linkListLimit, linkListLimit+1)
	}
	var rows int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM links`).Scan(&rows); err != nil {
		t.Fatalf("count links after over-limit read: %v", err)
	}
	if rows != linkListLimit+1 {
		t.Fatalf("link rows after over-limit read = %d, want %d (no truncation/mutation)", rows, linkListLimit+1)
	}
}

func linkListItemID(index int) string {
	return fmt.Sprintf("FID-link-%05d", index)
}

func stringPointer(value string) *string {
	return &value
}

func equalListItem(left, right ListItem) bool {
	return left.Direction == right.Direction &&
		left.Relation == right.Relation &&
		left.EntityType == right.EntityType &&
		left.EntityID == right.EntityID &&
		dereferenceListString(left.EntityTitle) == dereferenceListString(right.EntityTitle)
}

func dereferenceListString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
