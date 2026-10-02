package workflow

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/project"
)

func TestRecoveryFixturesKeepRepositoryAndProjectionIsolation(t *testing.T) {
	t.Parallel()
	templatePath := filepath.Join(workflowPristineWorkspace, ".fid", "fid.sqlite")
	before, err := os.ReadFile(templatePath)
	mustRecoveryTest(t, err)
	first, _ := recoveryFixture(t)
	second, _ := recoveryFixture(t)
	if first.WorkspaceRoot == second.WorkspaceRoot {
		t.Fatal("recovery fixtures share a workspace")
	}
	firstRun, err := GetRun(first.WorkspaceRoot, first.MigrationsDir, "RUN-prior")
	mustRecoveryTest(t, err)
	secondRun, err := GetRun(second.WorkspaceRoot, second.MigrationsDir, "RUN-prior")
	mustRecoveryTest(t, err)
	if firstRun.StartRepository.WorktreeID == secondRun.StartRepository.WorktreeID {
		t.Fatal("recovery fixtures share a repository binding")
	}
	for _, options := range []RunOptions{first, second} {
		observed, err := project.ObserveRepository(options.WorkspaceRoot)
		mustRecoveryTest(t, err)
		run, err := GetRun(options.WorkspaceRoot, options.MigrationsDir, "RUN-prior")
		mustRecoveryTest(t, err)
		if run.StartRepository.WorktreeID != observed.WorktreeID {
			t.Fatal("recovery fixture retained another workspace's repository")
		}
	}
	conn, err := openWorkflowDB(first.WorkspaceRoot)
	mustRecoveryTest(t, err)
	_, err = conn.Exec("UPDATE issues SET title = 'only-first' WHERE id = 'ISS-recovery'")
	closeWorkflowDB(conn)
	mustRecoveryTest(t, err)
	other, err := openWorkflowDB(second.WorkspaceRoot)
	mustRecoveryTest(t, err)
	defer closeWorkflowDB(other)
	var title string
	mustRecoveryTest(t, other.QueryRow("SELECT title FROM issues WHERE id = 'ISS-recovery'").Scan(&title))
	if title != "recovery fixture" {
		t.Fatalf("mutation crossed into sibling fixture: %q", title)
	}
	after, err := os.ReadFile(templatePath)
	mustRecoveryTest(t, err)
	if !bytes.Equal(before, after) {
		t.Fatal("recovery fixture changed pristine database bytes")
	}
}
