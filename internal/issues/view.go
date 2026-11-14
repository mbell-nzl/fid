package issues

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/questions"
	"github.com/mbell-nzl/fid/internal/workflow"
)

// IssueView is one issue as fid show shows it, with what surrounds it.
type IssueView struct {
	// Issue is issue.get's issue.
	Issue IssueDetail
	// State is how an epic reads, as fid status reads it, and nil for
	// tasks and subtasks.
	State *string
	// Epic is the issue's epic; a subtask's is its task's epic.
	Epic *EpicReference
	// Task is a subtask's task.
	Task *IssueReference
	// BlockedBy and Blocks are the issues on the other end of every blocks
	// link, read from either end, whatever their state, in fid status's
	// order.
	BlockedBy []IssueReference
	Blocks    []IssueReference
	// Progress counts the children; nil for a subtask, and for a task with
	// no subtasks.
	Progress *Progress
	// Children are an epic's tasks or a task's subtasks, in fid status's
	// order.
	Children   []ChildIssue
	Questions  []questions.Question
	OpenHold   *holds.Hold
	Scheduling *string
}

// IssueReference names another issue and its state.
type IssueReference struct {
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	Title          string `json:"title"`
	Status         string `json:"status"`
	StatusCategory string `json:"status_category"`
}

// EpicReference is an issue's epic, with its tasks' progress.
type EpicReference struct {
	IssueReference
	Progress Progress `json:"progress"`
}

// ChildIssue is an epic's task or a task's subtask, with the IDs of its
// blockers that aren't completed or cancelled.
type ChildIssue struct {
	IssueReference
	Sequence  *int     `json:"sequence"`
	Priority  *int     `json:"priority"`
	BlockedBy []string `json:"blocked_by"`
}

// childKinds names the kind of an issue's children.
var childKinds = map[string]string{"epic": "task", "task": "subtask"}

// ReadIssueView reads the issue whose ID is exactly options.ID, with its
// epic and task, its blocks links both ways, and its children.
func ReadIssueView(options GetOptions) (IssueView, error) {
	conn, registry, err := openIssueRead(options, "read issue view")
	if err != nil {
		return IssueView{}, err
	}
	defer closeIssueDB(conn)
	issue, err := readIssue(conn, registry, options.ID)
	if err != nil {
		return IssueView{}, err
	}
	view := IssueView{Issue: issueDetailFromIssue(issue), Children: []ChildIssue{}}
	if view.Epic, view.Task, err = issueParents(conn, registry, issue.ParentID); err != nil {
		return IssueView{}, err
	}
	if view.BlockedBy, view.Blocks, err = blockLinks(conn, registry, issue.ID); err != nil {
		return IssueView{}, err
	}
	if kind, ok := childKinds[issue.Kind]; ok {
		if view.Children, err = childIssues(conn, registry, issue.ID, kind); err != nil {
			return IssueView{}, err
		}
		blockers, err := openBlockers(conn, registry)
		if err != nil {
			return IssueView{}, err
		}
		for index := range view.Children {
			view.Children[index].BlockedBy = append([]string{}, blockers[view.Children[index].ID]...)
		}
	}
	progress, begun := childProgress(view.Children)
	if issue.Kind == "epic" || len(view.Children) > 0 {
		view.Progress = &progress
	}
	if issue.Kind == "epic" {
		state := epicState(issue.StatusCategory, begun)
		view.State = &state
	}
	view.Questions, err = questions.Read(conn, issue.ID)
	if err != nil {
		return IssueView{}, err
	}
	if issue.Kind == "task" {
		open, err := holds.Open(conn)
		if err != nil {
			return IssueView{}, err
		}
		for _, hold := range open {
			if hold.TaskID == issue.ID {
				view.OpenHold = &hold
				break
			}
		}
	}
	if issue.Kind == "epic" {
		var scheduling sql.NullString
		err := conn.QueryRow(`SELECT scheduling FROM epic_metadata WHERE issue_id = ?`, issue.ID).Scan(&scheduling)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return IssueView{}, fmt.Errorf("read epic scheduling: %w", err)
		}
		view.Scheduling = nullStringPointer(scheduling)
	}
	return view, nil
}

// issueParents reads the epic and task above an issue whose parent is
// parentID: a task's epic, or a subtask's task and that task's epic.
func issueParents(conn *sql.DB, registry *workflow.Registry, parentID *string) (*EpicReference, *IssueReference, error) {
	parent, grandparentID, err := issueReference(conn, registry, parentID)
	if err != nil || parent == nil {
		return nil, nil, err
	}
	var task *IssueReference
	if parent.Kind == "task" {
		task = parent
		if parent, _, err = issueReference(conn, registry, grandparentID); err != nil || parent == nil {
			return nil, task, err
		}
	}
	if parent.Kind != "epic" {
		return nil, task, nil
	}
	tasks, err := childIssues(conn, registry, parent.ID, "task")
	if err != nil {
		return nil, nil, err
	}
	progress, _ := childProgress(tasks)
	return &EpicReference{IssueReference: *parent, Progress: progress}, task, nil
}

