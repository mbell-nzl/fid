package issues

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issuestate"
	"github.com/mbell-nzl/fid/internal/journal"
)

func Get(options GetOptions) (Issue, error) {
	conn, registry, err := openIssueRead(options, "get issue")
	if err != nil {
		return Issue{}, err
	}
	defer closeIssueDB(conn)
	return readIssue(conn, registry, options.ID)
}

// openIssueRead checks a read of one issue, refreshes the projection as
// every read does, and opens it with its workflow registry. The caller
// closes conn.
func openIssueRead(options GetOptions, operation string) (*sql.DB, *issuestate.Registry, error) {
	switch {
	case strings.TrimSpace(options.WorkspaceRoot) == "":
		return nil, nil, fmt.Errorf("%s: workspace root is required", operation)
	case strings.TrimSpace(options.MigrationsDir) == "":
		return nil, nil, fmt.Errorf("%s: migrations dir is required", operation)
	case strings.TrimSpace(options.ID) == "":
		return nil, nil, fmt.Errorf("%s: issue id is required", operation)
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, nil, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, nil, err
	}
	registry, err := issuestate.Load(conn)
	if err != nil {
		closeIssueDB(conn)
		return nil, nil, fmt.Errorf("%s: load workflow registry: %w", operation, err)
	}
	return conn, registry, nil
}

// readIssue reads the issue whose ID is exactly id.
func readIssue(conn *sql.DB, registry *issuestate.Registry, id string) (Issue, error) {
	row := conn.QueryRow(`
		SELECT id, kind, title, status, workflow_id, parent_id, priority, sequence,
		       spec_md_path, revision, owner_role, owner_actor,
		       concerns_raised, open_concerns_count, triage_status, created_at, updated_at, lane
		FROM issues_enriched
		WHERE id = ?
	`, id)

	var (
		record                                         Issue
		parentID, specMDPath, ownerRole                sql.NullString
		ownerActor, triageStatus, createdAt, updatedAt sql.NullString
		priority, sequence                             sql.NullInt64
		lane                                           sql.NullString
	)
	if err := row.Scan(
		&record.ID,
		&record.Kind,
		&record.Title,
		&record.Status,
		&record.WorkflowID,
		&parentID,
		&priority,
		&sequence,
		&specMDPath,
		&record.Revision,
		&ownerRole,
		&ownerActor,
		&record.ConcernsRaised,
		&record.OpenConcernsCount,
		&triageStatus,
		&createdAt,
		&updatedAt, &lane,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Issue{}, &NotFoundError{IssueID: id}
		}
		return Issue{}, fmt.Errorf("query issue %q: %w", id, err)
	}

	var err error
	record.StatusCategory, err = registryStatusCategory(registry, record.WorkflowID, record.Status)
	if err != nil {
		return Issue{}, err
	}
	record.Labels, err = issueLabels(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	record.AcceptanceCriteria = []AcceptanceCriterion{}
	record.Lane = nullStringPointer(lane)
	record.ParentID = nullStringPointer(parentID)
	record.SpecMDPath = nullStringPointer(specMDPath)
	record.OwnerRole = nullStringPointer(ownerRole)
	record.OwnerActor = nullStringPointer(ownerActor)
	record.TriageStatus = nullStringPointer(triageStatus)
	record.CreatedAt = nullStringPointer(createdAt)
	record.UpdatedAt = nullStringPointer(updatedAt)
	record.Priority = nullIntPointer(priority)
	record.Sequence = nullIntPointer(sequence)
	record.ParentTitle, err = issueParentTitle(conn, record.ParentID)
	if err != nil {
		return Issue{}, err
	}
	record.AcceptanceCriteria, err = issueAcceptanceCriteria(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	record.Concerns, err = issueConcerns(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	record.Links, err = issueLinks(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	return record, nil
}
