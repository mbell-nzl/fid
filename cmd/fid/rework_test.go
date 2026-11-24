package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestReworkOrderingAndPersonTakeover(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Unblocking work", "--lane", "standard", "--priority", "0", "--actor", "human:mark", "--idempotency-key", "create-unblocking"), "issue.create")
	newTask := created["issue"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "link", "add", newTask, "blocks", task, "--from-type", "issue", "--to-type", "issue", "--actor", "human:mark", "--idempotency-key", "block"), "link.add")
	main := f.workspace
	f.workspace = gittest.AddWorktree(t, main)
	other, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:other", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if other["kind"] != "new" || other["task"].(map[string]any)["id"] != newTask {
		t.Fatalf("another worktree took rework: %v", other)
	}
	before := startEventCount(t, main)
	resumed, event := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:other", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if resumed["kind"] != "resume" || resumed["task"].(map[string]any)["id"] != newTask || event != nil || startEventCount(t, main) != before {
		t.Fatalf("resume must precede rework: %v", resumed)
	}
	f.workspace = main
	read, _ := lifecycleResult(t, f.run("--json", "next"), "next.get")
	if read["kind"] != "rework" || read["task"].(map[string]any)["id"] != task {
		t.Fatalf("rework must precede unblocking work: %v", read)
	}
	// A person may explicitly move the kept task to a third worktree.
	f.workspace = gittest.AddWorktree(t, main)
	taken, _ := lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:mark", "--reason", "Continue the fix here", "--idempotency-key", "take-rework"), "task.takeover")
	if taken["previous_hold_id"] == nil {
		t.Fatal(taken)
	}
	if state, _ := taskState(t, f, task); state != "in_progress" {
		t.Fatal(state)
	}
	lifecycleResult(t, f.run("--json", "heartbeat", "--hold", taken["hold_id"].(string), "--actor", "human:mark"), "hold.heartbeat")
	f.workspace = main
	old, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if old["kind"] != "none" {
		t.Fatalf("old worktree got task after takeover: %v", old)
	}
}

func TestMissingReworkWorktreeNeedsYou(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureAt(t, "standard", true)
	lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
	missing := f.workspace
	f.workspace = gittest.AddWorktree(t, missing)
	if err := os.RemoveAll(missing); err != nil {
		t.Fatal(err)
	}
	status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	found := false
	for _, raw := range status["needs_you"].([]any) {
		item := raw.(map[string]any)
		if item["kind"] == "worktree-missing" && strings.Contains(fmt.Sprint(item), task) && strings.Contains(fmt.Sprint(item), "fid takeover") {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing rework worktree not actionable: %v", status)
	}
	next, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:other", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if next["kind"] != "none" {
		t.Fatalf("missing worktree released rework: %v", next)
	}
}

func changesInput(stream string) map[string]any {
	in := reviewInput()
	in["stream"], in["verdict"] = stream, "CHANGES"
	in["reviewer"].(map[string]any)["session_id"] = "changes-" + stream
	in["findings"] = []any{reviewFinding("major")}
	return in
}

func taskState(t *testing.T, f *statusFixture, task string) (string, float64) {
	t.Helper()
	shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	i := shown["issue"].(map[string]any)
	return i["status"].(string), i["revision"].(float64)
}

