//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFirstLoopRefusesDirtyWorktree(t *testing.T) {
	t.Parallel()
	for _, role := range []string{"new", "review"} {
		t.Run(role, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, "approve")
			root := j.root
			diagnostic, wantBackups := "new work needs a clean worktree", 1
			if role == "review" {
				root = j.root + "-review"
				closureRun(t, j.root, "git", "worktree", "add", "--detach", root, "HEAD")
				diagnostic, wantBackups = "review worktree has uncommitted changes", 2
			}
			p22Write(t, root, "unfinished.txt", "Do not let an agent adopt this.\n")
			if out := j.run(t, 1); !strings.Contains(out, diagnostic) {
				t.Fatalf("missing pre-agent cleanliness refusal: %s", out)
			}
			if role == "new" {
				if _, err := os.Stat(filepath.Join(j.state, "built")); !os.IsNotExist(err) {
					t.Fatalf("builder ran in a dirty worktree: %v", err)
				}
			} else {
				for _, call := range j.calls(t) {
					if call["kind"] == "review" || call["kind"] == "ready" {
						t.Fatalf("dirty review worktree reached an agent or readiness: %v", call)
					}
				}
			}
			if outcomes := backupJourneyOutcomes(t, j.root); len(outcomes) != wantBackups {
				t.Fatalf("backups after cleanliness refusal = %v, want %d", outcomes, wantBackups)
			}
		})
	}
}

func TestFirstLoopInterruptReceipt(t *testing.T) {
	t.Parallel()
	j := newFirstLoop(t, "interrupt")
	j.run(t, 1)
	w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
	shown := w.call("show", j.task)
	if mode := p21PathString(t, shown, "result", "kept", "mode"); mode != "interrupted" {
		t.Fatalf("failed agent was not parked as interrupted: %v", shown)
	}
	if reason := p21PathString(t, shown, "result", "kept", "reason"); !strings.Contains(reason, "agent failed") {
		t.Fatalf("interrupted park lost failure reason: %v", shown)
	}
	w.wantNext("", false, "resume", j.task)
	if outcomes := backupJourneyOutcomes(t, j.root); len(outcomes) != 1 {
		t.Fatalf("backup after interruption = %v", outcomes)
	}
}
