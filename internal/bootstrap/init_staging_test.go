package bootstrap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// Init stages the ledger beside its target, in the git common directory, so
// the rename into place stays on one filesystem. The test looks while the
// staged ledger exists, before init publishes .fid.toml, from a main worktree
// and from a linked one.
func TestInitStagesTheLedgerInTheGitCommonDirectory(t *testing.T) {
	t.Parallel()
	for _, fromLinked := range []bool{false, true} {
		mainWorktree := gittest.NewWorkspace(t)
		worktrees := []string{mainWorktree}
		start := mainWorktree
		if fromLinked {
			start = gittest.AddWorktree(t, mainWorktree)
			worktrees = append(worktrees, start)
		}
		commonDir := filepath.Join(gittest.ResolvedPath(t, mainWorktree), ".git")
		inspected := 0
		_, err := initializeWorkspaceWithHooks(InitOptions{
			WorkspaceRoot: start,
			MigrationsDir: migrationFixtureDir(t),
			Prefix:        "FID",
		}, initHooks{before: func(stage initStage) error {
			if stage == initStageStaged || stage == initStageActivated {
				assertStagedLedgerInCommonDir(t, commonDir, worktrees)
				inspected++
			}
			return nil
		}})
		if err != nil {
			t.Fatalf("init (from linked %t): %v", fromLinked, err)
		}
		if inspected != 2 {
			t.Fatalf("init (from linked %t) inspected staging %d times, want 2", fromLinked, inspected)
		}
		if _, err := os.Stat(filepath.Join(commonDir, "fid", "config.toml")); err != nil {
			t.Fatalf("init (from linked %t) ledger: %v", fromLinked, err)
		}
		entries, err := os.ReadDir(commonDir)
		if err != nil {
			t.Fatalf("ReadDir(%q) error = %v", commonDir, err)
		}
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "fid-staging-") {
				t.Fatalf("init (from linked %t) left staging %q", fromLinked, entry.Name())
			}
		}
		for _, worktree := range worktrees {
			if worktree != start {
				assertOnlyGitIn(t, worktree)
				continue
			}
			id, err := config.ReadProjectIdentity(worktree)
			if err != nil || id == "" {
				t.Fatalf("initialized worktree identity = %q, %v", id, err)
			}
			entries, err := os.ReadDir(worktree)
			if err != nil || len(entries) != 2 {
				t.Fatalf("initialized worktree entries = %v, %v, want .git and .fid.toml", entries, err)
			}
		}
	}
}

// assertStagedLedgerInCommonDir checks that exactly one staging directory
// sits in commonDir, holding the staged ledger, that the ledger itself is not
// there yet, and that every worktree holds nothing but .git.
func assertStagedLedgerInCommonDir(t *testing.T, commonDir string, worktrees []string) {
	t.Helper()
	entries, err := os.ReadDir(commonDir)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", commonDir, err)
	}
	var staging []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "fid-staging-") {
			staging = append(staging, entry.Name())
		}
		if entry.Name() == "fid" {
			t.Fatalf("ledger %q exists while init is staging it", filepath.Join(commonDir, "fid"))
		}
	}
	if len(staging) != 1 {
		t.Fatalf("staging in %q = %v, want one fid-staging-* directory", commonDir, staging)
	}
	if _, err := os.Stat(filepath.Join(commonDir, staging[0], "fid", "config.toml")); err != nil {
		t.Fatalf("staged ledger in %q: %v", staging[0], err)
	}
	for _, worktree := range worktrees {
		assertOnlyGitIn(t, worktree)
	}
}

func assertOnlyGitIn(t *testing.T, worktree string) {
	t.Helper()
	entries, err := os.ReadDir(worktree)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", worktree, err)
	}
	for _, entry := range entries {
		if entry.Name() != ".git" {
			t.Fatalf("worktree %q holds %q, want only .git", worktree, entry.Name())
		}
	}
}
