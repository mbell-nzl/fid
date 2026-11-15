package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func reviewFixture(t *testing.T) *statusFixture {
	t.Helper()
	return reviewFixtureWith(t, nil)
}

func reviewFixtureWith(t *testing.T, prepare func(*statusFixture)) *statusFixture {
	t.Helper()
	f := requiredCheckFixture(t, "[]")
	if prepare != nil {
		prepare(f)
	}
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "review this change\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Task change")
	reviewHandoff(t, f, "review-handoff")
	return f
}

func reviewInput() map[string]any {
	return map[string]any{"reviewer": map[string]any{"actor": "agent:reviewer", "model": "claude/opus-5.5", "session_id": "review-session-1"}, "stream": "primary", "profiles": []string{"code", "test-adequacy"}, "verdict": "APPROVE", "purpose": "Persist the task change.", "findings": []any{}, "verification": "go test ./... passed", "provenance": map[string]any{"prompt": "Review the PR", "output": "## Review: APPROVE\n### Findings\nNone."}}
}

func recordReview(t *testing.T, f *statusFixture, input map[string]any, key, actor string) commandExecutionResult {
	t.Helper()
	return recordReviewTask(t, f, input, key, actor, "FID-9f13aa")
}

func recordReviewTask(t *testing.T, f *statusFixture, input map[string]any, key, actor, task string) commandExecutionResult {
	t.Helper()
	raw, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "review.json")
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	return f.run("--json", "review", "record", task, "--input", path, "--actor", actor, "--idempotency-key", key)
}

func TestReviewRoundCourierAndRetry(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	input := reviewInput()
	first := recordReview(t, f, input, "round", "agent:claude-builder")
	lifecycleResult(t, first, "review.record")
	before := startEventCount(t, f.workspace)
	if got := recordReview(t, f, input, "round", "agent:claude-builder"); got.err != nil || got.stdout != first.stdout || startEventCount(t, f.workspace) != before {
		t.Fatalf("retry=%#v", got)
	}
	listed, _ := lifecycleResult(t, f.run("--json", "review", "list", "FID-9f13aa"), "review.list")
	if !strings.Contains(first.stdout, "RND-") || !strings.Contains(first.stdout, "patch-id-verbatim-per-file-v1") || len(listed["rounds"].([]any)) != 1 {
		t.Fatalf("round=%#v list=%v", first, listed)
	}
}

// Reviewers remain independent from every builder, not just the last session.
func TestReviewAdmissionRefusals(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, want string
		prepare    func(*testing.T, *statusFixture, map[string]any)
	}{
		{"ended builder", "builder cannot", func(_ *testing.T, _ *statusFixture, in map[string]any) {
			in["reviewer"].(map[string]any)["actor"] = "agent:claude-builder"
		}},
		{"missing code", "missing required profile code", func(_ *testing.T, _ *statusFixture, in map[string]any) { in["profiles"] = []string{"test-adequacy"} }},
		{"separate standard", "Assured hand-off", func(_ *testing.T, _ *statusFixture, in map[string]any) { in["stream"] = "separate" }},
		{"approve major", "verdict disagrees", func(_ *testing.T, _ *statusFixture, in map[string]any) {
			in["findings"] = []any{reviewFinding("major")}
		}},
		{"changes no major", "verdict disagrees", func(_ *testing.T, _ *statusFixture, in map[string]any) { in["verdict"] = "CHANGES" }},
		{"malformed model", "model must", func(_ *testing.T, _ *statusFixture, in map[string]any) {
			in["reviewer"].(map[string]any)["model"] = "not-a-model"
		}},
		{"malformed severity", "input.findings[0].severity", func(_ *testing.T, _ *statusFixture, in map[string]any) {
			f := reviewFinding("blocker")
			in["findings"] = []any{f}
		}},
		{"unknown profile", "finding requires", func(_ *testing.T, _ *statusFixture, in map[string]any) {
			f := reviewFinding("minor")
			f["profile"] = "security"
			in["findings"] = []any{f}
		}},
		{"unknown member", "input.other_model", func(_ *testing.T, _ *statusFixture, in map[string]any) { in["other_model"] = true }},
		{"service reviewer", "work actor", func(_ *testing.T, _ *statusFixture, in map[string]any) {
			in["reviewer"].(map[string]any)["actor"] = "service:reviewer"
		}},
		{"missing session", "input.reviewer.session_id", func(_ *testing.T, _ *statusFixture, in map[string]any) {
			delete(in["reviewer"].(map[string]any), "session_id")
		}},
		{"missing objects", "fingerprint missing-object", func(t *testing.T, f *statusFixture, _ map[string]any) {
			head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
			object := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
			if err := os.Remove(filepath.Join(object, head[:2], head[2:])); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := reviewFixture(t)
			in := reviewInput()
			tc.prepare(t, f, in)
			before := reviewProjection(t, f)
			got := recordReview(t, f, in, "refused", "human:courier")
			if got.err == nil || !strings.Contains(got.stderr, tc.want) {
				t.Fatalf("refusal=%#v, want %s", got, tc.want)
			}
			if after := reviewProjection(t, f); after != before {
				t.Fatalf("refusal changed state: before=%s after=%s", before, after)
			}
		})
	}
}

