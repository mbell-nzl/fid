package links

import "testing"

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
