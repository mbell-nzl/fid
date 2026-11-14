package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
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
			// Use the real hand-off writer; both this and the state fixture append
			// Events, never SQL. Public rework and its V8 journey belong to 1.6f.
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
	// Model the future rework state only in this fixture. The product's public
	// transition refusal remains intact; no alternate production path is added.
	err := journal.WithLock(f.workspace, journal.LockOwner{Command: "second hand-off fixture"}, func(lock *journal.WriteLock) error {
		conn, err := reviewwork.OpenProjection(f.workspace)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		aggregate := journal.Aggregate{Type: "issue", ID: task}
		aggregate.Sequence, err = journal.NextAggregateSequence(filepath.Join(config.LedgerDir(f.workspace), "journal"), aggregate)
		if err != nil {
			return err
		}
		now := statusNow
		if f.clock != nil {
			now = f.clock()
		}
		e := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "issue.transitioned", CommandName: "issue.transition", OccurredAt: now, ActorV2: journal.EventActor{ActorID: "fixture", ActorType: "human", Producer: "test", ProvenanceRef: "review-rework-fixture"}, Aggregate: aggregate, Effects: []journal.Effect{
			{Op: "update", Table: "issues", PK: map[string]any{"id": task}, Values: map[string]any{"status": "in_progress"}, Preconditions: map[string]any{"status": "in_review"}},
			{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "review-rework-fixture"}, Values: map[string]any{"ref_type": "issue", "ref_id": task, "event": "status-changed", "actor": "fixture", "occurred_at": now.UTC().Format("2006-01-02T15:04:05Z")}},
		}}
		prepared, err := journal.PrepareEvent(lock, e)
		if err != nil {
			return err
		}
		e.Effects[1].Values["journal_event_id"] = prepared.EventID
		_, err = journal.CommitEvent(lock, conn, e)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:builder", "--reason", "Second hand-off fixture", "--model", "codex/gpt-6.1-sol"), "task.takeover")
}

func canonicalReviewFixture(t *testing.T) (*statusFixture, string) {
	t.Helper()
	f := newLifecycleFixture(t)
	now := time.Now().Add(time.Minute)
	f.clock = func() time.Time { return now }
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Review freshness", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "create"), "issue.create")
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
