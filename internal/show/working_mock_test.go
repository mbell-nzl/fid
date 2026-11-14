package show

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
)

func TestWorkingApprovedMocks(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, id, title, state, epicID, epicTitle string
		working                                   *holds.Working
	}{
		{"here", "FID-9f13aa", "Harden login rate limiting", "todo", "FID-e3a1c0", "Operating experience", &holds.Working{Here: true, Since: "2026-11-11T09:02:00Z", ExpiresAt: "2026-11-11T10:02:00Z"}},
		{"elsewhere", "FID-0d4e8b", "Document backup age", "todo", "FID-e3a1c0", "Operating experience", &holds.Working{Worktree: "/home/mark/projects/fid-wt3", Actor: "agent:codex", Since: "2026-11-11T08:40:00Z", ExpiresAt: "2026-11-11T09:40:00Z"}},
		{"interrupted", "FID-1f2e3d", "Contract fields and exit codes", "in_progress", "FID-8a1b2c", "Loop contract", &holds.Working{Worktree: "/home/mark/projects/fid-wt2", Actor: "agent:claude", Since: "2026-11-11T07:15:00Z", LastHeartbeatAt: "2026-11-11T07:15:00Z", Lapsed: true}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			done := 1
			if test.working.Lapsed {
				done = 0
			}
			result := Result{IssueView: issues.IssueView{
				Issue: issues.IssueDetail{IssueSummary: issues.IssueSummary{ID: test.id, Kind: "task", Title: test.title, Status: test.state, StatusCategory: "unstarted"}},
				Epic:  &issues.EpicReference{IssueReference: issues.IssueReference{ID: test.epicID, Kind: "epic", Title: test.epicTitle}, Progress: issues.Progress{Done: done, Total: 4}},
			}, Working: test.working, Now: time.Date(2026, 11, 11, 9, 10, 0, 0, time.UTC)}
			if test.working.Lapsed {
				result.Issue.StatusCategory = "started"
			}
			if test.name == "here" {
				priority := 1
				result.Issue.Priority = &priority
			}
			assertApprovedMock(t, test.name, result)
		})
	}
}

func TestSchedulingApprovedMock(t *testing.T) {
	t.Parallel()
	scheduling, state, priority := "sequential", issues.EpicInProgress, 2
	result := Result{IssueView: issues.IssueView{
		Issue: issues.IssueDetail{IssueSummary: issues.IssueSummary{ID: "FID-8a1b2c", Kind: "epic", Title: "Loop contract", Status: "in_progress", Priority: &priority}},
		State: &state, Scheduling: &scheduling, Progress: &issues.Progress{Total: 4},
		Children: []issues.ChildIssue{
			{IssueReference: issues.IssueReference{ID: "FID-1f2e3d", Kind: "task", Title: "Contract fields and exit codes", Status: "in_progress", StatusCategory: "started"}},
			{IssueReference: issues.IssueReference{ID: "FID-2f2e3d", Kind: "task", Title: "Golden tests for the contract", Status: "todo"}},
			{IssueReference: issues.IssueReference{ID: "FID-3f2e3d", Kind: "task", Title: "Scenario test for rework", Status: "todo"}},
			{IssueReference: issues.IssueReference{ID: "FID-4f2e3d", Kind: "task", Title: "Version the contract", Status: "todo"}},
		},
	}}
	assertApprovedMock(t, "sequential", result)
}

func assertApprovedMock(t *testing.T, name string, result Result) {
	t.Helper()
	want, err := os.ReadFile(filepath.Join("testdata", name+".txt"))
	if err != nil {
		t.Fatal(err)
	}
	if got := writeText(t, result, display.Renderer{}); got != string(want) {
		t.Fatalf("approved mock differs:\n%s\nwant:\n%s", got, want)
	}
}
