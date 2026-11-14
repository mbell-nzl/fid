package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issuestate"
	"github.com/mbell-nzl/fid/internal/journal"
)

type ReadyOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	EpicID        string
}

type ReadyResult struct {
	Task        *ReadyTask        `json:"task"`
	Explanation *ReadyExplanation `json:"explanation"`
}

type ReadyTask struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
	EpicID    *string `json:"epic_id"`
	EpicTitle *string `json:"epic_title"`
	Sequence  *int    `json:"sequence"`
	Priority  *int    `json:"priority"`
}

type ReadyExplanation struct {
	Code               string  `json:"code"`
	Scope              string  `json:"scope"`
	EpicID             *string `json:"epic_id"`
	BlockingTaskID     *string `json:"blocking_task_id"`
	BlockingTaskTitle  *string `json:"blocking_task_title"`
	BlockingTaskStatus *string `json:"blocking_task_status"`
}

type readyCandidate struct {
	task      ReadyTask
	createdAt string
}

type readyRowRecord struct {
	task       ReadyTask
	workflowID string
	status     string
	createdAt  string
	scheduling *string
}

func Ready(options ReadyOptions) (ReadyResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return ReadyResult{}, fmt.Errorf("ready issues: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return ReadyResult{}, fmt.Errorf("ready issues: migrations dir is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return ReadyResult{}, err
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return ReadyResult{}, err
	}
	defer closeIssueDB(conn)
	registry, err := issuestate.Load(conn)
	if err != nil {
		return ReadyResult{}, fmt.Errorf("ready issues: load workflow registry: %w", err)
	}

	rows, err := loadReadyRows(conn, strings.TrimSpace(options.EpicID))
	if err != nil {
		return ReadyResult{}, err
	}
	if strings.TrimSpace(options.EpicID) != "" {
		return selectReadyForEpic(registry, strings.TrimSpace(options.EpicID), rows), nil
	}
	return selectReadyProjectWide(registry, rows), nil
}

func loadReadyRows(conn *sql.DB, epicID string) ([]readyRowRecord, error) {
	args := []any{}
	query := `
		SELECT
			ie.id,
			ie.title,
			ie.status,
			ie.workflow_id,
			ie.parent_id,
			parent.title,
			ie.sequence,
			ie.priority,
			ie.created_at,
			em.scheduling
		FROM issues_enriched ie
		LEFT JOIN issues parent ON parent.id = ie.parent_id AND parent.kind = 'epic'
		LEFT JOIN epic_metadata em ON em.issue_id = ie.parent_id
		WHERE ie.kind = 'task'
	`
	if epicID != "" {
		query += "\n  AND ie.parent_id = ?"
		args = append(args, epicID)
	}
	query += "\nORDER BY ie.created_at ASC, ie.id ASC"

	sqlRows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query ready issues: %w", err)
	}
	defer closeIssueRows(sqlRows)

	records := []readyRowRecord{}
	for sqlRows.Next() {
		var (
			record                            readyRowRecord
			parentID, parentTitle, scheduling sql.NullString
			sequence, priority                sql.NullInt64
		)
		if err := sqlRows.Scan(
			&record.task.ID,
			&record.task.Title,
			&record.status,
			&record.workflowID,
			&parentID,
			&parentTitle,
			&sequence,
			&priority,
			&record.createdAt,
			&scheduling,
		); err != nil {
			return nil, fmt.Errorf("scan ready issue row: %w", err)
		}
		record.task.Status = record.status
		record.task.EpicID = nullStringPointer(parentID)
		record.task.EpicTitle = nullStringPointer(parentTitle)
		record.task.Sequence = nullIntPointer(sequence)
		record.task.Priority = nullIntPointer(priority)
		record.scheduling = nullStringPointer(scheduling)
		records = append(records, record)
	}
	if err := sqlRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ready issue rows: %w", err)
	}
	return records, nil
}

func selectReadyForEpic(registry *issuestate.Registry, epicID string, rows []readyRowRecord) ReadyResult {
	scheduling := epicScheduling(rows)
	if scheduling == "sequential" {
		return selectSequentialEpic(registry, epicID, rows)
	}
	if scheduling == "prioritised" {
		return selectPrioritisedEpic(registry, epicID, rows)
	}
	return ReadyResult{
		Explanation: &ReadyExplanation{
			Code:   "no_ready_tasks",
			Scope:  "epic",
			EpicID: readyStringPointer(epicID),
		},
	}
}

func selectReadyProjectWide(registry *issuestate.Registry, rows []readyRowRecord) ReadyResult {
	candidates := []readyCandidate{}
	epicRows := map[string][]readyRowRecord{}

	for _, row := range rows {
		if row.task.EpicID == nil {
			if isActionableReadyRow(registry, row) {
				candidates = append(candidates, readyCandidate{task: row.task, createdAt: row.createdAt})
			}
			continue
		}
		epicRows[*row.task.EpicID] = append(epicRows[*row.task.EpicID], row)
	}

	for epicID, groupedRows := range epicRows {
		var result ReadyResult
		switch epicScheduling(groupedRows) {
		case "sequential":
			result = selectSequentialEpic(registry, epicID, groupedRows)
		case "prioritised":
			result = selectPrioritisedEpic(registry, epicID, groupedRows)
		default:
			continue
		}
		if result.Task != nil {
			candidates = append(candidates, readyCandidate{
				task:      *result.Task,
				createdAt: readyCreatedAt(groupedRows, result.Task.ID),
			})
		}
	}

	if selected := selectTopReadyCandidate(candidates); selected != nil {
		return ReadyResult{Task: &selected.task}
	}
	return ReadyResult{
		Explanation: &ReadyExplanation{
			Code:  "no_ready_tasks",
			Scope: "project",
		},
	}
}