func reviewFinding(severity string) map[string]any {
	return map[string]any{"profile": "code", "severity": severity, "summary": "Retry conflicts lack coverage", "location": map[string]any{"file": "cmd/fid/merge_records_test.go", "line": 111}, "evidence": "Disabling the conflict guard left focused tests green.", "fix": "Retry a changed task and PR; assert refusal and unchanged state."}
}

func reviewProjection(t *testing.T, f *statusFixture) string {
	t.Helper()
	// Journal count catches any append, and SQL covers round/finding/retry writes.
	got := f.run("sql", `SELECT 'round' AS kind,id AS value FROM review_rounds UNION ALL SELECT 'finding',id FROM round_findings UNION ALL SELECT 'retry',result_json FROM idempotency_results UNION ALL SELECT 'closure',finding_id || ':' || coalesce(round_id,decision_id) || ':' || disposition || ':' || reason FROM finding_closures UNION ALL SELECT 'decision',id || ':' || binding_id || ':' || reason FROM person_decisions ORDER BY kind,value`)
	if got.err != nil {
		t.Fatal(got)
	}
	return fmt.Sprint(startEventCount(t, f.workspace)) + got.stdout
}

func TestReviewNoHandoffAndServiceRecorder(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ actor, want string }{{"human:courier", "needs a hand-off"}, {"service:ci", "authority-denied"}} {
		t.Run(tc.actor, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "[]")
			before := reviewProjection(t, f)
			got := recordReview(t, f, reviewInput(), "refuse", tc.actor)
			if got.err == nil || !strings.Contains(got.stderr, tc.want) || reviewProjection(t, f) != before {
				t.Fatalf("refused=%#v", got)
			}
		})
	}
}

func TestReviewExistingMajorAndConflictingRetry(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	in := reviewInput()
	in["verdict"] = "CHANGES"
	in["findings"] = []any{reviewFinding("major"), reviewFinding("minor"), reviewFinding("nit")}
	first := recordReview(t, f, in, "first", "human:courier")
	lifecycleResult(t, first, "review.record")
	shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-9f13aa"), "issue.show")
	if shown["issue"].(map[string]any)["status"] != "in_progress" {
		t.Fatalf("CHANGES did not return task to the builder: %v", shown)
	}
	before := reviewProjection(t, f)
	conflict := recordReview(t, f, reviewInput(), "first", "human:courier")
	assertJSONFailureCode(t, conflict.stderr, "review.record", "idempotency-conflict")
	got := recordReview(t, f, reviewInput(), "second", "human:courier")
	if got.err == nil || !strings.Contains(got.stderr, "verdict disagrees") || reviewProjection(t, f) != before {
		t.Fatalf("old major=%#v", got)
	}
	later := reviewInput()
	later["verdict"] = "CHANGES"
	lifecycleResult(t, recordReview(t, f, later, "third", "human:courier"), "review.record")
	listed, _ := lifecycleResult(t, f.run("--json", "review", "list", "FID-9f13aa"), "review.list")
	rounds := listed["rounds"].([]any)
	if len(rounds) != 2 {
		t.Fatal(listed)
	}
	for _, finding := range rounds[0].(map[string]any)["findings"].([]any) {
		if finding.(map[string]any)["open"] != true {
			t.Fatal(finding)
		}
	}
}

func reviewHandoff(t *testing.T, f *statusFixture, key string) {
	t.Helper()
	reviewHandoffTask(t, f, key, "FID-9f13aa")
}

