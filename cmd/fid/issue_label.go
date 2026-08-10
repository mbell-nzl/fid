package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newIssueLabelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage issue labels.",
	}

	cmd.AddCommand(newIssueLabelAddCommand())
	cmd.AddCommand(newIssueLabelRemoveCommand())
	return cmd
}

func newIssueLabelAddCommand() *cobra.Command {
	var actorRole string

	cmd := &cobra.Command{
		Use:   "add <issue-id> <label>",
		Short: "Add one label to an issue without replacing the rest.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.AddLabel(issues.AddLabelOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				Label:         args[1],
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "added label %q to %s\n", result.Label, result.IssueID)
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newIssueLabelRemoveCommand() *cobra.Command {
	var actorRole string

	cmd := &cobra.Command{
		Use:   "remove <issue-id> <label>",
		Short: "Remove one label from an issue without changing the rest.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.RemoveLabel(issues.RemoveLabelOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				Label:         args[1],
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "removed label %q from %s\n", result.Label, result.IssueID)
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}
