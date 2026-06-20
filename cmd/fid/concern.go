package main

import (
	"fmt"
	"os"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newConcernCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "concern",
		Short: "Manage concerns.",
	}

	cmd.AddCommand(newConcernRaiseCommand())
	return cmd
}

func newConcernRaiseCommand() *cobra.Command {
	var (
		issueID      string
		actorRole    string
		completionID string
		summary      string
		bodyMDPath   string
	)

	cmd := &cobra.Command{
		Use:   "raise",
		Short: "Raise a concern against a completed issue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			requiredFlags := []struct {
				name string
			}{
				{name: "on"},
				{name: "by"},
				{name: "in-completion"},
				{name: "summary"},
			}
			for _, required := range requiredFlags {
				if cmd.Flags().Changed(required.name) {
					continue
				}
				err := fmt.Errorf("--%s is required", required.name)
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.RaiseConcern(issues.RaiseConcernOptions{
				WorkspaceRoot:        workspaceRoot,
				MigrationsDir:        issueMigrationsDir(),
				IssueID:              issueID,
				ActorRole:            actorRole,
				RaisedInCompletionID: completionID,
				Summary:              summary,
				BodyMDPath:           issueCreateOptionalString(cmd, "body-md", bodyMDPath),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(
				cmd.ErrOrStderr(),
				"raised concern %s on %s: %q\n",
				result.ID,
				result.IssueID,
				result.Summary,
			)
		},
	}

	cmd.Flags().StringVar(&issueID, "on", "", "Issue ID")
	cmd.Flags().StringVar(&actorRole, "by", "", "Raising role")
	cmd.Flags().StringVar(&completionID, "in-completion", "", "Completion ID this concern was raised in")
	cmd.Flags().StringVar(&summary, "summary", "", "One-line concern summary")
	cmd.Flags().StringVar(&bodyMDPath, "body-md", "", "Path to detailed concern body")
	return cmd
}