func reviewHandoffTask(t *testing.T, f *statusFixture, key, task string) {
	t.Helper()
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	var hold, actor string
	if err := conn.QueryRow(`SELECT id,actor FROM holds WHERE task_id=? AND ended_at IS NULL`, task).Scan(&hold, &actor); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	lifecycleResult(t, f.run("--json", "handoff", "--hold", hold, "--actor", actor, "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", key), "task.handoff")
}

func reviewPolicy(t *testing.T, f *statusFixture, lane, paths string) {
	t.Helper()
	config := filepath.Join(f.workspace, ".fid.toml")
	raw, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Replace(string(raw), "test_adequacy = []", "test_adequacy = "+paths, 1)
	if text == string(raw) && paths != "[]" {
		t.Fatal("fixture has no test path setting")
	}
	if lane == "assured" {
		// The fixture already configures Assured checks.
		lifecycleResult(t, f.run("--json", "issue", "update", "FID-9f13aa", "--lane", "assured", "--actor", "human:mark", "--idempotency-key", "raise-lane"), "issue.update")
		f.clock = func() time.Time { return time.Now().Add(time.Minute) }
		lifecycleResult(t, f.run("--json", "takeover", "FID-9f13aa", "--actor", "human:policy-builder", "--reason", "Build at current fixture time"), "task.takeover")
	}
	writeBackupTestFile(t, config, text)
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "--allow-empty", "-m", "Review policy")
	runGitForTest(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", "HEAD")
}

func TestReviewRequiredProfileFromFrozenHandoff(t *testing.T) {
	t.Parallel()
	f := reviewFixtureWith(t, func(f *statusFixture) { reviewPolicy(t, f, "standard", "['**/*.md']") })
	// Current policy drops test review after the hand-off. Frozen paths still win.
	reviewPolicy(t, f, "standard", "[]")
	in := reviewInput()
	in["profiles"] = []string{"code"}
	before := reviewProjection(t, f)
	got := recordReview(t, f, in, "missing-test", "human:courier")
	if got.err == nil || !strings.Contains(got.stderr, "required profile test-adequacy") || reviewProjection(t, f) != before {
		t.Fatalf("profile refusal=%#v", got)
	}
}

func TestReviewSeparateSessionAndFrozenLane(t *testing.T) {
	t.Parallel()
	for _, reverse := range []bool{false, true} {
		t.Run(fmt.Sprint(reverse), func(t *testing.T) {
			t.Parallel()
			f := reviewFixtureWith(t, func(f *statusFixture) { reviewPolicy(t, f, "assured", "[]") })
			first, second := reviewInput(), reviewInput()
			if reverse {
				first["stream"] = "separate"
			} else {
				second["stream"] = "separate"
			}
			lifecycleResult(t, recordReview(t, f, first, "first", "human:courier"), "review.record")
			before := reviewProjection(t, f)
			got := recordReview(t, f, second, "same-session", "human:courier")
			if got.err == nil || !strings.Contains(got.stderr, "different session") || reviewProjection(t, f) != before {
				t.Fatalf("session refusal=%#v", got)
			}
			second["reviewer"].(map[string]any)["session_id"] = "fresh-session"
			lifecycleResult(t, recordReview(t, f, second, "second", "human:courier"), "review.record")
		})
	}
}

func TestReviewOtherModelFromEveryBuildHold(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		builders []string
		reviewer string
		want     bool
	}{
		{"same family", []string{"codex/gpt-6.1-sol"}, "codex/gpt-6-astra", false},
		{"other family", []string{"codex/gpt-6.1-sol"}, "claude/opus-5.5", true},
		{"person", []string{""}, "claude/opus-5.5", true},
		{"unknown", []string{"codex/gpt-6.1-sol"}, "", false},
		{"every family", []string{"codex/gpt-6.1-sol", "claude/opus-5.5"}, "codex/gpt-6-astra", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := reviewFixtureWith(t, func(f *statusFixture) {
				for i, model := range tc.builders {
					args := []string{"--json", "takeover", "FID-9f13aa", "--actor", fmt.Sprintf("human:builder%d", i), "--reason", "Continue build"}
					if model != "" {
						args = append(args, "--model", model)
					}
					lifecycleResult(t, f.run(args...), "task.takeover")
				}
			})
			in := reviewInput()
			if tc.reviewer == "" {
				delete(in["reviewer"].(map[string]any), "model")
			} else {
				in["reviewer"].(map[string]any)["model"] = tc.reviewer
			}
			lifecycleResult(t, recordReview(t, f, in, "round", "human:courier"), "review.record")
			shown, _ := lifecycleResult(t, f.run("--json", "review", "list", "FID-9f13aa"), "review.list")
			r := shown["rounds"].([]any)[0].(map[string]any)
			if r["other_model"] != tc.want || r["current"] != true {
				t.Fatalf("round=%v want other=%t", r, tc.want)
			}
		})
	}
}

