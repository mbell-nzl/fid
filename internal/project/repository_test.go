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
	// Until 1.1d5 fid excluded .fid/ paths; now this one counts like any other.
	writeFixtureFile(t, repository, ".fid/operational", "excluded\n")

	observation, err := ObserveRepository(repository)
	if err != nil {
		t.Fatalf("ObserveRepository() error = %v", err)
	}
	gotPaths := make([]string, 0, len(observation.Entries))
	for _, entry := range observation.Entries {
		gotPaths = append(gotPaths, entry.Path)
	}
	wantPaths := []string{".fid/operational", "deleted.txt", "link", "mode.txt", "old-name.txt", "renamed.txt", "staged.txt", "unstaged.txt", "untracked.txt"}
	if fmt.Sprint(gotPaths) != fmt.Sprint(wantPaths) {
		t.Fatalf("ObserveRepository() paths = %#v, want %#v", gotPaths, wantPaths)
	}
	if !observation.Dirty || observation.DirtyDigest == nil {
		t.Fatalf("ObserveRepository() = %#v, want a dirty canonical observation", observation)
	}
	const wantDigest = "sha256:1e6bbb5515d90e67323d890073e4338633671a9f09df0ae6e61ef964db206068"
	if *observation.DirtyDigest != wantDigest {
		t.Fatalf("ObserveRepository() dirty digest = %q, want %q", *observation.DirtyDigest, wantDigest)
	}
	if repeat, err := ObserveRepository(repository); err != nil || repeat.DirtyDigest == nil || *repeat.DirtyDigest != wantDigest {
		t.Fatalf("repeated ObserveRepository() = %#v, %v; want stable pinned digest %q", repeat, err, wantDigest)
	}
}

// The ledger is no longer in the worktree, so a .fid/ there is an ordinary
// path (D4 / F1, decided 6 November 2026): its untracked and edited files make
// the worktree dirty, and ignore rules apply to it as to any other path.
func TestObserveRepositoryCountsALeftoverFidDirectoryAsDirty(t *testing.T) {
	t.Parallel()
	repository := initRepositoryFixture(t)
	writeFixtureFile(t, repository, ".fid/config.toml", "leftover config\n")
	writeFixtureFile(t, repository, ".fid/journal/2026-10-01.jsonl", "leftover journal\n")
	writeFixtureFile(t, repository, ".fid/.gitignore", "fid.sqlite\n")
	writeFixtureFile(t, repository, ".fid/fid.sqlite", "ignored projection\n")

	observation, err := ObserveRepository(repository)
	if err != nil {
		t.Fatalf("ObserveRepository() error = %v", err)
	}
	gotPaths := make([]string, 0, len(observation.Entries))
	for _, entry := range observation.Entries {
		gotPaths = append(gotPaths, entry.Path)
	}
	wantPaths := []string{".fid/.gitignore", ".fid/config.toml", ".fid/journal/2026-10-01.jsonl"}
	if fmt.Sprint(gotPaths) != fmt.Sprint(wantPaths) {
		t.Fatalf("ObserveRepository() paths = %#v, want %#v", gotPaths, wantPaths)
	}
	if !observation.Dirty || observation.DirtyDigest == nil {
		t.Fatalf("ObserveRepository() = %#v, want a dirty observation", observation)
	}

	runGit(t, repository, "add", "--", ".fid")
	runGit(t, repository, "commit", "-qm", "track the leftover .fid")
	if clean, err := ObserveRepository(repository); err != nil || clean.Dirty || len(clean.Entries) != 0 {
		t.Fatalf("ObserveRepository() after commit = %#v, %v; want clean, with .fid/fid.sqlite still ignored", clean, err)
	}
	writeFixtureFile(t, repository, ".fid/config.toml", "edited config\n")
	edited, err := ObserveRepository(repository)
	if err != nil {
		t.Fatalf("ObserveRepository() after edit error = %v", err)
	}
	if len(edited.Entries) != 1 || edited.Entries[0].Path != ".fid/config.toml" || edited.Entries[0].Head == nil || edited.Entries[0].Worktree == nil || !edited.Dirty {
		t.Fatalf("ObserveRepository() entries = %#v, want the tracked .fid/config.toml edit", edited.Entries)
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
