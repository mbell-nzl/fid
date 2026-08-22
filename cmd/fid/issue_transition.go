package main

import (
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["issue transition"] = rootFlagCapability{Project: true, JSON: true}
}

func newIssueTransitionCommand() *cobra.Command {
	var (
		toState          string
		actorRole        string
		expectedRevision string
	)

	cmd := &cobra.Command{
		Use:   "transition <id>",
		Short: "Transition an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.transition", args[0], err)
			}
			expected, err := issueExpectedRevision(cmd, expectedRevision)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.transition", args[0], err)
			}

			result, err := issues.Transition(issues.TransitionOptions{
				WorkspaceRoot:    workspaceRoot,
				MigrationsDir:    issueMigrationsDir(),
				ID:               args[0],
				ToState:          toState,
				ActorRole:        actorRole,
				ExpectedRevision: expected,
			})
			if err != nil {
				return renderIssueMutationError(cmd, "issue.transition", args[0], err)
			}

			return renderIssueTransitionSuccess(cmd, result)
		},
	}

	cmd.Flags().StringVar(&toState, "to", "", "Target state")
	cmd.Flags().StringVar(&expectedRevision, "expected-revision", "", "Required current Issue revision")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.SetFlagErrorFunc(issueMutationFlagError("issue.transition"))
	return cmd
}