func TestPublicReworkJourney(t *testing.T) {
	t.Parallel()
	for _, stream := range []string{"primary", "separate"} {
		t.Run(stream, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithLane(t, "assured")
			_, revision := taskState(t, f, task)
			in := changesInput(stream)
			first := recordReviewTask(t, f, in, "changes", "human:courier", task)
			fields, _ := lifecycleResult(t, first, "review.record")
			finding := fields["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)["id"].(string)
			if state, rev := taskState(t, f, task); state != "in_progress" || rev != revision+1 {
				t.Fatalf("CHANGES left task %s revision %v; want in_progress revision %v", state, rev, revision+1)
			}
			shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
			kept := shown["kept"].(map[string]any)
			if kept["mode"] != "rework" || shown["working"] != nil || kept["here"] != true {
				t.Fatalf("rework reservation=%v", shown)
			}
			for _, args := range [][]string{{"next"}, {"show", task}, {"status"}} {
				text := f.run(args...)
				assertHumanCommand(t, text, 0)
				if text.err != nil || !strings.Contains(strings.ToLower(text.stdout), "rework") {
					t.Fatalf("%v did not describe rework: %#v", args, text)
				}
			}
			read, _ := lifecycleResult(t, f.run("--json", "next"), "next.get")
			if read["kind"] != "rework" || read["hold"] != nil || read["task"].(map[string]any)["id"] != task {
				t.Fatalf("read rework=%v", read)
			}
			// A second stream asking for changes records history without another
			// transition or changing the reservation's original rework time.
			other := map[string]string{"primary": "separate", "separate": "primary"}[stream]
			second := changesInput(other)
			second["findings"] = []any{}
			lifecycleResult(t, recordReviewTask(t, f, second, "second-stream", "human:courier", task), "review.record")
			if state, rev := taskState(t, f, task); state != "in_progress" || rev != revision+1 {
				t.Fatalf("second CHANGES moved task: %s/%v", state, rev)
			}
			taken, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:fixer", "--model", "claude/opus-5.5"), "next.claim")
			hold := taken["hold"].(map[string]any)
			if taken["kind"] != "rework" || hold["replaces"] != kept["hold_id"] || hold["model"] != "claude/opus-5.5" {
				t.Fatalf("taken rework=%v", taken)
			}
			count := startEventCount(t, f.workspace)
			again, event := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:fixer", "--model", "codex/gpt-6.1-sol"), "next.claim")
			if again["kind"] != "resume" || again["hold"].(map[string]any)["id"] != hold["id"] || event != nil || startEventCount(t, f.workspace) != count {
				t.Fatalf("repeated next=%v/%v", again, event)
			}
			if retry := recordReviewTask(t, f, in, "changes", "human:courier", task); retry.err != nil || retry.stdout != first.stdout || startEventCount(t, f.workspace) != count {
				t.Fatalf("late CHANGES retry=%#v", retry)
			}
			writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "Fixed the reviewed fault\n")
			runGitForTest(t, f.workspace, "add", ".")
			runGitForTest(t, f.workspace, "commit", "-m", "Fix review finding")
			lifecycleResult(t, f.run("--json", "check", "--hold", hold["id"].(string), "--actor", "agent:fixer", "--idempotency-key", "fix-check"), "check.execute")
			reviewHandoffTask(t, f, "fixed-handoff", task)
			approve := reviewInput()
			approve["stream"] = stream
			approve["reviewer"].(map[string]any)["session_id"] = "fixed-" + stream
			approve["dispositions"] = []any{disposition(finding, "resolved")}
			lifecycleResult(t, recordReviewTask(t, f, approve, "fixed-review", "human:courier", task), "review.record")
			if state, rev := taskState(t, f, task); state != "in_review" || rev != revision+2 {
				t.Fatalf("APPROVE moved task: %s/%v", state, rev)
			}
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
						t.Fatal(err)
					}
				}
				rounds, err := reviewwork.List(f.workspace, task)
				if err != nil || len(rounds) != 3 || rounds[0].Current || rounds[1].Current || !rounds[2].Current || rounds[0].Findings[0].Open || rounds[2].OtherModel {
					t.Fatalf("reworked reviews=%v/%v", rounds, err)
				}
			}
			// A rework builder is also disqualified from reviewing its task.
			approve["reviewer"].(map[string]any)["actor"] = "agent:fixer"
			approve["dispositions"] = []any{}
			got := recordReviewTask(t, f, approve, "self-review", "human:courier", task)
			if got.err == nil || !strings.Contains(got.stderr, "builder cannot") {
				t.Fatalf("self review=%v %s", got.err, got.stderr)
			}
		})
	}
}
