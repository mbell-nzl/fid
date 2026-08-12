package main

import (
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["issue update"] = rootFlagCapability{Project: true, JSON: true}
}

func newIssueUpdateCommand() *cobra.Command {
	var (
		title            string
		priority         int
		labels           []string
		specMDPath       string
		actorRole        string
		expectedRevision string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.update", args[0], err)
			}
			expected, err := issueExpectedRevision(cmd, expectedRevision)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.update", args[0], err)
			}

			result, err := issues.Update(issues.UpdateOptions{
				WorkspaceRoot:    workspaceRoot,
				MigrationsDir:    issueMigrationsDir(),
				ID:               args[0],
				ActorRole:        actorRole,
				Title:            issueCreateOptionalString(cmd, "title", title),
				Priority:         issueCreateOptionalInt(cmd, "priority", priority),
				Labels:           issueOptionalStringSlice(cmd, "labels", labels),
				SpecMDPath:       issueCreateOptionalString(cmd, "spec-md", specMDPath),
				ExpectedRevision: expected,
			})
			if err != nil {
				return renderIssueMutationError(cmd, "issue.update", args[0], err)
			}

			return renderIssueUpdateSuccess(cmd, result)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New title")
	cmd.Flags().IntVar(&priority, "priority", 0, "New priority")
	cmd.Flags().StringSliceVar(&labels, "labels", nil, "Replacement issue labels")
	cmd.Flags().StringVar(&specMDPath, "spec-md", "", "New spec markdown path")
	cmd.Flags().StringVar(&expectedRevision, "expected-revision", "", "Required current Issue revision")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.SetFlagErrorFunc(issueMutationFlagError("issue.update"))
	return cmd
}
