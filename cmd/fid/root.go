package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
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

type commandWorkingDirectoryKey struct{}

func execute() error {
	cmd, err := newRootCommand()
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func newRootCommand() (*cobra.Command, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("resolve command working directory: %w", err)
	}
	return newRootCommandAt(workingDirectory), nil
}

func newRootCommandAt(workingDirectory string) *cobra.Command {
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
	cmd.SetContext(context.WithValue(context.Background(), commandWorkingDirectoryKey{}, workingDirectory))
	cmd.PersistentPreRunE = rejectUnavailableRootFlags(cmd)
	cmd.PersistentFlags().String("project", "", "fid project root")
	cmd.PersistentFlags().Bool("json", false, "emit fid-cli.v1 JSON")

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand(newActivityCommand())
	cmd.AddCommand(newAttentionCommand())
	cmd.AddCommand(newBackupCommand())
	cmd.AddCommand(newCapabilitiesCommand())
	cmd.AddCommand(newActivationCommand())
	cmd.AddCommand(newArtifactCommand())
	cmd.AddCommand(newEnvelopeCommand())
	cmd.AddCommand(newExportCommand())
	cmd.AddCommand(newEvidenceCommand())
	cmd.AddCommand(newFindingCommand())
	cmd.AddCommand(newGateCommand())
	cmd.AddCommand(newClaimCommand())
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newConcernCommand())
	cmd.AddCommand(newDeviationCommand())
	cmd.AddCommand(newDecisionCommand())
	cmd.AddCommand(newHumanOutcomeCommand())
	cmd.AddCommand(newIssueCommand())
	cmd.AddCommand(newJournalCommand())
	cmd.AddCommand(newLabelCommand())
	cmd.AddCommand(newLinkCommand())
	cmd.AddCommand(newMaintenanceCommand())
	cmd.AddCommand(newReadyCommand())
	cmd.AddCommand(newReceiptCommand())
	cmd.AddCommand(newReviewCommand())
	cmd.AddCommand(newResumeCommand())
	cmd.AddCommand(newRequirementCommand())
	cmd.AddCommand(newRebuildCommand())
	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newSQLCommand())
	cmd.AddCommand(newStatsCommand())
	cmd.AddCommand(newStopCommand())
	cmd.AddCommand(newWorkflowCommand())
	installCommandGroups(cmd)
	installLeafGrammar(cmd)

	return cmd
}

func commandWorkingDirectory(cmd *cobra.Command) (string, error) {
	workingDirectory, ok := cmd.Root().Context().Value(commandWorkingDirectoryKey{}).(string)
	if !ok || strings.TrimSpace(workingDirectory) == "" {
		return "", fmt.Errorf("resolve workspace root: command working directory is unavailable")
	}
	return workingDirectory, nil
}

func rejectUnavailableRootFlags(root *cobra.Command) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, _ []string) error {
		if command.HasSubCommands() {
			return nil
		}
		capability := rootFlagCapabilities[rootCommandKey(command)]
		for _, flag := range []struct {
			name    string
			enabled bool
		}{
			{"project", capability.Project},
			{"json", capability.JSON},
		} {
			if root.PersistentFlags().Changed(flag.name) && !flag.enabled {
				if diagnosticLeaf(command) {
					err := fmt.Errorf("--%s is not available for %s until its fid-cli.v1 transport is migrated", flag.name, rootCommandKey(command))
					if _, writeErr := fmt.Fprintf(command.ErrOrStderr(), "error: %v\n", err); writeErr != nil {
						return fmt.Errorf("write root-flag rejection: %w", writeErr)
					}
					return err
				}
				spelling := "--" + flag.name
				return writeGrammarFailure(command, leafOperation(command), cli.CodeInvalidArgument,
					fmt.Sprintf("%s is not available for %s", spelling, rootCommandKey(command)),
					map[string]any{"argument": spelling, "value": spelling, "allowed": allowedLocalFlags(command)})
			}
		}
		return nil
	}
}

// installCommandGroups makes unknown words reach one validator at every depth.
// Cobra otherwise treats a non-runnable group with arguments as a help request.
func installCommandGroups(command *cobra.Command) {
	if !command.HasSubCommands() {
		return
	}
	command.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		allowed := []string{}
		for _, child := range cmd.Commands() {
			if child.IsAvailableCommand() || child.Name() == "help" {
				allowed = append(allowed, child.Name())
			}
		}
		sort.Strings(allowed)
		return writeGrammarFailure(cmd, rootCommandOperation(cmd), cli.CodeOperationUnavailable,
			fmt.Sprintf("unknown command %q for %q", args[0], cmd.CommandPath()),
			map[string]any{"argument": "command", "value": args[0], "allowed": allowed})
	}
	command.RunE = func(cmd *cobra.Command, _ []string) error { return cmd.Help() }
	for _, child := range command.Commands() {
		installCommandGroups(child)
	}
}

func rootCommandOperation(cmd *cobra.Command) string {
	return strings.ReplaceAll(rootCommandKey(cmd), " ", ".")
}

func rootCommandKey(command *cobra.Command) string {
	return strings.TrimPrefix(command.CommandPath(), "fid ")
}
