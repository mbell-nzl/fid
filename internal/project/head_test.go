package project

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestObserveHeadCleanlinessAndHiddenChanges(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"clean", "ignored", "untracked", "unstaged", "staged", "assume-unchanged", "skip-worktree"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			root := initRepositoryFixture(t)
			switch kind {
			case "ignored":
				writeFixtureFile(t, root, ".git/info/exclude", "scratch\n")
				writeFixtureFile(t, root, "scratch", "ignored\n")
			case "untracked":
				writeFixtureFile(t, root, "scratch", "new\n")
			case "unstaged", "staged":
				writeFixtureFile(t, root, "tracked.txt", "edited\n")
				if kind == "staged" {
					runGit(t, root, "add", "tracked.txt")
				}
			case "assume-unchanged", "skip-worktree":
				runGit(t, root, "update-index", "--"+kind, "tracked.txt")
				writeFixtureFile(t, root, "tracked.txt", "hidden\n")
			}
			observed, err := ObserveHead(root)
			if kind == "assume-unchanged" || kind == "skip-worktree" {
				var unavailable *RepositoryUnavailableError
				if !errors.As(err, &unavailable) || unavailable.Reason != "tracked-flags" {
					t.Fatalf("hidden change=%#v/%v", observed, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			dirty := kind != "clean" && kind != "ignored"
			if observed.Dirty != dirty {
				t.Fatalf("observation=%#v", observed)
			}
			if !dirty && (observed.HeadOID == "" || observed.TreeOID == "") {
				t.Fatalf("missing committed head=%#v", observed)
			}
		})
	}
}

//nolint:paralleltest // Checks the environment inherited by git via process-global PATH.
func TestObserveHeadKeepsEveryGitChildLocalAndReadOnly(t *testing.T) {
	root := initRepositoryFixture(t)
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	script := "#!/bin/sh\n[ \"$GIT_NO_LAZY_FETCH\" = 1 ] && [ \"$GIT_OPTIONAL_LOCKS\" = 0 ] || exit 91\nexec '" + strings.ReplaceAll(git, "'", "'\\''") + "' \"$@\"\n"
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	observed, err := ObserveHead(root)
	if err != nil || observed.Dirty || observed.HeadOID == "" || observed.TreeOID == "" {
		t.Fatalf("local head=%#v/%v", observed, err)
	}
}
