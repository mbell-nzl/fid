package workflow

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

// ListOptions supplies the read context and filters for workflow listing.
type ListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Kind          string
	Active        bool
}

// ListItem is one workflow with projected state, transition, and assignment summaries.
type ListItem struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	AppliesTo   *string      `json:"applies_to"`
	Active      bool         `json:"active"`
	States      int          `json:"states"`
	Transitions int          `json:"transitions"`
	AssignedTo  []Assignment `json:"assigned_to"`
}

// Assignment identifies one scope to which a workflow is assigned.
type Assignment struct {
	Scope string `json:"scope"`
	Kind  string `json:"kind,omitempty"`
}

// List reads workflows from the fresh projection in deterministic order.
func List(options ListOptions) ([]ListItem, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return nil, fmt.Errorf("list workflows: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return nil, fmt.Errorf("list workflows: migrations dir is required")
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}

	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer closeWorkflowDB(conn)

	items, err := listWorkflowSummaries(conn, options)
	if err != nil {
		return nil, err
	}
	for index := range items {
		assignments, err := listAssignments(conn, items[index].ID)
		if err != nil {
			return nil, err
		}
		items[index].AssignedTo = assignments
	}
	return items, nil
}

func listWorkflowSummaries(conn *sql.DB, options ListOptions) ([]ListItem, error) {
	rows, err := conn.Query(workflowListQuery(options), workflowListArgs(options)...)
	if err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]ListItem, 0)
	for rows.Next() {
		item, err := scanWorkflowListItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workflow list rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close workflow list rows: %w", err)
	}
	return items, nil
}

func workflowListQuery(options ListOptions) string {
	where := ""
	if strings.TrimSpace(options.Kind) != "" {
		where = "WHERE (w.applies_to = ? OR w.applies_to IS NULL)"
	}
	if options.Active {
		if where == "" {
			where = "WHERE w.active = 1"
		} else {
			where += " AND w.active = 1"
		}
	}
	return `
		SELECT w.id, w.name, w.applies_to, w.active,
			(SELECT COUNT(*) FROM workflow_states s WHERE s.workflow_id = w.id),
			(SELECT COUNT(*) FROM workflow_transitions t WHERE t.workflow_id = w.id)
		FROM workflows w
		` + where + `
		ORDER BY w.name, w.id`
}

func workflowListArgs(options ListOptions) []any {
	if strings.TrimSpace(options.Kind) == "" {
		return nil
	}
	return []any{strings.TrimSpace(options.Kind)}
}

func scanWorkflowListItem(rows *sql.Rows) (ListItem, error) {
	var (
		item      ListItem
		appliesTo sql.NullString
	)
	if err := rows.Scan(&item.ID, &item.Name, &appliesTo, &item.Active, &item.States, &item.Transitions); err != nil {
		return ListItem{}, fmt.Errorf("scan workflow list row: %w", err)
	}
	if appliesTo.Valid {
		item.AppliesTo = &appliesTo.String
	}
	return item, nil
}

func listAssignments(conn *sql.DB, workflowID string) ([]Assignment, error) {
	rows, err := conn.Query(`
		SELECT scope, kind
		FROM workflow_assignments
		WHERE workflow_id = ?
		ORDER BY scope = 'project', kind
	`, workflowID)
	if err != nil {
		return nil, fmt.Errorf("list workflow %q assignments: %w", workflowID, err)
	}
	defer func() { _ = rows.Close() }()

	assignments := make([]Assignment, 0)
	for rows.Next() {
		var (
			assignment Assignment
			kind       sql.NullString
		)
		if err := rows.Scan(&assignment.Scope, &kind); err != nil {
			return nil, fmt.Errorf("scan workflow %q assignment: %w", workflowID, err)
		}
		if kind.Valid {
			assignment.Kind = kind.String
		}
		assignments = append(assignments, assignment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workflow %q assignments: %w", workflowID, err)
	}
	return assignments, nil
}
