package next

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
)

var orderNow = time.Date(2026, 11, 10, 9, 0, 0, 0, time.UTC)

const here = "sha256:here"

// task is an open task for the order: its ID, epic ("" for none), status
// category, open blockers, and whether open work waits on it.
func task(id, epic, category string, unblocks bool, blockedBy ...string) issues.WorkTask {
	listed := issues.WorkTask{OverviewTask: issues.OverviewTask{ID: id, Title: "Title of " + id, Status: map[string]string{"unstarted": "todo", "started": "in_progress"}[category],
		StatusCategory: category, BlockedBy: append([]string{}, blockedBy...)}, Unblocks: unblocks}
	if epic != "" {
		listed.EpicID = &epic
	}
	return listed
}

// epic is an epic: its ID, status category, whether it is sequential, and
// its open blockers.
func epic(id, category string, sequential bool, blockedBy ...string) issues.WorkEpic {
	return issues.WorkEpic{ID: id, Title: "Title of " + id, StatusCategory: category, Sequential: sequential, BlockedBy: blockedBy}
}

// hold is an open hold on task in worktree, taken ago before orderNow.
func hold(id, task, worktree string, ago time.Duration) holds.Hold {
	return holds.New(id, task, worktree, "/w/"+strings.TrimPrefix(worktree, "sha256:"), "agent:builder", orderNow.Add(-ago), nil)
}

// describe writes the order as its kind and served task, then one line per
// entry: "T1 blocked" or "T2 not-next-in-sequence after T1" or
// "T3 - working HLD-1 here lapsed".
func describe(planned order) []string {
	served := "-"
	if planned.served >= 0 {
		served = planned.queue[planned.served].ID
	}
	lines := []string{planned.kind + " " + served}
	for _, entry := range planned.queue {
		line := entry.ID + " -"
		if entry.Reason != nil {
			line = entry.ID + " " + *entry.Reason
		}
		if entry.After != nil {
			line += " after " + *entry.After
		}
		if entry.Working != nil {
			line += " working " + entry.Working.HoldID
			if entry.Working.Here {
				line += " here"
			}
			if entry.Working.Lapsed {
				line += " lapsed"
			}
		}
		lines = append(lines, line)
	}
	return lines
}

