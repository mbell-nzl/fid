package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type rootFlagCapability struct {
	Project bool
	JSON    bool
}

// rootFlagCapabilities is the sole registry for leaves whose fid-cli.v1
// transport has migrated. Task 3 registers project.init after its transport is
// available; until then both global flags fail closed for every leaf.
var rootFlagCapabilities = map[string]rootFlagCapability{}

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
	cmd.PersistentPreRunE = rejectUnavailableRootFlags(cmd)
	cmd.PersistentFlags().String("project", "", "fid project root")
	cmd.PersistentFlags().Bool("json", false, "emit fid-cli.v1 JSON")

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand(newActivityCommand())
	cmd.AddCommand(newClaimCommand())
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newConcernCommand())
	cmd.AddCommand(newDecisionCommand())
	cmd.AddCommand(newDevCommand())
	cmd.AddCommand(newEscalationCommand())
	cmd.AddCommand(newIssueCommand())
	cmd.AddCommand(newJournalCommand())
	cmd.AddCommand(newLabelCommand())
	cmd.AddCommand(newLinkCommand())
	cmd.AddCommand(newMaintenanceCommand())
	cmd.AddCommand(newPolicyCommand())
	cmd.AddCommand(newReadyCommand())
	cmd.AddCommand(newRequirementCommand())
	cmd.AddCommand(newRebuildCommand())
	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newSQLCommand())
	cmd.AddCommand(newStatsCommand())
	cmd.AddCommand(newTriageCommand())
	cmd.AddCommand(newWorkflowCommand())

	return cmd
}

func rejectUnavailableRootFlags(root *cobra.Command) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, _ []string) error {
		capability := rootFlagCapabilities[rootCommandKey(command)]
		for _, flag := range []struct {
			name    string
			enabled bool
		}{
			{"project", capability.Project},
			{"json", capability.JSON},
		} {
			if root.PersistentFlags().Changed(flag.name) && !flag.enabled {
				err := fmt.Errorf("--%s is not available for %s until its fid-cli.v1 transport is migrated", flag.name, rootCommandKey(command))
				if _, writeErr := fmt.Fprintf(command.ErrOrStderr(), "error: %v\n", err); writeErr != nil {
					return fmt.Errorf("write root-flag rejection: %w", writeErr)
				}
				return err
			}
		}
		return nil
	}
}

func rootCommandKey(command *cobra.Command) string {
	return strings.TrimPrefix(command.CommandPath(), "fid ")
}
