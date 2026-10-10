package main

import (
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestP24ResumeIssuePosition(t *testing.T) {
	t.Parallel()
	resume := workflow.Resume{Issue: workflow.ResumeIssue{ID: "ISS-position", Status: "qa_pass", Revision: 4}}
	got := strings.Join(resumeContextHumanLines(resume), "\n")
	want := "Issue ISS-position\n  Status: qa_pass\n  Revision: 4\n  Next states: none.\n  Completion: "
	if !strings.HasPrefix(got, want) {
		t.Fatalf("Issue position: want prefix %q, got %q", want, got)
	}
}

func TestP24ResumeIssueNextStates(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name   string
		states []workflow.NextState
		want   string
	}{
		{"several", []workflow.NextState{{State: "cancelled", Category: "cancelled"}, {State: "qa_block", Category: "started"}, {State: "qa_pass", Category: "completed"}}, "cancelled (cancelled), qa_block (started), qa_pass (completed)"},
		{"empty", []workflow.NextState{}, "none."},
		{"nil", nil, "none."},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			resume := workflow.Resume{Issue: workflow.ResumeIssue{ID: "ISS-position", Status: "qa_review", Revision: 4, NextStates: row.states}}
			lines := resumeContextHumanLines(resume)
			if lines[2] != "  Revision: 4" || lines[3] != "  Next states: "+row.want || lines[4] != "  Completion: " {
				t.Fatalf("Issue lines = %q", lines)
			}
		})
	}
}

func TestP24IssueTransitionPositionHelp(t *testing.T) {
	t.Parallel()
	cmd := newIssueTransitionCommand()
	for name, want := range map[string]string{
		"to":                "Target workflow state; fid resume <issue-id> lists the Issue's next states",
		"expected-revision": "Required current Issue revision; fid resume <issue-id> shows it",
	} {
		if got := cmd.Flags().Lookup(name).Usage; got != want {
			t.Errorf("--%s usage = %q, want %q", name, got, want)
		}
	}
}
