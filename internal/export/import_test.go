package export

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestImportPublishesIdentityAndReplaysLostResponse(t *testing.T) {
	source := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize source workspace: %v", err)
	}
	addJournalBackedExportFixtureFacts(t, source)
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	target := t.TempDir()
	fixtureGit(t, target, "init", "-q")
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
	source := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize source workspace: %v", err)
	}
	addJournalBackedExportFixtureFacts(t, source)
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	populated := t.TempDir()
	if err := os.Mkdir(filepath.Join(populated, ".fid"), 0o755); err != nil {
		t.Fatalf("create populated target: %v", err)
	}
	marker := filepath.Join(populated, ".fid", "keep")
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
	root := t.TempDir()
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
	source := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize source workspace: %v", err)
	}
	addJournalBackedExportFixtureFacts(t, source)
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	target := t.TempDir()
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

func TestImportRequiresFreshDestinationObservationBeforeRecovery(t *testing.T) {
	source := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	addJournalBackedExportFixtureFacts(t, source)
	sourceRuns, err := workflow.ListRuns(workflow.RunListOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Status: "active", Limit: 1})
	if err != nil || len(sourceRuns) != 1 {
		t.Fatalf("source active Run = %#v, %v", sourceRuns, err)
	}
	if _, err := workflow.ObserveRun(workflow.RunOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), ActorID: "actor-harness", IdempotencyKey: "source-imported-observation"}, workflow.RunObserveRequest{RunID: sourceRuns[0].ID, ExpectedVersion: sourceRuns[0].Version, Liveness: workflow.Liveness{ExecutorID: sourceRuns[0].Executor.ID, State: "terminated", SourceRef: "source:historical", ObservedAt: time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}}); err != nil {
		t.Fatalf("observe source Run: %v", err)
	}
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "destination")
	if output, err := exec.Command("git", "clone", "-q", source, target).CombinedOutput(); err != nil {
		t.Fatalf("clone destination: %v: %s", err, output)
	}
	if _, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "destination-install"}); err != nil {
		t.Fatal(err)
	}
	runs, err := workflow.ListRuns(workflow.RunListOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), Status: "active", Limit: 1})
	if err != nil || len(runs) != 1 {
		t.Fatalf("imported active Run = %#v, %v", runs, err)
	}
	run := runs[0]
	inspection, err := workflow.InspectRun(target, testMigrationsDir(t), run.ID, target)
	if err != nil {
		t.Fatalf("inspect imported Run at destination: %v", err)
	}
	if inspection.Reconciliation.Repository.WorktreeID == run.StartRepository.WorktreeID {
		t.Fatal("imported inspection retained the historical Worktree ID")
	}
	database, err := fiddb.Open(filepath.Join(target, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	var importedObservationID string
	if err := database.QueryRow(`SELECT id FROM run_observations WHERE run_id = ? AND kind = 'liveness' ORDER BY observed_at DESC, id DESC LIMIT 1`, run.ID).Scan(&importedObservationID); err != nil {
		t.Fatal(err)
	}
	var registrationID string
	var registrationVersion int
	if err := database.QueryRow(`SELECT id, version FROM gate_registrations LIMIT 1`).Scan(&registrationID, &registrationVersion); err != nil {
		t.Fatal(err)
	}
	_ = database.Close()
	_, err = workflow.CloseRun(workflow.RunOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ActorID: "actor-builder", IdempotencyKey: "old-observation-close"}, workflow.RunCloseRequest{RunID: run.ID, ExpectedVersion: run.Version, ObservationID: importedObservationID, Repository: inspection.Reconciliation.Repository})
	if err == nil || !strings.Contains(err.Error(), "imported history") {
		t.Fatalf("close with imported observation error = %v, want imported-history rejection", err)
	}
	_, err = workflow.AbandonRun(workflow.RunOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ActorID: "actor-builder", IdempotencyKey: "old-observation-abandon"}, workflow.RunAbandonRequest{RunID: run.ID, ExpectedVersion: run.Version, ExpectedClaimVersion: 1, ObservationID: importedObservationID, Authority: workflow.RecoveryAuthority{Mode: "ordinary", ActorID: "actor-builder", Reason: "fresh local liveness required"}, Repository: inspection.Reconciliation.Repository})
	if err == nil || !strings.Contains(err.Error(), "imported history") {
		t.Fatalf("abandon with imported observation error = %v, want imported-history rejection", err)
	}
	_, err = evidence.RunGate(context.Background(), evidence.GateRunOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ActorID: "actor-builder", IdempotencyKey: "historical-gate", Request: evidence.GateRunRequest{RegistrationID: registrationID, RegistrationVersion: registrationVersion, ObligationID: "OBL-fixture", RunID: run.ID, Repository: evidenceRepository(inspection.Reconciliation.Repository)}})
	if err == nil || !strings.Contains(err.Error(), "repository binding is invalid") {
		t.Fatalf("execute imported historical Gate error = %v, want destination binding rejection", err)
	}
	_, err = workflow.RecoverRun(workflow.RunOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ActorID: "actor-builder", IdempotencyKey: "old-observation"}, workflow.RunRecoveryRequest{PriorRunID: run.ID, ExpectedRunVersion: run.Version, ExpectedClaimVersion: 1, ObservationID: importedObservationID, ReconciliationDigest: inspection.Reconciliation.Digest, AdoptedPaths: []string{}, Authority: workflow.RecoveryAuthority{Mode: "ordinary", ActorID: "actor-builder", Reason: "fresh local liveness required"}, Executor: workflow.Executor{ID: "executor-target", Source: "activated-harness"}, TTLSeconds: 60})
	if err == nil {
		t.Fatal("imported active Run recovered without a fresh destination observation")
	}
	observed, err := workflow.ObserveRun(workflow.RunOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ActorID: "actor-harness", IdempotencyKey: "destination-observe"}, workflow.RunObserveRequest{RunID: run.ID, ExpectedVersion: run.Version, Liveness: workflow.Liveness{ExecutorID: run.Executor.ID, State: "terminated", SourceRef: "destination:test", ObservedAt: time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
	if err != nil {
		t.Fatalf("observe imported Run at destination: %v", err)
	}
	if _, err := claims.ReleaseGeneric(claims.MutationOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ActorID: "actor-builder", IdempotencyKey: "destination-release"}, claims.ReleaseRequest{ClaimID: run.ClaimID, ExpectedVersion: 1, Reason: "operator-released"}); err != nil {
		t.Fatalf("release imported Claim: %v", err)
	}
	inspection, err = workflow.InspectRun(target, testMigrationsDir(t), run.ID, target)
	if err != nil {
		t.Fatalf("inspect fresh destination state: %v", err)
	}
	recovered, err := workflow.RecoverRun(workflow.RunOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ActorID: "actor-builder", IdempotencyKey: "destination-recover"}, workflow.RunRecoveryRequest{PriorRunID: run.ID, ExpectedRunVersion: observed.Run.Version, ExpectedClaimVersion: 2, ObservationID: observed.Observation.ID, ReconciliationDigest: inspection.Reconciliation.Digest, AdoptedPaths: []string{}, Authority: workflow.RecoveryAuthority{Mode: "ordinary", ActorID: "actor-builder", Reason: "fresh destination observation"}, Executor: workflow.Executor{ID: "executor-target", Source: "activated-harness"}, TTLSeconds: 60})
	if err != nil {
		t.Fatalf("recover imported Run at destination: %v", err)
	}
	if recovered.Claim.WorktreeID != inspection.Reconciliation.Repository.WorktreeID || recovered.Run.StartRepository.WorktreeID != recovered.Claim.WorktreeID {
		t.Fatalf("destination successor bindings = %#v / %#v", recovered.Claim, recovered.Run.StartRepository)
	}
	if err := os.WriteFile(filepath.Join(target, "fixture-artifact.txt"), []byte("artifact bytes\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gate, err := evidence.RunGate(context.Background(), evidence.GateRunOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ActorID: "actor-builder", IdempotencyKey: "destination-gate", Request: evidence.GateRunRequest{RegistrationID: registrationID, RegistrationVersion: registrationVersion, ObligationID: "OBL-fixture", RunID: recovered.Run.ID, Repository: evidenceRepository(recovered.Run.StartRepository)}})
	if err != nil || gate.Receipt == nil || gate.Receipt.RunID != recovered.Run.ID {
		t.Fatalf("execute successor Gate = %#v, %v", gate, err)
	}
}

func assertImportedJournalEquals(t *testing.T, source, target string) {
	t.Helper()
	sourceEntries, err := os.ReadDir(filepath.Join(source, ".fid", "journal"))
	if err != nil {
		t.Fatalf("read source journal: %v", err)
	}
	targetEntries, err := os.ReadDir(filepath.Join(target, ".fid", "journal"))
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
		data, err := os.ReadFile(filepath.Join(root, ".fid", "journal", entry.Name()))
		if err != nil {
			t.Fatalf("read journal %q: %v", entry.Name(), err)
		}
		all = append(all, data...)
	}
	return all
}
