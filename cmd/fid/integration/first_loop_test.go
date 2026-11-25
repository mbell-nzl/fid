//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

type firstLoop struct {
	root, binary, task, state string
	env                       []string
}

func newFirstLoop(t *testing.T, scenario string) firstLoop {
	t.Helper()
	root, binary := handoffJourneyRepo(t, true)
	p22Write(t, root, "docs/briefs/test.md", "Approved test brief: implement the retry.\n")
	raw := readP21File(t, filepath.Join(root, ".fid.toml"))
	p22Write(t, root, ".fid.toml", strings.Replace(string(raw), "test_adequacy=[]", "test_adequacy=['**/*_test.go']", 1))
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Brief and review paths")
	closureRun(t, root, "git", "push", "-q", "origin", "HEAD:trunk")
	w := loopWorker{t, binary, root, "agent:sol-builder"}
	task := w.create("task", "First loop task", "", "--spec-md", "docs/briefs/test.md")
	w.call("issue", "ac", "add", task, "--ordinal", "1", "--text", "Retries are covered", "--actor", "human:planner", "--idempotency-key", "criterion")
	state := t.TempDir()
	bin := filepath.Join(state, "bin")
	if err := os.Mkdir(bin, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"agent", "gh", "codex", "claude"} {
		if err := os.WriteFile(filepath.Join(bin, name), readP21File(t, "testdata/first-loop-stub.py"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	backup := filepath.Join(state, "backups.git")
	closureRun(t, root, "git", "init", "--bare", "-q", backup)
	configHome := filepath.Join(state, "config")
	if err := config.SetBackupTarget(filepath.Join(configHome, "fid/config.toml"), backup); err != nil {
		t.Fatal(err)
	}
	env := append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"XDG_CONFIG_HOME="+configHome, "XDG_CACHE_HOME="+filepath.Join(state, "cache"),
		"LOOP_FID="+binary, "LOOP_GH_REPO=mbell-nzl/fid", "LOOP_BUILDER_COMMAND="+filepath.Join(bin, "agent"),
		"LOOP_REVIEWER_COMMAND="+filepath.Join(bin, "agent"), "FIRST_LOOP_STATE="+state, "FIRST_LOOP_SCENARIO="+scenario, "FIRST_LOOP_TASK="+task)
	return firstLoop{root, binary, task, state, env}
}

func (j firstLoop) run(t *testing.T, want int) string {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), "python3", filepath.Join(integrationRepositoryRoot(), "loops/loop.py"))
	cmd.Dir, cmd.Env = j.root, j.env
	out, err := cmd.CombinedOutput()
	if got := processExitCode(t, err); got != want {
		t.Fatalf("loop exit=%d, want %d:\n%s", got, want, out)
	}
	return string(out)
}

func (j firstLoop) calls(t *testing.T) []map[string]any {
	t.Helper()
	raw := readP21File(t, filepath.Join(j.state, "calls.jsonl"))
	var calls []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		var call map[string]any
		if err := json.Unmarshal([]byte(line), &call); err != nil {
			t.Fatal(err)
		}
		calls = append(calls, call)
	}
	return calls
}

func TestFirstLoopBuildReviewRework(t *testing.T) {
	t.Parallel()
	for _, scenario := range []string{"approve", "changes"} {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, scenario)
			j.run(t, 0)
			var builds, reviews, creates, ready int
			for _, call := range j.calls(t) {
				switch call["kind"] {
				case "build":
					builds++
					context := call["context"].(map[string]any)
					issue := context["task"].(map[string]any)
					if issue["id"] != j.task || issue["title"] != "First loop task" || issue["spec_md_path"] != "docs/briefs/test.md" || !strings.Contains(call["prompt"].(string), "Approved test brief") || !strings.Contains(call["prompt"].(string), "Retries are covered") {
						t.Fatalf("lost task brief: %v", call)
					}
					if builds == 2 && !strings.Contains(call["prompt"].(string), "Missing retry case") {
						t.Fatal("fixer was not given open findings")
					}
				case "review":
					reviews++
					if call["cwd"] == j.root || call["actor"] != "agent:sol-reviewer" || call["head"] != call["expected_head"] || call["detached"] != true {
						t.Fatalf("review did not use the detached review worktree at the hand-off head: %v", call)
					}
				case "create":
					creates++
					if !strings.Contains(call["body"].(string), "Stub completion report") || call["draft"] != true {
						t.Fatalf("PR body/draft = %v", call)
					}
				case "ready":
					ready++
					if call["recorded"] != true {
						t.Fatal("PR was marked ready before APPROVE was recorded")
					}
				}
			}
			want := 1
			if scenario == "changes" {
				want = 2
			}
			if builds != want || reviews != want || creates != 1 || ready != 1 {
				t.Fatalf("builds/reviews/creates/ready = %d/%d/%d/%d", builds, reviews, creates, ready)
			}
			if outcomes := backupJourneyOutcomes(t, j.root); len(outcomes) != 2*want {
				t.Fatalf("backup after each build/review = %v", outcomes)
			}
			w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
			listed := w.call("review", "list", j.task)
			rounds := p21Path(t, listed, "result", "rounds").([]any)
			if len(rounds) != want || rounds[want-1].(map[string]any)["verdict"] != "APPROVE" {
				t.Fatalf("recorded reviews = %v", rounds)
			}
			if scenario == "changes" && rounds[0].(map[string]any)["findings"].([]any)[0].(map[string]any)["open"] != false {
				t.Fatal("second review did not close the original major")
			}
		})
	}
}

func TestFirstLoopParkAndResume(t *testing.T) {
	t.Parallel()
	for _, scenario := range []string{"interrupt", "waiting", "dirty"} {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, scenario)
			want := 1
			if scenario == "waiting" {
				want = 0
			}
			out := j.run(t, want)
			w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
			if scenario == "waiting" {
				if !strings.Contains(out, "waiting-for-you") {
					t.Fatal(out)
				}
				w.call("answer", j.task, "Include both retry edges.", "--actor", "human:mark")
			} else if scenario == "dirty" {
				if !strings.Contains(out, "uncommitted") {
					t.Fatal(out)
				}
			} else if !strings.Contains(out, "agent") {
				t.Fatal(out)
			}
			j.run(t, 0)
			builds := 0
			for _, call := range j.calls(t) {
				if call["kind"] != "build" {
					continue
				}
				builds++
				if builds == 2 {
					context := call["context"].(map[string]any)
					if context["work"].(map[string]any)["kind"] != "resume" {
						t.Fatalf("lost interrupted work: %v", call)
					}
					if scenario == "waiting" && (!strings.Contains(call["prompt"].(string), "Include both retry edges.") || !strings.Contains(call["prompt"].(string), "Which retries?")) {
						t.Fatal("person's question and answer not passed to resumed agent")
					}
				}
			}
			if builds != 2 || len(backupJourneyOutcomes(t, j.root)) != 3 {
				t.Fatalf("builds/backups = %d/%v", builds, backupJourneyOutcomes(t, j.root))
			}
		})
	}
}
