package bootstrap

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestBuildProjectInitializedEventMatchesV2Golden(t *testing.T) {
	event := BuildProjectInitializedEvent("project-1")
	event.EventID = "EVT-20260822-000001"
	event.OccurredAt = time.Date(2026, time.August, 22, 0, 0, 0, 0, time.UTC)

	got, err := journal.MarshalCanonicalJSON(event)
	if err != nil {
		t.Fatalf("MarshalCanonicalJSON() error = %v", err)
	}

	want := []byte(`{"actor":{"actor_id":"fid-bootstrap","actor_type":"service","producer":"fid","provenance_ref":"fid:init","session_ref":null},"aggregate":{"id":"project-1","sequence":1,"type":"project"},"causation_id":null,"command":"project.init","correlation_id":null,"effects":[],"event_id":"EVT-20260822-000001","event_schema":2,"event_type":"project.initialized","idempotency":null,"occurred_at":"2026-08-22T00:00:00Z","repository":null}`)
	if !bytes.Equal(got, want) {
		t.Fatalf("bootstrap Event = %s, want %s", got, want)
	}
}

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
