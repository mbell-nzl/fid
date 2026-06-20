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
		"issue.ac.verify.tl",
		"issue.ac.verify.human",
		"issue.transition.in_progress",
		"issue.transition.dev_complete",
		"issue.transition.failed_gate",
		"issue.transition.qa_review",
		"issue.transition.qa_pass",
		"issue.transition.qa_block",
		"issue.transition.done",
		"issue.transition.cancelled",
		"concern.raise",
		"concern.address",
		"maintenance.rebuild",
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

func containsAction(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
