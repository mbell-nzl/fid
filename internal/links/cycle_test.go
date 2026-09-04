package links

import (
	"errors"
	"testing"
)

func TestAddRejectsThreeNodeDependencyCycle(t *testing.T) {
	workspace := newLinksWorkspace(t)
	firstID := createLinksIssue(t, workspace, "First issue")
	secondID := createLinksIssue(t, workspace, "Second issue")
	thirdID := createLinksIssue(t, workspace, "Third issue")

	for _, edge := range [][2]string{{firstID, secondID}, {secondID, thirdID}} {
		_, err := Add(AddOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testLinksMigrationsDir(t),
			FromID:        edge[0],
			Relation:      "blocks",
			ToID:          edge[1],
			Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		})
		if err != nil {
			t.Fatalf("add %q blocks %q: %v", edge[0], edge[1], err)
		}
	}

	beforeEvents := countLinksJournalEvents(t, workspace)
	_, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromID:        thirdID,
		Relation:      "blocks",
		ToID:          firstID,
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	var cycle *DependencyCycleError
	if !errors.As(err, &cycle) {
		t.Fatalf("error = %T %v, want *DependencyCycleError", err, err)
	}
	if want := []string{thirdID, firstID, secondID, thirdID}; !sameStrings(cycle.Path, want) {
		t.Fatalf("cycle path = %q, want %q", cycle.Path, want)
	}
	if got := countLinksJournalEvents(t, workspace); got != beforeEvents {
		t.Fatalf("journal events after cycle rejection = %d, want %d", got, beforeEvents)
	}
}

func TestAddRejectsSelfDependencyCycle(t *testing.T) {
	workspace := newLinksWorkspace(t)
	issueID := createLinksIssue(t, workspace, "Self dependency")
	beforeEvents := countLinksJournalEvents(t, workspace)

	_, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromType:      "issue",
		FromID:        issueID,
		Relation:      "blocks",
		ToType:        "issue",
		ToID:          issueID,
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	var cycle *DependencyCycleError
	if !errors.As(err, &cycle) {
		t.Fatalf("error = %T %v, want *DependencyCycleError", err, err)
	}
	if want := []string{issueID, issueID}; !sameStrings(cycle.Path, want) {
		t.Fatalf("cycle path = %q, want %q", cycle.Path, want)
	}
	if got := countLinksJournalEvents(t, workspace); got != beforeEvents {
		t.Fatalf("journal events after self-cycle rejection = %d, want %d", got, beforeEvents)
	}
}

func TestDependencyCycleTraversalMatrix(t *testing.T) {
	t.Run("two-node cycle has ordered path and no append", func(t *testing.T) {
		workspace := newLinksWorkspace(t)
		firstID := createLinksIssue(t, workspace, "First issue")
		secondID := createLinksIssue(t, workspace, "Second issue")
		addDependency(t, workspace, firstID, "blocks", secondID)
		assertDependencyCycle(t, workspace, secondID, "blocks", firstID, []string{secondID, firstID, secondID})
	})

	t.Run("three-node cycle has ordered path and no append", func(t *testing.T) {
		workspace := newLinksWorkspace(t)
		firstID := createLinksIssue(t, workspace, "First issue")
		secondID := createLinksIssue(t, workspace, "Second issue")
		thirdID := createLinksIssue(t, workspace, "Third issue")
		addDependency(t, workspace, firstID, "blocks", secondID)
		addDependency(t, workspace, secondID, "blocks", thirdID)
		assertDependencyCycle(t, workspace, thirdID, "blocks", firstID, []string{thirdID, firstID, secondID, thirdID})
	})

	t.Run("already-cyclic storage terminates and rejects", func(t *testing.T) {
		workspace := newLinksWorkspace(t)
		firstID := createLinksIssue(t, workspace, "First issue")
		secondID := createLinksIssue(t, workspace, "Second issue")
		thirdID := createLinksIssue(t, workspace, "Third issue")
		insertStoredDependency(t, workspace, firstID, secondID)
		insertStoredDependency(t, workspace, secondID, thirdID)
		insertStoredDependency(t, workspace, thirdID, firstID)
		assertDependencyCycle(t, workspace, secondID, "blocked_by", firstID, []string{firstID, secondID, thirdID, firstID})
	})

	t.Run("acyclic diamond succeeds", func(t *testing.T) {
		workspace := newLinksWorkspace(t)
		firstID := createLinksIssue(t, workspace, "First issue")
		secondID := createLinksIssue(t, workspace, "Second issue")
		thirdID := createLinksIssue(t, workspace, "Third issue")
		fourthID := createLinksIssue(t, workspace, "Fourth issue")
		addDependency(t, workspace, firstID, "blocks", secondID)
		addDependency(t, workspace, firstID, "blocks", thirdID)
		addDependency(t, workspace, secondID, "blocks", fourthID)
		addDependency(t, workspace, thirdID, "blocks", fourthID)

		before := countLinksJournalEvents(t, workspace)
		addDependency(t, workspace, secondID, "blocks", thirdID)
		if got := countLinksJournalEvents(t, workspace); got != before+1 {
			t.Fatalf("journal events after acyclic diamond add = %d, want %d", got, before+1)
		}
	})
}

func addDependency(t *testing.T, workspace, fromID, relation, toID string) {
	t.Helper()
	if _, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromType:      "issue",
		FromID:        fromID,
		Relation:      relation,
		ToType:        "issue",
		ToID:          toID,
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	}); err != nil {
		t.Fatalf("Add(%q %s %q) error = %v", fromID, relation, toID, err)
	}
}

func assertDependencyCycle(t *testing.T, workspace, fromID, relation, toID string, wantPath []string) {
	t.Helper()
	before := countLinksJournalEvents(t, workspace)
	_, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		FromType:      "issue",
		FromID:        fromID,
		Relation:      relation,
		ToType:        "issue",
		ToID:          toID,
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	var cycle *DependencyCycleError
	if !errors.As(err, &cycle) {
		t.Fatalf("Add(%q %s %q) error = %T %v, want *DependencyCycleError", fromID, relation, toID, err, err)
	}
	if !sameStrings(cycle.Path, wantPath) {
		t.Fatalf("cycle path = %q, want %q", cycle.Path, wantPath)
	}
	if got := countLinksJournalEvents(t, workspace); got != before {
		t.Fatalf("journal events after cycle rejection = %d, want %d", got, before)
	}
}

func insertStoredDependency(t *testing.T, workspace, fromID, toID string) {
	t.Helper()
	conn := openLinksDB(t, workspace)
	defer closeLinksDB(t, conn)
	if _, err := conn.Exec(`
		INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
		VALUES ('issue', ?, 'blocks', 'issue', ?, '2026-08-15T00:00:00Z')
	`, fromID, toID); err != nil {
		t.Fatalf("insert pre-existing %q blocks %q: %v", fromID, toID, err)
	}
}

func sameStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for index := range got {
		if got[index] != want[index] {
			return false
		}
	}
	return true
}