func TestReviewRefusesReplacedBuilder(t *testing.T) {
	t.Parallel()
	f := reviewFixtureWith(t, func(f *statusFixture) {
		lifecycleResult(t, f.run("--json", "takeover", "FID-9f13aa", "--actor", "human:replacement", "--reason", "Continue"), "task.takeover")
	})
	before := reviewProjection(t, f)
	in := reviewInput()
	in["reviewer"].(map[string]any)["actor"] = "agent:claude-builder"
	got := recordReview(t, f, in, "replaced", "human:courier")
	if got.err == nil || !strings.Contains(got.stderr, "builder cannot") || reviewProjection(t, f) != before {
		t.Fatalf("replaced=%#v", got)
	}
}

func TestReviewRetryRejectsCorruptedStoredResult(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"replayed", "round", "occurred_at", "altered narrative"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			f := reviewFixture(t)
			in := reviewInput()
			lifecycleResult(t, recordReview(t, f, in, "round", "human:courier"), "review.record")
			mutateStoredRetry(t, f, "round", func(result map[string]any) {
				if field == "altered narrative" {
					result["round"].(map[string]any)["purpose"] = "Forged review narrative"
				} else {
					delete(result, field)
				}
			})
			before := reviewProjection(t, f)
			got := recordReview(t, f, in, "round", "human:courier")
			if got.err == nil || reviewProjection(t, f) != before {
				t.Fatalf("corrupt retry=%#v", got)
			}
		})
	}
}

func TestReviewRetryDoesNotNeedObjectsOrLiveHold(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	in := reviewInput()
	first := recordReview(t, f, in, "round", "human:courier")
	lifecycleResult(t, first, "review.record")
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
	before := reviewProjection(t, f)
	again := recordReview(t, f, in, "round", "human:courier")
	if again.err != nil || again.stdout != first.stdout || reviewProjection(t, f) != before {
		t.Fatalf("retry=%#v", again)
	}
}

func TestReviewMissingProjectUsesCommonRefusal(t *testing.T) {
	t.Parallel()
	root := gittest.NewWorkspace(t)
	raw, err := json.Marshal(reviewInput())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "review.json")
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"review", "list", "FID-missing"}, {"review", "record", "FID-missing", "--input", path, "--actor", "human:courier", "--idempotency-key", "key"}} {
		got := executeCommandForTest(root, append([]string{"--json"}, args...)...)
		assertJSONFailureCode(t, got.stderr, "review."+args[1], "project-not-found")
	}
}

func TestReviewRefusesOversizedJournalEventBeforeAppend(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	in := reviewInput()
	in["provenance"].(map[string]any)["output"] = strings.Repeat("x", 600000)
	before := reviewProjection(t, f)
	got := recordReview(t, f, in, "oversized", "human:courier")
	if got.err == nil {
		t.Fatal("recorded review whose duplicated Event exceeds the journal line budget")
	}
	if !strings.Contains(got.stderr, "review Event exceeds") || reviewProjection(t, f) != before {
		t.Fatalf("oversize refusal=%v %s", got.err, got.stderr)
	}
}

func TestReviewRecordFromStdin(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	in := reviewInput()
	in["findings"] = []any{reviewFinding("nit")}
	in["findings"].([]any)[0].(map[string]any)["location"] = nil
	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	root := newRootCommandAt(f.workspace)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetIn(bytes.NewReader(raw))
	root.SetArgs([]string{"--json", "review", "record", "FID-9f13aa", "--input", "-", "--actor", "human:courier", "--idempotency-key", "stdin"})
	err = root.Execute()
	result := commandExecutionResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
	fields, _ := lifecycleResult(t, result, "review.record")
	finding := fields["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)
	if finding["location"] != nil || finding["severity"] != "nit" {
		t.Fatal(finding)
	}
}
