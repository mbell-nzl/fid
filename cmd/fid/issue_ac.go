package main

import (
	"fmt"

	fidformat "github.com/mbell-nzl/fid/internal/format"
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
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("text") {
				err := fmt.Errorf("--text is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			if !cmd.Flags().Changed("ordinal") {
				err := fmt.Errorf("--ordinal is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
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
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(
				cmd.ErrOrStderr(),
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
	return cmd
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
	var formatName string

	cmd := &cobra.Command{
		Use:   "list <issue-id>",
		Short: "List acceptance criteria for an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.ListAcceptanceCriteria(issues.ListAcceptanceCriteriaOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, result); err != nil {
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
