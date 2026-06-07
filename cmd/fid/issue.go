package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

func newIssueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues.",
	}

	cmd.AddCommand(newIssueCreateCommand())
	cmd.AddCommand(newIssueGetCommand())
	return cmd
}

func newIssueCreateCommand() *cobra.Command {
	var kind string
	var title string
	var actorRole string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an issue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Create(issues.CreateOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				Kind:          kind,
				Title:         title,
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "created issue %s (%s) %q\n", result.ID, result.Kind, result.Title)
		},
	}

	cmd.Flags().StringVar(&kind, "kind", "", "Issue kind")
	cmd.Flags().StringVar(&title, "title", "", "Issue title")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newIssueGetCommand() *cobra.Command {
	var formatName string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an issue by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			issue, err := issues.Get(issues.GetOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ID:            args[0],
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			outputValue := any(issue)
			if fidformat.Normalize(formatName) == fidformat.Markdown {
				view, err := loadIssueMarkdownView(workspaceRoot, issue)
				if err != nil {
					if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
						return renderErr
					}
					return err
				}
				outputValue = view
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, outputValue); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}

func issueMigrationsDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "migrations"
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func renderIssueError(writer io.Writer, err error) error {
	var missingRoleErr *policy.MissingRoleError
	if errors.As(err, &missingRoleErr) {
		return writeIssueLine(writer, "error: %s\n", missingRoleErr.Error())
	}

	var unknownRoleErr *policy.UnknownRoleError
	if errors.As(err, &unknownRoleErr) {
		if err := writeIssueLine(writer, "error: %s\n", unknownRoleErr.Error()); err != nil {
			return err
		}
		return writeIssueLine(writer, "  known roles: %s\n", joinCommaSeparated(unknownRoleErr.KnownRoles))
	}

	handled, renderErr := renderJournalReplayError(func(format string, args ...any) error {
		return writeIssueLine(writer, format, args...)
	}, err, "read stopped - manual repair required")
	if handled {
		return renderErr
	}

	var unsupportedFormatErr *fidformat.UnsupportedFormatError
	if errors.As(err, &unsupportedFormatErr) {
		return writeIssueLine(writer, "error: %s\n", unsupportedFormatErr.Error())
	}

	var unsupportedKindErr *issues.UnsupportedKindError
	if errors.As(err, &unsupportedKindErr) {
		return writeIssueLine(writer, "error: %s\n", unsupportedKindErr.Error())
	}

	var notFoundErr *issues.NotFoundError
	if errors.As(err, &notFoundErr) {
		return writeIssueLine(writer, "error: %s\n", notFoundErr.Error())
	}

	return writeIssueLine(writer, "error: %v\n", err)
}

func writeIssueLine(writer io.Writer, format string, args ...any) error {
	if _, err := fmt.Fprintf(writer, format, args...); err != nil {
		return fmt.Errorf("write issue output: %w", err)
	}
	return nil
}

func joinCommaSeparated(values []string) string {
	if len(values) == 0 {
		return ""
	}
	result := values[0]
	for _, value := range values[1:] {
		result += ", " + value
	}
	return result
}

func loadIssueMarkdownView(workspaceRoot string, issue issues.Issue) (fidformat.IssueMarkdownView, error) {
	conn, err := fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return fidformat.IssueMarkdownView{}, err
	}
	defer closeIssueMarkdownDB(conn)

	view := fidformat.IssueMarkdownView{Issue: issue}
	view.ParentTitle, err = issueParentTitle(conn, issue.ParentID)
	if err != nil {
		return fidformat.IssueMarkdownView{}, err
	}

	view.Issue.Labels, err = issueLabels(conn, issue.ID)
	if err != nil {
		return fidformat.IssueMarkdownView{}, err
	}

	view.Issue.AcceptanceCriteria, err = issueAcceptanceCriteria(conn, issue)
	if err != nil {
		return fidformat.IssueMarkdownView{}, err
	}

	view.Links, err = issueMarkdownLinks(conn, issue.ID)
	if err != nil {
		return fidformat.IssueMarkdownView{}, err
	}

	return view, nil
}

func issueParentTitle(conn *sql.DB, parentID *string) (*string, error) {
	if parentID == nil {
		return nil, nil
	}

	var title string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, *parentID).Scan(&title); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query parent issue %q: %w", *parentID, err)
	}
	return &title, nil
}

