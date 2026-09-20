package main

import (
	"context"
	"fmt"
	"os"
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
	err = cmd.Execute()
	if err != nil && isCobraInputError(err) {
		failure, failureErr := cli.NewFailure(cli.CodeInvalidArgument, "invalid command request", nil, map[string]any{
			"argument": "request",
			"value":    nil,
			"allowed":  []string{"accepted command input"},
		})
		if failureErr != nil {
			return failureErr
		}
		return cli.NewFailureError(failure, err)
	}
	return err
}

func isCobraInputError(err error) bool {
	message := err.Error()
	return strings.HasPrefix(message, "unknown command ") ||
		strings.HasPrefix(message, "unknown flag: ") ||
		strings.HasPrefix(message, "accepts ") ||
		strings.HasPrefix(message, "requires at least ") ||
		strings.HasPrefix(message, "requires at most ")
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
