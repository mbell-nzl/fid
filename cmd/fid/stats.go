package main

import (
	"encoding/json"
	"fmt"

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
		by     string
		epicID string
	)

	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Aggregate run tokens.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
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

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(items); err != nil {
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
	return cmd
}

func newStatsCostCommand() *cobra.Command {
	var (
		by     string
		epicID string
	)

	cmd := &cobra.Command{
		Use:   "cost",
		Short: "Aggregate run costs.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
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

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(items); err != nil {
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
	return cmd
}