func issueLabels(conn *sql.DB, issueID string) ([]string, error) {
	rows, err := conn.Query(`SELECT label FROM issue_labels WHERE issue_id = ? ORDER BY label`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query issue labels for %q: %w", issueID, err)
	}
	defer closeIssueMarkdownRows(rows)

	labels := []string{}
	for rows.Next() {
		var label string
		if err := rows.Scan(&label); err != nil {
			return nil, fmt.Errorf("scan issue label for %q: %w", issueID, err)
		}
		labels = append(labels, label)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate issue labels for %q: %w", issueID, err)
	}
	return labels, nil
}

func issueAcceptanceCriteria(conn *sql.DB, issue issues.Issue) ([]issues.AcceptanceCriterion, error) {
	rows, err := conn.Query(`
		SELECT id, ordinal, text
		FROM acceptance_criteria
		WHERE issue_id = ?
		ORDER BY ordinal, id
	`, issue.ID)
	if err != nil {
		return nil, fmt.Errorf("query acceptance criteria for %q: %w", issue.ID, err)
	}
	defer closeIssueMarkdownRows(rows)

	criteria := []issues.AcceptanceCriterion{}
	for rows.Next() {
		var criterion issues.AcceptanceCriterion
		if err := rows.Scan(&criterion.ID, &criterion.Ordinal, &criterion.Text); err != nil {
			return nil, fmt.Errorf("scan acceptance criterion for %q: %w", issue.ID, err)
		}

		criterion.LatestVerifications, err = latestVerificationsByRole(conn, criterion.ID)
		if err != nil {
			return nil, err
		}
		if len(criterion.LatestVerifications) == 0 && issue.OwnerRole != nil {
			criterion.LatestVerifications[*issue.OwnerRole] = nil
		}
		criteria = append(criteria, criterion)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate acceptance criteria for %q: %w", issue.ID, err)
	}
	return criteria, nil
}

func latestVerificationsByRole(conn *sql.DB, criterionID string) (map[string]any, error) {
	rows, err := conn.Query(`
		SELECT verifier_role, status, verified_at
		FROM ac_verifications av
		WHERE av.ac_id = ?
		  AND av.id = (
		    SELECT inner_av.id
		    FROM ac_verifications inner_av
		    WHERE inner_av.ac_id = av.ac_id AND inner_av.verifier_role = av.verifier_role
		    ORDER BY inner_av.verified_at DESC, inner_av.id DESC
		    LIMIT 1
		  )
		ORDER BY verifier_role
	`, criterionID)
	if err != nil {
		return nil, fmt.Errorf("query acceptance criterion verifications for %q: %w", criterionID, err)
	}
	defer closeIssueMarkdownRows(rows)

	latest := map[string]any{}
	for rows.Next() {
		var role string
		var status string
		var verifiedAt string
		if err := rows.Scan(&role, &status, &verifiedAt); err != nil {
			return nil, fmt.Errorf("scan acceptance criterion verification for %q: %w", criterionID, err)
		}
		latest[role] = map[string]any{
			"status":      status,
			"verified_at": verifiedAt,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate acceptance criterion verifications for %q: %w", criterionID, err)
	}
	return latest, nil
}

func issueMarkdownLinks(conn *sql.DB, issueID string) ([]fidformat.IssueMarkdownLink, error) {
	rows, err := conn.Query(`
		SELECT l.relation, l.to_type, l.to_id, target_issue.title
		FROM links l
		LEFT JOIN issues target_issue
		  ON l.to_type = 'issue' AND target_issue.id = l.to_id
		WHERE l.from_type = 'issue' AND l.from_id = ?
		ORDER BY l.created_at, l.relation, l.to_type, l.to_id
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query issue links for %q: %w", issueID, err)
	}
	defer closeIssueMarkdownRows(rows)

	links := []fidformat.IssueMarkdownLink{}
	for rows.Next() {
		var link fidformat.IssueMarkdownLink
		var targetType string
		var targetTitle sql.NullString
		if err := rows.Scan(&link.Relation, &targetType, &link.TargetID, &targetTitle); err != nil {
			return nil, fmt.Errorf("scan issue link for %q: %w", issueID, err)
		}
		if targetType == "issue" && targetTitle.Valid {
			link.TargetTitle = targetTitle.String
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate issue links for %q: %w", issueID, err)
	}
	return links, nil
}

func closeIssueMarkdownRows(rows *sql.Rows) {
	_ = rows.Close()
}

func closeIssueMarkdownDB(conn *sql.DB) {
	_ = conn.Close()
}
