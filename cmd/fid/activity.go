package main

import (
	"fmt"

	internalactivity "github.com/mbell-nzl/fid/internal/activity"
	"github.com/mbell-nzl/fid/internal/cli"
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
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueReadError(cmd, jsonMode, "activity.list", err)
			}

			items, err := internalactivity.List(internalactivity.ListOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				On:            on,
				Since:         since,
				Event:         eventName,
			})
			if err != nil {
				return renderIssueReadError(cmd, jsonMode, "activity.list", err)
			}

			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "activity.list", map[string]any{"items": items}, nil)
			}
			for _, item := range items {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", item.OccurredAt, item.Event); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&on, "on", "", "Filter to one entity reference, e.g. issue/FID-a1b2c3")
	cmd.Flags().StringVar(&since, "since", "", "Include entries on or after this timestamp")
	cmd.Flags().StringVar(&eventName, "event", "", "Filter by projected activity event name")
	return cmd
}

func init() {
	rootFlagCapabilities["activity list"] = rootFlagCapability{Project: true, JSON: true}
}
