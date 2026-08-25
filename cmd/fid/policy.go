package main

import (
	"encoding/json"

	internalpolicy "github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

// The compatibility helper remains for replacement-bound in-process callers.
// It is deliberately not registered on the public command tree.
var _ = newPolicyShowCommand

func newPolicyShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show the legacy permissions matrix.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetEscapeHTML(false)
			return encoder.Encode(internalpolicy.MatrixRows())
		},
	}
	return cmd
}
