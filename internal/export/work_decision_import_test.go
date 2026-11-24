package export

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/mbell-nzl/fid/internal/workrecord"
)

func TestWorkDecisionsSurviveExportImport(t *testing.T) {
	t.Parallel()
	var accepted, pending string
	f := newExportSourceFixture(t, func(t *testing.T, root string) string {
		t.Helper()
		lane := "standard"
		for i := range 2 {
			task, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), Kind: "task", Title: "Agent discovery", Actor: "agent:orchestrator", Lane: &lane, IdempotencyKey: fmt.Sprintf("create-%d", i)})
			if err != nil {
				t.Fatal(err)
			}
			if i == 1 {
				pending = task.ID
				continue
			}
			accepted = task.ID
			if _, err := holdwork.Accept(holdwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), TaskID: accepted, Actor: "human:mark", Reason: "Approved discovery"}); err != nil {
				t.Fatal(err)
			}
			if _, err := holdwork.AllowRework(holdwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), TaskID: accepted, Actor: "human:mark", Reason: "Continue the fix", Rounds: 2}); err != nil {
				t.Fatal(err)
			}
		}
		return accepted
	})
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: f.archive, RepositoryRoot: target, MigrationsFS: testMigrationsFS(t), IdempotencyKey: "work-decisions"}); err != nil {
		t.Fatal(err)
	}
	for _, root := range []string{f.workspace, target} {
		for _, rebuild := range []bool{false, true} {
			if rebuild {
				if _, err := journal.RebuildWorkspace(root, testMigrationsFS(t)); err != nil {
					t.Fatal(err)
				}
			}
			conn, err := reviewwork.OpenProjection(root)
			if err != nil {
				t.Fatal(err)
			}
			unaccepted, err := workrecord.Unaccepted(conn)
			if err != nil {
				t.Fatal(err)
			}
			cap, latest, err := workrecord.Allowances(conn, accepted, 3)
			if err := conn.Close(); err != nil {
				t.Fatal(err)
			}
			if err != nil || !reflect.DeepEqual(unaccepted, map[string]string{pending: "agent:orchestrator"}) || cap != 5 || latest.Reason != "Continue the fix" {
				t.Fatalf("%s restored %v/%d/%v/%v", filepath.Base(root), unaccepted, cap, latest, err)
			}
		}
	}
	assertImportedJournalEquals(t, f.workspace, target)
}
