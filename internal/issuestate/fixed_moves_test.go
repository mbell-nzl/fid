package issuestate

import "testing"

// Every move of a task's or an epic's state has one owner; no command makes
// any other. A subtask's moves all belong to fid issue transition.
func TestMoveOwnerNamesEachCommandsMoves(t *testing.T) {
	t.Parallel()
	owned := map[[3]string]string{
		{"task", "", "todo"}:                 FamilyCreate,
		{"task", "todo", "in_progress"}:      FamilyStart,
		{"task", "in_progress", "in_review"}: FamilyHandoff,
		{"task", "in_review", "in_progress"}: FamilyRework,
		{"task", "in_review", "done"}:        FamilyDone,
		{"epic", "", "todo"}:                 FamilyCreate,
		{"epic", "todo", "in_progress"}:      FamilyStart,
		{"epic", "todo", "done"}:             FamilyEpicDone,
		{"epic", "in_progress", "done"}:      FamilyEpicDone,
	}
	for _, kind := range []string{"task", "epic"} {
		for _, from := range []string{"todo", "in_progress", "in_review"} {
			owned[[3]string{kind, from, "cancelled"}] = FamilyTransition
		}
	}
	delete(owned, [3]string{"epic", "in_review", "cancelled"})
	states := []string{"", "todo", "in_progress", "in_review", "done", "cancelled"}
	for _, kind := range []string{"task", "epic"} {
		for _, from := range states {
			for _, to := range states {
				family, ok := MoveOwner(kind, from, to)
				want, wantOK := owned[[3]string{kind, from, to}]
				if family != want || ok != wantOK {
					t.Errorf("MoveOwner(%s, %q, %s) = %q, %t; want %q, %t", kind, from, to, family, ok, want, wantOK)
				}
			}
		}
	}
	for _, move := range [][2]string{{"", "done"}, {"todo", "done"}, {"done", "todo"}} {
		if family, ok := MoveOwner("subtask", move[0], move[1]); !ok || family != map[bool]string{true: FamilyCreate, false: FamilyTransition}[move[0] == ""] {
			t.Errorf("MoveOwner(subtask, %q, %s) = %q, %t", move[0], move[1], family, ok)
		}
	}
	if !HasFixedStates("task") || !HasFixedStates("epic") || HasFixedStates("subtask") {
		t.Error("only tasks and epics have fixed states")
	}
}
