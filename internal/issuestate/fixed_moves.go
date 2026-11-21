package issuestate

// The fixed states of tasks and epics (ADR-0014; decisions record, 11
// November 2026, item 3). fid's commands move them, and each move belongs to
// one command's Event family. fid issue transition points at the owning
// command, and journal admission takes a move only from its owner's Event,
// on live append, replay and import. Subtasks keep plain transitions.

// The Event families of the commands that own a move.
const (
	FamilyCreate   = "issue.created"         // fid issue create: an Issue's first state
	FamilyStart    = "task.started"          // fid start: the task, and its epic's first start
	FamilyHandoff  = "task.handed-off"       // the hand-off
	FamilyRework   = "review.round-recorded" // a review round asking for changes
	FamilyDone     = "task.finished"         // fid done <task>
	FamilyEpicDone = "epic.closed"           // fid done <epic>, a person's decision
	// FamilyTransition is fid issue transition's: a person's cancel decision
	// for a task or an epic, and any move of a subtask.
	FamilyTransition = "issue.transitioned"
)

type fixedMove struct{ kind, from, to string }

// fixedMoves names the owner of every move of a task's or an epic's state.
// From "" is creation. No command makes a move missing here.
var fixedMoves = map[fixedMove]string{
	{"task", "", "todo"}:                 FamilyCreate,
	{"task", "todo", "in_progress"}:      FamilyStart,
	{"task", "in_progress", "in_review"}: FamilyHandoff,
	{"task", "in_review", "in_progress"}: FamilyRework,
	{"task", "in_review", "done"}:        FamilyDone,
	{"task", "todo", "cancelled"}:        FamilyTransition,
	{"task", "in_progress", "cancelled"}: FamilyTransition,
	{"task", "in_review", "cancelled"}:   FamilyTransition,
	{"task", "done", "cancelled"}:        FamilyTransition,
	{"epic", "", "todo"}:                 FamilyCreate,
	{"epic", "todo", "in_progress"}:      FamilyStart,
	{"epic", "todo", "done"}:             FamilyEpicDone,
	{"epic", "in_progress", "done"}:      FamilyEpicDone,
	{"epic", "todo", "cancelled"}:        FamilyTransition,
	{"epic", "in_progress", "cancelled"}: FamilyTransition,
	{"epic", "done", "cancelled"}:        FamilyTransition,
}

// HasFixedStates reports whether fid's commands own an Issue kind's moves:
// a task's and an epic's, not a subtask's.
func HasFixedStates(kind string) bool { return kind == "task" || kind == "epic" }

// MoveOwner returns the Event family whose command owns moving an Issue of
// a kind from one state to another; from "" is creation. owned is false when
// no command makes the move. A subtask is created by fid issue create in any
// state and moved by fid issue transition.
func MoveOwner(kind, from, to string) (family string, owned bool) {
	if kind == "subtask" {
		if from == "" {
			return FamilyCreate, true
		}
		return FamilyTransition, true
	}
	family, owned = fixedMoves[fixedMove{kind, from, to}]
	return family, owned
}

// StateOwners returns the set of families that move a task or an epic into
// a state, from any state.
func StateOwners(kind, to string) map[string]bool {
	owners := map[string]bool{}
	for move, family := range fixedMoves {
		if move.kind == kind && move.to == to {
			owners[family] = true
		}
	}
	return owners
}

// OwnerCommand names a move's owner in plain words.
func OwnerCommand(family string) string {
	return map[string]string{FamilyCreate: "fid issue create", FamilyStart: "fid start", FamilyHandoff: "the hand-off",
		FamilyRework: "a review round asking for changes", FamilyDone: "fid done", FamilyEpicDone: "fid done",
		FamilyTransition: "a person's cancellation"}[family]
}
