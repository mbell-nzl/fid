package export

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestImportProcessExitAtPublicationBoundaries(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, nil)
	for _, stage := range []ImportStage{ImportAfterValidation, ImportBeforeDurability, ImportBeforePublication, ImportAfterRename, ImportAfterPublication} {
		t.Run(string(stage), func(t *testing.T) {
			_, archive := fixture.clone(t)
			target := gittest.NewWorkspace(t)
			cmd := exec.Command(os.Args[0], "-test.run=^TestImportAbruptExitHelper$")
			cmd.Env = append(os.Environ(), "FID_IMPORT_CRASH_ARCHIVE="+archive, "FID_IMPORT_CRASH_TARGET="+target, "FID_IMPORT_CRASH_STAGE="+string(stage), "FID_IMPORT_CRASH_MIGRATIONS="+testMigrationsDir(t))
			output, err := cmd.CombinedOutput()
			var exited *exec.ExitError
			if !errors.As(err, &exited) || exited.ExitCode() != 86 {
				t.Fatalf("crash helper: %v / %s", err, output)
			}
			published := stage == ImportAfterRename || stage == ImportAfterPublication
			_, statErr := os.Stat(config.LedgerDir(target))
			if published && statErr != nil || !published && !os.IsNotExist(statErr) {
				t.Fatalf("published=%v, destination stat=%v", published, statErr)
			}
			result, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "crash"})
			if err != nil || result.Replayed != published {
				t.Fatalf("retry after abrupt process exit: %#v, %v", result, err)
			}
			if _, err := config.ReadImportBinding(target); err != nil {
				t.Fatalf("retry did not leave verified complete installation: %v", err)
			}
		})
	}
}

func TestImportAbruptExitHelper(t *testing.T) {
	t.Parallel()
	archive := os.Getenv("FID_IMPORT_CRASH_ARCHIVE")
	if archive == "" {
		return
	}
	_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: os.Getenv("FID_IMPORT_CRASH_TARGET"), MigrationsDir: os.Getenv("FID_IMPORT_CRASH_MIGRATIONS"), IdempotencyKey: "crash", Fault: func(stage ImportStage) error {
		if string(stage) == os.Getenv("FID_IMPORT_CRASH_STAGE") {
			os.Exit(86)
		}
		return nil
	}})
	t.Fatalf("crash boundary not reached: %v", err)
}

func TestImportRejectsChangedDigestAndInstalledHistory(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, nil)
	for _, corruption := range []string{"changed-digest", "journal-bytes", "journal-missing", "result-count", "symlink-workspace"} {
		t.Run(corruption, func(t *testing.T) {
			_, archive := fixture.clone(t)
			target := gittest.NewWorkspace(t)
			options := ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "history"}
			if _, err := Import(options); err != nil {
				t.Fatal(err)
			}
			markerPath := filepath.Join(config.LedgerDir(target), config.ImportBindingFileName)
			data, err := os.ReadFile(markerPath)
			if err != nil {
				t.Fatal(err)
			}
			var marker importMarker
			if err := json.Unmarshal(data, &marker); err != nil {
				t.Fatal(err)
			}
			journalPath := filepath.Join(config.LedgerDir(target), "journal", "00000001.jsonl")
			switch corruption {
			case "changed-digest":
				options.ArchivePath, _ = boundaryImportFixture(t)
			case "journal-bytes":
				data, err := os.ReadFile(journalPath)
				if err != nil {
					t.Fatal(err)
				}
				data[0] = ' '
				if err := os.WriteFile(journalPath, data, 0o600); err != nil {
					t.Fatal(err)
				}
			case "journal-missing":
				if err := os.Remove(journalPath); err != nil {
					t.Fatal(err)
				}
			case "result-count":
				marker.Result.EventCount++
				data, err = json.Marshal(marker)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(markerPath, data, 0o600); err != nil {
					t.Fatal(err)
				}
			case "symlink-workspace":
				moved := filepath.Join(t.TempDir(), "installed")
				if err := os.Rename(config.LedgerDir(target), moved); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(moved, config.LedgerDir(target)); err != nil {
					t.Fatal(err)
				}
			}
			before := workspaceTreeDigest(t, target)
			_, err = Import(options)
			var conflict *ImportTargetConflictError
			if !errors.As(err, &conflict) {
				t.Fatalf("corrupt or changed installation replay accepted: %v", err)
			}
			if after := workspaceTreeDigest(t, target); after != before {
				t.Fatal("rejected retry altered destination")
			}
		})
	}
}

func TestImportLockExcludesRepeatedContenders(t *testing.T) {
	t.Parallel()
	target := t.TempDir()
	err := withImportLock(target, func() error {
		for attempt := 0; attempt < 2; attempt++ {
			entered := false
			err := withImportLock(target, func() error { entered = true; return nil })
			if entered || err == nil {
				t.Errorf("contender %d entered while owner still held import lock", attempt+1)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := withImportLock(target, func() error { return nil }); err != nil {
		t.Fatalf("lock after owner released: %v", err)
	}
}

func TestImportBindingRejectsAlteredHistoryIdentity(t *testing.T) {
	t.Parallel()
	archive, target := boundaryImportFixture(t)
	if _, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "binding"}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(config.LedgerDir(target), config.ImportBindingFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var marker importMarker
	if err := json.Unmarshal(data, &marker); err != nil {
		t.Fatal(err)
	}
	marker.ImportedEventIDs[0] = "EVT-20260908-999999"
	data, err = json.Marshal(marker)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := config.ReadImportBinding(target); err == nil {
		t.Error("altered imported Event identity still qualifies installed binding")
	}
	if _, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "binding"}); err == nil {
		t.Error("altered imported Event identity still replays installation")
	}
}

func boundaryImportFixture(t *testing.T) (string, string) {
	t.Helper()
	source := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	return archive, gittest.NewWorkspace(t)
}
