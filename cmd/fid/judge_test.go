package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func judgeInput(verdict string) map[string]any {
	in := reviewInput()
	in["stream"], in["profiles"], in["verdict"] = "judge", []string{"merge"}, verdict
	in["reviewer"] = map[string]any{"actor": "agent:sol-judge", "model": "openai/gpt-6.1-sol", "session_id": "judge-session"}
	in["reason"] = "sealed sentinel: evidence supports this call"
	in["purpose"] = "sealed purpose sentinel"
	in["verification"] = "sealed verification sentinel"
	in["provenance"] = map[string]any{"prompt": "sealed prompt sentinel", "command": "sealed command sentinel", "output": "sealed output sentinel"}
	in["lane_raise"] = map[string]any{"lane": "assured", "reason": "sealed lane sentinel"}
	return in
}

func judgeReady(t *testing.T) (*statusFixture, string) {
	t.Helper()
	f, task := canonicalReviewFixture(t)
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "approved", "human:courier", task), "review.record")
	return f, task
}

func TestJudgeSealingHoldAndNoEffects(t *testing.T) {
	t.Parallel()
	f, task := judgeReady(t)
	recorded := recordReviewTask(t, f, judgeInput("HOLD"), "judge", "agent:sol-judge", task)
	lifecycleResult(t, recorded, "review.record")
	live, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	var liveState, liveLane string
	err = live.QueryRow(`SELECT status,lane FROM issues WHERE id=?`, task).Scan(&liveState, &liveLane)
	_ = live.Close()
	if err != nil || liveState != "in_review" || liveLane != "standard" {
		t.Fatalf("judge changed live task before rebuild: %s %s %v", liveState, liveLane, err)
	}

	if strings.Contains(recorded.stdout, "sentinel") || !strings.Contains(recorded.stdout, `"sealed":true`) {
		t.Fatalf("record leaked: %s", recorded.stdout)
	}
	assertJudgeSurfaces(t, f, task, true, "HOLD", false)
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
	assertJudgeSurfaces(t, f, task, true, "HOLD", false)
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	h, err := reviewrecord.LatestHandoff(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	fp, err := fingerprint.Compute(f.workspace, h.Base, h.Head)
	if err != nil {
		t.Fatal(err)
	}
	s, err := doneeval.GatherReviews(f.workspace, conn, h.TaskID, &fp)
	if err != nil || !doneeval.ReviewsDone(s) {
		t.Fatalf("judge blocked reviews: %v %+v", err, s)
	}
	var state, lane string
	if err := conn.QueryRow(`SELECT status,lane FROM issues WHERE id=?`, h.TaskID).Scan(&state, &lane); err != nil {
		t.Fatal(err)
	}
	if state != "in_review" || lane != "standard" {
		t.Fatalf("judge changed task: %s %s", state, lane)
	}
	held := f.run("--json", "hold", h.TaskID, "--actor", "human:mark", "--reason", "Wait for a closer look", "--idempotency-key", "hold-call")
	lifecycleResult(t, held, "task.hold")
	before := startEventCount(t, f.workspace)
	retry := f.run("--json", "hold", h.TaskID, "--actor", "human:mark", "--reason", "Wait for a closer look", "--idempotency-key", "hold-call")
	if retry.err != nil || held.stdout != retry.stdout || before != startEventCount(t, f.workspace) {
		t.Fatalf("hold retry: %+v", retry)
	}
	assertJudgeSurfaces(t, f, task, false, "HOLD", true)
	listed, _ := lifecycleResult(t, f.run("--json", "review", "list", task), "review.list")
	r := listed["rounds"].([]any)[1].(map[string]any)
	if raise, ok := r["lane_raise"].(map[string]any); !ok || raise["lane"] != "assured" || raise["reason"] != "sealed lane sentinel" {
		t.Fatalf("lost shadow lane recommendation: %v", r)
	}
	// Rebuild replays both canonical records and the first-call projection.
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
	assertJudgeSurfaces(t, f, task, false, "HOLD", true)
}

func assertJudgeSurfaces(t *testing.T, f *statusFixture, task string, sealed bool, verdict string, agree bool) {
	t.Helper()
	for _, args := range [][]string{{"review", "list", task}, {"show", task}, {"status"}} {
		for _, isJSON := range []bool{true, false} {
			runArgs := append([]string{}, args...)
			if isJSON {
				runArgs = append([]string{"--json"}, runArgs...)
			}
			got := f.run(runArgs...)
			if got.err != nil {
				t.Fatalf("%v: %+v", runArgs, got)
			}
			if sealed {
				if strings.Contains(got.stdout, "sentinel") || strings.Contains(got.stdout, `"verdict":"`+verdict+`"`) || strings.Contains(got.stdout, "Judge "+verdict) || !strings.Contains(got.stdout, "sealed") {
					t.Fatalf("%v not sealed: %s", runArgs, got.stdout)
				}
			} else if !strings.Contains(got.stdout, verdict) || !strings.Contains(got.stdout, "agree") {
				t.Fatalf("%v no comparison: %s", runArgs, got.stdout)
			}
		}
	}
	list, _ := lifecycleResult(t, f.run("--json", "review", "list", task), "review.list")
	rounds := list["rounds"].([]any)
	r := rounds[len(rounds)-1].(map[string]any)
	if (r["sealed"] == true) != sealed {
		t.Fatalf("sealed=%v", r)
	}
	if !sealed && r["agrees"] != agree {
		raw, _ := json.Marshal(r)
		t.Fatalf("comparison=%s", raw)
	}
}

func TestJudgeHoldRefusals(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, actor, reason string
		ready               bool
	}{{"agent", "agent:builder", "wait", true}, {"service", "service:bot", "wait", true}, {"blank", "human:mark", " ", true}, {"not review", "human:mark", "wait", false}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "[]")
			task := "FID-9f13aa"
			if tc.ready {
				f, task = judgeReady(t)
			}
			before := startEventCount(t, f.workspace)
			got := f.run("--json", "hold", task, "--actor", tc.actor, "--reason", tc.reason, "--idempotency-key", "refuse")
			code := "invalid-argument"
			if strings.HasPrefix(tc.actor, "agent:") || strings.HasPrefix(tc.actor, "service:") {
				code = "authority-denied"
			}
			assertJSONFailureCode(t, got.stderr, "task.hold", code)
			if got.err == nil || before != startEventCount(t, f.workspace) {
				t.Fatalf("hold accepted: %+v", got)
			}
		})
	}
}

