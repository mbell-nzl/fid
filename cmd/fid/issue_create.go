package main

import (
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["issue create"] = rootFlagCapability{Project: true, JSON: true}
}

func newIssueCreateCommand() *cobra.Command {
	var kind string
	var title string
	var actor string
	var idempotencyKey string
	var epicID string
	var parentID string
	var priority int
	var sequence int
	var specMDPath string
	var labels []string
	var scheduling string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an issue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireCanonicalFlag(cmd, "issue.create", "actor"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "issue.create", "idempotency-key"); err != nil {
				return err
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.create", "", err)
			}

			issueParentID, err := issueCreateParentID(kind, epicID, parentID)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.create", "", err)
			}

			result, err := issues.Create(issues.CreateOptions{
				WorkspaceRoot:  workspaceRoot,
				MigrationsDir:  issueMigrationsDir(),
				Kind:           kind,
				Title:          title,
				Actor:          actor,
				IdempotencyKey: idempotencyKey,
				ParentID:       issueParentID,
				Priority:       issueCreateOptionalInt(cmd, "priority", priority),
				Sequence:       issueCreateOptionalInt(cmd, "sequence", sequence),
				SpecMDPath:     issueCreateOptionalString(cmd, "spec-md", specMDPath),
				Labels:         labels,
				Scheduling:     scheduling,
			})
			if err != nil {
				return renderIssueMutationError(cmd, "issue.create", "", err)
			}

			return renderIssueCreateSuccess(cmd, result)
		},
	}

	cmd.Flags().StringVar(&kind, "kind", "", "Issue kind")
	cmd.Flags().StringVar(&title, "title", "", "Issue title")
	cmd.Flags().StringVar(&epicID, "epic", "", "Parent epic ID")
	cmd.Flags().StringVar(&parentID, "parent", "", "Parent task ID")
	cmd.Flags().IntVar(&priority, "priority", 0, "Issue priority")
	cmd.Flags().IntVar(&sequence, "sequence", 0, "Issue sequence")
	cmd.Flags().StringVar(&specMDPath, "spec-md", "", "Spec markdown path")
	cmd.Flags().StringSliceVar(&labels, "label", nil, "Issue label")
	cmd.Flags().StringVar(&actor, "actor", "", "Standalone actor attribution")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "Idempotency key")
	cmd.Flags().StringVar(&scheduling, "scheduling", "", "Epic scheduling")
	cmd.SetFlagErrorFunc(issueMutationFlagError("issue.create"))
	return cmd
}
