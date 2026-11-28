package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/next"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func selectionNext(t *testing.T, f *statusFixture, actor string) map[string]any {
	t.Helper()
	args := []string{"--json", "next", "--explain"}
	if actor != "" {
		args = append(args, "--claim", "--actor", actor, "--model", "claude/opus-5.5")
	}
	operation := "next.get"
	if actor != "" {
		operation = "next.claim"
	}
	got, _ := lifecycleResult(t, f.run(args...), operation)
	return got
}

func TestReviewWorkSelection(t *testing.T) {
	t.Parallel()
	for _, lane := range []string{"routine", "standard", "assured"} {
		t.Run(lane, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithLane(t, lane)
			builder := f.workspace
			f.workspace = gittest.AddWorktree(t, builder)
			got := selectionNext(t, f, "")
			if lane == "routine" {
				if item := selectionItem(t, f, "merge"); item == nil {
					t.Fatal("Routine without required reviews needs merge item")
				}
				if got["kind"] != "judge" {
					t.Fatalf("Routine with satisfied reviews needs shadow judge: %v", got)
				}
				return
			}
			assertReviewWork(t, got, task, "primary")
			text := f.run("next")
			if text.err != nil || !strings.Contains(text.stdout, "Review primary") || !strings.Contains(text.stdout, "code") || !strings.Contains(text.stdout, "https://github.com/mbell-nzl/fid/pull/63") {
				t.Fatalf("review text omits its instructions: %#v", text)
			}
			// The hand-off ended this actor's build hold. History still excludes it.
			skipped := selectionNext(t, f, "agent:builder")
			if skipped["kind"] != "none" || !strings.Contains(fmt.Sprint(skipped["queue"]), "built-by-actor") {
				t.Fatalf("ended builder got review work: %v", skipped)
			}
			in := reviewInput()
			lifecycleResult(t, recordReviewTask(t, f, in, "primary", "human:courier", task), "review.record")
			got = selectionNext(t, f, "")
			if lane == "assured" {
				assertReviewWork(t, got, task, "separate")
				in["stream"] = "separate"
				in["reviewer"].(map[string]any)["session_id"] = "separate-session"
				lifecycleResult(t, recordReviewTask(t, f, in, "separate", "human:courier", task), "review.record")
				got = selectionNext(t, f, "")
			}
			if got["kind"] != "judge" {
				t.Fatalf("approved task needs shadow judge: %v", got)
			}
			f.workspace = builder
			if got := selectionNext(t, f, ""); got["kind"] != "none" || got["waiting"].(map[string]any)["kind"] != "review" {
				t.Fatalf("builder reservation changed: %v", got)
			}
		})
	}
}

func assertReviewWork(t *testing.T, got map[string]any, task, stream string) {
	t.Helper()
	review, ok := got["review"].(map[string]any)
	if got["kind"] != "review" || !ok || review["stream"] != stream || got["task"].(map[string]any)["id"] != task {
		t.Fatalf("want %s review for %s: %v", stream, task, got)
	}
	h, ok := review["handoff"].(map[string]any)
	if !ok || h["id"] == "" || h["pr"] == "" || len(fmt.Sprint(h["head"])) != 40 || len(fmt.Sprint(h["base"])) != 40 || fmt.Sprint(review["required_profiles"]) != "[code]" {
		t.Fatalf("incomplete review work: %v", review)
	}
}

