package export

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// Every worktree of a repository shares its ledger, so an imported ledger's
// binding reads from each of them, still naming the worktree it was imported
// into. The binding travels with the ledger: moved into another repository,
// it still reads there, but that repository's root never stands in for the
// import's.
func TestImportBindingReadsFromEveryWorktreeThatSharesTheLedger(t *testing.T) {
	t.Parallel()
	archive, target := boundaryImportFixture(t)
	target = gittest.ResolvedPath(t, target)
	_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "binding"})
	if err != nil {
		t.Fatal(err)
	}
	linked := gittest.AddWorktree(t, target)
	binding, err := config.ReadImportBinding(linked)
	if err != nil {
		t.Fatalf("read the binding from a linked worktree: %v", err)
	}
	if binding.TargetRoot != target {
		t.Fatalf("binding target root = %q, want the import's target %q", binding.TargetRoot, target)
	}
	other := gittest.ResolvedPath(t, gittest.NewWorkspace(t))
	if err := os.Rename(config.LedgerDir(target), config.LedgerDir(other)); err != nil {
		t.Fatal(err)
	}
	if moved, err := config.ReadImportBinding(other); err != nil || moved.TargetRoot != target {
		t.Fatalf("binding of the moved ledger = %#v, %v; want the import's target %q", moved, err, target)
	}
}

// The ledger outlives the worktree an archive was imported into, and so does
// its provenance: after that worktree is removed, the main worktree still
// reads the binding, still naming the import's original destination.
func TestImportProvenanceSurvivesRemovingTheImportWorktree(t *testing.T) {
	t.Parallel()
	archive, mainWorktree := boundaryImportFixture(t)
	mainWorktree = gittest.ResolvedPath(t, mainWorktree)
	linked := gittest.ResolvedPath(t, gittest.AddWorktree(t, mainWorktree))
	result, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: linked, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "provenance"})
	if err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("git", "-C", mainWorktree, "worktree", "remove", linked).CombinedOutput(); err != nil {
		t.Fatalf("git worktree remove: %v: %s", err, output)
	}
	if _, err := os.Stat(config.LedgerDir(mainWorktree)); err != nil {
		t.Fatalf("shared ledger after removing the import worktree: %v", err)
	}
	binding, err := config.ReadImportBinding(mainWorktree)
	if err != nil {
		t.Fatalf("read import binding from the main worktree: %v", err)
	}
	if binding.TargetRoot != linked || binding.ProjectID != result.ProjectID || len(binding.ImportedEventIDs) == 0 {
		t.Fatalf("import binding = %#v, want the import into %q", binding, linked)
	}
}

// Import stages beside the ledger, in the git common directory: its snapshot
// of the archive, then the staged ledger it renames into place. It writes
// nothing in any worktree. The test looks while both exist, importing into a
// linked worktree, whose common directory is the main worktree's .git.
func TestImportStagesInTheGitCommonDirectory(t *testing.T) {
	t.Parallel()
	archive, mainWorktree := boundaryImportFixture(t)
	mainWorktree = gittest.ResolvedPath(t, mainWorktree)
	linked := gittest.ResolvedPath(t, gittest.AddWorktree(t, mainWorktree))
	commonDir := filepath.Join(mainWorktree, ".git")
	inspected := map[ImportStage]bool{}
	_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: linked, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "staging", Fault: func(stage ImportStage) error {
		switch stage {
		case ImportAfterValidation:
			assertImportStaging(t, commonDir, []string{mainWorktree, linked}, false)
		case ImportBeforePublication:
			assertImportStaging(t, commonDir, []string{mainWorktree, linked}, true)
		default:
			return nil
		}
		inspected[stage] = true
		return nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !inspected[ImportAfterValidation] || !inspected[ImportBeforePublication] {
		t.Fatalf("inspected stages = %v, want after-validation and before-publication", inspected)
	}
	if _, err := os.Stat(filepath.Join(commonDir, "fid", config.ImportBindingFileName)); err != nil {
		t.Fatalf("imported ledger: %v", err)
	}
	entries, err := os.ReadDir(commonDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "fid-staging-") {
			t.Fatalf("import left staging %q", entry.Name())
		}
	}
	for _, worktree := range []string{mainWorktree, linked} {
		assertOnlyGitIn(t, worktree)
	}
}

// assertImportStaging checks that commonDir holds import's archive snapshot,
// one fid-staging-* file, and, once staged, one fid-staging-* directory with
// the staged ledger and its binding; that the ledger is not there yet; and
// that every worktree holds nothing but .git.
func assertImportStaging(t *testing.T, commonDir string, worktrees []string, staged bool) {
	t.Helper()
	entries, err := os.ReadDir(commonDir)
	if err != nil {
		t.Fatal(err)
	}
	var snapshots, stages []string
	for _, entry := range entries {
		switch {
		case entry.Name() == "fid":
			t.Fatalf("ledger exists in %q before import publishes it", commonDir)
		case strings.HasPrefix(entry.Name(), "fid-staging-") && entry.Type().IsRegular():
			snapshots = append(snapshots, entry.Name())
		case strings.HasPrefix(entry.Name(), "fid-staging-") && entry.IsDir():
			stages = append(stages, entry.Name())
		}
	}
	wantStages := 0
	if staged {
		wantStages = 1
	}
	if len(snapshots) != 1 || len(stages) != wantStages {
		t.Fatalf("staging in %q: snapshots %v, stages %v; want one snapshot and %d stage", commonDir, snapshots, stages, wantStages)
	}
	if staged {
		if _, err := os.Stat(filepath.Join(commonDir, stages[0], "fid", config.ImportBindingFileName)); err != nil {
			t.Fatalf("staged ledger in %q: %v", stages[0], err)
		}
	}
	for _, worktree := range worktrees {
		assertOnlyGitIn(t, worktree)
	}
}

func assertOnlyGitIn(t *testing.T, worktree string) {
	t.Helper()
	entries, err := os.ReadDir(worktree)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Name() != ".git" {
			t.Fatalf("worktree %q holds %q, want only .git", worktree, entry.Name())
		}
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
