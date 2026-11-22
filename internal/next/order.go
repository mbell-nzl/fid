package next

import (
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

// What fid next hands out: a new task, the task this worktree already
// holds or was interrupted on, or nothing.
const (
	KindNew    = "new"
	KindResume = "resume"
	KindRework = "rework"
	KindReview = "review"
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

const ReasonWaitingForYou = "waiting-for-you"
const ReasonBuiltByActor = "built-by-actor"
const ReasonReworkCapReached = "rework-cap-reached"

// Entry is a task in the order, in fid status's task shape. Reason is why
// fid next passes over it, nil when nothing does; After names the task it
// waits for in its epic's sequence; Working is the open hold on it;
// EpicBlockedBy lists its epic's open blockers.
type Entry struct {
	review bool
	issues.OverviewTask
	Reason        *string     `json:"reason"`
	After         *string     `json:"after"`
	Working       *Working    `json:"working"`
	Kept          *holds.Kept `json:"kept"`
	EpicBlockedBy []string    `json:"epic_blocked_by"`
}

// Working is an open hold as the order shows it. Here says it is this
// worktree's; Lapsed says it had expired when the order was read.
type Working = holds.Working

// order is what fid next would hand one worktree.
type order struct {
	kind  string
	queue []Entry
	// served indexes the task handed out in queue, or is -1.
	served int
	// own is the worktree's open hold, if it has one.
	own *holds.Hold
	// held is each task's open hold.
	held      map[string]holds.Hold
	reviews   map[string]reviewCandidate
	ownReview *reviewrecord.Claim
}

// plan serves this worktree's build resume or rework, then its live review,
// then due reviews, unblocking work and other new work. Kept builder work
// waiting for a person or review occupies its worktree. Review reservations
// occupy a reviewer worktree while live, even if comparison objects are lost.
// The queue retains every unstarted task and every task with an open hold.
func plan(work issues.Work, open []holds.Hold, worktreeID string, now time.Time, reviewPlans ...map[string]reviewCandidate) order {
	var reviews map[string]reviewCandidate
	if len(reviewPlans) > 0 {
		reviews = reviewPlans[0]
	}
	return planInScope(work, open, worktreeID, now, reviews, "")
}

func planInScope(work issues.Work, open []holds.Hold, worktreeID string, now time.Time, candidates map[string]reviewCandidate, scope string) order {
	planned := order{queue: []Entry{}, served: -1, held: map[string]holds.Hold{}}
	planned.reviews = candidates
	outside := map[string]bool{}
	for _, task := range work.Tasks {
		outside[task.ID] = scope != "" && epicOf(task.OverviewTask) != scope
	}
	for _, candidate := range planned.reviews {
		if candidate.claim != nil && candidate.claim.WorktreeID == worktreeID {
			planned.ownReview = candidate.claim
		}
	}
	for _, hold := range open {
		planned.held[hold.TaskID] = hold
		if hold.WorktreeID == worktreeID {
			own := hold
			planned.own = &own
		}
	}
	resume := ""
	waiting := []issues.WorkTask{}
	unblocking := []issues.WorkTask{}
	rest := []issues.WorkTask{}
	reviews := []issues.WorkTask{}
	for _, task := range work.Tasks {
		_, held := planned.held[task.ID]
		switch {
		case planned.own != nil && task.ID == planned.own.TaskID && (planned.own.Live(now) || planned.own.Kept || planned.own.Waiting() || task.StatusCategory == "started"):
			resume = task.ID
			waiting = append([]issues.WorkTask{task}, waiting...)
		case planned.reviews[task.ID].claim != nil && planned.reviews[task.ID].claim.WorktreeID == worktreeID:
			waiting = append([]issues.WorkTask{task}, waiting...)
		case planned.reviews[task.ID].work != nil:
			reviews = append(reviews, task)
		case task.StatusCategory != "unstarted" && !held:
		case task.Unblocks:
			unblocking = append(unblocking, task)
		default:
			rest = append(rest, task)
		}
	}
	after := sequenceWaits(work, planned.held, now)
	for _, task := range append(append(append(waiting, reviews...), unblocking...), rest...) {
		entry := Entry{OverviewTask: task.OverviewTask, EpicBlockedBy: append([]string{}, work.Epics[epicOf(task.OverviewTask)].BlockedBy...)}
		if hold, held := planned.held[task.ID]; held {
			entry.Working = working(hold, worktreeID, now)
			if task.StatusCategory == "started" || hold.Kept || hold.Waiting() {
				entry.Kept = holds.DescribeKept(hold, worktreeID, now)
			}
		}
		if candidate := planned.reviews[task.ID]; candidate.work != nil && task.ID != resume {
			entry.review, entry.Working = true, nil
			if candidate.claim != nil {
				entry.Working = working(*claimHold(*candidate.claim), worktreeID, now)
			}
			if candidate.claim != nil && candidate.claim.WorktreeID != worktreeID {
				entry.Reason = pointer(ReasonHeldElsewhere)
			} else if candidate.built {
				entry.Reason = pointer(ReasonBuiltByActor)
			}
		} else if task.ID != resume || (planned.own != nil && planned.own.Waiting()) {
			if hold, ok := planned.held[task.ID]; ok && hold.InReview() {
				entry.Reason = pointer(ReasonHeldElsewhere)
			} else if ok && hold.Waiting() {
				entry.Reason = pointer(ReasonWaitingForYou)
			} else {
				entry.Reason, entry.After = skip(work, task, entry.Working, entry.Kept, after)
			}
		}
		if planned.reviews[task.ID].capped {
			entry.Reason = pointer(ReasonReworkCapReached)
		}
		if outside[task.ID] {
			entry.Reason, entry.After = pointer(ReasonOutsideScope), nil
		}
		if planned.served < 0 && entry.Reason == nil {
			planned.served = len(planned.queue)
		}
		planned.queue = append(planned.queue, entry)
	}
	// Scope cannot abandon occupied work and start a second task here.
	if outside[resume] || planned.ownReview != nil && outside[planned.ownReview.TaskID] {
		planned.kind, planned.served = KindNone, -1
		return planned
	}
	if planned.ownReview != nil && planned.reviews[planned.ownReview.TaskID].work == nil {
		planned.kind, planned.served = KindNone, -1
		return planned
	}
	if planned.own != nil && (planned.own.Waiting() || planned.own.InReview() || planned.reviews[planned.own.TaskID].capped) {
		planned.kind, planned.served = KindNone, -1
		return planned
	}
	switch {
	case resume != "" && planned.own.Rework():
		planned.kind = KindRework
	case resume != "":
		planned.kind = KindResume
	case planned.served >= 0 && planned.reviews[planned.queue[planned.served].ID].work != nil:
		planned.kind = KindReview
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
func skip(work issues.Work, task issues.WorkTask, hold *Working, kept *holds.Kept, after map[string]string) (*string, *string) {
	epic, inEpic := work.Epics[epicOf(task.OverviewTask)]
	reason := ""
	switch {
	case hold != nil && !hold.Here && !hold.Lapsed:
		reason = ReasonHeldElsewhere
	case kept != nil && !kept.Here:
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
			if hold, found := held[task.ID]; found && (hold.Live(now) || hold.Kept || hold.Waiting() || task.StatusCategory == "started") {
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
	return holds.Describe(hold, worktreeID, now)
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
