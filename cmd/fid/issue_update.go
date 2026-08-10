package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newIssueUpdateCommand() *cobra.Command {
	var (
		title      string
		priority   int
		labels     []string
		ownerRole  string
		ownerActor string
		specMDPath string
		estTokens  int
		dueDate    string
		actorRole  string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Update(issues.UpdateOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ID:            args[0],
				ActorRole:     actorRole,
				Title:         issueCreateOptionalString(cmd, "title", title),
				Priority:      issueCreateOptionalInt(cmd, "priority", priority),
				Labels:        issueOptionalStringSlice(cmd, "labels", labels),
				SpecMDPath:    issueCreateOptionalString(cmd, "spec-md", specMDPath),
				EstTokens:     issueCreateOptionalInt(cmd, "est-tokens", estTokens),
				DueDate:       issueCreateOptionalString(cmd, "due-date", dueDate),
				OwnerRole:     issueCreateOptionalString(cmd, "owner-role", ownerRole),
				OwnerActor:    issueCreateOptionalString(cmd, "owner-actor", ownerActor),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "updated %s:\n%s\n", result.ID, issueUpdateSummary(result.Changes))
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New title")
	cmd.Flags().IntVar(&priority, "priority", 0, "New priority")
	cmd.Flags().StringSliceVar(&labels, "labels", nil, "Replacement issue labels")
	cmd.Flags().StringVar(&ownerRole, "owner-role", "", "New owner role")
	cmd.Flags().StringVar(&ownerActor, "owner-actor", "", "New owner actor")
	cmd.Flags().StringVar(&specMDPath, "spec-md", "", "New spec markdown path")
	cmd.Flags().IntVar(&estTokens, "est-tokens", 0, "New estimated tokens")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "New due date")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}
