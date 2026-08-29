package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/spf13/cobra"
)

func canonicalExactArgs(operation string, count int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(count)(cmd, args); err != nil {
			return renderCanonicalRequestError(cmd, operation, err)
		}
		return nil
	}
}

func canonicalNoArgs(operation string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return renderCanonicalRequestError(cmd, operation, err)
		}
		return nil
	}
}

func canonicalFlagError(operation string) func(*cobra.Command, error) error {
	return func(cmd *cobra.Command, flagErr error) error {
		return renderCanonicalRequestError(cmd, operation, flagErr)
	}
}

func requireCanonicalFlag(cmd *cobra.Command, operation, name string) error {
	if cmd.Flags().Changed(name) {
		return nil
	}
	return renderCanonicalRequestError(cmd, operation, fmt.Errorf("--%s is required", name))
}

func renderCanonicalRequestError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := cli.NewFailure(cli.CodeInvalidArgument, "invalid command request", nil, map[string]any{
		"argument": "request",
		"value":    nil,
		"allowed":  []string{"accepted " + operation + " input"},
	})
	if err != nil {
		return err
	}
	if !jsonMode {
		if err := renderIssueError(cmd.ErrOrStderr(), cause); err != nil {
			return err
		}
		return cli.NewFailureError(failure, cause)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
}
