package main

import (
	"os"

	"github.com/spf13/cobra"
)

func execute() error {
	return newRootCommand().Execute()
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "fid",
		Short:         "Local-first, agent-first issue tracker for AI-assisted software development.",
		Long:          "fid is a local-first, agent-first issue tracker for AI-assisted software development.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand(newClaimCommand())
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newConcernCommand())
	cmd.AddCommand(newIssueCommand())
	cmd.AddCommand(newLabelCommand())
	cmd.AddCommand(newLinkCommand())
	cmd.AddCommand(newPolicyCommand())
	cmd.AddCommand(newReadyCommand())
	cmd.AddCommand(newRequirementCommand())
	cmd.AddCommand(newRebuildCommand())
	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newStatsCommand())

	return cmd
}
