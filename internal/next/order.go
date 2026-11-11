package next

import (
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
)

// What fid next hands out: a new task, the task this worktree already
// holds or was interrupted on, or nothing.
const (
	KindNew    = "new"
	KindResume = "resume"
	KindNone   = "none"
)

// Why fid next passes over a task. They are typed but informational: a loop
// never branches on them, and later tasks add more.
const (
	ReasonHeldElsewhere        = "held-elsewhere"
	ReasonInterruptedElsewhere = "interrupted-elsewhere"
	ReasonEpicClosed           = "epic-closed"
	ReasonEpicBlocked          = "epic-blocked"
	ReasonNotNextInSequence    = "not-next-in-sequence"
	ReasonBlocked              = "blocked"
)

// Entry is a task in the order, in fid status's task shape. Reason is why
// fid next passes over it, nil when nothing does; After names the task it
// waits for in its epic's sequence; Working is the open hold on it;
// EpicBlockedBy lists its epic's open blockers.
type Entry struct {
	issues.OverviewTask
	Reason        *string  `json:"reason"`
	After         *string  `json:"after"`
	Working       *Working `json:"working"`
	EpicBlockedBy []string `json:"epic_blocked_by"`
}

// Working is an open hold as the order shows it. Here says it is this
// worktree's; Lapsed says it had expired when the order was read.
type Working struct {
	HoldID    string `json:"hold_id"`
	Worktree  string `json:"worktree"`
	Here      bool   `json:"here"`
	Actor     string `json:"actor"`
	Since     string `json:"since"`
	ExpiresAt string `json:"expires_at"`
	Lapsed    bool   `json:"lapsed"`
}

// order is what fid next would hand one worktree.
type order struct {
	kind  string
	queue []Entry
	// served indexes the task handed out in queue, or is -1.
	served int
	// own is the worktree's open hold, if it has one.
	own *holds.Hold
	// held is each task's open hold.
	held map[string]holds.Hold
}

// plan walks the order for the worktree worktreeID at now, and serves the
// first task nothing stops:
//  1. this worktree's own hold, while it is live or its task has started
//     (a lapsed hold on a started task means the task was interrupted);
//  2. tasks that other open work waits on;
//  3. the other unstarted tasks, in fid status's order.
//
// The queue holds every unstarted task and every task with an open hold.
func plan(work issues.Work, open []holds.Hold, worktreeID string, now time.Time) order {
	planned := order{queue: []Entry{}, served: -1, held: map[string]holds.Hold{}}
	for _, hold := range open {
		planned.held[hold.TaskID] = hold
		if hold.WorktreeID == worktreeID {
			own := hold
			planned.own = &own
		}
	}
	resume := ""
	waiting := []issues.WorkTask{}
	rest := []issues.WorkTask{}
	for _, task := range work.Tasks {
		_, held := planned.held[task.ID]
		switch {
		case planned.own != nil && task.ID == planned.own.TaskID && (!planned.own.Lapsed(now) || task.StatusCategory == "started"):
			resume = task.ID
			waiting = append([]issues.WorkTask{task}, waiting...)
		case task.StatusCategory != "unstarted" && !held:
		case task.Unblocks:
			waiting = append(waiting, task)
		default:
			rest = append(rest, task)
		}
	}
	after := sequenceWaits(work, planned.held, now)
	for _, task := range append(waiting, rest...) {
		entry := Entry{OverviewTask: task.OverviewTask, EpicBlockedBy: append([]string{}, work.Epics[epicOf(task.OverviewTask)].BlockedBy...)}
		if hold, held := planned.held[task.ID]; held {
			entry.Working = working(hold, worktreeID, now)
		}
		if task.ID != resume {
			entry.Reason, entry.After = skip(work, task, entry.Working, after)
		}
		if planned.served < 0 && entry.Reason == nil {
			planned.served = len(planned.queue)
		}
		planned.queue = append(planned.queue, entry)
	}
	switch {
	case resume != "":
		planned.kind = KindResume
	case planned.served >= 0:
		planned.kind = KindNew
	default:
		planned.kind = KindNone
	}
	return planned
}

