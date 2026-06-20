package issues

import "testing"

func TestListConcernsFiltersByIssueAndStatus(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	firstIssueID, firstCompletionID := seedConcernRaiseIssueAndCompletion(t, workspace)

	firstConcern, err := RaiseConcern(RaiseConcernOptions{
		WorkspaceRoot:        workspace,
		MigrationsDir:        testIssueMigrationsDir(t),
		IssueID:              firstIssueID,
		ActorRole:            "developer",
		RaisedInCompletionID: firstCompletionID,
		Summary:              "First open concern",
	})
	if err != nil {
		t.Fatalf("RaiseConcern(first) error = %v", err)
	}

	addressCompletion, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       firstIssueID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Complete(address) error = %v", err)
	}

	if _, err := AddressConcern(AddressConcernOptions{
		WorkspaceRoot:           workspace,
		MigrationsDir:           testIssueMigrationsDir(t),
		ConcernID:               firstConcern.ID,
		ActorRole:               "developer",
		AddressedInCompletionID: addressCompletion.ID,
	}); err != nil {
		t.Fatalf("AddressConcern(first) error = %v", err)
	}

	secondIssueID, secondCompletionID := seedConcernRaiseIssueAndCompletion(t, workspace)
	secondConcern, err := RaiseConcern(RaiseConcernOptions{
		WorkspaceRoot:        workspace,
		MigrationsDir:        testIssueMigrationsDir(t),
		IssueID:              secondIssueID,
		ActorRole:            "developer",
		RaisedInCompletionID: secondCompletionID,
		Summary:              "Second open concern",
		BodyMDPath:           stringPointer("project-docs/concerns/CON-second.md"),
	})
	if err != nil {
		t.Fatalf("RaiseConcern(second) error = %v", err)
	}

	all, err := ListConcerns(ConcernListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("ListConcerns(all) error = %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("ListConcerns(all) len = %d, want 2", len(all))
	}
	secondConcernItem := concernListItemByID(t, all, secondConcern.ID)
	if secondConcernItem.IssueID != secondIssueID {
		t.Fatalf("ListConcerns(second).IssueID = %q, want %q", secondConcernItem.IssueID, secondIssueID)
	}
	if secondConcernItem.Status != "open" {
		t.Fatalf("ListConcerns(second).Status = %q, want open", secondConcernItem.Status)
	}
	if secondConcernItem.BodyMDPath == nil || *secondConcernItem.BodyMDPath != "project-docs/concerns/CON-second.md" {
		t.Fatalf("ListConcerns(second).BodyMDPath = %#v, want body path", secondConcernItem.BodyMDPath)
	}
	firstConcernItem := concernListItemByID(t, all, firstConcern.ID)
	if firstConcernItem.IssueID != firstIssueID {
		t.Fatalf("ListConcerns(first).IssueID = %q, want %q", firstConcernItem.IssueID, firstIssueID)
	}
	if firstConcernItem.Status != "addressed" {
		t.Fatalf("ListConcerns(first).Status = %q, want addressed", firstConcernItem.Status)
	}

	byIssue, err := ListConcerns(ConcernListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       firstIssueID,
	})
	if err != nil {
		t.Fatalf("ListConcerns(issue) error = %v", err)
	}
	if len(byIssue) != 1 || byIssue[0].ID != firstConcern.ID {
		t.Fatalf("ListConcerns(issue) = %#v, want only %q", byIssue, firstConcern.ID)
	}

	byStatus, err := ListConcerns(ConcernListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Status:        "open",
	})
	if err != nil {
		t.Fatalf("ListConcerns(status) error = %v", err)
	}
	if len(byStatus) != 1 || byStatus[0].ID != secondConcern.ID {
		t.Fatalf("ListConcerns(status) = %#v, want only %q", byStatus, secondConcern.ID)
	}

	combined, err := ListConcerns(ConcernListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       secondIssueID,
		Status:        "open",
	})
	if err != nil {
		t.Fatalf("ListConcerns(combined) error = %v", err)
	}
	if len(combined) != 1 || combined[0].ID != secondConcern.ID {
		t.Fatalf("ListConcerns(combined) = %#v, want only %q", combined, secondConcern.ID)
	}

	empty, err := ListConcerns(ConcernListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       secondIssueID,
		Status:        "addressed",
	})
	if err != nil {
		t.Fatalf("ListConcerns(empty) error = %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("ListConcerns(empty) len = %d, want 0", len(empty))
	}
}

func concernListItemByID(t *testing.T, items []ConcernListItem, concernID string) ConcernListItem {
	t.Helper()

	for _, item := range items {
		if item.ID == concernID {
			return item
		}
	}

	t.Fatalf("concern %q not found in %#v", concernID, items)
	return ConcernListItem{}
}
