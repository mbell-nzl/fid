package main

import "github.com/spf13/cobra"

func newIssueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues.",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(newIssueACCommand())
	cmd.AddCommand(newIssueCreateCommand())
	cmd.AddCommand(newIssueCompleteCommand())
	cmd.AddCommand(newIssueLabelCommand())
	cmd.AddCommand(newIssueUpdateCommand())
	cmd.AddCommand(newIssueTransitionCommand())
	cmd.AddCommand(newIssueGetCommand())
	cmd.AddCommand(newIssueListCommand())
	return cmd
}
