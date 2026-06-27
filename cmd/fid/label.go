package main

import (
	"fmt"
	"os"

	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/labels"
	"github.com/spf13/cobra"
)

func newLabelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage labels.",
	}

	cmd.AddCommand(newLabelAddCommand())
	cmd.AddCommand(newLabelListCommand())
	return cmd
}

func newLabelAddCommand() *cobra.Command {
	var (
		actorRole   string
		colour      string
		description string
	)

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Create a label.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := labels.Add(labels.AddOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				Name:          args[0],
				Colour:        optionalStringFlag(colour),
				Description:   optionalStringFlag(description),
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if _, err := fmt.Fprintf(cmd.ErrOrStderr(), "created label %q\n", result.Label.Name); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&colour, "colour", "", "Label colour")
	cmd.Flags().StringVar(&description, "description", "", "Label description")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func optionalStringFlag(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func newLabelListCommand() *cobra.Command {
	var formatName string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			items, err := labels.List(labels.ListOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, items); err != nil {
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