func TestPlanWalksTheOrderAndSkipsWithTypedReasons(t *testing.T) {
	t.Parallel()
	live, lapsed := 10*time.Minute, 2*time.Hour
	for _, test := range []struct {
		name  string
		epics []issues.WorkEpic
		tasks []issues.WorkTask
		holds []holds.Hold
		want  []string
	}{
		{"status's order, first served", nil,
			[]issues.WorkTask{task("T1", "", "unstarted", false), task("T2", "", "unstarted", false)}, nil,
			[]string{"new T1", "T1 -", "T2 -"}},
		{"nothing to do", nil, nil, nil, []string{"none -"}},
		{"unblocking work first, in order", nil,
			[]issues.WorkTask{task("T1", "", "unstarted", false), task("T2", "", "unstarted", true), task("T3", "", "unstarted", true, "T9")},
			nil, []string{"new T2", "T2 -", "T3 blocked", "T1 -"}},
		{"a started task with no hold isn't handed out", nil,
			[]issues.WorkTask{task("T1", "", "started", true), task("T2", "", "unstarted", false)}, nil,
			[]string{"new T2", "T2 -"}},
		{"blocked tasks, then nothing", nil,
			[]issues.WorkTask{task("T1", "", "unstarted", false, "T8"), task("T2", "", "unstarted", false, "T8", "T9")}, nil,
			[]string{"none -", "T1 blocked", "T2 blocked"}},
		{"closed and blocked epics", []issues.WorkEpic{epic("E1", "completed", false), epic("E2", "cancelled", false), epic("E3", "unstarted", false, "T9"), epic("E4", "started", false)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false), task("T2", "E2", "unstarted", false), task("T3", "E3", "unstarted", false), task("T4", "E4", "unstarted", false)}, nil,
			[]string{"new T4", "T1 epic-closed", "T2 epic-closed", "T3 epic-blocked", "T4 -"}},
		{"a sequential epic serves one task at a time", []issues.WorkEpic{epic("E1", "unstarted", true), epic("E2", "unstarted", false)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false), task("T2", "E1", "unstarted", false), task("T3", "E2", "unstarted", false), task("T4", "E2", "unstarted", false)}, nil,
			[]string{"new T1", "T1 -", "T2 not-next-in-sequence after T1", "T3 -", "T4 -"}},
		{"a started task holds its sequence", []issues.WorkEpic{epic("E1", "started", true)},
			[]issues.WorkTask{task("T1", "E1", "started", false), task("T2", "E1", "unstarted", false, "T1")},
			[]holds.Hold{hold("HLD-1", "T1", "sha256:wt2", live)},
			[]string{"none -", "T1 held-elsewhere working HLD-1", "T2 not-next-in-sequence after T1"}},
		{"links beat sequence", []issues.WorkEpic{epic("E1", "unstarted", true)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false, "T3"), task("T2", "E1", "unstarted", false), task("T3", "E1", "unstarted", true, "T4"), task("T4", "E1", "unstarted", true)}, nil,
			[]string{"new T4", "T3 blocked", "T4 -", "T1 blocked", "T2 not-next-in-sequence after T4"}},
		{"links that branch still serve one task of a sequential epic", []issues.WorkEpic{epic("E1", "unstarted", true)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false, "T2", "T3"), task("T2", "E1", "unstarted", true), task("T3", "E1", "unstarted", true), task("T4", "E1", "unstarted", false)}, nil,
			[]string{"new T2", "T2 -", "T3 not-next-in-sequence after T2", "T1 blocked", "T4 not-next-in-sequence after T2"}},
		{"a task held elsewhere occupies its sequential epic", []issues.WorkEpic{epic("E1", "unstarted", true)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false, "T2", "T3"), task("T2", "E1", "unstarted", true), task("T3", "E1", "unstarted", true), task("T4", "E1", "unstarted", false)},
			[]holds.Hold{hold("HLD-1", "T2", "sha256:wt2", live)},
			[]string{"none -", "T2 held-elsewhere working HLD-1", "T3 not-next-in-sequence after T2", "T1 not-next-in-sequence after T2", "T4 not-next-in-sequence after T2"}},
		{"a task interrupted elsewhere occupies its sequential epic", []issues.WorkEpic{epic("E1", "started", true)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false, "T2", "T3"), task("T2", "E1", "unstarted", true), task("T3", "E1", "started", true), task("T4", "E1", "unstarted", false)},
			[]holds.Hold{hold("HLD-1", "T3", "sha256:wt2", lapsed)},
			[]string{"none -", "T2 not-next-in-sequence after T3", "T3 interrupted-elsewhere working HLD-1 lapsed", "T1 not-next-in-sequence after T3", "T4 not-next-in-sequence after T3"}},
		{"a lapsed hold on an unstarted task doesn't occupy its sequential epic", []issues.WorkEpic{epic("E1", "unstarted", true)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false, "T2", "T3"), task("T2", "E1", "unstarted", true), task("T3", "E1", "unstarted", true), task("T4", "E1", "unstarted", false)},
			[]holds.Hold{hold("HLD-1", "T3", "sha256:wt2", lapsed)},
			[]string{"new T2", "T2 -", "T3 not-next-in-sequence after T2 working HLD-1 lapsed", "T1 blocked", "T4 not-next-in-sequence after T2"}},
		{"a hold here occupies its sequential epic", []issues.WorkEpic{epic("E1", "unstarted", true)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false, "T2", "T3"), task("T2", "E1", "unstarted", true), task("T3", "E1", "unstarted", true), task("T4", "E1", "unstarted", false)},
			[]holds.Hold{hold("HLD-1", "T3", here, live)},
			[]string{"resume T3", "T3 - working HLD-1 here", "T2 not-next-in-sequence after T3", "T1 not-next-in-sequence after T3", "T4 not-next-in-sequence after T3"}},
		{"a blocker outside the epic doesn't beat sequence", []issues.WorkEpic{epic("E1", "unstarted", true)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false, "T5"), task("T2", "E1", "unstarted", false), task("T5", "", "unstarted", true)}, nil,
			[]string{"new T5", "T5 -", "T1 blocked", "T2 not-next-in-sequence after T1"}},
		{"held and interrupted elsewhere; a lapsed unstarted hold frees its task", nil,
			[]issues.WorkTask{task("T1", "", "unstarted", false), task("T2", "", "started", false), task("T3", "", "unstarted", false, "T9"), task("T4", "", "unstarted", false)},
			[]holds.Hold{hold("HLD-1", "T1", "sha256:wt2", live), hold("HLD-2", "T2", "sha256:wt3", lapsed), hold("HLD-3", "T3", "sha256:wt4", live), hold("HLD-4", "T4", "sha256:wt5", lapsed)},
			[]string{"new T4", "T1 held-elsewhere working HLD-1", "T2 interrupted-elsewhere working HLD-2 lapsed", "T3 held-elsewhere working HLD-3", "T4 - working HLD-4 lapsed"}},
		{"a live hold here is resumed first, whatever stops it", nil,
			[]issues.WorkTask{task("T1", "", "unstarted", true), task("T2", "", "unstarted", false, "T9")},
			[]holds.Hold{hold("HLD-1", "T2", here, live)},
			[]string{"resume T2", "T2 - working HLD-1 here", "T1 -"}},
		{"a lapsed hold here on a started task is resumed", nil,
			[]issues.WorkTask{task("T1", "", "unstarted", false), task("T2", "", "started", false)},
			[]holds.Hold{hold("HLD-1", "T2", here, lapsed)},
			[]string{"resume T2", "T2 - working HLD-1 here lapsed", "T1 -"}},
		{"a lapsed hold here on an unstarted task frees it", nil,
			[]issues.WorkTask{task("T1", "", "unstarted", false), task("T2", "", "unstarted", true, "T9")},
			[]holds.Hold{hold("HLD-1", "T2", here, lapsed)},
			[]string{"new T1", "T2 blocked working HLD-1 here lapsed", "T1 -"}},
		{"a lapse is read from the clock", nil,
			[]issues.WorkTask{task("T1", "", "unstarted", false), task("T2", "", "unstarted", false)},
			[]holds.Hold{hold("HLD-1", "T1", "sha256:wt2", holds.Lifetime)},
			[]string{"new T1", "T1 - working HLD-1 lapsed", "T2 -"}},
		{"held elsewhere beats every other reason", []issues.WorkEpic{epic("E1", "completed", true)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false, "T9")},
			[]holds.Hold{hold("HLD-1", "T1", "sha256:wt2", live)},
			[]string{"none -", "T1 held-elsewhere working HLD-1"}},
		{"the sequence is named before a blocker", []issues.WorkEpic{epic("E1", "unstarted", true)},
			[]issues.WorkTask{task("T1", "E1", "unstarted", false), task("T2", "E1", "unstarted", false, "T9")}, nil,
			[]string{"new T1", "T1 -", "T2 not-next-in-sequence after T1"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			work := issues.Work{Epics: map[string]issues.WorkEpic{}, Tasks: test.tasks}
			for _, epic := range test.epics {
				work.Epics[epic.ID] = epic
			}
			if got := describe(plan(work, test.holds, here, orderNow)); !slices.Equal(got, test.want) {
				t.Fatalf("order =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(test.want, "\n"))
			}
		})
	}
}

// The order's working entries carry the hold's times in UTC, to the second.
func TestPlanShowsEachHoldAsWorking(t *testing.T) {
	t.Parallel()
	work := issues.Work{Epics: map[string]issues.WorkEpic{}, Tasks: []issues.WorkTask{task("T1", "", "unstarted", false)}}
	taken := hold("HLD-1", "T1", "sha256:wt2", 90*time.Second+500*time.Millisecond)
	planned := plan(work, []holds.Hold{taken}, here, orderNow)
	want := Working{HoldID: "HLD-1", Worktree: "/w/wt2", Here: false, Actor: "agent:builder", Since: "2026-11-10T08:58:29Z", ExpiresAt: "2026-11-10T09:58:29Z", LastHeartbeatAt: "2026-11-10T08:58:29Z", Lapsed: false}
	if got := planned.queue[0].Working; got == nil || *got != want {
		t.Fatalf("working = %s, want %+v", fmt.Sprint(got), want)
	}
}
