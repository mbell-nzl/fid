package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

// Epic states in the overview.
const (
	EpicInProgress = "in-progress"
	EpicToDo       = "to-do"
)

type OverviewOptions struct {
	WorkspaceRoot string
	MigrationsDir string
}

// Overview is the project's open epics and tasks, for fid status. Subtasks
// aren't in it.
type Overview struct {
	// InProgress holds each in-progress epic with its started tasks, then
	// the started tasks no listed epic holds.
	InProgress OverviewSection
	// ToDo holds every unstarted task: under its in-progress epic, then
	// under its to-do epic, then those no listed epic holds.
	ToDo OverviewSection
	// Empty says the project has no epics or tasks at all.
	Empty bool
}

type OverviewSection struct {
	Epics []OverviewEpic `json:"epics"`
	Tasks []OverviewTask `json:"tasks"`
}

type OverviewEpic struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Status         string         `json:"status"`
	StatusCategory string         `json:"status_category"`
	State          string         `json:"state"`
	Sequence       *int           `json:"sequence"`
	Priority       *int           `json:"priority"`
	Progress       EpicProgress   `json:"progress"`
	Tasks          []OverviewTask `json:"tasks"`
}

// EpicProgress counts an epic's completed tasks against those not
// cancelled.
type EpicProgress struct {
	Done  int `json:"done"`
	Total int `json:"total"`
}

type OverviewTask struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Status         string   `json:"status"`
	StatusCategory string   `json:"status_category"`
	EpicID         *string  `json:"epic_id"`
	Sequence       *int     `json:"sequence"`
	Priority       *int     `json:"priority"`
	BlockedBy      []string `json:"blocked_by"`
}

type overviewEpic struct {
	OverviewEpic
	started, unstarted []OverviewTask
}

// ReadOverview sorts the project's epics and tasks into in-progress and to-do
// work. An epic whose own state is completed or cancelled is hidden. An epic
// is in progress when its own state is started, or any of its tasks is
// started or completed; otherwise it is to do. Epics, and the tasks in each
// list, are ordered by sequence, then priority, then creation and ID, with
// a missing sequence or priority last.
func ReadOverview(options OverviewOptions) (Overview, error) {
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return Overview{}, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return Overview{}, err
	}
	defer closeIssueDB(conn)
	registry, err := workflow.Load(conn)
	if err != nil {
		return Overview{}, fmt.Errorf("read overview: load workflow registry: %w", err)
	}
	blockers, err := openBlockers(conn, registry)
	if err != nil {
		return Overview{}, err
	}
	epics, tasks, err := overviewRows(conn, registry, blockers)
	if err != nil {
		return Overview{}, err
	}
	return sortOverview(epics, tasks), nil
}

func sortOverview(epics []*overviewEpic, tasks []OverviewTask) Overview {
	overview := Overview{Empty: len(epics) == 0 && len(tasks) == 0,
		InProgress: OverviewSection{Epics: []OverviewEpic{}, Tasks: []OverviewTask{}},
		ToDo:       OverviewSection{Epics: []OverviewEpic{}, Tasks: []OverviewTask{}}}
	listed := map[string]*overviewEpic{}
	for _, epic := range epics {
		if epic.StatusCategory != "completed" && epic.StatusCategory != "cancelled" {
			listed[epic.ID] = epic
		}
	}
	for _, task := range tasks {
		epic := (*overviewEpic)(nil)
		if task.EpicID != nil {
			epic = listed[*task.EpicID]
		}
		if epic != nil {
			epic.count(task)
		}
		switch {
		case task.StatusCategory == "started" && epic != nil:
			epic.started = append(epic.started, task)
		case task.StatusCategory == "started":
			overview.InProgress.Tasks = append(overview.InProgress.Tasks, task)
		case task.StatusCategory == "unstarted" && epic != nil:
			epic.unstarted = append(epic.unstarted, task)
		case task.StatusCategory == "unstarted":
			overview.ToDo.Tasks = append(overview.ToDo.Tasks, task)
		}
	}
	toDo := []OverviewEpic{}
	for _, epic := range epics {
		if listed[epic.ID] == nil {
			continue
		}
		if epic.StatusCategory == "started" || epic.Progress.Done > 0 || len(epic.started) > 0 {
			epic.State = EpicInProgress
			overview.InProgress.Epics = append(overview.InProgress.Epics, epic.with(epic.started))
			if len(epic.unstarted) > 0 {
				overview.ToDo.Epics = append(overview.ToDo.Epics, epic.with(epic.unstarted))
			}
			continue
		}
		epic.State = EpicToDo
		toDo = append(toDo, epic.with(epic.unstarted))
	}
	overview.ToDo.Epics = append(overview.ToDo.Epics, toDo...)
	return overview
}

