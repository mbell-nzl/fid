package main

import (
	"github.com/spf13/cobra"
)

func newIssueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues.",
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