// issueReference reads the issue id names, if any, and its parent's ID.
func issueReference(conn *sql.DB, registry *workflow.Registry, id *string) (*IssueReference, *string, error) {
	if id == nil {
		return nil, nil, nil
	}
	var (
		reference  IssueReference
		workflowID string
		parentID   sql.NullString
	)
	err := conn.QueryRow(`SELECT id, kind, title, status, workflow_id, parent_id FROM issues WHERE id = ?`, *id).
		Scan(&reference.ID, &reference.Kind, &reference.Title, &reference.Status, &workflowID, &parentID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("query issue %q: %w", *id, err)
	}
	if reference.StatusCategory, err = registryStatusCategory(registry, workflowID, reference.Status); err != nil {
		return nil, nil, err
	}
	return &reference, nullStringPointer(parentID), nil
}

// blockLinks reads the issues that block id and those it blocks, each from
// a blocks link on one end or a blocked_by link on the other, in fid
// status's order.
func blockLinks(conn *sql.DB, registry *workflow.Registry, id string) ([]IssueReference, []IssueReference, error) {
	rows, err := conn.Query(`
		SELECT pair.direction, other.id, other.kind, other.title, other.status, other.workflow_id
		FROM (
			SELECT 'blocked_by' AS direction, to_id AS other FROM links
			WHERE relation = 'blocked_by' AND from_type = 'issue' AND to_type = 'issue' AND from_id = ?
			UNION SELECT 'blocked_by', from_id FROM links
			WHERE relation = 'blocks' AND from_type = 'issue' AND to_type = 'issue' AND to_id = ?
			UNION SELECT 'blocks', to_id FROM links
			WHERE relation = 'blocks' AND from_type = 'issue' AND to_type = 'issue' AND from_id = ?
			UNION SELECT 'blocks', from_id FROM links
			WHERE relation = 'blocked_by' AND from_type = 'issue' AND to_type = 'issue' AND to_id = ?
		) pair
		JOIN issues other ON other.id = pair.other
		ORDER BY other.sequence IS NULL, other.sequence, other.priority IS NULL, other.priority, other.created_at, other.id`, id, id, id, id)
	if err != nil {
		return nil, nil, fmt.Errorf("query blocks links for %q: %w", id, err)
	}
	defer closeIssueRows(rows)
	blockedBy, blocks := []IssueReference{}, []IssueReference{}
	for rows.Next() {
		var (
			reference             IssueReference
			direction, workflowID string
		)
		if err := rows.Scan(&direction, &reference.ID, &reference.Kind, &reference.Title, &reference.Status, &workflowID); err != nil {
			return nil, nil, fmt.Errorf("scan blocks link for %q: %w", id, err)
		}
		if reference.StatusCategory, err = registryStatusCategory(registry, workflowID, reference.Status); err != nil {
			return nil, nil, err
		}
		if direction == "blocks" {
			blocks = append(blocks, reference)
		} else {
			blockedBy = append(blockedBy, reference)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate blocks links for %q: %w", id, err)
	}
	return blockedBy, blocks, nil
}

// childIssues reads parentID's children of kind in fid status's order:
// sequence, then priority, then creation and ID, a missing sequence or
// priority last.
func childIssues(conn *sql.DB, registry *workflow.Registry, parentID, kind string) ([]ChildIssue, error) {
	rows, err := conn.Query(`
		SELECT id, kind, title, status, workflow_id, sequence, priority
		FROM issues
		WHERE parent_id = ? AND kind = ?
		ORDER BY sequence IS NULL, sequence, priority IS NULL, priority, created_at, id`, parentID, kind)
	if err != nil {
		return nil, fmt.Errorf("query children of %q: %w", parentID, err)
	}
	defer closeIssueRows(rows)
	children := []ChildIssue{}
	for rows.Next() {
		var (
			child              ChildIssue
			workflowID         string
			sequence, priority sql.NullInt64
		)
		if err := rows.Scan(&child.ID, &child.Kind, &child.Title, &child.Status, &workflowID, &sequence, &priority); err != nil {
			return nil, fmt.Errorf("scan child of %q: %w", parentID, err)
		}
		if child.StatusCategory, err = registryStatusCategory(registry, workflowID, child.Status); err != nil {
			return nil, err
		}
		child.Sequence, child.Priority, child.BlockedBy = nullIntPointer(sequence), nullIntPointer(priority), []string{}
		children = append(children, child)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate children of %q: %w", parentID, err)
	}
	return children, nil
}

// childProgress counts children as an epic's progress counts its tasks, and
// says whether any has started or completed.
func childProgress(children []ChildIssue) (Progress, bool) {
	progress, begun := Progress{}, false
	for _, child := range children {
		progress.count(child.StatusCategory)
		begun = begun || child.StatusCategory == "started" || child.StatusCategory == "completed"
	}
	return progress, begun
}
