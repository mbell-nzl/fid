package review

import (
	"fmt"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type ConcernListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Subject       string
	Status        string
	Severity      string
	Limit         int
}

func ListConcerns(options ConcernListOptions) ([]map[string]any, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}
	if options.Limit < 1 || options.Limit > 10000 || (options.Status != "" && !oneOf(options.Status, "open", "disposed")) || (options.Severity != "" && !oneOf(options.Severity, "blocker", "major", "minor", "nit")) {
		return nil, fmt.Errorf("concern list input is invalid")
	}
	subjectType, subjectID, err := splitSubjectFilter(options.Subject)
	if err != nil {
		return nil, err
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	rows, err := conn.Query(`SELECT id FROM concerns WHERE (?='' OR (subject_type=? AND subject_id=?)) AND (?='' OR status=?) AND (?='' OR severity=?) ORDER BY created_at,id LIMIT ?`, options.Subject, subjectType, subjectID, options.Status, options.Status, options.Severity, options.Severity, options.Limit+1)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := []map[string]any{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		value, err := loadConcern(conn, id)
		if err != nil {
			return nil, err
		}
		items = append(items, value)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) > options.Limit {
		return nil, &ResultLimitExceededError{Limit: options.Limit, Observed: len(items)}
	}
	return items, nil
}

// GetConcern returns the same record as the corresponding list item.
func GetConcern(workspaceRoot, migrationsDir, id string) (map[string]any, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return nil, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, fmt.Errorf("open concern projection: %w", err)
	}
	defer func() { _ = conn.Close() }()
	return loadConcern(conn, id)
}
