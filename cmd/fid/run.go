package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/runs"
	"github.com/spf13/cobra"
)

func newRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Manage recorded agent runs.",
	}

	cmd.AddCommand(newRunRecordCommand())
	return cmd
}

func newRunRecordCommand() *cobra.Command {
	var (
		issueID            string
		runRole            string
		model              string
		tokensIn           int
		tokensOut          int
		tokenBreakdownJSON string
		toolCallCount      int
		providerRunID      string
		providerCostUSD    float64
		completionID       string
		actorRole          string
	)

	cmd := &cobra.Command{
		Use:   "record",
		Short: "Record one agent run.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			requiredFlags := []string{"issue", "role", "model", "by"}
			for _, name := range requiredFlags {
				if cmd.Flags().Changed(name) {
					continue
				}
				err := fmt.Errorf("--%s is required", name)
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := runs.Record(runs.RecordOptions{
				WorkspaceRoot:      workspaceRoot,
				MigrationsDir:      issueMigrationsDir(),
				IssueID:            issueID,
				CompletionID:       changedOptionalString(cmd, "completion", completionID),
				ActorRole:          actorRole,
				RunRole:            runRole,
				Model:              model,
				TokensIn:           changedOptionalInt(cmd, "tokens-in", tokensIn),
				TokensOut:          changedOptionalInt(cmd, "tokens-out", tokensOut),
				TokenBreakdownJSON: changedOptionalString(cmd, "token-breakdown-json", tokenBreakdownJSON),
				ToolCallCount:      changedOptionalInt(cmd, "tool-call-count", toolCallCount),
				ProviderRunID:      changedOptionalString(cmd, "provider-run-id", providerRunID),
				ProviderCostUSD:    changedOptionalFloat(cmd, "provider-cost-usd", providerCostUSD),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if result.Warning != nil {
				if err := writeIssueLine(cmd.ErrOrStderr(), "%s\n", *result.Warning); err != nil {
					return err
				}
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "recorded agent run %s on %s\n", result.ID, result.IssueID); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  role: %s, model: %s\n", result.RunRole, result.Model); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  tokens: %d in, %d out\n", optionalInt(result.TokensIn), optionalInt(result.TokensOut)); err != nil {
				return err
			}
			if result.CostUSD != nil {
				return writeIssueLine(cmd.ErrOrStderr(), "  cost: $%s\n", formatRunCost(*result.CostUSD))
			}
			return writeIssueLine(cmd.ErrOrStderr(), "  cost: unknown\n")
		},
	}

	cmd.Flags().StringVar(&issueID, "issue", "", "Issue ID")
	cmd.Flags().StringVar(&runRole, "role", "", "Agent role that ran")
	cmd.Flags().StringVar(&model, "model", "", "Model used")
	cmd.Flags().IntVar(&tokensIn, "tokens-in", 0, "Input tokens")
	cmd.Flags().IntVar(&tokensOut, "tokens-out", 0, "Output tokens")
	cmd.Flags().StringVar(&tokenBreakdownJSON, "token-breakdown-json", "", "Provider-native token categories JSON")
	cmd.Flags().IntVar(&toolCallCount, "tool-call-count", 0, "Number of tool calls made by the run")
	cmd.Flags().StringVar(&providerRunID, "provider-run-id", "", "Provider request or run identifier")
	cmd.Flags().Float64Var(&providerCostUSD, "provider-cost-usd", 0, "Provider-billed cost in USD")
	cmd.Flags().StringVar(&completionID, "completion", "", "Associated completion ID")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func changedOptionalString(cmd *cobra.Command, name string, value string) *string {
	if !cmd.Flags().Changed(name) {
		return nil
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func changedOptionalInt(cmd *cobra.Command, name string, value int) *int {
	if !cmd.Flags().Changed(name) {
		return nil
	}
	cloned := value
	return &cloned
}

func changedOptionalFloat(cmd *cobra.Command, name string, value float64) *float64 {
	if !cmd.Flags().Changed(name) {
		return nil
	}
	cloned := value
	return &cloned
}

func optionalInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func formatRunCost(value float64) string {
	text := strconv.FormatFloat(value, 'f', 6, 64)
	text = strings.TrimRight(text, "0")
	text = strings.TrimRight(text, ".")
	if text == "" {
		return "0"
	}
	return text
}
