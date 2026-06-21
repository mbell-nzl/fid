package main

import (
	"fmt"
	"os"

	"github.com/mbell-nzl/fid/internal/claims"
	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/spf13/cobra"
)

func newClaimCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim",
		Short: "Inspect claims.",
	}

	cmd.AddCommand(newClaimListCommand())
	return cmd
}

func newClaimListCommand() *cobra.Command {
	var (
		staleOnly  bool
		claimedBy  string
		formatName string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List current claims.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			items, err := claims.List(claims.ListOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ClaimedBy:     claimedBy,
				StaleOnly:     staleOnly,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, items); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&staleOnly, "stale", false, "Show only stale claims")
	cmd.Flags().StringVar(&claimedBy, "by", "", "Filter by claiming role")
	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}
