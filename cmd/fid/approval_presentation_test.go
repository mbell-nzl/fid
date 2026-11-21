package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func TestApprovalPresentation(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"missing", "current", "stale"} {
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			f := approvalDisplayFixture(t)
			if state != "missing" {
				approveSuccess(t, approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true))
			}
			if state == "stale" {
				lifecycleResult(t, recordReview(t, f, changesInput("primary"), "changes", "human:courier"), "review.record")
				prepareSecondReviewHandoff(t, f, "FID-9f13aa")
				writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "edited change\n")
				runGitForTest(t, f.workspace, "add", ".")
				runGitForTest(t, f.workspace, "commit", "-m", "Edit approved change")
				reviewHandoffTask(t, f, "edited-handoff", "FID-9f13aa")
			}
			assertApprovalDisplay(t, f, state, true)
			approvalDisplayGoldens(t, f, state)
		})
	}
}

func approvalDisplayFixture(t *testing.T) *statusFixture {
	t.Helper()
	return reviewFixtureWith(t, func(f *statusFixture) {
		f.clock = func() time.Time { return statusNow }
		startSQL(t, f, `UPDATE issues SET lane='assured' WHERE id='FID-9f13aa'`)
	})
}

func assertApprovalDisplay(t *testing.T, f *statusFixture, state string, flag bool) {
	t.Helper()
	before := startEventCount(t, f.workspace)
	shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-9f13aa"), "issue.show")
	approval, ok := shown["approval"].(map[string]any)
	if !ok || approval["state"] != state || shown["no_other_model_review"] != flag {
		t.Fatalf("approval display=%v, want %s / no other model=%t", shown, state, flag)
	}
	if state != "missing" {
		for _, key := range []string{"decision_id", "actor", "at", "fingerprint"} {
			if approval[key] == nil || approval[key] == "" {
				t.Fatalf("approval missing %s: %v", key, approval)
			}
		}
	}
	status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	found := false
	for _, raw := range status["needs_you"].([]any) {
		item := raw.(map[string]any)
		if item["kind"] != "approval" {
			continue
		}
		found = true
		if !strings.Contains(item["summary"].(string), state) || strings.Contains(item["summary"].(string), "no other-model review") != flag || item["next_action"] != "fid approve FID-9f13aa --actor human:<name> --reason <text>" {
			t.Fatalf("approval item=%v", item)
		}
	}
	inReview := shown["issue"].(map[string]any)["status"] == "in_review"
	if found != (state != "current" && inReview) {
		t.Fatalf("NEEDS YOU approval found=%t for %s", found, state)
	}
	tasks := status["in_progress"].(map[string]any)["tasks"].([]any)
	for _, epic := range status["in_progress"].(map[string]any)["epics"].([]any) {
		tasks = append(tasks, epic.(map[string]any)["tasks"].([]any)...)
	}
	if len(tasks) != 1 || tasks[0].(map[string]any)["no_other_model_review"] != flag {
		t.Fatalf("status task approval=%v", tasks)
	}
	if approval, ok := tasks[0].(map[string]any)["approval"].(map[string]any); !ok || approval["state"] != state {
		t.Fatalf("status task approval=%v", tasks)
	}
	if startEventCount(t, f.workspace) != before {
		t.Fatal("display appended an Event")
	}
}

func TestApprovalReworkPresentation(t *testing.T) {
	t.Parallel()
	f := approvalDisplayFixture(t)
	lifecycleResult(t, recordReview(t, f, changesInput("primary"), "changes", "human:courier"), "review.record")
	prepareSecondReviewHandoff(t, f, "FID-9f13aa")
	// The latest hand-off still explains the decision while rework is in
	// progress, but only in-review tasks request approval under NEEDS YOU.
	assertApprovalDisplay(t, f, "missing", false)
}

func approvalDisplayGoldens(t *testing.T, f *statusFixture, name string) {
	t.Helper()
	for _, command := range []string{"show", "status"} {
		for _, jsonMode := range []bool{false, true} {
			args, ext := []string{command}, ".txt"
			if command == "show" {
				args = append(args, "FID-9f13aa")
			}
			if jsonMode {
				args, ext = append([]string{"--json"}, args...), ".json"
			}
			got := f.run(args...)
			if got.err != nil {
				t.Fatal(got)
			}
			out := regexp.MustCompile(`\b(?:[0-9a-f]{64}|[0-9a-f]{40})\b`).ReplaceAllString(got.stdout, "<fingerprint-or-commit>")
			out = regexp.MustCompile(`\b(?:RND|PDN)-[0-9a-f]{6}\b`).ReplaceAllStringFunc(out, func(id string) string { return id[:3] + "-<new>" })
			handoffGolden(t, f, filepath.Join("..", "approval-display", name+"-"+command+ext), out)
		}
	}
}

