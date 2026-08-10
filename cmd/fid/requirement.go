package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/requirements"
	"github.com/spf13/cobra"
)

func newRequirementCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "requirement",
		Short: "Manage requirement pointers.",
	}

	cmd.AddCommand(newRequirementSyncCommand())
	return cmd
}

func newRequirementSyncCommand() *cobra.Command {
	var (
		sourceDoc     string
		typeFilter    string
		requirementID string
		actorRole     string
	)

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync requirement pointers from a source markdown document.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := requirements.Sync(requirements.SyncOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				SourceDoc:     sourceDoc,
				ActorRole:     actorRole,
				TypeFilter:    typeFilter,
				RequirementID: requirementID,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if _, err := fmt.Fprintf(cmd.ErrOrStderr(),
				"synced requirements from %s\n  created: %d\n  updated: %d\n  unchanged: %d\n",
				result.SourceDoc,
				result.Created,
				result.Updated,
				result.Unchanged,
			); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&sourceDoc, "source-doc", "", "Path to the source markdown document")
	cmd.Flags().StringVar(&typeFilter, "type", "", "Optional requirement type filter: PRD or TRD")
	cmd.Flags().StringVar(&requirementID, "id", "", "Optional single requirement ID to sync")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}