func TestJudgeSelectionAndNewHandoff(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithLane(t, "assured")
	builder := f.workspace
	reviewer := gittest.AddWorktree(t, builder)
	f.workspace = reviewer
	if got := selectionNext(t, f, ""); got["kind"] != "review" {
		t.Fatalf("judge before reviews: %v", got)
	}
	in := reviewInput()
	lifecycleResult(t, recordReviewTask(t, f, in, "primary", "human:courier", task), "review.record")
	if got := selectionNext(t, f, ""); got["kind"] != "review" {
		t.Fatalf("judge before separate: %v", got)
	}
	in["stream"] = "separate"
	in["reviewer"].(map[string]any)["session_id"] = "separate"
	lifecycleResult(t, recordReviewTask(t, f, in, "separate", "human:courier", task), "review.record")
	for _, actor := range []string{"agent:builder", "agent:reviewer"} {
		if got := selectionNext(t, f, actor); got["kind"] != "none" {
			t.Fatalf("non-independent judge: %v", got)
		}
	}
	got := selectionNext(t, f, "agent:sol-judge")
	if got["kind"] != "judge" || got["task"].(map[string]any)["id"] != task {
		t.Fatalf("no judge: %v", got)
	}
	hold := got["hold"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "heartbeat", "--hold", hold, "--actor", "agent:sol-judge"), "hold.heartbeat")
	other := gittest.AddWorktree(t, builder)
	f.workspace = other
	if got := selectionNext(t, f, "agent:other-judge"); got["kind"] != "none" {
		t.Fatalf("double judge: %v", got)
	}
	f.workspace = reviewer
	j := judgeInput("HOLD")
	delete(j, "lane_raise")
	lifecycleResult(t, recordReviewTask(t, f, j, "judge", "agent:sol-judge", task), "review.record")
	if got := selectionNext(t, f, "agent:sol-judge"); got["kind"] != "none" {
		t.Fatalf("judge twice: %v", got)
	}
	// A subsequent ordinary CHANGES round, then a new hand-off, earns new work.
	in = changesInput("primary")
	fields, _ := lifecycleResult(t, recordReviewTask(t, f, in, "changes-after-judge", "human:courier", task), "review.record")
	finding := fields["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)["id"].(string)
	conn, err := reviewwork.OpenProjection(builder)
	if err != nil {
		t.Fatal(err)
	}
	limit, err := doneeval.ReadReworkLimit(conn, task)
	_ = conn.Close()
	if err != nil || limit == nil || limit.Count != 1 {
		t.Fatalf("judge consumed rework cap: %+v %v", limit, err)
	}
	f.workspace = builder
	next := selectionNext(t, f, "agent:builder")
	if next["kind"] != "rework" {
		t.Fatalf("judge consumed cap: %v", next)
	}
	hold = next["hold"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "check", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "new-check"), "check.execute")
	lifecycleResult(t, f.run("--json", "handoff", "--hold", hold, "--actor", "agent:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "new-handoff"), "task.handoff")
	f.workspace = reviewer
	in = reviewInput()
	in["dispositions"] = []any{disposition(finding, "resolved")}
	lifecycleResult(t, recordReviewTask(t, f, in, "fixed", "human:courier", task), "review.record")
	if got := selectionNext(t, f, "agent:sol-judge"); got["kind"] != "judge" {
		t.Fatalf("new hand-off no judge: %v", got)
	}
	setIssueStatus(t, builder, task, "done")
	if got := selectionNext(t, f, ""); got["kind"] == "judge" {
		t.Fatalf("judge after done: %v", got)
	}
}

func TestJudgeRevealTriggersAndFirstCall(t *testing.T) {
	t.Parallel()
	for _, trigger := range []string{"approve", "merge", "hold-then-approve", "hold-then-merge", "approve-before-judge"} {
		t.Run(trigger, func(t *testing.T) {
			t.Parallel()
			lane := "assured"
			if strings.Contains(trigger, "merge") {
				lane = "standard"
			}
			f, task, hold := doneBuildFixture(t, lane, false)
			lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
			lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "approved", "human:courier", task), "review.record")
			if lane == "assured" {
				in := reviewInput()
				in["stream"] = "separate"
				in["reviewer"].(map[string]any)["session_id"] = "separate-reveal"
				lifecycleResult(t, recordReviewTask(t, f, in, "separate-approved", "human:courier", task), "review.record")
			}
			if trigger == "approve-before-judge" {
				approveSuccess(t, approveForTest(f, task, task+"\n", true))
			}
			in := judgeInput("MERGE")
			delete(in, "lane_raise")
			lifecycleResult(t, recordReviewTask(t, f, in, "judge", "agent:sol-judge", task), "review.record")
			if strings.HasPrefix(trigger, "hold-then") {
				lifecycleResult(t, f.run("--json", "hold", task, "--actor", "human:mark", "--reason", "Wait", "--idempotency-key", "hold"), "task.hold")
			}
			if strings.Contains(trigger, "merge") {
				facts := doneFacts(t, f, task)
				facts.CI.Runs[0].Conclusion = "success"
				lifecycleResult(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
			} else if trigger != "approve-before-judge" {
				approveSuccess(t, approveForTest(f, task, task+"\n", true))
			}
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
						t.Fatal(err)
					}
				}
				listed, _ := lifecycleResult(t, f.run("--json", "review", "list", task), "review.list")
				rounds := listed["rounds"].([]any)
				r := rounds[len(rounds)-1].(map[string]any)
				if r["sealed"] == true || r["verdict"] != "MERGE" || r["agrees"] != (!strings.HasPrefix(trigger, "hold-then")) {
					t.Fatalf("trigger %s lost first call: %v", trigger, r)
				}
			}
		})
	}
}

