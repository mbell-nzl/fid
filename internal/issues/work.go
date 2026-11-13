package issues

import (
	"fmt"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/workflow"
)

// Work is what fid next orders: every epic, and the open tasks in fid
// status's order.
type Work struct {
	// Epics holds every epic, closed ones included, by ID.
	Epics map[string]WorkEpic
	// Tasks are the tasks that aren't completed or cancelled, started ones
	// included: an in-progress epic's, then a to-do epic's, then those no
	// open epic holds, each list by sequence, priority, creation and ID.
	Tasks []WorkTask
}

// WorkEpic is an epic as fid next reads it. BlockedBy lists its open
// blockers' IDs, in order.
type WorkEpic struct {
	ID             string
	Title          string
	Status         string
	StatusCategory string
	Sequential     bool
	Progress       Progress
	BlockedBy      []string
}

// WorkTask is a task in fid status's task shape. Unblocks says another open
// task waits on it, or an open epic with open tasks, which wait on it
// through their epic.
type WorkTask struct {
	OverviewTask
	Unblocks bool
}

// ReadWork reads the epics and open tasks from the projection on conn.
func ReadWork(conn fiddb.Queryer) (Work, error) {
	registry, err := workflow.LoadFrom(conn)
	if err != nil {
		return Work{}, fmt.Errorf("read work: load workflow registry: %w", err)
	}
	blockers, err := openBlockers(conn, registry)
	if err != nil {
		return Work{}, err
	}
	epics, tasks, err := overviewRows(conn, registry, blockers)
	if err != nil {
		return Work{}, err
	}
	sequential, err := sequentialEpics(conn)
	if err != nil {
		return Work{}, err
	}
	work := Work{Epics: map[string]WorkEpic{}, Tasks: []WorkTask{}}
	for _, epic := range epics {
		work.Epics[epic.ID] = WorkEpic{ID: epic.ID, Title: epic.Title, Status: epic.Status, StatusCategory: epic.StatusCategory,
			Sequential: sequential[epic.ID], BlockedBy: append([]string{}, blockers[epic.ID]...)}
	}
	waitedOn, begun := map[string]bool{}, map[string]bool{}
	for _, task := range tasks {
		if isOpen(task.StatusCategory) {
			waits := task.BlockedBy
			if epic, found := work.Epics[epicOf(task)]; found && isOpen(epic.StatusCategory) {
				waits = append(append([]string{}, waits...), epic.BlockedBy...)
			}
			for _, blocker := range waits {
				waitedOn[blocker] = true
			}
		}
		if epic, found := work.Epics[epicOf(task)]; found {
			epic.Progress.count(task.StatusCategory)
			work.Epics[epic.ID] = epic
			begun[epic.ID] = begun[epic.ID] || task.StatusCategory == "started" || task.StatusCategory == "completed"
		}
	}
	listed := map[string][]WorkTask{}
	rest := []WorkTask{}
	for _, task := range tasks {
		if !isOpen(task.StatusCategory) {
			continue
		}
		listedTask := WorkTask{OverviewTask: task, Unblocks: waitedOn[task.ID]}
		if epic, found := work.Epics[epicOf(task)]; found && isOpen(epic.StatusCategory) {
			listed[epic.ID] = append(listed[epic.ID], listedTask)
			continue
		}
		rest = append(rest, listedTask)
	}
	for _, state := range []string{EpicInProgress, EpicToDo} {
		for _, epic := range epics {
			if isOpen(epic.StatusCategory) && epicState(epic.StatusCategory, begun[epic.ID]) == state {
				work.Tasks = append(work.Tasks, listed[epic.ID]...)
			}
		}
	}
	work.Tasks = append(work.Tasks, rest...)
	return work, nil
}

// isOpen reports whether an issue in category is neither completed nor
// cancelled.
func isOpen(category string) bool {
	return category != "completed" && category != "cancelled"
}

func epicOf(task OverviewTask) string {
	if task.EpicID == nil {
		return ""
	}
	return *task.EpicID
}

// sequentialEpics names the epics whose scheduling is sequential. Other
// epics, prioritised or with no scheduling, go by priority.
func sequentialEpics(conn fiddb.Queryer) (map[string]bool, error) {
	rows, err := conn.Query(`SELECT issue_id FROM epic_metadata WHERE scheduling = 'sequential'`)
	if err != nil {
		return nil, fmt.Errorf("query epic scheduling: %w", err)
	}
	defer closeIssueRows(rows)
	sequential := map[string]bool{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan epic scheduling: %w", err)
		}
		sequential[id] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate epic scheduling: %w", err)
	}
	return sequential, nil
}
