package main

import (
	"encoding/json"

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
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show the active permissions matrix.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(internalpolicy.MatrixRows()); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	return cmd
}
