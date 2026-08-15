package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newIssueACCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ac",
		Short: "Manage issue acceptance criteria.",
	}

	cmd.AddCommand(newIssueACAddCommand())
	cmd.AddCommand(newIssueACListCommand())
	cmd.AddCommand(newIssueACVerifyCommand())
	return cmd
}

func newIssueACAddCommand() *cobra.Command {
	var (
		text      string
		ordinal   int
		actorRole string
	)

	cmd := &cobra.Command{
		Use:   "add <issue-id>",
		Short: "Add an acceptance criterion to a task or subtask.",
		Args:  canonicalExactArgs("issue.ac.add", 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireCanonicalFlag(cmd, "issue.ac.add", "text"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "issue.ac.add", "ordinal"); err != nil {
				return err
			}

			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.ac.add", args[0], err)
			}

			result, err := issues.AddAcceptanceCriterion(issues.AddAcceptanceCriterionOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				Text:          text,
				Ordinal:       ordinal,
				ActorRole:     actorRole,
			})
			if err != nil {
				return renderIssueMutationError(cmd, "issue.ac.add", args[0], err)
			}

			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "issue.ac.add", map[string]any{
					"acceptance_criterion": map[string]any{"id": result.ID, "issue_id": result.IssueID, "ordinal": result.Ordinal, "text": result.Text},
					"issue_revision":       result.IssueRevision,
				}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
			}
			return writeIssueLine(
				cmd.OutOrStdout(),
				"added AC %s on %s (ordinal %d): %q\n",
				result.ID,
				result.IssueID,
				result.Ordinal,
				result.Text,
			)
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Acceptance criterion text")
	cmd.Flags().IntVar(&ordinal, "ordinal", 0, "Acceptance criterion display order")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.SetFlagErrorFunc(canonicalFlagError("issue.ac.add"))
	return cmd
}

func init() {
	rootFlagCapabilities["issue ac add"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["issue ac list"] = rootFlagCapability{Project: true, JSON: true}
}

func newIssueACVerifyCommand() *cobra.Command {
	var (
		actorRole    string
		status       string
		evidenceType string
		evidenceRef  string
		notes        string
	)

	cmd := &cobra.Command{
		Use:   "verify <issue-id> <ac-id>",
		Short: "Verify an acceptance criterion.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			if !cmd.Flags().Changed("status") {
				err := fmt.Errorf("--status is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.VerifyAcceptanceCriterion(issues.VerifyAcceptanceCriterionOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				CriterionID:   args[1],
				ActorRole:     actorRole,
				Status:        status,
				EvidenceType:  issueCreateOptionalString(cmd, "evidence-type", evidenceType),
				EvidenceRef:   issueCreateOptionalString(cmd, "evidence-ref", evidenceRef),
				Notes:         issueCreateOptionalString(cmd, "notes", notes),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeIssueLine(
				cmd.ErrOrStderr(),
				"verified %s on %s by %s: %s\n",
				result.CriterionID,
				result.IssueID,
				result.ActorRole,
				result.Status,
			); err != nil {
				return err
			}
			if result.EvidenceType != nil && result.EvidenceRef != nil {
				return writeIssueLine(cmd.ErrOrStderr(), "  evidence: %s → %s\n", *result.EvidenceType, *result.EvidenceRef)
			}
			if result.EvidenceType != nil {
				return writeIssueLine(cmd.ErrOrStderr(), "  evidence: %s\n", *result.EvidenceType)
			}
			if result.EvidenceRef != nil {
				return writeIssueLine(cmd.ErrOrStderr(), "  evidence: %s\n", *result.EvidenceRef)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Verifying role")
	cmd.Flags().StringVar(&status, "status", "", "Verification status")
	cmd.Flags().StringVar(&evidenceType, "evidence-type", "", "Evidence type")
	cmd.Flags().StringVar(&evidenceRef, "evidence-ref", "", "Evidence reference")
	cmd.Flags().StringVar(&notes, "notes", "", "Verification notes")
	return cmd
}

func newIssueACListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <issue-id>",
		Short: "List acceptance criteria for an issue.",
		Args:  canonicalExactArgs("issue.ac.list", 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueReadError(cmd, true, "issue.ac.list", err)
			}

			result, err := issues.ListAcceptanceCriteria(issues.ListAcceptanceCriteriaOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
			})
			if err != nil {
				return renderIssueReadError(cmd, true, "issue.ac.list", err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				criteria := make([]map[string]any, 0, len(result.Criteria))
				for _, criterion := range result.Criteria {
					criteria = append(criteria, map[string]any{"id": criterion.ID, "ordinal": criterion.Ordinal, "text": criterion.Text})
				}
				return cli.WriteSuccess(cmd.OutOrStdout(), "issue.ac.list", map[string]any{"issue_id": result.IssueID, "acceptance_criteria": criteria}, nil)
			}
			for _, criterion := range result.Criteria {
				if err := writeIssueLine(cmd.OutOrStdout(), "%d. %s\n", criterion.Ordinal, criterion.Text); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.SetFlagErrorFunc(canonicalFlagError("issue.ac.list"))

	return cmd
}
