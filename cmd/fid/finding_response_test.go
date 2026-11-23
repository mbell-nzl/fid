package main

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func TestReviewDeliveryInvalidArgumentDetails(t *testing.T) {
	t.Parallel()
	f, ids, _ := closureFixture(t)
	for _, tc := range []struct {
		args              []string
		argument, allowed string
	}{
		{[]string{"finding", "answer", ids[0], "--actor", "agent:stranger", "--text", "Fixed"}, "finding", "an open finding and a builder or person"},
		{[]string{"finding", "defer", ids[2], "--actor", "agent:stranger", "--title", "Later"}, "finding", "an open finding and a builder or person"},
		{[]string{"review", "post", "RND-unknown"}, "round", "a recorded review round"},
	} {
		t.Run(strings.Join(tc.args[:2], "."), func(t *testing.T) {
			got := f.run(append([]string{"--json"}, tc.args...)...)
			var reply struct {
				Error struct {
					Details struct {
						Argument string   `json:"argument"`
						Allowed  []string `json:"allowed"`
					} `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(got.stderr), &reply); err != nil {
				t.Fatal(err)
			}
			if got.err == nil || reply.Error.Details.Argument != tc.argument || strings.Join(reply.Error.Details.Allowed, ",") != tc.allowed {
				t.Fatalf("wrong command guidance: %s", got.stderr)
			}
		})
	}
}

func TestFindingDeferralEpicAndLane(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	f.clock = time.Now
	create := func(kind, title string, args ...string) string {
		fields, _ := lifecycleResult(t, f.run(append([]string{"--json", "issue", "create", "--kind", kind, "--title", title, "--actor", "human:mark", "--idempotency-key", title}, args...)...), "issue.create")
		return fields["issue"].(map[string]any)["id"].(string)
	}
	epic := create("epic", "Source epic", "--scheduling", "prioritised")
	other := create("epic", "Other epic", "--scheduling", "prioritised")
	task := create("task", "Reviewed task", "--epic", epic, "--lane", "standard")
	taken, _ := lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:mark", "--reason", "Build"), "task.takeover")
	hold := taken["hold_id"].(string)
	lifecycleResult(t, f.run("--json", "start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "start"), "task.start")
	writeBackupTestFile(t, filepath.Join(f.workspace, "change.md"), "Change\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Change")
	lifecycleResult(t, f.run("--json", "handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "handoff"), "task.handoff")
	in := reviewInput()
	in["findings"] = []any{reviewFinding("minor"), reviewFinding("nit")}
	round, _ := lifecycleResult(t, recordReviewTask(t, f, in, "approved", "human:courier", task), "review.record")
	findings := round["round"].(map[string]any)["findings"].([]any)
	for i, wantEpic := range []string{epic, other} {
		id := findings[i].(map[string]any)["id"].(string)
		if i == 0 {
			for _, args := range [][]string{{"--lane", "routine"}, {"--epic", task}, {"--epic", "missing"}} {
				before := reviewProjection(t, f)
				got := f.run(append([]string{"--json", "finding", "defer", id, "--actor", "human:mark", "--title", "Later"}, args...)...)
				if got.err == nil || reviewProjection(t, f) != before {
					t.Fatal("invalid deferral options admitted")
				}
			}
		}
		args := []string{"--json", "finding", "defer", id, "--actor", "human:mark", "--title", "Later"}
		lane := "standard"
		if i == 1 {
			args = append(args, "--epic", other, "--lane", "assured")
			lane = "assured"
		}
		result, _ := lifecycleResult(t, f.run(args...), "finding.defer")
		deferred := result["response"].(map[string]any)["deferred_task"].(map[string]any)
		if deferred["epic_id"] != wantEpic || deferred["lane"] != lane {
			t.Fatalf("deferral options lost: %v", deferred)
		}
	}
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	var state string
	if err := conn.QueryRow(`SELECT status FROM issues WHERE id=?`, task).Scan(&state); err != nil || state != "in_review" {
		t.Fatalf("deferral reopened approved task: %s/%v", state, err)
	}
}

func TestFindingResponseGolden(t *testing.T) {
	t.Parallel()
	f, ids, task := closureFixture(t)
	normalize := func(s string) string {
		s = f.placeholders(s)
		s = regexp.MustCompile(`EVT-\d{8}-\d+`).ReplaceAllString(s, "EVT-<new>")
		s = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`).ReplaceAllString(s, "<TIME>")
		s = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`).ReplaceAllString(s, "<DATE>")
		for _, prefix := range []string{"RND", "FRS", "HOF", "FID"} {
			s = regexp.MustCompile(prefix+`-[0-9a-f]{6}`).ReplaceAllString(s, prefix+"-<new>")
		}
		return regexp.MustCompile(`\b(?:[0-9a-f]{64}|[0-9a-f]{40})\b`).ReplaceAllString(s, "<fingerprint-or-commit>")
	}
	for _, jsonMode := range []bool{false, true} {
		ext, prefix := "txt", []string{}
		if jsonMode {
			ext, prefix = "json", []string{"--json"}
		}
		for _, verb := range []string{"answer", "defer"} {
			index, flag, value := 0, "--text", "Regression added in the next change"
			if verb == "defer" {
				index, flag, value = 2, "--title", "Follow up the finding"
				if jsonMode {
					index = 3
				}
			}
			got := f.run(append(append([]string{}, prefix...), "finding", verb, ids[index], "--actor", "agent:builder", flag, value)...)
			if got.err != nil {
				t.Fatal(got)
			}
			assertGolden(t, filepath.Join("finding", verb+"."+ext), normalize(got.stdout))
		}
		listed := f.run(append(append([]string{}, prefix...), "review", "list", task)...)
		if listed.err != nil {
			t.Fatal(listed)
		}
		assertGolden(t, filepath.Join("finding", "responses-list."+ext), normalize(listed.stdout))
	}
}

func TestFindingAnswersInReviewerInput(t *testing.T) {
	t.Parallel()
	f, ids, task := closureFixture(t)
	lifecycleResult(t, f.run("--json", "finding", "answer", ids[0], "--actor", "agent:builder", "--text", "Regression added in this fix"), "finding.answer")
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "Fix\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Fix")
	prepareSecondReviewHandoff(t, f, task)
	reviewHandoffTask(t, f, "second", task)
	reviewer := gittest.AddWorktree(t, f.workspace)
	got := executeCommandForTest(reviewer, "--json", "next", "--claim", "--actor", "agent:next-reviewer", "--model", "claude/opus-5.5")
	if got.err != nil || !strings.Contains(got.stdout, "Regression added in this fix") || !strings.Contains(got.stdout, ids[0]) {
		t.Fatalf("reviewer input lacks answer: %#v", got)
	}
}

func TestFindingDeferralKeepsPathFloor(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	in := reviewInput()
	finding := reviewFinding("minor")
	finding["location"] = map[string]any{"file": ".fid.toml", "line": 1}
	in["findings"] = []any{finding}
	round, _ := lifecycleResult(t, recordReviewTask(t, f, in, "path-floor", "human:courier", task), "review.record")
	id := round["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)["id"].(string)
	result, _ := lifecycleResult(t, f.run("--json", "finding", "defer", id, "--actor", "agent:builder", "--title", "Fix policy", "--lane", "standard"), "finding.defer")
	if result["response"].(map[string]any)["deferred_task"].(map[string]any)["lane"] != "assured" {
		t.Fatal("finding's hard floor was lowered")
	}
}

func TestFindingAnswersStayOpen(t *testing.T) {
	t.Parallel()
	f, ids, task := closureFixture(t)
	for _, actor := range []string{"agent:builder", "human:mark"} {
		for _, index := range []int{0, 1, 2, 3} {
			lifecycleResult(t, f.run("--json", "finding", "answer", ids[index], "--actor", actor, "--text", "Fixed in the next change"), "finding.answer")
		}
	}
	for _, rebuild := range []bool{false, true} {
		if rebuild {
			if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
				t.Fatal(err)
			}
		}
		views, err := reviewwork.List(f.workspace, task)
		if err != nil || !views[0].Findings[0].Open {
			t.Fatalf("answer closed finding: %v / %v", views, err)
		}
		for _, finding := range views[0].Findings {
			if !finding.Open {
				t.Fatalf("answer closed %s", finding.ID)
			}
			if len(finding.Responses) != 2 || finding.Responses[0].Actor != "agent:builder" || finding.Responses[1].Actor != "human:mark" {
				t.Fatalf("answer order=%v", finding.Responses)
			}
		}
		for _, args := range [][]string{{"review", "list", task}, {"--json", "review", "list", task}} {
			got := f.run(args...)
			if got.err != nil || !strings.Contains(got.stdout, "Fixed in the next change") || !strings.Contains(got.stdout, "human:mark") {
				t.Fatalf("answer missing: %#v", got)
			}
		}
	}
	before := reviewProjection(t, f)
	if got := recordReviewTask(t, f, reviewInput(), "still-open", "human:courier", task); got.err == nil || reviewProjection(t, f) != before {
		t.Fatal("builder answer bypassed major gate")
	}
}

func TestFindingResponseRefusals(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, verb, actor, flag, value string
		index                          int
		closed                         bool
	}{
		{"nonbuilder", "answer", "agent:stranger", "--text", "Fixed", 0, false},
		{"service", "answer", "service:ci", "--text", "Fixed", 0, false},
		{"blank", "answer", "human:mark", "--text", " ", 0, false},
		{"closed", "answer", "human:mark", "--text", "Fixed", 0, true},
		{"major", "defer", "human:mark", "--title", "Later", 0, false},
		{"defer stranger", "defer", "agent:stranger", "--title", "Later", 2, false},
		{"defer closed", "defer", "human:mark", "--title", "Later", 2, true},
		{"blank title", "defer", "human:mark", "--title", " ", 2, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, ids, _ := closureFixture(t)
			if tc.closed {
				lifecycleResult(t, closeFinding(f, ids[tc.index], "human:mark", "Accepted", "close"), "finding.close")
			}
			before, count := reviewProjection(t, f), startEventCount(t, f.workspace)
			got := f.run("--json", "finding", tc.verb, ids[tc.index], "--actor", tc.actor, tc.flag, tc.value)
			if got.err == nil || before != reviewProjection(t, f) || count != startEventCount(t, f.workspace) {
				t.Fatalf("refusal changed state: %#v", got)
			}
		})
	}
}

func TestFindingDeferralAtomicAndUnaccepted(t *testing.T) {
	t.Parallel()
	for _, index := range []int{2, 3} {
		t.Run([]string{"minor", "nit"}[index-2], func(t *testing.T) {
			t.Parallel()
			f, ids, task := closureFixture(t)
			count := startEventCount(t, f.workspace)
			result, _ := lifecycleResult(t, f.run("--json", "finding", "defer", ids[index], "--actor", "agent:builder", "--title", "Follow up the finding"), "finding.defer")
			response := result["response"].(map[string]any)
			created := response["deferred_task"].(map[string]any)["id"].(string)
			if startEventCount(t, f.workspace) != count+1 {
				t.Fatal("deferral was not one Event")
			}
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
						t.Fatal(err)
					}
				}
				views, err := reviewwork.List(f.workspace, task)
				if err != nil || views[0].Findings[index].Open || views[0].Findings[index].Closure.Disposition != "deferred" {
					t.Fatalf("deferral not closed: %v / %v", views, err)
				}
				conn, err := reviewwork.OpenProjection(f.workspace)
				if err != nil {
					t.Fatal(err)
				}
				var lane, creator, criterion string
				var links int
				if err := conn.QueryRow(`SELECT i.lane,c.actor,a.text FROM issues i JOIN issue_creators c ON c.issue_id=i.id JOIN acceptance_criteria a ON a.issue_id=i.id WHERE i.id=?`, created).Scan(&lane, &creator, &criterion); err != nil {
					t.Fatal(err)
				}
				if err := conn.QueryRow(`SELECT count(*) FROM links WHERE (from_id=? AND relation='discovered_from' AND to_id=?) OR (from_id=? AND relation='discovered' AND to_id=?)`, created, task, task, created).Scan(&links); err != nil {
					t.Fatal(err)
				}
				closeCommandDB(t, conn)
				if lane != "standard" || creator != "agent:builder" || links != 2 || !strings.Contains(criterion, ids[index]) || !strings.Contains(criterion, "Retry conflicts lack coverage") || !strings.Contains(criterion, "Disabling the conflict guard") || !strings.Contains(criterion, "Retry a changed task") {
					t.Fatalf("discovery=%s / %s / %s / %d", lane, creator, criterion, links)
				}
				status := f.run("--json", "status")
				if status.err != nil || !strings.Contains(status.stdout, created) || !strings.Contains(status.stdout, "accept") {
					t.Fatalf("discovery not waiting for acceptance: %#v", status)
				}
			}
			next := f.run("--json", "next", "--explain")
			if next.err != nil || !strings.Contains(next.stdout, created) || !strings.Contains(next.stdout, "not-accepted") {
				t.Fatalf("discovery offered before acceptance: %#v", next)
			}
			before := reviewProjection(t, f)
			if got := f.run("--json", "finding", "answer", ids[index], "--actor", "agent:builder", "--text", "Too late"); got.err == nil {
				t.Fatal("answered deferred finding")
			}
			in := reviewInput()
			in["verdict"] = "CHANGES"
			in["dispositions"] = []any{disposition(ids[index], "resolved")}
			if got := recordReviewTask(t, f, in, "close-deferred", "human:courier", task); got.err == nil {
				t.Fatal("round replaced deferral")
			}
			if got := f.run("--json", "finding", "defer", ids[index], "--actor", "human:mark", "--title", "Again"); got.err == nil || reviewProjection(t, f) != before {
				t.Fatal("deferred twice")
			}
		})
	}
}
