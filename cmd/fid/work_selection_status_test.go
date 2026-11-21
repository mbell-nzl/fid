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

func selectionItem(t *testing.T, f *statusFixture, kind string) map[string]any {
	t.Helper()
	status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	var found map[string]any
	for _, raw := range status["needs_you"].([]any) {
		item := raw.(map[string]any)
		if item["kind"] == kind {
			if found != nil {
				t.Fatalf("duplicate %s items: %v", kind, status)
			}
			found = item
		}
	}
	return found
}

func TestReworkCapSelectionAndStatus(t *testing.T) {
	t.Parallel()
	for _, cap := range []int{0, 3} {
		t.Run(fmt.Sprint(cap), func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithClock(t, "standard", false, time.Now, func(f *statusFixture) {
				trustedShowPolicy(t, f, func(raw string) string {
					return strings.Replace(raw, "[lanes.standard]", fmt.Sprintf("[lanes.standard]\nrework_cap=%d", cap), 1)
				})
			})
			for count := 1; count <= cap+1; count++ {
				in := changesInput("primary")
				if count > 1 {
					in["findings"] = []any{}
				}
				lifecycleResult(t, recordReviewTask(t, f, in, fmt.Sprint(count), "human:courier", task), "review.record")
				got := selectionNext(t, f, "")
				if count <= cap {
					if got["kind"] != "rework" || selectionItem(t, f, "rework-cap") != nil {
						t.Fatalf("rework stopped at/below cap: %v", got)
					}
					if count == cap {
						selectionNext(t, f, "agent:builder")
					}
					continue
				}
				if got["kind"] != "none" || got["waiting"].(map[string]any)["kind"] != "rework-cap" || !strings.Contains(fmt.Sprint(got["queue"]), "rework-cap-reached") {
					t.Fatalf("automatic rework passed cap: %v", got)
				}
				item := selectionItem(t, f, "rework-cap")
				want := fmt.Sprintf("%s reached its rework cap (%d rounds asked for changes; cap %d)", task, count, cap)
				if item == nil || item["summary"] != want || !strings.Contains(fmt.Sprint(item["next_action"]), "fid takeover "+task) || !strings.Contains(fmt.Sprint(item["next_action"]), "cancelled") {
					t.Fatalf("cap needs a person's action: %v", item)
				}
			}
			builder := f.workspace
			f.workspace = gittest.AddWorktree(t, builder)
			lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Other work", "--actor", "human:mark", "--idempotency-key", "other"), "issue.create")
			if got := selectionNext(t, f, "agent:other"); got["kind"] != "new" {
				t.Fatalf("capped task stopped another worktree: %v", got)
			}
		})
	}
}

func TestApprovalAndMergeItemsFollowReviews(t *testing.T) {
	t.Parallel()
	for _, lane := range []string{"standard", "assured"} {
		t.Run(lane, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithLane(t, lane)
			if selectionItem(t, f, "approval") != nil || selectionItem(t, f, "merge") != nil {
				t.Fatal("status requested merge or approval before reviews")
			}
			in := reviewInput()
			in["reviewer"].(map[string]any)["model"] = "codex/gpt-6-astra"
			lifecycleResult(t, recordReviewTask(t, f, in, "primary", "human:courier", task), "review.record")
			if lane == "assured" {
				if selectionItem(t, f, "approval") != nil || selectionItem(t, f, "merge") != nil {
					t.Fatal("status ignored the separate review")
				}
				in["stream"] = "separate"
				in["reviewer"].(map[string]any)["session_id"] = "separate-session"
				lifecycleResult(t, recordReviewTask(t, f, in, "separate", "human:courier", task), "review.record")
				item := selectionItem(t, f, "approval")
				if item == nil || item["summary"] != task+": approve, then merge PR #63, no other-model review" || !strings.Contains(fmt.Sprint(item["next_action"]), "fid approve "+task) || selectionItem(t, f, "merge") != nil {
					t.Fatalf("approval item=%v", item)
				}
				approveSuccess(t, approveForTest(f, task, task+"\n", true))
			}
			item := selectionItem(t, f, "merge")
			if item == nil || item["summary"] != task+": merge PR #63" || item["next_action"] != "https://github.com/mbell-nzl/fid/pull/63" || selectionItem(t, f, "approval") != nil {
				t.Fatalf("merge item=%v", item)
			}
			facts := doneFacts(t, f, task)
			facts.CI.Runs = facts.CI.Runs[:0]
			got := runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "human:mark", "--idempotency-key", "merge")
			if got.err == nil || !strings.Contains(got.stderr, "done-ci-waiting") || selectionItem(t, f, "merge") != nil || selectionItem(t, f, "approval") != nil {
				t.Fatalf("recorded merge still needs merging: %#v", got)
			}
		})
	}
}

func TestReworkCapUsesRaisedLaneAndFrozenPolicy(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithClock(t, "standard", false, time.Now, func(f *statusFixture) {
		trustedShowPolicy(t, f, func(raw string) string {
			raw = strings.Replace(raw, "[lanes.assured]", "[lanes.assured]\npaths=['**/*.md']\nrework_cap=1", 1)
			return strings.Replace(raw, "[lanes.standard]", "[lanes.standard]\nrework_cap=9", 1)
		})
	})
	// Today's policy does not change the cap the raised Assured hand-off froze.
	trustedShowPolicy(t, f, func(raw string) string { return strings.Replace(raw, "rework_cap=1", "rework_cap=9", 1) })
	for i, stream := range []string{"primary", "separate"} {
		in := changesInput(stream)
		if i > 0 {
			in["findings"] = []any{}
		}
		lifecycleResult(t, recordReviewTask(t, f, in, stream, "human:courier", task), "review.record")
	}
	if item := selectionItem(t, f, "rework-cap"); item == nil || !strings.Contains(item["summary"].(string), "2 rounds asked for changes; cap 1") {
		t.Fatalf("cap ignored frozen raised lane or a review stream: %v", item)
	}
	if got := selectionNext(t, f, "agent:builder"); got["kind"] != "none" {
		t.Fatalf("claim served capped rework: %v", got)
	}
}

func TestReworkCapCountsEarlierHandoffs(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithClock(t, "standard", false, time.Now, func(f *statusFixture) {
		trustedShowPolicy(t, f, func(raw string) string {
			return strings.Replace(raw, "[lanes.standard]", "[lanes.standard]\nrework_cap=1", 1)
		})
	})
	lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "first", "human:courier", task), "review.record")
	prepareSecondReviewHandoff(t, f, task)
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "Fix for the first review\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Rework")
	reviewHandoffTask(t, f, "second-handoff", task)
	in := changesInput("primary")
	in["findings"] = []any{}
	lifecycleResult(t, recordReviewTask(t, f, in, "second", "human:courier", task), "review.record")
	if item := selectionItem(t, f, "rework-cap"); item == nil || !strings.Contains(item["summary"].(string), "2 rounds asked for changes; cap 1") {
		t.Fatalf("earlier hand-off stopped counting: %v", item)
	}
}

func TestStandardReviewUnavailablePreservesStatus(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
	head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
	objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
	if err := os.Remove(filepath.Join(objects, head[:2], head[2:])); err != nil {
		t.Fatal(err)
	}
	status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	if !strings.Contains(fmt.Sprint(status["warnings"]), "review-unavailable") || selectionItem(t, f, "merge") != nil {
		t.Fatalf("unavailable review lost its warning or offered merge: %v", status)
	}
}
