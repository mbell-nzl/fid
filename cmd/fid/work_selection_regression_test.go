package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestReviewUnavailableBeforeFirstRound(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		lane, object string
		testReview   bool
	}{
		{lane: "standard", object: "HEAD"},
		{lane: "standard", object: "HEAD^{tree}"},
		{lane: "assured", object: "HEAD"},
		{lane: "routine", object: "HEAD", testReview: true},
	} {
		t.Run(tc.lane+"/"+tc.object, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithClock(t, tc.lane, false, time.Now, func(f *statusFixture) {
				if tc.testReview {
					trustedShowPolicy(t, f, func(raw string) string {
						return strings.Replace(raw, "test_adequacy = []", "test_adequacy = ['**/*.md']", 1)
					})
				}
			})
			f.workspace = gittest.AddWorktree(t, f.workspace)
			object := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", tc.object))
			objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
			if err := os.Remove(filepath.Join(objects, object[:2], object[2:])); err != nil {
				t.Fatal(err)
			}
			before := startEventCount(t, f.workspace)
			status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
			kind := "review-unavailable"
			if tc.lane == "assured" {
				kind = "approval-unavailable"
			}
			warnings := 0
			for _, raw := range status["warnings"].([]any) {
				warning := raw.(map[string]any)
				if warning["kind"] == kind && warning["task_id"] == task {
					warnings++
					if !strings.Contains(fmt.Sprint(warning["problem"]), "fingerprint ") || !strings.Contains(fmt.Sprint(warning["next_action"]), "restore") {
						t.Errorf("warning lost the comparison failure or recovery action: %v", warning)
					}
				}
			}
			if warnings != 1 || !strings.Contains(fmt.Sprint(status["in_progress"]), task) {
				t.Errorf("want one %s warning and the task's board entry: %v", kind, status)
			}
			if selectionItem(t, f, "approval") != nil || selectionItem(t, f, "merge") != nil {
				t.Error("unavailable comparison offered approval or merge")
			}
			for _, actor := range []string{"", "agent:reviewer"} {
				if got := selectionNext(t, f, actor); got["kind"] != "none" {
					t.Errorf("unavailable first review was offered to %q: %v", actor, got)
				}
			}
			if startEventCount(t, f.workspace) != before {
				t.Fatal("unavailable review appended an Event")
			}
			created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Other work", "--actor", "human:mark", "--idempotency-key", "other"), "issue.create")
			other := created["issue"].(map[string]any)["id"]
			for _, actor := range []string{"", "agent:reviewer"} {
				if got := selectionNext(t, f, actor); got["kind"] != "new" || got["task"].(map[string]any)["id"] != other {
					t.Fatalf("unavailable review blocked other work for %q: %v", actor, got)
				}
			}
		})
	}
}

func TestStaleApprovalActionAfterFreshReviews(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithLane(t, "assured")
	for _, stream := range []string{"primary", "separate"} {
		in := reviewInput()
		in["stream"] = stream
		in["reviewer"].(map[string]any)["session_id"] = "original-" + stream
		lifecycleResult(t, recordReviewTask(t, f, in, "original-"+stream, "human:courier", task), "review.record")
	}
	approveSuccess(t, approveForTest(f, task, task+"\n", true))
	if selectionItem(t, f, "merge") == nil || selectionItem(t, f, "approval") != nil {
		t.Fatal("current approval did not offer merge")
	}
	changes, _ := lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
	finding := changes["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)["id"].(string)
	prepareSecondReviewHandoff(t, f, task)
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "edited approved change\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Edit approved change")
	reviewHandoffTask(t, f, "edited-handoff", task)
	for _, stream := range []string{"primary", "separate"} {
		if selectionItem(t, f, "approval") != nil || selectionItem(t, f, "merge") != nil {
			t.Fatalf("stale approval offered an action before %s review", stream)
		}
		in := reviewInput()
		in["stream"] = stream
		in["reviewer"].(map[string]any)["session_id"] = "fresh-" + stream
		if stream == "primary" {
			in["dispositions"] = []any{disposition(finding, "resolved")}
		}
		lifecycleResult(t, recordReviewTask(t, f, in, "fresh-"+stream, "human:courier", task), "review.record")
	}
	shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	if shown["approval"].(map[string]any)["state"] != "stale" {
		t.Fatalf("edited change did not stale the approval: %v", shown)
	}
	// selectionItem also refuses duplicate actions of the same kind.
	item := selectionItem(t, f, "approval")
	if item == nil || item["summary"] != task+": approve, then merge PR #63" || item["next_action"] != "fid approve "+task+" --actor human:<name> --reason <text>" || selectionItem(t, f, "merge") != nil {
		t.Fatalf("stale approval after both streams approved: %v", item)
	}
	approveSuccess(t, approveForTest(f, task, task+"\n", true))
	if selectionItem(t, f, "approval") != nil || selectionItem(t, f, "merge") == nil {
		t.Fatal("reapproval did not replace approval with merge")
	}
}

func TestRoutineReviewSelectionNeedsNoComparison(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithLane(t, "routine")
	f.workspace = gittest.AddWorktree(t, f.workspace)
	head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
	objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
	if err := os.Remove(filepath.Join(objects, head[:2], head[2:])); err != nil {
		t.Fatal(err)
	}
	if got := selectionNext(t, f, ""); got["kind"] != "none" {
		t.Fatalf("Routine without required profiles offered review: %v", got)
	}
	status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	if len(status["warnings"].([]any)) != 0 || selectionItem(t, f, "merge") == nil {
		t.Fatalf("Routine %s without review evidence required a comparison: %v", task, status)
	}
}
