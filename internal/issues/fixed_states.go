package issues

import "github.com/mbell-nzl/fid/internal/issuestate"

// CommandTransitionError directs callers to the command that owns an edge.
type CommandTransitionError struct{ IssueID, From, To, Rule, Message, Next string }

func (err *CommandTransitionError) Error() string { return err.Message }

// fixedTransition refuses a move of a task or an epic that another command
// owns (issuestate.StateOwners), pointing at that command. fid issue
// transition makes a person's cancellation and every move of a subtask.
func fixedTransition(current issueTransitionState, to string) error {
	owners := issuestate.StateOwners(current.Kind, to)
	if !issuestate.HasFixedStates(current.Kind) || owners[issuestate.FamilyTransition] {
		return nil
	}
	rule, message, next := "", "", ""
	switch {
	case owners[issuestate.FamilyDone]:
		rule, message, next = "use-done", current.ID+" completes with fid done", "fid done "+current.ID+" --actor <actor>"
	case owners[issuestate.FamilyEpicDone]:
		rule, message, next = "use-done", "Close the epic with fid done "+current.ID, `fid done `+current.ID+` --actor human:<name> --reason "<why>"`
	case owners[issuestate.FamilyStart]:
		rule, message, next = "use-start", current.ID+" starts with fid start", "fid start --hold <hold> --actor <actor> --idempotency-key <key>"
	case owners[issuestate.FamilyHandoff]:
		rule, message = "use-hand-off", current.ID+" goes to review through the hand-off"
	case current.Kind == "task" && to == "todo" && current.Status != "todo":
		rule, message = "task-already-started", current.ID+" cannot move back to to do"
	}
	if rule == "" {
		return nil
	}
	return &CommandTransitionError{IssueID: current.ID, From: current.Status, To: to, Rule: rule, Message: message, Next: next}
}