func TestReviewClaimLifecycle(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	builder := f.workspace
	reviewer := gittest.AddWorktree(t, builder)
	other := gittest.AddWorktree(t, builder)
	f.workspace = reviewer
	got := selectionNext(t, f, "agent:reviewer")
	assertReviewWork(t, got, task, "primary")
	hold := got["hold"].(map[string]any)["id"].(string)
	before := startEventCount(t, builder)
	if again := selectionNext(t, f, "agent:reviewer"); again["hold"].(map[string]any)["id"] != hold || startEventCount(t, builder) != before {
		t.Fatalf("claim repeat changed reservation: %v", again)
	}
	lifecycleResult(t, f.run("--json", "heartbeat", "--hold", hold, "--actor", "agent:reviewer"), "hold.heartbeat")
	for _, actor := range []string{"agent:other", "agent:reviewer"} {
		if actor == "agent:reviewer" {
			f.workspace = other
		}
		if got := f.run("--json", "heartbeat", "--hold", hold, "--actor", actor); got.err == nil || !strings.Contains(got.stderr, "hold-not-current") {
			t.Fatalf("foreign session renewed review: %#v", got)
		}
	}
	f.workspace = other
	if got := selectionNext(t, f, "agent:second"); got["kind"] != "none" || !strings.Contains(fmt.Sprint(got["queue"]), "held-elsewhere") {
		t.Fatalf("second reviewer got reserved work: %v", got)
	}
	// A lapsed review frees the stream, unlike a started build hold.
	f.clock = func() time.Time { return time.Now().Add(2 * time.Hour) }
	f.workspace = reviewer
	if expired := f.run("--json", "heartbeat", "--hold", hold, "--actor", "agent:reviewer"); expired.err == nil || !strings.Contains(expired.stderr, "hold-not-current") {
		t.Fatalf("lapsed review renewed: %#v", expired)
	}
	f.workspace = other
	got = selectionNext(t, f, "agent:second")
	assertReviewWork(t, got, task, "primary")
	second := got["hold"].(map[string]any)["id"].(string)
	in := changesInput("primary")
	in["reviewer"].(map[string]any)["actor"] = "agent:second"
	lifecycleResult(t, recordReviewTask(t, f, in, "changes", "human:courier", task), "review.record")
	if heartbeat := f.run("--json", "heartbeat", "--hold", second, "--actor", "agent:second"); heartbeat.err == nil || !strings.Contains(heartbeat.stderr, "hold-not-current") {
		t.Fatalf("recording did not end claim: %#v", heartbeat)
	}
	if got := selectionNext(t, f, ""); got["kind"] != "none" {
		t.Fatalf("review served while rework pending: %v", got)
	}
	f.workspace = builder
	if got := selectionNext(t, f, ""); got["kind"] != "rework" {
		t.Fatalf("claim lost builder reservation: %v", got)
	}
}

func TestConcurrentReviewClaimsAndLiveAdmission(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	roots := []string{gittest.AddWorktree(t, f.workspace), gittest.AddWorktree(t, f.workspace)}
	results := make(chan next.Result, 2)
	errors := make(chan error, 2)
	for i, root := range roots {
		go func() {
			r, err := next.Claim(next.Options{WorkspaceRoot: root, MigrationsFS: fidmigrations.FS, Clock: f.clock, Actor: fmt.Sprintf("agent:reviewer-%d", i), Model: "claude/opus-5.5"})
			results <- r
			errors <- err
		}()
	}
	count := 0
	for range roots {
		if err := <-errors; err != nil {
			t.Fatal(err)
		}
		if r := <-results; r.Kind == "review" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("%d reviewers took the same review", count)
	}
	e := mustReadCommandEventByType(t, f.workspace, "review.claimed")
	e.EventID, e.Aggregate.ID = "", "RCL-ffffff"
	e.Effects[0].PK["id"] = e.Aggregate.ID
	e.Effects[0].Values["actor"] = "agent:forged"
	before := snapshotWorkspacePersistence(t, f.workspace)
	err := journal.WithLock(f.workspace, journal.LockOwner{Command: "test malformed review"}, func(lock *journal.WriteLock) error {
		conn, err := reviewwork.OpenProjection(f.workspace)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		_, err = journal.CommitEvent(lock, conn, e)
		return err
	})
	if err == nil || !strings.Contains(err.Error(), "review reservation differs") || !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
		t.Fatalf("live malformed append changed state or failed elsewhere: %v", err)
	}
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
	f.workspace = roots[0]
	got := selectionNext(t, f, "")
	if got["kind"] == "review" {
		assertReviewWork(t, got, task, "primary")
	} else if got["kind"] != "none" {
		t.Fatalf("replay lost reservation: %v", got)
	}
}

func TestUnavailableReviewKeepsItsReviewerWorktree(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithLane(t, "assured")
	builder := f.workspace
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
	f.workspace = gittest.AddWorktree(t, builder)
	selectionNext(t, f, "agent:reviewer")
	lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "New work", "--actor", "human:mark", "--idempotency-key", "other"), "issue.create")
	head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
	objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
	if err := os.Remove(filepath.Join(objects, head[:2], head[2:])); err != nil {
		t.Fatal(err)
	}
	got := selectionNext(t, f, "agent:reviewer")
	if got["kind"] != "none" || got["waiting"].(map[string]any)["task_id"] != task {
		t.Fatalf("unavailable review let a worktree start new work: %v", got)
	}
}

func TestReviewSelectionPreservesMissingBuilderWorktree(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureAt(t, "standard", true)
	builder := f.workspace
	f.workspace = gittest.AddWorktree(t, builder)
	if err := os.RemoveAll(builder); err != nil {
		t.Fatal(err)
	}
	if item := selectionItem(t, f, "worktree-missing"); item == nil || item["subject"].(map[string]any)["id"] != task {
		t.Fatalf("review selection hid builder reservation: %v", item)
	}
}
