package main

import (
	"encoding/json"
	"fmt"

	internalactivity "github.com/mbell-nzl/fid/internal/activity"
	"github.com/spf13/cobra"
)

func newActivityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activity",
		Short: "Inspect the projected activity timeline.",
	}

	cmd.AddCommand(newActivityListCommand())
	return cmd
}

func newActivityListCommand() *cobra.Command {
	var (
		on        string
		since     string
		eventName string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projected activity entries.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			items, err := internalactivity.List(internalactivity.ListOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				On:            on,
				Since:         since,
				Event:         eventName,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(items); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&on, "on", "", "Filter to one entity reference, e.g. issue/FID-a1b2c3")
	cmd.Flags().StringVar(&since, "since", "", "Include entries on or after this timestamp")
	cmd.Flags().StringVar(&eventName, "event", "", "Filter by projected activity event name")
	return cmd
}
