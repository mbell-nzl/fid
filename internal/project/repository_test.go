package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
)

func TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths(t *testing.T) {
	t.Parallel()
	repository := initRepositoryFixture(t)
	path := "new\nfile.txt"
	if err := os.WriteFile(filepath.Join(repository, path), []byte("first"), 0o644); err != nil {
		t.Fatalf("write untracked file: %v", err)
	}

	observation, err := ObserveRepository(repository)
	if err != nil {
		t.Fatalf("ObserveRepository() error = %v", err)
	}
	if len(observation.Entries) != 1 || observation.Entries[0].Path != path || observation.Entries[0].Worktree == nil {
		t.Fatalf("ObserveRepository() entries = %#v, want the complete NUL-delimited untracked path", observation.Entries)
	}

	if err := os.WriteFile(filepath.Join(repository, path), []byte("second"), 0o644); err != nil {
		t.Fatalf("change untracked file: %v", err)
	}
	changed, err := ObserveRepository(repository)
	if err != nil {
		t.Fatalf("ObserveRepository() after change error = %v", err)
	}
	if changed.DirtyDigest == observation.DirtyDigest {
		t.Fatalf("ObserveRepository() digest = %v after changing untracked bytes, want a new digest", changed.DirtyDigest)
	}
}

func TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions(t *testing.T) {
	t.Parallel()
	repository := initRepositoryWithFiles(t, map[string]string{
		".gitignore":   "ignored.txt\n",
		"deleted.txt":  "deleted\n",
		"mode.txt":     "mode\n",
		"old-name.txt": "old\n",
		"staged.txt":   "staged-base\n",
		"unstaged.txt": "unstaged-base\n",
	})
	if err := os.Symlink("old-target", filepath.Join(repository, "link")); err != nil {
		t.Fatalf("create tracked symlink: %v", err)
	}
	runGit(t, repository, "add", "link")
	runGit(t, repository, "commit", "-qm", "link")
	writeFixtureFile(t, repository, "staged.txt", "staged-index\n")
	runGit(t, repository, "add", "staged.txt")
	writeFixtureFile(t, repository, "unstaged.txt", "unstaged-work\n")
	if err := os.Chmod(filepath.Join(repository, "mode.txt"), 0o755); err != nil {
		t.Fatalf("make executable: %v", err)
	}
	if err := os.Remove(filepath.Join(repository, "deleted.txt")); err != nil {
		t.Fatalf("delete tracked file: %v", err)
	}
	runGit(t, repository, "mv", "old-name.txt", "renamed.txt")
	if err := os.Remove(filepath.Join(repository, "link")); err != nil {
		t.Fatalf("remove tracked symlink: %v", err)
	}
	if err := os.Symlink("new-target", filepath.Join(repository, "link")); err != nil {
		t.Fatalf("replace tracked symlink: %v", err)
	}
	writeFixtureFile(t, repository, "untracked.txt", "untracked\n")
	writeFixtureFile(t, repository, "ignored.txt", "ignored\n")
	writeFixtureFile(t, repository, ".fid/operational", "excluded\n")

	observation, err := ObserveRepository(repository)
	if err != nil {
		t.Fatalf("ObserveRepository() error = %v", err)
	}
	gotPaths := make([]string, 0, len(observation.Entries))
	for _, entry := range observation.Entries {
		gotPaths = append(gotPaths, entry.Path)
	}
	wantPaths := []string{"deleted.txt", "link", "mode.txt", "old-name.txt", "renamed.txt", "staged.txt", "unstaged.txt", "untracked.txt"}
	if fmt.Sprint(gotPaths) != fmt.Sprint(wantPaths) {
		t.Fatalf("ObserveRepository() paths = %#v, want %#v", gotPaths, wantPaths)
	}
	if !observation.Dirty || observation.DirtyDigest == nil {
		t.Fatalf("ObserveRepository() = %#v, want a dirty canonical observation", observation)
	}
	const wantDigest = "sha256:78a4fa6c735d4f68e102172a4a0af0885638aec00891c2f6d79929b0d9af3cec"
	if *observation.DirtyDigest != wantDigest {
		t.Fatalf("ObserveRepository() dirty digest = %q, want %q", *observation.DirtyDigest, wantDigest)
	}
	if repeat, err := ObserveRepository(repository); err != nil || repeat.DirtyDigest == nil || *repeat.DirtyDigest != wantDigest {
		t.Fatalf("repeated ObserveRepository() = %#v, %v; want stable pinned digest %q", repeat, err, wantDigest)
	}
}

func initRepositoryFixture(t *testing.T) string {
	t.Helper()
	repository := t.TempDir()
	for _, arguments := range [][]string{{"init", "-q"}, {"config", "user.email", "fid@example.test"}, {"config", "user.name", "fid test"}} {
		command := exec.Command("git", arguments...)
		command.Dir = repository
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", arguments, err, output)
		}
	}
	if err := os.WriteFile(filepath.Join(repository, "tracked.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write tracked file: %v", err)
	}
	for _, arguments := range [][]string{{"add", "tracked.txt"}, {"commit", "-qm", "base"}} {
		command := exec.Command("git", arguments...)
		command.Dir = repository
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", arguments, err, output)
		}
	}
	return repository
}

func initRepositoryWithFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	repository := initRepositoryFixture(t)
	paths := make([]string, 0, len(files))
	for name, contents := range files {
		writeFixtureFile(t, repository, name, contents)
		paths = append(paths, name)
	}
	sort.Strings(paths)
	runGit(t, repository, append([]string{"add", "--"}, paths...)...)
	runGit(t, repository, "commit", "-qm", "fixture")
	return repository
}

func writeFixtureFile(t *testing.T, repository, name, contents string) {
	t.Helper()
	filename := filepath.Join(repository, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatalf("make parent for %q: %v", name, err)
	}
	if err := os.WriteFile(filename, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %q: %v", name, err)
	}
}

func runGit(t *testing.T, repository string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = repository
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
}
