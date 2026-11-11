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
// waits for in its epic's sequence; Working is the open hold on it.
type Entry struct {
	issues.OverviewTask
	Reason  *string  `json:"reason"`
	After   *string  `json:"after"`
	Working *Working `json:"working"`
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
	after := sequenceWaits(work)
	for _, task := range append(waiting, rest...) {
		entry := Entry{OverviewTask: task.OverviewTask}
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
// waits for: the epic's next unfinished task in sequence. That task, and
// any task in the epic it waits on through blocks links, are served first:
// links beat sequence.
func sequenceWaits(work issues.Work) map[string]string {
	byID := map[string]issues.WorkTask{}
	first := map[string]string{}
	for _, task := range work.Tasks {
		byID[task.ID] = task
		if epic := epicOf(task.OverviewTask); work.Epics[epic].Sequential && first[epic] == "" {
			first[epic] = task.ID
		}
	}
	waits := map[string]string{}
	for epic, next := range first {
		servable := map[string]bool{}
		for pending := []string{next}; len(pending) > 0; pending = pending[1:] {
			id := pending[0]
			if servable[id] {
				continue
			}
			servable[id] = true
			for _, blocker := range byID[id].BlockedBy {
				if task, found := byID[blocker]; found && epicOf(task.OverviewTask) == epic {
					pending = append(pending, blocker)
				}
			}
		}
		for _, task := range work.Tasks {
			if epicOf(task.OverviewTask) == epic && !servable[task.ID] {
				waits[task.ID] = next
			}
		}
	}
	return waits
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
