package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestKeptWorktreeHealth(t *testing.T) {
	t.Parallel()
	for _, test := range []struct{ name, mode string }{
		{"removed", "removed"},
		{"replaced", "replaced"},
		{"unverified", "unverified"},
		{"short", "unverified"},
		{"unverified-with-a-longer-temporary-path", "unverified"},
	} {
		t.Run(test.name, func(t *testing.T) {
			mode := test.mode
			f := newLifecycleFixture(t)
			nextWork(f)
			main := f.workspace
			other := fixedWidthHealthWorktree(t, main)
			f.workspace = other
			f.hold("HLD-aaaaaa", "FID-71a0c2", "", time.Minute)
			parkWaiting(t, f, "HLD-aaaaaa", waitingQuestion, "wait")
			f.workspace = main
			if mode == "unverified" {
				if err := os.WriteFile(filepath.Join(other, ".git"), []byte("gitdir: /missing/admin-directory\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := os.RemoveAll(other); err != nil {
					t.Fatal(err)
				}
				if mode == "replaced" {
					replacement := gittest.NewWorkspace(t)
					if err := os.Rename(replacement, other); err != nil {
						t.Fatal(err)
					}
				}
			}
			status := f.run("--json", "status")
			result, _ := lifecycleResult(t, status, "project.status")
			want := "worktree-missing"
			if mode == "unverified" {
				want = "worktree-unverified"
			}
			collection := result["needs_you"].([]any)
			if mode == "unverified" {
				collection = result["warnings"].([]any)
			}
			found := false
			for _, item := range collection {
				if item.(map[string]any)["kind"] == want {
					found = true
				}
			}
			if !found {
				t.Fatalf("%s missing from %v", want, result)
			}
			shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-71a0c2"), "issue.show")
			if shown["kept"] == nil {
				t.Fatal("worktree check released kept task")
			}
			// Normalise only the other worktree's generated location.
			waitingGolden(t, f, mode+"-status.json", strings.ReplaceAll(status.stdout, other, "/w/old"))
			human := f.run("status")
			if human.err != nil {
				t.Fatal(human.err)
			}
			waitingGolden(t, f, mode+"-status.txt", strings.ReplaceAll(human.stdout, other, "/w/old"))
		})
	}
}

// Move the real Git worktree before recording its hold: replacing its path
// after rendering cannot undo line breaks caused by a variable-width path.
func fixedWidthHealthWorktree(t *testing.T, repository string) string {
	t.Helper()
	original := gittest.AddWorktree(t, repository)
	parent := gittest.ResolvedPath(t, filepath.Dir(repository))
	const pathWidth = 61
	nameWidth := pathWidth - len(parent) - 1
	if nameWidth < 1 {
		t.Fatalf("health fixture parent %q is too long for a %d-column path", parent, pathWidth)
	}
	worktree := filepath.Join(parent, strings.Repeat("w", nameWidth))
	if output, err := exec.Command("git", "-C", repository, "worktree", "move", original, worktree).CombinedOutput(); err != nil {
		t.Fatalf("move health worktree: %v: %s", err, output)
	}
	if err := gittest.CheckRepository(worktree); err != nil {
		t.Fatal(err)
	}
	return worktree
}
