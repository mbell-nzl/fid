package bootstrap

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

func BenchmarkInitializeWorkspace(b *testing.B) {
	stages := []initStage{
		initStageValidated,
		initStageStaged,
		initStageSchemaInitialized,
		initStageGenesisAppended,
		initStageRebuilt,
		initStageCompared,
		initStageSynced,
		initStageActivated,
	}
	stageDurations := make(map[initStage]time.Duration, len(stages))
	var activationDuration time.Duration

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		b.StopTimer()
		workspaceRoot, err := os.MkdirTemp("", "fid-init-benchmark-")
		if err != nil {
			b.Fatal(err)
		}

		lastCheckpoint := time.Now()
		b.StartTimer()
		result, err := initializeWorkspaceWithHooks(InitOptions{
			WorkspaceRoot: workspaceRoot,
			MigrationsDir: filepath.Clean(filepath.Join("..", "..", "migrations")),
			Prefix:        "FID",
		}, initHooks{before: func(stage initStage) error {
			now := time.Now()
			stageDurations[stage] += now.Sub(lastCheckpoint)
			lastCheckpoint = now
			return nil
		}})
		completedAt := time.Now()
		b.StopTimer()
		activationDuration += completedAt.Sub(lastCheckpoint)

		if err != nil {
			b.Fatalf("InitializeWorkspace() error = %v", err)
		}
		paths := config.WorkspacePathsForRoot(workspaceRoot)
		if result.ProjectID == "" {
			b.Fatal("InitializeWorkspace() returned an empty project ID")
		}
		if _, err := os.Stat(paths.DatabasePath); err != nil {
			b.Fatalf("stat initialized database: %v", err)
		}
		if err := os.RemoveAll(workspaceRoot); err != nil {
			b.Fatalf("remove benchmark workspace: %v", err)
		}
	}

	for _, stage := range stages {
		b.ReportMetric(float64(stageDurations[stage].Nanoseconds())/float64(b.N), string(stage)+"-ns/op")
	}
	b.ReportMetric(float64(activationDuration.Nanoseconds())/float64(b.N), "activation-finish-ns/op")
}
