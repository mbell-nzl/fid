package main

import (
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/spf13/cobra"
)

func writeFailure(cmd *cobra.Command, jsonMode bool, operation string, cause error, code cli.Code, message string, subject *cli.Subject, details map[string]any) error {
	failure, err := cli.NewFailure(code, message, subject, details)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	return writeHumanFailure(cmd, operation, failure, cause)
}