// skip is why fid next passes over task, the most telling reason first, and
// the task it waits for in its epic's sequence. A hold here never stops a
// task: a live one is resumed, and a lapsed one frees an unstarted task.
func skip(work issues.Work, task issues.WorkTask, hold *Working, after map[string]string) (*string, *string) {
	epic, inEpic := work.Epics[epicOf(task.OverviewTask)]
	reason := ""
	switch {
	case hold != nil && !hold.Here && !hold.Lapsed:
		reason = ReasonHeldElsewhere
	case hold != nil && !hold.Here && task.StatusCategory == "started":
		reason = ReasonInterruptedElsewhere
	case inEpic && (epic.StatusCategory == "completed" || epic.StatusCategory == "cancelled"):
		reason = ReasonEpicClosed
	case inEpic && len(epic.BlockedBy) > 0:
		reason = ReasonEpicBlocked
	case after[task.ID] != "":
		next := after[task.ID]
		return pointer(ReasonNotNextInSequence), &next
	case len(task.BlockedBy) > 0:
		reason = ReasonBlocked
	default:
		return nil, nil
	}
	return &reason, nil
}

// sequenceWaits maps each task a sequential epic holds back to the task it
// waits for. A sequential epic serves one task at a time, across
// worktrees. While one of its tasks is held (a live hold, or a lapsed one
// on a started task), every other task waits for it. Otherwise the epic
// serves its next unfinished task in sequence or, links beating sequence,
// the first task in sequence that one waits on through blocks links in the
// epic and that nothing blocks; the rest wait for it. A task on that chain
// that is blocked isn't held back by the sequence, so it shows as blocked;
// when the chain has nothing to serve, the rest wait for the next task.
func sequenceWaits(work issues.Work, held map[string]holds.Hold, now time.Time) map[string]string {
	byID := map[string]issues.WorkTask{}
	inEpic := map[string][]issues.WorkTask{}
	for _, task := range work.Tasks {
		byID[task.ID] = task
		if epic := epicOf(task.OverviewTask); work.Epics[epic].Sequential {
			inEpic[epic] = append(inEpic[epic], task)
		}
	}
	waits := map[string]string{}
	for epic, tasks := range inEpic {
		serving, chain := "", map[string]bool{}
		for _, task := range tasks {
			if hold, found := held[task.ID]; found && (!hold.Lapsed(now) || task.StatusCategory == "started") {
				serving = task.ID
				break
			}
		}
		if serving == "" {
			chain = blockerChain(tasks[0].ID, epic, byID)
			serving = tasks[0].ID
			for _, task := range tasks {
				if chain[task.ID] && task.StatusCategory == "unstarted" && len(task.BlockedBy) == 0 {
					serving = task.ID
					break
				}
			}
		}
		for _, task := range tasks {
			if blocked := chain[task.ID] && len(task.BlockedBy) > 0; task.ID != serving && !blocked {
				waits[task.ID] = serving
			}
		}
	}
	return waits
}

// blockerChain is next and every task in epic it waits on through blocks
// links, directly or through others.
func blockerChain(next, epic string, byID map[string]issues.WorkTask) map[string]bool {
	chain := map[string]bool{}
	for pending := []string{next}; len(pending) > 0; pending = pending[1:] {
		id := pending[0]
		if chain[id] {
			continue
		}
		chain[id] = true
		for _, blocker := range byID[id].BlockedBy {
			if task, found := byID[blocker]; found && epicOf(task.OverviewTask) == epic {
				pending = append(pending, blocker)
			}
		}
	}
	return chain
}

func working(hold holds.Hold, worktreeID string, now time.Time) *Working {
	return &Working{HoldID: hold.ID, Worktree: hold.Worktree, Here: hold.WorktreeID == worktreeID, Actor: hold.Actor,
		Since: holds.Timestamp(hold.Since), ExpiresAt: holds.Timestamp(hold.ExpiresAt), Lapsed: hold.Lapsed(now)}
}

func epicOf(task issues.OverviewTask) string {
	if task.EpicID == nil {
		return ""
	}
	return *task.EpicID
}

func pointer(value string) *string {
	return &value
}
