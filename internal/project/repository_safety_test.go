package project

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestObserveRepositoryRejectsUnsafeImages(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"assume-unchanged", "skip-worktree", "unmerged", "missing-head", "gitlink", "unreadable", "fifo", "escaping-parent"} {
		t.Run(kind, func(t *testing.T) {
			root := initRepositoryFixture(t)
			switch kind {
			case "assume-unchanged", "skip-worktree":
				runGit(t, root, "update-index", "--"+kind, "tracked.txt")
			case "unmerged":
				blob, err := exec.Command("git", "-C", root, "rev-parse", "HEAD:tracked.txt").Output()
				requireRepositorySuccess(t, err)
				cmd := exec.Command("git", "-C", root, "update-index", "--index-info")
				cmd.Stdin = strings.NewReader("0 " + strings.Repeat("0", 40) + "\ttracked.txt\n100644 " + strings.TrimSpace(string(blob)) + " 1\ttracked.txt\n")
				requireRepositorySuccess(t, cmd.Run())
			case "missing-head":
				runGit(t, root, "symbolic-ref", "HEAD", "refs/heads/missing")
			case "gitlink":
				head, err := exec.Command("git", "-C", root, "rev-parse", "HEAD").Output()
				requireRepositorySuccess(t, err)
				runGit(t, root, "update-index", "--add", "--cacheinfo", "160000,"+strings.TrimSpace(string(head))+",module")
			case "unreadable":
				requireRepositorySuccess(t, os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("changed"), 0o644))
				requireRepositorySuccess(t, os.Chmod(filepath.Join(root, "tracked.txt"), 0))
				t.Cleanup(func() { _ = os.Chmod(filepath.Join(root, "tracked.txt"), 0o644) })
			case "fifo":
				requireRepositorySuccess(t, os.Remove(filepath.Join(root, "tracked.txt")))
				requireRepositorySuccess(t, syscall.Mkfifo(filepath.Join(root, "tracked.txt"), 0o600))
			case "escaping-parent":
				writeFixtureFile(t, root, "dir/child", "inside")
				runGit(t, root, "add", "dir/child")
				runGit(t, root, "commit", "-qm", "nested")
				requireRepositorySuccess(t, os.RemoveAll(filepath.Join(root, "dir")))
				outside := t.TempDir()
				writeFixtureFile(t, outside, "child", "outside secret")
				requireRepositorySuccess(t, os.Symlink(outside, filepath.Join(root, "dir")))
			}
			_, err := ObserveRepository(root)
			if !errors.Is(err, ErrUnsafeRepository) {
				t.Fatalf("observation error = %v, want unsafe repository", err)
			}
		})
	}
}

// Git supplies real NUL-delimited output. The wrapper adds a file immediately
// after untracked enumeration, deterministically simulating an external writer
// during capture without replacing the observer or its Git parser.
func TestObserveRepositoryRejectsChangeDuringCapture(t *testing.T) {
	root := initRepositoryFixture(t)
	git, err := exec.LookPath("git")
	requireRepositorySuccess(t, err)
	bin := t.TempDir()
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
	script := "#!/bin/sh\n" + quote(git) + " \"$@\"\nstatus=$?\ncase \"$*\" in\n*'ls-files --others'*)\nif [ ! -e " + quote(filepath.Join(bin, "changed")) + " ]; then\nprintf changed > " + quote(filepath.Join(root, "arrived.txt")) + "\nprintf done > " + quote(filepath.Join(bin, "changed")) + "\nfi;;\nesac\nexit \"$status\"\n"
	requireRepositorySuccess(t, os.WriteFile(filepath.Join(bin, "git"), []byte(script), 0o755))
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	_, err = ObserveRepository(root)
	if !errors.Is(err, ErrUnsafeRepository) {
		t.Fatalf("observation error = %v, want unsafe repository", err)
	}
	_, statErr := os.Stat(filepath.Join(root, "arrived.txt"))
	requireRepositorySuccess(t, statErr)
}

func TestRepositoryPathGrammarRejectsEscapes(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"", ".", "../outside", "/absolute", "dir/../outside", "bad\x00path", string([]byte{0xff})} {
		if !errors.Is(validPath(name), ErrUnsafeRepository) {
			t.Fatalf("accepted invalid path %q", name)
		}
	}
	for _, name := range []string{"normal", "dir/file", "new\nfile"} {
		requireRepositorySuccess(t, validPath(name))
	}
}

func TestObserveRepositoryPinsStagedChangeRestoredInWorktree(t *testing.T) {
	t.Parallel()
	root := initRepositoryFixture(t)
	writeFixtureFile(t, root, "tracked.txt", "staged\n")
	runGit(t, root, "add", "tracked.txt")
	writeFixtureFile(t, root, "tracked.txt", "base\n")
	observation, err := ObserveRepository(root)
	requireRepositorySuccess(t, err)
	// Independently encoded canonical JSON with SHA-256 of the exact base and
	// staged bytes. The index image remains material even when HEAD == worktree.
	const want = "sha256:9926607bb8970a8199f50aadbeb618aa94def7c761d00fc86cf3ddf7b2f77b48"
	if !observation.Dirty || observation.DirtyDigest == nil || *observation.DirtyDigest != want || len(observation.Entries) != 1 {
		t.Fatalf("staged/restored observation = %#v, want pinned dirty digest %s", observation, want)
	}
}

func requireRepositorySuccess(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
