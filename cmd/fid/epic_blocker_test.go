package main

import (
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/epicwork"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/mbell-nzl/fid/internal/workrecord"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestEpicDiscoveredBlockerExceptionLimits(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                                                 string
		subtask, approved, stale, linked, same, closed, want bool
	}{{"blocker", false, true, false, true, true, false, true}, {"subtask acceptance", true, true, false, true, true, false, true}, {"no link", false, true, false, false, true, false, false}, {"another epic", false, true, false, true, false, false, false}, {"unapproved", false, false, false, true, true, false, false}, {"stale", false, true, true, true, true, false, false}, {"closed target", false, true, false, true, true, true, false}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			target := epicDoneChild(t, f, epic, "target")
			if tc.approved {
				approveTestEpic(t, f, epic)
			}
			parent := epic
			if !tc.same {
				other, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "epic", "--scheduling", "prioritised", "--title", "Other", "--actor", "human:mark", "--idempotency-key", "other"), "issue.create")
				parent = other["issue"].(map[string]any)["id"].(string)
				approveTestEpic(t, f, parent)
			}
			args := []string{"--json", "issue", "create", "--title", "discovered", "--actor", "agent:builder", "--idempotency-key", "discovered"}
			if tc.subtask {
				args = append(args, "--kind", "subtask", "--parent", target)
			} else {
				args = append(args, "--kind", "task", "--lane", "assured", "--epic", parent)
			}
			created, _ := lifecycleResult(t, f.run(args...), "issue.create")
			id := created["issue"].(map[string]any)["id"].(string)
			if tc.linked {
				lifecycleResult(t, f.run("--json", "link", "add", id, "blocks", target, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:builder", "--idempotency-key", "blocks"), "link.add")
			}
			if tc.closed {
				lifecycleResult(t, f.run("--json", "issue", "transition", target, "--to", "cancelled", "--actor", "human:mark", "--reason", "Not needed"), "issue.transition")
			}
			if tc.stale {
				_, err := epicwork.Change(epicwork.Options{WorkspaceRoot: f.workspace, MigrationsFS: fidmigrations.FS, EpicID: epic, Actor: "agent:builder", IdempotencyKey: "stale", Validation: epicplan.Validation{Claim: "New validation", Verifier: "review:code"}}, "add")
				if err != nil {
					t.Fatal(err)
				}
			}
			q, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = q.Close() }()
			unaccepted, err := workrecord.Unaccepted(q)
			if err != nil {
				t.Fatal(err)
			}
			if (unaccepted[id] == "") != tc.want {
				t.Fatalf("unaccepted=%v want exception=%v", unaccepted, tc.want)
			}
			next, _ := lifecycleResult(t, f.run("--json", "next", "--epic", parent), "next.get")
			if !tc.subtask && tc.want {
				if next["kind"] != "new" || next["task"].(map[string]any)["id"] != id {
					t.Fatalf("blocker not served: %v", next)
				}
			}
			if tc.subtask && strings.Contains(f.run("--json", "next", "--explain").stdout, `"id":"`+id+`"`) {
				t.Fatal("subtask entered scheduling")
			}
		})
	}
}

func TestEpicUnapprovedLapsedUnstartedHoldFreesWorktree(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	now := time.Now()
	f.clock = func() time.Time { return now }
	epicDoneChild(t, f, epic, "planned")
	approveTestEpic(t, f, epic)
	lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "human:mark"), "next.claim")
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Outside epic", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "orphan"), "issue.create")
	orphan := created["issue"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "Changed", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "change"), "epic.validation.add")
	live, _ := lifecycleResult(t, f.run("--json", "next", "--actor", "human:mark"), "next.get")
	if live["kind"] != "none" {
		t.Fatalf("live hold abandoned: %v", live)
	}
	now = now.Add(holds.Lifetime + time.Minute)
	lapsed, _ := lifecycleResult(t, f.run("--json", "next", "--actor", "human:mark"), "next.get")
	if lapsed["kind"] != "new" || lapsed["task"].(map[string]any)["id"] != orphan {
		t.Fatalf("expired unstarted hold occupied worktree: %v", lapsed)
	}
}

func TestEpicStaleApprovalAllowsStartedResume(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	task := epicDoneChild(t, f, epic, "started")
	approveTestEpic(t, f, epic)
	r, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "human:mark", "--epic", epic), "next.claim")
	hold := r["hold"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "start"), "task.start")
	lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "Changed", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "change"), "epic.validation.add")
	r, _ = lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "human:mark"), "next.claim")
	if r["kind"] != "resume" || r["task"].(map[string]any)["id"] != task {
		t.Fatalf("started work blocked: %v", r)
	}
}
