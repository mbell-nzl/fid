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
			ActorRole:     "developer",
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
		ActorRole:     "developer",
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
