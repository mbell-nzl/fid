package donework

import (
	"errors"
	"os"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestEpicCloseRechecksNewChildUnderLock(t *testing.T) {
	t.Parallel()
	root := gittest.NewWorkspace(t)
	migrations := fidmigrations.FS
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: root, MigrationsFS: migrations}); err != nil {
		t.Fatal(err)
	}
	epic, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsFS: migrations, Kind: "epic", Title: "Stale eligibility", Scheduling: "prioritised", Actor: "human:mark", IdempotencyKey: "epic"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = closeEpic(Options{WorkspaceRoot: root, MigrationsFS: migrations, TaskID: epic.ID, Actor: "human:mark", Reason: "Finished", IdempotencyKey: "close"}, func() {
		_, createErr := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsFS: migrations, Kind: "task", Title: "Arrived before commit", ParentID: &epic.ID, Actor: "human:mark", IdempotencyKey: "child"})
		if createErr != nil {
			t.Fatal(createErr)
		}
	})
	var refused *Refusal
	if !errors.As(err, &refused) || refused.Code != "done-epic-open-tasks" {
		t.Fatalf("stale eligibility=%v", err)
	}
	issue, err := issues.Get(issues.GetOptions{WorkspaceRoot: root, MigrationsFS: migrations, ID: epic.ID})
	if err != nil || issue.Status != "todo" || issue.Revision != 1 {
		t.Fatalf("refused epic=%+v,%v", issue, err)
	}
	if _, err := journal.EnsureProjectionFresh(root, migrations); err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) { os.Exit(gittest.Main(m)) }