func selectSequentialEpic(registry *issuestate.Registry, epicID string, rows []readyRowRecord) ReadyResult {
	current := lowestSequenceUnfinished(registry, rows)
	if current == nil {
		return ReadyResult{
			Explanation: &ReadyExplanation{
				Code:   "no_ready_tasks",
				Scope:  "epic",
				EpicID: readyStringPointer(epicID),
			},
		}
	}
	if isActionableReadyRow(registry, *current) {
		return ReadyResult{Task: &current.task}
	}
	return ReadyResult{
		Explanation: &ReadyExplanation{
			Code:               "sequential_current_not_actionable",
			Scope:              "epic",
			EpicID:             readyStringPointer(epicID),
			BlockingTaskID:     readyStringPointer(current.task.ID),
			BlockingTaskTitle:  readyStringPointer(current.task.Title),
			BlockingTaskStatus: readyStringPointer(current.task.Status),
		},
	}
}

func selectPrioritisedEpic(registry *issuestate.Registry, epicID string, rows []readyRowRecord) ReadyResult {
	candidates := []readyCandidate{}
	for _, row := range rows {
		if !isActionableReadyRow(registry, row) {
			continue
		}
		candidates = append(candidates, readyCandidate{
			task:      row.task,
			createdAt: row.createdAt,
		})
	}
	if selected := selectTopReadyCandidate(candidates); selected != nil {
		return ReadyResult{Task: &selected.task}
	}
	return ReadyResult{
		Explanation: &ReadyExplanation{
			Code:   "no_ready_tasks",
			Scope:  "epic",
			EpicID: readyStringPointer(epicID),
		},
	}
}

func epicScheduling(rows []readyRowRecord) string {
	for _, row := range rows {
		if row.scheduling != nil {
			return *row.scheduling
		}
	}
	return ""
}

func lowestSequenceUnfinished(registry *issuestate.Registry, rows []readyRowRecord) *readyRowRecord {
	var current *readyRowRecord
	for _, row := range rows {
		if !isUnfinishedReadyRow(registry, row) {
			continue
		}
		if current == nil || compareSequenceRows(row, *current) < 0 {
			candidate := row
			current = &candidate
		}
	}
	return current
}

func compareSequenceRows(left, right readyRowRecord) int {
	leftSequence := optionalIntSortValue(left.task.Sequence)
	rightSequence := optionalIntSortValue(right.task.Sequence)
	switch {
	case leftSequence < rightSequence:
		return -1
	case leftSequence > rightSequence:
		return 1
	case left.createdAt < right.createdAt:
		return -1
	case left.createdAt > right.createdAt:
		return 1
	case left.task.ID < right.task.ID:
		return -1
	case left.task.ID > right.task.ID:
		return 1
	default:
		return 0
	}
}

func isActionableReadyRow(registry *issuestate.Registry, row readyRowRecord) bool {
	return readyStatusCategory(registry, row.workflowID, row.status) == "unstarted"
}

func isUnfinishedReadyRow(registry *issuestate.Registry, row readyRowRecord) bool {
	category := readyStatusCategory(registry, row.workflowID, row.status)
	return category != "completed" && category != "cancelled"
}

func readyStatusCategory(registry *issuestate.Registry, workflowID, status string) string {
	if category, ok := registry.CategoryForState(workflowID, status); ok {
		return category
	}
	return ""
}

func selectTopReadyCandidate(candidates []readyCandidate) *readyCandidate {
	var selected *readyCandidate
	for _, candidate := range candidates {
		if selected == nil || compareReadyCandidates(candidate, *selected) < 0 {
			chosen := candidate
			selected = &chosen
		}
	}
	return selected
}

func compareReadyCandidates(left, right readyCandidate) int {
	leftPriority := optionalIntSortValue(left.task.Priority)
	rightPriority := optionalIntSortValue(right.task.Priority)
	switch {
	case leftPriority < rightPriority:
		return -1
	case leftPriority > rightPriority:
		return 1
	case left.createdAt < right.createdAt:
		return -1
	case left.createdAt > right.createdAt:
		return 1
	case left.task.ID < right.task.ID:
		return -1
	case left.task.ID > right.task.ID:
		return 1
	default:
		return 0
	}
}

func readyCreatedAt(rows []readyRowRecord, issueID string) string {
	for _, row := range rows {
		if row.task.ID == issueID {
			return row.createdAt
		}
	}
	return ""
}

func optionalIntSortValue(value *int) int {
	if value == nil {
		return int(^uint(0) >> 1)
	}
	return *value
}

func readyStringPointer(value string) *string {
	cloned := value
	return &cloned
}