func TestJudgeAdmissionRefusals(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		prepare func(*testing.T, *statusFixture, string, map[string]any)
	}{
		{"before reviews", func(t *testing.T, f *statusFixture, task string, _ map[string]any) {
			startSQL(t, f, `DELETE FROM review_rounds WHERE task_id='`+task+`'`)
		}},
		{"builder", func(_ *testing.T, _ *statusFixture, _ string, in map[string]any) {
			in["reviewer"].(map[string]any)["actor"] = "agent:builder"
		}},
		{"reviewer", func(_ *testing.T, _ *statusFixture, _ string, in map[string]any) {
			in["reviewer"].(map[string]any)["actor"] = "agent:reviewer"
		}},
		{"review session", func(_ *testing.T, _ *statusFixture, _ string, in map[string]any) {
			in["reviewer"].(map[string]any)["session_id"] = "review-session-1"
		}},
		{"duplicate handoff", func(t *testing.T, f *statusFixture, task string, in map[string]any) {
			lifecycleResult(t, recordReviewTask(t, f, in, "first-judge", "agent:sol-judge", task), "review.record")
			in["reviewer"].(map[string]any)["session_id"] = "second-judge-session"
		}},
		{"done", func(t *testing.T, f *statusFixture, task string, _ map[string]any) {
			setIssueStatus(t, f.workspace, task, "done")
		}},
		{"findings", func(_ *testing.T, _ *statusFixture, _ string, in map[string]any) {
			in["findings"] = []any{reviewFinding("major")}
		}},
		{"wrong profile", func(_ *testing.T, _ *statusFixture, _ string, in map[string]any) { in["profiles"] = []string{"code"} }},
		{"review verdict", func(_ *testing.T, _ *statusFixture, _ string, in map[string]any) { in["verdict"] = "CHANGES" }},
		{"reason", func(_ *testing.T, _ *statusFixture, _ string, in map[string]any) { in["reason"] = " " }},
		{"lane not raised", func(_ *testing.T, _ *statusFixture, _ string, in map[string]any) {
			in["lane_raise"] = map[string]any{"lane": "standard", "reason": "No raise"}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, task := judgeReady(t)
			in := judgeInput("HOLD")
			tc.prepare(t, f, task, in)
			before := startEventCount(t, f.workspace)
			got := recordReviewTask(t, f, in, "refused-judge", "agent:sol-judge", task)
			if got.err == nil || before != startEventCount(t, f.workspace) {
				t.Fatalf("judge admitted: %+v", got)
			}
		})
	}
}
