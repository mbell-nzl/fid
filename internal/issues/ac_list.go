package issues

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	_ "modernc.org/sqlite"
)

const issueReadOnlyDriverName = "sqlite"

type ListAcceptanceCriteriaOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
}

type AcceptanceCriteriaList struct {
	IssueID  string                `json:"-"`
	Criteria []AcceptanceCriterion `json:"-"`
}

func (list AcceptanceCriteriaList) MarshalJSON() ([]byte, error) {
	criteria := list.Criteria
	if criteria == nil {
		criteria = []AcceptanceCriterion{}
	}
	return json.Marshal(criteria)
}

func ListAcceptanceCriteria(options ListAcceptanceCriteriaOptions) (AcceptanceCriteriaList, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return AcceptanceCriteriaList{}, fmt.Errorf("list acceptance criteria: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return AcceptanceCriteriaList{}, fmt.Errorf("list acceptance criteria: migrations dir is required")
	}
	if strings.TrimSpace(options.IssueID) == "" {
		return AcceptanceCriteriaList{}, fmt.Errorf("list acceptance criteria: issue id is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return AcceptanceCriteriaList{}, err
	}

	conn, err := openIssueReadOnlyDB(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return AcceptanceCriteriaList{}, err
	}
	defer closeIssueDB(conn)

	ownerRole, err := listAcceptanceCriteriaOwnerRole(conn, options.IssueID)
	if err != nil {
		return AcceptanceCriteriaList{}, err
	}

	criteria, err := issueAcceptanceCriteria(conn, options.IssueID, ownerRole)
	if err != nil {
		return AcceptanceCriteriaList{}, err
	}

	return AcceptanceCriteriaList{
		IssueID:  options.IssueID,
		Criteria: criteria,
	}, nil
}

func listAcceptanceCriteriaOwnerRole(conn *sql.DB, issueID string) (*string, error) {
	var ownerRole sql.NullString
	if err := conn.QueryRow(`SELECT owner_role FROM issues WHERE id = ?`, issueID).Scan(&ownerRole); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &NotFoundError{IssueID: issueID}
		}
		return nil, fmt.Errorf("query issue %q owner role: %w", issueID, err)
	}
	return nullStringPointer(ownerRole), nil
}

func openIssueReadOnlyDB(path string) (*sql.DB, error) {
	dsn, err := issueReadOnlyDSN(path)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open(issueReadOnlyDriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	return conn, nil
}

func issueReadOnlyDSN(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve sqlite database path: %w", err)
	}

	dsnPath := filepath.ToSlash(absolutePath)
	if !strings.HasPrefix(dsnPath, "/") {
		dsnPath = "/" + dsnPath
	}

	query := url.Values{}
	query.Set("mode", "ro")
	query.Add("_pragma", "foreign_keys(1)")
	query.Add("_pragma", "recursive_triggers(1)")

	return (&url.URL{
		Scheme:   "file",
		Path:     dsnPath,
		RawQuery: query.Encode(),
	}).String(), nil
}
