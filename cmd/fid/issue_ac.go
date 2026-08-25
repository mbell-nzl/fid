package main

import (
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newIssueACCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ac",
		Short: "Manage issue acceptance criteria.",
	}

	cmd.AddCommand(newIssueACAddCommand())
	cmd.AddCommand(newIssueACListCommand())
	cmd.AddCommand(newIssueACVerifyCommand())
	return cmd
}

func newIssueACAddCommand() *cobra.Command {
	var (
		text    string
		ordinal int
		actor   string
		key     string
	)

	cmd := &cobra.Command{
		Use:   "add <issue-id>",
		Short: "Add an acceptance criterion to a task or subtask.",
		Args:  canonicalExactArgs("issue.ac.add", 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireCanonicalFlag(cmd, "issue.ac.add", "text"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "issue.ac.add", "ordinal"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "issue.ac.add", "actor"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "issue.ac.add", "idempotency-key"); err != nil {
				return err
			}

			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueMutationError(cmd, "issue.ac.add", args[0], err)
			}

			result, err := issues.AddAcceptanceCriterion(issues.AddAcceptanceCriterionOptions{
				WorkspaceRoot:  workspaceRoot,
				MigrationsDir:  issueMigrationsDir(),
				IssueID:        args[0],
				Text:           text,
				Ordinal:        ordinal,
				Actor:          actor,
				IdempotencyKey: key,
			})
			if err != nil {
				return renderIssueMutationError(cmd, "issue.ac.add", args[0], err)
			}

			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "issue.ac.add", map[string]any{
					"acceptance_criterion": map[string]any{"id": result.ID, "issue_id": result.IssueID, "ordinal": result.Ordinal, "text": result.Text},
					"issue_revision":       result.IssueRevision,
					"idempotency_key":      result.IdempotencyKey,
					"replayed":             result.Replayed,
				}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
			}
			return writeIssueLine(
				cmd.OutOrStdout(),
				"added AC %s on %s (ordinal %d): %q\n",
				result.ID,
				result.IssueID,
				result.Ordinal,
				result.Text,
			)
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Acceptance criterion text")
	cmd.Flags().IntVar(&ordinal, "ordinal", 0, "Acceptance criterion display order")
	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("issue.ac.add"))
	return cmd
}

func init() {
	rootFlagCapabilities["issue ac add"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["issue ac list"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["issue ac verify"] = rootFlagCapability{Project: true, JSON: true}
}

func newIssueACVerifyCommand() *cobra.Command {
	var actor string

	cmd := &cobra.Command{
		Use:   "verify <issue-id> <ac-id>",
		Short: "Verify an acceptance criterion.",
		Args:  canonicalExactArgs("issue.ac.verify", 2),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writeCapabilityUnavailable(cmd, "issue.ac.verify", "reviews.coverage.v1", 23)
		},
	}

	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.SetFlagErrorFunc(canonicalFlagError("issue.ac.verify"))
	return cmd
}

func newIssueACListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <issue-id>",
		Short: "List acceptance criteria for an issue.",
		Args:  canonicalExactArgs("issue.ac.list", 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueReadError(cmd, true, "issue.ac.list", err)
			}

			result, err := issues.ListAcceptanceCriteria(issues.ListAcceptanceCriteriaOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
			})
			if err != nil {
				return renderIssueReadError(cmd, true, "issue.ac.list", err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				criteria := make([]map[string]any, 0, len(result.Criteria))
				for _, criterion := range result.Criteria {
					criteria = append(criteria, map[string]any{"id": criterion.ID, "ordinal": criterion.Ordinal, "text": criterion.Text})
				}
				return cli.WriteSuccess(cmd.OutOrStdout(), "issue.ac.list", map[string]any{"issue_id": result.IssueID, "acceptance_criteria": criteria}, nil)
			}
			for _, criterion := range result.Criteria {
				if err := writeIssueLine(cmd.OutOrStdout(), "%d. %s\n", criterion.Ordinal, criterion.Text); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.SetFlagErrorFunc(canonicalFlagError("issue.ac.list"))

	return cmd
}
