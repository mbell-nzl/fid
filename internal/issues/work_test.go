package issues

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func readWork(t *testing.T, workspace string) Work {
	t.Helper()
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)
	work, err := ReadWork(conn)
	if err != nil {
		t.Fatalf("ReadWork() error = %v", err)
	}
	return work
}

// ReadWork lists the open tasks in fid status's order, started ones
// included, with each epic's scheduling, progress and open blockers, and
// which tasks other open work waits on.
func TestReadWorkListsOpenTasksInStatusOrder(t *testing.T) {
	t.Parallel()
	n := intPointer
	rows := []overviewRow{
		{"E1", "epic", "", "todo", n(1), nil, 1},
		{"T1", "task", "E1", "todo", n(1), nil, 2},
		{"T2", "task", "E1", "todo", n(2), nil, 3},
		{"E2", "epic", "", "todo", n(2), nil, 4},
		{"T3", "task", "E2", "todo", nil, n(2), 5},
		{"T4", "task", "E2", "in_progress", nil, n(1), 6},
		{"T5", "task", "E2", "done", nil, n(1), 7},
		{"T6", "task", "E2", "cancelled", nil, nil, 8},
		{"X1", "subtask", "T3", "todo", nil, nil, 9},
		{"E3", "epic", "", "closed", n(3), nil, 10},
		{"T7", "task", "E3", "todo", nil, nil, 11},
		{"T8", "task", "", "todo", nil, n(1), 12},
		{"T9", "task", "", "in_progress", nil, nil, 13},
		{"T10", "task", "", "todo", nil, n(3), 14},
	}
	links := [][3]string{
		{"T8", "blocks", "T2"},     // T8 holds up open work
		{"T10", "blocks", "T5"},    // only finished work waits on T10
		{"T3", "blocks", "E1"},     // an epic waits on T3
		{"E1", "blocked_by", "T3"}, // the same pair from the other end
		{"T1", "blocked_by", "T6"}, // a cancelled blocker no longer blocks
		{"X1", "blocks", "T10"},    // any issue may block
		{"T9", "blocks", "X1"},     // only a subtask waits on T9
		{"T5", "blocks", "T2"},     // a finished blocker no longer blocks
		{"T4", "blocked_by", "T2"}, // a started task may wait too
		{"E2", "blocked_by", "T5"}, // a finished blocker of an epic
		{"T7", "blocked_by", "T4"}, // a task in a closed epic still waits
		{"T2", "blocked_by", "T8"}, // the T8 pair from the other end
		{"T1", "blocks", "T7"},     // T1 holds up T7
	}
	work := readWork(t, seedOverview(t, rows, links))
	describe := func(task WorkTask) string {
		epic := "-"
		if task.EpicID != nil {
			epic = *task.EpicID
		}
		return fmt.Sprintf("%s %s %s blocked:%s unblocks:%t", task.ID, task.StatusCategory, epic, strings.Join(task.BlockedBy, "+"), task.Unblocks)
	}
	got := []string{}
	for _, task := range work.Tasks {
		got = append(got, describe(task))
	}
	// E2 is in progress (T4 started), so its tasks come before E1's; E3 is
	// closed, so T7 joins the tasks with no epic.
	want := []string{
		"T4 started E2 blocked:T2 unblocks:true",
		"T3 unstarted E2 blocked: unblocks:true",
		"T1 unstarted E1 blocked: unblocks:true",
		"T2 unstarted E1 blocked:T8 unblocks:true",
		"T8 unstarted - blocked: unblocks:true",
		"T10 unstarted - blocked:X1 unblocks:false",
		"T7 unstarted E3 blocked:T1+T4 unblocks:false",
		"T9 started - blocked: unblocks:false",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("work =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	for id, want := range map[string]string{
		"E1": "todo unstarted sequential 0/2 blocked:T3",
		"E2": "todo unstarted prioritised 1/3 blocked:",
		"E3": "closed completed prioritised 0/1 blocked:",
	} {
		epic := work.Epics[id]
		scheduling := "prioritised"
		if epic.Sequential {
			scheduling = "sequential"
		}
		if got := fmt.Sprintf("%s %s %s %d/%d blocked:%s", epic.Status, epic.StatusCategory, scheduling, epic.Progress.Done, epic.Progress.Total, strings.Join(epic.BlockedBy, "+")); epic.ID != id || got != want {
			t.Errorf("epic %s = %s %s, want %s", id, epic.ID, got, want)
		}
	}
}
