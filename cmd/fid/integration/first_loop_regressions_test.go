//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func firstLoopReplaceStub(t *testing.T, j firstLoop, name, old, replacement string) {
	t.Helper()
	path := filepath.Join(j.state, "bin", name)
	source := string(readP21File(t, path))
	if strings.Count(source, old) != 1 {
		t.Fatalf("stub %s must contain exactly one %q", name, old)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(source, old, replacement, 1)), 0o700); err != nil {
		t.Fatal(err)
	}
}

func firstLoopRoutine(t *testing.T, j firstLoop) {
	t.Helper()
	policy := string(readP21File(t, filepath.Join(j.root, ".fid.toml"))) + "\n[lanes.routine]\npaths=['docs/**']\nchecks=['test']\n"
	p22Write(t, j.root, ".fid.toml", policy)
	closureRun(t, j.root, "git", "add", ".")
	closureRun(t, j.root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Routine docs policy")
	closureRun(t, j.root, "git", "push", "-q", "origin", "HEAD:trunk")
	w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
	w.call("issue", "update", j.task, "--lane", "routine", "--actor", "human:planner", "--idempotency-key", "routine")
	firstLoopReplaceStub(t, j, "agent", `Path("change_test.go").write_text("package example\n")`, `Path("docs/change.md").write_text("Updated documentation.\n")`)
	firstLoopReplaceStub(t, j, "agent", `Path("change.txt").write_text(os.environ["LOOP_SESSION"] + "\n")`, "# Documentation only.")
}

func TestFirstLoopAgentChangesBranch(t *testing.T) {
	t.Parallel()
	for _, branch := range []string{"trunk", "main", "other-task", "detached"} {
		t.Run(branch, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, "approve")
			before, _ := closureRun(t, j.root, "git", "ls-remote", "origin", "refs/heads/*")
			args := "['git', 'switch', '-c', " + strconv.Quote(branch) + ", 'origin/trunk']"
			if branch == "main" {
				args = "['git', 'switch', 'main']"
			} else if branch == "detached" {
				args = "['git', 'checkout', '--detach', 'origin/trunk']"
			}
			firstLoopReplaceStub(t, j, "agent", "    merging = subprocess.run",
				"    subprocess.run("+args+", check=True)\n    merging = subprocess.run")
			// Check the actual bare remote before asserting the loop's refusal.
			cmd := exec.CommandContext(t.Context(), "python3", filepath.Join(integrationRepositoryRoot(), "loops/loop.py"))
			cmd.Dir, cmd.Env = j.root, j.env
			out, err := cmd.CombinedOutput()
			after, _ := closureRun(t, j.root, "git", "ls-remote", "origin", "refs/heads/*")
			if after != before {
				t.Fatalf("loop pushed after agent changed branch:\nbefore %s\nafter %s\n%s", before, after, out)
			}
			if err == nil || processExitCode(t, err) != 1 || !strings.Contains(string(out), "task branch") {
				t.Fatalf("missing branch refusal: %s", out)
			}
		})
	}
}

func TestFirstLoopFailedAgentDescendant(t *testing.T) {
	t.Parallel()
	j := newFirstLoop(t, "interrupt")
	descendant := `from pathlib import Path
import sys,time
state = Path(sys.argv[1])
(state / 'descendant-started').touch()
while not (state / 'release-descendant').exists():
    time.sleep(0.01)
(state / 'descendant-survived').write_text('still running after park')
`
	injection := "        import time\n" +
		"        child = subprocess.Popen([sys.executable, '-c', " + strconv.Quote(descendant) + ", str(state)], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)\n" +
		"        (state / 'descendant-pid').write_text(str(child.pid))\n" +
		"        while not (state / 'descendant-started').exists():\n            time.sleep(0.01)\n        sys.exit(42)"
	firstLoopReplaceStub(t, j, "agent", "        sys.exit(42)", injection)
	t.Cleanup(func() {
		raw, err := os.ReadFile(filepath.Join(j.state, "descendant-pid"))
		if err == nil {
			pid, err := strconv.Atoi(string(raw))
			if err == nil {
				process, _ := os.FindProcess(pid)
				_ = process.Kill()
			}
		}
	})
	j.run(t, 1)
	w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
	w.wantNext("", false, "resume", j.task)
	if err := os.WriteFile(filepath.Join(j.state, "release-descendant"), []byte("release"), 0o600); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(1500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(filepath.Join(j.state, "descendant-survived")); err == nil {
			t.Fatal("failed agent's descendant continued after the interrupted build became resumable")
		}
		time.Sleep(10 * time.Millisecond)
	}
	j.run(t, 0)
}

func TestFirstLoopDraftReadiness(t *testing.T) {
	t.Parallel()
	for _, scenario := range []string{"routine", "ready-retry"} {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, "approve")
			wantExit, wantReviews := 1, 1
			if scenario == "routine" {
				wantExit, wantReviews = 0, 0
				firstLoopRoutine(t, j)
			} else {
				firstLoopReplaceStub(t, j, "gh", `if Path(sys.argv[0]).name == "gh":`, `if Path(sys.argv[0]).name == "gh":
    if args[:2] == ["pr", "ready"] and not (state / "ready-failed-once").exists():
        (state / "ready-failed-once").touch()
        print("temporary ready failure", file=sys.stderr)
        sys.exit(1)`)
			}
			out := j.run(t, wantExit)
			if scenario == "ready-retry" && !strings.Contains(out, "temporary ready failure") {
				t.Fatal(out)
			}
			j.run(t, 0)
			var builds, reviews, ready int
			for _, call := range j.calls(t) {
				switch call["kind"] {
				case "build":
					builds++
				case "review":
					reviews++
				case "ready":
					ready++
					if wantReviews != 0 && call["recorded"] != true {
						t.Fatal("PR became ready before the review was recorded")
					}
				}
			}
			if builds != 1 || reviews != wantReviews || ready != 1 {
				t.Fatalf("builds/reviews/ready = %d/%d/%d, want 1/%d/1", builds, reviews, ready, wantReviews)
			}
		})
	}
}

func TestFirstLoopReviewerChangesCode(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"commit", "dirty"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, "approve")
			change := `    Path("change.txt").write_text("Reviewer changed the code.\n")` + "\n"
			if mode == "commit" {
				change += `    subprocess.run(["git", "add", "change.txt"], check=True)` + "\n" +
					`    subprocess.run(["git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "test: reviewer change"], check=True)` + "\n"
			}
			firstLoopReplaceStub(t, j, "agent", `    Path(os.environ["LOOP_REVIEW"]).write_text(json.dumps(review))`,
				change+`    Path(os.environ["LOOP_REVIEW"]).write_text(json.dumps(review))`)
			out := j.run(t, 1)
			if !strings.Contains(out, "review agent changed the head or left uncommitted changes") {
				t.Fatalf("missing review integrity refusal: %s", out)
			}
			w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
			if rounds := p21Path(t, w.call("review", "list", j.task), "result", "rounds").([]any); len(rounds) != 0 {
				t.Fatalf("altered code received a recorded review: %v", rounds)
			}
			for _, call := range j.calls(t) {
				if call["kind"] == "ready" || call["kind"] == "comment" {
					t.Fatalf("altered review reached GitHub: %v", call)
				}
			}
			if outcomes := backupJourneyOutcomes(t, j.root); len(outcomes) != 2 {
				t.Fatalf("backups after build and failed review = %v", outcomes)
			}
		})
	}
}
