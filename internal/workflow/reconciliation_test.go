package workflow

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/project"
)

func TestReconcileRepositoryRejectsMissingBaseEvenWithUnchangedHead(t *testing.T) {
	t.Parallel()
	root := reconciliationRepositoryFixture(t)
	current, err := project.ObserveRepository(root)
	if err != nil {
		t.Fatal(err)
	}
	missing := "0000000000000000000000000000000000000001"
	for _, base := range []*string{nil, &missing} {
		start := repositoryFromObservation(current)
		start.BaseOID = base
		_, err := ReconcileRepository(start, current, nil, nil)
		var conflict *RecoveryStateConflictError
		if !errors.As(err, &conflict) {
			t.Fatalf("missing base %v: %v, want repository conflict", base, err)
		}
	}
}

func TestReconcileRepositoryRequiresExactDirtyAdoption(t *testing.T) {
	t.Parallel()
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

// A .fid/ in a worktree is an ordinary path (D4 / F1, decided 6 November
// 2026), so its changes are in scope or outside it like any other. Paths in
// .git stay protected.
func TestReconcileRepositoryTreatsALeftoverFidDirectoryLikeAnyPath(t *testing.T) {
	t.Parallel()
	repository := reconciliationRepositoryFixture(t)
	start, err := project.ObserveRepository(repository)
	if err != nil {
		t.Fatalf("observe clean repository: %v", err)
	}
	baseline := repositoryFromObservation(start)
	for _, name := range []string{".fid/config.toml", ".fid/fid.sqlite"} {
		path := filepath.Join(repository, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("make parent of %s: %v", name, err)
		}
		if err := os.WriteFile(path, []byte("leftover\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	current, err := project.ObserveRepository(repository)
	if err != nil {
		t.Fatalf("observe dirty repository: %v", err)
	}

	scoped, err := ReconcileRepository(baseline, current, nil, []string{".fid/config.toml"})
	if err != nil {
		t.Fatalf("ReconcileRepository() error = %v", err)
	}
	got := map[string]string{}
	for _, change := range scoped.Paths {
		got[change.Path] = change.Classification
	}
	want := map[string]string{".fid/config.toml": "in-scope", ".fid/fid.sqlite": "outside-scope"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("classifications = %v, want %v", got, want)
	}

	adoptable, err := ReconcileRepository(baseline, current, nil, []string{".fid/config.toml", ".fid/fid.sqlite"})
	if err != nil {
		t.Fatalf("ReconcileRepository() with both in scope error = %v", err)
	}
	if err := ValidateReconciliationAdoption(adoptable, []string{".fid/config.toml", ".fid/fid.sqlite"}, adoptable.Digest); err != nil {
		t.Fatalf("ValidateReconciliationAdoption() of the leftover .fid paths error = %v", err)
	}

	// Git never lists a path in .git, so this entry is made by hand.
	withGit := current
	withGit.Entries = append(append([]project.Entry(nil), current.Entries...), project.Entry{Path: ".git/fid/config.toml"})
	guarded, err := ReconcileRepository(baseline, withGit, nil, []string{".fid/config.toml", ".fid/fid.sqlite", ".git/fid/config.toml"})
	if err != nil {
		t.Fatalf("ReconcileRepository() with a .git entry error = %v", err)
	}
	for _, change := range guarded.Paths {
		if wantProtected := change.Path == ".git/fid/config.toml"; (change.Classification == "protected") != wantProtected {
			t.Fatalf("%s classified %q, want protected only for the .git path", change.Path, change.Classification)
		}
	}
}

func TestActivatedCapabilitiesAcceptRunCheckpoint(t *testing.T) {
	t.Parallel()
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
