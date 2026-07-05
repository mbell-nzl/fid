package issues

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func Get(options GetOptions) (Issue, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return Issue{}, fmt.Errorf("get issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return Issue{}, fmt.Errorf("get issue: migrations dir is required")
	}
	if strings.TrimSpace(options.ID) == "" {
		return Issue{}, fmt.Errorf("get issue: issue id is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return Issue{}, err
	}

	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return Issue{}, err
	}
	defer closeIssueDB(conn)

	row := conn.QueryRow(`
		SELECT id, kind, title, status, workflow_id, parent_id, priority, sequence,
		       spec_md_path, est_tokens, due_date, owner_role, owner_actor, in_inbox,
		       concerns_raised, open_concerns_count, triage_status, created_at, updated_at
		FROM issues_enriched
		WHERE id = ?
	`, options.ID)

	var (
		record                                         Issue
		parentID, specMDPath, dueDate, ownerRole       sql.NullString
		ownerActor, triageStatus, createdAt, updatedAt sql.NullString
		priority, sequence, estTokens                  sql.NullInt64
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
		&estTokens,
		&dueDate,
		&ownerRole,
		&ownerActor,
		&record.InInbox,
		&record.ConcernsRaised,
		&record.OpenConcernsCount,
		&triageStatus,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Issue{}, &NotFoundError{IssueID: options.ID}
		}
		return Issue{}, fmt.Errorf("query issue %q: %w", options.ID, err)
	}

	record.StatusCategory = statusCategory(record.Status)
	record.Labels, err = issueLabels(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	record.AcceptanceCriteria = []AcceptanceCriterion{}
	record.ParentID = nullStringPointer(parentID)
	record.SpecMDPath = nullStringPointer(specMDPath)
	record.DueDate = normalizeDateOnlyPointer(nullStringPointer(dueDate))
	record.OwnerRole = nullStringPointer(ownerRole)
	record.OwnerActor = nullStringPointer(ownerActor)
	record.TriageStatus = nullStringPointer(triageStatus)
	record.CreatedAt = nullStringPointer(createdAt)
	record.UpdatedAt = nullStringPointer(updatedAt)
	record.Priority = nullIntPointer(priority)
	record.Sequence = nullIntPointer(sequence)
	record.EstTokens = nullIntPointer(estTokens)
	record.ParentTitle, err = issueParentTitle(conn, record.ParentID)
	if err != nil {
		return Issue{}, err
	}
	record.AcceptanceCriteria, err = issueAcceptanceCriteria(conn, record.ID, record.OwnerRole)
	if err != nil {
		return Issue{}, err
	}
	record.Concerns, err = issueConcerns(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	record.Completions, err = issueCompletions(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	record.Links, err = issueLinks(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	record.Claims, err = issueClaims(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	record.Gates, err = issueGates(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}

	return record, nil
}
