package main

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["ready"] = rootFlagCapability{Project: true, JSON: true}
}

func newReadyCommand() *cobra.Command {
	var epicID string
	cmd := &cobra.Command{
		Use:   "ready",
		Short: "Show the next ready task.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			workspaceRoot, err := rebuildWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueReadError(cmd, jsonMode, "ready.get", err)
			}
			result, err := issues.Ready(issues.ReadyOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), EpicID: epicID})
			if err != nil {
				return renderIssueReadError(cmd, jsonMode, "ready.get", err)
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "ready.get", map[string]any{"task": result.Task, "explanation": result.Explanation}, nil)
			}
			if result.Task != nil {
				return writeIssueLine(cmd.OutOrStdout(), "%s | %s | %s\n", result.Task.ID, result.Task.Status, result.Task.Title)
			}
			return writeIssueLine(cmd.OutOrStdout(), "%s\n", readyExplanationText(*result.Explanation))
		},
	}
	cmd.Flags().StringVar(&epicID, "epic", "", "Scope to a single epic")
	return cmd
}

func readyExplanationText(explanation issues.ReadyExplanation) string {
	if explanation.Code == "sequential_current_not_actionable" && explanation.BlockingTaskID != nil {
		return fmt.Sprintf("current sequential task %s is not actionable", *explanation.BlockingTaskID)
	}
	if explanation.Scope == "epic" && explanation.EpicID != nil && strings.TrimSpace(*explanation.EpicID) != "" {
		return fmt.Sprintf("no ready tasks in epic %s", *explanation.EpicID)
	}
	return "no ready tasks in project"
}