func TestApprovalModelDisclosureAgrees(t *testing.T) {
	t.Parallel()
	for _, model := range []string{"codex/gpt-6-astra", "claude/opus-5.5"} {
		t.Run(model, func(t *testing.T) {
			t.Parallel()
			f := approvalDisplayFixture(t)
			for _, stream := range []string{"primary", "separate"} {
				in := reviewInput()
				in["stream"] = stream
				reviewer := in["reviewer"].(map[string]any)
				reviewer["model"], reviewer["session_id"] = model, stream+"-session"
				lifecycleResult(t, recordReview(t, f, in, stream, "human:courier"), "review.record")
			}
			want := strings.HasPrefix(model, "codex/")
			assertApprovalDisplay(t, f, "missing", want)
			prompt := approveForTest(f, "FID-9f13aa", "no\n", true)
			if prompt.err != nil || strings.Contains(prompt.stderr, "no other-model review") != want {
				t.Fatal(prompt)
			}
			conn, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = conn.Close() }()
			snapshot, err := doneeval.Gather(f.workspace, conn, "FID-9f13aa", doneeval.Target{}, doneeval.GitHubFacts{})
			if err != nil || doneeval.Evaluate(snapshot).NoOtherModelReview != want {
				t.Fatalf("done disclosure=%+v / %v", snapshot, err)
			}
			if !want {
				approvalDisplayGoldens(t, f, "cross-family")
			}
		})
	}
}

func TestApprovalDisplayUnavailablePolicyAndComparison(t *testing.T) {
	t.Parallel()
	for _, failure := range []string{"offline", "policy", "comparison", "fingerprint"} {
		t.Run(failure, func(t *testing.T) {
			t.Parallel()
			f := approvalDisplayFixture(t)
			if failure == "fingerprint" {
				lifecycleResult(t, recordReview(t, f, reviewInput(), "review", "human:courier"), "review.record")
			}
			approveSuccess(t, approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true))
			// No remote can be contacted successfully; reads must use local data.
			identityGit(t, f.workspace, "remote", "set-url", "origin", "/no/such/offline-remote")
			switch failure {
			case "policy":
				identityGit(t, f.workspace, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
			case "comparison":
				writeBackupTestFile(t, filepath.Join(f.workspace, ".git", "shallow"), strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))+"\n")
			case "fingerprint":
				head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
				objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
				if err := os.Remove(filepath.Join(objects, head[:2], head[2:])); err != nil {
					t.Fatal(err)
				}
			}
			if failure != "fingerprint" {
				assertApprovalDisplay(t, f, "current", true)
			} else {
				for _, command := range [][]string{{"show", "FID-9f13aa"}, {"status"}} {
					got := f.run(append([]string{"--json"}, command...)...)
					var envelope map[string]any
					if got.err != nil || json.Unmarshal([]byte(got.stdout), &envelope) != nil || !strings.Contains(got.stdout, `"approval":null`) || !strings.Contains(got.stdout, `"no_other_model_review":null`) || !strings.Contains(got.stdout, `"kind":"approval-unavailable"`) {
						t.Fatalf("unavailable approval=%#v", got)
					}
				}
			}
			approvalDisplayGoldens(t, f, failure)
		})
	}
}

func TestApprovalDisplayNotApplicable(t *testing.T) {
	t.Parallel()
	for _, handoff := range []bool{false, true} {
		t.Run(map[bool]string{false: "before handoff", true: "standard"}[handoff], func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "[]")
			if handoff {
				reviewHandoff(t, f, "standard-handoff")
			} else {
				startSQL(t, f, `UPDATE issues SET lane='assured' WHERE id='FID-9f13aa'`)
			}
			shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-9f13aa"), "issue.show")
			for _, key := range []string{"approval", "no_other_model_review"} {
				if value, exists := shown[key]; !exists || value != nil {
					t.Fatalf("%s=%v, want explicit null", key, value)
				}
			}
			for _, args := range [][]string{{"show", "FID-9f13aa"}, {"status"}} {
				if got := f.run(args...); got.err != nil || strings.Contains(got.stdout, "Approval") || strings.Contains(got.stdout, "no other-model review") || strings.Contains(got.stdout, "fid approve") {
					t.Fatalf("nonapplicable approval=%#v", got)
				}
			}
		})
	}
}
