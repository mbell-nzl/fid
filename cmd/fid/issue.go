package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newIssueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newIssueACCommand())
	cmd.AddCommand(newIssueCreateCommand())
	cmd.AddCommand(newIssueLabelCommand())
	cmd.AddCommand(newIssueUpdateCommand())
	cmd.AddCommand(newIssueTransitionCommand())
	cmd.AddCommand(newIssueGetCommand())
	cmd.AddCommand(newIssueListCommand())
	return cmd
}
