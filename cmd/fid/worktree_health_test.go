package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestKeptWorktreeHealth(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"removed", "replaced", "unverified"} {
		t.Run(mode, func(t *testing.T) {
			f := newLifecycleFixture(t)
			nextWork(f)
			main := f.workspace
			other := gittest.AddWorktree(t, main)
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
