package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/leaktest"
)

//nolint:paralleltest // Proves offline reads with a process-wide PATH containing git and no gh.
func TestMergedOfflineReadsGolden(t *testing.T) {
	f, task, hold := doneBuildFixture(t, "standard", false)
	now := f.clock()
	f.clock = func() time.Time { return now }
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	facts := doneFacts(t, f, task)
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	path := t.TempDir()
	if err := os.Symlink(git, filepath.Join(path, "git")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", path)
	if _, err := exec.LookPath("gh"); err == nil {
		t.Fatal("gh is on PATH")
	}
	f.setTarget()
	f.writeRecordAt(f.clock(), statusRecordLine{outcome: "pushed", ago: time.Minute})
	readMergedState(t, f, task, "", false)
	facts.CI.Runs = []github.WorkflowRun{}
	assertJSONFailureCode(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher").stderr, "task.done", "done-review-required")
	readMergedState(t, f, task, "waiting", true)
	facts.CI = ptrV8CI(completedCI(*facts.PR)) // A completed failure.
	assertJSONFailureCode(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher").stderr, "task.done", "done-review-required")
	readMergedState(t, f, task, "not-green", true)
	facts.CI.Runs[0].Attempt++
	facts.CI.Runs[0].Conclusion = "success"
	assertJSONFailureCode(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher").stderr, "task.done", "done-review-required")
	readMergedState(t, f, task, "green", false)
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "review", "human:courier", task), "review.record")
	lifecycleResult(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
	readMergedState(t, f, task, "", false)
}

func readMergedState(t *testing.T, f *statusFixture, task, ci string, golden bool) {
	t.Helper()
	before := startEventCount(t, f.workspace)
	for _, command := range []string{"status", "show"} {
		for _, jsonMode := range []bool{false, true} {
			args, ext := []string{command}, ".txt"
			if command == "show" {
				args = append(args, task)
			}
			if jsonMode {
				args, ext = append([]string{"--json"}, args...), ".json"
			}
			got := f.run(args...)
			if got.err != nil || got.stderr != "" {
				t.Fatalf("read=%#v", got)
			}
			wording := ""
			switch ci {
			case "waiting":
				wording = "merged, waiting for CI"
			case "not-green":
				wording = "merged, CI not green"
			}
			if jsonMode {
				if ci != "" && (!strings.Contains(got.stdout, `"ci":"`+ci+`"`) || !strings.Contains(got.stdout, `"merge_commit":"`+strings.Repeat("a", 40)+`"`)) {
					t.Fatalf("missing merge fields: %s", got.stdout)
				}
				if ci == "" && strings.Contains(got.stdout, `"merge":`) {
					t.Fatalf("unexpected merge state: %s", got.stdout)
				}
			} else {
				leaktest.AssertNoKernelWords(t, got.stdout)
				if wording != "" && !strings.Contains(got.stdout, wording) {
					t.Fatalf("missing %q: %s", wording, got.stdout)
				}
				if wording == "" && strings.Contains(got.stdout, "merged,") {
					t.Fatalf("unexpected merge wording: %s", got.stdout)
				}
			}
			if golden {
				out := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}T[^" ]+`).ReplaceAllString(got.stdout, "<time>")
				out = regexp.MustCompile(`\b[0-9]{2}:[0-9]{2}\b`).ReplaceAllString(out, "<clock>")
				out = regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`).ReplaceAllString(out, "<date>")
				out = regexp.MustCompile(`ART-[0-9a-f]{6}`).ReplaceAllString(out, "ART-<new>")
				out = regexp.MustCompile(`FID-[0-9a-f]{6}`).ReplaceAllString(out, "FID-<task>")
				handoffGolden(t, f, filepath.Join("..", "merged", ci+"-"+command+ext), out)
			}
		}
	}
	if startEventCount(t, f.workspace) != before {
		t.Fatal("offline read wrote Events")
	}
}

func TestMergedReadsIgnoreSupersededHandoffAndCancelledTask(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "routine", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	facts := doneFacts(t, f, task)
	facts.CI.Runs = []github.WorkflowRun{}
	assertJSONFailureCode(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher").stderr, "task.done", "done-ci-waiting")
	readMergedState(t, f, task, "waiting", false)
	in := changesInput("primary")
	lifecycleResult(t, recordReviewTask(t, f, in, "changes", "human:courier", task), "review.record")
	readMergedState(t, f, task, "", false)
	doneRehandoff(t, f, task, false)
	readMergedState(t, f, task, "", false)
	facts = doneFacts(t, f, task)
	facts.CI.Runs = []github.WorkflowRun{}
	assertJSONFailureCode(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher").stderr, "task.done", "done-open-major")
	readMergedState(t, f, task, "waiting", false)
	lifecycleResult(t, f.run("--json", "issue", "transition", task, "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed"), "issue.transition")
	readMergedState(t, f, task, "", false)
}
