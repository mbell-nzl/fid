package issues

import (
	"errors"
	"testing"
)

// fid issue transition points each move of a task or an epic that another
// command owns at that command; only a person's cancellation, and every
// move of a subtask, are its own. The rest is left to the workflow's edges.
func TestFixedTransitionPointsAtTheOwningCommand(t *testing.T) {
	t.Parallel()
	workflows := map[string]string{"task": "task-default", "epic": "epic-default", "subtask": "subtask-default"}
	states := []string{"todo", "in_progress", "in_review", "done", "cancelled", "unknown"}
	want := func(kind, from, to string) string {
		switch {
		case kind != "subtask" && from == "done" && to == "cancelled":
			return "finished-work-stays-done"
		case kind == "subtask":
			return ""
		case to == "done":
			return "use-done"
		case to == "in_progress":
			return "use-start"
		case kind == "task" && to == "in_review":
			return "use-hand-off"
		case kind == "task" && to == "todo" && from != "todo":
			return "task-already-started"
		}
		return ""
	}
	for kind, workflow := range workflows {
		for _, from := range states {
			for _, to := range states {
				err := fixedTransition(issueTransitionState{ID: "FID-a", Kind: kind, WorkflowID: workflow, Status: from}, to)
				rule := ""
				var refused *CommandTransitionError
				if errors.As(err, &refused) {
					rule = refused.Rule
					if refused.From != from || refused.To != to || refused.IssueID != "FID-a" || refused.Message == "" {
						t.Errorf("%s %s -> %s refusal = %+v", kind, from, to, refused)
					}
				} else if err != nil {
					t.Errorf("%s %s -> %s = %v", kind, from, to, err)
				}
				if rule != want(kind, from, to) {
					t.Errorf("%s %s -> %s rule = %q, want %q", kind, from, to, rule, want(kind, from, to))
				}
			}
		}
	}
}
