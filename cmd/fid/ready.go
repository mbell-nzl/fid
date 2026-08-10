package main

import (
	"fmt"
	"strings"

	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newReadyCommand() *cobra.Command {
	var (
		epicID     string
		formatName string
	)

	cmd := &cobra.Command{
		Use:   "ready",
		Short: "Show the next ready task.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Ready(issues.ReadyOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				EpicID:        epicID,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, fidformat.ReadyView{Task: result.Task}); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if result.Task == nil && result.Explanation != nil {
				if err := writeIssueLine(cmd.ErrOrStderr(), "note: %s\n", readyExplanationNote(*result.Explanation)); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&epicID, "epic", "", "Scope to a single epic")
	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}

func readyExplanationNote(explanation issues.ReadyExplanation) string {
	switch explanation.Code {
	case "sequential_current_not_actionable":
		scope := readyExplanationScope(explanation)
		taskID := readyExplanationValue(explanation.BlockingTaskID, "current task")
		taskTitle := readyExplanationQuoted(explanation.BlockingTaskTitle)
		status := readyExplanationValue(explanation.BlockingTaskStatus, "started")
		return fmt.Sprintf("no ready tasks in %s; current sequential task %s%s is %s", scope, taskID, taskTitle, status)
	case "no_ready_tasks":
		return fmt.Sprintf("no ready tasks in %s", readyExplanationScope(explanation))
	default:
		return fmt.Sprintf("no ready tasks in %s", readyExplanationScope(explanation))
	}
}

func readyExplanationScope(explanation issues.ReadyExplanation) string {
	if explanation.Scope == "epic" && explanation.EpicID != nil && strings.TrimSpace(*explanation.EpicID) != "" {
		return "epic " + *explanation.EpicID
	}
	return "project"
}

func readyExplanationQuoted(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return ""
	}
	return fmt.Sprintf(" %q", *value)
}

func readyExplanationValue(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback
	}
	return *value
}
