package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type ListOptions struct {
	WorkspaceRoot   string
	MigrationsDir   string
	Kind            string
	Status          string
	StatusCategory  string
	EpicID          string
	Label           string
	Priorities      []int
	InInbox         *bool
	HasOpenConcerns *bool
	Order           string
}

type ListItem struct {
	ID                string  `json:"id"`
	Kind              string  `json:"kind"`
	Title             string  `json:"title"`
	Status            string  `json:"status"`
	StatusCategory    string  `json:"status_category"`
	ParentID          *string `json:"parent_id"`
	Priority          *int    `json:"priority"`
	ConcernsRaised    bool    `json:"concerns_raised"`
	OpenConcernsCount int     `json:"open_concerns_count"`
	InInbox           bool    `json:"-"`
}

type UnsupportedListOrderError struct {
	Order          string
	SupportedOrder []string
}

func (err *UnsupportedListOrderError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("unsupported issue list order %q; supported orders: %s", err.Order, strings.Join(err.SupportedOrder, ", "))
}

var supportedListOrders = []string{"sequence", "priority", "created", "updated"}

func List(options ListOptions) ([]ListItem, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return nil, fmt.Errorf("list issues: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return nil, fmt.Errorf("list issues: migrations dir is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}

	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer closeIssueDB(conn)

	query, args, err := issueListQuery(options)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list issues: %w", err)
	}
	defer closeIssueRows(rows)

	items := []ListItem{}
	for rows.Next() {
		var (
			item     ListItem
			parentID sql.NullString
			priority sql.NullInt64
		)
		if err := rows.Scan(
			&item.ID,
			&item.Kind,
			&item.Title,
			&item.Status,
			&parentID,
			&priority,
			&item.ConcernsRaised,
			&item.OpenConcernsCount,
			&item.InInbox,
		); err != nil {
			return nil, fmt.Errorf("scan issue list row: %w", err)
		}

		item.StatusCategory = statusCategory(item.Status)
		item.ParentID = nullStringPointer(parentID)
		item.Priority = nullIntPointer(priority)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate issue list rows: %w", err)
	}

	return items, nil
}

func SupportedListOrders() []string {
	return append([]string(nil), supportedListOrders...)
}

func issueListQuery(options ListOptions) (string, []any, error) {
	whereParts := []string{}
	args := []any{}

	if kind := strings.TrimSpace(options.Kind); kind != "" {
		whereParts = append(whereParts, "ie.kind = ?")
		args = append(args, kind)
	}
	if status := strings.TrimSpace(options.Status); status != "" {
		whereParts = append(whereParts, "ie.status = ?")
		args = append(args, status)
	}
	if statusCategory := strings.TrimSpace(options.StatusCategory); statusCategory != "" {
		whereParts = append(whereParts, issueStatusCategorySQL()+" = ?")
		args = append(args, statusCategory)
	}
	if epicID := strings.TrimSpace(options.EpicID); epicID != "" {
		whereParts = append(whereParts, "ie.parent_id = ?")
		args = append(args, epicID)
	}
	if label := strings.TrimSpace(options.Label); label != "" {
		whereParts = append(whereParts, `EXISTS (
			SELECT 1
			FROM issue_labels il
			WHERE il.issue_id = ie.id AND il.label = ?
		)`)
		args = append(args, label)
	}
	if len(options.Priorities) > 0 {
		priorityPlaceholders := make([]string, 0, len(options.Priorities))
		for _, priority := range options.Priorities {
			priorityPlaceholders = append(priorityPlaceholders, "?")
			args = append(args, priority)
		}
		whereParts = append(whereParts, "ie.priority IN ("+strings.Join(priorityPlaceholders, ", ")+")")
	}
	if options.InInbox != nil {
		whereParts = append(whereParts, "ie.in_inbox = ?")
		args = append(args, *options.InInbox)
	}
	if options.HasOpenConcerns != nil {
		if *options.HasOpenConcerns {
			whereParts = append(whereParts, "ie.open_concerns_count > 0")
		} else {
			whereParts = append(whereParts, "ie.open_concerns_count = 0")
		}
	}

	query := strings.Builder{}
	query.WriteString(`
		SELECT
			ie.id,
			ie.kind,
			ie.title,
			ie.status,
			ie.parent_id,
			ie.priority,
			ie.concerns_raised,
			ie.open_concerns_count,
			ie.in_inbox
		FROM issues_enriched ie
	`)
	if len(whereParts) > 0 {
		query.WriteString("\nWHERE ")
		query.WriteString(strings.Join(whereParts, "\n  AND "))
	}

	orderClause, err := issueListOrderClause(options.Order)
	if err != nil {
		return "", nil, err
	}
	query.WriteString("\nORDER BY ")
	query.WriteString(orderClause)

	return query.String(), args, nil
}

func issueListOrderClause(order string) (string, error) {
	switch strings.TrimSpace(order) {
	case "", "created":
		return "ie.created_at ASC, ie.id ASC", nil
	case "updated":
		return "ie.updated_at ASC, ie.id ASC", nil
	case "priority":
		return "ie.priority IS NULL ASC, ie.priority ASC, ie.created_at ASC, ie.id ASC", nil
	case "sequence":
		return "ie.sequence IS NULL ASC, ie.sequence ASC, ie.created_at ASC, ie.id ASC", nil
	default:
		return "", &UnsupportedListOrderError{
			Order:          strings.TrimSpace(order),
			SupportedOrder: SupportedListOrders(),
		}
	}
}

func issueStatusCategorySQL() string {
	return `CASE
		WHEN ie.status IN ('todo') THEN 'unstarted'
		WHEN ie.status IN ('in_progress', 'dev_complete', 'qa_review', 'qa_block', 'failed_gate', 'wrapped') THEN 'started'
		WHEN ie.status IN ('qa_pass', 'done', 'closed') THEN 'completed'
		WHEN ie.status IN ('cancelled') THEN 'cancelled'
		ELSE 'unstarted'
	END`
}

func ParsePriorityFilter(value string) ([]int, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}

	parts := strings.Split(trimmed, ",")
	priorities := make([]int, 0, len(parts))
	for _, part := range parts {
		priorityValue := strings.TrimSpace(part)
		if priorityValue == "" {
			return nil, fmt.Errorf("priority filter contains an empty value")
		}
		priority, err := strconv.Atoi(priorityValue)
		if err != nil {
			return nil, fmt.Errorf("priority filter %q is not an integer", priorityValue)
		}
		priorities = append(priorities, priority)
	}

	return priorities, nil
}
