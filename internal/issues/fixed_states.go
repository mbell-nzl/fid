package issues

// CommandTransitionError directs callers to the command that owns an edge.
type CommandTransitionError struct{ IssueID, From, To, Rule, Message, Next string }

func (err *CommandTransitionError) Error() string { return err.Message }

func fixedTransition(current issueTransitionState, to string) error {
	if current.Kind == "subtask" {
		return nil
	}
	rule, message, next := "", "", ""
	switch {
	case to == "in_progress":
		rule, message, next = "use-start", current.ID+" starts with fid start", "fid start --hold <hold> --actor <actor> --idempotency-key <key>"
	case current.Kind == "task" && to == "in_review":
		rule, message = "use-hand-off", current.ID+" goes to review through the hand-off"
	case current.Kind == "task" && to == "todo" && current.Status != "todo":
		rule, message = "task-already-started", current.ID+" cannot move back to to do"
	}
	if rule == "" {
		return nil
	}
	return &CommandTransitionError{IssueID: current.ID, From: current.Status, To: to, Rule: rule, Message: message, Next: next}
}
