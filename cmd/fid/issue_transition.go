package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newIssueTransitionCommand() *cobra.Command {
	var (
		toState   string
		actorRole string
	)

	cmd := &cobra.Command{
		Use:   "transition <id>",
		Short: "Transition an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Transition(issues.TransitionOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ID:            args[0],
				ToState:       toState,
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "transitioned %s: %s → %s\n", result.ID, result.FromState, result.ToState)
		},
	}

	cmd.Flags().StringVar(&toState, "to", "", "Target state")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}
