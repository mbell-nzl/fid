package main

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/reviewpost"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestJudgeActorPeek(t *testing.T) {
	t.Parallel()
	f, task := judgeReady(t)
	f.workspace = gittest.AddWorktree(t, f.workspace)
	before := startEventCount(t, f.workspace)
	for _, tc := range []struct{ actor, kind string }{{"agent:reviewer", "none"}, {"agent:sol-judge", "judge"}} {
		got, _ := lifecycleResult(t, f.run("--json", "next", "--actor", tc.actor), "next.get")
		if got["kind"] != tc.kind {
			t.Fatalf("%s peek for %s: %v", tc.actor, task, got)
		}
	}
	if before != startEventCount(t, f.workspace) {
		t.Fatal("peek wrote a reservation")
	}
}

func TestPersonHoldBeforeReviews(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	lifecycleResult(t, f.run("--json", "hold", task, "--actor", "human:mark", "--reason", "Need an explanation"), "task.hold")
	got := f.run("status")
	if got.err != nil || !strings.Contains(got.stdout, "held by human:mark") || !strings.Contains(got.stdout, "Need an explanation") {
		t.Fatalf("hold not shown: %+v", got)
	}
}

func TestJudgeSealedSurfaces(t *testing.T) {
	t.Parallel()
	for _, surface := range []string{"list-json", "list-text", "show-json", "show-text", "status-json", "status-text", "record-json", "record-text", "delivery"} {
		t.Run(surface, func(t *testing.T) {
			t.Parallel()
			f, task := judgeReady(t)
			in := judgeInput("HOLD")
			var got commandExecutionResult
			if surface == "record-text" {
				raw, _ := json.Marshal(in)
				path := filepath.Join(t.TempDir(), "judge.json")
				writeBackupTestFile(t, path, string(raw))
				got = f.run("review", "record", task, "--input", path, "--actor", "agent:sol-judge", "--idempotency-key", "judge")
			} else {
				got = recordReviewTask(t, f, in, "judge", "agent:sol-judge", task)
			}
			if got.err != nil {
				t.Fatalf("record: %+v", got)
			}
			recorded, _ := lifecycleResult(t, f.run("--json", "review", "list", task), "review.list")
			round := recorded["rounds"].([]any)[1].(map[string]any)["id"].(string)
			if surface == "delivery" {
				conn, err := reviewwork.OpenProjection(f.workspace)
				if err != nil {
					t.Fatal(err)
				}
				pending, err := reviewpost.Pending(conn)
				_ = conn.Close()
				if err != nil || len(pending[round]) != 0 {
					t.Fatalf("sealed round queued for delivery: %v %v", pending, err)
				}
				calls := 0
				r, err := reviewpost.Round(t.Context(), f.workspace, fidmigrations.FS, round, func(context.Context, []string, []string) github.Result {
					calls++
					return github.Result{ExitCode: 1}
				})
				if err != nil || calls != 0 || len(r.Pending) != 0 || len(r.Posted) != 0 {
					t.Fatalf("sealed delivery=%+v %v calls=%d", r, err, calls)
				}
				return
			}
			if !strings.HasPrefix(surface, "record") {
				args := []string{}
				if strings.HasSuffix(surface, "json") {
					args = append(args, "--json")
				}
				switch {
				case strings.HasPrefix(surface, "list"):
					args = append(args, "review", "list", task)
				case strings.HasPrefix(surface, "show"):
					args = append(args, "show", task)
				default:
					args = append(args, "status")
				}
				got = f.run(args...)
			}
			if got.err != nil || strings.Contains(got.stdout, "sentinel") || strings.Contains(got.stdout, `"verdict":"HOLD"`) || strings.Contains(got.stdout, "Judge HOLD") || !strings.Contains(got.stdout, "sealed") {
				t.Fatalf("leaked %s: %+v", surface, got)
			}
		})
	}
}

