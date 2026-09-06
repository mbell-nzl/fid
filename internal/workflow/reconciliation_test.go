package workflow

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/project"
)

func TestReconcileRepositoryRequiresExactDirtyAdoption(t *testing.T) {
	repository := reconciliationRepositoryFixture(t)
	start, err := project.ObserveRepository(repository)
	if err != nil {
		t.Fatalf("observe clean repository: %v", err)
	}
	baseline := repositoryFromObservation(start)
	if err := os.WriteFile(filepath.Join(repository, "owned.txt"), []byte("owned\n"), 0o644); err != nil {
		t.Fatalf("write owned change: %v", err)
	}
	current, err := project.ObserveRepository(repository)
	if err != nil {
		t.Fatalf("observe dirty repository: %v", err)
	}
	reconciliation, err := ReconcileRepository(baseline, current, nil, []string{"owned.txt"})
	if err != nil {
		t.Fatalf("ReconcileRepository() error = %v", err)
	}
	if err := ValidateReconciliationAdoption(reconciliation, []string{"owned.txt"}, reconciliation.Digest); err != nil {
		t.Fatalf("ValidateReconciliationAdoption() exact adoption error = %v", err)
	}
	if err := ValidateReconciliationAdoption(reconciliation, nil, reconciliation.Digest); err == nil {
		t.Fatal("ValidateReconciliationAdoption() omitted path error = nil, want recovery-state-conflict")
	}
	if err := ValidateReconciliationAdoption(reconciliation, []string{"owned.txt"}, "sha256:stale"); err == nil {
		t.Fatal("ValidateReconciliationAdoption() stale digest error = nil, want recovery-state-conflict")
	}
}

func TestActivatedCapabilitiesAcceptRunCheckpoint(t *testing.T) {
	if !sortedUniqueCapabilities([]string{"run.checkpoint"}) {
		t.Fatal("activated capabilities reject run.checkpoint")
	}
}

func repositoryFromObservation(observation project.Observation) Repository {
	head := observation.HeadOID
	return Repository{BaseOID: &head, CaptureSource: "observed", Dirty: observation.Dirty, DirtyDigest: observation.DirtyDigest, HeadOID: observation.HeadOID, ProjectID: "PRJ-test", TreeOID: observation.TreeOID, WorktreeID: observation.WorktreeID}
}

func reconciliationRepositoryFixture(t *testing.T) string {
	t.Helper()
	repository := t.TempDir()
	for _, arguments := range [][]string{{"init", "-q"}, {"config", "user.email", "fid@example.test"}, {"config", "user.name", "fid test"}} {
		runReconciliationGit(t, repository, arguments...)
	}
	if err := os.WriteFile(filepath.Join(repository, "base.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runReconciliationGit(t, repository, "add", "base.txt")
	runReconciliationGit(t, repository, "commit", "-qm", "base")
	return repository
}

func runReconciliationGit(t *testing.T, repository string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = repository
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
}
