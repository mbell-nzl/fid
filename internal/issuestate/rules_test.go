package issuestate

import "testing"

func TestAllChildTasksDoneOrCancelled(t *testing.T) {
	t.Parallel()

	result := AllChildTasksDoneOrCancelled([]ChildTaskStatus{
		{IssueID: "FID-CHILD-1", State: "done"},
		{IssueID: "FID-CHILD-2", State: "qa_review"},
		{IssueID: "FID-CHILD-3", State: "cancelled"},
		{IssueID: "FID-CHILD-4", State: "in_progress"},
	})
	if result.Satisfied {
		t.Fatalf("expected non-final child task states to block the rule")
	}
	if len(result.Unmet) != 1 {
		t.Fatalf("len(result.Unmet) = %d, want 1", len(result.Unmet))
	}
	if result.Unmet[0].Rule != "all_child_tasks_done_or_cancelled" {
		t.Fatalf("unmet rule = %q, want all_child_tasks_done_or_cancelled", result.Unmet[0].Rule)
	}
	if len(result.Unmet[0].EntityIDs) != 2 || result.Unmet[0].EntityIDs[0] != "FID-CHILD-2" || result.Unmet[0].EntityIDs[1] != "FID-CHILD-4" {
		t.Fatalf("unmet entity IDs = %v, want [FID-CHILD-2 FID-CHILD-4]", result.Unmet[0].EntityIDs)
	}

	satisfied := AllChildTasksDoneOrCancelled([]ChildTaskStatus{
		{IssueID: "FID-CHILD-1", State: "done"},
		{IssueID: "FID-CHILD-2", State: "cancelled"},
	})
	if !satisfied.Satisfied {
		t.Fatalf("expected done or cancelled child tasks to satisfy the rule")
	}
	if len(satisfied.Unmet) != 0 {
		t.Fatalf("len(satisfied.Unmet) = %d, want 0", len(satisfied.Unmet))
	}
}