func TestJudgeStatusAfterObservedMerge(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "standard", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "approved", "human:courier", task), "review.record")
	in := judgeInput("HOLD")
	delete(in, "lane_raise")
	lifecycleResult(t, recordReviewTask(t, f, in, "judge", "agent:sol-judge", task), "review.record")
	facts := doneFacts(t, f, task)
	facts.CI.Runs = []github.WorkflowRun{}
	got := runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher")
	if got.err == nil || !strings.Contains(got.stderr, "done-ci-waiting") {
		t.Fatalf("expected merge waiting for CI: %+v", got)
	}
	assertJudgeSurfaces(t, f, task, false, "HOLD", false)
	lifecycleResult(t, f.run("--json", "hold", task, "--actor", "human:mark", "--reason", "Inspect the merged change"), "task.hold")
	status := f.run("status")
	if status.err != nil || !strings.Contains(status.stdout, "held by human:mark") || !strings.Contains(status.stdout, "Inspect the merged change") {
		t.Fatalf("hold after observed merge not shown: %+v", status)
	}
	assertJudgeSurfaces(t, f, task, false, "HOLD", false)

}

func TestPersonHoldAfterApproval(t *testing.T) {
	t.Parallel()
	f, task := judgeReady(t)
	approveSuccess(t, approveForTest(f, task, task+"\n", true))
	lifecycleResult(t, f.run("--json", "hold", task, "--actor", "human:mark", "--reason", "Wait despite earlier approval"), "task.hold")
	got := f.run("status")
	if got.err != nil || !strings.Contains(got.stdout, "held by human:mark") {
		t.Fatalf("hold after approval not shown: %+v", got)
	}
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	h, err := reviewrecord.LatestHandoff(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	call, err := reviewrecord.FirstCall(conn, h.ID)
	if err != nil || call == nil || call.Verdict != "MERGE" {
		t.Fatalf("hold replaced first call: %+v %v", call, err)
	}
}

func TestPersonHoldRetryRejectsChangedProjection(t *testing.T) {
	t.Parallel()
	f, task := judgeReady(t)
	args := []string{"--json", "hold", task, "--actor", "human:mark", "--reason", "Wait for evidence", "--idempotency-key", "hold-retry"}
	lifecycleResult(t, f.run(args...), "task.hold")
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	var raw string
	if err := conn.QueryRow(`SELECT result_json FROM idempotency_results WHERE operation_scope LIKE 'task.hold:%'`).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	var stored map[string]any
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		t.Fatal(err)
	}
	stored["fields"].(map[string]any)["reason"] = "A different reason"
	changed, err := json.Marshal(stored)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`UPDATE idempotency_results SET result_json=? WHERE operation_scope LIKE 'task.hold:%'`, string(changed)); err != nil {
		t.Fatal(err)
	}
	_ = conn.Close()
	before := startEventCount(t, f.workspace)
	got := f.run(args...)
	if got.err == nil || startEventCount(t, f.workspace) != before {
		t.Fatalf("returned corrupt hold retry: %+v", got)
	}
}

func TestJudgeAssuredMergeWithSurvivingApproval(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "assured", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	for _, stream := range []string{"primary", "separate"} {
		in := reviewInput()
		in["stream"] = stream
		in["reviewer"].(map[string]any)["session_id"] = stream + "-session"
		lifecycleResult(t, recordReviewTask(t, f, in, stream, "human:courier", task), "review.record")
	}
	approveSuccess(t, approveForTest(f, task, task+"\n", true))
	// A new hand-off of the same change retains its fingerprint-bound approval.

	changed, _ := lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "evidence-changes", "human:courier", task), "review.record")
	finding := changed["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)["id"].(string)
	taken, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "human:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	hold = taken["hold"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "check", "--hold", hold, "--actor", "human:builder", "--idempotency-key", "same-check"), "check.execute")
	lifecycleResult(t, f.run("--json", "handoff", "--hold", hold, "--actor", "human:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "same-handoff"), "task.handoff")
	resolved := reviewInput()
	resolved["dispositions"] = []any{disposition(finding, "resolved")}
	lifecycleResult(t, recordReviewTask(t, f, resolved, "resolved-evidence", "human:courier", task), "review.record")
	in := judgeInput("HOLD")
	delete(in, "lane_raise")
	lifecycleResult(t, recordReviewTask(t, f, in, "new-handoff-judge", "agent:sol-judge", task), "review.record")
	assertJudgeSurfaces(t, f, task, true, "HOLD", false)
	facts := doneFacts(t, f, task)
	facts.CI.Runs = []github.WorkflowRun{}
	got := runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher")
	if got.err == nil || !strings.Contains(got.stderr, "done-ci-waiting") {
		t.Fatalf("expected merge waiting: %+v", got)
	}
	assertJudgeSurfaces(t, f, task, false, "HOLD", false)
}
