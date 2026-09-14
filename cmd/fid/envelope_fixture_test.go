package main

import (
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/project"
)

func TestObservedEnvelopeFixtureClonesCaptureDistinctRepositories(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	for range 2 {
		root, _, registration := finalizedObservedEnvelopeWorkspaceWithDefinitions(t, gateDefinitionsForTest(nil))
		observation, err := project.ObserveRepository(root)
		if err != nil {
			t.Fatal(err)
		}
		if seen[observation.WorktreeID] {
			t.Fatal("distinct fixture roots reused a Worktree identity")
		}
		seen[observation.WorktreeID] = true
		for _, rebuild := range []bool{false, true} {
			if rebuild {
				result := executeCommandForTest(root, "--project", root, "--json", "rebuild")
				if result.err != nil {
					t.Fatalf("clone history rebuild: %#v", result)
				}
			}
			db, err := fiddb.Open(filepath.Join(root, config.WorkspaceDirName, config.DatabaseFileName))
			if err != nil {
				t.Fatal(err)
			}
			var worktree string
			err = db.QueryRow("SELECT worktree_id FROM gate_registrations WHERE id = ? AND version = ?", registration.ID, registration.Version).Scan(&worktree)
			closeErr := db.Close()
			if err != nil || closeErr != nil || worktree != observation.WorktreeID {
				t.Fatalf("rebuild=%t registration worktree=%s want=%s errors=%v/%v", rebuild, worktree, observation.WorktreeID, err, closeErr)
			}
		}
	}
}
