package labels

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func List(options ListOptions) ([]Label, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return nil, fmt.Errorf("label list: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return nil, fmt.Errorf("label list: migrations dir is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer closeLabelsConn(conn)

	rows, err := conn.Query(`
		SELECT name, colour, description
		FROM labels
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("label list: %w", err)
	}
	defer closeLabelRows(rows)

	labels := []Label{}
	for rows.Next() {
		var (
			label       Label
			colour      sql.NullString
			description sql.NullString
		)
		if err := rows.Scan(&label.Name, &colour, &description); err != nil {
			return nil, fmt.Errorf("scan label list row: %w", err)
		}
		label.Colour = nullStringPointer(colour)
		label.Description = nullStringPointer(description)
		labels = append(labels, label)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate label list rows: %w", err)
	}
	if len(labels) > ReadCollectionLimit {
		return nil, &ResultLimitExceededError{Collection: "labels", Limit: ReadCollectionLimit, Observed: len(labels)}
	}

	return labels, nil
}

func closeLabelRows(rows *sql.Rows) {
	_ = rows.Close()
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	copy := value.String
	return &copy
}
