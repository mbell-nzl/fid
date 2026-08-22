package bootstrap

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestInitializeFailureBeforeActivationLeavesNoWorkspace(t *testing.T) {
	for _, step := range []initStage{
		initStageValidated,
		initStageStaged,
		initStageSchemaInitialized,
		initStageGenesisAppended,
		initStageRebuilt,
		initStageCompared,
		initStageSynced,
		initStageActivated,
	} {
		step := step
		t.Run(string(step), func(t *testing.T) {
			workspaceRoot := t.TempDir()
			failure := errors.New("injected " + string(step) + " failure")

			_, err := initializeWorkspaceWithHooks(InitOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: migrationFixtureDir(t),
				Prefix:        "FID",
			}, initHooks{before: func(current initStage) error {
				if current == step {
					return failure
				}
				return nil
			}})
			if !errors.Is(err, failure) {
				t.Fatalf("InitializeWorkspace() error = %v, want injected failure", err)
			}

			assertNoActivationResidue(t, workspaceRoot)
		})
	}
}

func assertNoActivationResidue(t *testing.T, workspaceRoot string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(workspaceRoot, config.WorkspaceDirName)); !os.IsNotExist(err) {
		t.Fatalf("activated workspace state = %v, want absent", err)
	}

	entries, err := os.ReadDir(workspaceRoot)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", workspaceRoot, err)
	}
	for _, entry := range entries {
		if entry.IsDir() && isStagingDirectoryName(entry.Name()) {
			t.Fatalf("staging residue %q remains", entry.Name())
		}
	}
}

func migrationFixtureDir(t *testing.T) string {
	t.Helper()
	return filepath.Clean(filepath.Join("..", "..", "migrations"))
}
