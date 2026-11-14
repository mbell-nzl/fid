package export

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
)

func TestImportPublishesIdentityAndReplaysLostResponse(t *testing.T) {
	t.Parallel()
	source := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize source workspace: %v", err)
	}
	addLeanArtifactExportFacts(t, source)
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	target := gittest.NewWorkspace(t)
	keep := filepath.Join(target, "source.txt")
	if err := os.WriteFile(keep, []byte("preserve this source file\n"), 0o644); err != nil {
		t.Fatalf("write target source file: %v", err)
	}
	result, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "same-install"})
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if result.Replayed {
		t.Fatal("first import was marked replayed")
	}
	binding, err := config.ReadImportBinding(target)
	if err != nil || binding.ProjectID != result.ProjectID || len(binding.ImportedEventIDs) == 0 {
		t.Fatalf("installed import binding = %#v, %v", binding, err)
	}
	if data, err := os.ReadFile(keep); err != nil || string(data) != "preserve this source file\n" {
		t.Fatalf("target source file changed during import: %q, %v", data, err)
	}
	assertImportedJournalEquals(t, source, target)
	if _, err := issues.Create(issues.CreateOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: "After import", Actor: "human:mark", IdempotencyKey: "post-import-task"}); err != nil {
		t.Fatal(err)
	}
	before := workspaceTreeDigest(t, target)
	replayed, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "same-install"})
	if err != nil || !replayed.Replayed || replayed.ProjectID != result.ProjectID {
		t.Fatalf("same import replay = %#v, %v", replayed, err)
	}
	if after := workspaceTreeDigest(t, target); after != before {
		t.Fatal("same import replay changed installed facts")
	}
	_, err = Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "other-install"})
	var conflict *ImportTargetConflictError
	if !errors.As(err, &conflict) || conflict.Reason != "different-import" {
		t.Fatalf("different import error = %#v, want different-import conflict", err)
	}
}

func TestImportRejectsExistingAndUnsafeTargetsWithoutChangingThem(t *testing.T) {
	t.Parallel()
	source := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize source workspace: %v", err)
	}
	addLeanArtifactExportFacts(t, source)
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	populated := gittest.NewWorkspace(t)
	if err := os.Mkdir(config.LedgerDir(populated), 0o755); err != nil {
		t.Fatalf("create populated target: %v", err)
	}
	marker := filepath.Join(config.LedgerDir(populated), "keep")
	if err := os.WriteFile(marker, []byte("unchanged"), 0o600); err != nil {
		t.Fatalf("write populated target marker: %v", err)
	}
	_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: populated, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "conflict"})
	var conflict *ImportTargetConflictError
	if !errors.As(err, &conflict) || conflict.Reason != "populated" {
		t.Fatalf("populated import error = %#v, want populated conflict", err)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "unchanged" {
		t.Fatalf("populated target changed: %q, %v", data, err)
	}
	root := gittest.NewWorkspace(t)
	unsafe := filepath.Join(t.TempDir(), "target-link")
	if err := os.Symlink(root, unsafe); err != nil {
		t.Fatalf("create target symlink: %v", err)
	}
	_, err = Import(ImportOptions{ArchivePath: archive, RepositoryRoot: unsafe, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "unsafe"})
	if !errors.As(err, &conflict) || conflict.Reason != "unsafe-path" {
		t.Fatalf("unsafe import error = %#v, want unsafe-path conflict", err)
	}
}

func TestImportAfterPublicationFaultReplaysInstalledResult(t *testing.T) {
	t.Parallel()
	source := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize source workspace: %v", err)
	}
	addLeanArtifactExportFacts(t, source)
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	target := gittest.NewWorkspace(t)
	_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "lost-response", Fault: func(stage ImportStage) error {
		if stage == ImportAfterPublication {
			return ErrImportFault
		}
		return nil
	}})
	if !errors.Is(err, ErrImportFault) {
		t.Fatalf("post-publication import error = %v, want injected response-loss fault", err)
	}
	replayed, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "lost-response"})
	if err != nil || !replayed.Replayed {
		t.Fatalf("response-loss replay = %#v, %v", replayed, err)
	}
}

func TestImportSnapshotsInputPreservesExistingFilesAndRejectsTamperedMarker(t *testing.T) {
	t.Parallel()
	source := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	addLeanArtifactExportFacts(t, source)
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	archiveBytes, err := os.ReadFile(archive)
	if err != nil {
		t.Fatal(err)
	}
	target := gittest.NewWorkspace(t)
	if err := os.WriteFile(filepath.Join(target, ".fid-import.lock"), []byte("stale after process exit"), 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "snapshot", Fault: func(stage ImportStage) error {
		if stage == ImportAfterValidation {
			return os.WriteFile(archive, []byte("pathname replaced after snapshot"), 0o600)
		}
		return nil
	}})
	if err != nil {
		t.Fatalf("import from snapshot after source replacement: %v", err)
	}
	if data, err := os.ReadFile(filepath.Join(target, ".fid-import.lock")); err != nil || string(data) != "stale after process exit" {
		t.Fatalf("preexisting target file was changed: %q, %v", data, err)
	}
	tampered, err := json.Marshal(importMarker{ImportBinding: config.ImportBinding{ArchiveDigest: result.ArchiveDigest, IdempotencyKey: "snapshot", ProjectID: "FID-tampered", TargetRoot: target, ImportedEventIDs: []string{"EVT-tampered"}}, Result: ImportResult{ProjectID: "FID-tampered", ArchiveDigest: result.ArchiveDigest, TargetRoot: target}})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(config.LedgerDir(target), config.ImportBindingFileName), tampered, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(archive, archiveBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "snapshot"})
	var conflict *ImportTargetConflictError
	if !errors.As(err, &conflict) || conflict.Reason != "populated" {
		t.Fatalf("tampered installed marker error = %#v, want populated conflict", err)
	}
}

func assertImportedJournalEquals(t *testing.T, source, target string) {
	t.Helper()
	sourceEntries, err := os.ReadDir(filepath.Join(config.LedgerDir(source), "journal"))
	if err != nil {
		t.Fatalf("read source journal: %v", err)
	}
	targetEntries, err := os.ReadDir(filepath.Join(config.LedgerDir(target), "journal"))
	if err != nil {
		t.Fatalf("read target journal: %v", err)
	}
	if !bytes.Equal(journalBytes(t, source, sourceEntries), journalBytes(t, target, targetEntries)) {
		t.Fatal("ordered Event bytes differ after import")
	}
}

func journalBytes(t *testing.T, root string, entries []os.DirEntry) []byte {
	t.Helper()
	var all []byte
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(config.LedgerDir(root), "journal", entry.Name()))
		if err != nil {
			t.Fatalf("read journal %q: %v", entry.Name(), err)
		}
		all = append(all, data...)
	}
	return all
}