func (epic *overviewEpic) count(task OverviewTask) {
	if task.StatusCategory != "cancelled" {
		epic.Progress.Total++
	}
	if task.StatusCategory == "completed" {
		epic.Progress.Done++
	}
}

func (epic *overviewEpic) with(tasks []OverviewTask) OverviewEpic {
	listed := epic.OverviewEpic
	listed.Tasks = append([]OverviewTask{}, tasks...)
	return listed
}

func overviewRows(conn *sql.DB, registry *workflow.Registry, blockers map[string][]string) ([]*overviewEpic, []OverviewTask, error) {
	rows, err := conn.Query(`
		SELECT id, kind, parent_id, title, status, workflow_id, sequence, priority
		FROM issues
		WHERE kind IN ('epic', 'task')
		ORDER BY sequence IS NULL, sequence, priority IS NULL, priority, created_at, id`)
	if err != nil {
		return nil, nil, fmt.Errorf("query overview: %w", err)
	}
	defer closeIssueRows(rows)
	epics, tasks := []*overviewEpic{}, []OverviewTask{}
	for rows.Next() {
		var (
			task               OverviewTask
			kind, workflowID   string
			parentID           sql.NullString
			sequence, priority sql.NullInt64
		)
		if err := rows.Scan(&task.ID, &kind, &parentID, &task.Title, &task.Status, &workflowID, &sequence, &priority); err != nil {
			return nil, nil, fmt.Errorf("scan overview row: %w", err)
		}
		if task.StatusCategory, err = registryStatusCategory(registry, workflowID, task.Status); err != nil {
			return nil, nil, err
		}
		task.Sequence, task.Priority = nullIntPointer(sequence), nullIntPointer(priority)
		if kind == "epic" {
			epics = append(epics, &overviewEpic{OverviewEpic: OverviewEpic{ID: task.ID, Title: task.Title, Status: task.Status,
				StatusCategory: task.StatusCategory, Sequence: task.Sequence, Priority: task.Priority}})
			continue
		}
		task.EpicID = nullStringPointer(parentID)
		task.BlockedBy = append([]string{}, blockers[task.ID]...)
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate overview rows: %w", err)
	}
	return epics, tasks, nil
}

// openBlockers maps each issue to the issues that block it, by ID, leaving
// out blockers that are completed or cancelled. A dependency is a blocks
// link from the blocker or a blocked_by link to it.
func openBlockers(conn *sql.DB, registry *workflow.Registry) (map[string][]string, error) {
	rows, err := conn.Query(`
		SELECT DISTINCT pair.blocked, blocker.id, blocker.workflow_id, blocker.status
		FROM (
			SELECT from_id AS blocked, to_id AS blocker FROM links
			WHERE relation = 'blocked_by' AND from_type = 'issue' AND to_type = 'issue'
			UNION
			SELECT to_id, from_id FROM links
			WHERE relation = 'blocks' AND from_type = 'issue' AND to_type = 'issue'
		) pair
		JOIN issues blocker ON blocker.id = pair.blocker
		ORDER BY pair.blocked, blocker.id`)
	if err != nil {
		return nil, fmt.Errorf("query blockers: %w", err)
	}
	defer closeIssueRows(rows)
	blockers := map[string][]string{}
	for rows.Next() {
		var blocked, blocker, workflowID, status string
		if err := rows.Scan(&blocked, &blocker, &workflowID, &status); err != nil {
			return nil, fmt.Errorf("scan blocker: %w", err)
		}
		category, err := registryStatusCategory(registry, workflowID, status)
		if err != nil {
			return nil, err
		}
		if category != "completed" && category != "cancelled" {
			blockers[blocked] = append(blockers[blocked], blocker)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate blockers: %w", err)
	}
	return blockers, nil
}
