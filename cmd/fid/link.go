package main

import (
	"fmt"
	"os"

	"github.com/mbell-nzl/fid/internal/links"
	"github.com/spf13/cobra"
)

func newLinkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link",
		Short: "Manage typed links.",
	}

	cmd.AddCommand(newLinkAddCommand())
	cmd.AddCommand(newLinkRemoveCommand())
	return cmd
}

func newLinkAddCommand() *cobra.Command {
	var actorRole string

	cmd := &cobra.Command{
		Use:   "add <from-id> <relation> <to-id>",
		Short: "Add a typed link and its inverse.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := links.Add(links.AddOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				FromID:        args[0],
				Relation:      args[1],
				ToID:          args[2],
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if _, err := fmt.Fprintf(
				cmd.ErrOrStderr(),
				"linked %s %s %s\n  inverse: %s %s %s\n",
				result.Declared.FromID,
				result.Declared.Relation,
				result.Declared.ToID,
				result.Inverse.FromID,
				result.Inverse.Relation,
				result.Inverse.ToID,
			); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newLinkRemoveCommand() *cobra.Command {
	var actorRole string

	cmd := &cobra.Command{
		Use:   "remove <from-id> <relation> <to-id>",
		Short: "Remove a typed link and its inverse.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := links.Remove(links.RemoveOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				FromID:        args[0],
				Relation:      args[1],
				ToID:          args[2],
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if _, err := fmt.Fprintf(
				cmd.ErrOrStderr(),
				"removed link %s %s %s\n  inverse removed: %s %s %s\n",
				result.Declared.FromID,
				result.Declared.Relation,
				result.Declared.ToID,
				result.Inverse.FromID,
				result.Inverse.Relation,
				result.Inverse.ToID,
			); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}
