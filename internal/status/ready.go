package status

import (
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/next"
)

type workingTask struct {
	issues.OverviewTask
	Working *holds.Working `json:"working"`
	Kept    *holds.Kept    `json:"kept"`
}

type section[T any] struct {
	Epics []epic[T] `json:"epics"`
	Tasks []T       `json:"tasks"`
}

type epic[T any] struct {
	issues.OverviewEpic
	Tasks []T `json:"tasks"`
}

// mapSection adds read-only details without changing status's grouping or order.
func mapSection[T any](source issues.OverviewSection, task func(issues.OverviewTask) T) section[T] {
	result := section[T]{Epics: []epic[T]{}, Tasks: []T{}}
	for _, item := range source.Epics {
		listed := epic[T]{OverviewEpic: item, Tasks: []T{}}
		for _, child := range item.Tasks {
			listed.Tasks = append(listed.Tasks, task(child))
		}
		result.Epics = append(result.Epics, listed)
	}
	for _, item := range source.Tasks {
		result.Tasks = append(result.Tasks, task(item))
	}
	return result
}

func (result Result) readyFields() any {
	listed := mapSection(result.Ready, func(task issues.OverviewTask) next.Entry { return result.entries[task.ID] })
	return struct {
		Next *string `json:"next"`
		section[next.Entry]
	}{Next: result.Next, section: listed}
}

// readyVerdict uses next's verdict, with status's two display differences:
// a task already held here is working here, and other eligible rows are blank.
func (result Result) readyVerdict(task issues.OverviewTask) (string, display.Role) {
	entry := result.entries[task.ID]
	if hold := entry.Working; hold != nil && hold.Here && !hold.Lapsed {
		return "working here", display.RoleMuted
	}
	served := result.Next != nil && task.ID == *result.Next
	if !served && entry.Reason == nil {
		return "", display.RoleMuted
	}
	return next.Verdict(entry, served)
}

func (result Result) startedVerdict(task issues.OverviewTask) string {
	if kept := result.entries[task.ID].Kept; kept != nil {
		if kept.Mode == "waiting" {
			return "waiting for you"
		}
		return "interrupted in " + filepath.Base(kept.Worktree)
	}
	return display.StateText(task.Status)
}
