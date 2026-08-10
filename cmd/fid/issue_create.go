package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newIssueCreateCommand() *cobra.Command {
	var kind string
	var title string
	var actorRole string
	var epicID string
	var parentID string
	var priority int
	var sequence int
	var specMDPath string
	var estTokens int
	var dueDate string
	var labels []string
	var inbox bool
	var ownerRole string
	var ownerActor string
	var scheduling string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an issue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			issueParentID, err := issueCreateParentID(kind, epicID, parentID)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			result, err := issues.Create(issues.CreateOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				Kind:          kind,
				Title:         title,
				ActorRole:     actorRole,
				ParentID:      issueParentID,
				Priority:      issueCreateOptionalInt(cmd, "priority", priority),
				Sequence:      issueCreateOptionalInt(cmd, "sequence", sequence),
				SpecMDPath:    issueCreateOptionalString(cmd, "spec-md", specMDPath),
				EstTokens:     issueCreateOptionalInt(cmd, "est-tokens", estTokens),
				DueDate:       issueCreateOptionalString(cmd, "due-date", dueDate),
				Labels:        labels,
				InInbox:       inbox,
				OwnerRole:     issueCreateOptionalString(cmd, "owner-role", ownerRole),
				OwnerActor:    issueCreateOptionalString(cmd, "owner-actor", ownerActor),
				Scheduling:    scheduling,
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
	cmd.Flags().StringVar(&epicID, "epic", "", "Parent epic ID")
	cmd.Flags().StringVar(&parentID, "parent", "", "Parent task ID")
	cmd.Flags().IntVar(&priority, "priority", 0, "Issue priority")
	cmd.Flags().IntVar(&sequence, "sequence", 0, "Issue sequence")
	cmd.Flags().StringVar(&specMDPath, "spec-md", "", "Spec markdown path")
	cmd.Flags().IntVar(&estTokens, "est-tokens", 0, "Estimated tokens")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date")
	cmd.Flags().StringSliceVar(&labels, "label", nil, "Issue label")
	cmd.Flags().BoolVar(&inbox, "inbox", false, "Add issue to inbox")
	cmd.Flags().StringVar(&ownerRole, "owner-role", "", "Owner role")
	cmd.Flags().StringVar(&ownerActor, "owner-actor", "", "Owner actor")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.Flags().StringVar(&scheduling, "scheduling", "", "Epic scheduling")
	return cmd
}
