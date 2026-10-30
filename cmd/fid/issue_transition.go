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
		actor            string
		idempotencyKey   string
		expectedRevision string
	)

	cmd := &cobra.Command{
		Use:   "transition <id>",
		Short: "Transition an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireCanonicalFlag(cmd, "issue.transition", "actor"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "issue.transition", "idempotency-key"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "issue.transition", "to"); err != nil {
				return err
			}
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
				Actor:            actor,
				IdempotencyKey:   idempotencyKey,
				ExpectedRevision: expected,
			})
			if err != nil {
				return renderIssueMutationError(cmd, "issue.transition", args[0], err)
			}

			return renderIssueTransitionSuccess(cmd, result)
		},
	}

	cmd.Flags().StringVar(&toState, "to", "", "Target workflow state; fid resume <issue-id> lists the Issue's next states")
	cmd.Flags().StringVar(&expectedRevision, "expected-revision", "", "Required current Issue revision; fid resume <issue-id> shows it")
	cmd.Flags().StringVar(&actor, "actor", "", "Standalone actor attribution")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(issueMutationFlagError("issue.transition"))
	return cmd
}
