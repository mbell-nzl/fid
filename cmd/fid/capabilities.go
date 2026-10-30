package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["capabilities"] = rootFlagCapability{Project: true, JSON: true}
}

func newCapabilitiesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capabilities",
		Short: "Discover binary capability availability.",
		Args:  capabilitiesNoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderCapabilitiesError(cmd, err)
			}
			items, err := workflow.DiscoverCapabilities(workspaceRoot, issueMigrationsDir())
			if err != nil {
				return renderCapabilitiesError(cmd, err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "capabilities.get", map[string]any{"items": items}, nil)
			}
			for _, item := range items {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%t\n", item.CapabilityID, item.Available); err != nil {
					return fmt.Errorf("write capabilities: %w", err)
				}
			}
			return nil
		},
	}
	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return renderCapabilitiesInvalidRequest(cmd, "flags", err.Error())
	})
	return cmd
}

func capabilitiesNoArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return renderCapabilitiesInvalidRequest(cmd, "arguments", args[0])
}

func renderCapabilitiesInvalidRequest(cmd *cobra.Command, argument, value string) error {
	cause := fmt.Errorf("capabilities accepts no arguments")
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := cli.NewFailure(cli.CodeInvalidArgument, "capabilities accepts no arguments", nil, map[string]any{
		"allowed": []string{"no arguments"}, "argument": argument, "value": value,
	})
	if err != nil {
		return err
	}
	if !jsonMode {
		return writeHumanFailure(cmd, "capabilities.get", failure, cause)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "capabilities.get", failure, cause)
}

func renderCapabilitiesError(cmd *cobra.Command, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := cli.NewFailure(cli.CodeInternalError, "capability discovery failed", nil, map[string]any{})
	if err != nil {
		return err
	}
	if !jsonMode {
		return writeHumanFailure(cmd, "capabilities.get", failure, cause)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "capabilities.get", failure, cause)
}
