package main

import (
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newIssueLabelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage issue labels.",
	}

	cmd.AddCommand(newIssueLabelAddCommand())
	cmd.AddCommand(newIssueLabelRemoveCommand())
	return cmd
}

func newIssueLabelAddCommand() *cobra.Command {
	var actorRole string

	cmd := &cobra.Command{
		Use:   "add <issue-id> <label>",
		Short: "Add one label to an issue without replacing the rest.",
		Args:  canonicalExactArgs("issue.label.add", 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireCanonicalStringFlag(cmd, "issue.label.add", "by", actorRole); err != nil {
				return err
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.label.add", args[0], err)
			}

			result, err := issues.AddLabel(issues.AddLabelOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				Label:         args[1],
				ActorRole:     actorRole,
			})
			if err != nil {
				return renderIssueMutationError(cmd, "issue.label.add", args[0], err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "issue.label.add", map[string]any{"issue_id": result.IssueID, "label": result.Label, "issue_revision": result.IssueRevision}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
			}
			return writeIssueLine(cmd.OutOrStdout(), "added label %q to %s\n", result.Label, result.IssueID)
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.SetFlagErrorFunc(canonicalFlagError("issue.label.add"))
	return cmd
}

func newIssueLabelRemoveCommand() *cobra.Command {
	var actorRole string

	cmd := &cobra.Command{
		Use:   "remove <issue-id> <label>",
		Short: "Remove one label from an issue without changing the rest.",
		Args:  canonicalExactArgs("issue.label.remove", 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireCanonicalStringFlag(cmd, "issue.label.remove", "by", actorRole); err != nil {
				return err
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.label.remove", args[0], err)
			}

			result, err := issues.RemoveLabel(issues.RemoveLabelOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				Label:         args[1],
				ActorRole:     actorRole,
			})
			if err != nil {
				return renderIssueMutationError(cmd, "issue.label.remove", args[0], err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "issue.label.remove", map[string]any{"issue_id": result.IssueID, "label": result.Label, "issue_revision": result.IssueRevision}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
			}
			return writeIssueLine(cmd.OutOrStdout(), "removed label %q from %s\n", result.Label, result.IssueID)
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.SetFlagErrorFunc(canonicalFlagError("issue.label.remove"))
	return cmd
}

func init() {
	rootFlagCapabilities["issue label add"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["issue label remove"] = rootFlagCapability{Project: true, JSON: true}
}
