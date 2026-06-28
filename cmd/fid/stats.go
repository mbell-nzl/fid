package main

import (
	"fmt"
	"os"

	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/stats"
	"github.com/spf13/cobra"
)

func newStatsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Inspect telemetry aggregates.",
	}

	cmd.AddCommand(newStatsTokensCommand())
	cmd.AddCommand(newStatsCostCommand())
	return cmd
}

func newStatsTokensCommand() *cobra.Command {
	var (
		by         string
		epicID     string
		formatName string
	)

	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Aggregate run tokens.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			items, err := stats.Tokens(stats.TokensOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				By:            stats.GroupBy(by),
				EpicID:        epicID,
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

	cmd.Flags().StringVar(&by, "by", string(stats.GroupByRole), "Group by role, model, or epic")
	cmd.Flags().StringVar(&epicID, "epic", "", "Filter by owning epic ID")
	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}

func newStatsCostCommand() *cobra.Command {
	var (
		by         string
		epicID     string
		formatName string
	)

	cmd := &cobra.Command{
		Use:   "cost",
		Short: "Aggregate run costs.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			items, err := stats.Cost(stats.CostOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				By:            stats.GroupBy(by),
				EpicID:        epicID,
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

	cmd.Flags().StringVar(&by, "by", string(stats.GroupByRole), "Group by role, model, or epic")
	cmd.Flags().StringVar(&epicID, "epic", "", "Filter by owning epic ID")
	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}
