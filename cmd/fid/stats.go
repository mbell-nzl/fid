package main

import (
	"github.com/mbell-nzl/fid/internal/cli"
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

func init() {
	rootFlagCapabilities["stats tokens"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["stats cost"] = rootFlagCapability{Project: true, JSON: true}
}

func newStatsTokensCommand() *cobra.Command {
	var (
		by     string
		epicID string
	)

	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Aggregate run tokens.",
		Args:  canonicalNoArgs("stats.tokens"),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueReadError(cmd, true, "stats.tokens", err)
			}

			items, err := stats.Tokens(stats.TokensOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				By:            stats.GroupBy(by),
				EpicID:        epicID,
			})
			if err != nil {
				return renderIssueReadError(cmd, true, "stats.tokens", err)
			}
			return cli.WriteSuccess(cmd.OutOrStdout(), "stats.tokens", map[string]any{"items": tokenStatsResult(items)}, nil)
		},
	}

	cmd.Flags().StringVar(&by, "group-by", string(stats.GroupByActor), "Group by actor, model, or epic")
	cmd.Flags().StringVar(&epicID, "epic", "", "Filter by owning epic ID")
	cmd.SetFlagErrorFunc(canonicalFlagError("stats.tokens"))
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
		Args:  canonicalNoArgs("stats.cost"),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueReadError(cmd, true, "stats.cost", err)
			}

			items, err := stats.Cost(stats.CostOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				By:            stats.GroupBy(by),
				EpicID:        epicID,
			})
			if err != nil {
				return renderIssueReadError(cmd, true, "stats.cost", err)
			}
			return cli.WriteSuccess(cmd.OutOrStdout(), "stats.cost", map[string]any{"items": costStatsResult(items)}, nil)
		},
	}

	cmd.Flags().StringVar(&by, "group-by", string(stats.GroupByActor), "Group by actor, model, or epic")
	cmd.Flags().StringVar(&epicID, "epic", "", "Filter by owning epic ID")
	cmd.SetFlagErrorFunc(canonicalFlagError("stats.cost"))
	return cmd
}

func tokenStatsResult(items []stats.TokensRow) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		row := map[string]any{"actor_attribution": item.Role, "epic_id": item.EpicID, "model": item.Model, "run_count": item.RunCount, "total_tokens_in": item.TotalTokensIn, "total_tokens_out": item.TotalTokensOut}
		result = append(result, row)
	}
	return result
}

func costStatsResult(items []stats.CostRow) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		row := map[string]any{"actor_attribution": item.Role, "epic_id": item.EpicID, "model": item.Model, "run_count": item.RunCount, "total_cost_usd": item.TotalCostUSD, "total_tokens_in": item.TotalTokensIn, "total_tokens_out": item.TotalTokensOut}
		result = append(result, row)
	}
	return result
}
