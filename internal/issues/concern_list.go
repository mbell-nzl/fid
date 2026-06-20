package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type ConcernListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	Status        string
}

type ConcernListItem struct {
	ID                   string  `json:"id"`
	IssueID              string  `json:"issue_id"`
	Summary              string  `json:"summary"`
	Status               string  `json:"status"`
	RaisedByRole         string  `json:"raised_by_role"`
	RaisedInCompletionID string  `json:"raised_in_completion_id"`
	BodyMDPath           *string `json:"body_md_path"`
	CreatedAt            *string `json:"created_at"`
	UpdatedAt            *string `json:"updated_at"`
}

func ListConcerns(options ConcernListOptions) ([]ConcernListItem, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return nil, fmt.Errorf("list concerns: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return nil, fmt.Errorf("list concerns: migrations dir is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}

	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer closeIssueDB(conn)

	query, args := concernListQuery(options)
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list concerns: %w", err)
	}
	defer closeIssueRows(rows)

	items := []ConcernListItem{}
	for rows.Next() {
		var item ConcernListItem
		var bodyMDPath, createdAt, updatedAt sql.NullString
		if err := rows.Scan(
			&item.ID,
			&item.IssueID,
			&item.Summary,
			&item.Status,
			&item.RaisedByRole,
			&item.RaisedInCompletionID,
			&bodyMDPath,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan concern list row: %w", err)
		}
		item.BodyMDPath = nullStringPointer(bodyMDPath)
		item.CreatedAt = nullStringPointer(createdAt)
		item.UpdatedAt = nullStringPointer(updatedAt)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate concern list rows: %w", err)
	}

	return items, nil
}

func concernListQuery(options ConcernListOptions) (string, []any) {
	whereParts := []string{}
	args := []any{}

	if issueID := strings.TrimSpace(options.IssueID); issueID != "" {
		whereParts = append(whereParts, "issue_id = ?")
		args = append(args, issueID)
	}
	if status := strings.TrimSpace(options.Status); status != "" {
		whereParts = append(whereParts, "status = ?")
		args = append(args, status)
	}

	query := strings.Builder{}
	query.WriteString(`
		SELECT
			id,
			issue_id,
			summary,
			status,
			raised_by_role,
			raised_in_completion_id,
			body_md_path,
			created_at,
			updated_at
		FROM concerns
	`)
	if len(whereParts) > 0 {
		query.WriteString("\nWHERE ")
		query.WriteString(strings.Join(whereParts, "\n  AND "))
	}
	query.WriteString("\nORDER BY created_at DESC, id DESC")

	return query.String(), args
}
