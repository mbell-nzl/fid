package policy

import "testing"

func TestMatrixContainsExactlyTheImplementedV01WriteActions(t *testing.T) {
	t.Parallel()

	expectedActions := []string{
		"issue.create.epic",
		"issue.create.task",
		"issue.create.subtask",
		"issue.update",
		"issue.complete",
		"issue.ac.add",
		"issue.ac.verify.developer",
		"issue.ac.verify.qa_reviewer",
		"issue.ac.verify.qa_lead",
		"issue.ac.verify.tl_lite",
		"issue.ac.verify.tl",
		"issue.ac.verify.workflow_runner",
		"issue.ac.verify.human",
		"issue.transition.in_progress",
		"issue.transition.dev_complete",
		"issue.transition.failed_gate",
		"issue.transition.qa_review",
		"issue.transition.qa_pass",
		"issue.transition.qa_block",
		"issue.transition.wrapped",
		"issue.transition.closed",
		"issue.transition.done",
		"issue.transition.cancelled",
		"issue.claim.acquire",
		"issue.claim.renew",
		"issue.claim.force_acquire",
		"issue.claim.release",
		"concern.raise",
		"concern.address",
		"maintenance.rebuild",
		"requirement.sync",
		"link.add",
		"link.remove",
	}

	rows := MatrixRows()
	if len(rows) != len(expectedActions) {
		t.Fatalf("len(MatrixRows()) = %d, want %d", len(rows), len(expectedActions))
	}

	counts := make(map[string]int, len(rows))
	for _, row := range rows {
		counts[row.Action]++
	}

	for _, action := range expectedActions {
		if counts[action] != 1 {
			t.Fatalf("matrix entry count for %q = %d, want 1", action, counts[action])
		}
		if _, ok := MatrixEntry(action); !ok {
			t.Fatalf("MatrixEntry(%q) found = false, want true", action)
		}
	}

	for action, count := range counts {
		if count != 1 {
			t.Fatalf("matrix entry count for %q = %d, want 1", action, count)
		}
		if !containsAction(expectedActions, action) {
			t.Fatalf("matrix contains unexpected action %q", action)
		}
	}
}

func TestMatrixClaimRowsMatchSchemaSection13(t *testing.T) {
	t.Parallel()

	tests := []struct {
		action string
		want   MatrixRow
	}{
		{
			action: "issue.claim.acquire",
			want: MatrixRow{
				Action:         "issue.claim.acquire",
				Developer:      true,
				QAReviewer:     true,
				QALead:         true,
				TL:             true,
				WorkflowRunner: true,
				Human:          true,
			},
		},
		{
			action: "issue.claim.renew",
			want: MatrixRow{
				Action:         "issue.claim.renew",
				Developer:      true,
				QAReviewer:     true,
				QALead:         true,
				TL:             true,
				WorkflowRunner: true,
				Human:          true,
			},
		},
		{
			action: "issue.claim.force_acquire",
			want: MatrixRow{
				Action:         "issue.claim.force_acquire",
				TL:             true,
				WorkflowRunner: true,
				Human:          true,
			},
		},
		{
			action: "issue.claim.release",
			want: MatrixRow{
				Action:         "issue.claim.release",
				Developer:      true,
				QAReviewer:     true,
				QALead:         true,
				TL:             true,
				WorkflowRunner: true,
				Human:          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			got, ok := MatrixEntry(tt.action)
			if !ok {
				t.Fatalf("MatrixEntry(%q) found = false, want true", tt.action)
			}
			if got != tt.want {
				t.Fatalf("MatrixEntry(%q) = %#v, want %#v", tt.action, got, tt.want)
			}
		})
	}
}

func TestMatrixRequirementSyncRowMatchesSchemaSection13(t *testing.T) {
	t.Parallel()

	got, ok := MatrixEntry("requirement.sync")
	if !ok {
		t.Fatal(`MatrixEntry("requirement.sync") found = false, want true`)
	}

	want := MatrixRow{
		Action: "requirement.sync",
		TL:     true,
		Human:  true,
	}
	if got != want {
		t.Fatalf("MatrixEntry(%q) = %#v, want %#v", "requirement.sync", got, want)
	}
}

func TestMatrixLinkRowsMatchSchemaSection13(t *testing.T) {
	t.Parallel()

	tests := []struct {
		action string
		want   MatrixRow
	}{
		{
			action: "link.add",
			want: MatrixRow{
				Action:         "link.add",
				Developer:      true,
				QAReviewer:     true,
				QALead:         true,
				TLLite:         true,
				TL:             true,
				WorkflowRunner: true,
				Human:          true,
			},
		},
		{
			action: "link.remove",
			want: MatrixRow{
				Action:         "link.remove",
				Developer:      true,
				QAReviewer:     true,
				QALead:         true,
				TLLite:         true,
				TL:             true,
				WorkflowRunner: true,
				Human:          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			got, ok := MatrixEntry(tt.action)
			if !ok {
				t.Fatalf("MatrixEntry(%q) found = false, want true", tt.action)
			}
			if got != tt.want {
				t.Fatalf("MatrixEntry(%q) = %#v, want %#v", tt.action, got, tt.want)
			}
		})
	}
}

func containsAction(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
