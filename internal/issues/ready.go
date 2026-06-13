package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type ReadyOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	EpicID        string
}

type ReadyResult struct {
	Task        *ReadyTask
	Explanation *ReadyExplanation
}

type ReadyTask struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
	EpicID    *string `json:"epic_id"`
	EpicTitle *string `json:"epic_title"`
	Sequence  *int    `json:"sequence"`
	Priority  *int    `json:"priority"`
	EstTokens *int    `json:"est_tokens"`
}

type ReadyExplanation struct {
	Code               string
	Scope              string
	EpicID             *string
	BlockingTaskID     *string
	BlockingTaskTitle  *string
	BlockingTaskStatus *string
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

	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return ReadyResult{}, err
	}
	defer closeIssueDB(conn)

	rows, err := loadReadyRows(conn, strings.TrimSpace(options.EpicID))
	if err != nil {
		return ReadyResult{}, err
	}
	if strings.TrimSpace(options.EpicID) != "" {
		return selectReadyForEpic(strings.TrimSpace(options.EpicID), rows), nil
	}
	return selectReadyProjectWide(rows), nil
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
			ie.est_tokens,
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
			sequence, priority, estTokens     sql.NullInt64
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
			&estTokens,
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
		record.task.EstTokens = nullIntPointer(estTokens)
		record.scheduling = nullStringPointer(scheduling)
		records = append(records, record)
	}
	if err := sqlRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ready issue rows: %w", err)
	}
	return records, nil
}

func selectReadyForEpic(epicID string, rows []readyRowRecord) ReadyResult {
	scheduling := epicScheduling(rows)
	if scheduling == "sequential" {
		return selectSequentialEpic(epicID, rows)
	}
	if scheduling == "prioritised" {
		return selectPrioritisedEpic(epicID, rows)
	}
	return ReadyResult{
		Explanation: &ReadyExplanation{
			Code:   "no_ready_tasks",
			Scope:  "epic",
			EpicID: readyStringPointer(epicID),
		},
	}
}

func selectReadyProjectWide(rows []readyRowRecord) ReadyResult {
	candidates := []readyCandidate{}
	epicRows := map[string][]readyRowRecord{}

	for _, row := range rows {
		if row.task.EpicID == nil {
			if isActionableReadyRow(row) {
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
			result = selectSequentialEpic(epicID, groupedRows)
		case "prioritised":
			result = selectPrioritisedEpic(epicID, groupedRows)
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

func selectSequentialEpic(epicID string, rows []readyRowRecord) ReadyResult {
	current := lowestSequenceUnfinished(rows)
	if current == nil {
		return ReadyResult{
			Explanation: &ReadyExplanation{
				Code:   "no_ready_tasks",
				Scope:  "epic",
				EpicID: readyStringPointer(epicID),
			},
		}
	}
	if isActionableReadyRow(*current) {
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

func selectPrioritisedEpic(epicID string, rows []readyRowRecord) ReadyResult {
	candidates := []readyCandidate{}
	for _, row := range rows {
		if !isActionableReadyRow(row) {
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

func lowestSequenceUnfinished(rows []readyRowRecord) *readyRowRecord {
	var current *readyRowRecord
	for _, row := range rows {
		if !isUnfinishedReadyRow(row) {
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

func isActionableReadyRow(row readyRowRecord) bool {
	return readyStatusCategory(row.workflowID, row.status) == "unstarted"
}

func isUnfinishedReadyRow(row readyRowRecord) bool {
	category := readyStatusCategory(row.workflowID, row.status)
	return category != "completed" && category != "cancelled"
}

func readyStatusCategory(workflowID, status string) string {
	if category, ok := workflow.CategoryForState(workflowID, status); ok {
		return category
	}
	return statusCategory(status)
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
