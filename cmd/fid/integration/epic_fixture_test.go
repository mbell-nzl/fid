//go:build integration

package integration

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/epicwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

// Existing completion journeys start from an already approved empty plan.
// The new epic journey separately exercises the real terminal boundary.
func approveJourneyEpic(t *testing.T, root, epic string) {
	t.Helper()
	_, err := epicwork.Approve(epicwork.Options{WorkspaceRoot: root, MigrationsFS: fidmigrations.FS, EpicID: epic, Actor: "human:mark", Reason: "Approve fixture plan", Confirm: func(epicplan.Plan) (bool, error) { return true, nil }})
	if err != nil {
		t.Fatal(err)
	}
}
