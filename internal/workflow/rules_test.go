package workflow

import "testing"

func TestAllACsVerified(t *testing.T) {
	t.Parallel()

	t.Run("returns unmet acceptance criteria and missing roles", func(t *testing.T) {
		result := AllACsVerified(
			[]string{"developer", "qa_lead"},
			[]AcceptanceCriterionVerification{
				{AcceptanceCriterionID: "AC-1", PassingRoles: []string{"developer"}},
				{AcceptanceCriterionID: "AC-2", PassingRoles: []string{"developer", "qa_lead"}},
				{AcceptanceCriterionID: "AC-3", PassingRoles: nil},
			},
		)

		if result.Satisfied {
			t.Fatalf("expected unsatisfied result when required roles are missing")
		}
		if len(result.Unmet) != 2 {
			t.Fatalf("len(result.Unmet) = %d, want 2", len(result.Unmet))
		}

		first := result.Unmet[0]
		if first.Rule != "all_acs_verified" {
			t.Fatalf("first unmet rule = %q, want all_acs_verified", first.Rule)
		}
		if len(first.EntityIDs) != 1 || first.EntityIDs[0] != "AC-1" {
			t.Fatalf("first unmet entity IDs = %v, want [AC-1]", first.EntityIDs)
		}
		if len(first.MissingRoles) != 1 || first.MissingRoles[0] != "qa_lead" {
			t.Fatalf("first unmet missing roles = %v, want [qa_lead]", first.MissingRoles)
		}

		second := result.Unmet[1]
		if len(second.EntityIDs) != 1 || second.EntityIDs[0] != "AC-3" {
			t.Fatalf("second unmet entity IDs = %v, want [AC-3]", second.EntityIDs)
		}
		if len(second.MissingRoles) != 2 || second.MissingRoles[0] != "developer" || second.MissingRoles[1] != "qa_lead" {
			t.Fatalf("second unmet missing roles = %v, want [developer qa_lead]", second.MissingRoles)
		}
	})

	t.Run("returns satisfied when every acceptance criterion has passing verifications for the required roles", func(t *testing.T) {
		result := AllACsVerified(
			[]string{"developer"},
			[]AcceptanceCriterionVerification{
				{AcceptanceCriterionID: "AC-1", PassingRoles: []string{"developer"}},
				{AcceptanceCriterionID: "AC-2", PassingRoles: []string{"developer", "qa_lead"}},
			},
		)

		if !result.Satisfied {
			t.Fatalf("expected satisfied result when required roles are present")
		}
		if len(result.Unmet) != 0 {
			t.Fatalf("len(result.Unmet) = %d, want 0", len(result.Unmet))
		}
	})
}

func TestHasCompletionReport(t *testing.T) {
	t.Parallel()

	result := HasCompletionReport(CompletionReportStatus{IssueID: "FID-123", Exists: false})
	if result.Satisfied {
		t.Fatalf("expected missing completion report to be unsatisfied")
	}
	if len(result.Unmet) != 1 {
		t.Fatalf("len(result.Unmet) = %d, want 1", len(result.Unmet))
	}
	if result.Unmet[0].Rule != "has_completion_report" {
		t.Fatalf("unmet rule = %q, want has_completion_report", result.Unmet[0].Rule)
	}
	if len(result.Unmet[0].EntityIDs) != 1 || result.Unmet[0].EntityIDs[0] != "FID-123" {
		t.Fatalf("unmet entity IDs = %v, want [FID-123]", result.Unmet[0].EntityIDs)
	}

	satisfied := HasCompletionReport(CompletionReportStatus{IssueID: "FID-123", Exists: true})
	if !satisfied.Satisfied {
		t.Fatalf("expected completion report presence to satisfy the rule")
	}
	if len(satisfied.Unmet) != 0 {
		t.Fatalf("len(satisfied.Unmet) = %d, want 0", len(satisfied.Unmet))
	}
}

func TestConcernsAddressed(t *testing.T) {
	t.Parallel()

	result := ConcernsAddressed(ConcernStatus{OpenConcernIDs: []string{"CON-1", "CON-2"}})
	if result.Satisfied {
		t.Fatalf("expected open concerns to be unsatisfied")
	}
	if len(result.Unmet) != 1 {
		t.Fatalf("len(result.Unmet) = %d, want 1", len(result.Unmet))
	}
	if result.Unmet[0].Rule != "concerns_addressed" {
		t.Fatalf("unmet rule = %q, want concerns_addressed", result.Unmet[0].Rule)
	}
	if len(result.Unmet[0].EntityIDs) != 2 || result.Unmet[0].EntityIDs[0] != "CON-1" || result.Unmet[0].EntityIDs[1] != "CON-2" {
		t.Fatalf("unmet entity IDs = %v, want [CON-1 CON-2]", result.Unmet[0].EntityIDs)
	}

	satisfied := ConcernsAddressed(ConcernStatus{})
	if !satisfied.Satisfied {
		t.Fatalf("expected no open concerns to satisfy the rule")
	}
	if len(satisfied.Unmet) != 0 {
		t.Fatalf("len(satisfied.Unmet) = %d, want 0", len(satisfied.Unmet))
	}
}

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
