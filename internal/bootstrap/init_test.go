package bootstrap

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestBuildProjectInitializedEventMatchesV2Golden(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
			workspaceRoot := gittest.NewWorkspace(t)
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

func TestInitializeRefusesIdentityChangedWhileStaging(t *testing.T) {
	t.Parallel()
	for _, initialID := range []string{"", "project-11111111111111111111111111111111"} {
		name := "appears"
		if initialID != "" {
			name = "replaced"
		}
		t.Run(name, func(t *testing.T) {
			root := gittest.NewWorkspace(t)
			path := filepath.Join(root, config.ProjectIdentityFileName)
			if initialID != "" {
				if err := os.WriteFile(path, []byte("[project]\nid = \""+initialID+"\"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			other := []byte("# written while init stages\n[project]\nid = \"project-22222222222222222222222222222222\"\n")
			_, err := initializeWorkspaceWithHooks(InitOptions{
				WorkspaceRoot: root, MigrationsDir: migrationFixtureDir(t), Prefix: "FID",
			}, initHooks{before: func(stage initStage) error {
				if stage == initStageActivated {
					return os.WriteFile(path, other, 0o644)
				}
				return nil
			}})
			if err == nil || !strings.Contains(err.Error(), "changed during initialization") {
				t.Errorf("init error = %v, want changed-identity refusal", err)
			}
			if got, err := os.ReadFile(path); err != nil || !bytes.Equal(got, other) {
				t.Errorf("identity = %q, %v, want unchanged %q", got, err, other)
			}
			if _, err := os.Lstat(config.LedgerDir(root)); !os.IsNotExist(err) {
				t.Errorf("ledger state = %v, want absent", err)
			}
			entries, err := os.ReadDir(filepath.Dir(config.LedgerDir(root)))
			if err != nil {
				t.Fatal(err)
			}
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), config.LedgerStagingPrefix) {
					t.Errorf("staging residue %q remains", entry.Name())
				}
			}
		})
	}
}

func assertNoActivationResidue(t *testing.T, workspaceRoot string) {
	t.Helper()
	if _, err := os.Stat(config.LedgerDir(workspaceRoot)); !os.IsNotExist(err) {
		t.Fatalf("activated workspace state = %v, want absent", err)
	}

	commonDir := filepath.Dir(config.LedgerDir(workspaceRoot))
	entries, err := os.ReadDir(commonDir)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", commonDir, err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "fid-staging-") {
			t.Fatalf("staging residue %q remains", entry.Name())
		}
	}
	entries, err = os.ReadDir(workspaceRoot)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", workspaceRoot, err)
	}
	for _, entry := range entries {
		if entry.Name() != ".git" {
			t.Fatalf("init left %q in the worktree, want only .git", entry.Name())
		}
	}
}

func migrationFixtureDir(t *testing.T) string {
	t.Helper()
	return filepath.Clean(filepath.Join("..", "..", "migrations"))
}
