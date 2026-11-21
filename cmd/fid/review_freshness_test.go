package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func TestReviewFreshnessAfterJournalHandoff(t *testing.T) {
	t.Parallel()
	for _, rebase := range []bool{false, true} {
		t.Run(map[bool]string{false: "edit", true: "clean rebase"}[rebase], func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixture(t)
			in := reviewInput()
			in["verdict"] = "CHANGES"
			in["findings"] = []any{reviewFinding("major")}
			lifecycleResult(t, recordReviewTask(t, f, in, "first", "human:courier", task), "review.record")
			first, err := reviewwork.List(f.workspace, task)
			if err != nil || len(first) != 1 || !first[0].Current {
				t.Fatalf("first=%v/%v", first, err)
			}
			if rebase {
				runGitForTest(t, f.workspace, "checkout", "-q", "trunk")
				writeBackupTestFile(t, filepath.Join(f.workspace, "unrelated.txt"), "base change\n")
				runGitForTest(t, f.workspace, "add", ".")
				runGitForTest(t, f.workspace, "commit", "-m", "Move base")
				runGitForTest(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", "HEAD")
				runGitForTest(t, f.workspace, "checkout", "-q", "task")
				runGitForTest(t, f.workspace, "rebase", "trunk")
			} else {
				writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "edited review change\n")
				runGitForTest(t, f.workspace, "add", ".")
				runGitForTest(t, f.workspace, "commit", "-m", "Edit task")
			}
			prepareSecondReviewHandoff(t, f, task)
			// Both work return and the second hand-off use public commands.
			reviewHandoffTask(t, f, "second-handoff", task)
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
						t.Fatal(err)
					}
				}
				got, err := reviewwork.List(f.workspace, task)
				if err != nil || len(got) != 1 || got[0].Current != rebase || got[0].HandoffID != first[0].HandoffID || got[0].Binding != first[0].Binding {
					t.Fatalf("round after second hand-off=%v/%v, want current=%t", got, err, rebase)
				}
			}
			before := reviewProjection(t, f)
			refused := recordReviewTask(t, f, reviewInput(), "after-edit", "human:courier", task)
			if refused.err == nil || !strings.Contains(refused.stderr, "verdict disagrees") || reviewProjection(t, f) != before {
				t.Fatalf("old major after hand-off=%#v", refused)
			}
		})
	}
}

func prepareSecondReviewHandoff(t *testing.T, f *statusFixture, task string) {
	t.Helper()
	taken, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if taken["kind"] != "rework" || taken["task"].(map[string]any)["id"] != task {
		t.Fatalf("second hand-off requires public rework: %v", taken)
	}
}

func canonicalReviewFixture(t *testing.T) (*statusFixture, string) {
	t.Helper()
	return canonicalReviewFixtureWithLane(t, "standard")
}

func canonicalReviewFixtureWithLane(t *testing.T, lane string) (*statusFixture, string) {
	t.Helper()
	return canonicalReviewFixtureAt(t, lane, false)
}

func canonicalReviewFixtureAt(t *testing.T, lane string, linked bool) (*statusFixture, string) {
	t.Helper()
	now := time.Now().Add(time.Minute)
	return canonicalReviewFixtureWithClock(t, lane, linked, func() time.Time { return now })
}

func canonicalReviewFixtureWithClock(t *testing.T, lane string, linked bool, clock func() time.Time) (*statusFixture, string) {
	t.Helper()
	f := newLifecycleFixture(t)
	if linked {
		f.workspace = gittest.AddWorktree(t, f.workspace)
	}
	f.clock = clock
	if lane == "routine" {
		trustedShowPolicy(t, f, func(raw string) string {
			return strings.Replace(raw, "[lanes.routine]", "[lanes.routine]\npaths=['**/*.md']", 1)
		})
	}
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Review freshness", "--lane", lane, "--actor", "human:mark", "--idempotency-key", "create"), "issue.create")
	task := created["issue"].(map[string]any)["id"].(string)
	taken, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	hold := taken["hold"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start"), "task.start")
	runGitForTest(t, f.workspace, "checkout", "-b", "task")
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "review this change\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Task change")
	reviewHandoffTask(t, f, "first-handoff", task)

	return f, task
}

func TestReviewReadKnownBindingWithoutObjects(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "round", "human:courier", task), "review.record")
	head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
	objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
	path := filepath.Join(objects, head[:2], head[2:])
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(path, raw, 0o444); err != nil {
			t.Error(err)
		}
	})
	got, err := reviewwork.List(f.workspace, task)
	if err != nil || len(got) != 1 || !got[0].Current {
		t.Fatalf("known review binding=%v/%v", got, err)
	}
	refused := recordReviewTask(t, f, reviewInput(), "new-round", "human:courier", task)
	if refused.err == nil || !strings.Contains(refused.stderr, "fingerprint missing-object") {
		t.Fatalf("new round without objects=%#v", refused)
	}
}

func TestReviewListRebuildsMissingProjection(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "round", "human:courier", task), "review.record")
	before := startEventCount(t, f.workspace)
	for _, suffix := range []string{"", "-wal", "-shm"} {
		if err := os.Remove(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName+suffix)); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
	fields, _ := lifecycleResult(t, f.run("--json", "review", "list", task), "review.list")
	rounds := fields["rounds"].([]any)
	if len(rounds) != 1 || rounds[0].(map[string]any)["current"] != true || startEventCount(t, f.workspace) != before {
		t.Fatalf("rebuilt read=%v", fields)
	}
}
