package main

import (
	fidformat "github.com/mbell-nzl/fid/internal/format"
	internalpolicy "github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

func newPolicyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Inspect the active permissions policy.",
	}

	cmd.AddCommand(newPolicyShowCommand())
	return cmd
}

func newPolicyShowCommand() *cobra.Command {
	var formatName string

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show the active permissions matrix.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := fidformat.Write(cmd.OutOrStdout(), formatName, internalpolicy.MatrixRows()); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}
